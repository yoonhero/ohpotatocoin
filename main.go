package main

import (
	"github.com/yoonhero/ohpotatocoin/db"
	"github.com/yoonhero/ohpotatocoin/rest"
)

func main() {
	// port := os.Getenv("PORT")
	// close db to protect db file data
	defer db.CloseSqlDB()
	db.InitPostgresDB()

	// sv, err := strconv.Atoi(port)
	// utils.HandleErr(err)
	rest.Start(4000)
}
