package admin

import (
	"fmt"
	"go-admin/app/uc/service/admin"
	"go-admin/app/uc/service/admin/dto"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/uc/models"
	"go-admin/common/actions"
)

type UserPayAccount struct {
	api.Api
}

// GetPage 获取客户分类对应支付账户表列表
// @Summary 获取客户分类对应支付账户表列表
// @Description 获取客户分类对应支付账户表列表
// @Tags 客户分类对应支付账户表
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.UserPayAccount}} "{"code": 200, "data": [...]}"
// @Router /api/v1/uc/admin/user-pay-account [get]
// @Security Bearer
func (e UserPayAccount) GetPage(c *gin.Context) {
	req := dto.UserPayAccountGetPageReq{}
	s := admin.UserPayAccount{}
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
	list := make([]models.UserPayAccount, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取客户分类对应支付账户表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取客户分类对应支付账户表
// @Summary 获取客户分类对应支付账户表
// @Description 获取客户分类对应支付账户表
// @Tags 客户分类对应支付账户表
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.UserPayAccount} "{"code": 200, "data": [...]}"
// @Router /api/v1/uc/admin/user-pay-account/{id} [get]
// @Security Bearer
func (e UserPayAccount) Get(c *gin.Context) {
	req := dto.UserPayAccountGetReq{}
	s := admin.UserPayAccount{}
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
	var object models.UserPayAccount

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取客户分类对应支付账户表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建客户分类对应支付账户表
// @Summary 创建客户分类对应支付账户表
// @Description 创建客户分类对应支付账户表
// @Tags 客户分类对应支付账户表
// @Accept application/json
// @Product application/json
// @Param data body dto.UserPayAccountInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/uc/admin/user-pay-account [post]
// @Security Bearer
func (e UserPayAccount) Insert(c *gin.Context) {
	req := dto.UserPayAccountInsertReq{}
	s := admin.UserPayAccount{}
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
		e.Error(500, err, fmt.Sprintf("创建客户分类对应支付账户表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改客户分类对应支付账户表
// @Summary 修改客户分类对应支付账户表
// @Description 修改客户分类对应支付账户表
// @Tags 客户分类对应支付账户表
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.UserPayAccountUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/uc/admin/user-pay-account/{id} [put]
// @Security Bearer
func (e UserPayAccount) Update(c *gin.Context) {
	req := dto.UserPayAccountUpdateReq{}
	s := admin.UserPayAccount{}
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
		e.Error(500, err, fmt.Sprintf("修改客户分类对应支付账户表失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除客户分类对应支付账户表
// @Summary 删除客户分类对应支付账户表
// @Description 删除客户分类对应支付账户表
// @Tags 客户分类对应支付账户表
// @Param data body dto.UserPayAccountDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/uc/admin/user-pay-account [delete]
// @Security Bearer
func (e UserPayAccount) Delete(c *gin.Context) {
	s := admin.UserPayAccount{}
	req := dto.UserPayAccountDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除客户分类对应支付账户表失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}
