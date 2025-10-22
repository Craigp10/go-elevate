## Version 2

This version of go-elevate aims to solve the scheduling problem, where now multiple elevators are managed and rides need to be grouped and distributed to those elevators. This program will allow an allotment of elevators and 'people' that want to ride. A random floor will be generated for each person, those numbers will be grouped into the same elevator until a capacity of the elevator is reached. Elevator capacity can go over the limit but once reached another floor will not be scheduled on the elevator.

All elevators and people are assuming to start on floor 1.

### Inputs

Elevator Count int 
Number of Floors int 
People Count int 

### Constraints
rounding down...

- Elevator capacity - 5 
    - number of people sent to an elevator
- Elevators # - Floors / 4
    - elevators # must not be greater than # of floors / 4
- Floor # - 30
- People # - Elevators Capacity x Elevator Count 
    - number of people must not be greater than elevactor capacity x elevator count

This problem will build on the elevator interface. Upon execution will validate the inputs, generate the floors for each person, and determine the groups, then send them on their way. Use a simple hash algo.


Example program run
`go run main.go --version=v2 12 3 12`
