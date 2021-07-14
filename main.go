package main

import "fmt"

func main() {
	foods := []string{"potato", "pizza", "pasta"}
	fmt.Printf("%v\n", foods)
	foods = append(foods, "tomato")
	fmt.Printf("%v\n", foods)
}
