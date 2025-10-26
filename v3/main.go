package v3

import (
	"fmt"
	"log"
	"strconv"
)

type RideRequest struct {
	From int
	To   int
}

type Input struct {
	Elevators int
	Floors    int
	People    int
	Verbose   bool
}

func Run(verbose bool, args ...string) {
	// validate inputs
	inputs, err := validate(args)
	if err != nil {
		log.Fatal(err)
	}

	// Build additonal workers + setup queues
	rideReqQueue := newChanQueue(inputs.Elevators*inputs.Elevators, RideRequest{})
	rideFactory := newRideFactory(RideFactoryConfig{
		tickerTime: 10,
		queue:      rideReqQueue,
	})

	for range inputs.People {
		rideFactory.NewRide()
	}

	scheduler := NewScheduler(inputs.Floors, inputs.Elevators, inputs.Verbose, rideReqQueue)

	go rideFactory.Serve()

	go scheduler.Run()

	// Generate people w/ random floors
	floorsReq := make([]RideRequest, 0)

	// Send all floorRequest to the scheduler
	var i int
	for i < len(floorsReq) {
		// scheduler.RequestQueue.mut.Lock()
		// defer scheduler.RequestQueue.mut.Unlock()
		scheduler.RequestQueue.queue <- floorsReq[i]
		i++
	}

	// Below is done in scheduler
	// Remove any dups

	// Build sorted order -- prioritize same from, direction
	// Build buckets
	// for {
	// 	r := rand.Reader // Not sure how this is used?

	// 	waitTime, _ := rand.Int(r, big.NewInt(10))
	// 	generateNewRequest(scheduler)
	// 	time.Sleep(time.Duration(waitTime.Int64()))
	// }
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
