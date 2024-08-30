package stacks

import (
	"fmt"
	"net/http"
)

type StacksNetwork struct {
	coreAPIURL string
	version    TransactionVersion
	chainID    ChainID
}

func NewStacksMainnet() *StacksNetwork {
	return &StacksNetwork{
		coreAPIURL: "https://api.mainnet.hiro.so",
		version:    TransactionVersionMainnet,
		chainID:    ChainIDMainnet,
	}
}

func NewStacksTestnet() *StacksNetwork {
	return &StacksNetwork{
		coreAPIURL: "https://api.testnet.hiro.so",
		version:    TransactionVersionTestnet,
		chainID:    ChainIDTestnet,
	}
}

func (n *StacksNetwork) GetAccountAPIURL(address string) string {
	return fmt.Sprintf("%s/v2/accounts/%s?proof=0", n.coreAPIURL, address)
}

func (n *StacksNetwork) GetBroadcastAPIURL() string {
	return fmt.Sprintf("%s/v2/transactions", n.coreAPIURL)
}

func (n *StacksNetwork) GetTransferFeeEstimateAPIURL() string {
	return fmt.Sprintf("%s/v2/fees/transfer", n.coreAPIURL)
}

func (n *StacksNetwork) GetTransactionFeeEstimateAPIURL() string {
	return fmt.Sprintf("%s/v2/fees/transaction", n.coreAPIURL)
}

func (n *StacksNetwork) FetchFn(url string) (*http.Response, error) {
	return http.Get(url)
}
