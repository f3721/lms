package oc

import (
	"github.com/gin-gonic/gin"
	"github.com/monaco-io/request/response"
	"go-admin/common/client"
	"gorm.io/gorm"
)

var ocApiUrl string

type ApiClient struct {
	*client.BaseClient
}

func init() {
	ocApiUrl = client.ReadConfig("apiUrl.oc")
}

func Api(c *gin.Context) *ApiClient {
	return &ApiClient{
		&client.BaseClient{
			Context: c,
		},
	}
}

func ApiByDbContext(db *gorm.DB) *ApiClient {
	ctx := db.Statement.Context
	switch ctx.(type) {
	case *gin.Context:
		return Api(ctx.(*gin.Context))
	default:
		return nil
	}
}

func (c *ApiClient) GetByOrderId(orderId string) *response.Sugar {
	return client.Get(c.Context, ocApiUrl+"/inner/oc/order", client.HS{
		"orderIds": orderId,
	})
}

func (c *ApiClient) GetByCsNo(csNo string) *response.Sugar {
	return client.Get(c.Context, ocApiUrl+"/inner/oc/order/by-cs-no/"+csNo, client.HS{})
}

func (c *ApiClient) CheckOrderIsAfterPending(orderId string) *response.Sugar {
	return client.Get(c.Context, ocApiUrl+"/inner/oc/order/is-order-in-after-pending-review/"+orderId, client.HS{})
}

func (c *ApiClient) CheckExistUnCompletedOrder(companyId string) *response.Sugar {
	return client.Get(c.Context, ocApiUrl+"/inner/oc/order/check-uncompleted/"+companyId, client.HS{})
}

func (c *ApiClient) GetOrderListByUserName(userName string) *response.Sugar {
	return client.Get(c.Context, ocApiUrl+"/inner/oc/order/list", client.HS{
		"userName": userName,
	})
}
