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

// ISwapRouterExactInputParams is an auto generated low-level Go binding around an user-defined struct.
type ISwapRouterExactInputParams struct {
	Path             []byte
	Recipient        common.Address
	Deadline         *big.Int
	AmountIn         *big.Int
	AmountOutMinimum *big.Int
}

// ISwapRouterExactInputSingleParams is an auto generated low-level Go binding around an user-defined struct.
type ISwapRouterExactInputSingleParams struct {
	TokenIn           common.Address
	TokenOut          common.Address
	Fee               *big.Int
	Recipient         common.Address
	Deadline          *big.Int
	AmountIn          *big.Int
	AmountOutMinimum  *big.Int
	SqrtPriceLimitX96 *big.Int
}

// ISwapRouterExactOutputParams is an auto generated low-level Go binding around an user-defined struct.
type ISwapRouterExactOutputParams struct {
	Path            []byte
	Recipient       common.Address
	Deadline        *big.Int
	AmountOut       *big.Int
	AmountInMaximum *big.Int
}

// ISwapRouterExactOutputSingleParams is an auto generated low-level Go binding around an user-defined struct.
type ISwapRouterExactOutputSingleParams struct {
	TokenIn           common.Address
	TokenOut          common.Address
	Fee               *big.Int
	Recipient         common.Address
	Deadline          *big.Int
	AmountOut         *big.Int
	AmountInMaximum   *big.Int
	SqrtPriceLimitX96 *big.Int
}

// IswaprouterMetaData contains all meta data concerning the Iswaprouter contract.
var IswaprouterMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"exactInput\",\"inputs\":[{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structISwapRouter.ExactInputParams\",\"components\":[{\"name\":\"path\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amountIn\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amountOutMinimum\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[{\"name\":\"amountOut\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"exactInputSingle\",\"inputs\":[{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structISwapRouter.ExactInputSingleParams\",\"components\":[{\"name\":\"tokenIn\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenOut\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"fee\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amountIn\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amountOutMinimum\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"sqrtPriceLimitX96\",\"type\":\"uint160\",\"internalType\":\"uint160\"}]}],\"outputs\":[{\"name\":\"amountOut\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"exactOutput\",\"inputs\":[{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structISwapRouter.ExactOutputParams\",\"components\":[{\"name\":\"path\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amountOut\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amountInMaximum\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[{\"name\":\"amountIn\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"exactOutputSingle\",\"inputs\":[{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structISwapRouter.ExactOutputSingleParams\",\"components\":[{\"name\":\"tokenIn\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenOut\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"fee\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amountOut\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amountInMaximum\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"sqrtPriceLimitX96\",\"type\":\"uint160\",\"internalType\":\"uint160\"}]}],\"outputs\":[{\"name\":\"amountIn\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"uniswapV3SwapCallback\",\"inputs\":[{\"name\":\"amount0Delta\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"amount1Delta\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"}]",
}

// IswaprouterABI is the input ABI used to generate the binding from.
// Deprecated: Use IswaprouterMetaData.ABI instead.
var IswaprouterABI = IswaprouterMetaData.ABI

// Iswaprouter is an auto generated Go binding around an Ethereum contract.
type Iswaprouter struct {
	IswaprouterCaller     // Read-only binding to the contract
	IswaprouterTransactor // Write-only binding to the contract
	IswaprouterFilterer   // Log filterer for contract events
}

// IswaprouterCaller is an auto generated read-only Go binding around an Ethereum contract.
type IswaprouterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IswaprouterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IswaprouterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IswaprouterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IswaprouterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IswaprouterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IswaprouterSession struct {
	Contract     *Iswaprouter      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IswaprouterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IswaprouterCallerSession struct {
	Contract *IswaprouterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// IswaprouterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IswaprouterTransactorSession struct {
	Contract     *IswaprouterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// IswaprouterRaw is an auto generated low-level Go binding around an Ethereum contract.
type IswaprouterRaw struct {
	Contract *Iswaprouter // Generic contract binding to access the raw methods on
}

// IswaprouterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IswaprouterCallerRaw struct {
	Contract *IswaprouterCaller // Generic read-only contract binding to access the raw methods on
}

// IswaprouterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IswaprouterTransactorRaw struct {
	Contract *IswaprouterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIswaprouter creates a new instance of Iswaprouter, bound to a specific deployed contract.
func NewIswaprouter(address common.Address, backend bind.ContractBackend) (*Iswaprouter, error) {
	contract, err := bindIswaprouter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Iswaprouter{IswaprouterCaller: IswaprouterCaller{contract: contract}, IswaprouterTransactor: IswaprouterTransactor{contract: contract}, IswaprouterFilterer: IswaprouterFilterer{contract: contract}}, nil
}

// NewIswaprouterCaller creates a new read-only instance of Iswaprouter, bound to a specific deployed contract.
func NewIswaprouterCaller(address common.Address, caller bind.ContractCaller) (*IswaprouterCaller, error) {
	contract, err := bindIswaprouter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IswaprouterCaller{contract: contract}, nil
}

// NewIswaprouterTransactor creates a new write-only instance of Iswaprouter, bound to a specific deployed contract.
func NewIswaprouterTransactor(address common.Address, transactor bind.ContractTransactor) (*IswaprouterTransactor, error) {
	contract, err := bindIswaprouter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IswaprouterTransactor{contract: contract}, nil
}

// NewIswaprouterFilterer creates a new log filterer instance of Iswaprouter, bound to a specific deployed contract.
func NewIswaprouterFilterer(address common.Address, filterer bind.ContractFilterer) (*IswaprouterFilterer, error) {
	contract, err := bindIswaprouter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IswaprouterFilterer{contract: contract}, nil
}

// bindIswaprouter binds a generic wrapper to an already deployed contract.
func bindIswaprouter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IswaprouterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Iswaprouter *IswaprouterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Iswaprouter.Contract.IswaprouterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Iswaprouter *IswaprouterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Iswaprouter.Contract.IswaprouterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Iswaprouter *IswaprouterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Iswaprouter.Contract.IswaprouterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Iswaprouter *IswaprouterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Iswaprouter.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Iswaprouter *IswaprouterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Iswaprouter.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Iswaprouter *IswaprouterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Iswaprouter.Contract.contract.Transact(opts, method, params...)
}

// ExactInput is a paid mutator transaction binding the contract method 0xc04b8d59.
//
// Solidity: function exactInput((bytes,address,uint256,uint256,uint256) params) payable returns(uint256 amountOut)
func (_Iswaprouter *IswaprouterTransactor) ExactInput(opts *bind.TransactOpts, params ISwapRouterExactInputParams) (*types.Transaction, error) {
	return _Iswaprouter.contract.Transact(opts, "exactInput", params)
}

// ExactInput is a paid mutator transaction binding the contract method 0xc04b8d59.
//
// Solidity: function exactInput((bytes,address,uint256,uint256,uint256) params) payable returns(uint256 amountOut)
func (_Iswaprouter *IswaprouterSession) ExactInput(params ISwapRouterExactInputParams) (*types.Transaction, error) {
	return _Iswaprouter.Contract.ExactInput(&_Iswaprouter.TransactOpts, params)
}

// ExactInput is a paid mutator transaction binding the contract method 0xc04b8d59.
//
// Solidity: function exactInput((bytes,address,uint256,uint256,uint256) params) payable returns(uint256 amountOut)
func (_Iswaprouter *IswaprouterTransactorSession) ExactInput(params ISwapRouterExactInputParams) (*types.Transaction, error) {
	return _Iswaprouter.Contract.ExactInput(&_Iswaprouter.TransactOpts, params)
}

// ExactInputSingle is a paid mutator transaction binding the contract method 0x414bf389.
//
// Solidity: function exactInputSingle((address,address,uint24,address,uint256,uint256,uint256,uint160) params) payable returns(uint256 amountOut)
func (_Iswaprouter *IswaprouterTransactor) ExactInputSingle(opts *bind.TransactOpts, params ISwapRouterExactInputSingleParams) (*types.Transaction, error) {
	return _Iswaprouter.contract.Transact(opts, "exactInputSingle", params)
}

// ExactInputSingle is a paid mutator transaction binding the contract method 0x414bf389.
//
// Solidity: function exactInputSingle((address,address,uint24,address,uint256,uint256,uint256,uint160) params) payable returns(uint256 amountOut)
func (_Iswaprouter *IswaprouterSession) ExactInputSingle(params ISwapRouterExactInputSingleParams) (*types.Transaction, error) {
	return _Iswaprouter.Contract.ExactInputSingle(&_Iswaprouter.TransactOpts, params)
}

// ExactInputSingle is a paid mutator transaction binding the contract method 0x414bf389.
//
// Solidity: function exactInputSingle((address,address,uint24,address,uint256,uint256,uint256,uint160) params) payable returns(uint256 amountOut)
func (_Iswaprouter *IswaprouterTransactorSession) ExactInputSingle(params ISwapRouterExactInputSingleParams) (*types.Transaction, error) {
	return _Iswaprouter.Contract.ExactInputSingle(&_Iswaprouter.TransactOpts, params)
}

// ExactOutput is a paid mutator transaction binding the contract method 0xf28c0498.
//
// Solidity: function exactOutput((bytes,address,uint256,uint256,uint256) params) payable returns(uint256 amountIn)
func (_Iswaprouter *IswaprouterTransactor) ExactOutput(opts *bind.TransactOpts, params ISwapRouterExactOutputParams) (*types.Transaction, error) {
	return _Iswaprouter.contract.Transact(opts, "exactOutput", params)
}

// ExactOutput is a paid mutator transaction binding the contract method 0xf28c0498.
//
// Solidity: function exactOutput((bytes,address,uint256,uint256,uint256) params) payable returns(uint256 amountIn)
func (_Iswaprouter *IswaprouterSession) ExactOutput(params ISwapRouterExactOutputParams) (*types.Transaction, error) {
	return _Iswaprouter.Contract.ExactOutput(&_Iswaprouter.TransactOpts, params)
}

// ExactOutput is a paid mutator transaction binding the contract method 0xf28c0498.
//
// Solidity: function exactOutput((bytes,address,uint256,uint256,uint256) params) payable returns(uint256 amountIn)
func (_Iswaprouter *IswaprouterTransactorSession) ExactOutput(params ISwapRouterExactOutputParams) (*types.Transaction, error) {
	return _Iswaprouter.Contract.ExactOutput(&_Iswaprouter.TransactOpts, params)
}

// ExactOutputSingle is a paid mutator transaction binding the contract method 0xdb3e2198.
//
// Solidity: function exactOutputSingle((address,address,uint24,address,uint256,uint256,uint256,uint160) params) payable returns(uint256 amountIn)
func (_Iswaprouter *IswaprouterTransactor) ExactOutputSingle(opts *bind.TransactOpts, params ISwapRouterExactOutputSingleParams) (*types.Transaction, error) {
	return _Iswaprouter.contract.Transact(opts, "exactOutputSingle", params)
}

// ExactOutputSingle is a paid mutator transaction binding the contract method 0xdb3e2198.
//
// Solidity: function exactOutputSingle((address,address,uint24,address,uint256,uint256,uint256,uint160) params) payable returns(uint256 amountIn)
func (_Iswaprouter *IswaprouterSession) ExactOutputSingle(params ISwapRouterExactOutputSingleParams) (*types.Transaction, error) {
	return _Iswaprouter.Contract.ExactOutputSingle(&_Iswaprouter.TransactOpts, params)
}

// ExactOutputSingle is a paid mutator transaction binding the contract method 0xdb3e2198.
//
// Solidity: function exactOutputSingle((address,address,uint24,address,uint256,uint256,uint256,uint160) params) payable returns(uint256 amountIn)
func (_Iswaprouter *IswaprouterTransactorSession) ExactOutputSingle(params ISwapRouterExactOutputSingleParams) (*types.Transaction, error) {
	return _Iswaprouter.Contract.ExactOutputSingle(&_Iswaprouter.TransactOpts, params)
}

// UniswapV3SwapCallback is a paid mutator transaction binding the contract method 0xfa461e33.
//
// Solidity: function uniswapV3SwapCallback(int256 amount0Delta, int256 amount1Delta, bytes data) returns()
func (_Iswaprouter *IswaprouterTransactor) UniswapV3SwapCallback(opts *bind.TransactOpts, amount0Delta *big.Int, amount1Delta *big.Int, data []byte) (*types.Transaction, error) {
	return _Iswaprouter.contract.Transact(opts, "uniswapV3SwapCallback", amount0Delta, amount1Delta, data)
}

// UniswapV3SwapCallback is a paid mutator transaction binding the contract method 0xfa461e33.
//
// Solidity: function uniswapV3SwapCallback(int256 amount0Delta, int256 amount1Delta, bytes data) returns()
func (_Iswaprouter *IswaprouterSession) UniswapV3SwapCallback(amount0Delta *big.Int, amount1Delta *big.Int, data []byte) (*types.Transaction, error) {
	return _Iswaprouter.Contract.UniswapV3SwapCallback(&_Iswaprouter.TransactOpts, amount0Delta, amount1Delta, data)
}

// UniswapV3SwapCallback is a paid mutator transaction binding the contract method 0xfa461e33.
//
// Solidity: function uniswapV3SwapCallback(int256 amount0Delta, int256 amount1Delta, bytes data) returns()
func (_Iswaprouter *IswaprouterTransactorSession) UniswapV3SwapCallback(amount0Delta *big.Int, amount1Delta *big.Int, data []byte) (*types.Transaction, error) {
	return _Iswaprouter.Contract.UniswapV3SwapCallback(&_Iswaprouter.TransactOpts, amount0Delta, amount1Delta, data)
}
