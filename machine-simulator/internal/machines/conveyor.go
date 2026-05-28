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
	ID   string
	Name string
	FSM  RuntimeStateMachine

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
		ID:   NewMachineID(),
		Name: name,
		FSM: RuntimeStateMachine{
			State: StateStopped,
		},
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
	return c.FSM.State
}

func (c *Conveyor) Start() {
	c.FSM.Transition(
		StateRunning,
		0,
	)
}

func (c *Conveyor) Stop() {
	c.FSM.Transition(
		StateStopped,
		0,
	)
}

func (c *Conveyor) Tick() time.Duration {
	return 2 * time.Second
}

func (c *Conveyor) Cycle() {

	c.Update()
	c.publishTelemetry()
	c.checkAlarms()
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

	now := time.Now()

	if c.FSM.State == StateFault {

		c.DowntimeSeconds += 2

		if now.After(c.FSM.NextTransitionAt) {

			c.FSM.Transition(
				StateRecovering,
				3*time.Second,
			)
		}
		return
	}

	if c.FSM.State == StateRecovering {

		if now.After(c.FSM.NextTransitionAt) {

			c.FSM.Transition(
				StateRunning,
				0,
			)
		}
		return
	}

	if c.FSM.State == StateMaintenance {

		if now.After(c.FSM.NextTransitionAt) {

			c.FSM.Transition(
				StateRunning,
				0,
			)
		}
		return
	}

	if c.FSM.State != StateRunning {
		return
	}

	if rand.Float64() < 0.003 {

		c.FSM.Transition(
			StateMaintenance,
			5*time.Second,
		)
		return
	}

	c.Speed = 65 + rand.Float64()*15

	c.Temperature += rand.Float64()*2 - 1

	c.PowerConsumption = 4 + rand.Float64()*3

	c.ProductionCount += rand.Intn(8)

	if rand.Float64() < 0.02 {

		c.ErrorCount++
		c.FSM.Transition(
			StateFault,
			8*time.Second,
		)
	}

	c.RuntimeSeconds += 2
}

func (c *Conveyor) publishTelemetry() {

	t := telemetry.Telemetry{
		MachineID:       c.ID,
		MachineName:     c.Name,
		State:           string(c.FSM.State),
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
