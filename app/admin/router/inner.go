package router

import (
	"github.com/gin-gonic/gin"

	"go-admin/app/admin/apis"
)

func init() {
	routerInner = append(routerInner, registerInnerRouter)
}

// registerInnerRouter
func registerInnerRouter(v1 *gin.RouterGroup) {
	api := apis.SysUser{}
	r := v1.Group("/sys-user")
	{
		r.PUT("/permission", api.UpdatePermission)
	}

	dataApi := apis.SysDictData{}
	r = v1.Group("/sys-dict")
	{
		r.GET("/data", dataApi.GetAll)
	}
}
