package admin

import (
	"encoding/json"
	"errors"
	modelsPc "go-admin/app/pc/models"
	dtoStock "go-admin/common/dto/stock/dto"
	"go-admin/common/global"
	models2 "go-admin/common/models"
	"go-admin/common/utils"
	"strings"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"go-admin/app/wc/models"
	"go-admin/app/wc/service/admin/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type StockTransfer struct {
	service.Service
}

// 调拨单数据权限
func StockTransferPermission(p *actions.DataPermission) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db.Where("stock_transfer.from_warehouse_code in ?", utils.Split(p.AuthorityWarehouseAllocateId))
		db.Where("stock_transfer.to_warehouse_code in ?", utils.Split(p.AuthorityWarehouseAllocateId))
		db.Where("stock_transfer.vendor_id in ?", utils.SplitToInt(p.AuthorityVendorId))
		return db
	}
}

// GetPage 获取StockTransfer列表
func (e *StockTransfer) GetPage(c *dto.StockTransferGetPageReq, p *actions.DataPermission, outData *[]dto.StockTransferGetPageResp, count *int64) error {
	var err error
	var data models.StockTransfer
	var list = &[]models.StockTransfer{}
	pcPrefix := global.GetTenantPcDBNameWithDB(e.Orm)

	db := e.Orm.Model(&data).Preload("Vendor").Preload("ToWarehouse").Preload("FromLogicWarehouse").Preload("ToLogicWarehouse").Preload("FromWarehouse").
		Joins("LEFT JOIN stock_transfer_products tranProduct ON tranProduct.transfer_code = stock_transfer.transfer_code AND tranProduct.deleted_at IS NULL").
		Joins("LEFT JOIN "+pcPrefix+".product product ON product.sku_code = tranProduct.sku_code AND product.deleted_at IS NULL").
		Joins("LEFT JOIN "+pcPrefix+".goods goods ON goods.id = tranProduct.goods_id AND goods.deleted_at IS NULL").
		Select("stock_transfer.*").
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSizeNegative(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
			StockTransferPermission(p),
			dtoStock.GenCreatedAtTimeSearch(c.CreatedAtStart, c.CreatedAtEnd, "stock_transfer"),
			dtoStock.GenProductSearch(c.Sku, c.ProductNo, c.ProductName, "tranProduct"),
		).Order("stock_transfer.id DESC").Group("stock_transfer.transfer_code").Find(list)
	err = db.Error
	if err != nil {
		e.Log.Errorf("StockTransferService GetPage error:%s \r\n", err)
		return err
	}

	for _, item := range *list {
		dataItem := dto.StockTransferGetPageResp{}
		if err := utils.CopyDeep(&dataItem, &item); err != nil {
			e.Log.Errorf("StockTransferService GetPage error:%s \r\n", err)
			return err
		}
		dataItem.VendorName = dataItem.Vendor.NameZh
		dataItem.ToWarehouseName = dataItem.ToWarehouse.WarehouseName
		dataItem.FromWarehouseName = dataItem.FromWarehouse.WarehouseName
		dataItem.ToLogicWarehouseName = dataItem.ToLogicWarehouse.LogicWarehouseName
		dataItem.FromLogicWarehouseName = dataItem.FromLogicWarehouse.LogicWarehouseName

		dataItem.TypeName = utils.GetTFromMap(dataItem.Type, models.TransferTypeMap)
		dataItem.StatusName = utils.GetTFromMap(dataItem.Status, models.TransferStatusMap)
		dataItem.VerifyStatusName = utils.GetTFromMap(dataItem.VerifyStatus, models.TransferVerifyStatusMap)

		// 权限处理 todo
		dataItem.SetTransferRulesByStatus()

		*outData = append(*outData, dataItem)
	}

	err = e.Orm.Table("(?) as u", db.Limit(-1).Offset(-1)).Count(count).Error
	if err != nil {
		e.Log.Errorf("StockTransferService GetPage error:%s \r\n", err)
		return err
	}

	return nil
}

// Get 获取StockTransfer对象
func (e *StockTransfer) Get(d *dto.StockTransferGetReq, p *actions.DataPermission, outData *dto.StockTransferGetResp) error {
	var data models.StockTransfer
	var model = &models.StockTransfer{}

	err := e.Orm.Model(&data).Preload("StockTransferProducts").Preload("Vendor").Preload("ToWarehouse").Preload("FromLogicWarehouse").Preload("ToLogicWarehouse").Preload("FromWarehouse").
		Scopes(
			actions.Permission(data.TableName(), p),
			StockTransferPermission(p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetStockTransfer error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	if err := e.FormatDetailForTransfer(e.Orm, outData, model); err != nil {
		return err
	}
	return nil
}

// 调拨单格式化输出
func (e *StockTransfer) FormatDetailForTransfer(tx *gorm.DB, data *dto.StockTransferGetResp, model *models.StockTransfer) error {
	if err := utils.CopyDeep(data, model); err != nil {
		return err
	}
	return data.InitData(tx, model)
}

// Insert 创建StockTransfer对象
func (e *StockTransfer) Insert(c *dto.StockTransferInsertReq, action string, p *actions.DataPermission) error {
	var err error
	var data models.StockTransfer

	c.Generate(e.Orm, &data)

	if lo.IndexOf(utils.Split(p.AuthorityWarehouseAllocateId), data.ToWarehouseCode) == -1 {
		return errors.New("没有入库实体仓权限")
	}
	if lo.IndexOf(utils.Split(p.AuthorityWarehouseAllocateId), data.FromWarehouseCode) == -1 {
		return errors.New("没有出库实体仓权限")
	}
	if lo.IndexOf(utils.SplitToInt(p.AuthorityVendorId), data.VendorId) == -1 {
		return errors.New("没有货主权限")
	}
	if _, err = data.GenerateTransferCode(e.Orm); err != nil {
		e.Log.Errorf("StockTransferService Insert error:%s \r\n", err)
		return err
	}
	err = data.InsertTransfer(e.Orm, action)
	if err != nil {
		e.Log.Errorf("StockTransferService Insert error:%s \r\n", err)
		return err
	}
	//记录操作日志
	dataStr, err := e.transferLogDataHandle(&data)
	if err != nil {
		e.Log.Errorf("StockTransferService Insert error:%s \r\n", err)
		return err
	}
	opLog := models.OperateLogs{
		DataId:       data.TransferCode,
		ModelName:    models.TransferLogModelName,
		Type:         models.TransferLogModelInsert,
		Before:       "",
		Data:         dataStr,
		After:        dataStr,
		OperatorId:   c.CreateBy,
		OperatorName: c.CreateByName,
	}
	_ = opLog.InsertItem(e.Orm)
	return nil
}

// Update 修改StockTransfer对象
func (e *StockTransfer) Update(c *dto.StockTransferUpdateReq, p *actions.DataPermission, action string) error {
	var err error
	var data = models.StockTransfer{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
		StockTransferPermission(p),
	).First(&data, c.GetId())

	if data.Id == 0 {
		return errors.New("调拨单不存在或没有权限")
	}

	if data.Status != models.TransferStatus99 {
		return errors.New("调拨单状态不正确")
	}

	oldData := data
	stockTransferProducts := models.StockTransferProducts{}
	oldProducts, err := stockTransferProducts.GetByTransferCode(e.Orm, data.TransferCode)
	if err != nil {
		return err
	}

	c.Generate(e.Orm, &data)

	err = data.UpdateTransfer(e.Orm, action)
	if err != nil {
		e.Log.Errorf("StockTransferService Save error:%s \r\n", err)
		return err
	}

	//记录操作日志
	logBeforeData := struct {
		models.StockTransfer
		StockTransferProducts string `json:"stockTransferProducts"`
	}{}
	if err := utils.CopyDeep(&logBeforeData, &oldData); err != nil {
		e.Log.Errorf("StockTransferService Save error:%s \r\n", err)
		return err
	}
	oldProductStrBytes, _ := json.Marshal(&oldProducts)
	logBeforeData.StockTransferProducts = string(oldProductStrBytes)
	beforeDataStr, _ := json.Marshal(&logBeforeData)

	dataStr, err := e.transferLogDataHandle(&data)
	if err != nil {
		e.Log.Errorf("StockTransferService Save error:%s \r\n", err)
		return err
	}
	opLog := models.OperateLogs{
		DataId:       data.TransferCode,
		ModelName:    models.TransferLogModelName,
		Type:         models.TransferLogModelUpdate,
		Before:       string(beforeDataStr),
		Data:         dataStr,
		After:        dataStr,
		OperatorId:   c.UpdateBy,
		OperatorName: c.UpdateByName,
	}
	_ = opLog.InsertItem(e.Orm)
	return nil
}

// Remove 删除StockTransfer
func (e *StockTransfer) Remove(c *dto.StockTransferDeleteReq, p *actions.DataPermission) error {
	var data models.StockTransfer

	e.Orm.
		Scopes(
			actions.Permission(data.TableName(), p),
			StockTransferPermission(p),
		).First(&data, c.GetId())

	if data.Id == 0 {
		return errors.New("调拨单不存在或没有权限")
	}
	oldData := data

	if data.Status != models.TransferStatus99 {
		return errors.New("调拨单状态不正确")
	}
	data.Status = models.TransferStatus0
	data.UpdateBy = c.UpdateBy
	data.UpdateByName = c.UpdateByName

	db := e.Orm.Save(&data)

	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveStockTransfer error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	oldDataStr, _ := json.Marshal(oldData)
	dataStr, _ := json.Marshal(data)

	opLog := models.OperateLogs{
		DataId:       data.TransferCode,
		ModelName:    models.TransferLogModelName,
		Type:         models.TransferLogModelDelete,
		Before:       string(oldDataStr),
		Data:         string(dataStr),
		After:        string(dataStr),
		OperatorId:   c.UpdateBy,
		OperatorName: c.UpdateByName,
	}
	_ = opLog.InsertItem(e.Orm)
	return nil
}

// 日志 stockTransferProducts值转json字符串
func (e *StockTransfer) transferLogDataHandle(data *models.StockTransfer) (string, error) {
	logData := struct {
		models.StockTransfer
		StockTransferProducts string `json:"stockTransferProducts"`
	}{}
	if err := utils.CopyDeep(&logData, data); err != nil {
		return "", err
	}
	productStrBytes, _ := json.Marshal(&data.StockTransferProducts)
	logData.StockTransferProducts = string(productStrBytes)
	dataStr, _ := json.Marshal(&logData)
	return string(dataStr), nil
}

// Audit 审核StockTransfer对象
func (e *StockTransfer) Audit(c *dto.StockTransferAuditReq, p *actions.DataPermission) error {
	var err error
	var data = models.StockTransfer{}
	e.Orm.Preload("StockTransferProducts").Preload("ToWarehouse").Preload("FromLogicWarehouse").Preload("ToLogicWarehouse").Preload("FromWarehouse").Scopes(
		actions.Permission(data.TableName(), p),
		StockTransferPermission(p),
	).First(&data, c.GetId())

	if data.Id == 0 {
		return errors.New("调拨单不存在或没有权限")
	}

	if data.Status != models.TransferStatus98 {
		return errors.New("调拨单状态不正确")
	}

	if data.VerifyStatus != models.TransferVerifyStatus0 {
		return errors.New("调拨单审核状态不正确")
	}

	if len(data.StockTransferProducts) == 0 || data.FromWarehouse.Id == 0 || data.ToWarehouse.Id == 0 || data.FromLogicWarehouse.Id == 0 || data.ToLogicWarehouse.Id == 0 {
		return errors.New("调拨单关联产品、仓库信息异常")
	}

	oldData := data
	c.Generate(&data)

	switch c.VerifyStatus {
	case models.TransferVerifyStatus1:
		if err = e.AuditPass(&data); err != nil {
			e.Log.Errorf("StockTransferService Audit error:%s \r\n", err)
			return err
		}
	case models.TransferVerifyStatus2:
		if data.VerifyRemark == "" {
			return errors.New("调拨单审核描述不能为空")
		}
		if err = data.AuditReject(e.Orm); err != nil {
			e.Log.Errorf("StockTransferService Audit error:%s \r\n", err)
			return err
		}
	}

	beforeDataStr, _ := json.Marshal(&oldData)
	dataStr, _ := json.Marshal(&data)

	opLog := models.OperateLogs{
		DataId:       data.TransferCode,
		ModelName:    models.TransferLogModelName,
		Type:         models.TransferLogModelAudit,
		Before:       string(beforeDataStr),
		Data:         string(dataStr),
		After:        string(dataStr),
		OperatorId:   c.UpdateBy,
		OperatorName: c.UpdateByName,
	}
	_ = opLog.InsertItem(e.Orm)
	return nil
}

// 调拨单审核通过
func (e *StockTransfer) AuditPass(data *models.StockTransfer) error {
	data.SetOutboundStatus1()
	return e.Orm.Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit(clause.Associations).Save(data).Error; err != nil {
			return err
		}

		//创建出库单和入库单
		if err := CreateOutboundAndEntry(tx, data); err != nil {
			return err
		}
		return nil
	})
}

// 创建出库单和入库单
func CreateOutboundAndEntry(tx *gorm.DB, data *models.StockTransfer) error {
	stockOutbound, err := CreateOutbound(tx, data)
	if err != nil {
		return err
	}
	if _, err = CreateEntry(tx, data); err != nil {
		return err
	}
	//虚拟仓自动出库
	if data.FromWarehouse.CheckIsVirtual() {
		stockOutbound.SetConfirmOutboundStatus(data.UpdateBy, data.UpdateByName)
		if err := stockOutbound.ConfirmOutboundAutoForVirtual(tx); err != nil {
			return err
		}
		// 调拨单状态更改为已出库
		if err := data.SetOutboundCompleteStatus(tx); err != nil {
			return err
		}
	}

	return nil
}

// 创建出库单
func CreateOutbound(tx *gorm.DB, data *models.StockTransfer) (*models.StockOutbound, error) {
	productsReqSlice := []dto.StockOutboundProductsReq{}
	hasSub := 0
	if !data.FromWarehouse.CheckIsVirtual() {
		hasSub = 1
	}
	for _, item := range data.StockTransferProducts {
		productsReqSlice = append(productsReqSlice, dto.StockOutboundProductsReq{
			SkuCode:  item.SkuCode,
			Quantity: item.Quantity,
			GoodsId:  item.GoodsId,
			HasSub:   hasSub,
		})
	}
	insertReq := &dto.StockOutboundInsertReq{
		SourceCode:            data.TransferCode,
		Remark:                data.Remark,
		WarehouseCode:         data.FromWarehouseCode,
		LogicWarehouseCode:    data.FromLogicWarehouseCode,
		VendorId:              data.VendorId,
		StockOutboundProducts: productsReqSlice,
		ControlBy: models2.ControlBy{
			CreateBy:     data.CreateBy,
			CreateByName: data.CreateByName,
		},
	}

	stockOutbound, err := InsertOutbound(tx, insertReq, models.OutboundType0)
	if err != nil {
		return nil, err
	}

	//锁定并确定库位(非虚拟仓)
	if !data.FromWarehouse.CheckIsVirtual() {
		for _, item := range stockOutbound.StockOutboundProducts {
			if err := item.SplitProductAndLockLocation(tx, stockOutbound.OutboundCode, stockOutbound.Type); err != nil {
				return nil, err
			}
		}
	}
	return stockOutbound, nil
}

// 创建入库单
func CreateEntry(tx *gorm.DB, data *models.StockTransfer) (*models.StockEntry, error) {
	productsReqSlice := []dto.StockEntryProductsReq{}
	for _, item := range data.StockTransferProducts {
		productsReqSlice = append(productsReqSlice, dto.StockEntryProductsReq{
			SkuCode:  item.SkuCode,
			Quantity: item.Quantity,
			GoodsId:  item.ToGoodsId,
		})
	}
	insertReq := &dto.StockEntryInsertReq{
		SourceCode:         data.TransferCode,
		Remark:             data.Remark,
		WarehouseCode:      data.ToWarehouseCode,
		LogicWarehouseCode: data.ToLogicWarehouseCode,
		VendorId:           data.VendorId,
		StockEntryProducts: productsReqSlice,
		ControlBy: models2.ControlBy{
			CreateBy:     data.CreateBy,
			CreateByName: data.CreateByName,
		},
	}
	return InsertEntry(tx, insertReq, models.EntryType0)
}

func (e *StockTransfer) ValidateSkus(d *dto.StockTransferValidateSkusReq, p *actions.DataPermission, list *[]dto.StockTransferValidateSkusResp) error {
	whTo := models.Warehouse{}
	if err := whTo.GetByWarehouseCode(e.Orm, d.ToWarehouseCode); err != nil {
		return errors.New("入库仓不存在")
	}
	whFrom := models.Warehouse{}
	if err := whFrom.GetByWarehouseCode(e.Orm, d.FromWarehouseCode); err != nil {
		return errors.New("出库仓不存在")
	}
	vendor := models.Vendors{}
	if err := vendor.GetById(e.Orm, d.VendorId); err != nil {
		return errors.New("货主不存在")
	}
	skuSlice := lo.Uniq(utils.Split(strings.ReplaceAll(strings.Trim(d.SkuCodes, " "), "，", ",")))
	// 调用 pc的sku接口
	skuMap, skuSliceRes := models.GetSkuMapFromPcClient(e.Orm, skuSlice)
	skuErrSlice := lo.Without(skuSlice, skuSliceRes...)
	skuGoodsFromErrSlice := []string{}
	skuGoodsFromAuditErrSlice := []string{}
	skuGoodsToErrSlice := []string{}
	skuGoodsToAuditErrSlice := []string{}

	skuGoodsFromMap := map[string]modelsPc.Goods{}
	skuGoodsFromSlice := []string{}
	skuGoodsToMap := map[string]modelsPc.Goods{}
	skuGoodsToSlice := []string{}

	if !whFrom.CheckIsVirtual() {
		skuGoodsFromMap, skuGoodsFromSlice = models.GetGoodsInfoMapByThreeFromPcClient(e.Orm, skuSlice, d.FromWarehouseCode, d.VendorId, 0)
		skuGoodsFromErrSlice = lo.Without(skuSlice, skuGoodsFromSlice...)
		skuGoodsFromErrSlice = lo.Without(skuGoodsFromErrSlice, skuErrSlice...)
		for key, item := range skuGoodsFromMap {
			if item.ApproveStatus != 1 {
				skuGoodsFromAuditErrSlice = append(skuGoodsFromAuditErrSlice, key)
			}
		}
	}
	if !whTo.CheckIsVirtual() {
		skuGoodsToMap, skuGoodsToSlice = models.GetGoodsInfoMapByThreeFromPcClient(e.Orm, skuSlice, d.ToWarehouseCode, d.VendorId, 0)
		skuGoodsToErrSlice = lo.Without(skuSlice, skuGoodsToSlice...)
		skuGoodsToErrSlice = lo.Without(skuGoodsToErrSlice, skuErrSlice...)
		for key, item := range skuGoodsToMap {
			if item.ApproveStatus != 1 {
				skuGoodsToAuditErrSlice = append(skuGoodsToAuditErrSlice, key)
			}
		}
	}
	errStr := FormatValidateSkusErrInfoForTransfer(skuErrSlice, skuGoodsFromErrSlice, skuGoodsFromAuditErrSlice, skuGoodsToErrSlice, skuGoodsToAuditErrSlice)
	if errStr != "" {
		return errors.New(errStr)
	}
	tmpData := map[string]modelsPc.Goods{}
	stockTransferValidateSkusResp := dto.StockTransferValidateSkusResp{}
	if len(skuGoodsFromMap) > 0 {
		tmpData = skuGoodsFromMap
	} else if len(skuGoodsToMap) > 0 {
		tmpData = skuGoodsToMap
	}
	for _, item := range skuSlice {
		stockTransferValidateSkusResp.SkuCode = item
		stockTransferValidateSkusResp.VendorName = vendor.NameZh
		stockTransferValidateSkusResp.ProductName = skuMap[item].NameZh
		stockTransferValidateSkusResp.MfgModel = skuMap[item].MfgModel
		stockTransferValidateSkusResp.BrandName = skuMap[item].Brand.BrandZh
		stockTransferValidateSkusResp.SalesUom = skuMap[item].SalesUom
		stockTransferValidateSkusResp.VendorSkuCode = skuMap[item].SupplierSkuCode
		if goodsInfo, ok := tmpData[item]; ok {
			stockTransferValidateSkusResp.ProductNo = goodsInfo.ProductNo
		}
		*list = append(*list, stockTransferValidateSkusResp)
	}
	return nil
}

func FormatValidateSkusErrInfoForTransfer(skuErrSlice, skuGoodsFromErrSlice, skuGoodsFromAuditErrSlice, skuGoodsToErrSlice, skuGoodsToAuditErrSlice []string) string {
	errSlice := []string{}
	if len(skuErrSlice) > 0 {
		errSku := strings.Join(skuErrSlice, ",")
		errSlice = append(errSlice, "【"+errSku+"】不存在")
	}
	if len(skuGoodsFromErrSlice) > 0 {
		errSku := strings.Join(skuGoodsFromErrSlice, ",")
		errSlice = append(errSlice, "【"+errSku+"】没有对应货主和出库仓的商品关系")
	}
	if len(skuGoodsFromAuditErrSlice) > 0 {
		errSku := strings.Join(skuGoodsFromAuditErrSlice, ",")
		errSlice = append(errSlice, "【"+errSku+"】出库仓的商品关系待审核或审核不通过")
	}
	if len(skuGoodsToErrSlice) > 0 {
		errSku := strings.Join(skuGoodsToErrSlice, ",")
		errSlice = append(errSlice, "【"+errSku+"】没有对应货主和入库仓的商品关系")
	}
	if len(skuGoodsToAuditErrSlice) > 0 {
		errSku := strings.Join(skuGoodsToAuditErrSlice, ",")
		errSlice = append(errSlice, "【"+errSku+"】入库仓的商品关系待审核或审核不通过")
	}
	if len(errSlice) > 0 {
		return strings.Join(errSlice, "，") + "，请联系相关人员处理！"
	}
	return ""
}
