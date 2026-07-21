package main

import (
	"log/slog"
	"os"

	_ "github.com/lohasle/nimbus-cloud-framework-go/docs"
	"github.com/lohasle/nimbus-cloud-framework-go/internal/modules/member"
	"github.com/lohasle/nimbus-cloud-framework-go/internal/platform/config"
	"github.com/lohasle/nimbus-cloud-framework-go/internal/platform/database"
	"github.com/lohasle/nimbus-cloud-framework-go/internal/platform/middleware"
	"github.com/lohasle/nimbus-cloud-framework-go/internal/platform/modulehost"
)

// @title Nimbus Cloud Framework Go Member API
// @version 1.0
// @description APP member, level, group, tag, points and experience service.
// @BasePath /admin-api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg := config.Load()
	db, err := database.Connect(cfg)
	if err == nil {
		err = member.Migrate(db)
	}
	if err == nil {
		err = member.Seed(db, 1)
	}
	if err != nil {
		slog.Error("member database initialization failed", "error", err)
		os.Exit(1)
	}
	r := modulehost.New(member.ModuleName)
	member.Register(r.Group("/admin-api"), db, middleware.Auth(cfg))
	modulehost.Serve(member.ModuleName, 58087, r)
}
