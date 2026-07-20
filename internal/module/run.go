package module

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lohasle/nimbus-cloud-framework-go/internal/discovery"
	"github.com/lohasle/nimbus-cloud-framework-go/internal/httpx"
	"github.com/lohasle/nimbus-cloud-framework-go/internal/middleware"
)

func Run(name string, port uint64) {
	registry, err := discovery.New()
	if err != nil {
		slog.Error("service discovery initialization failed", "error", err)
		os.Exit(1)
	}
	if err = registry.Register("nimbus-"+name, port); err != nil {
		slog.Error("service registration failed", "error", err)
		os.Exit(1)
	}
	defer registry.Deregister("nimbus-"+name, port)

	r := gin.New()
	_ = r.SetTrustedProxies([]string{"127.0.0.1", "::1"})
	r.Use(gin.Recovery(), middleware.CORS(), middleware.RequestContext())
	r.GET("/health", Health(name))
	r.GET("/admin-api/"+name+"/health", Health(name))
	if name == "business" {
		for _, boundary := range []string{"application", "im", "app"} {
			r.GET("/admin-api/"+boundary+"/health", Health(boundary))
		}
	}

	server := &http.Server{Addr: ":" + strconv.FormatUint(port, 10), Handler: r, ReadHeaderTimeout: 5 * time.Second}
	go func() {
		if serveErr := server.ListenAndServe(); serveErr != nil && serveErr != http.ErrServerClosed {
			slog.Error("module server failed", "service", name, "error", serveErr)
			os.Exit(1)
		}
	}()
	slog.Info("module started", "service", name, "port", port)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)
}

// Health godoc
// @Summary Module health
// @Description Health endpoint for deployment probes and service-governance checks.
// @Tags Health
// @Produce json
// @Success 200 {object} httpx.Response
// @Router /health [get]
// @Router /admin-api/infra/health [get]
// @Router /admin-api/member/health [get]
// @Router /admin-api/pay/health [get]
// @Router /admin-api/business/health [get]
// @Router /admin-api/application/health [get]
// @Router /admin-api/im/health [get]
// @Router /admin-api/app/health [get]
func Health(name string) gin.HandlerFunc {
	return func(c *gin.Context) { httpx.OK(c, gin.H{"status": "UP", "service": name}) }
}
