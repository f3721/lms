package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"

	apis "go-admin/app/wc/apis/admin"
	"go-admin/common/actions"
	"go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerReminderListSkuRouter)
}

// registerReminderListSkuRouter
func registerReminderListSkuRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.ReminderListSku{}
	r := v1.Group("/admin/reminder-list-sku").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", actions.PermissionAction(), api.GetPage)
		//r.GET("/:id", actions.PermissionAction(), api.Get)
		//r.POST("", api.Insert)
		//r.PUT("/:id", actions.PermissionAction(), api.Update)
		//r.DELETE("", api.Delete)
	}
}
