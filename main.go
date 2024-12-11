package main

import (
	"github.com/guarzo/canifly/internal/cmd"
	"log"
)

func main() {
	if err := cmd.Start(); err != nil {
		log.Fatalf("Application failed to start: %v", err)
	}
}
