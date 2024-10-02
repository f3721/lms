package admin

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/pc/models"
	service "go-admin/app/pc/service/admin"
	"go-admin/app/pc/service/admin/dto"
	"go-admin/common/actions"
)

type CategoryAttribute struct {
	api.Api
}

// GetList 获取产线属性表列表
// @Summary 获取产线属性表列表
// @Description 获取产线属性表列表
// @Tags 产线属性表
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.CategoryAttribute}} "{"code": 200, "data": [...]}"
// @Router /api/v1/pc/category-attribute [get]
// @Security Bearer
func (e CategoryAttribute) GetList(c *gin.Context) {
	req := dto.CategoryAttributeGetPageReq{}
	s := service.CategoryAttribute{}
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
	list := make([]models.CategoryAttribute, 0)
	err = s.GetList(&req, p, &list)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取产线属性表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(list, "查询成功")
}

// Get 获取产线属性表
// @Summary 获取产线属性表
// @Description 获取产线属性表
// @Tags 产线属性表
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.CategoryAttribute} "{"code": 200, "data": [...]}"
// @Router /api/v1/pc/category-attribute/{id} [get]
// @Security Bearer
func (e CategoryAttribute) Get(c *gin.Context) {
	req := dto.CategoryAttributeGetReq{}
	s := service.CategoryAttribute{}
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
	var object models.CategoryAttribute

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取产线属性表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建产线属性表
// @Summary 创建产线属性表
// @Description 创建产线属性表
// @Tags 产线属性表
// @Accept application/json
// @Product application/json
// @Param data body dto.CategoryAttributeInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/pc/category-attribute [post]
// @Security Bearer
func (e CategoryAttribute) Insert(c *gin.Context) {
	req := dto.CategoryAttributeInsertReq{}
	s := service.CategoryAttribute{}
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
		e.Error(500, err, fmt.Sprintf("创建产线属性表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改产线属性表
// @Summary 修改产线属性表
// @Description 修改产线属性表
// @Tags 产线属性表
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.CategoryAttributeUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/pc/category-attribute/{id} [put]
// @Security Bearer
func (e CategoryAttribute) Update(c *gin.Context) {
	req := dto.CategoryAttributeUpdateReq{}
	s := service.CategoryAttribute{}
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
		e.Error(500, err, fmt.Sprintf("修改产线属性表失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除产线属性表
// @Summary 删除产线属性表
// @Description 删除产线属性表
// @Tags 产线属性表
// @Param data body dto.CategoryAttributeDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/pc/category-attribute [delete]
// @Security Bearer
func (e CategoryAttribute) Delete(c *gin.Context) {
	s := service.CategoryAttribute{}
	req := dto.CategoryAttributeDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除产线属性表失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}
