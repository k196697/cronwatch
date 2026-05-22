package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/user/cronwatch/internal/config"
	"github.com/user/cronwatch/internal/watcher"
)

const defaultConfigPath = "config.yaml"

func main() {
	configPath := flag.String("config", defaultConfigPath, "path to configuration file")
	version := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *version {
		fmt.Println("cronwatch v0.1.0")
		os.Exit(0)
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	w, err := watcher.New(cfg)
	if err != nil {
		log.Fatalf("failed to initialise watcher: %v", err)
	}

	if err := w.Run(); err != nil {
		log.Fatalf("watcher exited with error: %v", err)
	}
}
