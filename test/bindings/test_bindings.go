// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// TestABI is the input ABI used to generate the binding from.
const TestABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"value\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_value\",\"type\":\"string\"}],\"name\":\"setValue\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_value\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// TestBin is the compiled bytecode used for deploying new contracts.
const TestBin = `608060405234801561001057600080fd5b5060405161034d38038061034d83398101604052805101805161003a906000906020840190610041565b50506100dc565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061008257805160ff19168380011785556100af565b828001600101855582156100af579182015b828111156100af578251825591602001919060010190610094565b506100bb9291506100bf565b5090565b6100d991905b808211156100bb57600081556001016100c5565b90565b610262806100eb6000396000f30060806040526004361061004b5763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416633fa4f245811461005057806393a09352146100da575b600080fd5b34801561005c57600080fd5b506100656100fc565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561009f578181015183820152602001610087565b50505050905090810190601f1680156100cc5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b3480156100e657600080fd5b506100fa600480356024810191013561018a565b005b6000805460408051602060026001851615610100026000190190941693909304601f810184900484028201840190925281815292918301828280156101825780601f1061015757610100808354040283529160200191610182565b820191906000526020600020905b81548152906001019060200180831161016557829003601f168201915b505050505081565b6101966000838361019b565b505050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106101dc5782800160ff19823516178555610209565b82800160010185558215610209579182015b828111156102095782358255916020019190600101906101ee565b50610215929150610219565b5090565b61023391905b80821115610215576000815560010161021f565b905600a165627a7a72305820f793057130c5e6a2fd97245f60bb8fd877f3151c9a90de3a0dfd5752804ebe0c0029`

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

// SetValue is a paid mutator transaction binding the contract method 0x93a09352.
//
// Solidity: function setValue(_value string) returns()
func (_Test *TestTransactor) SetValue(opts *bind.TransactOpts, _value string) (*types.Transaction, error) {
	return _Test.contract.Transact(opts, "setValue", _value)
}

// SetValue is a paid mutator transaction binding the contract method 0x93a09352.
//
// Solidity: function setValue(_value string) returns()
func (_Test *TestSession) SetValue(_value string) (*types.Transaction, error) {
	return _Test.Contract.SetValue(&_Test.TransactOpts, _value)
}

// SetValue is a paid mutator transaction binding the contract method 0x93a09352.
//
// Solidity: function setValue(_value string) returns()
func (_Test *TestTransactorSession) SetValue(_value string) (*types.Transaction, error) {
	return _Test.Contract.SetValue(&_Test.TransactOpts, _value)
}
