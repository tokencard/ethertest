package ethertest

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/tokencard/ethertest/backends"
)

// TestRig ...
type TestRig struct {
	genesisAlloc core.GenesisAlloc
	contracts    map[string]*contract
}

// NewTestRig creates a new instance of a test rig
func NewTestRig() *TestRig {
	return &TestRig{
		genesisAlloc: core.GenesisAlloc{},
		contracts:    map[string]*contract{},
	}
}

// TestBackend is interface to an go-ethereum test backend
type TestBackend interface {
	bind.ContractBackend
	BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	Commit()
	Rollback()
	AdjustTime(adjustment time.Duration) error
}

// NewTestBackend creates a new instance of TestBackend
func (t *TestRig) NewTestBackend() TestBackend {
	return backends.NewSimulatedBackend(t.genesisAlloc, vm.Config{
		Debug:  true,
		Tracer: t,
	})
}

// AddGenesisAccountAllocation adds a GenesisAccount allocation to the test rig.
// When a new TestBackend is created, current genesis account allocations are used.
func (t *TestRig) AddGenesisAccountAllocation(a common.Address, balance *big.Int) *TestRig {
	t.genesisAlloc[a] = core.GenesisAccount{Balance: balance}
	return t
}

func (t *TestRig) AddCoverageForContracts(combinedJSON string, contractFiles ...string) *TestRig {
	sourceCode := map[string][]byte{}
	for _, path := range contractFiles {
		name := filepath.Base(path)
		source, err := ioutil.ReadFile(path)
		if err != nil {
			panic(fmt.Errorf("Could not read %q: %s", path, err.Error()))
		}
		sourceCode[name] = source
	}

	f, err := os.Open(combinedJSON)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	sc := &solcCombined{}
	err = json.NewDecoder(f).Decode(sc)

	if err != nil {
		panic(err)
	}

	for name := range sourceCode {
		_, found := sc.Sources[name]
		if !found {
			panic(fmt.Errorf("Could not find %q in the combined-json", name))
		}
	}

	for n, s := range sourceCode {
		ss := sc.Sources[n]
		for cn, scon := range sc.Contracts {
			if scon.BinRuntime != "" && strings.HasPrefix(cn+":", n) {
				con := newContract(s, ss, scon)
				t.contracts[cn] = con
			}
		}
	}

	return t
}

func (t *TestRig) ExpectMinimumCoverage(name string, expectedCoverage float64) {
	c, found := t.contracts[name]
	if !found {
		keys := []string{}
		for k := range t.contracts {
			keys = append(keys, k)
		}
		panic(fmt.Errorf("Could not find contract %q, available: %q", name, keys))
	}

	if c.percentageCovered() < expectedCoverage {
		fmt.Println()
		fmt.Printf("Coverage for %q:\n", name)
		c.Print()
		panic(fmt.Errorf("Contract %q has %.2f%% coverage (expected: %.2f%%)", name, c.percentageCovered(), expectedCoverage))
	}

	fmt.Printf("\nCoverage for %q: %.2f%%\n", name, c.percentageCovered())

}

func (t *TestRig) CaptureStart(from common.Address, to common.Address, call bool, input []byte, gas uint64, value *big.Int) error {
	return nil
}
func (t *TestRig) CaptureState(env *vm.EVM, pc uint64, op vm.OpCode, gas, cost uint64, memory *vm.Memory, stack *vm.Stack, contract *vm.Contract, depth int, err error) error {
	for _, c := range t.contracts {
		c.executed(contract.CodeHash, pc)
	}

	return nil
}
func (t *TestRig) CaptureFault(env *vm.EVM, pc uint64, op vm.OpCode, gas, cost uint64, memory *vm.Memory, stack *vm.Stack, contract *vm.Contract, depth int, err error) error {
	return nil
}

func (t *TestRig) CaptureEnd(output []byte, gasUsed uint64, tm time.Duration, err error) error {
	return nil
}
