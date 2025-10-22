package v1

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

type Elevator struct {
	Verbose bool
	ID      int
	Floor   int
}

func NewElevator(floor int) *Elevator {
	return &Elevator{
		Floor: floor,
	}
}

type Inputs struct {
	ElevatorFloor int
	ToFloor       int
	Verbose       bool
}

func (e *Elevator) Go(toFloor int) {
	dist := int(math.Abs(float64(e.Floor - toFloor)))
	direction := 1
	if e.Floor > toFloor {
		direction = -1
	}

	i := 0
	for i < dist {
		chg := i * direction
		if e.Verbose {
			fmt.Printf("Elevator %d on floor %d\n", e.ID, chg+e.Floor)
		}
		time.Sleep(1 * time.Second)
		i++
	}
	e.Floor = toFloor
}

func Run(verbose bool, args ...string) {
	// Assert our arguments as the correct inputs
	// inputs := args.(Inputs)
	elevatorFloor, err1 := strconv.Atoi(args[0])
	toFloor, err2 := strconv.Atoi(args[1])
	// verbose := args[2].(bool)
	if err1 != nil || err2 != nil {
		fmt.Println("Invalid inputs, not convertable to int")
	}

	inputs := &Inputs{
		ElevatorFloor: elevatorFloor,
		ToFloor:       toFloor,
		Verbose:       verbose,
	}

	e := &Elevator{
		Floor:   inputs.ElevatorFloor,
		Verbose: inputs.Verbose,
	}

	e.Go(inputs.ToFloor)
}
