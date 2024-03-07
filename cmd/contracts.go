package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func contractCMD(a *appState) *cobra.Command {
	db := newDBState()
	contract := &cobra.Command{
		Use:     "contract",
		Short:   "Manage the database",
		Aliases: []string{"db"},
		Example: strings.TrimSpace(fmt.Sprintf(`$ %s db [command]`, appName)),
	}

	feeCmd := &cobra.Command{
		Use:   "fee",
		Short: "Fee related operations",
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return db.closeSocket()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement fee related operations
			return nil
		},
	}

	messagesCmd := &cobra.Command{
		Use:     "set",
		Short:   "Get messages stored in the database",
		Aliases: []string{"m"},
	}
	messagesCmd.AddCommand(db.messagesList(a), db.messagesRelay(a), db.messagesRm(a), db.revertMessage(a))

	contract.AddCommand()
	return contract
}
