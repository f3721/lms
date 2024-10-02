package goStore

import (
	"fmt"
	"go-admin/common/client"

	"github.com/gin-gonic/gin"
	"github.com/monaco-io/request/response"
	"gorm.io/gorm"
)

var goStoreApiUrl string
var goStoreApiAuthUserName string
var goStoreApiAuthPassword string

type ApiClient struct {
	*client.BaseClient
}

func init() {
	goStoreApiUrl = client.ReadConfig("apiUrl.goStore")
	goStoreApiAuthUserName = client.ReadConfigByPrefix("goStore", "authName")
	goStoreApiAuthPassword = client.ReadConfigByPrefix("goStore", "authPassword")
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

func (c *ApiClient) GetShortLink(url string) *response.Sugar {
	header := client.HS{
		"X-Request-UserName":    goStoreApiAuthUserName,
		"X-Request-PasswordKey": goStoreApiAuthPassword,
	}
	fmt.Println("GetShortLink[2]:", url)
	return client.PostCustomHeader(c.Context, goStoreApiUrl+"/api/v1/short_link/add", map[string]string{"url": url}, header)
}
