package main

import (
	"fmt"
	"time"
)

func countToTen(c chan<- int) {
	for i := range [10]int{} {
		time.Sleep(1 * time.Second)
		c <- i
	}
	close(c)
}

func receive(c <-chan int) {
	for {
		a, ok := <-c
		if !ok {
			break
		}
		fmt.Println(a)
	}
}

func main() {
	// // close db to protect db file data
	// defer db.Close()

	// // rest or html server start
	// cli.Start()
	c := make(chan int)
	go countToTen(c)
	receive(c)

}
