package router

import (
	"github.com/gin-gonic/gin"
	apis "go-admin/app/pc/apis/admin"
)

func init() {
	routerInner = append(routerInner, registerInnerRouter)
}

// registerInnerRouter
func registerInnerRouter(v1 *gin.RouterGroup) {
	/**
	1、根据sku_code获取商品档案信息(其中supplier_sku_code、name_zh、sku_code、mfg_model、sales_uom、品牌名称等必有) (支持批量查询)
	2、根据sku_code、warehouse_code、vendor_id获取商品管理关系信息（其中id、product_no等必有）（注:已审核、状态正常、【是否上架不关心】）（支持批量查询)
	3、根据goods_id获取商品关系和档案信息(其中id、product_no、supplier_sku_code、name_zh、sku_code、mfg_model、sales_uom、品牌名称等必有) （注:这里不关心任何状态）（支持批量查询)
	4、需要一个通过goodsid批量查询商品信息的接口 要包含 ：商品图片 sku_code 和商品名
	**/
	api1 := apis.Product{}
	r1 := v1.Group("/admin/product")
	{
		r1.POST("/get-by-sku", api1.GetProductBySku)
		r1.POST("/get-category", api1.GetProductCategoryBySku)
	}

	api2 := apis.Goods{}
	r2 := v1.Group("/admin/goods")
	{
		r2.POST("/info", api2.GetGoodsInfo)
		r2.POST("/get-by-id", api2.GetGoodsById)
		r2.POST("/get-by-skucode", api2.GetGoodsBySkuCodeReq)
	}
}
