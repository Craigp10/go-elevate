package v3

import (
	"crypto/rand"
	"math/big"
	"time"
)

type RideFactoryConfig struct {
	tickerTime int
	queue      *ChanQueue[RideRequest]
}

type RideFactory struct {
	tickerTime int
	queue      *ChanQueue[RideRequest]
}

func newRideFactory(cfg RideFactoryConfig) RideFactory {
	return RideFactory{
		tickerTime: cfg.tickerTime,
		queue:      cfg.queue,
	}
}

func (rf *RideFactory) Serve() {

	for {
		time.Sleep(time.Second * time.Duration(rf.tickerTime))
		rf.NewRide()
	} //some ticker to create a new ride request over x time
}

func (rf *RideFactory) NewRide() {
	r := rand.Reader // Not sure how this is used?
	from, _ := rand.Int(r, big.NewInt(30))
	to, _ := rand.Int(r, big.NewInt(50))
	ride := RideRequest{
		From: int(from.Int64()),
		To:   int(to.Int64()),
	}
	rf.queue.mut.Lock()
	rf.queue.queue <- ride
	rf.queue.mut.Unlock()
}
