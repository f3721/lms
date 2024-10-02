package apis

import (
	"fmt"
	"go-admin/common/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/admin/models"
	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
	"go-admin/common/actions"
)

type SysUserLog struct {
	api.Api
}

// GetPage 获取系统用户日志列表
// @Summary 获取系统用户日志列表
// @Description 获取系统用户日志列表
// @Tags 系统用户日志
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.SysUserLog}} "{"code": 200, "data": [...]}"
// @Router /api/v1/sys-user-log [get]
// @Security Bearer
func (e SysUserLog) GetPage(c *gin.Context) {
    req := dto.SysUserLogGetPageReq{}
    s := service.SysUserLog{}
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
	list := make([]models.SysUserLog, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取系统用户日志失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取系统用户日志
// @Summary 获取系统用户日志
// @Description 获取系统用户日志
// @Tags 系统用户日志
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=utils.OperateLogDetailResp} "{"code": 200, "data": [...]}"
// @Router /api/v1/sys-user-log/{id} [get]
// @Security Bearer
func (e SysUserLog) Get(c *gin.Context) {
	req := dto.SysUserLogGetReq{}
	s := service.SysUserLog{}
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
	var object utils.OperateLogDetailResp

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取系统用户日志失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK( object, "查询成功")
}
