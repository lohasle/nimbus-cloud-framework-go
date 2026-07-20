package main

import (
	"log/slog"
	"os"

	"github.com/lohasle/nimbus-cloud-framework-go/internal/gateway"
)

// @title Nimbus Cloud Framework Go Gateway
// @version 1.0
// @description Public gateway for the Go microservice scaffold.
func main() {
	r, err := gateway.New()
	if err != nil {
		slog.Error("gateway initialization failed", "error", err)
		os.Exit(1)
	}
	if err = r.Run(":58080"); err != nil {
		slog.Error("gateway stopped", "error", err)
		os.Exit(1)
	}
}
