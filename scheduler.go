package main

import "fmt"

type Config struct {
	ElevatorsCount int
}

type Request struct {
	Origin      int
	Destination int
}

type Scheduler struct {
	Queue           chan Request
	Elevators       map[int]Elevator
	Floors          int
	ActiveQueue     []Request
	ActiveElevators chan int // elevator id
}

func (s *Scheduler) Run() {
	for {
		select {
		case req := <-s.Queue:
			s.ActiveQueue = append(s.ActiveQueue, req)
		case freeElevator := <-s.ActiveElevators:
			cur := s.Elevators[freeElevator]
			cur.Busy = true
			// Take top 5 of the queue...
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
