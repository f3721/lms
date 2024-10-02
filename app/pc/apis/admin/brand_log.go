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

type BrandLog struct {
	api.Api
}

// GetPage 获取品牌日志表列表
// @Summary 获取品牌日志表列表
// @Description 获取品牌日志表列表
// @Tags 品牌日志表
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Param brandId query int false "品牌ID"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.BrandLog}} "{"code": 200, "data": [...]}"
// @Router /api/v1/brand-log [get]
// @Security Bearer
func (e BrandLog) GetPage(c *gin.Context) {
	req := dto.BrandLogGetPageReq{}
	s := service.BrandLog{}
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
	list := make([]models.BrandLog, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取品牌日志表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取品牌日志表
// @Summary 获取品牌日志表
// @Description 获取品牌日志表
// @Tags 品牌日志表
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.BrandLog} "{"code": 200, "data": [...]}"
// @Router /api/v1/brand-log/{id} [get]
// @Security Bearer
func (e BrandLog) Get(c *gin.Context) {
	req := dto.BrandLogGetReq{}
	s := service.BrandLog{}
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
		e.Error(500, err, fmt.Sprintf("获取品牌日志表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}
