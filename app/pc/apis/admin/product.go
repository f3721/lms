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
	"strings"
)

type Product struct {
	api.Api
}

// GetPage 获取商品档案列表
// @Summary 获取商品档案列表
// @Description 获取商品档案列表
// @Tags 商品档案
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Param productId query int false "产品ID"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.Product}} "{"code": 200, "data": [...]}"
// @Router /api/v1/pc/product [get]
// @Security Bearer
func (e Product) GetPage(c *gin.Context) {
	req := dto.ProductGetPageReq{}
	s := service.Product{}
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
	list := make([]models.Product, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取商品档案失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取商品档案
// @Summary 获取商品档案
// @Description 获取商品档案
// @Tags 商品档案
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.Product} "{"code": 200, "data": [...]}"
// @Router /api/v1/pc/product/{id} [get]
// @Security Bearer
func (e Product) Get(c *gin.Context) {
	req := dto.ProductGetReq{}
	s := service.Product{}
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
	var object dto.GetInfoResp

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取商品档案失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建商品档案
// @Summary 创建商品档案
// @Description 创建商品档案
// @Tags 商品档案
// @Accept application/json
// @Product application/json
// @Param data body dto.ProductInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/pc/product [post]
// @Security Bearer
func (e Product) Insert(c *gin.Context) {
	req := dto.ProductInsertReq{}
	s := service.Product{}
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
		e.Error(500, err, fmt.Sprintf("创建商品档案失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改商品档案
// @Summary 修改商品档案
// @Description 修改商品档案
// @Tags 商品档案
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.ProductUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/pc/product/{id} [put]
// @Security Bearer
func (e Product) Update(c *gin.Context) {
	req := dto.ProductUpdateReq{}
	s := service.Product{}
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
		e.Error(500, err, fmt.Sprintf("修改商品档案失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除商品档案
// @Summary 删除商品档案
// @Description 删除商品档案
// @Tags 商品档案
// @Param data body dto.ProductDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/pc/product [delete]
// @Security Bearer
func (e Product) Delete(c *gin.Context) {
	s := service.Product{}
	req := dto.ProductDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除商品档案失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}

// BatchProApproval 批量审核图文
// @Summary 批量审核图文
// @Description 批量审核图文
// @Tags 商品档案
// @Param data body dto.BatchProApprovalReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "批量审核图文成功"}"
// @Router /api/v1/pc/product/batch-pro-approval [put]
// @Security Bearer
func (e Product) BatchProApproval(c *gin.Context) {
	s := service.Product{}
	req := dto.BatchProApprovalReq{}
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

	err = s.BatchProApproval(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("批量审核图文失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK("", "批量审核图文成功")
}

// ProApproval 单条图文审核
// @Summary 单条图文审核
// @Description 单条图文审核
// @Tags 单条图文审核
// @Param data body dto.ProductDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/pc/product/pro-approval/{id} [put]
// @Security Bearer
func (e Product) ProApproval(c *gin.Context) {
	s := service.Product{}
	req := dto.ProApprovalReq{}
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

	err = s.ProApproval(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("图文审核操作失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(12, "图文审核操作成功")
}

// ProductImportAdd 产品新增导入
// @Summary 产品新增导入
// @Description 产品新增导入
// @Tags 产品新增导入
// @Param data body dto.ExportReq true "body"
// @Router /api/v1/pc/product/product-import-add [post]
// @Security Bearer
func (e Product) ProductImportAdd(c *gin.Context) {
	s := service.Product{}
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
	err, title, exportData := s.ImportInsert(file, c)
	if err != nil {
		if len(exportData) > 0 {
			excelApp := excel.NewExcel()
			err = excelApp.ExportExcelByMap(c, title, exportData, "产品新增导入有错误", "Sheet1")
			if err != nil {
				e.Error(500, err, fmt.Sprintf("产品新增导入失败： %s", err.Error()))
			}
		} else {
			e.Error(500, err, err.Error())
		}
	} else {
		e.OK(200, "导入成功")
	}
}

// ProductImportUpdate 产品维护导入
// @Summary 产品维护导入
// @Description 产品维护导入
// @Tags 产品维护导入
// @Param data body dto.ExportReq true "body"
// @Router /api/v1/pc/product/product-import-update [post]
// @Security Bearer
func (e Product) ProductImportUpdate(c *gin.Context) {
	s := service.Product{}
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
	err, title, exportData := s.ImportUpdate(file, c)
	if err != nil {
		if len(exportData) > 0 {
			excelApp := excel.NewExcel()
			err = excelApp.ExportExcelByMap(c, title, exportData, "产品维护导入有错误", "Sheet1")
			if err != nil {
				e.Error(500, err, fmt.Sprintf("产品维护导入失败： %s", err.Error()))
			}
		} else {
			e.Error(500, err, err.Error())
		}
	} else {
		e.OK(200, "导入成功")
	}
}

// ProductExport 产品维护导出
// @Summary 产品维护导出
// @Description 产品维护导出
// @Tags 产品维护导出
// @Param data body dto.ProductGetPageReq true "body"
// @Router /api/v1/pc/product/product-export [get]
// @Security Bearer
func (e Product) ProductExport(c *gin.Context) {
	s := service.Product{}
	req := dto.ProductGetPageReq{}
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
	exportData, err := s.ProductExport(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("产品维护导出失败： %s", err.Error()))
		return
	}
	title := []map[string]string{
		{"skuCode": "产品SKU"},
		{"nameZh": "产品名称(中文)"},
		{"nameEn": "产品名称(英文)"},
		{"brandZh": "品牌(中文)"},
		{"brandEn": "品牌(英文)"},
		{"vendorName": "货主名称"},
		{"supplierSkuCode": "货主SKU"},
		{"level1CatName": "一级目录"},
		{"level2CatName": "二级目录"},
		{"level3CatName": "三级目录"},
		{"level4CatName": "四级目录"},
		{"mfgModel": "制造厂型号"},
		{"physicalUom": "物理单位"},
		{"salesUom": "售卖包装单位"},
		{"salesPhysicalFactor": "销售包装单位中含物理单位数量"},
		{"packWeight": "重量(kg)"},
		{"packLength": "售卖包装长(mm)"},
		{"packWidth": "售卖包装宽(mm)"},
		{"packHeight": "售卖包装高(mm)"},
		{"fragileFlag": "易碎标志(0:否;1:是)"},
		{"hazardFlag": "危险品标志(0:否;1:是)"},
		{"hazardClass": "危险品等级"},
		{"customMadeFlag": "定制品标志(0:否,1:是)"},
		{"bulkyFlag": "抛货标志(0:否; 1:是)"},
		{"assembleFlag": "拼装件标志(0:否; 1:是)"},
		{"isValuables": "是否贵重品(0:否; 1:是)"},
		{"isFluid": "是否液体(0:否; 1:是)"},
		{"consumptiveFlag": "耗材标志(0:否; 1:是)"},
		{"storageFlag": "保存期标志(0:否; 1:是)"},
		{"storageTime": "保存期限(月)"},
		{"refundFlag": "可退换货标志(0:否; 1:是)"},
	}
	excelApp := excel.NewExcel()
	err = excelApp.ExportExcelByStruct(c, title, exportData, "product", "product")
	if err != nil {
		e.Error(500, err, fmt.Sprintf("产品维护导出失败： %s", err.Error()))
	}
}

// AttributeExport 产品属性维护导出
// @Summary 产品属性维护导出
// @Description 产品属性维护导出
// @Tags 产品属性维护导出
// @Param data body dto.ExportReq true "body"
// @Router /api/v1/pc/product/product-attribute-export [get]
// @Security Bearer
func (e Product) AttributeExport(c *gin.Context) {
	s := service.Product{}
	req := dto.ProductGetPageReq{}
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

	title, exportData, err := s.AttributeExport(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("产品属性维护导出失败： %s", err.Error()))
		return
	}
	excelApp := excel.NewExcel()
	excelApp.ExportExcelByMap(c, title, exportData, "product-attribute", "Sheet1")
}

// ProductAttributeImport 产品属性维护导入
// @Summary 产品属性维护导入
// @Description 产品属性维护导入
// @Tags 产品属性维护导入
// @Router /api/v1/pc/product/product-attribute-import [post]
// @Security Bearer
func (e Product) ProductAttributeImport(c *gin.Context) {
	s := service.Product{}
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
	err, title, exportData := s.AttributeImportUpdate(file)
	if err != nil {
		if len(exportData) > 0 {
			excelApp := excel.NewExcel()
			err = excelApp.ExportExcelByMap(c, title, exportData, "产品属性维护导入有错误", "Sheet1")
			if err != nil {
				e.Error(500, err, fmt.Sprintf("产品属性维护导入失败： %s", err.Error()))
			}
		} else {
			e.Error(500, err, err.Error())
		}
	} else {
		e.OK(200, "导入成功")
	}
}

// GetProductAttribute 获取产品属性
// @Summary 获取产品属性
// @Description 获取产品属性
// @Tags 获取产品属性
// @Accept application/json
// @Param data body dto.GetProductExtAttributeReq true "data"
// @Router /api/v1/pc/product/product-attribute [post]
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Security Bearer
func (e Product) GetProductAttribute(c *gin.Context) {
	req := dto.GetProductExtAttributeReq{}
	s := service.ProductExtAttribute{}
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
	data, err := s.GetAttrs(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取商品档案失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(data, "查询成功")
}

// BatchUploadProductImage 批量上传产品图片
// @Summary 批量上传产品图片
// @Description 批量上传产品图片
// @Tags 批量上传产品图片
// @Accept application/json
// @Param data body dto.BatchUploadProductImageReq true "data"
// @Router /api/v1/pc/product/batch-upload-product-image [post]
// @Security Bearer
func (e Product) BatchUploadProductImage(c *gin.Context) {
	req := dto.BatchUploadProductImageReq{}
	s := service.Product{}
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
	errMsg, _ := s.BatchUploadProductImage(&req)
	if len(errMsg) > 0 {
		e.Error(500, err, fmt.Sprintf("批量上传图片失败，%s", strings.Join(errMsg, ";")))
		return
	}
	e.OK("", "上传成功")
}

// GetProductBySkuCode SKU查询产品档案信息
// @Summary SKU查询产品档案信息
// @Description SKU查询产品档案信息
// @Tags SKU查询产品档案信息
// @Param data body dto.GetProductBySkuCodeReq true "body"
// @Router /api/v1/pc/admin/product/skucode [get]
// @Security Bearer
func (e Product) GetProductBySkuCode(c *gin.Context) {
	s := service.Product{}
	req := dto.GetProductBySkuCodeReq{}
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
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取商品管理表失败，\r\n失败信息 %s", err.Error()))
		return
	}
	var data dto.GetProductBySkuCodeResp
	err = s.GetBySkuCode(&req, &data)
	e.OK(data, "成功")
}

// Sort 排序
// @Summary 排序
// @Description 排序
// @Tags 排序
// @Accept application/json
// @Product application/json
// @Param data body dto.ProductSortReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/pc/product/sort [post]
// @Security Bearer
func (e Product) Sort(c *gin.Context) {
	req := dto.ProductSortReq{}
	s := service.Product{}
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

	err = s.Sort(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("排序失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK("", "修改成功")
}

/** ------------------------------INNER部分--------------------------------------**/

// GetProductBySku SKU查询产品档案信息
// @Summary SKU查询产品档案信息
// @Description SKU查询产品档案信息
// @Tags SKU查询产品档案信息
// @Param data body dto.ExportReq true "body"
// @Router /inner/pc/admin/product/get-by-sku [post]
// @Security Bearer
func (e Product) GetProductBySku(c *gin.Context) {
	s := service.Product{}
	req := dto.InnerGetProductBySkuReq{}
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
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取商品管理表失败，\r\n失败信息 %s", err.Error()))
		return
	}
	var data []dto.InnerGetProductBySkuResp
	err = s.GetProductBySku(&req, &data)
	e.OK(data, "成功")
}

// GetProductCategoryBySku SKU查询产线
// @Summary SKU查询产线
// @Description SKU查询产线
// @Tags SKU查询产线
// @Param data body dto.ExportReq true "body"
// @Router /inner/pc/admin/product/get-category [post]
// @Security Bearer
func (e Product) GetProductCategoryBySku(c *gin.Context) {
	s := service.Product{}
	req := dto.InnerGetProductBySkuReq{}
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
	if err != nil {
		e.Error(500, err, fmt.Sprintf("SKU查询产线，\r\n失败信息 %s", err.Error()))
		return
	}
	var data []dto.InnerGetProductCategoryBySkuResp
	s.GetProductCategoryBySku(&req, &data)
	e.OK(data, "成功")
}

// -----------------------------------------------------------OPC商品同步--------------------------------------------------------------------

func (e Product) SupplierProductSync(c *gin.Context) {
	s := service.ProductSync{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	if err != nil {
		e.Error(500, err, fmt.Sprintf("产品同步失败，\r\n失败信息 %s", err.Error()))
		return
	}
	err = s.SupplierProductSync(c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	e.OK([]any{}, "成功")
}
