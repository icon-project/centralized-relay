package cmd

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

type DebugState struct {
	*dbState
	app        *appState
	chain      string
	fromHeight uint64
	toHeight   uint64
}

func newDebugState(a *appState) *DebugState {
	db := newDBState()
	return &DebugState{
		app:     a,
		dbState: db,
	}
}

func debugCmd(a *appState) *cobra.Command {
	state := newDebugState(a)
	debug := &cobra.Command{
		Use:     "debug",
		Short:   "Commands for troubleshooting the relayer",
		Aliases: []string{"dbg"},
		Example: strings.TrimSpace(fmt.Sprintf(`$ %s dbg [command]`, appName)),
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return state.closeSocket()
		},
	}

	heightCmd := &cobra.Command{
		Use:   "height",
		Short: "Get latest height of the chain",
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return state.closeSocket()
		},
	}
	heightCmd.AddCommand(state.getLatestHeight(a))

	blockCmd := &cobra.Command{
		Use:   "block",
		Short: "Get latest processed block of the chain",
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return state.closeSocket()
		},
	}
	blockCmd.AddCommand(state.getLatestProcessedBlock(a))

	queryCmd := &cobra.Command{
		Use:   "query",
		Short: "Query block range",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := state.getSocket(a)
			if err != nil {
				return err
			}
			defer client.Close()
			if state.server != nil {
				defer state.server.Close()
			}
			res, err := client.QueryBlockRange(state.chain, state.fromHeight, state.toHeight)
			if err != nil {
				return err
			}
			printLabels("Chain", "Sn", "Event Type", "height", "data")
			for _, msg := range res.Msgs {
				printValues(state.chain, msg.Sn.Text(10), msg.EventType, msg.MessageHeight, hex.EncodeToString(msg.Data))
			}
			return nil
		},
	}
	queryCmd.Flags().StringVar(&state.chain, "chain", "", "Chain ID")
	queryCmd.Flags().Uint64Var(&state.fromHeight, "from_height", 0, "From Height")
	queryCmd.Flags().Uint64Var(&state.toHeight, "to_height", 0, "To height")
	queryCmd.MarkFlagsRequiredTogether("chain", "from_height", "to_height")
	debug.AddCommand(heightCmd, blockCmd, queryCmd)

	return debug
}

func (c *DebugState) getLatestHeight(app *appState) *cobra.Command {
	getLatestHeight := &cobra.Command{
		Use:     "get",
		Short:   "Get the latest chain height",
		Aliases: []string{"g"},
		Example: strings.TrimSpace(fmt.Sprintf(`$ %s dbg height get --chain [chain-id]`, appName)),
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return c.closeSocket()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := c.getSocket(app)
			if err != nil {
				return err
			}
			defer client.Close()
			if c.server != nil {
				defer c.server.Close()
			}
			res, err := client.GetLatestHeight(c.chain)
			if err != nil {
				fmt.Println(err)
				return err
			}
			printLabels("Chain", "Latest Chain Height")
			printValues(c.chain, res.Height)
			return nil
		},
	}
	getLatestHeight.Flags().StringVar(&c.chain, "chain", "", "Chain ID")
	getLatestHeight.MarkFlagRequired("chain")
	return getLatestHeight
}

func (c *DebugState) getLatestProcessedBlock(app *appState) *cobra.Command {
	getLatestHeight := &cobra.Command{
		Use:     "get",
		Short:   "Get the latest chain height",
		Aliases: []string{"g"},
		Example: strings.TrimSpace(fmt.Sprintf(`$ %s dbg block get --chain [chain-id]`, appName)),
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return c.closeSocket()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := c.getSocket(app)
			if err != nil {
				return err
			}
			defer client.Close()
			if c.server != nil {
				defer c.server.Close()
			}
			res, err := client.GetLatestProcessedBlock(c.chain)
			if err != nil {
				fmt.Println(err)
				return err
			}
			printLabels("Chain", "Last Processed Block")
			printValues(c.chain, res.Height)
			return nil
		},
	}
	getLatestHeight.Flags().StringVar(&c.chain, "chain", "", "Chain ID")
	getLatestHeight.MarkFlagRequired("chain")
	return getLatestHeight
}
