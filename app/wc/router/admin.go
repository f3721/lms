package router

import (
	"go-admin/common/actions"
	"go-admin/common/middleware"

	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"

	apis "go-admin/app/wc/apis/admin"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerSysAdminApiRouter)
}

// registerSysApiRouter
func registerSysAdminApiRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {

	// 实体仓
	apiWh := apis.Warehouse{}
	rWh := v1.Group("/admin/warehouse").Use(authMiddleware.MiddlewareFunc())
	{
		rWh.GET("select", actions.PermissionAction(), apiWh.Select)
		rWh.Use(middleware.AuthCheckRole()).GET("", actions.PermissionAction(), apiWh.GetPage)
		rWh.Use(middleware.AuthCheckRole()).GET("/:id", actions.PermissionAction(), apiWh.Get)
		rWh.GET("/code/:warehouseCode", actions.PermissionAction(), apiWh.GetByCode)
		rWh.Use(middleware.AuthCheckRole()).POST("", actions.PermissionAction(), apiWh.Insert)
		rWh.Use(middleware.AuthCheckRole()).PUT("/:id", actions.PermissionAction(), apiWh.Update)
		rWh.Use(middleware.AuthCheckRole()).GET("cw-tree", actions.PermissionAction(), apiWh.GetCompanyWarehouseTree)
	}

	// 逻辑仓
	apiLwh := apis.LogicWarehouse{}
	rLwh := v1.Group("/admin/logic-warehouse").Use(authMiddleware.MiddlewareFunc())
	{
		//rLwh.GET("", actions.PermissionAction(), apiLwh.GetPage)
		rLwh.GET("/select", actions.PermissionAction(), apiLwh.Select)
		rLwh.Use(middleware.AuthCheckRole()).GET("", actions.PermissionAction(), apiLwh.GetPage)
		rLwh.Use(middleware.AuthCheckRole()).GET("/:id", actions.PermissionAction(), apiLwh.Get)
		rLwh.Use(middleware.AuthCheckRole()).POST("", actions.PermissionAction(), apiLwh.Insert)
		rLwh.Use(middleware.AuthCheckRole()).PUT("/:id", actions.PermissionAction(), apiLwh.Update)
	}

	// 货主
	apiV := apis.Vendors{}
	rV := v1.Group("/admin/vendors").Use(authMiddleware.MiddlewareFunc())
	{
		rV.Use(middleware.AuthCheckRole()).GET("", actions.PermissionAction(), apiV.GetPage)
		rV.Use(middleware.AuthCheckRole()).GET("/:id", actions.PermissionAction(), apiV.Get)
		rV.Use(middleware.AuthCheckRole()).POST("", actions.PermissionAction(), apiV.Insert)
		rV.Use(middleware.AuthCheckRole()).PUT("/:id", actions.PermissionAction(), apiV.Update)
		rV.GET("/select", actions.PermissionAction(), apiV.Select)
	}

	//供应商
	apiSu := apis.Supplier{}
	rSu := v1.Group("/admin/supplier").Use(authMiddleware.MiddlewareFunc())
	{
		rSu.Use(middleware.AuthCheckRole()).GET("", actions.PermissionAction(), apiSu.GetPage)
		rSu.Use(middleware.AuthCheckRole()).GET("/:id", actions.PermissionAction(), apiSu.Get)
		rSu.Use(middleware.AuthCheckRole()).POST("", actions.PermissionAction(), apiSu.Insert)
		rSu.Use(middleware.AuthCheckRole()).PUT("/:id", actions.PermissionAction(), apiSu.Update)
		rSu.GET("/select", actions.PermissionAction(), apiSu.Select)
	}

	// 省市区
	apiRegion := apis.Region{}
	rRegion := v1.Group("/admin/region")
	{
		rRegion.GET("", apiRegion.GetPage)
		rRegion.GET("/:id", apiRegion.Get)
	}

	// 操作日志
	apiOperateLogs := apis.OperateLogs{}
	rOperateLogs := v1.Group("/admin/operate-logs").Use(authMiddleware.MiddlewareFunc()) /*.Use(middleware.AuthCheckRole())*/
	{
		rOperateLogs.GET("", apiOperateLogs.GetPage)
		rOperateLogs.GET("/:id", apiOperateLogs.Get)
	}

	// 调拨单
	apiStockTransfer := apis.StockTransfer{}
	rStockTransfer := v1.Group("/admin/stock-transfer").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		rStockTransfer.GET("", actions.PermissionAction(), apiStockTransfer.GetPage)
		rStockTransfer.DELETE("", actions.PermissionAction(), apiStockTransfer.Delete)
		rStockTransfer.GET("/:id", actions.PermissionAction(), apiStockTransfer.Get)
		rStockTransfer.POST("", actions.PermissionAction(), apiStockTransfer.Insert)
		rStockTransfer.POST("/save-commit", actions.PermissionAction(), apiStockTransfer.SaveCommit)
		rStockTransfer.POST("/audit", actions.PermissionAction(), apiStockTransfer.Audit)
		rStockTransfer.PUT("/:id", actions.PermissionAction(), apiStockTransfer.Update)
		rStockTransfer.PUT("/update-commit/:id", actions.PermissionAction(), apiStockTransfer.UpdateCommit)
		rStockTransfer.GET("/validate-skus", apiStockTransfer.ValidateSkus)
	}

	// 入库单
	apiStockEntry := apis.StockEntry{}
	rStockEntry := v1.Group("/admin/stock-entry").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		rStockEntry.GET("", actions.PermissionAction(), apiStockEntry.GetPage)
		rStockEntry.GET("/:id", actions.PermissionAction(), apiStockEntry.Get)
		rStockEntry.POST("/add", actions.PermissionAction(), apiStockEntry.Add)
		rStockEntry.POST("/edit", actions.PermissionAction(), apiStockEntry.Edit)
		rStockEntry.PUT("/check-entry/:id", actions.PermissionAction(), apiStockEntry.CheckEntry)
		rStockEntry.POST("/stash", actions.PermissionAction(), apiStockEntry.Stash)
		rStockEntry.POST("/confirm-entry", actions.PermissionAction(), apiStockEntry.ConfirmEntry)
		rStockEntry.GET("/stock-locations", actions.PermissionAction(), apiStockEntry.StockLocationsSelect)
		rStockEntry.GET("/print-html/:id", actions.PermissionAction(), apiStockEntry.GetPrintHtml)
		rStockEntry.GET("/print-sku/:goodsId", actions.PermissionAction(), apiStockEntry.GetPrintSkuInfo)
		rStockEntry.POST("/print-skus", actions.PermissionAction(), apiStockEntry.GetPrintSkuInfos)
		rStockEntry.POST("/part-entry", actions.PermissionAction(), apiStockEntry.PartEntry)
		rStockEntry.GET("/export", actions.PermissionAction(), apiStockEntry.Export)
		rStockEntry.GET("/validate-skus", apiStockEntry.ValidateSkus)
	}
	v1.Group("/admin/stock-entry").GET("/sku-info/:goodsId", actions.PermissionAction(), apiStockEntry.GetSkuInfoByGoodsId)

	// 出库单
	apiStockOutbound := apis.StockOutbound{}
	rStockOutbound := v1.Group("/admin/stock-outbound").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		rStockOutbound.GET("", actions.PermissionAction(), apiStockOutbound.GetPage)
		rStockOutbound.GET("/export", actions.PermissionAction(), apiStockOutbound.Export)
		rStockOutbound.GET("/:id", actions.PermissionAction(), apiStockOutbound.Get)
		rStockOutbound.POST("/confirm-outbound", actions.PermissionAction(), apiStockOutbound.ConfirmOutbound)
		rStockOutbound.GET("/print-outbound/:id", actions.PermissionAction(), apiStockOutbound.PrintOutbound)
		rStockOutbound.GET("/print-picking/:id", actions.PermissionAction(), apiStockOutbound.PrintPicking)
		rStockOutbound.POST("/part-outbound", actions.PermissionAction(), apiStockOutbound.PartOutbound)
	}

	// 库存调整单
	apiStockControl := apis.StockControl{}
	rStockControl := v1.Group("admin/stock-control").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		rStockControl.GET("", actions.PermissionAction(), apiStockControl.GetPage)
		rStockControl.GET("/:id", actions.PermissionAction(), apiStockControl.Get)
		rStockControl.POST("", actions.PermissionAction(), apiStockControl.Insert)
		rStockControl.PUT("/:id", actions.PermissionAction(), apiStockControl.Update)
		rStockControl.DELETE("", actions.PermissionAction(), apiStockControl.Delete)
		rStockControl.POST("/save-commit", actions.PermissionAction(), apiStockControl.SaveCommit)
		rStockControl.PUT("/update-commit/:id", actions.PermissionAction(), apiStockControl.UpdateCommit)
		rStockControl.POST("/audit", actions.PermissionAction(), apiStockControl.Audit)
		rStockControl.GET("/validate-skus", apiStockControl.ValidateSkus)
		rStockControl.POST("/import", actions.PermissionAction(), apiStockControl.StockImport)
	}

	// 库存
	apiStockInfo := apis.StockInfo{}
	rStockInfo := v1.Group("/admin/stock-info").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		rStockInfo.GET("", actions.PermissionAction(), apiStockInfo.GetPage)
		rStockInfo.GET("/export", actions.PermissionAction(), apiStockInfo.Export)
	}

	// 库存日志
	apiStockLog := apis.StockLog{}
	rStockLog := v1.Group("/admin/stock-log").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		rStockLog.GET("", actions.PermissionAction(), apiStockLog.GetPage)
		rStockLog.GET("/export", actions.PermissionAction(), apiStockLog.Export)
	}

	apiStockLocationGoodsLog := apis.StockLocationGoodsLog{}
	rStockLocationGoodsLog := v1.Group("/admin/stock-location-goods-log").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		rStockLocationGoodsLog.GET("", actions.PermissionAction(), apiStockLocationGoodsLog.GetPage)
		rStockLocationGoodsLog.GET("/export", actions.PermissionAction(), apiStockLocationGoodsLog.Export)
	}

	// 库位管理
	apiSl := apis.StockLocation{}
	r := v1.Group("/stock-location").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", actions.PermissionAction(), apiSl.GetPage)
		r.GET("/:id", actions.PermissionAction(), apiSl.Get)
		r.POST("", apiSl.Insert)
		r.PUT("/:id", actions.PermissionAction(), apiSl.Update)
		r.POST("/import", actions.PermissionAction(), apiSl.Import)
	}

	// 库位商品管理
	apiSlg := apis.StockLocationGoods{}
	r = v1.Group("/stock-location-goods").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", actions.PermissionAction(), apiSlg.GetPage)
		r.GET("/:id", actions.PermissionAction(), apiSlg.Get)
		r.GET("/location-list", actions.PermissionAction(), apiSlg.GetSameLogicStockLocationList)
		r.POST("/transfer-stock", actions.PermissionAction(), apiSlg.TransferStock)
		r.GET("/export", actions.PermissionAction(), apiSlg.Export)
	}

	//公共打印 prints
	apiPrits := apis.PrintsJson{}
	rApiPrints := v1.Group("/admin/common-prints").Use(authMiddleware.MiddlewareFunc()) /*.Use(middleware.AuthCheckRole())*/
	{
		rApiPrints.GET("/:print-type/:id", actions.PermissionAction(), apiPrits.CommonPrints)
	}

	// 质检配置
	api := apis.QualityCheckConfig{}
	r = v1.Group("/admin/quality-check-config").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", actions.PermissionAction(), api.GetPage)
		r.GET("/:id", actions.PermissionAction(), api.Get)
		r.POST("", api.Insert)
		r.PUT("/:id", actions.PermissionAction(), api.Update)
		r.DELETE("", api.Delete)
		r.GET("init", actions.PermissionAction(), api.Init)
	}

	// 质检任务
	qualityCheck := apis.QualityCheck{}
	r = v1.Group("/admin/quality-check").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", actions.PermissionAction(), qualityCheck.GetPage)
		r.GET("/:id", actions.PermissionAction(), qualityCheck.Get)
		r.POST("/uploadquality/:id", actions.PermissionAction(), qualityCheck.UploadQuality)
		r.GET("/export", actions.PermissionAction(), qualityCheck.Export)
	}
}
