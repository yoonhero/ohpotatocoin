package main

import (
	"log"
	"os"
	"strconv"

	"github.com/yoonhero/ohpotatocoin/db"
	"github.com/yoonhero/ohpotatocoin/rest"
	"github.com/yoonhero/ohpotatocoin/utils"
)

func main() {
	os.Setenv("PORT", "4000")
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}
	// close db to protect db file data
	defer db.Close()
	// rest or html server start
	// cli.Start()
	sv, err := strconv.Atoi(port)
	utils.HandleErr(err)
	rest.Start(sv)
}
