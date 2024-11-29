package v2

import (
	"fmt"
	"math"
	"time"
)

type State string

var (
	STATE_IDLE    State = "IDLE"    // Not in movement
	STATE_PENDING State = "PENDING" // Moving to a pick up -- can receive more rides
	STATE_ACTIVE  State = "ACTIVE"  //
)

// Elevators handle moving to specified floors. Relatively dumb, don't handle any logic
// Besides moving to floors, and reporting back when finished.
type Elevator struct {
	ID    int
	Floor int
	Route []int
	State State
	// Direction           int        // -1 - down, 0 Idle, 1 up
	ActiveElevators chan<- int // Pass in ID, used for scheduler to assign work
}

func (e *Elevator) SetRoute(route []int) {
	e.Route = route
}

func (e *Elevator) Move(floor int) {
	dist := time.Duration(math.Abs(float64(floor - e.Floor)))
	time.Sleep(1 * dist * time.Second)
	fmt.Printf("Elevator %d reached floor %d\n", e.ID, floor)

	e.State = STATE_ACTIVE
	e.Go()
}

func (e *Elevator) Go() {
	for _, v := range e.Route {
		dist := time.Duration(math.Abs(float64(v - e.Floor)))
		time.Sleep(1 * dist * time.Second)
		fmt.Printf("Elevator %d reached floor %d\n", e.ID, v)
		e.Floor = v
	}

	e.State = STATE_IDLE

	e.ActiveElevators <- e.ID // AvailableElevators instead?
}
