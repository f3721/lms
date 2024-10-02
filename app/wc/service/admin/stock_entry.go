package admin

import (
	"encoding/json"
	"errors"
	"fmt"
	modelsCs "go-admin/app/oc/models"
	modelsPc "go-admin/app/pc/models"
	"go-admin/app/wc/models"
	"go-admin/app/wc/service/admin/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
	dtoStock "go-admin/common/dto/stock/dto"
	"go-admin/common/global"
	"go-admin/common/utils"
	ext "go-admin/config"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StockEntry struct {
	service.Service
}

// 入库单数据权限
func StockEntryPermission(p *actions.DataPermission) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db.Where("stock_entry.warehouse_code in ?", utils.Split(p.AuthorityWarehouseId))
		return db
	}
}

// GetPage 获取StockEntry列表
func (e *StockEntry) GetPage(c *dto.StockEntryGetPageReq, p *actions.DataPermission, outData *[]dto.StockEntryGetPageResp, count *int64) error {
	var err error
	var data models.StockEntry
	var list = &[]models.StockEntry{}
	pcPrefix := global.GetTenantPcDBNameWithDB(e.Orm)

	db := e.Orm.Model(&data).Preload("Warehouse").Preload("LogicWarehouse").Preload("StockEntryProducts").Preload("Supplier").
		Joins("LEFT JOIN stock_entry_products entryProduct ON entryProduct.entry_code = stock_entry.entry_code AND entryProduct.deleted_at IS NULL").
		Joins("LEFT JOIN "+pcPrefix+".product product ON product.sku_code = entryProduct.sku_code AND product.deleted_at IS NULL").
		Joins("LEFT JOIN "+pcPrefix+".goods goods ON goods.id = entryProduct.goods_id AND goods.deleted_at IS NULL").
		Joins("LEFT JOIN vendors vendors ON vendors.id = stock_entry.vendor_id AND vendors.deleted_at IS NULL").
		Select("stock_entry.*,vendors.name_zh").
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSizeNegative(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
			StockEntryPermission(p),
			dtoStock.GenCreatedAtTimeSearch(c.CreatedAtStart, c.CreatedAtEnd, "stock_entry"),
			dtoStock.GenProductSearch(c.Sku, c.ProductNo, c.ProductName, "entryProduct"),
			func(db *gorm.DB) *gorm.DB {
				if c.SourceType != "" {
					db.Where("stock_entry.type = ?", c.SourceType)
				}
				if c.SupplierId != 0 {
					db.Where("stock_entry.supplier_id = ?", c.SupplierId)
				}
				return db
			},
		).Group("stock_entry.entry_code").Order("stock_entry.id DESC").Find(list)
	err = db.Error
	if err != nil {
		e.Log.Errorf("StockEntryService GetPage error:%s \r\n", err)
		return err
	}

	if err := e.FormatPageOutput(list, outData); err != nil {
		e.Log.Errorf("StockEntryService GetPage error:%s \r\n", err)
		return err
	}

	err = e.Orm.Table("(?) as u", db.Limit(-1).Offset(-1)).Count(count).Error
	if err != nil {
		e.Log.Errorf("StockEntryService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

func (e *StockEntry) FormatPageOutput(list *[]models.StockEntry, outData *[]dto.StockEntryGetPageResp) error {
	for _, item := range *list {
		dataItem := dto.StockEntryGetPageResp{}
		if err := utils.CopyDeep(&dataItem, &item); err != nil {
			return err
		}
		dataItem.InitData()
		dataItem.DiffNum = e.GetStockEntryDiffNum(&item)
		dataItem.NameZh = item.Supplier.NameZh
		dataItem.SupplierName = item.Supplier.NameZh
		dataItem.SetStockEntryRulesByStatus()
		if dataItem.Status == models.EntryStatus1 {
			dataItem.UpdatedAt = time.Time{}
		}
		*outData = append(*outData, dataItem)
	}
	return nil
}

func (e *StockEntry) GetStockEntryDiffNum(entry *models.StockEntry) string {
	sum := lo.SumBy(entry.StockEntryProducts, func(item models.StockEntryProducts) int {
		return item.Quantity - item.ActQuantity
	})
	return strconv.Itoa(sum)
}

// Get 获取StockEntry对象
func (e *StockEntry) CheckStockEntry(c *gin.Context, d *dto.CheckStockEntryReq, p *actions.DataPermission) error {
	var model models.StockEntry
	var data = &models.StockEntry{}
	err := e.Orm.Model(&model).Preload("StockEntryProducts").Preload("StockEntryProducts.StockEntryProductsSub").Preload("Warehouse").Preload("LogicWarehouse").
		Scopes(
			actions.Permission(data.TableName(), p),
			StockEntryPermission(p),
		).
		First(data, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetStockEntry error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	//审核不通过
	if d.CheckStatus == models.CheckStatus3 {
		err = e.Orm.Model(&data).Updates(&models.StockEntry{
			CheckStatus: models.CheckStatus3, //失败
			CheckRemark: d.CheckRemark,
			Status:      models.EntryStatus0,
		}).Error
		if err != nil {
			return err
		}
		return nil
	}
	//准备参数
	productS := []*dto.StockEntryProducts{}
	for _, v := range data.StockEntryProducts {
		product := dto.StockEntryProducts{}
		//copier.Copy(&product, &v)
		product.Id = v.Id
		product.SkuCode = v.SkuCode
		product.GoodId = v.GoodsId
		product.Quantity = v.Quantity

		subs := []*dto.StockEntryLocationInfo{}
		for _, sub := range v.StockEntryProductsSub {
			StockEntryLocationInfo := &dto.StockEntryLocationInfo{
				Id:              sub.Id,
				StockLocationId: sub.StashLocationId,  //暂存的id，最终确认的stockLocationId
				ActQuantity:     sub.StashActQuantity, //暂存的数量，最终确认的数量
			}
			subs = append(subs, StockEntryLocationInfo)
		}
		product.StockEntryLocationInfo = subs
		productS = append(productS, &product)
	}

	//入库
	req := dto.StockEntryPartReq{
		Id:                 d.Id,
		EntryType:          2,
		StockEntryProducts: productS,
	}

	req.SetUpdateBy(user.GetUserId(c))
	req.SetUpdateByName(user.GetUserName(c))
	_, err = e.PartEntry(&req, p)
	if err != nil {
		return err
	}
	//成功
	err = e.Orm.Model(&data).Updates(&models.StockEntry{
		CheckStatus: "2", //成功
		EntryTime:   time.Now(),
	}).Error
	if err != nil {
		return err
	}
	return nil
}

// Get 获取StockEntry对象
func (e *StockEntry) Get(d *dto.StockEntryGetReq, p *actions.DataPermission, outData *dto.StockEntryGetResp) error {
	var data models.StockEntry
	var model = &models.StockEntry{}
	err := e.Orm.Model(&data).Preload("StockEntryProducts").Preload("Supplier").Preload("Warehouse").Preload("LogicWarehouse").
		Scopes(
			actions.Permission(data.TableName(), p),
			StockEntryPermission(p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetStockEntry error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}

	if err := e.FormatDetailForEntry(e.Orm, outData, model, d.Type); err != nil {
		return err
	}
	return nil
}

// 入库单格式化输出
func (e *StockEntry) FormatDetailForEntry(tx *gorm.DB, data *dto.StockEntryGetResp, stockEntry *models.StockEntry, Type string) error {
	if err := utils.CopyDeep(data, stockEntry); err != nil {
		return err
	}
	// if err := data.InitData(tx, stockEntry, Type); err != nil {
	// 	return err
	// }

	stockLocationIds := []int{}
	data.WarehouseName = data.Warehouse.WarehouseName
	data.LogicWarehouseName = data.LogicWarehouse.LogicWarehouseName
	data.TypeName = utils.GetTFromMap(data.Type, models.EntryTypeMap)
	data.StatusName = utils.GetTFromMap(data.Status, models.EntryStatusMap)
	data.SourceTypeName = utils.GetTFromMap(data.Type, models.EntrySourceTypeMap)
	// 获取goods product 信息通过pc
	goodsIdSlice := lo.Uniq(lo.Map(stockEntry.StockEntryProducts, func(item models.StockEntryProducts, _ int) int {
		return item.GoodsId
	}))
	goodsProductMap := models.GetGoodsProductInfoFromPc(tx, goodsIdSlice)

	// 获取vendor map
	vendorIdSlice := lo.Uniq(lo.Map(stockEntry.StockEntryProducts, func(item models.StockEntryProducts, _ int) int {
		return item.VendorId
	}))
	vendorIdMap := models.GetVendorsMapByIds(tx, vendorIdSlice)

	// 查询拆分入库记录
	entryProductIds := lo.Map(stockEntry.StockEntryProducts, func(item models.StockEntryProducts, _ int) int {
		return item.Id
	})
	entrySubs := []models.StockEntryProductsSub{}
	err := tx.Where("entry_product_id in ?", entryProductIds).Find(&entrySubs).Error
	if err != nil {
		return err
	}
	entrySubIdMap := map[int][]models.StockEntryProductsSub{}
	for _, sub := range entrySubs {
		entrySubIdMap[sub.EntryProductId] = append(entrySubIdMap[sub.EntryProductId], sub)
	}

	// 查询入库单数量
	sourceProducts, err := e.sourceProductsInfo(stockEntry.Type, stockEntry.SourceCode)
	if err != nil {
		return err
	}
	sourceProductsMap := lo.Associate(sourceProducts, func(item models.SourceProductsInfo) (string, models.SourceProductsInfo) {
		return item.SkuCode, item
	})

	// 商品处理
	for index, item := range stockEntry.StockEntryProducts {
		productGetResp := dto.StockEntryProductsGetResp{}
		if err := utils.CopyDeep(&productGetResp, &item); err != nil {
			return err
		}
		// 次品库位处理
		realLogicWarehouseCode := productGetResp.LogicWarehouseCode
		if productGetResp.CheckIsDefective() {
			defectiveLogicWarehouse := &models.LogicWarehouse{}
			_ = defectiveLogicWarehouse.GetDefectiveOrPassedLogicWarehouse(tx, productGetResp.LogicWarehouseCode, models.LwhType1)
			realLogicWarehouseCode = defectiveLogicWarehouse.LogicWarehouseCode
		}

		// 出库操作时，选择预选择合适的库位
		if Type == "confirm" {
			stockLocations, hasTop, err := models.GetStockLocationsForEntryHasId(tx, realLogicWarehouseCode, productGetResp.GoodsId, productGetResp.StashLocationId)
			if err == nil {
				stockEntryProductsLocationGetResp := dto.StockEntryProductsLocationGetResp{}
				for _, item := range *stockLocations {
					stockEntryProductsLocationGetResp.Regenerate(item)
					productGetResp.StockLocation = append(productGetResp.StockLocation, stockEntryProductsLocationGetResp)
				}
				if hasTop && len(*stockLocations) > 0 {
					productGetResp.StockLocationId = (*stockLocations)[0].Id
				}
			}
		}

		//新加一个，查正品的，库位的StockLocation，查实体仓，再逻辑仓
		if Type == "check" {
			stockLocations, err := models.GetStockLocationByLwhCodeWithLimit(tx, stockEntry.LogicWarehouseCode)
			if err != nil {
				return err
			}
			temp := []dto.StockEntryProductsLocationGetResp{}
			for _, location := range *stockLocations {
				temp = append(temp, dto.StockEntryProductsLocationGetResp{
					Id:           location.Id,
					LocationCode: location.LocationCode,
				})
			}
			productGetResp.StockLocation = temp
		}
		// 打印的时候不知道 为啥需要判断暂存数据 反正后续不需要了。先干掉 如果线上反馈有问题 再改回来
		//if Type == "confirm" || Type == "print" {
		if Type == "confirm" {
			if productGetResp.StashLocationId != 0 {
				productGetResp.StockLocationId = productGetResp.StashLocationId
			}
			if productGetResp.StashActQuantity != 0 {
				productGetResp.ActQuantity = productGetResp.StashActQuantity
			}
		}

		// 要查询的库位id
		stockLocationIds = append(stockLocationIds, productGetResp.StockLocationId)
		// 打印时取暂存库位id对应的库位code
		productGetResp.DiffNum = productGetResp.Quantity - productGetResp.ActQuantity
		productGetResp.Number = index + 1
		productGetResp.VendorName = vendorIdMap[productGetResp.VendorId]
		productGetResp.FillProductGoodsRespData(productGetResp.GoodsId, goodsProductMap)

		// 回显拆分行数据
		productGetResp.StockEntryProductsSub = entrySubIdMap[item.Id]
		// 兼容旧数据
		if data.Status == "2" && len(productGetResp.StockEntryProductsSub) == 0 {
			subInfo := models.StockEntryProductsSub{
				EntryCode:        data.EntryCode,
				EntryProductId:   item.Id,
				StockLocationId:  item.StockLocationId,
				ShouldQuantity:   item.Quantity,
				ActQuantity:      item.ActQuantity,
				StashLocationId:  item.StashLocationId,
				StashActQuantity: item.StashActQuantity,
				EntryTime:        data.UpdatedAt,
			}
			productGetResp.StockEntryProductsSub = []models.StockEntryProductsSub{subInfo}
		}

		// 统计当前SKU出库数量
		outBoundSku := sourceProductsMap[item.SkuCode]
		productGetResp.OutboundNum = outBoundSku.ActQuantity
		data.CheckStatusName = utils.GetTFromMap(data.CheckStatus, models.EntryCheckStatusMap)
		data.SupplierName = stockEntry.Supplier.NameZh

		// 最后赋值
		data.StockEntryProducts = append(data.StockEntryProducts, productGetResp)
	}

	stockLocationMap, err := models.GetStockLocationMapByIds(tx, stockLocationIds)
	if err != nil {
		return err
	}
	for index, item := range data.StockEntryProducts {
		if locationCode, ok := stockLocationMap[item.StockLocationId]; ok {
			data.StockEntryProducts[index].LocationCode = locationCode
		}
	}

	// 差异数量
	data.DiffNum = e.GetStockEntryDiffNum(stockEntry)

	// 补全地址信息
	if data.IsCsOrderEntry() || data.Type == models.EntryType3 {
		data.Address = data.Warehouse.Address
		data.District = data.Warehouse.District
		data.City = data.Warehouse.City
		data.Province = data.Warehouse.Province

		data.DistrictName = data.Warehouse.DistrictName
		data.CityName = data.Warehouse.CityName
		data.ProvinceName = data.Warehouse.ProvinceName
		orderInfo := models.GeOrderInfoFromOcByCsNo(tx, data.SourceCode)
		data.ExternalOrderNo = orderInfo.ContractNo

	} else {
		transfer := &models.StockTransfer{}
		_ = transfer.GetByTransferCode(tx, data.SourceCode)
		data.Address = transfer.Address
		data.District = transfer.District
		data.City = transfer.City
		data.Province = transfer.Province

		data.DistrictName = transfer.DistrictName
		data.CityName = transfer.CityName
		data.ProvinceName = transfer.ProvinceName
	}
	return nil
}

// CreateEntryForCsOrder 创建StockEntry对象（售后单）
func CreateEntryForCsOrder(tx *gorm.DB, c *dto.StockEntryInsertReq) (*models.StockEntry, error) {
	if err := CheckStockEntryInsertReqForCsOrder(tx, c); err != nil {
		return nil, err
	}
	return InsertEntry(tx, c, models.EntryType1)
}

func CheckStockEntryInsertReqForCsOrder(tx *gorm.DB, req *dto.StockEntryInsertReq) error {
	skuTemp := []string{}
	logicWarehouse := &models.LogicWarehouse{}
	if err := logicWarehouse.GetPassLogicWarehouseByWhCode(tx, req.WarehouseCode); err != nil {
		return err
	}
	req.LogicWarehouseCode = logicWarehouse.LogicWarehouseCode
	for _, item := range req.StockEntryProducts {
		if item.Quantity <= 0 {
			return errors.New("商品数量信息异常")
		}
		skuTemp = append(skuTemp, item.SkuCode)
	}
	skuTempLen := len(skuTemp)
	skuTemp = lo.Uniq(skuTemp)
	if skuTempLen > len(skuTemp) {
		return errors.New("重复的sku_code")
	}
	return nil
}

// CancelEntryForCancelCsOrder 取消StockEntry对象（取消售后单） todo
func CancelEntryForCancelCsOrder(tx *gorm.DB, c *dto.StockEntryCancelForCancelCsOrderReq) error {
	var data models.StockEntry
	wcPrefix := global.GetTenantWcDBNameWithDB(tx)

	_ = data.GetByEntrySourceCode(tx, c.CsNo)
	oldData := data
	if data.Id == 0 {
		return errors.New(c.CsNo + ",售后单对应的入库单不存在")
	}
	if data.Status != models.EntryStatus1 {
		return errors.New(c.CsNo + ",售后单对应的入库单状态不正确")
	}
	c.Generate(&data)

	if err := tx.Table(wcPrefix + "." + data.TableName()).Omit(clause.Associations).Save(&data).Error; err != nil {
		return errors.New(c.CsNo + ",售后单对应的入库单更新失败")
	}

	// 入库单日志
	beforeDataStr, _ := json.Marshal(&oldData)
	afterDataStr, _ := json.Marshal(&data)
	dataStr, _ := json.Marshal(c)
	opLog := models.OperateLogs{
		DataId:       data.EntryCode,
		ModelName:    models.EntryLogModelName,
		Type:         models.EntryLogModelTypeCancelForCancelCsOrder,
		Before:       string(beforeDataStr),
		Data:         string(dataStr),
		After:        string(afterDataStr),
		OperatorId:   c.UpdateBy,
		OperatorName: c.UpdateByName,
	}
	_ = opLog.InsertItem(tx)

	return nil
}

// Insert 创建StockEntry对象
func InsertEntry(tx *gorm.DB, c *dto.StockEntryInsertReq, Type string) (*models.StockEntry, error) {
	var data models.StockEntry
	c.Generate(&data)
	if err := data.InsertEntry(tx, Type); err != nil {
		return nil, err
	}
	//检验配置,仓库+订单类型+质检开关,增加质检配置，后面再继续封装
	s := StockEntry{}
	err := s.setQualityTask(tx, &data)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return &data, nil
}

// Insert 创建StockEntry对象
func InsertStockEntry(tx *gorm.DB, c *dto.AddStockEntryReq, EntryCode string) (*models.StockEntry, error) {
	var data models.StockEntry
	if err := c.Generate(tx, &data); err != nil {
		return nil, err
	}
	if err := data.InsertStockEntry(tx, EntryCode); err != nil {
		return nil, err
	}
	return &data, nil
}

// Update 修改StockEntry对象
func (e *StockEntry) Update(c *dto.StockEntryUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.StockEntry{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
		StockEntryPermission(p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("StockEntryService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除StockEntry
func (e *StockEntry) Remove(d *dto.StockEntryDeleteReq, p *actions.DataPermission) error {
	var data models.StockEntry

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveStockEntry error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}

// SubCreateOrUpdate 入库操作批量新建或更新库位信息
//
// tx: 事务db params: 入库库位列表
func (e *StockEntry) SubCreateOrUpdate(tx *gorm.DB, params []*models.StockEntryProductsSub) error {
	updateArr := []*models.StockEntryProductsSub{}
	createArr := []*models.StockEntryProductsSub{}
	for _, item := range params {
		if item.ActQuantity == 0 && item.StashActQuantity == 0 { // 如果全部0值，兜底跳过
			continue
		}
		// 区分更新还是新建
		if item.Id != 0 {
			updateArr = append(updateArr, item)
		} else {
			createArr = append(createArr, item)
		}
	}

	// 更新数据
	if len(updateArr) > 0 {
		for _, updateData := range updateArr {
			err := tx.Model(&models.StockEntryProductsSub{}).Where("id = ?", updateData.Id).Updates(&updateData).Error
			if err != nil {
				return err
			}
		}
	}

	// 新增数据
	if len(createArr) > 0 {
		err := tx.Model(&models.StockEntryProductsSub{}).Create(&createArr).Error
		if err != nil {
			return err
		}
	}

	return nil
}

// PartEntry 部分入库 | 目前只支持 大货入库
func (e *StockEntry) PartEntry(c *dto.StockEntryPartReq, p *actions.DataPermission) (html any, err error) {
	currTime := time.Now()

	// 查询入库单
	stockEntry := &models.StockEntry{}
	err = e.Orm.Preload("StockEntryProducts").Preload("Warehouse").Where("id = ?", c.Id).
		Scopes(
			actions.Permission(stockEntry.TableName(), p),
			StockEntryPermission(p),
		).First(stockEntry).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", err
	}
	if err == gorm.ErrRecordNotFound {
		return "", errors.New("入库单不存在或没有权限")
	}

	//对质检进行判断
	ss := QualityCheck{}
	ss.Orm = e.Orm
	err = ss.IsQualityNumOk(c, stockEntry)
	if err != nil {
		return nil, err
	}

	// 退货入库强制整单入库
	if stockEntry.IsCsOrderEntry() {
		c.EntryType = 2
	}

	// 校验入库单状态
	if lo.IndexOf([]string{"0", "2"}, stockEntry.Status) != -1 {
		return "", errors.New("入库单状态已作废或已完成")
	}
	if stockEntry.Warehouse.CheckIsVirtual() {
		return "", errors.New("虚拟仓无需手动入库")
	}

	oldData := stockEntry

	// 如果产品中有次品，校验配置了次品仓
	defectiveLogicWarehouse := &models.LogicWarehouse{}
	hasDefectiveList := lo.Filter(stockEntry.StockEntryProducts, func(x models.StockEntryProducts, _ int) bool {
		return x.CheckIsDefective()
	})
	if len(hasDefectiveList) > 0 {
		err = defectiveLogicWarehouse.GetDefectiveOrPassedLogicWarehouse(e.Orm, stockEntry.LogicWarehouseCode, models.LwhType1)
		if err != nil {
			return "", errors.New("入库单逻辑仓对应次品仓获取失败")
		}
	}

	// 查询当前实体仓下所有库位 | 校验用
	stockLocations := []*models.StockLocation{}
	err = e.Orm.Where("warehouse_code = ?", stockEntry.WarehouseCode).Where("status = ?", 1).Find(&stockLocations).Error
	if err != nil {
		return "", err
	}
	stockLocationArr := lo.Map(stockLocations, func(item *models.StockLocation, _ int) string {
		return strconv.Itoa(item.Id) + item.LogicWarehouseCode
	})

	// 入库单商品Map
	orgEntryProductsMap := lo.Associate(stockEntry.StockEntryProducts, func(item models.StockEntryProducts) (int, models.StockEntryProducts) {
		return item.Id, item
	})

	//检验是否需要质检订单
	num := e.checkQuality()
	if num > 0 {
		return "", fmt.Errorf("还有质检任务和数量需要完成质检：[%v]", num)
	}

	// 校验来源单数据 | 大货和订单
	if lo.IndexOf([]string{"0", "1"}, stockEntry.Type) != -1 {
		// 查询出库单或者退货单的商品Map
		sourceProducts, err := e.sourceProductsInfo(stockEntry.Type, stockEntry.SourceCode)
		if err != nil {
			return "", err
		}
		sourceProductsMap := lo.Associate(sourceProducts, func(item models.SourceProductsInfo) (string, models.SourceProductsInfo) {
			return item.SkuCode, item
		})
		for _, item := range c.StockEntryProducts {
			orgProduct := orgEntryProductsMap[item.Id]
			// 校验当前SKU的 入库数量总和是否 <= 来源SKU出库数量
			sourceSkuInfo := sourceProductsMap[item.SkuCode]
			hisAndCurrTotal := item.ActQuantityTotal + orgProduct.ActQuantity
			if hisAndCurrTotal > sourceSkuInfo.ActQuantity {
				return "", fmt.Errorf("%v,累计入库数量[%v]不能大于出库单已出库数量[%v]", item.SkuCode, hisAndCurrTotal, sourceSkuInfo.ActQuantity)
			}
		}
	}

	// 统计剩余应入库总数
	quantityArr := lo.Map(stockEntry.StockEntryProducts, func(item models.StockEntryProducts, _ int) int {
		return item.Quantity - item.ActQuantity
	})
	totalLeftQuantity := lo.Sum(quantityArr)

	// 全部入库校验参数中产品数量
	if c.EntryType == 2 && len(stockEntry.StockEntryProducts) != len(c.StockEntryProducts) {
		return "", errors.New("全部入库时,产品信息须默认全选")
	}

	// 传参规范校验
	if len(c.StockEntryProducts) == 0 {
		return "", errors.New("传参错误:商品库位信息为空[1],请检查! ")
	}
	goodIdsArr := lo.Map(c.StockEntryProducts, func(item *dto.StockEntryProducts, _ int) int {
		return item.Id
	})
	haskRepeats := utils.HasDuplicate(goodIdsArr)
	if haskRepeats {
		return "", errors.New("传参错误:商品列表中包含重复ID元素,请检查！")
	}

	// 校验商品参数&判断是否已经全部入库
	var allStored bool
	var currTotalActQuantity int
	for _, item := range c.StockEntryProducts {

		if len(item.StockEntryLocationInfo) > 10 {
			return "", errors.New("单个SKU,拆分行数量不能大于10,请检查！")
		}

		orgProduct := orgEntryProductsMap[item.Id]
		if orgProduct.Id == 0 {
			return "", errors.New("商品行ID传参错误,请检查！")
		}

		// 预埋商品ID
		item.GoodId = orgProduct.GoodsId

		// 退单入库不能拆分行
		if stockEntry.IsCsOrderEntry() && len(item.StockEntryLocationInfo) > 1 {
			return "", errors.New(item.SkuCode + ",订单入库, 不能拆分行入库")
		}

		// 传参规范校验
		if len(item.StockEntryLocationInfo) == 0 {
			return "", errors.New("传参错误:商品库位信息为空[2],请检查！")
		}

		// 校验拆分行参数
		item.ActQuantityTotal = 0
		for index, locationInfo := range item.StockEntryLocationInfo {

			// 校验实际入库数量应为大于0的整数 | 点击全部入库时支持行填写实际入库数量为0
			if locationInfo.ActQuantity <= 0 && c.EntryType != 2 {
				msg := fmt.Sprintf("%v,实际入库数量应为大于0的整数", item.SkuCode)
				if len(item.StockEntryLocationInfo) > 1 {
					msg = fmt.Sprintf("%v的拆分入库第%v行,实际入库数量应为大于0的整数", (index + 1), item.SkuCode)
				}
				return "", errors.New(msg)
			}

			// 校验库位编码必填
			if locationInfo.StockLocationId <= 0 {
				return "", errors.New(item.SkuCode + ",缺少库位编码")
			}

			// 校验库位信息是否正确
			logicWarehouseCodeNow := stockEntry.LogicWarehouseCode
			if orgProduct.CheckIsDefective() {
				logicWarehouseCodeNow = defectiveLogicWarehouse.LogicWarehouseCode
			}
			checkLocationStr := strconv.Itoa(locationInfo.StockLocationId) + logicWarehouseCodeNow
			if lo.IndexOf(stockLocationArr, checkLocationStr) == -1 {
				return "", errors.New(item.SkuCode + ",所选的库位不在当前逻辑仓库中")
			}

			// // 校验拆分库位不能在同一个库位
			// if tmpLocationId == locationInfo.StockLocationId {
			// 	return errors.New(item.SkuCode + ",当次拆分入库选择了相同库位，请检查！")
			// }
			// tmpLocationId = locationInfo.StockLocationId

			// 累加汇总实际入库数量
			item.ActQuantityTotal += locationInfo.ActQuantity
		}

		// 实际入库数量不能大于应入库数量
		item.Quantity = orgProduct.Quantity
		if item.ActQuantityTotal > (orgProduct.Quantity - orgProduct.ActQuantity) { // 当前实际入库数量 VS 剩余应入库数量(应入-已入)
			return "", errors.New(item.SkuCode + ",实际入库数量不能大于应入库数量或剩余应入库数量")
		}

		// 售后单的实际入库数量要等于应入库数量
		if stockEntry.IsCsOrderEntry() && item.ActQuantityTotal != (orgProduct.Quantity-orgProduct.ActQuantity) {
			return "", errors.New(item.SkuCode + ",售后单的实际入库数量要等于应入库数量")
		}

		// 统计本次入库之后的实际入库总数
		currTotalActQuantity += item.ActQuantityTotal

		// 校验当前SKU的 入库数量总和是否 <= 来源SKU出库数量
		//sourceSkuInfo := sourceProductsMap[item.SkuCode]
		//hisAndCurrTotal := item.ActQuantityTotal + orgProduct.ActQuantity
		//if hisAndCurrTotal > sourceSkuInfo.ActQuantity {
		//	return "", fmt.Errorf("%v,累计入库数量[%v]不能大于出库单已出库数量[%v]", item.SkuCode, hisAndCurrTotal, sourceSkuInfo.ActQuantity)
		//}
	}

	// 系统判断是否已经全部入库
	sysJudgmentAllStored := totalLeftQuantity == currTotalActQuantity
	if c.EntryType == 2 || sysJudgmentAllStored {
		allStored = true
	}

	// 事务落库
	err = e.Orm.Transaction(func(tx *gorm.DB) error {

		// -------- 保存处理 --------

		if c.EntryType == 0 {
			// 记录库位暂存信息
			var locationInfoArr []*models.StockEntryProductsSub
			for _, reqProduct := range c.StockEntryProducts {
				for _, item := range reqProduct.StockEntryLocationInfo {
					sub := models.StockEntryProductsSub{}
					sub.EntryCode = stockEntry.EntryCode       // 入库单号
					sub.EntryProductId = reqProduct.Id         // 商品ID
					sub.ShouldQuantity = reqProduct.Quantity   // 应入库总量
					sub.StashLocationId = item.StockLocationId // 暂存库位
					sub.StashActQuantity = item.ActQuantity    // 暂存入库数量
					if item.Id != 0 {
						sub.Id = item.Id
					}
					locationInfoArr = append(locationInfoArr, &sub)
				}
			}

			err := e.SubCreateOrUpdate(tx, locationInfoArr)
			if err != nil {
				return err
			}

			// 生成打印Html
			html, err = e.PrintEntryHtml(tx, stockEntry, c.StockEntryProducts, stockLocations)
			if err != nil {
				return err
			}

			// 结束执行
			return nil
		}

		// -------- 入库处理 --------

		// 大货入库更新调拨单状态 | 部分出库未入库，部分出库部分入库，全部出库部分入库
		if stockEntry.IsTransferEntry() {
			sourceOrder := models.StockTransfer{}
			err := sourceOrder.GetByTransferCode(tx, stockEntry.SourceCode)
			if err != nil {
				return err
			}

			// 调拨单状态校验 | 合法状态 : 4-部分出库未入库 5-部分出库部分入库 6-全部出库部分入库 2-已出库未入库
			if sourceOrder.Status == "3" {
				return errors.New("调拨类型已出库完成入库完成")
			}
			if lo.IndexOf([]string{"4", "5", "6", "2"}, sourceOrder.Status) == -1 {
				return errors.New("调拨类型的入库单需先出库再入库")
			}

			// 是否完成出库
			outboundCompleted := false
			if lo.IndexOf([]string{"2", "6"}, sourceOrder.Status) != -1 {
				outboundCompleted = true
			}

			// 用户选择全部入库 | 未全部出库
			if c.EntryType == 2 && !outboundCompleted {
				return errors.New("调拨类型 | [部分出库未入库]和[部分出库部分入库]状态下，不能全部入库")
			}

			// 系统判断全部入库 | 未全部出库
			if allStored && !outboundCompleted {
				return errors.New("调拨类型 | 系统检测到，未全部出库，但填写的入库数量，已经全部入库，请检查")
			}

			// // 调拨单未全部出库，逻辑判断全部入库，判断为部分入库
			// if !outboundCompleted {
			// 	allStored = false
			// }

			// 更新状态 | 5-部分出库部分入库 6-全部出库部分入库 2-已出库未入库 3-出库完成入库完成
			newStatus := sourceOrder.Status
			if !outboundCompleted && !allStored { // 未完成出库 + 部分入库 = 5-部分出库部分入库
				newStatus = "5"
			}
			if outboundCompleted && !allStored { // 已完成出库 + 部分入库 = 6-全部出库部分入库
				newStatus = "6"
			}
			if outboundCompleted && allStored { // 已完成出库 + 全部入库 = 3-出库完成入库完成
				newStatus = "3"
			}

			err = tx.Model(&models.StockTransfer{}).Where("id = ?", sourceOrder.Id).
				Where("status = ?", sourceOrder.Status).UpdateColumn("status", newStatus).Error
			if err != nil {
				return err
			}
		}

		// 售后入库更新退货单状态
		if stockEntry.IsCsOrderEntry() {
			csApplyModel := &modelsCs.CsApply{}
			if err := csApplyModel.ReturnCompleted(tx, stockEntry.SourceCode); err != nil {
				return err
			}
		}

		// 维护入库单状态, 最新入库时间, 完成入库时间
		if allStored {
			stockEntry.Status = "2"
			stockEntry.EntryEndTime = currTime
		} else {
			stockEntry.Status = "3"
		}
		if stockEntry.EntryTime.IsZero() {
			stockEntry.EntryTime = currTime
		}

		err = tx.Where("id = ?", stockEntry.Id).Save(&stockEntry).Error
		if err != nil {
			return err
		}

		// 维护出库单商品信息
		stockEntryDefectiveProducts := make([]*models.StockEntryDefectiveProduct, 0)
		var locationInfoArr []*models.StockEntryProductsSub // 插入库位的操作记录集合
		for _, reqProduct := range c.StockEntryProducts {
			orgProduct := orgEntryProductsMap[reqProduct.Id]

			// 维护入库商品实际入库数量，状态，最新入库时间，入库完成时间
			productStatus := 0
			if !allStored {
				productStatus = 1
			} else {
				productStatus = 2
			}
			productUpdate := map[string]any{
				"status":       productStatus,
				"act_quantity": gorm.Expr("act_quantity + ?", reqProduct.ActQuantityTotal),
				"entry_time":   currTime,
			}
			if !allStored {
				productUpdate["entry_end_time"] = currTime
			}

			err = tx.Model(&models.StockEntryProducts{}).
				Where("id = ?", orgProduct.Id).
				Where("act_quantity = ?", orgProduct.ActQuantity).
				Updates(&productUpdate).Error
			if err != nil {
				return err
			}

			// 增加商品库存 | 全部入库会有0值
			if reqProduct.ActQuantityTotal > 0 {
				stockInfo := &models.StockInfo{}
				if err := stockInfo.GetByGoodsIdAndLogicWarehouseCode(tx, orgProduct.GoodsId, orgProduct.LogicWarehouseCode); err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						stockInfo.SetBaseInfo(orgProduct.GoodsId, orgProduct.VendorId,
							stockEntry.UpdateBy, orgProduct.SkuCode, orgProduct.WarehouseCode,
							orgProduct.LogicWarehouseCode, stockEntry.UpdateByName)
					} else {
						return err
					}
				}
				err = stockInfo.AddStock(tx, reqProduct.ActQuantityTotal, &models.StockLog{
					DocketCode:   stockEntry.EntryCode,
					FromType:     models.StockLogFromType1,
					CreateBy:     stockEntry.UpdateBy,
					CreateByName: stockEntry.UpdateByName,
				})
				if err != nil {
					return err
				}
			}

			// 库位信息维护
			for _, reqLocationInfo := range reqProduct.StockEntryLocationInfo {

				// 全部入库会有0值,不进行库位记录
				if reqLocationInfo.ActQuantity <= 0 {
					continue
				}

				// 次品退货入库商品 这里的库位是次品仓库位 需要自动选择入库库位
				if orgProduct.CheckIsDefective() {
					// 处理次品入库的库位逻辑
					// 1，用户选择和传值的是次品库位；2，入库的是[系统选择的正品库位]；3，自动生成的次品出入库单是[系统选择的正品库位 -> 用户选择的次品库位]

					// A.查询真实库位&赋值调拨入库库位
					orgProduct.StockLocationId = reqLocationInfo.StockLocationId
					defectiveProduct, err := orgProduct.DefectiveHandle(tx)
					if err != nil {
						return err
					}

					// B.赋值真实库位数量
					defectiveProduct.ActQuantity = reqLocationInfo.ActQuantity
					stockEntryDefectiveProducts = append(stockEntryDefectiveProducts, defectiveProduct)

					// C.替换系统选择的入库库位
					reqLocationInfo.StockLocationId = defectiveProduct.PassedStockLocationId
				}

				// 增加库位库存
				locationGoods := &models.StockLocationGoods{}
				if err := locationGoods.GetByGoodsIdAndLocationId(tx, orgProduct.GoodsId, reqLocationInfo.StockLocationId); err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						locationGoods.SetBaseInfo(orgProduct.GoodsId, reqLocationInfo.StockLocationId, stockEntry.UpdateBy, stockEntry.UpdateByName)
					} else {
						return err
					}
				}
				if err := locationGoods.AddStock(tx, reqLocationInfo.ActQuantity, &models.StockLocationGoodsLog{
					DocketCode:   stockEntry.EntryCode,
					FromType:     models.StockLocationGoodsLogFromType1,
					CreateBy:     stockEntry.UpdateBy,
					CreateByName: stockEntry.UpdateByName,
				}); err != nil {
					return err
				}

				// 记录拆分库位入库记录|回显使用
				sub := &models.StockEntryProductsSub{}
				sub.EntryCode = stockEntry.EntryCode
				sub.EntryProductId = reqProduct.Id
				sub.ShouldQuantity = reqProduct.Quantity
				sub.StockLocationId = reqLocationInfo.StockLocationId // 实际库位
				sub.ActQuantity = reqLocationInfo.ActQuantity         // 实际数量
				sub.EntryTime = currTime                              // 入库时间
				if reqLocationInfo.Id != 0 {
					sub.Id = reqLocationInfo.Id
				}
				locationInfoArr = append(locationInfoArr, sub)
			}
		}
		// 批量插入入库库位操作记录
		err := e.SubCreateOrUpdate(tx, locationInfoArr)
		if err != nil {
			return err
		}

		// 创建调拨单For次品入库
		if len(stockEntryDefectiveProducts) > 0 {
			stockTransfer := &models.StockTransfer{}
			err := stockTransfer.InsertTransferForDefectiveEntry(tx, stockEntry, stockEntryDefectiveProducts)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	// 通知货主入库完成（外部调拨单） todo

	// 记录日志
	beforeDataStr, _ := json.Marshal(&oldData)
	afterDataStr, _ := json.Marshal(&stockEntry)
	dataStr, _ := json.Marshal(c)
	opLog := models.OperateLogs{
		DataId:       stockEntry.EntryCode,
		ModelName:    models.EntryLogModelName,
		Type:         models.EntryLogModelTypeConfirm,
		Before:       string(beforeDataStr),
		Data:         string(dataStr),
		After:        string(afterDataStr),
		OperatorId:   c.UpdateBy,
		OperatorName: c.UpdateByName,
	}
	_ = opLog.InsertItem(e.Orm)

	return html, nil
}

// 查看是否有需要检查的配置
func (e *StockEntry) checkQuality() int {
	return 0
}

// ConfirmEntry 确认入库 - 逐步废弃
func (e *StockEntry) ConfirmEntry(c *dto.StockEntryConfirmReq, p *actions.DataPermission) error {
	var data = models.StockEntry{}
	e.Orm.Preload("StockEntryProducts").Preload("Warehouse").
		Scopes(
			actions.Permission(data.TableName(), p),
			StockEntryPermission(p),
		).First(&data, c.GetId())
	if data.Id == 0 {
		return errors.New("入库单不存在或没有权限")
	}
	oldData := data
	if err := e.CheckConfirmEntry(e.Orm, &data); err != nil {
		return err
	}
	c.Generate(&data)

	if err := e.CheckConfirmEntryProducts(e.Orm, c.StockEntryProducts, &data); err != nil {
		return err
	}

	// 开启事务
	err := e.Orm.Transaction(func(tx *gorm.DB) error {
		return data.ConfirmEntryManual(tx)
	})
	if err != nil {
		e.Log.Errorf("StockEntryService ConfirmEntry error:%s \r\n", err)
		return err
	}
	// 通知货主入库完成（外部调拨单） todo

	// 记录日志
	beforeDataStr, _ := json.Marshal(&oldData)
	afterDataStr, _ := json.Marshal(&data)
	dataStr, _ := json.Marshal(c)
	opLog := models.OperateLogs{
		DataId:       data.EntryCode,
		ModelName:    models.EntryLogModelName,
		Type:         models.EntryLogModelTypeConfirm,
		Before:       string(beforeDataStr),
		Data:         string(dataStr),
		After:        string(afterDataStr),
		OperatorId:   c.UpdateBy,
		OperatorName: c.UpdateByName,
	}
	_ = opLog.InsertItem(e.Orm)
	return nil
}

func (e *StockEntry) AddStockEntry(db *gorm.DB, data *dto.AddStockEntryReq, p *actions.DataPermission) error {
	//参数校验
	err := e.CheckProductsForStockEntry(e.Orm, data, p)
	if err != nil {
		return err
		//return errors.New("参数校验失败，请检查重试!")
	}
	//新增采购入库
	tx := db.Begin()
	res, err := InsertStockEntry(tx, data, "")
	if err != nil {
		tx.Rollback()
		return err
	}

	//检验配置,仓库+订单类型+质检开关,增加质检配置，后面再继续封装
	err = e.setQualityTask(tx, res)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	//错误，需要回滚

	//记录日志
	afterDataStr, _ := json.Marshal(&data)
	dataStr, _ := json.Marshal(data)
	opLog := models.OperateLogs{
		DataId:       res.EntryCode,
		ModelName:    models.EntryLogModelName,
		Type:         models.EntryLogModelTypeConfirm,
		Data:         string(dataStr),
		After:        string(afterDataStr),
		OperatorId:   data.UpdateBy,
		OperatorName: data.UpdateByName,
	}
	_ = opLog.InsertItem(e.Orm)
	return err
}

func (e *StockEntry) setQualityTask(tx *gorm.DB, model *models.StockEntry) error {
	var QualityCheckNumProducts []*dto.QualityCheckNumProducts
	for _, product := range model.StockEntryProducts {
		item := &dto.QualityCheckNumProducts{
			SkuCode:  product.SkuCode,
			Quantity: product.Quantity,
		}
		QualityCheckNumProducts = append(QualityCheckNumProducts, item)
	}

	req := &dto.QualityCheckNumReq{
		WarehouseCode:           model.WarehouseCode,
		OrderType:               model.Type,
		QualityCheckNumProducts: QualityCheckNumProducts,
	}

	quaService := QualityCheckConfig{}
	quaService.Orm = e.Orm
	qualityRes, err := quaService.QualityCheckNum(req)
	if err != nil {
		return err
	}
	//如果无参数，结果无需配置
	if len(qualityRes) > 0 {
		curDate := time.Now().Format("20060102")
		var QualityItem []models.QualityCheck
		for _, qs := range qualityRes {
			Item := models.QualityCheck{
				SourceCode:         model.SourceCode,
				EntryCode:          model.EntryCode,
				WarehouseCode:      model.WarehouseCode,
				LogicWarehouseCode: model.LogicWarehouseCode,
				StayQualityNum:     qs.Quantity,
				Type:               qs.Type,
				SkuCode:            qs.SkuCode,
				QualityCheckCode:   fmt.Sprintf("QC%v", curDate),
				SourceName:         model.Supplier.NameZh,
			}
			Item.CreateByName = model.CreateByName
			if model.Type == models.EntryType0 { //后面查询
				Item.SourceName = ""
			}
			if model.Type == models.EntryType3 {
				Item.SourceName = model.Supplier.NameZh
			}
			Item.CreateBy = model.CreateBy
			Item.Unqualified = qs.Unqualified
			QualityItem = append(QualityItem, Item)
		}
		err = tx.Model(&models.QualityCheck{}).Debug().Omit("id").Create(&QualityItem).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *StockEntry) CheckProductsForStockEntry(tx *gorm.DB, data *dto.AddStockEntryReq, p *actions.DataPermission) error {
	if lo.IndexOf(utils.Split(p.AuthorityWarehouseId), data.WarehouseCode) == -1 {
		return errors.New("没有实体仓权限")
	}
	if lo.IndexOf(utils.SplitToInt(p.AuthorityVendorId), data.VendorId) == -1 {
		return errors.New("没有货主权限")
	}

	if len(data.AddStockEntryProducts) == 0 {
		return errors.New("商品明细不能为空")
	}

	vendorMap := models.GetVendorsMapByIds(e.Orm, []int{data.VendorId})
	if len(vendorMap) == 0 {
		return errors.New(fmt.Sprintf("%d，vender信息不存在", data.VendorId))
	}

	//有效的supplierId
	err := e.Orm.Where("status = 1").First(&models.Supplier{}, data.SupplierId).Error
	if err != nil {
		return errors.New(fmt.Sprintf("%d，供应商信息不存在", data.SupplierId))
	}

	tempSku := []string{}
	tempGoods := []string{} //goodsId+库位id校验，不能重复

	//skuGoodsToMap := map[string]modelsPc.Goods{}
	skuGoodsToSlice := []string{}
	for _, item := range data.AddStockEntryProducts {
		if item.SkuCode == "" {
			return errors.New("缺少SKU")
		}
		if existSkuFlag := lo.Contains(tempSku, item.SkuCode); existSkuFlag {
			return errors.New(item.SkuCode + "，重复的SKU")
		}

		tempSku = append(tempSku, item.SkuCode)
		//if item.Quantity <= 0 {
		//	return errors.New("入库数量为大于0的整数")
		//}
		//每次查询一个校验
		goods := models.GetGoodsProductInfoFromPc(e.Orm, []int{item.GoodsId})
		if len(goods) == 0 {
			return errors.New("goods数据错误")
		}

		//skuGoodsToMap, skuGoodsToSlice = models.GetGoodsInfoMapByThreeFromPcClient(tx, tempSku, data.WarehouseCode, data.VendorId, 1)
		_, skuGoodsToSlice = models.GetGoodsInfoMapByThreeFromPcClient(tx, tempSku, data.WarehouseCode, data.VendorId, 1)
		var skuGoodsToErrSlice []string = lo.Without(tempSku, skuGoodsToSlice...)
		if len(skuGoodsToErrSlice) != 0 {
			return errors.New("入库仓的商品关系存在异常")
		}

		for _, subInfo := range item.StockEntryProductsSub {
			if subInfo.StashLocationId == 0 {
				return errors.New("请输入暂存入库数量id")
			}
			err := e.Orm.First(&models.StockLocation{}, subInfo.StashLocationId).Error
			if err != nil {
				return errors.New("请检查库位Id")
			}

			if subInfo.StashActQuantity <= 0 {
				return errors.New("暂存入库数量为大于0的整数")
			}
			if existGoodsFlag := lo.Contains(tempGoods, fmt.Sprintf("%d-%d", item.GoodsId, subInfo.StashLocationId)); existGoodsFlag {
				return errors.New(fmt.Sprintf("商品和库位Id重复：%v", item.SkuCode))
			}
			tempGoods = append(tempGoods, fmt.Sprintf("%d-%d", item.GoodsId, subInfo.StashLocationId))
		}
	}
	return nil
}

// 先删除后新增方案
func (e *StockEntry) EditStockEntry(db *gorm.DB, c *dto.AddStockEntryReq, p *actions.DataPermission) error {
	//1.先查询
	//2.修改
	//3.保存
	if c.Id == 0 {
		return errors.New("请输入Id!")
	}
	var err error
	var skEntry = models.StockEntry{}
	err = e.Orm.Preload("StockEntryProducts").Preload("StockEntryProducts.StockEntryProductSub").Where("id = ?", c.Id).First(&skEntry).Error
	if err != nil {
		return errors.New("检查信息失败!")
	} else {
		if skEntry.Id < 1 {
			return errors.New("请检查入库单信息!")
		}
	}

	//参数校验
	err = e.CheckProductsForStockEntry(e.Orm, c, p)
	if err != nil {
		return errors.New("参数校验失败，请检查后重试!")
	}

	tx := db.Begin()
	//删除，加事务
	err = tx.Unscoped().Delete(&skEntry, c.Id).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	if skEntry.EntryCode != "" {
		err = tx.Unscoped().Where("entry_code = ?", skEntry.EntryCode).Delete(&models.StockEntryProducts{}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if skEntry.EntryCode != "" {
		err = tx.Unscoped().Where("entry_code = ?", skEntry.EntryCode).Delete(&models.StockEntryProductsSub{}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	c.Id = 0
	res, err := InsertStockEntry(tx, c, skEntry.EntryCode)
	if err != nil {
		tx.Rollback()
		return err
	}

	//检验配置,仓库+订单类型+质检开关,增加质检配置，后面再继续封装
	err = e.setQualityTask(tx, res)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	//日志
	data := skEntry
	afterDataStr, _ := json.Marshal(&data)
	dataStr, _ := json.Marshal(data)
	opLog := models.OperateLogs{
		DataId:       skEntry.EntryCode,
		ModelName:    models.EntryLogModelName,
		Type:         models.EntryLogModelTypeConfirm,
		Data:         string(dataStr),
		After:        string(afterDataStr),
		OperatorId:   data.UpdateBy,
		OperatorName: data.UpdateByName,
	}
	_ = opLog.InsertItem(e.Orm)
	return err
}

// CheckConfirmEntry 确认出库check
func (e *StockEntry) CheckConfirmEntry(tx *gorm.DB, data *models.StockEntry) error {
	if err := data.CheckBeforeConfirmEntry(tx); err != nil {
		return err
	}
	return nil
}

// CheckConfirmEntry 确认出库check
func (e *StockEntry) CheckConfirmEntryProducts(tx *gorm.DB, reqList []dto.StockEntryProductsConfirmReq, stockEntry *models.StockEntry) error {
	// 获取入库单逻辑仓对应次品仓
	defectiveLogicWarehouse := &models.LogicWarehouse{}
	var hasDefectiveList []models.StockEntryProducts = lo.Filter[models.StockEntryProducts](stockEntry.StockEntryProducts, func(x models.StockEntryProducts, _ int) bool {
		return x.CheckIsDefective()
	})
	if len(hasDefectiveList) > 0 {
		if err := defectiveLogicWarehouse.GetDefectiveOrPassedLogicWarehouse(tx, stockEntry.LogicWarehouseCode, models.LwhType1); err != nil {
			return errors.New("入库单逻辑仓对应次品仓获取失败")
		}
	}
	for index, item := range stockEntry.StockEntryProducts {
		reqItem := FindInProductsConfirmReqListForEntry(reqList, item.SkuCode)
		if reqItem == nil {
			return errors.New("参数异常,缺少" + item.SkuCode + "信息")
		}
		if reqItem.StockLocationId <= 0 {
			return errors.New(reqItem.SkuCode + ",缺少库位编码")
		}
		if reqItem.ActQuantity <= 0 {
			return errors.New(reqItem.SkuCode + ",实际入库数量应为为大于0的整数")
		}
		if item.Quantity < reqItem.ActQuantity {
			return errors.New(reqItem.SkuCode + ",实际入库数量不能大于应入库数量")
		}
		if stockEntry.IsCsOrderEntry() && item.Quantity != reqItem.ActQuantity {
			return errors.New(reqItem.SkuCode + ",售后单的实际入库数量要等于应入库数量")
		}
		// 库位信息检查
		logicWarehouseCodeNow := stockEntry.LogicWarehouseCode
		if item.CheckIsDefective() {
			logicWarehouseCodeNow = defectiveLogicWarehouse.LogicWarehouseCode
		}
		stockLocation := &models.StockLocation{}
		if !stockLocation.CheckStockLocation(tx, reqItem.StockLocationId, logicWarehouseCodeNow) {
			return errors.New(reqItem.SkuCode + ",所选的库位不在当前逻辑仓库中")
		}
		// 赋值实际入库数量
		stockEntry.StockEntryProducts[index].ActQuantity = reqItem.ActQuantity
		stockEntry.StockEntryProducts[index].StockLocationId = reqItem.StockLocationId
	}
	return nil
}

func FindInProductsConfirmReqListForEntry(list []dto.StockEntryProductsConfirmReq, sku string) *dto.StockEntryProductsConfirmReq {
	for _, item := range list {
		if item.SkuCode == sku {
			return &item
		}
	}
	return nil
}

func (e *StockEntry) Stash(c *dto.StockEntryStashReq, p *actions.DataPermission) error {
	var data = models.StockEntry{}
	e.Orm.Preload("StockEntryProducts").Preload("Warehouse").
		Scopes(
			actions.Permission(data.TableName(), p),
			StockEntryPermission(p),
		).First(&data, c.GetId())
	if data.Id == 0 {
		return errors.New("入库单不存在或没有权限")
	}
	oldData := data
	c.Generate(&data)

	if data.Status != models.EntryStatus1 {
		return errors.New("入库单状态不正确")
	}

	if err := e.CheckConfirmEntryProducts(e.Orm, c.StockEntryProducts, &data); err != nil {
		return err
	}

	// 开启事务
	err := e.Orm.Transaction(func(tx *gorm.DB) error {
		return data.Stash(e.Orm)
	})
	if err != nil {
		e.Log.Errorf("StockEntryService Stash error:%s \r\n", err)
		return err
	}

	beforeDataStr, _ := json.Marshal(&oldData)
	afterDataStr, _ := json.Marshal(&data)
	dataStr, _ := json.Marshal(c)
	opLog := models.OperateLogs{
		DataId:       data.EntryCode,
		ModelName:    models.EntryLogModelName,
		Type:         models.EntryLogModelTypeStash,
		Before:       string(beforeDataStr),
		Data:         string(dataStr),
		After:        string(afterDataStr),
		OperatorId:   c.UpdateBy,
		OperatorName: c.UpdateByName,
	}
	_ = opLog.InsertItem(e.Orm)
	return nil
}

func (e *StockEntry) StockLocationsSelect(c *dto.StockEntryProductsLocationGetReq, p *actions.DataPermission, outData *[]dto.StockEntryProductsLocationGetResp) error {
	var err error

	realLogicWarehouseCode := c.LogicWarehouseCode
	if c.IsDefective != 0 {
		defectiveLogicWarehouse := &models.LogicWarehouse{}
		_ = defectiveLogicWarehouse.GetDefectiveOrPassedLogicWarehouse(e.Orm, c.LogicWarehouseCode, models.LwhType1)
		realLogicWarehouseCode = defectiveLogicWarehouse.LogicWarehouseCode
	}
	locations, err := models.GetStockLocationByLwhCodeWithLocationCode(e.Orm, realLogicWarehouseCode, c.LocationCode)
	if err != nil {
		e.Log.Errorf("StockEntryService StockLocationsSelect error:%s \r\n", err)
		return err
	}
	stockEntryProductsLocationGetResp := dto.StockEntryProductsLocationGetResp{}
	for _, item := range *locations {
		stockEntryProductsLocationGetResp.Regenerate(item)
		*outData = append(*outData, stockEntryProductsLocationGetResp)
	}
	return nil
}

func (e *StockEntry) GetPrintHtml(c *dto.StockEntryPrintHtmlReq, p *actions.DataPermission, outData *dto.StockEntryPrintHtmlResp) error {
	req := &dto.StockEntryGetReq{
		Id:   c.Id,
		Type: "print",
	}
	resp := &dto.StockEntryGetResp{}
	if err := e.Get(req, p, resp); err != nil {
		e.Log.Errorf("StockEntryService GetPrintHtml error:%s \r\n", err)
		return err
	}
	outDataMap := utils.GetChunkProductsForPrint(resp, "stockEntryProducts", map[string]int{
		"productName": 12,
		"mfgModel":    6,
		"brandName":   6,
		"vendorName":  8,
	}, 39, 46)

	path := "static/viewtpl/wc/entry_print.html"
	resStr, err := utils.View(path, outDataMap)
	if err != nil {
		return err
	}
	outData.Html = resStr
	return nil
}

// 部分入库场景下的打印入库单
func (e *StockEntry) PrintEntryHtml(tx *gorm.DB, stockEntry *models.StockEntry, reqProducts []*dto.StockEntryProducts, stockLocations []*models.StockLocation) (any, error) {

	// 打印数据结构
	productsChunk := []map[string]string{}
	printData := map[string]any{
		"entryCode":     stockEntry.EntryCode,
		"sourceCode":    stockEntry.SourceCode,
		"warehouseName": "",
		"nowDateTime":   utils.TimeFormat(time.Now()),
		"productsChunk": [][]map[string]string{productsChunk},
	}

	// 查询仓库名称
	wh := &models.Warehouse{}
	err := wh.GetByWarehouseCode(tx, stockEntry.WarehouseCode)
	if err != nil {
		return "", err
	}
	printData["warehouseName"] = wh.WarehouseName

	// 查商品信息
	goodsIdSlice := lo.Uniq(lo.Map(stockEntry.StockEntryProducts, func(item models.StockEntryProducts, _ int) int {
		return item.GoodsId
	}))
	goodsProductMap := models.GetGoodsProductInfoFromPc(tx, goodsIdSlice)

	// 查询库位编号map
	stockLocationMap := lo.Associate(stockLocations, func(item *models.StockLocation) (int, *models.StockLocation) {
		return item.Id, item
	})

	// 获取vendor map
	vendorIdSlice := lo.Uniq(lo.Map(stockEntry.StockEntryProducts, func(item models.StockEntryProducts, _ int) int {
		return item.VendorId
	}))
	vendorIdMap := models.GetVendorsMapByIds(tx, vendorIdSlice)

	for index, product := range reqProducts {
		goodInfo := goodsProductMap[product.GoodId]
		for _, sub := range product.StockEntryLocationInfo {
			itemPrintProduct := map[string]string{}
			itemPrintProduct["number"] = strconv.Itoa(index + 1)
			itemPrintProduct["skuCode"] = product.SkuCode
			itemPrintProduct["productName"] = goodInfo.Product.NameZh
			itemPrintProduct["mfgModel"] = goodInfo.Product.MfgModel
			itemPrintProduct["brandName"] = goodInfo.Product.Brand.BrandZh
			itemPrintProduct["vendorName"] = vendorIdMap[goodInfo.Product.VendorId]
			itemPrintProduct["vendorSkuCode"] = goodInfo.Product.SupplierSkuCode
			itemPrintProduct["salesUom"] = goodInfo.Product.SalesUom
			itemPrintProduct["actQuantity"] = strconv.Itoa(sub.ActQuantity)
			itemPrintProduct["locationCode"] = stockLocationMap[sub.StockLocationId].LocationCode
			productsChunk = append(productsChunk, itemPrintProduct)
		}
	}
	printData["productsChunk"] = [][]map[string]string{productsChunk}

	// 补全商品打印数据
	path := "static/viewtpl/wc/entry_print.html"
	html, err := utils.View(path, printData)
	if err != nil {
		return "", err
	}
	return html, nil
}

// GetPrintSkuInfo 打印SKU标签
func (e *StockEntry) GetPrintSkuInfo(c *dto.StockEntryPrintSkuReq, p *actions.DataPermission, outData *dto.StockEntryPrintSkuResp) error {
	dataBase := &dto.StockEntryScanSkuBaseInfoResp{}
	err := e.GetSkuInfoByGoodsId(c, dataBase)
	if err != nil {
		return err
	}
	err = utils.CopyDeep(outData, dataBase)
	if err != nil {
		return err
	}

	// todo 二维码处理
	outData.QrCode = utils.GetShortLink(e.Orm, ext.ExtConfig.Wc.PrintSkuPage+"?tenantId="+e.Orm.Statement.Context.(*gin.Context).GetHeader("tenant-id")+"&goodsId="+strconv.Itoa(c.GoodsId))
	return nil
}

// GetPrintSkuInfos 批量打印SKU标签
func (e *StockEntry) GetPrintSkuInfos(c *dto.StockEntryPrintSkusReq, p *actions.DataPermission) (*[]dto.StockEntryPrintSkusResp, error) {
	res := []dto.StockEntryPrintSkusResp{}

	// 异常校验
	if len(c.GoodsIds) == 0 {
		return &res, nil
	}

	// 查询商品
	goodsMap := models.GetGoodsProductInfoFromPc(e.Orm, c.GoodsIds)

	// 查询货主
	vendorIds := lo.Map(lo.Values(goodsMap), func(item modelsPc.Goods, _ int) int {
		return item.VendorId
	})
	vendors, err := models.GetVendorListByIds(e.Orm, vendorIds)
	if err != nil {
		return nil, err
	}
	vendorsMap := lo.Associate(*vendors, func(item models.Vendors) (int, models.Vendors) {
		return item.Id, item
	})

	// 生成二维码
	for _, good := range goodsMap {
		vendor := vendorsMap[good.VendorId]
		item := dto.StockEntryPrintSkusResp{
			VendorCode:  vendor.Code,
			VendorName:  vendor.NameZh,
			SkuCode:     good.SkuCode,
			ProductNo:   good.ProductNo,
			BrandName:   good.Product.Brand.BrandZh,
			MfgModel:    good.Product.MfgModel,
			ProductName: good.Product.NameZh,
			SalesUom:    good.Product.SalesUom,
		}

		// 获取二维码
		url := ext.ExtConfig.Wc.PrintSkuPage + "?tenantId=" + e.Orm.Statement.Context.(*gin.Context).GetHeader("tenant-id") + "&goodsId=" + strconv.Itoa(good.Id)
		item.QrCode = utils.GetShortLink(e.Orm, url)
		res = append(res, item)
	}

	return &res, nil
}

func (e *StockEntry) GetSkuInfoByGoodsId(c *dto.StockEntryPrintSkuReq, outData *dto.StockEntryScanSkuBaseInfoResp) error {
	goodsMap := models.GetGoodsProductInfoFromPc(e.Orm, []int{c.GoodsId})
	vendor := &models.Vendors{}
	_ = vendor.GetById(e.Orm, goodsMap[c.GoodsId].VendorId)
	outData.VendorCode = vendor.Code
	outData.VendorName = vendor.NameZh
	outData.SkuCode = goodsMap[c.GoodsId].SkuCode
	outData.ProductNo = goodsMap[c.GoodsId].ProductNo
	outData.BrandName = goodsMap[c.GoodsId].Product.Brand.BrandZh
	outData.MfgModel = goodsMap[c.GoodsId].Product.MfgModel
	outData.ProductName = goodsMap[c.GoodsId].Product.NameZh
	outData.SalesUom = goodsMap[c.GoodsId].Product.SalesUom
	return nil
}

// sourceProductsInfo 根据来源单号查询来源商品的出库数量或退货数量
func (e *StockEntry) sourceProductsInfo(entryType string, sourceCode string) ([]models.SourceProductsInfo, error) {
	res := []models.SourceProductsInfo{}

	// 调拨单找出库单商品
	if entryType == "0" {
		err := e.Orm.Table("stock_outbound info").
			Joins("JOIN stock_outbound_products list ON info.outbound_code = list.outbound_code").
			Where("info.source_code = ?", sourceCode).
			Select("list.status as Status,list.sku_code as SkuCode,list.act_quantity as ActQuantity").Scan(&res).Error
		if err != nil {
			return nil, err
		}
	}

	// 售后退货单找出库商品
	if entryType == "1" {
		wcPrefix := global.GetTenantWcDBNameWithDB(e.Orm)
		ocPrefix := global.GetTenantOcDBNameWithDB(e.Orm)
		err := e.Orm.Table(wcPrefix+".stock_outbound info").Debug().
			Joins("JOIN "+ocPrefix+".cs_apply ca ON info.source_code = ca.order_id").
			Joins("JOIN stock_outbound_products list ON info.outbound_code = list.outbound_code").
			Where("ca.cs_no = ?", sourceCode).
			Select("list.status as Status,list.sku_code as SkuCode,list.act_quantity as ActQuantity").Scan(&res).Error
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (e *StockEntry) Export(c *dto.StockEntryGetPageReq, p *actions.DataPermission) (exportData []interface{}, err error) {
	//var err error
	var data models.StockEntry
	pcPrefix := global.GetTenantPcDBNameWithDB(e.Orm)
	//ocPrefix := global.GetTenantOcDBNameWithDB(e.Orm)
	var result []dto.StockEntryExportResp
	db := e.Orm.Model(&data).Preload("Warehouse").Preload("LogicWarehouse").Preload("StockEntryProducts").Preload("StockEntryProducts.StockEntryProductsSub").
		Joins("LEFT JOIN stock_entry_products entryProduct ON entryProduct.entry_code = stock_entry.entry_code AND entryProduct.deleted_at IS NULL").
		Joins("LEFT JOIN "+pcPrefix+".product product ON product.sku_code = entryProduct.sku_code AND product.deleted_at IS NULL").
		Joins("LEFT JOIN "+pcPrefix+".goods goods ON goods.id = entryProduct.goods_id AND goods.deleted_at IS NULL").
		//Joins("LEFT JOIN "+ocPrefix+".order_info oi ON oi.order_id = stock_entry.source_code").
		Joins("LEFT JOIN vendors vendors ON vendors.id = stock_entry.vendor_id AND vendors.deleted_at IS NULL").
		Select("stock_entry.*,vendors.name_zh as recipient").
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			//cDto.Paginate(c.GetPageSizeNegative(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
			StockEntryPermission(p),
			dtoStock.GenCreatedAtTimeSearch(c.CreatedAtStart, c.CreatedAtEnd, "stock_entry"),
			dtoStock.GenProductSearch(c.Sku, c.ProductNo, c.ProductName, "entryProduct"),
			dtoStock.GenEntryRecipientSearch(c.Recipient),
			func(db *gorm.DB) *gorm.DB {
				if c.SourceType != "" {
					db.Where("stock_entry.type = ?", c.SourceType)
				}
				db.Where("stock_entry.status not in ?", []int{0, 1})
				return db
			},
			dtoStock.GenEntryIdsSearch(c.Ids),
		).Group("stock_entry.entry_code").Order("stock_entry.id DESC").Find(&result)
	err = db.Error
	if err != nil {
		e.Log.Errorf("StockEntryService error:%s \r\n", err)
		//return err
	}

	//// 所有入库单的商品
	//StockEntryCode := lo.Uniq(lo.Map(result, func(item dto.StockEntryExportResp, _ int) string {
	//	return item.EntryCode
	//}))
	//
	//subProductsMap := make(map[string][]models.StockEntryProductsSubCustom, 0)
	//for _, product := range result {
	//	subProductsMap[product.EntryCode] = append(subProductsMap[product.EntryCode], product)
	//}

	// 获取goodsid 去重切片
	goodsIdMap := map[int]int{}
	vendorIdMap := map[int]int{}
	for _, item := range result {
		for _, product := range item.StockEntryProducts {
			goodsIdMap[product.GoodsId] = product.GoodsId
			vendorIdMap[product.VendorId] = product.VendorId
		}
	}
	goodsIdSlice := []int{}
	vendorIdSlice := []int{}
	for _, goodId := range goodsIdMap {
		goodsIdSlice = append(goodsIdSlice, goodId)
	}
	for _, vendorId := range vendorIdMap {
		vendorIdSlice = append(vendorIdSlice, vendorId)
	}

	// 获取goods product 信息通过pc
	goodsProductMap := models.GetGoodsProductInfoFromPc(e.Orm, goodsIdSlice)

	// 获取vendor map
	vendorMap := models.GetVendorsMapByIds(e.Orm, vendorIdSlice)

	// 库位map
	stl := models.StockLocation{}
	locationRes := stl.Find(e.Orm)
	LocationIdsMaps := map[int]string{}
	for _, v := range locationRes {
		LocationIdsMaps[v.Id] = v.LocationCode
	}
	//LocationIds = lo.Map(LocationIds, func(locationRes models.StockLocation, _ int) string {
	//	return locationRes.LocationCode
	//})

	for _, item := range result {
		for _, product := range item.StockEntryProducts {
			// 兼容旧数据
			if item.Status == "2" && len(product.StockEntryProductsSub) == 0 {
				oldData := []models.StockEntryProductsSub{{
					ShouldQuantity:  product.Quantity,
					ActQuantity:     product.ActQuantity,
					StockLocationId: product.StockLocationId,
					//EntryTime:       product.EntryTime,
				}}
				product.StockEntryProductsSub = oldData
			}
			for _, sub := range product.StockEntryProductsSub {
				entryProduct := dto.StockEntryExport{}
				entryProduct.TypeName = utils.GetTFromMap(item.Type, models.EntryTypeMap)
				entryProduct.StatusName = utils.GetTFromMap(item.Status, models.EntryStatusMap)
				entryProduct.WarehouseName = item.Warehouse.WarehouseName
				entryProduct.LogicWarehouseName = item.LogicWarehouse.LogicWarehouseName
				//entryProduct.VendorName = item.Recipient
				entryProduct.EntryCode = item.EntryCode
				entryProduct.SourceCode = item.SourceCode
				entryProduct.SourceTypeName = utils.GetTFromMap(item.Type, models.EntrySourceTypeMap)
				entryProduct.Remark = item.Remark
				entryProduct.SkuCode = product.SkuCode
				productInfo := goodsProductMap[product.GoodsId]
				entryProduct.ProductName = productInfo.Product.NameZh
				entryProduct.MfgModel = productInfo.Product.MfgModel
				entryProduct.BrandName = productInfo.Product.Brand.BrandZh
				entryProduct.SalesUom = productInfo.Product.SalesUom
				entryProduct.ProductNo = productInfo.ProductNo
				entryProduct.VendorSkuCode = productInfo.Product.SupplierSkuCode
				entryProduct.VendorName = vendorMap[product.VendorId]
				entryProduct.Quantity = product.Quantity
				entryProduct.ActQuantity = sub.ActQuantity
				entryProduct.DiffNum = product.Quantity - product.ActQuantity

				entryProduct.LocationCode = LocationIdsMaps[sub.StockLocationId]
				entryProduct.IsDefective = utils.GetTFromMap(strconv.Itoa(product.IsDefective), models.EntryProductIsDefectiveMap)
				if product.IsDefective == 1 {
					entryProduct.EntryTime = product.EntryTime.Format("2006-01-02 15:04:05")
				} else {
					if !sub.EntryTime.IsZero() {
						entryProduct.EntryTime = sub.EntryTime.Format("2006-01-02 15:04:05")
					}
				}
				exportData = append(exportData, entryProduct)
			}
		}
	}
	return
}

func (e *StockEntry) ValidateSkus(d *dto.StockEntryferValidateSkusReq, p *actions.DataPermission, list *[]dto.StockTransferValidateSkusResp) error {
	whTo := models.Warehouse{}
	if err := whTo.GetByWarehouseCode(e.Orm, d.WarehouseCode); err != nil {
		return errors.New("入库仓不存在")
	}
	//whFrom := models.Warehouse{}
	//if err := whFrom.GetByWarehouseCode(e.Orm, d.FromWarehouseCode); err != nil {
	//	//return errors.New("出库仓不存在")
	//}
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
	//skuGoodsFromSlice := []string{}
	skuGoodsToMap := map[string]modelsPc.Goods{}
	skuGoodsToSlice := []string{}

	//if !whFrom.CheckIsVirtual() {
	//	skuGoodsFromMap, skuGoodsFromSlice = models.GetGoodsInfoMapByThreeFromPcClient(e.Orm, skuSlice, d.FromWarehouseCode, d.VendorId, 0)
	//	skuGoodsFromErrSlice = lo.Without(skuSlice, skuGoodsFromSlice...)
	//	skuGoodsFromErrSlice = lo.Without(skuGoodsFromErrSlice, skuErrSlice...)
	//	for key, item := range skuGoodsFromMap {
	//		if item.ApproveStatus != 1 {
	//			skuGoodsFromAuditErrSlice = append(skuGoodsFromAuditErrSlice, key)
	//		}
	//	}
	//}
	if !whTo.CheckIsVirtual() {
		skuGoodsToMap, skuGoodsToSlice = models.GetGoodsInfoMapByThreeFromPcClient(e.Orm, skuSlice, d.WarehouseCode, d.VendorId, 0)
		skuGoodsToErrSlice = lo.Without(skuSlice, skuGoodsToSlice...)
		skuGoodsToErrSlice = lo.Without(skuGoodsToErrSlice, skuErrSlice...)
		for key, item := range skuGoodsToMap {
			if item.ApproveStatus != 1 {
				skuGoodsToAuditErrSlice = append(skuGoodsToAuditErrSlice, key)
			}
		}
	}
	errStr := FormatValidateSkusErrInfoForEntry(skuErrSlice, skuGoodsFromErrSlice, skuGoodsFromAuditErrSlice, skuGoodsToErrSlice, skuGoodsToAuditErrSlice)
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
			stockTransferValidateSkusResp.GoodId = goodsInfo.Id
		}

		if tmpData[item].Id != 0 {
			stockLocations, _, _ := models.GetStockLocationsForEntry(e.Orm, d.LogicWarehouseCode, tmpData[item].Id)
			if len(*stockLocations) > 0 {
				stockTransferValidateSkusResp.StockLocationId = (*stockLocations)[0].Id
			}

		}
		*list = append(*list, stockTransferValidateSkusResp)
	}
	if len(*list) == 0 {
		errors.New("该商品暂无信息")
	}

	return nil
}

func FormatValidateSkusErrInfoForEntry(skuErrSlice, skuGoodsFromErrSlice, skuGoodsFromAuditErrSlice, skuGoodsToErrSlice, skuGoodsToAuditErrSlice []string) string {
	errSlice := []string{}
	if len(skuErrSlice) > 0 {
		errSku := strings.Join(skuErrSlice, ",")
		errSlice = append(errSlice, "【"+errSku+"】不存在")
	}
	//if len(skuGoodsFromErrSlice) > 0 {
	//	errSku := strings.Join(skuGoodsFromErrSlice, ",")
	//	errSlice = append(errSlice, "【"+errSku+"】没有对应货主和出库仓的商品关系")
	//}
	//if len(skuGoodsFromAuditErrSlice) > 0 {
	//	errSku := strings.Join(skuGoodsFromAuditErrSlice, ",")
	//	errSlice = append(errSlice, "【"+errSku+"】出库仓的商品关系待审核或审核不通过")
	//}
	if len(skuGoodsToErrSlice) > 0 {
		errSku := strings.Join(skuGoodsToErrSlice, ",")
		errSlice = append(errSlice, "【"+errSku+"】没有对应货主和入库仓的商品关系")
	}
	//if len(skuGoodsToAuditErrSlice) > 0 {
	//	errSku := strings.Join(skuGoodsToAuditErrSlice, ",")
	//	errSlice = append(errSlice, "【"+errSku+"】入库仓的商品关系待审核或审核不通过")
	//}
	if len(errSlice) > 0 {
		return strings.Join(errSlice, "，") + "，请联系相关人员处理！"
	}
	return ""
}
