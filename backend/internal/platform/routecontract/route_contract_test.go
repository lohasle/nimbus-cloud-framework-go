package routecontract

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lohasle/nimbus-cloud-framework-go/internal/modules/infra"
	"github.com/lohasle/nimbus-cloud-framework-go/internal/modules/member"
	"github.com/lohasle/nimbus-cloud-framework-go/internal/modules/pay"
	"github.com/lohasle/nimbus-cloud-framework-go/internal/modules/system"
	"github.com/lohasle/nimbus-cloud-framework-go/internal/platform/config"
)

func TestVisibleAdminRouteContract(t *testing.T) {
	gin.SetMode(gin.TestMode)
	routes := map[string]struct{}{}
	addRoutes(routes, system.New(system.NewHandler(system.NewService(nil, config.Config{})), nil))

	for _, register := range []func(*gin.RouterGroup, gin.HandlerFunc){
		func(group *gin.RouterGroup, auth gin.HandlerFunc) { infra.Register(group, nil, auth) },
		func(group *gin.RouterGroup, auth gin.HandlerFunc) { member.Register(group, nil, auth) },
		func(group *gin.RouterGroup, auth gin.HandlerFunc) { pay.Register(group, nil, auth) },
	} {
		engine := gin.New()
		register(engine.Group("/admin-api"), func(c *gin.Context) { c.Next() })
		addRoutes(routes, engine)
	}

	want := []string{
		"POST /admin-api/system/auth/login",
		"GET /admin-api/system/user/page", "GET /admin-api/system/user/list",
		"GET /admin-api/system/user/export-excel", "POST /admin-api/system/user/import",
		"GET /admin-api/system/role/simple-list", "GET /admin-api/system/permission/list-user-roles",
		"GET /admin-api/system/area/tree",
		"GET /admin-api/infra/config/page", "GET /admin-api/infra/config/export-excel",
		"GET /admin-api/infra/file-config/page", "GET /admin-api/infra/api-access-log/page",
		"GET /admin-api/member/user/page", "GET /admin-api/member/point/record/page",
		"GET /admin-api/member/experience-record/page", "GET /admin-api/member/sign-in/record/page",
		"GET /admin-api/member/address/list",
		"GET /admin-api/pay/app/page", "GET /admin-api/pay/channel/get",
		"GET /admin-api/pay/order/page", "GET /admin-api/pay/refund/page",
		"GET /admin-api/pay/wallet/page", "GET /admin-api/pay/wallet-transaction/page",
	}
	for _, route := range want {
		if _, ok := routes[route]; !ok {
			t.Errorf("missing visible admin route %s", route)
		}
	}
}

func addRoutes(target map[string]struct{}, engine *gin.Engine) {
	for _, route := range engine.Routes() {
		target[route.Method+" "+route.Path] = struct{}{}
	}
}
