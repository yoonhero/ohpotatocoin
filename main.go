package main

import "fmt"

func main() {
	a := 2
	b := &a
	fmt.Println(*b, &a)
	fmt.Scanf("%d", &a)
	fmt.Println(a)
}
