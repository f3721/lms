package router

import (
	"github.com/gin-gonic/gin"
	apis "go-admin/app/oc/apis/admin"
)

func init() {
	routerInner = append(routerInner, registerInnerRouter)
}

// registerInnerRouter
func registerInnerRouter(v1 *gin.RouterGroup) {
	api := apis.OrderInfo{}
	r := v1.Group("/order")
	{
		r.GET("", api.GetByOrderId)
		r.GET("/check-uncompleted/:id", api.CheckExistUnCompletedOrder)
		r.GET("/list", api.GetOrderListByUserName)
	}

	apiCsApply := apis.CsApply{}
	r = v1.Group("/order")
	{
		r.GET("by-cs-no/:csNo", apiCsApply.GetOrderInfoByCsNo)
		r.GET("is-order-in-after-pending-review/:orderId", apiCsApply.IsOrderInAfterPendingReview)
	}

}
