package admin

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/uc/models"
	service "go-admin/app/uc/service/admin"
	"go-admin/app/uc/service/admin/dto"
	"go-admin/common/actions"
)

type UserLoginLog struct {
	api.Api
}

// GetPage 获取用户登录记录表 列表
// @Summary 获取用户登录记录表 列表
// @Description 获取用户登录记录表 列表
// @Tags 用户登录记录表
// @Param userId query int false "用户ID"
// @Param status query string false "状态"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.UserLoginLog}} "{"code": 200, "data": [...]}"
// @Router /api/v1/uc/admin/user/login-log [get]
// @Security Bearer
func (e UserLoginLog) GetPage(c *gin.Context) {
	req := dto.UserLoginLogGetPageReq{}
	s := service.UserLoginLog{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	p := actions.GetPermissionFromContext(c)
	list := make([]models.UserLoginLog, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取用户登录记录表 失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}
