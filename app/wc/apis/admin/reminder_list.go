package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"github.com/jinzhu/copier"
	"go-admin/app/wc/models"
	service "go-admin/app/wc/service/admin"
	"go-admin/app/wc/service/admin/dto"
	"go-admin/common/actions"
	"go-admin/common/excel"
)

type ReminderList struct {
	api.Api
}

// GetPage 获取补货清单列表
// @Summary 获取补货清单列表
// @Description 获取补货清单列表
// @Tags 补货清单
// @Param reminderRuleId query int false "sxyz_reminder_rule 补货规则表id"
// @Param companyId query int false "公司id"
// @Param warehouseCode query int false "仓库id"
// @Param VendorId query int false "货主id"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]dto.ReminderListData}} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/reminder-list [get]
// @Security Bearer
func (e ReminderList) GetPage(c *gin.Context) {
	req := dto.ReminderListGetPageReq{}
	s := service.ReminderList{}
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
	list := make([]*dto.ReminderListData, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取补货清单失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取补货清单
// @Summary 获取补货清单
// @Description 获取补货清单
// @Tags 补货清单
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=dto.ReminderListData} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/reminder-list/{id} [get]
// @Security Bearer
func (e ReminderList) Get(c *gin.Context) {
	req := dto.ReminderListGetReq{}
	s := service.ReminderList{}
	listSkuService := service.ReminderListSku{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		MakeService(&listSkuService.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	var object dto.ReminderListData

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取补货清单失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Export 导出补货清单
// @Summary 导出补货清单
// @Description 导出补货清单
// @Tags 补货清单
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.ReminderList} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/reminder-list/export/{id} [get]
// @Security Bearer
func (e ReminderList) Export(c *gin.Context) {
	req := dto.ReminderListGetReq{}
	s := service.ReminderList{}
	listSkuService := service.ReminderListSku{}
	vendorsService := service.Vendors{}
	warehouseService := service.Warehouse{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		MakeService(&listSkuService.Service).
		MakeService(&vendorsService.Service).
		MakeService(&warehouseService.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	var reminderList dto.ReminderListData

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &reminderList)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取补货清单失败，\r\n失败信息 %s", err.Error()))
		return
	}

	skuList := []models.ReminderListSku{}
	listSkuService.GetExportAll(&dto.ReminderListSkuGetPageReq{ReminderListId: int(reminderList.Id)}, p, &skuList)

	var exportSkulist []interface{}
	for _, sku := range skuList {
		addExportSkulist := dto.ExportSkuData{}
		copier.Copy(&addExportSkulist, sku)
		addExportSkulist.VendorName = reminderList.VendorName
		addExportSkulist.WarehouseName = reminderList.WarehouseName
		addExportSkulist.CompanyName = reminderList.CompanyName
		addExportSkulist.Date = reminderList.CreatedAt.Format("2006-01-02")
		exportSkulist = append(exportSkulist, addExportSkulist)
	}

	title := []map[string]string{
		{"companyName": "公司名称"},
		{"warehouseName": "仓库名称"},
		{"date": "日期"},
		{"skuCode": "商品编码"},
		{"supplierName": "供应商名称"},
		{"supplierInfoSku": "供应商商品编码"},
		{"warningValue": "预警值"},
		{"replenishmentValue": "补货值"},
		{"recommendReplenishmentValue": "推荐补货值"},
		{"genuineStock": "实际库存"},
		{"allStock": "总库存"},
		{"occupyStock": "占用库存"},
	}

	excelApp := excel.NewExcel()
	err = excelApp.ExportExcelByStruct(c, title, exportSkulist, "reminder-list", "sheet1")
	if err != nil {
		e.Error(500, err, fmt.Sprintf("导出部门失败： %s", err.Error()))
		return
	}
}

// Create 创建补货清单
// @Summary 自动创建补货清单
// @Description 自动创建补货清单
// @Tags 补货清单
// @Accept application/json
// @Product application/json
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/wc/admin/reminder-list/create [get]
// @Security Bearer
func (e ReminderList) Create(c *gin.Context) {

	//Service
	reminderListService := service.ReminderList{}

	_ = e.MakeContext(c).
		//MakeOrm().
		MakeService(&reminderListService.Service).
		Errors

	p := actions.GetPermissionFromContext(c)

	_, err := reminderListService.TenantsCreate(c, p)
	e.Logger.Error(err)
	e.OK(nil, "创建成功")
	return
}

//// Insert 创建补货清单
//// @Summary 创建补货清单
//// @Description 创建补货清单
//// @Tags 补货清单
//// @Accept application/json
//// @Product application/json
//// @Param data body dto.ReminderListInsertReq true "data"
//// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
//// @Router /api/v1/wc/admin/reminder-list [post]
//// @Security Bearer
//func (e ReminderList) Insert(c *gin.Context) {
//	req := dto.ReminderListInsertReq{}
//	s := service.ReminderList{}
//	err := e.MakeContext(c).
//		MakeOrm().
//		Bind(&req).
//		MakeService(&s.Service).
//		Errors
//	if err != nil {
//		e.Logger.Error(err)
//		e.Error(500, err, err.Error())
//		return
//	}
//	// 设置创建人
//	req.SetCreateBy(user.GetUserId(c))
//
//	err = s.Insert(&req)
//	if err != nil {
//		e.Error(500, err, fmt.Sprintf("创建补货清单失败，\r\n失败信息 %s", err.Error()))
//		return
//	}
//
//	e.OK(req.GetId(), "创建成功")
//}
//
//// Update 修改补货清单
//// @Summary 修改补货清单
//// @Description 修改补货清单
//// @Tags 补货清单
//// @Accept application/json
//// @Product application/json
//// @Param id path int true "id"
//// @Param data body dto.ReminderListUpdateReq true "body"
//// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
//// @Router /api/v1/wc/admin/reminder-list/{id} [put]
//// @Security Bearer
//func (e ReminderList) Update(c *gin.Context) {
//	req := dto.ReminderListUpdateReq{}
//	s := service.ReminderList{}
//	err := e.MakeContext(c).
//		MakeOrm().
//		Bind(&req).
//		MakeService(&s.Service).
//		Errors
//	if err != nil {
//		e.Logger.Error(err)
//		e.Error(500, err, err.Error())
//		return
//	}
//	req.SetUpdateBy(user.GetUserId(c))
//	p := actions.GetPermissionFromContext(c)
//
//	err = s.Update(&req, p)
//	if err != nil {
//		e.Error(500, err, fmt.Sprintf("修改补货清单失败，\r\n失败信息 %s", err.Error()))
//		return
//	}
//	e.OK(req.GetId(), "修改成功")
//}
//
//// Delete 删除补货清单
//// @Summary 删除补货清单
//// @Description 删除补货清单
//// @Tags 补货清单
//// @Param data body dto.ReminderListDeleteReq true "body"
//// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
//// @Router /api/v1/wc/admin/reminder-list [delete]
//// @Security Bearer
//func (e ReminderList) Delete(c *gin.Context) {
//	s := service.ReminderList{}
//	req := dto.ReminderListDeleteReq{}
//	err := e.MakeContext(c).
//		MakeOrm().
//		Bind(&req).
//		MakeService(&s.Service).
//		Errors
//	if err != nil {
//		e.Logger.Error(err)
//		e.Error(500, err, err.Error())
//		return
//	}
//
//	// req.SetUpdateBy(user.GetUserId(c))
//	p := actions.GetPermissionFromContext(c)
//
//	err = s.Remove(&req, p)
//	if err != nil {
//		e.Error(500, err, fmt.Sprintf("删除补货清单失败，\r\n失败信息 %s", err.Error()))
//		return
//	}
//	e.OK(req.GetId(), "删除成功")
//}
