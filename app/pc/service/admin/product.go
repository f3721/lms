package admin

import (
	"encoding/json"
	"errors"
	"fmt"
	vd "github.com/bytedance/go-tagexpr/v2/validator"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"github.com/jinzhu/copier"
	"github.com/samber/lo"
	dtoAdmin "go-admin/app/admin/service/dto"
	modelsWc "go-admin/app/wc/models"
	dtoWc "go-admin/app/wc/service/admin/dto"
	adminClient "go-admin/common/client/admin"
	wcClient "go-admin/common/client/wc"
	"go-admin/common/excel"
	"go-admin/common/global"
	cModels "go-admin/common/models"
	"go-admin/common/utils"
	"mime/multipart"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/pc/models"
	"go-admin/app/pc/service/admin/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type Product struct {
	service.Service
}

const SKU_PRE = "product_sku_pre"
const DefaultTax = "0.13"

var dictMap map[int]string

// GetPage 获取Product列表
func (e *Product) GetPage(c *dto.ProductGetPageReq, p *actions.DataPermission, list *[]models.Product, count *int64) error {
	var err error
	var data models.Product

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
			dto.MakeSearchCondition(c, e.Orm),
			actions.SysUserPermission(data.TableName(), p, 4), //货主权限
		).
		Joins("INNER JOIN brand ON product.brand_id = brand.id").
		Preload("Brand").
		Preload("MediaRelationship", func(db *gorm.DB) *gorm.DB {
			return db.Order("media_relationship.seq ASC")
		}).
		Preload("MediaRelationship.MediaInstant").
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("ProductService GetPage error:%s \r\n", err)
		return err
	}

	if *count > 0 {
		// 取出所有的货主ID
		var vendorIds []int
		_ = utils.StructColumn(&vendorIds, *list, "VendorId", "")
		if len(vendorIds) > 0 {
			// 货主map
			vendorResult := e.apiGetVendorInfoById(vendorIds)
			tmpList := *list
			for k, product := range tmpList {
				tmpList[k].VendorName = vendorResult[product.VendorId]
			}
			list = &tmpList
		}
	}

	return nil
}

// Get 获取Product对象
func (e *Product) Get(d *dto.ProductGetReq, p *actions.DataPermission, model *dto.GetInfoResp) error {
	var data models.Product
	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		Preload("Brand").
		Preload("MediaRelationship", func(db *gorm.DB) *gorm.DB {
			return db.Order("media_relationship.seq ASC")
		}).
		Preload("MediaRelationship.MediaInstant").
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetProduct error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	// 商品分类
	productCategoryData := make([]models.ProductCategory, 0)
	ProductCategoryService := ProductCategory{e.Service}
	category := Category{e.Service}
	_ = ProductCategoryService.GetCategoryBySku([]string{model.SkuCode}, &productCategoryData)
	if len(productCategoryData) > 0 {
		categoryPaths := make([]dto.CategoryPath, 0)
		for _, productCategory := range productCategoryData {
			categoryPath := category.GetCategoryPath(productCategory.CategoryId)
			var categoryName []string
			for _, path := range categoryPath {
				categoryName = append(categoryName, fmt.Sprintf("%s[%d]", path.NameZh, path.Id))
			}
			catePath := dto.CategoryPath{
				Id:         productCategory.Id,
				CategoryId: productCategory.CategoryId,
				Name:       strings.Join(categoryName, " > "),
				Value:      productCategory.CategoryId,
				PathName:   strings.Join(categoryName, " > "),
			}
			categoryPaths = append(categoryPaths, catePath)
			if productCategory.MainCateFlag == 1 {
				model.ProductCategoryPrimary = productCategory.CategoryId
			}
		}
		model.ProductCategory = categoryPaths
	}

	if model.ProductCategoryPrimary > 0 {
		productExtAttribute := ProductExtAttribute{e.Service}
		arrList, err := productExtAttribute.GetAttrs(&dto.GetProductExtAttributeReq{
			CategoryId: model.ProductCategoryPrimary,
			SkuCode:    model.SkuCode,
		})
		if err == nil && len(arrList) > 0 {
			model.AttrList = &arrList
		}
	}

	// 货主map
	vendorResult := e.apiGetVendorInfoById([]int{model.VendorId})
	model.VendorName = vendorResult[model.VendorId]

	return nil
}

// Insert 创建Product对象
func (e *Product) Insert(c *dto.ProductInsertReq) error {
	// 自动生成SKU
	skuCode := e.createSku()
	c.SkuCode = skuCode

	var err error
	var data models.Product
	c.Generate(&data)
	//数据插入校验
	err = e.validateForm(c, "add")
	if err != nil {
		return err
	}

	tx := e.Orm.Debug().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	err = tx.Create(&data).Error
	if err != nil {
		tx.Rollback()
		e.Log.Errorf("ProductService Insert error:%s \r\n", err)
		return err
	}

	// 判断是否上传了图片
	if len(c.ProductImage) > 0 {
		err = e.createProductImage(tx, c.ProductImage, &data)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// 创建商品分类
	err = e.createProductCategory(tx, c.ProductCategory, data.SkuCode, c.ProductCategoryPrimary)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 保存商品扩展属性
	_ = e.SaveProductExtAttribute(c)

	// 生成日志
	dataLog, _ := json.Marshal(&c)
	productLog := models.ProductLog{
		DataId:     data.Id,
		Type:       global.LogTypeCreate,
		Data:       string(dataLog),
		BeforeData: "",
		AfterData:  string(dataLog),
		ControlBy:  c.ControlBy,
	}
	_ = productLog.CreateLog("product", e.Orm)

	return nil
}

// ImportInsert 产品批量新增
func (e *Product) ImportInsert(file *multipart.FileHeader, c *gin.Context) (err error, errTitleList []map[string]string, errData []map[string]interface{}) {
	excelApp := excel.NewExcel()
	err, datas, title := excelApp.GetExcelData(file)
	if err != nil {
		return
	}
	var importData []dto.ProductImportData
	for _, m := range datas {
		importDto := dto.ProductImportData{}
		for k, v := range m {
			importDto.Generate(k, v)
			importDto.CreateBy = user.GetUserId(c)
			importDto.CreateByName = user.GetUserName(c)
		}
		importData = append(importData, importDto)
	}
	if len(importData) > 0 {
		// 取出所有货主名称调接口查
		var vendorName []string
		utils.StructColumn(&vendorName, importData, "VendorName", "")
		apiResult := e.apiGetVendorInfo(vendorName)
		errMsg := make(map[int]string, len(importData))
		successCount := 0
		for i, data := range importData {
			errStr := ""
			// 数据基础校验
			err = vd.Validate(data)
			if err != nil {
				errStr = err.Error()
			} else {
				vendorNames := vendorList(apiResult)
				if !utils.InArrayString(data.VendorName, vendorNames) {
					errStr = "货主名称不存在！"
				} else {
					err = e.validateImport(&data, "add")
					if err != nil {
						errStr = err.Error()
					} else {
						// 处理数据
						insertReq := dto.ProductInsertReq{}
						err = copier.Copy(&insertReq, &data)
						if err != nil {
							errStr = err.Error()
						} else {
							insertReq.Status = 1 // 待审核
							insertReq.VendorId = findVendorId(apiResult, data.VendorName)
							err = e.Insert(&insertReq)
							if err != nil {
								errStr = err.Error()
							}
						}
					}
				}
			}
			if errStr == "" {
				successCount += 1
			}
			errMsg[i] = errStr
		}

		if len(importData) > successCount {
			err = errors.New("有错误！")
			errTitleList, errData = excelApp.MergeErrMsgColumn(title, datas, errMsg)
		}
		return
	}
	return
}

// ImportUpdate 产品批量维护导入
func (e *Product) ImportUpdate(file *multipart.FileHeader, c *gin.Context) (err error, errTitleList []map[string]string, errData []map[string]interface{}) {
	excelApp := excel.NewExcel()
	err, datas, title := excelApp.GetExcelData(file)
	if err != nil {
		return
	}

	var importData []dto.ProductImportUpdateReq

	// 导入数据map转struct
	for _, m := range datas {
		importDto := dto.ProductImportUpdateReq{}
		for k, v := range m {
			importDto.MapToStruct(k, v)
		}
		importDto.UpdateBy = user.GetUserId(c)
		importDto.UpdateByName = user.GetUserName(c)
		importData = append(importData, importDto)
	}

	if len(importData) == 0 {
		err = errors.New("导入数据不能为空！")
		return
	}

	// api获取字典危险品等级列表
	dictMap = e.apiGetDictList("storage_time")
	// 数据校验
	errMsg := make(map[int]string, 0)
	successCount := 0
	for i, product := range importData {
		errStr := ""
		// 基础数据校验
		err = vd.Validate(product)
		if err != nil {
			errStr = err.Error()
		} else {
			// 导入校验
			productImportData := dto.ProductImportData{}
			_ = copier.Copy(&productImportData, &product)
			err = e.validateImport(&productImportData, "edit")
			if err != nil {
				errStr = err.Error()
			} else {
				// 数据保存校验
				var productInsertReq dto.ProductInsertReq
				_ = copier.Copy(&productInsertReq, &productImportData)
				err = e.validateForm(&productInsertReq, "edit")
				if err != nil {
					errStr = err.Error()
				} else {
					// 数据保存
					var data dto.ProductImportUpdateReq
					_ = copier.Copy(&data, &productInsertReq)
					err = e.importUpdate(&data)
					if err != nil {
						errStr = err.Error()
					}
				}
			}
		}
		if errStr == "" {
			successCount += 1
		}
		errMsg[i] = errStr
	}

	// 有错误返回
	if len(importData) > successCount {
		err = errors.New("有错误！")
		errTitleList, errData = excelApp.MergeErrMsgColumn(title, datas, errMsg)
		return
	}
	return
}

// AttributeImportUpdate 属性维护导入
func (e *Product) AttributeImportUpdate(file *multipart.FileHeader) (err error, errTitleList []map[string]string, errData []map[string]interface{}) {
	excelApp := excel.NewExcel()
	err, datas, title := excelApp.GetExcelData(file)
	if err != nil {
		return
	}
	errMsg := make(map[int]string, len(datas))
	successCount := 0
	for i, m := range datas {
		errStr := ""
		importDto := dto.ProductAttributeImportReq{}
		for k, v := range m {
			importDto.Generate(k, v)
		}
		// 数据基础校验
		err = vd.Validate(importDto)
		if err != nil {
			errStr = err.Error()
		} else {
			err = e.productAttributeImportValidate(&importDto)
			if err != nil {
				errStr = err.Error()
			} else {
				// 属性数据组装
				attrData := make([]*dto.ProductAttributeImportTemp, 0)
				result := make([]dto.AttrsKeyName, 0)
				categoryAttribute := CategoryAttribute{e.Service}
				err = categoryAttribute.GetAttrsByCategory([]int{importDto.CategoryId}, &result)
				if err != nil {
					errStr = err.Error()
				} else {
					if len(result) > 0 {
						for _, val := range result {
							if v, ok := m["attribute_"+strconv.Itoa(val.KeyName)].(string); ok {
								tmp := dto.ProductAttributeImportTemp{
									Key:            i,
									CategoryId:     importDto.CategoryId,
									SkuCode:        importDto.SkuCode,
									AttributeId:    strconv.Itoa(val.KeyName),
									AttributeValue: v,
								}
								attrData = append(attrData, &tmp)
							}
						}
					} else {
						errStr = "产线ID[终极目录]填写有误"
					}
				}

				// 保存数据
				if len(attrData) > 0 {
					productExtAttribute := ProductExtAttribute{e.Service}
					var errorMsg []string
					for _, productAttributeImportTemp := range attrData {
						err = productExtAttribute.Import(productAttributeImportTemp)
						if err != nil {
							errorMsg = append(errorMsg, err.Error())
						}
					}
					if len(errorMsg) > 0 {
						errStr = strings.Join(errorMsg, ";")
					}
				}
			}
		}
		if errStr == "" {
			successCount += 1
		}
		errMsg[i] = errStr
	}
	if len(datas) > successCount {
		err = errors.New("有错误！")
		errTitleList, errData = excelApp.MergeErrMsgColumn(title, datas, errMsg)
	}
	return
}

// ProductExport 产品维护导出
func (e *Product) ProductExport(c *dto.ProductGetPageReq, p *actions.DataPermission) (exportData []interface{}, err error) {
	var data []models.Product
	var model models.Product
	err = e.Orm.Model(&model).Scopes(
		cDto.MakeCondition(c.GetNeedSearch()),
		actions.Permission(model.TableName(), p),
		dto.MakeSearchCondition(c, e.Orm),
		actions.SysUserPermission(model.TableName(), p, 4), //货主权限
	).Preload("Brand").Find(&data).Error

	if err != nil {
		return
	}
	if len(data) > 0 {
		var skuCodes []string
		err = utils.StructColumn(&skuCodes, data, "SkuCode", "")
		if err != nil {
			return
		}
		categoryService := Category{e.Service}
		skuMap, result := categoryService.GetCategoryBySku(skuCodes)

		// api调用货主名称集合
		var vendorIds []int
		_ = utils.StructColumn(&vendorIds, data, "VendorId", "")
		vendorResult := e.apiGetVendorInfoById(vendorIds)

		// api获取字典危险品等级列表
		dictMap = e.apiGetDictList("hazard_class")

		for _, product := range data {
			if skuMap[product.SkuCode] != 0 {
				productResp := dto.ProductExportResp{}
				copier.Copy(&productResp, &product)
				productResp.BrandZh = product.Brand.BrandZh
				productResp.BrandEn = product.Brand.BrandEn
				productResp.Level1CatName = result[skuMap[product.SkuCode]].NameZh1
				productResp.Level2CatName = result[skuMap[product.SkuCode]].NameZh2
				productResp.Level3CatName = result[skuMap[product.SkuCode]].NameZh3
				productResp.Level4CatName = result[skuMap[product.SkuCode]].NameZh4
				productResp.VendorName = vendorResult[product.VendorId]
				if product.HazardClass == 0 {
					productResp.HazardClass = ""
				} else {
					productResp.HazardClass = dictMap[product.HazardClass]
				}
				exportData = append(exportData, productResp)
			}
		}
	}
	return
}

// AttributeExport 属性维护导出
func (e *Product) AttributeExport(c *dto.ProductGetPageReq) ([]map[string]string, []map[string]interface{}, error) {
	var err error
	headerMap := make([]map[string]string, 0)
	exportData := make([]map[string]interface{}, 0)

	lastId := 0
	if c.Level4Catid != 0 {
		lastId = c.Level4Catid
	} else {
		lastId = c.Level3Catid
	}

	if lastId > 0 {
		result := make([]dto.AttrsKeyName, 0)
		categoryAttribute := CategoryAttribute{e.Service}
		err = categoryAttribute.GetAttrsByCategory([]int{lastId}, &result)
		if err != nil {
			return nil, nil, err
		}
		var headerAttrs []map[string]string
		for _, val := range result {
			headerAttrs = append(headerAttrs, map[string]string{
				"attribute_" + strconv.Itoa(val.KeyName): val.Name,
			})
		}
		// 分类数据
		categoryData := map[string]interface{}{
			"categoryId": lastId,
		}
		category := Category{e.Service}
		categoryPath := category.GetCategoryPath(lastId)
		if len(categoryPath) > 0 {
			for i, category := range categoryPath {
				categoryData["categoryLevel"+strconv.Itoa(i+1)] = category.NameZh
			}
		}
		// 产品公有属性
		headerMap = []map[string]string{
			{"categoryLevel1": "一级目录"},
			{"categoryLevel2": "二级目录"},
			{"categoryLevel3": "三级目录"},
			{"categoryLevel4": "四级目录"},
			{"categoryId": "产线ID[终极目录]"},
			{"skuCode": "产品SKU"},
			{"nameZh": "产品名称(中文)"},
			{"mfgModel": "制造商型号"},
		}
		// 属性合并
		if len(headerAttrs) > 0 {
			for _, attr := range headerAttrs {
				for k, v := range attr {
					headerMap = append(headerMap, map[string]string{k: v})
				}
			}
		}
		list := make([]models.Product, 0)
		var count int64
		_ = e.getTotalProducts(c, &list, &count)
		productExportAttrResp := make([]dto.ProductExportAttrResp, 0)
		err = e.getProductsForExportAttr(&list, &productExportAttrResp)
		if err != nil {
			return nil, nil, err
		}
		if count > 0 {
			tmp := map[string]map[string]interface{}{}
			for _, resp := range productExportAttrResp {
				data := tmp[resp.SkuCode]
				if data == nil {
					data = make(map[string]interface{})
				}
				if _, exists := data["skuCode"]; !exists {
					data["skuCode"] = resp.SkuCode
				}
				if _, exists := data["NameZh"]; !exists {
					data["nameZh"] = resp.NameZh
				}
				if _, exists := data["MfgModel"]; !exists {
					data["mfgModel"] = resp.MfgModel
				}
				if _, exists := data["attribute_"+strconv.Itoa(resp.AttributeId)]; !exists {
					data["attribute_"+strconv.Itoa(resp.AttributeId)] = resp.AttributeValue
				}
				tmp[resp.SkuCode] = data
			}
			for _, val := range tmp {
				var element []map[string]interface{}
				tmps := make(map[string]interface{})
				for _, m := range headerMap {
					for k, _ := range m {
						value, exist := val[k]
						if !exist {
							value = ""
						}
						tmps[k] = value
					}
				}
				element = append(element, tmps)
				for i, m := range element {
					tmps = make(map[string]interface{}, 0)
					for mk, v := range m {
						value, exist := categoryData[mk]
						if !exist {
							value = v
						}
						tmps[mk] = value
					}
					element[i] = tmps
				}
				for _, e := range element {
					exportData = append(exportData, e)
				}
			}
		}
	}
	return headerMap, exportData, nil
}

// Update 修改Product对象
func (e *Product) Update(c *dto.ProductUpdateReq, p *actions.DataPermission) error {
	var data = models.Product{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())

	var form dto.ProductInsertReq
	copier.Copy(&form, &c)
	err := e.validateForm(&form, "edit")
	if err != nil {
		return err
	}

	err = e._update(c, data)
	if err != nil {
		return err
	}
	return nil
}

func (e *Product) apiGetVendorInfo(vendorName []string) map[int]string {
	vendorResult := wcClient.ApiByDbContext(e.Orm).GetVendorList(dtoWc.InnerVendorsGetListReq{
		NameZh: strings.Join(lo.Uniq(vendorName), ","),
	})
	vendorResultInfo := &struct {
		response.Response
		Data []modelsWc.Vendors
	}{}
	vendorResult.Scan(vendorResultInfo)
	vendorMap := make(map[int]string, len(vendorResultInfo.Data))
	for _, vendor := range vendorResultInfo.Data {
		vendorMap[vendor.Id] = vendor.NameZh
	}
	return vendorMap
}

func vendorList(vendors map[int]string) []string {
	vendorNames := make([]string, 0)
	for _, v := range vendors {
		vendorNames = append(vendorNames, v)
	}
	return vendorNames
}

func findVendorId(vendors map[int]string, vendorName string) int {
	vendorId := 0
	for i, s := range vendors {
		if s == vendorName {
			vendorId = i
		}
	}
	return vendorId
}

func (e *Product) apiGetVendorInfoById(vendorIds []int) map[int]string {
	// 货主map
	vendorResult := wcClient.ApiByDbContext(e.Orm).GetVendorList(dtoWc.InnerVendorsGetListReq{
		Ids: strings.Trim(strings.Join(strings.Fields(fmt.Sprint(vendorIds)), ","), "[]"),
	})
	vendorResultInfo := &struct {
		response.Response
		Data []modelsWc.Vendors
	}{}
	vendorResult.Scan(vendorResultInfo)
	vendorMap := make(map[int]string, len(vendorResultInfo.Data))
	for _, vendor := range vendorResultInfo.Data {
		vendorMap[vendor.Id] = vendor.NameZh
	}
	return vendorMap
}

func (e *Product) apiGetDictList(dictType string) map[int]string {
	dictResult := adminClient.ApiByDbContext(e.Orm).GetDictListByDictType(dictType)
	dictResultInfo := &struct {
		response.Response
		Data []dtoAdmin.SysDictDataGetAllResp
	}{}
	dictResult.Scan(dictResultInfo)
	dictMap := make(map[int]string, len(dictResultInfo.Data))
	for _, dict := range dictResultInfo.Data {
		key, _ := strconv.Atoi(dict.DictValue)
		dictMap[key] = dict.DictLabel
	}
	return dictMap
}

// GetBySkuCode 获取Product对象
func (e *Product) GetBySkuCode(c *dto.GetProductBySkuCodeReq, model *dto.GetProductBySkuCodeResp) error {
	var data models.Product
	err := e.Orm.Model(&data).
		Where("sku_code = ?", c.SkuCode).
		//Where("status = 2").
		First(model).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetProduct error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}

	if model.VendorId != 0 {
		// 货主map
		vendorResult := e.apiGetVendorInfoById([]int{model.VendorId})
		model.VendorName = vendorResult[model.VendorId]
	}

	return nil
}

func (e *Product) importUpdate(c *dto.ProductImportUpdateReq) error {
	var err error
	// 数据处理
	var product = models.Product{}
	err = e.Orm.Where("sku_code", c.SkuCode).First(&product).Error
	if err != nil {
		return err
	}

	tx := e.Orm.Debug().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	oldData := product
	c.Generate(&product)

	db := tx.Save(&product)
	if err = db.Error; err != nil {
		tx.Rollback()
		e.Log.Errorf("ProductService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}

	// 如果修改了货主SKU需要同步到GOODS表中
	if oldData.SupplierSkuCode != product.SupplierSkuCode {
		err := e.updateGoodsSupplerSkuCode(tx, product.SkuCode, product.VendorId, product.SupplierSkuCode)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// 删除商品分类
	err = e.removeProductCategory(tx, product.SkuCode)
	if err != nil {
		tx.Rollback()
		return errors.New("删除商品分类失败！")
	}

	// 创建商品分类
	err = e.createProductCategory(tx, c.ProductCategory, product.SkuCode, c.ProductCategoryPrimary)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 生成日志
	dataLog, _ := json.Marshal(&c)
	beforeDataStr, _ := json.Marshal(&oldData)
	afterDataStr, _ := json.Marshal(&product)
	productLog := models.ProductLog{
		DataId:     product.Id,
		Type:       global.LogTypeUpdate,
		Data:       string(dataLog),
		BeforeData: string(beforeDataStr),
		AfterData:  string(afterDataStr),
		ControlBy:  c.ControlBy,
	}
	_ = productLog.CreateLog("product", e.Orm)
	return nil
}

func (e *Product) _update(c *dto.ProductUpdateReq, data models.Product) error {
	var err error

	tx := e.Orm.Debug().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	oldData := data
	c.Generate(&data)

	db := tx.Save(&data)
	if err = db.Error; err != nil {
		tx.Rollback()
		e.Log.Errorf("ProductService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}

	// 如果修改了货主SKU需要同步到GOODS表中
	if oldData.SupplierSkuCode != data.SupplierSkuCode {
		err := e.updateGoodsSupplerSkuCode(tx, data.SkuCode, data.VendorId, data.SupplierSkuCode)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// 删除图片
	err = e.deleteMediaRelationship(tx, &data)
	if err != nil {
		tx.Rollback()
		return err
	}
	// 判断是否上传了图片
	if len(c.ProductImage) > 0 {
		err = e.createProductImage(tx, c.ProductImage, &data)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// 删除商品分类
	err = e.removeProductCategory(tx, data.SkuCode)
	if err != nil {
		tx.Rollback()
		return errors.New("删除商品分类失败！")
	}

	// 创建商品分类
	err = e.createProductCategory(tx, c.ProductCategory, data.SkuCode, c.ProductCategoryPrimary)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 保存商品扩展属性
	_ = e.UpdateProductExtAttribute(c)

	// 生成日志
	dataLog, _ := json.Marshal(&c)
	beforeDataStr, _ := json.Marshal(&oldData)
	afterDataStr, _ := json.Marshal(&data)
	productLog := models.ProductLog{
		DataId:     data.Id,
		Type:       global.LogTypeUpdate,
		Data:       string(dataLog),
		BeforeData: string(beforeDataStr),
		AfterData:  string(afterDataStr),
		ControlBy:  c.ControlBy,
	}
	_ = productLog.CreateLog("product", e.Orm)
	return nil
}

// BatchProApproval 批量审核图文
func (e *Product) BatchProApproval(d *dto.BatchProApprovalReq, p *actions.DataPermission) error {
	if len(d.Selected) <= 0 {
		return errors.New("请先勾选需要审核的行")
	}
	var model models.Product
	tx := e.Orm.Debug().Begin()
	for _, id := range d.Selected {
		var data = models.Product{}
		e.Orm.Model(&model).First(&data, id)
		if data.Status != 2 {
			oldData := data
			data.Status = 2
			err := tx.Save(&data).Error
			if err == nil {
				productUpdater := dto.ProductUpdater{
					ControlBy: cModels.ControlBy{
						UpdateBy:     d.UpdateBy,
						UpdateByName: d.UpdateByName,
					},
				}
				approveLog(productUpdater, d, &oldData, &data, tx)
			}
		}
	}
	tx.Commit()
	return nil
}

// ProApproval 单条图文审核
func (e *Product) ProApproval(d *dto.ProApprovalReq, p *actions.DataPermission) error {
	var data = models.Product{}
	err := e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, d.ProductId).Error

	if err != nil {
		return errors.New("数据不存在或无权查看！")
	}
	oldData := data
	data.Status = d.Status
	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("图文审核失败！")
	}
	productUpdater := dto.ProductUpdater{
		ControlBy: cModels.ControlBy{
			UpdateBy:     d.UpdateBy,
			UpdateByName: d.UpdateByName,
		},
	}
	approveLog(productUpdater, d, &oldData, &data, e.Orm)
	return nil
}

func approveLog(productUpdater dto.ProductUpdater, data any, beforeData *models.Product, afterData *models.Product, tx *gorm.DB) {
	dataLog, _ := json.Marshal(&data)
	beforeDataStr, _ := json.Marshal(beforeData)
	afterDataStr, _ := json.Marshal(afterData)
	productLog := models.ProductLog{
		DataId:     beforeData.Id,
		Type:       global.LogTypeUpdate,
		Data:       string(dataLog),
		BeforeData: string(beforeDataStr),
		AfterData:  string(afterDataStr),
		ControlBy: cModels.ControlBy{
			UpdateByName: productUpdater.UpdateByName,
			UpdateBy:     productUpdater.UpdateBy,
		},
	}
	_ = productLog.CreateLog("product", tx)
}

// Remove 删除Product
func (e *Product) Remove(d *dto.ProductDeleteReq, p *actions.DataPermission) error {
	var data models.Product

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveProduct error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}

// createProductImage 创建图片
func (e *Product) createProductImage(tx *gorm.DB, productImage []dto.MediaInstanceInsertReq, product *models.Product) error {
	mediaInstance := MediaInstance{service.Service{
		Orm: tx,
	}}
	var err error
	for _, productImg := range productImage {
		if productImg.MediaName != "" && productImg.MediaDir != "" {
			// 保存图片实例
			var media models.MediaInstance
			err = mediaInstance.Insert(&dto.MediaInstanceInsertReq{
				MediaName: productImg.MediaName,
				MediaDir:  productImg.MediaDir,
			}, &media)
			if err != nil {
				break
			}
			// 保存关联关系
			err = e.createMediaRelationship(tx, product, media.Id, productImg)
			if err != nil {
				break
			}
		}
	}
	return err
}

// Sort 批量排序
func (e *Product) Sort(c *dto.ProductSortReq, p *actions.DataPermission) error {
	var err error
	var data models.Product
	tx := e.Orm.Begin()
	for _, sort := range c.Sort {
		update := map[string]interface{}{
			"seq":            sort.Seq,
			"update_by":      c.UpdateBy,
			"update_by_name": c.UpdateByName,
		}
		tx.Model(&data).Where("id = ?", sort.ProductId).Updates(update)
	}
	err = tx.Commit().Error
	return err
}

// BatchUploadProductImage 批量上传产品图片
func (e *Product) BatchUploadProductImage(c *dto.BatchUploadProductImageReq) ([]string, error) {
	var errMsg []string
	tx := e.Orm
	mediaInstance := MediaInstance{service.Service{
		Orm: tx,
	}}

	for _, productImg := range c.ProductImage {
		// 图片名称截取SKU和序号 HLV007_1.jpg
		imageInfo := strings.Split(productImg.MediaName, ".")
		nameInfo := strings.Split(imageInfo[0], "_")

		// 判断sku_code是否存在
		var productInfo models.Product
		err := e.getBySkuCode(nameInfo[0], &productInfo)
		if err != nil || productInfo.SkuCode == "" {
			errMsg = append(errMsg, fmt.Sprintf("图片[%s]信息不存在", productImg.MediaName))
			break
		}
		// 保存图片实例
		var media models.MediaInstance
		err = mediaInstance.Insert(&dto.MediaInstanceInsertReq{
			MediaName: productImg.MediaName,
			MediaDir:  productImg.MediaDir,
		}, &media)
		if err != nil {
			errMsg = append(errMsg, fmt.Sprintf("图片[%s]保存失败", productImg.MediaName))
			break
		}
		if len(nameInfo) > 1 {
			seq, _ := strconv.Atoi(nameInfo[1])
			if seq > 0 {
				productImg.Seq = seq
			}
		}
		// 保存关联关系
		err = e.createMediaRelationship(tx, &productInfo, media.Id, productImg)
		if err != nil {
			errMsg = append(errMsg, fmt.Sprintf("图片[%s]保存失败", productImg.MediaName))
			break
		}
	}
	return errMsg, nil
}

// createProductCategory 创建商品分类
func (e *Product) createProductCategory(tx *gorm.DB, productCategorys []dto.ProductCategory, skuCode string, productCategoryPrimary int) error {
	productCategoryService := ProductCategory{service.Service{
		Orm: tx,
	}}
	var data []models.ProductCategory
	for _, productCategory := range productCategorys {
		mainCateFlag := 0
		if productCategoryPrimary == productCategory.CategoryId {
			mainCateFlag = 1
		}
		data = append(data, models.ProductCategory{
			SkuCode:      skuCode,
			CategoryId:   productCategory.CategoryId,
			MainCateFlag: mainCateFlag,
			ModelTime:    cModels.ModelTime{},
			ControlBy:    cModels.ControlBy{},
		})
	}
	err := productCategoryService.BatchInsert(&data)
	return err
}

// 删除商品分类
func (e *Product) removeProductCategory(tx *gorm.DB, skuCode string) error {
	productCategoryService := ProductCategory{service.Service{
		Orm: tx,
	}}
	err := productCategoryService.Remove(skuCode)
	return err
}

// 创建sku_code
func (e *Product) createSku() string {
	// 创建sku_code
	skuPre := utils.RandStr(3)
	var attributeConfig models.AttributeConfig
	err := e.findProductSkuPre(SKU_PRE, &attributeConfig)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err := e.createAttributeConfig(skuPre)
		if err != nil {
			e.findProductSkuPre(SKU_PRE, &attributeConfig)
		}
	}
	var product models.Product
	e.FindMaxSkuCode(attributeConfig.Value, &product)

	// 如果产品表无数据，取表里设置的默认前缀
	if product.SkuCode != "" {
		// 取出sku_code的数字部分，判断是否等于9999，
		num, _ := strconv.Atoi(product.SkuCode[3:])
		if num == 999 {
			newStr := skuPre
			// 判断新的前缀是否被使用
			err := e.FindMaxSkuCode(newStr, &product)
			if err != nil {
				e.createSku()
			}
			// 更新sku_code前缀
			attributeConfig.Key = newStr
			attributeConfig.Value = newStr
			e.SaveAttributeConfig(&attributeConfig)
			return newStr + utils.NumToString("0", 3, 1)
		}
		return product.SkuCode[0:3] + utils.NumToString("0", 3, num+1)
	} else {
		// 初始SKUCODE，取出默认前缀HLV + 四位数字 0001
		e.findProductSkuPre(SKU_PRE, &attributeConfig)
		return attributeConfig.Value + utils.NumToString("0", 3, 1)
	}
}

func (e *Product) findProductSkuPre(attributeType string, model *models.AttributeConfig) error {
	var data models.AttributeConfig
	err := e.Orm.Model(&data).Where("`type` = ?", attributeType).First(model).Error
	return err
}

func (e *Product) GetProductByBrandId(brandId int, model *models.Product) error {
	var data models.Product
	err := e.Orm.Model(&data).Where("`brand_id` = ?", brandId).Where("deleted_at IS NULL").First(model).Error
	return err
}

func (e *Product) FindMaxSkuCode(productSkuPre string, model *models.Product) error {
	var data models.Product
	err := e.Orm.Model(&data).Select("max(sku_code) as SkuCode").Where("`sku_code` like ?", "%"+productSkuPre+"%").First(model).Error
	return err
}

func (e *Product) createAttributeConfig(productSkuPre string) error {
	attributeConfig := AttributeConfig{e.Service}
	err := attributeConfig.Insert(&dto.AttributeConfigInsertReq{
		Type:      SKU_PRE,
		Key:       productSkuPre,
		Value:     productSkuPre,
		SortOrder: 0,
	})
	return err
}

func (e *Product) SaveAttributeConfig(model *models.AttributeConfig) error {
	attributeConfig := AttributeConfig{e.Service}
	err := attributeConfig.Update(&dto.AttributeConfigUpdateReq{
		Id:    model.Id,
		Type:  SKU_PRE,
		Key:   model.Key,
		Value: model.Value,
	})
	return err
}

// SaveProductExtAttribute 保存产品扩展属性
func (e *Product) SaveProductExtAttribute(c *dto.ProductInsertReq) error {
	if len(c.AttrList) > 0 {
		productExtAttribute := ProductExtAttribute{e.Service}
		for _, v := range c.AttrList {
			_ = productExtAttribute.Insert(&dto.ProductExtAttributeInsertReq{
				SkuCode:     c.SkuCode,
				AttributeId: v.KeyName,
				ValueZh:     v.Value,
				ValueEn:     "",
				Status:      1,
				ControlBy:   cModels.ControlBy{},
			})
		}
	}
	return nil
}

func (e *Product) UpdateProductExtAttribute(c *dto.ProductUpdateReq) error {
	if len(c.AttrList) > 0 {
		productExtAttribute := ProductExtAttribute{e.Service}
		err := productExtAttribute.Delete(c.SkuCode)
		if err != nil {
			return err
		}
		for _, v := range c.AttrList {
			_ = productExtAttribute.Insert(&dto.ProductExtAttributeInsertReq{
				SkuCode:     c.SkuCode,
				AttributeId: v.KeyName,
				ValueZh:     v.Value,
				ValueEn:     "",
				Status:      1,
				ControlBy:   cModels.ControlBy{},
			})
		}
	}
	return nil
}

// hasExists 验证商品唯一性：品牌＋型号＋包装单位
func (e *Product) productExists(c *dto.ProductInsertReq) bool {
	configure := map[string]interface{}{
		"brand_id":  c.BrandId,
		"mfg_model": c.MfgModel,
		"sales_uom": c.SalesUom,
	}
	var product models.Product
	err := e.Orm.Model(&models.Product{}).Scopes(dto.MakeCondition(configure, c.Id)).First(&product).Error
	if err != nil {
		return false
	}
	return true
}

// vendorIsBind 货主SKU已存在绑定！
func (e *Product) vendorIsBind(c *dto.ProductInsertReq, product *models.Product) bool {
	configure := map[string]interface{}{
		"vendor_id":         c.VendorId,
		"supplier_sku_code": c.SupplierSkuCode,
	}
	err := e.Orm.Scopes(dto.MakeCondition(configure, c.Id)).First(product).Error
	if err != nil {
		return false
	}
	return true
}

// skuIsBind 验证(SKU+货主+货主SKU)与驿站SKU无绑定关系!
func (e *Product) skuIsBind(c *dto.ProductInsertReq) bool {
	configure := map[string]interface{}{
		"sku_code":          c.SkuCode,
		"vendor_id":         c.VendorId,
		"supplier_sku_code": c.SupplierSkuCode,
	}
	var product models.Product
	err := e.Orm.Model(&models.Product{}).Scopes(dto.MakeCondition(configure, c.Id)).First(&product).Error
	if err != nil {
		return false
	}
	return true
}

func (e *Product) getBySkuCode(skuCode string, product *models.Product) error {
	var model models.Product
	err := e.Orm.Model(&model).Where("sku_code = ?", skuCode).First(&product).Error
	return err
}

// createMediaRelationship 绑定商品与图片实例
func (e *Product) createMediaRelationship(tx *gorm.DB, product *models.Product, mediaInstanceId int, mediaInstance dto.MediaInstanceInsertReq) error {
	var model models.MediaRelationship
	mediaRelationship := MediaRelationship{service.Service{
		Orm: tx,
	}}
	mediaRelationshipInsertReq := dto.MediaRelationshipInsertReq{
		MediaTypeId:    1, // 0分类 1商品档案
		BuszId:         product.SkuCode,
		MediaInstantId: mediaInstanceId,
		Watermark:      mediaInstance.WaterMark,
		Seq:            mediaInstance.Seq,
		ControlBy:      product.ControlBy,
	}
	err := mediaRelationship.Insert(&mediaRelationshipInsertReq, &model)
	return err
}

func (e *Product) deleteMediaRelationship(tx *gorm.DB, product *models.Product) error {
	mediaRelationship := MediaRelationship{service.Service{
		Orm: tx,
	}}
	err := mediaRelationship.Remove(1, product.SkuCode)
	if err != nil {
		e.Log.Errorf("Service RemoveMediaRelationship error:%s \r\n", err)
		return err
	}
	return nil
}

// validateForm 表单校验
func (e *Product) validateForm(c *dto.ProductInsertReq, formType string) error {
	if c.ProductCategoryPrimary <= 0 {
		return errors.New("请设置商品主分类！")
	}
	// 请至少选择一个分类！
	if len(c.ProductCategory) <= 0 {
		return errors.New("请至少选择一个分类！")
	}
	// 主产线的税率和产品税率不一致！
	var category dto.CategoryGetResp
	categoryService := Category{e.Service}
	categoryService.Get(&dto.CategoryGetReq{Id: c.ProductCategoryPrimary}, nil, &category)

	var categoryM models.Category
	if err := categoryService.GetCategoryByParentId(c.ProductCategoryPrimary, &categoryM); err == nil && categoryM.Id > 0 {
		return errors.New("主分类非末级分类！")
	}

	if err := categoryService.GetCategory(c.ProductCategoryPrimary, &categoryM); err != nil {
		return errors.New("分类不存在！")
	}

	if category.Tax == "" {
		category.Tax = DefaultTax
	}
	if category.Tax != c.Tax {
		return errors.New("主产线的税率和产品税率不一致！")
	}
	// 保存期限标准为是时，【保存期限(月)必填!保存期限大于1小于999】
	if c.StorageFlag == 1 && (c.StorageTime <= 0 || c.StorageTime > 999) {
		if c.StorageTime <= 0 || c.StorageTime >= 999 {
			return errors.New("保存期限标志为是时,保存期限值必要大于0小于999！")
		} else if _, ok := dictMap[c.StorageTime]; !ok {
			return errors.New("保存期限值填写错误，请参照批注！")
		}
	}
	// 危险品标志为是时，危险品等级必填!
	if c.HazardFlag == 1 && c.HazardClass <= 0 {
		if c.HazardClass <= 0 {
			return errors.New("危险品标志为是时，危险品等级必填！")
		}
		if c.HazardClass > 9 {
			return errors.New("危险品等级参数错误！")
		}
	}
	// 验证商品唯一性：品牌＋型号＋包装单位
	if e.productExists(c) {
		return errors.New("已存在相同品牌、相同型号、相同包装单位的商品！")
	}
	// 产品名称必须包含中文
	if !utils.ContainChinese(c.NameZh) {
		return errors.New("产品名称必须包含中文！")
	}
	// 货主SKU已存在绑定！
	var product models.Product
	if e.vendorIsBind(c, &product) && formType != "sync" {
		return errors.New("货主+货主SKU已存在绑定！")
	}
	if formType == "add" {
		// 新增商品拼装件标志不允许设置为是！
		if c.AssembleFlag == 1 {
			return errors.New("新增商品拼装件标志不允许设置为是！")
		}
	}
	// 编辑时拼装件标志不允许修改为否！
	if formType == "edit" {
		var product models.Product
		err := e.getBySkuCode(c.SkuCode, &product)
		if err != nil {
			return errors.New("商品查询失败！")
		}
		if c.AssembleFlag == 0 && product.AssembleFlag == 1 {
			return errors.New("拼装件标志不允许修改为否！")
		}
	}
	return nil
}

func (e *Product) validateImport(c *dto.ProductImportData, importType string) error {
	var err error
	// 产品名称必须包含中文
	if !utils.ContainChinese(c.NameZh) {
		return errors.New("产品名称必须包含中文！")
	}

	var product models.Product
	if importType == "edit" {
		reg := regexp.MustCompile(`^[A-Za-z]{3}\d{3}$`)
		res := reg.FindAllString(c.SkuCode, -1)
		if len(res) == 0 {
			return errors.New("skuCode格式错误，请检查！")
		}

		err = e.getBySkuCode(c.SkuCode, &product)
		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("数据不存在")
		}
		if err != nil {
			return err
		}
		c.Tax = product.Tax
		c.VendorId = product.VendorId
		c.Id = product.Id
	}

	category := Category{e.Service}

	var categoryLevel1Info models.Category
	err = category.GetCategoryByName(c.Level1CatName, 0, &categoryLevel1Info)
	if err != nil {
		return errors.New("分类1不存在！")
	}

	var categoryLevel2Info models.Category
	err = category.GetCategoryByName(c.Level2CatName, categoryLevel1Info.Id, &categoryLevel2Info)
	if err != nil {
		return errors.New("分类2不存在！")
	}

	var categoryLevel3Info models.Category
	err = category.GetCategoryByName(c.Level3CatName, categoryLevel2Info.Id, &categoryLevel3Info)
	if err != nil {
		return errors.New("分类3不存在！")
	}

	if c.Level4CatName != "" && categoryLevel3Info.Id > 0 {
		var categoryLevel4Info models.Category
		err = category.GetCategoryByName(c.Level4CatName, categoryLevel3Info.Id, &categoryLevel4Info)
		if err != nil {
			return errors.New("分类4不存在！")
		}
		if categoryLevel4Info.Id > 0 && categoryLevel4Info.Status == 1 {
			c.ProductCategory = []dto.ProductCategory{
				{
					CategoryId:   categoryLevel4Info.Id,
					CategoryName: categoryLevel4Info.NameZh,
				},
			}
			c.ProductCategoryPrimary = categoryLevel4Info.Id
		}
	} else if c.Level3CatName != "" && c.Level4CatName == "" {
		if categoryLevel3Info.Id > 0 && categoryLevel3Info.Status == 1 {
			c.ProductCategory = []dto.ProductCategory{
				{
					CategoryId:   categoryLevel3Info.Id,
					CategoryName: categoryLevel3Info.NameZh,
				},
			}
			c.ProductCategoryPrimary = categoryLevel3Info.Id
		}
	}

	if importType == "edit" {
		var data models.ProductCategory
		productCategoryService := ProductCategory{e.Service}
		err = productCategoryService.VerifyProductApproval(c.SkuCode, &data)
		if err != nil {
			return errors.New("SKU产线校验失败！")
		}
		if data.SkuCode == "" {
			return errors.New("未设置产线！")
		}
	}

	// '主产线的税率和产品税率不一致！'
	var categoryInfo dto.CategoryGetResp
	category.Get(&dto.CategoryGetReq{Id: c.ProductCategoryPrimary}, nil, &categoryInfo)
	if categoryInfo.Tax == "" {
		categoryInfo.Tax = DefaultTax
	}
	if categoryInfo.Tax != c.Tax {
		return errors.New("主产线的税率和产品税率不一致！")
	}
	//未找到品牌
	var brandInfo models.Brand
	brandService := Brand{e.Service}
	err = brandService.FindBrandInfoByName(c.BrandZh, c.BrandEn, &brandInfo)
	if err != nil {
		return errors.New("品牌不存在！")
	}
	c.BrandId = brandInfo.Id
	//未找到产线
	if c.ProductCategory[0].CategoryId == 0 {
		return errors.New("未找到产线！")
	}
	//制造厂长度不能超过40个字符！
	if len(c.MfgModel) > 40 {
		return errors.New("制造厂型号长度不能超过40个字符！")
	}
	// 售卖包装单位不存在
	var uommasterInfo models.Uommaster
	uommaster := Uommaster{e.Service}
	if err = uommaster.GetByName(c.SalesUom, &uommasterInfo); err != nil {
		return errors.New("售卖包装单位不存在！")
	}
	if c.PhysicalUom != "" {
		if err = uommaster.GetByName(c.PhysicalUom, &uommasterInfo); err != nil {
			return errors.New("物理单位不存在！")
		}
	}
	if c.StorageFlag == 1 && (c.StorageTime <= 0 || c.StorageTime >= 999) {
		if c.StorageTime <= 0 || c.StorageTime >= 999 {
			return errors.New("保存期限标志为是时,保存期限值必要大于0小于999！")
		} else if _, ok := dictMap[c.StorageTime]; !ok {
			return errors.New("保存期限值填写错误，请参照批注！")
		}
	}

	// 危险品标志为是时，危险品等级必填!
	if c.HazardFlag == 1 && c.HazardClass <= 0 {
		return errors.New("危险品标志为是时，危险品等级必填！")
	}
	// 验证商品唯一性：品牌＋型号＋包装单位
	productInsertReq := dto.ProductInsertReq{
		Id:       product.Id,
		BrandId:  c.BrandId,
		MfgModel: c.MfgModel,
		SalesUom: c.SalesUom,
	}
	if e.productExists(&productInsertReq) {
		return errors.New("已存在相同品牌、相同型号、相同包装单位的商品！")
	}

	return nil
}

// IsApprove 商品是否正在审核或审核未通过
func (e *Product) IsApprove(skuCode string) bool {
	var data models.Product
	result := e.Orm.Where("sku_code = ?", skuCode).Not("status = ?", 2).First(&data)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false
	}
	if result.Error != nil {
		return false
	}
	return true
}

// getProductsForExportAttr 获取附加属性
func (e *Product) getProductsForExportAttr(list *[]models.Product, data *[]dto.ProductExportAttrResp) error {
	var skuCodes []string
	utils.StructColumn(&skuCodes, *list, "SkuCode", "")
	if len(skuCodes) > 0 {
		err := e.Orm.Raw(`
			SELECT p.sku_code as SkuCode,p.name_zh as NameZh,p.mfg_model as MfgModel,pc.category_id as CategoryId,ca.attribute_id as AttributeId,ad.name_zh AS AttributeName,pea.value_zh as AttributeValue
			FROM product as p
			INNER JOIN product_category as pc ON(p.sku_code = pc.sku_code)
			INNER JOIN category_attribute as ca ON(pc.category_id = ca.category_id)
			INNER join attribute_def as ad ON(ca.attribute_id = ad.id)
			LEFT JOIN product_ext_attribute AS pea ON(pea.sku_code = pc.sku_code and ca.attribute_id = pea.attribute_id)
			WHERE p.sku_code IN ? ORDER BY p.sku_code ASC,pc.category_id ASC,ca.seq ASC
		`, skuCodes).Scan(data).Error
		return err
	}
	return nil
}

func (e *Product) getTotalProducts(c *dto.ProductGetPageReq, list *[]models.Product, count *int64) error {
	err := e.Orm.Model(&models.Product{}).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			dto.MakeSearchCondition(c, e.Orm),
		).Find(&list).Count(count).Error
	if err != nil {
		return err
	}
	return nil
}

func (e *Product) productAttributeImportValidate(c *dto.ProductAttributeImportReq) error {
	var err error

	reg := regexp.MustCompile(`^[A-Z]{3}\d{3}$`)
	res := reg.FindAllString(c.SkuCode, -1)
	if len(res) == 0 {
		return errors.New("skuCode存在有误，请检查正确后重新导入！")
	}

	if c.CategoryId == 0 {
		category := Category{e.Service}
		var categoryLevel1Info models.Category
		if c.CategoryLevel1 != "" {
			err = category.GetCategoryByName(c.CategoryLevel1, 0, &categoryLevel1Info)
			if err != nil {
				return errors.New("一级目录不存在！")
			}
		} else {
			return errors.New("一级目录必填！")
		}

		var categoryLevel2Info models.Category
		if c.CategoryLevel2 != "" {
			err = category.GetCategoryByName(c.CategoryLevel2, categoryLevel1Info.Id, &categoryLevel2Info)
			if err != nil {
				return errors.New("二级目录不存在！")
			}
		} else {
			return errors.New("二级目录必填！")
		}

		var categoryLevel3Info models.Category
		if c.CategoryLevel3 != "" {
			err = category.GetCategoryByName(c.CategoryLevel3, categoryLevel2Info.Id, &categoryLevel3Info)
			if err != nil {
				return errors.New("三级目录不存在！")
			}
		} else {
			return errors.New("三级目录必填！")
		}

		if c.CategoryLevel4 != "" && categoryLevel3Info.Id > 0 {
			var categoryLevel4Info models.Category
			err = category.GetCategoryByName(c.CategoryLevel4, categoryLevel3Info.Id, &categoryLevel4Info)
			if err != nil {
				return errors.New("四级目录不存在！")
			}
			if categoryLevel4Info.Id > 0 && categoryLevel4Info.Status == 1 {
				c.CategoryId = categoryLevel4Info.Id
			}
		} else if c.CategoryLevel3 != "" && c.CategoryLevel4 == "" {
			if categoryLevel3Info.Id > 0 && categoryLevel3Info.Status == 1 {
				c.CategoryId = categoryLevel3Info.Id
			}
		}
	}

	var product models.Product
	err = e.getRowBySku(c.SkuCode, &product)
	if err != nil {
		return errors.New("此SKU不存在(" + c.SkuCode + ")")
	} else {
		productCategoryService := ProductCategory{e.Service}
		isExist := productCategoryService.IsExistProductForCategory(c.SkuCode, c.CategoryId)
		if !isExist {
			return errors.New("产品与产线不匹配")
		}
	}
	return nil
}

func (e *Product) getRowBySku(skuCode string, data *models.Product) error {
	var model models.Product
	err := e.Orm.Model(&model).Where("sku_code = ?", skuCode).First(&data).Error
	if err != nil {
		return err
	}
	return nil
}

// FindSkuIsBind (SKU+货主+货主SKU)与驿站SKU是否存在绑定关系
func (e *Product) FindSkuIsBind(model *models.Goods) bool {
	var data models.Product
	_ = e.Orm.Model(&data).Scopes(dto.FindSkuIsBind(&dto.FindSkuIsBindReq{
		SkuCode:         model.SkuCode,
		VendorId:        model.VendorId,
		SupplierSkuCode: model.SupplierSkuCode,
	})).First(&data).Error
	if data.SkuCode != "" {
		return true
	}
	return false
}

//------------------------------------------------------INNER----------------------------------------------------------------

func (e *Product) GetProductBySku(c *dto.InnerGetProductBySkuReq, data *[]dto.InnerGetProductBySkuResp) error {
	var model models.Product
	err := e.Orm.Model(&model).Preload("Brand").Where("sku_code in ?", c.SkuCode).Find(&data).Error
	if err != nil {
		return err
	}
	return nil
}

func (e *Product) GetProductCategoryBySku(c *dto.InnerGetProductBySkuReq, data *[]dto.InnerGetProductCategoryBySkuResp) {
	productCategoryData := make([]models.ProductCategory, 0)
	ProductCategoryService := ProductCategory{e.Service}
	category := Category{e.Service}
	_ = ProductCategoryService.GetCategoryBySku(c.SkuCode, &productCategoryData)
	if len(productCategoryData) > 0 {
		productCategoryResp := make([]dto.InnerGetProductCategoryBySkuResp, 0)

		for _, productCategory := range productCategoryData {
			categoryPath := category.GetCategoryPath(productCategory.CategoryId)
			productCategoryResp = append(productCategoryResp, dto.InnerGetProductCategoryBySkuResp{
				SkuCode:         productCategory.SkuCode,
				ProductCategory: categoryPath,
			})
		}
		*data = productCategoryResp
	}
}

func (e *Product) updateGoodsSupplerSkuCode(tx *gorm.DB, skuCode string, vendorId int, supplierSkuCode string) error {
	err := tx.Model(&models.Goods{}).Where("sku_code", skuCode).Where("vendor_id", vendorId).Update("supplier_sku_Code", supplierSkuCode).Error
	return err
}
