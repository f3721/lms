package admin

import (
	"fmt"
	"go-admin/common/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/pc/models"
	service "go-admin/app/pc/service/admin"
	"go-admin/app/pc/service/admin/dto"
	"go-admin/common/actions"
)

type ProductLog struct {
	api.Api
}

// GetPage 获取商品档案日志表列表
// @Summary 获取商品档案日志表列表
// @Description 获取商品档案日志表列表
// @Tags 商品档案日志表
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.ProductLog}} "{"code": 200, "data": [...]}"
// @Router /api/v1/pc/product-log [get]
// @Security Bearer
func (e ProductLog) GetPage(c *gin.Context) {
	req := dto.ProductLogGetPageReq{}
	s := service.ProductLog{}
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
	list := make([]models.ProductLog, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取商品档案日志表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取商品档案日志表
// @Summary 获取商品档案日志表
// @Description 获取商品档案日志表
// @Tags 商品档案日志表
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.ProductLog} "{"code": 200, "data": [...]}"
// @Router /api/v1/pc/product-log/{id} [get]
// @Security Bearer
func (e ProductLog) Get(c *gin.Context) {
	req := dto.ProductLogGetReq{}
	s := service.ProductLog{}
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
		e.Error(500, err, fmt.Sprintf("获取商品档案日志表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建商品档案日志表
// @Summary 创建商品档案日志表
// @Description 创建商品档案日志表
// @Tags 商品档案日志表
// @Accept application/json
// @Product application/json
// @Param data body dto.ProductLogInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/pc/product-log [post]
// @Security Bearer
func (e ProductLog) Insert(c *gin.Context) {
	req := dto.ProductLogInsertReq{}
	s := service.ProductLog{}
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
		e.Error(500, err, fmt.Sprintf("创建商品档案日志表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改商品档案日志表
// @Summary 修改商品档案日志表
// @Description 修改商品档案日志表
// @Tags 商品档案日志表
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.ProductLogUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/pc/product-log/{id} [put]
// @Security Bearer
func (e ProductLog) Update(c *gin.Context) {
	req := dto.ProductLogUpdateReq{}
	s := service.ProductLog{}
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
		e.Error(500, err, fmt.Sprintf("修改商品档案日志表失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除商品档案日志表
// @Summary 删除商品档案日志表
// @Description 删除商品档案日志表
// @Tags 商品档案日志表
// @Param data body dto.ProductLogDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/pc/product-log [delete]
// @Security Bearer
func (e ProductLog) Delete(c *gin.Context) {
	s := service.ProductLog{}
	req := dto.ProductLogDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除商品档案日志表失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}
