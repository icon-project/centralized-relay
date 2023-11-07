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

// StorageMetaData contains all meta data concerning the Storage contract.
var StorageMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"targetNetwork\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"sn\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"msg\",\"type\":\"bytes\"}],\"name\":\"Message\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"admin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_to\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"_response\",\"type\":\"bool\"}],\"name\":\"getFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_xCall\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"srcNID\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"sn\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"msg\",\"type\":\"bytes\"}],\"name\":\"recvMessage\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"seenDeliveryVaaHashes\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_to\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_svc\",\"type\":\"string\"},{\"internalType\":\"int256\",\"name\":\"_sn\",\"type\":\"int256\"},{\"internalType\":\"bytes\",\"name\":\"_msg\",\"type\":\"bytes\"}],\"name\":\"sendMessage\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"networkId\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"messageFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"responseFee\",\"type\":\"uint256\"}],\"name\":\"setFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// StorageABI is the input ABI used to generate the binding from.
// Deprecated: Use StorageMetaData.ABI instead.
var StorageABI = StorageMetaData.ABI

// Storage is an auto generated Go binding around an Ethereum contract.
type Storage struct {
	StorageCaller     // Read-only binding to the contract
	StorageTransactor // Write-only binding to the contract
	StorageFilterer   // Log filterer for contract events
}

// StorageCaller is an auto generated read-only Go binding around an Ethereum contract.
type StorageCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StorageTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StorageTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StorageFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StorageFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StorageSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StorageSession struct {
	Contract     *Storage          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StorageCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StorageCallerSession struct {
	Contract *StorageCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// StorageTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StorageTransactorSession struct {
	Contract     *StorageTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// StorageRaw is an auto generated low-level Go binding around an Ethereum contract.
type StorageRaw struct {
	Contract *Storage // Generic contract binding to access the raw methods on
}

// StorageCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StorageCallerRaw struct {
	Contract *StorageCaller // Generic read-only contract binding to access the raw methods on
}

// StorageTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StorageTransactorRaw struct {
	Contract *StorageTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStorage creates a new instance of Storage, bound to a specific deployed contract.
func NewStorage(address common.Address, backend bind.ContractBackend) (*Storage, error) {
	contract, err := bindStorage(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Storage{StorageCaller: StorageCaller{contract: contract}, StorageTransactor: StorageTransactor{contract: contract}, StorageFilterer: StorageFilterer{contract: contract}}, nil
}

// NewStorageCaller creates a new read-only instance of Storage, bound to a specific deployed contract.
func NewStorageCaller(address common.Address, caller bind.ContractCaller) (*StorageCaller, error) {
	contract, err := bindStorage(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StorageCaller{contract: contract}, nil
}

// NewStorageTransactor creates a new write-only instance of Storage, bound to a specific deployed contract.
func NewStorageTransactor(address common.Address, transactor bind.ContractTransactor) (*StorageTransactor, error) {
	contract, err := bindStorage(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StorageTransactor{contract: contract}, nil
}

// NewStorageFilterer creates a new log filterer instance of Storage, bound to a specific deployed contract.
func NewStorageFilterer(address common.Address, filterer bind.ContractFilterer) (*StorageFilterer, error) {
	contract, err := bindStorage(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StorageFilterer{contract: contract}, nil
}

// bindStorage binds a generic wrapper to an already deployed contract.
func bindStorage(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := StorageMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Storage *StorageRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Storage.Contract.StorageCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Storage *StorageRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Storage.Contract.StorageTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Storage *StorageRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Storage.Contract.StorageTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Storage *StorageCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Storage.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Storage *StorageTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Storage.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Storage *StorageTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Storage.Contract.contract.Transact(opts, method, params...)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_Storage *StorageCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Storage.contract.Call(opts, &out, "admin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_Storage *StorageSession) Admin() (common.Address, error) {
	return _Storage.Contract.Admin(&_Storage.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_Storage *StorageCallerSession) Admin() (common.Address, error) {
	return _Storage.Contract.Admin(&_Storage.CallOpts)
}

// GetFee is a free data retrieval call binding the contract method 0x7d4c4f4a.
//
// Solidity: function getFee(string _to, bool _response) view returns(uint256 _fee)
func (_Storage *StorageCaller) GetFee(opts *bind.CallOpts, _to string, _response bool) (*big.Int, error) {
	var out []interface{}
	err := _Storage.contract.Call(opts, &out, "getFee", _to, _response)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetFee is a free data retrieval call binding the contract method 0x7d4c4f4a.
//
// Solidity: function getFee(string _to, bool _response) view returns(uint256 _fee)
func (_Storage *StorageSession) GetFee(_to string, _response bool) (*big.Int, error) {
	return _Storage.Contract.GetFee(&_Storage.CallOpts, _to, _response)
}

// GetFee is a free data retrieval call binding the contract method 0x7d4c4f4a.
//
// Solidity: function getFee(string _to, bool _response) view returns(uint256 _fee)
func (_Storage *StorageCallerSession) GetFee(_to string, _response bool) (*big.Int, error) {
	return _Storage.Contract.GetFee(&_Storage.CallOpts, _to, _response)
}

// SeenDeliveryVaaHashes is a free data retrieval call binding the contract method 0x180f6cc2.
//
// Solidity: function seenDeliveryVaaHashes(bytes32 ) view returns(bool)
func (_Storage *StorageCaller) SeenDeliveryVaaHashes(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _Storage.contract.Call(opts, &out, "seenDeliveryVaaHashes", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SeenDeliveryVaaHashes is a free data retrieval call binding the contract method 0x180f6cc2.
//
// Solidity: function seenDeliveryVaaHashes(bytes32 ) view returns(bool)
func (_Storage *StorageSession) SeenDeliveryVaaHashes(arg0 [32]byte) (bool, error) {
	return _Storage.Contract.SeenDeliveryVaaHashes(&_Storage.CallOpts, arg0)
}

// SeenDeliveryVaaHashes is a free data retrieval call binding the contract method 0x180f6cc2.
//
// Solidity: function seenDeliveryVaaHashes(bytes32 ) view returns(bool)
func (_Storage *StorageCallerSession) SeenDeliveryVaaHashes(arg0 [32]byte) (bool, error) {
	return _Storage.Contract.SeenDeliveryVaaHashes(&_Storage.CallOpts, arg0)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _xCall) returns()
func (_Storage *StorageTransactor) Initialize(opts *bind.TransactOpts, _xCall common.Address) (*types.Transaction, error) {
	return _Storage.contract.Transact(opts, "initialize", _xCall)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _xCall) returns()
func (_Storage *StorageSession) Initialize(_xCall common.Address) (*types.Transaction, error) {
	return _Storage.Contract.Initialize(&_Storage.TransactOpts, _xCall)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _xCall) returns()
func (_Storage *StorageTransactorSession) Initialize(_xCall common.Address) (*types.Transaction, error) {
	return _Storage.Contract.Initialize(&_Storage.TransactOpts, _xCall)
}

// RecvMessage is a paid mutator transaction binding the contract method 0x82d14c54.
//
// Solidity: function recvMessage(string srcNID, string sn, bytes msg) payable returns()
func (_Storage *StorageTransactor) RecvMessage(opts *bind.TransactOpts, srcNID string, sn string, msg []byte) (*types.Transaction, error) {
	return _Storage.contract.Transact(opts, "recvMessage", srcNID, sn, msg)
}

// RecvMessage is a paid mutator transaction binding the contract method 0x82d14c54.
//
// Solidity: function recvMessage(string srcNID, string sn, bytes msg) payable returns()
func (_Storage *StorageSession) RecvMessage(srcNID string, sn string, msg []byte) (*types.Transaction, error) {
	return _Storage.Contract.RecvMessage(&_Storage.TransactOpts, srcNID, sn, msg)
}

// RecvMessage is a paid mutator transaction binding the contract method 0x82d14c54.
//
// Solidity: function recvMessage(string srcNID, string sn, bytes msg) payable returns()
func (_Storage *StorageTransactorSession) RecvMessage(srcNID string, sn string, msg []byte) (*types.Transaction, error) {
	return _Storage.Contract.RecvMessage(&_Storage.TransactOpts, srcNID, sn, msg)
}

// SendMessage is a paid mutator transaction binding the contract method 0x522a901e.
//
// Solidity: function sendMessage(string _to, string _svc, int256 _sn, bytes _msg) payable returns()
func (_Storage *StorageTransactor) SendMessage(opts *bind.TransactOpts, _to string, _svc string, _sn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _Storage.contract.Transact(opts, "sendMessage", _to, _svc, _sn, _msg)
}

// SendMessage is a paid mutator transaction binding the contract method 0x522a901e.
//
// Solidity: function sendMessage(string _to, string _svc, int256 _sn, bytes _msg) payable returns()
func (_Storage *StorageSession) SendMessage(_to string, _svc string, _sn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _Storage.Contract.SendMessage(&_Storage.TransactOpts, _to, _svc, _sn, _msg)
}

// SendMessage is a paid mutator transaction binding the contract method 0x522a901e.
//
// Solidity: function sendMessage(string _to, string _svc, int256 _sn, bytes _msg) payable returns()
func (_Storage *StorageTransactorSession) SendMessage(_to string, _svc string, _sn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _Storage.Contract.SendMessage(&_Storage.TransactOpts, _to, _svc, _sn, _msg)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x704b6c02.
//
// Solidity: function setAdmin(address _address) returns()
func (_Storage *StorageTransactor) SetAdmin(opts *bind.TransactOpts, _address common.Address) (*types.Transaction, error) {
	return _Storage.contract.Transact(opts, "setAdmin", _address)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x704b6c02.
//
// Solidity: function setAdmin(address _address) returns()
func (_Storage *StorageSession) SetAdmin(_address common.Address) (*types.Transaction, error) {
	return _Storage.Contract.SetAdmin(&_Storage.TransactOpts, _address)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x704b6c02.
//
// Solidity: function setAdmin(address _address) returns()
func (_Storage *StorageTransactorSession) SetAdmin(_address common.Address) (*types.Transaction, error) {
	return _Storage.Contract.SetAdmin(&_Storage.TransactOpts, _address)
}

// SetFee is a paid mutator transaction binding the contract method 0x43f08a89.
//
// Solidity: function setFee(string networkId, uint256 messageFee, uint256 responseFee) returns()
func (_Storage *StorageTransactor) SetFee(opts *bind.TransactOpts, networkId string, messageFee *big.Int, responseFee *big.Int) (*types.Transaction, error) {
	return _Storage.contract.Transact(opts, "setFee", networkId, messageFee, responseFee)
}

// SetFee is a paid mutator transaction binding the contract method 0x43f08a89.
//
// Solidity: function setFee(string networkId, uint256 messageFee, uint256 responseFee) returns()
func (_Storage *StorageSession) SetFee(networkId string, messageFee *big.Int, responseFee *big.Int) (*types.Transaction, error) {
	return _Storage.Contract.SetFee(&_Storage.TransactOpts, networkId, messageFee, responseFee)
}

// SetFee is a paid mutator transaction binding the contract method 0x43f08a89.
//
// Solidity: function setFee(string networkId, uint256 messageFee, uint256 responseFee) returns()
func (_Storage *StorageTransactorSession) SetFee(networkId string, messageFee *big.Int, responseFee *big.Int) (*types.Transaction, error) {
	return _Storage.Contract.SetFee(&_Storage.TransactOpts, networkId, messageFee, responseFee)
}

// StorageInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Storage contract.
type StorageInitializedIterator struct {
	Event *StorageInitialized // Event containing the contract specifics and raw log

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
func (it *StorageInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StorageInitialized)
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
		it.Event = new(StorageInitialized)
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
func (it *StorageInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StorageInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StorageInitialized represents a Initialized event raised by the Storage contract.
type StorageInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Storage *StorageFilterer) FilterInitialized(opts *bind.FilterOpts) (*StorageInitializedIterator, error) {

	logs, sub, err := _Storage.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &StorageInitializedIterator{contract: _Storage.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Storage *StorageFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *StorageInitialized) (event.Subscription, error) {

	logs, sub, err := _Storage.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StorageInitialized)
				if err := _Storage.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Storage *StorageFilterer) ParseInitialized(log types.Log) (*StorageInitialized, error) {
	event := new(StorageInitialized)
	if err := _Storage.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StorageMessageIterator is returned from FilterMessage and is used to iterate over the raw logs and unpacked data for Message events raised by the Storage contract.
type StorageMessageIterator struct {
	Event *StorageMessage // Event containing the contract specifics and raw log

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
func (it *StorageMessageIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StorageMessage)
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
		it.Event = new(StorageMessage)
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
func (it *StorageMessageIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StorageMessageIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StorageMessage represents a Message event raised by the Storage contract.
type StorageMessage struct {
	TargetNetwork string
	Sn            *big.Int
	Msg           []byte
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterMessage is a free log retrieval operation binding the contract event 0x6dbbb5c83189670e066d281dfc37d9ded5132af5d6401cfc831c7499eb775f3d.
//
// Solidity: event Message(string targetNetwork, int256 sn, bytes msg)
func (_Storage *StorageFilterer) FilterMessage(opts *bind.FilterOpts) (*StorageMessageIterator, error) {

	logs, sub, err := _Storage.contract.FilterLogs(opts, "Message")
	if err != nil {
		return nil, err
	}
	return &StorageMessageIterator{contract: _Storage.contract, event: "Message", logs: logs, sub: sub}, nil
}

// WatchMessage is a free log subscription operation binding the contract event 0x6dbbb5c83189670e066d281dfc37d9ded5132af5d6401cfc831c7499eb775f3d.
//
// Solidity: event Message(string targetNetwork, int256 sn, bytes msg)
func (_Storage *StorageFilterer) WatchMessage(opts *bind.WatchOpts, sink chan<- *StorageMessage) (event.Subscription, error) {

	logs, sub, err := _Storage.contract.WatchLogs(opts, "Message")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StorageMessage)
				if err := _Storage.contract.UnpackLog(event, "Message", log); err != nil {
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

// ParseMessage is a log parse operation binding the contract event 0x6dbbb5c83189670e066d281dfc37d9ded5132af5d6401cfc831c7499eb775f3d.
//
// Solidity: event Message(string targetNetwork, int256 sn, bytes msg)
func (_Storage *StorageFilterer) ParseMessage(log types.Log) (*StorageMessage, error) {
	event := new(StorageMessage)
	if err := _Storage.contract.UnpackLog(event, "Message", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
