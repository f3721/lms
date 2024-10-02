package admin

import (
	"fmt"
	"go-admin/common/excel"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	service "go-admin/app/wc/service/admin"
	"go-admin/app/wc/service/admin/dto"
	"go-admin/common/actions"
)

type StockControl struct {
	api.Api
}

// GetPage 获取库存调整单列表
// @Summary 获取库存调整单列表
// @Description 获取库存调整单列表
// @Tags 库存调整单
// @Param controlCode query string false "调整单编码"
// @Param type query string false "调整类型: 0 调增 1 调减"
// @Param status query string false "状态:0-已作废 1-创建 2-已完成 99未提交"
// @Param warehouseCode query string false "实体仓code"
// @Param logicWarehouseCode query string false "逻辑仓code"
// @Param verifyStatus query string false "审核状态 0 待审核 1 审核通过 2 审核驳回 99初始化"
// @Param vendorId query int false "货主id"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]dto.StockControlGetPageResp}} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/stock-control [get]
// @Security Bearer
func (e StockControl) GetPage(c *gin.Context) {
	req := dto.StockControlGetPageReq{}
	s := service.StockControl{}
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
	list := make([]dto.StockControlGetPageResp, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取库存调整单失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取库存调整单
// @Summary 获取库存调整单
// @Description 获取库存调整单
// @Tags 库存调整单
// @Param id path int false "id"
// @Param type path string false "eg: edit"
// @Success 200 {object} response.Response{data=dto.StockControlGetResp} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/stock-control/{id} [get]
// @Security Bearer
func (e StockControl) Get(c *gin.Context) {
	req := dto.StockControlGetReq{}
	s := service.StockControl{}
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
	var object dto.StockControlGetResp

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取库存调整单失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建库存调整单
// @Summary 创建库存调整单
// @Description 创建库存调整单
// @Tags 库存调整单
// @Accept application/json
// @Product application/json
// @Param data body dto.StockControlInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/wc/admin/stock-control [post]
// @Security Bearer
func (e StockControl) Insert(c *gin.Context) {
	req := dto.StockControlInsertReq{}
	s := service.StockControl{}
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
		e.Error(500, err, fmt.Sprintf("创建库存调整单失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// SaveCommit 保存并提交调整单
// @Summary 保存并提交调整单
// @Description 保存并提交调整单
// @Tags 库存调整单
// @Accept application/json
// @Product application/json
// @Param data body dto.StockControlInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "保存并提交成功"}"
// @Router /api/v1/wc/admin/stock-control/save-commit [post]
// @Security Bearer
func (e StockControl) SaveCommit(c *gin.Context) {
	req := dto.StockControlInsertReq{}
	s := service.StockControl{}
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
		e.Error(500, err, fmt.Sprintf("保存并提交成功库存调整单失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "保存并提交成功")
}

// Update 修改库存调整单
// @Summary 修改库存调整单
// @Description 修改库存调整单
// @Tags 库存调整单
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.StockControlUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/wc/admin/stock-control/{id} [put]
// @Security Bearer
func (e StockControl) Update(c *gin.Context) {
	req := dto.StockControlUpdateReq{}
	s := service.StockControl{}
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
		e.Error(500, err, fmt.Sprintf("修改库存调整单失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// UpdateCommit 修改并提交库存调整单
// @Summary 修改并提交库存调整单
// @Description 修改并提交库存调整单
// @Tags 库存调整单
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.StockControlUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改并提交成功"}"
// @Router /api/v1/wc/admin/stock-control/update-commit/{id} [put]
// @Security Bearer
func (e StockControl) UpdateCommit(c *gin.Context) {
	req := dto.StockControlUpdateReq{}
	s := service.StockControl{}
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
		e.Error(500, err, fmt.Sprintf("修改并提交成功库存调整单失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改并提交成功")
}

// Delete 作废库存调整单
// @Summary 作废库存调整单
// @Description 作废库存调整单
// @Tags 库存调整单
// @Param data body dto.StockControlDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "作废成功"}"
// @Router /api/v1/wc/admin/stock-control [delete]
// @Security Bearer
func (e StockControl) Delete(c *gin.Context) {
	s := service.StockControl{}
	req := dto.StockControlDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("作废库存调整单失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "作废成功")
}

// Audit 审核库存调整单
// @Summary 审核库存调整单
// @Description 审核库存调整单
// @Tags 库存调整单
// @Accept application/json
// @Product application/json
// @Param data body dto.StockControlAuditReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "审核成功"}"
// @Router /api/v1/wc/admin/stock-control/audit [post]
// @Security Bearer
func (e StockControl) Audit(c *gin.Context) {
	req := dto.StockControlAuditReq{}
	s := service.StockControl{}
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
		e.Error(500, err, fmt.Sprintf("审核调整单失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "审核成功")
}

// ValidateSkus 调整单验证skus
// @Summary 调整单验证skus
// @Description 调整单验证skus
// @Tags 库存调整单
// @Param skuCodes query string false "多个用,分割"
// @Param vendorId query int false "货主id"
// @Param warehouseCode query string false "实体仓code"
// @Param logicWarehouseCode query string false "logicWarehouseCode"
// @Success 200 {object} response.Response{data=[]dto.StockControlValidateSkusResp} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/stock-control/validate-skus [get]
// @Security Bearer
func (e StockControl) ValidateSkus(c *gin.Context) {
	req := dto.StockControlValidateSkusReq{}
	s := service.StockControl{}
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
	var object []dto.StockControlValidateSkusResp

	p := actions.GetPermissionFromContext(c)
	err = s.ValidateSkus(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf(err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// StockImport 库存调整单新增导入
// @Summary 库存调整单新增导入
// @Description 库存调整单新增导入
// @Tags 库存调整单
// @Param file formData file false "导入文件"
// @Router /api/v1/wc/admin/stock-control/import [post]
// @Security Bearer
func (e StockControl) StockImport(c *gin.Context) {
	s := service.StockControl{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	// 获取文件
	file, err := c.FormFile("file")
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	p := actions.GetPermissionFromContext(c)
	err, title, exportData := s.StockImport(file, p)
	if err != nil {
		if len(exportData) > 0 {
			excelApp := excel.NewExcel()
			err = excelApp.ExportExcelByMap(c, title, exportData, "库存调整单导入错误", "Sheet1")
			if err != nil {
				e.Error(500, err, fmt.Sprintf("库存调整单导入失败： %s", err.Error()))
				return
			}
		} else {
			e.Error(500, err, err.Error())
			return
		}
	} else {
		e.OK(200, "导入成功")
	}
}
