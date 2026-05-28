package machines

import (
	"context"
	"math/rand"
	"smart-factory/machine-simulator/internal/alarms"
	"smart-factory/machine-simulator/internal/events"
	"smart-factory/machine-simulator/internal/telemetry"
	"time"
)

type Oven struct {
	ID    string
	Name  string
	State MachineState

	Temperature       float64
	TargetTemperature float64
	PowerConsumption  float64
	RuntimeSeconds    float64
	DowntimeSeconds   float64
	ErrorCount        int

	EventBus *events.Bus
}

func NewOven(
	name string,
	eventBus *events.Bus,
) *Oven {

	return &Oven{
		ID:                NewMachineID(),
		Name:              name,
		State:             StateStopped,
		Temperature:       25,
		TargetTemperature: 180,
		EventBus:          eventBus,
	}
}

func (o *Oven) GetName() string {
	return o.Name
}

func (o *Oven) GetState() MachineState {
	return o.State
}

func (o *Oven) Start() {
	o.State = StateRunning
}

func (c *Oven) Stop() {
	c.State = StateStopped
}

func (o *Oven) Run(ctx context.Context) {

	ticker := time.NewTicker(3 * time.Second)

	defer ticker.Stop()

	for {
		select {

		case <-ctx.Done():
			return

		case <-ticker.C:

			o.Update()

			o.publishTelemetry()

			o.checkAlarms()
		}
	}
}

func (o *Oven) Update() {

	if o.State == StateFault {

		o.DowntimeSeconds += 3

		if rand.Float64() < 0.10 {
			o.State = StateRecovering
		}
		return
	}

	if o.State == StateRecovering {

		time.Sleep(3 * time.Second)
		o.State = StateRunning
		return
	}

	if o.State == StateMaintenance {

		time.Sleep(5 * time.Second)
		o.State = StateRunning
		return
	}

	if o.State != StateRunning {
		return
	}

	if rand.Float64() < 0.002 {

		o.State = StateMaintenance
		return
	}

	delta := (o.TargetTemperature - o.Temperature) * 0.15

	o.Temperature += delta

	o.Temperature += rand.Float64()*4 - 2

	o.PowerConsumption = 12 + rand.Float64()*5

	if rand.Float64() < 0.015 {

		o.State = StateFault
		o.ErrorCount++
	}

	o.RuntimeSeconds += 3
}

func (o *Oven) publishTelemetry() {

	t := telemetry.Telemetry{
		MachineID:   o.ID,
		MachineName: o.Name,
		State:       string(o.State),

		Temperature: o.Temperature,
		Power:       o.PowerConsumption,

		RuntimeSeconds:  o.RuntimeSeconds,
		DowntimeSeconds: o.DowntimeSeconds,
		ErrorCount:      o.ErrorCount,

		Timestamp: time.Now(),
	}

	o.EventBus.Events <- events.Event{
		Type:    "telemetry",
		Payload: t,
	}
}

func (o *Oven) checkAlarms() {

	if o.Temperature > 220 {

		alarm := alarms.Alarm{
			MachineID:   o.ID,
			MachineName: o.Name,
			Message:     "Oven critical overheating",
			Severity:    alarms.SeverityCritical,
			Active:      true,
			Timestamp:   time.Now(),
		}

		o.EventBus.Events <- events.Event{
			Type:    "alarm",
			Payload: alarm,
		}
	}
}
