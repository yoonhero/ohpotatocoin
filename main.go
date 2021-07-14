package main

import "fmt"

func plus(a, b int, name string) (int, string) {
	return a + b, name
}

func plusAll(a ...int) int {
	var total int
	for _, item := range a {
		total += item
	}
	return total
}

func main() {
	//var name string = "hello"
	result := plusAll(2, 3, 4, 5, 6, 7, 8, 9)
	fmt.Println(result)
}
