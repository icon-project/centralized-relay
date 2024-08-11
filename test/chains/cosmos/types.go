package cosmos

import (
	"time"

	"github.com/docker/docker/client"
	"github.com/icon-project/centralized-relay/test/chains"
	"github.com/icon-project/centralized-relay/test/testsuite/testconfig"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/icon-project/icon-bridge/common/codec"
)

type CosmosRemotenet struct {
	cfg          chains.ChainConfig
	filepath     map[string]string
	IBCAddresses map[string]string     `json:"addresses"`
	Wallets      map[string]ibc.Wallet `json:"wallets"`
	log          *zap.Logger
	DockerClient *client.Client
	Network      string
	testconfig   *testconfig.Chain
	testName     string
	GrpcConn     *grpc.ClientConn
	Rpcclient    *rpchttp.HTTP
}

func (c *CosmosRemotenet) Config() chains.ChainConfig {
	return c.cfg
}

type Query struct {
	GetClientState            *map[string]interface{} `json:"get_client_state,omitempty"`
	GetAdmin                  *GetAdmin               `json:"get_admin,omitempty"`
	GetProtocolFee            *GetProtocolFee         `json:"get_protocol_fee,omitempty"`
	GetNextClientSequence     *map[string]interface{} `json:"get_next_client_sequence,omitempty"`
	HasPacketReceipt          *map[string]interface{} `json:"has_packet_receipt,omitempty"`
	GetNextSequenceReceive    *map[string]interface{} `json:"get_next_sequence_receive,omitempty"`
	GetConnection             *map[string]interface{} `json:"get_connection,omitempty"`
	GetChannel                *map[string]interface{} `json:"get_channel,omitempty"`
	GetNextConnectionSequence *map[string]interface{} `json:"get_next_connection_sequence,omitempty"`
	GetNextChannelSequence    *map[string]interface{} `json:"get_next_channel_sequence,omitempty"`
}

type SetAdmin struct {
	SetAdmin struct {
		Address string `json:"address"`
	} `json:"set_admin"`
}

type UpdateAdmin struct {
	UpdateAdmin struct {
		Address string `json:"address"`
	} `json:"update_admin"`
}

type XcallInit struct {
	TimeoutHeight int    `json:"timeout_height"`
	IbcHost       string `json:"ibc_host"`
}

type MockAppInit struct {
	IBCHost int    `json:"ibc_host"`
	PortId  string `json:"port_id"`
	Order   string `json:"order"`
	Denom   string `json:"denom"`
}

type DappInit struct {
	Address string `json:"address"`
}

type GetProtocolFee struct{}

type GetAdmin struct{}

type Result struct {
	NodeInfo struct {
		ProtocolVersion struct {
			P2P   string `json:"p2p"`
			Block string `json:"block"`
			App   string `json:"app"`
		} `json:"protocol_version"`
		ID         string `json:"id"`
		ListenAddr string `json:"listen_addr"`
		Network    string `json:"network"`
		Version    string `json:"version"`
		Channels   string `json:"channels"`
		Moniker    string `json:"moniker"`
		Other      struct {
			TxIndex    string `json:"tx_index"`
			RPCAddress string `json:"rpc_address"`
		} `json:"other"`
	} `json:"NodeInfo"`
	SyncInfo struct {
		LatestBlockHash     string    `json:"latest_block_hash"`
		LatestAppHash       string    `json:"latest_app_hash"`
		LatestBlockHeight   string    `json:"latest_block_height"`
		LatestBlockTime     time.Time `json:"latest_block_time"`
		EarliestBlockHash   string    `json:"earliest_block_hash"`
		EarliestAppHash     string    `json:"earliest_app_hash"`
		EarliestBlockHeight string    `json:"earliest_block_height"`
		EarliestBlockTime   time.Time `json:"earliest_block_time"`
		CatchingUp          bool      `json:"catching_up"`
	} `json:"SyncInfo"`
	ValidatorInfo struct {
		Address string `json:"Address"`
		PubKey  struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"PubKey"`
		VotingPower string `json:"VotingPower"`
	} `json:"ValidatorInfo"`
}

type LogEvent struct {
	MsgIndex int    `json:"msg_index"`
	Log      string `json:"log"`
	Events   []struct {
		Type       string `json:"type"`
		Attributes []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
			Index bool   `json:"index"`
		} `json:"attributes"`
	} `json:"events"`
}

type TxResul struct {
	Height    string        `json:"height"`
	Txhash    string        `json:"txhash"`
	Codespace string        `json:"codespace"`
	Code      int           `json:"code"`
	Data      string        `json:"data"`
	RawLog    string        `json:"raw_log"`
	Logs      []interface{} `json:"logs"`
	Info      string        `json:"info"`
	GasWanted string        `json:"gas_wanted"`
	GasUsed   string        `json:"gas_used"`
	Tx        interface{}   `json:"tx"`
	Timestamp string        `json:"timestamp"`
	Events    []struct {
		Type       string `json:"type"`
		Attributes []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
			Index bool   `json:"index"`
		} `json:"attributes"`
	} `json:"events"`
}

type CallServiceMessageType int64

const (
	CallServiceRequest  CallServiceMessageType = 0
	CallServiceResponse CallServiceMessageType = 1
)

type CSMessageResponseType int64

const (
	SUCCESS   CSMessageResponseType = 0
	FAILURE   CSMessageResponseType = -1
	IBC_ERROR CSMessageResponseType = -2
)

type CallServiceMessage struct {
	MessageType CallServiceMessageType
	Payload     []byte
}

func RlpEncodeRequest(request CallServiceMessageRequest, callType CallServiceMessageType) ([]byte, error) {

	data := CallServiceMessage{

		MessageType: CallServiceRequest,

		Payload: codec.RLP.MustMarshalToBytes(request),
	}

	return codec.RLP.MarshalToBytes(data)

}

func RlpDecode(raw_data []byte) (CallServiceMessage, error) {

	var callservicemessage = CallServiceMessage{}

	codec.RLP.UnmarshalFromBytes(raw_data, &callservicemessage)

	return CallServiceMessage{}, nil
}

type CallServiceMessageRequest struct {
	From       string `json:"from"`
	To         string `json:"to"`
	SequenceNo uint64 `json:"sequence_no"`
	Rollback   bool   `json:"rollback"`
	Data       []byte `json:"data"`
}

func (cr *CallServiceMessageRequest) RlpEncode() ([]byte, error) {

	return codec.RLP.MarshalToBytes(cr)
}

type CallServiceMessageResponse struct {
	SequenceNo uint64
	Code       CSMessageResponseType
	Msg        string
}

func (csr *CallServiceMessageResponse) RlpEncode() ([]byte, error) {
	return codec.RLP.MarshalToBytes(csr)
}

type QueryContractResponse struct {
	Contracts []string `json:"contracts"`
}

type CosmosTx struct {
	TxHash string `json:"txhash"`
	Code   int    `json:"code"`
	RawLog string `json:"raw_log"`
}

type CodeInfo struct {
	CodeID string `json:"code_id"`
}

type CodeInfosResponse struct {
	CodeInfos []CodeInfo `json:"code_infos"`
}
