package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	apis "go-admin/app/pc/apis/admin"
	"go-admin/common/actions"
	"go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, sysNoCheckRoleRouter, registerSysAdminApiRouter)
}

func sysNoCheckRoleRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	productApi := apis.Product{}
	r := v1.Group("/product")
	{
		r.GET("/supplier_product_sync", productApi.SupplierProductSync)
	}
}

// registerSysApiRouter
func registerSysAdminApiRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	// 品牌
	brandApi := apis.Brand{}
	r := v1.Group("/brand").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", actions.PermissionAction(), brandApi.GetPage)
		r.GET("/list", actions.PermissionAction(), brandApi.GetBrandListPage)
		r.GET("/:id", actions.PermissionAction(), brandApi.Get)
		r.POST("", brandApi.Insert)
		r.PUT("/:id", actions.PermissionAction(), brandApi.Update)
		r.DELETE("", brandApi.Delete)
	}

	// 品牌日志
	brandLogApi := apis.BrandLog{}
	r2 := v1.Group("/brand-log").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r2.GET("", actions.PermissionAction(), brandLogApi.GetPage)
		r2.GET("/:id", actions.PermissionAction(), brandLogApi.Get)
	}

	// 分类属性
	attributeDefApi := apis.AttributeDef{}
	r3 := v1.Group("/attribute-def").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r3.GET("", actions.PermissionAction(), attributeDefApi.GetPage)
		r3.GET("/:id", actions.PermissionAction(), attributeDefApi.Get)
		r3.POST("", attributeDefApi.Insert)
		r3.PUT("/:id", actions.PermissionAction(), attributeDefApi.Update)
		r3.DELETE("", attributeDefApi.Delete)
	}

	// 分类
	categoryApi := apis.Category{}
	r4 := v1.Group("/category").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r4.GET("", actions.PermissionAction(), categoryApi.GetPage)
		r4.GET("/:id", actions.PermissionAction(), categoryApi.Get)
		r4.GET("/category-path/:id", actions.PermissionAction(), categoryApi.GetCategoryPath)
		r4.GET("/list/:id", actions.PermissionAction(), categoryApi.GetList)
		r4.POST("", categoryApi.Insert)
		r4.PUT("/:id", actions.PermissionAction(), categoryApi.Update)
		r4.DELETE("", categoryApi.Delete)
		// 排序
		r4.POST("/sort", actions.PermissionAction(), categoryApi.Sort)
	}

	// 分类属性
	categoryAttributeApi := apis.CategoryAttribute{}
	r5 := v1.Group("/category-attribute").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r5.GET("", actions.PermissionAction(), categoryAttributeApi.GetList)
		r5.GET("/:id", actions.PermissionAction(), categoryAttributeApi.Get)
		r5.POST("", categoryAttributeApi.Insert)
		r5.PUT("/:id", actions.PermissionAction(), categoryAttributeApi.Update)
		r5.DELETE("", categoryAttributeApi.Delete)
	}

	// 分类日志
	categgoryLogApi := apis.CategoryLog{}
	r6 := v1.Group("/category-log").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r6.GET("", actions.PermissionAction(), categgoryLogApi.GetPage)
		r6.GET("/:id", actions.PermissionAction(), categgoryLogApi.Get)
	}

	// 商品档案
	productApi := apis.Product{}
	r7 := v1.Group("/product").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		// 列表
		r7.GET("", actions.PermissionAction(), productApi.GetPage)
		// 详情
		r7.GET("/:id", actions.PermissionAction(), productApi.Get)
		// 新增
		r7.POST("", productApi.Insert)
		// 更新
		r7.PUT("/:id", actions.PermissionAction(), productApi.Update)
		// 删除
		r7.DELETE("", productApi.Delete)
		// 批量图文审核
		r7.PUT("/batch-pro-approval", actions.PermissionAction(), productApi.BatchProApproval)
		// 单个审核
		r7.PUT("/pro-approval/:id", actions.PermissionAction(), productApi.ProApproval)
		// 产品新增导入
		r7.POST("/product-import-add", actions.PermissionAction(), productApi.ProductImportAdd)
		// 产品维护导入
		r7.POST("/product-import-update", actions.PermissionAction(), productApi.ProductImportUpdate)
		// 产品维护导出
		r7.GET("/product-export", actions.PermissionAction(), productApi.ProductExport)
		// 产品属性维护导出
		r7.GET("/product-attribute-export", actions.PermissionAction(), productApi.AttributeExport)
		// 产品属性维护导入
		r7.POST("/product-attribute-import", actions.PermissionAction(), productApi.ProductAttributeImport)
		// 产品属性查找
		r7.POST("/product-attribute", productApi.GetProductAttribute)
		// 批量上传产品图片
		r7.POST("/batch-upload-product-image", productApi.BatchUploadProductImage)
		// skucode获取基本信息
		r7.GET("/skucode/:skuCode", productApi.GetProductBySkuCode)
		// 排序
		r7.POST("/sort", productApi.Sort)
	}

	// 商品档案日志
	productLogApi := apis.ProductLog{}
	r8 := v1.Group("/product-log").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r8.GET("", actions.PermissionAction(), productLogApi.GetPage)
		r8.GET("/:id", actions.PermissionAction(), productLogApi.Get)
		r8.POST("", productLogApi.Insert)
		r8.PUT("/:id", actions.PermissionAction(), productLogApi.Update)
		r8.DELETE("", productLogApi.Delete)
	}

	// 商品管理
	goodsApi := apis.Goods{}
	r9 := v1.Group("/goods").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r9.GET("", actions.PermissionAction(), goodsApi.GetPage)
		r9.GET("/:id", actions.PermissionAction(), goodsApi.Get)
		r9.POST("", goodsApi.Insert)
		r9.PUT("/:id", actions.PermissionAction(), goodsApi.Update)
		r9.DELETE("", goodsApi.Delete)
		// 商品审核
		r9.POST("/approve", actions.PermissionAction(), goodsApi.Approve)
		// 商品批量上下架
		r9.POST("/online-offline", actions.PermissionAction(), goodsApi.OnlineOffline)
		// 商品导入
		r9.POST("/import", actions.PermissionAction(), goodsApi.GoodsImport)
		// 商品导出
		r9.GET("/export", actions.PermissionAction(), goodsApi.Export)
	}

	// 商品管理日志
	goodsLogApi := apis.GoodsLog{}
	r10 := v1.Group("/goods-log").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r10.GET("", actions.PermissionAction(), goodsLogApi.GetPage)
		r10.GET("/:id", actions.PermissionAction(), goodsLogApi.Get)
	}

	// 物理单位
	uommasterApi := apis.Uommaster{}
	r11 := v1.Group("/uommaster").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r11.GET("", actions.PermissionAction(), uommasterApi.GetPage)
		r11.GET("/:id", actions.PermissionAction(), uommasterApi.Get)
		r11.POST("", uommasterApi.Insert)
		r11.PUT("/:id", actions.PermissionAction(), uommasterApi.Update)
		r11.DELETE("", uommasterApi.Delete)
	}

	// 日志类型
	logTypeApi := apis.LogTypes{}
	r12 := v1.Group("/log-types").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r12.GET("", actions.PermissionAction(), logTypeApi.GetPage)
		r12.GET("/:id", actions.PermissionAction(), logTypeApi.Get)
		r12.POST("", logTypeApi.Insert)
		r12.PUT("/:id", actions.PermissionAction(), logTypeApi.Update)
		r12.DELETE("", logTypeApi.Delete)
	}
}
