package main

import "fmt"

func plus(a, b int, name string) (int, string) {
	return a + b, name
}

func main() {
	//var name string = "hello"
	result, name := plus(2, 2, "nico")
	fmt.Println(result, name)
}
