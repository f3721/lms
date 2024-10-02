package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"github.com/jinzhu/copier"
	"go-admin/app/wc/models"
	service "go-admin/app/wc/service/admin"
	"go-admin/app/wc/service/admin/dto"
	"go-admin/common/actions"
	"go-admin/common/excel"
)

type StockLocation struct {
	api.Api
}

// GetPage 获取库位列表
// @Summary 获取库位列表
// @Description 获取库位列表
// @Tags 库位
// @Param locationCode query string false "库位编码"
// @Param warehouseCode query string false "实体仓"
// @Param logicWarehouseCode query string false "逻辑仓"
// @Param status query int false "是否启用"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.StockLocation}} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/stock-location [get]
// @Security Bearer
func (e StockLocation) GetPage(c *gin.Context) {
	req := dto.StockLocationGetPageReq{}
	s := service.StockLocation{}
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
	list := make([]dto.StockLocationResp, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取库位失败，\r\n失败信息 %s", err.Error()))
		return
	}

	//var list2 []*dto.StockLocationResp
	//for _, v := range list {
	//	temp := dto.StockLocationResp{}
	//	_ = copier.Copy(&temp, &v)
	//	list2 = append(list2, &temp)
	//}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取库位
// @Summary 获取库位
// @Description 获取库位
// @Tags 库位
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=dto.StockLocationResp} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/stock-location/{id} [get]
// @Security Bearer
func (e StockLocation) Get(c *gin.Context) {
	req := dto.StockLocationGetReq{}
	s := service.StockLocation{}
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
	var object models.StockLocation

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取库位失败，\r\n失败信息 %s", err.Error()))
		return
	}

	data := dto.StockLocationResp{}
	_ = copier.Copy(&data, &object)

	e.OK(data, "查询成功")
}

// Insert 创建库位
// @Summary 创建库位
// @Description 创建库位
// @Tags 库位
// @Accept application/json
// @Product application/json
// @Param data body dto.StockLocationInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/wc/stock-location [post]
// @Security Bearer
func (e StockLocation) Insert(c *gin.Context) {
	req := dto.StockLocationInsertReq{}
	s := service.StockLocation{}
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
		e.Error(500, err, fmt.Sprintf("创建库位失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改库位
// @Summary 修改库位
// @Description 修改库位
// @Tags 库位
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.StockLocationUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/wc/stock-location/{id} [put]
// @Security Bearer
func (e StockLocation) Update(c *gin.Context) {
	req := dto.StockLocationUpdateReq{}
	s := service.StockLocation{}
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
		e.Error(500, err, fmt.Sprintf("修改库位失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK( req.GetId(), "修改成功")
}

// Import 导入库位
// @Summary 导入库位
// @Description 导入库位
// @Tags 库位
// @Accept application/json
// @Product application/json
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/wc/stock-location/import [post]
// @Security Bearer
func (e StockLocation) Import(c *gin.Context) {
	s := service.StockLocation{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
	}

	importFile, err := c.FormFile("file")
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	err, errTitleList, errData := s.Import(importFile)

	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
	} else {
		if len(errData) > 0 {
			// 有校验错误的行时调用 导出保存excel
			excelApp := excel.NewExcel()
			_ = excelApp.ExportExcelByMap(c, errTitleList, errData, "导入库位错误提示", "sheet1")
		} else {
			e.OK(nil, "导入成功")
		}
	}
}
