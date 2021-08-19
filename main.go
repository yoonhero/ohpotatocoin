package main

import (
	"os"
	"strconv"

	"github.com/yoonhero/ohpotatocoin/db"
	"github.com/yoonhero/ohpotatocoin/rest"
	"github.com/yoonhero/ohpotatocoin/utils"
)

func main() {
	// os.Setenv("PORT", "4000")
	port := os.Getenv("PORT")
	// close db to protect db file data
	defer db.CloseSqlDB()
	db.InitPostgresDB()

	sv, err := strconv.Atoi(port)
	utils.HandleErr(err)
	rest.Start(sv)
}
