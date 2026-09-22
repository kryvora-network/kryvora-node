package config

import (
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
