package backends

func (b *SimulatedBackend) Close() error {

	b.blockchain.Stop()

	return nil
}
