Go elevate is an application that is just for fun for me to practice designing a quick system. Maybe implementing it for practice and creating test cases around it. Maybe it go beyond a program and into a web app. Who knows

Anyways, the problem is this:

<b>Design an elevator system</b>

The requirements are

- The system must be able to handle n number of elevator inputs, max is 10 for now
- The elevators must operate in parallel
- When an elevator begins a `move`, it must complete all stop before changing directions (up or down)
- An elevator can have a max of 3 request fulfilled on each `move`, this may result in up to 5 stop.
- Once an elevator begins a move, it cannot receive more rides
  -- in other words, rides can only be scheduled on an elevator while it is on hold / moving to begin a `move`.

Thinking through this, we'll need to take advantage of Gos concurrency model to schedule these elevator rides off of the main thread, but also take advantage of a mutex to schedule rides and waitgroup to know when they are finished. A scheduler design pattern of some kind will need to be used.

Elevators will not manage their state, but use it to operate.
Report back when they have completed their trip
