package mall

import (
	"fmt"
	"go-admin/common/middleware/mall_handler"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/wc/models"
	service "go-admin/app/wc/service/mall"
	"go-admin/app/wc/service/mall/dto"
	"go-admin/common/actions"
)

type Warehouse struct {
	api.Api
}

// GetPageForAuthedUser 获取登录后用户的实体仓列表
// @Summary 获取登录后用户的实体仓列表
// @Description 获取登录后用户的实体仓列表
// @Tags Mall登录
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.Warehouse}} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/mall/warehouse/list [get]
// @Security Bearer
func (e Warehouse) GetPageForAuthedUser(c *gin.Context) {
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

	//获取公司用户ID
	companyId, err := mall_handler.GetUserCompanyID(c)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取用户公司ID失败: %s", err.Error()))
		return
	}
	p := actions.GetPermissionFromContext(c)
	list := make([]models.Warehouse, 0)
	var count int64

	err = s.GetPageWithCompanyId(companyId, &req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取实体仓失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// GetPage 获取实体仓列表
// @Summary 获取实体仓列表
// @Description 获取实体仓列表
// @Tags 实体仓
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.Warehouse}} "{"code": 200, "data": [...]}"
// @Router /api/v1/warehouse [get]
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
	list := make([]models.Warehouse, 0)
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
// @Success 200 {object} response.Response{data=models.Warehouse} "{"code": 200, "data": [...]}"
// @Router /api/v1/warehouse/{id} [get]
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
	var object models.Warehouse

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
// @Router /api/v1/warehouse [post]
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
// @Router /api/v1/warehouse/{id} [put]
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
	p := actions.GetPermissionFromContext(c)

	err = s.Update(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改实体仓失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除实体仓
// @Summary 删除实体仓
// @Description 删除实体仓
// @Tags 实体仓
// @Param data body dto.WarehouseDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/warehouse [delete]
// @Security Bearer
func (e Warehouse) Delete(c *gin.Context) {
	s := service.Warehouse{}
	req := dto.WarehouseDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除实体仓失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}
