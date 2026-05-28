package simulation

import (
	"context"
	"smart-factory/machine-simulator/internal/logging"
	"time"

	"go.uber.org/zap"
)

func StartScheduler(
	ctx context.Context,
	registry *Registry,
) {

	for _, machine := range registry.GetMachines() {

		go runMachine(
			ctx,
			machine,
		)
	}
}

func runMachine(
	ctx context.Context,
	machine interface {
		Cycle()
		Tick() time.Duration
		GetName() string
	},
) {

	ticker := time.NewTicker(machine.Tick())

	defer ticker.Stop()

	logging.Logger.Info(
		"machine runtime started",
		zap.String("machine", machine.GetName()),
	)

	for {
		select {

		case <-ctx.Done():
			return

		case <-ticker.C:

			machine.Cycle()
		}
	}
}
