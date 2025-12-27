package elevator

import (
	"elevator-system/pkg/helpers"
	"elevator-system/pkg/models"
	"math"
)

type Controller struct {
	Elevators []*Elevator
}

func NewController(n int) *Controller {
	elevators := make([]*Elevator, n)
	for i := 0; i < n; i++ {
		elevators[i] = NewElevator(int64(i))
	}
	return &Controller{
		Elevators: elevators,
	}
}

func (c *Controller) AddRequest(request models.ExternalRequest) {
	bestElevator := c.selectBestElevator(request)
	bestElevator.AddRequest(models.InternalRequest{
		Floor: request.Floor,
	})
}

func (c *Controller) selectBestElevator(request models.ExternalRequest) *Elevator {
	minDistance := int64(math.MaxInt64)
	var best *Elevator

	for _, e := range c.Elevators {
		dist := helpers.Abs(e.Floor, request.Floor)

		// Rule 1: Idle elevator
		if e.Direction == models.IDLE {
			if dist < minDistance {
				minDistance = dist
				best = e
			}
			continue
		}

		// Rule 2: Same direction
		if e.Direction == request.Direction {
			if (request.Direction == models.UP && e.Floor <= request.Floor) ||
				(request.Direction == models.DOWN && e.Floor >= request.Floor) {

				if dist < minDistance {
					minDistance = dist
					best = e
				}
			}
		}
	}

	if best == nil {
		best = c.Elevators[0]
	}
	return best
}
