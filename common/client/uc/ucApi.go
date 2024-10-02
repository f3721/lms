package uc

import (
	"github.com/gin-gonic/gin"
	"go-admin/app/uc/service/mall/dto"
	"go-admin/common/client"
	"gorm.io/gorm"
	"strconv"

	"github.com/monaco-io/request/response"
)

var ucApiUrl string

type ApiClient struct {
	*client.BaseClient
}

func init() {
	ucApiUrl = client.ReadConfig("apiUrl.uc")
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

func (c *ApiClient) CompanyIsAvailable(id int) *response.Sugar {
	return client.Get(c.Context, ucApiUrl+"/api/v1/uc/admin/company-info/is-available/"+strconv.Itoa(id), nil)
}

func (c *ApiClient) GetCompanyByIds(ids string) *response.Sugar {
	url := "/inner/uc/admin/company-info/select-list?pageSize=-1"
	if ids != "" {
		url = url + "&queryCompanyIds=" + ids
	}
	return client.Get(c.Context, ucApiUrl+url, nil)
}

func (c *ApiClient) GetCompanyByName(companyNames string) *response.Sugar {
	url := "/inner/uc/admin/company-info/select-list?pageSize=-1"
	if companyNames != "" {
		url = url + "&CompanyStatus=1&queryCompanyNames=" + companyNames
	}
	return client.Get(c.Context, ucApiUrl+url, nil)
}

func (c *ApiClient) GetCompanyInfoById(id int) *response.Sugar {
	url := "/inner/uc/admin/company-info/" + strconv.Itoa(id)

	return client.Get(c.Context, ucApiUrl+url, nil)
}

func (c *ApiClient) GetUserInfoById(id int) *response.Sugar {
	url := "/inner/uc/admin/user-info/" + strconv.Itoa(id)

	return client.Get(c.Context, ucApiUrl+url, nil)
}

func (c *ApiClient) GetUserCollect(d *dto.UserCollectGetGoodsIsCollected) *response.Sugar {
	url := "/inner/uc/user-collect/get-goods-is-collected?pageSize=-1"
	if d.UserId != 0 {
		url = url + "&userId=" + strconv.Itoa(d.UserId)
	}
	if d.GoodsIds != "" {
		url = url + "&goodsIds=" + d.GoodsIds
	}
	return client.Get(c.Context, ucApiUrl+url, nil)
}
