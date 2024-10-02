package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	"go-admin/app/other/apis"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerInittenantRouter)
}

// 需认证的路由代码
func registerInittenantRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.InitTenant{}
	r := v1.Group("/tenant")
	//.Use(authMiddleware.MiddlewareFunc())
	{
		r.POST("/init", api.InitTenant)
	}
}
