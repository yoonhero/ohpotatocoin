package main

import (
	"fmt"
	"time"
)

func countToTen(name string) {
	for i := range [10]int{} {
		fmt.Println(i)
		time.Sleep(1 * time.Second)
	}
}

func main() {
	// // close db to protect db file data
	// defer db.Close()

	// // rest or html server start
	// cli.Start()
	go countToTen("f")
	go countToTen("s")
}
