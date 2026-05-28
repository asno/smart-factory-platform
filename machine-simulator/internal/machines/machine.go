package machines

import (
	"time"

	"github.com/google/uuid"
)

type MachineState string

const (
	StateStopped     MachineState = "STOPPED"
	StateStarting    MachineState = "STARTING"
	StateRunning     MachineState = "RUNNING"
	StateIdle        MachineState = "IDLE"
	StateWarning     MachineState = "WARNING"
	StateFault       MachineState = "FAULT"
	StateRecovering  MachineState = "RECOVERING"
	StateMaintenance MachineState = "MAINTENANCE"
)

type MachineTelemetry struct {
	ID        string
	Name      string
	State     MachineState
	Timestamp time.Time
}

type Machine interface {
	Start()
	Stop()
	Update()
	GetTelemetry() MachineTelemetry
}

func NewMachineID() string {
	return uuid.New().String()
}
