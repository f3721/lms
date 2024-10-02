package admin

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	modelsOc "go-admin/app/oc/models"
	ocDto "go-admin/app/oc/service/admin/dto"
	ocClient "go-admin/common/client/oc"
	dtoStock "go-admin/common/dto/stock/dto"
	"go-admin/common/global"
	"go-admin/common/utils"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
	"gorm.io/gorm/clause"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/wc/models"
	"go-admin/app/wc/service/admin/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type StockOutbound struct {
	service.Service
}

// 出库单数据权限
func StockOutboundPermission(p *actions.DataPermission) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db.Where("stock_outbound.warehouse_code in ?", utils.Split(p.AuthorityWarehouseId))
		return db
	}
}

// GetPage 获取StockOutbound列表
func (e *StockOutbound) GetPage(c *dto.StockOutboundGetPageReq, p *actions.DataPermission, outData *[]dto.StockOutboundGetPageResp, count *int64) error {
	var err error
	var data models.StockOutbound
	pcPrefix := global.GetTenantPcDBNameWithDB(e.Orm)
	ocPrefix := global.GetTenantOcDBNameWithDB(e.Orm)

	db := e.Orm.Model(&data).Preload("Warehouse").Preload("LogicWarehouse").
		Joins("LEFT JOIN stock_outbound_products outboundProduct ON outboundProduct.outbound_code = stock_outbound.outbound_code AND outboundProduct.deleted_at IS NULL").
		Joins("LEFT JOIN "+pcPrefix+".product product ON product.sku_code = outboundProduct.sku_code AND product.deleted_at IS NULL").
		Joins("LEFT JOIN "+pcPrefix+".goods goods ON goods.id = outboundProduct.goods_id AND goods.deleted_at IS NULL").
		Joins("LEFT JOIN "+ocPrefix+".order_info oi ON oi.order_id = stock_outbound.source_code").
		Select("stock_outbound.*,oi.user_name as recipient").
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSizeNegative(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
			StockOutboundPermission(p),
			dtoStock.GenCreatedAtTimeSearch(c.CreatedAtStart, c.CreatedAtEnd, "stock_outbound"),
			dtoStock.GenProductSearch(c.Sku, c.ProductNo, c.ProductName, "outboundProduct"),
			dtoStock.GenRecipientSearch(c.Recipient),
			func(db *gorm.DB) *gorm.DB {
				if c.SourceType != "" {
					db.Where("stock_outbound.type = ?", c.SourceType)
				}
				return db
			},
		).Group("stock_outbound.outbound_code").Order("stock_outbound.id DESC").Find(outData)

	err = db.Error
	if err != nil {
		e.Log.Errorf("StockOutboundService GetPage error:%s \r\n", err)
		return err
	}
	var list = &[]dto.StockOutboundGetPageResp{}
	for _, item := range *outData {
		item.SetOutboundRulesByStatus()
		item.WarehouseName = item.Warehouse.WarehouseName
		item.LogicWarehouseName = item.LogicWarehouse.LogicWarehouseName
		item.TypeName = utils.GetTFromMap(item.Type, models.OutboundTypeMap)
		item.StatusName = utils.GetTFromMap(item.Status, models.OutboundStatusMap)
		item.SourceTypeName = utils.GetTFromMap(item.Type, models.OutboundSourceTypeMap)
		*list = append(*list, item)
	}
	*outData = *list

	err = e.Orm.Table("(?) as u", db.Limit(-1).Offset(-1)).Count(count).Error
	if err != nil {
		e.Log.Errorf("StockOutboundService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Export 出库明细维护导出
func (e *StockOutbound) Export(c *dto.StockOutboundGetPageReq, p *actions.DataPermission) (exportData []interface{}, err error) {
	var data models.StockOutbound
	pcPrefix := global.GetTenantPcDBNameWithDB(e.Orm)
	ocPrefix := global.GetTenantOcDBNameWithDB(e.Orm)
	var result []dto.StockOutboundExportResp
	db := e.Orm.Model(&data).Preload("Warehouse").Preload("LogicWarehouse").
		Joins("RIGHT JOIN stock_outbound_products outboundProduct ON outboundProduct.outbound_code = stock_outbound.outbound_code AND outboundProduct.deleted_at IS NULL").
		Joins("LEFT JOIN "+pcPrefix+".product product ON product.sku_code = outboundProduct.sku_code AND product.deleted_at IS NULL").
		Joins("LEFT JOIN "+pcPrefix+".goods goods ON goods.id = outboundProduct.goods_id AND goods.deleted_at IS NULL").
		Joins("LEFT JOIN "+ocPrefix+".order_info oi ON oi.order_id = stock_outbound.source_code").
		Select("stock_outbound.*,oi.user_name as recipient").
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			actions.Permission(data.TableName(), p),
			StockOutboundPermission(p),
			dtoStock.GenCreatedAtTimeSearch(c.CreatedAtStart, c.CreatedAtEnd, "stock_outbound"),
			dtoStock.GenProductSearch(c.Sku, c.ProductNo, c.ProductName, "outboundProduct"),
			dtoStock.GenRecipientSearch(c.Recipient),
			func(db *gorm.DB) *gorm.DB {
				if c.SourceType != "" {
					db.Where("stock_outbound.type = ?", c.SourceType)
				}
				db.Where("stock_outbound.status not in ?", []int{0, 1})
				return db
			},
			dtoStock.GenIdsSearch(c.Ids),
		).Group("stock_outbound.outbound_code").Order("stock_outbound.id DESC").Find(&result)

	err = db.Error
	if err != nil {
		e.Log.Errorf("StockOutboundService GetPage error:%s \r\n", err)
		return
	}

	// 所有出库单的商品
	outboundCode := lo.Uniq(lo.Map(result, func(item dto.StockOutboundExportResp, _ int) string {
		return item.OutboundCode
	}))
	OutboundProductSubCustom := models.OutboundProductSubCustom{}
	subProducts, err := OutboundProductSubCustom.GetList(e.Orm, outboundCode)
	if err != nil {
		return
	}
	subProductsMap := make(map[string][]models.OutboundProductSubCustom, 0)
	for _, product := range subProducts {
		subProductsMap[product.OutboundCode] = append(subProductsMap[product.OutboundCode], product)
	}
	// 获取goods product 信息通过pc
	goodsIdSlice := lo.Uniq(lo.Map(subProducts, func(item models.OutboundProductSubCustom, _ int) int {
		return item.GoodsId
	}))
	goodsProductMap := models.GetGoodsProductInfoFromPc(e.Orm, goodsIdSlice)
	// 获取vendor map
	vendorIdSlice := lo.Uniq(lo.Map(subProducts, func(item models.OutboundProductSubCustom, _ int) int {
		return item.VendorId
	}))
	// 查询库位部分出库历史记录
	subIds := lo.Map(subProducts, func(item models.OutboundProductSubCustom, _ int) int {
		return item.SubId
	})
	subLog := &models.StockOutboundProductsSubLog{}
	subList, err := subLog.ListbySubIds(e.Orm, subIds)
	if err != nil {
		return
	}
	subListMap := map[int][]*models.StockOutboundProductsSubLog{}
	for _, item := range subList {
		subListMap[item.SubId] = append(subListMap[item.SubId], item)
	}

	vendorIdMap := models.GetVendorsMapByIds(e.Orm, vendorIdSlice)
	for _, item := range result {
		// 出库商品明细
		for _, product := range subProductsMap[item.OutboundCode] {
			product.SubLog = subListMap[product.SubId]
			// 兼容旧数据
			if item.Status == "2" && len(product.SubLog) == 0 {
				oldSubLog := []*models.StockOutboundProductsSubLog{
					{
						SubId:               product.SubId,
						OutboundTime:        item.OutboundTime,
						LocationActQuantity: product.LocationActQuantity,
					},
				}
				product.SubLog = oldSubLog
			}
			for _, sub := range product.SubLog {
				outboundProduct := dto.StockOutboundExport{}
				outboundProduct.TypeName = utils.GetTFromMap(item.Type, models.OutboundTypeMap)
				outboundProduct.StatusName = utils.GetTFromMap(item.Status, models.OutboundStatusMap)
				outboundProduct.WarehouseName = item.Warehouse.WarehouseName
				outboundProduct.LogicWarehouseName = item.LogicWarehouse.LogicWarehouseName
				outboundProduct.Recipient = item.Recipient
				outboundProduct.OutboundCode = item.OutboundCode
				outboundProduct.SourceCode = item.SourceCode
				outboundProduct.SourceTypeName = utils.GetTFromMap(item.Type, models.OutboundSourceTypeMap)
				outboundProduct.Remark = item.Remark
				outboundProduct.SkuCode = product.SkuCode
				productInfo := goodsProductMap[product.GoodsId]
				outboundProduct.ProductName = productInfo.Product.NameZh
				outboundProduct.MfgModel = productInfo.Product.MfgModel
				outboundProduct.BrandName = productInfo.Product.Brand.BrandZh
				outboundProduct.SalesUom = productInfo.Product.SalesUom
				outboundProduct.ProductNo = productInfo.ProductNo
				outboundProduct.VendorSkuCode = productInfo.Product.SupplierSkuCode
				outboundProduct.VendorName = vendorIdMap[product.VendorId]
				outboundProduct.LocationQuantity = product.LocationQuantity
				outboundProduct.LocationActQuantity = sub.LocationActQuantity
				outboundProduct.DiffNum = product.LocationQuantity - product.LocationActQuantity
				outboundProduct.LocationCode = product.StockLocation.LocationCode
				outboundProduct.OutboundTime = sub.OutboundTime.Format("2006-01-02 15:04:05")
				if outboundProduct.OutboundTime == "0001-01-01 00:00:00" {
					outboundProduct.OutboundTime = ""
				}
				exportData = append(exportData, outboundProduct)
			}
		}
	}
	return
}

// Get 获取StockOutbound对象
func (e *StockOutbound) Get(d *dto.StockOutboundGetReq, p *actions.DataPermission, outData *dto.StockOutboundGetResp) error {
	model, err := e.GetBaseInfo(d, p)
	if err != nil {
		return err
	}
	if err := utils.CopyDeep(outData, model); err != nil {
		e.Log.Errorf("Service GetStockOutbound error:%s \r\n", err)
		return err
	}
	return FormatDetailForOutbound(e.Orm, outData, d.Type)
}

func (e *StockOutbound) GetBaseInfo(d *dto.StockOutboundGetReq, p *actions.DataPermission) (*models.StockOutbound, error) {
	var data models.StockOutbound
	var model = &models.StockOutbound{}

	err := e.Orm.Model(&data).Preload("Warehouse").Preload("LogicWarehouse").
		Scopes(
			actions.Permission(data.TableName(), p),
			StockOutboundPermission(p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetStockOutbound error:%s \r\n", err)
		return nil, err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return nil, err
	}
	return model, nil
}

func (e *StockOutbound) GetNoLocation(d *dto.StockOutboundGetReq, p *actions.DataPermission, outData *dto.StockOutboundNoLocationResp) error {
	model, err := e.GetBaseInfo(d, p)
	if err != nil {
		return err
	}
	if err := utils.CopyDeep(outData, model); err != nil {
		e.Log.Errorf("Service StockOutbound GetNoLocation error:%s \r\n", err)
		return err
	}
	if err := outData.InitProductData(e.Orm); err != nil {
		return err
	}
	FormatOutboundData(e.Orm, outData.StockOutboundCommonResp)
	// 拼接地址相关信息
	FormatAddrForStockOutbound(e.Orm, outData.StockOutboundCommonResp)
	return nil
}

func FormatAddrForStockOutbound(tx *gorm.DB, data *dto.StockOutboundCommonResp) {
	data.AddressFullName = data.ProvinceName + data.CityName + data.DistrictName + data.Address
	data.ReceiveAddrInfo.AddressFullName = data.ReceiveAddrInfo.ProvinceName + data.ReceiveAddrInfo.CityName + data.ReceiveAddrInfo.DistrictName + data.ReceiveAddrInfo.Address
}

// 出库单格式化输出
func FormatDetailForOutbound(tx *gorm.DB, data *dto.StockOutboundGetResp, Type string) error {
	if err := data.InitProductData(tx, Type); err != nil {
		return err
	}
	FormatOutboundData(tx, data.StockOutboundCommonResp)
	return nil
}

func FormatOutboundData(tx *gorm.DB, data *dto.StockOutboundCommonResp) {
	data.WarehouseName = data.Warehouse.WarehouseName
	data.LogicWarehouseName = data.LogicWarehouse.LogicWarehouseName
	data.TypeName = utils.GetTFromMap(data.Type, models.OutboundTypeMap)
	data.StatusName = utils.GetTFromMap(data.Status, models.OutboundStatusMap)
	data.SourceTypeName = utils.GetTFromMap(data.Type, models.OutboundSourceTypeMap)

	data.Address = data.Warehouse.Address
	data.District = data.Warehouse.District
	data.City = data.Warehouse.City
	data.Province = data.Warehouse.Province
	data.DistrictName = data.Warehouse.DistrictName
	data.CityName = data.Warehouse.CityName
	data.ProvinceName = data.Warehouse.ProvinceName

	// 收货地址处理
	if data.IsOrderOutbound() {
		orderInfo := models.GeOrderInfoFromOc(tx, data.SourceCode)
		data.ExternalOrderNo = orderInfo.ContractNo

		data.Remark = orderInfo.Remark
		data.ReceiveAddrInfo.Address = orderInfo.Address
		data.ReceiveAddrInfo.District = orderInfo.AreaId
		data.ReceiveAddrInfo.City = orderInfo.CityId
		data.ReceiveAddrInfo.Province = orderInfo.ProvinceId
		data.ReceiveAddrInfo.Mobile = orderInfo.Mobile
		data.ReceiveAddrInfo.Linkman = orderInfo.Consignee

		data.ReceiveAddrInfo.UserName = orderInfo.UserName
		data.ReceiveAddrInfo.UserCompanyName = orderInfo.UserCompanyName
		// todo 用户接口 获取用户手机号
		userInfo := models.GetUserInfoByIdFromUc(tx, orderInfo.UserId)
		data.ReceiveAddrInfo.UserPhone = userInfo.UserPhone
		data.ReceiveAddrInfo.UserDepartment = strings.Join(userInfo.UserDepartments.Names, "-")

		regionIds := []int{
			data.ReceiveAddrInfo.Province,
			data.ReceiveAddrInfo.City,
			data.ReceiveAddrInfo.District,
		}
		regionsMap := models.GeRegionMapByIds(tx, regionIds)
		data.ReceiveAddrInfo.ProvinceName = regionsMap[data.ReceiveAddrInfo.Province]
		data.ReceiveAddrInfo.CityName = regionsMap[data.ReceiveAddrInfo.City]
		data.ReceiveAddrInfo.DistrictName = regionsMap[data.ReceiveAddrInfo.District]
		data.ReceiveAddrInfo.AddressFullName = data.ReceiveAddrInfo.ProvinceName + " " +
			data.ReceiveAddrInfo.CityName + " " +
			data.ReceiveAddrInfo.DistrictName + " " +
			data.ReceiveAddrInfo.Address
	} else {
		transfer := &models.StockTransfer{}
		_ = transfer.GetByTransferCode(tx, data.SourceCode)
		data.ReceiveAddrInfo.Address = transfer.Address
		data.ReceiveAddrInfo.District = transfer.District
		data.ReceiveAddrInfo.City = transfer.City
		data.ReceiveAddrInfo.Province = transfer.Province
		data.ReceiveAddrInfo.Mobile = transfer.Mobile
		data.ReceiveAddrInfo.Linkman = transfer.Linkman

		data.ReceiveAddrInfo.ProvinceName = transfer.ProvinceName
		data.ReceiveAddrInfo.CityName = transfer.CityName
		data.ReceiveAddrInfo.DistrictName = transfer.DistrictName
	}
}

// CreateOutboundForOrder 创建StockOutbound对象（领用单）
func CreateOutboundForOrder(tx *gorm.DB, c *dto.StockOutboundInsertForOrderReq) (*models.StockOutbound, error) {
	var data models.StockOutbound
	if err := CheckStockOutboundInsertReqForOrder(tx, c); err != nil {
		return nil, err
	}
	if err := c.Generate(tx, &data); err != nil {
		return nil, err
	}

	if err := data.InsertOutbound(tx, models.OutboundType1); err != nil {
		return nil, err
	}
	for _, item := range data.StockOutboundProducts {
		if err := item.SplitProductAndLockLocation(tx, data.OutboundCode, data.Type); err != nil {
			return nil, err
		}
	}
	return &data, nil
}

func CheckStockOutboundInsertReqForOrder(tx *gorm.DB, req *dto.StockOutboundInsertForOrderReq) error {
	warehouse := &models.Warehouse{}
	if err := warehouse.GetByWarehouseCode(tx, req.WarehouseCode); err != nil {
		return err
	}
	if warehouse.CheckIsVirtual() {
		return errors.New("虚拟仓无法创建领用单类型的出库单")
	}
	stockOutbound := &models.StockOutbound{}
	stockOutbound.GetBySourceCode(tx, req.SourceCode)
	if stockOutbound.Id != 0 {
		return errors.New("领用单对应的出库单已存在，无法重复创建")
	}
	skuTemp := []string{}
	for _, item := range req.StockOutboundProducts {
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

// CancelOutboundForCsOrder 取消出库单-售后所有商品时调用
func CancelOutboundForCsOrder(tx *gorm.DB, c *dto.StockOutboundPartCancelForCsOrderReq) error {

	// 查询出库单
	var data = models.StockOutbound{}
	wcPrefix := global.GetTenantWcDBNameWithDB(tx)
	if err := CheckOutboundForCsOrder(tx, &data, c.OrderId); err != nil {
		return err
	}
	oldData := data

	// 部分发货情况下-订单发货状态:全部发货
	if data.Status == "3" {
		ocPrefix := global.GetTenantOcDBNameWithDB(tx)
		err := tx.Table(ocPrefix+".order_info").Where("order_id = ?", c.OrderId).Update("send_status", "2").Error
		if err != nil {
			return err
		}
	}

	// 变更出库单状态
	c.GenerateOutBound(&data)
	if err := tx.Table(wcPrefix + "." + data.TableName()).Omit(clause.Associations).Save(&data).Error; err != nil {
		return errors.New(c.OrderId + ",订单对应的出库单更新失败")
	}
	// 商品行+库位行取消处理
	_, err := PartCancelOutboundHandleProduct(tx, data.OutboundCode, c)
	if err != nil {
		return err
	}

	// 出库单日志
	beforeDataStr, _ := json.Marshal(&oldData)
	afterDataStr, _ := json.Marshal(&data)
	dataStr, _ := json.Marshal(c)
	opLog := models.OperateLogs{
		DataId:       data.OutboundCode,
		ModelName:    models.OutboundLogModelName,
		Type:         models.OutboundLogModelTypeCancelForCsOrder,
		Before:       string(beforeDataStr),
		Data:         string(dataStr),
		After:        string(afterDataStr),
		OperatorId:   c.UpdateBy,
		OperatorName: c.UpdateByName,
	}
	_ = opLog.InsertItem(tx)

	return nil
}

func CheckOutboundForCsOrder(tx *gorm.DB, data *models.StockOutbound, orderId string) error {
	data.GetBySourceCode(tx, orderId)
	if data.Id == 0 {
		return errors.New(orderId + ",订单对应的出库单未找到")
	}
	if lo.IndexOf([]string{"1", "3"}, data.Status) == -1 {
		return errors.New(orderId + ",订单对应的出库单状态不正确,非未出库或部分发货状态")
	}
	return nil
}

// PartCancelOutboundForCsOrder 取消出库单-售后部分商品时调用
func PartCancelOutboundForCsOrder(tx *gorm.DB, c *dto.StockOutboundPartCancelForCsOrderReq) error {
	var data = models.StockOutbound{}
	if err := CheckStockOutboundPartCancelForCsOrderReq(c); err != nil {
		return err
	}
	if err := CheckOutboundForCsOrder(tx, &data, c.OrderId); err != nil {
		return err
	}

	// 商品行+库位行取消处理
	oldData, err := PartCancelOutboundHandleProduct(tx, data.OutboundCode, c)
	if err != nil {
		return err
	}

	stockOutboundProducts := &models.StockOutboundProducts{}
	afterData, _ := stockOutboundProducts.GetList(tx, data.OutboundCode)

	// 出库单日志
	beforeDataStr, _ := json.Marshal(&oldData)
	afterDataStr, _ := json.Marshal(&afterData)
	dataStr, _ := json.Marshal(c)
	opLog := models.OperateLogs{
		DataId:       data.OutboundCode,
		ModelName:    models.OutboundLogModelName,
		Type:         models.OutboundLogModelTypePartCancelForCsOrder,
		Before:       string(beforeDataStr),
		Data:         string(dataStr),
		After:        string(afterDataStr),
		OperatorId:   c.UpdateBy,
		OperatorName: c.UpdateByName,
	}
	_ = opLog.InsertItem(tx)
	return nil
}

func CheckStockOutboundPartCancelForCsOrderReq(req *dto.StockOutboundPartCancelForCsOrderReq) error {
	goodsTemp := []int{}
	for _, item := range req.CsOrderProducts {
		if item.Quantity <= 0 {
			return errors.New("商品数量信息异常")
		}
		goodsTemp = append(goodsTemp, item.GoodsId)
	}
	goodsTempLen := len(goodsTemp)
	goodsTemp = lo.Uniq(goodsTemp)
	if goodsTempLen > len(goodsTemp) {
		return errors.New("重复的goodsId")
	}
	return nil
}

// 订单取消 出库单产品处理
func PartCancelOutboundHandleProduct(tx *gorm.DB, outboundCode string, req *dto.StockOutboundPartCancelForCsOrderReq) ([]models.StockOutboundProducts, error) {
	// 1. 查询出库单商品
	stockOutboundProducts := &models.StockOutboundProducts{}
	products, err := stockOutboundProducts.GetList(tx, outboundCode)
	if err != nil {
		return nil, errors.New("获取出库单商品失败")
	}
	productsMap := lo.Associate(products, func(f models.StockOutboundProducts) (int, models.StockOutboundProducts) {
		return f.GoodsId, f
	})

	// 2. 查询出库单[商品库位]信息列表
	OutboundProductSubCustom := &models.OutboundProductSubCustom{}
	subProducts, err := OutboundProductSubCustom.GetList(tx, []string{outboundCode})
	if err != nil || len(subProducts) == 0 {
		return nil, errors.New("获取出库单商品（库位）失败")
	}

	// 3. 商品行处理
	for _, item := range req.CsOrderProducts {
		product := productsMap[item.GoodsId]
		if product.Id == 0 {
			return nil, errors.New("对应的原始出单中商品不存在")
		}

		// 校验取消数量
		leftUnOutQuantity := product.Quantity - product.ActQuantity
		if item.Quantity > leftUnOutQuantity {
			return nil, errors.New("部分取消数量不能大于剩余未出库数量")
		}

		// 新逻辑：直接扣减数量
		req.GenerateProduct(&product)
		if err := product.SubQuantity(tx, item.Quantity); err != nil {
			return nil, err
		}

		stockInfo := &models.StockInfo{}
		if !stockInfo.CheckLockStockByGoodsIdAndLogicWarehouseCode(tx, item.Quantity, product.GoodsId, product.LogicWarehouseCode) {
			return nil, errors.New(product.SkuCode + "逻辑仓:" + product.LogicWarehouseCode + ",锁定库存不足")
		}
		if err := stockInfo.Unlock(tx, item.Quantity, &models.StockLockLog{
			DocketCode: outboundCode,
			FromType:   models.StockLockLogFromType2,
			Remark:     models.RemarkStockLockLogCsOrderCancelPart,
		}); err != nil {
			return nil, errors.New(product.SkuCode + ",解锁失败")
		}

	}

	// 4. 库位行处理
	for _, item := range req.CsOrderProducts {
		quantity := item.Quantity
		matchSubProducts := lo.Filter(subProducts, func(x models.OutboundProductSubCustom, _ int) bool {
			return x.GoodsId == item.GoodsId
		})
		sort.Slice(matchSubProducts, func(i, j int) bool {
			return i > j
		})
		for _, matchItem := range matchSubProducts {

			// 处理库位数据-新逻辑-有部分出库的概念
			locationReduceNum := 0                                                     // 库位扣减数量
			unlockNum := 0                                                             // 解锁锁库数量
			subLeftUnOut := matchItem.LocationQuantity - matchItem.LocationActQuantity // 剩余未出库数量
			stockOutboundProductsSub := &models.StockOutboundProductsSub{}

			// 跳过可扣为0的
			if subLeftUnOut == 0 {
				continue
			}

			// 新逻辑：直接扣减数量
			if quantity >= subLeftUnOut { // 取消数量 >= 当前库位数量 = 扣减剩余未出库数量
				locationReduceNum = subLeftUnOut
				unlockNum = subLeftUnOut
				if err := stockOutboundProductsSub.SubLocationQuantity(tx, matchItem.SubId, locationReduceNum); err != nil {
					return nil, err
				}
			} else { // 取消数量 < 当前库位数量 = 扣减值为取消数量
				locationReduceNum = quantity
				unlockNum = quantity
				if err := stockOutboundProductsSub.SubLocationQuantity(tx, matchItem.SubId, locationReduceNum); err != nil {
					return nil, err
				}
			}

			// 库位锁定库存释放
			locationGoods := &models.StockLocationGoods{}
			if !locationGoods.CheckLockStockByGoodsIdAndLocationId(tx, unlockNum, matchItem.GoodsId, matchItem.StockLocationId) {
				return nil, errors.New(matchItem.SkuCode + "库位:" + matchItem.StockLocation.LocationCode + ",锁定库位库存不足")
			}
			if err := locationGoods.Unlock(tx, unlockNum, &models.StockLocationGoodsLockLog{
				DocketCode: outboundCode,
				FromType:   models.StockLocationGoodsLockLogFromType2,
				Remark:     models.RemarkLocationProductsLockCsOrderCancelPart,
			}); err != nil {
				return nil, errors.New(matchItem.SkuCode + ",解锁库位库存失败")
			}

			// 逐个递减已经处理过的数量
			quantity -= subLeftUnOut
			if quantity <= 0 {
				break
			}
		}
		if quantity > 0 {
			return nil, fmt.Errorf("商品ID[%v]取消数量大于出库总数量", item.GoodsId)
		}
	}

	return products, nil
}

// Insert 创建StockOutbound对象
func InsertOutbound(tx *gorm.DB, c *dto.StockOutboundInsertReq, Type string) (*models.StockOutbound, error) {
	var data models.StockOutbound
	c.Generate(&data)
	if err := data.InsertOutbound(tx, Type); err != nil {
		return nil, err
	}
	return &data, nil
}

// Update 修改StockOutbound对象
func (e *StockOutbound) Update(c *dto.StockOutboundUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.StockOutbound{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
		StockOutboundPermission(p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("StockOutboundService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除StockOutbound
func (e *StockOutbound) Remove(d *dto.StockOutboundDeleteReq, p *actions.DataPermission) error {
	var data models.StockOutbound

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveStockOutbound error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}

// PartOutbound 部分出库
func (e *StockOutbound) PartOutbound(c *dto.StockPartOutboundReq, p *actions.DataPermission) error {
	currTime := time.Now()

	// 1. 查询出库单信息
	var stockOutbound = &models.StockOutbound{}
	e.Orm.Preload("StockOutboundProducts").Preload("Warehouse").
		Scopes(
			actions.Permission(stockOutbound.TableName(), p),
			StockOutboundPermission(p),
		).First(stockOutbound, c.Id)
	if stockOutbound.Id == 0 {
		return errors.New("出库单不存在或没有权限")
	}
	oldData := stockOutbound

	// 校验出库单状态
	if err := e.CheckConfirmOutbound(stockOutbound); err != nil {
		return err
	}

	// 出库单商品Map
	orgProductMap := lo.Associate(stockOutbound.StockOutboundProducts, func(item models.StockOutboundProducts) (int, *models.StockOutboundProducts) {
		return item.Id, &item
	})

	// 2. 查询出库单库位数据Map
	stockOutboundSub := &models.StockOutboundProductsSub{}
	subList, err := stockOutboundSub.ListByOutProductIds(e.Orm, lo.Keys(orgProductMap))
	if err != nil {
		return err
	}
	subListMap := lo.Associate(*subList, func(item models.StockOutboundProductsSub) (int, models.StockOutboundProductsSub) {
		return item.Id, item
	})

	// 3. 查询出库单下的所有库位信息列表 | 只回显编号用
	stockLocation := &models.StockLocation{}
	locationList, err := stockLocation.StockLocationListByLwsCode(e.Orm, stockOutbound.LogicWarehouseCode)
	if err != nil {
		return err
	}
	locationMap := lo.Associate(*locationList, func(item *models.StockLocation) (int, *models.StockLocation) {
		return item.Id, item
	})

	// 4. 批量查询商品库存数量
	goodIds := lo.Map(c.StockOutboundProducts, func(item *dto.StockOutboundProductsConfirmReq, _ int) int {
		return item.GoodsId
	})
	stockInfo := models.StockInfo{}
	stockInfoList, err := stockInfo.ListByLwsCodeAndGoogIds(e.Orm, stockOutbound.LogicWarehouseCode, goodIds)
	if err != nil {
		return err
	}
	stockInfoMap := lo.Associate(*stockInfoList, func(item *models.StockInfo) (int, *models.StockInfo) {
		return item.GoodsId, item
	})
	fmt.Println(stockInfoMap)

	// 5. 批量查询库位库存数量
	locationIds := lo.Map(c.StockOutboundProducts, func(item *dto.StockOutboundProductsConfirmReq, _ int) int {
		return item.StockLocationId
	})
	stockLocationGoods := models.StockLocationGoods{}
	stockLocationGoodList, err := stockLocationGoods.ListByGoodsIdsAndLocationIds(e.Orm, goodIds, locationIds)
	if err != nil {
		return err
	}
	stockLocationGoodMap := lo.Associate(*stockLocationGoodList, func(item models.StockLocationGoods) (string, models.StockLocationGoods) {
		return strconv.Itoa(item.GoodsId) + "_" + strconv.Itoa(item.LocationId), item
	})

	// 6. 单个商品聚合本次实际出库数量 | id => ActQuantity
	groupActQuantity := map[int]int{}
	skuActQuantity := map[int]int{}

	// 7. 校验传参数据
	for _, reqProduct := range c.StockOutboundProducts { // 同一个商品多个库位
		orgProduct := orgProductMap[reqProduct.Id]
		if orgProduct.Id == 0 {
			return errors.New("出库单商品明细ID传参错误")
		}
		if orgProduct.GoodsId != reqProduct.GoodsId {
			return errors.New("商品ID传参错误")
		}
		if reqProduct.StockLocationId <= 0 {
			return errors.New(reqProduct.SkuCode + "库位必选")
		}
		if lo.IndexOf(lo.Keys(locationMap), reqProduct.StockLocationId) == -1 {
			return errors.New(reqProduct.SkuCode + "库位不在当前仓库范围，请检查")
		}

		locationInfo := locationMap[reqProduct.StockLocationId]
		if c.OutboundType != 2 && reqProduct.LocationActQuantity <= 0 { // 全部出库，允许出库数量为0
			return errors.New(reqProduct.SkuCode + "库位:" + locationInfo.LocationCode + "实际出库数量应为大于0的整数")
		}

		// 聚合单个商品不同库位[Id]的实际出库数量
		groupActQuantity[reqProduct.Id] += reqProduct.LocationActQuantity
		// 聚合单个商品不同库位[GoodsId]的实际出库数量
		skuActQuantity[reqProduct.GoodsId] += reqProduct.LocationActQuantity

		// 判断单个库位不能大于当前库位应出库数量
		subInfo := subListMap[reqProduct.SubId]
		if subInfo.Id == 0 {
			return errors.New("库位ID传参错误")
		}
		currQuantity := subInfo.LocationQuantity - subInfo.LocationActQuantity
		if reqProduct.LocationActQuantity > currQuantity {
			return fmt.Errorf("%v库位:%v实际出库数量不能大于当前应出库数量[%v]", reqProduct.SkuCode, locationInfo.LocationCode, currQuantity)
		}

		// 校验库位库存锁库是否充足
		key := strconv.Itoa(reqProduct.GoodsId) + "_" + strconv.Itoa(reqProduct.StockLocationId)
		stockLocationInfo := stockLocationGoodMap[key]
		if reqProduct.LocationActQuantity > stockLocationInfo.LockStock {
			return fmt.Errorf("SKU:%v库位:%v锁库数量不足[%v]", reqProduct.SkuCode, locationInfo.LocationCode, stockLocationInfo.LockStock)
		}
	}

	// 8. 校验商品锁库&汇总之后的出库数量
	quantityArr := lo.Map(stockOutbound.StockOutboundProducts, func(item models.StockOutboundProducts, _ int) int {
		return item.Quantity - item.ActQuantity
	})
	totalLeftQuantity := lo.Sum(quantityArr)
	currTotalQuantity := 0
	for productId, actQuantity := range groupActQuantity { // 商品ID汇总维度
		orgProduct := orgProductMap[productId]

		// 统计连续部分出库之后是否达成全部出库的条件
		currTotalQuantity += actQuantity

		// 校验汇总的实际出库数量不能大于当前应该出库数量 | 包含已经部分出库的
		currNeedQuantity := orgProduct.Quantity - orgProduct.ActQuantity
		if actQuantity > currNeedQuantity {
			return fmt.Errorf("SKU:%v汇总实际出库不能大于剩余未出库数量[%v]", orgProduct.SkuCode, currNeedQuantity)
		}

		// 用户选择全部出库 | 订单出库类型 | 数量不等于剩余未出库数量，不能全部出库
		if c.OutboundType == 2 && stockOutbound.IsOrderOutbound() && actQuantity != currNeedQuantity {
			return fmt.Errorf("SKU:%v当前出库数量[%v]没有满足应出库数量[%v],不能选择全部出库", orgProduct.SkuCode, actQuantity, currNeedQuantity)
		}

		// 校验商品库存锁库是否充足
		stockInfo := stockInfoMap[orgProduct.GoodsId]
		if actQuantity > stockInfo.LockStock {
			return fmt.Errorf("SKU:%v商品库存锁库数量不足[%v]", orgProduct.SkuCode, stockInfo.LockStock)
		}
	}

	// 系统判断是否全部出库
	allOutbound := false
	sysJudgmentAllOutbound := totalLeftQuantity == currTotalQuantity
	if c.OutboundType == 2 || sysJudgmentAllOutbound {
		allOutbound = true
	}

	// 9. 事务落库
	err = e.Orm.Transaction(func(tx *gorm.DB) error {
		// 修改调拨单/订单状态
		if stockOutbound.IsTransferOutbound() {
			transfer := &models.StockTransfer{}
			if err := transfer.GetByTransferCodeWithOptions(tx, stockOutbound.SourceCode, func(tx *gorm.DB) *gorm.DB {
				return tx.Preload("ToWarehouse")
			}); err != nil {
				return err
			}
			// 调拨单状态检查
			if lo.IndexOf([]string{"1", "4", "5"}, transfer.Status) == -1 {
				return fmt.Errorf("调拨单[%v]状态不符合(1,4,5),请检查[%v]", transfer.TransferCode, transfer.Status)
			}

			changeMap := map[string]map[bool]string{
				"1": { // 未出库未入库
					true:  "2", // 已出库未入库
					false: "4", // 部分出库未入库
				},
				"4": { // 部分出库未入库
					true:  "2", // 已出库未入库
					false: "4", // 部分出库未入库
				},
				"5": { // 部分出库部分入库
					true:  "6", // 全部出库部分入库
					false: "5", // 部分出库部分入库
				},
			}
			transferNewStatus := changeMap[transfer.Status][allOutbound]
			err = tx.Model(transfer).Where("status = ?", transfer.Status).Debug().Update("status", transferNewStatus).Error
			if err != nil {
				return err
			}

			// 虚拟出库仓自动入库
			if transfer.ToWarehouse.CheckIsVirtual() {
				// 校验必须是全部入库
				if !allOutbound {
					return errors.New("虚拟仓入库,必须是全部入库")
				}

				// 自动入库
				stockEntry := &models.StockEntry{}
				if err := stockEntry.GetByEntrySourceCodeWithOptions(tx, stockOutbound.SourceCode, func(tx *gorm.DB) *gorm.DB {
					return tx.Preload("StockEntryProducts")
				}); err != nil {
					return err
				}
				stockEntry.SetConfirmEntryStatus(stockOutbound.UpdateBy, stockOutbound.UpdateByName)
				if err := stockEntry.ConfirmEntryAutoForVirtual(tx); err != nil {
					return err
				}
				if err := transfer.SetEntryCompleteStatus(tx); err != nil {
					return err
				}
			}
		}

		// 修改订单状态
		if stockOutbound.IsOrderOutbound() {
			modelOc := &modelsOc.OrderInfo{}
			if err := modelOc.OutboundUpdateOrder(tx, stockOutbound.SourceCode, allOutbound, skuActQuantity); err != nil {
				return err
			}
		}

		// 修改出库单, 状态,最新出库时间，最后出库时间
		if allOutbound {
			stockOutbound.Status = "2"
			stockOutbound.OutboundEndTime = currTime
		} else {
			stockOutbound.Status = "3"
		}
		if stockOutbound.OutboundTime.IsZero() {
			stockOutbound.OutboundTime = currTime
		}
		err = tx.Where("id = ?", stockOutbound.Id).Save(&stockOutbound).Error
		if err != nil {
			return err
		}

		// 商品明细+商品库存更新
		for productId, actQuantity := range groupActQuantity { // 商品维度循环处理
			// 全部出库可以传0
			if actQuantity <= 0 {
				continue
			}

			// 更新出库单商品明细信息 | 状态，实际入库，首次入库时间，入库完成时间
			orgProduct := orgProductMap[productId]
			detailUpdate := map[string]any{
				"status":        "1",
				"act_quantity":  gorm.Expr("act_quantity + ?", actQuantity),
				"outbound_time": currTime,
			}
			if allOutbound {
				detailUpdate["status"] = "2"
				detailUpdate["outbound_end_time"] = currTime
			}
			err := tx.Model(&orgProduct).Where("act_quantity = ?", orgProduct.ActQuantity).Updates(detailUpdate).Error
			if err != nil {
				return err
			}

			// 扣减商品库存
			stockInfo := stockInfoMap[orgProduct.GoodsId]
			if err := stockInfo.SubLockStock(tx, actQuantity, &models.StockLog{
				DocketCode:   stockOutbound.OutboundCode,
				FromType:     models.StockLogFromType0,
				CreateBy:     stockOutbound.UpdateBy,
				CreateByName: stockOutbound.UpdateByName,
			}); err != nil {
				return err
			}

			// 完成出库释放多余商品锁库
			diffQuantity := orgProduct.Quantity - (orgProduct.ActQuantity + actQuantity)
			if allOutbound && diffQuantity > 0 { // 全部出库 && 存在差异
				if err := stockInfo.Unlock(tx, diffQuantity, &models.StockLockLog{
					DocketCode: stockOutbound.OutboundCode,
					FromType:   stockOutbound.Type,
					Remark:     models.RemarkManualConfirmOutboundProductsNotEqual,
				}); err != nil {
					return err
				}
			}
		}

		// 商品SUB+库位库存更新
		for _, reqProduct := range c.StockOutboundProducts { // SUB维度循环处理
			// 全部出库可以传0
			if reqProduct.LocationActQuantity <= 0 {
				continue
			}

			subInfo := subListMap[reqProduct.SubId]

			// 更新出库商品库位信息
			subUpdate := map[string]any{
				"location_act_quantity": gorm.Expr("location_act_quantity + ?", reqProduct.LocationActQuantity),
			}
			err := tx.Model(&subInfo).
				Where("id = ?", reqProduct.SubId).
				Where("location_act_quantity = ?", subInfo.LocationActQuantity).Updates(subUpdate).Error
			if err != nil {
				return err
			}

			// 记录部分出库日志
			subLog := &models.StockOutboundProductsSubLog{
				SubId:               subInfo.Id,
				LocationActQuantity: reqProduct.LocationActQuantity,
				OutboundTime:        currTime,
			}
			err = subLog.Create(tx)
			if err != nil {
				return err
			}

			// 扣减库位库存
			key := strconv.Itoa(reqProduct.GoodsId) + "_" + strconv.Itoa(reqProduct.StockLocationId)
			locationGoods := stockLocationGoodMap[key]
			if err := locationGoods.SubLockStock(tx, reqProduct.LocationActQuantity, &models.StockLocationGoodsLog{
				DocketCode:   stockOutbound.OutboundCode,
				FromType:     models.StockLocationGoodsLogFromType0,
				CreateBy:     stockOutbound.UpdateBy,
				CreateByName: stockOutbound.UpdateByName,
			}); err != nil {
				return err
			}

			// 实际出库数量不等于应出库数量(这里只可能是调拨单类型手动出库)，释放多锁住库位库存
			LocationDiffQuantity := subInfo.LocationQuantity - (subInfo.LocationActQuantity + reqProduct.LocationActQuantity)
			if allOutbound && LocationDiffQuantity > 0 { // 全部出库 && 库位实际出库总量 < 应出库总量
				if err := locationGoods.Unlock(tx, LocationDiffQuantity, &models.StockLocationGoodsLockLog{
					DocketCode: stockOutbound.OutboundCode,
					FromType:   stockOutbound.Type,
					Remark:     models.RemarkManualConfirmOutboundLocationProductsNotEqual,
				}); err != nil {
					return err
				}
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	// 记录接口日志
	beforeDataStr, _ := json.Marshal(&oldData)
	afterDataStr, _ := json.Marshal(&stockOutbound)
	dataStr, _ := json.Marshal(c)
	opLog := models.OperateLogs{
		DataId:       stockOutbound.OutboundCode,
		ModelName:    models.OutboundLogModelName,
		Type:         models.OutboundLogModelTypeConfirm,
		Before:       string(beforeDataStr),
		Data:         string(dataStr),
		After:        string(afterDataStr),
		OperatorId:   c.UpdateBy,
		OperatorName: c.UpdateByName,
	}
	_ = opLog.InsertItem(e.Orm)
	return nil
}

// ConfirmOutbound 确认出库 - 逐步废弃
func (e *StockOutbound) ConfirmOutbound(c *dto.StockOutboundConfirmReq, p *actions.DataPermission) error {

	// 查询出库单信息
	var stockOutbound = models.StockOutbound{}
	e.Orm.Preload("StockOutboundProducts").Preload("Warehouse").
		Scopes(
			actions.Permission(stockOutbound.TableName(), p),
			StockOutboundPermission(p),
		).First(&stockOutbound, c.GetId())
	if stockOutbound.Id == 0 {
		return errors.New("出库单不存在或没有权限")
	}
	oldData := stockOutbound
	if err := e.CheckConfirmOutbound(&stockOutbound); err != nil {
		return err
	}
	c.Generate(&stockOutbound)

	// 查询出库单产品信息&库位信息
	OutboundProductSubCustom := models.OutboundProductSubCustom{}
	subProducts, err := OutboundProductSubCustom.GetList(e.Orm, []string{stockOutbound.OutboundCode})

	if err != nil {
		return err
	}

	// 校验传参信息
	if err := e.CheckConfirmOutboundProducts(&subProducts, c.StockOutboundProducts, &stockOutbound); err != nil {
		return err
	}

	// 开启事务
	tx := e.Orm.Begin()
	defer func() {
		if r := recover(); r != nil {
			e.Log.Errorf("StockOutboundService ConfirmOutbound panic:%s \r\n", r)
			tx.Rollback()
		}
	}()
	if err := stockOutbound.ConfirmOutboundManual(tx, subProducts); err != nil {
		tx.Rollback()
		e.Log.Errorf("StockOutboundService ConfirmOutbound error:%s \r\n", err)
		return err
	}
	if err := tx.Commit().Error; err != nil {
		e.Log.Errorf("StockOutboundService ConfirmOutbound error:%s \r\n", err)
		return err
	}

	beforeDataStr, _ := json.Marshal(&oldData)
	afterDataStr, _ := json.Marshal(&stockOutbound)
	dataStr, _ := json.Marshal(c)
	opLog := models.OperateLogs{
		DataId:       stockOutbound.OutboundCode,
		ModelName:    models.OutboundLogModelName,
		Type:         models.OutboundLogModelTypeConfirm,
		Before:       string(beforeDataStr),
		Data:         string(dataStr),
		After:        string(afterDataStr),
		OperatorId:   c.UpdateBy,
		OperatorName: c.UpdateByName,
	}
	_ = opLog.InsertItem(e.Orm)
	return nil
}

// CheckConfirmOutbound 出库单状态检查
func (e *StockOutbound) CheckConfirmOutbound(data *models.StockOutbound) error {
	if lo.IndexOf([]string{"1", "3"}, data.Status) == -1 {
		return errors.New("出库单状态不正确")
	}
	if data.Warehouse.CheckIsVirtual() {
		return errors.New("虚拟仓无需手动出库")
	}
	if data.IsOrderOutbound() {
		// 检查是否有售后 todo
		if models.CheckOrderIsAfterPendingFromOc(e.Orm, data.SourceCode) {
			return errors.New("有正在处理中的售后，无法出库")
		}
	}
	return nil
}

// CheckConfirmOutbound 确认出库check
func (e *StockOutbound) CheckConfirmOutboundProducts(list *[]models.OutboundProductSubCustom, reqList []dto.StockOutboundProductsConfirmReq, stockOutbound *models.StockOutbound) error {
	for index, item := range *list {
		reqItem := FindInProductsConfirmReqListForOutbound(reqList, item.SkuCode, item.StockLocationId)
		if reqItem == nil {
			return errors.New("参数异常,缺少" + item.SkuCode + ",库位id:" + strconv.Itoa(item.StockLocationId) + "信息")
		}
		if reqItem.LocationActQuantity <= 0 {
			return errors.New(reqItem.SkuCode + ",库位:" + item.StockLocation.LocationCode + "实际出库数量应为为大于0的整数")
		}
		if item.LocationQuantity < reqItem.LocationActQuantity {
			return errors.New(reqItem.SkuCode + ",库位:" + item.StockLocation.LocationCode + "实际出库数量不能大于应出库数量")
		}
		if stockOutbound.IsOrderOutbound() && item.LocationQuantity != reqItem.LocationActQuantity {
			return errors.New(reqItem.SkuCode + ",库位:" + item.StockLocation.LocationCode + "领用单的实际出库数量要等于应出库数量")
		}
		// 赋值实际出库数量
		(*list)[index].LocationActQuantity = reqItem.LocationActQuantity
	}
	return nil
}

func FindInProductsConfirmReqListForOutbound(list []dto.StockOutboundProductsConfirmReq, sku string, location int) *dto.StockOutboundProductsConfirmReq {
	for _, item := range list {
		if item.SkuCode == sku && item.StockLocationId == location {
			return &item
		}
	}
	return nil
}

// 出库单打印html
func (e *StockOutbound) PrintOutbound(d *dto.StockOutboundPrintReq, p *actions.DataPermission, outData *dto.StockOutboundPrintResp) error {
	req := &dto.StockOutboundGetReq{
		Id: d.Id,
	}
	resp := &dto.StockOutboundNoLocationResp{}
	if err := e.GetNoLocation(req, p, resp); err != nil {
		e.Log.Errorf("StockOutboundService PrintOutbound error:%s \r\n", err)
		return err
	}
	outDataMap := utils.GetChunkProductsForPrint(resp, "stockOutboundProducts", map[string]int{
		"productName": 10,
		"mfgModel":    8,
		"brandName":   6,
		"vendorName":  8,
		"productNo":   10,
	}, 28, 50)
	path := "static/viewtpl/wc/outbound_print.html"
	resStr, err := utils.View(path, outDataMap)
	if err != nil {
		return err
	}
	outData.Html = resStr
	return nil
}

// 拣货单打印html
func (e *StockOutbound) PrintPicking(d *dto.StockOutboundPrintReq, p *actions.DataPermission, outData *dto.StockOutboundPrintResp) error {
	req := &dto.StockOutboundGetReq{
		Id: d.Id,
	}
	resp := &dto.StockOutboundGetResp{}
	if err := e.Get(req, p, resp); err != nil {
		e.Log.Errorf("StockOutboundService PrintPicking error:%s \r\n", err)
		return err
	}
	// 部分出库-打印剩余未出库数量
	for _, item := range resp.StockOutboundProducts {
		item.OutboundProductSubCustom.LocationQuantity = item.OutboundProductSubCustom.LocationQuantity - item.OutboundProductSubCustom.LocationActQuantity
	}

	outDataMap := utils.GetChunkProductsForPrint(resp, "stockOutboundProducts", map[string]int{
		"productName":  12,
		"mfgModel":     6,
		"brandName":    6,
		"locationCode": 10,
		"vendorName":   8,
	}, 39, 46)
	path := "static/viewtpl/wc/outbound_print_picking.html"
	if resp.Type == "1" { // 订单出库新增字段
		path = "static/viewtpl/wc/outbound_print_picking_1.html"
	}
	resStr, err := utils.View(path, outDataMap)
	fmt.Println(resStr)
	if err != nil {
		return err
	}
	outData.Html = resStr
	outData.Data = outDataMap
	return nil
}

// GetOrderListByUserName 获取领用人的订单
func (e *StockOutbound) GetOrderListByUserName(userName string) map[string]string {
	orderResult := ocClient.ApiByDbContext(e.Orm).GetOrderListByUserName(userName)
	orderResultList := &struct {
		response.Response
		Data []ocDto.OrderListResp
	}{}
	orderResult.Scan(orderResultList)
	orderListMap := make(map[string]string, len(orderResultList.Data))
	for _, order := range orderResultList.Data {
		orderListMap[order.OrderId] = order.UserName
	}
	return orderListMap
}
