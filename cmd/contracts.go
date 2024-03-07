package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

type contractState struct {
	*dbState
	chain  string
	msgFee uint64
	resFee uint64
}

func newContractState() *contractState {
	db := newDBState()
	return &contractState{
		dbState: db,
	}
}

func contractCMD(a *appState) *cobra.Command {
	state := newContractState()
	contract := &cobra.Command{
		Use:     "contract",
		Short:   "Manage the contracts",
		Aliases: []string{"c"},
		Example: strings.TrimSpace(fmt.Sprintf(`$ %s db [command]`, appName)),
	}

	feeCmd := &cobra.Command{
		Use:   "fee",
		Short: "Fee related operations",
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return state.closeSocket()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement fee related operations
			return nil
		},
	}

	setFeeCmd := &cobra.Command{
		Use:     "set",
		Short:   "Set the fee for the chain",
		Aliases: []string{"s"},
	}

	getFeeCmd := &cobra.Command{
		Use:     "get",
		Short:   "Get the fee set for the chain",
		Aliases: []string{"g"},
	}

	feeCmd.AddCommand(setFeeCmd, getFeeCmd)

	contract.AddCommand(feeCmd)
	return contract
}
