package cmd

import (
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/icon-project/centralized-relay/relayer"
	"github.com/icon-project/centralized-relay/relayer/lvldb"
	"github.com/icon-project/centralized-relay/relayer/socket"
	"github.com/spf13/cobra"
)

type dbState struct {
	chain      string
	height     uint64
	sn         uint64
	txHash     string
	page       uint
	limit      uint
	server     *socket.Server
	fromHeight uint64
	toHeight   uint64
}

func newDBState() *dbState {
	return new(dbState)
}

func dbCmd(a *appState) *cobra.Command {
	db := newDBState()
	dbCMD := &cobra.Command{
		Use:     "database",
		Short:   "Manage the database",
		Aliases: []string{"db"},
		Example: strings.TrimSpace(fmt.Sprintf(`$ %s db [command]`, appName)),
	}

	pruneCmd := &cobra.Command{
		Use:   "prune",
		Short: "Prune the database",
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return db.closeSocket()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := db.getSocket(a)
			if err != nil {
				return err
			}
			result, err := client.PruneDB()
			if err != nil {
				return err
			}
			printLabels("Status")
			printValues(result.Status)
			return nil
		},
	}

	messagesCmd := &cobra.Command{
		Use:     "messages",
		Short:   "Get messages stored in the database",
		Aliases: []string{"m"},
	}
	messagesCmd.AddCommand(db.messagesList(a), db.messagesRelay(a), db.messagesRm(a), db.revertMessage(a))

	blockCmd := &cobra.Command{
		Use:     "block",
		Short:   "Get block info stored in the database",
		Aliases: []string{"b"},
	}
	blockCmd.AddCommand(db.blockInfo(a))

	dbCMD.AddCommand(messagesCmd, blockCmd, pruneCmd)
	return dbCMD
}

func (d *dbState) messagesList(app *appState) *cobra.Command {
	list := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List messages stored in the database",
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return d.closeSocket()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := d.getSocket(app)
			if err != nil {
				return err
			}
			defer client.Close()
			messages, err := client.GetMessageList(d.chain, d.limit)
			if err != nil {
				return err
			}

			printLabels("Sn", "Src", "Dst", "Height", "Event", "Retry")
			// Print messages
			for _, msg := range messages.Message {
				fmt.Printf("%-10d %-10s %-10s %-10d %-10s %-10d \n",
					msg.Sn, msg.Src, msg.Dst, msg.MessageHeight, msg.EventType, msg.Retry)
			}
			return nil
		},
	}
	d.dbMessageFlagsListFlags(list)
	return list
}

func (d *dbState) messagesRelay(app *appState) *cobra.Command {
	rly := &cobra.Command{
		Use:     "relay",
		Aliases: []string{"rly"},
		Short:   "Relay message",
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return d.closeSocket()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := d.getSocket(app)
			if err != nil {
				return err
			}
			messages, err := client.RelayMessage(d.chain, d.height, d.txHash)
			if err != nil {
				return err
			}
			printLabels("Sn", "Src", "Dst", "Height", "Event")
			for _, msg := range messages {
				printValues(msg.Sn, msg.Src, msg.Dst, msg.MessageHeight, msg.EventType)
			}
			return nil
		},
	}
	d.messageChainFlag(rly, true)
	d.messageHeightFlag(rly)
	d.messageTxHashFlag(rly)
	return rly
}

func (d *dbState) messagesRm(app *appState) *cobra.Command {
	rm := &cobra.Command{
		Use:   "rm",
		Short: "Remove messages stored in the database",
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return d.closeSocket()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := d.getSocket(app)
			if err != nil {
				return err
			}
			defer client.Close()

			result, err := client.MessageRemove(d.chain, new(big.Int).SetUint64(d.sn))
			if err != nil {
				return err
			}
			printLabels("Sn", "Src", "Dst", "Height", "Event")
			printValues(result.Sn, result.Chain, result.Dst, result.Height, result.Event)
			return nil
		},
	}
	d.messageMsgIDFlag(rm, true)
	d.messageChainFlag(rm, true)
	return rm
}

func (d *dbState) messageMsgIDFlag(cmd *cobra.Command, markRequired bool) {
	cmd.Flags().Uint64Var(&d.sn, "sn", 0, "message sn to select")
	if markRequired {
		if err := cmd.MarkFlagRequired("sn"); err != nil {
			panic(err)
		}
	}
}

func (d *dbState) messageHeightFlag(cmd *cobra.Command) {
	cmd.Flags().Uint64Var(&d.height, "height", 0, "block height")
}

func (d *dbState) messageTxHashFlag(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&d.txHash, "tx_hash", "t", "", "tx hash")
}

func (d *dbState) messageChainFlag(cmd *cobra.Command, markRequired bool) {
	cmd.Flags().StringVarP(&d.chain, "chain", "c", "", "message chain to select")
	if markRequired {
		if err := cmd.MarkFlagRequired("chain"); err != nil {
			panic(err)
		}
	}
}

func (d *dbState) dbMessageFlagsListFlags(cmd *cobra.Command) {
	// limit numberof results
	cmd.Flags().UintVarP(&d.limit, "limit", "l", 10, "limit number of results")
	// filter by chain
	cmd.Flags().StringVarP(&d.chain, "chain", "c", "", "filter by chain")
	// offset results
	cmd.Flags().UintVarP(&d.page, "page", "p", 1, "page number")

	// make chain arg required
	if err := cmd.MarkFlagRequired("chain"); err != nil {
		panic(err)
	}
}

func (d *dbState) blockInfo(app *appState) *cobra.Command {
	block := &cobra.Command{
		Use:     "view",
		Aliases: []string{"get"},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return d.closeSocket()
		},
		Short: "Show blocks stored in the database",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := d.getSocket(app)
			if err != nil {
				return err
			}
			defer client.Close()
			blocks, err := client.GetBlock(d.chain)
			if err != nil {
				return err
			}
			printLabels("NID", "Height")
			for _, block := range blocks {
				printValues(block.Chain, block.CheckPointHeight)
			}
			return nil
		},
	}
	d.messageChainFlag(block, false)
	return block
}

func (d *dbState) revertMessage(app *appState) *cobra.Command {
	revert := &cobra.Command{
		Use:     "revert",
		Aliases: []string{"rv"},
		Short:   "Revert message",
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return d.closeSocket()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := d.getSocket(app)
			if err != nil {
				return err
			}
			result, err := client.RevertMessage(d.chain, d.sn)
			if err != nil {
				return err
			}
			fmt.Fprintln(os.Stdout, result)
			return nil
		},
	}
	d.messageMsgIDFlag(revert, true)
	d.messageChainFlag(revert, true)
	return revert
}

// getRelayer returns the relayer instance
func (d *dbState) getRelayer(app *appState) (*relayer.Relayer, error) {
	db, err := lvldb.NewLvlDB(app.dbPath)
	if err != nil {
		return nil, err
	}
	rly, err := relayer.NewRelayer(app.log, db, app.config.Chains.GetAll(), false)
	if err != nil {
		fmt.Printf("failed to create relayer: %s\n", err)
		return nil, err
	}
	return rly, nil
}

func printLabels(labels ...any) {
	padStr := `%-10s`
	var labelCell string
	var border []any
	for range labels {
		labelCell += padStr + " "
		border = append(border, strings.Repeat("-", 10))
	}
	labelCell += "\n"
	fmt.Printf(labelCell, labels...)
	fmt.Printf(labelCell, border...)
}

func printValues(values ...any) {
	padStr := `%-10v`
	var valueCell string
	for range values {
		valueCell += padStr + " "
	}
	valueCell += "\n"
	fmt.Printf(valueCell, values...)
}

func (d *dbState) getSocket(app *appState) (*socket.Client, error) {
	client, err := socket.NewClient()
	if err != nil {
		if errors.Is(err, socket.ErrSocketClosed) {
			rly, err := d.getRelayer(app)
			if err != nil {
				return nil, err
			}
			server, err := socket.NewSocket(rly)
			if err != nil {
				return nil, err
			}
			d.server = server
			go server.Listen()
		}
		return socket.NewClient()
	}
	return client, nil
}

// PostRunE is a function that is called after the command is run
func (d *dbState) closeSocket() error {
	if d.server != nil {
		return d.server.Close()
	}
	return nil
}
