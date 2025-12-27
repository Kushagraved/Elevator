package elevator

import (
	"container/heap"
	"elevator-system/pkg/helpers"
	"elevator-system/pkg/models"
	"fmt"
	"sync"
)

type Elevator struct {
	ID        int64
	Floor     int64
	Direction models.Direction

	UpQueue   helpers.MinHeap
	DownQueue helpers.MaxHeap
	mutex     sync.Mutex
}

func NewElevator(id int64) *Elevator {
	return &Elevator{
		ID:        id,
		Floor:     0,
		Direction: models.IDLE,
	}
}

func (e *Elevator) AddRequest(request models.InternalRequest) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if request.Floor > e.Floor {
		heap.Push(&e.UpQueue, request.Floor)
	} else if request.Floor < e.Floor {
		heap.Push(&e.DownQueue, request.Floor)
	}
}

// Step Next Logical Move
func (e *Elevator) Step() {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	// No, pending work
	if e.UpQueue.Len() == 0 && e.DownQueue.Len() == 0 {
		e.Direction = models.IDLE
		return
	}
	switch e.Direction {
	case models.UP, models.IDLE:
		if e.UpQueue.Len() > 0 {
			next := heap.Pop(&e.UpQueue).(int64)
			e.moveTo(next)
			return
		}
	case models.DOWN:
		if e.DownQueue.Len() > 0 {
			next := heap.Pop(&e.DownQueue).(int64)
			e.moveTo(next)
			return
		}
	}
}

func (e *Elevator) moveTo(floor int64) {
	if floor > e.Floor {
		e.Direction = models.UP
	} else if floor < e.Floor {
		e.Direction = models.DOWN
	}
	fmt.Printf("Elevator %d: %d -> %d\n", e.ID, e.Floor, floor)
	e.Floor = floor
}
