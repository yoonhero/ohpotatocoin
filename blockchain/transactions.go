package blockchain

import (
	"errors"
	"time"

	"github.com/yoonhero/ohpotatocoin/utils"
)

const (
	minerReward int = 50
)

// not confirmed transactions list
type mempool struct {
	Txs []*Tx
}

var Mempool *mempool = &mempool{}

// transaction struct
type Tx struct {
	ID        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
}

// hashing the transaction to get ID
func (t *Tx) getId() {
	t.ID = utils.Hash(t)
}

// transaction input value
type TxIn struct {
	TxID  string
	Index int
	Owner string `json:"owner"`
}

// transaction output value
type TxOut struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

// unspent transaction structure
type UTxOut struct {
	TxID   string
	Index  int
	Amount int
}

// give coin when mining
func makeCoinbaseTx(address string) *Tx {
	txIns := []*TxIn{
		{"", -1, "COINBASE"},
	}
	txOuts := []*TxOut{
		{address, minerReward},
	}

	tx := Tx{
		ID:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return &tx
}

// make transaction
func makeTx(from, to string, amount int) (*Tx, error) {
	if Blockchain().BalancByAddress(from) < amount {
		return nil, errors.New("Not enough money")
	}
	var txOuts []*TxOut
	var txIns []*TxIn
	total := 0
	uTxOuts := Blockchain().UTxOutsByAddress(from)
	for _, uTxOut := range uTxOuts {
		if total > amount {
			break
		}
		txIn := &TxIn{uTxOut.TxID, uTxOut.Index, from}
		txIns = append(txIns, txIn)
		total += uTxOut.Amount
	}
	if change := total - amount; change != 0 {
		changeTxOut := &TxOut{from, change}
		txOuts = append(txOuts, changeTxOut)
	}
	txOut := &TxOut{to, amount}
	txOuts = append(txOuts, txOut)
	tx := &Tx{
		ID:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return tx, nil
}

// add transaction
func (m *mempool) AddTx(to string, amount int) error {
	tx, err := makeTx("yoonhero", to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}

func (m *mempool) TxToConfirm() []*Tx {
	coinbase := makeCoinbaseTx("yoonhero")
	txs := m.Txs
	txs = append(txs, coinbase)
	m.Txs = nil
	return txs
}
