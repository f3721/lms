package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	"go-admin/common/middleware"

	apis "go-admin/app/wc/apis/admin"
	"go-admin/common/actions"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerReminderRuleRouter)
}

// registerReminderRuleRouter
func registerReminderRuleRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.ReminderRule{}
	r := v1.Group("/admin/reminder-rule").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", actions.PermissionAction(), api.GetPage)
		r.GET("/:id", actions.PermissionAction(), api.Get)
		r.POST("", actions.PermissionAction(), api.Insert)
		r.PUT("/:id", actions.PermissionAction(), api.Update)
		//r.DELETE("", api.Delete)
	}
}
