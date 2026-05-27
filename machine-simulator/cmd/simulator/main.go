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
	"smart-factory/machine-simulator/internal/config"
	"smart-factory/machine-simulator/internal/simulation"
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
	cfg, err := config.Load("configs/simulator.yaml")

	if err != nil {
		panic(err)
	}

	for _, c := range cfg.Machines.Conveyors {

		conveyor := machines.NewConveyor(
			c.Name,
			eventBus,
		)

		conveyor.Start()

		go conveyor.Run(ctx)
	}

	for _, o := range cfg.Machines.Ovens {

		oven := machines.NewOven(
			o.Name,
			eventBus,
		)

		oven.Start()

		go oven.Run(ctx)
	}

	for _, p := range cfg.Machines.Pumps {

		pump := machines.NewPump(
			p.Name,
			eventBus,
		)

		pump.Start()

		go pump.Run(ctx)
	}

	go simulation.StartHeartbeat(
		ctx,
		"machine-simulator",
	)
	

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