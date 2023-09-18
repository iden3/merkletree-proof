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

// IOnchainCredentialStatusResolverCredentialStatus is an auto generated low-level Go binding around an user-defined struct.
type IOnchainCredentialStatusResolverCredentialStatus struct {
	Issuer IOnchainCredentialStatusResolverIdentityStateRoots
	Mtp    IOnchainCredentialStatusResolverProof
}

// IOnchainCredentialStatusResolverIdentityStateRoots is an auto generated low-level Go binding around an user-defined struct.
type IOnchainCredentialStatusResolverIdentityStateRoots struct {
	State              *big.Int
	ClaimsTreeRoot     *big.Int
	RevocationTreeRoot *big.Int
	RootOfRoots        *big.Int
}

// IOnchainCredentialStatusResolverProof is an auto generated low-level Go binding around an user-defined struct.
type IOnchainCredentialStatusResolverProof struct {
	Root         *big.Int
	Existence    bool
	Siblings     []*big.Int
	Index        *big.Int
	Value        *big.Int
	AuxExistence bool
	AuxIndex     *big.Int
	AuxValue     *big.Int
}

// IdentityTreeStoreMetaData contains all meta data concerning the IdentityTreeStore contract.
var IdentityTreeStoreMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"state\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"preimage\",\"type\":\"uint256[]\"}],\"name\":\"_hashFunc\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getNode\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"name\":\"getRevocationStatus\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"state\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"claimsTreeRoot\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"revocationTreeRoot\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"rootOfRoots\",\"type\":\"uint256\"}],\"internalType\":\"structIOnchainCredentialStatusResolver.IdentityStateRoots\",\"name\":\"issuer\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"root\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"existence\",\"type\":\"bool\"},{\"internalType\":\"uint256[]\",\"name\":\"siblings\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"auxExistence\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"auxIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"auxValue\",\"type\":\"uint256\"}],\"internalType\":\"structIOnchainCredentialStatusResolver.Proof\",\"name\":\"mtp\",\"type\":\"tuple\"}],\"internalType\":\"structIOnchainCredentialStatusResolver.CredentialStatus\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"state\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"name\":\"getRevocationStatusByIdAndState\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"state\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"claimsTreeRoot\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"revocationTreeRoot\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"rootOfRoots\",\"type\":\"uint256\"}],\"internalType\":\"structIOnchainCredentialStatusResolver.IdentityStateRoots\",\"name\":\"issuer\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"root\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"existence\",\"type\":\"bool\"},{\"internalType\":\"uint256[]\",\"name\":\"siblings\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"auxExistence\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"auxIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"auxValue\",\"type\":\"uint256\"}],\"internalType\":\"structIOnchainCredentialStatusResolver.Proof\",\"name\":\"mtp\",\"type\":\"tuple\"}],\"internalType\":\"structIOnchainCredentialStatusResolver.CredentialStatus\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[][]\",\"name\":\"preimage\",\"type\":\"uint256[][]\"}],\"name\":\"saveNodes\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// IdentityTreeStoreABI is the input ABI used to generate the binding from.
// Deprecated: Use IdentityTreeStoreMetaData.ABI instead.
var IdentityTreeStoreABI = IdentityTreeStoreMetaData.ABI

// IdentityTreeStore is an auto generated Go binding around an Ethereum contract.
type IdentityTreeStore struct {
	IdentityTreeStoreCaller     // Read-only binding to the contract
	IdentityTreeStoreTransactor // Write-only binding to the contract
	IdentityTreeStoreFilterer   // Log filterer for contract events
}

// IdentityTreeStoreCaller is an auto generated read-only Go binding around an Ethereum contract.
type IdentityTreeStoreCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IdentityTreeStoreTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IdentityTreeStoreTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IdentityTreeStoreFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IdentityTreeStoreFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IdentityTreeStoreSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IdentityTreeStoreSession struct {
	Contract     *IdentityTreeStore // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// IdentityTreeStoreCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IdentityTreeStoreCallerSession struct {
	Contract *IdentityTreeStoreCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// IdentityTreeStoreTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IdentityTreeStoreTransactorSession struct {
	Contract     *IdentityTreeStoreTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// IdentityTreeStoreRaw is an auto generated low-level Go binding around an Ethereum contract.
type IdentityTreeStoreRaw struct {
	Contract *IdentityTreeStore // Generic contract binding to access the raw methods on
}

// IdentityTreeStoreCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IdentityTreeStoreCallerRaw struct {
	Contract *IdentityTreeStoreCaller // Generic read-only contract binding to access the raw methods on
}

// IdentityTreeStoreTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IdentityTreeStoreTransactorRaw struct {
	Contract *IdentityTreeStoreTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIdentityTreeStore creates a new instance of IdentityTreeStore, bound to a specific deployed contract.
func NewIdentityTreeStore(address common.Address, backend bind.ContractBackend) (*IdentityTreeStore, error) {
	contract, err := bindIdentityTreeStore(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IdentityTreeStore{IdentityTreeStoreCaller: IdentityTreeStoreCaller{contract: contract}, IdentityTreeStoreTransactor: IdentityTreeStoreTransactor{contract: contract}, IdentityTreeStoreFilterer: IdentityTreeStoreFilterer{contract: contract}}, nil
}

// NewIdentityTreeStoreCaller creates a new read-only instance of IdentityTreeStore, bound to a specific deployed contract.
func NewIdentityTreeStoreCaller(address common.Address, caller bind.ContractCaller) (*IdentityTreeStoreCaller, error) {
	contract, err := bindIdentityTreeStore(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IdentityTreeStoreCaller{contract: contract}, nil
}

// NewIdentityTreeStoreTransactor creates a new write-only instance of IdentityTreeStore, bound to a specific deployed contract.
func NewIdentityTreeStoreTransactor(address common.Address, transactor bind.ContractTransactor) (*IdentityTreeStoreTransactor, error) {
	contract, err := bindIdentityTreeStore(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IdentityTreeStoreTransactor{contract: contract}, nil
}

// NewIdentityTreeStoreFilterer creates a new log filterer instance of IdentityTreeStore, bound to a specific deployed contract.
func NewIdentityTreeStoreFilterer(address common.Address, filterer bind.ContractFilterer) (*IdentityTreeStoreFilterer, error) {
	contract, err := bindIdentityTreeStore(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IdentityTreeStoreFilterer{contract: contract}, nil
}

// bindIdentityTreeStore binds a generic wrapper to an already deployed contract.
func bindIdentityTreeStore(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IdentityTreeStoreMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IdentityTreeStore *IdentityTreeStoreRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IdentityTreeStore.Contract.IdentityTreeStoreCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IdentityTreeStore *IdentityTreeStoreRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IdentityTreeStore.Contract.IdentityTreeStoreTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IdentityTreeStore *IdentityTreeStoreRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IdentityTreeStore.Contract.IdentityTreeStoreTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IdentityTreeStore *IdentityTreeStoreCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IdentityTreeStore.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IdentityTreeStore *IdentityTreeStoreTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IdentityTreeStore.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IdentityTreeStore *IdentityTreeStoreTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IdentityTreeStore.Contract.contract.Transact(opts, method, params...)
}

// HashFunc is a free data retrieval call binding the contract method 0x7f68025c.
//
// Solidity: function _hashFunc(uint256[] preimage) pure returns(uint256)
func (_IdentityTreeStore *IdentityTreeStoreCaller) HashFunc(opts *bind.CallOpts, preimage []*big.Int) (*big.Int, error) {
	var out []interface{}
	err := _IdentityTreeStore.contract.Call(opts, &out, "_hashFunc", preimage)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// HashFunc is a free data retrieval call binding the contract method 0x7f68025c.
//
// Solidity: function _hashFunc(uint256[] preimage) pure returns(uint256)
func (_IdentityTreeStore *IdentityTreeStoreSession) HashFunc(preimage []*big.Int) (*big.Int, error) {
	return _IdentityTreeStore.Contract.HashFunc(&_IdentityTreeStore.CallOpts, preimage)
}

// HashFunc is a free data retrieval call binding the contract method 0x7f68025c.
//
// Solidity: function _hashFunc(uint256[] preimage) pure returns(uint256)
func (_IdentityTreeStore *IdentityTreeStoreCallerSession) HashFunc(preimage []*big.Int) (*big.Int, error) {
	return _IdentityTreeStore.Contract.HashFunc(&_IdentityTreeStore.CallOpts, preimage)
}

// GetNode is a free data retrieval call binding the contract method 0x4f0f4aa9.
//
// Solidity: function getNode(uint256 id) view returns(uint256[])
func (_IdentityTreeStore *IdentityTreeStoreCaller) GetNode(opts *bind.CallOpts, id *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _IdentityTreeStore.contract.Call(opts, &out, "getNode", id)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetNode is a free data retrieval call binding the contract method 0x4f0f4aa9.
//
// Solidity: function getNode(uint256 id) view returns(uint256[])
func (_IdentityTreeStore *IdentityTreeStoreSession) GetNode(id *big.Int) ([]*big.Int, error) {
	return _IdentityTreeStore.Contract.GetNode(&_IdentityTreeStore.CallOpts, id)
}

// GetNode is a free data retrieval call binding the contract method 0x4f0f4aa9.
//
// Solidity: function getNode(uint256 id) view returns(uint256[])
func (_IdentityTreeStore *IdentityTreeStoreCallerSession) GetNode(id *big.Int) ([]*big.Int, error) {
	return _IdentityTreeStore.Contract.GetNode(&_IdentityTreeStore.CallOpts, id)
}

// GetRevocationStatus is a free data retrieval call binding the contract method 0x110c96a7.
//
// Solidity: function getRevocationStatus(uint256 id, uint64 nonce) view returns(((uint256,uint256,uint256,uint256),(uint256,bool,uint256[],uint256,uint256,bool,uint256,uint256)))
func (_IdentityTreeStore *IdentityTreeStoreCaller) GetRevocationStatus(opts *bind.CallOpts, id *big.Int, nonce uint64) (IOnchainCredentialStatusResolverCredentialStatus, error) {
	var out []interface{}
	err := _IdentityTreeStore.contract.Call(opts, &out, "getRevocationStatus", id, nonce)

	if err != nil {
		return *new(IOnchainCredentialStatusResolverCredentialStatus), err
	}

	out0 := *abi.ConvertType(out[0], new(IOnchainCredentialStatusResolverCredentialStatus)).(*IOnchainCredentialStatusResolverCredentialStatus)

	return out0, err

}

// GetRevocationStatus is a free data retrieval call binding the contract method 0x110c96a7.
//
// Solidity: function getRevocationStatus(uint256 id, uint64 nonce) view returns(((uint256,uint256,uint256,uint256),(uint256,bool,uint256[],uint256,uint256,bool,uint256,uint256)))
func (_IdentityTreeStore *IdentityTreeStoreSession) GetRevocationStatus(id *big.Int, nonce uint64) (IOnchainCredentialStatusResolverCredentialStatus, error) {
	return _IdentityTreeStore.Contract.GetRevocationStatus(&_IdentityTreeStore.CallOpts, id, nonce)
}

// GetRevocationStatus is a free data retrieval call binding the contract method 0x110c96a7.
//
// Solidity: function getRevocationStatus(uint256 id, uint64 nonce) view returns(((uint256,uint256,uint256,uint256),(uint256,bool,uint256[],uint256,uint256,bool,uint256,uint256)))
func (_IdentityTreeStore *IdentityTreeStoreCallerSession) GetRevocationStatus(id *big.Int, nonce uint64) (IOnchainCredentialStatusResolverCredentialStatus, error) {
	return _IdentityTreeStore.Contract.GetRevocationStatus(&_IdentityTreeStore.CallOpts, id, nonce)
}

// GetRevocationStatusByIdAndState is a free data retrieval call binding the contract method 0xaad72921.
//
// Solidity: function getRevocationStatusByIdAndState(uint256 id, uint256 state, uint64 nonce) view returns(((uint256,uint256,uint256,uint256),(uint256,bool,uint256[],uint256,uint256,bool,uint256,uint256)))
func (_IdentityTreeStore *IdentityTreeStoreCaller) GetRevocationStatusByIdAndState(opts *bind.CallOpts, id *big.Int, state *big.Int, nonce uint64) (IOnchainCredentialStatusResolverCredentialStatus, error) {
	var out []interface{}
	err := _IdentityTreeStore.contract.Call(opts, &out, "getRevocationStatusByIdAndState", id, state, nonce)

	if err != nil {
		return *new(IOnchainCredentialStatusResolverCredentialStatus), err
	}

	out0 := *abi.ConvertType(out[0], new(IOnchainCredentialStatusResolverCredentialStatus)).(*IOnchainCredentialStatusResolverCredentialStatus)

	return out0, err

}

// GetRevocationStatusByIdAndState is a free data retrieval call binding the contract method 0xaad72921.
//
// Solidity: function getRevocationStatusByIdAndState(uint256 id, uint256 state, uint64 nonce) view returns(((uint256,uint256,uint256,uint256),(uint256,bool,uint256[],uint256,uint256,bool,uint256,uint256)))
func (_IdentityTreeStore *IdentityTreeStoreSession) GetRevocationStatusByIdAndState(id *big.Int, state *big.Int, nonce uint64) (IOnchainCredentialStatusResolverCredentialStatus, error) {
	return _IdentityTreeStore.Contract.GetRevocationStatusByIdAndState(&_IdentityTreeStore.CallOpts, id, state, nonce)
}

// GetRevocationStatusByIdAndState is a free data retrieval call binding the contract method 0xaad72921.
//
// Solidity: function getRevocationStatusByIdAndState(uint256 id, uint256 state, uint64 nonce) view returns(((uint256,uint256,uint256,uint256),(uint256,bool,uint256[],uint256,uint256,bool,uint256,uint256)))
func (_IdentityTreeStore *IdentityTreeStoreCallerSession) GetRevocationStatusByIdAndState(id *big.Int, state *big.Int, nonce uint64) (IOnchainCredentialStatusResolverCredentialStatus, error) {
	return _IdentityTreeStore.Contract.GetRevocationStatusByIdAndState(&_IdentityTreeStore.CallOpts, id, state, nonce)
}

// SaveNodes is a paid mutator transaction binding the contract method 0x879957f7.
//
// Solidity: function saveNodes(uint256[][] preimage) returns()
func (_IdentityTreeStore *IdentityTreeStoreTransactor) SaveNodes(opts *bind.TransactOpts, preimage [][]*big.Int) (*types.Transaction, error) {
	return _IdentityTreeStore.contract.Transact(opts, "saveNodes", preimage)
}

// SaveNodes is a paid mutator transaction binding the contract method 0x879957f7.
//
// Solidity: function saveNodes(uint256[][] preimage) returns()
func (_IdentityTreeStore *IdentityTreeStoreSession) SaveNodes(preimage [][]*big.Int) (*types.Transaction, error) {
	return _IdentityTreeStore.Contract.SaveNodes(&_IdentityTreeStore.TransactOpts, preimage)
}

// SaveNodes is a paid mutator transaction binding the contract method 0x879957f7.
//
// Solidity: function saveNodes(uint256[][] preimage) returns()
func (_IdentityTreeStore *IdentityTreeStoreTransactorSession) SaveNodes(preimage [][]*big.Int) (*types.Transaction, error) {
	return _IdentityTreeStore.Contract.SaveNodes(&_IdentityTreeStore.TransactOpts, preimage)
}
