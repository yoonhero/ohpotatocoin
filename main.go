package main

import (
	"fmt"
	"time"
)

func countToTen(c chan<- int) {
	for i := range [10]int{} {
		fmt.Printf(">> sending %d << \n", i)
		c <- i
		fmt.Printf(">> sent %d << \n", i)
	}
	close(c)
}

func receive(c <-chan int) {
	for {
		time.Sleep(10 * time.Second)
		a, ok := <-c
		if !ok {
			break
		}
		fmt.Printf("|| received %d || \n", a)
	}
}

func main() {
	// // close db to protect db file data
	// defer db.Close()

	// // rest or html server start
	// cli.Start()
	c := make(chan int, 5)
	go countToTen(c)
	receive(c)

}
