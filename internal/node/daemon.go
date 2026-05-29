package node

import (
	"context"
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
	d.mu.Lock()
	d.running = true
	d.mu.Unlock()

	log.Printf("[node] daemon started on %s", d.cfg.ListenAddr)

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
	log.Printf("[node] sending heartbeat ping to %s", d.cfg.HubEndpoint)
	// Transient network retry check
	if !d.running {
		return
	}
}
