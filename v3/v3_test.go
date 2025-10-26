package v3

import (
	"testing"

	"github.com/stretchr/testify/require" // need to go get
)

func Test_Scheduler_Constructor(t *testing.T) {

	rideRequestChan := newChanQueue(10, RideRequest{})
	sch := NewScheduler(10, 2, true, rideRequestChan)
	require.NotNil(sch)

	require.NotNil(t, sch.RequestQueue)
	require.Len(t, 2, sch.ElevatorsRegistry)
	require.Equal(t, 10, sch.Floors)
	require.NotNil(t, sch.AvailableElevators)
	require.NotNil(t, sch.ActiveElevators)
	require.NotNil(t, sch.IdleElevator)
	require.NotNil(t, sch.idleChan)
	require.NotNil(t, sch.upRides)
	require.NotNil(t, sch.downRides)
	require.NotNil(t, sch.upQueue)
	require.NotNil(t, sch.downQueue)
}
