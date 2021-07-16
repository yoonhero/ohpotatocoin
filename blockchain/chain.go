// persistence of block
// connected to DB to save data
// using bolt DB (bitcoin levelDB)

package blockchain

import (
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
	CurrentDifficulty int    `json:"currentDifficulty`
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

func (b *blockchain) recalculateDifficulty() int {
	allBlocks := b.Blocks()
	newestBlock := allBlocks[0]
	lastRecalculatedBlock := allBlocks[difficultyInterval-1]
	actualTime := (newestBlock.Timestamp / 60) - (lastRecalculatedBlock.Timestamp / 60)
	expectedTime := difficultyInterval * blockInterval
	if actualTime <= (expectedTime - allowedRange) {
		return b.CurrentDifficulty + 1
	} else if actualTime >= (expectedTime + allowedRange) {
		return b.CurrentDifficulty - 1
	}
	return b.CurrentDifficulty
}

func (b *blockchain) difficulty() int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		// recalculate the difficulty
		return b.recalculateDifficulty()
	} else {
		return b.CurrentDifficulty
	}
}

func Blockchain() *blockchain {
	// if var blockchain is nil
	// add first block
	if b == nil {
		// do only once a time
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

// func (b *Block) calculateHash() {
// 	// create hash
// 	hash := sha256.Sum256([]byte(b.Data + b.PrevHash))

// 	// change block hash by pointer
// 	// Sprintf("%x", hash) is format hash
// 	b.Hash = fmt.Sprintf("%x", hash)
// }

// func getLastHash() string {
// 	// get len(blocks)
// 	totalBlocks := len(GetBlockchain().blocks)
// 	// if blocks are empty
// 	// return nothing
// 	if totalBlocks == 0 {
// 		return ""
// 	}

// 	// or return last block
// 	return GetBlockchain().blocks[totalBlocks-1].Hash
// }

// func createBlock(data string) *Block {
// 	// variable newBlock
// 	// {data:data, hash:"", prevHash: getLastHash()}
// 	newBlock := Block{data, "", getLastHash(), len(GetBlockchain().blocks) + 1}
// 	// newblock calculate hash
// 	newBlock.calculateHash()

// 	// return new block
// 	return &newBlock
// }

// func (b *blockchain) AddBlock(data string) {
// 	// add block to blockchain.blocks slice
// 	b.blocks = append(b.blocks, createBlock(data))
// }

// func (b *blockchain) AllBlocks() []*Block {
// 	// return all blocks
// 	return GetBlockchain().blocks
// }

// // make new error
// var ErrNotFound = errors.New("Block Not Found")

// // find block by height (id)
// func (b *blockchain) GetBlock(height int) (*Block, error) {
// 	// if block not exist
// 	if height > len(b.blocks) || height == 0 {
// 		// return nil and error
// 		return nil, ErrNotFound
// 	}
// 	return b.blocks[height-1], nil
// }
