package main

import "time"

func main() {
	// ele := Elevator{
	// 	Route: []int{1, 3, 5, 6, 1},
	// 	Floor: 0,
	// }

	// ele.Move()

	sch := New(10, 2)

	go sch.Run()

	for i := 0; i < 10; i++ {
		if i == 4 {
			continue
		}
		req := Request{
			Origin:      i + 1,
			Destination: 10 - i - 1,
		}
		sch.Queue <- req
		time.Sleep(2 * time.Second)

	}
	time.Sleep(20 * time.Second)
}

func RegisterElevators(numberOfElevators int) map[int]Elevator {
	var elevators map[int]Elevator
	for i := range elevators {
		e := Elevator{
			Id: i,
		}

		elevators[i] = e
	}

	return elevators
}
