package telemetry

import "time"

type Telemetry struct {
	MachineID   string
	MachineName string
	State       string

	Temperature float64
	Speed       float64
	Power       float64

	RuntimeSeconds  float64
	DowntimeSeconds float64
	ErrorCount      int

	ProductionCount int

	Timestamp time.Time
}
