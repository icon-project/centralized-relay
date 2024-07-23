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
	blockCmd.AddCommand(state.getLatestProcessedBlock(a), state.setLatestProcessedBlock(a))

	queryCmd := &cobra.Command{
		Use:   "query",
		Short: "Query block range",
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return state.closeSocket()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := state.getSocket(a)
			if err != nil {
				return err
			}
			defer client.Close()
			res, err := client.QueryBlockRange(state.chain, state.fromHeight, state.toHeight)
			if err != nil {
				fmt.Println(err)
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
	queryCmd.Flags().Uint64Var(&state.fromHeight, "fromHeight", 0, "From Height")
	queryCmd.Flags().Uint64Var(&state.toHeight, "toHeight", 0, "To height")
	queryCmd.MarkFlagsRequiredTogether("chain", "fromHeight", "toHeight")
	// queryCmd.AddCommand(state.getBlockRange(a))

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

func (c *DebugState) setLatestProcessedBlock(app *appState) *cobra.Command {
	getLatestHeight := &cobra.Command{
		Use:     "set",
		Short:   "updates the last saved block height",
		Aliases: []string{"s"},
		Example: strings.TrimSpace(fmt.Sprintf(`$ %s dbg block set --chain [chain-id] --height [height]`, appName)),
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return c.closeSocket()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := c.getSocket(app)
			if err != nil {
				return err
			}
			defer client.Close()
			res, err := client.SetLatestHeight(c.chain, c.height)
			if err != nil {
				return err
			}
			printLabels("Chain", "Height")
			printValues(c.chain, res.Height)
			return nil
		},
	}
	getLatestHeight.Flags().StringVar(&c.chain, "chain", "", "Chain ID")
	getLatestHeight.Flags().Uint64Var(&c.height, "height", 0, "Block Height")
	getLatestHeight.MarkFlagsRequiredTogether("chain", "height")
	return getLatestHeight
}

func (c *DebugState) getBlockRange(app *appState) *cobra.Command {
	getLatestHeight := &cobra.Command{
		Use:     "get",
		Short:   "Get the messages from block range",
		Aliases: []string{"g"},
		Example: strings.TrimSpace(fmt.Sprintf(`$ %s debug fee get [chain-id]`, appName)),
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return c.closeSocket()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := c.getSocket(app)
			if err != nil {
				return err
			}
			defer client.Close()
			res, err := client.QueryBlockRange(c.chain, c.fromHeight, c.toHeight)
			if err != nil {
				fmt.Println(err)
				return err
			}
			printLabels("Chain", "Sn", "Event Type", "height", "data")
			for _, msg := range res.Msgs {
				printValues(c.chain, msg.Sn.Text(10), msg.EventType, msg.MessageHeight, hex.EncodeToString(msg.Data))
			}
			return nil
		},
	}
	getLatestHeight.Flags().StringVar(&c.chain, "chain", "", "Chain ID")
	getLatestHeight.Flags().Uint64Var(&c.fromHeight, "fromHeight", 0, "From Height")
	getLatestHeight.Flags().Uint64Var(&c.toHeight, "toHeight", 0, "To height")
	getLatestHeight.MarkFlagsRequiredTogether("chain", "fromHeight", "toHeight")
	return getLatestHeight
}
