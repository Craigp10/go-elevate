Go elevate is an application that is just for fun for me to practice designing a quick system. Maybe implementing it for practice and creating test cases around it. Maybe it go beyond a program and into a web app. Who knows

Anyways, the problem is this:

<b>Design an elevator system</b>

There are a few problems to solve within. The first is managing an elevator, the second is scaling this to manage multiple elevators in a static system. The Third problem is managing multiple elevators in a dynamic decouple system.

This project will be created through 3 versions to incrently solve this problem.

1. The first approach will solve the problem of the elevator, where we can call a program to
   send an elevator to multiple floors and return once the elevator has completed its route
2. The second approach will solve the problem of managing multiple elevators, where static elevators and floors are declared to run the program. Upon execution the routes are determined and elevators are ran asynchronously.
3. The thrid approach will be building off of the second, implementing a scheduling system to handle dynamic inputs.
4. Fourth, if it comes to it will be to handle improved concurrency and design patterns to simplify code.
