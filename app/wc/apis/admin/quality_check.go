package admin

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"go-admin/common/excel"

	service "go-admin/app/wc/service/admin"
	"go-admin/app/wc/service/admin/dto"
	"go-admin/common/actions"
)

type QualityCheck struct {
	api.Api
}

// GetPage 获取质检任务列表
// @Summary 获取质检任务列表
// @Description 获取质检任务列表
// @Tags 质检任务
// @Param id query int false "id"
// @Param qualityCheckCode query string false "质检单号"
// @Param sourceCode query string false "来源单据code"
// @Param entryCode query string false "入库单编码"
// @Param sourceName query string false "来源方"
// @Param status query int false "状态:0-已作废 1-创建 2-已完成"
// @Param warehouseCode query string false "实体仓code"
// @Param logicWarehouseCode query string false "逻辑仓code"
// @Param qualityRes query int false "质检结果"
// @Param type query int false "质检类型,1全检，2抽检"
// @Param qualityStatus query int false "质检进度：0未质检，1部分质检，2全部质检"
// @Param skuCode query string false "skuCode"
// @Param goodsName query string false "商品名称"
// @Param startDate query string false "开始日期"
// @Param endDate query string false "结束日期"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.QualityCheck}} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/quality-check [get]
// @Security Bearer
func (e QualityCheck) GetPage(c *gin.Context) {
	req := dto.QualityCheckGetPageReq{}
	s := service.QualityCheck{}
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
	list := make([]dto.QualityCheckRes, 0)
	var count int64
	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取质检任务失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取质检任务
// @Summary 获取质检任务
// @Description 获取质检任务
// @Tags 质检任务
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.QualityCheck} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/quality-check/{id} [get]
// @Security Bearer
func (e QualityCheck) Get(c *gin.Context) {
	req := dto.QualityCheckGetReq{}
	s := service.QualityCheck{}
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
	var object dto.QualityCheckRes
	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取质检任务失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(object, "查询成功")
}

// Get 上传质检结果
// @Summary 上传质检结果
// @Description 上传质检结果
// @Tags 上传质检结果
// @Param id path int false "id"
// @Accept application/json
// @Product application/json
// @Param data body dto.QualityCheckUpdateReq true "data"
// @Success 200 {object} response.Response{data=models.QualityCheck} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/quality-check/uploadquality/{id} [get]
// @Security Bearer
func (e QualityCheck) UploadQuality(c *gin.Context) {
	req := dto.QualityCheckUpdateReq{}
	s := service.QualityCheck{}
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
	err = s.UploadQualityRes(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("上传质检任务失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK("", "上传成功")
}

// Get 质检导出
// @Summary 质检导出
// @Description 质检导出
// @Tags 质检导出
// @Param id query int false "id"
// @Param qualityCheckCode query string false "质检单号"
// @Param sourceCode query string false "来源单据code"
// @Param entryCode query string false "入库单编码"
// @Param sourceName query string false "来源方"
// @Param status query int false "状态:0-已作废 1-创建 2-已完成"
// @Param warehouseCode query string false "实体仓code"
// @Param logicWarehouseCode query string false "逻辑仓code"
// @Param qualityRes query int false "质检结果"
// @Param type query int false "质检类型,1全检，2抽检"
// @Param qualityStatus query int false "质检进度：0未质检，1部分质检，2全部质检"
// @Param skuCode query string false "skuCode"
// @Param goodsName query string false "商品名称"
// @Param startDate query string false "开始日期"
// @Param endDate query string false "结束日期"
// @Success 200 {object} response.Response{data=models.QualityCheck} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/quality-check/{id} [get]
// @Security Bearer
func (e QualityCheck) Export(c *gin.Context) {
	req := dto.QualityCheckGetPageReq{}
	s := service.QualityCheck{}
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
	outData, err := s.Export(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("质检导出数据失败： %s", err.Error()))
		return
	}
	title := []map[string]string{
		{"ID": "ID"},
		{"qualityCheckCode": "质检单号"},
		{"sourceCode": "来源单号"},
		{"entryCode": "入库单号"},
		{"warehouseCode": "实体仓"},
		{"logicWarehouseCode": "逻辑仓"},
		{"sourceName": "来源方"},
		{"status": "状态"},
		{"type": "质检类型"},
		{"skuCode": "驿站SKU"},
		{"qualityStatus": "质检进度"},
		{"stayQualityNum": "待质检数量"},
		{"quantityNum": "质检数量"},
		{"qualityRes": "质检结果"},
		{"qualityTime": "质检时间"},
	}
	exportData := []map[string]interface{}{}
	byteJson, _ := json.Marshal(outData)
	_ = json.Unmarshal(byteJson, &exportData)
	excelApp := excel.NewExcel()
	err = excelApp.ExportExcelByMap(e.Context, title, exportData, "质检单-导出", "质检单明细")
	if err != nil {
		e.Error(500, err, fmt.Sprintf("质检单导出失败： %s", err.Error()))
	}
}
