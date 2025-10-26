package v3

import (
	"fmt"
	v1 "go-elevate/v1"
	v2 "go-elevate/v2"
	"math"
	"time"
)

type State string

var (
	STATE_IDLE    State = "IDLE"    // Not in movement
	STATE_PENDING State = "PENDING" // Moving to a pick up -- can receive more rides -- ignore for now, coming soon
	STATE_ACTIVE  State = "ACTIVE"  //
)

// Elevators handle moving to specified floors. Relatively dumb, don't handle any logic
// Besides moving to floors, and reporting back when finished.
type Elevator struct {
	*v2.Elevator
	idleElevators *ChanQueue[int] // Pass in ID, used for scheduler to assign work
	// Direction           int        // -1 - down, 0 Idle, 1 up
	State State
}

func newElevator(verbose bool, id int, idleChan *ChanQueue[int]) *Elevator {
	return &Elevator{
		idleElevators: idleChan,
		State:         STATE_IDLE,
		Elevator: &v2.Elevator{
			Route: make([]int, 0),
			Elevator: &v1.Elevator{
				ID:      1,
				Verbose: verbose,
				Floor:   0,
			},
		},
	}
}

func (e *Elevator) SetRoute(route []int) {
	e.Elevator.Route = route
}

// func (e *Elevator) Move(floor int) {
// Unused function -- to be used when elevators can pick up routes as moving... idea move to begin of pickup w/o changing state for scheduler checking
// 	dist := time.Duration(math.Abs(float64(floor - e.Floor)))
// 	time.Sleep(1 * dist * time.Second)
// 	fmt.Printf("Elevator %d reached floor %d\n", e.ID, floor)

// 	e.State = STATE_ACTIVE
// 	e.Go()
// }

// Go sends an eleavator on it's 'route', changing state to active and runng through all floors it is assigned.
// It can no longer pick up work.
func (e *Elevator) Go() {
	e.State = STATE_ACTIVE
	for _, v := range e.Route {
		dist := time.Duration(math.Abs(float64(v - e.Floor)))
		time.Sleep(1 * dist * time.Second)
		fmt.Printf("Elevator %d reached floor %d\n", e.ID, v)
		e.Floor = v
	}

	e.State = STATE_IDLE

	e.idleElevators.mut.Lock()
	// e.idleElevators <- e.ID
}
