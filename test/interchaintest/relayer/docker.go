package relayer

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/icon-project/centralized-relay/test/interchaintest/_internal/dockerutil"
	"github.com/icon-project/centralized-relay/test/interchaintest/ibc"
	"github.com/icon-project/centralized-relay/test/interchaintest/testutil"
	iccrypto "github.com/icon-project/icon-bridge/common/crypto"
	"io"
	"time"

	"github.com/docker/docker/api/types"
	volumetypes "github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	_wallet "github.com/icon-project/icon-bridge/common/wallet"
	"go.uber.org/zap"
)

const (
	defaultRlyHomeDirectory = "/home/relayer"
)

// DockerRelayer provides a common base for relayer implementations
// that run on Docker.
type DockerRelayer struct {
	log *zap.Logger

	// c defines all the commands to run inside the container.
	c RelayerCommander

	networkID  string
	client     *client.Client
	volumeName string

	testName string

	customImage *ibc.DockerImage
	pullImage   bool

	// The ID of the container created by StartRelayer.
	containerLifecycle *dockerutil.ContainerLifecycle

	// wallets contains a mapping of chainID to relayer wallet
	wallets map[string]ibc.Wallet

	homeDir string
}

var _ ibc.Relayer = (*DockerRelayer)(nil)

// NewDockerRelayer returns a new DockerRelayer.
func NewDockerRelayer(ctx context.Context, log *zap.Logger, testName string, cli *client.Client, networkID string, c RelayerCommander, options ...RelayerOption) (*DockerRelayer, error) {
	r := DockerRelayer{
		log: log,

		c: c,

		networkID: networkID,
		client:    cli,

		// pull true by default, can be overridden with options
		pullImage: true,

		testName: testName,

		wallets: map[string]ibc.Wallet{},
		homeDir: defaultRlyHomeDirectory,
	}

	for _, opt := range options {
		switch o := opt.(type) {
		case RelayerOptionDockerImage:
			r.customImage = &o.DockerImage
		case RelayerOptionImagePull:
			r.pullImage = o.Pull
		case RelayerOptionHomeDir:
			r.homeDir = o.HomeDir
		}
	}

	containerImage := r.containerImage()
	if err := r.pullContainerImageIfNecessary(containerImage); err != nil {
		return nil, fmt.Errorf("pulling container image %s: %w", containerImage.Ref(), err)
	}

	v, err := cli.VolumeCreate(ctx, volumetypes.VolumeCreateBody{
		// Have to leave Driver unspecified for Docker Desktop compatibility.

		Labels: map[string]string{dockerutil.CleanupLabel: testName},
	})
	if err != nil {
		return nil, fmt.Errorf("creating volume: %w", err)
	}
	r.volumeName = v.Name

	// The volume is created owned by root,
	// but we configure the relayer to run as a non-root user,
	// so set the node home (where the volume is mounted) to be owned
	// by the relayer user.
	if err := dockerutil.SetVolumeOwner(ctx, dockerutil.VolumeOwnerOptions{
		Log: r.log,

		Client: r.client,

		VolumeName: r.volumeName,
		ImageRef:   containerImage.Ref(),
		TestName:   testName,
		UidGid:     containerImage.UidGid,
	}); err != nil {
		return nil, fmt.Errorf("set volume owner: %w", err)
	}

	//if init := r.c.Init(r.HomeDir()); len(init) > 0 {
	//	// Initialization should complete immediately,
	//	// but add a 1-minute timeout in case Docker hangs on a developer workstation.
	//	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	//	defer cancel()
	//
	//	// Using a nop reporter here because it keeps the API simpler,
	//	// and the init command is typically not of high interest.
	//	res := r.Exec(ctx, ibc.NopRelayerExecReporter{}, init, nil)
	//	if res.Err != nil {
	//		return nil, res.Err
	//	}
	//}

	return &r, nil
}

// WriteFileToHomeDir writes the given contents to a file at the relative path specified. The file is relative
// to the home directory in the relayer container.
func (r *DockerRelayer) WriteFileToHomeDir(ctx context.Context, relativePath string, contents []byte) error {
	fw := dockerutil.NewFileWriter(r.log, r.client, r.testName)
	if err := fw.WriteFile(ctx, r.volumeName, relativePath, contents); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}

// ReadFileFromHomeDir reads a file at the relative path specified and returns the contents. The file is
// relative to the home directory in the relayer container.
func (r *DockerRelayer) ReadFileFromHomeDir(ctx context.Context, relativePath string) ([]byte, error) {
	fr := dockerutil.NewFileRetriever(r.log, r.client, r.testName)
	bytes, err := fr.SingleFileContent(ctx, r.volumeName, relativePath)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve %s: %w", relativePath, err)
	}
	return bytes, nil
}

// Modify a toml config file in relayer home directory
func (r *DockerRelayer) ModifyTomlConfigFile(ctx context.Context, relativePath string, modification testutil.Toml) error {
	return testutil.ModifyTomlConfigFile(ctx, r.log, r.client, r.testName, r.volumeName, relativePath, modification)
}

// AddWallet adds a stores a wallet for the given chain ID.
func (r *DockerRelayer) AddWallet(chainID string, wallet ibc.Wallet) {
	r.wallets[chainID] = wallet
}

func (r *DockerRelayer) AddChainConfiguration(ctx context.Context, rep ibc.RelayerExecReporter, chainConfig ibc.ChainConfig, keyName, rpcAddr, grpcAddr string) error {
	panic("not implemented")
}

func (r *DockerRelayer) GetWallet(chainID string) (ibc.Wallet, bool) {
	wallet, ok := r.wallets[chainID]
	return wallet, ok
}

func (r *DockerRelayer) Flush(ctx context.Context, rep ibc.RelayerExecReporter, pathName, channelID string) error {
	cmd := r.c.Flush(pathName, channelID, r.HomeDir())
	res := r.Exec(ctx, rep, cmd, nil)
	return res.Err
}

func (r *DockerRelayer) ExecBin(ctx context.Context, rep ibc.RelayerExecReporter, command string, params ...interface{}) ibc.RelayerExecResult {

	cmd := r.c.RelayerCommand(command, params)
	return r.Exec(ctx, rep, cmd, nil)
}

func (r *DockerRelayer) Exec(ctx context.Context, rep ibc.RelayerExecReporter, cmd []string, env []string) ibc.RelayerExecResult {
	job := dockerutil.NewImage(r.log, r.client, r.networkID, r.testName, r.containerImage().Repository, r.containerImage().Version)
	opts := dockerutil.ContainerOptions{
		Env:   env,
		Binds: r.Bind(),
	}
	startedAt := time.Now()
	res := job.Run(ctx, cmd, opts)

	defer func() {
		rep.TrackRelayerExec(
			r.Name(),
			cmd,
			string(res.Stdout), string(res.Stderr),
			res.ExitCode,
			startedAt, time.Now(),
			res.Err,
		)
	}()

	result := ibc.RelayerExecResult{
		Err:      res.Err,
		ExitCode: res.ExitCode,
		Stdout:   res.Stdout,
		Stderr:   res.Stderr,
	}

	fmt.Println(res.Err)
	fmt.Println(res.ExitCode)
	fmt.Println(string(res.Stdout))
	fmt.Println(string(res.Stderr))

	return result
}

func (r *DockerRelayer) RestoreKey(ctx context.Context, rep ibc.RelayerExecReporter, cfg ibc.ChainConfig, keyName, mnemonic string) error {
	chainID := cfg.ChainID
	coinType := cfg.CoinType
	cmd := r.c.RestoreKey(chainID, keyName, coinType, mnemonic, r.HomeDir())

	// Restoring a key should be near-instantaneous, so add a 1-minute timeout
	// to detect if Docker has hung.
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	res := r.Exec(ctx, rep, cmd, nil)
	if res.Err != nil {
		return res.Err
	}
	addrBytes := r.c.ParseRestoreKeyOutput(string(res.Stdout), string(res.Stderr))

	r.wallets[chainID] = r.c.CreateWallet("", addrBytes, mnemonic)

	return nil
}

func (r *DockerRelayer) StartRelayer(ctx context.Context, rep ibc.RelayerExecReporter) error {
	if r.containerLifecycle != nil {
		return fmt.Errorf("tried to start relayer again without stopping first")
	}

	containerImage := r.containerImage()
	//joinedPaths := strings.Join(pathNames, ".")
	containerName := fmt.Sprintf("%s", r.c.Name())

	cmd := r.c.StartRelayer(r.HomeDir())

	r.containerLifecycle = dockerutil.NewContainerLifecycle(r.log, r.client, containerName)

	if err := r.containerLifecycle.CreateContainer(
		ctx, r.testName, r.networkID, containerImage, nil,
		r.Bind(), r.HostName(""), cmd,
	); err != nil {
		return err
	}

	return r.containerLifecycle.StartContainer(ctx)
}

func (r *DockerRelayer) stopRelayer(ctx context.Context, rep ibc.RelayerExecReporter) error {
	if r.containerLifecycle == nil {
		return fmt.Errorf("tried to stop relayer again without starting first")
	}
	if err := r.containerLifecycle.StopContainer(ctx); err != nil {
		return err
	}

	stdoutBuf := new(bytes.Buffer)
	stderrBuf := new(bytes.Buffer)
	containerID := r.containerLifecycle.ContainerID()
	rc, err := r.client.ContainerLogs(ctx, containerID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       "50",
	})
	if err != nil {
		return fmt.Errorf("StopRelayer: retrieving ContainerLogs: %w", err)
	}
	defer func() { _ = rc.Close() }()

	// Logs are multiplexed into one stream; see docs for ContainerLogs.
	_, err = stdcopy.StdCopy(stdoutBuf, stderrBuf, rc)
	if err != nil {
		return fmt.Errorf("StopRelayer: demuxing logs: %w", err)
	}
	_ = rc.Close()

	stdout := stdoutBuf.String()
	stderr := stderrBuf.String()

	c, err := r.client.ContainerInspect(ctx, containerID)
	if err != nil {
		return fmt.Errorf("StopRelayer: inspecting container: %w", err)
	}

	startedAt, err := time.Parse(c.State.StartedAt, time.RFC3339)
	if err != nil {
		startedAt, err = time.Parse(c.State.StartedAt, time.RFC3339Nano)
	}
	if startedAt.IsZero() {
		r.log.Info("Failed to parse container StartedAt", zap.Error(err))
		startedAt = time.Unix(0, 0)
	}

	var finishedAt time.Time

	finishedAt, err = time.Parse(c.State.FinishedAt, time.RFC3339)
	if err != nil {
		err = fmt.Errorf("failed to parse container FinishedAt: %w", err)
		finishedAt = time.Now().UTC()
	}

	rep.TrackRelayerExec(
		c.Name,
		c.Args,
		stdout, stderr,
		c.State.ExitCode,
		startedAt,
		finishedAt,
		err,
	)
	return nil
}

// CleanUp cleans up the relayer's home directory and any containers that were
func (r *DockerRelayer) cleanUp(ctx context.Context) error {
	if r.containerLifecycle != nil {
		return fmt.Errorf("tried to clean up relayer without stopping first")
	}
	if err := r.containerLifecycle.RemoveContainer(ctx); err != nil {
		return err
	}
	r.containerLifecycle = nil
	return nil
}

func (r *DockerRelayer) StopRelayer(ctx context.Context, rep ibc.RelayerExecReporter) error {
	return r.stopRelayer(ctx, rep)
}

// RestartRelayer restarts the relayer with the same paths as before.
func (r *DockerRelayer) StopRelayerContainer(ctx context.Context, rep ibc.RelayerExecReporter) error {
	if r.containerLifecycle == nil {
		return fmt.Errorf("tried to restart relayer without starting first")
	}
	return r.stopRelayer(ctx, rep)
}

// WriteBlockHeight writes the block height to the relayer's home directory.
func (r *DockerRelayer) WriteBlockHeight(ctx context.Context, chainID string, height uint64) error {
	return r.WriteFileToHomeDir(ctx, fmt.Sprintf(".relayer/%s/latest_height", chainID), []byte(fmt.Sprintf("%d", height)))
}

func (r *DockerRelayer) GetKeystore(chain string, wallet ibc.Wallet) ([]byte, error) {
	switch chain {
	case "evm":
		keyName := wallet.KeyName()
		privateKey, _ := crypto.HexToECDSA(wallet.Mnemonic())

		ks := keystore.NewKeyStore("/tmp", keystore.StandardScryptN, keystore.StandardScryptP)

		// Create a new account in the keystore
		account, _ := ks.ImportECDSA(privateKey, wallet.KeyName())

		keyJSON, err := ks.Export(account, keyName, keyName)
		if err != nil {
			return nil, fmt.Errorf("failed to parse keystore: %w", err)
		}
		return keyJSON, nil

	case "icon":
		keyName := wallet.KeyName()
		pk, _ := hex.DecodeString(wallet.Mnemonic())
		privateKey, _ := iccrypto.ParsePrivateKey(pk)

		ks, err := _wallet.EncryptKeyAsKeyStore(privateKey, []byte(keyName))
		if err != nil {
			return nil, fmt.Errorf("failed to parse keystore: %w", err)
		}

		return ks, nil

	default:
		return nil, fmt.Errorf("unsupported chain: %s", chain)
	}

}

func (r *DockerRelayer) CreateConfig(ctx context.Context, configYAML []byte) error {
	path := ".centralized-relay/config/config.yaml"
	fw := dockerutil.NewFileWriter(r.log, r.client, r.testName)
	if err := fw.WriteFile(ctx, r.volumeName, path, configYAML); err != nil {
		return fmt.Errorf("failed to restore config: %w", err)
	}
	return nil
}
func (r *DockerRelayer) RestoreKeystore(ctx context.Context, keyJSON []byte, chainID string, name string) error {

	ksPath := fmt.Sprintf(".centralized-relay/keys/%s/%s", chainID, name)
	fw := dockerutil.NewFileWriter(r.log, r.client, r.testName)
	if err := fw.WriteFile(ctx, r.volumeName, ksPath, keyJSON); err != nil {
		return fmt.Errorf("failed to restore keystore: %w", err)
	}
	return nil
}

// RestartRelayer restarts the relayer with the same paths as before.
func (r *DockerRelayer) RestartRelayerContainer(ctx context.Context) error {
	if r.containerLifecycle == nil {
		return fmt.Errorf("tried to restart relayer without starting first")
	}
	return r.containerLifecycle.StartContainer(ctx)
}

func (r *DockerRelayer) containerImage() ibc.DockerImage {
	if r.customImage != nil {
		return *r.customImage
	}
	return ibc.DockerImage{
		Repository: r.c.DefaultContainerImage(),
		Version:    r.c.DefaultContainerVersion(),
		UidGid:     r.c.DockerUser(),
	}
}

func (r *DockerRelayer) pullContainerImageIfNecessary(containerImage ibc.DockerImage) error {
	if !r.pullImage {
		return nil
	}

	rc, err := r.client.ImagePull(context.TODO(), containerImage.Ref(), types.ImagePullOptions{})
	if err != nil {
		return err
	}

	_, _ = io.Copy(io.Discard, rc)
	_ = rc.Close()
	return nil
}

func (r *DockerRelayer) Name() string {
	return r.c.Name() + "-" + dockerutil.SanitizeContainerName(r.testName)
}

// Bind returns the home folder bind point for running the node.
func (r *DockerRelayer) Bind() []string {
	return []string{r.volumeName + ":" + r.HomeDir()}
}

// HomeDir returns the home directory of the relayer on the underlying Docker container's filesystem.
func (r *DockerRelayer) HomeDir() string {
	return r.homeDir
}

func (r *DockerRelayer) HostName(pathName string) string {
	return dockerutil.CondenseHostName(fmt.Sprintf("%s-%s", r.c.Name(), pathName))
}

func (r *DockerRelayer) UseDockerNetwork() bool {
	return true
}

func (r *DockerRelayer) SetClientContractHash(ctx context.Context, rep ibc.RelayerExecReporter, cfg ibc.ChainConfig, hash string) error {
	panic("[rly/SetClientContractHash] Implement me")
}

type RelayerCommander interface {
	// Name is the name of the relayer, e.g. "rly" or "hermes".
	Name() string

	DefaultContainerImage() string
	DefaultContainerVersion() string

	// The Docker user to use in created container.
	// For interchaintest, must be of the format: uid:gid.
	DockerUser() string

	// ParseAddKeyOutput processes the output of AddKey
	// to produce the wallet that was created.
	ParseAddKeyOutput(stdout, stderr string) (ibc.Wallet, error)

	// ParseRestoreKeyOutput extracts the address from the output of RestoreKey.
	ParseRestoreKeyOutput(stdout, stderr string) string

	// Init is the command to run on the first call to AddChainConfiguration.
	// If the returned command is nil or empty, nothing will be executed.
	Init(homeDir string) []string

	// The remaining methods produce the command to run inside the container.
	Flush(pathName, channelID, homeDir string) []string
	RestoreKey(chainID, keyName, coinType, mnemonic, homeDir string) []string
	StartRelayer(homeDir string, pathNames ...string) []string
	CreateWallet(keyName, address, mnemonic string) ibc.Wallet

	RelayerCommand(command string, params ...interface{}) []string
}
