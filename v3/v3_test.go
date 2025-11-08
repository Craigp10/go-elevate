package v3

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require" // need to go get
)

func Test_Scheduler_Constructor(t *testing.T) {

	rideRequestChan := newChanQueue(10, RideRequest{})
	sch := NewScheduler(10, 2, true, rideRequestChan)
	require.NotNil(t, sch)

	require.NotNil(t, sch.RequestQueue)
	require.Len(t, sch.ElevatorsRegistry, 2)
	require.Equal(t, 10, sch.Floors)
	require.NotNil(t, sch.AvailableElevators)
	require.NotNil(t, sch.ActiveElevators)
	require.NotNil(t, sch.upRides)
	require.NotNil(t, sch.downRides)
	require.NotNil(t, sch.upQueue)
	require.NotNil(t, sch.downQueue)
}

func Test_rideFactory(t *testing.T) {
	ctx := context.Background()
	rideQueue := newChanQueue(100, RideRequest{})
	cfg := RideFactoryConfig{
		tickerTime: 5,
		queue:      rideQueue,
	}
	newFactory := newRideFactory(cfg)

	require.NotNil(t, newFactory)

	newFactory.NewRide()
	go newFactory.Serve(ctx)
	time.Sleep(6 * time.Second)

	require.Equal(t, 2, rideQueue.Length())
}
