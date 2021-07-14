package main

import (
	"fmt"

	"github.com/yoonhero/ohpotatocoin/blockchain"
)

func main() {
	chain := blockchain.GetBlockchain()
	fmt.Println(chain)
}
