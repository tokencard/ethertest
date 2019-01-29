package backends

import (
	"github.com/tokencard/ethertest/peek"
)

func (b *SimulatedBackend) Close() error {
	_, err := peek.Call("blockchain.stateCache.db.cleans", b, "Close")
	return err
}
