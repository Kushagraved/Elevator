package main

import (
	"container/heap"
	"fmt"
	"sync"
)

type Direction int

const (
	UP Direction = iota
	DOWN
	IDLE
)

type Request struct {
	Floor     int
	Direction Direction
}

/* ---------- Priority Queues ---------- */

type MinHeap []int

func (h MinHeap) Len() int            { return len(h) }
func (h MinHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h MinHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *MinHeap) Pop() interface{} {
	x := (*h)[len(*h)-1]
	*h = (*h)[:len(*h)-1]
	return x
}

type MaxHeap []int

func (h MaxHeap) Len() int            { return len(h) }
func (h MaxHeap) Less(i, j int) bool  { return h[i] > h[j] }
func (h MaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *MaxHeap) Pop() interface{} {
	x := (*h)[len(*h)-1]
	*h = (*h)[:len(*h)-1]
	return x
}

/* ---------- Elevator ---------- */

type Elevator struct {
	ID           int
	CurrentFloor int
	Direction    Direction
	upQueue      MinHeap
	downQueue    MaxHeap
	mu           sync.Mutex
}

func NewElevator(id int) *Elevator {
	return &Elevator{ID: id, Direction: IDLE}
}

func (e *Elevator) AddRequest(floor int) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if floor > e.CurrentFloor {
		heap.Push(&e.upQueue, floor)
	} else if floor < e.CurrentFloor {
		heap.Push(&e.downQueue, floor)
	}
}

func (e *Elevator) Step() {
	e.mu.Lock()
	defer e.mu.Unlock()

	if (e.Direction == UP || e.Direction == IDLE) && e.upQueue.Len() > 0 {
		e.moveTo(heap.Pop(&e.upQueue).(int))
		return
	}
	if e.downQueue.Len() > 0 {
		e.moveTo(heap.Pop(&e.downQueue).(int))
		return
	}
	e.Direction = IDLE
}

func (e *Elevator) moveTo(floor int) {
	if floor > e.CurrentFloor {
		e.Direction = UP
	} else {
		e.Direction = DOWN
	}
	fmt.Printf("Elevator %d: %d -> %d\n", e.ID, e.CurrentFloor, floor)
	e.CurrentFloor = floor
}

/* ---------- Controller ---------- */

type ElevatorController struct {
	Elevators []*Elevator
}

func NewController(n int) *ElevatorController {
	e := make([]*Elevator, n)
	for i := 0; i < n; i++ {
		e[i] = NewElevator(i)
	}
	return &ElevatorController{Elevators: e}
}

func (c *ElevatorController) RequestElevator(req Request) {
	best := c.selectBestElevator(req)
	best.AddRequest(req.Floor)
}

func (c *ElevatorController) selectBestElevator(req Request) *Elevator {
	var best *Elevator
	minDist := int(^uint(0) >> 1)

	for _, e := range c.Elevators {
		dist := abs(e.CurrentFloor - req.Floor)

		if e.Direction == IDLE ||
			(e.Direction == req.Direction &&
				((req.Direction == UP && e.CurrentFloor <= req.Floor) ||
					(req.Direction == DOWN && e.CurrentFloor >= req.Floor))) {

			if dist < minDist {
				minDist = dist
				best = e
			}
		}
	}
	if best == nil {
		best = c.Elevators[0]
	}
	return best
}

/* ---------- Utils ---------- */

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

/* ---------- Main ---------- */

func main() {
	controller := NewController(2)

	controller.RequestElevator(Request{Floor: 5, Direction: UP})
	controller.RequestElevator(Request{Floor: 2, Direction: DOWN})
	controller.RequestElevator(Request{Floor: 8, Direction: UP})

	for i := 0; i < 5; i++ {
		for _, e := range controller.Elevators {
			e.Step()
		}
	}
}
