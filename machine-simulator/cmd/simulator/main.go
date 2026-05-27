package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"smart-factory/machine-simulator/internal/alarms"
	"smart-factory/machine-simulator/internal/events"
	"smart-factory/machine-simulator/internal/logging"
	"smart-factory/machine-simulator/internal/machines"
	"smart-factory/machine-simulator/internal/telemetry"
	"go.uber.org/zap"
)

func main() {

	err := logging.Init()

	if err != nil {
		panic(err)
	}

	defer logging.Logger.Sync()

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	eventBus := events.NewBus(100)

	conveyor := machines.NewConveyor(
		"Conveyor-01",
		eventBus,
	)

	conveyor.Start()

	go conveyor.Run(ctx)

	go processEvents(eventBus)

	waitForShutdown(cancel)
}

func processEvents(eventBus *events.Bus) {

	for event := range eventBus.Events {

		switch event.Type {

		case "telemetry":

			t := event.Payload.(telemetry.Telemetry)

			logging.Logger.Info(
				"telemetry received",
				zap.String("machine", t.MachineName),
				zap.String("state", t.State),
				zap.Float64("temperature", t.Temperature),
				zap.Float64("speed", t.Speed),
				zap.Float64("power", t.Power),
				zap.Int("production_count", t.ProductionCount),
			)

		case "alarm":

			a := event.Payload.(alarms.Alarm)

			logging.Logger.Warn(
				"alarm triggered",
				zap.String("machine", a.MachineName),
				zap.String("message", a.Message),
				zap.String("severity", string(a.Severity)),
			)
		}
	}
}

func waitForShutdown(cancel context.CancelFunc) {

	sigChan := make(chan os.Signal, 1)

	signal.Notify(
		sigChan,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-sigChan

	logging.Logger.Info("shutdown signal received")

	cancel()
}