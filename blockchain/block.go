package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/yoonhero/ohpotatocoin/db"
	"github.com/yoonhero/ohpotatocoin/utils"
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

// persist data
func (b *Block) persist() {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}

// restore data
func (b *Block) restore(data []byte) {
	utils.FromBytes(b, data)
}

var ErrNotFound = errors.New("Block not Found")

// find block by hash
func FindBlock(hash string) (*Block, error) {
	blockBytes := db.Block(hash)

	// if that block don't exist
	if blockBytes == nil {
		// return nil with error
		return nil, ErrNotFound
	}

	block := &Block{}
	// restore the block data
	block.restore(blockBytes)

	return block, nil
}

// create block
func createBlock(data string, prevHash string, height int) *Block {
	block := &Block{
		Data:     data,
		Hash:     "",
		PrevHash: prevHash,
		Height:   height,
	}
	payload := block.Data + block.PrevHash + fmt.Sprint(block.Height)

	// hashing the payload
	block.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(payload)))

	// persist the block
	block.persist()

	return block
}
