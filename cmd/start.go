package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/icon-project/centralized-relay/relayer"
	"github.com/icon-project/centralized-relay/relayer/lvldb"
	"github.com/icon-project/centralized-relay/relayer/socket"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// startCmd represents the start command
func startCmd(a *appState) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "start path_name",
		Aliases: []string{"st"},
		Short:   "Start the listening relayer on a given path",
		Args:    withUsage(cobra.MinimumNArgs(0)),
		Example: strings.TrimSpace(fmt.Sprintf(`
			$ %s start # start all the registered chains
		`, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			a.log.Info("Starting relayer", zap.String("version", Version))
			chains := a.config.Chains.GetAll()

			flushInterval, err := cmd.Flags().GetDuration(flagFlushInterval)
			if err != nil {
				return err
			}

			fresh, err := cmd.Flags().GetBool(flagFresh)
			if err != nil {
				return err
			}

			db, err := lvldb.NewLvlDB(a.dbPath)
			if err != nil {
				return err
			}
			sigs := make(chan os.Signal, 1)
			rly, err := relayer.NewRelayer(a.log, db, chains, fresh, sigs)
			if err != nil {
				return fmt.Errorf("error creating new relayer %v", err)
			}

			rlyErrCh, err := rly.Start(cmd.Context(), flushInterval, fresh)
			if err != nil {
				return err
			}
			listener, err := socket.NewSocket(rly)
			if err != nil {
				return err
			}
			go listener.Listen()
			defer listener.Close()

			//signal handlers
			signal.Notify(sigs, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM,
				syscall.SIGINT, syscall.SIGSEGV, syscall.SIGTRAP, os.Interrupt)

			// spin file watcher
			fileChangedChan := make(chan struct{})

			go func(doneChan chan struct{}) {
				watchFile(a.configPath, fileChangedChan)
			}(fileChangedChan)

			for {
				select {
				case sig := <-sigs:
					switch sig {
					case syscall.SIGHUP:
						reloadconfig(a, rly)
					default:
						listener.Close()
						a.log.Warn("Relayer signal handled")
						os.Exit(0)
					}
				case err = <-rlyErrCh:
					if !errors.Is(err, context.Canceled) {
						a.log.Warn("Relayer start error", zap.Error(err))
						return err
					}
				case <-fileChangedChan:
					reloadconfig(a, rly)
				}
			}
		},
	}
	cmd = flushIntervalFlag(a.viper, cmd)
	cmd = freshFlag(a.viper, cmd)
	return cmd
}

func reloadconfig(a *appState, rly *relayer.Relayer) {
	file, err := os.ReadFile(a.configPath)
	if err != nil {
		a.log.Warn("error reading config file:", zap.Error(err))
	}

	cfgWrapper := &ConfigInputWrapper{}
	err = yaml.Unmarshal(file, cfgWrapper)
	if err != nil {
		a.log.Warn("error unmarshalling config: file:", zap.Error(err))
	}
	for _, cr := range rly.GetAllChainsRuntime() {
		for chainName, pcfg := range cfgWrapper.ProviderConfigs {
			if cr.Provider.NID() == chainName {
				cr.Provider.ReloadConfigs(pcfg.Value)
			}
		}
	}
	a.log.Info("Relayer config reloaded")
}

func watchFile(filePath string, cn chan struct{}) error {
	initialStat, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	for {
		stat, err := os.Stat(filePath)
		if err != nil {
			return err
		}
		if stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
			cn <- struct{}{}
			initialStat = stat
		}
		time.Sleep(5 * time.Second)
	}
}
