package cmd

import (
	"fmt"
	"strings"

	"github.com/icon-project/centralized-relay/relayer"
	"github.com/icon-project/centralized-relay/relayer/store"
	"github.com/icon-project/centralized-relay/relayer/types"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type dbState struct {
	chain string
	sn    uint64
	page  uint
	limit uint
}

func dbCmd(a *appState) *cobra.Command {
	dbCMD := &cobra.Command{
		Use:     "database",
		Short:   "Manage the database",
		Aliases: []string{"db"},
		Example: strings.TrimSpace(fmt.Sprintf(`$ %s db [command]`, appName)),
	}
	db := new(dbState)

	pruneCmd := &cobra.Command{
		Use:   "prune",
		Short: "Prune the database",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Pruning the database...")
			if err := a.db.ClearStore(); err != nil {
				a.log.Error("failed to prune database", zap.Error(err))
			}
		},
	}

	messagesCmd := &cobra.Command{
		Use:     "messages",
		Short:   "Get messages stored in the database",
		Aliases: []string{"m"},
	}
	messagesCmd.AddCommand(db.messagesList(a))
	messagesCmd.AddCommand(db.messagesRm(a))
	messagesCmd.AddCommand(db.messagesRelay(a))

	blockCmd := &cobra.Command{
		Use:     "block",
		Short:   "Get block info stored in the database",
		Aliases: []string{"b"},
	}
	blockCmd.AddCommand(db.blockInfo(a))

	dbCMD.AddCommand(messagesCmd, pruneCmd, blockCmd)
	return dbCMD
}

func (d *dbState) messagesList(app *appState) *cobra.Command {
	list := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List messages stored in the database",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Listing messages stored in the database...")
			rly, err := d.GetRelayer(app)
			if err != nil {
				return err
			}
			pg := store.NewPagination().WithPage(d.page, d.limit)
			messages, err := rly.GetMessageStore().GetMessages(d.chain, pg)
			if err != nil {
				return err
			}
			totalMessages := len(messages)
			if totalMessages == 0 {
				fmt.Println("No messages found in the database")
				return nil
			}
			printLabels("Sn", "Src", "Dst", "Height", "Event", "Retry", "Data")
			// Print messages
			for _, msg := range messages {
				fmt.Printf("%-10d %-10s %-10s %-10d %-10s %-10d %-10s\n",
					msg.Sn, msg.Src, msg.Dst, msg.MessageHeight, msg.EventType, msg.Retry, string(msg.Data))
			}
			// Print total number of messages
			fmt.Printf("Total: %d\n", totalMessages)
			// Current and total pages of messages
			fmt.Printf("Page: %d/%d\n", d.page, pg.CalculateTotalPages(totalMessages))
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
		RunE: func(cmd *cobra.Command, args []string) error {
			app.log.Debug("Relaying messages stored in the database...")
			rly, err := d.GetRelayer(app)
			if err != nil {
				return err
			}
			key := types.MessageKey{Src: d.chain, Sn: d.sn}
			message, err := rly.GetMessageStore().GetMessage(key)
			if err != nil {
				return err
			}
			message.SetIsProcessing(true)
			if err = rly.GetMessageStore().StoreMessage(message); err != nil {
				return err
			}
			srcChain, err := rly.FindChainRuntime(message.Src)
			if err != nil {
				return err
			}
			dstChain, err := rly.FindChainRuntime(message.Dst)
			if err != nil {
				return err
			}
			// skipping filters because we are relaying messages manually
			rly.RouteMessage(cmd.Context(), message, dstChain, srcChain)
			return nil
		},
	}
	d.messageMsgIDFlag(rly)
	d.messageChainFlag(rly)
	return rly
}

func (d *dbState) messagesRm(app *appState) *cobra.Command {
	rm := &cobra.Command{
		Use:   "rm",
		Short: "Remove messages stored in the database",
		RunE: func(cmd *cobra.Command, args []string) error {
			app.log.Debug("removing messages stored in the database...")
			rly, err := d.GetRelayer(app)
			if err != nil {
				return err
			}
			key := types.MessageKey{Src: d.chain, Sn: d.sn}
			message, err := rly.GetMessageStore().GetMessage(key)
			if err != nil {
				return err
			}
			app.log.Debug("message", zap.Any("message", message))
			return rly.GetMessageStore().DeleteMessage(key)
		},
	}
	d.messageMsgIDFlag(rm)
	d.messageChainFlag(rm)
	return rm
}

func (d *dbState) messageMsgIDFlag(cmd *cobra.Command) {
	cmd.Flags().Uint64Var(&d.sn, "sn", 0, "message sn to select")
	if err := cmd.MarkFlagRequired("sn"); err != nil {
		panic(err)
	}
}

func (d *dbState) messageChainFlag(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&d.chain, "chain", "c", "", "message chain to select")
	if err := cmd.MarkFlagRequired("chain"); err != nil {
		panic(err)
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
		Short:   "Show blocks stored in the database",
		RunE: func(cmd *cobra.Command, args []string) error {
			app.log.Debug("Show blocks stored in the database...")
			rly, err := d.GetRelayer(app)
			if err != nil {
				return err
			}
			block := rly.GetBlockStore()
			height, err := block.GetLastStoredBlock(d.chain)
			if err != nil {
				return err
			}
			printLabels("NID", "Height")
			printValues(d.chain, height)
			return nil
		},
	}
	d.messageChainFlag(block)
	return block
}

// GetRelayer returns the relayer instance
func (d *dbState) GetRelayer(app *appState) (*relayer.Relayer, error) {
	rly, err := relayer.NewRelayer(app.log, app.db, app.config.Chains.GetAll(), false)
	if err != nil {
		app.log.Fatal("failed to create relayer", zap.Error(err))
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
	padStr := `%-10s`
	padInt := `%-10d`
	var valueCell string
	for _, val := range values {
		if _, ok := val.(string); ok {
			valueCell += padStr + " "
		} else if _, ok := val.(int); ok {
			valueCell += padInt + " "
		} else if _, ok := val.(uint); ok {
			valueCell += padInt + " "
		} else if _, ok := val.(uint64); ok {
			valueCell += padInt + " "
		}
	}
	valueCell += "\n"
	fmt.Printf(valueCell, values...)
}
