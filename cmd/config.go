package cmd

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"os"
	"reflect"
	"strings"

	jsoniter "github.com/json-iterator/go"

	"github.com/icon-project/centralized-relay/relayer/chains/steller"
	"github.com/icon-project/centralized-relay/relayer/chains/sui"
	"github.com/icon-project/centralized-relay/relayer/chains/wasm"

	"github.com/icon-project/centralized-relay/relayer"
	"github.com/icon-project/centralized-relay/relayer/chains/evm"
	"github.com/icon-project/centralized-relay/relayer/chains/icon"
	"github.com/icon-project/centralized-relay/relayer/kms"
	"github.com/icon-project/centralized-relay/relayer/provider"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

func configCmd(a *appState) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "config",
		Aliases: []string{"cfg"},
		Short:   "Manage configuration file",
	}

	cmd.AddCommand(
		configShowCmd(a),
		configInitCmd(a),
	)
	return cmd
}

// Command for printing current configuration
func configShowCmd(a *appState) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show",
		Aliases: []string{"s", "list", "l"},
		Short:   "Prints current configuration",
		Args:    withUsage(cobra.NoArgs),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s config show --home %s
$ %s cfg list`, appName, a.homePath, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			home, err := cmd.Flags().GetString(flagHome)
			if err != nil {
				return err
			}

			cfgPath := a.configPath
			if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
				if _, err := os.Stat(home); os.IsNotExist(err) {
					return fmt.Errorf("home path does not exist: %s", home)
				}
				return fmt.Errorf("config does not exist: %s", cfgPath)
			}

			jsn, err := cmd.Flags().GetBool(flagJSON)
			if err != nil {
				return err
			}
			yml, err := cmd.Flags().GetBool(flagYAML)
			if err != nil {
				return err
			}
			switch {
			case yml && jsn:
				return fmt.Errorf("can't pass both --json and --yaml, must pick one")
			case jsn:
				out, err := jsoniter.Marshal(a.config.Wrapped())
				if err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(out))
				return nil
			default:
				out, err := yaml.Marshal(a.config.Wrapped())
				if err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(out))
				return nil
			}
		},
	}

	return yamlFlag(a.viper, jsonFlag(a.viper, cmd))
}

// Command for initializing an empty config at the --home location
func configInitCmd(a *appState) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init",
		Aliases: []string{"i"},
		Short:   "Creates a default home directory at path defined by --home",
		Args:    withUsage(cobra.NoArgs),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s config init --home %s
$ %s cfg i`, appName, a.homePath, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := os.MkdirAll(a.homePath, os.ModePerm); err != nil {
				return err
			}
			// Then create the file...
			if _, err := os.Stat(a.configPath); err == nil {
				return fmt.Errorf("config already exists: %s", a.configPath)
			}
			f, err := os.Create(a.configPath)
			if err != nil {
				return err
			}
			defer f.Close()

			// And write the default config to that location...
			if _, err = f.Write(defaultConfigYAML()); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}

// GlobalConfig describes any global relayer settings
type GlobalConfig struct {
	Timeout     string        `yaml:"timeout" json:"timeout"`
	KMSKeyID    string        `yaml:"kms-key-id" json:"kms-key-id"`
	ClusterMode ClusterConfig `yaml:"cluster-mode" json:"cluster-mode"`
}

// SetClusterMode sets the cluster mode for the global config
type ClusterConfig struct {
	Enabled    bool   `yaml:"enabled" json:"enabled"`
	Key        string `yaml:"key" json:"key"`
	privateKey *ecdsa.PrivateKey
}

func (c ClusterConfig) IsEnabled() bool {
	return c.Enabled
}

func (c ClusterConfig) SignMessage(msg []byte) ([]byte, error) {
	return c.privateKey.Sign(rand.Reader, msg, crypto.SHA256)
}

// verify message signature
func (c ClusterConfig) VerifySignature(msg, sig []byte) error {
	if c.privateKey == nil {
		return errors.New("private key is nil")
	}
	if !ecdsa.VerifyASN1(&c.privateKey.PublicKey, msg, sig) {
		return errors.New("signature verification failed")
	}
	return nil
}

// newDefaultGlobalConfig returns a global config with defaults set
func newDefaultGlobalConfig() *GlobalConfig {
	return &GlobalConfig{
		Timeout:     "10s",
		KMSKeyID:    "",
		ClusterMode: ClusterConfig{},
	}
}

type Config struct {
	Global *GlobalConfig  `yaml:"global" json:"global"`
	Chains relayer.Chains `yaml:"chains" json:"chains"`
}

func (c *Config) Save(cfgPath string) error {
	out, err := yaml.Marshal(c.Wrapped())
	if err != nil {
		return err
	}
	return os.WriteFile(cfgPath, out, 0o600)
}

// validateConfig is used to validate the GlobalConfig values
func (c *Config) validateConfig() error {
	// validating config
	return nil
}

// ConfigOutputWrapper is an intermediary type for writing the config to disk and stdout
type ConfigOutputWrapper struct {
	Global          *GlobalConfig   `yaml:"global" json:"global"`
	ProviderConfigs ProviderConfigs `yaml:"chains" json:"chains"`
}

// ConfigInputWrapper is an intermediary type for parsing the config.yaml file
type ConfigInputWrapper struct {
	Global          *GlobalConfig                         `yaml:"global"`
	ProviderConfigs map[string]*ProviderConfigYAMLWrapper `yaml:"chains"`
}

// RuntimeConfig converts the input disk config into the relayer runtime config.
func (c *ConfigInputWrapper) RuntimeConfig(ctx context.Context, a *appState) (*Config, error) {
	// build providers for each chain
	chains := make(relayer.Chains)
	kmsProvider, err := kms.NewKMSConfig(ctx, &c.Global.KMSKeyID)
	if err != nil {
		return nil, err
	}
	if c.Global.ClusterMode.Enabled && c.Global.ClusterMode.Key != "" {
		path := a.homePath + "/keystore/cluster/" + c.Global.ClusterMode.Key
		if _, err := os.Stat(path); err != nil {
			return nil, fmt.Errorf("cluster key not found: %s", path)
		}
		keyBytes, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("error reading cluster key: %w", err)
		}
		key, err := kmsProvider.Decrypt(ctx, keyBytes)
		if err != nil {
			return nil, fmt.Errorf("error decrypting cluster key: %w", err)
		}
		privKey := new(ecdsa.PrivateKey)
		privKey.D = new(big.Int).SetBytes(key)
		privKey.PublicKey = ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     privKey.X,
			Y:     privKey.Y,
		}
		c.Global.ClusterMode.privateKey = privKey
	}
	for chainName, pcfg := range c.ProviderConfigs {
		prov, err := pcfg.Value.(provider.Config).NewProvider(ctx,
			a.log.With(zap.Stringp("provider_type", &pcfg.Type)),
			a.homePath, a.debug, chainName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to build ChainProviders: %w", err)
		}
		prov.Config().(provider.ClusterConfig).SetClusterMode(c.Global.ClusterMode.Enabled)
		a.kms = kmsProvider
		if err := prov.Init(ctx, a.homePath, kmsProvider); err != nil {
			return nil, fmt.Errorf("failed to initialize provider: %w", err)
		}
		chain := relayer.NewChain(a.log, prov, a.debug)
		chains[chain.ChainProvider.NID()] = chain
	}
	a.cluster = c.Global.ClusterMode

	return &Config{
		Global: c.Global,
		Chains: chains,
	}, nil
}

type ProviderConfigs map[string]*ProviderConfigWrapper

// ProviderConfigWrapper is an intermediary type for parsing arbitrary ProviderConfigs from json files and writing to json/yaml files
type ProviderConfigWrapper struct {
	Type  string          `yaml:"type"  json:"type"`
	Value provider.Config `yaml:"value" json:"value"`
}

// ProviderConfigYAMLWrapper is an intermediary type for parsing arbitrary ProviderConfigs from yaml files
type ProviderConfigYAMLWrapper struct {
	Type  string `yaml:"type"`
	Value any    `yaml:"-"`
}

// UnmarshalJSON adds support for unmarshalling data from an arbitrary ProviderConfig
// NOTE: Add new ProviderConfig types in the map here with the key set equal to the type of ChainProvider (e.g. cosmos, substrate, etc.)
func (pcw *ProviderConfigWrapper) UnmarshalJSON(data []byte) error {
	customTypes := map[string]reflect.Type{
		"icon":   reflect.TypeOf(icon.Config{}),
		"evm":    reflect.TypeOf(evm.Config{}),
		"cosmos": reflect.TypeOf(wasm.Config{}),
		"sui":    reflect.TypeOf(sui.Config{}),
	}
	val, err := UnmarshalJSONProviderConfig(data, customTypes)
	if err != nil {
		return err
	}
	pc := val.(provider.Config)
	pcw.Value = pc
	return nil
}

// UnmarshalYAML adds support for unmarshalling data from arbitrary ProviderConfig entries found in the config file
// NOTE: Add logic for new ProviderConfig types in a switch case here
func (iw *ProviderConfigYAMLWrapper) UnmarshalYAML(n *yaml.Node) error {
	type inputWrapper ProviderConfigYAMLWrapper
	type T struct {
		*inputWrapper `yaml:",inline"`
		Wrapper       yaml.Node `yaml:"value"`
	}

	obj := &T{inputWrapper: (*inputWrapper)(iw)}
	if err := n.Decode(obj); err != nil {
		return err
	}

	switch iw.Type {
	case "icon":
		iw.Value = new(icon.Config)
	case "evm":
		iw.Value = new(evm.Config)
	case "cosmos":
		iw.Value = new(wasm.Config)
	case "stellar":
		iw.Value = new(steller.Config)
	case "sui":
		iw.Value = new(sui.Config)
	default:
		return fmt.Errorf("%s is an invalid chain type, check your config file", iw.Type)
	}

	return obj.Wrapper.Decode(iw.Value)
}

// UnmarshalJSONProviderConfig contains the custom unmarshalling logic for ProviderConfig structs
func UnmarshalJSONProviderConfig(data []byte, customTypes map[string]reflect.Type) (any, error) {
	m := map[string]any{
		"icon":   reflect.TypeOf(icon.Config{}),
		"evm":    reflect.TypeOf(evm.Config{}),
		"cosmos": reflect.TypeOf(wasm.Config{}),
		"sui":    reflect.TypeOf(sui.Config{}),
	}
	if err := jsoniter.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	typeName := m["type"].(string)
	var provCfg provider.Config
	if ty, found := customTypes[typeName]; found {
		provCfg = reflect.New(ty).Interface().(provider.Config)
	}

	valueBytes, err := jsoniter.Marshal(m["value"])
	if err != nil {
		return nil, err
	}

	return provCfg, jsoniter.Unmarshal(valueBytes, &provCfg)
}

// Note: chainId and chainName is basically the same
// Wrapped converts the Config struct into a ConfigOutputWrapper struct
func (c *Config) Wrapped() *ConfigOutputWrapper {
	providers := make(ProviderConfigs)
	for _, chain := range c.Chains {
		pcfgw := &ProviderConfigWrapper{
			Type:  chain.ChainProvider.Type(),
			Value: chain.ChainProvider.Config(),
		}
		providers[chain.ChainProvider.Name()] = pcfgw
	}
	return &ConfigOutputWrapper{Global: c.Global, ProviderConfigs: providers}
}

func defaultConfigYAML() []byte {
	return DefaultConfig().MustYAML()
}

func DefaultConfig() *Config {
	return &Config{
		Global: newDefaultGlobalConfig(),
		Chains: make(relayer.Chains),
	}
}

func (c Config) MustYAML() []byte {
	out, err := yaml.Marshal(c)
	if err != nil {
		panic(err)
	}
	return out
}

// AddChain adds an additional chain to the config
func (c *Config) AddChain(chain *relayer.Chain) error {
	nId := chain.ChainProvider.NID()
	if nId == "" {
		return fmt.Errorf("chain ID cannot be empty")
	}
	chn, err := c.Chains.Get(nId)
	if chn != nil || err == nil {
		return fmt.Errorf("chain with NID %s already exists in config", nId)
	}
	c.Chains[chain.ChainProvider.Name()] = chain
	return nil
}
