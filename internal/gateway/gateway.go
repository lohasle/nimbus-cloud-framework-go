package gateway

import (
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lohasle/nimbus-cloud-framework-go/internal/discovery"
	"github.com/lohasle/nimbus-cloud-framework-go/internal/httpx"
	"github.com/lohasle/nimbus-cloud-framework-go/internal/middleware"
)

type Gateway struct {
	registry  *discovery.Registry
	fallbacks map[string]string
}

func New() (*gin.Engine, error) {
	registry, err := discovery.New()
	if err != nil {
		return nil, err
	}
	gateway := &Gateway{registry: registry, fallbacks: map[string]string{
		"system":   env("NIMBUS_SYSTEM_URL", "http://127.0.0.1:58081"),
		"infra":    env("NIMBUS_INFRA_URL", "http://127.0.0.1:58082"),
		"pay":      env("NIMBUS_PAY_URL", "http://127.0.0.1:58085"),
		"member":   env("NIMBUS_MEMBER_URL", "http://127.0.0.1:58087"),
		"business": env("NIMBUS_BUSINESS_URL", "http://127.0.0.1:58090"),
	}}
	r := gin.New()
	_ = r.SetTrustedProxies([]string{"127.0.0.1", "::1"})
	r.Use(gin.Recovery(), middleware.CORS(), middleware.RequestContext())
	r.GET("/health", gateway.health)
	r.Any("/admin-api/*path", gateway.proxy)
	return r, nil
}

// health godoc
// @Summary Gateway health
// @Description Health endpoint for the public API gateway.
// @Tags Health
// @Produce json
// @Success 200 {object} httpx.Response
// @Router /health [get]
func (g *Gateway) health(c *gin.Context) { httpx.OK(c, gin.H{"status": "UP", "service": "gateway"}) }

func (g *Gateway) proxy(c *gin.Context) {
	segments := strings.Split(strings.TrimPrefix(c.Param("path"), "/"), "/")
	if len(segments) == 0 || segments[0] == "" {
		httpx.Fail(c, http.StatusNotFound, 404, "缺少服务路径")
		return
	}
	moduleName := segments[0]
	targetKey := moduleName
	if moduleName == "application" || moduleName == "im" || moduleName == "app" {
		targetKey = "business"
	}
	fallback, ok := g.fallbacks[targetKey]
	if !ok {
		httpx.Fail(c, http.StatusNotFound, 404, "未配置服务:"+moduleName)
		return
	}
	target := g.registry.Resolve("nimbus-"+targetKey, fallback)
	targetURL, err := url.Parse(target)
	if err != nil {
		httpx.Fail(c, http.StatusBadGateway, 502, "服务地址无效")
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, proxyErr error) {
		slog.Error("gateway proxy failed", "target", target, "error", proxyErr)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(`{"code":502,"data":null,"msg":"上游服务不可用"}`))
	}
	proxy.ServeHTTP(c.Writer, c.Request)
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
