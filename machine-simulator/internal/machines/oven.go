package machines

import (
	"context"
	"math/rand"
	"time"
	"smart-factory/machine-simulator/internal/alarms"
	"smart-factory/machine-simulator/internal/events"
	"smart-factory/machine-simulator/internal/telemetry"
)

type Oven struct {
	ID    string
	Name  string
	State MachineState

	Temperature       float64
	TargetTemperature float64
	PowerConsumption  float64

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

func (o *Oven) Start() {
	o.State = StateRunning
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

	if o.State != StateRunning {
		return
	}

	delta := (o.TargetTemperature - o.Temperature) * 0.15

	o.Temperature += delta

	o.Temperature += rand.Float64()*4 - 2

	o.PowerConsumption = 12 + rand.Float64()*5
}

func (o *Oven) publishTelemetry() {

	t := telemetry.Telemetry{
		MachineID:   o.ID,
		MachineName: o.Name,
		State:       string(o.State),

		Temperature: o.Temperature,
		Power:       o.PowerConsumption,

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