package main

import (
	"github.com/yoonhero/ohpotatocoin/db"
	"github.com/yoonhero/ohpotatocoin/rest"
)

func main() {
	// close db to protect db file data
	defer db.Close()
	// rest or html server start
	// cli.Start()
	rest.Start(4000)
}
