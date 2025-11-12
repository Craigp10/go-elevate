package v3

import (
	"context"
	"fmt"
	"log"
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
	Verbose   bool
}

func Run(verbose bool, args ...string) {
	ctx := context.Background()

	// validate inputs
	inputs, err := validate(args)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("validated inputs: %+v\n", inputs)
	// Build additonal workers + setup queues
	rideReqQueue := make(chan RideRequest, inputs.People*inputs.People)
	rideFactory := newRideFactory(RideFactoryConfig{
		tickerTime: 10,
		queue:      rideReqQueue,
		verbose:    verbose,
	})

	fmt.Println("created queues")

	for i := 0; i < inputs.People; i++ {
		// blocking to queue up rides
		rideFactory.NewRide(int64(inputs.Floors))
	}

	fmt.Println("created rides: ", len(rideReqQueue))
	scheduler := NewScheduler(inputs.Floors, inputs.Elevators, inputs.Verbose, rideReqQueue)
	go scheduler.Run(inputs.Elevators)
	go rideFactory.Serve(ctx, int64(inputs.Floors))

	time.Sleep(2 * time.Second)
	// TODO -- add closer listener instead of loop
	for {
		time.Sleep(2 * time.Second)
		fmt.Printf("available elevators %d\n", len(scheduler.AvailableElevators))
	}
	scheduler.Close()

	close(rideFactory.queue)
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
