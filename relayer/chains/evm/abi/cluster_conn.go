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
	ABI: "[{\"type\":\"function\",\"name\":\"admin\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"claimFees\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"connSn\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getFee\",\"inputs\":[{\"name\":\"to\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"response\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[{\"name\":\"fee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getReceipt\",\"inputs\":[{\"name\":\"srcNetwork\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_connSn\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRequiredValidatorCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_relayer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_xCall\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isValidator\",\"inputs\":[{\"name\":\"signer\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isValidatorProcessed\",\"inputs\":[{\"name\":\"processedSigners\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"signer\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"listValidators\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"recoverSigner\",\"inputs\":[{\"name\":\"messageHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"recvMessageWithSignatures\",\"inputs\":[{\"name\":\"srcNetwork\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_connSn\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_msg\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_signedMessages\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"relayer\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"revertMessage\",\"inputs\":[{\"name\":\"sn\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"sendMessage\",\"inputs\":[{\"name\":\"to\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_svc\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"sn\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"_msg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"setAdmin\",\"inputs\":[{\"name\":\"_address\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setFee\",\"inputs\":[{\"name\":\"networkId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"messageFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"responseFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRelayer\",\"inputs\":[{\"name\":\"_address\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRequiredValidatorCount\",\"inputs\":[{\"name\":\"_count\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateValidators\",\"inputs\":[{\"name\":\"_validators\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"_threshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Message\",\"inputs\":[{\"name\":\"targetNetwork\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"sn\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"_msg\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorSetAdded\",\"inputs\":[{\"name\":\"_validator\",\"type\":\"bytes[]\",\"indexed\":false,\"internalType\":\"bytes[]\"},{\"name\":\"_threshold\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]}]",
	Bin: "0x608060405234801561001057600080fd5b506125a2806100206000396000f3fe60806040526004361061011f5760003560e01c80637d4c4f4a116100a0578063a914a35211610064578063a914a35214610328578063cbd856fb14610348578063d294f09314610368578063f851a4401461037d578063facd743b1461039b57600080fd5b80637d4c4f4a146102625780638406c079146102905780639664da0e146102c257806397aba7f9146102f257806399f1fca71461031257600080fd5b8063522a901e116100e7578063522a901e146101c657806362534d93146101d95780636548e9bc1461020057806368d4e54414610220578063704b6c021461024257600080fd5b80630a523d68146101245780632d3fb82314610146578063308b74df1461016657806343f08a8914610186578063485cc955146101a6575b600080fd5b34801561013057600080fd5b5061014461013f366004611b07565b6103bb565b005b34801561015257600080fd5b50610144610161366004611bd3565b61070c565b34801561017257600080fd5b50610144610181366004611c02565b6107ff565b34801561019257600080fd5b506101446101a1366004611c24565b6108a7565b3480156101b257600080fd5b506101446101c1366004611c94565b610987565b6101446101d4366004611ccd565b610ac7565b3480156101e557600080fd5b5060085460405160ff90911681526020015b60405180910390f35b34801561020c57600080fd5b5061014461021b366004611d70565b610cb7565b34801561022c57600080fd5b50610235610d6b565b6040516101f79190611d8d565b34801561024e57600080fd5b5061014461025d366004611d70565b610dcd565b34801561026e57600080fd5b5061028261027d366004611dda565b610e81565b6040519081526020016101f7565b34801561029c57600080fd5b506004546001600160a01b03165b6040516001600160a01b0390911681526020016101f7565b3480156102ce57600080fd5b506102e26102dd366004611e25565b610ef0565b60405190151581526020016101f7565b3480156102fe57600080fd5b506102aa61030d366004611e69565b610f28565b34801561031e57600080fd5b5061028260065481565b34801561033457600080fd5b506102e2610343366004611ed2565b61106d565b34801561035457600080fd5b50610144610363366004611f7d565b6110c8565b34801561037457600080fd5b506101446112ab565b34801561038957600080fd5b506005546001600160a01b03166102aa565b3480156103a757600080fd5b506102e26103b6366004611d70565b611379565b306001600160a01b0316638406c0796040518163ffffffff1660e01b8152600401602060405180830381865afa1580156103f9573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061041d919061203f565b6001600160a01b0316336001600160a01b0316146104565760405162461bcd60e51b815260040161044d9061205c565b60405180910390fd5b60085460ff168110156104ab5760405162461bcd60e51b815260206004820152601c60248201527f4e6f7420656e6f756768207369676e6174757265732070617373656400000000604482015260640161044d565b60035460408051630e717cff60e21b815290516000926001600160a01b0316916339c5f3fc91600480830192869291908290030181865afa1580156104f4573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f1916820160405261051c91908101906120a5565b9050600061052d88888888866113d5565b9050600080846001600160401b0381111561054a5761054a6119fc565b604051908082528060200260200182016040528015610573578160200160208202803683370190505b50905060005b858110156106935760006105e5858989858181106105995761059961211b565b90506020028101906105ab9190612131565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250610f2892505050565b90506001600160a01b0381166106315760405162461bcd60e51b8152602060048201526011602482015270496e76616c6964207369676e617475726560781b604482015260640161044d565b61063b838261106d565b15801561064c575061064c8161146f565b1561068a57808385815181106106645761066461211b565b6001600160a01b0390921660209283029190910190910152836106868161218d565b9450505b50600101610579565b5060085460ff168210156106f45760405162461bcd60e51b815260206004820152602260248201527f4e6f7420656e6f7567682076616c6964207369676e6174757265732070617373604482015261195960f21b606482015260840161044d565b6107008a8a8a8a6114ce565b50505050505050505050565b306001600160a01b0316638406c0796040518163ffffffff1660e01b8152600401602060405180830381865afa15801561074a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061076e919061203f565b6001600160a01b0316336001600160a01b03161461079e5760405162461bcd60e51b815260040161044d9061205c565b60035460405163b070f9e560e01b8152600481018390526001600160a01b039091169063b070f9e590602401600060405180830381600087803b1580156107e457600080fd5b505af11580156107f8573d6000803e3d6000fd5b5050505050565b306001600160a01b031663f851a4406040518163ffffffff1660e01b8152600401602060405180830381865afa15801561083d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610861919061203f565b6001600160a01b0316336001600160a01b0316146108915760405162461bcd60e51b815260040161044d906121a6565b6008805460ff191660ff92909216919091179055565b306001600160a01b0316638406c0796040518163ffffffff1660e01b8152600401602060405180830381865afa1580156108e5573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610909919061203f565b6001600160a01b0316336001600160a01b0316146109395760405162461bcd60e51b815260040161044d9061205c565b816000858560405161094c9291906121c9565b90815260200160405180910390208190555080600185856040516109719291906121c9565b9081526040519081900360200190205550505050565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a008054600160401b810460ff1615906001600160401b03166000811580156109cc5750825b90506000826001600160401b031660011480156109e85750303b155b9050811580156109f6575080155b15610a145760405163f92ee8a960e01b815260040160405180910390fd5b845467ffffffffffffffff191660011785558315610a3e57845460ff60401b1916600160401b1785555b600380546001600160a01b038089166001600160a01b0319928316179092556005805482163317905560048054928a16929091169190911790558315610abe57845460ff60401b19168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b50505050505050565b6003546001600160a01b03163314610b215760405162461bcd60e51b815260206004820152601f60248201527f4f6e6c79205863616c6c2063616e2063616c6c2073656e644d65737361676500604482015260640161044d565b600080841315610b9b57604051633ea627a560e11b81523090637d4c4f4a90610b53908b908b90600190600401612202565b602060405180830381865afa158015610b70573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b949190612228565b9050610c0f565b83600003610c0f57604051633ea627a560e11b81523090637d4c4f4a90610bcb908b908b90600090600401612202565b602060405180830381865afa158015610be8573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c0c9190612228565b90505b80341015610c575760405162461bcd60e51b8152602060048201526015602482015274119959481a5cc81b9bdd0814dd59999a58da595b9d605a1b604482015260640161044d565b60068054906000610c678361218d565b91905055507f37be353f216cf7e33639101fd610c542e6a0c0109173fa1c1d8b04d34edb7c1b88886006548686604051610ca5959493929190612241565b60405180910390a15050505050505050565b306001600160a01b031663f851a4406040518163ffffffff1660e01b8152600401602060405180830381865afa158015610cf5573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610d19919061203f565b6001600160a01b0316336001600160a01b031614610d495760405162461bcd60e51b815260040161044d906121a6565b600480546001600160a01b0319166001600160a01b0392909216919091179055565b60606007805480602002602001604051908101604052809291908181526020018280548015610dc357602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311610da5575b5050505050905090565b306001600160a01b031663f851a4406040518163ffffffff1660e01b8152600401602060405180830381865afa158015610e0b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610e2f919061203f565b6001600160a01b0316336001600160a01b031614610e5f5760405162461bcd60e51b815260040161044d906121a6565b600580546001600160a01b0319166001600160a01b0392909216919091179055565b600080600084604051610e94919061227a565b908152604051908190036020019020549050821515600103610ee7576000600185604051610ec2919061227a565b908152604051908190036020019020549050610ede8183612296565b92505050610eea565b90505b92915050565b6000600283604051610f02919061227a565b908152604080516020928190038301902060009485529091529091205460ff1692915050565b60008151604114610f7b5760405162461bcd60e51b815260206004820152601860248201527f496e76616c6964207369676e6174757265206c656e6774680000000000000000604482015260640161044d565b60208201516040830151606084015160001a601b811015610fa457610fa1601b826122a9565b90505b8060ff16601b1480610fb957508060ff16601c145b6110055760405162461bcd60e51b815260206004820152601b60248201527f496e76616c6964207369676e6174757265202776272076616c75650000000000604482015260640161044d565b60408051600081526020810180835288905260ff831691810191909152606081018490526080810183905260019060a0016020604051602081039080840390855afa158015611058573d6000803e3d6000fd5b5050604051601f190151979650505050505050565b6000805b83518110156110be57826001600160a01b03168482815181106110965761109661211b565b60200260200101516001600160a01b0316036110b6576001915050610eea565b600101611071565b5060009392505050565b306001600160a01b031663f851a4406040518163ffffffff1660e01b8152600401602060405180830381865afa158015611106573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061112a919061203f565b6001600160a01b0316336001600160a01b03161461115a5760405162461bcd60e51b815260040161044d906121a6565b611166600760006119ca565b60005b82518110156112105760006111968483815181106111895761118961211b565b60200260200101516115e0565b90506111a181611379565b1580156111b657506001600160a01b03811615155b1561120757600780546001810182556000919091527fa66cc928b5edb82af9bd49922954155ab7b0942694bea4ce44661d9a8736c6880180546001600160a01b0319166001600160a01b0383161790555b50600101611169565b5060075460ff8216111561125e5760405162461bcd60e51b81526020600482015260156024820152744e6f7420656e6f7567682076616c696461746f727360581b604482015260640161044d565b6008805460ff191660ff83161790556040517fee022091c8c62247fb714803aafa9265a5b04fdd39f40668cc38a4170b21a2909061129f90849084906122ee565b60405180910390a15050565b306001600160a01b0316638406c0796040518163ffffffff1660e01b8152600401602060405180830381865afa1580156112e9573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061130d919061203f565b6001600160a01b0316336001600160a01b03161461133d5760405162461bcd60e51b815260040161044d9061205c565b6004546040516001600160a01b03909116904780156108fc02916000818181858888f19350505050158015611376573d6000803e3d6000fd5b50565b6000805b6007548110156113cf57826001600160a01b0316600782815481106113a4576113a461211b565b6000918252602090912001546001600160a01b0316036113c75750600192915050565b60010161137d565b50919050565b60008061145c6113e48861163f565b6113ed8861164a565b61142c88888080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061165892505050565b6114358b61163f565b604051602001611448949392919061235e565b6040516020818303038152906040526116c7565b8051602090910120979650505050505050565b6000805b6007548110156114c557826001600160a01b03166007828154811061149a5761149a61211b565b6000918252602090912001546001600160a01b0316036114bd5750600192915050565b600101611473565b50600092915050565b6002846040516114de919061227a565b90815260408051602092819003830190206000868152925290205460ff161561153d5760405162461bcd60e51b81526020600482015260116024820152704475706c6963617465204d65737361676560781b604482015260640161044d565b600160028560405161154f919061227a565b9081526040805160209281900383018120600088815293529120805460ff19169215159290921790915560035463bbc22efd60e01b82526001600160a01b03169063bbc22efd906115a8908790869086906004016123b5565b600060405180830381600087803b1580156115c257600080fd5b505af11580156115d6573d6000803e3d6000fd5b5050505050505050565b600081516041146116335760405162461bcd60e51b815260206004820152601960248201527f496e76616c6964207075626c6963206b6579206c656e67746800000000000000604482015260640161044d565b50604060219091012090565b6060610eea82611658565b6060610eea611658836116fd565b60608082516001148015611686575060808360008151811061167c5761167c61211b565b016020015160f81c105b15611692575081610eea565b61169e835160806117a7565b836040516020016116b09291906123e5565b604051602081830303815290604052905092915050565b60606116d5825160c06117a7565b826040516020016116e79291906123e5565b6040516020818303038152906040529050919050565b60608160000361172a57604080516001808252818301909252906020820181803683370190505092915050565b608060015b602081101561176157818410156117525761174a848261195e565b949350505050565b60089190911b9060010161172f565b508083101561178e576040805160208101859052015b604051602081830303815290604052915050919050565b6040516000602082015260218101849052604101611777565b606080603884101561181157604080516001808252818301909252906020820181803683370190505090506117dc8385612296565b601f1a60f81b816000815181106117f5576117f561211b565b60200101906001600160f81b031916908160001a905350610ee7565b600060015b611820818761242a565b15611846578161182f8161218d565b925061183f90506101008261243e565b9050611816565b611851826001612296565b6001600160401b03811115611868576118686119fc565b6040519080825280601f01601f191660200182016040528015611892576020820181803683370190505b50925061189f8583612296565b6118aa906037612296565b601f1a60f81b836000815181106118c3576118c361211b565b60200101906001600160f81b031916908160001a905350600190505b818111611954576101006118f38284612455565b6118ff9061010061254c565b611909908861242a565b6119139190612558565b601f1a60f81b83828151811061192b5761192b61211b565b60200101906001600160f81b031916908160001a9053508061194c8161218d565b9150506118df565b5050905092915050565b60606000826001600160401b0381111561197a5761197a6119fc565b6040519080825280601f01601f1916602001820160405280156119a4576020820181803683370190505b50905060208101836020035b60208110156119545785811a8253600191820191016119b0565b508054600082559060005260206000209081019061137691905b808211156119f857600081556001016119e4565b5090565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f191681016001600160401b0381118282101715611a3a57611a3a6119fc565b604052919050565b60006001600160401b03821115611a5b57611a5b6119fc565b50601f01601f191660200190565b600082601f830112611a7a57600080fd5b8135611a8d611a8882611a42565b611a12565b818152846020838601011115611aa257600080fd5b816020850160208301376000918101602001919091529392505050565b60008083601f840112611ad157600080fd5b5081356001600160401b03811115611ae857600080fd5b602083019150836020828501011115611b0057600080fd5b9250929050565b60008060008060008060808789031215611b2057600080fd5b86356001600160401b0380821115611b3757600080fd5b611b438a838b01611a69565b9750602089013596506040890135915080821115611b6057600080fd5b611b6c8a838b01611abf565b90965094506060890135915080821115611b8557600080fd5b818901915089601f830112611b9957600080fd5b813581811115611ba857600080fd5b8a60208260051b8501011115611bbd57600080fd5b6020830194508093505050509295509295509295565b600060208284031215611be557600080fd5b5035919050565b803560ff81168114611bfd57600080fd5b919050565b600060208284031215611c1457600080fd5b611c1d82611bec565b9392505050565b60008060008060608587031215611c3a57600080fd5b84356001600160401b03811115611c5057600080fd5b611c5c87828801611abf565b90989097506020870135966040013595509350505050565b6001600160a01b038116811461137657600080fd5b8035611bfd81611c74565b60008060408385031215611ca757600080fd5b8235611cb281611c74565b91506020830135611cc281611c74565b809150509250929050565b60008060008060008060006080888a031215611ce857600080fd5b87356001600160401b0380821115611cff57600080fd5b611d0b8b838c01611abf565b909950975060208a0135915080821115611d2457600080fd5b611d308b838c01611abf565b909750955060408a0135945060608a0135915080821115611d5057600080fd5b50611d5d8a828b01611abf565b989b979a50959850939692959293505050565b600060208284031215611d8257600080fd5b8135610ee781611c74565b6020808252825182820181905260009190848201906040850190845b81811015611dce5783516001600160a01b031683529284019291840191600101611da9565b50909695505050505050565b60008060408385031215611ded57600080fd5b82356001600160401b03811115611e0357600080fd5b611e0f85828601611a69565b92505060208301358015158114611cc257600080fd5b60008060408385031215611e3857600080fd5b82356001600160401b03811115611e4e57600080fd5b611e5a85828601611a69565b95602094909401359450505050565b60008060408385031215611e7c57600080fd5b8235915060208301356001600160401b03811115611e9957600080fd5b611ea585828601611a69565b9150509250929050565b60006001600160401b03821115611ec857611ec86119fc565b5060051b60200190565b60008060408385031215611ee557600080fd5b82356001600160401b03811115611efb57600080fd5b8301601f81018513611f0c57600080fd5b80356020611f1c611a8883611eaf565b82815260059290921b83018101918181019088841115611f3b57600080fd5b938201935b83851015611f62578435611f5381611c74565b82529382019390820190611f40565b9550611f719050868201611c89565b93505050509250929050565b60008060408385031215611f9057600080fd5b82356001600160401b0380821115611fa757600080fd5b818501915085601f830112611fbb57600080fd5b81356020611fcb611a8883611eaf565b82815260059290921b84018101918181019089841115611fea57600080fd5b8286015b84811015612022578035868111156120065760008081fd5b6120148c86838b0101611a69565b845250918301918301611fee565b5096506120329050878201611bec565b9450505050509250929050565b60006020828403121561205157600080fd5b8151610ee781611c74565b6020808252600b908201526a27b7363ca932b630bcb2b960a91b604082015260600190565b60005b8381101561209c578181015183820152602001612084565b50506000910152565b6000602082840312156120b757600080fd5b81516001600160401b038111156120cd57600080fd5b8201601f810184136120de57600080fd5b80516120ec611a8882611a42565b81815285602083850101111561210157600080fd5b612112826020830160208601612081565b95945050505050565b634e487b7160e01b600052603260045260246000fd5b6000808335601e1984360301811261214857600080fd5b8301803591506001600160401b0382111561216257600080fd5b602001915036819003821315611b0057600080fd5b634e487b7160e01b600052601160045260246000fd5b60006001820161219f5761219f612177565b5060010190565b60208082526009908201526827b7363ca0b236b4b760b91b604082015260600190565b8183823760009101908152919050565b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b6040815260006122166040830185876121d9565b90508215156020830152949350505050565b60006020828403121561223a57600080fd5b5051919050565b6060815260006122556060830187896121d9565b856020840152828103604084015261226e8185876121d9565b98975050505050505050565b6000825161228c818460208701612081565b9190910192915050565b80820180821115610eea57610eea612177565b60ff8181168382160190811115610eea57610eea612177565b600081518084526122da816020860160208601612081565b601f01601f19169290920160200192915050565b6000604082016040835280855180835260608501915060608160051b8601019250602080880160005b8381101561234557605f198887030185526123338683516122c2565b95509382019390820190600101612317565b50505050508091505060ff831660208301529392505050565b60008551612370818460208a01612081565b855190830190612384818360208a01612081565b8551910190612397818360208901612081565b84519101906123aa818360208801612081565b019695505050505050565b6040815260006123c860408301866122c2565b82810360208401526123db8185876121d9565b9695505050505050565b600083516123f7818460208801612081565b83519083019061240b818360208801612081565b01949350505050565b634e487b7160e01b600052601260045260246000fd5b60008261243957612439612414565b500490565b8082028115828204841417610eea57610eea612177565b81810381811115610eea57610eea612177565b600181815b808511156124a357816000190482111561248957612489612177565b8085161561249657918102915b93841c939080029061246d565b509250929050565b6000826124ba57506001610eea565b816124c757506000610eea565b81600181146124dd57600281146124e757612503565b6001915050610eea565b60ff8411156124f8576124f8612177565b50506001821b610eea565b5060208310610133831016604e8410600b8410161715612526575081810a610eea565b6125308383612468565b806000190482111561254457612544612177565b029392505050565b6000611c1d83836124ab565b60008261256757612567612414565b50069056fea2646970667358221220a01d7651401f730f7c158660dc39f4cd736734b3cc8ad423bcf502b6b146123564736f6c63430008180033",
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
