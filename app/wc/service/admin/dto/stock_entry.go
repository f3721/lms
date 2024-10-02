package dto

import (
	StockDto "go-admin/common/dto/stock/dto"
	"go-admin/common/utils"
	"time"

	"github.com/jinzhu/copier"
	"github.com/samber/lo"

	"gorm.io/gorm"

	"go-admin/app/wc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type StockEntryGetPageReq struct {
	dto.Pagination     `search:"-"`
	Ids                []int     `form:"ids[]"  search:"-"`
	EntryCode          string    `form:"entryCode"  search:"type:exact;column:entry_code;table:stock_entry" comment:"入库单编码"`
	Type               string    `form:"type"  search:"type:exact;column:type;table:stock_entry" comment:"入库类型:  0 大货入库  1 退货入库"`
	SourceType         string    `form:"sourceType"  search:"-"`
	Status             string    `form:"status"  search:"type:exact;column:status;table:stock_entry" comment:"状态:0-已作废 1-创建 2-已完成"`
	SourceCode         string    `form:"sourceCode"  search:"type:exact;column:source_code;table:stock_entry" comment:"来源单据code"`
	EntryTime          time.Time `form:"entryTime"  search:"type:exact;column:entry_time;table:stock_entry" comment:"入库时间"`
	WarehouseCode      string    `form:"warehouseCode"  search:"type:exact;column:warehouse_code;table:stock_entry" comment:"实体仓code"`
	LogicWarehouseCode string    `form:"logicWarehouseCode"  search:"type:exact;column:logic_warehouse_code;table:stock_entry" comment:"逻辑仓code"`
	CreatedAtStart     time.Time `form:"createdAtStart"  search:"-" comment:"" time_format:"2006-01-02 15:04:05"`
	CreatedAtEnd       time.Time `form:"createdAtEnd"  search:"-" comment:"" time_format:"2006-01-02 15:04:05"`
	VendorId           int       `form:"vendorId"  search:"type:exact;column:vendor_id;table:stock_entry" comment:"货主id"`
	SupplierId         int       `form:"supplierId"  search:"type:exact;column:supplier_id;table:stock_entry" comment:"供应商id"`
	Recipient          string    `form:"recipient" search:"-"`
	StockDto.ProductSearch
}

type AddStockEntryReq1 struct {
	//入库单
	StockEntryInsertReq
	//供应商 - vendors
	VendorId int    `json:"vendorId" comment:"货主id"`
	NameZh   string `json:"name_zh" comment:"货主名称"`

	//实体仓
	Province int    `json:"province"`
	City     int    `json:"city"`
	District int    `json:"district" comment:"区"`
	Address  string `json:"address"`
	Remark   string `json:"Remark"`

	//sku - stock_entry_products

	//其它类型
	TypeName           string `json:"typeName"`
	SourceTypeName     string `json:"sourceTypeName"`
	StatusName         string `json:"statusName"`
	WarehouseName      string `json:"warehouseName"`
	LogicWarehouseName string `json:"logicWarehouseName"`
	DiffNum            string `json:"diffNum"`
	//Address            string                      `json:"address"`
	//District           int                         `json:"district"`
	//City               int                         `json:"city"`
	//Province           int                         `json:"province"`
	DistrictName       string                      `json:"districtName"`
	CityName           string                      `json:"cityName"`
	ProvinceName       string                      `json:"provinceName"`
	ExternalOrderNo    string                      `json:"externalOrderNo"`
	StockEntryProducts []StockEntryProductsGetResp `json:"stockEntryProducts" copier:"-"`

	common.ControlBy
}

type StockEntry struct {
	EntryCode string `json:"entryCode"  comment:"入库单编码"`
	Type      string `json:"type" binding:"required" comment:"入库类型:  0 大货入库  1 退货入库"`
	VendorId  int    `json:"vendorId" comment:"货主id"`
	//SourceType         string    `json:"sourceType"`
	//Status string `json:"status"  comment:"状态:0-已作废 1-创建 2-已完成"`
	//SourceCode         string    `json:"sourceCode"  comment:"来源单据code"`
	//EntryTime          time.Time `json:"entryTime"  comment:"入库时间"`
	WarehouseCode      string `json:"warehouseCode"  comment:"实体仓code"`
	LogicWarehouseCode string `json:"logicWarehouseCode"  comment:"逻辑仓code"`
}

func (m *StockEntryGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type StockEntryGetPageResp struct {
	models.StockEntry
	SupplierName       string                    `json:"supplierName"`
	TypeName           string                    `json:"typeName"`
	SourceTypeName     string                    `json:"sourceTypeName"`
	WarehouseName      string                    `json:"warehouseName"`
	LogicWarehouseName string                    `json:"logicWarehouseName"`
	DiffNum            string                    `json:"diffNum"`
	NameZh             string                    `json:"nameZh"`
	CheckStatus        string                    `json:"checkStatus"`
	CheckStatusName    string                    `json:"checkStatusName"`
	StatusName         string                    `json:"statusName"`
	OptionRulesEntry   StockDto.OptionRulesEntry `json:"optionRules"`
	//OptionRulesTransfer StockDto.OptionRulesTransfer `json:"optionRules"`
}

func (m *StockEntryGetPageResp) SetEntryRulesByStatus() {
	// casbin 设置 todo
	m.OptionRulesEntry.Confirm = true
	m.OptionRulesEntry.Detail = true
	m.OptionRulesEntry.Print = true
	switch m.Status {
	case models.EntryStatus0:
		m.OptionRulesEntry.Confirm.False()
		m.OptionRulesEntry.Detail.True()
		m.OptionRulesEntry.Print.False()
	case models.EntryStatus1:
		m.OptionRulesEntry.Confirm.True()
		m.OptionRulesEntry.Detail.True()
		m.OptionRulesEntry.Print.False()
	case models.EntryStatus2:
		m.OptionRulesEntry.Confirm.False()
		m.OptionRulesEntry.Print.False()
		m.OptionRulesEntry.Detail.True()
	}
}

func (s *StockEntryGetPageResp) InitData() {
	s.SetEntryRulesByStatus()
	s.WarehouseName = s.Warehouse.WarehouseName
	s.LogicWarehouseName = s.LogicWarehouse.LogicWarehouseName
	s.TypeName = utils.GetTFromMap(s.Type, models.EntryTypeMap)
	s.StatusName = utils.GetTFromMap(s.Status, models.EntryStatusMap)
	s.CheckStatusName = utils.GetTFromMap(s.CheckStatus, models.EntryCheckStatusMap)
	s.SourceTypeName = utils.GetTFromMap(s.Type, models.EntrySourceTypeMap)
}

type StockEntryInsertReq struct {
	SourceCode         string `json:"sourceCode" comment:"来源单据code"`
	Remark             string `json:"remark" comment:"备注"`
	WarehouseCode      string `json:"warehouseCode" comment:"实体仓code"`
	LogicWarehouseCode string `json:"logicWarehouseCode" comment:"逻辑仓code"`
	VendorId           int    `json:"vendorId" comment:"货主id"`
	common.ControlBy
	StockEntryProducts []StockEntryProductsReq `json:"stockEntryProducts"`
}

type AddStockEntryReq struct {
	Id int `json:"id"`

	//EntryCode string `json:"entryCode"`
	//SkuCode            string    `json:"skuCode"`
	//SourceCode         string    `json:"sourceCode" comment:"来源单号"`
	//Type               string `json:"type" binding:"required" comment:"入库类型:  0 大货入库  1 退货入库  2 其他 3采购入库"`
	//Status             string `json:"status" binding:"required" comment:"状态:0-已作废 1-未出库未入库 2-已出库未入库 4-部分出库未入库 5-部分出库部分入库 6-全部出库部分入库  3-出库完成入库完成 98已提交 99未提交"`
	SupplierId         int    `json:"supplierId" binding:"required" comment:"供应商id"`
	WarehouseCode      string `json:"warehouseCode" binding:"required" comment:"出库实体仓code"`
	LogicWarehouseCode string `json:"logicWarehouseCode" binding:"required" comment:"出库逻辑仓code"`
	Remark             string `json:"remark" comment:"备注"`
	VendorId           int    `json:"vendorId" binding:"required" comment:"货主id"`

	AddStockEntryProducts []AddStockEntryProductsReq `json:"stockEntryProducts"`
	common.ControlBy
}

type StockEntryResp struct {
	models.StockEntry
	TypeName           string                      `json:"typeName"`
	SourceTypeName     string                      `json:"sourceTypeName"`
	StatusName         string                      `json:"statusName"`
	WarehouseName      string                      `json:"warehouseName"`
	LogicWarehouseName string                      `json:"logicWarehouseName"`
	DiffNum            string                      `json:"diffNum"`
	Address            string                      `json:"address"`
	District           int                         `json:"district"`
	City               int                         `json:"city"`
	Province           int                         `json:"province"`
	DistrictName       string                      `json:"districtName"`
	CityName           string                      `json:"cityName"`
	ProvinceName       string                      `json:"provinceName"`
	ExternalOrderNo    string                      `json:"externalOrderNo"`
	StockEntryProducts []StockEntryProductsGetResp `json:"stockEntryProducts" copier:"-"`
}

type EditStockEntryReq struct {
	Id int `json:"id"`
	//EntryCode string `json:"entryCode" binding:"required"`
	//SourceCode         string `json:"sourceCode" comment:"来源单号"`
	//Type               string `json:"type" binding:"required" comment:"入库类型:  0 大货入库  1 退货入库  2 其他  3采购入库"`
	//Status             string `json:"status" binding:"required" comment:"状态:0-已作废 1-未出库未入库 2-已出库未入库 4-部分出库未入库 5-部分出库部分入库 6-全部出库部分入库  3-出库完成入库完成 98已提交 99未提交"`
	//CheckStatus        string `json:"checkStatus" comment:"审核状态：1已提交，2待审核，3入库完成"`
	VendorId           int    `json:"vendorId" binding:"required" comment:"货主id"`
	SupplierId         int    `json:"supplierId" binding:"required" comment:"供应商id"`
	WarehouseCode      string `json:"warehouseCode" binding:"required" comment:"出库实体仓code"`
	LogicWarehouseCode string `json:"logicWarehouseCode" binding:"required" comment:"出库逻辑仓code"`
	Remark             string `json:"remark" comment:"备注"`
	StockLocationId    int    `json:"stockLocationId" comment:"StockLocationId"`

	AddStockEntryProducts []AddStockEntryProductsReq `json:"stockEntryProducts"`
	//AddStockEntryProductsSub []StockEntryProductsSubReq `json:"stockEntryProductsSub"`
	common.ControlBy
}

type CheckStockEntryReq struct {
	Id          int    `uri:"id" comment:"id"` // id
	CheckRemark string `json:"checkRemark"`
	//CheckStatus string `json:"checkStatus" vd:"len($)>0; msg:'审核通过与否必填'"`
	CheckStatus string `json:"checkStatus" vd:"@:len($)>0; msg:'供应商编码长度在0-10之间'"`

	common.ControlBy
}

func (s *CheckStockEntryReq) GetId() interface{} {
	return s.Id
}

func (s *EditStockEntryReq) GetId() interface{} {
	return s.Id
}

//func (s *AddStockEntryReq) GetId() interface{} {
//	return s.Id
//}

func (s *StockEntryInsertReq) Generate(model *models.StockEntry) {
	model.Status = models.EntryStatus1
	model.SourceCode = s.SourceCode
	model.Remark = s.Remark
	model.WarehouseCode = s.WarehouseCode
	model.LogicWarehouseCode = s.LogicWarehouseCode
	model.VendorId = s.VendorId
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.CreateByName = s.CreateByName
	for _, item := range s.StockEntryProducts {
		item.SetCreateBy(s.CreateBy)
		item.SetCreateByName(s.CreateByName)
		modelStockEntryProducts := models.StockEntryProducts{}
		item.Generate(&modelStockEntryProducts)
		model.StockEntryProducts = append(model.StockEntryProducts, modelStockEntryProducts)
	}
}

func (s *AddStockEntryReq) Generate(tx *gorm.DB, model *models.StockEntry) error {
	//currTime := time.Now()
	copier.Copy(model, s)
	if model.EntryCode == "" {
		if _, err := model.GenerateEntryCode(tx); err != nil {
			return err
		}
	}
	//model.EntryTime = currTime
	model.CreateBy = s.CreateBy
	model.CreateByName = s.CreateByName
	model.Type = "3"   //采购入库
	model.Status = "1" //创建
	model.CheckStatus = "1"
	for _, item := range s.AddStockEntryProducts {
		modelStockEntryProducts := models.StockEntryProducts{}
		copier.Copy(&modelStockEntryProducts, &item)
		modelStockEntryProducts.EntryCode = model.EntryCode
		//modelStockEntryProducts.EntryTime = currTime
		modelStockEntryProducts.VendorId = model.VendorId
		modelStockEntryProducts.WarehouseCode = model.WarehouseCode
		modelStockEntryProducts.LogicWarehouseCode = model.LogicWarehouseCode
		//model.StockEntryProducts = append(model.StockEntryProducts, modelStockEntryProducts)

		//先处理sub
		subS := []models.StockEntryProductsSub{}
		tmpQuantity := 0
		for _, itemsub := range item.StockEntryProductsSub {
			modelStockEntryProductsSub := models.StockEntryProductsSub{}
			copier.Copy(&modelStockEntryProductsSub, &itemsub)
			//modelStockEntryProductsSub.EntryTime = currTime
			modelStockEntryProductsSub.EntryCode = model.EntryCode
			tmpQuantity += itemsub.StashActQuantity
			subS = append(subS, modelStockEntryProductsSub)
		}

		modelStockEntryProducts.Quantity = tmpQuantity
		modelStockEntryProducts.StockEntryProductsSub = subS
		model.StockEntryProducts = append(model.StockEntryProducts, modelStockEntryProducts)
	}
	return nil
}

type StockEntryCancelForCancelCsOrderReq struct {
	CsNo   string `json:"csNo" comment:"售后单号"`
	Remark string `json:"remark" comment:"备注"`
	common.ControlBy
}

func (s *StockEntryCancelForCancelCsOrderReq) Generate(model *models.StockEntry) {
	model.Status = models.EntryStatus0
	model.Remark = s.Remark
	model.UpdateBy = s.UpdateBy
	model.UpdateByName = s.UpdateByName
}

type StockEntryUpdateReq struct {
	Id                 int       `uri:"id" comment:"id"` // id
	EntryCode          string    `json:"entryCode" comment:"入库单编码"`
	Type               string    `json:"type" comment:"入库类型:  0 大货入库  1 退货入库  2 其他"`
	Status             string    `json:"status" comment:"状态:0-已作废 1-创建 2-已完成"`
	SourceCode         string    `json:"sourceCode" comment:"来源单据code"`
	Remark             string    `json:"remark" comment:"备注"`
	EntryTime          time.Time `json:"entryTime" comment:"入库时间"`
	WarehouseCode      string    `json:"warehouseCode" comment:"实体仓code"`
	LogicWarehouseCode string    `json:"logicWarehouseCode" comment:"逻辑仓code"`
	VendorId           int       `json:"vendorId" comment:"货主id"`
	CreateByName       string    `json:"createByName" comment:""`
	UpdateByName       string    `json:"updateByName" comment:""`
	common.ControlBy
}

func (s *StockEntryUpdateReq) Generate(model *models.StockEntry) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.EntryCode = s.EntryCode
	model.Type = s.Type
	model.Status = s.Status
	model.SourceCode = s.SourceCode
	model.Remark = s.Remark
	model.EntryTime = s.EntryTime
	model.WarehouseCode = s.WarehouseCode
	model.LogicWarehouseCode = s.LogicWarehouseCode
	model.VendorId = s.VendorId
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
	model.CreateByName = s.CreateByName
	model.UpdateByName = s.UpdateByName
}

func (s *StockEntryUpdateReq) GetId() interface{} {
	return s.Id
}

// StockEntryGetReq 功能获取请求参数
type StockEntryGetReq struct {
	Id   int    `uri:"id"`
	Type string `form:"type"`
}

func (s *StockEntryGetReq) GetId() interface{} {
	return s.Id
}

type StockEntryGetResp struct {
	models.StockEntry
	TypeName           string                      `json:"typeName"`
	SourceTypeName     string                      `json:"sourceTypeName"`
	SupplierName       string                      `json:"SupplierName"`
	CheckStatusName    string                      `json:"checkStatusName"`
	StatusName         string                      `json:"statusName"`
	WarehouseName      string                      `json:"warehouseName"`
	LogicWarehouseName string                      `json:"logicWarehouseName"`
	DiffNum            string                      `json:"diffNum"`
	Address            string                      `json:"address"`
	District           int                         `json:"district"`
	City               int                         `json:"city"`
	Province           int                         `json:"province"`
	DistrictName       string                      `json:"districtName"`
	CityName           string                      `json:"cityName"`
	ProvinceName       string                      `json:"provinceName"`
	ExternalOrderNo    string                      `json:"externalOrderNo"`
	StockEntryProducts []StockEntryProductsGetResp `json:"stockEntryProducts" copier:"-"`
}

type OutputStockEntryProducts struct {
	Number       int    `json:"number"`
	SkuCode      string `json:"skuCode"`
	LocationCode string `json:"locationCode"`
	ProductName  string `json:"productName"`
	MfgModel     string `json:"mfgModel"`
	BrandName    string `json:"brandName"`
	VendorName   string `json:"vendorName"`
	ActQuantity  int    `json:"actQuantity"`
	SalesUom     string `json:"salesUom"`
}

// 逐步废弃-禁止迭代使用
func (s *StockEntryGetResp) InitData(tx *gorm.DB, data *models.StockEntry, Type string) error {
	stockLocationIds := []int{}
	s.WarehouseName = s.Warehouse.WarehouseName
	s.LogicWarehouseName = s.LogicWarehouse.LogicWarehouseName
	s.TypeName = utils.GetTFromMap(s.Type, models.EntryTypeMap)
	s.StatusName = utils.GetTFromMap(s.Status, models.EntryStatusMap)
	s.SourceTypeName = utils.GetTFromMap(s.Type, models.EntrySourceTypeMap)
	// 获取goods product 信息通过pc
	goodsIdSlice := lo.Uniq(lo.Map(data.StockEntryProducts, func(item models.StockEntryProducts, _ int) int {
		return item.GoodsId
	}))
	goodsProductMap := models.GetGoodsProductInfoFromPc(tx, goodsIdSlice)

	// 获取vendor map
	vendorIdSlice := lo.Uniq(lo.Map(data.StockEntryProducts, func(item models.StockEntryProducts, _ int) int {
		return item.VendorId
	}))
	vendorIdMap := models.GetVendorsMapByIds(tx, vendorIdSlice)

	// 查询拆分入库记录
	entryProductIds := lo.Map(data.StockEntryProducts, func(item models.StockEntryProducts, _ int) int {
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

	for index, item := range data.StockEntryProducts {
		productGetResp := StockEntryProductsGetResp{}
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
				stockEntryProductsLocationGetResp := StockEntryProductsLocationGetResp{}
				for _, item := range *stockLocations {
					stockEntryProductsLocationGetResp.Regenerate(item)
					productGetResp.StockLocation = append(productGetResp.StockLocation, stockEntryProductsLocationGetResp)
				}
				if hasTop && len(*stockLocations) > 0 {
					productGetResp.StockLocationId = (*stockLocations)[0].Id
				}
			}
		}

		if Type == "confirm" || Type == "print" {
			// if data.IsCsOrderEntry() {
			// 	productGetResp.ActQuantity = productGetResp.Quantity
			// }
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

		s.StockEntryProducts = append(s.StockEntryProducts, productGetResp)
	}
	stockLocationMap, err := models.GetStockLocationMapByIds(tx, stockLocationIds)

	for index, item := range s.StockEntryProducts {
		if locationCode, ok := stockLocationMap[item.StockLocationId]; ok {
			s.StockEntryProducts[index].LocationCode = locationCode
		}
	}

	if err != nil {
		return err
	}
	return nil
}

type StockEntryProductsGetResp struct {
	models.StockEntryProducts
	StockDto.ProductGoodsResp
	StockLocation         []StockEntryProductsLocationGetResp `json:"stockLocations"  copier:"-"`
	StockEntryProductsSub []models.StockEntryProductsSub      `json:"stockEntryProductsSub" copier:"-"`
	LocationCode          string                              `json:"locationCode"`
	DiffNum               int                                 `json:"diffNum"`
	Number                int                                 `json:"number"`
	OutboundNum           int                                 `json:"outboundNum"`
}

type StockEntryProductsLocationGetResp struct {
	Id           int    `json:"id"`
	LocationCode string `json:"locationCode"`
}

func (s *StockEntryProductsLocationGetResp) Regenerate(location models.StockLocation) {
	s.Id = location.Id
	s.LocationCode = location.LocationCode
}

type StockEntryProductsLocationGetReq struct {
	LocationCode       string `form:"locationCode"`
	LogicWarehouseCode string `form:"logicWarehouseCode"`
	IsDefective        int    `form:"isDefective"`
}

// StockEntryDeleteReq 功能删除请求参数
type StockEntryDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *StockEntryDeleteReq) GetId() interface{} {
	return s.Ids
}

// 确认入库入参
type StockEntryConfirmReq struct {
	Id                 int                            `json:"id"`
	StockEntryProducts []StockEntryProductsConfirmReq `json:"stockEntryProducts"`
	common.ControlBy
}

func (s *StockEntryConfirmReq) GetId() interface{} {
	return s.Id
}

func (s *StockEntryConfirmReq) Generate(model *models.StockEntry) {
	model.Status = models.EntryStatus2
	model.EntryTime = time.Now()
	model.UpdateBy = s.UpdateBy
	model.UpdateByName = s.UpdateByName
	for index := range model.StockEntryProducts {
		model.StockEntryProducts[index].UpdateBy = s.UpdateBy
		model.StockEntryProducts[index].UpdateByName = s.UpdateByName

	}
}

// 部分入库入参
type StockEntryPartReq struct {
	Id                 int                   `json:"id"`                 // 入库单ID
	EntryType          int                   `json:"type" `              // 0-保存 , 1-部分入库, 2-全部入库
	StockEntryProducts []*StockEntryProducts `json:"stockEntryProducts"` // 产品信息列表
	common.ControlBy
}
type StockEntryProducts struct {
	Id                     int                       `json:"id"`               // 入库单产品明细id
	SkuCode                string                    `json:"skuCode"`          // 入库单商品sku
	GoodId                 int                       `json:"goodId"`           // 入库单商品ID
	Quantity               int                       `json:"quantity"`         // 应入库数量
	ActQuantityTotal       int                       `json:"actQuantityTotal"` // 总实际入库数量
	WarehouseCode          string                    `json:"warehouseCode"`
	StockEntryLocationInfo []*StockEntryLocationInfo `json:"stockLocationInfo"` // 库位ID+实际入库数量
}

type StockEntryLocationInfo struct {
	Id              int `json:"id"`              // 拆分行ID
	StockLocationId int `json:"stockLocationId"` // 库位ID
	ActQuantity     int `json:"actQuantity"`     // 实际入库数量
}

type StockEntryStashReq struct {
	Id                 int                            `json:"id"`
	StockEntryProducts []StockEntryProductsConfirmReq `json:"stockEntryProducts"`
	common.ControlBy
}

func (s *StockEntryStashReq) GetId() interface{} {
	return s.Id
}

func (s *StockEntryStashReq) Generate(model *models.StockEntry) {
	for index, _ := range model.StockEntryProducts {
		model.StockEntryProducts[index].UpdateBy = s.UpdateBy
		model.StockEntryProducts[index].UpdateByName = s.UpdateByName

	}
}

type StockEntryPrintHtmlReq struct {
	Id int `uri:"id"`
}

func (s *StockEntryPrintHtmlReq) GetId() interface{} {
	return s.Id
}

type StockEntryPrintHtmlResp struct {
	Html string `json:"html"`
}

type StockEntryPrintSkuReq struct {
	GoodsId int `uri:"goodsId"`
}
type StockEntryPrintSkusReq struct {
	GoodsIds []int `json:"goodsIds"`
}

func (s *StockEntryPrintSkuReq) GetGoodsId() interface{} {
	return s.GoodsId
}

type StockEntryPrintSkuResp struct {
	StockEntryScanSkuBaseInfoResp
	QrCode string `json:"qrCode"`
}

type StockEntryScanSkuBaseInfoResp struct {
	VendorCode  string `json:"vendorCode"`
	VendorName  string `json:"vendorName"`
	SkuCode     string `json:"skuCode"`
	ProductNo   string `json:"productNo"`
	BrandName   string `json:"brandName"`
	MfgModel    string `json:"mfgModel"`
	ProductName string `json:"productName"`
	SalesUom    string `json:"salesUom"`
}

type StockEntryPrintSkusResp struct {
	VendorCode  string `json:"vendorCode"`
	VendorName  string `json:"vendorName"`
	SkuCode     string `json:"skuCode"`
	ProductNo   string `json:"productNo"`
	BrandName   string `json:"brandName"`
	MfgModel    string `json:"mfgModel"`
	ProductName string `json:"productName"`
	SalesUom    string `json:"salesUom"`
	QrCode      string `json:"qrCode"`
}

type StockEntryExportResp struct {
	models.StockEntry
	TypeName                  string                      `json:"typeName"`
	SourceTypeName            string                      `json:"sourceTypeName"`
	StatusName                string                      `json:"statusName"`
	WarehouseName             string                      `json:"warehouseName"`
	LogicWarehouseName        string                      `json:"logicWarehouseName"`
	Recipient                 string                      `json:"recipient"` // 领用人
	OptionRulesEntry          StockDto.OptionRulesEntry   `json:"optionRules" gorm:"-"`
	NameZh                    string                      `json:"nameZh" comment:货主id"`
	StockEntryProductsGetResp []StockEntryProductsGetResp `json:"stockEntryProducts" copier:"-"  gorm:"-"`
}

type StockEntryExport struct {
	TypeName           string `json:"typeName"`
	SourceTypeName     string `json:"sourceTypeName"`
	StatusName         string `json:"statusName"`
	WarehouseName      string `json:"warehouseName"`
	LogicWarehouseName string `json:"logicWarehouseName"`
	Recipient          string `json:"recipient"`
	EntryCode          string `json:"entryCode"`
	SourceCode         string `json:"sourceCode"`
	Remark             string `json:"remark"`
	EntryTime          string `json:"entryTime"`
	SkuCode            string `json:"skuCode"`
	ProductName        string `json:"productName"`
	MfgModel           string `json:"mfgModel"`
	BrandName          string `json:"brandName"`
	SalesUom           string `json:"salesUom"`
	VendorName         string `json:"vendorName"`
	ProductNo          string `json:"productNo"`
	VendorSkuCode      string `json:"vendorSkuCode"`
	Quantity           int    `json:"Quantity"`
	ActQuantity        int    `json:"ActQuantity"`
	LocationCode       string `json:"locationCode"`
	DiffNum            int    `json:"diffNum"`
	IsDefective        string `json:"isDefective"`
}

type StockEntryferValidateSkusReq struct {
	SkuCodes string `form:"skuCodes" vd:"@:len($)>0; msg:'skuCodes不能为空'"`
	VendorId int    `form:"vendorId" vd:"@:$>0; msg:'vendorId不能为空'"`
	//StockLocationId    int    `form:"stockLocationId" vd:"?"`
	WarehouseCode      string `form:"warehouseCode" vd:"@:len($)>0; msg:'实体仓不能为空'"`
	LogicWarehouseCode string `form:"logicWarehouseCode" vd:"@:len($)>0; msg:'逻辑仓不能为空'"`
}

func (m *StockEntryGetPageResp) SetStockEntryRulesByStatus() {
	// 审核按钮 | 已提交状态下
	if m.CheckStatus == "1" {
		m.OptionRulesEntry.Audit.True()
	}
	if m.Type == "3" {
		m.OptionRulesEntry.Confirm.False()
	}

}
