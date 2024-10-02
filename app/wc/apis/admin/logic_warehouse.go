package admin

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	service "go-admin/app/wc/service/admin"
	"go-admin/app/wc/service/admin/dto"
	"go-admin/common/actions"
)

type LogicWarehouse struct {
	api.Api
}

// GetPage 获取逻辑仓列表
// @Summary 获取逻辑仓列表
// @Description 获取逻辑仓列表
// @Tags 逻辑仓
// @Param logicWarehouseCode query string false "逻辑仓库编码"
// @Param logicWarehouseName query string false "逻辑仓库名称"
// @Param warehouseCode query string false "逻辑仓库对应实体仓code"
// @Param type query string false "0 正品仓 1次品仓"
// @Param status query string false "是否使用 0-否，1-是"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.LogicWarehouse}} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/logic-warehouse [get]
// @Security Bearer
func (e LogicWarehouse) GetPage(c *gin.Context) {
	req := dto.LogicWarehouseGetPageReq{}
	s := service.LogicWarehouse{}
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
	list := make([]dto.LogicWarehouseGetPageResp, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取逻辑仓失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取逻辑仓
// @Summary 获取逻辑仓
// @Description 获取逻辑仓
// @Tags 逻辑仓
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.LogicWarehouse} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/logic-warehouse/{id} [get]
// @Security Bearer
func (e LogicWarehouse) Get(c *gin.Context) {
	req := dto.LogicWarehouseGetReq{}
	s := service.LogicWarehouse{}
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
	var object dto.LogicWarehouseGetResp

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取逻辑仓失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建逻辑仓
// @Summary 创建逻辑仓
// @Description 创建逻辑仓
// @Tags 逻辑仓
// @Accept application/json
// @Product application/json
// @Param data body dto.LogicWarehouseInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/wc/admin/logic-warehouse [post]
// @Security Bearer
func (e LogicWarehouse) Insert(c *gin.Context) {
	req := dto.LogicWarehouseInsertReq{}
	s := service.LogicWarehouse{}
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
		e.Error(500, err, fmt.Sprintf("创建逻辑仓失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改逻辑仓
// @Summary 修改逻辑仓
// @Description 修改逻辑仓
// @Tags 逻辑仓
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.LogicWarehouseUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/wc/admin/logic-warehouse/{id} [put]
// @Security Bearer
func (e LogicWarehouse) Update(c *gin.Context) {
	req := dto.LogicWarehouseUpdateReq{}
	s := service.LogicWarehouse{}
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
		e.Error(500, err, fmt.Sprintf("修改逻辑仓失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Select 获取逻辑仓下拉数据
// @Summary 获取逻辑仓下拉数据
// @Description 获取逻辑仓下拉数据
// @Tags 逻辑仓
// @Param logicWarehouseCode query string false "逻辑仓库编码"
// @Param logicWarehouseName query string false "逻辑仓库名称"
// @Param warehouseCode query string false "逻辑仓库对应实体仓code"
// @Param type query string false "0 正品仓 1次品仓"
// @Param isTransfer query string false "调拨单获取仓库时为1"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]dto.LogicWarehouseSelectResp}} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/logic-warehouse/select [get]
// @Security Bearer
func (e LogicWarehouse) Select(c *gin.Context) {
	req := dto.LogicWarehouseSelectReq{}
	s := service.LogicWarehouse{}
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
	list := make([]dto.LogicWarehouseSelectResp, 0)
	var count int64

	err = s.Select(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取逻辑仓下拉数据，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}
