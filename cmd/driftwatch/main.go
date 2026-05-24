package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/driftwatch/driftwatch/internal/tfstate"
)

func main() {
	statePath := flag.String("state", "terraform.tfstate", "path to Terraform state file")
	flag.Parse()

	if _, err := os.Stat(*statePath); os.IsNotExist(err) {
		log.Fatalf("state file not found: %s", *statePath)
	}

	state, err := tfstate.ParseFile(*statePath)
	if err != nil {
		log.Fatalf("failed to parse state: %v", err)
	}

	fmt.Printf("Terraform state version: %d\n", state.Version)
	fmt.Printf("Resources found: %d\n", len(state.Resources))
	for _, r := range state.Resources {
		fmt.Printf("  [%s] %s\n", r.Type, r.Name)
	}
}
