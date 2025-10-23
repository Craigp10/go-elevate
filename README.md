Go elevate is an application for fun, to practice designing a quick low level program. Maybe it will go beyond a program and into a web app. Who knows!

Anyways, the problem is this:

<b>Design an elevator system</b>

There are a few problems to solve within. The first is managing a single elevator, the second is scaling the system to manage multiple elevators in a static system, where inputs are pre-defined and concurrency is not heavy. The Third problem is managing multiple elevators in a dynamic decoupled system, where elevators operate thread safe and elevators are managed by dynamic random inputs (real scenario).

This project will be created through 3 versions to incrently solve this problem.

1. The first approach will solve the problem of the elevator, where we can call a program to
   send an elevator to multiple floors and return once the elevator has completed its route. The
   elevator is dumb and moves to exactly where it is told to go... No validations on the floors.

2. The second approach will solve the problem of managing multiple elevators, where static elevators 
   and floors are declared to run the program. Upon execution the routes are determined and elevators are ran asynchronously.

3. The third approach will be building off of the second, implementing a scheduling system to handle dynamic inputs.

4. Fourth, if it comes to it will be to handle improved concurrency and design patterns to simplify code.
