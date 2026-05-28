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
	ID   string
	Name string
	FSM  RuntimeStateMachine

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
		FSM:               RuntimeStateMachine{State: StateStopped},
		Temperature:       25,
		TargetTemperature: 180,
		EventBus:          eventBus,
	}
}

func (o *Oven) GetName() string {
	return o.Name
}

func (o *Oven) GetState() MachineState {
	return o.FSM.State
}

func (o *Oven) Start() {
	o.FSM.Transition(
		StateRunning,
		0,
	)
}

func (o *Oven) Stop() {
	o.FSM.Transition(
		StateStopped,
		0,
	)
}

func (o *Oven) Tick() time.Duration {
	return 3 * time.Second
}

func (o *Oven) Cycle() {

	o.Update()
	o.publishTelemetry()
	o.checkAlarms()
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

	now := time.Now()

	if o.FSM.State == StateFault {

		o.DowntimeSeconds += 3

		if now.After(o.FSM.NextTransitionAt) {

			o.FSM.Transition(
				StateRecovering,
				3*time.Second,
			)
		}
		return
	}

	if o.FSM.State == StateRecovering {

		if now.After(o.FSM.NextTransitionAt) {

			o.FSM.Transition(
				StateRunning,
				3*time.Second,
			)
		}
		return
	}

	if o.FSM.State == StateMaintenance {

		if now.After(o.FSM.NextTransitionAt) {

			o.FSM.Transition(
				StateRunning,
				5*time.Second,
			)
		}
		return
	}

	if o.FSM.State != StateRunning {
		return
	}

	if rand.Float64() < 0.002 {

		o.FSM.Transition(
			StateMaintenance,
			6*time.Second,
		)
		return
	}

	delta := (o.TargetTemperature - o.Temperature) * 0.15

	o.Temperature += delta

	o.Temperature += rand.Float64()*4 - 2

	o.PowerConsumption = 12 + rand.Float64()*5

	if rand.Float64() < 0.015 {

		o.ErrorCount++
		o.FSM.Transition(
			StateFault,
			8*time.Second,
		)
	}

	o.RuntimeSeconds += 3
}

func (o *Oven) publishTelemetry() {

	t := telemetry.Telemetry{
		MachineID:   o.ID,
		MachineName: o.Name,
		State:       string(o.FSM.State),

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
