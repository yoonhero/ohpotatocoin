// persistence of block
// connected to DB to save data
// using bolt DB (bitcoin levelDB)

package blockchain

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/yoonhero/ohpotatocoin/db"
	"github.com/yoonhero/ohpotatocoin/utils"
)

const (
	defaultDifficulty  int = 2
	difficultyInterval int = 5
	blockInterval      int = 2
	allowedRange       int = 2
)

// type blockchain
// blocks is slice of []Block
type blockchain struct {
	NewestHash        string `json:"newestHash"`
	Height            int    `json:"height"`
	CurrentDifficulty int    `json:"currentDifficulty"`
	m                 sync.Mutex
}

// variable blockchain pointers
var b *blockchain

// variable struct that play func only one time
var once sync.Once

func (b *blockchain) restore(data []byte) {
	// decoder := gob.NewDecoder(bytes.NewReader(data))
	// decoder.Decode(b)
	utils.FromBytes(b, data)
}

// add block to blockchain
func (b *blockchain) AddBlock() *Block {
	// createBlock
	block := createBlock(b.NewestHash, b.Height+1, getDifficulty(b))

	// set newesthash new block's hash
	b.NewestHash = block.Hash
	// set height new block's height
	b.Height = block.Height

	b.CurrentDifficulty = block.Difficulty

	// persist the blockchain
	persistBlockchain(b)
	return block
}

// all blocks
func Blocks(b *blockchain) []*Block {
	b.m.Lock()
	defer b.m.Unlock()
	var blocks []*Block

	// start newesthash and its prevhash and find block
	// if prevhash dont exist = genesis block break
	hashCursor := b.NewestHash

	for {
		block, _ := FindBlock(hashCursor)
		blocks = append(blocks, block)
		if block.PrevHash != "" {
			hashCursor = block.PrevHash
		} else {
			break
		}
	}
	return blocks
}

// persist the blockchain data
func persistBlockchain(b *blockchain) {
	db.SaveCheckpoint(utils.ToBytes(b))
}

func Blockchain() *blockchain {
	// run only one time
	once.Do(func() {
		// initial blockchain struct
		b = &blockchain{Height: 0}

		// search for checkpoint on the db
		checkpoint := db.Checkpoint()

		if checkpoint == nil {
			// if blockchain don't exist create block
			b.AddBlock()
		} else {
			// restore data from db
			b.restore(checkpoint)
		}
	})
	// return type blockchain struct
	return b
}

func Status(b *blockchain, rw http.ResponseWriter) {
	b.m.Lock()
	defer b.m.Unlock()

	utils.HandleErr(json.NewEncoder(rw).Encode(b))
}

// all transactions
func Txs(b *blockchain) []*Tx {
	var txs []*Tx
	for _, block := range Blocks(b) {
		txs = append(txs, block.Transactions...)
	}
	return txs
}

// find specific transaction
func FindTx(b *blockchain, targetID string) *Tx {
	// for loop to find
	for _, tx := range Txs(b) {
		if tx.ID == targetID {
			return tx
		}
	}
	return nil
}

// recalculate difficulty of block by timestamp
func recalculateDifficulty(b *blockchain) int {
	// get all blocks
	allBlocks := Blocks(b)
	newestBlock := allBlocks[0]
	lastRecalculatedBlock := allBlocks[difficultyInterval-1]
	actualTime := (newestBlock.Timestamp / 60) - (lastRecalculatedBlock.Timestamp / 60)
	expectedTime := difficultyInterval * blockInterval
	if actualTime <= (expectedTime - allowedRange) {
		// if acuaultime < 8 difficulty + 1
		return b.CurrentDifficulty + 1
	} else if actualTime >= (expectedTime + allowedRange) {
		// if actualtime >= 12 difficulty - 1
		return b.CurrentDifficulty - 1
	}
	return b.CurrentDifficulty
}

func getDifficulty(b *blockchain) int {
	// if genesis block or not
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		// recalculate the difficulty
		return recalculateDifficulty(b)
	} else {
		return b.CurrentDifficulty
	}
}

// unspent transaction out by address
func UTxOutsByAddress(address string, b *blockchain) []*UTxOut {
	var uTxOuts []*UTxOut
	creatorTxs := make(map[string]bool)
	// for loop all blocks
	for _, block := range Blocks(b) {
		// for loop block transactions
		for _, tx := range block.Transactions {
			// for loop transactions input
			for _, input := range tx.TxIns {
				// input signature is coinbase and break loop
				if input.Signature == "COINBASE" {
					break
				}
				// same address with txouts.address
				if FindTx(b, input.TxID).TxOuts[input.Index].Address == address {
					creatorTxs[input.TxID] = true
				}
			}
			// for loop transactions output
			for index, output := range tx.TxOuts {
				// if output is owned by address
				if output.Address == address {
					// if it didn't spent yet
					if _, ok := creatorTxs[tx.ID]; !ok {
						uTxOut := &UTxOut{tx.ID, index, output.Amount}
						// if that transaction doesn't on the mempool append it
						if !isOnMempool(uTxOut) {
							uTxOuts = append(uTxOuts, uTxOut)
						}
					}
				}
			}
		}
	}

	return uTxOuts
}

// get balance of address
func BalanceByAddress(address string, b *blockchain) int {
	txOuts := UTxOutsByAddress(address, b)
	var amount int
	for _, txOut := range txOuts {
		amount += txOut.Amount
	}
	return amount
}

func (b *blockchain) Replace(newBlocks []*Block) {
	b.m.Lock()
	defer b.m.Unlock()
	b.CurrentDifficulty = newBlocks[0].Difficulty
	b.Height = len(newBlocks)
	b.NewestHash = newBlocks[0].Hash
	persistBlockchain(b)
	db.EmptyBlocks()
	for _, block := range newBlocks {
		persistBlock(block)
	}
}

func (b *blockchain) LockBlockchain() {
	b.m.Lock()
	defer b.m.Unlock()
}

func (b *blockchain) AddPeerBlock(block *Block) {
	b.m.Lock()
	defer b.m.Unlock()

	b.Height += 1
	b.CurrentDifficulty = block.Difficulty
	b.NewestHash = block.Hash

	persistBlockchain(b)
	persistBlock(block)

	// mempool
}
