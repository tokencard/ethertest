Ethertest
===

Ethertest is a simple Go testing framework for smart contracts based on [geth](https://github.com/ethereum/go-ethereum).

It enables writing isolated and repeatable tests without cross-test side effects.


## Getting Started
Ethertest enables you to interact with smart contracts deployed to a full features in-memory blockchain.
It helps you set up repeatable conditions on the the blockchain, interact with deployed contracts and inspect the results of the interaction.

Before testing the contracts, you'll have to compile smart contract and generate `combined-json` file containing `bin-runtime,srcmap-runtime,ast,srcmap,bin`.
This is necessary to create a code coverage report:

```sh
  solc --optimize --overwrite --bin --abi --combined-json bin-runtime,srcmap-runtime,ast,srcmap,bin -o <build_path> <your_contract>.sol
```


The next step is to generate `geth` bindings for the contract - this will give you a type-safe interface for deploying / interacting with the contract on the blockchain:

```sh
  abigen --abi <path to .abi file>  --bin <path to .bin file> --pkg <go package name> --type=<go type name> --out <package directory>/<go file name>
```

Now you're ready to write a test.
A sample `gotest` style




## Test Rig
Test Rig is the top-level factory for new instances of Test Backend.
New instance of TestRig can be created using the `NewTestRig` function:

```go
  testRig := ethertest.NewTestRig()
```

### Genesis Account Allocation
When a new test backend (blockchain simulation) is created all accounts have 0 ETH.
This can be changed adding genesis account allocation before creating the test backend.
```go
  testRig.AddGenesisAccountAllocation(<ether address>, <amount in WEI (*big.Int)>)
```

Genesis allocation is memorized in the test rig, so every creation of a new test backend will contain all added allocations.

## Accounts
To simplify Ethereum account private/public key creation and handling for the test purposes, Ethertest offers the type `*ethertest.Account`

### Creating Account

```go
  acct := etheretest.NewAccount()
```

### Account Address

```go
  addr := acct.Address()
```

### Account ETH Transfer

```go
  err := acct.Transfer(backend, <ether address>, <amount in WEI (*big.Int)>)
```

```go
  acct.MustTransfer(backend, <ether address>, <amount in WEI (*big.Int)>)
```

### Account Balance

```go
  balance := acct.Balance(backend)
```

### Account TransactOpts

```go
  to := acct.TransactOpts([<modifiers>])
```

Modifiers:
```go
  ethertest.WithGasPrice(<amount in WEI (*big.Int)>)
```

## Test Backend
### NewTestBackend options

## Code Coverage
## Gas Usage
## Tracing
