package dto

import (
	StockDto "go-admin/common/dto/stock/dto"
	"go-admin/common/utils"
	"time"

	"github.com/samber/lo"
	"gorm.io/gorm"

	"go-admin/app/wc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type StockTransferGetPageReq struct {
	dto.Pagination         `search:"-"`
	TransferCode           string    `form:"transferCode"  search:"type:exact;column:transfer_code;table:stock_transfer" comment:"调拨单编码"`
	Type                   string    `form:"type"  search:"type:exact;column:type;table:stock_transfer" comment:"调拨类型: 0 -正常调拨 1-次品调拨 2-正转次调拨 3-次转正调拨"`
	Status                 string    `form:"status"  search:"type:exact;column:status;table:stock_transfer" comment:"状态:0-已作废 1-未出库未入库 2-已出库未入库 3-出库完成入库完成 98已提交 99未提交"`
	FromWarehouseCode      string    `form:"fromWarehouseCode"  search:"type:exact;column:from_warehouse_code;table:stock_transfer" comment:"出库实体仓code"`
	FromLogicWarehouseCode string    `form:"fromLogicWarehouseCode"  search:"type:exact;column:from_logic_warehouse_code;table:stock_transfer" comment:"出库逻辑仓code"`
	ToWarehouseCode        string    `form:"toWarehouseCode"  search:"type:exact;column:to_warehouse_code;table:stock_transfer" comment:"入库实体仓code"`
	ToLogicWarehouseCode   string    `form:"toLogicWarehouseCode"  search:"type:exact;column:to_logic_warehouse_code;table:stock_transfer" comment:"入库逻辑仓code"`
	SourceCode             string    `form:"sourceCode"  search:"type:exact;column:source_code;table:stock_transfer" comment:"来源单号"`
	VerifyStatus           string    `form:"verifyStatus"  search:"type:exact;column:verify_status;table:stock_transfer" comment:"审核状态 0 待审核 1 审核通过 2 审核驳回 99初始化"`
	VendorId               int       `form:"vendorId"  search:"type:exact;column:vendor_id;table:stock_transfer" comment:"货主id"`
	CreatedAtStart         time.Time `form:"createdAtStart"  search:"-" comment:"" time_format:"2006-01-02 15:04:05"`
	CreatedAtEnd           time.Time `form:"createdAtEnd"  search:"-" comment:"" time_format:"2006-01-02 15:04:05"`
	StockDto.ProductSearch
}

func (m *StockTransferGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type StockTransferGetPageResp struct {
	models.StockTransfer
	VendorName             string                       `json:"vendorName"`
	TypeName               string                       `json:"typeName"`
	StatusName             string                       `json:"statusName"`
	ToWarehouseName        string                       `json:"toWarehouseName"`
	ToLogicWarehouseName   string                       `json:"toLogicWarehouseName"`
	FromWarehouseName      string                       `json:"fromWarehouseName"`
	FromLogicWarehouseName string                       `json:"fromLogicWarehouseName"`
	VerifyStatusName       string                       `json:"verifyStatusName"`
	OptionRulesTransfer    StockDto.OptionRulesTransfer `json:"optionRules"`
}

func (m *StockTransferGetPageResp) Set() interface{} {
	return *m
}

func (m *StockTransferGetPageResp) SetTransferRulesByStatus() {
	// 查看按钮 | 全部状态展示
	m.OptionRulesTransfer.Detail = true
	// 更新按钮 | 待提交状态下
	if m.Status == models.TransferStatus99 {
		m.OptionRulesTransfer.Update.True()
	}
	// 删除按钮 | 待提交状态下
	if m.Status == models.TransferStatus99 {
		m.OptionRulesTransfer.Detail.True()
	}
	// 审核按钮 | 已提交状态下
	if m.Status == models.TransferStatus98 {
		m.OptionRulesTransfer.Audit.True()
	}

	// m.OptionRulesTransfer.Update = true
	// m.OptionRulesTransfer.Delete = true
	// m.OptionRulesTransfer.Detail = true
	// m.OptionRulesTransfer.Audit = false
	// switch m.Status {
	// case models.TransferStatus0:
	// 	m.OptionRulesTransfer.Update.False()
	// 	m.OptionRulesTransfer.Delete.False()
	// 	m.OptionRulesTransfer.Detail.True()
	// 	m.OptionRulesTransfer.Audit.False()
	// case models.TransferStatus1:
	// 	m.OptionRulesTransfer.Update.False()
	// 	m.OptionRulesTransfer.Delete.False()
	// 	m.OptionRulesTransfer.Detail.True()
	// 	m.OptionRulesTransfer.Audit.False()
	// case models.TransferStatus2:
	// 	m.OptionRulesTransfer.Update.False()
	// 	m.OptionRulesTransfer.Delete.False()
	// 	m.OptionRulesTransfer.Detail.True()
	// 	m.OptionRulesTransfer.Audit.False()
	// case models.TransferStatus3:
	// 	m.OptionRulesTransfer.Update.False()
	// 	m.OptionRulesTransfer.Delete.False()
	// 	m.OptionRulesTransfer.Detail.True()
	// 	m.OptionRulesTransfer.Audit.False()
	// case models.TransferStatus98:
	// 	m.OptionRulesTransfer.Update.False()
	// 	m.OptionRulesTransfer.Delete.False()
	// 	m.OptionRulesTransfer.Detail.True()
	// 	m.OptionRulesTransfer.Audit.True()
	// case models.TransferStatus99:
	// 	m.OptionRulesTransfer.Update.True()
	// 	m.OptionRulesTransfer.Delete.True()
	// 	m.OptionRulesTransfer.Detail.True()
	// 	m.OptionRulesTransfer.Audit.False()
	// }
}

type StockTransferInsertReq struct {
	Id                     int    `json:"-" comment:"id"` // id
	Type                   string `json:"type" comment:"调拨类型: 0 -正常调拨 1-次品调拨 2-正转次调拨 3-次转正调拨" vd:"$=='0' || $=='1' || $=='2' || $=='3'; msg:'type为0、1、2、3'"`
	FromWarehouseCode      string `json:"fromWarehouseCode" comment:"出库实体仓code" vd:"@:len($)>0; msg:'出库实体仓库不能为空'"`
	FromLogicWarehouseCode string `json:"fromLogicWarehouseCode" comment:"出库逻辑仓code" vd:"@:len($)>0; msg:'出库逻辑仓库不能为空'"`
	ToWarehouseCode        string `json:"toWarehouseCode" comment:"入库实体仓code" vd:"@:len($)>0; msg:'入库实体仓库不能为空'"`
	ToLogicWarehouseCode   string `json:"toLogicWarehouseCode" comment:"入库逻辑仓code" vd:"@:len($)>0; msg:'入库逻辑仓库不能为空'"`
	Remark                 string `json:"remark" comment:"备注"`
	LogisticsRemark        string `json:"logisticsRemark" comment:"物流备注"`
	Mobile                 string `json:"mobile" comment:"" vd:"regexp('^1[0-9]{10}$'); msg:'收货人手机格式不正确'"`
	Linkman                string `json:"linkman" comment:"联系人" vd:"@:len($)>0; msg:'联系人不能为空'"`
	Address                string `json:"address" comment:"详细地址"  vd:"@:len($)>0; msg:'详细地址不能为空'"`
	District               int    `json:"district" comment:"区"  vd:"$>0; msg:'收货人区不能为空'"`
	City                   int    `json:"city" comment:"市"  vd:"$>0; msg:'收货人市不能为空'"`
	Province               int    `json:"province" comment:"省"  vd:"$>0; msg:'收货人省不能为空'"`
	VendorId               int    `json:"vendorId" comment:"货主id" vd:"$>0; msg:'商品所属货主不能为空'"`
	common.ControlBy
	StockTransferProducts []StockTransferProductsReq `json:"stockTransferProducts"`
}

func (s *StockTransferInsertReq) Generate(tx *gorm.DB, model *models.StockTransfer) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.Type = s.Type
	model.FromWarehouseCode = s.FromWarehouseCode
	model.FromLogicWarehouseCode = s.FromLogicWarehouseCode
	model.ToWarehouseCode = s.ToWarehouseCode
	model.ToLogicWarehouseCode = s.ToLogicWarehouseCode
	model.Remark = s.Remark
	model.LogisticsRemark = s.LogisticsRemark
	model.Mobile = s.Mobile
	model.Linkman = s.Linkman
	model.Address = s.Address
	model.District = s.District
	model.City = s.City
	model.Province = s.Province
	model.VendorId = s.VendorId
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.CreateByName = s.CreateByName
	for _, item := range s.StockTransferProducts {
		item.SetCreateBy(s.CreateBy)
		item.SetCreateByName(s.CreateByName)
		modelStockTransferProducts := models.StockTransferProducts{}
		item.Generate(&modelStockTransferProducts)
		model.StockTransferProducts = append(model.StockTransferProducts, modelStockTransferProducts)
	}
	regionsMap := models.GeRegionMapByIds(tx, []int{s.Province, s.City, s.District})
	model.GenerateRegionName(s.Province, s.City, s.District, regionsMap)
}

func (s *StockTransferInsertReq) GetId() interface{} {
	return s.Id
}

type StockTransferUpdateReq struct {
	Id                     int    `uri:"id" comment:"id"` // id
	Type                   string `json:"type" comment:"调拨类型: 0 -正常调拨 1-次品调拨 2-正转次调拨 3-次转正调拨" vd:"$=='0' || $=='1' || $=='2' || $=='3'; msg:'type为0、1、2、3'"`
	FromWarehouseCode      string `json:"fromWarehouseCode" comment:"出库实体仓code" vd:"@:len($)>0; msg:'出库实体仓库不能为空'"`
	FromLogicWarehouseCode string `json:"fromLogicWarehouseCode" comment:"出库逻辑仓code" vd:"@:len($)>0; msg:'出库逻辑仓库不能为空'"`
	ToWarehouseCode        string `json:"toWarehouseCode" comment:"入库实体仓code" vd:"@:len($)>0; msg:'入库实体仓库不能为空'"`
	ToLogicWarehouseCode   string `json:"toLogicWarehouseCode" comment:"入库逻辑仓code" vd:"@:len($)>0; msg:'入库逻辑仓库不能为空'"`
	Remark                 string `json:"remark" comment:"备注"`
	LogisticsRemark        string `json:"logisticsRemark" comment:"物流备注"`
	Mobile                 string `json:"mobile" comment:"" vd:"regexp('^1[0-9]{10}$'); msg:'收货人手机格式不正确'"`
	Linkman                string `json:"linkman" comment:"联系人" vd:"@:len($)>0; msg:'联系人不能为空'"`
	Address                string `json:"address" comment:"详细地址"  vd:"@:len($)>0; msg:'详细地址不能为空'"`
	District               int    `json:"district" comment:"区"  vd:"$>0; msg:'收货人区不能为空'"`
	City                   int    `json:"city" comment:"市"  vd:"$>0; msg:'收货人市不能为空'"`
	Province               int    `json:"province" comment:"省"  vd:"$>0; msg:'收货人省不能为空'"`
	VendorId               int    `json:"vendorId" comment:"货主id" vd:"$>0; msg:'商品所属货主不能为空'"`
	common.ControlBy
	StockTransferProducts []StockTransferProductsReq `json:"stockTransferProducts"`
}

func (s *StockTransferUpdateReq) Generate(tx *gorm.DB, model *models.StockTransfer) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.Type = s.Type
	model.FromWarehouseCode = s.FromWarehouseCode
	model.FromLogicWarehouseCode = s.FromLogicWarehouseCode
	model.ToWarehouseCode = s.ToWarehouseCode
	model.ToLogicWarehouseCode = s.ToLogicWarehouseCode
	model.Remark = s.Remark
	model.LogisticsRemark = s.LogisticsRemark
	model.Mobile = s.Mobile
	model.Linkman = s.Linkman
	model.Address = s.Address
	model.District = s.District
	model.City = s.City
	model.Province = s.Province
	model.VendorId = s.VendorId
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
	model.UpdateByName = s.UpdateByName
	for _, item := range s.StockTransferProducts {
		item.SetUpdateBy(s.UpdateBy)
		item.SetUpdateByName(s.UpdateByName)
		modelStockTransferProducts := models.StockTransferProducts{}
		item.GenerateForUpdate(&modelStockTransferProducts)
		model.StockTransferProducts = append(model.StockTransferProducts, modelStockTransferProducts)
	}
	regionsMap := models.GeRegionMapByIds(tx, []int{s.Province, s.City, s.District})
	model.GenerateRegionName(s.Province, s.City, s.District, regionsMap)
}

func (s *StockTransferUpdateReq) GetId() interface{} {
	return s.Id
}

// StockTransferGetReq 功能获取请求参数
type StockTransferGetReq struct {
	Id int `uri:"id"`
}

func (s *StockTransferGetReq) GetId() interface{} {
	return s.Id
}

type StockTransferGetResp struct {
	models.StockTransfer
	VendorName             string `json:"vendorName"`
	TypeName               string `json:"typeName"`
	StatusName             string `json:"statusName"`
	ToWarehouseName        string `json:"toWarehouseName"`
	ToLogicWarehouseName   string `json:"toLogicWarehouseName"`
	FromWarehouseName      string `json:"fromWarehouseName"`
	FromLogicWarehouseName string `json:"fromLogicWarehouseName"`
	VerifyStatusName       string `json:"verifyStatusName"`
	AutoInName             string `json:"autoInName"`
	AutoOutName            string `json:"autoOutName"`

	StockTransferProducts []StockTransferProductsGetResp `json:"stockTransferProducts" copier:"-"`
}

func (s *StockTransferGetResp) InitData(tx *gorm.DB, model *models.StockTransfer) error {
	// 获取goods product 信息通过pc
	goodsIdSlice := lo.Uniq(lo.Map(model.StockTransferProducts, func(item models.StockTransferProducts, _ int) int {
		return item.GoodsId
	}))

	goodsProductMap := models.GetGoodsProductInfoFromPc(tx, goodsIdSlice)

	// 获取出库单商品出库数量
	stockOutboundProducts := &models.StockOutboundProducts{}
	outboundProductsActQuantityMap := stockOutboundProducts.GetActQuantityMapBySourceCode(tx, s.TransferCode)

	// 获取入库单商品入库数量
	stockEntryProducts := &models.StockEntryProducts{}
	entryProductsActQuantityMap := stockEntryProducts.GetActQuantityMapBySourceCode(tx, s.TransferCode)

	for _, item := range model.StockTransferProducts {
		productGetResp := StockTransferProductsGetResp{}
		if err := utils.CopyDeep(&productGetResp, &item); err != nil {
			return err
		}
		productGetResp.Number += 1
		productGetResp.FromActQuantity = outboundProductsActQuantityMap[productGetResp.SkuCode]
		productGetResp.ToActQuantity = entryProductsActQuantityMap[productGetResp.SkuCode]
		productGetResp.DiffActQuantity = utils.AbsToInt(productGetResp.FromActQuantity - productGetResp.ToActQuantity)
		productGetResp.VendorName = s.Vendor.NameZh
		productGetResp.FillProductGoodsRespData(productGetResp.GoodsId, goodsProductMap)
		s.StockTransferProducts = append(s.StockTransferProducts, productGetResp)
	}
	s.VendorName = s.Vendor.NameZh
	s.TypeName = utils.GetTFromMap(s.Type, models.TransferTypeMap)
	s.StatusName = utils.GetTFromMap(s.Status, models.TransferStatusMap)
	s.VerifyStatusName = utils.GetTFromMap(s.VerifyStatus, models.TransferVerifyStatusMap)
	s.AutoInName = models.TransferAutoInMap[s.AutoIn]
	s.AutoOutName = models.TransferAutoInMap[s.AutoOut]

	s.ToWarehouseName = s.ToWarehouse.WarehouseName
	s.FromWarehouseName = s.FromWarehouse.WarehouseName
	s.ToLogicWarehouseName = s.ToLogicWarehouse.LogicWarehouseName
	s.FromLogicWarehouseName = s.FromLogicWarehouse.LogicWarehouseName

	return nil
}

type StockTransferProductsGetResp struct {
	models.StockTransferProducts
	StockDto.ProductGoodsResp
	FromActQuantity int `json:"fromActQuantity"`
	ToActQuantity   int `json:"toActQuantity"`
	DiffActQuantity int `json:"diffActQuantity"`
	Number          int `json:"number"`
}

// StockTransferDeleteReq 功能删除请求参数
type StockTransferDeleteReq struct {
	Id int `json:"id"`
	common.ControlBy
}

func (s *StockTransferDeleteReq) GetId() interface{} {
	return s.Id
}

type StockTransferAuditReq struct {
	Id           int    `json:"id"`
	VerifyRemark string `json:"verifyRemark"`
	VerifyStatus string `json:"verifyStatus" vd:"$=='1' || $=='2'; msg:'verifyStatus为1、2'"`
	common.ControlBy
}

func (s *StockTransferAuditReq) Generate(model *models.StockTransfer) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.VerifyStatus = s.VerifyStatus
	model.VerifyRemark = s.VerifyRemark
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
	model.UpdateByName = s.UpdateByName
	model.VerifyUid = s.UpdateBy
	model.VerifyTime = time.Now()
}

func (s *StockTransferAuditReq) GetId() interface{} {
	return s.Id
}

// 调拨单验证skus

type StockTransferValidateSkusReq struct {
	SkuCodes          string `form:"skuCodes" vd:"@:len($)>0; msg:'skuCodes不能为空'"`
	VendorId          int    `form:"vendorId" vd:"@:$>0; msg:'vendorId不能为空'"`
	FromWarehouseCode string `form:"fromWarehouseCode" vd:"@:len($)>0; msg:'fromWarehouseCode不能为空'"`
	ToWarehouseCode   string `form:"toWarehouseCode" vd:"@:len($)>0; msg:'toWarehouseCode不能为空'"`
}

type StockTransferValidateSkusResp struct {
	StockDto.ProductGoodsResp
	VendorName      string `json:"vendorName"`
	SkuCode         string `json:"skuCode"`
	StockLocationId int    `json:"stockLocationId"`
}
