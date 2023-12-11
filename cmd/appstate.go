package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/icon-project/centralized-relay/relayer/lvldb"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type appState struct {
	// Log is the root logger of the application.
	// Consumers are expected to store and use local copies of the logger
	// after modifying with the .With method.
	log *zap.Logger

	viper *viper.Viper

	homePath   string
	configPath string
	dbPath     string
	debug      bool
	config     *Config
	db         *lvldb.LVLDB
}

// loadConfigFile reads config file into a.Config if file is present.
func (a *appState) loadConfigFile(ctx context.Context) error {
	if _, err := os.Stat(a.configPath); err != nil {
		// don't return error if file doesn't exist
		return nil
	}

	// read the config file bytes
	file, err := os.ReadFile(a.configPath)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	// unmarshall them into the wrapper struct
	cfgWrapper := &ConfigInputWrapper{}
	err = yaml.Unmarshal(file, cfgWrapper)
	if err != nil {
		return fmt.Errorf("error unmarshalling config: %w", err)
	}

	// retrieve the runtime configuration from the disk configuration.
	newCfg, err := cfgWrapper.RuntimeConfig(ctx, a)
	if err != nil {
		return err
	}

	// validate runtime configuration
	if err := newCfg.validateConfig(); err != nil {
		return fmt.Errorf("error parsing chain config: %w", err)
	}

	// save runtime configuration in app state
	a.config = newCfg

	return nil
}
