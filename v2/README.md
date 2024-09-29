## Version 2

This version solves the problem of managing multiple elevators, statically increasing the previous programs scalability but the growing pains are solved yet. Just more scale. This program will allow an allotment of elevators and 'people' that want to enter. A random floor will be generated for each person, those numbers will be grouped into the same elevator until a capacity of the elevator is reached.

All elevators are assuming to start on floor 1

### Inputs

Elevator #
Floors #
People #

### Constraints ?

- Elevator capacity - 5
- Elevators # - Floors / 4
- Floor # - 30
- People # - Elevators Capacity \* Elevator #

This problem will build on the elevator interface. Upon execution will validate the inputs, generate the floors for each person, and determine the groups, then send them on their way. Use a simple hash algo.
