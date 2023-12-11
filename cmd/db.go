package cmd

import (
	"fmt"
	"strings"

	"github.com/icon-project/centralized-relay/relayer"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type dbState struct {
	app *appState
}

func dbCmd(a *appState) *cobra.Command {
	dbCMD := &cobra.Command{
		Use:     "database",
		Short:   "Manage the database",
		Aliases: []string{"db"},
		Example: strings.TrimSpace(fmt.Sprintf(`$ %s db [command]`, appName)),
	}

	var db &dbState

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
	messagesCmd.AddCommand(db.dbMessagesList())
	messagesCmd.AddCommand(db.dbMessagesRm())
	messagesCmd.AddCommand(db.dbMessagesRelay())
	dbCMD.AddCommand(messagesCmd, pruneCmd)
	return dbCMD
}

func (r *dbState) dbMessagesList() *cobra.Command {
	list := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List messages stored in the database",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Listing messages stored in the database...")
			chainID := cmd.Flag("chain").Value.String()
			src := cmd.Flag("src").Value.String()
			dst := cmd.Flag("dst").Value.String()
			limit, err := cmd.Flags().GetInt("limit")
			if err != nil {
				fmt.Println(err)
			}
			page, err := cmd.Flags().GetInt("page")
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("chainID: ", chainID, "src: ", src, "dst: ", dst, "limit: ", limit, "page: ", page)
		},
	}
	r.dbMessageFlagsListFlags(list)
	return list
}

func (r *dbState) dbMessagesRelay() *cobra.Command {
	rly := &cobra.Command{
		Use:     "relay",
		Aliases: []string{"rly"},
		Short:   "Relay message",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Relaying messages stored in the database...")
		},
	}
	r.dbMessageMsgIDFlag(rly)
	return rly
}

func (r *dbState) dbMessagesRm() *cobra.Command {
	rm := &cobra.Command{
		Use:   "rm",
		Short: "Remove messages stored in the database",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("removing messages stored in the database...")
		},
	}
	r.dbMessageMsgIDFlag(rm)
	return rm
}

func (r *dbState) dbMessageMsgIDFlag(cmd *cobra.Command) *string {
	return cmd.Flags().String("sn", "", "message sn to select")
}

func (r *dbState) dbMessageFlagsListFlags(cmd *cobra.Command) *cobra.Command {
	// limit numberof results
	cmd.Flags().IntP("limit", "l", 100, "limit number of results")
	// filter by chain
	cmd.Flags().StringP("chain", "c", "", "filter by chain")
	// filter by src
	cmd.Flags().String("src", "", "filter by src chain")
	// filter by dst
	cmd.Flags().String("dst", "", "filter by dst chain")
	// offset results
	cmd.Flags().IntP("page", "p", 0, "offset results")
	return cmd
}

func (r *dbState) dbMessageRemove(cmd *cobra.Command, args []string) *cobra.Command {
	rm := &cobra.Command{
		Use:   "rm",
		Short: "Remove a message from the database",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Removing a message from the database...")
		},
	}
	cmd.AddCommand(rm)
	return cmd
}

// GetRelayer returns the relayer instance
func (a *dbState) GetRelayer() (*relayer.Relayer, error) {
	rly, err := relayer.NewRelayer(a.log, a.db, a.config.Chains.GetAll(), false)
	if err != nil {
		a.log.Fatal("failed to create relayer", zap.Error(err))
		return nil, err
	}
	return rly, nil
}