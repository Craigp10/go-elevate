package main

func main() {
	ele := Elevator{
		Route: []int{1, 3, 5, 6, 1},
		Floor: 0,
	}

	ele.Move()
}
