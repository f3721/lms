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

type ManageMenu struct {
	api.Api
}

// GetPage 获取ManageMenu列表
// @Summary 获取ManageMenu列表
// @Description 获取ManageMenu列表
// @Tags ManageMenu
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.ManageMenu}} "{"code": 200, "data": [...]}"
// @Router /api/v1/uc/manage-menu [get]
// @Security Bearer
func (e ManageMenu) GetPage(c *gin.Context) {
	req := dto.ManageMenuGetPageReq{}
	s := service.ManageMenu{}
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
	list := make([]models.ManageMenu, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取ManageMenu失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取ManageMenu
// @Summary 获取ManageMenu
// @Description 获取ManageMenu
// @Tags ManageMenu
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.ManageMenu} "{"code": 200, "data": [...]}"
// @Router /api/v1/uc/manage-menu/{id} [get]
// @Security Bearer
func (e ManageMenu) Get(c *gin.Context) {
	req := dto.ManageMenuGetReq{}
	s := service.ManageMenu{}
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
	var object models.ManageMenu

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取ManageMenu失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建ManageMenu
// @Summary 创建ManageMenu
// @Description 创建ManageMenu
// @Tags ManageMenu
// @Accept application/json
// @Product application/json
// @Param data body dto.ManageMenuInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/uc/manage-menu [post]
// @Security Bearer
func (e ManageMenu) Insert(c *gin.Context) {
	req := dto.ManageMenuInsertReq{}
	s := service.ManageMenu{}
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
		e.Error(500, err, fmt.Sprintf("创建ManageMenu失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改ManageMenu
// @Summary 修改ManageMenu
// @Description 修改ManageMenu
// @Tags ManageMenu
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.ManageMenuUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/uc/manage-menu/{id} [put]
// @Security Bearer
func (e ManageMenu) Update(c *gin.Context) {
	req := dto.ManageMenuUpdateReq{}
	s := service.ManageMenu{}
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
		e.Error(500, err, fmt.Sprintf("修改ManageMenu失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除ManageMenu
// @Summary 删除ManageMenu
// @Description 删除ManageMenu
// @Tags ManageMenu
// @Param data body dto.ManageMenuDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/uc/manage-menu [delete]
// @Security Bearer
func (e ManageMenu) Delete(c *gin.Context) {
	s := service.ManageMenu{}
	req := dto.ManageMenuDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除ManageMenu失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}
