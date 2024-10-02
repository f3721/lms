package admin

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/wc/models"
	service "go-admin/app/wc/service/admin"
	"go-admin/app/wc/service/admin/dto"
	"go-admin/common/actions"
)

type QualityCheckConfig struct {
	api.Api
}

// GetPage 获取质检配置列表
// @Summary 获取质检配置列表
// @Description 获取质检配置列表
// @Tags 质检配置
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.QualityCheckConfig}} "{"code": 200, "data": [...]}"
// @Router /api/v1/quality-check-config [get]
// @Security Bearer
func (e QualityCheckConfig) GetPage(c *gin.Context) {
	req := dto.QualityCheckConfigGetPageReq{}
	s := service.QualityCheckConfig{}
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
	list := make([]*dto.QualityCheckConfigGetPageResp, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取质检配置失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取质检配置
// @Summary 获取质检配置
// @Description 获取质检配置
// @Tags 质检配置
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.QualityCheckConfig} "{"code": 200, "data": [...]}"
// @Router /api/v1/quality-check-config/{id} [get]
// @Security Bearer
func (e QualityCheckConfig) Get(c *gin.Context) {
	req := dto.QualityCheckConfigGetReq{}
	s := service.QualityCheckConfig{}
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
	var object models.QualityCheckConfig

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取质检配置失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建质检配置
// @Summary 创建质检配置
// @Description 创建质检配置
// @Tags 质检配置
// @Accept application/json
// @Product application/json
// @Param data body dto.QualityCheckConfigInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/quality-check-config [post]
// @Security Bearer
func (e QualityCheckConfig) Insert(c *gin.Context) {
	req := dto.QualityCheckConfigInsertReq{}
	s := service.QualityCheckConfig{}
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

	// 设置创建人
	req.SetCreateBy(user.GetUserId(c))
	req.SetCreateByName(user.GetUserName(c))

	err = s.Insert(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建质检配置失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改质检配置
// @Summary 修改质检配置
// @Description 修改质检配置
// @Tags 质检配置
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.QualityCheckConfigInsertReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/quality-check-config/{id} [put]
// @Security Bearer
func (e QualityCheckConfig) Update(c *gin.Context) {
	req := dto.QualityCheckConfigInsertReq{}
	s := service.QualityCheckConfig{}
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

	req.SetUpdateBy(user.GetUserId(c))
	req.SetUpdateByName(user.GetUserName(c))
	p := actions.GetPermissionFromContext(c)

	err = s.Update(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改质检配置失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除质检配置
// @Summary 删除质检配置
// @Description 删除质检配置
// @Tags 质检配置
// @Param data body dto.QualityCheckConfigDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/quality-check-config [delete]
// @Security Bearer
func (e QualityCheckConfig) Delete(c *gin.Context) {
	s := service.QualityCheckConfig{}
	req := dto.QualityCheckConfigDeleteReq{}
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

	// req.SetUpdateBy(user.GetUserId(c))
	p := actions.GetPermissionFromContext(c)

	err = s.Remove(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("删除质检配置失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}

// GetPage 获取初始化信息
// @Summary 获取初始化信息
// @Description 获取初始化信息
// @Tags 质检配置
// @Success 200 {object} response.Response{data=response.Page{list=[]models.QualityCheckConfig}} "{"code": 200, "data": [...]}"
// @Router /api/v1/quality-check-config [get]
// @Security Bearer
func (e QualityCheckConfig) Init(c *gin.Context) {
	s := service.QualityCheckConfig{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	res, err := s.Init()
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取初始化信息失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(res, "获取初始化信息成功")
}
