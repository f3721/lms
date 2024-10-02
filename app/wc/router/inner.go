package router

import (
	"github.com/gin-gonic/gin"
	apis "go-admin/app/wc/apis/admin"
)

func init() {
	routerInner = append(routerInner, registerInnerRouter)
}

// registerInnerRouter
func registerInnerRouter(v1 *gin.RouterGroup) {
	apiV := apis.Vendors{}
	rV := v1.Group("/admin/vendors")
	{
		rV.GET("/list", apiV.InnerGetList)
	}

	apiSu := apis.Supplier{}
	rSu := v1.Group("/admin/supplier")
	{
		rSu.GET("/list", apiSu.InnerGetList)
	}

	apiWh := apis.Warehouse{}
	rWh := v1.Group("/admin/warehouse")
	{
		rWh.POST("/list-by-name-company-id", apiWh.InnerGetListByNameAndCompanyId)
		rWh.GET("/list", apiWh.InnerGetList)
	}

	apiStockInfo := apis.StockInfo{}
	rStockInfo := v1.Group("/admin/stock-info")
	{
		rStockInfo.POST("get", apiStockInfo.InnerGetByGoodsIdAndLwhCode)
		rStockInfo.POST("", apiStockInfo.InnerGetByGoodsIdAndWarehouseCode)
	}

	apiRegion := apis.Region{}
	rRegion := v1.Group("/admin/region")
	{
		rRegion.GET("get-info", apiRegion.InnerGetByIds)
	}

}
