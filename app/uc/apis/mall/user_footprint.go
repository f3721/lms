package mall

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/uc/models"
	service "go-admin/app/uc/service/mall"
	"go-admin/app/uc/service/mall/dto"
	"go-admin/common/actions"
)

type UserFootprint struct {
	api.Api
}

// GetPage 获取用户足迹表列表
// @Summary 获取用户足迹表列表
// @Description 获取用户足迹表列表
// @Tags 商城-用户足迹
// @Param goodsId query int false "goodsId"
// @Param skuCode query string false "skuCode"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.UserFootprint}} "{"code": 200, "data": [...]}"
// @Router /api/v1/uc/mall/user-footprint [get]
// @Security Bearer
func (e UserFootprint) GetPage(c *gin.Context) {
	req := dto.UserFootprintGetPageReq{}
	s := service.UserFootprint{}
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
	list := make([]*dto.UserFootprintGetListPageRes, 0)
	var count int64

	req.UserId = user.GetUserId(c)
	err = s.GetListPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取用户足迹表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取用户足迹表
// @Summary 获取用户足迹表
// @Description 获取用户足迹表
// @Tags 商城-用户足迹
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.UserFootprint} "{"code": 200, "data": [...]}"
// @Router /api/v1/uc/mall/user-footprint/{id} [get]
// @Security Bearer
func (e UserFootprint) Get(c *gin.Context) {
	req := dto.UserFootprintGetReq{}
	s := service.UserFootprint{}
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
	var object models.UserFootprint

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取用户足迹表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建用户足迹表
// @Summary 创建用户足迹表
// @Description 创建用户足迹表
// @Tags 商城-用户足迹
// @Accept application/json
// @Product application/json
// @Param data body dto.UserFootprintInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/uc/mall/user-footprint [post]
// @Security Bearer
func (e UserFootprint) Insert(c *gin.Context) {
	req := dto.UserFootprintInsertReq{}
	s := service.UserFootprint{}
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

	req.UserId = user.GetUserId(c)
	err = s.Insert(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建用户足迹表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(nil, "创建成功")
}

// Update 修改用户足迹表
// @Summary 修改用户足迹表
// @Description 修改用户足迹表
// @Tags 商城-用户足迹
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.UserFootprintUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/uc/mall/user-footprint/{id} [put]
// @Security Bearer
func (e UserFootprint) Update(c *gin.Context) {
	req := dto.UserFootprintUpdateReq{}
	s := service.UserFootprint{}
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
	p := actions.GetPermissionFromContext(c)

	err = s.Update(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改用户足迹表失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除用户足迹表
// @Summary 删除用户足迹表
// @Description 删除用户足迹表
// @Tags 商城-用户足迹
// @Param data body dto.UserFootprintDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/uc/mall/user-footprint [delete]
// @Security Bearer
func (e UserFootprint) Delete(c *gin.Context) {
	s := service.UserFootprint{}
	req := dto.UserFootprintDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除用户足迹表失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}
