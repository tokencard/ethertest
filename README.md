Ethertest
===

This is a set of tools that will help you test smart contracts using Golang.

Ethertest will help you run repeatable and isolated tests, assert code coverage of your Solidity, record traces of contract execution and estimate Gas usage of external methods.

Core of Ethertest is based on [Geth](https://github.com/ethereum/go-ethereum), making execution EVM code as close to running on a Geth node as possible.

## Overview

Ethertest consists of three main components **TestRig**, **TestBackend** and **Account**:

### TestRig

TestRig is factory of new in-memory block chain instances (TestBackend).
TestRig records all code coverage and tracing information for all TestBackend
instances it has created, so it is best kept as a singleton in the test package.

### TestBackend

TestBackend is an in-memory blockchain your tests can interact with.
Ideally you should create one TestBackend per test, making each test
isolated from side effects of other tests.

TestBackend implements `github.com/ethereum/go-ethereum/accounts/abi/bind/ContractBackend` interface, so it can be used by `abigen` generated contract bindings.
In addition it will give you access to the current blockchain account balances, transaction receipts, and option to either commit (mine one block) or roll back (go one block back in time).

After the test is done, `Close()` method should be executed on TestBackend.
This will free allocated caches and stop go routines.

### Account

Account encapsules private/public key for an Ethereum address. Every time `NewAccount()` function is called, a new random private/public key is created.

Account has functions to get the account address, get balance of the account address, transfer funds to another address and create `github.com/ethereum/go-ethereum/accounts/abi/bind.TransactOpts` needed to call transaction methods on contract bindings.

## Code Coverage

To enable code coverage, TestRig needs `AST` and `srcmap` and `bin` generated by solidity compiler in a file called `combined-json`.
Following command will generate `combined-json` file:

```sh
  solc --optimize --overwrite --bin --abi --combined-json bin-runtime,srcmap-runtime,ast,srcmap,bin -o <build_path> <your_contract>.sol
```

Every contract that need code coverage has to be registered with the TestRig:
```go
  tr.AddCoverageForContracts("<path to combined.json>", "<path to the solidity source file>")
```

After all tests have finished, code coverage can be asserted with:
```go
  testRig.ExpectMinimumCoverage("<sol file name>:<contract name>", <expected coverage percent as float64>)
```

If the coverage of the contract is lower than expected, the method will print a coloured source of the contract (green for executed, red for not executed) and panic with a message stating expected and current code coverage.


## Genesis Account Allocation
When a new TestBackend is created all accounts have 0 ETH, making the whole blockchain unusable.
This can be changed by adding genesis account allocation to `TestRig` before creating the TestBackend:

```go
  testRig.AddGenesisAccountAllocation(<ether address>, <amount in WEI (*big.Int)>)
```

Genesis allocation is memorized in the TestRig, so every creation of a new TestBackend will contain all added allocations.

## Gas Usage
After all have finished, gas usage of the contracts can be printed by calling `PrintGasUsage` method of TestRig:

```go
  testRig.PrintGasUsage(os.Stdout)
```

Output will be similar to:

```
Gas Usage for "test.sol:Test"
+------------------+-------+-------+-------+
|  FUNCTION NAME   |  MIN  |  MED  |  MAX  |
+------------------+-------+-------+-------+
| setValue(string) | 33109 | 33109 | 33109 |
+------------------+-------+-------+-------+
```

Where MIN, MED and MAX are minimum, median (50 percentile), and maximum gas spent
calling each function in a transaction context.


## LastExecuted

When a transaction fails, it is sometimes useful to find out what was the last line of
code executed. Method `LastExecuted()` on the TestRig will return a string containing file name, line number and the appropriate source code snippet.
