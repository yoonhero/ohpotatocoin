package cli

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/yoonhero/ohpotatocoin/explorer"
	"github.com/yoonhero/ohpotatocoin/rest"
)

// type portsFlags []string

// func (i *portsFlags) Set(value string) error {
// 	*i = append(*i, value)
// 	return nil
// }

// Define a type named "intslice" as a slice of ints
type intslice []int

// Now, for our new type, implement the two methods of
// the flag.Value interface...
// The first method is String() string
func (i *intslice) String() string {
	return fmt.Sprintf("%d", *i)
}

// The second method is Set(value string) error
func (i *intslice) Set(value string) error {
	tmp, err := strconv.Atoi(value)
	if err != nil {
		*i = append(*i, -1)
	} else {
		*i = append(*i, tmp)
	}
	return nil
}

var myints intslice

// introduce usages
func usage() {
	fmt.Printf("Welcome to 오감자 코인 \n")
	fmt.Printf("Please use the following commands: \n\n")
	fmt.Printf("-port=4000:           Start the PORT of the server\n")
	fmt.Printf("-mode=rest:           Choose between 'html' or 'rest' \n\n")
	// fmt.Printf("-port=3000 -port=4000:   if run both mode set int list")
	// exit program not to occur err
	runtime.Goexit()
}

func Start() {

	// if don't type like rest or explorer
	if len(os.Args) < 2 {
		usage()
	}
	// flag
	// defualt val is 4000
	// get val in -port
	// on err -> print Set port of the server
	port := flag.Int("port", 4000, "Set port of the server")

	// default val is rest
	// get val in -mode
	mode := flag.String("mode", "rest", "Choose between 'html' and 'rest'")

	// flag.Var(&myints, "port", "List of integers")

	flag.Parse()

	switch *mode {
	case "rest":
		// start rest api
		rest.Start(*port)
	case "html":
		// start html explorer
		explorer.Start(*port)
	default:
		usage()
	}

}
