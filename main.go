package main

import (
	"github.com/yoonhero/ohpotatocoin/cli"
	"github.com/yoonhero/ohpotatocoin/db"
)

func main() {
	// port := os.Getenv("PORT")
	// // close db to protect db file data
	// defer db.Close()
	// // rest or html server start
	// // cli.Start()
	// sv, err := strconv.Atoi(port)
	// utils.HandleErr(err)
	// rest.Start(sv)
	// db.Start()

	defer db.Close()
	db.InitDB()
	cli.Start()
}
