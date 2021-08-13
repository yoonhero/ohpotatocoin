package blockchain

type fakeDB struct {
	fakeLoadChain func() []byte
	fakeFindBlock func() []byte
}

func (f fakeDB) FindBlock(hash string) []byte {
	return f.fakeFindBlock()
}

func (f fakeDB) LoadChain() []byte {
	return f.fakeLoadChain()
}

func (fakeDB) SaveBlock(hash string, data []byte) {
	return
}

func (fakeDB) SaveChain(data []byte) {
	return
}

func (fakeDB) DeleteAllBlocks() {
	return
}
