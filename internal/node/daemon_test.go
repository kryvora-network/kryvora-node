package node

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/kryvora-network/kryvora-node/internal/config"
	"github.com/kryvora-network/kryvora-node/internal/storage"
)

func TestDaemonLifecycle(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "daemon-test-*")
	if err != nil {
		t.Fatalf("temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	store, err := storage.Open(tempDir)
	if err != nil {
		t.Fatalf("storage open: %v", err)
	}
	defer store.Close()

	cfg := config.DefaultConfig()
	cfg.HeartbeatInterval = 20 * time.Millisecond
	cfg.BootstrapToken = "kn_test_genesis_bootstrap_token_0xBCf2cD"

	d := New(cfg, store)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err = d.Start(ctx)
	if err != nil && err != context.DeadlineExceeded {
		t.Fatalf("unexpected start error: %v", err)
	}

	if err := d.Stop(); err != nil {
		t.Errorf("stop error: %v", err)
	}
}

func TestDaemonMissingAuthFails(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "daemon-auth-test-*")
	if err != nil {
		t.Fatalf("temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	store, err := storage.Open(tempDir)
	if err != nil {
		t.Fatalf("storage open: %v", err)
	}
	defer store.Close()

	cfg := config.DefaultConfig()
	cfg.BootstrapToken = "" // unauthenticated

	d := New(cfg, store)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err = d.Start(ctx)
	if err == nil {
		t.Fatal("expected daemon to fail without bootstrap token, but it started successfully")
	}
}
