package main

import (
	"log/slog"
	"os"

	_ "github.com/lohasle/nimbus-cloud-framework-go/docs"
	"github.com/lohasle/nimbus-cloud-framework-go/internal/modules/pay"
	"github.com/lohasle/nimbus-cloud-framework-go/internal/platform/config"
	"github.com/lohasle/nimbus-cloud-framework-go/internal/platform/database"
	"github.com/lohasle/nimbus-cloud-framework-go/internal/platform/middleware"
	"github.com/lohasle/nimbus-cloud-framework-go/internal/platform/modulehost"
)

// @title Nimbus Cloud Framework Go Pay API
// @version 1.0
// @description Payment applications, channels, orders and refunds service.
// @BasePath /admin-api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg := config.Load()
	db, err := database.Connect(cfg)
	if err == nil {
		err = pay.Migrate(db)
	}
	if err == nil {
		err = pay.Seed(db, 1)
	}
	if err != nil {
		slog.Error("pay database initialization failed", "error", err)
		os.Exit(1)
	}
	r := modulehost.New(pay.ModuleName)
	pay.Register(r.Group("/admin-api"), db, middleware.Auth(cfg))
	modulehost.Serve(pay.ModuleName, 58085, r)
}
