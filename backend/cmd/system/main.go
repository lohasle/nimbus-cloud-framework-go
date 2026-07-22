package main

import (
	"log/slog"
	"os"

	_ "github.com/lohasle/nimbus-cloud-framework-go/docs"
	"github.com/lohasle/nimbus-cloud-framework-go/internal/modules/system"
	"github.com/lohasle/nimbus-cloud-framework-go/internal/platform/config"
	"github.com/lohasle/nimbus-cloud-framework-go/internal/platform/database"
	"github.com/lohasle/nimbus-cloud-framework-go/internal/platform/discovery"
)

// @title Nimbus Cloud Framework Go System API
// @version 1.0
// @description Tenant, operations-user and authentication service.
// @BasePath /admin-api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg := config.Load()
	if os.Getenv("NIMBUS_HTTP_ADDR") == "" {
		cfg.HTTPAddr = ":58081"
	}
	db, err := database.Open(cfg)
	if err != nil {
		slog.Error("database initialization failed", "error", err)
		os.Exit(1)
	}
	registry, err := discovery.New()
	if err != nil || registry.Register("nimbus-system", 58081) != nil {
		slog.Error("service registration failed", "error", err)
		os.Exit(1)
	}
	defer registry.Deregister("nimbus-system", 58081)
	service := system.NewService(db, cfg)
	if err = system.New(system.NewHandler(service), db).Run(cfg.HTTPAddr); err != nil {
		slog.Error("system service stopped", "error", err)
		os.Exit(1)
	}
}
