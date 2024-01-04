package cosmos

import (
	"context"
	"fmt"
	abciTypes "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	txTypes "github.com/cosmos/cosmos-sdk/types/tx"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/icon-project/centralized-relay/relayer/provider"
	"github.com/icon-project/centralized-relay/relayer/types"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"io"
	"os"
	"sync"
	"time"
)

const (
	ChainType string = "wasm"
)

type Provider struct {
	logger *zap.Logger
	config *ProviderConfig

	client *grpc.ClientConn

	Input  io.Reader
	Output io.Writer
	//Cdc       Codec

	nextAccountSeq uint64
	txMu           sync.Mutex

	// metrics to monitor the provider
	//TotalFees   sdk.Coins
	totalFeesMu sync.Mutex
}

type ProviderConfig struct {
	ChainName string `json:"-" yaml:"-"`
	ChainID   string `json:"chain_id" yaml:"chain-id"`
	NID       string `json:"nid" yaml:"nid"`

	Keystore string `json:"keystore" yaml:"keystore"`
	Password string `json:"password" yaml:"password"`

	RPCUrl string `json:"rpc-url" yaml:"rpc-url"`

	ContractAddress string `json:"contract-address" yaml:"contract-address"`

	AccountPrefix string `json:"account-prefix" yaml:"account-prefix"`

	GasAdjustment float64 `json:"gas-adjustment" yaml:"gas-adjustment"`
	GasPrices     string  `json:"gas-prices" yaml:"gas-prices"`
	MinGasAmount  uint64  `json:"min-gas-amount" yaml:"min-gas-amount"`
	MaxGasAmount  uint64  `json:"max-gas-amount" yaml:"max-gas-amount"`

	Timeout          string `json:"timeout" yaml:"timeout"`
	BlockTimeout     string `json:"block-timeout" yaml:"block-timeout"`
	SignModeStr      string `json:"sign-mode" yaml:"sign-mode"`
	SigningAlgorithm string `json:"signing-algorithm" yaml:"signing-algorithm"`
}

func (pc ProviderConfig) NewProvider(logger *zap.Logger, homePath string, debug bool, chainName string) (provider.ChainProvider, error) {
	if err := pc.Validate(); err != nil {
		return nil, err
	}

	pc.ChainName = chainName

	cp := &Provider{
		logger: logger,
		Input:  os.Stdin,
		Output: os.Stdout,
	}

	return cp, nil
}

func (pc ProviderConfig) Validate() error {
	if _, err := time.ParseDuration(pc.Timeout); err != nil {
		return fmt.Errorf("invalid Timeout: %w", err)
	}
	return nil
}

func (p *Provider) QueryLatestHeight(ctx context.Context) (uint64, error) {
	serviceClient := cmtservice.NewServiceClient(p.client)
	res, err := serviceClient.GetLatestBlock(ctx, &cmtservice.GetLatestBlockRequest{})
	if err != nil {
		return 0, err
	}

	return uint64(res.GetBlock().Header.Height), nil
}

func (p *Provider) QueryTransactionReceipt(ctx context.Context, txHash string) (*types.Receipt, error) {
	serviceClient := txTypes.NewServiceClient(p.client)
	res, err := serviceClient.GetTx(ctx, &txTypes.GetTxRequest{Hash: txHash})
	if err != nil {
		return nil, err
	}
	return &types.Receipt{
		TxHash: txHash,
		Height: uint64(res.TxResponse.Height),
		Status: abciTypes.CodeTypeOK == res.TxResponse.Code,
	}, nil
}

func (p *Provider) NID() string {
	return p.config.NID
}

func (p *Provider) ChainName() string {
	return p.config.ChainName
}

func (p *Provider) Init(ctx context.Context) error {
	return nil
}

func (p *Provider) Type() string {
	return ChainType
}

func (p *Provider) ProviderConfig() provider.ProviderConfig {
	return *p.config
}

func (p *Provider) Listener(ctx context.Context, lastSavedHeight uint64, blockInfo chan types.BlockInfo) error {
	return nil
}

func (p *Provider) Route(ctx context.Context, message *types.Message, callback types.TxResponseFunc) error {
	return nil
}

func (p *Provider) ShouldReceiveMessage(ctx context.Context, message types.Message) (bool, error) {
	return true, nil
}

func (p *Provider) ShouldSendMessage(ctx context.Context, message types.Message) (bool, error) {
	return true, nil
}

func (p *Provider) MessageReceived(ctx context.Context, key types.MessageKey) (bool, error) {
	return false, nil
}

func (p *Provider) QueryBalance(ctx context.Context, addr string) (*types.Coin, error) {
	queryClient := bankTypes.NewQueryClient(p.client)

	res, err := queryClient.Balance(ctx, &bankTypes.QueryBalanceRequest{
		Address: addr,
		Denom:   "s",
	})
	if err != nil {
		return nil, err
	}

	return &types.Coin{
		Denom:  res.GetBalance().Denom,
		Amount: res.GetBalance().Amount.Uint64(),
	}, nil
}

func (p *Provider) GenerateMessage(ctx context.Context, messageKey *types.MessageKeyWithMessageHeight) (*types.Message, error) {
	return nil, nil
}

func (p *Provider) FinalityBlock(ctx context.Context) uint64 {
	return 0
}
