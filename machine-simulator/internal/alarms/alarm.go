package alarms

import "time"

type Severity string

const (
	SeverityLow      Severity = "LOW"
	SeverityMedium   Severity = "MEDIUM"
	SeverityHigh     Severity = "HIGH"
	SeverityCritical Severity = "CRITICAL"
)

type Alarm struct {
	ID          string
	MachineID   string
	MachineName string

	Message  string
	Severity Severity

	Active bool

	Timestamp time.Time
}