package wc

import (
	"github.com/gin-gonic/gin"
	"github.com/monaco-io/request/response"
	dtoWc "go-admin/app/wc/service/admin/dto"
	"go-admin/common/client"
	"gorm.io/gorm"
)

var wcApiUrl string

type ApiClient struct {
	*client.BaseClient
}

func init() {
	wcApiUrl = client.ReadConfig("apiUrl.wc")
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

func (c *ApiClient) GetWarehouseList(req dtoWc.InnerWarehouseGetListReq) *response.Sugar {
	query := client.HS{}
	if req.WarehouseCode != "" {
		query["warehouseCode"] = req.WarehouseCode
	}
	if req.Status != "" {
		query["status"] = req.Status
	}
	return client.Get(c.Context, wcApiUrl+"/inner/wc/admin/warehouse/list", query)
}

func (c *ApiClient) GetVendorList(req dtoWc.InnerVendorsGetListReq) *response.Sugar {
	query := client.HS{}
	if req.Ids != "" {
		query["ids"] = req.Ids
	}
	if req.NameZh != "" {
		query["nameZh"] = req.NameZh
	}
	if req.Status != "" {
		query["status"] = req.Status
	}
	return client.Get(c.Context, wcApiUrl+"/inner/wc/admin/vendors/list", query)
}

func (c *ApiClient) GetStockListByGoodsIdAndWarehouseCode(req dtoWc.InnerStockInfoGetByGoodsIdAndWarehouseCodeReq) *response.Sugar {
	return client.PostJson(c.Context, wcApiUrl+"/inner/wc/admin/stock-info", req)
}

func (c *ApiClient) GetRegionByIds(ids string) *response.Sugar {
	query := client.HS{}
	query["ids"] = ids
	return client.Get(c.Context, wcApiUrl+"/inner/wc/admin/region/get-info", query)
}
