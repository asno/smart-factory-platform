package machines

import (
	"context"
	"math/rand"
	"smart-factory/machine-simulator/internal/alarms"
	"smart-factory/machine-simulator/internal/events"
	"smart-factory/machine-simulator/internal/telemetry"
	"time"
)

type Conveyor struct {
	ID    string
	Name  string
	State MachineState

	Speed            float64
	Temperature      float64
	PowerConsumption float64
	RuntimeSeconds   float64
	DowntimeSeconds  float64
	ErrorCount       int

	ProductionCount int

	EventBus *events.Bus
}

func NewConveyor(
	name string,
	eventBus *events.Bus,
) *Conveyor {

	return &Conveyor{
		ID:               NewMachineID(),
		Name:             name,
		State:            StateStopped,
		EventBus:         eventBus,
		Temperature:      25,
		Speed:            0,
		PowerConsumption: 0,
	}
}

func (c *Conveyor) GetName() string {
	return c.Name
}

func (c *Conveyor) GetState() MachineState {
	return c.State
}

func (c *Conveyor) Start() {
	c.State = StateRunning
}

func (c *Conveyor) Stop() {
	c.State = StateStopped
}

func (c *Conveyor) Run(ctx context.Context) {

	baseInterval := 2 * time.Second

	jitter := time.Duration(rand.Intn(500)) * time.Millisecond

	ticker := time.NewTicker(baseInterval + jitter)

	defer ticker.Stop()

	for {
		select {

		case <-ctx.Done():
			return

		case <-ticker.C:

			c.Update()

			c.publishTelemetry()

			c.checkAlarms()
		}
	}
}

func (c *Conveyor) Update() {

	if c.State == StateFault {
		c.DowntimeSeconds += 2

		if rand.Float64() < 0.15 {
			c.State = StateRecovering
		}
		return
	}

	if c.State == StateRecovering {
		time.Sleep(2 * time.Second)
		c.State = StateRunning
		return
	}

	if c.State != StateRunning {
		return
	}

	c.Speed = 65 + rand.Float64()*15

	c.Temperature += rand.Float64()*2 - 1

	c.PowerConsumption = 4 + rand.Float64()*3

	c.ProductionCount += rand.Intn(8)

	c.RuntimeSeconds += 2

	if rand.Float64() < 0.02 {
		c.State = StateFault
		c.ErrorCount++
	}

	if rand.Float64() < 0.003 {
		c.State = StateMaintenance
		time.Sleep(3 * time.Second)
		c.State = StateRunning
		return
	}
}

func (c *Conveyor) publishTelemetry() {

	t := telemetry.Telemetry{
		MachineID:       c.ID,
		MachineName:     c.Name,
		State:           string(c.State),
		Temperature:     c.Temperature,
		Speed:           c.Speed,
		Power:           c.PowerConsumption,
		ProductionCount: c.ProductionCount,
		RuntimeSeconds:  c.RuntimeSeconds,
		DowntimeSeconds: c.DowntimeSeconds,
		ErrorCount:      c.ErrorCount,
		Timestamp:       time.Now(),
	}

	c.EventBus.Events <- events.Event{
		Type:    "telemetry",
		Payload: t,
	}
}

func (c *Conveyor) checkAlarms() {

	if c.Temperature > 80 {

		alarm := alarms.Alarm{
			MachineID:   c.ID,
			MachineName: c.Name,
			Message:     "Conveyor overheating",
			Severity:    alarms.SeverityHigh,
			Active:      true,
			Timestamp:   time.Now(),
		}

		c.EventBus.Events <- events.Event{
			Type:    "alarm",
			Payload: alarm,
		}
	}
}
