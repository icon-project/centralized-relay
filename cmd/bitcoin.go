package cmd

import (
	"github.com/icon-project/centralized-relay/relayer/chains/bitcoin"
	"github.com/spf13/cobra"
)

func bitcoinCmd(a *appState) *cobra.Command {
	bitcoinCmd := &cobra.Command{
		Use:   "bitcoin",
		Short: "Run Bitcoin Relayer",
		Run: func(cmd *cobra.Command, args []string) {
			bitcoin.RunApp()
		},
	}

	return bitcoinCmd
}
