// persistence of block
// connected to DB to save data
// using bolt DB (bitcoin levelDB)

package blockchain

import (
	"fmt"
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

// persist the blockchain data
func (b *blockchain) persist() {
	db.SaveCheckpoint(utils.ToBytes(b))
}

// add block to blockchain
func (b *blockchain) AddBlock() {
	// createBlock
	block := createBlock(b.NewestHash, b.Height+1)

	// set newesthash new block's hash
	b.NewestHash = block.Hash
	// set height new block's height
	b.Height = block.Height

	b.CurrentDifficulty = block.Difficulty

	// persist the blockchain
	b.persist()
}

// return all blocks
func (b *blockchain) Blocks() []*Block {
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

// recalculate difficulty of block by timestamp
func (b *blockchain) recalculateDifficulty() int {
	// get all blocks
	allBlocks := b.Blocks()
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

func (b *blockchain) difficulty() int {
	// if genesis block or not
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		// recalculate the difficulty
		return b.recalculateDifficulty()
	} else {
		return b.CurrentDifficulty
	}
}

// unspent transaction out by address
func (b *blockchain) UTxOutsByAddress(address string) []*UTxOut {
	var uTxOuts []*UTxOut
	creatorTxs := make(map[string]bool)
	for _, block := range b.Blocks() {
		for _, tx := range block.Transactions {
			for _, input := range tx.TxIns {
				fmt.Println(input)
				if input.Owner == address {
					creatorTxs[input.TxID] = true
				}
			}
			for index, output := range tx.TxOuts {
				fmt.Println(output)
				if output.Owner == address {
					if _, ok := creatorTxs[tx.ID]; !ok {
						uTxOut := &UTxOut{tx.ID, index, output.Amount}
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

// get balance
func (b *blockchain) BalancByAddress(address string) int {
	txOuts := b.UTxOutsByAddress(address)
	var amount int
	for _, txOut := range txOuts {
		amount += txOut.Amount
	}
	return amount
}

func Blockchain() *blockchain {
	// if var blockchain is nil
	// add first block
	if b == nil {
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
	}
	// return type blockchain struct
	return b
}
