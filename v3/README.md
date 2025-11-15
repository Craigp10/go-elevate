## Go Elevate Version 3

This version of go-elevate will solve the problem of handling multiple elevators w/ on demand ride request. Stimulating a more realistic <b>elevator system</b>. 

No longer a single elevator but multiple elevators that are managed and listen to live updates/waiting users. What does that mean? It means people are no longer all on the same floor, grouped into the same elevator by buckets. People are distributed across floors and elevators will be routed to waiting floors while en route.

### Inputs

Elevator #
Floors #
People # - number of people at start

### Constraints ?

- Elevator capacity - 5
- Elevators # - Floors / 4
- Floor # - 30
- People # - Elevator Capacity * Elevator #

### New concepts

There is no longer a limit on people input because the system can now handle excessive people due to dynamic elevator routing. How? A capacity will be set for the elevator, 3 to a route. While an elevator is moving in a direction, it can pick up any additional floors, as long as it's active floors is not at 3.

### Approach

- Routing - A scheduler that can listen to elevator pick ups through channels
- Elevator state - Elevators will mantain more state than just floors. Route []int, capacity, chan.
- Smart routes - Properly bucket users based on their going to floor and/or from Floor.
- The elevators must operate in parallel, asynchronously
- When an elevator begins a `move`, it must complete all stop before changing directions (up or down), it can add stops
<!-- - Once an elevator begins a move, it cannot receive more rides -- in other words, rides can only be scheduled on an elevator while it is on hold / moving to begin a `move`. -->

Thinking through this, we'll need to take advantage of Gos concurrency model to schedule these elevator rides off of the main thread, but also take advantage of a mutex to schedule rides and waitgroup to know when they are finished. A scheduler design pattern of some kind will need to be used.

idea in my head

Scheduler listens to channel for 'free elevators'. When an elevator becomes free it takes routes up to 3, in a single direction from where it is at atm. When an elevator 'frees' up it will provides it's id through the channel, used on the map for the scheduler to move it. The scheduler will also be listenign to another channel where users can request floors, from-to. (this is where it gets complicated and I may need help here. Deciding how to schedule these in terms of direction is difficult... maybe just put next three on soonest elevator for now.)


Completed -- State of program
Currently the program works as intended, but missing some features to consider...
1. A direction is blocked until it reaches atleast 5 stops so an elevator can begin for it.
2. The program doesn't properly shutdown ad doesn't run continously atm, needs listener for closer to properly shut down.<D-s>
3. Improve logging and design for channels.
