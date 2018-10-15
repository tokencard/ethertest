package ethertest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/olekukonko/tablewriter"
	"github.com/tokencard/ethertest/backends"
	"github.com/tokencard/ethertest/stats"
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

type interceptingBackend struct {
	TestBackend
	sentTransactions []*types.Transaction
	tr               *TestRig
}

func (ib *interceptingBackend) Commit() {
	ib.TestBackend.Commit()

	for _, t := range ib.sentTransactions {
		r, err := ib.TransactionReceipt(context.Background(), t.Hash())
		if err != nil {
			panic(err)
		}

		for _, c := range ib.tr.contracts {
			to := t.To()
			if to != nil {
				c.transactionCommited(*to, t.Data(), r.GasUsed)
			}
		}

	}
	ib.sentTransactions = nil
}

func (ib *interceptingBackend) SendTransaction(ctx context.Context, tx *types.Transaction) error {

	err := ib.TestBackend.SendTransaction(ctx, tx)
	if err != nil {
		return err
	}
	ib.sentTransactions = append(ib.sentTransactions, tx)
	return nil
}

type backendOption func(*backendOptions)

type backendOptions struct {
	blockchainTime time.Time
	blockGasLimit  uint64
}

// WithBlockchainTime sets the initial time on the blockchain.
// If not set, it will default to 1970-01-01T00:00:00Z
// Warning: every commit() will increase the time by 15 seconds.
// Once the blockchain time takes over the current time,
// simulated backend will not accept new blocks.
func WithBlockchainTime(t time.Time) func(*backendOptions) {
	return func(opt *backendOptions) {
		opt.blockchainTime = t
	}
}

// WithBlockGasLimit sets the block limit.
// If not set, it will default to 7981579.
func WithBlockGasLimit(limit uint64) func(*backendOptions) {
	return func(opt *backendOptions) {
		opt.blockGasLimit = limit
	}
}

// NewTestBackend creates a new instance of TestBackend
func (t *TestRig) NewTestBackend(opts ...backendOption) TestBackend {

	backendOptions := &backendOptions{
		blockGasLimit:  7981579,
		blockchainTime: time.Unix(0, 0),
	}
	for _, opt := range opts {
		opt(backendOptions)
	}

	sb := backends.NewSimulatedBackend(t.genesisAlloc, backendOptions.blockGasLimit, vm.Config{
		Debug:  true,
		Tracer: t,
	}, backendOptions.blockchainTime)

	return &interceptingBackend{
		TestBackend: sb,
		tr:          t,
	}
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
		sourceIndex := sc.findSourceIndex(n)
		if sourceIndex < 0 {
			panic(fmt.Errorf("Could not find %q in the source-index", n))
		}
		ss := sc.Sources[n]
		for cn, scon := range sc.Contracts {
			if scon.BinRuntime != "" && strings.HasPrefix(cn+":", n) {
				con, err := newContract(cn, s, ss, scon, sourceIndex)
				if err != nil {
					panic(err)
				}
				t.contracts[cn] = con
			}
		}
	}

	return t
}

func (t *TestRig) PrintGasUsage(w io.Writer) {
	for _, c := range t.contracts {

		if !c.hasAnyGasInformation() {
			continue
		}

		tw := tablewriter.NewWriter(w)
		fmt.Fprintf(w, "Gas Usage for %q\n", c.name)
		tw.SetHeader([]string{"Function Name", "Min", "Med", "Max"})

		functions := []*Function{}
		for _, f := range c.functions {
			functions = append(functions, f)
		}

		sort.Slice(functions, func(i int, j int) bool {
			return functions[i].name < functions[j].name
		})

		for _, f := range functions {
			tw.Append([]string{
				f.name,
				fmt.Sprintf("%d", stats.Uint64Min(f.gasUsed)),
				fmt.Sprintf("%d", stats.Uint64Median(f.gasUsed)),
				fmt.Sprintf("%d", stats.Uint64Max(f.gasUsed)),
			},
			)
		}
		tw.Render()
		fmt.Fprintln(w)
	}

}

func (t *TestRig) CoverageOf(name string) float64 {
	c, found := t.contracts[name]
	if !found {
		keys := []string{}
		for k := range t.contracts {
			keys = append(keys, k)
		}

		panic(fmt.Errorf("Could not find contract %q, available: %q", name, keys))
	}
	return c.percentageCovered()
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
		c.executed(pc, contract.Address(), contract)
	}

	return nil
}

func (t *TestRig) CaptureFault(env *vm.EVM, pc uint64, op vm.OpCode, gas, cost uint64, memory *vm.Memory, stack *vm.Stack, contract *vm.Contract, depth int, err error) error {
	return nil
}

func (t *TestRig) CaptureEnd(output []byte, gasUsed uint64, tm time.Duration, err error) error {
	return nil
}
