package utils

import (
	"context"
	"time"
)

func InvokeFunctionWithInterval(ctx context.Context, duration time.Duration, functionToInvoke func()) {
	ticker := time.NewTicker(duration)
	for {
		select {
		case <-ticker.C:
			functionToInvoke()
		case <-ctx.Done():
			return
		}

	}
}
