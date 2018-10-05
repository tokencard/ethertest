package ethertest

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

// OneEth is 10^18 Wei
var OneEth = big.NewInt(1000000000000000000)

// EthToWei converts ETH value to Wei
func EthToWei(amount int) *big.Int {
	r := big.NewInt(1).Set(OneEth)
	return r.Mul(r, big.NewInt(int64(amount)))
}

// IsSuccessful returns boolean indicating if the transaction was successful
func IsSuccessful(be TestBackend, tx *types.Transaction) (bool, error) {
	r, err := be.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		return false, err
	}
	return r.Status == types.ReceiptStatusSuccessful, nil
}
