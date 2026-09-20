package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg == nil {
		t.Fatal("expected non-nil default config")
	}
	if cfg.ListenAddr != "0.0.0.0:4177" {
		t.Errorf("unexpected listen addr: %s", cfg.ListenAddr)
	}
	if cfg.HubEndpoint != "https://hub.kryvora.network:4188" {
		t.Errorf("unexpected hub endpoint: %s", cfg.HubEndpoint)
	}
	if cfg.HeartbeatInterval != 30*time.Second {
		t.Errorf("unexpected heartbeat interval: %v", cfg.HeartbeatInterval)
	}
	if err := cfg.Validate(); err != nil {
		t.Errorf("default config should be valid: %v", err)
	}
}

func TestConfigValidation(t *testing.T) {
	cases := []struct {
		name    string
		mutate  func(*Config)
		wantErr bool
	}{
		{
			name: "empty listen addr",
			mutate: func(c *Config) {
				c.ListenAddr = ""
			},
			wantErr: true,
		},
		{
			name: "empty hub endpoint",
			mutate: func(c *Config) {
				c.HubEndpoint = ""
			},
			wantErr: true,
		},
		{
			name: "empty data dir",
			mutate: func(c *Config) {
				c.DataDir = ""
			},
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := DefaultConfig()
			tc.mutate(cfg)
			err := cfg.Validate()
			if (err != nil) != tc.wantErr {
				t.Fatalf("Validate() error = %v, wantErr = %v", err, tc.wantErr)
			}
		})
	}
}

func TestConfigLoadEmptyPath(t *testing.T) {
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("expected nil error on empty path, got: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected default config")
	}
}

func TestConfigLoadWithBootstrapToken(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	content := `node:
  listen_addr: "127.0.0.1:4177"
  data_dir: "/tmp/kryvora-test"
  worker_threads: 8

auth:
  bootstrap_token: "kn_genesis_token_test_0x1234"

hub:
  endpoint: "https://hub.kryvora.network:4188"
  heartbeat_interval_sec: 15
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.ListenAddr != "127.0.0.1:4177" {
		t.Errorf("expected 127.0.0.1:4177, got %s", cfg.ListenAddr)
	}
	if cfg.DataDir != "/tmp/kryvora-test" {
		t.Errorf("expected /tmp/kryvora-test, got %s", cfg.DataDir)
	}
	if cfg.WorkerThreads != 8 {
		t.Errorf("expected 8 worker threads, got %d", cfg.WorkerThreads)
	}
	if cfg.BootstrapToken != "kn_genesis_token_test_0x1234" {
		t.Errorf("expected kn_genesis_token_test_0x1234, got %s", cfg.BootstrapToken)
	}
	if cfg.HeartbeatInterval != 15*time.Second {
		t.Errorf("expected 15s interval, got %v", cfg.HeartbeatInterval)
	}
}

func TestConfigLoadEnvOverride(t *testing.T) {
	os.Setenv("KRYVORA_BOOTSTRAP_TOKEN", "kn_env_override_token_999")
	defer os.Unsetenv("KRYVORA_BOOTSTRAP_TOKEN")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.BootstrapToken != "kn_env_override_token_999" {
		t.Errorf("expected env token, got %s", cfg.BootstrapToken)
	}
}
