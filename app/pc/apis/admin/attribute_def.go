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

type AttributeDef struct {
	api.Api
}

// GetPage 获取属性主表列表
// @Summary 获取属性主表列表
// @Description 获取属性主表列表
// @Tags 属性主表
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.AttributeDef}} "{"code": 200, "data": [...]}"
// @Router /api/v1/pc/attribute-def [get]
// @Security Bearer
func (e AttributeDef) GetPage(c *gin.Context) {
    req := dto.AttributeDefGetPageReq{}
    s := service.AttributeDef{}
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
	list := make([]models.AttributeDef, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取属性主表失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取属性主表
// @Summary 获取属性主表
// @Description 获取属性主表
// @Tags 属性主表
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.AttributeDef} "{"code": 200, "data": [...]}"
// @Router /api/v1/pc/attribute-def/{id} [get]
// @Security Bearer
func (e AttributeDef) Get(c *gin.Context) {
	req := dto.AttributeDefGetReq{}
	s := service.AttributeDef{}
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
	var object models.AttributeDef

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取属性主表失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK( object, "查询成功")
}

// Insert 创建属性主表
// @Summary 创建属性主表
// @Description 创建属性主表
// @Tags 属性主表
// @Accept application/json
// @Product application/json
// @Param data body dto.AttributeDefInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/pc/attribute-def [post]
// @Security Bearer
func (e AttributeDef) Insert(c *gin.Context) {
    req := dto.AttributeDefInsertReq{}
    s := service.AttributeDef{}
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
		e.Error(500, err, fmt.Sprintf("创建属性主表失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改属性主表
// @Summary 修改属性主表
// @Description 修改属性主表
// @Tags 属性主表
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.AttributeDefUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/pc/attribute-def/{id} [put]
// @Security Bearer
func (e AttributeDef) Update(c *gin.Context) {
    req := dto.AttributeDefUpdateReq{}
    s := service.AttributeDef{}
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
		e.Error(500, err, fmt.Sprintf("修改属性主表失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "修改成功")
}

// Delete 删除属性主表
// @Summary 删除属性主表
// @Description 删除属性主表
// @Tags 属性主表
// @Param data body dto.AttributeDefDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/pc/attribute-def [delete]
// @Security Bearer
func (e AttributeDef) Delete(c *gin.Context) {
    s := service.AttributeDef{}
    req := dto.AttributeDefDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除属性主表失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "删除成功")
}
