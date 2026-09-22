package telemetry

import (
	"math/rand"
	"testing"
)

func TestCalculateBackoff(t *testing.T) {
	cfg := DefaultBackoffConfig()
	rng := rand.New(rand.NewSource(42))

	for attempt := 0; attempt < 10; attempt++ {
		backoff := CalculateBackoff(attempt, cfg, rng)
		if backoff < cfg.BaseInterval {
			t.Errorf("attempt %d: backoff %v is below base interval %v", attempt, backoff, cfg.BaseInterval)
		}
		if backoff > cfg.MaxInterval {
			t.Errorf("attempt %d: backoff %v exceeds max interval %v", attempt, backoff, cfg.MaxInterval)
		}
	}
}
