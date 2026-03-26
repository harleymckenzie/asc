package awsutil

import (
	"context"
	"fmt"
	"time"
)

// WaitConfig configures the polling behaviour for WaitForStatus.
type WaitConfig struct {
	ResourceName string
	PollInterval time.Duration
	MaxWait      time.Duration
	StatusFunc   func(ctx context.Context) (string, error)
	IsTerminal   func(status string) bool
}

// WaitForStatus polls StatusFunc at PollInterval until IsTerminal returns true
// or MaxWait is exceeded. Returns the final status.
func WaitForStatus(ctx context.Context, config WaitConfig) (string, error) {
	start := time.Now()
	ticker := time.NewTicker(config.PollInterval)
	defer ticker.Stop()

	// Check immediately before first tick
	status, err := config.StatusFunc(ctx)
	if err != nil {
		return "", fmt.Errorf("get status: %w", err)
	}
	fmt.Printf("Status: %s (0s elapsed)\n", status)
	if config.IsTerminal(status) {
		return status, nil
	}

	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-ticker.C:
			elapsed := time.Since(start).Truncate(time.Second)
			if elapsed > config.MaxWait {
				return "", fmt.Errorf("timeout after %s waiting for %s", config.MaxWait, config.ResourceName)
			}
			status, err := config.StatusFunc(ctx)
			if err != nil {
				return "", fmt.Errorf("get status: %w", err)
			}
			fmt.Printf("Status: %s (%s elapsed)\n", status, elapsed)
			if config.IsTerminal(status) {
				return status, nil
			}
		}
	}
}
