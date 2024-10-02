package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	"go-admin/common/actions"
	"go-admin/common/middleware"

	apis "go-admin/app/oc/apis/admin"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerSysAdminApiRouter)
	routerNoCheckRole = append(routerNoCheckRole, registerSysAdminNoCheckApiRouter)
}

// registerSysApiRouter
func registerSysAdminApiRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	apiStatements := apis.Statements{}
	apiOTS := apis.OrderToStatements{}
	r := v1.Group("/statements").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", actions.PermissionAction(), apiStatements.GetPage)
		r.GET("/:id", actions.PermissionAction(), apiStatements.Get)
		r.GET("/details/:id", actions.PermissionAction(), apiOTS.GetDetails)
		r.GET("/details/export/:id", actions.PermissionAction(), apiOTS.DetailsExport)
	}

	apiOderInfo := apis.OrderInfo{}
	r = v1.Group("/order").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", actions.PermissionAction(), apiOderInfo.GetPage)
		r.GET("/:id", actions.PermissionAction(), apiOderInfo.Get)
		r.POST("", apiOderInfo.Insert)
		r.PUT("/:id", actions.PermissionAction(), apiOderInfo.Update)
		r.PUT("/shipping/:id", actions.PermissionAction(), apiOderInfo.UpdateShipping)
		r.PUT("/product/:id", actions.PermissionAction(), apiOderInfo.UpdateProduct)
		r.PUT("/receipt/:id", actions.PermissionAction(), apiOderInfo.Receipt)

		r.GET("/receipt-image/:orderId", actions.PermissionAction(), apiOderInfo.GetReceiptImage)
		r.PUT("/receipt-image/:orderId", actions.PermissionAction(), apiOderInfo.SaveReceiptImage)

		r.PUT("/cancel/:id", actions.PermissionAction(), apiOderInfo.Cancel)
		r.PUT("/confirm", actions.PermissionAction(), apiOderInfo.Confirm)
		r.POST("/addProduct", actions.PermissionAction(), apiOderInfo.AddProduct)
		r.POST("/check-over-budget", actions.PermissionAction(), apiOderInfo.CheckIsOverBudget)
	}

	// 订单日志
	apiOderInfoLog := apis.OrderInfoLog{}
	r = v1.Group("/order-info-log").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", actions.PermissionAction(), apiOderInfoLog.GetPage)
		r.GET("/:id", actions.PermissionAction(), apiOderInfoLog.Get)
	}

	apiCsApply := apis.CsApply{}
	r = v1.Group("/admin/cs-apply").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", actions.PermissionAction(), apiCsApply.GetPage)
		r.GET("/:csNo", actions.PermissionAction(), apiCsApply.Get)
		r.POST("", actions.PermissionAction(), apiCsApply.Insert)
		//r.PUT("/:id", actions.PermissionAction(), api.Update)
		//批量审核
		r.PUT("all-audit", actions.PermissionAction(), apiCsApply.AllAudit)
		//关闭售后单
		r.PUT("cancel", actions.PermissionAction(), apiCsApply.Cancel)
		//拒绝售后单
		r.PUT("un-do", actions.PermissionAction(), apiCsApply.Undo)
		//确认售后单
		r.PUT("confirm", actions.PermissionAction(), apiCsApply.Confirm)

		r.GET("sale-products/:orderId", actions.PermissionAction(), apiCsApply.GetSaleProducts)
		r.DELETE("", apiCsApply.Delete)
	}

	apiCsApplyDetail := apis.CsApplyDetail{}
	r = v1.Group("/admin/cs-apply-detail").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("cs-apply-info/:csNo", actions.PermissionAction(), apiCsApplyDetail.GetByCsNo)
	}

	apiCsApplyLog := apis.CsApplyLog{}
	r = v1.Group("/admin/cs-apply-log").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", actions.PermissionAction(), apiCsApplyLog.GetPage)
		//r.GET("/:id", actions.PermissionAction(), api.Get)
		//r.POST("", api.Insert)
		//r.PUT("/:id", actions.PermissionAction(), api.Update)
		//r.DELETE("", api.Delete)
	}
}

func registerSysAdminNoCheckApiRouter(v1 *gin.RouterGroup) {
	apiStatements := apis.Statements{}
	r := v1.Group("/statements")
	{
		r.GET("/init", apiStatements.InitStatements)
	}

	apiOderInfo := apis.OrderInfo{}
	r = v1.Group("/order")
	{
		r.GET("/auto-confirm", apiOderInfo.AutoConfirm)
		r.GET("/auto-sign-for", apiOderInfo.AutoSignFor)
		r.GET("/out-of-stock", apiOderInfo.AutOfStock)
	}
}
