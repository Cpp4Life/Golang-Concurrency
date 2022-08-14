package main

import (
	"fmt"
	"sync"
	"time"
)

// The Dining Philosophers problem is well known in computer science circles.
// Five philosophers, numbered from 0 through 4, live in a house where the
// table is laid for them; each philosopher has their own place at the table.
// Their only difficulty - besides those of philosophy - is that the dish
// served is a very difficult kind of spaghetti which has to be eaten with
// two forks. There are two forks next to each plate, so that presents no
// difficulty. As a consequence, however, this means that no two neighbours
// may be eating simultaneously.

// constants
const hunger = 3

// variables - philosophers
var (
	philosophers = []string{"Aristotle", "Kant", "Descartes", "Spinoza", "Marx"}
	wg           sync.WaitGroup
	sleepTime    = 1 * time.Second
	eatTime      = 2 * time.Second
	thinkTime    = 1 * time.Second
)

func diningProblem(philosopher string, leftFork, rightFork *sync.Mutex) {
	defer wg.Done()

	// print a message
	fmt.Println(philosopher, "is seated.")
	time.Sleep(sleepTime)

	for i := hunger; i > 0; i-- {
		fmt.Println(philosopher, "is hungry.")
		time.Sleep(sleepTime)

		// lock both forks
		leftFork.Lock()
		fmt.Printf("\t%s picked up the fork to his left.\n", philosopher)
		rightFork.Lock()
		fmt.Printf("\t%s picked up the fork to his right.\n", philosopher)

		// print a message
		fmt.Println(philosopher, "has both forks, and is eating.")
		time.Sleep(eatTime)

		// give the philosopher some time to think
		fmt.Println(philosopher, "is thinking.")
		time.Sleep(thinkTime)

		// unlock the mutexes
		rightFork.Unlock()
		fmt.Printf("\t%s put down the fork on his right.\n", philosopher)
		leftFork.Unlock()
		fmt.Printf("\t%s put down the fork on his left.\n", philosopher)
		time.Sleep(sleepTime)
	}

	// print out done message
	fmt.Println(philosopher, "is satisfied.")
	time.Sleep(sleepTime)

	fmt.Println(philosopher, "has left the table.")

}

func main() {
	// print introduction
	fmt.Println("The Dining Philosophers Problem")
	fmt.Println("-------------------------------")

	// add 5 (the number of philosophers) to the wait group
	wg.Add(len(philosophers))

	// we need to create a mutex for the very first fork (the one to
	// the left of the first philosopher). We create it as a pointer,
	// since a sync.Mutex must not be copied after its initial use.
	forkLeft := &sync.Mutex{}

	// spawn one goroutine for each philosopher
	for i := 0; i < len(philosophers); i++ {

		// create a mutex for the right fork
		forkRight := &sync.Mutex{}

		// call a goroutine with the current philosopher, and both mutexes
		go diningProblem(philosophers[i], forkLeft, forkRight)

		// create the next philosopher's left fork (which is the
		// current philosopher's right fork). Note that we are not
		// copying a mutex here; we are making forkLeft equal to the pointer
		// to an existing mutex, so it points to the same location in memory,
		// and does not copy it.
		forkLeft = forkRight
	}

	wg.Wait()

	fmt.Println("The table is empty")
}
