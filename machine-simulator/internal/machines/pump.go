package machines

import (
	"context"
	"math/rand"
	"smart-factory/machine-simulator/internal/alarms"
	"smart-factory/machine-simulator/internal/events"
	"smart-factory/machine-simulator/internal/telemetry"
	"time"
)

type Pump struct {
	ID    string
	Name  string
	State MachineState

	Pressure         float64
	FlowRate         float64
	PowerConsumption float64
	RuntimeSeconds   float64
	DowntimeSeconds  float64
	ErrorCount       int

	EventBus *events.Bus
}

func NewPump(
	name string,
	eventBus *events.Bus,
) *Pump {

	return &Pump{
		ID:       NewMachineID(),
		Name:     name,
		State:    StateStopped,
		EventBus: eventBus,
		Pressure: 2.0,
		FlowRate: 0,
	}
}

func (p *Pump) GetName() string {
	return p.Name
}

func (p *Pump) GetState() MachineState {
	return p.State
}

func (p *Pump) Start() {
	p.State = StateRunning
}

func (c *Pump) Stop() {
	c.State = StateStopped
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

	if p.State == StateFault {

		p.DowntimeSeconds += 2

		if rand.Float64() < 0.12 {
			p.State = StateRecovering
		}

		return
	}

	if p.State == StateRecovering {

		time.Sleep(2 * time.Second)
		p.State = StateRunning
		return
	}

	if p.State == StateMaintenance {

		time.Sleep(4 * time.Second)
		p.State = StateRunning
		return
	}

	if p.State != StateRunning {
		return
	}

	if rand.Float64() < 0.003 {

		p.State = StateMaintenance
		return
	}

	p.Pressure = 2 + rand.Float64()*3

	p.FlowRate = 50 + rand.Float64()*15

	p.PowerConsumption = 6 + rand.Float64()*4

	if rand.Float64() < 0.01 {

		p.State = StateFault
		p.ErrorCount++
	}

	p.RuntimeSeconds += 2
}

func (p *Pump) publishTelemetry() {

	t := telemetry.Telemetry{
		MachineID:   p.ID,
		MachineName: p.Name,
		State:       string(p.State),

		Power: p.PowerConsumption,

		RuntimeSeconds:  p.RuntimeSeconds,
		DowntimeSeconds: p.DowntimeSeconds,
		ErrorCount:      p.ErrorCount,

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
