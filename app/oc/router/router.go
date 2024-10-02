package router

import (
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
)

var (
	routerNoCheckRole     = make([]func(*gin.RouterGroup), 0)
	routerCheckRole       = make([]func(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware), 0)
	routerInner           = make([]func(*gin.RouterGroup), 0)
	routerMallNoCheckRole = make([]func(*gin.RouterGroup), 0)
	routerMallCheckRole   = make([]func(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware), 0)
)

// InitAdminRouter 后台用户路由管理
func InitAdminRouter(r *gin.Engine, authMiddleware *jwt.GinJWTMiddleware) *gin.Engine {
	// 无需认证的路由
	examplesNoCheckRoleRouter(r)
	// 需要认证的路由
	examplesCheckRoleRouter(r, authMiddleware)
	// 服务内部路由
	examplesInnerRouter(r)

	return r
}

// InitMallRouter 前台用户路由管理
func InitMallRouter(r *gin.Engine, authMiddleware *jwt.GinJWTMiddleware) *gin.Engine {

	// 前端无需认证的路由
	examplesMallNoCheckRoleRouter(r)
	// 前端需要认证的路由
	examplesMallCheckRoleRouter(r, authMiddleware)

	return r
}

// 无需认证的路由示例
func examplesNoCheckRoleRouter(r *gin.Engine) {
	// 可根据业务需求来设置接口版本
	v1 := r.Group("/api/v1/oc")
	for _, f := range routerNoCheckRole {
		f(v1)
	}

	// 健康检测
	v1.GET("/health", func(c *gin.Context) {
		c.String(200, "health")
	})
}

// 需要认证的路由示例
func examplesCheckRoleRouter(r *gin.Engine, authMiddleware *jwtauth.GinJWTMiddleware) {
	// 可根据业务需求来设置接口版本
	v1 := r.Group("/api/v1/oc")
	for _, f := range routerCheckRole {
		f(v1, authMiddleware)
	}
}

// 服务内部路由示例
func examplesInnerRouter(r *gin.Engine) {
	v1 := r.Group("/inner/oc")
	for _, f := range routerInner {
		f(v1)
	}
}

// 前端 无需认证的路由示例
func examplesMallNoCheckRoleRouter(r *gin.Engine) {
	// 可根据业务需求来设置接口版本
	v1 := r.Group("/api/v1/oc/mall")
	for _, f := range routerMallNoCheckRole {
		f(v1)
	}
}

// 前端 需要认证的路由示例
func examplesMallCheckRoleRouter(r *gin.Engine, authMiddleware *jwtauth.GinJWTMiddleware) {
	// 可根据业务需求来设置接口版本
	v1 := r.Group("/api/v1/oc/mall")
	for _, f := range routerMallCheckRole {
		f(v1, authMiddleware)
	}
}
