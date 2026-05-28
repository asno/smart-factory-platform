package alarms

import (
	"smart-factory/machine-simulator/internal/logging"

	"go.uber.org/zap"
)

func Dispatch(a Alarm) {

	logging.Logger.Warn(
		"alarm triggered",

		zap.String("machine", a.MachineName),
		zap.String("message", a.Message),
		zap.String("severity", string(a.Severity)),
	)
}
