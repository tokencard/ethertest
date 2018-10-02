package ethertest

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

var ErrTransactionFailed = errors.New("Transaction Failed")

func NewWallet() *Wallet {
	k, err := crypto.GenerateKey()

	// will only happen if RNG fails big time
	if err != nil {
		panic(err)
	}

	return &Wallet{k}
}

type Wallet struct {
	pk *ecdsa.PrivateKey
}

func (w *Wallet) Address() common.Address {
	return crypto.PubkeyToAddress(w.pk.PublicKey)
}

func (w *Wallet) MustTransfer(be TestBackend, to common.Address, amount *big.Int) {
	err := w.Transfer(be, to, amount)
	if err != nil {
		panic(err)
	}
}

func (w *Wallet) Transfer(be TestBackend, to common.Address, amount *big.Int) error {
	n, err := be.PendingNonceAt(context.Background(), w.Address())
	if err != nil {
		return err
	}

	gasPrice, err := be.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}

	tx := types.NewTransaction(n, to, amount, 41000, gasPrice, nil)

	signed, err := types.SignTx(tx, types.HomesteadSigner{}, w.pk)
	if err != nil {
		return err
	}

	err = be.SendTransaction(context.Background(), signed)
	if err != nil {
		return err
	}

	be.Commit()

	rcpt, err := be.TransactionReceipt(context.Background(), signed.Hash())
	if err != nil {
		return err
	}
	if rcpt.Status != types.ReceiptStatusSuccessful {
		return ErrTransactionFailed
	}
	return nil
}

func (w *Wallet) TransactOpts() *bind.TransactOpts {
	return bind.NewKeyedTransactor(w.pk)
}

func (w *Wallet) TransactOptsWithGasLimit(gasLimit uint64) *bind.TransactOpts {
	to := bind.NewKeyedTransactor(w.pk)
	to.GasLimit = gasLimit
	return to
}

func (w *Wallet) Balance(be TestBackend) *big.Int {
	b, err := be.BalanceAt(context.Background(), w.Address(), nil)
	if err != nil {
		panic(err)
	}
	return b
}
