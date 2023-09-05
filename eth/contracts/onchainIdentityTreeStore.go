// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

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

// OnchainIdentityTreeStoreMetaData contains all meta data concerning the OnchainIdentityTreeStore contract.
var OnchainIdentityTreeStoreMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"preimage\",\"type\":\"uint256[]\"}],\"name\":\"addNode\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[][]\",\"name\":\"preimage\",\"type\":\"uint256[][]\"}],\"name\":\"addNodes\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getNode\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"preimage\",\"type\":\"uint256[]\"}],\"name\":\"hash\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
}

// OnchainIdentityTreeStoreABI is the input ABI used to generate the binding from.
// Deprecated: Use OnchainIdentityTreeStoreMetaData.ABI instead.
var OnchainIdentityTreeStoreABI = OnchainIdentityTreeStoreMetaData.ABI

// OnchainIdentityTreeStore is an auto generated Go binding around an Ethereum contract.
type OnchainIdentityTreeStore struct {
	OnchainIdentityTreeStoreCaller     // Read-only binding to the contract
	OnchainIdentityTreeStoreTransactor // Write-only binding to the contract
	OnchainIdentityTreeStoreFilterer   // Log filterer for contract events
}

// OnchainIdentityTreeStoreCaller is an auto generated read-only Go binding around an Ethereum contract.
type OnchainIdentityTreeStoreCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OnchainIdentityTreeStoreTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OnchainIdentityTreeStoreTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OnchainIdentityTreeStoreFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OnchainIdentityTreeStoreFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OnchainIdentityTreeStoreSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OnchainIdentityTreeStoreSession struct {
	Contract     *OnchainIdentityTreeStore // Generic contract binding to set the session for
	CallOpts     bind.CallOpts             // Call options to use throughout this session
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// OnchainIdentityTreeStoreCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OnchainIdentityTreeStoreCallerSession struct {
	Contract *OnchainIdentityTreeStoreCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                   // Call options to use throughout this session
}

// OnchainIdentityTreeStoreTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OnchainIdentityTreeStoreTransactorSession struct {
	Contract     *OnchainIdentityTreeStoreTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                   // Transaction auth options to use throughout this session
}

// OnchainIdentityTreeStoreRaw is an auto generated low-level Go binding around an Ethereum contract.
type OnchainIdentityTreeStoreRaw struct {
	Contract *OnchainIdentityTreeStore // Generic contract binding to access the raw methods on
}

// OnchainIdentityTreeStoreCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OnchainIdentityTreeStoreCallerRaw struct {
	Contract *OnchainIdentityTreeStoreCaller // Generic read-only contract binding to access the raw methods on
}

// OnchainIdentityTreeStoreTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OnchainIdentityTreeStoreTransactorRaw struct {
	Contract *OnchainIdentityTreeStoreTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOnchainIdentityTreeStore creates a new instance of OnchainIdentityTreeStore, bound to a specific deployed contract.
func NewOnchainIdentityTreeStore(address common.Address, backend bind.ContractBackend) (*OnchainIdentityTreeStore, error) {
	contract, err := bindOnchainIdentityTreeStore(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OnchainIdentityTreeStore{OnchainIdentityTreeStoreCaller: OnchainIdentityTreeStoreCaller{contract: contract}, OnchainIdentityTreeStoreTransactor: OnchainIdentityTreeStoreTransactor{contract: contract}, OnchainIdentityTreeStoreFilterer: OnchainIdentityTreeStoreFilterer{contract: contract}}, nil
}

// NewOnchainIdentityTreeStoreCaller creates a new read-only instance of OnchainIdentityTreeStore, bound to a specific deployed contract.
func NewOnchainIdentityTreeStoreCaller(address common.Address, caller bind.ContractCaller) (*OnchainIdentityTreeStoreCaller, error) {
	contract, err := bindOnchainIdentityTreeStore(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OnchainIdentityTreeStoreCaller{contract: contract}, nil
}

// NewOnchainIdentityTreeStoreTransactor creates a new write-only instance of OnchainIdentityTreeStore, bound to a specific deployed contract.
func NewOnchainIdentityTreeStoreTransactor(address common.Address, transactor bind.ContractTransactor) (*OnchainIdentityTreeStoreTransactor, error) {
	contract, err := bindOnchainIdentityTreeStore(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OnchainIdentityTreeStoreTransactor{contract: contract}, nil
}

// NewOnchainIdentityTreeStoreFilterer creates a new log filterer instance of OnchainIdentityTreeStore, bound to a specific deployed contract.
func NewOnchainIdentityTreeStoreFilterer(address common.Address, filterer bind.ContractFilterer) (*OnchainIdentityTreeStoreFilterer, error) {
	contract, err := bindOnchainIdentityTreeStore(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OnchainIdentityTreeStoreFilterer{contract: contract}, nil
}

// bindOnchainIdentityTreeStore binds a generic wrapper to an already deployed contract.
func bindOnchainIdentityTreeStore(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OnchainIdentityTreeStoreMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OnchainIdentityTreeStore *OnchainIdentityTreeStoreRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OnchainIdentityTreeStore.Contract.OnchainIdentityTreeStoreCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OnchainIdentityTreeStore *OnchainIdentityTreeStoreRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OnchainIdentityTreeStore.Contract.OnchainIdentityTreeStoreTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OnchainIdentityTreeStore *OnchainIdentityTreeStoreRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OnchainIdentityTreeStore.Contract.OnchainIdentityTreeStoreTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OnchainIdentityTreeStore *OnchainIdentityTreeStoreCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OnchainIdentityTreeStore.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OnchainIdentityTreeStore *OnchainIdentityTreeStoreTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OnchainIdentityTreeStore.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OnchainIdentityTreeStore *OnchainIdentityTreeStoreTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OnchainIdentityTreeStore.Contract.contract.Transact(opts, method, params...)
}

// GetNode is a free data retrieval call binding the contract method 0x4f0f4aa9.
//
// Solidity: function getNode(uint256 id) view returns(uint256[])
func (_OnchainIdentityTreeStore *OnchainIdentityTreeStoreCaller) GetNode(opts *bind.CallOpts, id *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _OnchainIdentityTreeStore.contract.Call(opts, &out, "getNode", id)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetNode is a free data retrieval call binding the contract method 0x4f0f4aa9.
//
// Solidity: function getNode(uint256 id) view returns(uint256[])
func (_OnchainIdentityTreeStore *OnchainIdentityTreeStoreSession) GetNode(id *big.Int) ([]*big.Int, error) {
	return _OnchainIdentityTreeStore.Contract.GetNode(&_OnchainIdentityTreeStore.CallOpts, id)
}

// GetNode is a free data retrieval call binding the contract method 0x4f0f4aa9.
//
// Solidity: function getNode(uint256 id) view returns(uint256[])
func (_OnchainIdentityTreeStore *OnchainIdentityTreeStoreCallerSession) GetNode(id *big.Int) ([]*big.Int, error) {
	return _OnchainIdentityTreeStore.Contract.GetNode(&_OnchainIdentityTreeStore.CallOpts, id)
}

// Hash is a free data retrieval call binding the contract method 0x40ec6e49.
//
// Solidity: function hash(uint256[] preimage) pure returns(uint256)
func (_OnchainIdentityTreeStore *OnchainIdentityTreeStoreCaller) Hash(opts *bind.CallOpts, preimage []*big.Int) (*big.Int, error) {
	var out []interface{}
	err := _OnchainIdentityTreeStore.contract.Call(opts, &out, "hash", preimage)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Hash is a free data retrieval call binding the contract method 0x40ec6e49.
//
// Solidity: function hash(uint256[] preimage) pure returns(uint256)
func (_OnchainIdentityTreeStore *OnchainIdentityTreeStoreSession) Hash(preimage []*big.Int) (*big.Int, error) {
	return _OnchainIdentityTreeStore.Contract.Hash(&_OnchainIdentityTreeStore.CallOpts, preimage)
}

// Hash is a free data retrieval call binding the contract method 0x40ec6e49.
//
// Solidity: function hash(uint256[] preimage) pure returns(uint256)
func (_OnchainIdentityTreeStore *OnchainIdentityTreeStoreCallerSession) Hash(preimage []*big.Int) (*big.Int, error) {
	return _OnchainIdentityTreeStore.Contract.Hash(&_OnchainIdentityTreeStore.CallOpts, preimage)
}

// AddNode is a paid mutator transaction binding the contract method 0x34e80e71.
//
// Solidity: function addNode(uint256[] preimage) returns()
func (_OnchainIdentityTreeStore *OnchainIdentityTreeStoreTransactor) AddNode(opts *bind.TransactOpts, preimage []*big.Int) (*types.Transaction, error) {
	return _OnchainIdentityTreeStore.contract.Transact(opts, "addNode", preimage)
}

// AddNode is a paid mutator transaction binding the contract method 0x34e80e71.
//
// Solidity: function addNode(uint256[] preimage) returns()
func (_OnchainIdentityTreeStore *OnchainIdentityTreeStoreSession) AddNode(preimage []*big.Int) (*types.Transaction, error) {
	return _OnchainIdentityTreeStore.Contract.AddNode(&_OnchainIdentityTreeStore.TransactOpts, preimage)
}

// AddNode is a paid mutator transaction binding the contract method 0x34e80e71.
//
// Solidity: function addNode(uint256[] preimage) returns()
func (_OnchainIdentityTreeStore *OnchainIdentityTreeStoreTransactorSession) AddNode(preimage []*big.Int) (*types.Transaction, error) {
	return _OnchainIdentityTreeStore.Contract.AddNode(&_OnchainIdentityTreeStore.TransactOpts, preimage)
}

// AddNodes is a paid mutator transaction binding the contract method 0xa00f738d.
//
// Solidity: function addNodes(uint256[][] preimage) returns()
func (_OnchainIdentityTreeStore *OnchainIdentityTreeStoreTransactor) AddNodes(opts *bind.TransactOpts, preimage [][]*big.Int) (*types.Transaction, error) {
	return _OnchainIdentityTreeStore.contract.Transact(opts, "addNodes", preimage)
}

// AddNodes is a paid mutator transaction binding the contract method 0xa00f738d.
//
// Solidity: function addNodes(uint256[][] preimage) returns()
func (_OnchainIdentityTreeStore *OnchainIdentityTreeStoreSession) AddNodes(preimage [][]*big.Int) (*types.Transaction, error) {
	return _OnchainIdentityTreeStore.Contract.AddNodes(&_OnchainIdentityTreeStore.TransactOpts, preimage)
}

// AddNodes is a paid mutator transaction binding the contract method 0xa00f738d.
//
// Solidity: function addNodes(uint256[][] preimage) returns()
func (_OnchainIdentityTreeStore *OnchainIdentityTreeStoreTransactorSession) AddNodes(preimage [][]*big.Int) (*types.Transaction, error) {
	return _OnchainIdentityTreeStore.Contract.AddNodes(&_OnchainIdentityTreeStore.TransactOpts, preimage)
}
