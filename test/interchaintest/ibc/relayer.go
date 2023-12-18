package ibc

import (
	"context"
	"time"
)

type Order int

const (
	Invalid Order = iota
	Ordered
	Unordered
)

type Relayer interface {
	// restore a mnemonic to be used as a relayer wallet for a chain
	RestoreKey(ctx context.Context, rep RelayerExecReporter, cfg ChainConfig, keyName, mnemonic string) error

	// GetWallet returns a Wallet for that relayer on the given chain and a boolean indicating if it was found.
	GetWallet(chainID string) (Wallet, bool)

	// After configuration is initialized, begin relaying.
	// This method is intended to create a background worker that runs the relayer.
	// You must call StopRelayer to cleanly stop the relaying.
	StartRelayer(ctx context.Context, rep RelayerExecReporter) error

	// StopRelayer stops a relayer that started work through StartRelayer.
	StopRelayer(ctx context.Context, rep RelayerExecReporter) error

	// Exec runs an arbitrary relayer command.
	// If the Relayer implementation runs in Docker,
	// whether the invoked command is run in a one-off container or execing into an already running container
	// is an implementation detail.
	//
	// "env" are environment variables in the format "MY_ENV_VAR=value"
	Exec(ctx context.Context, rep RelayerExecReporter, cmd []string, env []string) RelayerExecResult

	ExecBin(ctx context.Context, rep RelayerExecReporter, command string, params ...interface{}) RelayerExecResult

	HomeDir() string
	RestartRelayerContainer(context.Context) error
	StopRelayerContainer(context.Context, RelayerExecReporter) error
	WriteBlockHeight(context.Context, string, uint64) error
	GetKeystore(chain string, wallet Wallet) ([]byte, error)

	RestoreKeystore(ctx context.Context, keyJSON []byte, chainID string, name string) error

	CreateConfig(ctx context.Context, configYAML []byte) error
}

//var _ Relayer = (*relayer.DockerRelayer)(nil)

// RelayerMap is a mapping from test names to a relayer set for that test.
type RelayerMap map[string]map[Wallet]bool

// AddRelayer adds the given relayer to the relayer set for the given test name.
func (r RelayerMap) AddRelayer(testName string, relayer Relayer, chainID string) {
	if _, ok := r[testName]; !ok {
		r[testName] = make(map[Wallet]bool)
	}
	wallet, exists := relayer.GetWallet(chainID)
	r[testName][wallet] = exists
}

// containsRelayer returns true if the given relayer is in the relayer set for the given test name.
func (r RelayerMap) ContainsRelayer(testName string, wallet Wallet) bool {
	_, ok := r[testName]
	return ok
}

// ExecReporter is the interface of a narrow type returned by testreporter.RelayerExecReporter.
// This avoids a direct dependency on the testreporter package,
// and it avoids the relayer needing to be aware of a *testing.T.
type RelayerExecReporter interface {
	TrackRelayerExec(
		// The name of the docker container in which this relayer command executed,
		// or empty if it did not run in docker.
		containerName string,

		// The command line passed to this invocation of the relayer.
		command []string,

		// The standard output and standard error that the relayer produced during this invocation.
		stdout, stderr string,

		// The exit code of executing the command.
		// This field may not be applicable for e.g. an in-process relayer implementation.
		exitCode int,

		// When the command started and finished.
		startedAt, finishedAt time.Time,

		// Any error that occurred during execution.
		// This indicates a failure to execute,
		// e.g. the relayer binary not being found, failure communicating with Docker, etc.
		// If the process completed with a non-zero exit code,
		// those details should be indicated between stdout, stderr, and exitCode.
		err error,
	)
}

// NopRelayerExecReporter is a no-op RelayerExecReporter.
type NopRelayerExecReporter struct{}

func (NopRelayerExecReporter) TrackRelayerExec(string, []string, string, string, int, time.Time, time.Time, error) {
}

// RelyaerExecResult holds the details of a call to Relayer.Exec.
type RelayerExecResult struct {
	// This type is a redeclaration of dockerutil.ContainerExecResult.
	// While most relayer implementations are in Docker,
	// the dockerutil package is and will continue to be _internal,
	// so we need an externally importable type for third-party Relayer implementations.
	//
	// A type alias would be a potential fit here
	// (i.e. type RelayerExecResult = dockerutil.ContainerExecResult)
	// but that would be slightly misleading as not all implementations are in Docker;
	// and the type is small enough and has no methods associated,
	// so a redeclaration keeps things simple for external implementers.

	// Err is only set when there is a failure to execute.
	// A successful execution that exits non-zero will have a nil Err
	// and an appropriate ExitCode.
	Err error

	ExitCode       int
	Stdout, Stderr []byte
}
