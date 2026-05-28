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
	ID   string
	Name string
	FSM  RuntimeStateMachine

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
		FSM:      RuntimeStateMachine{State: StateStopped},
		EventBus: eventBus,
		Pressure: 2.0,
		FlowRate: 0,
	}
}

func (p *Pump) GetName() string {
	return p.Name
}

func (p *Pump) GetState() MachineState {
	return p.FSM.State
}

func (p *Pump) Start() {
	p.FSM.Transition(
		StateRunning,
		0,
	)
}

func (p *Pump) Stop() {
	p.FSM.Transition(
		StateStopped,
		0,
	)
}

func (p *Pump) Tick() time.Duration {
	return 2 * time.Second
}

func (p *Pump) Cycle() {

	p.Update()
	p.publishTelemetry()
	p.checkAlarms()
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

	now := time.Now()

	if p.FSM.State == StateFault {

		p.DowntimeSeconds += 2

		if now.After(p.FSM.NextTransitionAt) {

			p.FSM.Transition(
				StateRecovering,
				3*time.Second,
			)
		}
		return
	}

	if p.FSM.State == StateRecovering {

		if now.After(p.FSM.NextTransitionAt) {

			p.FSM.Transition(
				StateRunning,
				3*time.Second,
			)
		}
		return
	}

	if p.FSM.State == StateMaintenance {

		if now.After(p.FSM.NextTransitionAt) {

			p.FSM.Transition(
				StateRunning,
				5*time.Second,
			)
		}
		return
	}

	if p.FSM.State != StateRunning {
		return
	}

	if rand.Float64() < 0.003 {

		p.FSM.Transition(
			StateMaintenance,
			5*time.Second,
		)
		return
	}

	p.Pressure = 2 + rand.Float64()*3

	p.FlowRate = 50 + rand.Float64()*15

	p.PowerConsumption = 6 + rand.Float64()*4

	if rand.Float64() < 0.01 {

		p.ErrorCount++
		p.FSM.Transition(
			StateFault,
			8*time.Second,
		)
	}

	p.RuntimeSeconds += 2
}

func (p *Pump) publishTelemetry() {

	t := telemetry.Telemetry{
		MachineID:   p.ID,
		MachineName: p.Name,
		State:       string(p.FSM.State),

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
