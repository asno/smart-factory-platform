package simulation

import (
	"context"
	"runtime"
	"smart-factory/machine-simulator/internal/logging"
	"time"

	"go.uber.org/zap"
)

func StartRuntimeSupervisor(ctx context.Context) {

	ticker := time.NewTicker(15 * time.Second)

	defer ticker.Stop()

	for {
		select {

		case <-ctx.Done():
			return

		case <-ticker.C:

			logging.Logger.Info(
				"runtime metrics",
				zap.Int("goroutines", runtime.NumGoroutine()),
				zap.Int("cpu_count", runtime.NumCPU()),
			)
		}
	}
}
