package telemetry

import (
	"math/rand"
	"time"
)

type BackoffConfig struct {
	BaseInterval time.Duration
	MaxInterval  time.Duration
	Multiplier   float64
}

func DefaultBackoffConfig() BackoffConfig {
	return BackoffConfig{
		BaseInterval: 1 * time.Second,
		MaxInterval:  60 * time.Second,
		Multiplier:   1.5,
	}
}

func CalculateBackoff(attempt int, cfg BackoffConfig, rng *rand.Rand) time.Duration {
	if attempt <= 0 {
		return cfg.BaseInterval
	}
	if rng == nil {
		rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	factor := 1.0
	for i := 0; i < attempt && factor < float64(cfg.MaxInterval/cfg.BaseInterval); i++ {
		factor *= cfg.Multiplier
	}

	ceiling := float64(cfg.BaseInterval) * factor
	if ceiling > float64(cfg.MaxInterval) {
		ceiling = float64(cfg.MaxInterval)
	}

	jittered := rng.Float64() * ceiling
	if jittered < float64(cfg.BaseInterval) {
		jittered = float64(cfg.BaseInterval)
	}

	return time.Duration(jittered)
}
