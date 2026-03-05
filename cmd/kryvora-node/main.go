package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kryvora-network/kryvora-node/internal/config"
	"github.com/kryvora-network/kryvora-node/internal/node"
	"github.com/kryvora-network/kryvora-node/internal/storage"
)

const Version = "0.2.0"

func main() {
	configPath := flag.String("config", "config.example.yaml", "Path to configuration file")
	port := flag.String("port", "4177", "Local listen port")
	showVersion := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("kryvora-node version %s\n", Version)
		os.Exit(0)
	}

	cfg := config.DefaultConfig()
	if *port != "" {
		cfg.ListenAddr = "0.0.0.0:" + *port
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("invalid configuration: %v", err)
	}

	store, err := storage.Open(cfg.DataDir)
	if err != nil {
		log.Fatalf("failed to initialize storage: %v", err)
	}
	defer store.Close()

	daemon := node.New(cfg, store)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	log.Printf("Starting kryvora-node v%s", Version)
	if err := daemon.Start(ctx); err != nil {
		log.Fatalf("daemon error: %v", err)
	}
}
