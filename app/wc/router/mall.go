package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	apis "go-admin/app/wc/apis/mall"
	"go-admin/common/actions"
)

func init() {
	routerMallCheckRole = append(routerMallCheckRole, registerMallApiRouter)
}

// registerMallApiRouter
func registerMallApiRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	//用户验证
	r := v1.Group("warehouse").Use(authMiddleware.MiddlewareFunc())
	{
		r.POST("test", func(context *gin.Context) {
			userId := user.GetUserIdStr(context)
			context.JSON(200, gin.H{"status": "user validated, test success! userId is " + userId})
		})
	}

	apiWarehouse := apis.Warehouse{}
	rWarehouse := v1.Group("/warehouse").Use(authMiddleware.MiddlewareFunc())
	{
		rWarehouse.GET("/list", actions.PermissionAction(), apiWarehouse.GetPageForAuthedUser)
		rWarehouse.GET("", actions.PermissionAction(), apiWarehouse.GetPage)
		rWarehouse.GET("/:id", actions.PermissionAction(), apiWarehouse.Get)
		rWarehouse.POST("", apiWarehouse.Insert)
		rWarehouse.PUT("/:id", actions.PermissionAction(), apiWarehouse.Update)
		rWarehouse.DELETE("", apiWarehouse.Delete)
	}
}
