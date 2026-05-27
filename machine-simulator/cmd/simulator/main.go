package main

import (
	"time"
	"smart-factory/machine-simulator/internal/logging"
	"smart-factory/machine-simulator/internal/machines"
	"go.uber.org/zap"
)

func main() {

	err := logging.Init()

	if err != nil {
		panic(err)
	}

	defer logging.Logger.Sync()

	conveyor := machines.NewConveyor("Conveyor-01")

	conveyor.Start()

	ticker := time.NewTicker(2 * time.Second)

	defer ticker.Stop()

	for range ticker.C {

		conveyor.Update()

		logging.Logger.Info(
			"machine telemetry",
			zap.String("machine", conveyor.Name),
			zap.String("state", string(conveyor.State)),
			zap.Float64("speed", conveyor.Speed),
			zap.Int("production_count", conveyor.ProductionCount),
		)
	}
}