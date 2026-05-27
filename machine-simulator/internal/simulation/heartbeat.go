package simulation

import (
	"context"
	"time"
	"smart-factory/machine-simulator/internal/logging"
	"go.uber.org/zap"
)

func StartHeartbeat(
	ctx context.Context,
	systemName string,
) {

	ticker := time.NewTicker(10 * time.Second)

	defer ticker.Stop()

	for {
		select {

		case <-ctx.Done():
			return

		case <-ticker.C:

			logging.Logger.Info(
				"system heartbeat",
				zap.String("system", systemName),
				zap.String("status", "ONLINE"),
			)
		}
	}
}