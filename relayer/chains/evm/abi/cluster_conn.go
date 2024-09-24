// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package abi

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// ClusterConnectionMetaData contains all meta data concerning the ClusterConnection contract.
var ClusterConnectionMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"addSigner\",\"inputs\":[{\"name\":\"_newSigner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"admin\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"claimFees\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"connSn\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getFee\",\"inputs\":[{\"name\":\"to\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"response\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[{\"name\":\"fee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getReceipt\",\"inputs\":[{\"name\":\"srcNetwork\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_connSn\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_relayer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_xCall\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isSigner\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"recoverSigner\",\"inputs\":[{\"name\":\"messageHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"recvMessage\",\"inputs\":[{\"name\":\"srcNetwork\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_connSn\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_msg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"recvMessageWithSignatures\",\"inputs\":[{\"name\":\"srcNetwork\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_connSn\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_msg\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_signedMessages\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"removeSigner\",\"inputs\":[{\"name\":\"_signer\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revertMessage\",\"inputs\":[{\"name\":\"sn\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"sendMessage\",\"inputs\":[{\"name\":\"to\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"svc\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"sn\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"_msg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"setAdmin\",\"inputs\":[{\"name\":\"_address\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setFee\",\"inputs\":[{\"name\":\"networkId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"messageFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"responseFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Message\",\"inputs\":[{\"name\":\"targetNetwork\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"sn\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"_msg\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]}]",
	Bin: "0x608060405234801561001057600080fd5b50611c47806100206000396000f3fe6080604052600436106100f35760003560e01c80637df73e271161008a578063b58b4cec11610059578063b58b4cec146102ae578063d294f093146102ce578063eb12d61e146102e3578063f851a4401461030357600080fd5b80637df73e27146102005780639664da0e1461024057806397aba7f91461026057806399f1fca71461029857600080fd5b8063485cc955116100c6578063485cc9551461017a578063522a901e1461019a578063704b6c02146101ad5780637d4c4f4a146101cd57600080fd5b80630a523d68146100f85780630e316ab71461011a5780632d3fb8231461013a57806343f08a891461015a575b600080fd5b34801561010457600080fd5b50610118610113366004611617565b610321565b005b34801561012657600080fd5b506101186101353660046116f8565b61058b565b34801561014657600080fd5b50610118610155366004611715565b6107f8565b34801561016657600080fd5b5061011861017536600461172e565b6108e9565b34801561018657600080fd5b5061011861019536600461177e565b6109c9565b6101186101a83660046117b7565b610b38565b3480156101b957600080fd5b506101186101c83660046116f8565b610d28565b3480156101d957600080fd5b506101ed6101e836600461185a565b610ddc565b6040519081526020015b60405180910390f35b34801561020c57600080fd5b5061023061021b3660046116f8565b60036020526000908152604090205460ff1681565b60405190151581526020016101f7565b34801561024c57600080fd5b5061023061025b3660046118a5565b610e4b565b34801561026c57600080fd5b5061028061027b3660046118e9565b610e83565b6040516001600160a01b0390911681526020016101f7565b3480156102a457600080fd5b506101ed60065481565b3480156102ba57600080fd5b506101186102c9366004611943565b61101c565b3480156102da57600080fd5b506101186111bf565b3480156102ef57600080fd5b506101186102fe3660046116f8565b61128d565b34801561030f57600080fd5b506005546001600160a01b0316610280565b306001600160a01b031663f851a4406040518163ffffffff1660e01b8152600401602060405180830381865afa15801561035f573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061038391906119b2565b6001600160a01b0316336001600160a01b0316146103bc5760405162461bcd60e51b81526004016103b3906119cf565b60405180910390fd5b806104025760405162461bcd60e51b8152602060048201526016602482015275139bc81cda59db985d1d5c995cc81c1c9bdd9a59195960521b60448201526064016103b3565b600084846040516104149291906119f4565b60405190819003902090506000826001600160401b038111156104395761043961151d565b604051908082528060200260200182016040528015610462578160200160208202803683370190505b5090506000805b848110156105625760006104d58588888581811061048957610489611a04565b905060200281019061049b9190611a1a565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250610e8392505050565b90506001600160a01b0381166105215760405162461bcd60e51b8152602060048201526011602482015270496e76616c6964207369676e617475726560781b60448201526064016103b3565b8084838151811061053457610534611a04565b6001600160a01b03909216602092830291909101909101528261055681611a76565b93505050600101610469565b50600061056e836113ef565b1115610580576105808989898961101c565b505050505050505050565b306001600160a01b031663f851a4406040518163ffffffff1660e01b8152600401602060405180830381865afa1580156105c9573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105ed91906119b2565b6001600160a01b0316336001600160a01b03161461061d5760405162461bcd60e51b81526004016103b3906119cf565b6001600160a01b03811660009081526003602052604090205460ff166106855760405162461bcd60e51b815260206004820152601760248201527f41646472657373206973206e6f7420616e206f776e657200000000000000000060448201526064016103b3565b6007546001106106d75760405162461bcd60e51b815260206004820152601e60248201527f4174206c65617374206f6e65206f776e6572206973207265717569726564000060448201526064016103b3565b60005b6007548110156107f457816001600160a01b03166007828154811061070157610701611a04565b6000918252602090912001546001600160a01b0316036107ec576007805461072b90600190611a8f565b8154811061073b5761073b611a04565b600091825260209091200154600780546001600160a01b03909216918390811061076757610767611a04565b9060005260206000200160006101000a8154816001600160a01b0302191690836001600160a01b0316021790555060078054806107a6576107a6611aa2565b60008281526020808220830160001990810180546001600160a01b03191690559092019092556001600160a01b03841682526003905260409020805460ff191690555050565b6001016106da565b5050565b306001600160a01b031663f851a4406040518163ffffffff1660e01b8152600401602060405180830381865afa158015610836573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061085a91906119b2565b6001600160a01b0316336001600160a01b03161461088a5760405162461bcd60e51b81526004016103b3906119cf565b6004805460405163b070f9e560e01b81529182018390526001600160a01b03169063b070f9e590602401600060405180830381600087803b1580156108ce57600080fd5b505af11580156108e2573d6000803e3d6000fd5b5050505050565b306001600160a01b031663f851a4406040518163ffffffff1660e01b8152600401602060405180830381865afa158015610927573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061094b91906119b2565b6001600160a01b0316336001600160a01b03161461097b5760405162461bcd60e51b81526004016103b3906119cf565b816000858560405161098e9291906119f4565b90815260200160405180910390208190555080600185856040516109b39291906119f4565b9081526040519081900360200190205550505050565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a008054600160401b810460ff1615906001600160401b0316600081158015610a0e5750825b90506000826001600160401b03166001148015610a2a5750303b155b905081158015610a38575080155b15610a565760405163f92ee8a960e01b815260040160405180910390fd5b845467ffffffffffffffff191660011785558315610a8057845460ff60401b1916600160401b1785555b600480546001600160a01b038089166001600160a01b031992831617909255600780546001810182556000919091527fa66cc928b5edb82af9bd49922954155ab7b0942694bea4ce44661d9a8736c688018054928a16928216831790556005805490911690911790558315610b2f57845460ff60401b19168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b50505050505050565b6004546001600160a01b03163314610b925760405162461bcd60e51b815260206004820152601f60248201527f4f6e6c79205863616c6c2063616e2063616c6c2073656e644d6573736167650060448201526064016103b3565b600080841315610c0c57604051633ea627a560e11b81523090637d4c4f4a90610bc4908b908b90600190600401611ae1565b602060405180830381865afa158015610be1573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c059190611b07565b9050610c80565b83600003610c8057604051633ea627a560e11b81523090637d4c4f4a90610c3c908b908b90600090600401611ae1565b602060405180830381865afa158015610c59573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c7d9190611b07565b90505b80341015610cc85760405162461bcd60e51b8152602060048201526015602482015274119959481a5cc81b9bdd0814dd59999a58da595b9d605a1b60448201526064016103b3565b60068054906000610cd883611a76565b91905055507f37be353f216cf7e33639101fd610c542e6a0c0109173fa1c1d8b04d34edb7c1b88886006548686604051610d16959493929190611b20565b60405180910390a15050505050505050565b306001600160a01b031663f851a4406040518163ffffffff1660e01b8152600401602060405180830381865afa158015610d66573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610d8a91906119b2565b6001600160a01b0316336001600160a01b031614610dba5760405162461bcd60e51b81526004016103b3906119cf565b600580546001600160a01b0319166001600160a01b0392909216919091179055565b600080600084604051610def9190611b7d565b908152604051908190036020019020549050821515600103610e42576000600185604051610e1d9190611b7d565b908152604051908190036020019020549050610e398183611b99565b92505050610e45565b90505b92915050565b6000600283604051610e5d9190611b7d565b908152604080516020928190038301902060009485529091529091205460ff1692915050565b60008151604114610ed65760405162461bcd60e51b815260206004820152601860248201527f496e76616c6964207369676e6174757265206c656e677468000000000000000060448201526064016103b3565b60208201516040830151606084015160001a601b811015610eff57610efc601b82611bac565b90505b8060ff16601b1480610f1457508060ff16601c145b610f605760405162461bcd60e51b815260206004820152601b60248201527f496e76616c6964207369676e6174757265202776272076616c7565000000000060448201526064016103b3565b6001610fb9876040517f19457468657265756d205369676e6564204d6573736167653a0a3332000000006020820152603c8101829052600090605c01604051602081830303815290604052805190602001209050919050565b6040805160008152602081018083529290925260ff841690820152606081018590526080810184905260a0016020604051602081039080840390855afa158015611007573d6000803e3d6000fd5b5050604051601f190151979650505050505050565b306001600160a01b031663f851a4406040518163ffffffff1660e01b8152600401602060405180830381865afa15801561105a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061107e91906119b2565b6001600160a01b0316336001600160a01b0316146110ae5760405162461bcd60e51b81526004016103b3906119cf565b6002846040516110be9190611b7d565b90815260408051602092819003830190206000868152925290205460ff161561111d5760405162461bcd60e51b81526020600482015260116024820152704475706c6963617465204d65737361676560781b60448201526064016103b3565b600160028560405161112f9190611b7d565b9081526040805160209281900383018120600088815293529120805460ff1916921515929092179091556004805463bbc22efd60e01b83526001600160a01b03169163bbc22efd916111879188918791879101611bc5565b600060405180830381600087803b1580156111a157600080fd5b505af11580156111b5573d6000803e3d6000fd5b5050505050505050565b306001600160a01b031663f851a4406040518163ffffffff1660e01b8152600401602060405180830381865afa1580156111fd573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061122191906119b2565b6001600160a01b0316336001600160a01b0316146112515760405162461bcd60e51b81526004016103b3906119cf565b6005546040516001600160a01b03909116904780156108fc02916000818181858888f1935050505015801561128a573d6000803e3d6000fd5b50565b306001600160a01b031663f851a4406040518163ffffffff1660e01b8152600401602060405180830381865afa1580156112cb573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112ef91906119b2565b6001600160a01b0316336001600160a01b03161461131f5760405162461bcd60e51b81526004016103b3906119cf565b6001600160a01b03811660009081526003602052604090205460ff16156113885760405162461bcd60e51b815260206004820152601c60248201527f4164647265737320697320616c726561647920616e207369676e65720000000060448201526064016103b3565b6007805460018181019092557fa66cc928b5edb82af9bd49922954155ab7b0942694bea4ce44661d9a8736c6880180546001600160a01b039093166001600160a01b031990931683179055600091825260036020526040909120805460ff19169091179055565b60008060009050600083516001600160401b038111156114115761141161151d565b60405190808252806020026020018201604052801561143a578160200160208202803683370190505b50905060005b845181101561151457600160005b848110156114b25786838151811061146857611468611a04565b60200260200101516001600160a01b031684828151811061148b5761148b611a04565b60200260200101516001600160a01b0316036114aa57600091506114b2565b60010161144e565b50801561150b578582815181106114cb576114cb611a04565b60200260200101518385815181106114e5576114e5611a04565b6001600160a01b03909216602092830291909101909101528361150781611a76565b9450505b50600101611440565b50909392505050565b634e487b7160e01b600052604160045260246000fd5b60006001600160401b038084111561154d5761154d61151d565b604051601f8501601f19908116603f011681019082821181831017156115755761157561151d565b8160405280935085815286868601111561158e57600080fd5b858560208301376000602087830101525050509392505050565b600082601f8301126115b957600080fd5b6115c883833560208501611533565b9392505050565b60008083601f8401126115e157600080fd5b5081356001600160401b038111156115f857600080fd5b60208301915083602082850101111561161057600080fd5b9250929050565b6000806000806000806080878903121561163057600080fd5b86356001600160401b038082111561164757600080fd5b6116538a838b016115a8565b975060208901359650604089013591508082111561167057600080fd5b61167c8a838b016115cf565b9096509450606089013591508082111561169557600080fd5b818901915089601f8301126116a957600080fd5b8135818111156116b857600080fd5b8a60208260051b85010111156116cd57600080fd5b6020830194508093505050509295509295509295565b6001600160a01b038116811461128a57600080fd5b60006020828403121561170a57600080fd5b8135610e42816116e3565b60006020828403121561172757600080fd5b5035919050565b6000806000806060858703121561174457600080fd5b84356001600160401b0381111561175a57600080fd5b611766878288016115cf565b90989097506020870135966040013595509350505050565b6000806040838503121561179157600080fd5b823561179c816116e3565b915060208301356117ac816116e3565b809150509250929050565b60008060008060008060006080888a0312156117d257600080fd5b87356001600160401b03808211156117e957600080fd5b6117f58b838c016115cf565b909950975060208a013591508082111561180e57600080fd5b61181a8b838c016115cf565b909750955060408a0135945060608a013591508082111561183a57600080fd5b506118478a828b016115cf565b989b979a50959850939692959293505050565b6000806040838503121561186d57600080fd5b82356001600160401b0381111561188357600080fd5b61188f858286016115a8565b925050602083013580151581146117ac57600080fd5b600080604083850312156118b857600080fd5b82356001600160401b038111156118ce57600080fd5b6118da858286016115a8565b95602094909401359450505050565b600080604083850312156118fc57600080fd5b8235915060208301356001600160401b0381111561191957600080fd5b8301601f8101851361192a57600080fd5b61193985823560208401611533565b9150509250929050565b6000806000806060858703121561195957600080fd5b84356001600160401b038082111561197057600080fd5b61197c888389016115a8565b955060208701359450604087013591508082111561199957600080fd5b506119a6878288016115cf565b95989497509550505050565b6000602082840312156119c457600080fd5b8151610e42816116e3565b6020808252600b908201526a27b7363ca932b630bcb2b960a91b604082015260600190565b8183823760009101908152919050565b634e487b7160e01b600052603260045260246000fd5b6000808335601e19843603018112611a3157600080fd5b8301803591506001600160401b03821115611a4b57600080fd5b60200191503681900382131561161057600080fd5b634e487b7160e01b600052601160045260246000fd5b600060018201611a8857611a88611a60565b5060010190565b81810381811115610e4557610e45611a60565b634e487b7160e01b600052603160045260246000fd5b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b604081526000611af5604083018587611ab8565b90508215156020830152949350505050565b600060208284031215611b1957600080fd5b5051919050565b606081526000611b34606083018789611ab8565b8560208401528281036040840152611b4d818587611ab8565b98975050505050505050565b60005b83811015611b74578181015183820152602001611b5c565b50506000910152565b60008251611b8f818460208701611b59565b9190910192915050565b80820180821115610e4557610e45611a60565b60ff8181168382160190811115610e4557610e45611a60565b6040815260008451806040840152611be4816060850160208901611b59565b601f01601f1916820182810360609081016020850152611c079082018587611ab8565b969550505050505056fea2646970667358221220ad774c60ed3f291d2291fc2b50ed8a7a5ddb71c9b4390fd2403db27bb73d46ae64736f6c63430008180033",
}

// ClusterConnectionABI is the input ABI used to generate the binding from.
// Deprecated: Use ClusterConnectionMetaData.ABI instead.
var ClusterConnectionABI = ClusterConnectionMetaData.ABI

// ClusterConnectionBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ClusterConnectionMetaData.Bin instead.
var ClusterConnectionBin = ClusterConnectionMetaData.Bin

// DeployClusterConnection deploys a new Ethereum contract, binding an instance of ClusterConnection to it.
func DeployClusterConnection(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ClusterConnection, error) {
	parsed, err := ClusterConnectionMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ClusterConnectionBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ClusterConnection{ClusterConnectionCaller: ClusterConnectionCaller{contract: contract}, ClusterConnectionTransactor: ClusterConnectionTransactor{contract: contract}, ClusterConnectionFilterer: ClusterConnectionFilterer{contract: contract}}, nil
}

// ClusterConnection is an auto generated Go binding around an Ethereum contract.
type ClusterConnection struct {
	ClusterConnectionCaller     // Read-only binding to the contract
	ClusterConnectionTransactor // Write-only binding to the contract
	ClusterConnectionFilterer   // Log filterer for contract events
}

// ClusterConnectionCaller is an auto generated read-only Go binding around an Ethereum contract.
type ClusterConnectionCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ClusterConnectionTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ClusterConnectionTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ClusterConnectionFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ClusterConnectionFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ClusterConnectionSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ClusterConnectionSession struct {
	Contract     *ClusterConnection // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ClusterConnectionCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ClusterConnectionCallerSession struct {
	Contract *ClusterConnectionCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// ClusterConnectionTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ClusterConnectionTransactorSession struct {
	Contract     *ClusterConnectionTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// ClusterConnectionRaw is an auto generated low-level Go binding around an Ethereum contract.
type ClusterConnectionRaw struct {
	Contract *ClusterConnection // Generic contract binding to access the raw methods on
}

// ClusterConnectionCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ClusterConnectionCallerRaw struct {
	Contract *ClusterConnectionCaller // Generic read-only contract binding to access the raw methods on
}

// ClusterConnectionTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ClusterConnectionTransactorRaw struct {
	Contract *ClusterConnectionTransactor // Generic write-only contract binding to access the raw methods on
}

// NewClusterConnection creates a new instance of ClusterConnection, bound to a specific deployed contract.
func NewClusterConnection(address common.Address, backend bind.ContractBackend) (*ClusterConnection, error) {
	contract, err := bindClusterConnection(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ClusterConnection{ClusterConnectionCaller: ClusterConnectionCaller{contract: contract}, ClusterConnectionTransactor: ClusterConnectionTransactor{contract: contract}, ClusterConnectionFilterer: ClusterConnectionFilterer{contract: contract}}, nil
}

// NewClusterConnectionCaller creates a new read-only instance of ClusterConnection, bound to a specific deployed contract.
func NewClusterConnectionCaller(address common.Address, caller bind.ContractCaller) (*ClusterConnectionCaller, error) {
	contract, err := bindClusterConnection(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ClusterConnectionCaller{contract: contract}, nil
}

// NewClusterConnectionTransactor creates a new write-only instance of ClusterConnection, bound to a specific deployed contract.
func NewClusterConnectionTransactor(address common.Address, transactor bind.ContractTransactor) (*ClusterConnectionTransactor, error) {
	contract, err := bindClusterConnection(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ClusterConnectionTransactor{contract: contract}, nil
}

// NewClusterConnectionFilterer creates a new log filterer instance of ClusterConnection, bound to a specific deployed contract.
func NewClusterConnectionFilterer(address common.Address, filterer bind.ContractFilterer) (*ClusterConnectionFilterer, error) {
	contract, err := bindClusterConnection(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ClusterConnectionFilterer{contract: contract}, nil
}

// bindClusterConnection binds a generic wrapper to an already deployed contract.
func bindClusterConnection(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ClusterConnectionMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ClusterConnection *ClusterConnectionRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ClusterConnection.Contract.ClusterConnectionCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ClusterConnection *ClusterConnectionRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ClusterConnection.Contract.ClusterConnectionTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ClusterConnection *ClusterConnectionRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ClusterConnection.Contract.ClusterConnectionTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ClusterConnection *ClusterConnectionCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ClusterConnection.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ClusterConnection *ClusterConnectionTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ClusterConnection.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ClusterConnection *ClusterConnectionTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ClusterConnection.Contract.contract.Transact(opts, method, params...)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_ClusterConnection *ClusterConnectionCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ClusterConnection.contract.Call(opts, &out, "admin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_ClusterConnection *ClusterConnectionSession) Admin() (common.Address, error) {
	return _ClusterConnection.Contract.Admin(&_ClusterConnection.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_ClusterConnection *ClusterConnectionCallerSession) Admin() (common.Address, error) {
	return _ClusterConnection.Contract.Admin(&_ClusterConnection.CallOpts)
}

// ConnSn is a free data retrieval call binding the contract method 0x99f1fca7.
//
// Solidity: function connSn() view returns(uint256)
func (_ClusterConnection *ClusterConnectionCaller) ConnSn(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ClusterConnection.contract.Call(opts, &out, "connSn")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ConnSn is a free data retrieval call binding the contract method 0x99f1fca7.
//
// Solidity: function connSn() view returns(uint256)
func (_ClusterConnection *ClusterConnectionSession) ConnSn() (*big.Int, error) {
	return _ClusterConnection.Contract.ConnSn(&_ClusterConnection.CallOpts)
}

// ConnSn is a free data retrieval call binding the contract method 0x99f1fca7.
//
// Solidity: function connSn() view returns(uint256)
func (_ClusterConnection *ClusterConnectionCallerSession) ConnSn() (*big.Int, error) {
	return _ClusterConnection.Contract.ConnSn(&_ClusterConnection.CallOpts)
}

// GetFee is a free data retrieval call binding the contract method 0x7d4c4f4a.
//
// Solidity: function getFee(string to, bool response) view returns(uint256 fee)
func (_ClusterConnection *ClusterConnectionCaller) GetFee(opts *bind.CallOpts, to string, response bool) (*big.Int, error) {
	var out []interface{}
	err := _ClusterConnection.contract.Call(opts, &out, "getFee", to, response)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetFee is a free data retrieval call binding the contract method 0x7d4c4f4a.
//
// Solidity: function getFee(string to, bool response) view returns(uint256 fee)
func (_ClusterConnection *ClusterConnectionSession) GetFee(to string, response bool) (*big.Int, error) {
	return _ClusterConnection.Contract.GetFee(&_ClusterConnection.CallOpts, to, response)
}

// GetFee is a free data retrieval call binding the contract method 0x7d4c4f4a.
//
// Solidity: function getFee(string to, bool response) view returns(uint256 fee)
func (_ClusterConnection *ClusterConnectionCallerSession) GetFee(to string, response bool) (*big.Int, error) {
	return _ClusterConnection.Contract.GetFee(&_ClusterConnection.CallOpts, to, response)
}

// GetReceipt is a free data retrieval call binding the contract method 0x9664da0e.
//
// Solidity: function getReceipt(string srcNetwork, uint256 _connSn) view returns(bool)
func (_ClusterConnection *ClusterConnectionCaller) GetReceipt(opts *bind.CallOpts, srcNetwork string, _connSn *big.Int) (bool, error) {
	var out []interface{}
	err := _ClusterConnection.contract.Call(opts, &out, "getReceipt", srcNetwork, _connSn)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetReceipt is a free data retrieval call binding the contract method 0x9664da0e.
//
// Solidity: function getReceipt(string srcNetwork, uint256 _connSn) view returns(bool)
func (_ClusterConnection *ClusterConnectionSession) GetReceipt(srcNetwork string, _connSn *big.Int) (bool, error) {
	return _ClusterConnection.Contract.GetReceipt(&_ClusterConnection.CallOpts, srcNetwork, _connSn)
}

// GetReceipt is a free data retrieval call binding the contract method 0x9664da0e.
//
// Solidity: function getReceipt(string srcNetwork, uint256 _connSn) view returns(bool)
func (_ClusterConnection *ClusterConnectionCallerSession) GetReceipt(srcNetwork string, _connSn *big.Int) (bool, error) {
	return _ClusterConnection.Contract.GetReceipt(&_ClusterConnection.CallOpts, srcNetwork, _connSn)
}

// IsSigner is a free data retrieval call binding the contract method 0x7df73e27.
//
// Solidity: function isSigner(address ) view returns(bool)
func (_ClusterConnection *ClusterConnectionCaller) IsSigner(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _ClusterConnection.contract.Call(opts, &out, "isSigner", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsSigner is a free data retrieval call binding the contract method 0x7df73e27.
//
// Solidity: function isSigner(address ) view returns(bool)
func (_ClusterConnection *ClusterConnectionSession) IsSigner(arg0 common.Address) (bool, error) {
	return _ClusterConnection.Contract.IsSigner(&_ClusterConnection.CallOpts, arg0)
}

// IsSigner is a free data retrieval call binding the contract method 0x7df73e27.
//
// Solidity: function isSigner(address ) view returns(bool)
func (_ClusterConnection *ClusterConnectionCallerSession) IsSigner(arg0 common.Address) (bool, error) {
	return _ClusterConnection.Contract.IsSigner(&_ClusterConnection.CallOpts, arg0)
}

// RecoverSigner is a free data retrieval call binding the contract method 0x97aba7f9.
//
// Solidity: function recoverSigner(bytes32 messageHash, bytes signature) pure returns(address)
func (_ClusterConnection *ClusterConnectionCaller) RecoverSigner(opts *bind.CallOpts, messageHash [32]byte, signature []byte) (common.Address, error) {
	var out []interface{}
	err := _ClusterConnection.contract.Call(opts, &out, "recoverSigner", messageHash, signature)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RecoverSigner is a free data retrieval call binding the contract method 0x97aba7f9.
//
// Solidity: function recoverSigner(bytes32 messageHash, bytes signature) pure returns(address)
func (_ClusterConnection *ClusterConnectionSession) RecoverSigner(messageHash [32]byte, signature []byte) (common.Address, error) {
	return _ClusterConnection.Contract.RecoverSigner(&_ClusterConnection.CallOpts, messageHash, signature)
}

// RecoverSigner is a free data retrieval call binding the contract method 0x97aba7f9.
//
// Solidity: function recoverSigner(bytes32 messageHash, bytes signature) pure returns(address)
func (_ClusterConnection *ClusterConnectionCallerSession) RecoverSigner(messageHash [32]byte, signature []byte) (common.Address, error) {
	return _ClusterConnection.Contract.RecoverSigner(&_ClusterConnection.CallOpts, messageHash, signature)
}

// AddSigner is a paid mutator transaction binding the contract method 0xeb12d61e.
//
// Solidity: function addSigner(address _newSigner) returns()
func (_ClusterConnection *ClusterConnectionTransactor) AddSigner(opts *bind.TransactOpts, _newSigner common.Address) (*types.Transaction, error) {
	return _ClusterConnection.contract.Transact(opts, "addSigner", _newSigner)
}

// AddSigner is a paid mutator transaction binding the contract method 0xeb12d61e.
//
// Solidity: function addSigner(address _newSigner) returns()
func (_ClusterConnection *ClusterConnectionSession) AddSigner(_newSigner common.Address) (*types.Transaction, error) {
	return _ClusterConnection.Contract.AddSigner(&_ClusterConnection.TransactOpts, _newSigner)
}

// AddSigner is a paid mutator transaction binding the contract method 0xeb12d61e.
//
// Solidity: function addSigner(address _newSigner) returns()
func (_ClusterConnection *ClusterConnectionTransactorSession) AddSigner(_newSigner common.Address) (*types.Transaction, error) {
	return _ClusterConnection.Contract.AddSigner(&_ClusterConnection.TransactOpts, _newSigner)
}

// ClaimFees is a paid mutator transaction binding the contract method 0xd294f093.
//
// Solidity: function claimFees() returns()
func (_ClusterConnection *ClusterConnectionTransactor) ClaimFees(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ClusterConnection.contract.Transact(opts, "claimFees")
}

// ClaimFees is a paid mutator transaction binding the contract method 0xd294f093.
//
// Solidity: function claimFees() returns()
func (_ClusterConnection *ClusterConnectionSession) ClaimFees() (*types.Transaction, error) {
	return _ClusterConnection.Contract.ClaimFees(&_ClusterConnection.TransactOpts)
}

// ClaimFees is a paid mutator transaction binding the contract method 0xd294f093.
//
// Solidity: function claimFees() returns()
func (_ClusterConnection *ClusterConnectionTransactorSession) ClaimFees() (*types.Transaction, error) {
	return _ClusterConnection.Contract.ClaimFees(&_ClusterConnection.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _relayer, address _xCall) returns()
func (_ClusterConnection *ClusterConnectionTransactor) Initialize(opts *bind.TransactOpts, _relayer common.Address, _xCall common.Address) (*types.Transaction, error) {
	return _ClusterConnection.contract.Transact(opts, "initialize", _relayer, _xCall)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _relayer, address _xCall) returns()
func (_ClusterConnection *ClusterConnectionSession) Initialize(_relayer common.Address, _xCall common.Address) (*types.Transaction, error) {
	return _ClusterConnection.Contract.Initialize(&_ClusterConnection.TransactOpts, _relayer, _xCall)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _relayer, address _xCall) returns()
func (_ClusterConnection *ClusterConnectionTransactorSession) Initialize(_relayer common.Address, _xCall common.Address) (*types.Transaction, error) {
	return _ClusterConnection.Contract.Initialize(&_ClusterConnection.TransactOpts, _relayer, _xCall)
}

// RecvMessage is a paid mutator transaction binding the contract method 0xb58b4cec.
//
// Solidity: function recvMessage(string srcNetwork, uint256 _connSn, bytes _msg) returns()
func (_ClusterConnection *ClusterConnectionTransactor) RecvMessage(opts *bind.TransactOpts, srcNetwork string, _connSn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _ClusterConnection.contract.Transact(opts, "recvMessage", srcNetwork, _connSn, _msg)
}

// RecvMessage is a paid mutator transaction binding the contract method 0xb58b4cec.
//
// Solidity: function recvMessage(string srcNetwork, uint256 _connSn, bytes _msg) returns()
func (_ClusterConnection *ClusterConnectionSession) RecvMessage(srcNetwork string, _connSn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _ClusterConnection.Contract.RecvMessage(&_ClusterConnection.TransactOpts, srcNetwork, _connSn, _msg)
}

// RecvMessage is a paid mutator transaction binding the contract method 0xb58b4cec.
//
// Solidity: function recvMessage(string srcNetwork, uint256 _connSn, bytes _msg) returns()
func (_ClusterConnection *ClusterConnectionTransactorSession) RecvMessage(srcNetwork string, _connSn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _ClusterConnection.Contract.RecvMessage(&_ClusterConnection.TransactOpts, srcNetwork, _connSn, _msg)
}

// RecvMessageWithSignatures is a paid mutator transaction binding the contract method 0x0a523d68.
//
// Solidity: function recvMessageWithSignatures(string srcNetwork, uint256 _connSn, bytes _msg, bytes[] _signedMessages) returns()
func (_ClusterConnection *ClusterConnectionTransactor) RecvMessageWithSignatures(opts *bind.TransactOpts, srcNetwork string, _connSn *big.Int, _msg []byte, _signedMessages [][]byte) (*types.Transaction, error) {
	return _ClusterConnection.contract.Transact(opts, "recvMessageWithSignatures", srcNetwork, _connSn, _msg, _signedMessages)
}

// RecvMessageWithSignatures is a paid mutator transaction binding the contract method 0x0a523d68.
//
// Solidity: function recvMessageWithSignatures(string srcNetwork, uint256 _connSn, bytes _msg, bytes[] _signedMessages) returns()
func (_ClusterConnection *ClusterConnectionSession) RecvMessageWithSignatures(srcNetwork string, _connSn *big.Int, _msg []byte, _signedMessages [][]byte) (*types.Transaction, error) {
	return _ClusterConnection.Contract.RecvMessageWithSignatures(&_ClusterConnection.TransactOpts, srcNetwork, _connSn, _msg, _signedMessages)
}

// RecvMessageWithSignatures is a paid mutator transaction binding the contract method 0x0a523d68.
//
// Solidity: function recvMessageWithSignatures(string srcNetwork, uint256 _connSn, bytes _msg, bytes[] _signedMessages) returns()
func (_ClusterConnection *ClusterConnectionTransactorSession) RecvMessageWithSignatures(srcNetwork string, _connSn *big.Int, _msg []byte, _signedMessages [][]byte) (*types.Transaction, error) {
	return _ClusterConnection.Contract.RecvMessageWithSignatures(&_ClusterConnection.TransactOpts, srcNetwork, _connSn, _msg, _signedMessages)
}

// RemoveSigner is a paid mutator transaction binding the contract method 0x0e316ab7.
//
// Solidity: function removeSigner(address _signer) returns()
func (_ClusterConnection *ClusterConnectionTransactor) RemoveSigner(opts *bind.TransactOpts, _signer common.Address) (*types.Transaction, error) {
	return _ClusterConnection.contract.Transact(opts, "removeSigner", _signer)
}

// RemoveSigner is a paid mutator transaction binding the contract method 0x0e316ab7.
//
// Solidity: function removeSigner(address _signer) returns()
func (_ClusterConnection *ClusterConnectionSession) RemoveSigner(_signer common.Address) (*types.Transaction, error) {
	return _ClusterConnection.Contract.RemoveSigner(&_ClusterConnection.TransactOpts, _signer)
}

// RemoveSigner is a paid mutator transaction binding the contract method 0x0e316ab7.
//
// Solidity: function removeSigner(address _signer) returns()
func (_ClusterConnection *ClusterConnectionTransactorSession) RemoveSigner(_signer common.Address) (*types.Transaction, error) {
	return _ClusterConnection.Contract.RemoveSigner(&_ClusterConnection.TransactOpts, _signer)
}

// RevertMessage is a paid mutator transaction binding the contract method 0x2d3fb823.
//
// Solidity: function revertMessage(uint256 sn) returns()
func (_ClusterConnection *ClusterConnectionTransactor) RevertMessage(opts *bind.TransactOpts, sn *big.Int) (*types.Transaction, error) {
	return _ClusterConnection.contract.Transact(opts, "revertMessage", sn)
}

// RevertMessage is a paid mutator transaction binding the contract method 0x2d3fb823.
//
// Solidity: function revertMessage(uint256 sn) returns()
func (_ClusterConnection *ClusterConnectionSession) RevertMessage(sn *big.Int) (*types.Transaction, error) {
	return _ClusterConnection.Contract.RevertMessage(&_ClusterConnection.TransactOpts, sn)
}

// RevertMessage is a paid mutator transaction binding the contract method 0x2d3fb823.
//
// Solidity: function revertMessage(uint256 sn) returns()
func (_ClusterConnection *ClusterConnectionTransactorSession) RevertMessage(sn *big.Int) (*types.Transaction, error) {
	return _ClusterConnection.Contract.RevertMessage(&_ClusterConnection.TransactOpts, sn)
}

// SendMessage is a paid mutator transaction binding the contract method 0x522a901e.
//
// Solidity: function sendMessage(string to, string svc, int256 sn, bytes _msg) payable returns()
func (_ClusterConnection *ClusterConnectionTransactor) SendMessage(opts *bind.TransactOpts, to string, svc string, sn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _ClusterConnection.contract.Transact(opts, "sendMessage", to, svc, sn, _msg)
}

// SendMessage is a paid mutator transaction binding the contract method 0x522a901e.
//
// Solidity: function sendMessage(string to, string svc, int256 sn, bytes _msg) payable returns()
func (_ClusterConnection *ClusterConnectionSession) SendMessage(to string, svc string, sn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _ClusterConnection.Contract.SendMessage(&_ClusterConnection.TransactOpts, to, svc, sn, _msg)
}

// SendMessage is a paid mutator transaction binding the contract method 0x522a901e.
//
// Solidity: function sendMessage(string to, string svc, int256 sn, bytes _msg) payable returns()
func (_ClusterConnection *ClusterConnectionTransactorSession) SendMessage(to string, svc string, sn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _ClusterConnection.Contract.SendMessage(&_ClusterConnection.TransactOpts, to, svc, sn, _msg)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x704b6c02.
//
// Solidity: function setAdmin(address _address) returns()
func (_ClusterConnection *ClusterConnectionTransactor) SetAdmin(opts *bind.TransactOpts, _address common.Address) (*types.Transaction, error) {
	return _ClusterConnection.contract.Transact(opts, "setAdmin", _address)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x704b6c02.
//
// Solidity: function setAdmin(address _address) returns()
func (_ClusterConnection *ClusterConnectionSession) SetAdmin(_address common.Address) (*types.Transaction, error) {
	return _ClusterConnection.Contract.SetAdmin(&_ClusterConnection.TransactOpts, _address)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x704b6c02.
//
// Solidity: function setAdmin(address _address) returns()
func (_ClusterConnection *ClusterConnectionTransactorSession) SetAdmin(_address common.Address) (*types.Transaction, error) {
	return _ClusterConnection.Contract.SetAdmin(&_ClusterConnection.TransactOpts, _address)
}

// SetFee is a paid mutator transaction binding the contract method 0x43f08a89.
//
// Solidity: function setFee(string networkId, uint256 messageFee, uint256 responseFee) returns()
func (_ClusterConnection *ClusterConnectionTransactor) SetFee(opts *bind.TransactOpts, networkId string, messageFee *big.Int, responseFee *big.Int) (*types.Transaction, error) {
	return _ClusterConnection.contract.Transact(opts, "setFee", networkId, messageFee, responseFee)
}

// SetFee is a paid mutator transaction binding the contract method 0x43f08a89.
//
// Solidity: function setFee(string networkId, uint256 messageFee, uint256 responseFee) returns()
func (_ClusterConnection *ClusterConnectionSession) SetFee(networkId string, messageFee *big.Int, responseFee *big.Int) (*types.Transaction, error) {
	return _ClusterConnection.Contract.SetFee(&_ClusterConnection.TransactOpts, networkId, messageFee, responseFee)
}

// SetFee is a paid mutator transaction binding the contract method 0x43f08a89.
//
// Solidity: function setFee(string networkId, uint256 messageFee, uint256 responseFee) returns()
func (_ClusterConnection *ClusterConnectionTransactorSession) SetFee(networkId string, messageFee *big.Int, responseFee *big.Int) (*types.Transaction, error) {
	return _ClusterConnection.Contract.SetFee(&_ClusterConnection.TransactOpts, networkId, messageFee, responseFee)
}

// ClusterConnectionInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ClusterConnection contract.
type ClusterConnectionInitializedIterator struct {
	Event *ClusterConnectionInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ClusterConnectionInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ClusterConnectionInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ClusterConnectionInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ClusterConnectionInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ClusterConnectionInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ClusterConnectionInitialized represents a Initialized event raised by the ClusterConnection contract.
type ClusterConnectionInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_ClusterConnection *ClusterConnectionFilterer) FilterInitialized(opts *bind.FilterOpts) (*ClusterConnectionInitializedIterator, error) {

	logs, sub, err := _ClusterConnection.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ClusterConnectionInitializedIterator{contract: _ClusterConnection.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_ClusterConnection *ClusterConnectionFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ClusterConnectionInitialized) (event.Subscription, error) {

	logs, sub, err := _ClusterConnection.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ClusterConnectionInitialized)
				if err := _ClusterConnection.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_ClusterConnection *ClusterConnectionFilterer) ParseInitialized(log types.Log) (*ClusterConnectionInitialized, error) {
	event := new(ClusterConnectionInitialized)
	if err := _ClusterConnection.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ClusterConnectionMessageIterator is returned from FilterMessage and is used to iterate over the raw logs and unpacked data for Message events raised by the ClusterConnection contract.
type ClusterConnectionMessageIterator struct {
	Event *ClusterConnectionMessage // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ClusterConnectionMessageIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ClusterConnectionMessage)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ClusterConnectionMessage)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ClusterConnectionMessageIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ClusterConnectionMessageIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ClusterConnectionMessage represents a Message event raised by the ClusterConnection contract.
type ClusterConnectionMessage struct {
	TargetNetwork string
	Sn            *big.Int
	Msg           []byte
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterMessage is a free log retrieval operation binding the contract event 0x37be353f216cf7e33639101fd610c542e6a0c0109173fa1c1d8b04d34edb7c1b.
//
// Solidity: event Message(string targetNetwork, uint256 sn, bytes _msg)
func (_ClusterConnection *ClusterConnectionFilterer) FilterMessage(opts *bind.FilterOpts) (*ClusterConnectionMessageIterator, error) {

	logs, sub, err := _ClusterConnection.contract.FilterLogs(opts, "Message")
	if err != nil {
		return nil, err
	}
	return &ClusterConnectionMessageIterator{contract: _ClusterConnection.contract, event: "Message", logs: logs, sub: sub}, nil
}

// WatchMessage is a free log subscription operation binding the contract event 0x37be353f216cf7e33639101fd610c542e6a0c0109173fa1c1d8b04d34edb7c1b.
//
// Solidity: event Message(string targetNetwork, uint256 sn, bytes _msg)
func (_ClusterConnection *ClusterConnectionFilterer) WatchMessage(opts *bind.WatchOpts, sink chan<- *ClusterConnectionMessage) (event.Subscription, error) {

	logs, sub, err := _ClusterConnection.contract.WatchLogs(opts, "Message")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ClusterConnectionMessage)
				if err := _ClusterConnection.contract.UnpackLog(event, "Message", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseMessage is a log parse operation binding the contract event 0x37be353f216cf7e33639101fd610c542e6a0c0109173fa1c1d8b04d34edb7c1b.
//
// Solidity: event Message(string targetNetwork, uint256 sn, bytes _msg)
func (_ClusterConnection *ClusterConnectionFilterer) ParseMessage(log types.Log) (*ClusterConnectionMessage, error) {
	event := new(ClusterConnectionMessage)
	if err := _ClusterConnection.contract.UnpackLog(event, "Message", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
