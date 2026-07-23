package main

import (
	"log/slog"
	"os"

	_ "github.com/lohasle/nimbus-cloud-framework-go/docs"
	"github.com/lohasle/nimbus-cloud-framework-go/internal/modules/system"
	"github.com/lohasle/nimbus-cloud-framework-go/internal/platform/config"
	"github.com/lohasle/nimbus-cloud-framework-go/internal/platform/database"
	"github.com/lohasle/nimbus-cloud-framework-go/internal/platform/modulehost"
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
	db, err := database.Open(cfg)
	if err != nil {
		slog.Error("database initialization failed", "error", err)
		os.Exit(1)
	}
	service := system.NewService(db, cfg)
	modulehost.Serve(system.ModuleName, 58081, system.New(system.NewHandler(service), db))
}
