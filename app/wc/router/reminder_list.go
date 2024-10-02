package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	apis "go-admin/app/wc/apis/admin"
	"go-admin/common/actions"
	"go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerReminderListRouter)
}

// registerReminderListRouter
func registerReminderListRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.ReminderList{}
	r := v1.Group("/admin/reminder-list").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", actions.PermissionAction(), api.GetPage)
		r.GET("/:id", actions.PermissionAction(), api.Get)
		r.GET("/export/:id", actions.PermissionAction(), api.Export)

		//r.POST("", api.Insert)
		//r.PUT("/:id", actions.PermissionAction(), api.Update)
		//r.DELETE("", api.Delete)
	}

	r = v1.Group("/admin/reminder-list")
	{
		r.GET("create", actions.PermissionAction(), api.Create)
	}

}
