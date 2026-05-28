package telemetry

import (
	"smart-factory/machine-simulator/internal/logging"

	"go.uber.org/zap"
)

func Dispatch(t Telemetry) {

	logging.Logger.Info(
		"telemetry received",
		zap.String("machine", t.MachineName),
		zap.String("state", t.State),

		zap.Float64("temperature", t.Temperature),
		zap.Float64("speed", t.Speed),
		zap.Float64("power", t.Power),

		zap.Int("production_count", t.ProductionCount),

		zap.Float64("runtime_seconds", t.RuntimeSeconds),
		zap.Float64("downtime_seconds", t.DowntimeSeconds),

		zap.Int("error_count", t.ErrorCount),
	)
}
