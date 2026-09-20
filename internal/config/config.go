package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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
	BootstrapToken    string        `json:"bootstrap_token"`
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
		BootstrapToken:    "",
	}
}

func Load(path string) (*Config, error) {
	cfg := DefaultConfig()

	// Environment variable fallback/initial value
	if token := os.Getenv("KRYVORA_BOOTSTRAP_TOKEN"); token != "" {
		cfg.BootstrapToken = token
	}

	if path == "" {
		return cfg, nil
	}

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("failed to open config: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		val = strings.Trim(val, `"'`)

		switch key {
		case "listen_addr":
			if val != "" {
				cfg.ListenAddr = val
			}
		case "data_dir":
			if val != "" {
				cfg.DataDir = val
			}
		case "endpoint", "hub_endpoint":
			if val != "" {
				cfg.HubEndpoint = val
			}
		case "bootstrap_token":
			if val != "" {
				cfg.BootstrapToken = val
			}
		case "worker_threads":
			if n, err := strconv.Atoi(val); err == nil && n > 0 {
				cfg.WorkerThreads = n
			}
		case "heartbeat_interval_sec":
			if n, err := strconv.Atoi(val); err == nil && n > 0 {
				cfg.HeartbeatInterval = time.Duration(n) * time.Second
			}
		case "enabled":
			if val == "false" {
				cfg.TelemetryEnabled = false
			} else if val == "true" {
				cfg.TelemetryEnabled = true
			}
		}
	}

	// Environment variable overrides file if provided
	if token := os.Getenv("KRYVORA_BOOTSTRAP_TOKEN"); token != "" {
		cfg.BootstrapToken = token
	}

	return cfg, scanner.Err()
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
