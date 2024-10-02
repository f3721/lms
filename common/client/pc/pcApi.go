package pc

import (
	"github.com/gin-gonic/gin"
	"github.com/monaco-io/request/response"
	"go-admin/app/pc/service/admin/dto"
	"go-admin/common/client"
	"gorm.io/gorm"
)

var pcApiUrl string

type ApiClient struct {
	*client.BaseClient
}

func init() {
	pcApiUrl = client.ReadConfig("apiUrl.pc")
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

//// GetUserName 具体业务逻辑 GET 请求
//func (c *ApiClient) GetUserName(userId int) *response.Sugar {
//	return client.Get(c.Context, pcApiUrl+"/dept", client.HS{
//		"userId": strconv.Itoa(userId),
//	})
//}
//
//// GetCompanyName 具体业务逻辑 POST 请求
//func (c *ApiClient) GetCompanyName(companyId int) *response.Sugar {
//	return client.Post(c.Context, pcApiUrl+"/test", client.HS{
//		"companyId": strconv.Itoa(companyId),
//	}, nil)
//}
//
//// GetCompanyInfo 具体业务逻辑 POST 请求
//func (c *ApiClient) GetCompanyInfo(companyName string) *response.Sugar {
//	return client.Post(c.Context, pcApiUrl+"/test", client.HS{
//		"companyName": companyName,
//	}, nil)
//}
//
//func (c *ApiClient) GetWarehouseInfoByCompanyId(companyId int) *response.Sugar {
//	return client.Post(c.Context, pcApiUrl+"/test", client.HS{
//		"companyId": strconv.Itoa(companyId),
//	}, nil)
//}

func (c *ApiClient) GetProductBySku(skus []string) *response.Sugar {
	req := dto.InnerGetProductBySkuReq{}
	req.SkuCode = skus
	return client.PostJson(c.Context, pcApiUrl+"/inner/pc/admin/product/get-by-sku", req)
}

func (c *ApiClient) GetGoodsBySkuAndVendorAndWarehouse(skus []string, whCode string, vendorId, approveStatus int) *response.Sugar {
	req := dto.GoodsInfoReq{}
	GoodsInfo := dto.GoodsInfo{}
	GoodsInfoSlice := []dto.GoodsInfo{}
	GoodsInfo.VendorId = vendorId
	GoodsInfo.WarehouseCode = whCode
	GoodsInfo.Status = 1
	GoodsInfo.ApproveStatus = approveStatus
	for _, item := range skus {
		GoodsInfo.SkuCode = item
		GoodsInfoSlice = append(GoodsInfoSlice, GoodsInfo)
	}
	req.Query = GoodsInfoSlice
	return client.PostJson(c.Context, pcApiUrl+"/inner/pc/admin/goods/info", req)
}

func (c *ApiClient) GetGoodsById(GoodsSlice []int) *response.Sugar {
	req := dto.GetGoodsByIdReq{
		Ids: GoodsSlice,
	}
	return client.PostJson(c.Context, pcApiUrl+"/inner/pc/admin/goods/get-by-id?", req)
}

// GetGoodsBySkuCodeReq sku[]+warehousecode 查询商品信息
func (c *ApiClient) GetGoodsBySkuCodeReq(query dto.GetGoodsBySkuCodeReq) *response.Sugar {
	return client.PostJson(c.Context, pcApiUrl+"/inner/pc/admin/goods/get-by-skucode", query)
}

func (c *ApiClient) GetCategoryBySku(skus []string) *response.Sugar {
	req := dto.InnerGetProductBySkuReq{
		SkuCode: skus,
	}
	return client.PostJson(c.Context, pcApiUrl+"/inner/pc/admin/product/get-category?", req)
}
