package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	apis "go-admin/app/oc/apis/mall"
	"go-admin/common/actions"
)

func init() {
	routerMallCheckRole = append(routerMallCheckRole, registerMallApiRouter)
	routerMallNoCheckRole = append(routerMallNoCheckRole, registerMalNoCheckApiRouter)
}

// registerMallApiRouter
func registerMallApiRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	apiOderInfo := apis.OrderInfo{}
	r := v1.Group("/order").Use(authMiddleware.MiddlewareFunc())
	{
		r.GET("", actions.PermissionAction(), apiOderInfo.GetPage)
		r.GET("/count", actions.PermissionAction(), apiOderInfo.GetPageCount)
		r.GET("/:id", actions.PermissionAction(), apiOderInfo.Get)
		r.PUT("/update-po/:id", actions.PermissionAction(), apiOderInfo.UpdatePo)
		r.PUT("/cancel/:id", actions.PermissionAction(), apiOderInfo.Cancel)
		r.PUT("/buy-again/:id", actions.PermissionAction(), apiOderInfo.BuyAgain)
		r.GET("/export", actions.PermissionAction(), apiOderInfo.Export)
		r.PUT("/delete/:id", actions.PermissionAction(), apiOderInfo.Delete)
		r.POST("", apiOderInfo.Insert)
		r.PUT("/receipt/:id", actions.PermissionAction(), apiOderInfo.Receipt)
	}
	apiPreso := apis.Preso{}
	r = v1.Group("/preso").Use(authMiddleware.MiddlewareFunc())
	{
		r.POST("/submit-approval", apiPreso.SubmitApproval)
		r.GET("", actions.PermissionAction(), apiPreso.GetPage)
		r.GET("/count", actions.PermissionAction(), apiPreso.GetPageCount)
		r.GET("/count-approval", actions.PermissionAction(), apiPreso.GetApprovePageCount)
		r.GET("/:id", actions.PermissionAction(), apiPreso.Get)
		r.POST("/finish-approval", actions.PermissionAction(), apiPreso.FinishApproval)
		r.POST("/batch-approval", actions.PermissionAction(), apiPreso.BatchApproval)
		r.POST("/withdraw", actions.PermissionAction(), apiPreso.Withdraw)
		r.POST("/save-file", actions.PermissionAction(), apiPreso.SaveFile)
		r.GET("/export/:id", actions.PermissionAction(), apiPreso.Export)
		r.PUT("/buy-again/:id", actions.PermissionAction(), apiPreso.BuyAgain)
		r.DELETE("/file", actions.PermissionAction(), apiPreso.DeleteFile)
	}

	apiCsApply := apis.CsApply{}
	r = v1.Group("/cs-apply").Use(authMiddleware.MiddlewareFunc())
	{
		r.GET("", actions.PermissionAction(), apiCsApply.GetPage)
		r.GET("/:id", actions.PermissionAction(), apiCsApply.Get)
		r.POST("", apiCsApply.Insert)
		r.PUT("/:id", actions.PermissionAction(), apiCsApply.Update)
		r.DELETE("", apiCsApply.Delete)

		// 关闭售后单
		r.PUT("cancel", actions.PermissionAction(), apiCsApply.Cancel)
		r.GET("group-by-order", actions.PermissionAction(), apiCsApply.GetPageGroupByOrder)
		r.GET("get-page-details", actions.PermissionAction(), apiCsApply.GetPageDetails)
	}
}

func registerMalNoCheckApiRouter(v1 *gin.RouterGroup) {
	apiPreso := apis.Preso{}
	r := v1.Group("/preso")
	{
		r.GET("/expire", actions.PermissionAction(), apiPreso.Expire)
		r.GET("/cron-approve", actions.PermissionAction(), apiPreso.CronApprove)
	}
}
