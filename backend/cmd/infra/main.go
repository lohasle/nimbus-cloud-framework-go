package main

import (
	"log/slog"
	"os"

	_ "github.com/lohasle/nimbus-cloud-framework-go/docs"
	"github.com/lohasle/nimbus-cloud-framework-go/internal/modules/infra"
	"github.com/lohasle/nimbus-cloud-framework-go/internal/platform/accesslog"
	"github.com/lohasle/nimbus-cloud-framework-go/internal/platform/config"
	"github.com/lohasle/nimbus-cloud-framework-go/internal/platform/database"
	"github.com/lohasle/nimbus-cloud-framework-go/internal/platform/middleware"
	"github.com/lohasle/nimbus-cloud-framework-go/internal/platform/modulehost"
)

// @title Nimbus Cloud Framework Go Infra API
// @version 1.0
// @description Configuration, file storage and infrastructure log service.
// @BasePath /admin-api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg := config.Load()
	db, err := database.Connect(cfg)
	if err == nil {
		err = infra.Migrate(db)
	}
	if err == nil {
		err = infra.Seed(db, 1)
	}
	if err != nil {
		slog.Error("infra database initialization failed", "error", err)
		os.Exit(1)
	}
	r := modulehost.New(infra.ModuleName, accesslog.Recorder(db, "nimbus-infra"))
	infra.Register(r.Group("/admin-api"), db, middleware.Auth(cfg))
	modulehost.Serve(infra.ModuleName, 58082, r)
}
