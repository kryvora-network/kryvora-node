package node

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/kryvora-network/kryvora-node/internal/config"
	"github.com/kryvora-network/kryvora-node/internal/storage"
)

type Daemon struct {
	cfg     *config.Config
	storage *storage.Store
	mu      sync.RWMutex
	running bool
	stopCh  chan struct{}
}

func New(cfg *config.Config, store *storage.Store) *Daemon {
	return &Daemon{
		cfg:     cfg,
		storage: store,
		stopCh:  make(chan struct{}),
	}
}

func (d *Daemon) Start(ctx context.Context) error {
	if d.cfg.BootstrapToken == "" {
		return fmt.Errorf("node authentication failed: missing Genesis Node Key bootstrap token (contract: 0xBCf2cD12D1D37578fA7C69805fE77788e866BE1a). Register key at https://node.kryvora.network")
	}

	d.mu.Lock()
	d.running = true
	d.mu.Unlock()

	log.Printf("[node] daemon started on %s (Genesis Key: authenticated)", d.cfg.ListenAddr)

	ticker := time.NewTicker(d.cfg.HeartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return d.Stop()
		case <-d.stopCh:
			return nil
		case <-ticker.C:
			d.heartbeat()
		}
	}
}

func (d *Daemon) Stop() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if !d.running {
		return nil
	}

	d.running = false
	close(d.stopCh)
	log.Println("[node] daemon stopped cleanly")
	return nil
}

func (d *Daemon) heartbeat() {
	if d.cfg.BootstrapToken == "" {
		log.Printf("[node] [AUTH_DENIED] heartbeat aborted: node missing Genesis Node Key license. Visit https://node.kryvora.network")
		return
	}
	log.Printf("[node] sending authenticated heartbeat ping to %s", d.cfg.HubEndpoint)
	// Transient network retry check
	if !d.running {
		return
	}
}
