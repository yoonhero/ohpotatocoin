package main

import (
	"fmt"
	"time"
)

func countToTen(c chan int) {
	for i := range [10]int{} {
		c <- i
		time.Sleep(1 * time.Second)
	}
}

func main() {
	// // close db to protect db file data
	// defer db.Close()

	// // rest or html server start
	// cli.Start()
	c := make(chan int)
	go countToTen(c)
	for {
		a := <-c
		fmt.Println(a)
	}

}
