package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lohasle/nimbus-cloud-framework-go/internal/platform/config"
	"github.com/lohasle/nimbus-cloud-framework-go/internal/platform/httpx"
)

// Auth validates the shared admin JWT and exposes user_id and tenant_id to a
// downstream bounded context. Individual services never trust client headers
// for identity or tenant isolation.
func Auth(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw := strings.TrimSpace(strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer"))
		uid, tenantID, err := parseToken(raw, cfg.JWTSecret)
		if err != nil {
			httpx.Fail(c, http.StatusUnauthorized, 401, "登录状态已失效")
			c.Abort()
			return
		}
		c.Set("user_id", uid)
		c.Set("tenant_id", tenantID)
		c.Next()
	}
}

func parseToken(raw string, secret string) (uint64, uint64, error) {
	token, err := jwt.Parse(raw, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return 0, 0, errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, 0, errors.New("invalid claims")
	}
	subject, err := claims.GetSubject()
	if err != nil {
		return 0, 0, err
	}
	var uid uint64
	if _, err = fmt.Sscan(subject, &uid); err != nil {
		return 0, 0, err
	}
	tenant, ok := claims["tenant_id"].(float64)
	if !ok {
		return 0, 0, errors.New("invalid tenant claim")
	}
	return uid, uint64(tenant), nil
}
