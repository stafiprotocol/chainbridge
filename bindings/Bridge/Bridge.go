// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package Bridge

import (
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
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// BridgeProposal is an auto generated low-level Go binding around an user-defined struct.
type BridgeProposal struct {
	Status        uint8
	YesVotes      *big.Int
	YesVotesTotal uint8
	ProposedBlock *big.Int
}

// BridgeABI is the input ABI used to generate the binding from.
const BridgeABI = "[{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"chainID\",\"type\":\"uint8\"},{\"internalType\":\"address[]\",\"name\":\"initialRelayers\",\"type\":\"address[]\"},{\"internalType\":\"uint256\",\"name\":\"initialRelayerThreshold\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expiry\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint8\",\"name\":\"destinationChainID\",\"type\":\"uint8\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"resourceID\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"depositNonce\",\"type\":\"uint64\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint8\",\"name\":\"originChainID\",\"type\":\"uint8\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"depositNonce\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"enumBridge.ProposalStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"dataHash\",\"type\":\"bytes32\"}],\"name\":\"ProposalEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint8\",\"name\":\"originChainID\",\"type\":\"uint8\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"depositNonce\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"enumBridge.ProposalStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"dataHash\",\"type\":\"bytes32\"}],\"name\":\"ProposalVote\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"relayer\",\"type\":\"address\"}],\"name\":\"RelayerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"relayer\",\"type\":\"address\"}],\"name\":\"RelayerRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"newThreshold\",\"type\":\"uint256\"}],\"name\":\"RelayerThresholdChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_RELAYERS\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"RELAYER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_chainID\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"name\":\"_depositCounts\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_expiry\",\"outputs\":[{\"internalType\":\"uint40\",\"name\":\"\",\"type\":\"uint40\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_fee\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"\",\"type\":\"uint128\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint72\",\"name\":\"destNonce\",\"type\":\"uint72\"},{\"internalType\":\"bytes32\",\"name\":\"dataHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"relayer\",\"type\":\"address\"}],\"name\":\"_hasVotedOnProposal\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"name\":\"_oldDepositNonce\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint72\",\"name\":\"\",\"type\":\"uint72\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"_proposals\",\"outputs\":[{\"internalType\":\"enumBridge.ProposalStatus\",\"name\":\"_status\",\"type\":\"uint8\"},{\"internalType\":\"uint200\",\"name\":\"_yesVotes\",\"type\":\"uint200\"},{\"internalType\":\"uint8\",\"name\":\"_yesVotesTotal\",\"type\":\"uint8\"},{\"internalType\":\"uint40\",\"name\":\"_proposedBlock\",\"type\":\"uint40\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_relayerThreshold\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"_resourceIDToHandlerAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_totalRelayers\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"relayerAddress\",\"type\":\"address\"}],\"name\":\"adminAddRelayer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newFee\",\"type\":\"uint256\"}],\"name\":\"adminChangeFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newThreshold\",\"type\":\"uint256\"}],\"name\":\"adminChangeRelayerThreshold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"adminPauseTransfers\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"relayerAddress\",\"type\":\"address\"}],\"name\":\"adminRemoveRelayer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"chainId\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"depositCount\",\"type\":\"uint64\"}],\"name\":\"adminSetDepositCount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"handlerAddress\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"resourceID\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"bytes4\",\"name\":\"depositFunctionSig\",\"type\":\"bytes4\"},{\"internalType\":\"uint256\",\"name\":\"depositFunctionDepositerOffset\",\"type\":\"uint256\"},{\"internalType\":\"bytes4\",\"name\":\"executeFunctionSig\",\"type\":\"bytes4\"}],\"name\":\"adminSetGenericResource\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"chainId\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"oldDepositNonce\",\"type\":\"uint64\"}],\"name\":\"adminSetOldDepositNonce\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"handlerAddress\",\"type\":\"address\"},{\"internalType\":\"bytes32[]\",\"name\":\"resourceIDs\",\"type\":\"bytes32[]\"},{\"internalType\":\"address[]\",\"name\":\"tokenContractAddresses\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"burnablTokenContractAddresses\",\"type\":\"address[]\"}],\"name\":\"adminSetResourceAndBurnable\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"adminUnpauseTransfers\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"handlerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountOrTokenID\",\"type\":\"uint256\"}],\"name\":\"adminWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"chainID\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"depositNonce\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"dataHash\",\"type\":\"bytes32\"}],\"name\":\"cancelProposal\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"destinationChainID\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"resourceID\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"originChainID\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"depositNonce\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"dataHash\",\"type\":\"bytes32\"}],\"name\":\"getProposal\",\"outputs\":[{\"components\":[{\"internalType\":\"enumBridge.ProposalStatus\",\"name\":\"_status\",\"type\":\"uint8\"},{\"internalType\":\"uint200\",\"name\":\"_yesVotes\",\"type\":\"uint200\"},{\"internalType\":\"uint8\",\"name\":\"_yesVotesTotal\",\"type\":\"uint8\"},{\"internalType\":\"uint40\",\"name\":\"_proposedBlock\",\"type\":\"uint40\"}],\"internalType\":\"structBridge.Proposal\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"getRoleMember\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleMemberCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"getRoleMemberIndex\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"relayer\",\"type\":\"address\"}],\"name\":\"isRelayer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"renounceAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable[]\",\"name\":\"addrs\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"name\":\"transferFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"chainID\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"depositNonce\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"resourceID\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"voteProposal\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// Bridge is an auto generated Go binding around an Ethereum contract.
type Bridge struct {
	BridgeCaller     // Read-only binding to the contract
	BridgeTransactor // Write-only binding to the contract
	BridgeFilterer   // Log filterer for contract events
}

// BridgeCaller is an auto generated read-only Go binding around an Ethereum contract.
type BridgeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BridgeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BridgeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BridgeSession struct {
	Contract     *Bridge           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BridgeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BridgeCallerSession struct {
	Contract *BridgeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// BridgeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BridgeTransactorSession struct {
	Contract     *BridgeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BridgeRaw is an auto generated low-level Go binding around an Ethereum contract.
type BridgeRaw struct {
	Contract *Bridge // Generic contract binding to access the raw methods on
}

// BridgeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BridgeCallerRaw struct {
	Contract *BridgeCaller // Generic read-only contract binding to access the raw methods on
}

// BridgeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BridgeTransactorRaw struct {
	Contract *BridgeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBridge creates a new instance of Bridge, bound to a specific deployed contract.
func NewBridge(address common.Address, backend bind.ContractBackend) (*Bridge, error) {
	contract, err := bindBridge(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Bridge{BridgeCaller: BridgeCaller{contract: contract}, BridgeTransactor: BridgeTransactor{contract: contract}, BridgeFilterer: BridgeFilterer{contract: contract}}, nil
}

// NewBridgeCaller creates a new read-only instance of Bridge, bound to a specific deployed contract.
func NewBridgeCaller(address common.Address, caller bind.ContractCaller) (*BridgeCaller, error) {
	contract, err := bindBridge(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BridgeCaller{contract: contract}, nil
}

// NewBridgeTransactor creates a new write-only instance of Bridge, bound to a specific deployed contract.
func NewBridgeTransactor(address common.Address, transactor bind.ContractTransactor) (*BridgeTransactor, error) {
	contract, err := bindBridge(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BridgeTransactor{contract: contract}, nil
}

// NewBridgeFilterer creates a new log filterer instance of Bridge, bound to a specific deployed contract.
func NewBridgeFilterer(address common.Address, filterer bind.ContractFilterer) (*BridgeFilterer, error) {
	contract, err := bindBridge(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BridgeFilterer{contract: contract}, nil
}

// bindBridge binds a generic wrapper to an already deployed contract.
func bindBridge(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BridgeABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Bridge *BridgeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Bridge.Contract.BridgeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Bridge *BridgeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bridge.Contract.BridgeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Bridge *BridgeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Bridge.Contract.BridgeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Bridge *BridgeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Bridge.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Bridge *BridgeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bridge.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Bridge *BridgeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Bridge.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Bridge *BridgeCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Bridge *BridgeSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Bridge.Contract.DEFAULTADMINROLE(&_Bridge.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Bridge *BridgeCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Bridge.Contract.DEFAULTADMINROLE(&_Bridge.CallOpts)
}

// MAXRELAYERS is a free data retrieval call binding the contract method 0x9debb3bd.
//
// Solidity: function MAX_RELAYERS() view returns(uint256)
func (_Bridge *BridgeCaller) MAXRELAYERS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "MAX_RELAYERS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXRELAYERS is a free data retrieval call binding the contract method 0x9debb3bd.
//
// Solidity: function MAX_RELAYERS() view returns(uint256)
func (_Bridge *BridgeSession) MAXRELAYERS() (*big.Int, error) {
	return _Bridge.Contract.MAXRELAYERS(&_Bridge.CallOpts)
}

// MAXRELAYERS is a free data retrieval call binding the contract method 0x9debb3bd.
//
// Solidity: function MAX_RELAYERS() view returns(uint256)
func (_Bridge *BridgeCallerSession) MAXRELAYERS() (*big.Int, error) {
	return _Bridge.Contract.MAXRELAYERS(&_Bridge.CallOpts)
}

// RELAYERROLE is a free data retrieval call binding the contract method 0x926d7d7f.
//
// Solidity: function RELAYER_ROLE() view returns(bytes32)
func (_Bridge *BridgeCaller) RELAYERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "RELAYER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// RELAYERROLE is a free data retrieval call binding the contract method 0x926d7d7f.
//
// Solidity: function RELAYER_ROLE() view returns(bytes32)
func (_Bridge *BridgeSession) RELAYERROLE() ([32]byte, error) {
	return _Bridge.Contract.RELAYERROLE(&_Bridge.CallOpts)
}

// RELAYERROLE is a free data retrieval call binding the contract method 0x926d7d7f.
//
// Solidity: function RELAYER_ROLE() view returns(bytes32)
func (_Bridge *BridgeCallerSession) RELAYERROLE() ([32]byte, error) {
	return _Bridge.Contract.RELAYERROLE(&_Bridge.CallOpts)
}

// ChainID is a free data retrieval call binding the contract method 0xbeab7131.
//
// Solidity: function _chainID() view returns(uint8)
func (_Bridge *BridgeCaller) ChainID(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "_chainID")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// ChainID is a free data retrieval call binding the contract method 0xbeab7131.
//
// Solidity: function _chainID() view returns(uint8)
func (_Bridge *BridgeSession) ChainID() (uint8, error) {
	return _Bridge.Contract.ChainID(&_Bridge.CallOpts)
}

// ChainID is a free data retrieval call binding the contract method 0xbeab7131.
//
// Solidity: function _chainID() view returns(uint8)
func (_Bridge *BridgeCallerSession) ChainID() (uint8, error) {
	return _Bridge.Contract.ChainID(&_Bridge.CallOpts)
}

// DepositCounts is a free data retrieval call binding the contract method 0x4b0b919d.
//
// Solidity: function _depositCounts(uint8 ) view returns(uint64)
func (_Bridge *BridgeCaller) DepositCounts(opts *bind.CallOpts, arg0 uint8) (uint64, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "_depositCounts", arg0)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// DepositCounts is a free data retrieval call binding the contract method 0x4b0b919d.
//
// Solidity: function _depositCounts(uint8 ) view returns(uint64)
func (_Bridge *BridgeSession) DepositCounts(arg0 uint8) (uint64, error) {
	return _Bridge.Contract.DepositCounts(&_Bridge.CallOpts, arg0)
}

// DepositCounts is a free data retrieval call binding the contract method 0x4b0b919d.
//
// Solidity: function _depositCounts(uint8 ) view returns(uint64)
func (_Bridge *BridgeCallerSession) DepositCounts(arg0 uint8) (uint64, error) {
	return _Bridge.Contract.DepositCounts(&_Bridge.CallOpts, arg0)
}

// Expiry is a free data retrieval call binding the contract method 0xc5ec8970.
//
// Solidity: function _expiry() view returns(uint40)
func (_Bridge *BridgeCaller) Expiry(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "_expiry")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Expiry is a free data retrieval call binding the contract method 0xc5ec8970.
//
// Solidity: function _expiry() view returns(uint40)
func (_Bridge *BridgeSession) Expiry() (*big.Int, error) {
	return _Bridge.Contract.Expiry(&_Bridge.CallOpts)
}

// Expiry is a free data retrieval call binding the contract method 0xc5ec8970.
//
// Solidity: function _expiry() view returns(uint40)
func (_Bridge *BridgeCallerSession) Expiry() (*big.Int, error) {
	return _Bridge.Contract.Expiry(&_Bridge.CallOpts)
}

// Fee is a free data retrieval call binding the contract method 0xc5b37c22.
//
// Solidity: function _fee() view returns(uint128)
func (_Bridge *BridgeCaller) Fee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "_fee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Fee is a free data retrieval call binding the contract method 0xc5b37c22.
//
// Solidity: function _fee() view returns(uint128)
func (_Bridge *BridgeSession) Fee() (*big.Int, error) {
	return _Bridge.Contract.Fee(&_Bridge.CallOpts)
}

// Fee is a free data retrieval call binding the contract method 0xc5b37c22.
//
// Solidity: function _fee() view returns(uint128)
func (_Bridge *BridgeCallerSession) Fee() (*big.Int, error) {
	return _Bridge.Contract.Fee(&_Bridge.CallOpts)
}

// HasVotedOnProposal is a free data retrieval call binding the contract method 0x7febe63f.
//
// Solidity: function _hasVotedOnProposal(uint72 destNonce, bytes32 dataHash, address relayer) view returns(bool)
func (_Bridge *BridgeCaller) HasVotedOnProposal(opts *bind.CallOpts, destNonce *big.Int, dataHash [32]byte, relayer common.Address) (bool, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "_hasVotedOnProposal", destNonce, dataHash, relayer)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasVotedOnProposal is a free data retrieval call binding the contract method 0x7febe63f.
//
// Solidity: function _hasVotedOnProposal(uint72 destNonce, bytes32 dataHash, address relayer) view returns(bool)
func (_Bridge *BridgeSession) HasVotedOnProposal(destNonce *big.Int, dataHash [32]byte, relayer common.Address) (bool, error) {
	return _Bridge.Contract.HasVotedOnProposal(&_Bridge.CallOpts, destNonce, dataHash, relayer)
}

// HasVotedOnProposal is a free data retrieval call binding the contract method 0x7febe63f.
//
// Solidity: function _hasVotedOnProposal(uint72 destNonce, bytes32 dataHash, address relayer) view returns(bool)
func (_Bridge *BridgeCallerSession) HasVotedOnProposal(destNonce *big.Int, dataHash [32]byte, relayer common.Address) (bool, error) {
	return _Bridge.Contract.HasVotedOnProposal(&_Bridge.CallOpts, destNonce, dataHash, relayer)
}

// OldDepositNonce is a free data retrieval call binding the contract method 0x263dcade.
//
// Solidity: function _oldDepositNonce(uint8 ) view returns(uint64)
func (_Bridge *BridgeCaller) OldDepositNonce(opts *bind.CallOpts, arg0 uint8) (uint64, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "_oldDepositNonce", arg0)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// OldDepositNonce is a free data retrieval call binding the contract method 0x263dcade.
//
// Solidity: function _oldDepositNonce(uint8 ) view returns(uint64)
func (_Bridge *BridgeSession) OldDepositNonce(arg0 uint8) (uint64, error) {
	return _Bridge.Contract.OldDepositNonce(&_Bridge.CallOpts, arg0)
}

// OldDepositNonce is a free data retrieval call binding the contract method 0x263dcade.
//
// Solidity: function _oldDepositNonce(uint8 ) view returns(uint64)
func (_Bridge *BridgeCallerSession) OldDepositNonce(arg0 uint8) (uint64, error) {
	return _Bridge.Contract.OldDepositNonce(&_Bridge.CallOpts, arg0)
}

// Proposals is a free data retrieval call binding the contract method 0x50598719.
//
// Solidity: function _proposals(uint72 , bytes32 ) view returns(uint8 _status, uint200 _yesVotes, uint8 _yesVotesTotal, uint40 _proposedBlock)
func (_Bridge *BridgeCaller) Proposals(opts *bind.CallOpts, arg0 *big.Int, arg1 [32]byte) (struct {
	Status        uint8
	YesVotes      *big.Int
	YesVotesTotal uint8
	ProposedBlock *big.Int
}, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "_proposals", arg0, arg1)

	outstruct := new(struct {
		Status        uint8
		YesVotes      *big.Int
		YesVotesTotal uint8
		ProposedBlock *big.Int
	})

	outstruct.Status = out[0].(uint8)
	outstruct.YesVotes = out[1].(*big.Int)
	outstruct.YesVotesTotal = out[2].(uint8)
	outstruct.ProposedBlock = out[3].(*big.Int)

	return *outstruct, err

}

// Proposals is a free data retrieval call binding the contract method 0x50598719.
//
// Solidity: function _proposals(uint72 , bytes32 ) view returns(uint8 _status, uint200 _yesVotes, uint8 _yesVotesTotal, uint40 _proposedBlock)
func (_Bridge *BridgeSession) Proposals(arg0 *big.Int, arg1 [32]byte) (struct {
	Status        uint8
	YesVotes      *big.Int
	YesVotesTotal uint8
	ProposedBlock *big.Int
}, error) {
	return _Bridge.Contract.Proposals(&_Bridge.CallOpts, arg0, arg1)
}

// Proposals is a free data retrieval call binding the contract method 0x50598719.
//
// Solidity: function _proposals(uint72 , bytes32 ) view returns(uint8 _status, uint200 _yesVotes, uint8 _yesVotesTotal, uint40 _proposedBlock)
func (_Bridge *BridgeCallerSession) Proposals(arg0 *big.Int, arg1 [32]byte) (struct {
	Status        uint8
	YesVotes      *big.Int
	YesVotesTotal uint8
	ProposedBlock *big.Int
}, error) {
	return _Bridge.Contract.Proposals(&_Bridge.CallOpts, arg0, arg1)
}

// RelayerThreshold is a free data retrieval call binding the contract method 0xd7a9cd79.
//
// Solidity: function _relayerThreshold() view returns(uint8)
func (_Bridge *BridgeCaller) RelayerThreshold(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "_relayerThreshold")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// RelayerThreshold is a free data retrieval call binding the contract method 0xd7a9cd79.
//
// Solidity: function _relayerThreshold() view returns(uint8)
func (_Bridge *BridgeSession) RelayerThreshold() (uint8, error) {
	return _Bridge.Contract.RelayerThreshold(&_Bridge.CallOpts)
}

// RelayerThreshold is a free data retrieval call binding the contract method 0xd7a9cd79.
//
// Solidity: function _relayerThreshold() view returns(uint8)
func (_Bridge *BridgeCallerSession) RelayerThreshold() (uint8, error) {
	return _Bridge.Contract.RelayerThreshold(&_Bridge.CallOpts)
}

// ResourceIDToHandlerAddress is a free data retrieval call binding the contract method 0x84db809f.
//
// Solidity: function _resourceIDToHandlerAddress(bytes32 ) view returns(address)
func (_Bridge *BridgeCaller) ResourceIDToHandlerAddress(opts *bind.CallOpts, arg0 [32]byte) (common.Address, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "_resourceIDToHandlerAddress", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ResourceIDToHandlerAddress is a free data retrieval call binding the contract method 0x84db809f.
//
// Solidity: function _resourceIDToHandlerAddress(bytes32 ) view returns(address)
func (_Bridge *BridgeSession) ResourceIDToHandlerAddress(arg0 [32]byte) (common.Address, error) {
	return _Bridge.Contract.ResourceIDToHandlerAddress(&_Bridge.CallOpts, arg0)
}

// ResourceIDToHandlerAddress is a free data retrieval call binding the contract method 0x84db809f.
//
// Solidity: function _resourceIDToHandlerAddress(bytes32 ) view returns(address)
func (_Bridge *BridgeCallerSession) ResourceIDToHandlerAddress(arg0 [32]byte) (common.Address, error) {
	return _Bridge.Contract.ResourceIDToHandlerAddress(&_Bridge.CallOpts, arg0)
}

// TotalRelayers is a free data retrieval call binding the contract method 0x802aabe8.
//
// Solidity: function _totalRelayers() view returns(uint256)
func (_Bridge *BridgeCaller) TotalRelayers(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "_totalRelayers")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalRelayers is a free data retrieval call binding the contract method 0x802aabe8.
//
// Solidity: function _totalRelayers() view returns(uint256)
func (_Bridge *BridgeSession) TotalRelayers() (*big.Int, error) {
	return _Bridge.Contract.TotalRelayers(&_Bridge.CallOpts)
}

// TotalRelayers is a free data retrieval call binding the contract method 0x802aabe8.
//
// Solidity: function _totalRelayers() view returns(uint256)
func (_Bridge *BridgeCallerSession) TotalRelayers() (*big.Int, error) {
	return _Bridge.Contract.TotalRelayers(&_Bridge.CallOpts)
}

// GetProposal is a free data retrieval call binding the contract method 0xa9cf69fa.
//
// Solidity: function getProposal(uint8 originChainID, uint64 depositNonce, bytes32 dataHash) view returns((uint8,uint200,uint8,uint40))
func (_Bridge *BridgeCaller) GetProposal(opts *bind.CallOpts, originChainID uint8, depositNonce uint64, dataHash [32]byte) (BridgeProposal, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getProposal", originChainID, depositNonce, dataHash)

	if err != nil {
		return *new(BridgeProposal), err
	}

	out0 := *abi.ConvertType(out[0], new(BridgeProposal)).(*BridgeProposal)

	return out0, err

}

// GetProposal is a free data retrieval call binding the contract method 0xa9cf69fa.
//
// Solidity: function getProposal(uint8 originChainID, uint64 depositNonce, bytes32 dataHash) view returns((uint8,uint200,uint8,uint40))
func (_Bridge *BridgeSession) GetProposal(originChainID uint8, depositNonce uint64, dataHash [32]byte) (BridgeProposal, error) {
	return _Bridge.Contract.GetProposal(&_Bridge.CallOpts, originChainID, depositNonce, dataHash)
}

// GetProposal is a free data retrieval call binding the contract method 0xa9cf69fa.
//
// Solidity: function getProposal(uint8 originChainID, uint64 depositNonce, bytes32 dataHash) view returns((uint8,uint200,uint8,uint40))
func (_Bridge *BridgeCallerSession) GetProposal(originChainID uint8, depositNonce uint64, dataHash [32]byte) (BridgeProposal, error) {
	return _Bridge.Contract.GetProposal(&_Bridge.CallOpts, originChainID, depositNonce, dataHash)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Bridge *BridgeCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Bridge *BridgeSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Bridge.Contract.GetRoleAdmin(&_Bridge.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Bridge *BridgeCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Bridge.Contract.GetRoleAdmin(&_Bridge.CallOpts, role)
}

// GetRoleMember is a free data retrieval call binding the contract method 0x9010d07c.
//
// Solidity: function getRoleMember(bytes32 role, uint256 index) view returns(address)
func (_Bridge *BridgeCaller) GetRoleMember(opts *bind.CallOpts, role [32]byte, index *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getRoleMember", role, index)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetRoleMember is a free data retrieval call binding the contract method 0x9010d07c.
//
// Solidity: function getRoleMember(bytes32 role, uint256 index) view returns(address)
func (_Bridge *BridgeSession) GetRoleMember(role [32]byte, index *big.Int) (common.Address, error) {
	return _Bridge.Contract.GetRoleMember(&_Bridge.CallOpts, role, index)
}

// GetRoleMember is a free data retrieval call binding the contract method 0x9010d07c.
//
// Solidity: function getRoleMember(bytes32 role, uint256 index) view returns(address)
func (_Bridge *BridgeCallerSession) GetRoleMember(role [32]byte, index *big.Int) (common.Address, error) {
	return _Bridge.Contract.GetRoleMember(&_Bridge.CallOpts, role, index)
}

// GetRoleMemberCount is a free data retrieval call binding the contract method 0xca15c873.
//
// Solidity: function getRoleMemberCount(bytes32 role) view returns(uint256)
func (_Bridge *BridgeCaller) GetRoleMemberCount(opts *bind.CallOpts, role [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getRoleMemberCount", role)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRoleMemberCount is a free data retrieval call binding the contract method 0xca15c873.
//
// Solidity: function getRoleMemberCount(bytes32 role) view returns(uint256)
func (_Bridge *BridgeSession) GetRoleMemberCount(role [32]byte) (*big.Int, error) {
	return _Bridge.Contract.GetRoleMemberCount(&_Bridge.CallOpts, role)
}

// GetRoleMemberCount is a free data retrieval call binding the contract method 0xca15c873.
//
// Solidity: function getRoleMemberCount(bytes32 role) view returns(uint256)
func (_Bridge *BridgeCallerSession) GetRoleMemberCount(role [32]byte) (*big.Int, error) {
	return _Bridge.Contract.GetRoleMemberCount(&_Bridge.CallOpts, role)
}

// GetRoleMemberIndex is a free data retrieval call binding the contract method 0x4e0df3f6.
//
// Solidity: function getRoleMemberIndex(bytes32 role, address account) view returns(uint256)
func (_Bridge *BridgeCaller) GetRoleMemberIndex(opts *bind.CallOpts, role [32]byte, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getRoleMemberIndex", role, account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRoleMemberIndex is a free data retrieval call binding the contract method 0x4e0df3f6.
//
// Solidity: function getRoleMemberIndex(bytes32 role, address account) view returns(uint256)
func (_Bridge *BridgeSession) GetRoleMemberIndex(role [32]byte, account common.Address) (*big.Int, error) {
	return _Bridge.Contract.GetRoleMemberIndex(&_Bridge.CallOpts, role, account)
}

// GetRoleMemberIndex is a free data retrieval call binding the contract method 0x4e0df3f6.
//
// Solidity: function getRoleMemberIndex(bytes32 role, address account) view returns(uint256)
func (_Bridge *BridgeCallerSession) GetRoleMemberIndex(role [32]byte, account common.Address) (*big.Int, error) {
	return _Bridge.Contract.GetRoleMemberIndex(&_Bridge.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Bridge *BridgeCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Bridge *BridgeSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Bridge.Contract.HasRole(&_Bridge.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Bridge *BridgeCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Bridge.Contract.HasRole(&_Bridge.CallOpts, role, account)
}

// IsRelayer is a free data retrieval call binding the contract method 0x541d5548.
//
// Solidity: function isRelayer(address relayer) view returns(bool)
func (_Bridge *BridgeCaller) IsRelayer(opts *bind.CallOpts, relayer common.Address) (bool, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "isRelayer", relayer)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsRelayer is a free data retrieval call binding the contract method 0x541d5548.
//
// Solidity: function isRelayer(address relayer) view returns(bool)
func (_Bridge *BridgeSession) IsRelayer(relayer common.Address) (bool, error) {
	return _Bridge.Contract.IsRelayer(&_Bridge.CallOpts, relayer)
}

// IsRelayer is a free data retrieval call binding the contract method 0x541d5548.
//
// Solidity: function isRelayer(address relayer) view returns(bool)
func (_Bridge *BridgeCallerSession) IsRelayer(relayer common.Address) (bool, error) {
	return _Bridge.Contract.IsRelayer(&_Bridge.CallOpts, relayer)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Bridge *BridgeCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Bridge *BridgeSession) Paused() (bool, error) {
	return _Bridge.Contract.Paused(&_Bridge.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Bridge *BridgeCallerSession) Paused() (bool, error) {
	return _Bridge.Contract.Paused(&_Bridge.CallOpts)
}

// AdminAddRelayer is a paid mutator transaction binding the contract method 0xcdb0f73a.
//
// Solidity: function adminAddRelayer(address relayerAddress) returns()
func (_Bridge *BridgeTransactor) AdminAddRelayer(opts *bind.TransactOpts, relayerAddress common.Address) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "adminAddRelayer", relayerAddress)
}

// AdminAddRelayer is a paid mutator transaction binding the contract method 0xcdb0f73a.
//
// Solidity: function adminAddRelayer(address relayerAddress) returns()
func (_Bridge *BridgeSession) AdminAddRelayer(relayerAddress common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.AdminAddRelayer(&_Bridge.TransactOpts, relayerAddress)
}

// AdminAddRelayer is a paid mutator transaction binding the contract method 0xcdb0f73a.
//
// Solidity: function adminAddRelayer(address relayerAddress) returns()
func (_Bridge *BridgeTransactorSession) AdminAddRelayer(relayerAddress common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.AdminAddRelayer(&_Bridge.TransactOpts, relayerAddress)
}

// AdminChangeFee is a paid mutator transaction binding the contract method 0x91c404ac.
//
// Solidity: function adminChangeFee(uint256 newFee) returns()
func (_Bridge *BridgeTransactor) AdminChangeFee(opts *bind.TransactOpts, newFee *big.Int) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "adminChangeFee", newFee)
}

// AdminChangeFee is a paid mutator transaction binding the contract method 0x91c404ac.
//
// Solidity: function adminChangeFee(uint256 newFee) returns()
func (_Bridge *BridgeSession) AdminChangeFee(newFee *big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.AdminChangeFee(&_Bridge.TransactOpts, newFee)
}

// AdminChangeFee is a paid mutator transaction binding the contract method 0x91c404ac.
//
// Solidity: function adminChangeFee(uint256 newFee) returns()
func (_Bridge *BridgeTransactorSession) AdminChangeFee(newFee *big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.AdminChangeFee(&_Bridge.TransactOpts, newFee)
}

// AdminChangeRelayerThreshold is a paid mutator transaction binding the contract method 0x4e056005.
//
// Solidity: function adminChangeRelayerThreshold(uint256 newThreshold) returns()
func (_Bridge *BridgeTransactor) AdminChangeRelayerThreshold(opts *bind.TransactOpts, newThreshold *big.Int) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "adminChangeRelayerThreshold", newThreshold)
}

// AdminChangeRelayerThreshold is a paid mutator transaction binding the contract method 0x4e056005.
//
// Solidity: function adminChangeRelayerThreshold(uint256 newThreshold) returns()
func (_Bridge *BridgeSession) AdminChangeRelayerThreshold(newThreshold *big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.AdminChangeRelayerThreshold(&_Bridge.TransactOpts, newThreshold)
}

// AdminChangeRelayerThreshold is a paid mutator transaction binding the contract method 0x4e056005.
//
// Solidity: function adminChangeRelayerThreshold(uint256 newThreshold) returns()
func (_Bridge *BridgeTransactorSession) AdminChangeRelayerThreshold(newThreshold *big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.AdminChangeRelayerThreshold(&_Bridge.TransactOpts, newThreshold)
}

// AdminPauseTransfers is a paid mutator transaction binding the contract method 0x80ae1c28.
//
// Solidity: function adminPauseTransfers() returns()
func (_Bridge *BridgeTransactor) AdminPauseTransfers(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "adminPauseTransfers")
}

// AdminPauseTransfers is a paid mutator transaction binding the contract method 0x80ae1c28.
//
// Solidity: function adminPauseTransfers() returns()
func (_Bridge *BridgeSession) AdminPauseTransfers() (*types.Transaction, error) {
	return _Bridge.Contract.AdminPauseTransfers(&_Bridge.TransactOpts)
}

// AdminPauseTransfers is a paid mutator transaction binding the contract method 0x80ae1c28.
//
// Solidity: function adminPauseTransfers() returns()
func (_Bridge *BridgeTransactorSession) AdminPauseTransfers() (*types.Transaction, error) {
	return _Bridge.Contract.AdminPauseTransfers(&_Bridge.TransactOpts)
}

// AdminRemoveRelayer is a paid mutator transaction binding the contract method 0x9d82dd63.
//
// Solidity: function adminRemoveRelayer(address relayerAddress) returns()
func (_Bridge *BridgeTransactor) AdminRemoveRelayer(opts *bind.TransactOpts, relayerAddress common.Address) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "adminRemoveRelayer", relayerAddress)
}

// AdminRemoveRelayer is a paid mutator transaction binding the contract method 0x9d82dd63.
//
// Solidity: function adminRemoveRelayer(address relayerAddress) returns()
func (_Bridge *BridgeSession) AdminRemoveRelayer(relayerAddress common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.AdminRemoveRelayer(&_Bridge.TransactOpts, relayerAddress)
}

// AdminRemoveRelayer is a paid mutator transaction binding the contract method 0x9d82dd63.
//
// Solidity: function adminRemoveRelayer(address relayerAddress) returns()
func (_Bridge *BridgeTransactorSession) AdminRemoveRelayer(relayerAddress common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.AdminRemoveRelayer(&_Bridge.TransactOpts, relayerAddress)
}

// AdminSetDepositCount is a paid mutator transaction binding the contract method 0x1209998c.
//
// Solidity: function adminSetDepositCount(uint8 chainId, uint64 depositCount) returns()
func (_Bridge *BridgeTransactor) AdminSetDepositCount(opts *bind.TransactOpts, chainId uint8, depositCount uint64) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "adminSetDepositCount", chainId, depositCount)
}

// AdminSetDepositCount is a paid mutator transaction binding the contract method 0x1209998c.
//
// Solidity: function adminSetDepositCount(uint8 chainId, uint64 depositCount) returns()
func (_Bridge *BridgeSession) AdminSetDepositCount(chainId uint8, depositCount uint64) (*types.Transaction, error) {
	return _Bridge.Contract.AdminSetDepositCount(&_Bridge.TransactOpts, chainId, depositCount)
}

// AdminSetDepositCount is a paid mutator transaction binding the contract method 0x1209998c.
//
// Solidity: function adminSetDepositCount(uint8 chainId, uint64 depositCount) returns()
func (_Bridge *BridgeTransactorSession) AdminSetDepositCount(chainId uint8, depositCount uint64) (*types.Transaction, error) {
	return _Bridge.Contract.AdminSetDepositCount(&_Bridge.TransactOpts, chainId, depositCount)
}

// AdminSetGenericResource is a paid mutator transaction binding the contract method 0x5a1ad87c.
//
// Solidity: function adminSetGenericResource(address handlerAddress, bytes32 resourceID, address contractAddress, bytes4 depositFunctionSig, uint256 depositFunctionDepositerOffset, bytes4 executeFunctionSig) returns()
func (_Bridge *BridgeTransactor) AdminSetGenericResource(opts *bind.TransactOpts, handlerAddress common.Address, resourceID [32]byte, contractAddress common.Address, depositFunctionSig [4]byte, depositFunctionDepositerOffset *big.Int, executeFunctionSig [4]byte) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "adminSetGenericResource", handlerAddress, resourceID, contractAddress, depositFunctionSig, depositFunctionDepositerOffset, executeFunctionSig)
}

// AdminSetGenericResource is a paid mutator transaction binding the contract method 0x5a1ad87c.
//
// Solidity: function adminSetGenericResource(address handlerAddress, bytes32 resourceID, address contractAddress, bytes4 depositFunctionSig, uint256 depositFunctionDepositerOffset, bytes4 executeFunctionSig) returns()
func (_Bridge *BridgeSession) AdminSetGenericResource(handlerAddress common.Address, resourceID [32]byte, contractAddress common.Address, depositFunctionSig [4]byte, depositFunctionDepositerOffset *big.Int, executeFunctionSig [4]byte) (*types.Transaction, error) {
	return _Bridge.Contract.AdminSetGenericResource(&_Bridge.TransactOpts, handlerAddress, resourceID, contractAddress, depositFunctionSig, depositFunctionDepositerOffset, executeFunctionSig)
}

// AdminSetGenericResource is a paid mutator transaction binding the contract method 0x5a1ad87c.
//
// Solidity: function adminSetGenericResource(address handlerAddress, bytes32 resourceID, address contractAddress, bytes4 depositFunctionSig, uint256 depositFunctionDepositerOffset, bytes4 executeFunctionSig) returns()
func (_Bridge *BridgeTransactorSession) AdminSetGenericResource(handlerAddress common.Address, resourceID [32]byte, contractAddress common.Address, depositFunctionSig [4]byte, depositFunctionDepositerOffset *big.Int, executeFunctionSig [4]byte) (*types.Transaction, error) {
	return _Bridge.Contract.AdminSetGenericResource(&_Bridge.TransactOpts, handlerAddress, resourceID, contractAddress, depositFunctionSig, depositFunctionDepositerOffset, executeFunctionSig)
}

// AdminSetOldDepositNonce is a paid mutator transaction binding the contract method 0x8a412d67.
//
// Solidity: function adminSetOldDepositNonce(uint8 chainId, uint64 oldDepositNonce) returns()
func (_Bridge *BridgeTransactor) AdminSetOldDepositNonce(opts *bind.TransactOpts, chainId uint8, oldDepositNonce uint64) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "adminSetOldDepositNonce", chainId, oldDepositNonce)
}

// AdminSetOldDepositNonce is a paid mutator transaction binding the contract method 0x8a412d67.
//
// Solidity: function adminSetOldDepositNonce(uint8 chainId, uint64 oldDepositNonce) returns()
func (_Bridge *BridgeSession) AdminSetOldDepositNonce(chainId uint8, oldDepositNonce uint64) (*types.Transaction, error) {
	return _Bridge.Contract.AdminSetOldDepositNonce(&_Bridge.TransactOpts, chainId, oldDepositNonce)
}

// AdminSetOldDepositNonce is a paid mutator transaction binding the contract method 0x8a412d67.
//
// Solidity: function adminSetOldDepositNonce(uint8 chainId, uint64 oldDepositNonce) returns()
func (_Bridge *BridgeTransactorSession) AdminSetOldDepositNonce(chainId uint8, oldDepositNonce uint64) (*types.Transaction, error) {
	return _Bridge.Contract.AdminSetOldDepositNonce(&_Bridge.TransactOpts, chainId, oldDepositNonce)
}

// AdminSetResourceAndBurnable is a paid mutator transaction binding the contract method 0x0d724dbe.
//
// Solidity: function adminSetResourceAndBurnable(address handlerAddress, bytes32[] resourceIDs, address[] tokenContractAddresses, address[] burnablTokenContractAddresses) returns()
func (_Bridge *BridgeTransactor) AdminSetResourceAndBurnable(opts *bind.TransactOpts, handlerAddress common.Address, resourceIDs [][32]byte, tokenContractAddresses []common.Address, burnablTokenContractAddresses []common.Address) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "adminSetResourceAndBurnable", handlerAddress, resourceIDs, tokenContractAddresses, burnablTokenContractAddresses)
}

// AdminSetResourceAndBurnable is a paid mutator transaction binding the contract method 0x0d724dbe.
//
// Solidity: function adminSetResourceAndBurnable(address handlerAddress, bytes32[] resourceIDs, address[] tokenContractAddresses, address[] burnablTokenContractAddresses) returns()
func (_Bridge *BridgeSession) AdminSetResourceAndBurnable(handlerAddress common.Address, resourceIDs [][32]byte, tokenContractAddresses []common.Address, burnablTokenContractAddresses []common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.AdminSetResourceAndBurnable(&_Bridge.TransactOpts, handlerAddress, resourceIDs, tokenContractAddresses, burnablTokenContractAddresses)
}

// AdminSetResourceAndBurnable is a paid mutator transaction binding the contract method 0x0d724dbe.
//
// Solidity: function adminSetResourceAndBurnable(address handlerAddress, bytes32[] resourceIDs, address[] tokenContractAddresses, address[] burnablTokenContractAddresses) returns()
func (_Bridge *BridgeTransactorSession) AdminSetResourceAndBurnable(handlerAddress common.Address, resourceIDs [][32]byte, tokenContractAddresses []common.Address, burnablTokenContractAddresses []common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.AdminSetResourceAndBurnable(&_Bridge.TransactOpts, handlerAddress, resourceIDs, tokenContractAddresses, burnablTokenContractAddresses)
}

// AdminUnpauseTransfers is a paid mutator transaction binding the contract method 0xffaac0eb.
//
// Solidity: function adminUnpauseTransfers() returns()
func (_Bridge *BridgeTransactor) AdminUnpauseTransfers(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "adminUnpauseTransfers")
}

// AdminUnpauseTransfers is a paid mutator transaction binding the contract method 0xffaac0eb.
//
// Solidity: function adminUnpauseTransfers() returns()
func (_Bridge *BridgeSession) AdminUnpauseTransfers() (*types.Transaction, error) {
	return _Bridge.Contract.AdminUnpauseTransfers(&_Bridge.TransactOpts)
}

// AdminUnpauseTransfers is a paid mutator transaction binding the contract method 0xffaac0eb.
//
// Solidity: function adminUnpauseTransfers() returns()
func (_Bridge *BridgeTransactorSession) AdminUnpauseTransfers() (*types.Transaction, error) {
	return _Bridge.Contract.AdminUnpauseTransfers(&_Bridge.TransactOpts)
}

// AdminWithdraw is a paid mutator transaction binding the contract method 0x780cf004.
//
// Solidity: function adminWithdraw(address handlerAddress, address tokenAddress, address recipient, uint256 amountOrTokenID) returns()
func (_Bridge *BridgeTransactor) AdminWithdraw(opts *bind.TransactOpts, handlerAddress common.Address, tokenAddress common.Address, recipient common.Address, amountOrTokenID *big.Int) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "adminWithdraw", handlerAddress, tokenAddress, recipient, amountOrTokenID)
}

// AdminWithdraw is a paid mutator transaction binding the contract method 0x780cf004.
//
// Solidity: function adminWithdraw(address handlerAddress, address tokenAddress, address recipient, uint256 amountOrTokenID) returns()
func (_Bridge *BridgeSession) AdminWithdraw(handlerAddress common.Address, tokenAddress common.Address, recipient common.Address, amountOrTokenID *big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.AdminWithdraw(&_Bridge.TransactOpts, handlerAddress, tokenAddress, recipient, amountOrTokenID)
}

// AdminWithdraw is a paid mutator transaction binding the contract method 0x780cf004.
//
// Solidity: function adminWithdraw(address handlerAddress, address tokenAddress, address recipient, uint256 amountOrTokenID) returns()
func (_Bridge *BridgeTransactorSession) AdminWithdraw(handlerAddress common.Address, tokenAddress common.Address, recipient common.Address, amountOrTokenID *big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.AdminWithdraw(&_Bridge.TransactOpts, handlerAddress, tokenAddress, recipient, amountOrTokenID)
}

// CancelProposal is a paid mutator transaction binding the contract method 0x17f03ce5.
//
// Solidity: function cancelProposal(uint8 chainID, uint64 depositNonce, bytes32 dataHash) returns()
func (_Bridge *BridgeTransactor) CancelProposal(opts *bind.TransactOpts, chainID uint8, depositNonce uint64, dataHash [32]byte) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "cancelProposal", chainID, depositNonce, dataHash)
}

// CancelProposal is a paid mutator transaction binding the contract method 0x17f03ce5.
//
// Solidity: function cancelProposal(uint8 chainID, uint64 depositNonce, bytes32 dataHash) returns()
func (_Bridge *BridgeSession) CancelProposal(chainID uint8, depositNonce uint64, dataHash [32]byte) (*types.Transaction, error) {
	return _Bridge.Contract.CancelProposal(&_Bridge.TransactOpts, chainID, depositNonce, dataHash)
}

// CancelProposal is a paid mutator transaction binding the contract method 0x17f03ce5.
//
// Solidity: function cancelProposal(uint8 chainID, uint64 depositNonce, bytes32 dataHash) returns()
func (_Bridge *BridgeTransactorSession) CancelProposal(chainID uint8, depositNonce uint64, dataHash [32]byte) (*types.Transaction, error) {
	return _Bridge.Contract.CancelProposal(&_Bridge.TransactOpts, chainID, depositNonce, dataHash)
}

// Deposit is a paid mutator transaction binding the contract method 0x05e2ca17.
//
// Solidity: function deposit(uint8 destinationChainID, bytes32 resourceID, bytes data) payable returns()
func (_Bridge *BridgeTransactor) Deposit(opts *bind.TransactOpts, destinationChainID uint8, resourceID [32]byte, data []byte) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "deposit", destinationChainID, resourceID, data)
}

// Deposit is a paid mutator transaction binding the contract method 0x05e2ca17.
//
// Solidity: function deposit(uint8 destinationChainID, bytes32 resourceID, bytes data) payable returns()
func (_Bridge *BridgeSession) Deposit(destinationChainID uint8, resourceID [32]byte, data []byte) (*types.Transaction, error) {
	return _Bridge.Contract.Deposit(&_Bridge.TransactOpts, destinationChainID, resourceID, data)
}

// Deposit is a paid mutator transaction binding the contract method 0x05e2ca17.
//
// Solidity: function deposit(uint8 destinationChainID, bytes32 resourceID, bytes data) payable returns()
func (_Bridge *BridgeTransactorSession) Deposit(destinationChainID uint8, resourceID [32]byte, data []byte) (*types.Transaction, error) {
	return _Bridge.Contract.Deposit(&_Bridge.TransactOpts, destinationChainID, resourceID, data)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Bridge *BridgeTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Bridge *BridgeSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.GrantRole(&_Bridge.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Bridge *BridgeTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.GrantRole(&_Bridge.TransactOpts, role, account)
}

// RenounceAdmin is a paid mutator transaction binding the contract method 0x5e1fab0f.
//
// Solidity: function renounceAdmin(address newAdmin) returns()
func (_Bridge *BridgeTransactor) RenounceAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "renounceAdmin", newAdmin)
}

// RenounceAdmin is a paid mutator transaction binding the contract method 0x5e1fab0f.
//
// Solidity: function renounceAdmin(address newAdmin) returns()
func (_Bridge *BridgeSession) RenounceAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.RenounceAdmin(&_Bridge.TransactOpts, newAdmin)
}

// RenounceAdmin is a paid mutator transaction binding the contract method 0x5e1fab0f.
//
// Solidity: function renounceAdmin(address newAdmin) returns()
func (_Bridge *BridgeTransactorSession) RenounceAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.RenounceAdmin(&_Bridge.TransactOpts, newAdmin)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Bridge *BridgeTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Bridge *BridgeSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.RenounceRole(&_Bridge.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Bridge *BridgeTransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.RenounceRole(&_Bridge.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Bridge *BridgeTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Bridge *BridgeSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.RevokeRole(&_Bridge.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Bridge *BridgeTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.RevokeRole(&_Bridge.TransactOpts, role, account)
}

// TransferFunds is a paid mutator transaction binding the contract method 0x4603ae38.
//
// Solidity: function transferFunds(address[] addrs, uint256[] amounts) returns()
func (_Bridge *BridgeTransactor) TransferFunds(opts *bind.TransactOpts, addrs []common.Address, amounts []*big.Int) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "transferFunds", addrs, amounts)
}

// TransferFunds is a paid mutator transaction binding the contract method 0x4603ae38.
//
// Solidity: function transferFunds(address[] addrs, uint256[] amounts) returns()
func (_Bridge *BridgeSession) TransferFunds(addrs []common.Address, amounts []*big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.TransferFunds(&_Bridge.TransactOpts, addrs, amounts)
}

// TransferFunds is a paid mutator transaction binding the contract method 0x4603ae38.
//
// Solidity: function transferFunds(address[] addrs, uint256[] amounts) returns()
func (_Bridge *BridgeTransactorSession) TransferFunds(addrs []common.Address, amounts []*big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.TransferFunds(&_Bridge.TransactOpts, addrs, amounts)
}

// VoteProposal is a paid mutator transaction binding the contract method 0xc0331b3e.
//
// Solidity: function voteProposal(uint8 chainID, uint64 depositNonce, bytes32 resourceID, bytes data) returns()
func (_Bridge *BridgeTransactor) VoteProposal(opts *bind.TransactOpts, chainID uint8, depositNonce uint64, resourceID [32]byte, data []byte) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "voteProposal", chainID, depositNonce, resourceID, data)
}

// VoteProposal is a paid mutator transaction binding the contract method 0xc0331b3e.
//
// Solidity: function voteProposal(uint8 chainID, uint64 depositNonce, bytes32 resourceID, bytes data) returns()
func (_Bridge *BridgeSession) VoteProposal(chainID uint8, depositNonce uint64, resourceID [32]byte, data []byte) (*types.Transaction, error) {
	return _Bridge.Contract.VoteProposal(&_Bridge.TransactOpts, chainID, depositNonce, resourceID, data)
}

// VoteProposal is a paid mutator transaction binding the contract method 0xc0331b3e.
//
// Solidity: function voteProposal(uint8 chainID, uint64 depositNonce, bytes32 resourceID, bytes data) returns()
func (_Bridge *BridgeTransactorSession) VoteProposal(chainID uint8, depositNonce uint64, resourceID [32]byte, data []byte) (*types.Transaction, error) {
	return _Bridge.Contract.VoteProposal(&_Bridge.TransactOpts, chainID, depositNonce, resourceID, data)
}

// BridgeDepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the Bridge contract.
type BridgeDepositIterator struct {
	Event *BridgeDeposit // Event containing the contract specifics and raw log

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
func (it *BridgeDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeDeposit)
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
		it.Event = new(BridgeDeposit)
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
func (it *BridgeDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeDeposit represents a Deposit event raised by the Bridge contract.
type BridgeDeposit struct {
	DestinationChainID uint8
	ResourceID         [32]byte
	DepositNonce       uint64
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0xdbb69440df8433824a026ef190652f29929eb64b4d1d5d2a69be8afe3e6eaed8.
//
// Solidity: event Deposit(uint8 indexed destinationChainID, bytes32 indexed resourceID, uint64 indexed depositNonce)
func (_Bridge *BridgeFilterer) FilterDeposit(opts *bind.FilterOpts, destinationChainID []uint8, resourceID [][32]byte, depositNonce []uint64) (*BridgeDepositIterator, error) {

	var destinationChainIDRule []interface{}
	for _, destinationChainIDItem := range destinationChainID {
		destinationChainIDRule = append(destinationChainIDRule, destinationChainIDItem)
	}
	var resourceIDRule []interface{}
	for _, resourceIDItem := range resourceID {
		resourceIDRule = append(resourceIDRule, resourceIDItem)
	}
	var depositNonceRule []interface{}
	for _, depositNonceItem := range depositNonce {
		depositNonceRule = append(depositNonceRule, depositNonceItem)
	}

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "Deposit", destinationChainIDRule, resourceIDRule, depositNonceRule)
	if err != nil {
		return nil, err
	}
	return &BridgeDepositIterator{contract: _Bridge.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0xdbb69440df8433824a026ef190652f29929eb64b4d1d5d2a69be8afe3e6eaed8.
//
// Solidity: event Deposit(uint8 indexed destinationChainID, bytes32 indexed resourceID, uint64 indexed depositNonce)
func (_Bridge *BridgeFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *BridgeDeposit, destinationChainID []uint8, resourceID [][32]byte, depositNonce []uint64) (event.Subscription, error) {

	var destinationChainIDRule []interface{}
	for _, destinationChainIDItem := range destinationChainID {
		destinationChainIDRule = append(destinationChainIDRule, destinationChainIDItem)
	}
	var resourceIDRule []interface{}
	for _, resourceIDItem := range resourceID {
		resourceIDRule = append(resourceIDRule, resourceIDItem)
	}
	var depositNonceRule []interface{}
	for _, depositNonceItem := range depositNonce {
		depositNonceRule = append(depositNonceRule, depositNonceItem)
	}

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "Deposit", destinationChainIDRule, resourceIDRule, depositNonceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeDeposit)
				if err := _Bridge.contract.UnpackLog(event, "Deposit", log); err != nil {
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

// ParseDeposit is a log parse operation binding the contract event 0xdbb69440df8433824a026ef190652f29929eb64b4d1d5d2a69be8afe3e6eaed8.
//
// Solidity: event Deposit(uint8 indexed destinationChainID, bytes32 indexed resourceID, uint64 indexed depositNonce)
func (_Bridge *BridgeFilterer) ParseDeposit(log types.Log) (*BridgeDeposit, error) {
	event := new(BridgeDeposit)
	if err := _Bridge.contract.UnpackLog(event, "Deposit", log); err != nil {
		return nil, err
	}
	return event, nil
}

// BridgePausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the Bridge contract.
type BridgePausedIterator struct {
	Event *BridgePaused // Event containing the contract specifics and raw log

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
func (it *BridgePausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgePaused)
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
		it.Event = new(BridgePaused)
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
func (it *BridgePausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgePausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgePaused represents a Paused event raised by the Bridge contract.
type BridgePaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Bridge *BridgeFilterer) FilterPaused(opts *bind.FilterOpts) (*BridgePausedIterator, error) {

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &BridgePausedIterator{contract: _Bridge.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Bridge *BridgeFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *BridgePaused) (event.Subscription, error) {

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgePaused)
				if err := _Bridge.contract.UnpackLog(event, "Paused", log); err != nil {
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

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Bridge *BridgeFilterer) ParsePaused(log types.Log) (*BridgePaused, error) {
	event := new(BridgePaused)
	if err := _Bridge.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	return event, nil
}

// BridgeProposalEventIterator is returned from FilterProposalEvent and is used to iterate over the raw logs and unpacked data for ProposalEvent events raised by the Bridge contract.
type BridgeProposalEventIterator struct {
	Event *BridgeProposalEvent // Event containing the contract specifics and raw log

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
func (it *BridgeProposalEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeProposalEvent)
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
		it.Event = new(BridgeProposalEvent)
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
func (it *BridgeProposalEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeProposalEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeProposalEvent represents a ProposalEvent event raised by the Bridge contract.
type BridgeProposalEvent struct {
	OriginChainID uint8
	DepositNonce  uint64
	Status        uint8
	DataHash      [32]byte
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterProposalEvent is a free log retrieval operation binding the contract event 0x968626a768e76ba1363efe44e322a6c4900c5f084e0b45f35e294dfddaa9e0d5.
//
// Solidity: event ProposalEvent(uint8 indexed originChainID, uint64 indexed depositNonce, uint8 indexed status, bytes32 dataHash)
func (_Bridge *BridgeFilterer) FilterProposalEvent(opts *bind.FilterOpts, originChainID []uint8, depositNonce []uint64, status []uint8) (*BridgeProposalEventIterator, error) {

	var originChainIDRule []interface{}
	for _, originChainIDItem := range originChainID {
		originChainIDRule = append(originChainIDRule, originChainIDItem)
	}
	var depositNonceRule []interface{}
	for _, depositNonceItem := range depositNonce {
		depositNonceRule = append(depositNonceRule, depositNonceItem)
	}
	var statusRule []interface{}
	for _, statusItem := range status {
		statusRule = append(statusRule, statusItem)
	}

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "ProposalEvent", originChainIDRule, depositNonceRule, statusRule)
	if err != nil {
		return nil, err
	}
	return &BridgeProposalEventIterator{contract: _Bridge.contract, event: "ProposalEvent", logs: logs, sub: sub}, nil
}

// WatchProposalEvent is a free log subscription operation binding the contract event 0x968626a768e76ba1363efe44e322a6c4900c5f084e0b45f35e294dfddaa9e0d5.
//
// Solidity: event ProposalEvent(uint8 indexed originChainID, uint64 indexed depositNonce, uint8 indexed status, bytes32 dataHash)
func (_Bridge *BridgeFilterer) WatchProposalEvent(opts *bind.WatchOpts, sink chan<- *BridgeProposalEvent, originChainID []uint8, depositNonce []uint64, status []uint8) (event.Subscription, error) {

	var originChainIDRule []interface{}
	for _, originChainIDItem := range originChainID {
		originChainIDRule = append(originChainIDRule, originChainIDItem)
	}
	var depositNonceRule []interface{}
	for _, depositNonceItem := range depositNonce {
		depositNonceRule = append(depositNonceRule, depositNonceItem)
	}
	var statusRule []interface{}
	for _, statusItem := range status {
		statusRule = append(statusRule, statusItem)
	}

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "ProposalEvent", originChainIDRule, depositNonceRule, statusRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeProposalEvent)
				if err := _Bridge.contract.UnpackLog(event, "ProposalEvent", log); err != nil {
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

// ParseProposalEvent is a log parse operation binding the contract event 0x968626a768e76ba1363efe44e322a6c4900c5f084e0b45f35e294dfddaa9e0d5.
//
// Solidity: event ProposalEvent(uint8 indexed originChainID, uint64 indexed depositNonce, uint8 indexed status, bytes32 dataHash)
func (_Bridge *BridgeFilterer) ParseProposalEvent(log types.Log) (*BridgeProposalEvent, error) {
	event := new(BridgeProposalEvent)
	if err := _Bridge.contract.UnpackLog(event, "ProposalEvent", log); err != nil {
		return nil, err
	}
	return event, nil
}

// BridgeProposalVoteIterator is returned from FilterProposalVote and is used to iterate over the raw logs and unpacked data for ProposalVote events raised by the Bridge contract.
type BridgeProposalVoteIterator struct {
	Event *BridgeProposalVote // Event containing the contract specifics and raw log

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
func (it *BridgeProposalVoteIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeProposalVote)
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
		it.Event = new(BridgeProposalVote)
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
func (it *BridgeProposalVoteIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeProposalVoteIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeProposalVote represents a ProposalVote event raised by the Bridge contract.
type BridgeProposalVote struct {
	OriginChainID uint8
	DepositNonce  uint64
	Status        uint8
	DataHash      [32]byte
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterProposalVote is a free log retrieval operation binding the contract event 0x25f8daaa4635a7729927ba3f5b3d59cc3320aca7c32c9db4e7ca7b9574343640.
//
// Solidity: event ProposalVote(uint8 indexed originChainID, uint64 indexed depositNonce, uint8 indexed status, bytes32 dataHash)
func (_Bridge *BridgeFilterer) FilterProposalVote(opts *bind.FilterOpts, originChainID []uint8, depositNonce []uint64, status []uint8) (*BridgeProposalVoteIterator, error) {

	var originChainIDRule []interface{}
	for _, originChainIDItem := range originChainID {
		originChainIDRule = append(originChainIDRule, originChainIDItem)
	}
	var depositNonceRule []interface{}
	for _, depositNonceItem := range depositNonce {
		depositNonceRule = append(depositNonceRule, depositNonceItem)
	}
	var statusRule []interface{}
	for _, statusItem := range status {
		statusRule = append(statusRule, statusItem)
	}

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "ProposalVote", originChainIDRule, depositNonceRule, statusRule)
	if err != nil {
		return nil, err
	}
	return &BridgeProposalVoteIterator{contract: _Bridge.contract, event: "ProposalVote", logs: logs, sub: sub}, nil
}

// WatchProposalVote is a free log subscription operation binding the contract event 0x25f8daaa4635a7729927ba3f5b3d59cc3320aca7c32c9db4e7ca7b9574343640.
//
// Solidity: event ProposalVote(uint8 indexed originChainID, uint64 indexed depositNonce, uint8 indexed status, bytes32 dataHash)
func (_Bridge *BridgeFilterer) WatchProposalVote(opts *bind.WatchOpts, sink chan<- *BridgeProposalVote, originChainID []uint8, depositNonce []uint64, status []uint8) (event.Subscription, error) {

	var originChainIDRule []interface{}
	for _, originChainIDItem := range originChainID {
		originChainIDRule = append(originChainIDRule, originChainIDItem)
	}
	var depositNonceRule []interface{}
	for _, depositNonceItem := range depositNonce {
		depositNonceRule = append(depositNonceRule, depositNonceItem)
	}
	var statusRule []interface{}
	for _, statusItem := range status {
		statusRule = append(statusRule, statusItem)
	}

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "ProposalVote", originChainIDRule, depositNonceRule, statusRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeProposalVote)
				if err := _Bridge.contract.UnpackLog(event, "ProposalVote", log); err != nil {
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

// ParseProposalVote is a log parse operation binding the contract event 0x25f8daaa4635a7729927ba3f5b3d59cc3320aca7c32c9db4e7ca7b9574343640.
//
// Solidity: event ProposalVote(uint8 indexed originChainID, uint64 indexed depositNonce, uint8 indexed status, bytes32 dataHash)
func (_Bridge *BridgeFilterer) ParseProposalVote(log types.Log) (*BridgeProposalVote, error) {
	event := new(BridgeProposalVote)
	if err := _Bridge.contract.UnpackLog(event, "ProposalVote", log); err != nil {
		return nil, err
	}
	return event, nil
}

// BridgeRelayerAddedIterator is returned from FilterRelayerAdded and is used to iterate over the raw logs and unpacked data for RelayerAdded events raised by the Bridge contract.
type BridgeRelayerAddedIterator struct {
	Event *BridgeRelayerAdded // Event containing the contract specifics and raw log

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
func (it *BridgeRelayerAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeRelayerAdded)
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
		it.Event = new(BridgeRelayerAdded)
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
func (it *BridgeRelayerAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeRelayerAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeRelayerAdded represents a RelayerAdded event raised by the Bridge contract.
type BridgeRelayerAdded struct {
	Relayer common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRelayerAdded is a free log retrieval operation binding the contract event 0x03580ee9f53a62b7cb409a2cb56f9be87747dd15017afc5cef6eef321e4fb2c5.
//
// Solidity: event RelayerAdded(address indexed relayer)
func (_Bridge *BridgeFilterer) FilterRelayerAdded(opts *bind.FilterOpts, relayer []common.Address) (*BridgeRelayerAddedIterator, error) {

	var relayerRule []interface{}
	for _, relayerItem := range relayer {
		relayerRule = append(relayerRule, relayerItem)
	}

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "RelayerAdded", relayerRule)
	if err != nil {
		return nil, err
	}
	return &BridgeRelayerAddedIterator{contract: _Bridge.contract, event: "RelayerAdded", logs: logs, sub: sub}, nil
}

// WatchRelayerAdded is a free log subscription operation binding the contract event 0x03580ee9f53a62b7cb409a2cb56f9be87747dd15017afc5cef6eef321e4fb2c5.
//
// Solidity: event RelayerAdded(address indexed relayer)
func (_Bridge *BridgeFilterer) WatchRelayerAdded(opts *bind.WatchOpts, sink chan<- *BridgeRelayerAdded, relayer []common.Address) (event.Subscription, error) {

	var relayerRule []interface{}
	for _, relayerItem := range relayer {
		relayerRule = append(relayerRule, relayerItem)
	}

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "RelayerAdded", relayerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeRelayerAdded)
				if err := _Bridge.contract.UnpackLog(event, "RelayerAdded", log); err != nil {
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

// ParseRelayerAdded is a log parse operation binding the contract event 0x03580ee9f53a62b7cb409a2cb56f9be87747dd15017afc5cef6eef321e4fb2c5.
//
// Solidity: event RelayerAdded(address indexed relayer)
func (_Bridge *BridgeFilterer) ParseRelayerAdded(log types.Log) (*BridgeRelayerAdded, error) {
	event := new(BridgeRelayerAdded)
	if err := _Bridge.contract.UnpackLog(event, "RelayerAdded", log); err != nil {
		return nil, err
	}
	return event, nil
}

// BridgeRelayerRemovedIterator is returned from FilterRelayerRemoved and is used to iterate over the raw logs and unpacked data for RelayerRemoved events raised by the Bridge contract.
type BridgeRelayerRemovedIterator struct {
	Event *BridgeRelayerRemoved // Event containing the contract specifics and raw log

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
func (it *BridgeRelayerRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeRelayerRemoved)
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
		it.Event = new(BridgeRelayerRemoved)
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
func (it *BridgeRelayerRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeRelayerRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeRelayerRemoved represents a RelayerRemoved event raised by the Bridge contract.
type BridgeRelayerRemoved struct {
	Relayer common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRelayerRemoved is a free log retrieval operation binding the contract event 0x10e1f7ce9fd7d1b90a66d13a2ab3cb8dd7f29f3f8d520b143b063ccfbab6906b.
//
// Solidity: event RelayerRemoved(address indexed relayer)
func (_Bridge *BridgeFilterer) FilterRelayerRemoved(opts *bind.FilterOpts, relayer []common.Address) (*BridgeRelayerRemovedIterator, error) {

	var relayerRule []interface{}
	for _, relayerItem := range relayer {
		relayerRule = append(relayerRule, relayerItem)
	}

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "RelayerRemoved", relayerRule)
	if err != nil {
		return nil, err
	}
	return &BridgeRelayerRemovedIterator{contract: _Bridge.contract, event: "RelayerRemoved", logs: logs, sub: sub}, nil
}

// WatchRelayerRemoved is a free log subscription operation binding the contract event 0x10e1f7ce9fd7d1b90a66d13a2ab3cb8dd7f29f3f8d520b143b063ccfbab6906b.
//
// Solidity: event RelayerRemoved(address indexed relayer)
func (_Bridge *BridgeFilterer) WatchRelayerRemoved(opts *bind.WatchOpts, sink chan<- *BridgeRelayerRemoved, relayer []common.Address) (event.Subscription, error) {

	var relayerRule []interface{}
	for _, relayerItem := range relayer {
		relayerRule = append(relayerRule, relayerItem)
	}

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "RelayerRemoved", relayerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeRelayerRemoved)
				if err := _Bridge.contract.UnpackLog(event, "RelayerRemoved", log); err != nil {
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

// ParseRelayerRemoved is a log parse operation binding the contract event 0x10e1f7ce9fd7d1b90a66d13a2ab3cb8dd7f29f3f8d520b143b063ccfbab6906b.
//
// Solidity: event RelayerRemoved(address indexed relayer)
func (_Bridge *BridgeFilterer) ParseRelayerRemoved(log types.Log) (*BridgeRelayerRemoved, error) {
	event := new(BridgeRelayerRemoved)
	if err := _Bridge.contract.UnpackLog(event, "RelayerRemoved", log); err != nil {
		return nil, err
	}
	return event, nil
}

// BridgeRelayerThresholdChangedIterator is returned from FilterRelayerThresholdChanged and is used to iterate over the raw logs and unpacked data for RelayerThresholdChanged events raised by the Bridge contract.
type BridgeRelayerThresholdChangedIterator struct {
	Event *BridgeRelayerThresholdChanged // Event containing the contract specifics and raw log

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
func (it *BridgeRelayerThresholdChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeRelayerThresholdChanged)
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
		it.Event = new(BridgeRelayerThresholdChanged)
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
func (it *BridgeRelayerThresholdChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeRelayerThresholdChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeRelayerThresholdChanged represents a RelayerThresholdChanged event raised by the Bridge contract.
type BridgeRelayerThresholdChanged struct {
	NewThreshold *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterRelayerThresholdChanged is a free log retrieval operation binding the contract event 0xa20d6b84cd798a24038be305eff8a45ca82ef54a2aa2082005d8e14c0a4746c8.
//
// Solidity: event RelayerThresholdChanged(uint256 indexed newThreshold)
func (_Bridge *BridgeFilterer) FilterRelayerThresholdChanged(opts *bind.FilterOpts, newThreshold []*big.Int) (*BridgeRelayerThresholdChangedIterator, error) {

	var newThresholdRule []interface{}
	for _, newThresholdItem := range newThreshold {
		newThresholdRule = append(newThresholdRule, newThresholdItem)
	}

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "RelayerThresholdChanged", newThresholdRule)
	if err != nil {
		return nil, err
	}
	return &BridgeRelayerThresholdChangedIterator{contract: _Bridge.contract, event: "RelayerThresholdChanged", logs: logs, sub: sub}, nil
}

// WatchRelayerThresholdChanged is a free log subscription operation binding the contract event 0xa20d6b84cd798a24038be305eff8a45ca82ef54a2aa2082005d8e14c0a4746c8.
//
// Solidity: event RelayerThresholdChanged(uint256 indexed newThreshold)
func (_Bridge *BridgeFilterer) WatchRelayerThresholdChanged(opts *bind.WatchOpts, sink chan<- *BridgeRelayerThresholdChanged, newThreshold []*big.Int) (event.Subscription, error) {

	var newThresholdRule []interface{}
	for _, newThresholdItem := range newThreshold {
		newThresholdRule = append(newThresholdRule, newThresholdItem)
	}

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "RelayerThresholdChanged", newThresholdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeRelayerThresholdChanged)
				if err := _Bridge.contract.UnpackLog(event, "RelayerThresholdChanged", log); err != nil {
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

// ParseRelayerThresholdChanged is a log parse operation binding the contract event 0xa20d6b84cd798a24038be305eff8a45ca82ef54a2aa2082005d8e14c0a4746c8.
//
// Solidity: event RelayerThresholdChanged(uint256 indexed newThreshold)
func (_Bridge *BridgeFilterer) ParseRelayerThresholdChanged(log types.Log) (*BridgeRelayerThresholdChanged, error) {
	event := new(BridgeRelayerThresholdChanged)
	if err := _Bridge.contract.UnpackLog(event, "RelayerThresholdChanged", log); err != nil {
		return nil, err
	}
	return event, nil
}

// BridgeRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the Bridge contract.
type BridgeRoleGrantedIterator struct {
	Event *BridgeRoleGranted // Event containing the contract specifics and raw log

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
func (it *BridgeRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeRoleGranted)
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
		it.Event = new(BridgeRoleGranted)
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
func (it *BridgeRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeRoleGranted represents a RoleGranted event raised by the Bridge contract.
type BridgeRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Bridge *BridgeFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*BridgeRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &BridgeRoleGrantedIterator{contract: _Bridge.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Bridge *BridgeFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *BridgeRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeRoleGranted)
				if err := _Bridge.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Bridge *BridgeFilterer) ParseRoleGranted(log types.Log) (*BridgeRoleGranted, error) {
	event := new(BridgeRoleGranted)
	if err := _Bridge.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	return event, nil
}

// BridgeRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the Bridge contract.
type BridgeRoleRevokedIterator struct {
	Event *BridgeRoleRevoked // Event containing the contract specifics and raw log

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
func (it *BridgeRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeRoleRevoked)
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
		it.Event = new(BridgeRoleRevoked)
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
func (it *BridgeRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeRoleRevoked represents a RoleRevoked event raised by the Bridge contract.
type BridgeRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Bridge *BridgeFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*BridgeRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &BridgeRoleRevokedIterator{contract: _Bridge.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Bridge *BridgeFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *BridgeRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeRoleRevoked)
				if err := _Bridge.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Bridge *BridgeFilterer) ParseRoleRevoked(log types.Log) (*BridgeRoleRevoked, error) {
	event := new(BridgeRoleRevoked)
	if err := _Bridge.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	return event, nil
}

// BridgeUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the Bridge contract.
type BridgeUnpausedIterator struct {
	Event *BridgeUnpaused // Event containing the contract specifics and raw log

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
func (it *BridgeUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeUnpaused)
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
		it.Event = new(BridgeUnpaused)
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
func (it *BridgeUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeUnpaused represents a Unpaused event raised by the Bridge contract.
type BridgeUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Bridge *BridgeFilterer) FilterUnpaused(opts *bind.FilterOpts) (*BridgeUnpausedIterator, error) {

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &BridgeUnpausedIterator{contract: _Bridge.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Bridge *BridgeFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *BridgeUnpaused) (event.Subscription, error) {

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeUnpaused)
				if err := _Bridge.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Bridge *BridgeFilterer) ParseUnpaused(log types.Log) (*BridgeUnpaused, error) {
	event := new(BridgeUnpaused)
	if err := _Bridge.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	return event, nil
}
