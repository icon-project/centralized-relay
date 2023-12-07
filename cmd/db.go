package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

type dbState struct {
	chainID string
	src     string
	dst     string
	limit   int
	offset  int
	page    int
	total   int
}

func dbCmd(a *appState) *cobra.Command {
	dbCMD := &cobra.Command{
		Use:     "db",
		Short:   "Manage the database",
		Aliases: []string{"d"},
		Example: strings.TrimSpace(fmt.Sprintf(`$ %s db [command]`, appName)),
	}

	pruneCmd := &cobra.Command{
		Use:   "prune",
		Short: "Prune the database",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Pruning the database...")
			if err := a.db.ClearStore(); err != nil {
				fmt.Println(err)
			}
		},
	}

	messagesCmd := &cobra.Command{
		Use:   "messages",
		Short: "Get messages stored in the database",
		Run: func(cmd *cobra.Command, args []string) {
			dbFlags(cmd)
			dbMessageRemove(cmd, args)
			fmt.Println("Getting messages stored in the database...")
			// TODO:
		},
	}

	relayCmd := &cobra.Command{
		Use:   "relay",
		Short: "Get relay stored in the database",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Getting relay stored in the database...")
		},
	}

	dbCMD.AddCommand(pruneCmd, messagesCmd, relayCmd)
	return dbCMD
}

func dbFlags(cmd *cobra.Command) {
	// limit numberof results
	cmd.Flags().IntP("limit", "l", 100, "limit number of results")
	// filter by chain
	cmd.Flags().StringP("chain", "c", "", "filter by chain")
	// filter by src
	cmd.Flags().StringP("src", "s", "", "filter by src chain")
	// filter by dst
	cmd.Flags().StringP("dst", "d", "", "filter by dst chain")
}

func dbMessageFlags(cmd *cobra.Command) {
	// flag msg id
	cmd.Flags().StringP("msg-id", "m", "", "message id to get")
}

func dbRelayFlags(cmd *cobra.Command) {
	// flag msg id
	cmd.Flags().StringP("msg-id", "m", "", "message id to relay")
}

func dbMessageRemove(cmd *cobra.Command, args []string) *cobra.Command {
	rm := &cobra.Command{
		Use:   "rm",
		Short: "Remove a message from the database",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Removing a message from the database...")
		},
	}
	return rm
}
