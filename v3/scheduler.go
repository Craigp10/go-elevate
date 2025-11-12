package v3

import (
	"fmt"
	"slices"
	"sync"

	"github.com/google/uuid"
)

// tbd on this idea
// type CommandQueue struct {
// 	UpQueue   <-chan RideQueue
// 	DownQueue <-chan RideQueue
// 	SchQueue  chan<- RideRequest
// }

// type ChanQueue[T RideRequest | RideQueue | int] struct {
// 	queue chan T
// }

// func newChanQueue[T RideRequest | RideQueue | int](size int, queue T) *ChanQueue[T] {
// 	return &ChanQueue[T]{
// 		queue: make(chan T, size),
// 	}
// }

// If we want to extend a scheduler to handle multiple floors, consider mutex's around these queues
type Scheduler struct {
	RequestQueue       chan RideRequest // schQueue on CommandQueue
	ElevatorsRegistry  map[string]*Elevator
	Floors             int
	AvailableElevators chan string // elevator id
	ActiveElevators    chan int    // size 1
	upRides            RideArray   // limit 6 stops
	downRides          RideArray   // limit 6 stops
	upQueue            chan RideQueue
	downQueue          chan RideQueue
}

type RideQueue struct {
	RideRequest
	Direction int
}

type RideArray struct {
	items []RideQueue
	mut   *sync.Mutex
}

func (ra *RideArray) Clear() {
	ra.items = []RideQueue{}
}

func newRideArray(rides ...RideQueue) RideArray {
	items := make([]RideQueue, 0)
	return RideArray{
		items: append(items, rides...),
		mut:   &sync.Mutex{},
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

func (ra *RideArray) Sort(direction int) {
	slices.SortFunc(ra.items, func(a, b RideQueue) int {
		if direction < 0 {
			return b.Direction - a.Direction
		} else {
			return a.Direction + b.Direction
		}
	})
}

func (s *Scheduler) Run(elevatorCount int) {
	for {
		select {
		case req := <-s.RequestQueue:
			fmt.Println("New ride request", req)
			// Any additional work that we want for a new 'request' coming in...

			ride := RideQueue{
				RideRequest: req,
				Direction:   setDirection(req),
			}

			if ride.Direction > 0 {
				s.upQueue <- ride
			} else {
				s.downQueue <- ride
			}

		case downRide := <-s.downQueue:
			if s.downRides.Contains(downRide) {
				continue
			}

			if s.downRides.Length() >= 4 {
				s.downRides.mut.Lock()
				defer s.downRides.mut.Unlock()

				fmt.Printf("down ride sending: %+v\n", s.downRides)

				nextElevatorId := <-s.AvailableElevators
				nextElevator := s.ElevatorsRegistry[nextElevatorId]
				s.downRides.Sort(-1)
				nextElevator.SetRoute(s.downRides.Flatten()...)
				s.downRides.Clear()

				go nextElevator.Go()
			} else {
				s.downRides.items = append(s.downRides.items, downRide)
			}

		case upRide := <-s.upQueue:
			if s.upRides.Contains(upRide) {
				continue
			}

			if s.upRides.Length() >= 4 {
				s.upRides.mut.Lock()
				defer s.upRides.mut.Unlock()

				fmt.Printf("up ride sending: %+v\n", s.upRides)

				nextElevatorId := <-s.AvailableElevators
				nextElevator := s.ElevatorsRegistry[nextElevatorId]
				s.downRides.Sort(1)
				nextElevator.SetRoute(s.upRides.Flatten()...)
				s.upRides.Clear()

				go nextElevator.Go()
			} else {
				s.upRides.items = append(s.upRides.items, upRide)
			}
			// TODO -- allows rides to be taken when an elevator frees up.
			// case availableElevator := <-s.AvailableElevators:
			// 	// take bigger direction queue
			// 	fmt.Println("Sending elevator")
			// 	route := s.nextElevatorRoute()

			// 	elevator := s.ElevatorsRegistry[availableElevator]
			// 	elevator.SetRoute(route...)
			// 	go elevator.Go()
		}
	}
}

func (s *Scheduler) Close() {
	close(s.downQueue)
	close(s.upQueue)
	close(s.ActiveElevators)
	close(s.AvailableElevators)
}

func (s *Scheduler) nextElevatorRoute() []int {
	if s.downRides.Length() > s.upRides.Length() {
		return s.downRides.Flatten()
	}

	return s.upRides.Flatten()
}

func NewScheduler(floors int, elevatorCount int, verbose bool, rideRequestQueue chan RideRequest) *Scheduler {
	elevators := make(map[string]*Elevator, elevatorCount)
	idleElevatorsChan := make(chan string, elevatorCount)

	for i := range elevatorCount {
		id := uuid.New().String()
		elevators[id] = newElevator(verbose, id, idleElevatorsChan)
		i++
		idleElevatorsChan <- id
	}

	return &Scheduler{
		RequestQueue:       rideRequestQueue,
		upRides:            newRideArray(),
		downRides:          newRideArray(),
		upQueue:            make(chan RideQueue, 2*floors*floors),
		downQueue:          make(chan RideQueue, 2*floors*floors),
		Floors:             floors,
		ElevatorsRegistry:  elevators,
		AvailableElevators: idleElevatorsChan,
		ActiveElevators:    make(chan int, elevatorCount),
	}
}

func setDirection(req RideRequest) int {
	if req.To > req.From {
		return 1
	}

	return -1
}
