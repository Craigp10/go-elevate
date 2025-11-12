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
	queue      chan RideRequest
}

type RideFactory struct {
	verbose    bool
	tickerTime int
	queue      chan RideRequest
	cancel     context.CancelFunc
}

func newRideFactory(cfg RideFactoryConfig) RideFactory {
	return RideFactory{
		tickerTime: cfg.tickerTime,
		queue:      cfg.queue,
		verbose:    cfg.verbose,
	}
}

func (rf *RideFactory) Serve(ctx context.Context, floorCount int64) {
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
			rf.NewRide(floorCount)
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

func (rf *RideFactory) NewRide(floorCount int64) {

	r := rand.Reader // Not sure how this is used?
	// from, to := big.NewInt(1), big.NewInt(1)
	// for from.Cmp(to) != 0 {
	from, _ := rand.Int(r, big.NewInt(floorCount))
	to, _ := rand.Int(r, big.NewInt(floorCount))
	// }

	ride := RideRequest{
		From: int(from.Int64()),
		To:   int(to.Int64()),
	}

	if rf.verbose {
		fmt.Printf("new ride -- From: %d, To: %d", ride.From, ride.To)
	}

	rf.queue <- ride
}
