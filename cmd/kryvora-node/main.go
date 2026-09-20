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

const Version = "0.2.1"

func main() {
	configPath := flag.String("config", "config.example.yaml", "Path to configuration file")
	port := flag.String("port", "4177", "Local listen port")
	token := flag.String("token", "", "Genesis Node Key bootstrap token (overrides config)")
	showVersion := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("kryvora-node version %s\n", Version)
		os.Exit(0)
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}
	if *port != "" {
		cfg.ListenAddr = "0.0.0.0:" + *port
	}
	if *token != "" {
		cfg.BootstrapToken = *token
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("invalid configuration: %v", err)
	}

	if cfg.BootstrapToken == "" {
		log.Fatalf("\n================================================================================\n" +
			"[FATAL] AUTHENTICATION REQUIRED: Genesis Node Key not configured.\n" +
			"--------------------------------------------------------------------------------\n" +
			"Kryvora DePIN nodes require a verified Genesis Node Key on Arbitrum One.\n" +
			"Contract: 0xBCf2cD12D1D37578fA7C69805fE77788e866BE1a\n\n" +
			"To activate your node:\n" +
			"  1. Mint or link your Genesis Node Key at: https://node.kryvora.network\n" +
			"  2. Add your bootstrap token to config.yaml under 'auth.bootstrap_token'\n" +
			"     or set KRYVORA_BOOTSTRAP_TOKEN in your environment.\n" +
			"================================================================================\n")
	}

	store, err := storage.Open(cfg.DataDir)
	if err != nil {
		log.Fatalf("failed to initialize storage: %v", err)
	}
	defer store.Close()

	daemon := node.New(cfg, store)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	log.Printf("Starting kryvora-node v%s (Genesis Key: authenticated)", Version)
	if err := daemon.Start(ctx); err != nil {
		log.Fatalf("daemon error: %v", err)
	}
}
