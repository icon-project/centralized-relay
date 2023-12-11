package cmd

import (
	"fmt"
	"strings"

	"github.com/icon-project/centralized-relay/relayer"
	"github.com/icon-project/centralized-relay/relayer/store"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type dbState struct{}

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
	messagesCmd.AddCommand(db.messageRemove(a))

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
			chainID := cmd.Flag("chain").Value.String()
			limit, err := cmd.Flags().GetUint("limit")
			if err != nil {
				fmt.Println(err)
			}
			page, err := cmd.Flags().GetUint("page")
			if err != nil {
				fmt.Println(err)
			}
			rly, err := d.GetRelayer(app)
			if err != nil {
				return err
			}
			pg := store.NewPagination().WithPage(page, limit)
			messages, err := rly.GetMessageStore().GetMessages(chainID, pg)
			if err != nil {
				return err
			}
			totalMessages := len(messages)
			if totalMessages == 0 {
				fmt.Println("No messages found in the database")
				return nil
			}
			fmt.Printf("%-10s %-10s %-10s %-10s %-10s %-10s %-10s %-10s\n", "Sn", "Src", "Dst", "Height", "Event", "Retry", "Data", "Time")
			// print messages respecting pagination
			for _, msg := range messages {
				fmt.Printf("%-10d %-10s %-10s %-10d %-10s %-10d %-10s %-10s\n",
					msg.Sn, msg.Src, msg.Dst, msg.MessageHeight, msg.EventType, msg.Retry, string(msg.Data), msg.GetTime())
			}
			// Print total number of messages
			fmt.Printf("Total: %d\n", totalMessages)
			// Current and total pages of messages
			fmt.Printf("Page: %d/%d\n", page, pg.CalculateTotalPages(totalMessages))
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
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Relaying messages stored in the database...")
		},
	}
	d.messageMsgIDFlag(rly)
	return rly
}

func (d *dbState) messagesRm(app *appState) *cobra.Command {
	rm := &cobra.Command{
		Use:   "rm",
		Short: "Remove messages stored in the database",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("removing messages stored in the database...")
		},
	}
	d.messageMsgIDFlag(rm)
	return rm
}

func (r *dbState) messageMsgIDFlag(cmd *cobra.Command) {
	cmd.Flags().Int("sn", 0, "message sn to select")
	if err := cmd.MarkFlagRequired("sn"); err != nil {
		panic(err)
	}
}

func (r *dbState) dbMessageFlagsListFlags(cmd *cobra.Command) {
	// limit numberof results
	cmd.Flags().UintP("limit", "l", 10, "limit number of results")
	// filter by chain
	cmd.Flags().StringP("chain", "c", "", "filter by chain")
	// offset results
	cmd.Flags().UintP("page", "p", 1, "page number")

	// make chain arg required
	if err := cmd.MarkFlagRequired("chain"); err != nil {
		panic(err)
	}
}

func (r *dbState) messageRemove(app *appState) *cobra.Command {
	rm := &cobra.Command{
		Use:     "rm",
		Aliases: []string{"r"},
		Short:   "Remove a message from the database",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Removing a message from the database...")
		},
	}
	r.messageMsgIDFlag(rm)
	return rm
}

func (d *dbState) blockInfo(app *appState) *cobra.Command {
	block := &cobra.Command{
		Use:     "view",
		Aliases: []string{"get"},
		Short:   "Show blocks stored in the database",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Show blocks stored in the database...")
			rly, err := d.GetRelayer(app)
			if err != nil {
				return err
			}
			block := rly.GetBlockStore()
			chainID := cmd.Flag("chain").Value.String()
			height, err := block.GetLastStoredBlock(chainID)
			if err != nil {
				return err
			}
			fmt.Printf("Block height: %d\n", height)
			return nil
		},
	}
	block.Flags().StringP("chain", "c", "", "ChainID to filter by")
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
