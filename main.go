package main

import "fmt"

func main() {
	name := "HI my name is ysh"
	for _, letter := range name {
		fmt.Printf("%x\n", letter)
	}
}
