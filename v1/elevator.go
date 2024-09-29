package v1

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

type Elevator struct {
	Floor int
}

type Inputs struct {
	ElevatorFloor int
	ToFloor       int
}

func (e *Elevator) Go(toFloor int) {
	dist := int(math.Abs(float64(e.Floor - toFloor)))
	i := 1
	for i < dist+1 {
		fmt.Printf("On floor %d\n", i+e.Floor)
		time.Sleep(1 * time.Second)
		i++
	}
	e.Floor = toFloor
	fmt.Println("Elevator has Reached floor!")
}

func Run(args ...interface{}) {
	// Assert our arguments as the correct inputs
	// inputs := args.(Inputs)
	elevatorFloor, err1 := strconv.Atoi(args[0].(string))
	toFloor, err2 := strconv.Atoi(args[1].(string))
	if err1 != nil || err2 != nil {
		fmt.Println("Invalid inputs, not convertable to int")
	}

	inputs := &Inputs{
		ElevatorFloor: elevatorFloor,
		ToFloor:       toFloor,
	}

	e := &Elevator{
		Floor: inputs.ElevatorFloor,
	}

	e.Go(inputs.ToFloor)
}
