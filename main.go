package main

import "time"

func main() {
	// ele := Elevator{
	// 	Route: []int{1, 3, 5, 6, 1},
	// 	Floor: 0,
	// }

	// ele.Move()

	sch := Scheduler{
		Elevators: RegisterElevators(2),
		Floors:    10,
		Queue:     make(chan Request),
	}

	// sch.Queue <- 4
	// sch.Queue <- 3
	// sch.Queue <- 2
	// sch.Queue <- 6
	// sch.Queue <- 8
	// sch.Queue <- 5
	// sch.Queue <- 10
	go sch.Run()

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
