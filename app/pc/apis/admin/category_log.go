package admin

import (
	"fmt"
	"go-admin/common/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/pc/models"
	service "go-admin/app/pc/service/admin"
	"go-admin/app/pc/service/admin/dto"
	"go-admin/common/actions"
)

type CategoryLog struct {
	api.Api
}

// GetPage 获取产线分类日志表列表
// @Summary 获取产线分类日志表列表
// @Description 获取产线分类日志表列表
// @Tags 产线分类日志表
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.CategoryLog}} "{"code": 200, "data": [...]}"
// @Router /api/v1/pc/category-log [get]
// @Security Bearer
func (e CategoryLog) GetPage(c *gin.Context) {
	req := dto.CategoryLogGetPageReq{}
	s := service.CategoryLog{}
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
	list := make([]models.CategoryLog, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取产线分类日志表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取产线分类日志表
// @Summary 获取产线分类日志表
// @Description 获取产线分类日志表
// @Tags 产线分类日志表
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.CategoryLog} "{"code": 200, "data": [...]}"
// @Router /api/v1/pc/category-log/{id} [get]
// @Security Bearer
func (e CategoryLog) Get(c *gin.Context) {
	req := dto.CategoryLogGetReq{}
	s := service.CategoryLog{}
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
		e.Error(500, err, fmt.Sprintf("获取产线分类日志表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}
