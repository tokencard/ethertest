// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

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
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// TestABI is the input ABI used to generate the binding from.
const TestABI = "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_value\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_value\",\"type\":\"string\"}],\"name\":\"setValue\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"value\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"willFail\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"}]"

// TestBin is the compiled bytecode used for deploying new contracts.
var TestBin = "0x608060405234801561001057600080fd5b506040516104973803806104978339818101604052602081101561003357600080fd5b810190808051604051939291908464010000000082111561005357600080fd5b90830190602082018581111561006857600080fd5b825164010000000081118282018810171561008257600080fd5b82525081516020918201929091019080838360005b838110156100af578181015183820152602001610097565b50505050905090810190601f1680156100dc5780820380516001836020036101000a031916815260200191505b50604052505081516100f6915060009060208401906100fd565b5050610198565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061013e57805160ff191683800117855561016b565b8280016001018555821561016b579182015b8281111561016b578251825591602001919060010190610150565b5061017792915061017b565b5090565b61019591905b808211156101775760008155600101610181565b90565b6102f0806101a76000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80633fa4f24514610046578063625676a2146100c357806393a09352146100cd575b600080fd5b61004e61013d565b6040805160208082528351818301528351919283929083019185019080838360005b83811015610088578181015183820152602001610070565b50505050905090810190601f1680156100b55780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6100cb6101cb565b005b6100cb600480360360208110156100e357600080fd5b8101906020810181356401000000008111156100fe57600080fd5b82018360208201111561011057600080fd5b8035906020019184600183028401116401000000008311171561013257600080fd5b5090925090506101d5565b6000805460408051602060026001851615610100026000190190941693909304601f810184900484028201840190925281815292918301828280156101c35780601f10610198576101008083540402835291602001916101c3565b820191906000526020600020905b8154815290600101906020018083116101a657829003601f168201915b505050505081565b6101d36101e6565b565b6101e16000838361021f565b505050565b6040805162461bcd60e51b81526020600482015260096024820152681dda5b1b0819985a5b60ba1b604482015290519081900360640190fd5b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106102605782800160ff1982351617855561028d565b8280016001018555821561028d579182015b8281111561028d578235825591602001919060010190610272565b5061029992915061029d565b5090565b6102b791905b8082111561029957600081556001016102a3565b9056fea2646970667358221220330a41b89c04be6a8f212a5f4b93ff2c06e4eb01c48df30188361de6eba7dcb764736f6c63430006040033"

// DeployTest deploys a new Ethereum contract, binding an instance of Test to it.
func DeployTest(auth *bind.TransactOpts, backend bind.ContractBackend, _value string) (common.Address, *types.Transaction, *Test, error) {
	parsed, err := abi.JSON(strings.NewReader(TestABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(TestBin), backend, _value)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Test{TestCaller: TestCaller{contract: contract}, TestTransactor: TestTransactor{contract: contract}, TestFilterer: TestFilterer{contract: contract}}, nil
}

// Test is an auto generated Go binding around an Ethereum contract.
type Test struct {
	TestCaller     // Read-only binding to the contract
	TestTransactor // Write-only binding to the contract
	TestFilterer   // Log filterer for contract events
}

// TestCaller is an auto generated read-only Go binding around an Ethereum contract.
type TestCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TestTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TestFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TestSession struct {
	Contract     *Test             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TestCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TestCallerSession struct {
	Contract *TestCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// TestTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TestTransactorSession struct {
	Contract     *TestTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TestRaw is an auto generated low-level Go binding around an Ethereum contract.
type TestRaw struct {
	Contract *Test // Generic contract binding to access the raw methods on
}

// TestCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TestCallerRaw struct {
	Contract *TestCaller // Generic read-only contract binding to access the raw methods on
}

// TestTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TestTransactorRaw struct {
	Contract *TestTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTest creates a new instance of Test, bound to a specific deployed contract.
func NewTest(address common.Address, backend bind.ContractBackend) (*Test, error) {
	contract, err := bindTest(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Test{TestCaller: TestCaller{contract: contract}, TestTransactor: TestTransactor{contract: contract}, TestFilterer: TestFilterer{contract: contract}}, nil
}

// NewTestCaller creates a new read-only instance of Test, bound to a specific deployed contract.
func NewTestCaller(address common.Address, caller bind.ContractCaller) (*TestCaller, error) {
	contract, err := bindTest(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestCaller{contract: contract}, nil
}

// NewTestTransactor creates a new write-only instance of Test, bound to a specific deployed contract.
func NewTestTransactor(address common.Address, transactor bind.ContractTransactor) (*TestTransactor, error) {
	contract, err := bindTest(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestTransactor{contract: contract}, nil
}

// NewTestFilterer creates a new log filterer instance of Test, bound to a specific deployed contract.
func NewTestFilterer(address common.Address, filterer bind.ContractFilterer) (*TestFilterer, error) {
	contract, err := bindTest(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestFilterer{contract: contract}, nil
}

// bindTest binds a generic wrapper to an already deployed contract.
func bindTest(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TestABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Test *TestRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Test.Contract.TestCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Test *TestRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Test.Contract.TestTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Test *TestRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Test.Contract.TestTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Test *TestCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Test.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Test *TestTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Test.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Test *TestTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Test.Contract.contract.Transact(opts, method, params...)
}

// Value is a free data retrieval call binding the contract method 0x3fa4f245.
//
// Solidity: function value() constant returns(string)
func (_Test *TestCaller) Value(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _Test.contract.Call(opts, out, "value")
	return *ret0, err
}

// Value is a free data retrieval call binding the contract method 0x3fa4f245.
//
// Solidity: function value() constant returns(string)
func (_Test *TestSession) Value() (string, error) {
	return _Test.Contract.Value(&_Test.CallOpts)
}

// Value is a free data retrieval call binding the contract method 0x3fa4f245.
//
// Solidity: function value() constant returns(string)
func (_Test *TestCallerSession) Value() (string, error) {
	return _Test.Contract.Value(&_Test.CallOpts)
}

// WillFail is a free data retrieval call binding the contract method 0x625676a2.
//
// Solidity: function willFail() constant returns()
func (_Test *TestCaller) WillFail(opts *bind.CallOpts) error {
	var ()
	out := &[]interface{}{}
	err := _Test.contract.Call(opts, out, "willFail")
	return err
}

// WillFail is a free data retrieval call binding the contract method 0x625676a2.
//
// Solidity: function willFail() constant returns()
func (_Test *TestSession) WillFail() error {
	return _Test.Contract.WillFail(&_Test.CallOpts)
}

// WillFail is a free data retrieval call binding the contract method 0x625676a2.
//
// Solidity: function willFail() constant returns()
func (_Test *TestCallerSession) WillFail() error {
	return _Test.Contract.WillFail(&_Test.CallOpts)
}

// SetValue is a paid mutator transaction binding the contract method 0x93a09352.
//
// Solidity: function setValue(string _value) returns()
func (_Test *TestTransactor) SetValue(opts *bind.TransactOpts, _value string) (*types.Transaction, error) {
	return _Test.contract.Transact(opts, "setValue", _value)
}

// SetValue is a paid mutator transaction binding the contract method 0x93a09352.
//
// Solidity: function setValue(string _value) returns()
func (_Test *TestSession) SetValue(_value string) (*types.Transaction, error) {
	return _Test.Contract.SetValue(&_Test.TransactOpts, _value)
}

// SetValue is a paid mutator transaction binding the contract method 0x93a09352.
//
// Solidity: function setValue(string _value) returns()
func (_Test *TestTransactorSession) SetValue(_value string) (*types.Transaction, error) {
	return _Test.Contract.SetValue(&_Test.TransactOpts, _value)
}
