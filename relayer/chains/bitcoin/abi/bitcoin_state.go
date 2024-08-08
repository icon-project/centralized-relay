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

// IBitcoinStateTokenInfo is an auto generated low-level Go binding around an user-defined struct.
type IBitcoinStateTokenInfo struct {
	Name   string
	Symbol string
}

// StdInvariantFuzzArtifactSelector is an auto generated low-level Go binding around an user-defined struct.
type StdInvariantFuzzArtifactSelector struct {
	Artifact  string
	Selectors [][4]byte
}

// StdInvariantFuzzInterface is an auto generated low-level Go binding around an user-defined struct.
type StdInvariantFuzzInterface struct {
	Addr      common.Address
	Artifacts []string
}

// StdInvariantFuzzSelector is an auto generated low-level Go binding around an user-defined struct.
type StdInvariantFuzzSelector struct {
	Addr      common.Address
	Selectors [][4]byte
}

// BitcoinStateMetaData contains all meta data concerning the BitcoinState contract.
var BitcoinStateMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"IS_TEST\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"_PERMIT_TYPEHASH\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"accountBalances\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"addConnection\",\"inputs\":[{\"name\":\"connection_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"bitcoinNid\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"claimTokens\",\"inputs\":[{\"name\":\"token0\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"token1\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"computeTokenAddress\",\"inputs\":[{\"name\":\"tokenName\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"connections\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"connectionsEndpoints\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"excludeArtifacts\",\"inputs\":[],\"outputs\":[{\"name\":\"excludedArtifacts_\",\"type\":\"string[]\",\"internalType\":\"string[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"excludeContracts\",\"inputs\":[],\"outputs\":[{\"name\":\"excludedContracts_\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"excludeSelectors\",\"inputs\":[],\"outputs\":[{\"name\":\"excludedSelectors_\",\"type\":\"tuple[]\",\"internalType\":\"structStdInvariant.FuzzSelector[]\",\"components\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"selectors\",\"type\":\"bytes4[]\",\"internalType\":\"bytes4[]\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"excludeSenders\",\"inputs\":[],\"outputs\":[{\"name\":\"excludedSenders_\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"failed\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSignData\",\"inputs\":[{\"name\":\"requester_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"handleCallMessage\",\"inputs\":[{\"name\":\"_from\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_protocols\",\"type\":\"string[]\",\"internalType\":\"string[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initPool\",\"inputs\":[{\"name\":\"data_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"xcall_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"uinswapV3Router_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nonfungiblePositionManager_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"connections\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"migrate\",\"inputs\":[{\"name\":\"_data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"migrateComplete\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"nftOwners\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"nonFungibleManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"params\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIBitcoinState.TokenInfo\",\"components\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol\",\"type\":\"string\",\"internalType\":\"string\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"removeConnection\",\"inputs\":[{\"name\":\"connection_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"removeLiquidity\",\"inputs\":[{\"name\":\"data_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"routerV2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"targetArtifactSelectors\",\"inputs\":[],\"outputs\":[{\"name\":\"targetedArtifactSelectors_\",\"type\":\"tuple[]\",\"internalType\":\"structStdInvariant.FuzzArtifactSelector[]\",\"components\":[{\"name\":\"artifact\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"selectors\",\"type\":\"bytes4[]\",\"internalType\":\"bytes4[]\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"targetArtifacts\",\"inputs\":[],\"outputs\":[{\"name\":\"targetedArtifacts_\",\"type\":\"string[]\",\"internalType\":\"string[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"targetContracts\",\"inputs\":[],\"outputs\":[{\"name\":\"targetedContracts_\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"targetInterfaces\",\"inputs\":[],\"outputs\":[{\"name\":\"targetedInterfaces_\",\"type\":\"tuple[]\",\"internalType\":\"structStdInvariant.FuzzInterface[]\",\"components\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"artifacts\",\"type\":\"string[]\",\"internalType\":\"string[]\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"targetSelectors\",\"inputs\":[],\"outputs\":[{\"name\":\"targetedSelectors_\",\"type\":\"tuple[]\",\"internalType\":\"structStdInvariant.FuzzSelector[]\",\"components\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"selectors\",\"type\":\"bytes4[]\",\"internalType\":\"bytes4[]\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"targetSenders\",\"inputs\":[],\"outputs\":[{\"name\":\"targetedSenders_\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"tokens\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"xcallService\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"AddConnection\",\"inputs\":[{\"name\":\"connection_\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AddSelector\",\"inputs\":[{\"name\":\"selector_\",\"type\":\"bytes4\",\"indexed\":false,\"internalType\":\"bytes4\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RemoveConnection\",\"inputs\":[{\"name\":\"connection_\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RemoveSelector\",\"inputs\":[{\"name\":\"selector_\",\"type\":\"bytes4\",\"indexed\":false,\"internalType\":\"bytes4\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RequestExecuted\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"stateRoot\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"data\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"log\",\"inputs\":[{\"name\":\"\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"log_address\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"log_array\",\"inputs\":[{\"name\":\"val\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"log_array\",\"inputs\":[{\"name\":\"val\",\"type\":\"int256[]\",\"indexed\":false,\"internalType\":\"int256[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"log_array\",\"inputs\":[{\"name\":\"val\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"log_bytes\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"log_bytes32\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"log_int\",\"inputs\":[{\"name\":\"\",\"type\":\"int256\",\"indexed\":false,\"internalType\":\"int256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"log_named_address\",\"inputs\":[{\"name\":\"key\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"val\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"log_named_array\",\"inputs\":[{\"name\":\"key\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"val\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"log_named_array\",\"inputs\":[{\"name\":\"key\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"val\",\"type\":\"int256[]\",\"indexed\":false,\"internalType\":\"int256[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"log_named_array\",\"inputs\":[{\"name\":\"key\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"val\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"log_named_bytes\",\"inputs\":[{\"name\":\"key\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"val\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"log_named_bytes32\",\"inputs\":[{\"name\":\"key\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"val\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"log_named_decimal_int\",\"inputs\":[{\"name\":\"key\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"val\",\"type\":\"int256\",\"indexed\":false,\"internalType\":\"int256\"},{\"name\":\"decimals\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"log_named_decimal_uint\",\"inputs\":[{\"name\":\"key\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"val\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"decimals\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"log_named_int\",\"inputs\":[{\"name\":\"key\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"val\",\"type\":\"int256\",\"indexed\":false,\"internalType\":\"int256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"log_named_string\",\"inputs\":[{\"name\":\"key\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"val\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"log_named_uint\",\"inputs\":[{\"name\":\"key\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"val\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"log_string\",\"inputs\":[{\"name\":\"\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"log_uint\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"logs\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false}]",
}

// BitcoinStateABI is the input ABI used to generate the binding from.
// Deprecated: Use BitcoinStateMetaData.ABI instead.
var BitcoinStateABI = BitcoinStateMetaData.ABI

// BitcoinState is an auto generated Go binding around an Ethereum contract.
type BitcoinState struct {
	BitcoinStateCaller     // Read-only binding to the contract
	BitcoinStateTransactor // Write-only binding to the contract
	BitcoinStateFilterer   // Log filterer for contract events
}

// BitcoinStateCaller is an auto generated read-only Go binding around an Ethereum contract.
type BitcoinStateCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BitcoinStateTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BitcoinStateTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BitcoinStateFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BitcoinStateFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BitcoinStateSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BitcoinStateSession struct {
	Contract     *BitcoinState     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BitcoinStateCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BitcoinStateCallerSession struct {
	Contract *BitcoinStateCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// BitcoinStateTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BitcoinStateTransactorSession struct {
	Contract     *BitcoinStateTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// BitcoinStateRaw is an auto generated low-level Go binding around an Ethereum contract.
type BitcoinStateRaw struct {
	Contract *BitcoinState // Generic contract binding to access the raw methods on
}

// BitcoinStateCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BitcoinStateCallerRaw struct {
	Contract *BitcoinStateCaller // Generic read-only contract binding to access the raw methods on
}

// BitcoinStateTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BitcoinStateTransactorRaw struct {
	Contract *BitcoinStateTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBitcoinState creates a new instance of BitcoinState, bound to a specific deployed contract.
func NewBitcoinState(address common.Address, backend bind.ContractBackend) (*BitcoinState, error) {
	contract, err := bindBitcoinState(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BitcoinState{BitcoinStateCaller: BitcoinStateCaller{contract: contract}, BitcoinStateTransactor: BitcoinStateTransactor{contract: contract}, BitcoinStateFilterer: BitcoinStateFilterer{contract: contract}}, nil
}

// NewBitcoinStateCaller creates a new read-only instance of BitcoinState, bound to a specific deployed contract.
func NewBitcoinStateCaller(address common.Address, caller bind.ContractCaller) (*BitcoinStateCaller, error) {
	contract, err := bindBitcoinState(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BitcoinStateCaller{contract: contract}, nil
}

// NewBitcoinStateTransactor creates a new write-only instance of BitcoinState, bound to a specific deployed contract.
func NewBitcoinStateTransactor(address common.Address, transactor bind.ContractTransactor) (*BitcoinStateTransactor, error) {
	contract, err := bindBitcoinState(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BitcoinStateTransactor{contract: contract}, nil
}

// NewBitcoinStateFilterer creates a new log filterer instance of BitcoinState, bound to a specific deployed contract.
func NewBitcoinStateFilterer(address common.Address, filterer bind.ContractFilterer) (*BitcoinStateFilterer, error) {
	contract, err := bindBitcoinState(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BitcoinStateFilterer{contract: contract}, nil
}

// bindBitcoinState binds a generic wrapper to an already deployed contract.
func bindBitcoinState(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BitcoinStateMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BitcoinState *BitcoinStateRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BitcoinState.Contract.BitcoinStateCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BitcoinState *BitcoinStateRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BitcoinState.Contract.BitcoinStateTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BitcoinState *BitcoinStateRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BitcoinState.Contract.BitcoinStateTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BitcoinState *BitcoinStateCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BitcoinState.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BitcoinState *BitcoinStateTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BitcoinState.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BitcoinState *BitcoinStateTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BitcoinState.Contract.contract.Transact(opts, method, params...)
}

// ISTEST is a free data retrieval call binding the contract method 0xfa7626d4.
//
// Solidity: function IS_TEST() view returns(bool)
func (_BitcoinState *BitcoinStateCaller) ISTEST(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _BitcoinState.contract.Call(opts, &out, "IS_TEST")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ISTEST is a free data retrieval call binding the contract method 0xfa7626d4.
//
// Solidity: function IS_TEST() view returns(bool)
func (_BitcoinState *BitcoinStateSession) ISTEST() (bool, error) {
	return _BitcoinState.Contract.ISTEST(&_BitcoinState.CallOpts)
}

// ISTEST is a free data retrieval call binding the contract method 0xfa7626d4.
//
// Solidity: function IS_TEST() view returns(bool)
func (_BitcoinState *BitcoinStateCallerSession) ISTEST() (bool, error) {
	return _BitcoinState.Contract.ISTEST(&_BitcoinState.CallOpts)
}

// PERMITTYPEHASH is a free data retrieval call binding the contract method 0x982aaf6b.
//
// Solidity: function _PERMIT_TYPEHASH() view returns(bytes32)
func (_BitcoinState *BitcoinStateCaller) PERMITTYPEHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _BitcoinState.contract.Call(opts, &out, "_PERMIT_TYPEHASH")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// PERMITTYPEHASH is a free data retrieval call binding the contract method 0x982aaf6b.
//
// Solidity: function _PERMIT_TYPEHASH() view returns(bytes32)
func (_BitcoinState *BitcoinStateSession) PERMITTYPEHASH() ([32]byte, error) {
	return _BitcoinState.Contract.PERMITTYPEHASH(&_BitcoinState.CallOpts)
}

// PERMITTYPEHASH is a free data retrieval call binding the contract method 0x982aaf6b.
//
// Solidity: function _PERMIT_TYPEHASH() view returns(bytes32)
func (_BitcoinState *BitcoinStateCallerSession) PERMITTYPEHASH() ([32]byte, error) {
	return _BitcoinState.Contract.PERMITTYPEHASH(&_BitcoinState.CallOpts)
}

// AccountBalances is a free data retrieval call binding the contract method 0x1242e5fd.
//
// Solidity: function accountBalances(address , address ) view returns(uint256)
func (_BitcoinState *BitcoinStateCaller) AccountBalances(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _BitcoinState.contract.Call(opts, &out, "accountBalances", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AccountBalances is a free data retrieval call binding the contract method 0x1242e5fd.
//
// Solidity: function accountBalances(address , address ) view returns(uint256)
func (_BitcoinState *BitcoinStateSession) AccountBalances(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _BitcoinState.Contract.AccountBalances(&_BitcoinState.CallOpts, arg0, arg1)
}

// AccountBalances is a free data retrieval call binding the contract method 0x1242e5fd.
//
// Solidity: function accountBalances(address , address ) view returns(uint256)
func (_BitcoinState *BitcoinStateCallerSession) AccountBalances(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _BitcoinState.Contract.AccountBalances(&_BitcoinState.CallOpts, arg0, arg1)
}

// BitcoinNid is a free data retrieval call binding the contract method 0xf3cc4515.
//
// Solidity: function bitcoinNid() view returns(string)
func (_BitcoinState *BitcoinStateCaller) BitcoinNid(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _BitcoinState.contract.Call(opts, &out, "bitcoinNid")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// BitcoinNid is a free data retrieval call binding the contract method 0xf3cc4515.
//
// Solidity: function bitcoinNid() view returns(string)
func (_BitcoinState *BitcoinStateSession) BitcoinNid() (string, error) {
	return _BitcoinState.Contract.BitcoinNid(&_BitcoinState.CallOpts)
}

// BitcoinNid is a free data retrieval call binding the contract method 0xf3cc4515.
//
// Solidity: function bitcoinNid() view returns(string)
func (_BitcoinState *BitcoinStateCallerSession) BitcoinNid() (string, error) {
	return _BitcoinState.Contract.BitcoinNid(&_BitcoinState.CallOpts)
}

// ClaimTokens is a free data retrieval call binding the contract method 0x69ffa08a.
//
// Solidity: function claimTokens(address token0, address token1) pure returns()
func (_BitcoinState *BitcoinStateCaller) ClaimTokens(opts *bind.CallOpts, token0 common.Address, token1 common.Address) error {
	var out []interface{}
	err := _BitcoinState.contract.Call(opts, &out, "claimTokens", token0, token1)

	if err != nil {
		return err
	}

	return err

}

// ClaimTokens is a free data retrieval call binding the contract method 0x69ffa08a.
//
// Solidity: function claimTokens(address token0, address token1) pure returns()
func (_BitcoinState *BitcoinStateSession) ClaimTokens(token0 common.Address, token1 common.Address) error {
	return _BitcoinState.Contract.ClaimTokens(&_BitcoinState.CallOpts, token0, token1)
}

// ClaimTokens is a free data retrieval call binding the contract method 0x69ffa08a.
//
// Solidity: function claimTokens(address token0, address token1) pure returns()
func (_BitcoinState *BitcoinStateCallerSession) ClaimTokens(token0 common.Address, token1 common.Address) error {
	return _BitcoinState.Contract.ClaimTokens(&_BitcoinState.CallOpts, token0, token1)
}

// ComputeTokenAddress is a free data retrieval call binding the contract method 0x4aa6af7b.
//
// Solidity: function computeTokenAddress(string tokenName) view returns(address)
func (_BitcoinState *BitcoinStateCaller) ComputeTokenAddress(opts *bind.CallOpts, tokenName string) (common.Address, error) {
	var out []interface{}
	err := _BitcoinState.contract.Call(opts, &out, "computeTokenAddress", tokenName)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ComputeTokenAddress is a free data retrieval call binding the contract method 0x4aa6af7b.
//
// Solidity: function computeTokenAddress(string tokenName) view returns(address)
func (_BitcoinState *BitcoinStateSession) ComputeTokenAddress(tokenName string) (common.Address, error) {
	return _BitcoinState.Contract.ComputeTokenAddress(&_BitcoinState.CallOpts, tokenName)
}

// ComputeTokenAddress is a free data retrieval call binding the contract method 0x4aa6af7b.
//
// Solidity: function computeTokenAddress(string tokenName) view returns(address)
func (_BitcoinState *BitcoinStateCallerSession) ComputeTokenAddress(tokenName string) (common.Address, error) {
	return _BitcoinState.Contract.ComputeTokenAddress(&_BitcoinState.CallOpts, tokenName)
}

// Connections is a free data retrieval call binding the contract method 0xc0896578.
//
// Solidity: function connections(address ) view returns(bool)
func (_BitcoinState *BitcoinStateCaller) Connections(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _BitcoinState.contract.Call(opts, &out, "connections", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Connections is a free data retrieval call binding the contract method 0xc0896578.
//
// Solidity: function connections(address ) view returns(bool)
func (_BitcoinState *BitcoinStateSession) Connections(arg0 common.Address) (bool, error) {
	return _BitcoinState.Contract.Connections(&_BitcoinState.CallOpts, arg0)
}

// Connections is a free data retrieval call binding the contract method 0xc0896578.
//
// Solidity: function connections(address ) view returns(bool)
func (_BitcoinState *BitcoinStateCallerSession) Connections(arg0 common.Address) (bool, error) {
	return _BitcoinState.Contract.Connections(&_BitcoinState.CallOpts, arg0)
}

// ConnectionsEndpoints is a free data retrieval call binding the contract method 0xb687db2e.
//
// Solidity: function connectionsEndpoints(uint256 ) view returns(address)
func (_BitcoinState *BitcoinStateCaller) ConnectionsEndpoints(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _BitcoinState.contract.Call(opts, &out, "connectionsEndpoints", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ConnectionsEndpoints is a free data retrieval call binding the contract method 0xb687db2e.
//
// Solidity: function connectionsEndpoints(uint256 ) view returns(address)
func (_BitcoinState *BitcoinStateSession) ConnectionsEndpoints(arg0 *big.Int) (common.Address, error) {
	return _BitcoinState.Contract.ConnectionsEndpoints(&_BitcoinState.CallOpts, arg0)
}

// ConnectionsEndpoints is a free data retrieval call binding the contract method 0xb687db2e.
//
// Solidity: function connectionsEndpoints(uint256 ) view returns(address)
func (_BitcoinState *BitcoinStateCallerSession) ConnectionsEndpoints(arg0 *big.Int) (common.Address, error) {
	return _BitcoinState.Contract.ConnectionsEndpoints(&_BitcoinState.CallOpts, arg0)
}

// ExcludeArtifacts is a free data retrieval call binding the contract method 0xb5508aa9.
//
// Solidity: function excludeArtifacts() view returns(string[] excludedArtifacts_)
func (_BitcoinState *BitcoinStateCaller) ExcludeArtifacts(opts *bind.CallOpts) ([]string, error) {
	var out []interface{}
	err := _BitcoinState.contract.Call(opts, &out, "excludeArtifacts")

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// ExcludeArtifacts is a free data retrieval call binding the contract method 0xb5508aa9.
//
// Solidity: function excludeArtifacts() view returns(string[] excludedArtifacts_)
func (_BitcoinState *BitcoinStateSession) ExcludeArtifacts() ([]string, error) {
	return _BitcoinState.Contract.ExcludeArtifacts(&_BitcoinState.CallOpts)
}

// ExcludeArtifacts is a free data retrieval call binding the contract method 0xb5508aa9.
//
// Solidity: function excludeArtifacts() view returns(string[] excludedArtifacts_)
func (_BitcoinState *BitcoinStateCallerSession) ExcludeArtifacts() ([]string, error) {
	return _BitcoinState.Contract.ExcludeArtifacts(&_BitcoinState.CallOpts)
}

// ExcludeContracts is a free data retrieval call binding the contract method 0xe20c9f71.
//
// Solidity: function excludeContracts() view returns(address[] excludedContracts_)
func (_BitcoinState *BitcoinStateCaller) ExcludeContracts(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _BitcoinState.contract.Call(opts, &out, "excludeContracts")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// ExcludeContracts is a free data retrieval call binding the contract method 0xe20c9f71.
//
// Solidity: function excludeContracts() view returns(address[] excludedContracts_)
func (_BitcoinState *BitcoinStateSession) ExcludeContracts() ([]common.Address, error) {
	return _BitcoinState.Contract.ExcludeContracts(&_BitcoinState.CallOpts)
}

// ExcludeContracts is a free data retrieval call binding the contract method 0xe20c9f71.
//
// Solidity: function excludeContracts() view returns(address[] excludedContracts_)
func (_BitcoinState *BitcoinStateCallerSession) ExcludeContracts() ([]common.Address, error) {
	return _BitcoinState.Contract.ExcludeContracts(&_BitcoinState.CallOpts)
}

// ExcludeSelectors is a free data retrieval call binding the contract method 0xb0464fdc.
//
// Solidity: function excludeSelectors() view returns((address,bytes4[])[] excludedSelectors_)
func (_BitcoinState *BitcoinStateCaller) ExcludeSelectors(opts *bind.CallOpts) ([]StdInvariantFuzzSelector, error) {
	var out []interface{}
	err := _BitcoinState.contract.Call(opts, &out, "excludeSelectors")

	if err != nil {
		return *new([]StdInvariantFuzzSelector), err
	}

	out0 := *abi.ConvertType(out[0], new([]StdInvariantFuzzSelector)).(*[]StdInvariantFuzzSelector)

	return out0, err

}

// ExcludeSelectors is a free data retrieval call binding the contract method 0xb0464fdc.
//
// Solidity: function excludeSelectors() view returns((address,bytes4[])[] excludedSelectors_)
func (_BitcoinState *BitcoinStateSession) ExcludeSelectors() ([]StdInvariantFuzzSelector, error) {
	return _BitcoinState.Contract.ExcludeSelectors(&_BitcoinState.CallOpts)
}

// ExcludeSelectors is a free data retrieval call binding the contract method 0xb0464fdc.
//
// Solidity: function excludeSelectors() view returns((address,bytes4[])[] excludedSelectors_)
func (_BitcoinState *BitcoinStateCallerSession) ExcludeSelectors() ([]StdInvariantFuzzSelector, error) {
	return _BitcoinState.Contract.ExcludeSelectors(&_BitcoinState.CallOpts)
}

// ExcludeSenders is a free data retrieval call binding the contract method 0x1ed7831c.
//
// Solidity: function excludeSenders() view returns(address[] excludedSenders_)
func (_BitcoinState *BitcoinStateCaller) ExcludeSenders(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _BitcoinState.contract.Call(opts, &out, "excludeSenders")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// ExcludeSenders is a free data retrieval call binding the contract method 0x1ed7831c.
//
// Solidity: function excludeSenders() view returns(address[] excludedSenders_)
func (_BitcoinState *BitcoinStateSession) ExcludeSenders() ([]common.Address, error) {
	return _BitcoinState.Contract.ExcludeSenders(&_BitcoinState.CallOpts)
}

// ExcludeSenders is a free data retrieval call binding the contract method 0x1ed7831c.
//
// Solidity: function excludeSenders() view returns(address[] excludedSenders_)
func (_BitcoinState *BitcoinStateCallerSession) ExcludeSenders() ([]common.Address, error) {
	return _BitcoinState.Contract.ExcludeSenders(&_BitcoinState.CallOpts)
}

// Failed is a free data retrieval call binding the contract method 0xba414fa6.
//
// Solidity: function failed() view returns(bool)
func (_BitcoinState *BitcoinStateCaller) Failed(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _BitcoinState.contract.Call(opts, &out, "failed")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Failed is a free data retrieval call binding the contract method 0xba414fa6.
//
// Solidity: function failed() view returns(bool)
func (_BitcoinState *BitcoinStateSession) Failed() (bool, error) {
	return _BitcoinState.Contract.Failed(&_BitcoinState.CallOpts)
}

// Failed is a free data retrieval call binding the contract method 0xba414fa6.
//
// Solidity: function failed() view returns(bool)
func (_BitcoinState *BitcoinStateCallerSession) Failed() (bool, error) {
	return _BitcoinState.Contract.Failed(&_BitcoinState.CallOpts)
}

// InitPool is a free data retrieval call binding the contract method 0xca38b326.
//
// Solidity: function initPool(bytes data_) pure returns()
func (_BitcoinState *BitcoinStateCaller) InitPool(opts *bind.CallOpts, data_ []byte) error {
	var out []interface{}
	err := _BitcoinState.contract.Call(opts, &out, "initPool", data_)

	if err != nil {
		return err
	}

	return err

}

// InitPool is a free data retrieval call binding the contract method 0xca38b326.
//
// Solidity: function initPool(bytes data_) pure returns()
func (_BitcoinState *BitcoinStateSession) InitPool(data_ []byte) error {
	return _BitcoinState.Contract.InitPool(&_BitcoinState.CallOpts, data_)
}

// InitPool is a free data retrieval call binding the contract method 0xca38b326.
//
// Solidity: function initPool(bytes data_) pure returns()
func (_BitcoinState *BitcoinStateCallerSession) InitPool(data_ []byte) error {
	return _BitcoinState.Contract.InitPool(&_BitcoinState.CallOpts, data_)
}

// NftOwners is a free data retrieval call binding the contract method 0xbbd94c2f.
//
// Solidity: function nftOwners(uint256 ) view returns(address)
func (_BitcoinState *BitcoinStateCaller) NftOwners(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _BitcoinState.contract.Call(opts, &out, "nftOwners", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// NftOwners is a free data retrieval call binding the contract method 0xbbd94c2f.
//
// Solidity: function nftOwners(uint256 ) view returns(address)
func (_BitcoinState *BitcoinStateSession) NftOwners(arg0 *big.Int) (common.Address, error) {
	return _BitcoinState.Contract.NftOwners(&_BitcoinState.CallOpts, arg0)
}

// NftOwners is a free data retrieval call binding the contract method 0xbbd94c2f.
//
// Solidity: function nftOwners(uint256 ) view returns(address)
func (_BitcoinState *BitcoinStateCallerSession) NftOwners(arg0 *big.Int) (common.Address, error) {
	return _BitcoinState.Contract.NftOwners(&_BitcoinState.CallOpts, arg0)
}

// NonFungibleManager is a free data retrieval call binding the contract method 0xecc28165.
//
// Solidity: function nonFungibleManager() view returns(address)
func (_BitcoinState *BitcoinStateCaller) NonFungibleManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BitcoinState.contract.Call(opts, &out, "nonFungibleManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// NonFungibleManager is a free data retrieval call binding the contract method 0xecc28165.
//
// Solidity: function nonFungibleManager() view returns(address)
func (_BitcoinState *BitcoinStateSession) NonFungibleManager() (common.Address, error) {
	return _BitcoinState.Contract.NonFungibleManager(&_BitcoinState.CallOpts)
}

// NonFungibleManager is a free data retrieval call binding the contract method 0xecc28165.
//
// Solidity: function nonFungibleManager() view returns(address)
func (_BitcoinState *BitcoinStateCallerSession) NonFungibleManager() (common.Address, error) {
	return _BitcoinState.Contract.NonFungibleManager(&_BitcoinState.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_BitcoinState *BitcoinStateCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BitcoinState.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_BitcoinState *BitcoinStateSession) Owner() (common.Address, error) {
	return _BitcoinState.Contract.Owner(&_BitcoinState.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_BitcoinState *BitcoinStateCallerSession) Owner() (common.Address, error) {
	return _BitcoinState.Contract.Owner(&_BitcoinState.CallOpts)
}

// Params is a free data retrieval call binding the contract method 0xcff0ab96.
//
// Solidity: function params() view returns((string,string))
func (_BitcoinState *BitcoinStateCaller) Params(opts *bind.CallOpts) (IBitcoinStateTokenInfo, error) {
	var out []interface{}
	err := _BitcoinState.contract.Call(opts, &out, "params")

	if err != nil {
		return *new(IBitcoinStateTokenInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(IBitcoinStateTokenInfo)).(*IBitcoinStateTokenInfo)

	return out0, err

}

// Params is a free data retrieval call binding the contract method 0xcff0ab96.
//
// Solidity: function params() view returns((string,string))
func (_BitcoinState *BitcoinStateSession) Params() (IBitcoinStateTokenInfo, error) {
	return _BitcoinState.Contract.Params(&_BitcoinState.CallOpts)
}

// Params is a free data retrieval call binding the contract method 0xcff0ab96.
//
// Solidity: function params() view returns((string,string))
func (_BitcoinState *BitcoinStateCallerSession) Params() (IBitcoinStateTokenInfo, error) {
	return _BitcoinState.Contract.Params(&_BitcoinState.CallOpts)
}

// RemoveLiquidity is a free data retrieval call binding the contract method 0x028318cf.
//
// Solidity: function removeLiquidity(bytes data_) pure returns()
func (_BitcoinState *BitcoinStateCaller) RemoveLiquidity(opts *bind.CallOpts, data_ []byte) error {
	var out []interface{}
	err := _BitcoinState.contract.Call(opts, &out, "removeLiquidity", data_)

	if err != nil {
		return err
	}

	return err

}

// RemoveLiquidity is a free data retrieval call binding the contract method 0x028318cf.
//
// Solidity: function removeLiquidity(bytes data_) pure returns()
func (_BitcoinState *BitcoinStateSession) RemoveLiquidity(data_ []byte) error {
	return _BitcoinState.Contract.RemoveLiquidity(&_BitcoinState.CallOpts, data_)
}

// RemoveLiquidity is a free data retrieval call binding the contract method 0x028318cf.
//
// Solidity: function removeLiquidity(bytes data_) pure returns()
func (_BitcoinState *BitcoinStateCallerSession) RemoveLiquidity(data_ []byte) error {
	return _BitcoinState.Contract.RemoveLiquidity(&_BitcoinState.CallOpts, data_)
}

// RouterV2 is a free data retrieval call binding the contract method 0x502f7446.
//
// Solidity: function routerV2() view returns(address)
func (_BitcoinState *BitcoinStateCaller) RouterV2(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BitcoinState.contract.Call(opts, &out, "routerV2")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RouterV2 is a free data retrieval call binding the contract method 0x502f7446.
//
// Solidity: function routerV2() view returns(address)
func (_BitcoinState *BitcoinStateSession) RouterV2() (common.Address, error) {
	return _BitcoinState.Contract.RouterV2(&_BitcoinState.CallOpts)
}

// RouterV2 is a free data retrieval call binding the contract method 0x502f7446.
//
// Solidity: function routerV2() view returns(address)
func (_BitcoinState *BitcoinStateCallerSession) RouterV2() (common.Address, error) {
	return _BitcoinState.Contract.RouterV2(&_BitcoinState.CallOpts)
}

// TargetArtifactSelectors is a free data retrieval call binding the contract method 0x66d9a9a0.
//
// Solidity: function targetArtifactSelectors() view returns((string,bytes4[])[] targetedArtifactSelectors_)
func (_BitcoinState *BitcoinStateCaller) TargetArtifactSelectors(opts *bind.CallOpts) ([]StdInvariantFuzzArtifactSelector, error) {
	var out []interface{}
	err := _BitcoinState.contract.Call(opts, &out, "targetArtifactSelectors")

	if err != nil {
		return *new([]StdInvariantFuzzArtifactSelector), err
	}

	out0 := *abi.ConvertType(out[0], new([]StdInvariantFuzzArtifactSelector)).(*[]StdInvariantFuzzArtifactSelector)

	return out0, err

}

// TargetArtifactSelectors is a free data retrieval call binding the contract method 0x66d9a9a0.
//
// Solidity: function targetArtifactSelectors() view returns((string,bytes4[])[] targetedArtifactSelectors_)
func (_BitcoinState *BitcoinStateSession) TargetArtifactSelectors() ([]StdInvariantFuzzArtifactSelector, error) {
	return _BitcoinState.Contract.TargetArtifactSelectors(&_BitcoinState.CallOpts)
}

// TargetArtifactSelectors is a free data retrieval call binding the contract method 0x66d9a9a0.
//
// Solidity: function targetArtifactSelectors() view returns((string,bytes4[])[] targetedArtifactSelectors_)
func (_BitcoinState *BitcoinStateCallerSession) TargetArtifactSelectors() ([]StdInvariantFuzzArtifactSelector, error) {
	return _BitcoinState.Contract.TargetArtifactSelectors(&_BitcoinState.CallOpts)
}

// TargetArtifacts is a free data retrieval call binding the contract method 0x85226c81.
//
// Solidity: function targetArtifacts() view returns(string[] targetedArtifacts_)
func (_BitcoinState *BitcoinStateCaller) TargetArtifacts(opts *bind.CallOpts) ([]string, error) {
	var out []interface{}
	err := _BitcoinState.contract.Call(opts, &out, "targetArtifacts")

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// TargetArtifacts is a free data retrieval call binding the contract method 0x85226c81.
//
// Solidity: function targetArtifacts() view returns(string[] targetedArtifacts_)
func (_BitcoinState *BitcoinStateSession) TargetArtifacts() ([]string, error) {
	return _BitcoinState.Contract.TargetArtifacts(&_BitcoinState.CallOpts)
}

// TargetArtifacts is a free data retrieval call binding the contract method 0x85226c81.
//
// Solidity: function targetArtifacts() view returns(string[] targetedArtifacts_)
func (_BitcoinState *BitcoinStateCallerSession) TargetArtifacts() ([]string, error) {
	return _BitcoinState.Contract.TargetArtifacts(&_BitcoinState.CallOpts)
}

// TargetContracts is a free data retrieval call binding the contract method 0x3f7286f4.
//
// Solidity: function targetContracts() view returns(address[] targetedContracts_)
func (_BitcoinState *BitcoinStateCaller) TargetContracts(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _BitcoinState.contract.Call(opts, &out, "targetContracts")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// TargetContracts is a free data retrieval call binding the contract method 0x3f7286f4.
//
// Solidity: function targetContracts() view returns(address[] targetedContracts_)
func (_BitcoinState *BitcoinStateSession) TargetContracts() ([]common.Address, error) {
	return _BitcoinState.Contract.TargetContracts(&_BitcoinState.CallOpts)
}

// TargetContracts is a free data retrieval call binding the contract method 0x3f7286f4.
//
// Solidity: function targetContracts() view returns(address[] targetedContracts_)
func (_BitcoinState *BitcoinStateCallerSession) TargetContracts() ([]common.Address, error) {
	return _BitcoinState.Contract.TargetContracts(&_BitcoinState.CallOpts)
}

// TargetInterfaces is a free data retrieval call binding the contract method 0x2ade3880.
//
// Solidity: function targetInterfaces() view returns((address,string[])[] targetedInterfaces_)
func (_BitcoinState *BitcoinStateCaller) TargetInterfaces(opts *bind.CallOpts) ([]StdInvariantFuzzInterface, error) {
	var out []interface{}
	err := _BitcoinState.contract.Call(opts, &out, "targetInterfaces")

	if err != nil {
		return *new([]StdInvariantFuzzInterface), err
	}

	out0 := *abi.ConvertType(out[0], new([]StdInvariantFuzzInterface)).(*[]StdInvariantFuzzInterface)

	return out0, err

}

// TargetInterfaces is a free data retrieval call binding the contract method 0x2ade3880.
//
// Solidity: function targetInterfaces() view returns((address,string[])[] targetedInterfaces_)
func (_BitcoinState *BitcoinStateSession) TargetInterfaces() ([]StdInvariantFuzzInterface, error) {
	return _BitcoinState.Contract.TargetInterfaces(&_BitcoinState.CallOpts)
}

// TargetInterfaces is a free data retrieval call binding the contract method 0x2ade3880.
//
// Solidity: function targetInterfaces() view returns((address,string[])[] targetedInterfaces_)
func (_BitcoinState *BitcoinStateCallerSession) TargetInterfaces() ([]StdInvariantFuzzInterface, error) {
	return _BitcoinState.Contract.TargetInterfaces(&_BitcoinState.CallOpts)
}

// TargetSelectors is a free data retrieval call binding the contract method 0x916a17c6.
//
// Solidity: function targetSelectors() view returns((address,bytes4[])[] targetedSelectors_)
func (_BitcoinState *BitcoinStateCaller) TargetSelectors(opts *bind.CallOpts) ([]StdInvariantFuzzSelector, error) {
	var out []interface{}
	err := _BitcoinState.contract.Call(opts, &out, "targetSelectors")

	if err != nil {
		return *new([]StdInvariantFuzzSelector), err
	}

	out0 := *abi.ConvertType(out[0], new([]StdInvariantFuzzSelector)).(*[]StdInvariantFuzzSelector)

	return out0, err

}

// TargetSelectors is a free data retrieval call binding the contract method 0x916a17c6.
//
// Solidity: function targetSelectors() view returns((address,bytes4[])[] targetedSelectors_)
func (_BitcoinState *BitcoinStateSession) TargetSelectors() ([]StdInvariantFuzzSelector, error) {
	return _BitcoinState.Contract.TargetSelectors(&_BitcoinState.CallOpts)
}

// TargetSelectors is a free data retrieval call binding the contract method 0x916a17c6.
//
// Solidity: function targetSelectors() view returns((address,bytes4[])[] targetedSelectors_)
func (_BitcoinState *BitcoinStateCallerSession) TargetSelectors() ([]StdInvariantFuzzSelector, error) {
	return _BitcoinState.Contract.TargetSelectors(&_BitcoinState.CallOpts)
}

// TargetSenders is a free data retrieval call binding the contract method 0x3e5e3c23.
//
// Solidity: function targetSenders() view returns(address[] targetedSenders_)
func (_BitcoinState *BitcoinStateCaller) TargetSenders(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _BitcoinState.contract.Call(opts, &out, "targetSenders")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// TargetSenders is a free data retrieval call binding the contract method 0x3e5e3c23.
//
// Solidity: function targetSenders() view returns(address[] targetedSenders_)
func (_BitcoinState *BitcoinStateSession) TargetSenders() ([]common.Address, error) {
	return _BitcoinState.Contract.TargetSenders(&_BitcoinState.CallOpts)
}

// TargetSenders is a free data retrieval call binding the contract method 0x3e5e3c23.
//
// Solidity: function targetSenders() view returns(address[] targetedSenders_)
func (_BitcoinState *BitcoinStateCallerSession) TargetSenders() ([]common.Address, error) {
	return _BitcoinState.Contract.TargetSenders(&_BitcoinState.CallOpts)
}

// Tokens is a free data retrieval call binding the contract method 0x904194a3.
//
// Solidity: function tokens(bytes32 ) view returns(address)
func (_BitcoinState *BitcoinStateCaller) Tokens(opts *bind.CallOpts, arg0 [32]byte) (common.Address, error) {
	var out []interface{}
	err := _BitcoinState.contract.Call(opts, &out, "tokens", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Tokens is a free data retrieval call binding the contract method 0x904194a3.
//
// Solidity: function tokens(bytes32 ) view returns(address)
func (_BitcoinState *BitcoinStateSession) Tokens(arg0 [32]byte) (common.Address, error) {
	return _BitcoinState.Contract.Tokens(&_BitcoinState.CallOpts, arg0)
}

// Tokens is a free data retrieval call binding the contract method 0x904194a3.
//
// Solidity: function tokens(bytes32 ) view returns(address)
func (_BitcoinState *BitcoinStateCallerSession) Tokens(arg0 [32]byte) (common.Address, error) {
	return _BitcoinState.Contract.Tokens(&_BitcoinState.CallOpts, arg0)
}

// XcallService is a free data retrieval call binding the contract method 0x7bf07164.
//
// Solidity: function xcallService() view returns(address)
func (_BitcoinState *BitcoinStateCaller) XcallService(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BitcoinState.contract.Call(opts, &out, "xcallService")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// XcallService is a free data retrieval call binding the contract method 0x7bf07164.
//
// Solidity: function xcallService() view returns(address)
func (_BitcoinState *BitcoinStateSession) XcallService() (common.Address, error) {
	return _BitcoinState.Contract.XcallService(&_BitcoinState.CallOpts)
}

// XcallService is a free data retrieval call binding the contract method 0x7bf07164.
//
// Solidity: function xcallService() view returns(address)
func (_BitcoinState *BitcoinStateCallerSession) XcallService() (common.Address, error) {
	return _BitcoinState.Contract.XcallService(&_BitcoinState.CallOpts)
}

// AddConnection is a paid mutator transaction binding the contract method 0x677dea1d.
//
// Solidity: function addConnection(address connection_) returns()
func (_BitcoinState *BitcoinStateTransactor) AddConnection(opts *bind.TransactOpts, connection_ common.Address) (*types.Transaction, error) {
	return _BitcoinState.contract.Transact(opts, "addConnection", connection_)
}

// AddConnection is a paid mutator transaction binding the contract method 0x677dea1d.
//
// Solidity: function addConnection(address connection_) returns()
func (_BitcoinState *BitcoinStateSession) AddConnection(connection_ common.Address) (*types.Transaction, error) {
	return _BitcoinState.Contract.AddConnection(&_BitcoinState.TransactOpts, connection_)
}

// AddConnection is a paid mutator transaction binding the contract method 0x677dea1d.
//
// Solidity: function addConnection(address connection_) returns()
func (_BitcoinState *BitcoinStateTransactorSession) AddConnection(connection_ common.Address) (*types.Transaction, error) {
	return _BitcoinState.Contract.AddConnection(&_BitcoinState.TransactOpts, connection_)
}

// GetSignData is a paid mutator transaction binding the contract method 0xbcc2f19d.
//
// Solidity: function getSignData(address requester_, bytes data_) returns(bytes32)
func (_BitcoinState *BitcoinStateTransactor) GetSignData(opts *bind.TransactOpts, requester_ common.Address, data_ []byte) (*types.Transaction, error) {
	return _BitcoinState.contract.Transact(opts, "getSignData", requester_, data_)
}

// GetSignData is a paid mutator transaction binding the contract method 0xbcc2f19d.
//
// Solidity: function getSignData(address requester_, bytes data_) returns(bytes32)
func (_BitcoinState *BitcoinStateSession) GetSignData(requester_ common.Address, data_ []byte) (*types.Transaction, error) {
	return _BitcoinState.Contract.GetSignData(&_BitcoinState.TransactOpts, requester_, data_)
}

// GetSignData is a paid mutator transaction binding the contract method 0xbcc2f19d.
//
// Solidity: function getSignData(address requester_, bytes data_) returns(bytes32)
func (_BitcoinState *BitcoinStateTransactorSession) GetSignData(requester_ common.Address, data_ []byte) (*types.Transaction, error) {
	return _BitcoinState.Contract.GetSignData(&_BitcoinState.TransactOpts, requester_, data_)
}

// HandleCallMessage is a paid mutator transaction binding the contract method 0x5d6a16f5.
//
// Solidity: function handleCallMessage(string _from, bytes _data, string[] _protocols) returns()
func (_BitcoinState *BitcoinStateTransactor) HandleCallMessage(opts *bind.TransactOpts, _from string, _data []byte, _protocols []string) (*types.Transaction, error) {
	return _BitcoinState.contract.Transact(opts, "handleCallMessage", _from, _data, _protocols)
}

// HandleCallMessage is a paid mutator transaction binding the contract method 0x5d6a16f5.
//
// Solidity: function handleCallMessage(string _from, bytes _data, string[] _protocols) returns()
func (_BitcoinState *BitcoinStateSession) HandleCallMessage(_from string, _data []byte, _protocols []string) (*types.Transaction, error) {
	return _BitcoinState.Contract.HandleCallMessage(&_BitcoinState.TransactOpts, _from, _data, _protocols)
}

// HandleCallMessage is a paid mutator transaction binding the contract method 0x5d6a16f5.
//
// Solidity: function handleCallMessage(string _from, bytes _data, string[] _protocols) returns()
func (_BitcoinState *BitcoinStateTransactorSession) HandleCallMessage(_from string, _data []byte, _protocols []string) (*types.Transaction, error) {
	return _BitcoinState.Contract.HandleCallMessage(&_BitcoinState.TransactOpts, _from, _data, _protocols)
}

// Initialize is a paid mutator transaction binding the contract method 0xe6bfbfd8.
//
// Solidity: function initialize(address xcall_, address uinswapV3Router_, address nonfungiblePositionManager_, address[] connections) returns()
func (_BitcoinState *BitcoinStateTransactor) Initialize(opts *bind.TransactOpts, xcall_ common.Address, uinswapV3Router_ common.Address, nonfungiblePositionManager_ common.Address, connections []common.Address) (*types.Transaction, error) {
	return _BitcoinState.contract.Transact(opts, "initialize", xcall_, uinswapV3Router_, nonfungiblePositionManager_, connections)
}

// Initialize is a paid mutator transaction binding the contract method 0xe6bfbfd8.
//
// Solidity: function initialize(address xcall_, address uinswapV3Router_, address nonfungiblePositionManager_, address[] connections) returns()
func (_BitcoinState *BitcoinStateSession) Initialize(xcall_ common.Address, uinswapV3Router_ common.Address, nonfungiblePositionManager_ common.Address, connections []common.Address) (*types.Transaction, error) {
	return _BitcoinState.Contract.Initialize(&_BitcoinState.TransactOpts, xcall_, uinswapV3Router_, nonfungiblePositionManager_, connections)
}

// Initialize is a paid mutator transaction binding the contract method 0xe6bfbfd8.
//
// Solidity: function initialize(address xcall_, address uinswapV3Router_, address nonfungiblePositionManager_, address[] connections) returns()
func (_BitcoinState *BitcoinStateTransactorSession) Initialize(xcall_ common.Address, uinswapV3Router_ common.Address, nonfungiblePositionManager_ common.Address, connections []common.Address) (*types.Transaction, error) {
	return _BitcoinState.Contract.Initialize(&_BitcoinState.TransactOpts, xcall_, uinswapV3Router_, nonfungiblePositionManager_, connections)
}

// Migrate is a paid mutator transaction binding the contract method 0x8932a90d.
//
// Solidity: function migrate(bytes _data) returns()
func (_BitcoinState *BitcoinStateTransactor) Migrate(opts *bind.TransactOpts, _data []byte) (*types.Transaction, error) {
	return _BitcoinState.contract.Transact(opts, "migrate", _data)
}

// Migrate is a paid mutator transaction binding the contract method 0x8932a90d.
//
// Solidity: function migrate(bytes _data) returns()
func (_BitcoinState *BitcoinStateSession) Migrate(_data []byte) (*types.Transaction, error) {
	return _BitcoinState.Contract.Migrate(&_BitcoinState.TransactOpts, _data)
}

// Migrate is a paid mutator transaction binding the contract method 0x8932a90d.
//
// Solidity: function migrate(bytes _data) returns()
func (_BitcoinState *BitcoinStateTransactorSession) Migrate(_data []byte) (*types.Transaction, error) {
	return _BitcoinState.Contract.Migrate(&_BitcoinState.TransactOpts, _data)
}

// MigrateComplete is a paid mutator transaction binding the contract method 0xf0ad3762.
//
// Solidity: function migrateComplete() returns()
func (_BitcoinState *BitcoinStateTransactor) MigrateComplete(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BitcoinState.contract.Transact(opts, "migrateComplete")
}

// MigrateComplete is a paid mutator transaction binding the contract method 0xf0ad3762.
//
// Solidity: function migrateComplete() returns()
func (_BitcoinState *BitcoinStateSession) MigrateComplete() (*types.Transaction, error) {
	return _BitcoinState.Contract.MigrateComplete(&_BitcoinState.TransactOpts)
}

// MigrateComplete is a paid mutator transaction binding the contract method 0xf0ad3762.
//
// Solidity: function migrateComplete() returns()
func (_BitcoinState *BitcoinStateTransactorSession) MigrateComplete() (*types.Transaction, error) {
	return _BitcoinState.Contract.MigrateComplete(&_BitcoinState.TransactOpts)
}

// RemoveConnection is a paid mutator transaction binding the contract method 0x65301f0d.
//
// Solidity: function removeConnection(address connection_) returns()
func (_BitcoinState *BitcoinStateTransactor) RemoveConnection(opts *bind.TransactOpts, connection_ common.Address) (*types.Transaction, error) {
	return _BitcoinState.contract.Transact(opts, "removeConnection", connection_)
}

// RemoveConnection is a paid mutator transaction binding the contract method 0x65301f0d.
//
// Solidity: function removeConnection(address connection_) returns()
func (_BitcoinState *BitcoinStateSession) RemoveConnection(connection_ common.Address) (*types.Transaction, error) {
	return _BitcoinState.Contract.RemoveConnection(&_BitcoinState.TransactOpts, connection_)
}

// RemoveConnection is a paid mutator transaction binding the contract method 0x65301f0d.
//
// Solidity: function removeConnection(address connection_) returns()
func (_BitcoinState *BitcoinStateTransactorSession) RemoveConnection(connection_ common.Address) (*types.Transaction, error) {
	return _BitcoinState.Contract.RemoveConnection(&_BitcoinState.TransactOpts, connection_)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BitcoinState *BitcoinStateTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BitcoinState.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BitcoinState *BitcoinStateSession) RenounceOwnership() (*types.Transaction, error) {
	return _BitcoinState.Contract.RenounceOwnership(&_BitcoinState.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BitcoinState *BitcoinStateTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _BitcoinState.Contract.RenounceOwnership(&_BitcoinState.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_BitcoinState *BitcoinStateTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _BitcoinState.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_BitcoinState *BitcoinStateSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _BitcoinState.Contract.TransferOwnership(&_BitcoinState.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_BitcoinState *BitcoinStateTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _BitcoinState.Contract.TransferOwnership(&_BitcoinState.TransactOpts, newOwner)
}

// BitcoinStateAddConnectionIterator is returned from FilterAddConnection and is used to iterate over the raw logs and unpacked data for AddConnection events raised by the BitcoinState contract.
type BitcoinStateAddConnectionIterator struct {
	Event *BitcoinStateAddConnection // Event containing the contract specifics and raw log

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
func (it *BitcoinStateAddConnectionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BitcoinStateAddConnection)
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
		it.Event = new(BitcoinStateAddConnection)
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
func (it *BitcoinStateAddConnectionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BitcoinStateAddConnectionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BitcoinStateAddConnection represents a AddConnection event raised by the BitcoinState contract.
type BitcoinStateAddConnection struct {
	Connection common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterAddConnection is a free log retrieval operation binding the contract event 0x5cd3d8da8ef00a8b5228348fe7683dae605751d4867ba7ea9fdc8260b9e2c7d3.
//
// Solidity: event AddConnection(address connection_)
func (_BitcoinState *BitcoinStateFilterer) FilterAddConnection(opts *bind.FilterOpts) (*BitcoinStateAddConnectionIterator, error) {

	logs, sub, err := _BitcoinState.contract.FilterLogs(opts, "AddConnection")
	if err != nil {
		return nil, err
	}
	return &BitcoinStateAddConnectionIterator{contract: _BitcoinState.contract, event: "AddConnection", logs: logs, sub: sub}, nil
}

// WatchAddConnection is a free log subscription operation binding the contract event 0x5cd3d8da8ef00a8b5228348fe7683dae605751d4867ba7ea9fdc8260b9e2c7d3.
//
// Solidity: event AddConnection(address connection_)
func (_BitcoinState *BitcoinStateFilterer) WatchAddConnection(opts *bind.WatchOpts, sink chan<- *BitcoinStateAddConnection) (event.Subscription, error) {

	logs, sub, err := _BitcoinState.contract.WatchLogs(opts, "AddConnection")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BitcoinStateAddConnection)
				if err := _BitcoinState.contract.UnpackLog(event, "AddConnection", log); err != nil {
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

// ParseAddConnection is a log parse operation binding the contract event 0x5cd3d8da8ef00a8b5228348fe7683dae605751d4867ba7ea9fdc8260b9e2c7d3.
//
// Solidity: event AddConnection(address connection_)
func (_BitcoinState *BitcoinStateFilterer) ParseAddConnection(log types.Log) (*BitcoinStateAddConnection, error) {
	event := new(BitcoinStateAddConnection)
	if err := _BitcoinState.contract.UnpackLog(event, "AddConnection", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BitcoinStateAddSelectorIterator is returned from FilterAddSelector and is used to iterate over the raw logs and unpacked data for AddSelector events raised by the BitcoinState contract.
type BitcoinStateAddSelectorIterator struct {
	Event *BitcoinStateAddSelector // Event containing the contract specifics and raw log

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
func (it *BitcoinStateAddSelectorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BitcoinStateAddSelector)
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
		it.Event = new(BitcoinStateAddSelector)
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
func (it *BitcoinStateAddSelectorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BitcoinStateAddSelectorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BitcoinStateAddSelector represents a AddSelector event raised by the BitcoinState contract.
type BitcoinStateAddSelector struct {
	Selector  [4]byte
	Recipient common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAddSelector is a free log retrieval operation binding the contract event 0x95c0d09571a33725bdfaf52708354735c2650abb90fe193d20429e324a2b3696.
//
// Solidity: event AddSelector(bytes4 selector_, address recipient)
func (_BitcoinState *BitcoinStateFilterer) FilterAddSelector(opts *bind.FilterOpts) (*BitcoinStateAddSelectorIterator, error) {

	logs, sub, err := _BitcoinState.contract.FilterLogs(opts, "AddSelector")
	if err != nil {
		return nil, err
	}
	return &BitcoinStateAddSelectorIterator{contract: _BitcoinState.contract, event: "AddSelector", logs: logs, sub: sub}, nil
}

// WatchAddSelector is a free log subscription operation binding the contract event 0x95c0d09571a33725bdfaf52708354735c2650abb90fe193d20429e324a2b3696.
//
// Solidity: event AddSelector(bytes4 selector_, address recipient)
func (_BitcoinState *BitcoinStateFilterer) WatchAddSelector(opts *bind.WatchOpts, sink chan<- *BitcoinStateAddSelector) (event.Subscription, error) {

	logs, sub, err := _BitcoinState.contract.WatchLogs(opts, "AddSelector")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BitcoinStateAddSelector)
				if err := _BitcoinState.contract.UnpackLog(event, "AddSelector", log); err != nil {
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

// ParseAddSelector is a log parse operation binding the contract event 0x95c0d09571a33725bdfaf52708354735c2650abb90fe193d20429e324a2b3696.
//
// Solidity: event AddSelector(bytes4 selector_, address recipient)
func (_BitcoinState *BitcoinStateFilterer) ParseAddSelector(log types.Log) (*BitcoinStateAddSelector, error) {
	event := new(BitcoinStateAddSelector)
	if err := _BitcoinState.contract.UnpackLog(event, "AddSelector", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BitcoinStateInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the BitcoinState contract.
type BitcoinStateInitializedIterator struct {
	Event *BitcoinStateInitialized // Event containing the contract specifics and raw log

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
func (it *BitcoinStateInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BitcoinStateInitialized)
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
		it.Event = new(BitcoinStateInitialized)
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
func (it *BitcoinStateInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BitcoinStateInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BitcoinStateInitialized represents a Initialized event raised by the BitcoinState contract.
type BitcoinStateInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_BitcoinState *BitcoinStateFilterer) FilterInitialized(opts *bind.FilterOpts) (*BitcoinStateInitializedIterator, error) {

	logs, sub, err := _BitcoinState.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &BitcoinStateInitializedIterator{contract: _BitcoinState.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_BitcoinState *BitcoinStateFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *BitcoinStateInitialized) (event.Subscription, error) {

	logs, sub, err := _BitcoinState.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BitcoinStateInitialized)
				if err := _BitcoinState.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_BitcoinState *BitcoinStateFilterer) ParseInitialized(log types.Log) (*BitcoinStateInitialized, error) {
	event := new(BitcoinStateInitialized)
	if err := _BitcoinState.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BitcoinStateOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the BitcoinState contract.
type BitcoinStateOwnershipTransferredIterator struct {
	Event *BitcoinStateOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *BitcoinStateOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BitcoinStateOwnershipTransferred)
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
		it.Event = new(BitcoinStateOwnershipTransferred)
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
func (it *BitcoinStateOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BitcoinStateOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BitcoinStateOwnershipTransferred represents a OwnershipTransferred event raised by the BitcoinState contract.
type BitcoinStateOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_BitcoinState *BitcoinStateFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*BitcoinStateOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _BitcoinState.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &BitcoinStateOwnershipTransferredIterator{contract: _BitcoinState.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_BitcoinState *BitcoinStateFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BitcoinStateOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _BitcoinState.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BitcoinStateOwnershipTransferred)
				if err := _BitcoinState.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_BitcoinState *BitcoinStateFilterer) ParseOwnershipTransferred(log types.Log) (*BitcoinStateOwnershipTransferred, error) {
	event := new(BitcoinStateOwnershipTransferred)
	if err := _BitcoinState.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BitcoinStateRemoveConnectionIterator is returned from FilterRemoveConnection and is used to iterate over the raw logs and unpacked data for RemoveConnection events raised by the BitcoinState contract.
type BitcoinStateRemoveConnectionIterator struct {
	Event *BitcoinStateRemoveConnection // Event containing the contract specifics and raw log

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
func (it *BitcoinStateRemoveConnectionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BitcoinStateRemoveConnection)
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
		it.Event = new(BitcoinStateRemoveConnection)
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
func (it *BitcoinStateRemoveConnectionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BitcoinStateRemoveConnectionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BitcoinStateRemoveConnection represents a RemoveConnection event raised by the BitcoinState contract.
type BitcoinStateRemoveConnection struct {
	Connection common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterRemoveConnection is a free log retrieval operation binding the contract event 0xace8d11a44b7aa536cc46a77b519166c001adc485ba8cfa404e1aa252b07db38.
//
// Solidity: event RemoveConnection(address connection_)
func (_BitcoinState *BitcoinStateFilterer) FilterRemoveConnection(opts *bind.FilterOpts) (*BitcoinStateRemoveConnectionIterator, error) {

	logs, sub, err := _BitcoinState.contract.FilterLogs(opts, "RemoveConnection")
	if err != nil {
		return nil, err
	}
	return &BitcoinStateRemoveConnectionIterator{contract: _BitcoinState.contract, event: "RemoveConnection", logs: logs, sub: sub}, nil
}

// WatchRemoveConnection is a free log subscription operation binding the contract event 0xace8d11a44b7aa536cc46a77b519166c001adc485ba8cfa404e1aa252b07db38.
//
// Solidity: event RemoveConnection(address connection_)
func (_BitcoinState *BitcoinStateFilterer) WatchRemoveConnection(opts *bind.WatchOpts, sink chan<- *BitcoinStateRemoveConnection) (event.Subscription, error) {

	logs, sub, err := _BitcoinState.contract.WatchLogs(opts, "RemoveConnection")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BitcoinStateRemoveConnection)
				if err := _BitcoinState.contract.UnpackLog(event, "RemoveConnection", log); err != nil {
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

// ParseRemoveConnection is a log parse operation binding the contract event 0xace8d11a44b7aa536cc46a77b519166c001adc485ba8cfa404e1aa252b07db38.
//
// Solidity: event RemoveConnection(address connection_)
func (_BitcoinState *BitcoinStateFilterer) ParseRemoveConnection(log types.Log) (*BitcoinStateRemoveConnection, error) {
	event := new(BitcoinStateRemoveConnection)
	if err := _BitcoinState.contract.UnpackLog(event, "RemoveConnection", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BitcoinStateRemoveSelectorIterator is returned from FilterRemoveSelector and is used to iterate over the raw logs and unpacked data for RemoveSelector events raised by the BitcoinState contract.
type BitcoinStateRemoveSelectorIterator struct {
	Event *BitcoinStateRemoveSelector // Event containing the contract specifics and raw log

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
func (it *BitcoinStateRemoveSelectorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BitcoinStateRemoveSelector)
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
		it.Event = new(BitcoinStateRemoveSelector)
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
func (it *BitcoinStateRemoveSelectorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BitcoinStateRemoveSelectorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BitcoinStateRemoveSelector represents a RemoveSelector event raised by the BitcoinState contract.
type BitcoinStateRemoveSelector struct {
	Selector [4]byte
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterRemoveSelector is a free log retrieval operation binding the contract event 0x85a48e474e38192938033da06ebe84a59fb194958d1be091ac2fa9a6630da31f.
//
// Solidity: event RemoveSelector(bytes4 selector_)
func (_BitcoinState *BitcoinStateFilterer) FilterRemoveSelector(opts *bind.FilterOpts) (*BitcoinStateRemoveSelectorIterator, error) {

	logs, sub, err := _BitcoinState.contract.FilterLogs(opts, "RemoveSelector")
	if err != nil {
		return nil, err
	}
	return &BitcoinStateRemoveSelectorIterator{contract: _BitcoinState.contract, event: "RemoveSelector", logs: logs, sub: sub}, nil
}

// WatchRemoveSelector is a free log subscription operation binding the contract event 0x85a48e474e38192938033da06ebe84a59fb194958d1be091ac2fa9a6630da31f.
//
// Solidity: event RemoveSelector(bytes4 selector_)
func (_BitcoinState *BitcoinStateFilterer) WatchRemoveSelector(opts *bind.WatchOpts, sink chan<- *BitcoinStateRemoveSelector) (event.Subscription, error) {

	logs, sub, err := _BitcoinState.contract.WatchLogs(opts, "RemoveSelector")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BitcoinStateRemoveSelector)
				if err := _BitcoinState.contract.UnpackLog(event, "RemoveSelector", log); err != nil {
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

// ParseRemoveSelector is a log parse operation binding the contract event 0x85a48e474e38192938033da06ebe84a59fb194958d1be091ac2fa9a6630da31f.
//
// Solidity: event RemoveSelector(bytes4 selector_)
func (_BitcoinState *BitcoinStateFilterer) ParseRemoveSelector(log types.Log) (*BitcoinStateRemoveSelector, error) {
	event := new(BitcoinStateRemoveSelector)
	if err := _BitcoinState.contract.UnpackLog(event, "RemoveSelector", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BitcoinStateRequestExecutedIterator is returned from FilterRequestExecuted and is used to iterate over the raw logs and unpacked data for RequestExecuted events raised by the BitcoinState contract.
type BitcoinStateRequestExecutedIterator struct {
	Event *BitcoinStateRequestExecuted // Event containing the contract specifics and raw log

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
func (it *BitcoinStateRequestExecutedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BitcoinStateRequestExecuted)
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
		it.Event = new(BitcoinStateRequestExecuted)
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
func (it *BitcoinStateRequestExecutedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BitcoinStateRequestExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BitcoinStateRequestExecuted represents a RequestExecuted event raised by the BitcoinState contract.
type BitcoinStateRequestExecuted struct {
	Id        *big.Int
	StateRoot [32]byte
	Data      []byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRequestExecuted is a free log retrieval operation binding the contract event 0x9c343316d67a8e28446ef883ab491ece3ff70d3eeaa9fbd13a362a0afd690721.
//
// Solidity: event RequestExecuted(uint256 id, bytes32 stateRoot, bytes data)
func (_BitcoinState *BitcoinStateFilterer) FilterRequestExecuted(opts *bind.FilterOpts) (*BitcoinStateRequestExecutedIterator, error) {

	logs, sub, err := _BitcoinState.contract.FilterLogs(opts, "RequestExecuted")
	if err != nil {
		return nil, err
	}
	return &BitcoinStateRequestExecutedIterator{contract: _BitcoinState.contract, event: "RequestExecuted", logs: logs, sub: sub}, nil
}

// WatchRequestExecuted is a free log subscription operation binding the contract event 0x9c343316d67a8e28446ef883ab491ece3ff70d3eeaa9fbd13a362a0afd690721.
//
// Solidity: event RequestExecuted(uint256 id, bytes32 stateRoot, bytes data)
func (_BitcoinState *BitcoinStateFilterer) WatchRequestExecuted(opts *bind.WatchOpts, sink chan<- *BitcoinStateRequestExecuted) (event.Subscription, error) {

	logs, sub, err := _BitcoinState.contract.WatchLogs(opts, "RequestExecuted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BitcoinStateRequestExecuted)
				if err := _BitcoinState.contract.UnpackLog(event, "RequestExecuted", log); err != nil {
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

// ParseRequestExecuted is a log parse operation binding the contract event 0x9c343316d67a8e28446ef883ab491ece3ff70d3eeaa9fbd13a362a0afd690721.
//
// Solidity: event RequestExecuted(uint256 id, bytes32 stateRoot, bytes data)
func (_BitcoinState *BitcoinStateFilterer) ParseRequestExecuted(log types.Log) (*BitcoinStateRequestExecuted, error) {
	event := new(BitcoinStateRequestExecuted)
	if err := _BitcoinState.contract.UnpackLog(event, "RequestExecuted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BitcoinStateLogIterator is returned from FilterLog and is used to iterate over the raw logs and unpacked data for Log events raised by the BitcoinState contract.
type BitcoinStateLogIterator struct {
	Event *BitcoinStateLog // Event containing the contract specifics and raw log

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
func (it *BitcoinStateLogIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BitcoinStateLog)
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
		it.Event = new(BitcoinStateLog)
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
func (it *BitcoinStateLogIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BitcoinStateLogIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BitcoinStateLog represents a Log event raised by the BitcoinState contract.
type BitcoinStateLog struct {
	Arg0 string
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLog is a free log retrieval operation binding the contract event 0x41304facd9323d75b11bcdd609cb38effffdb05710f7caf0e9b16c6d9d709f50.
//
// Solidity: event log(string arg0)
func (_BitcoinState *BitcoinStateFilterer) FilterLog(opts *bind.FilterOpts) (*BitcoinStateLogIterator, error) {

	logs, sub, err := _BitcoinState.contract.FilterLogs(opts, "log")
	if err != nil {
		return nil, err
	}
	return &BitcoinStateLogIterator{contract: _BitcoinState.contract, event: "log", logs: logs, sub: sub}, nil
}

// WatchLog is a free log subscription operation binding the contract event 0x41304facd9323d75b11bcdd609cb38effffdb05710f7caf0e9b16c6d9d709f50.
//
// Solidity: event log(string arg0)
func (_BitcoinState *BitcoinStateFilterer) WatchLog(opts *bind.WatchOpts, sink chan<- *BitcoinStateLog) (event.Subscription, error) {

	logs, sub, err := _BitcoinState.contract.WatchLogs(opts, "log")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BitcoinStateLog)
				if err := _BitcoinState.contract.UnpackLog(event, "log", log); err != nil {
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

// ParseLog is a log parse operation binding the contract event 0x41304facd9323d75b11bcdd609cb38effffdb05710f7caf0e9b16c6d9d709f50.
//
// Solidity: event log(string arg0)
func (_BitcoinState *BitcoinStateFilterer) ParseLog(log types.Log) (*BitcoinStateLog, error) {
	event := new(BitcoinStateLog)
	if err := _BitcoinState.contract.UnpackLog(event, "log", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BitcoinStateLogAddressIterator is returned from FilterLogAddress and is used to iterate over the raw logs and unpacked data for LogAddress events raised by the BitcoinState contract.
type BitcoinStateLogAddressIterator struct {
	Event *BitcoinStateLogAddress // Event containing the contract specifics and raw log

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
func (it *BitcoinStateLogAddressIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BitcoinStateLogAddress)
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
		it.Event = new(BitcoinStateLogAddress)
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
func (it *BitcoinStateLogAddressIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BitcoinStateLogAddressIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BitcoinStateLogAddress represents a LogAddress event raised by the BitcoinState contract.
type BitcoinStateLogAddress struct {
	Arg0 common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogAddress is a free log retrieval operation binding the contract event 0x7ae74c527414ae135fd97047b12921a5ec3911b804197855d67e25c7b75ee6f3.
//
// Solidity: event log_address(address arg0)
func (_BitcoinState *BitcoinStateFilterer) FilterLogAddress(opts *bind.FilterOpts) (*BitcoinStateLogAddressIterator, error) {

	logs, sub, err := _BitcoinState.contract.FilterLogs(opts, "log_address")
	if err != nil {
		return nil, err
	}
	return &BitcoinStateLogAddressIterator{contract: _BitcoinState.contract, event: "log_address", logs: logs, sub: sub}, nil
}

// WatchLogAddress is a free log subscription operation binding the contract event 0x7ae74c527414ae135fd97047b12921a5ec3911b804197855d67e25c7b75ee6f3.
//
// Solidity: event log_address(address arg0)
func (_BitcoinState *BitcoinStateFilterer) WatchLogAddress(opts *bind.WatchOpts, sink chan<- *BitcoinStateLogAddress) (event.Subscription, error) {

	logs, sub, err := _BitcoinState.contract.WatchLogs(opts, "log_address")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BitcoinStateLogAddress)
				if err := _BitcoinState.contract.UnpackLog(event, "log_address", log); err != nil {
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

// ParseLogAddress is a log parse operation binding the contract event 0x7ae74c527414ae135fd97047b12921a5ec3911b804197855d67e25c7b75ee6f3.
//
// Solidity: event log_address(address arg0)
func (_BitcoinState *BitcoinStateFilterer) ParseLogAddress(log types.Log) (*BitcoinStateLogAddress, error) {
	event := new(BitcoinStateLogAddress)
	if err := _BitcoinState.contract.UnpackLog(event, "log_address", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BitcoinStateLogArrayIterator is returned from FilterLogArray and is used to iterate over the raw logs and unpacked data for LogArray events raised by the BitcoinState contract.
type BitcoinStateLogArrayIterator struct {
	Event *BitcoinStateLogArray // Event containing the contract specifics and raw log

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
func (it *BitcoinStateLogArrayIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BitcoinStateLogArray)
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
		it.Event = new(BitcoinStateLogArray)
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
func (it *BitcoinStateLogArrayIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BitcoinStateLogArrayIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BitcoinStateLogArray represents a LogArray event raised by the BitcoinState contract.
type BitcoinStateLogArray struct {
	Val []*big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogArray is a free log retrieval operation binding the contract event 0xfb102865d50addddf69da9b5aa1bced66c80cf869a5c8d0471a467e18ce9cab1.
//
// Solidity: event log_array(uint256[] val)
func (_BitcoinState *BitcoinStateFilterer) FilterLogArray(opts *bind.FilterOpts) (*BitcoinStateLogArrayIterator, error) {

	logs, sub, err := _BitcoinState.contract.FilterLogs(opts, "log_array")
	if err != nil {
		return nil, err
	}
	return &BitcoinStateLogArrayIterator{contract: _BitcoinState.contract, event: "log_array", logs: logs, sub: sub}, nil
}

// WatchLogArray is a free log subscription operation binding the contract event 0xfb102865d50addddf69da9b5aa1bced66c80cf869a5c8d0471a467e18ce9cab1.
//
// Solidity: event log_array(uint256[] val)
func (_BitcoinState *BitcoinStateFilterer) WatchLogArray(opts *bind.WatchOpts, sink chan<- *BitcoinStateLogArray) (event.Subscription, error) {

	logs, sub, err := _BitcoinState.contract.WatchLogs(opts, "log_array")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BitcoinStateLogArray)
				if err := _BitcoinState.contract.UnpackLog(event, "log_array", log); err != nil {
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

// ParseLogArray is a log parse operation binding the contract event 0xfb102865d50addddf69da9b5aa1bced66c80cf869a5c8d0471a467e18ce9cab1.
//
// Solidity: event log_array(uint256[] val)
func (_BitcoinState *BitcoinStateFilterer) ParseLogArray(log types.Log) (*BitcoinStateLogArray, error) {
	event := new(BitcoinStateLogArray)
	if err := _BitcoinState.contract.UnpackLog(event, "log_array", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BitcoinStateLogArray0Iterator is returned from FilterLogArray0 and is used to iterate over the raw logs and unpacked data for LogArray0 events raised by the BitcoinState contract.
type BitcoinStateLogArray0Iterator struct {
	Event *BitcoinStateLogArray0 // Event containing the contract specifics and raw log

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
func (it *BitcoinStateLogArray0Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BitcoinStateLogArray0)
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
		it.Event = new(BitcoinStateLogArray0)
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
func (it *BitcoinStateLogArray0Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BitcoinStateLogArray0Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BitcoinStateLogArray0 represents a LogArray0 event raised by the BitcoinState contract.
type BitcoinStateLogArray0 struct {
	Val []*big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogArray0 is a free log retrieval operation binding the contract event 0x890a82679b470f2bd82816ed9b161f97d8b967f37fa3647c21d5bf39749e2dd5.
//
// Solidity: event log_array(int256[] val)
func (_BitcoinState *BitcoinStateFilterer) FilterLogArray0(opts *bind.FilterOpts) (*BitcoinStateLogArray0Iterator, error) {

	logs, sub, err := _BitcoinState.contract.FilterLogs(opts, "log_array0")
	if err != nil {
		return nil, err
	}
	return &BitcoinStateLogArray0Iterator{contract: _BitcoinState.contract, event: "log_array0", logs: logs, sub: sub}, nil
}

// WatchLogArray0 is a free log subscription operation binding the contract event 0x890a82679b470f2bd82816ed9b161f97d8b967f37fa3647c21d5bf39749e2dd5.
//
// Solidity: event log_array(int256[] val)
func (_BitcoinState *BitcoinStateFilterer) WatchLogArray0(opts *bind.WatchOpts, sink chan<- *BitcoinStateLogArray0) (event.Subscription, error) {

	logs, sub, err := _BitcoinState.contract.WatchLogs(opts, "log_array0")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BitcoinStateLogArray0)
				if err := _BitcoinState.contract.UnpackLog(event, "log_array0", log); err != nil {
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

// ParseLogArray0 is a log parse operation binding the contract event 0x890a82679b470f2bd82816ed9b161f97d8b967f37fa3647c21d5bf39749e2dd5.
//
// Solidity: event log_array(int256[] val)
func (_BitcoinState *BitcoinStateFilterer) ParseLogArray0(log types.Log) (*BitcoinStateLogArray0, error) {
	event := new(BitcoinStateLogArray0)
	if err := _BitcoinState.contract.UnpackLog(event, "log_array0", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BitcoinStateLogArray1Iterator is returned from FilterLogArray1 and is used to iterate over the raw logs and unpacked data for LogArray1 events raised by the BitcoinState contract.
type BitcoinStateLogArray1Iterator struct {
	Event *BitcoinStateLogArray1 // Event containing the contract specifics and raw log

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
func (it *BitcoinStateLogArray1Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BitcoinStateLogArray1)
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
		it.Event = new(BitcoinStateLogArray1)
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
func (it *BitcoinStateLogArray1Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BitcoinStateLogArray1Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BitcoinStateLogArray1 represents a LogArray1 event raised by the BitcoinState contract.
type BitcoinStateLogArray1 struct {
	Val []common.Address
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogArray1 is a free log retrieval operation binding the contract event 0x40e1840f5769073d61bd01372d9b75baa9842d5629a0c99ff103be1178a8e9e2.
//
// Solidity: event log_array(address[] val)
func (_BitcoinState *BitcoinStateFilterer) FilterLogArray1(opts *bind.FilterOpts) (*BitcoinStateLogArray1Iterator, error) {

	logs, sub, err := _BitcoinState.contract.FilterLogs(opts, "log_array1")
	if err != nil {
		return nil, err
	}
	return &BitcoinStateLogArray1Iterator{contract: _BitcoinState.contract, event: "log_array1", logs: logs, sub: sub}, nil
}

// WatchLogArray1 is a free log subscription operation binding the contract event 0x40e1840f5769073d61bd01372d9b75baa9842d5629a0c99ff103be1178a8e9e2.
//
// Solidity: event log_array(address[] val)
func (_BitcoinState *BitcoinStateFilterer) WatchLogArray1(opts *bind.WatchOpts, sink chan<- *BitcoinStateLogArray1) (event.Subscription, error) {

	logs, sub, err := _BitcoinState.contract.WatchLogs(opts, "log_array1")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BitcoinStateLogArray1)
				if err := _BitcoinState.contract.UnpackLog(event, "log_array1", log); err != nil {
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

// ParseLogArray1 is a log parse operation binding the contract event 0x40e1840f5769073d61bd01372d9b75baa9842d5629a0c99ff103be1178a8e9e2.
//
// Solidity: event log_array(address[] val)
func (_BitcoinState *BitcoinStateFilterer) ParseLogArray1(log types.Log) (*BitcoinStateLogArray1, error) {
	event := new(BitcoinStateLogArray1)
	if err := _BitcoinState.contract.UnpackLog(event, "log_array1", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BitcoinStateLogBytesIterator is returned from FilterLogBytes and is used to iterate over the raw logs and unpacked data for LogBytes events raised by the BitcoinState contract.
type BitcoinStateLogBytesIterator struct {
	Event *BitcoinStateLogBytes // Event containing the contract specifics and raw log

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
func (it *BitcoinStateLogBytesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BitcoinStateLogBytes)
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
		it.Event = new(BitcoinStateLogBytes)
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
func (it *BitcoinStateLogBytesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BitcoinStateLogBytesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BitcoinStateLogBytes represents a LogBytes event raised by the BitcoinState contract.
type BitcoinStateLogBytes struct {
	Arg0 []byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogBytes is a free log retrieval operation binding the contract event 0x23b62ad0584d24a75f0bf3560391ef5659ec6db1269c56e11aa241d637f19b20.
//
// Solidity: event log_bytes(bytes arg0)
func (_BitcoinState *BitcoinStateFilterer) FilterLogBytes(opts *bind.FilterOpts) (*BitcoinStateLogBytesIterator, error) {

	logs, sub, err := _BitcoinState.contract.FilterLogs(opts, "log_bytes")
	if err != nil {
		return nil, err
	}
	return &BitcoinStateLogBytesIterator{contract: _BitcoinState.contract, event: "log_bytes", logs: logs, sub: sub}, nil
}

// WatchLogBytes is a free log subscription operation binding the contract event 0x23b62ad0584d24a75f0bf3560391ef5659ec6db1269c56e11aa241d637f19b20.
//
// Solidity: event log_bytes(bytes arg0)
func (_BitcoinState *BitcoinStateFilterer) WatchLogBytes(opts *bind.WatchOpts, sink chan<- *BitcoinStateLogBytes) (event.Subscription, error) {

	logs, sub, err := _BitcoinState.contract.WatchLogs(opts, "log_bytes")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BitcoinStateLogBytes)
				if err := _BitcoinState.contract.UnpackLog(event, "log_bytes", log); err != nil {
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

// ParseLogBytes is a log parse operation binding the contract event 0x23b62ad0584d24a75f0bf3560391ef5659ec6db1269c56e11aa241d637f19b20.
//
// Solidity: event log_bytes(bytes arg0)
func (_BitcoinState *BitcoinStateFilterer) ParseLogBytes(log types.Log) (*BitcoinStateLogBytes, error) {
	event := new(BitcoinStateLogBytes)
	if err := _BitcoinState.contract.UnpackLog(event, "log_bytes", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BitcoinStateLogBytes32Iterator is returned from FilterLogBytes32 and is used to iterate over the raw logs and unpacked data for LogBytes32 events raised by the BitcoinState contract.
type BitcoinStateLogBytes32Iterator struct {
	Event *BitcoinStateLogBytes32 // Event containing the contract specifics and raw log

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
func (it *BitcoinStateLogBytes32Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BitcoinStateLogBytes32)
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
		it.Event = new(BitcoinStateLogBytes32)
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
func (it *BitcoinStateLogBytes32Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BitcoinStateLogBytes32Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BitcoinStateLogBytes32 represents a LogBytes32 event raised by the BitcoinState contract.
type BitcoinStateLogBytes32 struct {
	Arg0 [32]byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogBytes32 is a free log retrieval operation binding the contract event 0xe81699b85113eea1c73e10588b2b035e55893369632173afd43feb192fac64e3.
//
// Solidity: event log_bytes32(bytes32 arg0)
func (_BitcoinState *BitcoinStateFilterer) FilterLogBytes32(opts *bind.FilterOpts) (*BitcoinStateLogBytes32Iterator, error) {

	logs, sub, err := _BitcoinState.contract.FilterLogs(opts, "log_bytes32")
	if err != nil {
		return nil, err
	}
	return &BitcoinStateLogBytes32Iterator{contract: _BitcoinState.contract, event: "log_bytes32", logs: logs, sub: sub}, nil
}

// WatchLogBytes32 is a free log subscription operation binding the contract event 0xe81699b85113eea1c73e10588b2b035e55893369632173afd43feb192fac64e3.
//
// Solidity: event log_bytes32(bytes32 arg0)
func (_BitcoinState *BitcoinStateFilterer) WatchLogBytes32(opts *bind.WatchOpts, sink chan<- *BitcoinStateLogBytes32) (event.Subscription, error) {

	logs, sub, err := _BitcoinState.contract.WatchLogs(opts, "log_bytes32")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BitcoinStateLogBytes32)
				if err := _BitcoinState.contract.UnpackLog(event, "log_bytes32", log); err != nil {
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

// ParseLogBytes32 is a log parse operation binding the contract event 0xe81699b85113eea1c73e10588b2b035e55893369632173afd43feb192fac64e3.
//
// Solidity: event log_bytes32(bytes32 arg0)
func (_BitcoinState *BitcoinStateFilterer) ParseLogBytes32(log types.Log) (*BitcoinStateLogBytes32, error) {
	event := new(BitcoinStateLogBytes32)
	if err := _BitcoinState.contract.UnpackLog(event, "log_bytes32", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BitcoinStateLogIntIterator is returned from FilterLogInt and is used to iterate over the raw logs and unpacked data for LogInt events raised by the BitcoinState contract.
type BitcoinStateLogIntIterator struct {
	Event *BitcoinStateLogInt // Event containing the contract specifics and raw log

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
func (it *BitcoinStateLogIntIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BitcoinStateLogInt)
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
		it.Event = new(BitcoinStateLogInt)
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
func (it *BitcoinStateLogIntIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BitcoinStateLogIntIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BitcoinStateLogInt represents a LogInt event raised by the BitcoinState contract.
type BitcoinStateLogInt struct {
	Arg0 *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogInt is a free log retrieval operation binding the contract event 0x0eb5d52624c8d28ada9fc55a8c502ed5aa3fbe2fb6e91b71b5f376882b1d2fb8.
//
// Solidity: event log_int(int256 arg0)
func (_BitcoinState *BitcoinStateFilterer) FilterLogInt(opts *bind.FilterOpts) (*BitcoinStateLogIntIterator, error) {

	logs, sub, err := _BitcoinState.contract.FilterLogs(opts, "log_int")
	if err != nil {
		return nil, err
	}
	return &BitcoinStateLogIntIterator{contract: _BitcoinState.contract, event: "log_int", logs: logs, sub: sub}, nil
}

// WatchLogInt is a free log subscription operation binding the contract event 0x0eb5d52624c8d28ada9fc55a8c502ed5aa3fbe2fb6e91b71b5f376882b1d2fb8.
//
// Solidity: event log_int(int256 arg0)
func (_BitcoinState *BitcoinStateFilterer) WatchLogInt(opts *bind.WatchOpts, sink chan<- *BitcoinStateLogInt) (event.Subscription, error) {

	logs, sub, err := _BitcoinState.contract.WatchLogs(opts, "log_int")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BitcoinStateLogInt)
				if err := _BitcoinState.contract.UnpackLog(event, "log_int", log); err != nil {
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

// ParseLogInt is a log parse operation binding the contract event 0x0eb5d52624c8d28ada9fc55a8c502ed5aa3fbe2fb6e91b71b5f376882b1d2fb8.
//
// Solidity: event log_int(int256 arg0)
func (_BitcoinState *BitcoinStateFilterer) ParseLogInt(log types.Log) (*BitcoinStateLogInt, error) {
	event := new(BitcoinStateLogInt)
	if err := _BitcoinState.contract.UnpackLog(event, "log_int", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BitcoinStateLogNamedAddressIterator is returned from FilterLogNamedAddress and is used to iterate over the raw logs and unpacked data for LogNamedAddress events raised by the BitcoinState contract.
type BitcoinStateLogNamedAddressIterator struct {
	Event *BitcoinStateLogNamedAddress // Event containing the contract specifics and raw log

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
func (it *BitcoinStateLogNamedAddressIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BitcoinStateLogNamedAddress)
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
		it.Event = new(BitcoinStateLogNamedAddress)
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
func (it *BitcoinStateLogNamedAddressIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BitcoinStateLogNamedAddressIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BitcoinStateLogNamedAddress represents a LogNamedAddress event raised by the BitcoinState contract.
type BitcoinStateLogNamedAddress struct {
	Key string
	Val common.Address
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedAddress is a free log retrieval operation binding the contract event 0x9c4e8541ca8f0dc1c413f9108f66d82d3cecb1bddbce437a61caa3175c4cc96f.
//
// Solidity: event log_named_address(string key, address val)
func (_BitcoinState *BitcoinStateFilterer) FilterLogNamedAddress(opts *bind.FilterOpts) (*BitcoinStateLogNamedAddressIterator, error) {

	logs, sub, err := _BitcoinState.contract.FilterLogs(opts, "log_named_address")
	if err != nil {
		return nil, err
	}
	return &BitcoinStateLogNamedAddressIterator{contract: _BitcoinState.contract, event: "log_named_address", logs: logs, sub: sub}, nil
}

// WatchLogNamedAddress is a free log subscription operation binding the contract event 0x9c4e8541ca8f0dc1c413f9108f66d82d3cecb1bddbce437a61caa3175c4cc96f.
//
// Solidity: event log_named_address(string key, address val)
func (_BitcoinState *BitcoinStateFilterer) WatchLogNamedAddress(opts *bind.WatchOpts, sink chan<- *BitcoinStateLogNamedAddress) (event.Subscription, error) {

	logs, sub, err := _BitcoinState.contract.WatchLogs(opts, "log_named_address")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BitcoinStateLogNamedAddress)
				if err := _BitcoinState.contract.UnpackLog(event, "log_named_address", log); err != nil {
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

// ParseLogNamedAddress is a log parse operation binding the contract event 0x9c4e8541ca8f0dc1c413f9108f66d82d3cecb1bddbce437a61caa3175c4cc96f.
//
// Solidity: event log_named_address(string key, address val)
func (_BitcoinState *BitcoinStateFilterer) ParseLogNamedAddress(log types.Log) (*BitcoinStateLogNamedAddress, error) {
	event := new(BitcoinStateLogNamedAddress)
	if err := _BitcoinState.contract.UnpackLog(event, "log_named_address", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BitcoinStateLogNamedArrayIterator is returned from FilterLogNamedArray and is used to iterate over the raw logs and unpacked data for LogNamedArray events raised by the BitcoinState contract.
type BitcoinStateLogNamedArrayIterator struct {
	Event *BitcoinStateLogNamedArray // Event containing the contract specifics and raw log

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
func (it *BitcoinStateLogNamedArrayIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BitcoinStateLogNamedArray)
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
		it.Event = new(BitcoinStateLogNamedArray)
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
func (it *BitcoinStateLogNamedArrayIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BitcoinStateLogNamedArrayIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BitcoinStateLogNamedArray represents a LogNamedArray event raised by the BitcoinState contract.
type BitcoinStateLogNamedArray struct {
	Key string
	Val []*big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedArray is a free log retrieval operation binding the contract event 0x00aaa39c9ffb5f567a4534380c737075702e1f7f14107fc95328e3b56c0325fb.
//
// Solidity: event log_named_array(string key, uint256[] val)
func (_BitcoinState *BitcoinStateFilterer) FilterLogNamedArray(opts *bind.FilterOpts) (*BitcoinStateLogNamedArrayIterator, error) {

	logs, sub, err := _BitcoinState.contract.FilterLogs(opts, "log_named_array")
	if err != nil {
		return nil, err
	}
	return &BitcoinStateLogNamedArrayIterator{contract: _BitcoinState.contract, event: "log_named_array", logs: logs, sub: sub}, nil
}

// WatchLogNamedArray is a free log subscription operation binding the contract event 0x00aaa39c9ffb5f567a4534380c737075702e1f7f14107fc95328e3b56c0325fb.
//
// Solidity: event log_named_array(string key, uint256[] val)
func (_BitcoinState *BitcoinStateFilterer) WatchLogNamedArray(opts *bind.WatchOpts, sink chan<- *BitcoinStateLogNamedArray) (event.Subscription, error) {

	logs, sub, err := _BitcoinState.contract.WatchLogs(opts, "log_named_array")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BitcoinStateLogNamedArray)
				if err := _BitcoinState.contract.UnpackLog(event, "log_named_array", log); err != nil {
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

// ParseLogNamedArray is a log parse operation binding the contract event 0x00aaa39c9ffb5f567a4534380c737075702e1f7f14107fc95328e3b56c0325fb.
//
// Solidity: event log_named_array(string key, uint256[] val)
func (_BitcoinState *BitcoinStateFilterer) ParseLogNamedArray(log types.Log) (*BitcoinStateLogNamedArray, error) {
	event := new(BitcoinStateLogNamedArray)
	if err := _BitcoinState.contract.UnpackLog(event, "log_named_array", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BitcoinStateLogNamedArray0Iterator is returned from FilterLogNamedArray0 and is used to iterate over the raw logs and unpacked data for LogNamedArray0 events raised by the BitcoinState contract.
type BitcoinStateLogNamedArray0Iterator struct {
	Event *BitcoinStateLogNamedArray0 // Event containing the contract specifics and raw log

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
func (it *BitcoinStateLogNamedArray0Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BitcoinStateLogNamedArray0)
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
		it.Event = new(BitcoinStateLogNamedArray0)
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
func (it *BitcoinStateLogNamedArray0Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BitcoinStateLogNamedArray0Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BitcoinStateLogNamedArray0 represents a LogNamedArray0 event raised by the BitcoinState contract.
type BitcoinStateLogNamedArray0 struct {
	Key string
	Val []*big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedArray0 is a free log retrieval operation binding the contract event 0xa73eda09662f46dde729be4611385ff34fe6c44fbbc6f7e17b042b59a3445b57.
//
// Solidity: event log_named_array(string key, int256[] val)
func (_BitcoinState *BitcoinStateFilterer) FilterLogNamedArray0(opts *bind.FilterOpts) (*BitcoinStateLogNamedArray0Iterator, error) {

	logs, sub, err := _BitcoinState.contract.FilterLogs(opts, "log_named_array0")
	if err != nil {
		return nil, err
	}
	return &BitcoinStateLogNamedArray0Iterator{contract: _BitcoinState.contract, event: "log_named_array0", logs: logs, sub: sub}, nil
}

// WatchLogNamedArray0 is a free log subscription operation binding the contract event 0xa73eda09662f46dde729be4611385ff34fe6c44fbbc6f7e17b042b59a3445b57.
//
// Solidity: event log_named_array(string key, int256[] val)
func (_BitcoinState *BitcoinStateFilterer) WatchLogNamedArray0(opts *bind.WatchOpts, sink chan<- *BitcoinStateLogNamedArray0) (event.Subscription, error) {

	logs, sub, err := _BitcoinState.contract.WatchLogs(opts, "log_named_array0")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BitcoinStateLogNamedArray0)
				if err := _BitcoinState.contract.UnpackLog(event, "log_named_array0", log); err != nil {
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

// ParseLogNamedArray0 is a log parse operation binding the contract event 0xa73eda09662f46dde729be4611385ff34fe6c44fbbc6f7e17b042b59a3445b57.
//
// Solidity: event log_named_array(string key, int256[] val)
func (_BitcoinState *BitcoinStateFilterer) ParseLogNamedArray0(log types.Log) (*BitcoinStateLogNamedArray0, error) {
	event := new(BitcoinStateLogNamedArray0)
	if err := _BitcoinState.contract.UnpackLog(event, "log_named_array0", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BitcoinStateLogNamedArray1Iterator is returned from FilterLogNamedArray1 and is used to iterate over the raw logs and unpacked data for LogNamedArray1 events raised by the BitcoinState contract.
type BitcoinStateLogNamedArray1Iterator struct {
	Event *BitcoinStateLogNamedArray1 // Event containing the contract specifics and raw log

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
func (it *BitcoinStateLogNamedArray1Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BitcoinStateLogNamedArray1)
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
		it.Event = new(BitcoinStateLogNamedArray1)
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
func (it *BitcoinStateLogNamedArray1Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BitcoinStateLogNamedArray1Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BitcoinStateLogNamedArray1 represents a LogNamedArray1 event raised by the BitcoinState contract.
type BitcoinStateLogNamedArray1 struct {
	Key string
	Val []common.Address
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedArray1 is a free log retrieval operation binding the contract event 0x3bcfb2ae2e8d132dd1fce7cf278a9a19756a9fceabe470df3bdabb4bc577d1bd.
//
// Solidity: event log_named_array(string key, address[] val)
func (_BitcoinState *BitcoinStateFilterer) FilterLogNamedArray1(opts *bind.FilterOpts) (*BitcoinStateLogNamedArray1Iterator, error) {

	logs, sub, err := _BitcoinState.contract.FilterLogs(opts, "log_named_array1")
	if err != nil {
		return nil, err
	}
	return &BitcoinStateLogNamedArray1Iterator{contract: _BitcoinState.contract, event: "log_named_array1", logs: logs, sub: sub}, nil
}

// WatchLogNamedArray1 is a free log subscription operation binding the contract event 0x3bcfb2ae2e8d132dd1fce7cf278a9a19756a9fceabe470df3bdabb4bc577d1bd.
//
// Solidity: event log_named_array(string key, address[] val)
func (_BitcoinState *BitcoinStateFilterer) WatchLogNamedArray1(opts *bind.WatchOpts, sink chan<- *BitcoinStateLogNamedArray1) (event.Subscription, error) {

	logs, sub, err := _BitcoinState.contract.WatchLogs(opts, "log_named_array1")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BitcoinStateLogNamedArray1)
				if err := _BitcoinState.contract.UnpackLog(event, "log_named_array1", log); err != nil {
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

// ParseLogNamedArray1 is a log parse operation binding the contract event 0x3bcfb2ae2e8d132dd1fce7cf278a9a19756a9fceabe470df3bdabb4bc577d1bd.
//
// Solidity: event log_named_array(string key, address[] val)
func (_BitcoinState *BitcoinStateFilterer) ParseLogNamedArray1(log types.Log) (*BitcoinStateLogNamedArray1, error) {
	event := new(BitcoinStateLogNamedArray1)
	if err := _BitcoinState.contract.UnpackLog(event, "log_named_array1", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BitcoinStateLogNamedBytesIterator is returned from FilterLogNamedBytes and is used to iterate over the raw logs and unpacked data for LogNamedBytes events raised by the BitcoinState contract.
type BitcoinStateLogNamedBytesIterator struct {
	Event *BitcoinStateLogNamedBytes // Event containing the contract specifics and raw log

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
func (it *BitcoinStateLogNamedBytesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BitcoinStateLogNamedBytes)
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
		it.Event = new(BitcoinStateLogNamedBytes)
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
func (it *BitcoinStateLogNamedBytesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BitcoinStateLogNamedBytesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BitcoinStateLogNamedBytes represents a LogNamedBytes event raised by the BitcoinState contract.
type BitcoinStateLogNamedBytes struct {
	Key string
	Val []byte
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedBytes is a free log retrieval operation binding the contract event 0xd26e16cad4548705e4c9e2d94f98ee91c289085ee425594fd5635fa2964ccf18.
//
// Solidity: event log_named_bytes(string key, bytes val)
func (_BitcoinState *BitcoinStateFilterer) FilterLogNamedBytes(opts *bind.FilterOpts) (*BitcoinStateLogNamedBytesIterator, error) {

	logs, sub, err := _BitcoinState.contract.FilterLogs(opts, "log_named_bytes")
	if err != nil {
		return nil, err
	}
	return &BitcoinStateLogNamedBytesIterator{contract: _BitcoinState.contract, event: "log_named_bytes", logs: logs, sub: sub}, nil
}

// WatchLogNamedBytes is a free log subscription operation binding the contract event 0xd26e16cad4548705e4c9e2d94f98ee91c289085ee425594fd5635fa2964ccf18.
//
// Solidity: event log_named_bytes(string key, bytes val)
func (_BitcoinState *BitcoinStateFilterer) WatchLogNamedBytes(opts *bind.WatchOpts, sink chan<- *BitcoinStateLogNamedBytes) (event.Subscription, error) {

	logs, sub, err := _BitcoinState.contract.WatchLogs(opts, "log_named_bytes")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BitcoinStateLogNamedBytes)
				if err := _BitcoinState.contract.UnpackLog(event, "log_named_bytes", log); err != nil {
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

// ParseLogNamedBytes is a log parse operation binding the contract event 0xd26e16cad4548705e4c9e2d94f98ee91c289085ee425594fd5635fa2964ccf18.
//
// Solidity: event log_named_bytes(string key, bytes val)
func (_BitcoinState *BitcoinStateFilterer) ParseLogNamedBytes(log types.Log) (*BitcoinStateLogNamedBytes, error) {
	event := new(BitcoinStateLogNamedBytes)
	if err := _BitcoinState.contract.UnpackLog(event, "log_named_bytes", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BitcoinStateLogNamedBytes32Iterator is returned from FilterLogNamedBytes32 and is used to iterate over the raw logs and unpacked data for LogNamedBytes32 events raised by the BitcoinState contract.
type BitcoinStateLogNamedBytes32Iterator struct {
	Event *BitcoinStateLogNamedBytes32 // Event containing the contract specifics and raw log

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
func (it *BitcoinStateLogNamedBytes32Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BitcoinStateLogNamedBytes32)
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
		it.Event = new(BitcoinStateLogNamedBytes32)
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
func (it *BitcoinStateLogNamedBytes32Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BitcoinStateLogNamedBytes32Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BitcoinStateLogNamedBytes32 represents a LogNamedBytes32 event raised by the BitcoinState contract.
type BitcoinStateLogNamedBytes32 struct {
	Key string
	Val [32]byte
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedBytes32 is a free log retrieval operation binding the contract event 0xafb795c9c61e4fe7468c386f925d7a5429ecad9c0495ddb8d38d690614d32f99.
//
// Solidity: event log_named_bytes32(string key, bytes32 val)
func (_BitcoinState *BitcoinStateFilterer) FilterLogNamedBytes32(opts *bind.FilterOpts) (*BitcoinStateLogNamedBytes32Iterator, error) {

	logs, sub, err := _BitcoinState.contract.FilterLogs(opts, "log_named_bytes32")
	if err != nil {
		return nil, err
	}
	return &BitcoinStateLogNamedBytes32Iterator{contract: _BitcoinState.contract, event: "log_named_bytes32", logs: logs, sub: sub}, nil
}

// WatchLogNamedBytes32 is a free log subscription operation binding the contract event 0xafb795c9c61e4fe7468c386f925d7a5429ecad9c0495ddb8d38d690614d32f99.
//
// Solidity: event log_named_bytes32(string key, bytes32 val)
func (_BitcoinState *BitcoinStateFilterer) WatchLogNamedBytes32(opts *bind.WatchOpts, sink chan<- *BitcoinStateLogNamedBytes32) (event.Subscription, error) {

	logs, sub, err := _BitcoinState.contract.WatchLogs(opts, "log_named_bytes32")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BitcoinStateLogNamedBytes32)
				if err := _BitcoinState.contract.UnpackLog(event, "log_named_bytes32", log); err != nil {
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

// ParseLogNamedBytes32 is a log parse operation binding the contract event 0xafb795c9c61e4fe7468c386f925d7a5429ecad9c0495ddb8d38d690614d32f99.
//
// Solidity: event log_named_bytes32(string key, bytes32 val)
func (_BitcoinState *BitcoinStateFilterer) ParseLogNamedBytes32(log types.Log) (*BitcoinStateLogNamedBytes32, error) {
	event := new(BitcoinStateLogNamedBytes32)
	if err := _BitcoinState.contract.UnpackLog(event, "log_named_bytes32", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BitcoinStateLogNamedDecimalIntIterator is returned from FilterLogNamedDecimalInt and is used to iterate over the raw logs and unpacked data for LogNamedDecimalInt events raised by the BitcoinState contract.
type BitcoinStateLogNamedDecimalIntIterator struct {
	Event *BitcoinStateLogNamedDecimalInt // Event containing the contract specifics and raw log

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
func (it *BitcoinStateLogNamedDecimalIntIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BitcoinStateLogNamedDecimalInt)
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
		it.Event = new(BitcoinStateLogNamedDecimalInt)
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
func (it *BitcoinStateLogNamedDecimalIntIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BitcoinStateLogNamedDecimalIntIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BitcoinStateLogNamedDecimalInt represents a LogNamedDecimalInt event raised by the BitcoinState contract.
type BitcoinStateLogNamedDecimalInt struct {
	Key      string
	Val      *big.Int
	Decimals *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterLogNamedDecimalInt is a free log retrieval operation binding the contract event 0x5da6ce9d51151ba10c09a559ef24d520b9dac5c5b8810ae8434e4d0d86411a95.
//
// Solidity: event log_named_decimal_int(string key, int256 val, uint256 decimals)
func (_BitcoinState *BitcoinStateFilterer) FilterLogNamedDecimalInt(opts *bind.FilterOpts) (*BitcoinStateLogNamedDecimalIntIterator, error) {

	logs, sub, err := _BitcoinState.contract.FilterLogs(opts, "log_named_decimal_int")
	if err != nil {
		return nil, err
	}
	return &BitcoinStateLogNamedDecimalIntIterator{contract: _BitcoinState.contract, event: "log_named_decimal_int", logs: logs, sub: sub}, nil
}

// WatchLogNamedDecimalInt is a free log subscription operation binding the contract event 0x5da6ce9d51151ba10c09a559ef24d520b9dac5c5b8810ae8434e4d0d86411a95.
//
// Solidity: event log_named_decimal_int(string key, int256 val, uint256 decimals)
func (_BitcoinState *BitcoinStateFilterer) WatchLogNamedDecimalInt(opts *bind.WatchOpts, sink chan<- *BitcoinStateLogNamedDecimalInt) (event.Subscription, error) {

	logs, sub, err := _BitcoinState.contract.WatchLogs(opts, "log_named_decimal_int")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BitcoinStateLogNamedDecimalInt)
				if err := _BitcoinState.contract.UnpackLog(event, "log_named_decimal_int", log); err != nil {
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

// ParseLogNamedDecimalInt is a log parse operation binding the contract event 0x5da6ce9d51151ba10c09a559ef24d520b9dac5c5b8810ae8434e4d0d86411a95.
//
// Solidity: event log_named_decimal_int(string key, int256 val, uint256 decimals)
func (_BitcoinState *BitcoinStateFilterer) ParseLogNamedDecimalInt(log types.Log) (*BitcoinStateLogNamedDecimalInt, error) {
	event := new(BitcoinStateLogNamedDecimalInt)
	if err := _BitcoinState.contract.UnpackLog(event, "log_named_decimal_int", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BitcoinStateLogNamedDecimalUintIterator is returned from FilterLogNamedDecimalUint and is used to iterate over the raw logs and unpacked data for LogNamedDecimalUint events raised by the BitcoinState contract.
type BitcoinStateLogNamedDecimalUintIterator struct {
	Event *BitcoinStateLogNamedDecimalUint // Event containing the contract specifics and raw log

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
func (it *BitcoinStateLogNamedDecimalUintIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BitcoinStateLogNamedDecimalUint)
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
		it.Event = new(BitcoinStateLogNamedDecimalUint)
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
func (it *BitcoinStateLogNamedDecimalUintIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BitcoinStateLogNamedDecimalUintIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BitcoinStateLogNamedDecimalUint represents a LogNamedDecimalUint event raised by the BitcoinState contract.
type BitcoinStateLogNamedDecimalUint struct {
	Key      string
	Val      *big.Int
	Decimals *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterLogNamedDecimalUint is a free log retrieval operation binding the contract event 0xeb8ba43ced7537421946bd43e828b8b2b8428927aa8f801c13d934bf11aca57b.
//
// Solidity: event log_named_decimal_uint(string key, uint256 val, uint256 decimals)
func (_BitcoinState *BitcoinStateFilterer) FilterLogNamedDecimalUint(opts *bind.FilterOpts) (*BitcoinStateLogNamedDecimalUintIterator, error) {

	logs, sub, err := _BitcoinState.contract.FilterLogs(opts, "log_named_decimal_uint")
	if err != nil {
		return nil, err
	}
	return &BitcoinStateLogNamedDecimalUintIterator{contract: _BitcoinState.contract, event: "log_named_decimal_uint", logs: logs, sub: sub}, nil
}

// WatchLogNamedDecimalUint is a free log subscription operation binding the contract event 0xeb8ba43ced7537421946bd43e828b8b2b8428927aa8f801c13d934bf11aca57b.
//
// Solidity: event log_named_decimal_uint(string key, uint256 val, uint256 decimals)
func (_BitcoinState *BitcoinStateFilterer) WatchLogNamedDecimalUint(opts *bind.WatchOpts, sink chan<- *BitcoinStateLogNamedDecimalUint) (event.Subscription, error) {

	logs, sub, err := _BitcoinState.contract.WatchLogs(opts, "log_named_decimal_uint")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BitcoinStateLogNamedDecimalUint)
				if err := _BitcoinState.contract.UnpackLog(event, "log_named_decimal_uint", log); err != nil {
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

// ParseLogNamedDecimalUint is a log parse operation binding the contract event 0xeb8ba43ced7537421946bd43e828b8b2b8428927aa8f801c13d934bf11aca57b.
//
// Solidity: event log_named_decimal_uint(string key, uint256 val, uint256 decimals)
func (_BitcoinState *BitcoinStateFilterer) ParseLogNamedDecimalUint(log types.Log) (*BitcoinStateLogNamedDecimalUint, error) {
	event := new(BitcoinStateLogNamedDecimalUint)
	if err := _BitcoinState.contract.UnpackLog(event, "log_named_decimal_uint", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BitcoinStateLogNamedIntIterator is returned from FilterLogNamedInt and is used to iterate over the raw logs and unpacked data for LogNamedInt events raised by the BitcoinState contract.
type BitcoinStateLogNamedIntIterator struct {
	Event *BitcoinStateLogNamedInt // Event containing the contract specifics and raw log

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
func (it *BitcoinStateLogNamedIntIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BitcoinStateLogNamedInt)
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
		it.Event = new(BitcoinStateLogNamedInt)
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
func (it *BitcoinStateLogNamedIntIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BitcoinStateLogNamedIntIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BitcoinStateLogNamedInt represents a LogNamedInt event raised by the BitcoinState contract.
type BitcoinStateLogNamedInt struct {
	Key string
	Val *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedInt is a free log retrieval operation binding the contract event 0x2fe632779174374378442a8e978bccfbdcc1d6b2b0d81f7e8eb776ab2286f168.
//
// Solidity: event log_named_int(string key, int256 val)
func (_BitcoinState *BitcoinStateFilterer) FilterLogNamedInt(opts *bind.FilterOpts) (*BitcoinStateLogNamedIntIterator, error) {

	logs, sub, err := _BitcoinState.contract.FilterLogs(opts, "log_named_int")
	if err != nil {
		return nil, err
	}
	return &BitcoinStateLogNamedIntIterator{contract: _BitcoinState.contract, event: "log_named_int", logs: logs, sub: sub}, nil
}

// WatchLogNamedInt is a free log subscription operation binding the contract event 0x2fe632779174374378442a8e978bccfbdcc1d6b2b0d81f7e8eb776ab2286f168.
//
// Solidity: event log_named_int(string key, int256 val)
func (_BitcoinState *BitcoinStateFilterer) WatchLogNamedInt(opts *bind.WatchOpts, sink chan<- *BitcoinStateLogNamedInt) (event.Subscription, error) {

	logs, sub, err := _BitcoinState.contract.WatchLogs(opts, "log_named_int")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BitcoinStateLogNamedInt)
				if err := _BitcoinState.contract.UnpackLog(event, "log_named_int", log); err != nil {
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

// ParseLogNamedInt is a log parse operation binding the contract event 0x2fe632779174374378442a8e978bccfbdcc1d6b2b0d81f7e8eb776ab2286f168.
//
// Solidity: event log_named_int(string key, int256 val)
func (_BitcoinState *BitcoinStateFilterer) ParseLogNamedInt(log types.Log) (*BitcoinStateLogNamedInt, error) {
	event := new(BitcoinStateLogNamedInt)
	if err := _BitcoinState.contract.UnpackLog(event, "log_named_int", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BitcoinStateLogNamedStringIterator is returned from FilterLogNamedString and is used to iterate over the raw logs and unpacked data for LogNamedString events raised by the BitcoinState contract.
type BitcoinStateLogNamedStringIterator struct {
	Event *BitcoinStateLogNamedString // Event containing the contract specifics and raw log

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
func (it *BitcoinStateLogNamedStringIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BitcoinStateLogNamedString)
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
		it.Event = new(BitcoinStateLogNamedString)
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
func (it *BitcoinStateLogNamedStringIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BitcoinStateLogNamedStringIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BitcoinStateLogNamedString represents a LogNamedString event raised by the BitcoinState contract.
type BitcoinStateLogNamedString struct {
	Key string
	Val string
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedString is a free log retrieval operation binding the contract event 0x280f4446b28a1372417dda658d30b95b2992b12ac9c7f378535f29a97acf3583.
//
// Solidity: event log_named_string(string key, string val)
func (_BitcoinState *BitcoinStateFilterer) FilterLogNamedString(opts *bind.FilterOpts) (*BitcoinStateLogNamedStringIterator, error) {

	logs, sub, err := _BitcoinState.contract.FilterLogs(opts, "log_named_string")
	if err != nil {
		return nil, err
	}
	return &BitcoinStateLogNamedStringIterator{contract: _BitcoinState.contract, event: "log_named_string", logs: logs, sub: sub}, nil
}

// WatchLogNamedString is a free log subscription operation binding the contract event 0x280f4446b28a1372417dda658d30b95b2992b12ac9c7f378535f29a97acf3583.
//
// Solidity: event log_named_string(string key, string val)
func (_BitcoinState *BitcoinStateFilterer) WatchLogNamedString(opts *bind.WatchOpts, sink chan<- *BitcoinStateLogNamedString) (event.Subscription, error) {

	logs, sub, err := _BitcoinState.contract.WatchLogs(opts, "log_named_string")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BitcoinStateLogNamedString)
				if err := _BitcoinState.contract.UnpackLog(event, "log_named_string", log); err != nil {
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

// ParseLogNamedString is a log parse operation binding the contract event 0x280f4446b28a1372417dda658d30b95b2992b12ac9c7f378535f29a97acf3583.
//
// Solidity: event log_named_string(string key, string val)
func (_BitcoinState *BitcoinStateFilterer) ParseLogNamedString(log types.Log) (*BitcoinStateLogNamedString, error) {
	event := new(BitcoinStateLogNamedString)
	if err := _BitcoinState.contract.UnpackLog(event, "log_named_string", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BitcoinStateLogNamedUintIterator is returned from FilterLogNamedUint and is used to iterate over the raw logs and unpacked data for LogNamedUint events raised by the BitcoinState contract.
type BitcoinStateLogNamedUintIterator struct {
	Event *BitcoinStateLogNamedUint // Event containing the contract specifics and raw log

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
func (it *BitcoinStateLogNamedUintIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BitcoinStateLogNamedUint)
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
		it.Event = new(BitcoinStateLogNamedUint)
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
func (it *BitcoinStateLogNamedUintIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BitcoinStateLogNamedUintIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BitcoinStateLogNamedUint represents a LogNamedUint event raised by the BitcoinState contract.
type BitcoinStateLogNamedUint struct {
	Key string
	Val *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedUint is a free log retrieval operation binding the contract event 0xb2de2fbe801a0df6c0cbddfd448ba3c41d48a040ca35c56c8196ef0fcae721a8.
//
// Solidity: event log_named_uint(string key, uint256 val)
func (_BitcoinState *BitcoinStateFilterer) FilterLogNamedUint(opts *bind.FilterOpts) (*BitcoinStateLogNamedUintIterator, error) {

	logs, sub, err := _BitcoinState.contract.FilterLogs(opts, "log_named_uint")
	if err != nil {
		return nil, err
	}
	return &BitcoinStateLogNamedUintIterator{contract: _BitcoinState.contract, event: "log_named_uint", logs: logs, sub: sub}, nil
}

// WatchLogNamedUint is a free log subscription operation binding the contract event 0xb2de2fbe801a0df6c0cbddfd448ba3c41d48a040ca35c56c8196ef0fcae721a8.
//
// Solidity: event log_named_uint(string key, uint256 val)
func (_BitcoinState *BitcoinStateFilterer) WatchLogNamedUint(opts *bind.WatchOpts, sink chan<- *BitcoinStateLogNamedUint) (event.Subscription, error) {

	logs, sub, err := _BitcoinState.contract.WatchLogs(opts, "log_named_uint")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BitcoinStateLogNamedUint)
				if err := _BitcoinState.contract.UnpackLog(event, "log_named_uint", log); err != nil {
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

// ParseLogNamedUint is a log parse operation binding the contract event 0xb2de2fbe801a0df6c0cbddfd448ba3c41d48a040ca35c56c8196ef0fcae721a8.
//
// Solidity: event log_named_uint(string key, uint256 val)
func (_BitcoinState *BitcoinStateFilterer) ParseLogNamedUint(log types.Log) (*BitcoinStateLogNamedUint, error) {
	event := new(BitcoinStateLogNamedUint)
	if err := _BitcoinState.contract.UnpackLog(event, "log_named_uint", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BitcoinStateLogStringIterator is returned from FilterLogString and is used to iterate over the raw logs and unpacked data for LogString events raised by the BitcoinState contract.
type BitcoinStateLogStringIterator struct {
	Event *BitcoinStateLogString // Event containing the contract specifics and raw log

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
func (it *BitcoinStateLogStringIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BitcoinStateLogString)
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
		it.Event = new(BitcoinStateLogString)
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
func (it *BitcoinStateLogStringIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BitcoinStateLogStringIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BitcoinStateLogString represents a LogString event raised by the BitcoinState contract.
type BitcoinStateLogString struct {
	Arg0 string
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogString is a free log retrieval operation binding the contract event 0x0b2e13ff20ac7b474198655583edf70dedd2c1dc980e329c4fbb2fc0748b796b.
//
// Solidity: event log_string(string arg0)
func (_BitcoinState *BitcoinStateFilterer) FilterLogString(opts *bind.FilterOpts) (*BitcoinStateLogStringIterator, error) {

	logs, sub, err := _BitcoinState.contract.FilterLogs(opts, "log_string")
	if err != nil {
		return nil, err
	}
	return &BitcoinStateLogStringIterator{contract: _BitcoinState.contract, event: "log_string", logs: logs, sub: sub}, nil
}

// WatchLogString is a free log subscription operation binding the contract event 0x0b2e13ff20ac7b474198655583edf70dedd2c1dc980e329c4fbb2fc0748b796b.
//
// Solidity: event log_string(string arg0)
func (_BitcoinState *BitcoinStateFilterer) WatchLogString(opts *bind.WatchOpts, sink chan<- *BitcoinStateLogString) (event.Subscription, error) {

	logs, sub, err := _BitcoinState.contract.WatchLogs(opts, "log_string")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BitcoinStateLogString)
				if err := _BitcoinState.contract.UnpackLog(event, "log_string", log); err != nil {
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

// ParseLogString is a log parse operation binding the contract event 0x0b2e13ff20ac7b474198655583edf70dedd2c1dc980e329c4fbb2fc0748b796b.
//
// Solidity: event log_string(string arg0)
func (_BitcoinState *BitcoinStateFilterer) ParseLogString(log types.Log) (*BitcoinStateLogString, error) {
	event := new(BitcoinStateLogString)
	if err := _BitcoinState.contract.UnpackLog(event, "log_string", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BitcoinStateLogUintIterator is returned from FilterLogUint and is used to iterate over the raw logs and unpacked data for LogUint events raised by the BitcoinState contract.
type BitcoinStateLogUintIterator struct {
	Event *BitcoinStateLogUint // Event containing the contract specifics and raw log

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
func (it *BitcoinStateLogUintIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BitcoinStateLogUint)
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
		it.Event = new(BitcoinStateLogUint)
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
func (it *BitcoinStateLogUintIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BitcoinStateLogUintIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BitcoinStateLogUint represents a LogUint event raised by the BitcoinState contract.
type BitcoinStateLogUint struct {
	Arg0 *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogUint is a free log retrieval operation binding the contract event 0x2cab9790510fd8bdfbd2115288db33fec66691d476efc5427cfd4c0969301755.
//
// Solidity: event log_uint(uint256 arg0)
func (_BitcoinState *BitcoinStateFilterer) FilterLogUint(opts *bind.FilterOpts) (*BitcoinStateLogUintIterator, error) {

	logs, sub, err := _BitcoinState.contract.FilterLogs(opts, "log_uint")
	if err != nil {
		return nil, err
	}
	return &BitcoinStateLogUintIterator{contract: _BitcoinState.contract, event: "log_uint", logs: logs, sub: sub}, nil
}

// WatchLogUint is a free log subscription operation binding the contract event 0x2cab9790510fd8bdfbd2115288db33fec66691d476efc5427cfd4c0969301755.
//
// Solidity: event log_uint(uint256 arg0)
func (_BitcoinState *BitcoinStateFilterer) WatchLogUint(opts *bind.WatchOpts, sink chan<- *BitcoinStateLogUint) (event.Subscription, error) {

	logs, sub, err := _BitcoinState.contract.WatchLogs(opts, "log_uint")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BitcoinStateLogUint)
				if err := _BitcoinState.contract.UnpackLog(event, "log_uint", log); err != nil {
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

// ParseLogUint is a log parse operation binding the contract event 0x2cab9790510fd8bdfbd2115288db33fec66691d476efc5427cfd4c0969301755.
//
// Solidity: event log_uint(uint256 arg0)
func (_BitcoinState *BitcoinStateFilterer) ParseLogUint(log types.Log) (*BitcoinStateLogUint, error) {
	event := new(BitcoinStateLogUint)
	if err := _BitcoinState.contract.UnpackLog(event, "log_uint", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BitcoinStateLogsIterator is returned from FilterLogs and is used to iterate over the raw logs and unpacked data for Logs events raised by the BitcoinState contract.
type BitcoinStateLogsIterator struct {
	Event *BitcoinStateLogs // Event containing the contract specifics and raw log

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
func (it *BitcoinStateLogsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BitcoinStateLogs)
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
		it.Event = new(BitcoinStateLogs)
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
func (it *BitcoinStateLogsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BitcoinStateLogsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BitcoinStateLogs represents a Logs event raised by the BitcoinState contract.
type BitcoinStateLogs struct {
	Arg0 []byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogs is a free log retrieval operation binding the contract event 0xe7950ede0394b9f2ce4a5a1bf5a7e1852411f7e6661b4308c913c4bfd11027e4.
//
// Solidity: event logs(bytes arg0)
func (_BitcoinState *BitcoinStateFilterer) FilterLogs(opts *bind.FilterOpts) (*BitcoinStateLogsIterator, error) {

	logs, sub, err := _BitcoinState.contract.FilterLogs(opts, "logs")
	if err != nil {
		return nil, err
	}
	return &BitcoinStateLogsIterator{contract: _BitcoinState.contract, event: "logs", logs: logs, sub: sub}, nil
}

// WatchLogs is a free log subscription operation binding the contract event 0xe7950ede0394b9f2ce4a5a1bf5a7e1852411f7e6661b4308c913c4bfd11027e4.
//
// Solidity: event logs(bytes arg0)
func (_BitcoinState *BitcoinStateFilterer) WatchLogs(opts *bind.WatchOpts, sink chan<- *BitcoinStateLogs) (event.Subscription, error) {

	logs, sub, err := _BitcoinState.contract.WatchLogs(opts, "logs")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BitcoinStateLogs)
				if err := _BitcoinState.contract.UnpackLog(event, "logs", log); err != nil {
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

// ParseLogs is a log parse operation binding the contract event 0xe7950ede0394b9f2ce4a5a1bf5a7e1852411f7e6661b4308c913c4bfd11027e4.
//
// Solidity: event logs(bytes arg0)
func (_BitcoinState *BitcoinStateFilterer) ParseLogs(log types.Log) (*BitcoinStateLogs, error) {
	event := new(BitcoinStateLogs)
	if err := _BitcoinState.contract.UnpackLog(event, "logs", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
