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

// If we want to extend a scheduler to handle multiple floors, consider mutex's around these queues
type Scheduler struct {
	RequestQueue      *ChanQueue[RideRequest] // schQueue on CommandQueue
	ElevatorsRegistry map[int]*Elevator
	Floors            int
	// ActiveQueue        chan RideQueue
	AvailableElevators *ChanQueue[int] // elevator id
	ActiveElevators    *ChanQueue[int] // size 1
	IdleElevator       *Elevator
	idleChan           *ChanQueue[int]
	upRides            RideArray // limit 6 stops
	downRides          RideArray // limit 6 stops
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
			fmt.Println("New ride request")
			// Any additional work that we want for a new 'request' coming in...

			ride := RideQueue{
				RideRequest: req,
				Direction:   setDirection(req),
			}

			if ride.Direction > 0 {
				s.downQueue.mut.Lock()
				defer s.downQueue.mut.Unlock()
				s.upQueue.queue <- ride
			} else {
				s.downQueue.mut.Lock()
				defer s.downQueue.mut.Unlock()
				s.downQueue.queue <- ride
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
		// case req := <-s.ActiveQueue:
		// 	// New request passed in

		// 	if arr.Contains(req.RideRequest) {
		// 		return
		// 	}

		// 	arr.items = append(arr.items, req)

		// Come back to this... This handles idle elevators.. probably b/c there was no ride.
		// for _, ele := range s.ElevatorsRegister {
		// 	if ele.State == STATE_ACTIVE {
		// 		continue
		// 	} else if ele.State == STATE_IDLE {
		// 		ele.Move(req.From)
		// 	} else {

		// 	}

		// Attempts to put ride on active elevator... this is broken come back too
		// case activeElevator := <-s.ActiveElevators: // Busy elevators 'check in'
		// 	if arr.Length() == 0 {
		// 		fmt.Println("Available Elevator -- No active request")
		// 		// Somehow need to wait here until we get in a request
		// 		// s.IdleElevator = freeElevator
		// 		continue
		// 	}

		// 	// fmt.Println("Available Elevator", availableElevator, s.AvailableElevators)
		// 	cur := s.ElevatorsRegister[activeElevator]

		// 	// cur.State =
		// 	// cur.Route = []int{first.From, first.To}

		// 	for len(cur.Route) <= 4 {
		// 		next := <-s.ActiveQueue
		// 		if next.Direction != cur.State {
		// 			s.ActiveQueue <- next
		// 		}
		// 		cur.Route = append(cur.Route, next.From, next.To)
		// 	}

		// 	// Will need to flatten the route
		// 	go cur.Go()
		case availableElevator := <-s.AvailableElevators.queue:
			// take bigger direction queue
			route := s.nextElevatorRoute()

			elevator := s.ElevatorsRegistry[availableElevator]
			elevator.SetRoute(route)
			go elevator.Go()
		}
	}
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

	for i := range elevators {
		elevators[i+1] = newElevator(verbose, i, idleElevatorsChan)
		i++
	}

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
