package v3

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

type RideFactoryConfig struct {
	verbose    bool
	tickerTime int
	queue      *ChanQueue[RideRequest]
}

type RideFactory struct {
	verbose    bool
	tickerTime int
	queue      *ChanQueue[RideRequest]
	cancel     context.CancelFunc
}

func newRideFactory(cfg RideFactoryConfig) RideFactory {
	return RideFactory{
		tickerTime: cfg.tickerTime,
		queue:      cfg.queue,
		verbose:    cfg.verbose,
	}
}

func (rf *RideFactory) Serve(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	rf.cancel = cancel

	ticker := time.NewTicker(time.Duration(rf.tickerTime) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("RideFactory shutting down gracefully...")
			return
		case <-ticker.C:
			rf.NewRide()
		}
	}
}

func (rf *RideFactory) Close() {
	if rf.cancel != nil {
		rf.cancel()
		return
	}

	fmt.Println("Nothing to close...")
}

func (rf *RideFactory) NewRide() {
	if rf.verbose {
		fmt.Println("new ride")
	}

	r := rand.Reader // Not sure how this is used?
	from, _ := rand.Int(r, big.NewInt(30))
	to, _ := rand.Int(r, big.NewInt(50))
	ride := RideRequest{
		From: int(from.Int64()),
		To:   int(to.Int64()),
	}
	rf.queue.mut.Lock()
	defer rf.queue.mut.Unlock()
	rf.queue.queue <- ride
}
