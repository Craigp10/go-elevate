package main

import (
	"flag"
	"fmt"
	v1 "go-elevate/v1"
	"os"
)

// func main() {
// 	// ele := Elevator{
// 	// 	Route: []int{1, 3, 5, 6, 1},
// 	// 	Floor: 0,
// 	// }

// 	// ele.Move()

// 	sch := New(10, 2)

// 	go sch.Run()

// 	for i := 0; i < 10; i++ {
// 		if i == 4 {
// 			continue
// 		}
// 		req := Request{
// 			Origin:      i + 1,
// 			Destination: 10 - i - 1,
// 		}
// 		sch.Queue <- req
// 		time.Sleep(2 * time.Second)

// 	}
// 	time.Sleep(20 * time.Second)
// }

// func RegisterElevators(numberOfElevators int) map[int]Elevator {
// 	var elevators map[int]Elevator
// 	for i := range elevators {
// 		e := Elevator{
// 			Id: i,
// 		}

// 		elevators[i] = e
// 	}

//		return elevators
//	}
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
		// TODO

	default:
		fmt.Printf("Unsupported version: %s\n", *version)
		os.Exit(1)
	}
}
