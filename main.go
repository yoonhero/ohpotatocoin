package main

import (
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

	switch os.Args[1] {
	// if os.Args[1] == 'explorer'
	case "explorer":
		fmt.Println("Start Explorer")

	// if os.Args[1] == "rest"
	case "rest":
		fmt.Println("Start REST API")
	default:
		usage()
	}
}
