package main

import (
	"flag"
	"fmt"
	v1 "go-elevate/v1"
	v2 "go-elevate/v2"
	"os"
)

type Elevator interface {
	Go(int)
}

type Runner interface {
	Run(args ...interface{}) error
}

func main() {

	version := flag.String("version", "v1", "Specify the version (v1, v2, etc.)")

	flag.Parse()

	args := flag.Args()

	switch *version {
	case "v1":
		if len(args) != 2 {
			fmt.Println("Version v1 requires exactly 2 arguments.")
			fmt.Println("Usage: --version=v1 <arg1> <arg2>")
			os.Exit(1)
		}
		arg1, arg2 := args[0], args[1]
		// runner = Version1{}

		v1.Run(arg1, arg2)

	case "v2":
		if len(args) != 3 {
			fmt.Println("Version v2 requires exactly 3 arguments.")
			fmt.Println("Usage: --version=v2 <Floors int> <Elevators int> <People int>")
			os.Exit(1)
		}

		v2.Run(args...)

	default:
		fmt.Printf("Unsupported version: %s\n", *version)
		os.Exit(1)
	}
}
