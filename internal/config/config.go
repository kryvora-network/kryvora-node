package config

import (
	"fmt"
	"time"
)

type Config struct {
	NodeID            string        `json:"node_id"`
	ListenAddr        string        `json:"listen_addr"`
	DataDir           string        `json:"data_dir"`
	HubEndpoint       string        `json:"hub_endpoint"`
	HeartbeatInterval time.Duration `json:"heartbeat_interval"`
	WorkerThreads     int           `json:"worker_threads"`
	TelemetryEnabled  bool          `json:"telemetry_enabled"`
}

func DefaultConfig() *Config {
	return &Config{
		NodeID:            "",
		ListenAddr:        "0.0.0.0:4177",
		DataDir:           "/var/lib/kryvora",
		HubEndpoint:       "https://hub.kryvora.network:4188",
		HeartbeatInterval: 30 * time.Second,
		WorkerThreads:     4,
		TelemetryEnabled:  true,
	}
}

func (c *Config) Validate() error {
	if c.ListenAddr == "" {
		return fmt.Errorf("listen_addr must be specified")
	}
	if c.HubEndpoint == "" {
		return fmt.Errorf("hub_endpoint must be specified")
	}
	if c.DataDir == "" {
		return fmt.Errorf("data_dir must be specified")
	}
	return nil
}
