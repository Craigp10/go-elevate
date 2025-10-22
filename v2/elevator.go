package v2

import (
	"crypto/rand"
	"fmt"
	v1 "go-elevate/v1"
	"log"
	"math/big"
	"slices"
	"strconv"
	"sync"
)

type Input struct {
	Elevators int
	Floors    int
	People    int
	Verbose   bool
}

// Extend elevator from v1
type Elevator struct {
	*v1.Elevator
	toFloors []int
}

func NewElevator(id int, verbose bool) *Elevator {
	return &Elevator{
		toFloors: make([]int, 0),
		Elevator: &v1.Elevator{
			Floor:   0,
			ID:      int(id),
			Verbose: verbose,
		},
	}
}

func (e *Elevator) ToFloors(newFloor ...int) {
	e.toFloors = append(e.toFloors, newFloor...)
}

func (e *Elevator) Go(wg *sync.WaitGroup) {
	slices.SortFunc(e.toFloors, func(a, b int) int {
		return a - b
	})

	for _, floor := range e.toFloors {
		e.Elevator.Go(floor) // Floor is reset each time... Doesn't count for new floor.
		fmt.Printf("Elevator %d has reached floor %d -- Dropping off people\n", e.ID, floor)
	}
	wg.Done()
}

func Run(verbose bool, args ...string) {
	Inputs, err := validate(args)
	if err != nil {
		log.Fatal(err)
	}

	// Generate elevators
	elevators := make([]*Elevator, Inputs.Elevators)
	var i int
	for range elevators {
		elevators[i] = NewElevator(i+1, verbose)
		i++
	}
	// Generate people w/ random floors
	floors := make([]int, 0)
	for range Inputs.People {
		r := rand.Reader
		inte, _ := rand.Int(r, big.NewInt(int64(Inputs.Floors)))
		floors = append(floors, int(inte.Int64()))
	}

	slices.SortFunc(floors, func(a, b int) int {
		return a - b
	})

	// Hash people by groups -- simplifying to just sort and insert
	mappedPeople := make([][]int, len(elevators))
	bucket := 0
	for i := range floors {
		if i != 0 && i%5 == 0 {
			bucket++
		}
		mappedPeople[bucket] = append(mappedPeople[bucket], floors[i])
		i++
	}

	wg := sync.WaitGroup{}

	for j, ele := range elevators {
		wg.Add(1)
		go func() {
			ele.ToFloors(mappedPeople[j]...)
			fmt.Println("elevator running: ", ele.ID)
			ele.Go(&wg) // Need to now handle a start floor.... Does go allow for forget whats its call... argument raising, changing signature.
		}()
	}
	wg.Wait()
}

func validate(args []string) (*Input, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("invalid # of arguments: %d", len(args))
	}

	floors, err := strconv.Atoi(args[0])
	if err != nil || floors > 30 {
		return nil, fmt.Errorf("invalid floors input floors: %d, error: <D-s>'%v'", floors, err)
	}

	elevators, err := strconv.Atoi(args[1])
	if err != nil || elevators > int(floors/4) {
		return nil, fmt.Errorf("invalid elevator input elevators: %d, error: '%v'", elevators, err)
	}

	people, err := strconv.Atoi(args[2])
	if err != nil || people > 5*elevators {
		return nil, fmt.Errorf("invalid people input people: %d, error: '%v'", people, err)
	}

	return &Input{
		Elevators: elevators,
		Floors:    floors,
		People:    people,
	}, nil
}
