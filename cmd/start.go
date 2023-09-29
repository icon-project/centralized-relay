package cmd

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/icon-project/centralized-relay/relayer"
	"github.com/icon-project/centralized-relay/relayer/lvldb"
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
			chains := make(map[string]*relayer.Chain)

			chainIDs := make([]string, 0, len(chains))
			for chainID := range chains {
				chainIDs = append(chainIDs, chainID)
			}

			// get chain configurations
			chains, err := a.config.Chains.Gets(chainIDs...)
			if err != nil {
				return err
			}

			flushInterval, err := cmd.Flags().GetDuration(flagFlushInterval)
			if err != nil {
				return err
			}

			fresh, err := cmd.Flags().GetBool(flagFresh)
			if err != nil {
				return err
			}

			db, err := lvldb.NewLvlDB(filepath.Join(defaultHome, defaultDBName))
			if err != nil {
				return err
			}

			rlyErrCh, err := relayer.Start(
				cmd.Context(),
				a.log,
				chains,
				flushInterval,
				fresh,
				db,
			)
			if err != nil {
				return err
			}

			// Block until the error channel sends a message.
			// The context being canceled will cause the relayer to stop,
			// so we don't want to separately monitor the ctx.Done channel,
			// because we would risk returning before the relayer cleans up.
			if err := <-rlyErrCh; err != nil && !errors.Is(err, context.Canceled) {
				a.log.Warn(
					"Relayer start error",
					zap.Error(err),
				)
				return err
			}
			return nil
		},
	}
	cmd = flushIntervalFlag(a.viper, cmd)
	cmd = freshFlag(a.viper, cmd)
	return cmd
}
