package cmd

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/gofrs/flock"
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

func (a *appState) performConfigLockingOperation(ctx context.Context, operation func() error) error {
	lockFilePath := path.Join(a.homePath, "config.lock")
	fileLock := flock.New(lockFilePath)
	_, err := fileLock.TryLock()
	if err != nil {
		return fmt.Errorf("failed to acquire config lock: %w", err)
	}
	defer func() {
		if err := fileLock.Unlock(); err != nil {
			a.log.Error("error unlocking config file lock, please manually delete",
				zap.String("filepath", lockFilePath),
			)
		}
	}()

	// load config from file and validate it. don't want to miss
	// any changes that may have been made while unlocked.
	if err := a.loadConfigFile(ctx); err != nil {
		return fmt.Errorf("failed to initialize config from file: %w", err)
	}

	// perform the operation that requires config flock.
	if err := operation(); err != nil {
		return err
	}

	// validate config after changes have been made.
	if err := a.config.validateConfig(); err != nil {
		return fmt.Errorf("error parsing chain config: %w", err)
	}

	// marshal the new config
	out, err := yaml.Marshal(a.config.Wrapped())
	if err != nil {
		return err
	}

	cfgPath := a.configPath

	// Overwrite the config file.
	if err := os.WriteFile(cfgPath, out, 0600); err != nil {
		return fmt.Errorf("failed to write config file at %s: %w", cfgPath, err)
	}

	return nil
}

// func (a *appState) configPath() string {
// 	return path.Join(a.homePath, "config.yaml")
// }
