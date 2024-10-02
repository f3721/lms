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

type Uommaster struct {
	api.Api
}

// GetPage 获取Uommaster列表
// @Summary 获取Uommaster列表
// @Description 获取Uommaster列表
// @Tags Uommaster
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.Uommaster}} "{"code": 200, "data": [...]}"
// @Router /api/v1/pc/uommaster [get]
// @Security Bearer
func (e Uommaster) GetPage(c *gin.Context) {
    req := dto.UommasterGetPageReq{}
    s := service.Uommaster{}
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
	list := make([]models.Uommaster, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取Uommaster失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取Uommaster
// @Summary 获取Uommaster
// @Description 获取Uommaster
// @Tags Uommaster
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.Uommaster} "{"code": 200, "data": [...]}"
// @Router /api/v1/pc/uommaster/{id} [get]
// @Security Bearer
func (e Uommaster) Get(c *gin.Context) {
	req := dto.UommasterGetReq{}
	s := service.Uommaster{}
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
	var object models.Uommaster

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取Uommaster失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK( object, "查询成功")
}

// Insert 创建Uommaster
// @Summary 创建Uommaster
// @Description 创建Uommaster
// @Tags Uommaster
// @Accept application/json
// @Product application/json
// @Param data body dto.UommasterInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/pc/uommaster [post]
// @Security Bearer
func (e Uommaster) Insert(c *gin.Context) {
    req := dto.UommasterInsertReq{}
    s := service.Uommaster{}
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
		e.Error(500, err, fmt.Sprintf("创建Uommaster失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改Uommaster
// @Summary 修改Uommaster
// @Description 修改Uommaster
// @Tags Uommaster
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.UommasterUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/pc/uommaster/{id} [put]
// @Security Bearer
func (e Uommaster) Update(c *gin.Context) {
    req := dto.UommasterUpdateReq{}
    s := service.Uommaster{}
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
		e.Error(500, err, fmt.Sprintf("修改Uommaster失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "修改成功")
}

// Delete 删除Uommaster
// @Summary 删除Uommaster
// @Description 删除Uommaster
// @Tags Uommaster
// @Param data body dto.UommasterDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/pc/uommaster [delete]
// @Security Bearer
func (e Uommaster) Delete(c *gin.Context) {
    s := service.Uommaster{}
    req := dto.UommasterDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除Uommaster失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "删除成功")
}
