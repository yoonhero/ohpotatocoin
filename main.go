package main

import (
	"flag"
	"fmt"
	"os"
)

// introduce usages
func usage() {
	fmt.Printf("Welcome to 오감자 코인 \n")
	fmt.Printf("Please use the following commands: \n\n")
	fmt.Printf("explorer:    Start the HTML Explorer\n")
	fmt.Printf("rest:        Start the Rest API (recommended)\n")
	// exit program not to occur err
	os.Exit(0)
}

func main() {

	// if don't type like rest or explorer
	if len(os.Args) < 2 {
		usage()
	}

	// variable flag set
	rest := flag.NewFlagSet("rest", flag.ExitOnError)

	// variable which change -port to int var
	// if err occus -> "Sets the port of the server"
	portFlag := rest.Int("port", 4000, "Sets the port of the server")

	switch os.Args[1] {
	// if os.Args[1] == 'explorer'
	case "explorer":
		fmt.Println("Start Explorer")

	// if os.Args[1] == "rest"
	case "rest":
		rest.Parse(os.Args[2:])

	default:
		usage()
	}

	// if parsed(boolean)
	if rest.Parsed() {
		fmt.Println(*portFlag)
		fmt.Println("Start server")
	}
}
