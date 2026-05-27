package machines

import (
	"math/rand"
	"time"
)

type Conveyor struct {
	ID              string
	Name            string
	State           MachineState
	Speed           float64
	ProductionCount int
}

func NewConveyor(name string) *Conveyor {
	return &Conveyor{
		ID:    NewMachineID(),
		Name:  name,
		State: StateStopped,
		Speed: 0,
	}
}

func (c *Conveyor) Start() {
	c.State = StateRunning
}

func (c *Conveyor) Stop() {
	c.State = StateStopped
	c.Speed = 0
}

func (c *Conveyor) Update() {
	if c.State != StateRunning {
		return
	}

	c.Speed = 60 + rand.Float64()*20
	c.ProductionCount += rand.Intn(5)
}

func (c *Conveyor) GetTelemetry() MachineTelemetry {
	return MachineTelemetry{
		ID:        c.ID,
		Name:      c.Name,
		State:     c.State,
		Timestamp: time.Now(),
	}
}