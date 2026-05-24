package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/yourorg/driftwatch/internal/config"
	"github.com/yourorg/driftwatch/internal/drift"
	"github.com/yourorg/driftwatch/internal/provider"
	"github.com/yourorg/driftwatch/internal/tfstate"
)

func main() {
	cfgPath := flag.String("config", "driftwatch.yaml", "path to driftwatch config file")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	state, err := tfstate.ParseFile(cfg.Statefile)
	if err != nil {
		log.Fatalf("state: %v", err)
	}

	var p provider.Provider
	switch cfg.Provider {
	case config.ProviderAWS:
		p = provider.NewAWSProvider()
	case config.ProviderGCP:
		p = provider.NewGCPProvider()
	case config.ProviderAzure:
		p = provider.NewAzureProvider()
	default:
		log.Fatalf("unsupported provider: %s", cfg.Provider)
	}

	detector := drift.NewDetector(p)
	report, err := detector.Detect(state)
	if err != nil {
		log.Fatalf("detect: %v", err)
	}

	if cfg.Output == "json" {
		fmt.Println(report.JSON())
	} else {
		fmt.Println(report.Text())
	}

	if report.HasDrift() {
		os.Exit(1)
	}
}
