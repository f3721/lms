package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	"go-admin/app/admin/apis"
	"go-admin/common/middleware"
)

func init() {
	routerMallCheckRole = append(routerMallCheckRole, registerMallApiRouter)
	routerMallNoCheckRole = append(routerMallNoCheckRole, registerMalNoCheckApiRouter)
}

// registerMallApiRouter
func registerMallApiRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {

	dataApi := apis.SysDictData{}
	opSelect := v1.Group("/dict-data").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		opSelect.GET("/option-select", dataApi.GetAll)
	}
}

func registerMalNoCheckApiRouter(v1 *gin.RouterGroup) {
	//apiPreso := apis.Preso{}
	//r := v1.Group("/preso")
	//{
	//	r.GET("/expire", actions.PermissionAction(), apiPreso.Expire)
	//}
}
