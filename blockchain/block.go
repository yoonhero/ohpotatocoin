package blockchain

import (
	"crypto/sha256"
	"fmt"
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

func createBlock(data string, prevHash string, height int) *Block {
	block := &Block{
		Data:     data,
		Hash:     "",
		PrevHash: prevHash,
		Height:   height,
	}
	payload := block.Data + block.PrevHash + fmt.Sprint(block.Height)
	block.Hash = fmt.Sprintf("%s", sha256.Sum256([]byte(payload)))
	return block
}
