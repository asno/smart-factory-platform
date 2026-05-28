package simulation

import (
	"context"

	"smart-factory/machine-simulator/internal/logging"

	"go.uber.org/zap"
)

func StartScheduler(
	ctx context.Context,
	registry *Registry,
) {

	for _, machine := range registry.GetMachines() {

		logging.Logger.Info(
			"starting machine runtime",
			zap.String("machine", machine.GetName()),
		)

		go machine.Run(ctx)
	}
}
