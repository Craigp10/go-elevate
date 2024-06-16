package main

import (
	"fmt"
	"math"
	"time"
)

// Elevators handle moving to specified floors. Relatively dumb, don't handle any logic
// Besides moving to floors, and reporting back when finished.
type Elevator struct {
	Id              int
	Floor           int
	Route           []int
	Busy            bool
	ActiveElevators chan<- int
}

func (e *Elevator) SetRoute(route []int) {
	e.Route = route
}

func (e *Elevator) Move() {

	for _, v := range e.Route {
		dist := time.Duration(math.Abs(float64(v - e.Floor)))
		time.Sleep(1 * dist * time.Second)
		fmt.Printf("Elevator %d reached floor %d\n", e.Id, v)
		e.Floor = v
	}

	e.Busy = false
	e.ActiveElevators <- e.Id
}

// ValidateRoute validates the route for the elevator to path.
// Atm it only checks length of route and removes if first floor is current floor.
func (e *Elevator) ValidateRoute() bool {

	if len(e.Route) == 0 {
		return true
	}

	if e.Route[0] == e.Floor {
		e.Route = e.Route[1:]
	}

	return true
}
