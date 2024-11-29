package v2

import (
	"fmt"
)

// If we want to extend a scheduler to handle multiple floors, consider mutex's around these queues
type Scheduler struct {
	RequestQueue       chan RideRequest
	ElevatorsRegister  map[int]*Elevator
	Floors             int
	ActiveQueue        chan RideQueue
	AvailableElevators chan int // elevator id
	ActiveElevators    chan int // size 1
	PendingEle         *Elevator
}

type RideQueue struct {
	RideRequest
	Direction int
}

type RideArray struct {
	items []RideQueue
}

func (ra *RideArray) Contains(rq RideRequest) bool {
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

	return floors
}

func (ra *RideArray) Length() int {
	return len(ra.items)
}

func (s *Scheduler) Run() {
	var arr RideArray
	for {
		select {
		case req := <-s.RequestQueue:
			fmt.Println("New to queue")
			// Any additional work that we want for a new 'request' coming in...

			s.ActiveQueue <- RideQueue{
				RideRequest: req,
				Direction:   setDirection(req),
			}

		case req := <-s.ActiveQueue:
			// New request passed in

			if arr.Contains(req.RideRequest) {
				return
			}

			arr.items = append(arr.items, req)

			for _, ele := range s.ElevatorsRegister {
				if ele.State == STATE_ACTIVE {
					continue
				} else if ele.State == STATE_IDLE {
					ele.Move(req.From)
				} else {

				}

				ele.State = req.Direction

			}

		case freeElevator := <-s.ActiveElevators: // Busy elevators 'check in'
			if arr.Length == 0 {
				fmt.Println("Available Elevator -- No active request")
				continue
			}

			fmt.Println("Available Elevator", freeElevator, s.ActiveQueue)
			cur := s.ElevatorsRegister[freeElevator]

			// cur.State =
			// cur.Route = []int{first.From, first.To}

			for len(cur.Route) <= 4 {
				next := <-s.ActiveQueue
				if next.Direction != cur.State {
					s.ActiveQueue <- next
				}
				cur.Route = append(cur.Route, next.From, next.To)
			}

			// Will need to flatten the route
			go cur.Go()
		}
	}
}

func NewScheduler(floors int, elevatorCount int) Scheduler {
	elevators := make(map[int]*Elevator, elevatorCount)
	for i := range elevators {
		elevators[i+1] = &Elevator{
			ID:    i + 1,
			Route: make([]int, 0),
			Floor: 0,
			State: 0,
		}
		i++
	}
	q := make(chan RideRequest, 2*floors*floors) // Floor**2 options * 2 for each direction
	aq := make(chan RideQueue, 2*floors*floors)  // Floor**2 options * 2 for each direction
	return Scheduler{
		RequestQueue:       q,
		ActiveQueue:        aq,
		Floors:             floors,
		ElevatorsRegister:  elevators,
		AvailableElevators: make(chan int, elevatorCount),
		ActiveElevators:    make(chan int, elevatorCount),
	}
}

func setDirection(req RideRequest) int {
	if req.To > req.From {
		return 1
	}

	return -1
}
