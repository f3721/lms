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

type Address struct {
	api.Api
}

// GetPage 获取用户收货地址列表
// @Summary 获取用户收货地址列表
// @Description 获取用户收货地址列表
// @Tags 用户收货地址
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.Address}} "{"code": 200, "data": [...]}"
// @Router /api/v1/uc/mall/address [get]
// @Security Bearer
func (e Address) GetPage(c *gin.Context) {
	req := dto.AddressGetPageReq{}
	s := service.Address{}
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
	req.UserId = user.GetUserId(c)
	list := make([]models.Address, 0)
	var count int64

	// 商城-我的地址默认只显示 收货地址类型
	if req.AddressType == 0 {
		req.AddressType = 1
	}

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取用户收货地址失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取用户收货地址
// @Summary 获取用户收货地址
// @Description 获取用户收货地址
// @Tags 用户收货地址
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.Address} "{"code": 200, "data": [...]}"
// @Router /api/v1/uc/mall/address/{id} [get]
// @Security Bearer
func (e Address) Get(c *gin.Context) {
	req := dto.AddressGetReq{}
	s := service.Address{}
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
	var object models.Address

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取用户收货地址失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建用户收货地址
// @Summary 创建用户收货地址
// @Description 创建用户收货地址
// @Tags 用户收货地址
// @Accept application/json
// @Product application/json
// @Param data body dto.AddressInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/uc/mall/address [post]
// @Security Bearer
func (e Address) Insert(c *gin.Context) {
	req := dto.AddressInsertReq{}
	s := service.Address{}
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
		e.Error(500, err, fmt.Sprintf("创建用户收货地址失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改用户收货地址
// @Summary 修改用户收货地址
// @Description 修改用户收货地址
// @Tags 用户收货地址
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.AddressUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/uc/mall/address/{id} [put]
// @Security Bearer
func (e Address) Update(c *gin.Context) {
	req := dto.AddressUpdateReq{}
	s := service.Address{}
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
	req.UserId = user.GetUserId(c)

	err = s.Update(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改用户收货地址失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除用户收货地址
// @Summary 删除用户收货地址
// @Description 删除用户收货地址
// @Tags 用户收货地址
// @Param data body dto.AddressDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/uc/mall/address [delete]
// @Security Bearer
func (e Address) Delete(c *gin.Context) {
	s := service.Address{}
	req := dto.AddressDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除用户收货地址失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}

// Default 设置默认地址
// @Summary 设置默认地址
// @Description 设置默认地址
// @Tags 用户收货地址
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.AddressDefaultReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "设置成功"}"
// @Router /api/v1/uc/mall/address/default/{id} [put]
// @Security Bearer
func (e Address) Default(c *gin.Context) {
	req := dto.AddressDefaultReq{}
	s := service.Address{}
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
	req.UserId = user.GetUserId(c)

	err = s.Default(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("设置默认地址失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "设置成功")
}
