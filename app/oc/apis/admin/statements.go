package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"github.com/jinzhu/copier"
	service "go-admin/app/oc/service/admin"
	"go-admin/app/oc/service/admin/dto"
	"go-admin/common/actions"
	"go-admin/common/excel"
	//"go-admin/common/utils"
)

type Statements struct {
	api.Api
}

type OrderToStatements struct {
	api.Api
}

// GetPage 获取对账单列表
// @Summary 获取对账单列表
// @Description 获取对账单列表
// @Tags 对账单
// @Param companyId query int false "公司id"
// @Param startTime query string false "创建时间开始"
// @Param endTime query string false "创建时间结束"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]dto.StatementsGetPageResp}} "{"code": 200, "data": [...]}"
// @Router /api/v1/oc/statements [get]
// @Security Bearer
func (e Statements) GetPage(c *gin.Context) {
	req := dto.StatementsGetPageReq{}
	s := service.Statements{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.Form).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	p := actions.GetPermissionFromContext(c)
	list := make([]dto.StatementsGetPageResp, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取对账单失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取对账单
// @Summary 获取对账单
// @Description 获取对账单
// @Tags 对账单
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.Statements} "{"code": 200, "data": [...]}"
// @Router /api/v1/oc/statements/{id} [get]
// @Security Bearer
func (e Statements) Get(c *gin.Context) {
	req := dto.StatementsGetReq{}
	s := service.Statements{}
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
	var object dto.StatementsGetPageResp

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取对账单失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// GetDetails 获取对账单详情
// @Summary 获取对账单详情
// @Description 获取对账单详情
// @Tags 对账单
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=[]dto.OrderToStatementsListResp} "{"code": 200, "data": [...]}"
// @Router /api/v1/oc/statements/details/{id} [get]
// @Security Bearer
func (e OrderToStatements) GetDetails(c *gin.Context) {
	req := dto.OrderToStatementsGetReq{}
	s := service.OrderToStatements{}
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

	companySkuSwitchKeyword := "sku_classification"
	CompanySwitch, _ := s.GetCompanySwitchInfo(req.Id, companySkuSwitchKeyword)

	list := make([]dto.OrderToStatementsListResp, 0)

	p := actions.GetPermissionFromContext(c)
	var count int64
	err = s.GetDetails(&req, p, &list, &count, CompanySwitch.SwitchStatus, false)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取对账单详情失败，\r\n失败信息 %s", err.Error()))
		return
	}

	for i, _ := range list {
		list[i].Key = i + 1
		if list[i].CreateFrom == "XCX" {
			list[i].CreateFrom = "小程序"
		}
		if list[i].ParentCompanyDepartmentName != "" {
			list[i].FullCompanyName = list[i].ParentCompanyDepartmentName + " > " + list[i].CompanyDepartmentName
		} else {
			list[i].FullCompanyName = list[i].CompanyDepartmentName
		}
	}

	var companySwitchPageResp dto.OrderToStatementsListPageResp
	companySwitchPageResp.CompanySwitch.Keyword = companySkuSwitchKeyword
	companySwitchPageResp.CompanySwitch.SwitchStatus = CompanySwitch.SwitchStatus
	companySwitchPageResp.List = list
	companySwitchPageResp.Count = count
	companySwitchPageResp.PageIndex = req.GetPageIndex()
	companySwitchPageResp.PageSize = req.GetPageSize()

	e.OK(companySwitchPageResp, "查询成功")
}

// DetailsExport 导出对账单
// @Summary 导出对账单
// @Description 导出对账单
// @Tags 对账单
// @Param id path int false "id"
// @Router /api/v1/oc/statements/details/export/{id} [get]
// @Security Bearer
func (e OrderToStatements) DetailsExport(c *gin.Context) {
	req := dto.OrderToStatementsGetReq{}
	s := service.OrderToStatements{}
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

	companySkuSwitchKeyword := "sku_classification"
	CompanySwitch, _ := s.GetCompanySwitchInfo(req.Id, companySkuSwitchKeyword)

	list := make([]dto.OrderToStatementsListResp, 0)

	p := actions.GetPermissionFromContext(c)
	var count int64
	err = s.GetDetails(&req, p, &list, &count, CompanySwitch.SwitchStatus, true)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("对账单导出失败，\r\n失败信息 %s", err.Error()))
		return
	}
	var exportData []interface{}

	for i, resp := range list {
		resp.Key = i + 1
		if resp.CreateFrom == "XCX" {
			resp.CreateFrom = "小程序"
		}
		if resp.ParentCompanyDepartmentName != "" {
			resp.FullCompanyName = resp.ParentCompanyDepartmentName + " > " + resp.CompanyDepartmentName
		} else {
			resp.FullCompanyName = resp.CompanyDepartmentName
		}

		exportOrderToStatementsList := dto.ExportOrderToStatementsList{}
		err := copier.Copy(&exportOrderToStatementsList, resp)
		if err != nil {
			return
		}

		exportData = append(exportData, exportOrderToStatementsList)
	}

	var title []map[string]string
	if CompanySwitch.SwitchStatus == "1" {
		title = []map[string]string{
			{"key": "序号"},
			{"orderId": "订单编号"},
			{"finalTotalAmount": "订单总金额"},
			{"skuCode": "SKU"},
			{"productName": "商品名称"},
			{"brandName": "品牌"},
			{"productModel": "型号"},
			{"productNo": "物料编码"},
			{"vendorName": "货主"},
			{"supplierSkuCode": "货主SKU"},
			{"salePrice": "销售价"},
			{"finalQuantity": "商品数量"},
			{"finalSubTotalAmount": "行项目小计金额"},
			{"userName": "领用人"},
			{"userId": "用户ID"},
			{"userPhone": "手机号"},
			{"fullCompanyName": "部门"},
			{"createFrom": "来源"},
			{"skuClassificationName": "客户对应产线"},
			{"payAccount": "客户支付账号"},
		}
	} else {
		title = []map[string]string{
			{"key": "序号"},
			{"orderId": "订单编号"},
			{"finalTotalAmount": "订单总金额"},
			{"skuCode": "SKU"},
			{"productName": "商品名称"},
			{"brandName": "品牌"},
			{"productModel": "型号"},
			{"productNo": "物料编码"},
			{"vendorName": "货主"},
			{"supplierSkuCode": "货主SKU"},
			{"salePrice": "销售价"},
			{"finalQuantity": "商品数量"},
			{"finalSubTotalAmount": "行项目小计金额"},
			{"userName": "领用人"},
			{"userId": "用户ID"},
			{"userPhone": "手机号"},
			{"fullCompanyName": "部门"},
			{"createFrom": "来源"},
		}
	}

	excelApp := excel.NewExcel()
	err = excelApp.ExportExcelByStruct(c, title, exportData, "对账单", "Sheet1")
	if err != nil {
		e.Error(500, err, fmt.Sprintf("对账单导出失败： %s", err.Error()))
	}
}

// InitStatements 生成对账单
func (e Statements) InitStatements(c *gin.Context) {
	s := service.Statements{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	echoMsg, err := s.InitStatements(c)
	if err != nil {
		c.String(500, err.Error())
	} else {
		c.String(200, echoMsg)
	}
}
