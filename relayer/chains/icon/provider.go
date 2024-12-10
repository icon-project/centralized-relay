package icon

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/gorilla/websocket"
	"github.com/icon-project/centralized-relay/relayer/chains/icon/types"
	"github.com/icon-project/centralized-relay/relayer/events"
	"github.com/icon-project/centralized-relay/relayer/kms"
	"github.com/icon-project/centralized-relay/relayer/provider"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/icon-project/goloop/module"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Config struct {
	provider.CommonConfig `json:",inline" yaml:",inline"`
	StepMin               int64 `json:"step-min" yaml:"step-min"`
	StepLimit             int64 `json:"step-limit" yaml:"step-limit"`
	StepAdjustment        int64 `json:"step-adjustment" yaml:"step-adjustment"`
}

// NewProvider returns new Icon provider
func (c *Config) NewProvider(ctx context.Context, log *zap.Logger, homepath string, debug bool, chainName string) (provider.ChainProvider, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	if err := c.sanitize(); err != nil {
		return nil, err
	}

	client := NewClient(ctx, c.RPCUrl, log)
	NetworkInfo, err := client.GetNetworkInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get network id: %v", err)
	}

	c.ChainName = chainName
	c.HomeDir = homepath

	return &Provider{
		log:       log.With(zap.Stringp("nid ", &c.NID), zap.Stringp("name", &c.ChainName)),
		client:    client,
		cfg:       c,
		networkID: NetworkInfo.NetworkID,
		contracts: c.eventMap(),
	}, nil
}

func (c *Config) Validate() error {
	if c.RPCUrl == "" {
		return fmt.Errorf("icon provider rpc endpoint is empty")
	}

	if err := c.Contracts.Validate(); err != nil {
		return fmt.Errorf("contracts are not valid: %s", err)
	}

	// TODO: validation for keystore
	// TODO: contractaddress validation
	// TODO: account should have some balance no balance then use another accoutn

	return nil
}

func (c *Config) sanitize() error {
	if c.StepAdjustment == 0 {
		c.StepAdjustment = 50
	}
	if c.Decimals == 0 {
		c.Decimals = providerTypes.DefaultCoinDecimals
	}
	return nil
}

func (p *Config) SetWallet(addr string) {
	p.Address = addr
}

func (p *Config) GetWallet() string {
	return p.Address
}

// Enabled returns true if the chain is enabled
func (c *Config) Enabled() bool {
	return !c.Disabled
}

type Provider struct {
	log                 *zap.Logger
	cfg                 *Config
	wallet              module.Wallet
	client              *Client
	kms                 kms.KMS
	contracts           map[string]providerTypes.EventMap
	networkID           types.HexInt
	LastSavedHeightFunc func() uint64
}

func (p *Provider) GetLastProcessedBlockHeight(ctx context.Context) (uint64, error) {
	return p.GetLastSavedBlockHeight(), nil
}

func (p *Provider) QueryBlockMessages(ctx context.Context, fromHeight, toHeight uint64) ([]*providerTypes.Message, error) {
	var messages []*providerTypes.Message
	eventReq := &types.EventRequest{
		Height:           types.NewHexInt(int64(fromHeight)),
		EventFilter:      p.GetMonitorEventFilters(),
		Logs:             types.NewHexInt(1),
		ProgressInterval: types.NewHexInt(25),
	}
	ctxMonitorBlock, cancelMonitorBlock := context.WithTimeout(ctx, 5*time.Second)
	defer cancelMonitorBlock()
	err := p.client.MonitorEvent(ctxMonitorBlock, eventReq, nil, func(v *types.EventNotification, outgoing chan *providerTypes.BlockInfo) error {
		if !errors.Is(ctx.Err(), context.Canceled) {
			msgs, err := p.parseMessageEvent(v)
			if err != nil {
				p.log.Error("failed to parse message event", zap.Error(err))
				return err
			}
			shouldBreak := false
			for _, msg := range msgs {
				if msg.MessageHeight > toHeight {
					shouldBreak = true
				} else {
					p.log.Info("Found eventlog",
						zap.Uint64("height", msg.MessageHeight),
						zap.String("target_network", msg.Dst),
						zap.Uint64("sn", msg.Sn.Uint64()),
						zap.String("tx_hash", v.Hash.String()),
						zap.String("event_type", msg.EventType),
					)
					messages = append(messages, msg)
				}
			}
			if shouldBreak {
				cancelMonitorBlock()
			}
		}
		return nil
	}, func(conn *websocket.Conn, err error) {})
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return messages, nil
		}
		return nil, err
	}
	return messages, nil
}

func (p *Provider) NID() string {
	return p.cfg.NID
}

func (p *Provider) Init(ctx context.Context, homepath string, kms kms.KMS) error {
	p.kms = kms
	return nil
}

func (p *Provider) Type() string {
	return "icon"
}

func (p *Provider) Config() provider.Config {
	return p.cfg
}

func (p *Provider) Name() string {
	return p.cfg.ChainName
}

func (p *Provider) NetworkID() types.HexInt {
	return p.networkID
}

func (p *Provider) Wallet() (module.Wallet, error) {
	if p.wallet == nil {
		if err := p.RestoreKeystore(context.Background()); err != nil {
			return nil, err
		}
	}
	return p.wallet, nil
}

func (p *Provider) FinalityBlock(ctx context.Context) uint64 {
	return p.cfg.FinalityBlock
}

// MessageReceived checks if the message is received
func (p *Provider) MessageReceived(ctx context.Context, messageKey *providerTypes.MessageKey) (bool, error) {
	switch messageKey.EventType {
	case events.EmitMessage:
		callParam := p.prepareCallParams(MethodGetReceipts, p.cfg.Contracts[providerTypes.ConnectionContract], map[string]interface{}{
			"srcNetwork": messageKey.Src,
			"_connSn":    types.NewHexInt(messageKey.Sn.Int64()),
		})
		var status types.HexInt
		return status == types.NewHexInt(1), p.client.Call(callParam, &status)
	case events.CallMessage:
		return false, nil
	case events.RollbackMessage:
		return false, nil
	default:
		return true, fmt.Errorf("unknown event type")
	}
}

// ClusterMessageReceived checks if the message is received
func (p *Provider) ClusterMessageReceived(ctx context.Context, message *providerTypes.Message) (bool, error) {
	srcHeight := message.WrappedSourceHeight
	if srcHeight == nil {
		srcHeight = big.NewInt(int64(message.MessageHeight))
	}
	callParam := p.prepareCallParams(MethodClusterMsgReceived, p.cfg.Contracts[providerTypes.AggregationContract], map[string]interface{}{
		"relayer":            p.cfg.GetWallet(),
		"srcNetwork":         message.Src,
		"srcContractAddress": message.SrcConnAddress,
		"srcSn":              types.NewHexInt(message.Sn.Int64()),
		"srcHeight":          srcHeight,
		"dstNetwork":         message.Dst,
		"dstContractAddress": message.DstConnAddress,
		"data":               types.NewHexBytes(message.Data),
	})
	var status types.HexInt
	err := p.client.Call(callParam, &status)
	return status == types.NewHexInt(1), err
}

// ClusterMessageAcknowledged checks if the message is acknowledged or not
func (p *Provider) ClusterMessageAcknowledged(ctx context.Context, message *providerTypes.Message) (bool, error) {
	srcHeight := message.WrappedSourceHeight
	if srcHeight == nil {
		srcHeight = big.NewInt(int64(message.MessageHeight))
	}
	callParam := p.prepareCallParams(MethodPacketAcknowledged, p.cfg.Contracts[providerTypes.AggregationContract], map[string]interface{}{
		"srcNetwork":         message.Src,
		"srcContractAddress": message.SrcConnAddress,
		"srcSn":              types.NewHexInt(message.Sn.Int64()),
		"srcHeight":          types.NewHexInt(srcHeight.Int64()),
		"dstNetwork":         message.Dst,
		"dstContractAddress": message.DstConnAddress,
		"data":               types.NewHexBytes(message.Data),
	})
	var status types.HexInt
	err := p.client.Call(callParam, &status)
	return status == types.NewHexInt(1), err
}

// ReverseMessage reverts a message
func (p *Provider) RevertMessage(ctx context.Context, sn *big.Int) error {
	params := map[string]interface{}{"sn": types.NewHexInt(sn.Int64())}
	message := p.NewIconMessage(types.Address(p.cfg.Contracts[providerTypes.ConnectionContract]), params, MethodRevertMessage)
	txHash, err := p.SendTransaction(ctx, message)
	if err != nil {
		return err
	}
	txr, err := p.client.WaitForResults(ctx, &types.TransactionHashParam{Hash: types.NewHexBytes(txHash)})
	if err != nil {
		return err
	}
	if txr.Status != types.NewHexInt(1) {
		return fmt.Errorf("failed: %s", txr.TxHash)
	}
	return nil
}

// SetAdmin sets the admin address of the bridge contract
func (p *Provider) SetAdmin(ctx context.Context, admin string) error {
	callParam := map[string]interface{}{
		"_relayer": admin,
	}
	message := p.NewIconMessage(types.Address(p.cfg.Contracts[providerTypes.ConnectionContract]), callParam, MethodSetAdmin)

	txHash, err := p.SendTransaction(ctx, message)
	if err != nil {
		return fmt.Errorf("SetAdmin: %v", err)
	}
	txr, err := p.client.WaitForResults(ctx, &types.TransactionHashParam{Hash: types.HexBytes(txHash)})
	if err != nil {
		return fmt.Errorf("SetAdmin: WaitForResults: %v", err)
	}
	if txr.Status != types.NewHexInt(1) {
		return fmt.Errorf("SetAdmin: failed to set admin: %s", txr.TxHash)
	}
	return nil
}

// GetFee
func (p *Provider) GetFee(ctx context.Context, networkID string, responseFee bool) (uint64, error) {
	callParam := p.prepareCallParams(MethodGetFee, p.cfg.Contracts[providerTypes.ConnectionContract], map[string]interface{}{
		"to":       networkID,
		"response": types.NewHexInt(1),
	})

	var status types.HexInt
	if err := p.client.Call(callParam, &status); err != nil {
		return 0, fmt.Errorf("GetFee: %v", err)
	}
	fee, err := status.BigInt()
	if err != nil {
		return 0, fmt.Errorf("GetFee: %v", err)
	}
	return fee.Uint64(), nil
}

// SetFees
func (p *Provider) SetFee(ctx context.Context, networkID string, msgFee, resFee *big.Int) error {
	callParam := map[string]interface{}{
		"networkId":   networkID,
		"messageFee":  types.NewHexInt(msgFee.Int64()),
		"responseFee": types.NewHexInt(resFee.Int64()),
	}

	msg := p.NewIconMessage(types.Address(p.cfg.Contracts[providerTypes.ConnectionContract]), callParam, MethodSetFee)
	txHash, err := p.SendTransaction(ctx, msg)
	if err != nil {
		return fmt.Errorf("SetFee: %v", err)
	}
	txr, err := p.client.WaitForResults(ctx, &types.TransactionHashParam{Hash: types.NewHexBytes(txHash)})
	if err != nil {
		return fmt.Errorf("SetFee: WaitForResults: %v", err)
	}
	if txr.Status != types.NewHexInt(1) {
		return fmt.Errorf("SetFee: failed to claim fees: %s", txr.TxHash)
	}
	return nil
}

// ClaimFees
func (p *Provider) ClaimFee(ctx context.Context) error {
	msg := p.NewIconMessage(types.Address(p.cfg.Contracts[providerTypes.ConnectionContract]), map[string]interface{}{}, MethodClaimFees)
	txHash, err := p.SendTransaction(ctx, msg)
	if err != nil {
		return fmt.Errorf("ClaimFees: %v", err)
	}
	txr, err := p.client.WaitForResults(ctx, &types.TransactionHashParam{Hash: types.NewHexBytes(txHash)})
	if err != nil {
		return fmt.Errorf("ClaimFees: WaitForResults: %v", err)
	}
	if txr.Status != types.NewHexInt(1) {
		return fmt.Errorf("ClaimFees: failed to claim fees: %s", txr.TxHash)
	}
	return nil
}

// ExecuteRollback
func (p *Provider) ExecuteRollback(ctx context.Context, sn int64) error {
	params := map[string]interface{}{"_sn": types.NewHexInt(sn)}
	message := p.NewIconMessage(types.Address(p.cfg.Contracts[providerTypes.XcallContract]), params, MethodExecuteRollback)
	txHash, err := p.SendTransaction(ctx, message)
	if err != nil {
		return err
	}
	txr, err := p.client.WaitForResults(ctx, &types.TransactionHashParam{Hash: types.NewHexBytes(txHash)})
	if err != nil {
		return err
	}
	if txr.Status != types.NewHexInt(1) {
		return fmt.Errorf("failed: %s", txr.TxHash)
	}
	return nil
}

// SetLastSavedBlockHeightFunc sets the function to save the last saved block height
func (p *Provider) SetLastSavedHeightFunc(f func() uint64) {
	p.LastSavedHeightFunc = f
}

// GetLastSavedBlockHeight returns the last saved block height
func (p *Provider) GetLastSavedBlockHeight() uint64 {
	return p.LastSavedHeightFunc()
}

func (p *Config) GetConnContract() string {
	return p.Contracts[providerTypes.ConnectionContract]
}
