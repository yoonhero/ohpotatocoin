package main

import (
	"github.com/yoonhero/ohpotatocoin/cli"
	"github.com/yoonhero/ohpotatocoin/db"
)

func main() {
	// close db to protect db file data
	defer db.Close()

	// rest or html server start
	cli.Start()
}
