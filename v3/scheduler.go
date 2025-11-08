package v3

import (
	"fmt"
	"slices"
	"sync"
	"time"
)

// tbd on this idea
// type CommandQueue struct {
// 	UpQueue   <-chan RideQueue
// 	DownQueue <-chan RideQueue
// 	SchQueue  chan<- RideRequest
// }

type ChanQueue[T RideRequest | RideQueue | int] struct {
	queue chan T
	mut   *sync.Mutex
}

func TestIt[T comparable](v, a T) bool {
	return v == a
}

func newChanQueue[T RideRequest | RideQueue | int](size int, queue T) *ChanQueue[T] {
	return &ChanQueue[T]{
		queue: make(chan T, size),
		mut:   &sync.Mutex{},
	}
}

func (cq *ChanQueue[T]) Length() int {
	return len(cq.queue)
}

// If we want to extend a scheduler to handle multiple floors, consider mutex's around these queues
type Scheduler struct {
	RequestQueue       *ChanQueue[RideRequest] // schQueue on CommandQueue
	ElevatorsRegistry  map[int]*Elevator
	Floors             int
	AvailableElevators *ChanQueue[int] // elevator id
	ActiveElevators    *ChanQueue[int] // size 1
	upRides            RideArray       // limit 6 stops
	downRides          RideArray       // limit 6 stops
	upQueue            *ChanQueue[RideQueue]
	downQueue          *ChanQueue[RideQueue]
}

type RideQueue struct {
	RideRequest
	Direction int
}

type RideArray struct {
	items []RideQueue
}

func newRideArray(rides ...RideQueue) RideArray {
	items := make([]RideQueue, 0)
	return RideArray{
		items: append(items, rides...),
	}
}

func (ra *RideArray) Contains(rq RideQueue) bool {
	for _, r := range ra.items {
		if r.RideRequest.From == rq.From && r.RideRequest.To == rq.To {
			return true
		}
	}

	return false
}

func (ra *RideArray) Flatten() []int {
	m := make(map[int]bool)

	for _, item := range ra.items {
		m[item.To] = true
		m[item.From] = true
	}

	floors := make([]int, 0)
	i := 0
	for k := range m {
		floors = append(floors, k)
		i++
	}

	slices.SortFunc(floors, func(a, b int) int {
		return a - b
	})

	return floors
}

func (ra *RideArray) Length() int {
	return len(ra.items)
}

func (s *Scheduler) Run() {
	for {
		select {
		case req := <-s.RequestQueue.queue:
			fmt.Println("New ride request", req)
			// Any additional work that we want for a new 'request' coming in...

			ride := RideQueue{
				RideRequest: req,
				Direction:   setDirection(req),
			}

			if ride.Direction > 0 {
				s.upQueue.mut.Lock()
				s.upQueue.queue <- ride
				s.upQueue.mut.Unlock()
			} else {
				s.downQueue.mut.Lock()
				s.downQueue.queue <- ride
				s.downQueue.mut.Unlock()
			}

		case downRide := <-s.downQueue.queue:
			if s.downRides.Contains(downRide) {
				continue
			}

			var wg sync.WaitGroup
			wg.Add(1)
			if s.downRides.Length() >= 5 {
				go func(rides *RideArray) {
					for {
						if rides.Length() == 0 {
							wg.Done()
						}
						time.Sleep(1 * time.Second)
					}
				}(&s.downRides)
				wg.Wait()
			}

			s.downRides.items = append(s.downRides.items, downRide)
		case upRide := <-s.upQueue.queue:
			if s.upRides.Contains(upRide) {
				continue
			}

			var wg sync.WaitGroup
			wg.Add(1)
			if s.upRides.Length() >= 5 {
				go func(rides *RideArray) {
					for {
						if rides.Length() == 0 {
							wg.Done()
						}
						time.Sleep(1 * time.Second)
					}
				}(&s.upRides)
				wg.Wait()
			}

			s.upRides.items = append(s.upRides.items, upRide)
		case availableElevator := <-s.AvailableElevators.queue:
			// take bigger direction queue
			route := s.nextElevatorRoute()

			elevator := s.ElevatorsRegistry[availableElevator]
			elevator.SetRoute(route...)
			go elevator.Go()
		}
	}
}

func (s *Scheduler) Close() {
	close(s.downQueue.queue)
	close(s.upQueue.queue)
	close(s.ActiveElevators.queue)
	close(s.AvailableElevators.queue)
}

func (s *Scheduler) nextElevatorRoute() []int {
	if s.downRides.Length() > s.upRides.Length() {
		return s.downRides.Flatten()
	}

	return s.upRides.Flatten()
}

func NewScheduler(floors int, elevatorCount int, verbose bool, rideRequestQueue *ChanQueue[RideRequest]) Scheduler {
	elevators := make(map[int]*Elevator, elevatorCount)
	var a int
	idleElevatorsChan := newChanQueue(elevatorCount, a)

	for i := range elevatorCount {
		elevators[i+1] = newElevator(verbose, i, idleElevatorsChan)
		i++
	}
	fmt.Println(elevators)
	return Scheduler{
		RequestQueue:       newChanQueue(elevatorCount*elevatorCount, RideRequest{}),
		upRides:            newRideArray(),
		downRides:          newRideArray(),
		upQueue:            newChanQueue(2*floors*floors, RideQueue{}),
		downQueue:          newChanQueue(2*floors*floors, RideQueue{}),
		Floors:             floors,
		ElevatorsRegistry:  elevators,
		AvailableElevators: idleElevatorsChan,
		ActiveElevators:    newChanQueue(elevatorCount, a),
	}
}

func setDirection(req RideRequest) int {
	if req.To > req.From {
		return 1
	}

	return -1
}
