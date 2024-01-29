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

// ConnectionMetaData contains all meta data concerning the Connection contract.
var ConnectionMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"admin\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"claimFees\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"connSn\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getFee\",\"inputs\":[{\"name\":\"to\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"response\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[{\"name\":\"fee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getReceipt\",\"inputs\":[{\"name\":\"srcNetwork\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_connSn\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_relayer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_xCall\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"recvMessage\",\"inputs\":[{\"name\":\"srcNetwork\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_connSn\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_msg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revertMessage\",\"inputs\":[{\"name\":\"sn\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"sendMessage\",\"inputs\":[{\"name\":\"to\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"svc\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"sn\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"_msg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"setAdmin\",\"inputs\":[{\"name\":\"_address\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setFee\",\"inputs\":[{\"name\":\"networkId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"messageFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"responseFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Message\",\"inputs\":[{\"name\":\"targetNetwork\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"sn\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"_msg\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false}]",
}

// ConnectionABI is the input ABI used to generate the binding from.
// Deprecated: Use ConnectionMetaData.ABI instead.
var ConnectionABI = ConnectionMetaData.ABI

// Connection is an auto generated Go binding around an Ethereum contract.
type Connection struct {
	ConnectionCaller     // Read-only binding to the contract
	ConnectionTransactor // Write-only binding to the contract
	ConnectionFilterer   // Log filterer for contract events
}

// ConnectionCaller is an auto generated read-only Go binding around an Ethereum contract.
type ConnectionCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConnectionTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ConnectionTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConnectionFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ConnectionFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConnectionSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ConnectionSession struct {
	Contract     *Connection       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ConnectionCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ConnectionCallerSession struct {
	Contract *ConnectionCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// ConnectionTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ConnectionTransactorSession struct {
	Contract     *ConnectionTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// ConnectionRaw is an auto generated low-level Go binding around an Ethereum contract.
type ConnectionRaw struct {
	Contract *Connection // Generic contract binding to access the raw methods on
}

// ConnectionCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ConnectionCallerRaw struct {
	Contract *ConnectionCaller // Generic read-only contract binding to access the raw methods on
}

// ConnectionTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ConnectionTransactorRaw struct {
	Contract *ConnectionTransactor // Generic write-only contract binding to access the raw methods on
}

// NewConnection creates a new instance of Connection, bound to a specific deployed contract.
func NewConnection(address common.Address, backend bind.ContractBackend) (*Connection, error) {
	contract, err := bindConnection(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Connection{ConnectionCaller: ConnectionCaller{contract: contract}, ConnectionTransactor: ConnectionTransactor{contract: contract}, ConnectionFilterer: ConnectionFilterer{contract: contract}}, nil
}

// NewConnectionCaller creates a new read-only instance of Connection, bound to a specific deployed contract.
func NewConnectionCaller(address common.Address, caller bind.ContractCaller) (*ConnectionCaller, error) {
	contract, err := bindConnection(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ConnectionCaller{contract: contract}, nil
}

// NewConnectionTransactor creates a new write-only instance of Connection, bound to a specific deployed contract.
func NewConnectionTransactor(address common.Address, transactor bind.ContractTransactor) (*ConnectionTransactor, error) {
	contract, err := bindConnection(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ConnectionTransactor{contract: contract}, nil
}

// NewConnectionFilterer creates a new log filterer instance of Connection, bound to a specific deployed contract.
func NewConnectionFilterer(address common.Address, filterer bind.ContractFilterer) (*ConnectionFilterer, error) {
	contract, err := bindConnection(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ConnectionFilterer{contract: contract}, nil
}

// bindConnection binds a generic wrapper to an already deployed contract.
func bindConnection(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ConnectionMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Connection *ConnectionRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Connection.Contract.ConnectionCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Connection *ConnectionRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Connection.Contract.ConnectionTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Connection *ConnectionRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Connection.Contract.ConnectionTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Connection *ConnectionCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Connection.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Connection *ConnectionTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Connection.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Connection *ConnectionTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Connection.Contract.contract.Transact(opts, method, params...)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_Connection *ConnectionCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Connection.contract.Call(opts, &out, "admin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_Connection *ConnectionSession) Admin() (common.Address, error) {
	return _Connection.Contract.Admin(&_Connection.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_Connection *ConnectionCallerSession) Admin() (common.Address, error) {
	return _Connection.Contract.Admin(&_Connection.CallOpts)
}

// ConnSn is a free data retrieval call binding the contract method 0x99f1fca7.
//
// Solidity: function connSn() view returns(uint256)
func (_Connection *ConnectionCaller) ConnSn(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Connection.contract.Call(opts, &out, "connSn")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ConnSn is a free data retrieval call binding the contract method 0x99f1fca7.
//
// Solidity: function connSn() view returns(uint256)
func (_Connection *ConnectionSession) ConnSn() (*big.Int, error) {
	return _Connection.Contract.ConnSn(&_Connection.CallOpts)
}

// ConnSn is a free data retrieval call binding the contract method 0x99f1fca7.
//
// Solidity: function connSn() view returns(uint256)
func (_Connection *ConnectionCallerSession) ConnSn() (*big.Int, error) {
	return _Connection.Contract.ConnSn(&_Connection.CallOpts)
}

// GetFee is a free data retrieval call binding the contract method 0x7d4c4f4a.
//
// Solidity: function getFee(string to, bool response) view returns(uint256 fee)
func (_Connection *ConnectionCaller) GetFee(opts *bind.CallOpts, to string, response bool) (*big.Int, error) {
	var out []interface{}
	err := _Connection.contract.Call(opts, &out, "getFee", to, response)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetFee is a free data retrieval call binding the contract method 0x7d4c4f4a.
//
// Solidity: function getFee(string to, bool response) view returns(uint256 fee)
func (_Connection *ConnectionSession) GetFee(to string, response bool) (*big.Int, error) {
	return _Connection.Contract.GetFee(&_Connection.CallOpts, to, response)
}

// GetFee is a free data retrieval call binding the contract method 0x7d4c4f4a.
//
// Solidity: function getFee(string to, bool response) view returns(uint256 fee)
func (_Connection *ConnectionCallerSession) GetFee(to string, response bool) (*big.Int, error) {
	return _Connection.Contract.GetFee(&_Connection.CallOpts, to, response)
}

// GetReceipt is a free data retrieval call binding the contract method 0x9664da0e.
//
// Solidity: function getReceipt(string srcNetwork, uint256 _connSn) view returns(bool)
func (_Connection *ConnectionCaller) GetReceipt(opts *bind.CallOpts, srcNetwork string, _connSn *big.Int) (bool, error) {
	var out []interface{}
	err := _Connection.contract.Call(opts, &out, "getReceipt", srcNetwork, _connSn)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetReceipt is a free data retrieval call binding the contract method 0x9664da0e.
//
// Solidity: function getReceipt(string srcNetwork, uint256 _connSn) view returns(bool)
func (_Connection *ConnectionSession) GetReceipt(srcNetwork string, _connSn *big.Int) (bool, error) {
	return _Connection.Contract.GetReceipt(&_Connection.CallOpts, srcNetwork, _connSn)
}

// GetReceipt is a free data retrieval call binding the contract method 0x9664da0e.
//
// Solidity: function getReceipt(string srcNetwork, uint256 _connSn) view returns(bool)
func (_Connection *ConnectionCallerSession) GetReceipt(srcNetwork string, _connSn *big.Int) (bool, error) {
	return _Connection.Contract.GetReceipt(&_Connection.CallOpts, srcNetwork, _connSn)
}

// ClaimFees is a paid mutator transaction binding the contract method 0xd294f093.
//
// Solidity: function claimFees() returns()
func (_Connection *ConnectionTransactor) ClaimFees(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Connection.contract.Transact(opts, "claimFees")
}

// ClaimFees is a paid mutator transaction binding the contract method 0xd294f093.
//
// Solidity: function claimFees() returns()
func (_Connection *ConnectionSession) ClaimFees() (*types.Transaction, error) {
	return _Connection.Contract.ClaimFees(&_Connection.TransactOpts)
}

// ClaimFees is a paid mutator transaction binding the contract method 0xd294f093.
//
// Solidity: function claimFees() returns()
func (_Connection *ConnectionTransactorSession) ClaimFees() (*types.Transaction, error) {
	return _Connection.Contract.ClaimFees(&_Connection.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _relayer, address _xCall) returns()
func (_Connection *ConnectionTransactor) Initialize(opts *bind.TransactOpts, _relayer common.Address, _xCall common.Address) (*types.Transaction, error) {
	return _Connection.contract.Transact(opts, "initialize", _relayer, _xCall)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _relayer, address _xCall) returns()
func (_Connection *ConnectionSession) Initialize(_relayer common.Address, _xCall common.Address) (*types.Transaction, error) {
	return _Connection.Contract.Initialize(&_Connection.TransactOpts, _relayer, _xCall)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _relayer, address _xCall) returns()
func (_Connection *ConnectionTransactorSession) Initialize(_relayer common.Address, _xCall common.Address) (*types.Transaction, error) {
	return _Connection.Contract.Initialize(&_Connection.TransactOpts, _relayer, _xCall)
}

// RecvMessage is a paid mutator transaction binding the contract method 0xb58b4cec.
//
// Solidity: function recvMessage(string srcNetwork, uint256 _connSn, bytes _msg) returns()
func (_Connection *ConnectionTransactor) RecvMessage(opts *bind.TransactOpts, srcNetwork string, _connSn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _Connection.contract.Transact(opts, "recvMessage", srcNetwork, _connSn, _msg)
}

// RecvMessage is a paid mutator transaction binding the contract method 0xb58b4cec.
//
// Solidity: function recvMessage(string srcNetwork, uint256 _connSn, bytes _msg) returns()
func (_Connection *ConnectionSession) RecvMessage(srcNetwork string, _connSn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _Connection.Contract.RecvMessage(&_Connection.TransactOpts, srcNetwork, _connSn, _msg)
}

// RecvMessage is a paid mutator transaction binding the contract method 0xb58b4cec.
//
// Solidity: function recvMessage(string srcNetwork, uint256 _connSn, bytes _msg) returns()
func (_Connection *ConnectionTransactorSession) RecvMessage(srcNetwork string, _connSn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _Connection.Contract.RecvMessage(&_Connection.TransactOpts, srcNetwork, _connSn, _msg)
}

// RevertMessage is a paid mutator transaction binding the contract method 0x2d3fb823.
//
// Solidity: function revertMessage(uint256 sn) returns()
func (_Connection *ConnectionTransactor) RevertMessage(opts *bind.TransactOpts, sn *big.Int) (*types.Transaction, error) {
	return _Connection.contract.Transact(opts, "revertMessage", sn)
}

// RevertMessage is a paid mutator transaction binding the contract method 0x2d3fb823.
//
// Solidity: function revertMessage(uint256 sn) returns()
func (_Connection *ConnectionSession) RevertMessage(sn *big.Int) (*types.Transaction, error) {
	return _Connection.Contract.RevertMessage(&_Connection.TransactOpts, sn)
}

// RevertMessage is a paid mutator transaction binding the contract method 0x2d3fb823.
//
// Solidity: function revertMessage(uint256 sn) returns()
func (_Connection *ConnectionTransactorSession) RevertMessage(sn *big.Int) (*types.Transaction, error) {
	return _Connection.Contract.RevertMessage(&_Connection.TransactOpts, sn)
}

// SendMessage is a paid mutator transaction binding the contract method 0x522a901e.
//
// Solidity: function sendMessage(string to, string svc, int256 sn, bytes _msg) payable returns()
func (_Connection *ConnectionTransactor) SendMessage(opts *bind.TransactOpts, to string, svc string, sn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _Connection.contract.Transact(opts, "sendMessage", to, svc, sn, _msg)
}

// SendMessage is a paid mutator transaction binding the contract method 0x522a901e.
//
// Solidity: function sendMessage(string to, string svc, int256 sn, bytes _msg) payable returns()
func (_Connection *ConnectionSession) SendMessage(to string, svc string, sn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _Connection.Contract.SendMessage(&_Connection.TransactOpts, to, svc, sn, _msg)
}

// SendMessage is a paid mutator transaction binding the contract method 0x522a901e.
//
// Solidity: function sendMessage(string to, string svc, int256 sn, bytes _msg) payable returns()
func (_Connection *ConnectionTransactorSession) SendMessage(to string, svc string, sn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _Connection.Contract.SendMessage(&_Connection.TransactOpts, to, svc, sn, _msg)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x704b6c02.
//
// Solidity: function setAdmin(address _address) returns()
func (_Connection *ConnectionTransactor) SetAdmin(opts *bind.TransactOpts, _address common.Address) (*types.Transaction, error) {
	return _Connection.contract.Transact(opts, "setAdmin", _address)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x704b6c02.
//
// Solidity: function setAdmin(address _address) returns()
func (_Connection *ConnectionSession) SetAdmin(_address common.Address) (*types.Transaction, error) {
	return _Connection.Contract.SetAdmin(&_Connection.TransactOpts, _address)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x704b6c02.
//
// Solidity: function setAdmin(address _address) returns()
func (_Connection *ConnectionTransactorSession) SetAdmin(_address common.Address) (*types.Transaction, error) {
	return _Connection.Contract.SetAdmin(&_Connection.TransactOpts, _address)
}

// SetFee is a paid mutator transaction binding the contract method 0x43f08a89.
//
// Solidity: function setFee(string networkId, uint256 messageFee, uint256 responseFee) returns()
func (_Connection *ConnectionTransactor) SetFee(opts *bind.TransactOpts, networkId string, messageFee *big.Int, responseFee *big.Int) (*types.Transaction, error) {
	return _Connection.contract.Transact(opts, "setFee", networkId, messageFee, responseFee)
}

// SetFee is a paid mutator transaction binding the contract method 0x43f08a89.
//
// Solidity: function setFee(string networkId, uint256 messageFee, uint256 responseFee) returns()
func (_Connection *ConnectionSession) SetFee(networkId string, messageFee *big.Int, responseFee *big.Int) (*types.Transaction, error) {
	return _Connection.Contract.SetFee(&_Connection.TransactOpts, networkId, messageFee, responseFee)
}

// SetFee is a paid mutator transaction binding the contract method 0x43f08a89.
//
// Solidity: function setFee(string networkId, uint256 messageFee, uint256 responseFee) returns()
func (_Connection *ConnectionTransactorSession) SetFee(networkId string, messageFee *big.Int, responseFee *big.Int) (*types.Transaction, error) {
	return _Connection.Contract.SetFee(&_Connection.TransactOpts, networkId, messageFee, responseFee)
}

// ConnectionInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Connection contract.
type ConnectionInitializedIterator struct {
	Event *ConnectionInitialized // Event containing the contract specifics and raw log

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
func (it *ConnectionInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConnectionInitialized)
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
		it.Event = new(ConnectionInitialized)
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
func (it *ConnectionInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConnectionInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConnectionInitialized represents a Initialized event raised by the Connection contract.
type ConnectionInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Connection *ConnectionFilterer) FilterInitialized(opts *bind.FilterOpts) (*ConnectionInitializedIterator, error) {

	logs, sub, err := _Connection.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ConnectionInitializedIterator{contract: _Connection.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Connection *ConnectionFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ConnectionInitialized) (event.Subscription, error) {

	logs, sub, err := _Connection.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConnectionInitialized)
				if err := _Connection.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Connection *ConnectionFilterer) ParseInitialized(log types.Log) (*ConnectionInitialized, error) {
	event := new(ConnectionInitialized)
	if err := _Connection.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConnectionMessageIterator is returned from FilterMessage and is used to iterate over the raw logs and unpacked data for Message events raised by the Connection contract.
type ConnectionMessageIterator struct {
	Event *ConnectionMessage // Event containing the contract specifics and raw log

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
func (it *ConnectionMessageIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConnectionMessage)
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
		it.Event = new(ConnectionMessage)
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
func (it *ConnectionMessageIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConnectionMessageIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConnectionMessage represents a Message event raised by the Connection contract.
type ConnectionMessage struct {
	TargetNetwork string
	Sn            *big.Int
	Msg           []byte
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterMessage is a free log retrieval operation binding the contract event 0x37be353f216cf7e33639101fd610c542e6a0c0109173fa1c1d8b04d34edb7c1b.
//
// Solidity: event Message(string targetNetwork, uint256 sn, bytes _msg)
func (_Connection *ConnectionFilterer) FilterMessage(opts *bind.FilterOpts) (*ConnectionMessageIterator, error) {

	logs, sub, err := _Connection.contract.FilterLogs(opts, "Message")
	if err != nil {
		return nil, err
	}
	return &ConnectionMessageIterator{contract: _Connection.contract, event: "Message", logs: logs, sub: sub}, nil
}

// WatchMessage is a free log subscription operation binding the contract event 0x37be353f216cf7e33639101fd610c542e6a0c0109173fa1c1d8b04d34edb7c1b.
//
// Solidity: event Message(string targetNetwork, uint256 sn, bytes _msg)
func (_Connection *ConnectionFilterer) WatchMessage(opts *bind.WatchOpts, sink chan<- *ConnectionMessage) (event.Subscription, error) {

	logs, sub, err := _Connection.contract.WatchLogs(opts, "Message")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConnectionMessage)
				if err := _Connection.contract.UnpackLog(event, "Message", log); err != nil {
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
func (_Connection *ConnectionFilterer) ParseMessage(log types.Log) (*ConnectionMessage, error) {
	event := new(ConnectionMessage)
	if err := _Connection.contract.UnpackLog(event, "Message", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
