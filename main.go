package main

import (
	"fmt"

	"github.com/yoonhero/ohpotatocoin/person"
)

func main() {
	yoonhero := person.Person{}
	yoonhero.SetDetails("nico", 12)
	yoonhero.Name()
	fmt.Println(yoonhero)
}
