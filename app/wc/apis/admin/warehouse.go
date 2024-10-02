package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/wc/models"
	service "go-admin/app/wc/service/admin"
	"go-admin/app/wc/service/admin/dto"
	"go-admin/common/actions"
)

type Warehouse struct {
	api.Api
}

// GetPage 获取实体仓列表
// @Summary 获取实体仓列表
// @Description 获取实体仓列表
// @Tags 实体仓
// @Param warehouseCode query string false "仓库编码"
// @Param WarehouseName query string false "仓库名称"
// @Param companyId query int false "仓库对应公司id"
// @Param isVirtual query string false "是否为虚拟仓 0-否，1-是"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]dto.WarehouseGetPageResp}} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/warehouse [get]
// @Security Bearer
func (e Warehouse) GetPage(c *gin.Context) {
	req := dto.WarehouseGetPageReq{}
	s := service.Warehouse{}
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
	list := make([]dto.WarehouseGetPageResp, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取实体仓失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取实体仓
// @Summary 获取实体仓
// @Description 获取实体仓
// @Tags 实体仓
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=dto.WarehouseGetResp} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/warehouse/{id} [get]
// @Security Bearer
func (e Warehouse) Get(c *gin.Context) {
	req := dto.WarehouseGetReq{}
	s := service.Warehouse{}
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
	var object dto.WarehouseGetResp

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取实体仓失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建实体仓
// @Summary 创建实体仓
// @Description 创建实体仓
// @Tags 实体仓
// @Accept application/json
// @Product application/json
// @Param data body dto.WarehouseInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/wc/admin/warehouse [post]
// @Security Bearer
func (e Warehouse) Insert(c *gin.Context) {
	req := dto.WarehouseInsertReq{}
	s := service.Warehouse{}
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
		e.Error(500, err, fmt.Sprintf("创建实体仓失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改实体仓
// @Summary 修改实体仓
// @Description 修改实体仓
// @Tags 实体仓
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.WarehouseUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/wc/admin/warehouse/{id} [put]
// @Security Bearer
func (e Warehouse) Update(c *gin.Context) {
	req := dto.WarehouseUpdateReq{}
	s := service.Warehouse{}
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
		e.Error(500, err, fmt.Sprintf("修改实体仓失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Select 获取实体仓下拉数据
// @Summary 获取实体仓下拉数据
// @Description 获取实体仓下拉数据
// @Tags 实体仓
// @Param warehouseCode query string false "仓库编码"
// @Param WarehouseName query string false "仓库名称"
// @Param companyId query int false "仓库对应公司id"
// @Param isVirtual query string false "是否为虚拟仓 0-否，1-是"
// @Param isTransfer query string false "调拨单获取仓库时为1"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]dto.WarehouseSelectResp}} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/warehouse/select [get]
// @Security Bearer
func (e Warehouse) Select(c *gin.Context) {
	req := dto.WarehouseSelectReq{}
	s := service.Warehouse{}
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
	list := make([]dto.WarehouseSelectResp, 0)
	var count int64

	err = s.Select(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取实体仓失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// GetByCode 获取实体仓通过code
// @Summary 获取实体仓通过code
// @Description 获取实体仓通过code
// @Tags 实体仓
// @Param code path int false "code"
// @Success 200 {object} response.Response{data=models.Warehouse} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/warehouse/code/{code} [get]
// @Security Bearer
func (e Warehouse) GetByCode(c *gin.Context) {
	req := dto.WarehouseGetByCodeReq{}
	s := service.Warehouse{}
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
	var object models.Warehouse

	p := actions.GetPermissionFromContext(c)
	err = s.GetByWarehouseCode(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取实体仓失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// GetCompanyWarehouseTree 获取公司实体仓树
// @Summary 获取公司实体仓树
// @Description 获取公司实体仓树
// @Tags 实体仓
// @Success 200 {object} response.Response{data=models.Warehouse} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/warehouse/cw-tree [get]
// @Security Bearer
func (e Warehouse) GetCompanyWarehouseTree(c *gin.Context) {
	s := service.Warehouse{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	object, err := s.GetCompanyWarehouseTree()
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取公司仓库树失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// InnerGetListByNameAndCompanyId Inner获取实体仓列表通过公司Id和仓库名
// @Summary Inner获取实体仓列表通过公司Id和仓库名
// @Description Inner获取实体仓列表通过公司Id和仓库名
// @Tags Inner实体仓
// @Accept application/json
// @Product application/json
// @Param data body []dto.InnerWarehouseGetListReq true "body"
// @Success 200 {object} response.Response{data=[]models.Warehouse} "{"code": 200, "data": [...]}"
// @Router /inner/wc/admin/warehouse/list-by-name-company-id [post]
func (e Warehouse) InnerGetListByNameAndCompanyId(c *gin.Context) {
	req := dto.InnerWarehouseGetListByNameAndCompanyIdReq{}
	s := service.Warehouse{}
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

	list := make([]models.Warehouse, 0)

	err = s.InnerGetListByNameAndCompanyId(&req, &list)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("Inner获取实体仓列表通过公司Id和仓库名，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(list, "成功")
}

// InnerGetList Inner获取实体仓列表
// @Summary Inner获取实体仓列表
// @Description Inner获取实体仓列表
// @Tags Inner实体仓
// @Param warehouseCode query string false "仓库编码 多个用','分割"
// @Param WarehouseName query string false "仓库名称 多个用','分割"
// @Param companyId query int false "仓库对应公司id"
// @Param isVirtual query string false "是否为虚拟仓 0-否，1-是"
// @Param status query string false "是否为虚拟仓 0-否，1-是"
// @Success 200 {object} response.Response{data=[]dto.WarehouseGetPageResp} "{"code": 200, "data": [...]}"
// @Router /inner/wc/admin/warehouse/list [get]
func (e Warehouse) InnerGetList(c *gin.Context) {
	req := dto.InnerWarehouseGetListReq{}
	s := service.Warehouse{}
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

	list := make([]dto.WarehouseGetPageResp, 0)

	err = s.InnerGetList(&req, &list)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("Inner获取实体仓列表，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(list, "成功")
}
