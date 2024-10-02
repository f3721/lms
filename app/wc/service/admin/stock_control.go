package admin

import (
	"encoding/json"
	"errors"
	vd "github.com/bytedance/go-tagexpr/v2/validator"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"github.com/samber/lo"
	modelsPc "go-admin/app/pc/models"
	"go-admin/app/wc/models"
	"go-admin/app/wc/service/admin/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
	dtoStock "go-admin/common/dto/stock/dto"
	"go-admin/common/excel"
	"go-admin/common/global"
	"go-admin/common/utils"
	"gorm.io/gorm"
	"mime/multipart"
	"strconv"
	"strings"
)

type StockControl struct {
	service.Service
}

// 调整单数据权限

func StockControlPermission(p *actions.DataPermission) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db.Where("stock_control.warehouse_code in ?", utils.Split(p.AuthorityWarehouseId))
		db.Where("stock_control.vendor_id in ?", utils.SplitToInt(p.AuthorityVendorId))
		return db
	}
}

// GetPage 获取StockControl列表
func (e *StockControl) GetPage(c *dto.StockControlGetPageReq, p *actions.DataPermission, outData *[]dto.StockControlGetPageResp, count *int64) error {
	var err error
	var data models.StockControl
	pcPrefix := global.GetTenantPcDBNameWithDB(e.Orm)

	db := e.Orm.Model(&data).Preload("Warehouse").Preload("LogicWarehouse").Preload("Vendor").
		Joins("LEFT JOIN stock_control_products controlProduct ON controlProduct.control_code = stock_control.control_code AND controlProduct.deleted_at IS NULL").
		Joins("LEFT JOIN "+pcPrefix+".product product ON product.sku_code = controlProduct.sku_code AND product.deleted_at IS NULL").
		Joins("LEFT JOIN "+pcPrefix+".goods goods ON goods.id = controlProduct.goods_id AND goods.deleted_at IS NULL").
		Select("stock_control.*, sum(controlProduct.quantity) as control_total_num").
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSizeNegative(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
			StockControlPermission(p),
			dtoStock.GenCreatedAtTimeSearch(c.CreatedAtStart, c.CreatedAtEnd, "stock_control"),
			dtoStock.GenProductSearch(c.Sku, c.ProductNo, c.ProductName, "controlProduct"),
		).Group("stock_control.control_code").Order("stock_control.id DESC").Find(outData)
	err = db.Error
	if err != nil {
		e.Log.Errorf("StockControlService GetPage error:%s \r\n", err)
		return err
	}
	for index, _ := range *outData {
		(*outData)[index].InitData()
	}

	err = e.Orm.Table("(?) as u", db.Limit(-1).Offset(-1)).Count(count).Error
	if err != nil {
		e.Log.Errorf("StockControlService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取StockControl对象
func (e *StockControl) Get(d *dto.StockControlGetReq, p *actions.DataPermission, outData *dto.StockControlGetResp) error {
	var data models.StockControl
	var model = &models.StockControl{}
	err := e.Orm.Model(&data).Preload("StockControlProducts").Preload("StockControlProducts.StockLocation").Preload("Warehouse").Preload("LogicWarehouse").Preload("Vendor").
		Scopes(
			actions.Permission(data.TableName(), p),
			StockControlPermission(p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetStockControl error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}

	return FormatDetailFoControl(e.Orm, outData, model, d.Type)
}

// Insert 创建StockControl对象
func (e *StockControl) Insert(c *dto.StockControlInsertReq, action string, p *actions.DataPermission) error {
	var err error
	var data models.StockControl
	c.Generate(&data)

	if lo.IndexOf(utils.Split(p.AuthorityWarehouseId), data.WarehouseCode) == -1 {
		return errors.New("没有实体仓权限")
	}
	if lo.IndexOf(utils.SplitToInt(p.AuthorityVendorId), data.VendorId) == -1 {
		return errors.New("没有货主权限")
	}
	if _, err := data.GenerateControlCode(e.Orm); err != nil {
		e.Log.Errorf("StockControlService Insert error:%s \r\n", err)
		return err
	}
	err = data.InsertControl(e.Orm, action)
	if err != nil {
		e.Log.Errorf("StockControlService Insert error:%s \r\n", err)
		return err
	}
	//记录操作日志
	dataStr, err := e.ControlLogDataHandle(&data)
	if err != nil {
		e.Log.Errorf("StockControlService Insert error:%s \r\n", err)
		return err
	}
	opLog := models.OperateLogs{
		DataId:       data.ControlCode,
		ModelName:    models.ControlLogModelName,
		Type:         models.ControlLogModelInsert,
		Before:       "",
		Data:         dataStr,
		After:        dataStr,
		OperatorId:   c.CreateBy,
		OperatorName: c.CreateByName,
	}
	_ = opLog.InsertItem(e.Orm)
	return nil
}

// Update 修改StockControl对象
func (e *StockControl) Update(c *dto.StockControlUpdateReq, p *actions.DataPermission, action string) error {
	var err error
	var data = models.StockControl{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
		StockControlPermission(p),
	).First(&data, c.GetId())

	if data.Id == 0 {
		return errors.New("调整单不存在或没有权限")
	}

	if data.Status != models.ControlStatus99 {
		return errors.New("调整单状态不正确")
	}

	oldData := data
	stockControlProducts := models.StockControlProducts{}
	oldProducts, err := stockControlProducts.GetByControlCode(e.Orm, data.ControlCode)
	if err != nil {
		return err
	}

	c.Generate(&data)
	err = data.UpdateControl(e.Orm, action)
	if err != nil {
		e.Log.Errorf("StockControlService Update error:%s \r\n", err)
		return err
	}
	//记录操作日志
	logBeforeData := struct {
		models.StockControl
		StockControlProducts string `json:"stockControlProducts"`
	}{}
	if err := utils.CopyDeep(&logBeforeData, &oldData); err != nil {
		e.Log.Errorf("StockControlService Save error:%s \r\n", err)
		return err
	}
	oldProductStrBytes, _ := json.Marshal(&oldProducts)
	logBeforeData.StockControlProducts = string(oldProductStrBytes)
	beforeDataStr, _ := json.Marshal(&logBeforeData)

	dataStr, err := e.ControlLogDataHandle(&data)
	if err != nil {
		e.Log.Errorf("StockControlService Update error:%s \r\n", err)
		return err
	}
	opLog := models.OperateLogs{
		DataId:       data.ControlCode,
		ModelName:    models.ControlLogModelName,
		Type:         models.ControlLogModelUpdate,
		Before:       string(beforeDataStr),
		Data:         dataStr,
		After:        dataStr,
		OperatorId:   c.UpdateBy,
		OperatorName: c.UpdateByName,
	}
	_ = opLog.InsertItem(e.Orm)
	return nil
}

// Remove 删除StockControl
func (e *StockControl) Remove(d *dto.StockControlDeleteReq, p *actions.DataPermission) error {
	var data models.StockControl
	e.Orm.
		Scopes(
			actions.Permission(data.TableName(), p),
			StockControlPermission(p),
		).First(&data, d.GetId())
	if data.Id == 0 {
		return errors.New("调整单不存在或没有权限")
	}
	oldData := data

	if data.Status != models.ControlStatus99 {
		return errors.New("调整单状态不正确")
	}
	data.Status = models.ControlStatus0
	data.UpdateBy = d.UpdateBy
	data.UpdateByName = d.UpdateByName

	db := e.Orm.Save(&data)

	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveStockControl error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	oldDataStr, _ := json.Marshal(oldData)
	dataStr, _ := json.Marshal(data)

	opLog := models.OperateLogs{
		DataId:       data.ControlCode,
		ModelName:    models.ControlLogModelName,
		Type:         models.ControlLogModelDelete,
		Before:       string(oldDataStr),
		Data:         string(dataStr),
		After:        string(dataStr),
		OperatorId:   d.UpdateBy,
		OperatorName: d.UpdateByName,
	}
	_ = opLog.InsertItem(e.Orm)
	return nil
}

// 入库单格式化输出

func FormatDetailFoControl(tx *gorm.DB, data *dto.StockControlGetResp, model *models.StockControl, Type string) error {
	if err := utils.CopyDeep(data, model); err != nil {
		return err
	}
	return data.InitData(tx, model, Type)
}

func (e *StockControl) ControlLogDataHandle(data *models.StockControl) (string, error) {
	logData := struct {
		models.StockControl
		StockControlProducts string `json:"stockControlProducts"`
	}{}
	if err := utils.CopyDeep(&logData, data); err != nil {
		return "", err
	}
	productStrBytes, _ := json.Marshal(&data.StockControlProducts)
	logData.StockControlProducts = string(productStrBytes)
	dataStr, _ := json.Marshal(&logData)
	return string(dataStr), nil
}

// Audit 审核StockControl对象
func (e *StockControl) Audit(c *dto.StockControlAuditReq, p *actions.DataPermission) error {
	var err error
	var data = models.StockControl{}
	e.Orm.Preload("StockControlProducts").Preload("StockControlProducts.StockLocation").Scopes(
		actions.Permission(data.TableName(), p),
		StockControlPermission(p),
	).First(&data, c.GetId())

	if data.Id == 0 {
		return errors.New("调整单不存在或没有权限")
	}

	if data.Status != models.ControlStatus1 {
		return errors.New("调整单状态不正确")
	}

	if data.VerifyStatus != models.ControlVerifyStatus0 {
		return errors.New("调整单审核状态不正确")
	}

	if len(data.StockControlProducts) == 0 {
		return errors.New("调整单关联产品信息异常")
	}

	oldData := data
	c.Generate(&data)

	switch c.VerifyStatus {
	case models.ControlVerifyStatus1:
		if err = data.AuditPass(e.Orm); err != nil {
			e.Log.Errorf("StockControlService Audit error:%s \r\n", err)
			return err
		}
	case models.ControlVerifyStatus2:
		if data.VerifyRemark == "" {
			return errors.New("审核描述不能为空")
		}
		if err = data.AuditReject(e.Orm); err != nil {
			e.Log.Errorf("StockControlService Audit error:%s \r\n", err)
			return err
		}
	}

	beforeDataStr, _ := json.Marshal(&oldData)
	dataStr, _ := json.Marshal(&data)

	opLog := models.OperateLogs{
		DataId:       data.ControlCode,
		ModelName:    models.ControlLogModelName,
		Type:         models.ControlLogModelAudit,
		Before:       string(beforeDataStr),
		Data:         string(dataStr),
		After:        string(dataStr),
		OperatorId:   c.UpdateBy,
		OperatorName: c.UpdateByName,
	}
	_ = opLog.InsertItem(e.Orm)
	return nil
}

func (e *StockControl) ValidateSkus(d *dto.StockControlValidateSkusReq, p *actions.DataPermission, list *[]dto.StockControlValidateSkusResp) error {
	whTo := models.Warehouse{}
	if err := whTo.GetByWarehouseCode(e.Orm, d.WarehouseCode); err != nil {
		return errors.New("实体仓不存在")
	}
	vendor := models.Vendors{}
	if err := vendor.GetById(e.Orm, d.VendorId); err != nil {
		return errors.New("货主不存在")
	}
	skuSliceOri := utils.Split(strings.ReplaceAll(strings.Trim(d.SkuCodes, " "), "，", ","))
	skuSlice := lo.Uniq(skuSliceOri)
	// 调用 pc的sku接口
	skuMap, skuSliceRes := models.GetSkuMapFromPcClient(e.Orm, skuSlice)
	skuErrSlice := lo.Without(skuSlice, skuSliceRes...)
	skuGoodsErrSlice := []string{}
	skuGoodsAuditErrSlice := []string{}

	skuGoodsMap := map[string]modelsPc.Goods{}
	skuGoodsSlice := []string{}
	skuGoodsMap, skuGoodsSlice = models.GetGoodsInfoMapByThreeFromPcClient(e.Orm, skuSlice, d.WarehouseCode, d.VendorId, 0)

	stockInfo := &models.StockInfo{}
	stockInfos, _ := stockInfo.GetStockInfos(e.Orm, d.VendorId, d.LogicWarehouseCode, skuSlice)

	skuGoodsErrSlice = lo.Without(skuSlice, skuGoodsSlice...)
	skuGoodsErrSlice = lo.Without(skuGoodsErrSlice, skuErrSlice...)

	for key, item := range skuGoodsMap {
		if item.ApproveStatus != 1 {
			skuGoodsAuditErrSlice = append(skuGoodsAuditErrSlice, key)
		}
	}
	errStr := FormatValidateSkusErrInfoForControl(skuErrSlice, skuGoodsErrSlice, skuGoodsAuditErrSlice)
	if errStr != "" {
		return errors.New(errStr)
	}
	//获取 库位信息
	skuLocationsMap := map[string][]models.StockLocation{}
	skuLocationsIsTopMap := map[string]bool{}
	for _, item := range skuGoodsMap {
		if stations, isTop, err := models.GetStockLocationsForEntry(e.Orm, d.LogicWarehouseCode, item.Id); err == nil {
			skuLocationsMap[item.SkuCode] = *stations
			skuLocationsIsTopMap[item.SkuCode] = isTop
		} else {
			skuLocationsMap[item.SkuCode] = []models.StockLocation{}
			skuLocationsIsTopMap[item.SkuCode] = false
		}
	}

	stockControlValidateSkusResp := dto.StockControlValidateSkusResp{}
	var stockLocationId int
	for _, item := range skuSliceOri {
		stockControlValidateSkusResp.SkuCode = item
		stockControlValidateSkusResp.VendorName = vendor.NameZh
		stockControlValidateSkusResp.ProductName = skuMap[item].NameZh
		stockControlValidateSkusResp.MfgModel = skuMap[item].MfgModel
		stockControlValidateSkusResp.BrandName = skuMap[item].Brand.BrandZh
		stockControlValidateSkusResp.SalesUom = skuMap[item].SalesUom
		stockControlValidateSkusResp.VendorSkuCode = skuMap[item].SupplierSkuCode
		stockControlValidateSkusResp.ProductNo = skuGoodsMap[item].ProductNo
		stockControlValidateSkusResp.CurrentQuantity = d.CurrentQuantity

		stockControlValidateSkusResp.StockLocation = []dto.StockControlProductsLocationGetResp{}
		stockControlProductsLocationGetResp := dto.StockControlProductsLocationGetResp{}
		if skuLocationsIsTopMap[item] {
			stockLocationId = skuLocationsMap[item][0].Id
		}
		skuLocationIds := []int{}
		for _, itemLocation := range skuLocationsMap[item] {
			skuLocationIds = append(skuLocationIds, itemLocation.Id)
		}

		stockLocationGoods, err := models.GetLocationGoodsByGoodsIdAndLocationIdSlice(e.Orm, skuGoodsMap[item].Id, skuLocationIds)
		if err != nil {
			return err
		}

		for _, itemLocation := range skuLocationsMap[item] {
			stockControlProductsLocationGetResp.Regenerate(itemLocation)
			if itemLocation.Id == d.StockLocationId {
				stockLocationId = d.StockLocationId
			}
			for _, slg := range *stockLocationGoods {
				stockControlProductsLocationGetResp.Stock = 0
				stockControlProductsLocationGetResp.LockStock = 0
				if stockControlProductsLocationGetResp.Id == slg.LocationId {
					stockControlProductsLocationGetResp.Stock = slg.Stock
					stockControlProductsLocationGetResp.LockStock = slg.LockStock
					break
				}
			}
			stockControlValidateSkusResp.StockLocation = append(stockControlValidateSkusResp.StockLocation, stockControlProductsLocationGetResp)
		}
		stockControlValidateSkusResp.StockLocationId = stockLocationId

		for _, sv := range *stockInfos {
			if sv.SkuCode == item {
				stockControlValidateSkusResp.Stock = sv.Stock
				stockControlValidateSkusResp.LockStock = sv.LockStock
				stockControlValidateSkusResp.TotalStock = sv.Stock + sv.LockStock
				break
			}
		}

		for _, sv := range stockControlValidateSkusResp.StockLocation {
			if sv.Id == stockControlValidateSkusResp.StockLocationId {
				// 盘后数量 < 库位库存 盘减
				stockControlProductsTotal := sv.Stock + sv.LockStock
				if d.CurrentQuantity < stockControlProductsTotal {
					stockControlValidateSkusResp.Type = "1"
					stockControlValidateSkusResp.Quantity = stockControlProductsTotal - d.CurrentQuantity
				} else {
					stockControlValidateSkusResp.Type = "0"
					stockControlValidateSkusResp.Quantity = d.CurrentQuantity - stockControlProductsTotal
					if stockControlValidateSkusResp.Quantity == 0 {
						stockControlValidateSkusResp.Type = "2"
					}
				}

				// 盘后数量 < 锁库库位库存 报错提醒
				if d.CurrentQuantity < sv.LockStock {
					stockControlValidateSkusResp.ErrorMsg = "库位可用库存不足"
				}
			}
		}

		*list = append(*list, stockControlValidateSkusResp)
	}
	return nil
}

func FormatValidateSkusErrInfoForControl(skuErrSlice, skuGoodsErrSlice, skuGoodsAuditErrSlice []string) string {
	errSlice := []string{}
	if len(skuErrSlice) > 0 {
		errSku := strings.Join(skuErrSlice, ",")
		errSlice = append(errSlice, "【"+errSku+"】不存在")
	}
	if len(skuGoodsErrSlice) > 0 {
		errSku := strings.Join(skuGoodsErrSlice, ",")
		errSlice = append(errSlice, "【"+errSku+"】没有对应货主和实体仓的商品关系")
	}
	if len(skuGoodsAuditErrSlice) > 0 {
		errSku := strings.Join(skuGoodsAuditErrSlice, ",")
		errSlice = append(errSlice, "【"+errSku+"】商品关系待审核或审核不通过")
	}
	if len(errSlice) > 0 {
		return strings.Join(errSlice, "，") + "，请联系相关人员处理！"
	}
	return ""
}

// StockImport 批量新增维护导入
func (e *StockControl) StockImport(file *multipart.FileHeader, p *actions.DataPermission) (err error, errTitleList []map[string]string, errData []map[string]interface{}) {
	excelApp := excel.NewExcel()
	err, datas, titleList := excelApp.GetExcelData(file)
	if err != nil {
		return
	}
	if len(datas) == 0 {
		err = errors.New("导入数据不能为空！")
		return
	}
	if len(datas) > 0 && len(datas) <= 20000 {
		stockControlProductsReq, baseStockControl, errMsgList := checkImport(e.Orm, datas)

		errMsgCount := 0
		for _, emv := range errMsgList {
			if emv != "" {
				errMsgCount += 1
			}
		}

		if errMsgCount > 0 {
			// 有校验错误的行时调用 导出保存excel
			err = errors.New("有错误！")
			errTitleList, errData = excelApp.MergeErrMsgColumn(titleList, datas, errMsgList)
		} else {
			if lo.IndexOf(utils.Split(p.AuthorityWarehouseId), baseStockControl.BaseWarehouseCode) == -1 {
				err = errors.New("没有实体仓权限")
				return
			}
			if lo.IndexOf(utils.SplitToInt(p.AuthorityVendorId), baseStockControl.BaseVendorId) == -1 {
				err = errors.New("没有货主权限")
				return
			}

			// TODO 全部成功 插入数据库
			insertStockData := dto.StockControlInsertReq{}
			insertStockData.WarehouseCode = baseStockControl.BaseWarehouseCode
			insertStockData.LogicWarehouseCode = baseStockControl.BaseLogicWarehouseCode
			insertStockData.VendorId = baseStockControl.BaseVendorId
			for _, scrv := range stockControlProductsReq {
				insertStockData.StockControlProducts = append(insertStockData.StockControlProducts, *scrv)
			}

			// 设置创建人
			ctx := e.Orm.Statement.Context.(*gin.Context)
			insertStockData.SetCreateBy(user.GetUserId(ctx))
			insertStockData.SetCreateByName(user.GetUserName(ctx))

			var data models.StockControl
			if _, GenerateErr := data.GenerateControlCode(e.Orm); GenerateErr != nil {
				e.Log.Errorf("StockControlService Insert error:%s \r\n", GenerateErr)
				err = errors.New(GenerateErr.Error())
				return
			}

			insertStockData.Generate(&data)
			err = data.InsertImmediately(e.Orm, "Commit")
			if err != nil {
				e.Log.Errorf("StockControlService Insert error:%s \r\n", err)
				return
			}
		}
		return
	} else {
		err = errors.New("导入数据不能超过2万行！")
		return
	}
}

func checkImport(tx *gorm.DB, datas []map[string]interface{}) (StockControlProductsReq map[string]*dto.StockControlProductsReq, baseStockControl *dto.BaseStockControl, errMsgList map[int]string) {
	var err error
	errMsgList = make(map[int]string)
	var importReq *dto.ImportReq
	importReq = new(dto.ImportReq)
	baseStockControl = new(dto.BaseStockControl)
	StockControlProductsReq = make(map[string]*dto.StockControlProductsReq)

	var errStr string
	var skuCodes []string
	var skuCodeList map[string]int
	skuCodeList = make(map[string]int)
	var uniqueCodeList map[string]int
	uniqueCodeList = make(map[string]int)

	var stockLocationCodeToId map[string]int
	stockLocationCodeToId = make(map[string]int)

	for i, stockControls := range datas {
		errStr = ""
		var importStockControl dto.ImportStockControl
		importStockControl.Trim(stockControls)

		// 基本数据校验
		err = vd.Validate(importStockControl)
		if err != nil {
			errStr = err.Error()
		}
		// 商品所属货主、实体仓、逻辑仓必须唯一
		if i == 0 {
			baseStockControl.BaseVendorName = importStockControl.VendorName
			modelVendors := models.Vendors{}
			var vendor, vendorErr = modelVendors.GetByNameZh(tx, baseStockControl.BaseVendorName)
			if vendorErr != nil {
				errStr = "货主不存在;"
				errMsgList[i] = errStr
				continue
			}
			baseStockControl.BaseVendorId = vendor.Id

			baseStockControl.BaseWarehouseName = importStockControl.WarehouseName
			modelWarehouse := models.Warehouse{}
			var warehouse, warehouseErr = modelWarehouse.GetWarehoseByName(tx, baseStockControl.BaseWarehouseName)
			if warehouseErr != nil {
				errStr = "实体仓不存在;"
				errMsgList[i] = errStr
				continue
			}
			if warehouse.Status != "1" {
				errStr = "实体仓状态不正确;"
				errMsgList[i] = errStr
				continue
			}
			baseStockControl.BaseWarehouseCode = warehouse.WarehouseCode

			baseStockControl.BaseLogicWarehouseName = importStockControl.LogicWarehouseName
			modelLogicWarehouse := models.LogicWarehouse{}
			var logicWarehouse, logicWarehouseErr = modelLogicWarehouse.GetLogicWarehouseByWarehouseName(tx, baseStockControl.BaseLogicWarehouseName)
			if logicWarehouseErr != nil {
				errStr = "逻辑仓不存在;"
				errMsgList[i] = errStr
				continue
			}
			if logicWarehouse.Status != "1" {
				errStr = "逻辑仓状态不正确;"
				errMsgList[i] = errStr
				continue
			}
			baseStockControl.BaseLogicWarehouseCode = logicWarehouse.LogicWarehouseCode

			stockLocations, stockLocationsErr := models.GetStockLocationByLwhCodeWithLimit(tx, baseStockControl.BaseLogicWarehouseCode)
			if stockLocationsErr != nil {
				errStr = "逻辑仓下无可用库位编号;"
				errMsgList[i] = errStr
				continue
			}
			stockControlProductsLocationGetResp := dto.StockControlProductsLocationGetResp{}
			for _, item := range *stockLocations {
				if item.Status != "1" {
					continue
				}
				stockControlProductsLocationGetResp.Regenerate(item)
				baseStockControl.StockLocation = append(baseStockControl.StockLocation, stockControlProductsLocationGetResp)
			}
		} else {
			if baseStockControl.BaseVendorName != importStockControl.VendorName {
				errStr = "商品所属货主不唯一;"
				errMsgList[i] = errStr
				continue
			}
			if baseStockControl.BaseWarehouseName != importStockControl.WarehouseName {
				errStr = "实体仓不唯一;"
				errMsgList[i] = errStr
				continue
			}
			if baseStockControl.BaseLogicWarehouseName != importStockControl.LogicWarehouseName {
				errStr = "逻辑仓不唯一;"
				errMsgList[i] = errStr
				continue
			}
			if len(baseStockControl.StockLocation) == 0 {
				errStr = "逻辑仓下无可用库位编号;"
				errMsgList[i] = errStr
				continue
			}
		}

		uniqueCode := utils.Md5Uc(importStockControl.SkuCode + ":" + importStockControl.LocationCode)
		if value, ok := uniqueCodeList[uniqueCode]; ok {
			errStr = importStockControl.SkuCode + "+库位与第" + strconv.Itoa(value+3) + "行数据重复;"
			errMsgList[i] = errStr
			continue
		}
		uniqueCodeList[uniqueCode] = i

		stockLocationFlag := false
		for _, slv := range baseStockControl.StockLocation {
			stockLocationCodeToId[slv.LocationCode] = slv.Id
			if slv.LocationCode == importStockControl.LocationCode {
				stockLocationFlag = true
				break
			}
		}
		if stockLocationFlag == false {
			errStr = "逻辑仓和库位编号不一致;"
		}

		if importStockControl.CurrentQuantity < 0 {
			errStr = "盘后数量为大于等于0的整数;"
			errMsgList[i] = errStr
			continue
		}

		if errStr == "" {
			skuCodes = append(skuCodes, importStockControl.SkuCode)
			skuCodeList[importStockControl.SkuCode] = i
			importStockControl.Key = i
			importReq.Data = append(importReq.Data, importStockControl)
		}

		StockControlProductsReq[uniqueCode] = &dto.StockControlProductsReq{}
		StockControlProductsReq[uniqueCode].VendorId = baseStockControl.BaseVendorId
		StockControlProductsReq[uniqueCode].WarehouseCode = baseStockControl.BaseWarehouseCode
		StockControlProductsReq[uniqueCode].LogicWarehouseCode = baseStockControl.BaseLogicWarehouseCode
		StockControlProductsReq[uniqueCode].SkuCode = importStockControl.SkuCode
		StockControlProductsReq[uniqueCode].CurrentQuantity = importStockControl.CurrentQuantity
		StockControlProductsReq[uniqueCode].StockLocationId = stockLocationCodeToId[importStockControl.LocationCode]
		errMsgList[i] = errStr
	}

	var existSkus []string
	modelGoods := modelsPc.Goods{}
	goods, _ := modelGoods.GetGoodsBySku(tx, skuCodes, baseStockControl.BaseWarehouseCode, baseStockControl.BaseVendorId)

	var skuGoodsInfo map[string]int
	skuGoodsInfo = make(map[string]int)
	for _, gv := range *goods {
		skuGoodsInfo[gv.SkuCode] = gv.Id
		existSkus = append(existSkus, gv.SkuCode)

		for idx, itm := range StockControlProductsReq {
			if gv.SkuCode == itm.SkuCode {
				StockControlProductsReq[idx].GoodsId = gv.Id
			}
		}
	}

	for sck, scv := range skuCodeList {
		errStr = ""
		if !utils.InArrayString(sck, existSkus) {
			errMsgList[scv] += "SKU不存在;"
			continue
		}
	}

	modelStockInfo := models.StockInfo{}
	stockInfos, _ := modelStockInfo.GetStockInfos(tx, baseStockControl.BaseVendorId, baseStockControl.BaseLogicWarehouseCode, existSkus)

	// 统计库存表的库存数量
	goodsCurrentQuantity := map[int]*dto.GoodsCurrentQuantity{}
	var goodsIds []int
	for _, sv := range *stockInfos {
		goodsIds = append(goodsIds, sv.GoodsId)
	}

	modelStockLocationGoods := models.StockLocationGoods{}
	stockLocationGoods, _ := modelStockLocationGoods.GetByGoodsId(tx, goodsIds)
	for slgk, slgv := range *stockLocationGoods {
		goodsCurrentQuantity[slgk] = &dto.GoodsCurrentQuantity{
			Keys:            make([]int, 0),
			SkuCode:         "",
			GoodsId:         slgv.GoodsId,
			LocationId:      slgv.LocationId,
			CurrentQuantity: slgv.Stock + slgv.LockStock,
			Stock:           slgv.Stock,
			LockStock:       slgv.LockStock,
			Quantity:        0,
		}

		for _, idv := range importReq.Data {
			if skuGoodsInfo[idv.SkuCode] == slgv.GoodsId && stockLocationCodeToId[idv.LocationCode] == slgv.LocationId {
				goodsCurrentQuantity[slgk].Keys = append(goodsCurrentQuantity[slgk].Keys, idv.Key)
				goodsCurrentQuantity[slgk].SkuCode = idv.SkuCode
				goodsCurrentQuantity[slgk].CurrentQuantity = idv.CurrentQuantity
				goodsCurrentQuantity[slgk].Stock = idv.CurrentQuantity - slgv.LockStock
				goodsCurrentQuantity[slgk].LockStock = slgv.LockStock
				goodsCurrentQuantity[slgk].Quantity = idv.CurrentQuantity - slgv.LockStock - slgv.Stock
				// 判断锁库库位库存是否大于当前商品的库位盘后数量
				if slgv.LockStock > idv.CurrentQuantity {
					errMsgList[idv.Key] += "库位可用库存不足;"
				} else {
					uniqueCode := utils.Md5Uc(idv.SkuCode + ":" + idv.LocationCode)
					if slgv.LockStock+slgv.Stock > idv.CurrentQuantity { // 盘亏
						StockControlProductsReq[uniqueCode].Type = "1"
					} else if slgv.LockStock+slgv.Stock == idv.CurrentQuantity { // 无差异
						StockControlProductsReq[uniqueCode].Type = "2"
					} else { // 盘盈
						StockControlProductsReq[uniqueCode].Type = "0"
					}
					StockControlProductsReq[uniqueCode].Quantity = utils.AbsToInt(idv.CurrentQuantity - slgv.LockStock - slgv.Stock)
				}
			}
		}
	}

	locationAvailableStockTotal := map[int]int{}
	for _, gcqv := range goodsCurrentQuantity { // 统计每个good的异动数量
		locationAvailableStockTotal[gcqv.GoodsId] += gcqv.Quantity
	}

	for _, sv := range *stockInfos {
		if _, ok := locationAvailableStockTotal[sv.GoodsId]; ok {
			// 库位的可用库存 + 异动数量 < 0 则表示库存不足
			if sv.Stock+locationAvailableStockTotal[sv.GoodsId] < 0 {
				// 遍历所有的库位数据 将其标记
				for _, gcqv2 := range goodsCurrentQuantity {
					// 判断是否goods id是否一致
					if gcqv2.GoodsId == sv.GoodsId {
						for _, key := range gcqv2.Keys {
							errMsgList[key] += "可用库存不足；"
						}
					}
				}
			}
		}
	}

	return
}
