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

type StockTransfer struct {
	api.Api
}

// GetPage 获取调拨单列表
// @Summary 获取调拨单列表
// @Description 获取调拨单列表
// @Tags 调拨单
// @Param transferCode query string false "调拨单编码"
// @Param type query string false "调拨类型: 0 -正常调拨 1-次品调拨 2-正转次调拨 3-次转正调拨"
// @Param status query string false "状态:0-已作废 1-未出库未入库 2-已出库未入库 3-出库完成入库完成 98已提交 99未提交"
// @Param fromWarehouseCode query string false "出库实体仓code"
// @Param fromLogicWarehouseCode query string false "出库逻辑仓code"
// @Param toWarehouseCode query string false "入库实体仓code"
// @Param toLogicWarehouseCode query string false "入库逻辑仓code"
// @Param sourceCode query string false "来源单号"
// @Param verifyStatus query string false "审核状态 0 待审核 1 审核通过 2 审核驳回 99初始化"
// @Param vendorId query int false "货主id"
// @Param createdAtStart query string false "创建时间开始"
// @Param createdAtEnd query string false "创建时间结束"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]dto.StockTransferGetPageResp}} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/stock-transfer [get]
// @Security Bearer
func (e StockTransfer) GetPage(c *gin.Context) {
	req := dto.StockTransferGetPageReq{}
	s := service.StockTransfer{}
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
	list := make([]dto.StockTransferGetPageResp, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取调拨单失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取调拨单
// @Summary 获取调拨单
// @Description 获取调拨单
// @Tags 调拨单
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=dto.StockTransferGetResp} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/stock-transfer/{id} [get]
// @Security Bearer
func (e StockTransfer) Get(c *gin.Context) {
	req := dto.StockTransferGetReq{}
	s := service.StockTransfer{}
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
	var object dto.StockTransferGetResp

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取调拨单失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 保存调拨单
// @Summary 保存调拨单
// @Description 保存调拨单
// @Tags 调拨单
// @Accept application/json
// @Product application/json
// @Param data body dto.StockTransferInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "创建成功"}"
// @Router /api/v1/wc/admin/stock-transfer [post]
// @Security Bearer
func (e StockTransfer) Insert(c *gin.Context) {
	req := dto.StockTransferInsertReq{}
	s := service.StockTransfer{}
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
	p := actions.GetPermissionFromContext(c)

	err = s.Insert(&req, "", p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("保存调拨单失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// SaveCommit 保存并提交调拨单
// @Summary 保存并提交调拨单
// @Description 保存并提交调拨单
// @Tags 调拨单
// @Accept application/json
// @Product application/json
// @Param data body dto.StockTransferInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "保存并提交成功"}"
// @Router /api/v1/wc/admin/stock-transfer/save-commit [post]
// @Security Bearer
func (e StockTransfer) SaveCommit(c *gin.Context) {
	req := dto.StockTransferInsertReq{}
	s := service.StockTransfer{}
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
	p := actions.GetPermissionFromContext(c)

	err = s.Insert(&req, "Commit", p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("保存并提交调拨单失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "保存并提交成功")
}

// Update 修改调拨单
// @Summary 修改调拨单
// @Description 修改调拨单
// @Tags 调拨单
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.StockTransferUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/wc/admin/stock-transfer/{id} [put]
// @Security Bearer
func (e StockTransfer) Update(c *gin.Context) {
	req := dto.StockTransferUpdateReq{}
	s := service.StockTransfer{}
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

	err = s.Update(&req, p, "")
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改调拨单失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// UpdateCommit 修改并提交调拨单
// @Summary 修改并提交调拨单
// @Description 修改并提交调拨单
// @Tags 调拨单
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.StockTransferUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改并提交调拨单"}"
// @Router /api/v1/wc/admin/stock-transfer/update-commit/{id} [put]
// @Security Bearer
func (e StockTransfer) UpdateCommit(c *gin.Context) {
	req := dto.StockTransferUpdateReq{}
	s := service.StockTransfer{}
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

	err = s.Update(&req, p, "Commit")
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改并提交调拨单失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改并提交调拨单")
}

// Delete 作废调拨单
// @Summary 作废调拨单
// @Description 作废调拨单
// @Tags 调拨单
// @Param data body dto.StockTransferDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "作废成功"}"
// @Router /api/v1/wc/admin/stock-transfer [delete]
// @Security Bearer
func (e StockTransfer) Delete(c *gin.Context) {
	s := service.StockTransfer{}
	req := dto.StockTransferDeleteReq{}
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

	err = s.Remove(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("作废调拨单失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "作废成功")
}

// Audit 审核调拨单
// @Summary 审核调拨单
// @Description 审核调拨单
// @Tags 调拨单
// @Accept application/json
// @Product application/json
// @Param data body dto.StockTransferAuditReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "审核成功"}"
// @Router /api/v1/wc/admin/stock-transfer/audit [post]
// @Security Bearer
func (e StockTransfer) Audit(c *gin.Context) {
	req := dto.StockTransferAuditReq{}
	s := service.StockTransfer{}
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

	err = s.Audit(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("审核调拨单失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "审核成功")
}

// ValidateSkus 调拨单验证skus
// @Summary 调拨单验证skus
// @Description 调拨单验证skus
// @Tags 调拨单
// @Param skuCodes query string false "多个用,分割"
// @Param vendorId query int false "货主id"
// @Param fromWarehouseCode query string false "出库仓code"
// @Param toWarehouseCode query string false "入库仓code"
// @Success 200 {object} response.Response{data=[]dto.StockTransferValidateSkusResp} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/stock-transfer/validate-skus [get]
// @Security Bearer
func (e StockTransfer) ValidateSkus(c *gin.Context) {
	req := dto.StockTransferValidateSkusReq{}
	s := service.StockTransfer{}
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
	var object []dto.StockTransferValidateSkusResp

	p := actions.GetPermissionFromContext(c)
	err = s.ValidateSkus(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf(err.Error()))
		return
	}

	e.OK(object, "查询成功")
}
