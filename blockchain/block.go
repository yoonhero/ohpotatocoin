package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"

	"github.com/yoonhero/ohpotatocoin/db"
	"github.com/yoonhero/ohpotatocoin/utils"
)

const difficulty int = 2

// type block
// data is data for block
// hash is sha256.Sum256([]byte(Data+PrevHash))
// prevHash is previous block's hash
// Height is id of block
type Block struct {
	Data       string `json:"data"`
	Hash       string `json:"hash"`
	PrevHash   string `json:"prevHash,omitempty"`
	Height     int    `json:"height"`
	Difficulty int    `json:"difficulty"`
	Nonce      int    `json:"none"`
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

func (b *Block) mine() {
	target := strings.Repeat("0", b.Difficulty)
	for {
		blockAsString := fmt.Sprint(b)
		hash := fmt.Sprintf("%x", sha256.Sum256([]byte(blockAsString)))
		fmt.Printf("Block as String:%s\nHash:%s\nTarget:%s\nNonce:%d\n\n\n", blockAsString, hash, target, b.Nonce)
		if strings.HasPrefix(hash, target) {
			b.Hash = hash
			break
		} else {
			b.Nonce++
		}
	}
}

// create block
func createBlock(data string, prevHash string, height int) *Block {
	block := &Block{
		Data:       data,
		Hash:       "",
		PrevHash:   prevHash,
		Height:     height,
		Difficulty: difficulty,
		Nonce:      0,
	}

	// mining the block
	block.mine()

	// persist the block
	block.persist()

	return block
}
