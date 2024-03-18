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

// XcallMetaData contains all meta data concerning the Xcall contract.
var XcallMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"admin\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"executeCall\",\"inputs\":[{\"name\":\"_reqId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"executeRollback\",\"inputs\":[{\"name\":\"_sn\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getDefaultConnection\",\"inputs\":[{\"name\":\"_nid\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getFee\",\"inputs\":[{\"name\":\"_net\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_rollback\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"_sources\",\"type\":\"string[]\",\"internalType\":\"string[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getFee\",\"inputs\":[{\"name\":\"_net\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_rollback\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNetworkAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNetworkId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getProtocolFee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getProtocolFeeHandler\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"handleBTPError\",\"inputs\":[{\"name\":\"_src\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_svc\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_sn\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_code\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_msg\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"handleBTPMessage\",\"inputs\":[{\"name\":\"_from\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_svc\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_sn\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_msg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"handleError\",\"inputs\":[{\"name\":\"_sn\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"handleMessage\",\"inputs\":[{\"name\":\"_from\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_msg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_nid\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"sendCallMessage\",\"inputs\":[{\"name\":\"_to\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_rollback\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"sendCallMessage\",\"inputs\":[{\"name\":\"_to\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_rollback\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"sources\",\"type\":\"string[]\",\"internalType\":\"string[]\"},{\"name\":\"destinations\",\"type\":\"string[]\",\"internalType\":\"string[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"setAdmin\",\"inputs\":[{\"name\":\"_address\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setDefaultConnection\",\"inputs\":[{\"name\":\"_nid\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"connection\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setProtocolFee\",\"inputs\":[{\"name\":\"_value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setProtocolFeeHandler\",\"inputs\":[{\"name\":\"_addr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"tryHandleCallMessage\",\"inputs\":[{\"name\":\"toAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"from\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"protocols\",\"type\":\"string[]\",\"internalType\":\"string[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"verifySuccess\",\"inputs\":[{\"name\":\"_sn\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"CallExecuted\",\"inputs\":[{\"name\":\"_reqId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"_code\",\"type\":\"int256\",\"indexed\":false,\"internalType\":\"int256\"},{\"name\":\"_msg\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CallMessage\",\"inputs\":[{\"name\":\"_from\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"_to\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"_sn\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"_reqId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"_data\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CallMessageSent\",\"inputs\":[{\"name\":\"_from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"_to\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"_sn\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ResponseMessage\",\"inputs\":[{\"name\":\"_sn\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"_code\",\"type\":\"int256\",\"indexed\":false,\"internalType\":\"int256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RollbackExecuted\",\"inputs\":[{\"name\":\"_sn\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RollbackMessage\",\"inputs\":[{\"name\":\"_sn\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false}]",
}

// XcallABI is the input ABI used to generate the binding from.
// Deprecated: Use XcallMetaData.ABI instead.
var XcallABI = XcallMetaData.ABI

// Xcall is an auto generated Go binding around an Ethereum contract.
type Xcall struct {
	XcallCaller     // Read-only binding to the contract
	XcallTransactor // Write-only binding to the contract
	XcallFilterer   // Log filterer for contract events
}

// XcallCaller is an auto generated read-only Go binding around an Ethereum contract.
type XcallCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// XcallTransactor is an auto generated write-only Go binding around an Ethereum contract.
type XcallTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// XcallFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type XcallFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// XcallSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type XcallSession struct {
	Contract     *Xcall            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// XcallCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type XcallCallerSession struct {
	Contract *XcallCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// XcallTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type XcallTransactorSession struct {
	Contract     *XcallTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// XcallRaw is an auto generated low-level Go binding around an Ethereum contract.
type XcallRaw struct {
	Contract *Xcall // Generic contract binding to access the raw methods on
}

// XcallCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type XcallCallerRaw struct {
	Contract *XcallCaller // Generic read-only contract binding to access the raw methods on
}

// XcallTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type XcallTransactorRaw struct {
	Contract *XcallTransactor // Generic write-only contract binding to access the raw methods on
}

// NewXcall creates a new instance of Xcall, bound to a specific deployed contract.
func NewXcall(address common.Address, backend bind.ContractBackend) (*Xcall, error) {
	contract, err := bindXcall(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Xcall{XcallCaller: XcallCaller{contract: contract}, XcallTransactor: XcallTransactor{contract: contract}, XcallFilterer: XcallFilterer{contract: contract}}, nil
}

// NewXcallCaller creates a new read-only instance of Xcall, bound to a specific deployed contract.
func NewXcallCaller(address common.Address, caller bind.ContractCaller) (*XcallCaller, error) {
	contract, err := bindXcall(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &XcallCaller{contract: contract}, nil
}

// NewXcallTransactor creates a new write-only instance of Xcall, bound to a specific deployed contract.
func NewXcallTransactor(address common.Address, transactor bind.ContractTransactor) (*XcallTransactor, error) {
	contract, err := bindXcall(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &XcallTransactor{contract: contract}, nil
}

// NewXcallFilterer creates a new log filterer instance of Xcall, bound to a specific deployed contract.
func NewXcallFilterer(address common.Address, filterer bind.ContractFilterer) (*XcallFilterer, error) {
	contract, err := bindXcall(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &XcallFilterer{contract: contract}, nil
}

// bindXcall binds a generic wrapper to an already deployed contract.
func bindXcall(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := XcallMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Xcall *XcallRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Xcall.Contract.XcallCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Xcall *XcallRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Xcall.Contract.XcallTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Xcall *XcallRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Xcall.Contract.XcallTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Xcall *XcallCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Xcall.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Xcall *XcallTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Xcall.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Xcall *XcallTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Xcall.Contract.contract.Transact(opts, method, params...)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_Xcall *XcallCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Xcall.contract.Call(opts, &out, "admin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_Xcall *XcallSession) Admin() (common.Address, error) {
	return _Xcall.Contract.Admin(&_Xcall.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_Xcall *XcallCallerSession) Admin() (common.Address, error) {
	return _Xcall.Contract.Admin(&_Xcall.CallOpts)
}

// GetDefaultConnection is a free data retrieval call binding the contract method 0x9e553a4f.
//
// Solidity: function getDefaultConnection(string _nid) view returns(address)
func (_Xcall *XcallCaller) GetDefaultConnection(opts *bind.CallOpts, _nid string) (common.Address, error) {
	var out []interface{}
	err := _Xcall.contract.Call(opts, &out, "getDefaultConnection", _nid)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetDefaultConnection is a free data retrieval call binding the contract method 0x9e553a4f.
//
// Solidity: function getDefaultConnection(string _nid) view returns(address)
func (_Xcall *XcallSession) GetDefaultConnection(_nid string) (common.Address, error) {
	return _Xcall.Contract.GetDefaultConnection(&_Xcall.CallOpts, _nid)
}

// GetDefaultConnection is a free data retrieval call binding the contract method 0x9e553a4f.
//
// Solidity: function getDefaultConnection(string _nid) view returns(address)
func (_Xcall *XcallCallerSession) GetDefaultConnection(_nid string) (common.Address, error) {
	return _Xcall.Contract.GetDefaultConnection(&_Xcall.CallOpts, _nid)
}

// GetFee is a free data retrieval call binding the contract method 0x304a70b5.
//
// Solidity: function getFee(string _net, bool _rollback, string[] _sources) view returns(uint256)
func (_Xcall *XcallCaller) GetFee(opts *bind.CallOpts, _net string, _rollback bool, _sources []string) (*big.Int, error) {
	var out []interface{}
	err := _Xcall.contract.Call(opts, &out, "getFee", _net, _rollback, _sources)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetFee is a free data retrieval call binding the contract method 0x304a70b5.
//
// Solidity: function getFee(string _net, bool _rollback, string[] _sources) view returns(uint256)
func (_Xcall *XcallSession) GetFee(_net string, _rollback bool, _sources []string) (*big.Int, error) {
	return _Xcall.Contract.GetFee(&_Xcall.CallOpts, _net, _rollback, _sources)
}

// GetFee is a free data retrieval call binding the contract method 0x304a70b5.
//
// Solidity: function getFee(string _net, bool _rollback, string[] _sources) view returns(uint256)
func (_Xcall *XcallCallerSession) GetFee(_net string, _rollback bool, _sources []string) (*big.Int, error) {
	return _Xcall.Contract.GetFee(&_Xcall.CallOpts, _net, _rollback, _sources)
}

// GetFee0 is a free data retrieval call binding the contract method 0x7d4c4f4a.
//
// Solidity: function getFee(string _net, bool _rollback) view returns(uint256)
func (_Xcall *XcallCaller) GetFee0(opts *bind.CallOpts, _net string, _rollback bool) (*big.Int, error) {
	var out []interface{}
	err := _Xcall.contract.Call(opts, &out, "getFee0", _net, _rollback)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetFee0 is a free data retrieval call binding the contract method 0x7d4c4f4a.
//
// Solidity: function getFee(string _net, bool _rollback) view returns(uint256)
func (_Xcall *XcallSession) GetFee0(_net string, _rollback bool) (*big.Int, error) {
	return _Xcall.Contract.GetFee0(&_Xcall.CallOpts, _net, _rollback)
}

// GetFee0 is a free data retrieval call binding the contract method 0x7d4c4f4a.
//
// Solidity: function getFee(string _net, bool _rollback) view returns(uint256)
func (_Xcall *XcallCallerSession) GetFee0(_net string, _rollback bool) (*big.Int, error) {
	return _Xcall.Contract.GetFee0(&_Xcall.CallOpts, _net, _rollback)
}

// GetNetworkAddress is a free data retrieval call binding the contract method 0x6bf459cb.
//
// Solidity: function getNetworkAddress() view returns(string)
func (_Xcall *XcallCaller) GetNetworkAddress(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Xcall.contract.Call(opts, &out, "getNetworkAddress")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetNetworkAddress is a free data retrieval call binding the contract method 0x6bf459cb.
//
// Solidity: function getNetworkAddress() view returns(string)
func (_Xcall *XcallSession) GetNetworkAddress() (string, error) {
	return _Xcall.Contract.GetNetworkAddress(&_Xcall.CallOpts)
}

// GetNetworkAddress is a free data retrieval call binding the contract method 0x6bf459cb.
//
// Solidity: function getNetworkAddress() view returns(string)
func (_Xcall *XcallCallerSession) GetNetworkAddress() (string, error) {
	return _Xcall.Contract.GetNetworkAddress(&_Xcall.CallOpts)
}

// GetNetworkId is a free data retrieval call binding the contract method 0x39c5f3fc.
//
// Solidity: function getNetworkId() view returns(string)
func (_Xcall *XcallCaller) GetNetworkId(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Xcall.contract.Call(opts, &out, "getNetworkId")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetNetworkId is a free data retrieval call binding the contract method 0x39c5f3fc.
//
// Solidity: function getNetworkId() view returns(string)
func (_Xcall *XcallSession) GetNetworkId() (string, error) {
	return _Xcall.Contract.GetNetworkId(&_Xcall.CallOpts)
}

// GetNetworkId is a free data retrieval call binding the contract method 0x39c5f3fc.
//
// Solidity: function getNetworkId() view returns(string)
func (_Xcall *XcallCallerSession) GetNetworkId() (string, error) {
	return _Xcall.Contract.GetNetworkId(&_Xcall.CallOpts)
}

// GetProtocolFee is a free data retrieval call binding the contract method 0xa5a41031.
//
// Solidity: function getProtocolFee() view returns(uint256)
func (_Xcall *XcallCaller) GetProtocolFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Xcall.contract.Call(opts, &out, "getProtocolFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetProtocolFee is a free data retrieval call binding the contract method 0xa5a41031.
//
// Solidity: function getProtocolFee() view returns(uint256)
func (_Xcall *XcallSession) GetProtocolFee() (*big.Int, error) {
	return _Xcall.Contract.GetProtocolFee(&_Xcall.CallOpts)
}

// GetProtocolFee is a free data retrieval call binding the contract method 0xa5a41031.
//
// Solidity: function getProtocolFee() view returns(uint256)
func (_Xcall *XcallCallerSession) GetProtocolFee() (*big.Int, error) {
	return _Xcall.Contract.GetProtocolFee(&_Xcall.CallOpts)
}

// GetProtocolFeeHandler is a free data retrieval call binding the contract method 0x2eb71414.
//
// Solidity: function getProtocolFeeHandler() view returns(address)
func (_Xcall *XcallCaller) GetProtocolFeeHandler(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Xcall.contract.Call(opts, &out, "getProtocolFeeHandler")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetProtocolFeeHandler is a free data retrieval call binding the contract method 0x2eb71414.
//
// Solidity: function getProtocolFeeHandler() view returns(address)
func (_Xcall *XcallSession) GetProtocolFeeHandler() (common.Address, error) {
	return _Xcall.Contract.GetProtocolFeeHandler(&_Xcall.CallOpts)
}

// GetProtocolFeeHandler is a free data retrieval call binding the contract method 0x2eb71414.
//
// Solidity: function getProtocolFeeHandler() view returns(address)
func (_Xcall *XcallCallerSession) GetProtocolFeeHandler() (common.Address, error) {
	return _Xcall.Contract.GetProtocolFeeHandler(&_Xcall.CallOpts)
}

// VerifySuccess is a free data retrieval call binding the contract method 0xec05386b.
//
// Solidity: function verifySuccess(uint256 _sn) view returns(bool)
func (_Xcall *XcallCaller) VerifySuccess(opts *bind.CallOpts, _sn *big.Int) (bool, error) {
	var out []interface{}
	err := _Xcall.contract.Call(opts, &out, "verifySuccess", _sn)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VerifySuccess is a free data retrieval call binding the contract method 0xec05386b.
//
// Solidity: function verifySuccess(uint256 _sn) view returns(bool)
func (_Xcall *XcallSession) VerifySuccess(_sn *big.Int) (bool, error) {
	return _Xcall.Contract.VerifySuccess(&_Xcall.CallOpts, _sn)
}

// VerifySuccess is a free data retrieval call binding the contract method 0xec05386b.
//
// Solidity: function verifySuccess(uint256 _sn) view returns(bool)
func (_Xcall *XcallCallerSession) VerifySuccess(_sn *big.Int) (bool, error) {
	return _Xcall.Contract.VerifySuccess(&_Xcall.CallOpts, _sn)
}

// ExecuteCall is a paid mutator transaction binding the contract method 0xbda8ce21.
//
// Solidity: function executeCall(uint256 _reqId, bytes _data) returns()
func (_Xcall *XcallTransactor) ExecuteCall(opts *bind.TransactOpts, _reqId *big.Int, _data []byte) (*types.Transaction, error) {
	return _Xcall.contract.Transact(opts, "executeCall", _reqId, _data)
}

// ExecuteCall is a paid mutator transaction binding the contract method 0xbda8ce21.
//
// Solidity: function executeCall(uint256 _reqId, bytes _data) returns()
func (_Xcall *XcallSession) ExecuteCall(_reqId *big.Int, _data []byte) (*types.Transaction, error) {
	return _Xcall.Contract.ExecuteCall(&_Xcall.TransactOpts, _reqId, _data)
}

// ExecuteCall is a paid mutator transaction binding the contract method 0xbda8ce21.
//
// Solidity: function executeCall(uint256 _reqId, bytes _data) returns()
func (_Xcall *XcallTransactorSession) ExecuteCall(_reqId *big.Int, _data []byte) (*types.Transaction, error) {
	return _Xcall.Contract.ExecuteCall(&_Xcall.TransactOpts, _reqId, _data)
}

// ExecuteRollback is a paid mutator transaction binding the contract method 0x2a84e1b0.
//
// Solidity: function executeRollback(uint256 _sn) returns()
func (_Xcall *XcallTransactor) ExecuteRollback(opts *bind.TransactOpts, _sn *big.Int) (*types.Transaction, error) {
	return _Xcall.contract.Transact(opts, "executeRollback", _sn)
}

// ExecuteRollback is a paid mutator transaction binding the contract method 0x2a84e1b0.
//
// Solidity: function executeRollback(uint256 _sn) returns()
func (_Xcall *XcallSession) ExecuteRollback(_sn *big.Int) (*types.Transaction, error) {
	return _Xcall.Contract.ExecuteRollback(&_Xcall.TransactOpts, _sn)
}

// ExecuteRollback is a paid mutator transaction binding the contract method 0x2a84e1b0.
//
// Solidity: function executeRollback(uint256 _sn) returns()
func (_Xcall *XcallTransactorSession) ExecuteRollback(_sn *big.Int) (*types.Transaction, error) {
	return _Xcall.Contract.ExecuteRollback(&_Xcall.TransactOpts, _sn)
}

// HandleBTPError is a paid mutator transaction binding the contract method 0x0a823dea.
//
// Solidity: function handleBTPError(string _src, string _svc, uint256 _sn, uint256 _code, string _msg) returns()
func (_Xcall *XcallTransactor) HandleBTPError(opts *bind.TransactOpts, _src string, _svc string, _sn *big.Int, _code *big.Int, _msg string) (*types.Transaction, error) {
	return _Xcall.contract.Transact(opts, "handleBTPError", _src, _svc, _sn, _code, _msg)
}

// HandleBTPError is a paid mutator transaction binding the contract method 0x0a823dea.
//
// Solidity: function handleBTPError(string _src, string _svc, uint256 _sn, uint256 _code, string _msg) returns()
func (_Xcall *XcallSession) HandleBTPError(_src string, _svc string, _sn *big.Int, _code *big.Int, _msg string) (*types.Transaction, error) {
	return _Xcall.Contract.HandleBTPError(&_Xcall.TransactOpts, _src, _svc, _sn, _code, _msg)
}

// HandleBTPError is a paid mutator transaction binding the contract method 0x0a823dea.
//
// Solidity: function handleBTPError(string _src, string _svc, uint256 _sn, uint256 _code, string _msg) returns()
func (_Xcall *XcallTransactorSession) HandleBTPError(_src string, _svc string, _sn *big.Int, _code *big.Int, _msg string) (*types.Transaction, error) {
	return _Xcall.Contract.HandleBTPError(&_Xcall.TransactOpts, _src, _svc, _sn, _code, _msg)
}

// HandleBTPMessage is a paid mutator transaction binding the contract method 0xb70eeb8d.
//
// Solidity: function handleBTPMessage(string _from, string _svc, uint256 _sn, bytes _msg) returns()
func (_Xcall *XcallTransactor) HandleBTPMessage(opts *bind.TransactOpts, _from string, _svc string, _sn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _Xcall.contract.Transact(opts, "handleBTPMessage", _from, _svc, _sn, _msg)
}

// HandleBTPMessage is a paid mutator transaction binding the contract method 0xb70eeb8d.
//
// Solidity: function handleBTPMessage(string _from, string _svc, uint256 _sn, bytes _msg) returns()
func (_Xcall *XcallSession) HandleBTPMessage(_from string, _svc string, _sn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _Xcall.Contract.HandleBTPMessage(&_Xcall.TransactOpts, _from, _svc, _sn, _msg)
}

// HandleBTPMessage is a paid mutator transaction binding the contract method 0xb70eeb8d.
//
// Solidity: function handleBTPMessage(string _from, string _svc, uint256 _sn, bytes _msg) returns()
func (_Xcall *XcallTransactorSession) HandleBTPMessage(_from string, _svc string, _sn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _Xcall.Contract.HandleBTPMessage(&_Xcall.TransactOpts, _from, _svc, _sn, _msg)
}

// HandleError is a paid mutator transaction binding the contract method 0xb070f9e5.
//
// Solidity: function handleError(uint256 _sn) returns()
func (_Xcall *XcallTransactor) HandleError(opts *bind.TransactOpts, _sn *big.Int) (*types.Transaction, error) {
	return _Xcall.contract.Transact(opts, "handleError", _sn)
}

// HandleError is a paid mutator transaction binding the contract method 0xb070f9e5.
//
// Solidity: function handleError(uint256 _sn) returns()
func (_Xcall *XcallSession) HandleError(_sn *big.Int) (*types.Transaction, error) {
	return _Xcall.Contract.HandleError(&_Xcall.TransactOpts, _sn)
}

// HandleError is a paid mutator transaction binding the contract method 0xb070f9e5.
//
// Solidity: function handleError(uint256 _sn) returns()
func (_Xcall *XcallTransactorSession) HandleError(_sn *big.Int) (*types.Transaction, error) {
	return _Xcall.Contract.HandleError(&_Xcall.TransactOpts, _sn)
}

// HandleMessage is a paid mutator transaction binding the contract method 0xbbc22efd.
//
// Solidity: function handleMessage(string _from, bytes _msg) returns()
func (_Xcall *XcallTransactor) HandleMessage(opts *bind.TransactOpts, _from string, _msg []byte) (*types.Transaction, error) {
	return _Xcall.contract.Transact(opts, "handleMessage", _from, _msg)
}

// HandleMessage is a paid mutator transaction binding the contract method 0xbbc22efd.
//
// Solidity: function handleMessage(string _from, bytes _msg) returns()
func (_Xcall *XcallSession) HandleMessage(_from string, _msg []byte) (*types.Transaction, error) {
	return _Xcall.Contract.HandleMessage(&_Xcall.TransactOpts, _from, _msg)
}

// HandleMessage is a paid mutator transaction binding the contract method 0xbbc22efd.
//
// Solidity: function handleMessage(string _from, bytes _msg) returns()
func (_Xcall *XcallTransactorSession) HandleMessage(_from string, _msg []byte) (*types.Transaction, error) {
	return _Xcall.Contract.HandleMessage(&_Xcall.TransactOpts, _from, _msg)
}

// Initialize is a paid mutator transaction binding the contract method 0xf62d1888.
//
// Solidity: function initialize(string _nid) returns()
func (_Xcall *XcallTransactor) Initialize(opts *bind.TransactOpts, _nid string) (*types.Transaction, error) {
	return _Xcall.contract.Transact(opts, "initialize", _nid)
}

// Initialize is a paid mutator transaction binding the contract method 0xf62d1888.
//
// Solidity: function initialize(string _nid) returns()
func (_Xcall *XcallSession) Initialize(_nid string) (*types.Transaction, error) {
	return _Xcall.Contract.Initialize(&_Xcall.TransactOpts, _nid)
}

// Initialize is a paid mutator transaction binding the contract method 0xf62d1888.
//
// Solidity: function initialize(string _nid) returns()
func (_Xcall *XcallTransactorSession) Initialize(_nid string) (*types.Transaction, error) {
	return _Xcall.Contract.Initialize(&_Xcall.TransactOpts, _nid)
}

// SendCallMessage is a paid mutator transaction binding the contract method 0x8ef378b8.
//
// Solidity: function sendCallMessage(string _to, bytes _data, bytes _rollback) payable returns(uint256)
func (_Xcall *XcallTransactor) SendCallMessage(opts *bind.TransactOpts, _to string, _data []byte, _rollback []byte) (*types.Transaction, error) {
	return _Xcall.contract.Transact(opts, "sendCallMessage", _to, _data, _rollback)
}

// SendCallMessage is a paid mutator transaction binding the contract method 0x8ef378b8.
//
// Solidity: function sendCallMessage(string _to, bytes _data, bytes _rollback) payable returns(uint256)
func (_Xcall *XcallSession) SendCallMessage(_to string, _data []byte, _rollback []byte) (*types.Transaction, error) {
	return _Xcall.Contract.SendCallMessage(&_Xcall.TransactOpts, _to, _data, _rollback)
}

// SendCallMessage is a paid mutator transaction binding the contract method 0x8ef378b8.
//
// Solidity: function sendCallMessage(string _to, bytes _data, bytes _rollback) payable returns(uint256)
func (_Xcall *XcallTransactorSession) SendCallMessage(_to string, _data []byte, _rollback []byte) (*types.Transaction, error) {
	return _Xcall.Contract.SendCallMessage(&_Xcall.TransactOpts, _to, _data, _rollback)
}

// SendCallMessage0 is a paid mutator transaction binding the contract method 0xedc6afff.
//
// Solidity: function sendCallMessage(string _to, bytes _data, bytes _rollback, string[] sources, string[] destinations) payable returns(uint256)
func (_Xcall *XcallTransactor) SendCallMessage0(opts *bind.TransactOpts, _to string, _data []byte, _rollback []byte, sources []string, destinations []string) (*types.Transaction, error) {
	return _Xcall.contract.Transact(opts, "sendCallMessage0", _to, _data, _rollback, sources, destinations)
}

// SendCallMessage0 is a paid mutator transaction binding the contract method 0xedc6afff.
//
// Solidity: function sendCallMessage(string _to, bytes _data, bytes _rollback, string[] sources, string[] destinations) payable returns(uint256)
func (_Xcall *XcallSession) SendCallMessage0(_to string, _data []byte, _rollback []byte, sources []string, destinations []string) (*types.Transaction, error) {
	return _Xcall.Contract.SendCallMessage0(&_Xcall.TransactOpts, _to, _data, _rollback, sources, destinations)
}

// SendCallMessage0 is a paid mutator transaction binding the contract method 0xedc6afff.
//
// Solidity: function sendCallMessage(string _to, bytes _data, bytes _rollback, string[] sources, string[] destinations) payable returns(uint256)
func (_Xcall *XcallTransactorSession) SendCallMessage0(_to string, _data []byte, _rollback []byte, sources []string, destinations []string) (*types.Transaction, error) {
	return _Xcall.Contract.SendCallMessage0(&_Xcall.TransactOpts, _to, _data, _rollback, sources, destinations)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x704b6c02.
//
// Solidity: function setAdmin(address _address) returns()
func (_Xcall *XcallTransactor) SetAdmin(opts *bind.TransactOpts, _address common.Address) (*types.Transaction, error) {
	return _Xcall.contract.Transact(opts, "setAdmin", _address)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x704b6c02.
//
// Solidity: function setAdmin(address _address) returns()
func (_Xcall *XcallSession) SetAdmin(_address common.Address) (*types.Transaction, error) {
	return _Xcall.Contract.SetAdmin(&_Xcall.TransactOpts, _address)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x704b6c02.
//
// Solidity: function setAdmin(address _address) returns()
func (_Xcall *XcallTransactorSession) SetAdmin(_address common.Address) (*types.Transaction, error) {
	return _Xcall.Contract.SetAdmin(&_Xcall.TransactOpts, _address)
}

// SetDefaultConnection is a paid mutator transaction binding the contract method 0x64f03757.
//
// Solidity: function setDefaultConnection(string _nid, address connection) returns()
func (_Xcall *XcallTransactor) SetDefaultConnection(opts *bind.TransactOpts, _nid string, connection common.Address) (*types.Transaction, error) {
	return _Xcall.contract.Transact(opts, "setDefaultConnection", _nid, connection)
}

// SetDefaultConnection is a paid mutator transaction binding the contract method 0x64f03757.
//
// Solidity: function setDefaultConnection(string _nid, address connection) returns()
func (_Xcall *XcallSession) SetDefaultConnection(_nid string, connection common.Address) (*types.Transaction, error) {
	return _Xcall.Contract.SetDefaultConnection(&_Xcall.TransactOpts, _nid, connection)
}

// SetDefaultConnection is a paid mutator transaction binding the contract method 0x64f03757.
//
// Solidity: function setDefaultConnection(string _nid, address connection) returns()
func (_Xcall *XcallTransactorSession) SetDefaultConnection(_nid string, connection common.Address) (*types.Transaction, error) {
	return _Xcall.Contract.SetDefaultConnection(&_Xcall.TransactOpts, _nid, connection)
}

// SetProtocolFee is a paid mutator transaction binding the contract method 0x787dce3d.
//
// Solidity: function setProtocolFee(uint256 _value) returns()
func (_Xcall *XcallTransactor) SetProtocolFee(opts *bind.TransactOpts, _value *big.Int) (*types.Transaction, error) {
	return _Xcall.contract.Transact(opts, "setProtocolFee", _value)
}

// SetProtocolFee is a paid mutator transaction binding the contract method 0x787dce3d.
//
// Solidity: function setProtocolFee(uint256 _value) returns()
func (_Xcall *XcallSession) SetProtocolFee(_value *big.Int) (*types.Transaction, error) {
	return _Xcall.Contract.SetProtocolFee(&_Xcall.TransactOpts, _value)
}

// SetProtocolFee is a paid mutator transaction binding the contract method 0x787dce3d.
//
// Solidity: function setProtocolFee(uint256 _value) returns()
func (_Xcall *XcallTransactorSession) SetProtocolFee(_value *big.Int) (*types.Transaction, error) {
	return _Xcall.Contract.SetProtocolFee(&_Xcall.TransactOpts, _value)
}

// SetProtocolFeeHandler is a paid mutator transaction binding the contract method 0x502bf8e0.
//
// Solidity: function setProtocolFeeHandler(address _addr) returns()
func (_Xcall *XcallTransactor) SetProtocolFeeHandler(opts *bind.TransactOpts, _addr common.Address) (*types.Transaction, error) {
	return _Xcall.contract.Transact(opts, "setProtocolFeeHandler", _addr)
}

// SetProtocolFeeHandler is a paid mutator transaction binding the contract method 0x502bf8e0.
//
// Solidity: function setProtocolFeeHandler(address _addr) returns()
func (_Xcall *XcallSession) SetProtocolFeeHandler(_addr common.Address) (*types.Transaction, error) {
	return _Xcall.Contract.SetProtocolFeeHandler(&_Xcall.TransactOpts, _addr)
}

// SetProtocolFeeHandler is a paid mutator transaction binding the contract method 0x502bf8e0.
//
// Solidity: function setProtocolFeeHandler(address _addr) returns()
func (_Xcall *XcallTransactorSession) SetProtocolFeeHandler(_addr common.Address) (*types.Transaction, error) {
	return _Xcall.Contract.SetProtocolFeeHandler(&_Xcall.TransactOpts, _addr)
}

// TryHandleCallMessage is a paid mutator transaction binding the contract method 0x04df28a9.
//
// Solidity: function tryHandleCallMessage(address toAddr, string to, string from, bytes data, string[] protocols) returns()
func (_Xcall *XcallTransactor) TryHandleCallMessage(opts *bind.TransactOpts, toAddr common.Address, to string, from string, data []byte, protocols []string) (*types.Transaction, error) {
	return _Xcall.contract.Transact(opts, "tryHandleCallMessage", toAddr, to, from, data, protocols)
}

// TryHandleCallMessage is a paid mutator transaction binding the contract method 0x04df28a9.
//
// Solidity: function tryHandleCallMessage(address toAddr, string to, string from, bytes data, string[] protocols) returns()
func (_Xcall *XcallSession) TryHandleCallMessage(toAddr common.Address, to string, from string, data []byte, protocols []string) (*types.Transaction, error) {
	return _Xcall.Contract.TryHandleCallMessage(&_Xcall.TransactOpts, toAddr, to, from, data, protocols)
}

// TryHandleCallMessage is a paid mutator transaction binding the contract method 0x04df28a9.
//
// Solidity: function tryHandleCallMessage(address toAddr, string to, string from, bytes data, string[] protocols) returns()
func (_Xcall *XcallTransactorSession) TryHandleCallMessage(toAddr common.Address, to string, from string, data []byte, protocols []string) (*types.Transaction, error) {
	return _Xcall.Contract.TryHandleCallMessage(&_Xcall.TransactOpts, toAddr, to, from, data, protocols)
}

// XcallCallExecutedIterator is returned from FilterCallExecuted and is used to iterate over the raw logs and unpacked data for CallExecuted events raised by the Xcall contract.
type XcallCallExecutedIterator struct {
	Event *XcallCallExecuted // Event containing the contract specifics and raw log

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
func (it *XcallCallExecutedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(XcallCallExecuted)
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
		it.Event = new(XcallCallExecuted)
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
func (it *XcallCallExecutedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *XcallCallExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// XcallCallExecuted represents a CallExecuted event raised by the Xcall contract.
type XcallCallExecuted struct {
	ReqId *big.Int
	Code  *big.Int
	Msg   string
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterCallExecuted is a free log retrieval operation binding the contract event 0xc7391e04887f8b3c16fa20877e028e8163139a478c8447e7d449eba1905caa51.
//
// Solidity: event CallExecuted(uint256 indexed _reqId, int256 _code, string _msg)
func (_Xcall *XcallFilterer) FilterCallExecuted(opts *bind.FilterOpts, _reqId []*big.Int) (*XcallCallExecutedIterator, error) {

	var _reqIdRule []interface{}
	for _, _reqIdItem := range _reqId {
		_reqIdRule = append(_reqIdRule, _reqIdItem)
	}

	logs, sub, err := _Xcall.contract.FilterLogs(opts, "CallExecuted", _reqIdRule)
	if err != nil {
		return nil, err
	}
	return &XcallCallExecutedIterator{contract: _Xcall.contract, event: "CallExecuted", logs: logs, sub: sub}, nil
}

// WatchCallExecuted is a free log subscription operation binding the contract event 0xc7391e04887f8b3c16fa20877e028e8163139a478c8447e7d449eba1905caa51.
//
// Solidity: event CallExecuted(uint256 indexed _reqId, int256 _code, string _msg)
func (_Xcall *XcallFilterer) WatchCallExecuted(opts *bind.WatchOpts, sink chan<- *XcallCallExecuted, _reqId []*big.Int) (event.Subscription, error) {

	var _reqIdRule []interface{}
	for _, _reqIdItem := range _reqId {
		_reqIdRule = append(_reqIdRule, _reqIdItem)
	}

	logs, sub, err := _Xcall.contract.WatchLogs(opts, "CallExecuted", _reqIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(XcallCallExecuted)
				if err := _Xcall.contract.UnpackLog(event, "CallExecuted", log); err != nil {
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
func (_Xcall *XcallFilterer) ParseCallExecuted(log types.Log) (*XcallCallExecuted, error) {
	event := new(XcallCallExecuted)
	if err := _Xcall.contract.UnpackLog(event, "CallExecuted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// XcallCallMessageIterator is returned from FilterCallMessage and is used to iterate over the raw logs and unpacked data for CallMessage events raised by the Xcall contract.
type XcallCallMessageIterator struct {
	Event *XcallCallMessage // Event containing the contract specifics and raw log

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
func (it *XcallCallMessageIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(XcallCallMessage)
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
		it.Event = new(XcallCallMessage)
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
func (it *XcallCallMessageIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *XcallCallMessageIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// XcallCallMessage represents a CallMessage event raised by the Xcall contract.
type XcallCallMessage struct {
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
func (_Xcall *XcallFilterer) FilterCallMessage(opts *bind.FilterOpts, _from []string, _to []string, _sn []*big.Int) (*XcallCallMessageIterator, error) {

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

	logs, sub, err := _Xcall.contract.FilterLogs(opts, "CallMessage", _fromRule, _toRule, _snRule)
	if err != nil {
		return nil, err
	}
	return &XcallCallMessageIterator{contract: _Xcall.contract, event: "CallMessage", logs: logs, sub: sub}, nil
}

// WatchCallMessage is a free log subscription operation binding the contract event 0x2cbc78425621c181f9f8a25fc06e44a0ac2b67cd6a31f8ed7918934187f8cc59.
//
// Solidity: event CallMessage(string indexed _from, string indexed _to, uint256 indexed _sn, uint256 _reqId, bytes _data)
func (_Xcall *XcallFilterer) WatchCallMessage(opts *bind.WatchOpts, sink chan<- *XcallCallMessage, _from []string, _to []string, _sn []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _Xcall.contract.WatchLogs(opts, "CallMessage", _fromRule, _toRule, _snRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(XcallCallMessage)
				if err := _Xcall.contract.UnpackLog(event, "CallMessage", log); err != nil {
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
func (_Xcall *XcallFilterer) ParseCallMessage(log types.Log) (*XcallCallMessage, error) {
	event := new(XcallCallMessage)
	if err := _Xcall.contract.UnpackLog(event, "CallMessage", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// XcallCallMessageSentIterator is returned from FilterCallMessageSent and is used to iterate over the raw logs and unpacked data for CallMessageSent events raised by the Xcall contract.
type XcallCallMessageSentIterator struct {
	Event *XcallCallMessageSent // Event containing the contract specifics and raw log

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
func (it *XcallCallMessageSentIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(XcallCallMessageSent)
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
		it.Event = new(XcallCallMessageSent)
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
func (it *XcallCallMessageSentIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *XcallCallMessageSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// XcallCallMessageSent represents a CallMessageSent event raised by the Xcall contract.
type XcallCallMessageSent struct {
	From common.Address
	To   common.Hash
	Sn   *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterCallMessageSent is a free log retrieval operation binding the contract event 0x69e53ea70fdf945f6d035b3979748bc999151691fb1dc69d66f8017f8840ae28.
//
// Solidity: event CallMessageSent(address indexed _from, string indexed _to, uint256 indexed _sn)
func (_Xcall *XcallFilterer) FilterCallMessageSent(opts *bind.FilterOpts, _from []common.Address, _to []string, _sn []*big.Int) (*XcallCallMessageSentIterator, error) {

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

	logs, sub, err := _Xcall.contract.FilterLogs(opts, "CallMessageSent", _fromRule, _toRule, _snRule)
	if err != nil {
		return nil, err
	}
	return &XcallCallMessageSentIterator{contract: _Xcall.contract, event: "CallMessageSent", logs: logs, sub: sub}, nil
}

// WatchCallMessageSent is a free log subscription operation binding the contract event 0x69e53ea70fdf945f6d035b3979748bc999151691fb1dc69d66f8017f8840ae28.
//
// Solidity: event CallMessageSent(address indexed _from, string indexed _to, uint256 indexed _sn)
func (_Xcall *XcallFilterer) WatchCallMessageSent(opts *bind.WatchOpts, sink chan<- *XcallCallMessageSent, _from []common.Address, _to []string, _sn []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _Xcall.contract.WatchLogs(opts, "CallMessageSent", _fromRule, _toRule, _snRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(XcallCallMessageSent)
				if err := _Xcall.contract.UnpackLog(event, "CallMessageSent", log); err != nil {
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
func (_Xcall *XcallFilterer) ParseCallMessageSent(log types.Log) (*XcallCallMessageSent, error) {
	event := new(XcallCallMessageSent)
	if err := _Xcall.contract.UnpackLog(event, "CallMessageSent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// XcallInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Xcall contract.
type XcallInitializedIterator struct {
	Event *XcallInitialized // Event containing the contract specifics and raw log

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
func (it *XcallInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(XcallInitialized)
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
		it.Event = new(XcallInitialized)
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
func (it *XcallInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *XcallInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// XcallInitialized represents a Initialized event raised by the Xcall contract.
type XcallInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Xcall *XcallFilterer) FilterInitialized(opts *bind.FilterOpts) (*XcallInitializedIterator, error) {

	logs, sub, err := _Xcall.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &XcallInitializedIterator{contract: _Xcall.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Xcall *XcallFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *XcallInitialized) (event.Subscription, error) {

	logs, sub, err := _Xcall.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(XcallInitialized)
				if err := _Xcall.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Xcall *XcallFilterer) ParseInitialized(log types.Log) (*XcallInitialized, error) {
	event := new(XcallInitialized)
	if err := _Xcall.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// XcallResponseMessageIterator is returned from FilterResponseMessage and is used to iterate over the raw logs and unpacked data for ResponseMessage events raised by the Xcall contract.
type XcallResponseMessageIterator struct {
	Event *XcallResponseMessage // Event containing the contract specifics and raw log

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
func (it *XcallResponseMessageIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(XcallResponseMessage)
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
		it.Event = new(XcallResponseMessage)
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
func (it *XcallResponseMessageIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *XcallResponseMessageIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// XcallResponseMessage represents a ResponseMessage event raised by the Xcall contract.
type XcallResponseMessage struct {
	Sn   *big.Int
	Code *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterResponseMessage is a free log retrieval operation binding the contract event 0xbeacafd006c5e60667f6f04aec3a498f81c8e94142b4e95b5a5a763de43ca0ab.
//
// Solidity: event ResponseMessage(uint256 indexed _sn, int256 _code)
func (_Xcall *XcallFilterer) FilterResponseMessage(opts *bind.FilterOpts, _sn []*big.Int) (*XcallResponseMessageIterator, error) {

	var _snRule []interface{}
	for _, _snItem := range _sn {
		_snRule = append(_snRule, _snItem)
	}

	logs, sub, err := _Xcall.contract.FilterLogs(opts, "ResponseMessage", _snRule)
	if err != nil {
		return nil, err
	}
	return &XcallResponseMessageIterator{contract: _Xcall.contract, event: "ResponseMessage", logs: logs, sub: sub}, nil
}

// WatchResponseMessage is a free log subscription operation binding the contract event 0xbeacafd006c5e60667f6f04aec3a498f81c8e94142b4e95b5a5a763de43ca0ab.
//
// Solidity: event ResponseMessage(uint256 indexed _sn, int256 _code)
func (_Xcall *XcallFilterer) WatchResponseMessage(opts *bind.WatchOpts, sink chan<- *XcallResponseMessage, _sn []*big.Int) (event.Subscription, error) {

	var _snRule []interface{}
	for _, _snItem := range _sn {
		_snRule = append(_snRule, _snItem)
	}

	logs, sub, err := _Xcall.contract.WatchLogs(opts, "ResponseMessage", _snRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(XcallResponseMessage)
				if err := _Xcall.contract.UnpackLog(event, "ResponseMessage", log); err != nil {
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
func (_Xcall *XcallFilterer) ParseResponseMessage(log types.Log) (*XcallResponseMessage, error) {
	event := new(XcallResponseMessage)
	if err := _Xcall.contract.UnpackLog(event, "ResponseMessage", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// XcallRollbackExecutedIterator is returned from FilterRollbackExecuted and is used to iterate over the raw logs and unpacked data for RollbackExecuted events raised by the Xcall contract.
type XcallRollbackExecutedIterator struct {
	Event *XcallRollbackExecuted // Event containing the contract specifics and raw log

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
func (it *XcallRollbackExecutedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(XcallRollbackExecuted)
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
		it.Event = new(XcallRollbackExecuted)
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
func (it *XcallRollbackExecutedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *XcallRollbackExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// XcallRollbackExecuted represents a RollbackExecuted event raised by the Xcall contract.
type XcallRollbackExecuted struct {
	Sn  *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterRollbackExecuted is a free log retrieval operation binding the contract event 0x08f0ac7aef6da8bbe43bee8b1444a1883f1359566618bc379ce5abba44883837.
//
// Solidity: event RollbackExecuted(uint256 indexed _sn)
func (_Xcall *XcallFilterer) FilterRollbackExecuted(opts *bind.FilterOpts, _sn []*big.Int) (*XcallRollbackExecutedIterator, error) {

	var _snRule []interface{}
	for _, _snItem := range _sn {
		_snRule = append(_snRule, _snItem)
	}

	logs, sub, err := _Xcall.contract.FilterLogs(opts, "RollbackExecuted", _snRule)
	if err != nil {
		return nil, err
	}
	return &XcallRollbackExecutedIterator{contract: _Xcall.contract, event: "RollbackExecuted", logs: logs, sub: sub}, nil
}

// WatchRollbackExecuted is a free log subscription operation binding the contract event 0x08f0ac7aef6da8bbe43bee8b1444a1883f1359566618bc379ce5abba44883837.
//
// Solidity: event RollbackExecuted(uint256 indexed _sn)
func (_Xcall *XcallFilterer) WatchRollbackExecuted(opts *bind.WatchOpts, sink chan<- *XcallRollbackExecuted, _sn []*big.Int) (event.Subscription, error) {

	var _snRule []interface{}
	for _, _snItem := range _sn {
		_snRule = append(_snRule, _snItem)
	}

	logs, sub, err := _Xcall.contract.WatchLogs(opts, "RollbackExecuted", _snRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(XcallRollbackExecuted)
				if err := _Xcall.contract.UnpackLog(event, "RollbackExecuted", log); err != nil {
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
func (_Xcall *XcallFilterer) ParseRollbackExecuted(log types.Log) (*XcallRollbackExecuted, error) {
	event := new(XcallRollbackExecuted)
	if err := _Xcall.contract.UnpackLog(event, "RollbackExecuted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// XcallRollbackMessageIterator is returned from FilterRollbackMessage and is used to iterate over the raw logs and unpacked data for RollbackMessage events raised by the Xcall contract.
type XcallRollbackMessageIterator struct {
	Event *XcallRollbackMessage // Event containing the contract specifics and raw log

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
func (it *XcallRollbackMessageIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(XcallRollbackMessage)
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
		it.Event = new(XcallRollbackMessage)
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
func (it *XcallRollbackMessageIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *XcallRollbackMessageIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// XcallRollbackMessage represents a RollbackMessage event raised by the Xcall contract.
type XcallRollbackMessage struct {
	Sn  *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterRollbackMessage is a free log retrieval operation binding the contract event 0x38934ab923f985814047679ba041577b8203ddd15fe9910d3fc6a7aa6001e9c7.
//
// Solidity: event RollbackMessage(uint256 indexed _sn)
func (_Xcall *XcallFilterer) FilterRollbackMessage(opts *bind.FilterOpts, _sn []*big.Int) (*XcallRollbackMessageIterator, error) {

	var _snRule []interface{}
	for _, _snItem := range _sn {
		_snRule = append(_snRule, _snItem)
	}

	logs, sub, err := _Xcall.contract.FilterLogs(opts, "RollbackMessage", _snRule)
	if err != nil {
		return nil, err
	}
	return &XcallRollbackMessageIterator{contract: _Xcall.contract, event: "RollbackMessage", logs: logs, sub: sub}, nil
}

// WatchRollbackMessage is a free log subscription operation binding the contract event 0x38934ab923f985814047679ba041577b8203ddd15fe9910d3fc6a7aa6001e9c7.
//
// Solidity: event RollbackMessage(uint256 indexed _sn)
func (_Xcall *XcallFilterer) WatchRollbackMessage(opts *bind.WatchOpts, sink chan<- *XcallRollbackMessage, _sn []*big.Int) (event.Subscription, error) {

	var _snRule []interface{}
	for _, _snItem := range _sn {
		_snRule = append(_snRule, _snItem)
	}

	logs, sub, err := _Xcall.contract.WatchLogs(opts, "RollbackMessage", _snRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(XcallRollbackMessage)
				if err := _Xcall.contract.UnpackLog(event, "RollbackMessage", log); err != nil {
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
func (_Xcall *XcallFilterer) ParseRollbackMessage(log types.Log) (*XcallRollbackMessage, error) {
	event := new(XcallRollbackMessage)
	if err := _Xcall.contract.UnpackLog(event, "RollbackMessage", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
