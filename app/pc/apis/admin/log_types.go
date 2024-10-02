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

type LogTypes struct {
	api.Api
}

// GetPage 获取LOG字段类型表列表
// @Summary 获取LOG字段类型表列表
// @Description 获取LOG字段类型表列表
// @Tags LOG字段类型表
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.LogTypes}} "{"code": 200, "data": [...]}"
// @Router /api/v1/pc/log-types [get]
// @Security Bearer
func (e LogTypes) GetPage(c *gin.Context) {
    req := dto.LogTypesGetPageReq{}
    s := service.LogTypes{}
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
	list := make([]models.LogTypes, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取LOG字段类型表失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取LOG字段类型表
// @Summary 获取LOG字段类型表
// @Description 获取LOG字段类型表
// @Tags LOG字段类型表
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.LogTypes} "{"code": 200, "data": [...]}"
// @Router /api/v1/log-types/{id} [get]
// @Security Bearer
func (e LogTypes) Get(c *gin.Context) {
	req := dto.LogTypesGetReq{}
	s := service.LogTypes{}
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
	var object models.LogTypes

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取LOG字段类型表失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK( object, "查询成功")
}

// Insert 创建LOG字段类型表
// @Summary 创建LOG字段类型表
// @Description 创建LOG字段类型表
// @Tags LOG字段类型表
// @Accept application/json
// @Product application/json
// @Param data body dto.LogTypesInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/log-types [post]
// @Security Bearer
func (e LogTypes) Insert(c *gin.Context) {
    req := dto.LogTypesInsertReq{}
    s := service.LogTypes{}
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

	err = s.Insert(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建LOG字段类型表失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改LOG字段类型表
// @Summary 修改LOG字段类型表
// @Description 修改LOG字段类型表
// @Tags LOG字段类型表
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.LogTypesUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/log-types/{id} [put]
// @Security Bearer
func (e LogTypes) Update(c *gin.Context) {
    req := dto.LogTypesUpdateReq{}
    s := service.LogTypes{}
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
		e.Error(500, err, fmt.Sprintf("修改LOG字段类型表失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "修改成功")
}

// Delete 删除LOG字段类型表
// @Summary 删除LOG字段类型表
// @Description 删除LOG字段类型表
// @Tags LOG字段类型表
// @Param data body dto.LogTypesDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/log-types [delete]
// @Security Bearer
func (e LogTypes) Delete(c *gin.Context) {
    s := service.LogTypes{}
    req := dto.LogTypesDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除LOG字段类型表失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "删除成功")
}
