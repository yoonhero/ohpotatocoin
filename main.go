package main

import (
	"github.com/yoonhero/ohpotatocoin/cli"
	"github.com/yoonhero/ohpotatocoin/db"
)

func main() {
	defer db.Close()
	cli.Start()
}
