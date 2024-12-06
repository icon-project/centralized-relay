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
	ABI: "[{\"type\":\"function\",\"name\":\"admin\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"claimFees\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"connSn\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getFee\",\"inputs\":[{\"name\":\"to\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"response\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[{\"name\":\"fee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getReceipt\",\"inputs\":[{\"name\":\"srcNetwork\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_connSn\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRequiredValidatorCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_relayer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_xCall\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isValidator\",\"inputs\":[{\"name\":\"signer\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isValidatorProcessed\",\"inputs\":[{\"name\":\"processedSigners\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"signer\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"listValidators\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"recoverSigner\",\"inputs\":[{\"name\":\"messageHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"recvMessageWithSignatures\",\"inputs\":[{\"name\":\"srcNetwork\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_connSn\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_msg\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"dstNetwork\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_signedMessages\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"relayer\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"revertMessage\",\"inputs\":[{\"name\":\"sn\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"sendMessage\",\"inputs\":[{\"name\":\"to\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_svc\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"sn\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"_msg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"setAdmin\",\"inputs\":[{\"name\":\"_address\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setFee\",\"inputs\":[{\"name\":\"networkId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"messageFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"responseFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRelayer\",\"inputs\":[{\"name\":\"_address\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRequiredValidatorCount\",\"inputs\":[{\"name\":\"_count\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateValidators\",\"inputs\":[{\"name\":\"_validators\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"_threshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Message\",\"inputs\":[{\"name\":\"targetNetwork\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"sn\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"_msg\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorSetAdded\",\"inputs\":[{\"name\":\"_validator\",\"type\":\"bytes[]\",\"indexed\":false,\"internalType\":\"bytes[]\"},{\"name\":\"_threshold\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]}]",
	Bin: "0x608060405234801561001057600080fd5b506124d6806100206000396000f3fe60806040526004361061011f5760003560e01c80638406c079116100a0578063bbf8dfc311610064578063bbf8dfc314610328578063cbd856fb14610348578063d294f09314610368578063f851a4401461037d578063facd743b1461039b57600080fd5b80638406c079146102705780639664da0e146102a257806397aba7f9146102d257806399f1fca7146102f2578063a914a3521461030857600080fd5b806362534d93116100e757806362534d93146101b95780636548e9bc146101e057806368d4e54414610200578063704b6c02146102225780637d4c4f4a1461024257600080fd5b80632d3fb82314610124578063308b74df1461014657806343f08a8914610166578063485cc95514610186578063522a901e146101a6575b600080fd5b34801561013057600080fd5b5061014461013f366004611989565b6103bb565b005b34801561015257600080fd5b506101446101613660046119b8565b6104b7565b34801561017257600080fd5b50610144610181366004611a22565b61055f565b34801561019257600080fd5b506101446101a1366004611a92565b61063f565b6101446101b4366004611acb565b61077f565b3480156101c557600080fd5b5060085460405160ff90911681526020015b60405180910390f35b3480156101ec57600080fd5b506101446101fb366004611b6e565b61096f565b34801561020c57600080fd5b50610215610a23565b6040516101d79190611b8b565b34801561022e57600080fd5b5061014461023d366004611b6e565b610a85565b34801561024e57600080fd5b5061026261025d366004611c8d565b610b39565b6040519081526020016101d7565b34801561027c57600080fd5b506004546001600160a01b03165b6040516001600160a01b0390911681526020016101d7565b3480156102ae57600080fd5b506102c26102bd366004611cd8565b610ba8565b60405190151581526020016101d7565b3480156102de57600080fd5b5061028a6102ed366004611d1c565b610be0565b3480156102fe57600080fd5b5061026260065481565b34801561031457600080fd5b506102c2610323366004611d85565b610d25565b34801561033457600080fd5b50610144610343366004611e35565b610d80565b34801561035457600080fd5b50610144610363366004611f27565b611055565b34801561037457600080fd5b50610144611238565b34801561038957600080fd5b506005546001600160a01b031661028a565b3480156103a757600080fd5b506102c26103b6366004611b6e565b611306565b306001600160a01b0316638406c0796040518163ffffffff1660e01b8152600401602060405180830381865afa1580156103f9573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061041d9190611fe9565b6001600160a01b0316336001600160a01b0316146104565760405162461bcd60e51b815260040161044d90612006565b60405180910390fd5b60035460405163b070f9e560e01b8152600481018390526001600160a01b039091169063b070f9e590602401600060405180830381600087803b15801561049c57600080fd5b505af11580156104b0573d6000803e3d6000fd5b5050505050565b306001600160a01b031663f851a4406040518163ffffffff1660e01b8152600401602060405180830381865afa1580156104f5573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105199190611fe9565b6001600160a01b0316336001600160a01b0316146105495760405162461bcd60e51b815260040161044d9061202b565b6008805460ff191660ff92909216919091179055565b306001600160a01b0316638406c0796040518163ffffffff1660e01b8152600401602060405180830381865afa15801561059d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105c19190611fe9565b6001600160a01b0316336001600160a01b0316146105f15760405162461bcd60e51b815260040161044d90612006565b816000858560405161060492919061204e565b908152602001604051809103902081905550806001858560405161062992919061204e565b9081526040519081900360200190205550505050565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a008054600160401b810460ff1615906001600160401b03166000811580156106845750825b90506000826001600160401b031660011480156106a05750303b155b9050811580156106ae575080155b156106cc5760405163f92ee8a960e01b815260040160405180910390fd5b845467ffffffffffffffff1916600117855583156106f657845460ff60401b1916600160401b1785555b600380546001600160a01b038089166001600160a01b0319928316179092556005805482163317905560048054928a1692909116919091179055831561077657845460ff60401b19168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b50505050505050565b6003546001600160a01b031633146107d95760405162461bcd60e51b815260206004820152601f60248201527f4f6e6c79205863616c6c2063616e2063616c6c2073656e644d65737361676500604482015260640161044d565b60008084131561085357604051633ea627a560e11b81523090637d4c4f4a9061080b908b908b90600190600401612087565b602060405180830381865afa158015610828573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061084c91906120ad565b90506108c7565b836000036108c757604051633ea627a560e11b81523090637d4c4f4a90610883908b908b90600090600401612087565b602060405180830381865afa1580156108a0573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108c491906120ad565b90505b8034101561090f5760405162461bcd60e51b8152602060048201526015602482015274119959481a5cc81b9bdd0814dd59999a58da595b9d605a1b604482015260640161044d565b6006805490600061091f836120dc565b91905055507f37be353f216cf7e33639101fd610c542e6a0c0109173fa1c1d8b04d34edb7c1b8888600654868660405161095d9594939291906120f5565b60405180910390a15050505050505050565b306001600160a01b031663f851a4406040518163ffffffff1660e01b8152600401602060405180830381865afa1580156109ad573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906109d19190611fe9565b6001600160a01b0316336001600160a01b031614610a015760405162461bcd60e51b815260040161044d9061202b565b600480546001600160a01b0319166001600160a01b0392909216919091179055565b60606007805480602002602001604051908101604052809291908181526020018280548015610a7b57602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311610a5d575b5050505050905090565b306001600160a01b031663f851a4406040518163ffffffff1660e01b8152600401602060405180830381865afa158015610ac3573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ae79190611fe9565b6001600160a01b0316336001600160a01b031614610b175760405162461bcd60e51b815260040161044d9061202b565b600580546001600160a01b0319166001600160a01b0392909216919091179055565b600080600084604051610b4c9190612152565b908152604051908190036020019020549050821515600103610b9f576000600185604051610b7a9190612152565b908152604051908190036020019020549050610b96818361216e565b92505050610ba2565b90505b92915050565b6000600283604051610bba9190612152565b908152604080516020928190038301902060009485529091529091205460ff1692915050565b60008151604114610c335760405162461bcd60e51b815260206004820152601860248201527f496e76616c6964207369676e6174757265206c656e6774680000000000000000604482015260640161044d565b60208201516040830151606084015160001a601b811015610c5c57610c59601b82612181565b90505b8060ff16601b1480610c7157508060ff16601c145b610cbd5760405162461bcd60e51b815260206004820152601b60248201527f496e76616c6964207369676e6174757265202776272076616c75650000000000604482015260640161044d565b60408051600081526020810180835288905260ff831691810191909152606081018490526080810183905260019060a0016020604051602081039080840390855afa158015610d10573d6000803e3d6000fd5b5050604051601f190151979650505050505050565b6000805b8351811015610d7657826001600160a01b0316848281518110610d4e57610d4e61219a565b60200260200101516001600160a01b031603610d6e576001915050610ba2565b600101610d29565b5060009392505050565b306001600160a01b0316638406c0796040518163ffffffff1660e01b8152600401602060405180830381865afa158015610dbe573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610de29190611fe9565b6001600160a01b0316336001600160a01b031614610e125760405162461bcd60e51b815260040161044d90612006565b60085460ff16811015610e675760405162461bcd60e51b815260206004820152601c60248201527f4e6f7420656e6f756768207369676e6174757265732070617373656400000000604482015260640161044d565b6000610e768888888888611362565b9050600080836001600160401b03811115610e9357610e93611bd8565b604051908082528060200260200182016040528015610ebc578160200160208202803683370190505b50905060005b84811015610fdc576000610f2e85888885818110610ee257610ee261219a565b9050602002810190610ef491906121b0565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250610be092505050565b90506001600160a01b038116610f7a5760405162461bcd60e51b8152602060048201526011602482015270496e76616c6964207369676e617475726560781b604482015260640161044d565b610f848382610d25565b158015610f955750610f95816113fc565b15610fd35780838581518110610fad57610fad61219a565b6001600160a01b039092166020928302919091019091015283610fcf816120dc565b9450505b50600101610ec2565b5060085460ff1682101561103d5760405162461bcd60e51b815260206004820152602260248201527f4e6f7420656e6f7567682076616c6964207369676e6174757265732070617373604482015261195960f21b606482015260840161044d565b6110498a8a8a8a61145b565b50505050505050505050565b306001600160a01b031663f851a4406040518163ffffffff1660e01b8152600401602060405180830381865afa158015611093573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906110b79190611fe9565b6001600160a01b0316336001600160a01b0316146110e75760405162461bcd60e51b815260040161044d9061202b565b6110f360076000611957565b60005b825181101561119d5760006111238483815181106111165761111661219a565b602002602001015161156d565b905061112e81611306565b15801561114357506001600160a01b03811615155b1561119457600780546001810182556000919091527fa66cc928b5edb82af9bd49922954155ab7b0942694bea4ce44661d9a8736c6880180546001600160a01b0319166001600160a01b0383161790555b506001016110f6565b5060075460ff821611156111eb5760405162461bcd60e51b81526020600482015260156024820152744e6f7420656e6f7567682076616c696461746f727360581b604482015260640161044d565b6008805460ff191660ff83161790556040517fee022091c8c62247fb714803aafa9265a5b04fdd39f40668cc38a4170b21a2909061122c9084908490612222565b60405180910390a15050565b306001600160a01b0316638406c0796040518163ffffffff1660e01b8152600401602060405180830381865afa158015611276573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061129a9190611fe9565b6001600160a01b0316336001600160a01b0316146112ca5760405162461bcd60e51b815260040161044d90612006565b6004546040516001600160a01b03909116904780156108fc02916000818181858888f19350505050158015611303573d6000803e3d6000fd5b50565b6000805b60075481101561135c57826001600160a01b0316600782815481106113315761133161219a565b6000918252602090912001546001600160a01b0316036113545750600192915050565b60010161130a565b50919050565b6000806113e9611371886115cc565b61137a886115d7565b6113b988888080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506115e592505050565b6113c28b6115cc565b6040516020016113d59493929190612292565b604051602081830303815290604052611654565b8051602090910120979650505050505050565b6000805b60075481101561145257826001600160a01b0316600782815481106114275761142761219a565b6000918252602090912001546001600160a01b03160361144a5750600192915050565b600101611400565b50600092915050565b60028460405161146b9190612152565b90815260408051602092819003830190206000868152925290205460ff16156114ca5760405162461bcd60e51b81526020600482015260116024820152704475706c6963617465204d65737361676560781b604482015260640161044d565b60016002856040516114dc9190612152565b9081526040805160209281900383018120600088815293529120805460ff19169215159290921790915560035463bbc22efd60e01b82526001600160a01b03169063bbc22efd90611535908790869086906004016122e9565b600060405180830381600087803b15801561154f57600080fd5b505af1158015611563573d6000803e3d6000fd5b5050505050505050565b600081516041146115c05760405162461bcd60e51b815260206004820152601960248201527f496e76616c6964207075626c6963206b6579206c656e67746800000000000000604482015260640161044d565b50604060219091012090565b6060610ba2826115e5565b6060610ba26115e58361168a565b6060808251600114801561161357506080836000815181106116095761160961219a565b016020015160f81c105b1561161f575081610ba2565b61162b83516080611734565b8360405160200161163d929190612319565b604051602081830303815290604052905092915050565b6060611662825160c0611734565b82604051602001611674929190612319565b6040516020818303038152906040529050919050565b6060816000036116b757604080516001808252818301909252906020820181803683370190505092915050565b608060015b60208110156116ee57818410156116df576116d784826118eb565b949350505050565b60089190911b906001016116bc565b508083101561171b576040805160208101859052015b604051602081830303815290604052915050919050565b6040516000602082015260218101849052604101611704565b606080603884101561179e5760408051600180825281830190925290602082018180368337019050509050611769838561216e565b601f1a60f81b816000815181106117825761178261219a565b60200101906001600160f81b031916908160001a905350610b9f565b600060015b6117ad818761235e565b156117d357816117bc816120dc565b92506117cc905061010082612372565b90506117a3565b6117de82600161216e565b6001600160401b038111156117f5576117f5611bd8565b6040519080825280601f01601f19166020018201604052801561181f576020820181803683370190505b50925061182c858361216e565b61183790603761216e565b601f1a60f81b836000815181106118505761185061219a565b60200101906001600160f81b031916908160001a905350600190505b8181116118e1576101006118808284612389565b61188c90610100612480565b611896908861235e565b6118a0919061248c565b601f1a60f81b8382815181106118b8576118b861219a565b60200101906001600160f81b031916908160001a905350806118d9816120dc565b91505061186c565b5050905092915050565b60606000826001600160401b0381111561190757611907611bd8565b6040519080825280601f01601f191660200182016040528015611931576020820181803683370190505b50905060208101836020035b60208110156118e15785811a82536001918201910161193d565b508054600082559060005260206000209081019061130391905b808211156119855760008155600101611971565b5090565b60006020828403121561199b57600080fd5b5035919050565b803560ff811681146119b357600080fd5b919050565b6000602082840312156119ca57600080fd5b6119d3826119a2565b9392505050565b60008083601f8401126119ec57600080fd5b5081356001600160401b03811115611a0357600080fd5b602083019150836020828501011115611a1b57600080fd5b9250929050565b60008060008060608587031215611a3857600080fd5b84356001600160401b03811115611a4e57600080fd5b611a5a878288016119da565b90989097506020870135966040013595509350505050565b6001600160a01b038116811461130357600080fd5b80356119b381611a72565b60008060408385031215611aa557600080fd5b8235611ab081611a72565b91506020830135611ac081611a72565b809150509250929050565b60008060008060008060006080888a031215611ae657600080fd5b87356001600160401b0380821115611afd57600080fd5b611b098b838c016119da565b909950975060208a0135915080821115611b2257600080fd5b611b2e8b838c016119da565b909750955060408a0135945060608a0135915080821115611b4e57600080fd5b50611b5b8a828b016119da565b989b979a50959850939692959293505050565b600060208284031215611b8057600080fd5b8135610b9f81611a72565b6020808252825182820181905260009190848201906040850190845b81811015611bcc5783516001600160a01b031683529284019291840191600101611ba7565b50909695505050505050565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f191681016001600160401b0381118282101715611c1657611c16611bd8565b604052919050565b600082601f830112611c2f57600080fd5b81356001600160401b03811115611c4857611c48611bd8565b611c5b601f8201601f1916602001611bee565b818152846020838601011115611c7057600080fd5b816020850160208301376000918101602001919091529392505050565b60008060408385031215611ca057600080fd5b82356001600160401b03811115611cb657600080fd5b611cc285828601611c1e565b92505060208301358015158114611ac057600080fd5b60008060408385031215611ceb57600080fd5b82356001600160401b03811115611d0157600080fd5b611d0d85828601611c1e565b95602094909401359450505050565b60008060408385031215611d2f57600080fd5b8235915060208301356001600160401b03811115611d4c57600080fd5b611d5885828601611c1e565b9150509250929050565b60006001600160401b03821115611d7b57611d7b611bd8565b5060051b60200190565b60008060408385031215611d9857600080fd5b82356001600160401b03811115611dae57600080fd5b8301601f81018513611dbf57600080fd5b80356020611dd4611dcf83611d62565b611bee565b82815260059290921b83018101918181019088841115611df357600080fd5b938201935b83851015611e1a578435611e0b81611a72565b82529382019390820190611df8565b9550611e299050868201611a87565b93505050509250929050565b600080600080600080600060a0888a031215611e5057600080fd5b87356001600160401b0380821115611e6757600080fd5b611e738b838c01611c1e565b985060208a0135975060408a0135915080821115611e9057600080fd5b611e9c8b838c016119da565b909750955060608a0135915080821115611eb557600080fd5b611ec18b838c01611c1e565b945060808a0135915080821115611ed757600080fd5b818a0191508a601f830112611eeb57600080fd5b813581811115611efa57600080fd5b8b60208260051b8501011115611f0f57600080fd5b60208301945080935050505092959891949750929550565b60008060408385031215611f3a57600080fd5b82356001600160401b0380821115611f5157600080fd5b818501915085601f830112611f6557600080fd5b81356020611f75611dcf83611d62565b82815260059290921b84018101918181019089841115611f9457600080fd5b8286015b84811015611fcc57803586811115611fb05760008081fd5b611fbe8c86838b0101611c1e565b845250918301918301611f98565b509650611fdc90508782016119a2565b9450505050509250929050565b600060208284031215611ffb57600080fd5b8151610b9f81611a72565b6020808252600b908201526a27b7363ca932b630bcb2b960a91b604082015260600190565b60208082526009908201526827b7363ca0b236b4b760b91b604082015260600190565b8183823760009101908152919050565b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b60408152600061209b60408301858761205e565b90508215156020830152949350505050565b6000602082840312156120bf57600080fd5b5051919050565b634e487b7160e01b600052601160045260246000fd5b6000600182016120ee576120ee6120c6565b5060010190565b60608152600061210960608301878961205e565b856020840152828103604084015261212281858761205e565b98975050505050505050565b60005b83811015612149578181015183820152602001612131565b50506000910152565b6000825161216481846020870161212e565b9190910192915050565b80820180821115610ba257610ba26120c6565b60ff8181168382160190811115610ba257610ba26120c6565b634e487b7160e01b600052603260045260246000fd5b6000808335601e198436030181126121c757600080fd5b8301803591506001600160401b038211156121e157600080fd5b602001915036819003821315611a1b57600080fd5b6000815180845261220e81602086016020860161212e565b601f01601f19169290920160200192915050565b6000604082016040835280855180835260608501915060608160051b8601019250602080880160005b8381101561227957605f198887030185526122678683516121f6565b9550938201939082019060010161224b565b50505050508091505060ff831660208301529392505050565b600085516122a4818460208a0161212e565b8551908301906122b8818360208a0161212e565b85519101906122cb81836020890161212e565b84519101906122de81836020880161212e565b019695505050505050565b6040815260006122fc60408301866121f6565b828103602084015261230f81858761205e565b9695505050505050565b6000835161232b81846020880161212e565b83519083019061233f81836020880161212e565b01949350505050565b634e487b7160e01b600052601260045260246000fd5b60008261236d5761236d612348565b500490565b8082028115828204841417610ba257610ba26120c6565b81810381811115610ba257610ba26120c6565b600181815b808511156123d75781600019048211156123bd576123bd6120c6565b808516156123ca57918102915b93841c93908002906123a1565b509250929050565b6000826123ee57506001610ba2565b816123fb57506000610ba2565b8160018114612411576002811461241b57612437565b6001915050610ba2565b60ff84111561242c5761242c6120c6565b50506001821b610ba2565b5060208310610133831016604e8410600b841016171561245a575081810a610ba2565b612464838361239c565b8060001904821115612478576124786120c6565b029392505050565b60006119d383836123df565b60008261249b5761249b612348565b50069056fea26469706673582212201c02b614b483cbee9ff08d8199d28a616e755281b7e28bf99532f312eb1ee56e64736f6c63430008180033",
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

// GetRequiredValidatorCount is a free data retrieval call binding the contract method 0x62534d93.
//
// Solidity: function getRequiredValidatorCount() view returns(uint8)
func (_ClusterConnection *ClusterConnectionCaller) GetRequiredValidatorCount(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _ClusterConnection.contract.Call(opts, &out, "getRequiredValidatorCount")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetRequiredValidatorCount is a free data retrieval call binding the contract method 0x62534d93.
//
// Solidity: function getRequiredValidatorCount() view returns(uint8)
func (_ClusterConnection *ClusterConnectionSession) GetRequiredValidatorCount() (uint8, error) {
	return _ClusterConnection.Contract.GetRequiredValidatorCount(&_ClusterConnection.CallOpts)
}

// GetRequiredValidatorCount is a free data retrieval call binding the contract method 0x62534d93.
//
// Solidity: function getRequiredValidatorCount() view returns(uint8)
func (_ClusterConnection *ClusterConnectionCallerSession) GetRequiredValidatorCount() (uint8, error) {
	return _ClusterConnection.Contract.GetRequiredValidatorCount(&_ClusterConnection.CallOpts)
}

// IsValidator is a free data retrieval call binding the contract method 0xfacd743b.
//
// Solidity: function isValidator(address signer) view returns(bool)
func (_ClusterConnection *ClusterConnectionCaller) IsValidator(opts *bind.CallOpts, signer common.Address) (bool, error) {
	var out []interface{}
	err := _ClusterConnection.contract.Call(opts, &out, "isValidator", signer)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsValidator is a free data retrieval call binding the contract method 0xfacd743b.
//
// Solidity: function isValidator(address signer) view returns(bool)
func (_ClusterConnection *ClusterConnectionSession) IsValidator(signer common.Address) (bool, error) {
	return _ClusterConnection.Contract.IsValidator(&_ClusterConnection.CallOpts, signer)
}

// IsValidator is a free data retrieval call binding the contract method 0xfacd743b.
//
// Solidity: function isValidator(address signer) view returns(bool)
func (_ClusterConnection *ClusterConnectionCallerSession) IsValidator(signer common.Address) (bool, error) {
	return _ClusterConnection.Contract.IsValidator(&_ClusterConnection.CallOpts, signer)
}

// IsValidatorProcessed is a free data retrieval call binding the contract method 0xa914a352.
//
// Solidity: function isValidatorProcessed(address[] processedSigners, address signer) pure returns(bool)
func (_ClusterConnection *ClusterConnectionCaller) IsValidatorProcessed(opts *bind.CallOpts, processedSigners []common.Address, signer common.Address) (bool, error) {
	var out []interface{}
	err := _ClusterConnection.contract.Call(opts, &out, "isValidatorProcessed", processedSigners, signer)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsValidatorProcessed is a free data retrieval call binding the contract method 0xa914a352.
//
// Solidity: function isValidatorProcessed(address[] processedSigners, address signer) pure returns(bool)
func (_ClusterConnection *ClusterConnectionSession) IsValidatorProcessed(processedSigners []common.Address, signer common.Address) (bool, error) {
	return _ClusterConnection.Contract.IsValidatorProcessed(&_ClusterConnection.CallOpts, processedSigners, signer)
}

// IsValidatorProcessed is a free data retrieval call binding the contract method 0xa914a352.
//
// Solidity: function isValidatorProcessed(address[] processedSigners, address signer) pure returns(bool)
func (_ClusterConnection *ClusterConnectionCallerSession) IsValidatorProcessed(processedSigners []common.Address, signer common.Address) (bool, error) {
	return _ClusterConnection.Contract.IsValidatorProcessed(&_ClusterConnection.CallOpts, processedSigners, signer)
}

// ListValidators is a free data retrieval call binding the contract method 0x68d4e544.
//
// Solidity: function listValidators() view returns(address[])
func (_ClusterConnection *ClusterConnectionCaller) ListValidators(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _ClusterConnection.contract.Call(opts, &out, "listValidators")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// ListValidators is a free data retrieval call binding the contract method 0x68d4e544.
//
// Solidity: function listValidators() view returns(address[])
func (_ClusterConnection *ClusterConnectionSession) ListValidators() ([]common.Address, error) {
	return _ClusterConnection.Contract.ListValidators(&_ClusterConnection.CallOpts)
}

// ListValidators is a free data retrieval call binding the contract method 0x68d4e544.
//
// Solidity: function listValidators() view returns(address[])
func (_ClusterConnection *ClusterConnectionCallerSession) ListValidators() ([]common.Address, error) {
	return _ClusterConnection.Contract.ListValidators(&_ClusterConnection.CallOpts)
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

// Relayer is a free data retrieval call binding the contract method 0x8406c079.
//
// Solidity: function relayer() view returns(address)
func (_ClusterConnection *ClusterConnectionCaller) Relayer(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ClusterConnection.contract.Call(opts, &out, "relayer")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Relayer is a free data retrieval call binding the contract method 0x8406c079.
//
// Solidity: function relayer() view returns(address)
func (_ClusterConnection *ClusterConnectionSession) Relayer() (common.Address, error) {
	return _ClusterConnection.Contract.Relayer(&_ClusterConnection.CallOpts)
}

// Relayer is a free data retrieval call binding the contract method 0x8406c079.
//
// Solidity: function relayer() view returns(address)
func (_ClusterConnection *ClusterConnectionCallerSession) Relayer() (common.Address, error) {
	return _ClusterConnection.Contract.Relayer(&_ClusterConnection.CallOpts)
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

// RecvMessageWithSignatures is a paid mutator transaction binding the contract method 0xbbf8dfc3.
//
// Solidity: function recvMessageWithSignatures(string srcNetwork, uint256 _connSn, bytes _msg, string dstNetwork, bytes[] _signedMessages) returns()
func (_ClusterConnection *ClusterConnectionTransactor) RecvMessageWithSignatures(opts *bind.TransactOpts, srcNetwork string, _connSn *big.Int, _msg []byte, dstNetwork string, _signedMessages [][]byte) (*types.Transaction, error) {
	return _ClusterConnection.contract.Transact(opts, "recvMessageWithSignatures", srcNetwork, _connSn, _msg, dstNetwork, _signedMessages)
}

// RecvMessageWithSignatures is a paid mutator transaction binding the contract method 0xbbf8dfc3.
//
// Solidity: function recvMessageWithSignatures(string srcNetwork, uint256 _connSn, bytes _msg, string dstNetwork, bytes[] _signedMessages) returns()
func (_ClusterConnection *ClusterConnectionSession) RecvMessageWithSignatures(srcNetwork string, _connSn *big.Int, _msg []byte, dstNetwork string, _signedMessages [][]byte) (*types.Transaction, error) {
	return _ClusterConnection.Contract.RecvMessageWithSignatures(&_ClusterConnection.TransactOpts, srcNetwork, _connSn, _msg, dstNetwork, _signedMessages)
}

// RecvMessageWithSignatures is a paid mutator transaction binding the contract method 0xbbf8dfc3.
//
// Solidity: function recvMessageWithSignatures(string srcNetwork, uint256 _connSn, bytes _msg, string dstNetwork, bytes[] _signedMessages) returns()
func (_ClusterConnection *ClusterConnectionTransactorSession) RecvMessageWithSignatures(srcNetwork string, _connSn *big.Int, _msg []byte, dstNetwork string, _signedMessages [][]byte) (*types.Transaction, error) {
	return _ClusterConnection.Contract.RecvMessageWithSignatures(&_ClusterConnection.TransactOpts, srcNetwork, _connSn, _msg, dstNetwork, _signedMessages)
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
// Solidity: function sendMessage(string to, string _svc, int256 sn, bytes _msg) payable returns()
func (_ClusterConnection *ClusterConnectionTransactor) SendMessage(opts *bind.TransactOpts, to string, _svc string, sn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _ClusterConnection.contract.Transact(opts, "sendMessage", to, _svc, sn, _msg)
}

// SendMessage is a paid mutator transaction binding the contract method 0x522a901e.
//
// Solidity: function sendMessage(string to, string _svc, int256 sn, bytes _msg) payable returns()
func (_ClusterConnection *ClusterConnectionSession) SendMessage(to string, _svc string, sn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _ClusterConnection.Contract.SendMessage(&_ClusterConnection.TransactOpts, to, _svc, sn, _msg)
}

// SendMessage is a paid mutator transaction binding the contract method 0x522a901e.
//
// Solidity: function sendMessage(string to, string _svc, int256 sn, bytes _msg) payable returns()
func (_ClusterConnection *ClusterConnectionTransactorSession) SendMessage(to string, _svc string, sn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _ClusterConnection.Contract.SendMessage(&_ClusterConnection.TransactOpts, to, _svc, sn, _msg)
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

// SetRelayer is a paid mutator transaction binding the contract method 0x6548e9bc.
//
// Solidity: function setRelayer(address _address) returns()
func (_ClusterConnection *ClusterConnectionTransactor) SetRelayer(opts *bind.TransactOpts, _address common.Address) (*types.Transaction, error) {
	return _ClusterConnection.contract.Transact(opts, "setRelayer", _address)
}

// SetRelayer is a paid mutator transaction binding the contract method 0x6548e9bc.
//
// Solidity: function setRelayer(address _address) returns()
func (_ClusterConnection *ClusterConnectionSession) SetRelayer(_address common.Address) (*types.Transaction, error) {
	return _ClusterConnection.Contract.SetRelayer(&_ClusterConnection.TransactOpts, _address)
}

// SetRelayer is a paid mutator transaction binding the contract method 0x6548e9bc.
//
// Solidity: function setRelayer(address _address) returns()
func (_ClusterConnection *ClusterConnectionTransactorSession) SetRelayer(_address common.Address) (*types.Transaction, error) {
	return _ClusterConnection.Contract.SetRelayer(&_ClusterConnection.TransactOpts, _address)
}

// SetRequiredValidatorCount is a paid mutator transaction binding the contract method 0x308b74df.
//
// Solidity: function setRequiredValidatorCount(uint8 _count) returns()
func (_ClusterConnection *ClusterConnectionTransactor) SetRequiredValidatorCount(opts *bind.TransactOpts, _count uint8) (*types.Transaction, error) {
	return _ClusterConnection.contract.Transact(opts, "setRequiredValidatorCount", _count)
}

// SetRequiredValidatorCount is a paid mutator transaction binding the contract method 0x308b74df.
//
// Solidity: function setRequiredValidatorCount(uint8 _count) returns()
func (_ClusterConnection *ClusterConnectionSession) SetRequiredValidatorCount(_count uint8) (*types.Transaction, error) {
	return _ClusterConnection.Contract.SetRequiredValidatorCount(&_ClusterConnection.TransactOpts, _count)
}

// SetRequiredValidatorCount is a paid mutator transaction binding the contract method 0x308b74df.
//
// Solidity: function setRequiredValidatorCount(uint8 _count) returns()
func (_ClusterConnection *ClusterConnectionTransactorSession) SetRequiredValidatorCount(_count uint8) (*types.Transaction, error) {
	return _ClusterConnection.Contract.SetRequiredValidatorCount(&_ClusterConnection.TransactOpts, _count)
}

// UpdateValidators is a paid mutator transaction binding the contract method 0xcbd856fb.
//
// Solidity: function updateValidators(bytes[] _validators, uint8 _threshold) returns()
func (_ClusterConnection *ClusterConnectionTransactor) UpdateValidators(opts *bind.TransactOpts, _validators [][]byte, _threshold uint8) (*types.Transaction, error) {
	return _ClusterConnection.contract.Transact(opts, "updateValidators", _validators, _threshold)
}

// UpdateValidators is a paid mutator transaction binding the contract method 0xcbd856fb.
//
// Solidity: function updateValidators(bytes[] _validators, uint8 _threshold) returns()
func (_ClusterConnection *ClusterConnectionSession) UpdateValidators(_validators [][]byte, _threshold uint8) (*types.Transaction, error) {
	return _ClusterConnection.Contract.UpdateValidators(&_ClusterConnection.TransactOpts, _validators, _threshold)
}

// UpdateValidators is a paid mutator transaction binding the contract method 0xcbd856fb.
//
// Solidity: function updateValidators(bytes[] _validators, uint8 _threshold) returns()
func (_ClusterConnection *ClusterConnectionTransactorSession) UpdateValidators(_validators [][]byte, _threshold uint8) (*types.Transaction, error) {
	return _ClusterConnection.Contract.UpdateValidators(&_ClusterConnection.TransactOpts, _validators, _threshold)
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

// ClusterConnectionValidatorSetAddedIterator is returned from FilterValidatorSetAdded and is used to iterate over the raw logs and unpacked data for ValidatorSetAdded events raised by the ClusterConnection contract.
type ClusterConnectionValidatorSetAddedIterator struct {
	Event *ClusterConnectionValidatorSetAdded // Event containing the contract specifics and raw log

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
func (it *ClusterConnectionValidatorSetAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ClusterConnectionValidatorSetAdded)
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
		it.Event = new(ClusterConnectionValidatorSetAdded)
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
func (it *ClusterConnectionValidatorSetAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ClusterConnectionValidatorSetAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ClusterConnectionValidatorSetAdded represents a ValidatorSetAdded event raised by the ClusterConnection contract.
type ClusterConnectionValidatorSetAdded struct {
	Validator [][]byte
	Threshold uint8
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterValidatorSetAdded is a free log retrieval operation binding the contract event 0xee022091c8c62247fb714803aafa9265a5b04fdd39f40668cc38a4170b21a290.
//
// Solidity: event ValidatorSetAdded(bytes[] _validator, uint8 _threshold)
func (_ClusterConnection *ClusterConnectionFilterer) FilterValidatorSetAdded(opts *bind.FilterOpts) (*ClusterConnectionValidatorSetAddedIterator, error) {

	logs, sub, err := _ClusterConnection.contract.FilterLogs(opts, "ValidatorSetAdded")
	if err != nil {
		return nil, err
	}
	return &ClusterConnectionValidatorSetAddedIterator{contract: _ClusterConnection.contract, event: "ValidatorSetAdded", logs: logs, sub: sub}, nil
}

// WatchValidatorSetAdded is a free log subscription operation binding the contract event 0xee022091c8c62247fb714803aafa9265a5b04fdd39f40668cc38a4170b21a290.
//
// Solidity: event ValidatorSetAdded(bytes[] _validator, uint8 _threshold)
func (_ClusterConnection *ClusterConnectionFilterer) WatchValidatorSetAdded(opts *bind.WatchOpts, sink chan<- *ClusterConnectionValidatorSetAdded) (event.Subscription, error) {

	logs, sub, err := _ClusterConnection.contract.WatchLogs(opts, "ValidatorSetAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ClusterConnectionValidatorSetAdded)
				if err := _ClusterConnection.contract.UnpackLog(event, "ValidatorSetAdded", log); err != nil {
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

// ParseValidatorSetAdded is a log parse operation binding the contract event 0xee022091c8c62247fb714803aafa9265a5b04fdd39f40668cc38a4170b21a290.
//
// Solidity: event ValidatorSetAdded(bytes[] _validator, uint8 _threshold)
func (_ClusterConnection *ClusterConnectionFilterer) ParseValidatorSetAdded(log types.Log) (*ClusterConnectionValidatorSetAdded, error) {
	event := new(ClusterConnectionValidatorSetAdded)
	if err := _ClusterConnection.contract.UnpackLog(event, "ValidatorSetAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
