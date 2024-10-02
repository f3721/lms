package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"go-admin/common/global"
)

func WithContextDb(c *gin.Context) {
	tenantKey := "tenant-id"
	tenantId := c.GetHeader(tenantKey)
	if tenantId == "" {
		tenantId = c.Query(tenantKey)
		if tenantId == "" {
			tenantId = c.PostForm(tenantKey)
		}
	}
	if tenantId != "" {
		//id, err := pkg.StringToInt(tenantId)
		//if err != nil {
		//	response.Error(c, 500, nil, "tenant-id非数值型")
		//	return
		//}
		tenant, err := global.GetTenant(tenantId)
		if err != nil {
			response.Error(c, 500, nil, "该tenant-id，全局未找到")
			return
		}
		tenantDBPrefix := tenant.TenantDBPrefix()
		db := sdk.Runtime.GetDbByKey(tenantDBPrefix)
		if db == nil {
			response.Error(c, 500, nil, "该租户数据库未配置")
			return
		}

		c.Set("db", db.WithContext(c))
	} else {
		db := sdk.Runtime.GetDbByKey(c.Request.Host)
		if db != nil {
			c.Set("db", db.WithContext(c))
		} else {
			db = sdk.Runtime.GetDbByKey("base")
			c.Set("db", db.WithContext(c))
		}
	}
	c.Next()
}
