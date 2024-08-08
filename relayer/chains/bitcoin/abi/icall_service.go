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

// CallServiceMetaData contains all meta data concerning the CallService contract.
var CallServiceMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"admin\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"executeCall\",\"inputs\":[{\"name\":\"_reqId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"executeMessage\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"from\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"protocols\",\"type\":\"string[]\",\"internalType\":\"string[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"executeRollback\",\"inputs\":[{\"name\":\"_sn\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getDefaultConnection\",\"inputs\":[{\"name\":\"_nid\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getFee\",\"inputs\":[{\"name\":\"_net\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_rollback\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"_sources\",\"type\":\"string[]\",\"internalType\":\"string[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getFee\",\"inputs\":[{\"name\":\"_net\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_rollback\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNetworkAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNetworkId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getProtocolFee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getProtocolFeeHandler\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"handleBTPError\",\"inputs\":[{\"name\":\"_src\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_svc\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_sn\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_code\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_msg\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"handleBTPMessage\",\"inputs\":[{\"name\":\"_from\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_svc\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_sn\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_msg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"handleError\",\"inputs\":[{\"name\":\"_sn\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"handleMessage\",\"inputs\":[{\"name\":\"_from\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_msg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_nid\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"sendCall\",\"inputs\":[{\"name\":\"_to\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"sendCallMessage\",\"inputs\":[{\"name\":\"_to\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_rollback\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"sendCallMessage\",\"inputs\":[{\"name\":\"_to\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_rollback\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"sources\",\"type\":\"string[]\",\"internalType\":\"string[]\"},{\"name\":\"destinations\",\"type\":\"string[]\",\"internalType\":\"string[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"setAdmin\",\"inputs\":[{\"name\":\"_address\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setDefaultConnection\",\"inputs\":[{\"name\":\"_nid\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"connection\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setProtocolFee\",\"inputs\":[{\"name\":\"_value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setProtocolFeeHandler\",\"inputs\":[{\"name\":\"_addr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"verifySuccess\",\"inputs\":[{\"name\":\"_sn\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"CallExecuted\",\"inputs\":[{\"name\":\"_reqId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"_code\",\"type\":\"int256\",\"indexed\":false,\"internalType\":\"int256\"},{\"name\":\"_msg\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CallMessage\",\"inputs\":[{\"name\":\"_from\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"_to\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"_sn\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"_reqId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"_data\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CallMessageSent\",\"inputs\":[{\"name\":\"_from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"_to\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"_sn\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ResponseMessage\",\"inputs\":[{\"name\":\"_sn\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"_code\",\"type\":\"int256\",\"indexed\":false,\"internalType\":\"int256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RollbackExecuted\",\"inputs\":[{\"name\":\"_sn\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RollbackMessage\",\"inputs\":[{\"name\":\"_sn\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]}]",
}

// CallServiceABI is the input ABI used to generate the binding from.
// Deprecated: Use CallServiceMetaData.ABI instead.
var CallServiceABI = CallServiceMetaData.ABI

// CallService is an auto generated Go binding around an Ethereum contract.
type CallService struct {
	CallServiceCaller     // Read-only binding to the contract
	CallServiceTransactor // Write-only binding to the contract
	CallServiceFilterer   // Log filterer for contract events
}

// CallServiceCaller is an auto generated read-only Go binding around an Ethereum contract.
type CallServiceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CallServiceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CallServiceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CallServiceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CallServiceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CallServiceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CallServiceSession struct {
	Contract     *CallService      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CallServiceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CallServiceCallerSession struct {
	Contract *CallServiceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// CallServiceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CallServiceTransactorSession struct {
	Contract     *CallServiceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// CallServiceRaw is an auto generated low-level Go binding around an Ethereum contract.
type CallServiceRaw struct {
	Contract *CallService // Generic contract binding to access the raw methods on
}

// CallServiceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CallServiceCallerRaw struct {
	Contract *CallServiceCaller // Generic read-only contract binding to access the raw methods on
}

// CallServiceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CallServiceTransactorRaw struct {
	Contract *CallServiceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCallService creates a new instance of CallService, bound to a specific deployed contract.
func NewCallService(address common.Address, backend bind.ContractBackend) (*CallService, error) {
	contract, err := bindCallService(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CallService{CallServiceCaller: CallServiceCaller{contract: contract}, CallServiceTransactor: CallServiceTransactor{contract: contract}, CallServiceFilterer: CallServiceFilterer{contract: contract}}, nil
}

// NewCallServiceCaller creates a new read-only instance of CallService, bound to a specific deployed contract.
func NewCallServiceCaller(address common.Address, caller bind.ContractCaller) (*CallServiceCaller, error) {
	contract, err := bindCallService(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CallServiceCaller{contract: contract}, nil
}

// NewCallServiceTransactor creates a new write-only instance of CallService, bound to a specific deployed contract.
func NewCallServiceTransactor(address common.Address, transactor bind.ContractTransactor) (*CallServiceTransactor, error) {
	contract, err := bindCallService(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CallServiceTransactor{contract: contract}, nil
}

// NewCallServiceFilterer creates a new log filterer instance of CallService, bound to a specific deployed contract.
func NewCallServiceFilterer(address common.Address, filterer bind.ContractFilterer) (*CallServiceFilterer, error) {
	contract, err := bindCallService(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CallServiceFilterer{contract: contract}, nil
}

// bindCallService binds a generic wrapper to an already deployed contract.
func bindCallService(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CallServiceMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CallService *CallServiceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CallService.Contract.CallServiceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CallService *CallServiceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CallService.Contract.CallServiceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CallService *CallServiceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CallService.Contract.CallServiceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CallService *CallServiceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CallService.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CallService *CallServiceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CallService.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CallService *CallServiceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CallService.Contract.contract.Transact(opts, method, params...)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_CallService *CallServiceCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CallService.contract.Call(opts, &out, "admin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_CallService *CallServiceSession) Admin() (common.Address, error) {
	return _CallService.Contract.Admin(&_CallService.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_CallService *CallServiceCallerSession) Admin() (common.Address, error) {
	return _CallService.Contract.Admin(&_CallService.CallOpts)
}

// GetDefaultConnection is a free data retrieval call binding the contract method 0x9e553a4f.
//
// Solidity: function getDefaultConnection(string _nid) view returns(address)
func (_CallService *CallServiceCaller) GetDefaultConnection(opts *bind.CallOpts, _nid string) (common.Address, error) {
	var out []interface{}
	err := _CallService.contract.Call(opts, &out, "getDefaultConnection", _nid)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetDefaultConnection is a free data retrieval call binding the contract method 0x9e553a4f.
//
// Solidity: function getDefaultConnection(string _nid) view returns(address)
func (_CallService *CallServiceSession) GetDefaultConnection(_nid string) (common.Address, error) {
	return _CallService.Contract.GetDefaultConnection(&_CallService.CallOpts, _nid)
}

// GetDefaultConnection is a free data retrieval call binding the contract method 0x9e553a4f.
//
// Solidity: function getDefaultConnection(string _nid) view returns(address)
func (_CallService *CallServiceCallerSession) GetDefaultConnection(_nid string) (common.Address, error) {
	return _CallService.Contract.GetDefaultConnection(&_CallService.CallOpts, _nid)
}

// GetFee is a free data retrieval call binding the contract method 0x304a70b5.
//
// Solidity: function getFee(string _net, bool _rollback, string[] _sources) view returns(uint256)
func (_CallService *CallServiceCaller) GetFee(opts *bind.CallOpts, _net string, _rollback bool, _sources []string) (*big.Int, error) {
	var out []interface{}
	err := _CallService.contract.Call(opts, &out, "getFee", _net, _rollback, _sources)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetFee is a free data retrieval call binding the contract method 0x304a70b5.
//
// Solidity: function getFee(string _net, bool _rollback, string[] _sources) view returns(uint256)
func (_CallService *CallServiceSession) GetFee(_net string, _rollback bool, _sources []string) (*big.Int, error) {
	return _CallService.Contract.GetFee(&_CallService.CallOpts, _net, _rollback, _sources)
}

// GetFee is a free data retrieval call binding the contract method 0x304a70b5.
//
// Solidity: function getFee(string _net, bool _rollback, string[] _sources) view returns(uint256)
func (_CallService *CallServiceCallerSession) GetFee(_net string, _rollback bool, _sources []string) (*big.Int, error) {
	return _CallService.Contract.GetFee(&_CallService.CallOpts, _net, _rollback, _sources)
}

// GetFee0 is a free data retrieval call binding the contract method 0x7d4c4f4a.
//
// Solidity: function getFee(string _net, bool _rollback) view returns(uint256)
func (_CallService *CallServiceCaller) GetFee0(opts *bind.CallOpts, _net string, _rollback bool) (*big.Int, error) {
	var out []interface{}
	err := _CallService.contract.Call(opts, &out, "getFee0", _net, _rollback)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetFee0 is a free data retrieval call binding the contract method 0x7d4c4f4a.
//
// Solidity: function getFee(string _net, bool _rollback) view returns(uint256)
func (_CallService *CallServiceSession) GetFee0(_net string, _rollback bool) (*big.Int, error) {
	return _CallService.Contract.GetFee0(&_CallService.CallOpts, _net, _rollback)
}

// GetFee0 is a free data retrieval call binding the contract method 0x7d4c4f4a.
//
// Solidity: function getFee(string _net, bool _rollback) view returns(uint256)
func (_CallService *CallServiceCallerSession) GetFee0(_net string, _rollback bool) (*big.Int, error) {
	return _CallService.Contract.GetFee0(&_CallService.CallOpts, _net, _rollback)
}

// GetNetworkAddress is a free data retrieval call binding the contract method 0x6bf459cb.
//
// Solidity: function getNetworkAddress() view returns(string)
func (_CallService *CallServiceCaller) GetNetworkAddress(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _CallService.contract.Call(opts, &out, "getNetworkAddress")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetNetworkAddress is a free data retrieval call binding the contract method 0x6bf459cb.
//
// Solidity: function getNetworkAddress() view returns(string)
func (_CallService *CallServiceSession) GetNetworkAddress() (string, error) {
	return _CallService.Contract.GetNetworkAddress(&_CallService.CallOpts)
}

// GetNetworkAddress is a free data retrieval call binding the contract method 0x6bf459cb.
//
// Solidity: function getNetworkAddress() view returns(string)
func (_CallService *CallServiceCallerSession) GetNetworkAddress() (string, error) {
	return _CallService.Contract.GetNetworkAddress(&_CallService.CallOpts)
}

// GetNetworkId is a free data retrieval call binding the contract method 0x39c5f3fc.
//
// Solidity: function getNetworkId() view returns(string)
func (_CallService *CallServiceCaller) GetNetworkId(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _CallService.contract.Call(opts, &out, "getNetworkId")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetNetworkId is a free data retrieval call binding the contract method 0x39c5f3fc.
//
// Solidity: function getNetworkId() view returns(string)
func (_CallService *CallServiceSession) GetNetworkId() (string, error) {
	return _CallService.Contract.GetNetworkId(&_CallService.CallOpts)
}

// GetNetworkId is a free data retrieval call binding the contract method 0x39c5f3fc.
//
// Solidity: function getNetworkId() view returns(string)
func (_CallService *CallServiceCallerSession) GetNetworkId() (string, error) {
	return _CallService.Contract.GetNetworkId(&_CallService.CallOpts)
}

// GetProtocolFee is a free data retrieval call binding the contract method 0xa5a41031.
//
// Solidity: function getProtocolFee() view returns(uint256)
func (_CallService *CallServiceCaller) GetProtocolFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CallService.contract.Call(opts, &out, "getProtocolFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetProtocolFee is a free data retrieval call binding the contract method 0xa5a41031.
//
// Solidity: function getProtocolFee() view returns(uint256)
func (_CallService *CallServiceSession) GetProtocolFee() (*big.Int, error) {
	return _CallService.Contract.GetProtocolFee(&_CallService.CallOpts)
}

// GetProtocolFee is a free data retrieval call binding the contract method 0xa5a41031.
//
// Solidity: function getProtocolFee() view returns(uint256)
func (_CallService *CallServiceCallerSession) GetProtocolFee() (*big.Int, error) {
	return _CallService.Contract.GetProtocolFee(&_CallService.CallOpts)
}

// GetProtocolFeeHandler is a free data retrieval call binding the contract method 0x2eb71414.
//
// Solidity: function getProtocolFeeHandler() view returns(address)
func (_CallService *CallServiceCaller) GetProtocolFeeHandler(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CallService.contract.Call(opts, &out, "getProtocolFeeHandler")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetProtocolFeeHandler is a free data retrieval call binding the contract method 0x2eb71414.
//
// Solidity: function getProtocolFeeHandler() view returns(address)
func (_CallService *CallServiceSession) GetProtocolFeeHandler() (common.Address, error) {
	return _CallService.Contract.GetProtocolFeeHandler(&_CallService.CallOpts)
}

// GetProtocolFeeHandler is a free data retrieval call binding the contract method 0x2eb71414.
//
// Solidity: function getProtocolFeeHandler() view returns(address)
func (_CallService *CallServiceCallerSession) GetProtocolFeeHandler() (common.Address, error) {
	return _CallService.Contract.GetProtocolFeeHandler(&_CallService.CallOpts)
}

// VerifySuccess is a free data retrieval call binding the contract method 0xec05386b.
//
// Solidity: function verifySuccess(uint256 _sn) view returns(bool)
func (_CallService *CallServiceCaller) VerifySuccess(opts *bind.CallOpts, _sn *big.Int) (bool, error) {
	var out []interface{}
	err := _CallService.contract.Call(opts, &out, "verifySuccess", _sn)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VerifySuccess is a free data retrieval call binding the contract method 0xec05386b.
//
// Solidity: function verifySuccess(uint256 _sn) view returns(bool)
func (_CallService *CallServiceSession) VerifySuccess(_sn *big.Int) (bool, error) {
	return _CallService.Contract.VerifySuccess(&_CallService.CallOpts, _sn)
}

// VerifySuccess is a free data retrieval call binding the contract method 0xec05386b.
//
// Solidity: function verifySuccess(uint256 _sn) view returns(bool)
func (_CallService *CallServiceCallerSession) VerifySuccess(_sn *big.Int) (bool, error) {
	return _CallService.Contract.VerifySuccess(&_CallService.CallOpts, _sn)
}

// ExecuteCall is a paid mutator transaction binding the contract method 0xbda8ce21.
//
// Solidity: function executeCall(uint256 _reqId, bytes _data) returns()
func (_CallService *CallServiceTransactor) ExecuteCall(opts *bind.TransactOpts, _reqId *big.Int, _data []byte) (*types.Transaction, error) {
	return _CallService.contract.Transact(opts, "executeCall", _reqId, _data)
}

// ExecuteCall is a paid mutator transaction binding the contract method 0xbda8ce21.
//
// Solidity: function executeCall(uint256 _reqId, bytes _data) returns()
func (_CallService *CallServiceSession) ExecuteCall(_reqId *big.Int, _data []byte) (*types.Transaction, error) {
	return _CallService.Contract.ExecuteCall(&_CallService.TransactOpts, _reqId, _data)
}

// ExecuteCall is a paid mutator transaction binding the contract method 0xbda8ce21.
//
// Solidity: function executeCall(uint256 _reqId, bytes _data) returns()
func (_CallService *CallServiceTransactorSession) ExecuteCall(_reqId *big.Int, _data []byte) (*types.Transaction, error) {
	return _CallService.Contract.ExecuteCall(&_CallService.TransactOpts, _reqId, _data)
}

// ExecuteMessage is a paid mutator transaction binding the contract method 0x313bf398.
//
// Solidity: function executeMessage(address to, string from, bytes data, string[] protocols) returns()
func (_CallService *CallServiceTransactor) ExecuteMessage(opts *bind.TransactOpts, to common.Address, from string, data []byte, protocols []string) (*types.Transaction, error) {
	return _CallService.contract.Transact(opts, "executeMessage", to, from, data, protocols)
}

// ExecuteMessage is a paid mutator transaction binding the contract method 0x313bf398.
//
// Solidity: function executeMessage(address to, string from, bytes data, string[] protocols) returns()
func (_CallService *CallServiceSession) ExecuteMessage(to common.Address, from string, data []byte, protocols []string) (*types.Transaction, error) {
	return _CallService.Contract.ExecuteMessage(&_CallService.TransactOpts, to, from, data, protocols)
}

// ExecuteMessage is a paid mutator transaction binding the contract method 0x313bf398.
//
// Solidity: function executeMessage(address to, string from, bytes data, string[] protocols) returns()
func (_CallService *CallServiceTransactorSession) ExecuteMessage(to common.Address, from string, data []byte, protocols []string) (*types.Transaction, error) {
	return _CallService.Contract.ExecuteMessage(&_CallService.TransactOpts, to, from, data, protocols)
}

// ExecuteRollback is a paid mutator transaction binding the contract method 0x2a84e1b0.
//
// Solidity: function executeRollback(uint256 _sn) returns()
func (_CallService *CallServiceTransactor) ExecuteRollback(opts *bind.TransactOpts, _sn *big.Int) (*types.Transaction, error) {
	return _CallService.contract.Transact(opts, "executeRollback", _sn)
}

// ExecuteRollback is a paid mutator transaction binding the contract method 0x2a84e1b0.
//
// Solidity: function executeRollback(uint256 _sn) returns()
func (_CallService *CallServiceSession) ExecuteRollback(_sn *big.Int) (*types.Transaction, error) {
	return _CallService.Contract.ExecuteRollback(&_CallService.TransactOpts, _sn)
}

// ExecuteRollback is a paid mutator transaction binding the contract method 0x2a84e1b0.
//
// Solidity: function executeRollback(uint256 _sn) returns()
func (_CallService *CallServiceTransactorSession) ExecuteRollback(_sn *big.Int) (*types.Transaction, error) {
	return _CallService.Contract.ExecuteRollback(&_CallService.TransactOpts, _sn)
}

// HandleBTPError is a paid mutator transaction binding the contract method 0x0a823dea.
//
// Solidity: function handleBTPError(string _src, string _svc, uint256 _sn, uint256 _code, string _msg) returns()
func (_CallService *CallServiceTransactor) HandleBTPError(opts *bind.TransactOpts, _src string, _svc string, _sn *big.Int, _code *big.Int, _msg string) (*types.Transaction, error) {
	return _CallService.contract.Transact(opts, "handleBTPError", _src, _svc, _sn, _code, _msg)
}

// HandleBTPError is a paid mutator transaction binding the contract method 0x0a823dea.
//
// Solidity: function handleBTPError(string _src, string _svc, uint256 _sn, uint256 _code, string _msg) returns()
func (_CallService *CallServiceSession) HandleBTPError(_src string, _svc string, _sn *big.Int, _code *big.Int, _msg string) (*types.Transaction, error) {
	return _CallService.Contract.HandleBTPError(&_CallService.TransactOpts, _src, _svc, _sn, _code, _msg)
}

// HandleBTPError is a paid mutator transaction binding the contract method 0x0a823dea.
//
// Solidity: function handleBTPError(string _src, string _svc, uint256 _sn, uint256 _code, string _msg) returns()
func (_CallService *CallServiceTransactorSession) HandleBTPError(_src string, _svc string, _sn *big.Int, _code *big.Int, _msg string) (*types.Transaction, error) {
	return _CallService.Contract.HandleBTPError(&_CallService.TransactOpts, _src, _svc, _sn, _code, _msg)
}

// HandleBTPMessage is a paid mutator transaction binding the contract method 0xb70eeb8d.
//
// Solidity: function handleBTPMessage(string _from, string _svc, uint256 _sn, bytes _msg) returns()
func (_CallService *CallServiceTransactor) HandleBTPMessage(opts *bind.TransactOpts, _from string, _svc string, _sn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _CallService.contract.Transact(opts, "handleBTPMessage", _from, _svc, _sn, _msg)
}

// HandleBTPMessage is a paid mutator transaction binding the contract method 0xb70eeb8d.
//
// Solidity: function handleBTPMessage(string _from, string _svc, uint256 _sn, bytes _msg) returns()
func (_CallService *CallServiceSession) HandleBTPMessage(_from string, _svc string, _sn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _CallService.Contract.HandleBTPMessage(&_CallService.TransactOpts, _from, _svc, _sn, _msg)
}

// HandleBTPMessage is a paid mutator transaction binding the contract method 0xb70eeb8d.
//
// Solidity: function handleBTPMessage(string _from, string _svc, uint256 _sn, bytes _msg) returns()
func (_CallService *CallServiceTransactorSession) HandleBTPMessage(_from string, _svc string, _sn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _CallService.Contract.HandleBTPMessage(&_CallService.TransactOpts, _from, _svc, _sn, _msg)
}

// HandleError is a paid mutator transaction binding the contract method 0xb070f9e5.
//
// Solidity: function handleError(uint256 _sn) returns()
func (_CallService *CallServiceTransactor) HandleError(opts *bind.TransactOpts, _sn *big.Int) (*types.Transaction, error) {
	return _CallService.contract.Transact(opts, "handleError", _sn)
}

// HandleError is a paid mutator transaction binding the contract method 0xb070f9e5.
//
// Solidity: function handleError(uint256 _sn) returns()
func (_CallService *CallServiceSession) HandleError(_sn *big.Int) (*types.Transaction, error) {
	return _CallService.Contract.HandleError(&_CallService.TransactOpts, _sn)
}

// HandleError is a paid mutator transaction binding the contract method 0xb070f9e5.
//
// Solidity: function handleError(uint256 _sn) returns()
func (_CallService *CallServiceTransactorSession) HandleError(_sn *big.Int) (*types.Transaction, error) {
	return _CallService.Contract.HandleError(&_CallService.TransactOpts, _sn)
}

// HandleMessage is a paid mutator transaction binding the contract method 0xbbc22efd.
//
// Solidity: function handleMessage(string _from, bytes _msg) returns()
func (_CallService *CallServiceTransactor) HandleMessage(opts *bind.TransactOpts, _from string, _msg []byte) (*types.Transaction, error) {
	return _CallService.contract.Transact(opts, "handleMessage", _from, _msg)
}

// HandleMessage is a paid mutator transaction binding the contract method 0xbbc22efd.
//
// Solidity: function handleMessage(string _from, bytes _msg) returns()
func (_CallService *CallServiceSession) HandleMessage(_from string, _msg []byte) (*types.Transaction, error) {
	return _CallService.Contract.HandleMessage(&_CallService.TransactOpts, _from, _msg)
}

// HandleMessage is a paid mutator transaction binding the contract method 0xbbc22efd.
//
// Solidity: function handleMessage(string _from, bytes _msg) returns()
func (_CallService *CallServiceTransactorSession) HandleMessage(_from string, _msg []byte) (*types.Transaction, error) {
	return _CallService.Contract.HandleMessage(&_CallService.TransactOpts, _from, _msg)
}

// Initialize is a paid mutator transaction binding the contract method 0xf62d1888.
//
// Solidity: function initialize(string _nid) returns()
func (_CallService *CallServiceTransactor) Initialize(opts *bind.TransactOpts, _nid string) (*types.Transaction, error) {
	return _CallService.contract.Transact(opts, "initialize", _nid)
}

// Initialize is a paid mutator transaction binding the contract method 0xf62d1888.
//
// Solidity: function initialize(string _nid) returns()
func (_CallService *CallServiceSession) Initialize(_nid string) (*types.Transaction, error) {
	return _CallService.Contract.Initialize(&_CallService.TransactOpts, _nid)
}

// Initialize is a paid mutator transaction binding the contract method 0xf62d1888.
//
// Solidity: function initialize(string _nid) returns()
func (_CallService *CallServiceTransactorSession) Initialize(_nid string) (*types.Transaction, error) {
	return _CallService.Contract.Initialize(&_CallService.TransactOpts, _nid)
}

// SendCall is a paid mutator transaction binding the contract method 0x17fd7a33.
//
// Solidity: function sendCall(string _to, bytes _data) payable returns(uint256)
func (_CallService *CallServiceTransactor) SendCall(opts *bind.TransactOpts, _to string, _data []byte) (*types.Transaction, error) {
	return _CallService.contract.Transact(opts, "sendCall", _to, _data)
}

// SendCall is a paid mutator transaction binding the contract method 0x17fd7a33.
//
// Solidity: function sendCall(string _to, bytes _data) payable returns(uint256)
func (_CallService *CallServiceSession) SendCall(_to string, _data []byte) (*types.Transaction, error) {
	return _CallService.Contract.SendCall(&_CallService.TransactOpts, _to, _data)
}

// SendCall is a paid mutator transaction binding the contract method 0x17fd7a33.
//
// Solidity: function sendCall(string _to, bytes _data) payable returns(uint256)
func (_CallService *CallServiceTransactorSession) SendCall(_to string, _data []byte) (*types.Transaction, error) {
	return _CallService.Contract.SendCall(&_CallService.TransactOpts, _to, _data)
}

// SendCallMessage is a paid mutator transaction binding the contract method 0x8ef378b8.
//
// Solidity: function sendCallMessage(string _to, bytes _data, bytes _rollback) payable returns(uint256)
func (_CallService *CallServiceTransactor) SendCallMessage(opts *bind.TransactOpts, _to string, _data []byte, _rollback []byte) (*types.Transaction, error) {
	return _CallService.contract.Transact(opts, "sendCallMessage", _to, _data, _rollback)
}

// SendCallMessage is a paid mutator transaction binding the contract method 0x8ef378b8.
//
// Solidity: function sendCallMessage(string _to, bytes _data, bytes _rollback) payable returns(uint256)
func (_CallService *CallServiceSession) SendCallMessage(_to string, _data []byte, _rollback []byte) (*types.Transaction, error) {
	return _CallService.Contract.SendCallMessage(&_CallService.TransactOpts, _to, _data, _rollback)
}

// SendCallMessage is a paid mutator transaction binding the contract method 0x8ef378b8.
//
// Solidity: function sendCallMessage(string _to, bytes _data, bytes _rollback) payable returns(uint256)
func (_CallService *CallServiceTransactorSession) SendCallMessage(_to string, _data []byte, _rollback []byte) (*types.Transaction, error) {
	return _CallService.Contract.SendCallMessage(&_CallService.TransactOpts, _to, _data, _rollback)
}

// SendCallMessage0 is a paid mutator transaction binding the contract method 0xedc6afff.
//
// Solidity: function sendCallMessage(string _to, bytes _data, bytes _rollback, string[] sources, string[] destinations) payable returns(uint256)
func (_CallService *CallServiceTransactor) SendCallMessage0(opts *bind.TransactOpts, _to string, _data []byte, _rollback []byte, sources []string, destinations []string) (*types.Transaction, error) {
	return _CallService.contract.Transact(opts, "sendCallMessage0", _to, _data, _rollback, sources, destinations)
}

// SendCallMessage0 is a paid mutator transaction binding the contract method 0xedc6afff.
//
// Solidity: function sendCallMessage(string _to, bytes _data, bytes _rollback, string[] sources, string[] destinations) payable returns(uint256)
func (_CallService *CallServiceSession) SendCallMessage0(_to string, _data []byte, _rollback []byte, sources []string, destinations []string) (*types.Transaction, error) {
	return _CallService.Contract.SendCallMessage0(&_CallService.TransactOpts, _to, _data, _rollback, sources, destinations)
}

// SendCallMessage0 is a paid mutator transaction binding the contract method 0xedc6afff.
//
// Solidity: function sendCallMessage(string _to, bytes _data, bytes _rollback, string[] sources, string[] destinations) payable returns(uint256)
func (_CallService *CallServiceTransactorSession) SendCallMessage0(_to string, _data []byte, _rollback []byte, sources []string, destinations []string) (*types.Transaction, error) {
	return _CallService.Contract.SendCallMessage0(&_CallService.TransactOpts, _to, _data, _rollback, sources, destinations)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x704b6c02.
//
// Solidity: function setAdmin(address _address) returns()
func (_CallService *CallServiceTransactor) SetAdmin(opts *bind.TransactOpts, _address common.Address) (*types.Transaction, error) {
	return _CallService.contract.Transact(opts, "setAdmin", _address)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x704b6c02.
//
// Solidity: function setAdmin(address _address) returns()
func (_CallService *CallServiceSession) SetAdmin(_address common.Address) (*types.Transaction, error) {
	return _CallService.Contract.SetAdmin(&_CallService.TransactOpts, _address)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x704b6c02.
//
// Solidity: function setAdmin(address _address) returns()
func (_CallService *CallServiceTransactorSession) SetAdmin(_address common.Address) (*types.Transaction, error) {
	return _CallService.Contract.SetAdmin(&_CallService.TransactOpts, _address)
}

// SetDefaultConnection is a paid mutator transaction binding the contract method 0x64f03757.
//
// Solidity: function setDefaultConnection(string _nid, address connection) returns()
func (_CallService *CallServiceTransactor) SetDefaultConnection(opts *bind.TransactOpts, _nid string, connection common.Address) (*types.Transaction, error) {
	return _CallService.contract.Transact(opts, "setDefaultConnection", _nid, connection)
}

// SetDefaultConnection is a paid mutator transaction binding the contract method 0x64f03757.
//
// Solidity: function setDefaultConnection(string _nid, address connection) returns()
func (_CallService *CallServiceSession) SetDefaultConnection(_nid string, connection common.Address) (*types.Transaction, error) {
	return _CallService.Contract.SetDefaultConnection(&_CallService.TransactOpts, _nid, connection)
}

// SetDefaultConnection is a paid mutator transaction binding the contract method 0x64f03757.
//
// Solidity: function setDefaultConnection(string _nid, address connection) returns()
func (_CallService *CallServiceTransactorSession) SetDefaultConnection(_nid string, connection common.Address) (*types.Transaction, error) {
	return _CallService.Contract.SetDefaultConnection(&_CallService.TransactOpts, _nid, connection)
}

// SetProtocolFee is a paid mutator transaction binding the contract method 0x787dce3d.
//
// Solidity: function setProtocolFee(uint256 _value) returns()
func (_CallService *CallServiceTransactor) SetProtocolFee(opts *bind.TransactOpts, _value *big.Int) (*types.Transaction, error) {
	return _CallService.contract.Transact(opts, "setProtocolFee", _value)
}

// SetProtocolFee is a paid mutator transaction binding the contract method 0x787dce3d.
//
// Solidity: function setProtocolFee(uint256 _value) returns()
func (_CallService *CallServiceSession) SetProtocolFee(_value *big.Int) (*types.Transaction, error) {
	return _CallService.Contract.SetProtocolFee(&_CallService.TransactOpts, _value)
}

// SetProtocolFee is a paid mutator transaction binding the contract method 0x787dce3d.
//
// Solidity: function setProtocolFee(uint256 _value) returns()
func (_CallService *CallServiceTransactorSession) SetProtocolFee(_value *big.Int) (*types.Transaction, error) {
	return _CallService.Contract.SetProtocolFee(&_CallService.TransactOpts, _value)
}

// SetProtocolFeeHandler is a paid mutator transaction binding the contract method 0x502bf8e0.
//
// Solidity: function setProtocolFeeHandler(address _addr) returns()
func (_CallService *CallServiceTransactor) SetProtocolFeeHandler(opts *bind.TransactOpts, _addr common.Address) (*types.Transaction, error) {
	return _CallService.contract.Transact(opts, "setProtocolFeeHandler", _addr)
}

// SetProtocolFeeHandler is a paid mutator transaction binding the contract method 0x502bf8e0.
//
// Solidity: function setProtocolFeeHandler(address _addr) returns()
func (_CallService *CallServiceSession) SetProtocolFeeHandler(_addr common.Address) (*types.Transaction, error) {
	return _CallService.Contract.SetProtocolFeeHandler(&_CallService.TransactOpts, _addr)
}

// SetProtocolFeeHandler is a paid mutator transaction binding the contract method 0x502bf8e0.
//
// Solidity: function setProtocolFeeHandler(address _addr) returns()
func (_CallService *CallServiceTransactorSession) SetProtocolFeeHandler(_addr common.Address) (*types.Transaction, error) {
	return _CallService.Contract.SetProtocolFeeHandler(&_CallService.TransactOpts, _addr)
}

// CallServiceCallExecutedIterator is returned from FilterCallExecuted and is used to iterate over the raw logs and unpacked data for CallExecuted events raised by the CallService contract.
type CallServiceCallExecutedIterator struct {
	Event *CallServiceCallExecuted // Event containing the contract specifics and raw log

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
func (it *CallServiceCallExecutedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CallServiceCallExecuted)
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
		it.Event = new(CallServiceCallExecuted)
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
func (it *CallServiceCallExecutedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CallServiceCallExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CallServiceCallExecuted represents a CallExecuted event raised by the CallService contract.
type CallServiceCallExecuted struct {
	ReqId *big.Int
	Code  *big.Int
	Msg   string
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterCallExecuted is a free log retrieval operation binding the contract event 0xc7391e04887f8b3c16fa20877e028e8163139a478c8447e7d449eba1905caa51.
//
// Solidity: event CallExecuted(uint256 indexed _reqId, int256 _code, string _msg)
func (_CallService *CallServiceFilterer) FilterCallExecuted(opts *bind.FilterOpts, _reqId []*big.Int) (*CallServiceCallExecutedIterator, error) {

	var _reqIdRule []interface{}
	for _, _reqIdItem := range _reqId {
		_reqIdRule = append(_reqIdRule, _reqIdItem)
	}

	logs, sub, err := _CallService.contract.FilterLogs(opts, "CallExecuted", _reqIdRule)
	if err != nil {
		return nil, err
	}
	return &CallServiceCallExecutedIterator{contract: _CallService.contract, event: "CallExecuted", logs: logs, sub: sub}, nil
}

// WatchCallExecuted is a free log subscription operation binding the contract event 0xc7391e04887f8b3c16fa20877e028e8163139a478c8447e7d449eba1905caa51.
//
// Solidity: event CallExecuted(uint256 indexed _reqId, int256 _code, string _msg)
func (_CallService *CallServiceFilterer) WatchCallExecuted(opts *bind.WatchOpts, sink chan<- *CallServiceCallExecuted, _reqId []*big.Int) (event.Subscription, error) {

	var _reqIdRule []interface{}
	for _, _reqIdItem := range _reqId {
		_reqIdRule = append(_reqIdRule, _reqIdItem)
	}

	logs, sub, err := _CallService.contract.WatchLogs(opts, "CallExecuted", _reqIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CallServiceCallExecuted)
				if err := _CallService.contract.UnpackLog(event, "CallExecuted", log); err != nil {
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

// ParseCallExecuted is a log parse operation binding the contract event 0xc7391e04887f8b3c16fa20877e028e8163139a478c8447e7d449eba1905caa51.
//
// Solidity: event CallExecuted(uint256 indexed _reqId, int256 _code, string _msg)
func (_CallService *CallServiceFilterer) ParseCallExecuted(log types.Log) (*CallServiceCallExecuted, error) {
	event := new(CallServiceCallExecuted)
	if err := _CallService.contract.UnpackLog(event, "CallExecuted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CallServiceCallMessageIterator is returned from FilterCallMessage and is used to iterate over the raw logs and unpacked data for CallMessage events raised by the CallService contract.
type CallServiceCallMessageIterator struct {
	Event *CallServiceCallMessage // Event containing the contract specifics and raw log

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
func (it *CallServiceCallMessageIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CallServiceCallMessage)
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
		it.Event = new(CallServiceCallMessage)
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
func (it *CallServiceCallMessageIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CallServiceCallMessageIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CallServiceCallMessage represents a CallMessage event raised by the CallService contract.
type CallServiceCallMessage struct {
	From  common.Hash
	To    common.Hash
	Sn    *big.Int
	ReqId *big.Int
	Data  []byte
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterCallMessage is a free log retrieval operation binding the contract event 0x2cbc78425621c181f9f8a25fc06e44a0ac2b67cd6a31f8ed7918934187f8cc59.
//
// Solidity: event CallMessage(string indexed _from, string indexed _to, uint256 indexed _sn, uint256 _reqId, bytes _data)
func (_CallService *CallServiceFilterer) FilterCallMessage(opts *bind.FilterOpts, _from []string, _to []string, _sn []*big.Int) (*CallServiceCallMessageIterator, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}
	var _snRule []interface{}
	for _, _snItem := range _sn {
		_snRule = append(_snRule, _snItem)
	}

	logs, sub, err := _CallService.contract.FilterLogs(opts, "CallMessage", _fromRule, _toRule, _snRule)
	if err != nil {
		return nil, err
	}
	return &CallServiceCallMessageIterator{contract: _CallService.contract, event: "CallMessage", logs: logs, sub: sub}, nil
}

// WatchCallMessage is a free log subscription operation binding the contract event 0x2cbc78425621c181f9f8a25fc06e44a0ac2b67cd6a31f8ed7918934187f8cc59.
//
// Solidity: event CallMessage(string indexed _from, string indexed _to, uint256 indexed _sn, uint256 _reqId, bytes _data)
func (_CallService *CallServiceFilterer) WatchCallMessage(opts *bind.WatchOpts, sink chan<- *CallServiceCallMessage, _from []string, _to []string, _sn []*big.Int) (event.Subscription, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}
	var _snRule []interface{}
	for _, _snItem := range _sn {
		_snRule = append(_snRule, _snItem)
	}

	logs, sub, err := _CallService.contract.WatchLogs(opts, "CallMessage", _fromRule, _toRule, _snRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CallServiceCallMessage)
				if err := _CallService.contract.UnpackLog(event, "CallMessage", log); err != nil {
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

// ParseCallMessage is a log parse operation binding the contract event 0x2cbc78425621c181f9f8a25fc06e44a0ac2b67cd6a31f8ed7918934187f8cc59.
//
// Solidity: event CallMessage(string indexed _from, string indexed _to, uint256 indexed _sn, uint256 _reqId, bytes _data)
func (_CallService *CallServiceFilterer) ParseCallMessage(log types.Log) (*CallServiceCallMessage, error) {
	event := new(CallServiceCallMessage)
	if err := _CallService.contract.UnpackLog(event, "CallMessage", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CallServiceCallMessageSentIterator is returned from FilterCallMessageSent and is used to iterate over the raw logs and unpacked data for CallMessageSent events raised by the CallService contract.
type CallServiceCallMessageSentIterator struct {
	Event *CallServiceCallMessageSent // Event containing the contract specifics and raw log

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
func (it *CallServiceCallMessageSentIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CallServiceCallMessageSent)
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
		it.Event = new(CallServiceCallMessageSent)
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
func (it *CallServiceCallMessageSentIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CallServiceCallMessageSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CallServiceCallMessageSent represents a CallMessageSent event raised by the CallService contract.
type CallServiceCallMessageSent struct {
	From common.Address
	To   common.Hash
	Sn   *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterCallMessageSent is a free log retrieval operation binding the contract event 0x69e53ea70fdf945f6d035b3979748bc999151691fb1dc69d66f8017f8840ae28.
//
// Solidity: event CallMessageSent(address indexed _from, string indexed _to, uint256 indexed _sn)
func (_CallService *CallServiceFilterer) FilterCallMessageSent(opts *bind.FilterOpts, _from []common.Address, _to []string, _sn []*big.Int) (*CallServiceCallMessageSentIterator, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}
	var _snRule []interface{}
	for _, _snItem := range _sn {
		_snRule = append(_snRule, _snItem)
	}

	logs, sub, err := _CallService.contract.FilterLogs(opts, "CallMessageSent", _fromRule, _toRule, _snRule)
	if err != nil {
		return nil, err
	}
	return &CallServiceCallMessageSentIterator{contract: _CallService.contract, event: "CallMessageSent", logs: logs, sub: sub}, nil
}

// WatchCallMessageSent is a free log subscription operation binding the contract event 0x69e53ea70fdf945f6d035b3979748bc999151691fb1dc69d66f8017f8840ae28.
//
// Solidity: event CallMessageSent(address indexed _from, string indexed _to, uint256 indexed _sn)
func (_CallService *CallServiceFilterer) WatchCallMessageSent(opts *bind.WatchOpts, sink chan<- *CallServiceCallMessageSent, _from []common.Address, _to []string, _sn []*big.Int) (event.Subscription, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}
	var _snRule []interface{}
	for _, _snItem := range _sn {
		_snRule = append(_snRule, _snItem)
	}

	logs, sub, err := _CallService.contract.WatchLogs(opts, "CallMessageSent", _fromRule, _toRule, _snRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CallServiceCallMessageSent)
				if err := _CallService.contract.UnpackLog(event, "CallMessageSent", log); err != nil {
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

// ParseCallMessageSent is a log parse operation binding the contract event 0x69e53ea70fdf945f6d035b3979748bc999151691fb1dc69d66f8017f8840ae28.
//
// Solidity: event CallMessageSent(address indexed _from, string indexed _to, uint256 indexed _sn)
func (_CallService *CallServiceFilterer) ParseCallMessageSent(log types.Log) (*CallServiceCallMessageSent, error) {
	event := new(CallServiceCallMessageSent)
	if err := _CallService.contract.UnpackLog(event, "CallMessageSent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CallServiceInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the CallService contract.
type CallServiceInitializedIterator struct {
	Event *CallServiceInitialized // Event containing the contract specifics and raw log

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
func (it *CallServiceInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CallServiceInitialized)
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
		it.Event = new(CallServiceInitialized)
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
func (it *CallServiceInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CallServiceInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CallServiceInitialized represents a Initialized event raised by the CallService contract.
type CallServiceInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_CallService *CallServiceFilterer) FilterInitialized(opts *bind.FilterOpts) (*CallServiceInitializedIterator, error) {

	logs, sub, err := _CallService.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &CallServiceInitializedIterator{contract: _CallService.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_CallService *CallServiceFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *CallServiceInitialized) (event.Subscription, error) {

	logs, sub, err := _CallService.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CallServiceInitialized)
				if err := _CallService.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_CallService *CallServiceFilterer) ParseInitialized(log types.Log) (*CallServiceInitialized, error) {
	event := new(CallServiceInitialized)
	if err := _CallService.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CallServiceResponseMessageIterator is returned from FilterResponseMessage and is used to iterate over the raw logs and unpacked data for ResponseMessage events raised by the CallService contract.
type CallServiceResponseMessageIterator struct {
	Event *CallServiceResponseMessage // Event containing the contract specifics and raw log

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
func (it *CallServiceResponseMessageIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CallServiceResponseMessage)
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
		it.Event = new(CallServiceResponseMessage)
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
func (it *CallServiceResponseMessageIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CallServiceResponseMessageIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CallServiceResponseMessage represents a ResponseMessage event raised by the CallService contract.
type CallServiceResponseMessage struct {
	Sn   *big.Int
	Code *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterResponseMessage is a free log retrieval operation binding the contract event 0xbeacafd006c5e60667f6f04aec3a498f81c8e94142b4e95b5a5a763de43ca0ab.
//
// Solidity: event ResponseMessage(uint256 indexed _sn, int256 _code)
func (_CallService *CallServiceFilterer) FilterResponseMessage(opts *bind.FilterOpts, _sn []*big.Int) (*CallServiceResponseMessageIterator, error) {

	var _snRule []interface{}
	for _, _snItem := range _sn {
		_snRule = append(_snRule, _snItem)
	}

	logs, sub, err := _CallService.contract.FilterLogs(opts, "ResponseMessage", _snRule)
	if err != nil {
		return nil, err
	}
	return &CallServiceResponseMessageIterator{contract: _CallService.contract, event: "ResponseMessage", logs: logs, sub: sub}, nil
}

// WatchResponseMessage is a free log subscription operation binding the contract event 0xbeacafd006c5e60667f6f04aec3a498f81c8e94142b4e95b5a5a763de43ca0ab.
//
// Solidity: event ResponseMessage(uint256 indexed _sn, int256 _code)
func (_CallService *CallServiceFilterer) WatchResponseMessage(opts *bind.WatchOpts, sink chan<- *CallServiceResponseMessage, _sn []*big.Int) (event.Subscription, error) {

	var _snRule []interface{}
	for _, _snItem := range _sn {
		_snRule = append(_snRule, _snItem)
	}

	logs, sub, err := _CallService.contract.WatchLogs(opts, "ResponseMessage", _snRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CallServiceResponseMessage)
				if err := _CallService.contract.UnpackLog(event, "ResponseMessage", log); err != nil {
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

// ParseResponseMessage is a log parse operation binding the contract event 0xbeacafd006c5e60667f6f04aec3a498f81c8e94142b4e95b5a5a763de43ca0ab.
//
// Solidity: event ResponseMessage(uint256 indexed _sn, int256 _code)
func (_CallService *CallServiceFilterer) ParseResponseMessage(log types.Log) (*CallServiceResponseMessage, error) {
	event := new(CallServiceResponseMessage)
	if err := _CallService.contract.UnpackLog(event, "ResponseMessage", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CallServiceRollbackExecutedIterator is returned from FilterRollbackExecuted and is used to iterate over the raw logs and unpacked data for RollbackExecuted events raised by the CallService contract.
type CallServiceRollbackExecutedIterator struct {
	Event *CallServiceRollbackExecuted // Event containing the contract specifics and raw log

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
func (it *CallServiceRollbackExecutedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CallServiceRollbackExecuted)
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
		it.Event = new(CallServiceRollbackExecuted)
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
func (it *CallServiceRollbackExecutedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CallServiceRollbackExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CallServiceRollbackExecuted represents a RollbackExecuted event raised by the CallService contract.
type CallServiceRollbackExecuted struct {
	Sn  *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterRollbackExecuted is a free log retrieval operation binding the contract event 0x08f0ac7aef6da8bbe43bee8b1444a1883f1359566618bc379ce5abba44883837.
//
// Solidity: event RollbackExecuted(uint256 indexed _sn)
func (_CallService *CallServiceFilterer) FilterRollbackExecuted(opts *bind.FilterOpts, _sn []*big.Int) (*CallServiceRollbackExecutedIterator, error) {

	var _snRule []interface{}
	for _, _snItem := range _sn {
		_snRule = append(_snRule, _snItem)
	}

	logs, sub, err := _CallService.contract.FilterLogs(opts, "RollbackExecuted", _snRule)
	if err != nil {
		return nil, err
	}
	return &CallServiceRollbackExecutedIterator{contract: _CallService.contract, event: "RollbackExecuted", logs: logs, sub: sub}, nil
}

// WatchRollbackExecuted is a free log subscription operation binding the contract event 0x08f0ac7aef6da8bbe43bee8b1444a1883f1359566618bc379ce5abba44883837.
//
// Solidity: event RollbackExecuted(uint256 indexed _sn)
func (_CallService *CallServiceFilterer) WatchRollbackExecuted(opts *bind.WatchOpts, sink chan<- *CallServiceRollbackExecuted, _sn []*big.Int) (event.Subscription, error) {

	var _snRule []interface{}
	for _, _snItem := range _sn {
		_snRule = append(_snRule, _snItem)
	}

	logs, sub, err := _CallService.contract.WatchLogs(opts, "RollbackExecuted", _snRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CallServiceRollbackExecuted)
				if err := _CallService.contract.UnpackLog(event, "RollbackExecuted", log); err != nil {
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

// ParseRollbackExecuted is a log parse operation binding the contract event 0x08f0ac7aef6da8bbe43bee8b1444a1883f1359566618bc379ce5abba44883837.
//
// Solidity: event RollbackExecuted(uint256 indexed _sn)
func (_CallService *CallServiceFilterer) ParseRollbackExecuted(log types.Log) (*CallServiceRollbackExecuted, error) {
	event := new(CallServiceRollbackExecuted)
	if err := _CallService.contract.UnpackLog(event, "RollbackExecuted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CallServiceRollbackMessageIterator is returned from FilterRollbackMessage and is used to iterate over the raw logs and unpacked data for RollbackMessage events raised by the CallService contract.
type CallServiceRollbackMessageIterator struct {
	Event *CallServiceRollbackMessage // Event containing the contract specifics and raw log

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
func (it *CallServiceRollbackMessageIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CallServiceRollbackMessage)
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
		it.Event = new(CallServiceRollbackMessage)
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
func (it *CallServiceRollbackMessageIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CallServiceRollbackMessageIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CallServiceRollbackMessage represents a RollbackMessage event raised by the CallService contract.
type CallServiceRollbackMessage struct {
	Sn  *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterRollbackMessage is a free log retrieval operation binding the contract event 0x38934ab923f985814047679ba041577b8203ddd15fe9910d3fc6a7aa6001e9c7.
//
// Solidity: event RollbackMessage(uint256 indexed _sn)
func (_CallService *CallServiceFilterer) FilterRollbackMessage(opts *bind.FilterOpts, _sn []*big.Int) (*CallServiceRollbackMessageIterator, error) {

	var _snRule []interface{}
	for _, _snItem := range _sn {
		_snRule = append(_snRule, _snItem)
	}

	logs, sub, err := _CallService.contract.FilterLogs(opts, "RollbackMessage", _snRule)
	if err != nil {
		return nil, err
	}
	return &CallServiceRollbackMessageIterator{contract: _CallService.contract, event: "RollbackMessage", logs: logs, sub: sub}, nil
}

// WatchRollbackMessage is a free log subscription operation binding the contract event 0x38934ab923f985814047679ba041577b8203ddd15fe9910d3fc6a7aa6001e9c7.
//
// Solidity: event RollbackMessage(uint256 indexed _sn)
func (_CallService *CallServiceFilterer) WatchRollbackMessage(opts *bind.WatchOpts, sink chan<- *CallServiceRollbackMessage, _sn []*big.Int) (event.Subscription, error) {

	var _snRule []interface{}
	for _, _snItem := range _sn {
		_snRule = append(_snRule, _snItem)
	}

	logs, sub, err := _CallService.contract.WatchLogs(opts, "RollbackMessage", _snRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CallServiceRollbackMessage)
				if err := _CallService.contract.UnpackLog(event, "RollbackMessage", log); err != nil {
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

// ParseRollbackMessage is a log parse operation binding the contract event 0x38934ab923f985814047679ba041577b8203ddd15fe9910d3fc6a7aa6001e9c7.
//
// Solidity: event RollbackMessage(uint256 indexed _sn)
func (_CallService *CallServiceFilterer) ParseRollbackMessage(log types.Log) (*CallServiceRollbackMessage, error) {
	event := new(CallServiceRollbackMessage)
	if err := _CallService.contract.UnpackLog(event, "RollbackMessage", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
