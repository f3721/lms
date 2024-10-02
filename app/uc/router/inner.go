package router

import (
	"github.com/gin-gonic/gin"
	apis "go-admin/app/uc/apis/admin"
	apisMall "go-admin/app/uc/apis/mall"
)

func init() {
	routerInner = append(routerInner, registerInnerRouter)
}

// registerInnerRouter
func registerInnerRouter(v1 *gin.RouterGroup) {
	api := apis.CompanyInfo{}

	rCompany := v1.Group("/admin/company-info")
	{
		rCompany.GET("/:id", api.GetInner)
		rCompany.GET("get-by-name", api.GetByName)
		rCompany.GET("is-available", api.IsAvailable)
		rCompany.GET("select-list", api.GetInnerSelectList)
	}

	apiUser := apis.UserInfo{}
	rUser := v1.Group("/admin/user-info")
	{
		rUser.GET("/:id", apiUser.Get)
	}

	// mall-Inner
	apisMallUserCollect := apisMall.UserCollect{}
	r := v1.Group("/user-collect")
	{
		r.GET("get-goods-is-collected", apisMallUserCollect.GetGoodsIsCollected)
	}

}
