package cmd

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/spf13/cobra"
)

type contractState struct {
	*dbState
	app    *appState
	chain  string
	msgFee int64
	resFee int64
}

func newContractState(a *appState) *contractState {
	db := newDBState()
	return &contractState{
		dbState: db,
		app:     a,
	}
}

func contractCMD(a *appState) *cobra.Command {
	state := newContractState(a)
	contract := &cobra.Command{
		Use:     "contract",
		Short:   "Manage the contracts related to the chain",
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

	feeCmd.AddCommand(state.getFee(feeCmd), state.setFee(feeCmd), state.claimFee(feeCmd))

	deployCmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy centralized connection contract",
	}

	contract.AddCommand(feeCmd, deployCmd)

	return contract
}

// getFeeCmd gets the fee for the chain
func (c *contractState) getFee(cmd *cobra.Command) *cobra.Command {
	getFeeCmd := &cobra.Command{
		Use:     "get",
		Short:   "Get the fee set for the chain",
		Aliases: []string{"g"},
		Example: strings.TrimSpace(fmt.Sprintf(`$ %s contract fee get [chain-id]`, appName)),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := c.getSocket(c.app)
			if err != nil {
				return err
			}
			defer client.Close()
			res, err := client.GetFee(c.chain, true)
			if err != nil {
				return err
			}
			printLabels("Chain", "Fee")
			printValues(c.chain, res.Fee)
			return nil
		},
	}
	cmd.Flags().StringVar(&c.chain, "chain", "", "Chain ID")
	if err := cmd.MarkFlagRequired("chain"); err != nil {
		panic(err)
	}
	return getFeeCmd
}

// setFeeCmd sets the fee for the chain
func (c *contractState) setFee(cmd *cobra.Command) *cobra.Command {
	setFeeCmd := &cobra.Command{
		Use:     "set",
		Short:   "Set the fee for the chain",
		Aliases: []string{"s"},
		Example: strings.TrimSpace(fmt.Sprintf(`$ %s contract fee set [chain-id] [msg-fee] [res-fee]`, appName)),
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := c.getSocket(c.app)
			if err != nil {
				return err
			}
			defer client.Close()
			if err := client.SetFee(c.chain, big.NewInt(c.msgFee), big.NewInt(c.resFee)); err != nil {
				return err
			}
			printLabels("Status")
			printValues("Success")
			return nil
		},
	}
	cmd.Flags().StringVar(&c.chain, "chain", "", "Chain ID")
	cmd.Flags().Int64Var(&c.msgFee, "msg-fee", 0, "Message Fee")
	cmd.Flags().Int64Var(&c.resFee, "res-fee", 0, "Response Fee")
	cmd.MarkFlagsRequiredTogether("chain", "msg-fee", "res-fee")
	return setFeeCmd
}

// claimFeeCmd claims the fee for the chain
func (c *contractState) claimFee(cmd *cobra.Command) *cobra.Command {
	claimFeeCmd := &cobra.Command{
		Use:     "claim",
		Short:   "Claim the fee for the chain",
		Aliases: []string{"c"},
		Example: strings.TrimSpace(fmt.Sprintf(`$ %s contract fee claim [chain-id]`, appName)),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := c.getSocket(c.app)
			if err != nil {
				return err
			}
			defer client.Close()
			if err := client.ClaimFee(c.chain); err != nil {
				return err
			}
			printLabels("Status")
			printValues("Success")
			return nil
		},
	}
	cmd.Flags().StringVar(&c.chain, "chain", "", "Chain ID")
	if err := cmd.MarkFlagRequired("chain"); err != nil {
		panic(err)
	}
	return claimFeeCmd
}
