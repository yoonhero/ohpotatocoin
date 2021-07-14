package main

import "fmt"

func main() {
	foods := [3]string{"potato", "pizza", "pasta"}
	for i := 0; i < len(foods); i++ {
		fmt.Println(foods[i])
	}
}
