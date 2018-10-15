package ethertest_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/tokencard/ethertest"
	"github.com/tokencard/ethertest/test/bindings"
)

func TestContract(t *testing.T) {
	var tr = ethertest.NewTestRig()
	var owner = ethertest.NewAccount()

	tr.AddGenesisAccountAllocation(owner.Address(), ethertest.EthToWei(100))
	tr.AddCoverageForContracts("./test/build/test/combined.json", "test/contracts/test.sol")

	require := require.New(t)
	be := tr.NewTestBackend(ethertest.WithBlockGasLimit(8000000), ethertest.WithBlockchainTime(time.Now().Add(-24*time.Hour)))
	_, tx, testBinding, err := bindings.DeployTest(owner.TransactOpts(), be, "initial value")
	require.Nil(err)
	be.Commit()

	receipt, err := be.TransactionReceipt(context.Background(), tx.Hash())
	require.Nil(err)
	require.Equal(types.ReceiptStatusSuccessful, receipt.Status)

	tx, err = testBinding.SetValue(owner.TransactOpts(), "new value")
	require.Nil(err)
	be.Commit()

	successful, err := ethertest.IsSuccessful(be, tx)
	require.Nil(err)
	require.True(successful)

	value, err := testBinding.Value(nil)
	require.Nil(err)
	require.Equal("new value", value)

	tr.ExpectMinimumCoverage("test.sol:Test", 100.0)
	tr.PrintGasUsage(os.Stdout)
}
