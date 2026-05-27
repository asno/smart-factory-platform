package telemetry

import "time"

type Telemetry struct {
	MachineID   string
	MachineName string
	State       string

	Temperature float64
	Speed       float64
	Power       float64

	ProductionCount int

	Timestamp time.Time
}