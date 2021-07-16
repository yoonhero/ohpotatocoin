package main

import "github.com/yoonhero/ohpotatocoin/blockchain"

func main() {
	blockchain.Blockchain().AddBlock("First")
	blockchain.Blockchain().AddBlock("Second")
	blockchain.Blockchain().AddBlock("Third")
}
