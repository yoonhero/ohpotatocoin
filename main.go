package main

import "fmt"

type person struct {
	name string
	age  int
}

func (p person) sayHello() {
	fmt.Printf("Hello my name is %s and I'm %d", p.name, p.age-1)
}

func main() {
	yoonhero := person{name: "ysh", age: 12}
	yoonhero.sayHello()
}
