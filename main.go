package main

import (
	"github.com/yoonhero/ohpotatocoin/explorer"
	"github.com/yoonhero/ohpotatocoin/rest"
)

func main() {
	go explorer.Start(5000)
	rest.Start(4000)
}
