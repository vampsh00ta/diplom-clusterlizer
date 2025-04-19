package main

import (
	"clusterlizer/internal/app/publicapi"
	"flag"
	"fmt"
	"log"
	"os"
)

const localConfigPath = "./local_config.yaml"

func main() {
	configPath := flag.String("f", localConfigPath, "Path to config")
	flag.Parse()

	if _, err := os.Stat(*configPath); os.IsNotExist(err) {
		fmt.Printf("Config file '%s' not found\n", *configPath)
		os.Exit(1)
	}

	// Configuration
	cfg, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	publicapi.Run(cfg)
}

func loadConfig(cfgPath string) (*publicapi.Config, error) {
	cfg, err := publicapi.NewDefaultConfig(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	return cfg, nil
}
