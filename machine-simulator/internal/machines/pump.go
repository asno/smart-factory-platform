package machines

import (
	"context"
	"math/rand"
	"time"
	"smart-factory/machine-simulator/internal/alarms"
	"smart-factory/machine-simulator/internal/events"
	"smart-factory/machine-simulator/internal/telemetry"
)

type Pump struct {
	ID    string
	Name  string
	State MachineState

	Pressure         float64
	FlowRate         float64
	PowerConsumption float64

	EventBus *events.Bus
}

func NewPump(
	name string,
	eventBus *events.Bus,
) *Pump {

	return &Pump{
		ID:        NewMachineID(),
		Name:      name,
		State:     StateStopped,
		EventBus:  eventBus,
		Pressure:  2.0,
		FlowRate:  0,
	}
}

func (p *Pump) Start() {
	p.State = StateRunning
}

func (p *Pump) Run(ctx context.Context) {

	ticker := time.NewTicker(2 * time.Second)

	defer ticker.Stop()

	for {
		select {

		case <-ctx.Done():
			return

		case <-ticker.C:

			p.Update()

			p.publishTelemetry()

			p.checkAlarms()
		}
	}
}

func (p *Pump) Update() {

	if p.State != StateRunning {
		return
	}

	p.Pressure = 2 + rand.Float64()*3

	p.FlowRate = 50 + rand.Float64()*15

	p.PowerConsumption = 6 + rand.Float64()*4

	if rand.Float64() < 0.005 {
		p.State = StateFault
	}
}

func (p *Pump) publishTelemetry() {

	t := telemetry.Telemetry{
		MachineID:   p.ID,
		MachineName: p.Name,
		State:       string(p.State),

		Power: p.PowerConsumption,

		Timestamp: time.Now(),
	}

	p.EventBus.Events <- events.Event{
		Type:    "telemetry",
		Payload: t,
	}
}

func (p *Pump) checkAlarms() {

	if p.Pressure > 4.5 {

		alarm := alarms.Alarm{
			MachineID:   p.ID,
			MachineName: p.Name,
			Message:     "Pump overpressure",
			Severity:    alarms.SeverityHigh,
			Active:      true,
			Timestamp:   time.Now(),
		}

		p.EventBus.Events <- events.Event{
			Type:    "alarm",
			Payload: alarm,
		}
	}
}