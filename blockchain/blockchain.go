package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"sync"
)

// type block
// data is data for block
// hash is sha256.Sum256([]byte(Data+PrevHash))
// prevHash is previous block's hash
// Height is id of block
type Block struct {
	Data     string `json:"data"`
	Hash     string `json:"hash"`
	PrevHash string `json:"prevHash,omitempty"`
	Height   int    `json:"height"`
}

// type blockchain
// blocks is slice of []Block
type blockchain struct {
	blocks []*Block
}

// variable blockchain pointers
var b *blockchain

// variable struct that play func only one time
var once sync.Once

func (b *Block) calculateHash() {
	// create hash
	hash := sha256.Sum256([]byte(b.Data + b.PrevHash))

	// change block hash by pointer
	// Sprintf("%x", hash) is format hash
	b.Hash = fmt.Sprintf("%x", hash)
}

func getLastHash() string {
	// get len(blocks)
	totalBlocks := len(GetBlockchain().blocks)
	// if blocks are empty
	// return nothing
	if totalBlocks == 0 {
		return ""
	}

	// or return last block
	return GetBlockchain().blocks[totalBlocks-1].Hash
}

func createBlock(data string) *Block {
	// variable newBlock
	// {data:data, hash:"", prevHash: getLastHash()}
	newBlock := Block{data, "", getLastHash(), len(GetBlockchain().blocks) + 1}
	// newblock calculate hash
	newBlock.calculateHash()

	// return new block
	return &newBlock
}

func (b *blockchain) AddBlock(data string) {
	// add block to blockchain.blocks slice
	b.blocks = append(b.blocks, createBlock(data))
}

func GetBlockchain() *blockchain {
	// if var blockchain is nil
	// add first block
	if b == nil {
		// do only once a time
		once.Do(func() {
			b = &blockchain{}
			b.AddBlock("Genesis")
		})
	}

	// return type blockchain struct
	return b
}

func (b *blockchain) AllBlocks() []*Block {
	// return all blocks
	return GetBlockchain().blocks
}

// make new error
var ErrNotFound = errors.New("Block Not Found")

// find block by height (id)
func (b *blockchain) GetBlock(height int) (*Block, error) {
	// if block not exist
	if height > len(b.blocks) || height == 0 {
		// return nil and error
		return nil, ErrNotFound
	}
	return b.blocks[height-1], nil
}
