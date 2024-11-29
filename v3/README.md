## Go Elevate Version 3

The goal of this version is to create an elevator <b>system</b>. No longer a single elevator but multiple elevators that are managed and listen to live updates/waiting users. What does that mean? It means people are no longer all on the same floor, grouped into the same elevator by buckets. People are distributed across floors and elevators will be routed to waiting floors while en route.

### Inputs

Elevator #
Floors #
People #

### Constraints ?

- Elevator capacity - 5
- Elevators # - Floors / 4
- Floor # - 50
- People # - No limit

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
