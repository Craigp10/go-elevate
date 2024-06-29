package main

import "fmt"

type Config struct {
	ElevatorsCount int
}

type Request struct {
	Origin      int
	Destination int
}

// If we want to extend a scheduler to handle multiple floors, consider mutex's around these queues
type Scheduler struct {
	Queue           chan Request
	Elevators       map[int]Elevator
	Floors          int
	ActiveQueue     []Request
	ActiveElevators chan int // elevator id
	MovingElevator  chan int // size 1
}

func (s *Scheduler) Run() {
	for {
		select {
		case req := <-s.Queue:
			fmt.Println("New to queue")
			s.ActiveQueue = append(s.ActiveQueue, req)

			for i := 0; i < len(s.Elevators); i++ {
				if s.Elevators[i].Busy {
					continue
				}
				// else
				// How can we check for non working elevators without having race conditions with below?
				// We could lock the resource? Would need to look more into it.
				// Probably exclude a check here...
			}

		case freeElevator := <-s.ActiveElevators: // Busy elevators 'check in'
			fmt.Println("freeElevator", freeElevator, s.ActiveQueue)
			cur := s.Elevators[freeElevator]
			cur.Busy = true
			// Take top 3 of the queue...
			i := 0
			var routeR []Request
			if len(s.ActiveQueue) <= 3 {
				routeR = s.ActiveQueue[:]
				s.ActiveQueue = make([]Request, 2*s.Floors*s.Floors)
			} else {
				routeR = s.ActiveQueue[:3]
				s.ActiveQueue = s.ActiveQueue[3:]
				i++
			}
			var route []int
			for i := range routeR {
				route = append(route, routeR[i].Origin, routeR[i].Destination)
			}
			fmt.Println(route)
			// Will need to flatten the route
			cur.SetRoute(route)
			go cur.Move()
		}
	}
}

func New(floors int, elevatorCount int) Scheduler {
	elevators := make(map[int]Elevator, elevatorCount)
	for i := range elevatorCount {
		id := i + 1
		elevator := Elevator{Id: id}
		elevators[id] = elevator
	}
	q := make(chan Request, 2*floors*floors) // Floor**2 options * 2 for each direction
	aq := make([]Request, 2*floors*floors)   // Floor**2 options * 2 for each direction
	return Scheduler{
		Queue:       q,
		ActiveQueue: aq,
		Floors:      floors,
		Elevators:   elevators,
	}
}

func (s *Scheduler) StartMove(id int, route []int) {
	curr := s.Elevators[id]
	go curr.Move()
}
