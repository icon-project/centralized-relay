package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/icon-project/centralized-relay/relayer"
	"github.com/icon-project/centralized-relay/relayer/lvldb"
	"github.com/icon-project/centralized-relay/relayer/socket"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
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
			a.log.Info("Starting relayer", zap.String("version", relayer.Version))
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
			rly, err := relayer.NewRelayer(a.log, db, chains, fresh)
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

			// Block until the error channel sends a message.
			// The context being canceled will cause the relayer to stop,
			// so we don't want to separately monitor the ctx.Done channel,
			// because we would risk returning before the relayer cleans up.
			if err := <-rlyErrCh; err != nil && !errors.Is(err, context.Canceled) {
				a.log.Warn("Relayer start error", zap.Error(err))
				return err
			}
			return nil
		},
	}
	cmd = flushIntervalFlag(a.viper, cmd)
	cmd = freshFlag(a.viper, cmd)
	return cmd
}
