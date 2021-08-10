package main

import (
	"os"
	"strconv"

	"github.com/yoonhero/ohpotatocoin/db"
	"github.com/yoonhero/ohpotatocoin/rest"
	"github.com/yoonhero/ohpotatocoin/utils"
)

func main() {
	port := os.Getenv("PORT")
	// close db to protect db file data
	defer db.Close()
	// rest or html server start
	// cli.Start()
	sv, err := strconv.Atoi(port)
	utils.HandleErr(err)
	rest.Start(sv)
}
