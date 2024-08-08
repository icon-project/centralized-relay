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

// CentralizedConnectionMetaData contains all meta data concerning the CentralizedConnection contract.
var CentralizedConnectionMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"admin\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"claimFees\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"connSn\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getFee\",\"inputs\":[{\"name\":\"to\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"response\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[{\"name\":\"fee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getReceipt\",\"inputs\":[{\"name\":\"srcNetwork\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_connSn\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_relayer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_xCall\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"recvMessage\",\"inputs\":[{\"name\":\"srcNetwork\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_connSn\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_msg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revertMessage\",\"inputs\":[{\"name\":\"sn\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"sendMessage\",\"inputs\":[{\"name\":\"to\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"svc\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"sn\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"_msg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"setAdmin\",\"inputs\":[{\"name\":\"_address\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setFee\",\"inputs\":[{\"name\":\"networkId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"messageFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"responseFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Message\",\"inputs\":[{\"name\":\"targetNetwork\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"sn\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"_msg\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]}]",
}

// CentralizedConnectionABI is the input ABI used to generate the binding from.
// Deprecated: Use CentralizedConnectionMetaData.ABI instead.
var CentralizedConnectionABI = CentralizedConnectionMetaData.ABI

// CentralizedConnection is an auto generated Go binding around an Ethereum contract.
type CentralizedConnection struct {
	CentralizedConnectionCaller     // Read-only binding to the contract
	CentralizedConnectionTransactor // Write-only binding to the contract
	CentralizedConnectionFilterer   // Log filterer for contract events
}

// CentralizedConnectionCaller is an auto generated read-only Go binding around an Ethereum contract.
type CentralizedConnectionCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CentralizedConnectionTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CentralizedConnectionTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CentralizedConnectionFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CentralizedConnectionFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CentralizedConnectionSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CentralizedConnectionSession struct {
	Contract     *CentralizedConnection // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// CentralizedConnectionCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CentralizedConnectionCallerSession struct {
	Contract *CentralizedConnectionCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// CentralizedConnectionTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CentralizedConnectionTransactorSession struct {
	Contract     *CentralizedConnectionTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// CentralizedConnectionRaw is an auto generated low-level Go binding around an Ethereum contract.
type CentralizedConnectionRaw struct {
	Contract *CentralizedConnection // Generic contract binding to access the raw methods on
}

// CentralizedConnectionCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CentralizedConnectionCallerRaw struct {
	Contract *CentralizedConnectionCaller // Generic read-only contract binding to access the raw methods on
}

// CentralizedConnectionTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CentralizedConnectionTransactorRaw struct {
	Contract *CentralizedConnectionTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCentralizedConnection creates a new instance of CentralizedConnection, bound to a specific deployed contract.
func NewCentralizedConnection(address common.Address, backend bind.ContractBackend) (*CentralizedConnection, error) {
	contract, err := bindCentralizedConnection(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CentralizedConnection{CentralizedConnectionCaller: CentralizedConnectionCaller{contract: contract}, CentralizedConnectionTransactor: CentralizedConnectionTransactor{contract: contract}, CentralizedConnectionFilterer: CentralizedConnectionFilterer{contract: contract}}, nil
}

// NewCentralizedConnectionCaller creates a new read-only instance of CentralizedConnection, bound to a specific deployed contract.
func NewCentralizedConnectionCaller(address common.Address, caller bind.ContractCaller) (*CentralizedConnectionCaller, error) {
	contract, err := bindCentralizedConnection(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CentralizedConnectionCaller{contract: contract}, nil
}

// NewCentralizedConnectionTransactor creates a new write-only instance of CentralizedConnection, bound to a specific deployed contract.
func NewCentralizedConnectionTransactor(address common.Address, transactor bind.ContractTransactor) (*CentralizedConnectionTransactor, error) {
	contract, err := bindCentralizedConnection(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CentralizedConnectionTransactor{contract: contract}, nil
}

// NewCentralizedConnectionFilterer creates a new log filterer instance of CentralizedConnection, bound to a specific deployed contract.
func NewCentralizedConnectionFilterer(address common.Address, filterer bind.ContractFilterer) (*CentralizedConnectionFilterer, error) {
	contract, err := bindCentralizedConnection(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CentralizedConnectionFilterer{contract: contract}, nil
}

// bindCentralizedConnection binds a generic wrapper to an already deployed contract.
func bindCentralizedConnection(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CentralizedConnectionMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CentralizedConnection *CentralizedConnectionRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CentralizedConnection.Contract.CentralizedConnectionCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CentralizedConnection *CentralizedConnectionRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CentralizedConnection.Contract.CentralizedConnectionTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CentralizedConnection *CentralizedConnectionRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CentralizedConnection.Contract.CentralizedConnectionTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CentralizedConnection *CentralizedConnectionCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CentralizedConnection.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CentralizedConnection *CentralizedConnectionTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CentralizedConnection.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CentralizedConnection *CentralizedConnectionTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CentralizedConnection.Contract.contract.Transact(opts, method, params...)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_CentralizedConnection *CentralizedConnectionCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CentralizedConnection.contract.Call(opts, &out, "admin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_CentralizedConnection *CentralizedConnectionSession) Admin() (common.Address, error) {
	return _CentralizedConnection.Contract.Admin(&_CentralizedConnection.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_CentralizedConnection *CentralizedConnectionCallerSession) Admin() (common.Address, error) {
	return _CentralizedConnection.Contract.Admin(&_CentralizedConnection.CallOpts)
}

// ConnSn is a free data retrieval call binding the contract method 0x99f1fca7.
//
// Solidity: function connSn() view returns(uint256)
func (_CentralizedConnection *CentralizedConnectionCaller) ConnSn(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CentralizedConnection.contract.Call(opts, &out, "connSn")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ConnSn is a free data retrieval call binding the contract method 0x99f1fca7.
//
// Solidity: function connSn() view returns(uint256)
func (_CentralizedConnection *CentralizedConnectionSession) ConnSn() (*big.Int, error) {
	return _CentralizedConnection.Contract.ConnSn(&_CentralizedConnection.CallOpts)
}

// ConnSn is a free data retrieval call binding the contract method 0x99f1fca7.
//
// Solidity: function connSn() view returns(uint256)
func (_CentralizedConnection *CentralizedConnectionCallerSession) ConnSn() (*big.Int, error) {
	return _CentralizedConnection.Contract.ConnSn(&_CentralizedConnection.CallOpts)
}

// GetFee is a free data retrieval call binding the contract method 0x7d4c4f4a.
//
// Solidity: function getFee(string to, bool response) view returns(uint256 fee)
func (_CentralizedConnection *CentralizedConnectionCaller) GetFee(opts *bind.CallOpts, to string, response bool) (*big.Int, error) {
	var out []interface{}
	err := _CentralizedConnection.contract.Call(opts, &out, "getFee", to, response)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetFee is a free data retrieval call binding the contract method 0x7d4c4f4a.
//
// Solidity: function getFee(string to, bool response) view returns(uint256 fee)
func (_CentralizedConnection *CentralizedConnectionSession) GetFee(to string, response bool) (*big.Int, error) {
	return _CentralizedConnection.Contract.GetFee(&_CentralizedConnection.CallOpts, to, response)
}

// GetFee is a free data retrieval call binding the contract method 0x7d4c4f4a.
//
// Solidity: function getFee(string to, bool response) view returns(uint256 fee)
func (_CentralizedConnection *CentralizedConnectionCallerSession) GetFee(to string, response bool) (*big.Int, error) {
	return _CentralizedConnection.Contract.GetFee(&_CentralizedConnection.CallOpts, to, response)
}

// GetReceipt is a free data retrieval call binding the contract method 0x9664da0e.
//
// Solidity: function getReceipt(string srcNetwork, uint256 _connSn) view returns(bool)
func (_CentralizedConnection *CentralizedConnectionCaller) GetReceipt(opts *bind.CallOpts, srcNetwork string, _connSn *big.Int) (bool, error) {
	var out []interface{}
	err := _CentralizedConnection.contract.Call(opts, &out, "getReceipt", srcNetwork, _connSn)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetReceipt is a free data retrieval call binding the contract method 0x9664da0e.
//
// Solidity: function getReceipt(string srcNetwork, uint256 _connSn) view returns(bool)
func (_CentralizedConnection *CentralizedConnectionSession) GetReceipt(srcNetwork string, _connSn *big.Int) (bool, error) {
	return _CentralizedConnection.Contract.GetReceipt(&_CentralizedConnection.CallOpts, srcNetwork, _connSn)
}

// GetReceipt is a free data retrieval call binding the contract method 0x9664da0e.
//
// Solidity: function getReceipt(string srcNetwork, uint256 _connSn) view returns(bool)
func (_CentralizedConnection *CentralizedConnectionCallerSession) GetReceipt(srcNetwork string, _connSn *big.Int) (bool, error) {
	return _CentralizedConnection.Contract.GetReceipt(&_CentralizedConnection.CallOpts, srcNetwork, _connSn)
}

// ClaimFees is a paid mutator transaction binding the contract method 0xd294f093.
//
// Solidity: function claimFees() returns()
func (_CentralizedConnection *CentralizedConnectionTransactor) ClaimFees(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CentralizedConnection.contract.Transact(opts, "claimFees")
}

// ClaimFees is a paid mutator transaction binding the contract method 0xd294f093.
//
// Solidity: function claimFees() returns()
func (_CentralizedConnection *CentralizedConnectionSession) ClaimFees() (*types.Transaction, error) {
	return _CentralizedConnection.Contract.ClaimFees(&_CentralizedConnection.TransactOpts)
}

// ClaimFees is a paid mutator transaction binding the contract method 0xd294f093.
//
// Solidity: function claimFees() returns()
func (_CentralizedConnection *CentralizedConnectionTransactorSession) ClaimFees() (*types.Transaction, error) {
	return _CentralizedConnection.Contract.ClaimFees(&_CentralizedConnection.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _relayer, address _xCall) returns()
func (_CentralizedConnection *CentralizedConnectionTransactor) Initialize(opts *bind.TransactOpts, _relayer common.Address, _xCall common.Address) (*types.Transaction, error) {
	return _CentralizedConnection.contract.Transact(opts, "initialize", _relayer, _xCall)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _relayer, address _xCall) returns()
func (_CentralizedConnection *CentralizedConnectionSession) Initialize(_relayer common.Address, _xCall common.Address) (*types.Transaction, error) {
	return _CentralizedConnection.Contract.Initialize(&_CentralizedConnection.TransactOpts, _relayer, _xCall)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _relayer, address _xCall) returns()
func (_CentralizedConnection *CentralizedConnectionTransactorSession) Initialize(_relayer common.Address, _xCall common.Address) (*types.Transaction, error) {
	return _CentralizedConnection.Contract.Initialize(&_CentralizedConnection.TransactOpts, _relayer, _xCall)
}

// RecvMessage is a paid mutator transaction binding the contract method 0xb58b4cec.
//
// Solidity: function recvMessage(string srcNetwork, uint256 _connSn, bytes _msg) returns()
func (_CentralizedConnection *CentralizedConnectionTransactor) RecvMessage(opts *bind.TransactOpts, srcNetwork string, _connSn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _CentralizedConnection.contract.Transact(opts, "recvMessage", srcNetwork, _connSn, _msg)
}

// RecvMessage is a paid mutator transaction binding the contract method 0xb58b4cec.
//
// Solidity: function recvMessage(string srcNetwork, uint256 _connSn, bytes _msg) returns()
func (_CentralizedConnection *CentralizedConnectionSession) RecvMessage(srcNetwork string, _connSn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _CentralizedConnection.Contract.RecvMessage(&_CentralizedConnection.TransactOpts, srcNetwork, _connSn, _msg)
}

// RecvMessage is a paid mutator transaction binding the contract method 0xb58b4cec.
//
// Solidity: function recvMessage(string srcNetwork, uint256 _connSn, bytes _msg) returns()
func (_CentralizedConnection *CentralizedConnectionTransactorSession) RecvMessage(srcNetwork string, _connSn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _CentralizedConnection.Contract.RecvMessage(&_CentralizedConnection.TransactOpts, srcNetwork, _connSn, _msg)
}

// RevertMessage is a paid mutator transaction binding the contract method 0x2d3fb823.
//
// Solidity: function revertMessage(uint256 sn) returns()
func (_CentralizedConnection *CentralizedConnectionTransactor) RevertMessage(opts *bind.TransactOpts, sn *big.Int) (*types.Transaction, error) {
	return _CentralizedConnection.contract.Transact(opts, "revertMessage", sn)
}

// RevertMessage is a paid mutator transaction binding the contract method 0x2d3fb823.
//
// Solidity: function revertMessage(uint256 sn) returns()
func (_CentralizedConnection *CentralizedConnectionSession) RevertMessage(sn *big.Int) (*types.Transaction, error) {
	return _CentralizedConnection.Contract.RevertMessage(&_CentralizedConnection.TransactOpts, sn)
}

// RevertMessage is a paid mutator transaction binding the contract method 0x2d3fb823.
//
// Solidity: function revertMessage(uint256 sn) returns()
func (_CentralizedConnection *CentralizedConnectionTransactorSession) RevertMessage(sn *big.Int) (*types.Transaction, error) {
	return _CentralizedConnection.Contract.RevertMessage(&_CentralizedConnection.TransactOpts, sn)
}

// SendMessage is a paid mutator transaction binding the contract method 0x522a901e.
//
// Solidity: function sendMessage(string to, string svc, int256 sn, bytes _msg) payable returns()
func (_CentralizedConnection *CentralizedConnectionTransactor) SendMessage(opts *bind.TransactOpts, to string, svc string, sn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _CentralizedConnection.contract.Transact(opts, "sendMessage", to, svc, sn, _msg)
}

// SendMessage is a paid mutator transaction binding the contract method 0x522a901e.
//
// Solidity: function sendMessage(string to, string svc, int256 sn, bytes _msg) payable returns()
func (_CentralizedConnection *CentralizedConnectionSession) SendMessage(to string, svc string, sn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _CentralizedConnection.Contract.SendMessage(&_CentralizedConnection.TransactOpts, to, svc, sn, _msg)
}

// SendMessage is a paid mutator transaction binding the contract method 0x522a901e.
//
// Solidity: function sendMessage(string to, string svc, int256 sn, bytes _msg) payable returns()
func (_CentralizedConnection *CentralizedConnectionTransactorSession) SendMessage(to string, svc string, sn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _CentralizedConnection.Contract.SendMessage(&_CentralizedConnection.TransactOpts, to, svc, sn, _msg)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x704b6c02.
//
// Solidity: function setAdmin(address _address) returns()
func (_CentralizedConnection *CentralizedConnectionTransactor) SetAdmin(opts *bind.TransactOpts, _address common.Address) (*types.Transaction, error) {
	return _CentralizedConnection.contract.Transact(opts, "setAdmin", _address)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x704b6c02.
//
// Solidity: function setAdmin(address _address) returns()
func (_CentralizedConnection *CentralizedConnectionSession) SetAdmin(_address common.Address) (*types.Transaction, error) {
	return _CentralizedConnection.Contract.SetAdmin(&_CentralizedConnection.TransactOpts, _address)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x704b6c02.
//
// Solidity: function setAdmin(address _address) returns()
func (_CentralizedConnection *CentralizedConnectionTransactorSession) SetAdmin(_address common.Address) (*types.Transaction, error) {
	return _CentralizedConnection.Contract.SetAdmin(&_CentralizedConnection.TransactOpts, _address)
}

// SetFee is a paid mutator transaction binding the contract method 0x43f08a89.
//
// Solidity: function setFee(string networkId, uint256 messageFee, uint256 responseFee) returns()
func (_CentralizedConnection *CentralizedConnectionTransactor) SetFee(opts *bind.TransactOpts, networkId string, messageFee *big.Int, responseFee *big.Int) (*types.Transaction, error) {
	return _CentralizedConnection.contract.Transact(opts, "setFee", networkId, messageFee, responseFee)
}

// SetFee is a paid mutator transaction binding the contract method 0x43f08a89.
//
// Solidity: function setFee(string networkId, uint256 messageFee, uint256 responseFee) returns()
func (_CentralizedConnection *CentralizedConnectionSession) SetFee(networkId string, messageFee *big.Int, responseFee *big.Int) (*types.Transaction, error) {
	return _CentralizedConnection.Contract.SetFee(&_CentralizedConnection.TransactOpts, networkId, messageFee, responseFee)
}

// SetFee is a paid mutator transaction binding the contract method 0x43f08a89.
//
// Solidity: function setFee(string networkId, uint256 messageFee, uint256 responseFee) returns()
func (_CentralizedConnection *CentralizedConnectionTransactorSession) SetFee(networkId string, messageFee *big.Int, responseFee *big.Int) (*types.Transaction, error) {
	return _CentralizedConnection.Contract.SetFee(&_CentralizedConnection.TransactOpts, networkId, messageFee, responseFee)
}

// CentralizedConnectionInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the CentralizedConnection contract.
type CentralizedConnectionInitializedIterator struct {
	Event *CentralizedConnectionInitialized // Event containing the contract specifics and raw log

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
func (it *CentralizedConnectionInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CentralizedConnectionInitialized)
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
		it.Event = new(CentralizedConnectionInitialized)
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
func (it *CentralizedConnectionInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CentralizedConnectionInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CentralizedConnectionInitialized represents a Initialized event raised by the CentralizedConnection contract.
type CentralizedConnectionInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_CentralizedConnection *CentralizedConnectionFilterer) FilterInitialized(opts *bind.FilterOpts) (*CentralizedConnectionInitializedIterator, error) {

	logs, sub, err := _CentralizedConnection.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &CentralizedConnectionInitializedIterator{contract: _CentralizedConnection.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_CentralizedConnection *CentralizedConnectionFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *CentralizedConnectionInitialized) (event.Subscription, error) {

	logs, sub, err := _CentralizedConnection.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CentralizedConnectionInitialized)
				if err := _CentralizedConnection.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_CentralizedConnection *CentralizedConnectionFilterer) ParseInitialized(log types.Log) (*CentralizedConnectionInitialized, error) {
	event := new(CentralizedConnectionInitialized)
	if err := _CentralizedConnection.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CentralizedConnectionMessageIterator is returned from FilterMessage and is used to iterate over the raw logs and unpacked data for Message events raised by the CentralizedConnection contract.
type CentralizedConnectionMessageIterator struct {
	Event *CentralizedConnectionMessage // Event containing the contract specifics and raw log

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
func (it *CentralizedConnectionMessageIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CentralizedConnectionMessage)
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
		it.Event = new(CentralizedConnectionMessage)
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
func (it *CentralizedConnectionMessageIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CentralizedConnectionMessageIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CentralizedConnectionMessage represents a Message event raised by the CentralizedConnection contract.
type CentralizedConnectionMessage struct {
	TargetNetwork string
	Sn            *big.Int
	Msg           []byte
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterMessage is a free log retrieval operation binding the contract event 0x37be353f216cf7e33639101fd610c542e6a0c0109173fa1c1d8b04d34edb7c1b.
//
// Solidity: event Message(string targetNetwork, uint256 sn, bytes _msg)
func (_CentralizedConnection *CentralizedConnectionFilterer) FilterMessage(opts *bind.FilterOpts) (*CentralizedConnectionMessageIterator, error) {

	logs, sub, err := _CentralizedConnection.contract.FilterLogs(opts, "Message")
	if err != nil {
		return nil, err
	}
	return &CentralizedConnectionMessageIterator{contract: _CentralizedConnection.contract, event: "Message", logs: logs, sub: sub}, nil
}

// WatchMessage is a free log subscription operation binding the contract event 0x37be353f216cf7e33639101fd610c542e6a0c0109173fa1c1d8b04d34edb7c1b.
//
// Solidity: event Message(string targetNetwork, uint256 sn, bytes _msg)
func (_CentralizedConnection *CentralizedConnectionFilterer) WatchMessage(opts *bind.WatchOpts, sink chan<- *CentralizedConnectionMessage) (event.Subscription, error) {

	logs, sub, err := _CentralizedConnection.contract.WatchLogs(opts, "Message")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CentralizedConnectionMessage)
				if err := _CentralizedConnection.contract.UnpackLog(event, "Message", log); err != nil {
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
func (_CentralizedConnection *CentralizedConnectionFilterer) ParseMessage(log types.Log) (*CentralizedConnectionMessage, error) {
	event := new(CentralizedConnectionMessage)
	if err := _CentralizedConnection.contract.UnpackLog(event, "Message", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
