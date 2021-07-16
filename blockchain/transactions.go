package blockchain

type Tx struct {
	Id        string
	Timestamp int
	TxIns     []*TxIn
	TxOuts    []*TxOut
}

type TxIn struct {
	Owner  string
	Amount int
}

type TxOut struct {
	Owner  string
	Amount int
}
