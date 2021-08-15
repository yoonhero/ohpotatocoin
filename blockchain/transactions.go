package blockchain

import (
	"encoding/hex"
	"errors"
	"sync"
	"time"

	"github.com/yoonhero/ohpotatocoin/utils"
	"github.com/yoonhero/ohpotatocoin/wallet"
)

const (
	minerReward int = 1
)

// not confirmed transactions list
type mempool struct {
	Txs map[string]*Tx
	m   sync.Mutex
}

var m *mempool = &mempool{}
var memOnce sync.Once

func Mempool() *mempool {
	memOnce.Do(func() {
		m = &mempool{
			Txs: make(map[string]*Tx),
		}
	})
	return m
}

// transaction struct
type Tx struct {
	ID        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
}

// transaction input value
type TxIn struct {
	TxID      string `json:"txId"`  // set unspent transactoin output id
	Index     int    `json:"index"` // set unspent transaction output index
	Signature string `json:"signature"`
}

// transaction output value
type TxOut struct {
	Address string `json:"address"`
	Amount  int    `json:"amount"`
}

// unspent transaction structure
type UTxOut struct {
	TxID   string `json:"txId"`
	Index  int    `json:"index"`
	Amount int    `json:"amount"`
}

// hashing the transaction to get ID
func (t *Tx) getId() {
	t.ID = utils.Hash(t)
}

// sign the transaction
func (t *Tx) sign(keyAsBytes []byte) {
	for _, txIn := range t.TxIns {
		txIn.Signature = wallet.Sign(t.ID, keyAsBytes)
	}
}

// validate the transaction
func validate(tx *Tx) bool {
	valid := true
	for _, txIn := range tx.TxIns {
		// find prev transaction and it exists or not
		prevTx := FindTx(Blockchain(), txIn.TxID)
		if prevTx == nil {
			valid = false
			break
		}

		// validate the private key and public key
		address := prevTx.TxOuts[txIn.Index].Address
		valid = wallet.Verify(txIn.Signature, tx.ID, address)

		if !valid {
			break
		}
	}
	return valid
}

// recognize that transaction is on mempool or not
func isOnMempool(uTxOut *UTxOut) (exists bool) {
Outer:
	for _, tx := range Mempool().Txs {
		for _, input := range tx.TxIns {
			if input.TxID == uTxOut.TxID && input.Index == uTxOut.Index {
				exists = true
				break Outer
			}
		}
	}
	return
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

var ErrorNoMonery = errors.New("not enough monery")
var ErrorNotValid = errors.New("Tx Invalid")

// make transaction
func makeTx(privkey, to string, amount int) (*Tx, error) {
	privkeyAsByte, err := hex.DecodeString(privkey)
	utils.HandleErr(err)
	w := wallet.RestApiWallet(privkeyAsByte)
	from := w.Address
	// if from's balance < amount return
	if BalanceByAddress(from, Blockchain()) < amount {
		return nil, ErrorNoMonery
	}

	var txOuts []*TxOut
	var txIns []*TxIn
	total := 0
	uTxOuts := UTxOutsByAddress(from, Blockchain())

	for _, uTxOut := range uTxOuts {
		// total is more than amount and return
		if total > amount {
			break
		}

		// make transaction input
		txIn := &TxIn{uTxOut.TxID, uTxOut.Index, from}
		txIns = append(txIns, txIn)

		// plus total
		total += uTxOut.Amount
	}

	// if there is change and return transaction output
	if change := total - amount; change != 0 {
		changeTxOut := &TxOut{from, change}
		txOuts = append(txOuts, changeTxOut)
	}

	// make tranasaction output
	txOut := &TxOut{to, amount}
	txOuts = append(txOuts, txOut)

	// make transaction
	tx := &Tx{
		ID:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()

	// sign transaction
	tx.sign(privkeyAsByte)

	// validate transaction
	valid := validate(tx)

	if !valid {
		return nil, ErrorNotValid
	}

	return tx, nil

}

// add transaction
func (m *mempool) AddTx(from, to string, amount int) (*Tx, error) {
	m.m.Lock()
	defer m.m.Unlock()
	tx, err := makeTx(from, to, amount)
	if err != nil {
		return nil, err
	}
	m.Txs[tx.ID] = tx
	return tx, nil
}

// transaction confirm
func (m *mempool) TxToConfirm(from string) []*Tx {
	m.m.Lock()
	defer m.m.Unlock()
	coinbase := makeCoinbaseTx(from)
	var txs []*Tx
	for _, tx := range m.Txs {
		txs = append(txs, tx)
	}
	txs = append(txs, coinbase)
	m.Txs = make(map[string]*Tx)
	return txs
}

func (m *mempool) AddPeerTx(tx *Tx) {
	m.m.Lock()
	defer m.m.Unlock()

	m.Txs[tx.ID] = tx
}
