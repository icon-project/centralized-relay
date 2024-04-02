package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

type contractState struct {
	*dbState
	app     *appState
	chain   string
	network string
	msgFee  uint64
	resFee  uint64
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
	}

	feeCmd.AddCommand(state.getFee(), state.setFee(), state.claimFee())

	deployCmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy centralized connection contract",
	}

	contract.AddCommand(feeCmd, deployCmd)

	return contract
}

// getFeeCmd gets the fee for the chain
func (c *contractState) getFee() *cobra.Command {
	getFeeCmd := &cobra.Command{
		Use:     "get",
		Short:   "Get the fee set for the chain",
		Aliases: []string{"g"},
		Example: strings.TrimSpace(fmt.Sprintf(`$ %s contract fee get [chain-id]`, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := c.getSocket(c.app)
			if err != nil {
				return err
			}
			defer client.Close()
			defer c.closeSocket()
			res, err := client.GetFee(c.chain, c.network, true)
			if err != nil {
				return err
			}
			printLabels("Chain", "Fee")
			printValues(c.network, res.Fee)
			return nil
		},
	}
	getFeeCmd.Flags().StringVar(&c.chain, "chain", "", "Chain ID")
	getFeeCmd.Flags().StringVar(&c.network, "network", "", "Network ID")
	getFeeCmd.MarkFlagsRequiredTogether("chain", "network")
	return getFeeCmd
}

// setFeeCmd sets the fee for the chain
func (c *contractState) setFee() *cobra.Command {
	setFeeCmd := &cobra.Command{
		Use:     "set",
		Short:   "Set the fee for the chain",
		Aliases: []string{"s"},
		Example: strings.TrimSpace(fmt.Sprintf(`$ %s contract fee set [chain-id] [msg-fee] [res-fee]`, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := c.getSocket(c.app)
			if err != nil {
				return err
			}
			defer client.Close()
			defer c.closeSocket()
			res, err := client.SetFee(c.chain, c.network, c.msgFee, c.resFee)
			if err != nil {
				return err
			}
			printLabels("Status")
			printValues(res.Status)
			return nil
		},
	}
	setFeeCmd.Flags().StringVar(&c.chain, "chain", "", "Chain NID")
	setFeeCmd.Flags().StringVar(&c.network, "network", "", "Network ID")
	setFeeCmd.Flags().Uint64Var(&c.msgFee, "msg-fee", 0, "Message Fee")
	setFeeCmd.Flags().Uint64Var(&c.resFee, "res-fee", 0, "Response Fee")
	setFeeCmd.MarkFlagsRequiredTogether("chain", "network", "msg-fee", "res-fee")
	return setFeeCmd
}

// claimFeeCmd claims the fee for the chain
func (c *contractState) claimFee() *cobra.Command {
	claimFeeCmd := &cobra.Command{
		Use:     "claim",
		Short:   "Claim the fee for the chain",
		Aliases: []string{"cm"},
		Example: strings.TrimSpace(fmt.Sprintf(`$ %s contract fee claim [chain-nid]`, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := c.getSocket(c.app)
			if err != nil {
				return err
			}
			defer client.Close()
			defer c.closeSocket()
			res, err := client.ClaimFee(c.chain)
			if err != nil {
				return err
			}
			printLabels("Status")
			printValues(res.Status)
			return nil
		},
	}
	claimFeeCmd.Flags().StringVar(&c.chain, "chain", "", "Chain NID")
	if err := claimFeeCmd.MarkFlagRequired("chain"); err != nil {
		panic(err)
	}
	return claimFeeCmd
}
