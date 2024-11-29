package v2

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"time"
)

type RideRequest struct {
	From int
	To   int
}

type Input struct {
	Elevators int
	Floors    int
	People    int
}

func Run(args ...string) {
	Inputs, err := validate(args)
	if err != nil {
		log.Fatal(err)
	}

	scheduler := NewScheduler(Inputs.Floors, Inputs.Elevators)

	go scheduler.Run()

	// Generate people w/ random floors
	floorsReq := make([]RideRequest, 0)
	for range Inputs.People {
		r := rand.Reader // Not sure how this is used?

		from, _ := rand.Int(r, big.NewInt(30))
		to, _ := rand.Int(r, big.NewInt(50))

		floorsReq = append(floorsReq, RideRequest{
			From: int(from.Int64()),
			To:   int(to.Int64()),
		})
	}

	// Send all floorRequest to the scheduler
	var i int
	for i < len(floorsReq) {
		scheduler.RequestQueue <- floorsReq[i]
		i++
	}

	// Below is done in scheduler
	// Remove any dups

	// Build sorted order -- prioritize same from, direction
	// Build buckets
	for {
		r := rand.Reader // Not sure how this is used?

		waitTime, _ := rand.Int(r, big.NewInt(10))
		generateNewRequest(scheduler)
		time.Sleep(time.Duration(waitTime.Int64()))
	}
}

func validate(args []string) (*Input, error) {
	floors, err := strconv.Atoi(args[0])
	if err != nil || floors > 50 || floors < 1 {
		return nil, fmt.Errorf("invalid floors input '%v'", err)
	}

	elevators, err := strconv.Atoi(args[1])
	if err != nil || elevators > int(floors/4) || elevators < 1 {
		return nil, fmt.Errorf("invalid elevator input '%v'", err)
	}

	people, err := strconv.Atoi(args[2])
	if err != nil || people < 1 {
		return nil, fmt.Errorf("invalid people input '%v'", err)
	}

	return &Input{
		Elevators: elevators,
		Floors:    floors,
		People:    people,
	}, nil
}

func generateNewRequest(s Scheduler) {
	r := rand.Reader // Not sure how this is used?
	from, _ := rand.Int(r, big.NewInt(30))
	to, _ := rand.Int(r, big.NewInt(50))
	req := RideRequest{
		From: int(from.Int64()),
		To:   int(to.Int64()),
	}
	s.RequestQueue <- req
}
