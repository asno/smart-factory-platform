package main

import (
	"context"
	"os"
	"os/signal"
	"smart-factory/machine-simulator/internal/alarms"
	"smart-factory/machine-simulator/internal/config"
	"smart-factory/machine-simulator/internal/events"
	"smart-factory/machine-simulator/internal/logging"
	"smart-factory/machine-simulator/internal/machines"
	"smart-factory/machine-simulator/internal/simulation"
	"smart-factory/machine-simulator/internal/telemetry"
	"syscall"
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

	registry := simulation.NewRegistry()

	for _, c := range cfg.Machines.Conveyors {

		conveyor := machines.NewConveyor(
			c.Name,
			eventBus,
		)

		conveyor.Start()
		registry.Add(conveyor)
	}

	for _, o := range cfg.Machines.Ovens {

		oven := machines.NewOven(
			o.Name,
			eventBus,
		)

		oven.Start()
		registry.Add(oven)
	}

	for _, p := range cfg.Machines.Pumps {

		pump := machines.NewPump(
			p.Name,
			eventBus,
		)

		pump.Start()
		registry.Add(pump)
	}

	simulation.StartScheduler(
		ctx,
		registry,
	)

	go simulation.StartRuntimeSupervisor(ctx)

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

			telemetry.Dispatch(t)

		case "alarm":

			a := event.Payload.(alarms.Alarm)

			alarms.Dispatch(a)
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
