package admin

import (
	"github.com/gin-gonic/gin"
	"go-admin/common/client"
	"gorm.io/gorm"
	"strconv"

	"github.com/monaco-io/request/response"
)

var adminApiUrl string

func init() {
	adminApiUrl = client.ReadConfig("apiUrl.admin")
}

type ApiClient struct {
	*client.BaseClient
}

func Api(c *gin.Context) *ApiClient {
	return &ApiClient{
		&client.BaseClient{Context: c},
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

func (c *ApiClient) Test(userId int64) *response.Sugar {
	return client.Get(c.Context, adminApiUrl+"/admin/get/part", map[string]string{
		"userId": strconv.FormatInt(userId, 10),
	})
}

type ApiClientPermissionUpdateRequest struct {
	UserID                    int    `json:"userId"`
	AuthorityCompanyID        string `json:"authorityCompanyId"`
	AuthorityWarehouseID      string `json:"authorityWarehouseId"`
	AuthorityWarehouseAllocID string `json:"authorityWarehouseAllocateId"`
	AuthorityVendorID         string `json:"authorityVendorId"`
	UpdateBy                  int    `json:"updateBy"`
	UpdateByName              string `json:"updateByName"`
}

// UpdatePermission 更新用户权限 更新那个权限传哪个
//
// 参数:
//   - userId: 用户ID
//   - AuthorityCompanyId: 用户公司权限 非必填
//   - AuthorityWarehouseId: 用户仓库权限 非必填
//   - AuthorityWarehouseAllocateId: 用户仓库调拨权限 非必填
//   - AuthorityVendorId: 用户货主权限 非必填
//
// 返回值:
//   - resp: 请求响应对象
func (c *ApiClient) UpdatePermission(req ApiClientPermissionUpdateRequest) *response.Sugar {
	return client.Put(c.Context, adminApiUrl+"/inner/admin/sys-user/permission", req)
}

// GetDictListByDictType 获取字典
func (c *ApiClient) GetDictListByDictType(dictType string) *response.Sugar {
	return client.Get(c.Context, adminApiUrl+"/inner/admin/sys-dict/data", map[string]string{
		"dictType": dictType,
	})
}
