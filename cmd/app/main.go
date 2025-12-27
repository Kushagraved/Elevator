package main

import (
	"elevator-system/pkg/elevator"
	"elevator-system/pkg/models"
	"time"
)

func main() {
	// Create controller with 2 elevators
	controller := elevator.NewController(2)

	// -------- External pickup requests --------
	controller.AddRequest(models.ExternalRequest{
		Floor:     5,
		Direction: models.UP,
	})

	controller.AddRequest(models.ExternalRequest{
		Floor:     2,
		Direction: models.DOWN,
	})

	controller.AddRequest(models.ExternalRequest{
		Floor:     8,
		Direction: models.UP,
	})

	controller.AddRequest(models.ExternalRequest{
		Floor:     1,
		Direction: models.UP,
	})

	//var wg sync.WaitGroup
	for _, e := range controller.Elevators {
		//wg.Add(1)
		go func(el *elevator.Elevator) {
			//defer wg.Done()
			for {
				el.Step()
				time.Sleep(500 * time.Millisecond)
			}
		}(e)
	}

	//wg.Wait()

	select {} // keep main al
}
