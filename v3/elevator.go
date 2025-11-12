package v3

import (
	"fmt"
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
	idleElevators chan string // Pass in ID, used for scheduler to assign work
	// Direction           int        // -1 - down, 0 Idle, 1 up
	State State
	ID    string
}

func newElevator(verbose bool, id string, idleChan chan string) *Elevator {
	return &Elevator{
		idleElevators: idleChan,
		State:         STATE_IDLE,
		Elevator:      v2.NewElevator(0, verbose),
		ID:            id,
	}
}

// Go sends an eleavator on it's 'route', changing state to active and runng through all floors it is assigned.
// It can no longer pick up work.
func (e *Elevator) Go() {
	fmt.Println("elevator going", e.Route())
	e.State = STATE_ACTIVE
	for _, v := range e.Route() {
		dist := time.Duration(math.Abs(float64(v - e.Floor)))
		time.Sleep(2 * dist * time.Second)
		fmt.Printf("Elevator %s reached floor %d\n", e.ID, v)
		e.Floor = v
	}

	e.State = STATE_IDLE

	e.idleElevators <- e.ID
}
