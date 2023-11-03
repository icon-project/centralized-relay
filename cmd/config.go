package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/icon-project/centralized-relay/relayer"
	"github.com/icon-project/centralized-relay/relayer/chains/evm"
	"github.com/icon-project/centralized-relay/relayer/chains/icon"
	"github.com/icon-project/centralized-relay/relayer/provider"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// GlobalConfig describes any global relayer settings
type GlobalConfig struct {
	APIListenPort  string `yaml:"api-listen-addr" json:"api-listen-addr"`
	Timeout        string `yaml:"timeout" json:"timeout"`
	Memo           string `yaml:"memo" json:"memo"`
	LightCacheSize int    `yaml:"light-cache-size" json:"light-cache-size"`
}

// newDefaultGlobalConfig returns a global config with defaults set
func newDefaultGlobalConfig(memo string) *GlobalConfig {
	return &GlobalConfig{
		APIListenPort:  ":5183",
		Timeout:        "10s",
		LightCacheSize: 20,
		Memo:           memo,
	}
}

type Config struct {
	Global GlobalConfig   `yaml:"global" json:"global"`
	Chains relayer.Chains `yaml:"chains" json:"chains"`
}

// validateConfig is used to validate the GlobalConfig values
func (c *Config) validateConfig() error {
	// validating config
	return nil
}

// ConfigOutputWrapper is an intermediary type for writing the config to disk and stdout
type ConfigOutputWrapper struct {
	Global          GlobalConfig    `yaml:"global" json:"global"`
	ProviderConfigs ProviderConfigs `yaml:"chains" json:"chains"`
}

// ConfigInputWrapper is an intermediary type for parsing the config.yaml file
type ConfigInputWrapper struct {
	Global          GlobalConfig                          `yaml:"global"`
	ProviderConfigs map[string]*ProviderConfigYAMLWrapper `yaml:"chains"`
}

// RuntimeConfig converts the input disk config into the relayer runtime config.
func (c *ConfigInputWrapper) RuntimeConfig(ctx context.Context, a *appState) (*Config, error) {
	// build providers for each chain
	chains := make(relayer.Chains)
	for chainName, pcfg := range c.ProviderConfigs {
		prov, err := pcfg.Value.(provider.ProviderConfig).NewProvider(
			a.log.With(zap.String("provider_type", pcfg.Type)),
			a.homePath, a.debug, chainName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to build ChainProviders: %w", err)
		}

		if err := prov.Init(ctx); err != nil {
			return nil, fmt.Errorf("failed to initialize provider: %w", err)
		}

		chain := relayer.NewChain(a.log, prov, a.debug)
		chains[chainName] = chain
	}

	return &Config{
		Global: c.Global,
		Chains: chains,
	}, nil
}

type ProviderConfigs map[string]*ProviderConfigWrapper

// ProviderConfigWrapper is an intermediary type for parsing arbitrary ProviderConfigs from json files and writing to json/yaml files
type ProviderConfigWrapper struct {
	Type  string                  `yaml:"type"  json:"type"`
	Value provider.ProviderConfig `yaml:"value" json:"value"`
}

// ProviderConfigYAMLWrapper is an intermediary type for parsing arbitrary ProviderConfigs from yaml files
type ProviderConfigYAMLWrapper struct {
	Type  string `yaml:"type"`
	Value any    `yaml:"-"`
}

// UnmarshalJSON adds support for unmarshalling data from an arbitrary ProviderConfig
// NOTE: Add new ProviderConfig types in the map here with the key set equal to the type of ChainProvider (e.g. cosmos, substrate, etc.)
func (pcw *ProviderConfigWrapper) UnmarshalJSON(data []byte) error {
	customTypes := map[string]reflect.Type{}
	val, err := UnmarshalJSONProviderConfig(data, customTypes)
	if err != nil {
		return err
	}
	pc := val.(provider.ProviderConfig)
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
		iw.Value = new(icon.IconProviderConfig)
	case "evm":
		iw.Value = new(evm.EVMProviderConfig)
	default:
		return fmt.Errorf("%s is an invalid chain type, check your config file", iw.Type)
	}

	return obj.Wrapper.Decode(iw.Value)
}

// UnmarshalJSONProviderConfig contains the custom unmarshalling logic for ProviderConfig structs
func UnmarshalJSONProviderConfig(data []byte, customTypes map[string]reflect.Type) (any, error) {
	m := map[string]any{
		"icon": reflect.TypeOf(icon.IconProviderConfig{}),
	}
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	typeName := m["type"].(string)
	var provCfg provider.ProviderConfig
	if ty, found := customTypes[typeName]; found {
		provCfg = reflect.New(ty).Interface().(provider.ProviderConfig)
	}

	valueBytes, err := json.Marshal(m["value"])
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(valueBytes, &provCfg); err != nil {
		return nil, err
	}

	return provCfg, nil
}
