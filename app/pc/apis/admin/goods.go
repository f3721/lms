package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"go-admin/app/pc/models"
	service "go-admin/app/pc/service/admin"
	"go-admin/app/pc/service/admin/dto"
	"go-admin/common/actions"
	"go-admin/common/excel"
)

type Goods struct {
	api.Api
}

// GetPage 获取商品管理表列表
// @Summary 获取商品管理表列表
// @Description 获取商品管理表列表
// @Tags 商品管理表
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.Goods}} "{"code": 200, "data": [...]}"
// @Router /api/v1/pc/goods [get]
// @Security Bearer
func (e Goods) GetPage(c *gin.Context) {
	req := dto.GoodsGetPageReq{}
	s := service.Goods{}
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
	list := make([]dto.GoodsGetPageResp, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取商品管理表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取商品管理表
// @Summary 获取商品管理表
// @Description 获取商品管理表
// @Tags 商品管理表
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.Goods} "{"code": 200, "data": [...]}"
// @Router /api/v1/pc/goods/{id} [get]
// @Security Bearer
func (e Goods) Get(c *gin.Context) {
	req := dto.GoodsGetReq{}
	s := service.Goods{}
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
	var object dto.GoodsGetResp

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取商品管理表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建商品管理表
// @Summary 创建商品管理表
// @Description 创建商品管理表
// @Tags 商品管理表
// @Accept application/json
// @Product application/json
// @Param data body dto.GoodsInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/pc/goods [post]
// @Security Bearer
func (e Goods) Insert(c *gin.Context) {
	req := dto.GoodsInsertReq{}
	s := service.Goods{}
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
	// 设置创建人
	req.SetCreateBy(user.GetUserId(c))
	req.SetCreateByName(user.GetUserName(c))
	err = s.Insert(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建商品管理表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改商品管理表
// @Summary 修改商品管理表
// @Description 修改商品管理表
// @Tags 商品管理表
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.GoodsUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/pc/goods/{id} [put]
// @Security Bearer
func (e Goods) Update(c *gin.Context) {
	req := dto.GoodsUpdateReq{}
	s := service.Goods{}
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
		e.Error(500, err, fmt.Sprintf("修改商品管理表失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除商品管理表
// @Summary 删除商品管理表
// @Description 删除商品管理表
// @Tags 商品管理表
// @Param data body dto.GoodsDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/pc/goods [delete]
// @Security Bearer
func (e Goods) Delete(c *gin.Context) {
	s := service.Goods{}
	req := dto.GoodsDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除商品管理表失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}

// Approve 商品审核
// @Summary 商品审核
// @Description 商品审核
// @Tags 商品审核
// @Param data body dto.GoodsApproveReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "商品审核成功"}"
// @Router /api/v1/pc/goods/approve [post]
// @Security Bearer
func (e Goods) Approve(c *gin.Context) {
	s := service.Goods{}
	req := dto.GoodsApproveReq{}
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

	msg, err := s.Approve(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("商品审核失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(msg, "商品审核成功")
}

// OnlineOffline 商品批量上下架
// @Summary 商品批量上下架
// @Description 商品批量上下架
// @Tags 商品批量上下架
// @Param data body dto.OlineOfflineReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "商品批量上下架成功"}"
// @Router /api/v1/pc/goods/online-offline [post]
// @Security Bearer
func (e Goods) OnlineOffline(c *gin.Context) {
	s := service.Goods{}
	req := dto.OlineOfflineReq{}
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

	msg, err := s.OnlineOffline(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("商品审核失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(msg, "商品批量上下架成功")
}

// GoodsImport 产品新增导入
// @Summary 产品新增导入
// @Description 产品新增导入
// @Tags 产品新增导入
// @Router /api/v1/pc/goods/import [post]
// @Security Bearer
func (e Goods) GoodsImport(c *gin.Context) {
	s := service.Goods{}
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
	err, title, exportData := s.GoodsImport(file, p)
	if err != nil {
		if len(exportData) > 0 {
			excelApp := excel.NewExcel()
			err = excelApp.ExportExcelByMap(c, title, exportData, "商品管理导入错误", "Sheet1")
			if err != nil {
				e.Error(500, err, fmt.Sprintf("商品管理导入失败： %s", err.Error()))
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

// Export 商品导出
// @Summary 商品导出
// @Description 商品导出
// @Tags 商品导出
// @Param data body dto.GoodsGetPageReq true "body"
// @Router /api/v1/pc/goods/export [get]
// @Security Bearer
func (e Goods) Export(c *gin.Context) {
	s := service.Goods{}
	req := dto.GoodsGetPageReq{}
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
	exportData, err := s.Export(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("产品维护导出失败： %s", err.Error()))
		return
	}
	title := []map[string]string{
		{"skuCode": "产品SKU"},
		{"goodsName": "商品名称"},
		{"brandZh": "品牌"},
		{"mfgModel": "型号"},
		{"companyName": "公司名称"},
		{"warehouseName": "仓库名称"},
		{"vendorName": "货主名称"},
		{"supplierSkuCode": "货主SKU"},
		{"marketPrice": "销售价"},
		{"priceModifyReason": "价格调整备注"},
		{"productNo": "物料编码"},
		{"onlineStatus": "上下架状态"},
		{"status": "启用状态"},
	}
	excelApp := excel.NewExcel()

	err = excelApp.ExportExcelByStruct(c, title, exportData, "商品管理", "goods")
	if err != nil {
		e.Error(500, err, fmt.Sprintf("产品维护导出失败： %s", err.Error()))
	}
}

/** ------------------------------INNER部分--------------------------------------**/

// GetGoodsInfo 获取商品信息
// @Summary 获取商品信息
// @Description 获取商品信息
// @Tags 商品管理表
// @Param data body dto.GoodsDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /inner/pc/info [post]
// @Security Bearer
func (e Goods) GetGoodsInfo(c *gin.Context) {
	s := service.Goods{}
	req := dto.GoodsInfoReq{}
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
	var data []models.Goods
	err = s.GetGoodsInfo(&req, &data)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取商品管理表失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(data, "成功")
}

// GetGoodsById ID获取商品信息
// @Summary ID获取商品信息
// @Description ID获取商品信息
// @Tags 商品管理表
// @Param data body dto.GetGoodsByIdReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /inner/pc/get-by-id [post]
// @Security Bearer
func (e Goods) GetGoodsById(c *gin.Context) {
	s := service.Goods{}
	req := dto.GetGoodsByIdReq{}
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
	var data []models.Goods
	err = s.GetGoodsById(&req, &data)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取商品管理表失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(data, "成功")
}

// GetGoodsBySkuCodeReq 仓库Code + sku查商品信息
// @Summary 仓库Code + sku查商品信息
// @Description 仓库Code + sku查商品信息
// @Tags 商品管理表
// @Param data body dto.GoodsDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /inner/pc/get-by-skucode [post]
// @Security Bearer
func (e Goods) GetGoodsBySkuCodeReq(c *gin.Context) {
	s := service.Goods{}
	req := dto.GetGoodsBySkuCodeReq{}
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
	var data []models.Goods
	err = s.GetGoodsBySkuCodeReq(&req, &data)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取商品管理表失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(data, "成功")
}
