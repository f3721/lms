package admin

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	service "go-admin/app/uc/service/admin"
	"go-admin/app/uc/service/admin/dto"
	"go-admin/common/actions"
	"go-admin/common/models"
)

type OperateLogs struct {
	api.Api
}

// GetPage 获取操作日志列表
// @Summary 获取操作日志列表
// @Description 获取操作日志列表
// @Tags uc操作日志
// @Param dataId query string false "数据id eg:WH0001"
// @Param modelName query string false "模型name"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.OperateLogs}} "{"code": 200, "data": [...]}"
// @Router  /api/v1/uc/admin/operate-logs [get]
// @Security Bearer
func (e OperateLogs) GetPage(c *gin.Context) {
	req := dto.OperateLogsGetPageReq{}
	s := service.OperateLogs{}
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
	list := make([]models.OperateLogs, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取操作日志失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取操作日志
// @Summary 获取操作日志
// @Description 获取操作日志
// @Tags uc操作日志
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.OperateLogDetailResp} "{"code": 200, "data": [...]}"
// @Router  /api/v1/uc/admin/operate-logs/{id} [get]
// @Security Bearer
func (e OperateLogs) Get(c *gin.Context) {
	req := dto.OperateLogsGetReq{}
	s := service.OperateLogs{}
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
	var object models.OperateLogDetailResp

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取操作日志失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}
