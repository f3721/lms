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

type StockOutboundGetPageReq struct {
	dto.Pagination     `search:"-"`
	Ids                []int     `form:"ids[]"  search:"-"`
	OutboundCode       string    `form:"outboundCode"  search:"type:exact;column:outbound_code;table:stock_outbound" comment:"出库单编码"`
	Type               string    `form:"type"  search:"type:exact;column:type;table:stock_outbound" comment:"出库类型:  0 大货出库  1 订单出库  2 其他"`
	SourceType         string    `form:"sourceType"  search:"-"`
	Status             string    `form:"status"  search:"type:exact;column:status;table:stock_outbound" comment:"状态:0-已作废 1-创建 2-已完成"`
	SourceCode         string    `form:"sourceCode"  search:"type:exact;column:source_code;table:stock_outbound" comment:"来源单据code"`
	OutboundTime       time.Time `form:"outboundTime"  search:"type:exact;column:outbound_time;table:stock_outbound" comment:"出库时间"`
	WarehouseCode      string    `form:"warehouseCode"  search:"type:exact;column:warehouse_code;table:stock_outbound" comment:"实体仓code"`
	LogicWarehouseCode string    `form:"logicWarehouseCode"  search:"type:exact;column:logic_warehouse_code;table:stock_outbound" comment:"逻辑仓code"`
	CreatedAtStart     time.Time `form:"createdAtStart"  search:"-" comment:"" time_format:"2006-01-02 15:04:05"`
	CreatedAtEnd       time.Time `form:"createdAtEnd"  search:"-" comment:"" time_format:"2006-01-02 15:04:05"`
	VendorId           int       `form:"vendorId"  search:"type:exact;column:vendor_id;table:stock_outbound" comment:"货主id"`
	Recipient          string    `form:"recipient" search:"-"`
	StockDto.ProductSearch
}

func (m *StockOutboundGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type StockOutboundGetPageResp struct {
	models.StockOutbound
	TypeName            string                       `json:"typeName"`
	SourceTypeName      string                       `json:"sourceTypeName"`
	StatusName          string                       `json:"statusName"`
	WarehouseName       string                       `json:"warehouseName"`
	LogicWarehouseName  string                       `json:"logicWarehouseName"`
	Recipient           string                       `json:"recipient"` // 领用人
	OptionRulesOutbound StockDto.OptionRulesOutbound `json:"optionRules" gorm:"-"`
}

type StockOutboundExportResp struct {
	models.StockOutbound
	TypeName                     string                         `json:"typeName"`
	SourceTypeName               string                         `json:"sourceTypeName"`
	StatusName                   string                         `json:"statusName"`
	WarehouseName                string                         `json:"warehouseName"`
	LogicWarehouseName           string                         `json:"logicWarehouseName"`
	Recipient                    string                         `json:"recipient"` // 领用人
	OptionRulesOutbound          StockDto.OptionRulesOutbound   `json:"optionRules" gorm:"-"`
	StockOutboundProductsGetResp []StockOutboundProductsGetResp `json:"stockOutboundProducts" copier:"-"  gorm:"-"`
}

type StockOutboundExport struct {
	TypeName            string `json:"typeName"`
	SourceTypeName      string `json:"sourceTypeName"`
	StatusName          string `json:"statusName"`
	WarehouseName       string `json:"warehouseName"`
	LogicWarehouseName  string `json:"logicWarehouseName"`
	Recipient           string `json:"recipient"`
	OutboundCode        string `json:"outboundCode"`
	SourceCode          string `json:"sourceCode"`
	Remark              string `json:"remark"`
	OutboundTime        string `json:"outboundTime"`
	SkuCode             string `json:"skuCode"`
	ProductName         string `json:"productName"`
	MfgModel            string `json:"mfgModel"`
	BrandName           string `json:"brandName"`
	SalesUom            string `json:"salesUom"`
	VendorName          string `json:"vendorName"`
	ProductNo           string `json:"productNo"`
	VendorSkuCode       string `json:"vendorSkuCode"`
	LocationQuantity    int    `json:"locationQuantity"`
	LocationActQuantity int    `json:"locationActQuantity"`
	LocationCode        string `json:"locationCode"`
	DiffNum             int    `json:"diffNum"`
}

func (m *StockOutboundGetPageResp) SetOutboundRulesByStatus() {
	// casbin 设置 todo
	m.OptionRulesOutbound.Confirm = true
	m.OptionRulesOutbound.Detail = true
	m.OptionRulesOutbound.Print = true
	switch m.Status {
	case models.OutboundStatus0:
		m.OptionRulesOutbound.Confirm.False()
		m.OptionRulesOutbound.Detail.True()
		m.OptionRulesOutbound.Print.False()
	case models.OutboundStatus1:
		m.OptionRulesOutbound.Confirm.True()
		m.OptionRulesOutbound.Detail.True()
		m.OptionRulesOutbound.Print.False()
	case models.OutboundStatus2:
		m.OptionRulesOutbound.Confirm.False()
		m.OptionRulesOutbound.Detail.True()
		if m.IsTransferOutbound() {
			m.OptionRulesOutbound.Print.False()
		}
	}
}

type StockOutboundInsertReq struct {
	SourceCode            string                     `json:"sourceCode" comment:"来源单据code"`
	Remark                string                     `json:"remark" comment:"备注"`
	WarehouseCode         string                     `json:"warehouseCode" comment:"实体仓code"`
	LogicWarehouseCode    string                     `json:"logicWarehouseCode" comment:"逻辑仓code"`
	VendorId              int                        `json:"vendorId" comment:"货主id"`
	StockOutboundProducts []StockOutboundProductsReq `json:"stockOutboundProducts"`
	common.ControlBy
}

func (s *StockOutboundInsertReq) Generate(model *models.StockOutbound) {
	model.Status = models.OutboundStatus1
	model.SourceCode = s.SourceCode
	model.Remark = s.Remark
	model.WarehouseCode = s.WarehouseCode
	model.LogicWarehouseCode = s.LogicWarehouseCode
	model.VendorId = s.VendorId
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.CreateByName = s.CreateByName
	for _, item := range s.StockOutboundProducts {
		item.SetCreateBy(s.CreateBy)
		item.SetCreateByName(s.CreateByName)
		modelStockOutboundProducts := models.StockOutboundProducts{}
		item.Generate(&modelStockOutboundProducts)
		model.StockOutboundProducts = append(model.StockOutboundProducts, modelStockOutboundProducts)
	}
}

type StockOutboundInsertForOrderReq struct {
	SourceCode            string                             `json:"sourceCode" comment:"来源单据code"`
	Remark                string                             `json:"remark" comment:"备注"`
	WarehouseCode         string                             `json:"warehouseCode" comment:"实体仓code"`
	StockOutboundProducts []StockOutboundProductsForOrderReq `json:"stockOutboundProducts"`
	common.ControlBy
}

func (s *StockOutboundInsertForOrderReq) Generate(tx *gorm.DB, model *models.StockOutbound) error {
	model.Status = models.OutboundStatus1
	model.SourceCode = s.SourceCode
	model.Remark = s.Remark
	model.WarehouseCode = s.WarehouseCode
	// 逻辑仓获取
	lwh := &models.LogicWarehouse{}
	if err := lwh.GetPassLogicWarehouseByWhCode(tx, s.WarehouseCode); err != nil {
		return err
	}
	model.LogicWarehouseCode = lwh.LogicWarehouseCode

	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.CreateByName = s.CreateByName
	for _, item := range s.StockOutboundProducts {
		item.SetCreateBy(s.CreateBy)
		item.SetCreateByName(s.CreateByName)
		modelStockOutboundProducts := models.StockOutboundProducts{}
		item.Generate(&modelStockOutboundProducts)
		model.StockOutboundProducts = append(model.StockOutboundProducts, modelStockOutboundProducts)
	}
	return nil
}

// // 逐步作废
// type StockOutboundCancelForCsOrderReq struct {
// 	CsNo            string                                         `json:"csNo" comment:"售后单号"`
// 	OrderId         string                                         `json:"orderId" comment:"订单号"`
// 	CsOrderProducts []StockOutboundProductsPartCancelForCsOrderReq `json:"csOrderProducts"`
// 	Remark          string                                         `json:"remark" comment:"备注"`
// 	common.ControlBy
// }

// // 逐步作废
// func (s *StockOutboundCancelForCsOrderReq) Generate(model *models.StockOutbound) {
// 	if model.Status == "1" { // 未出库-作废
// 		model.Status = models.OutboundStatus0
// 	}
// 	if model.Status == "3" { // 部分出库-已发货
// 		model.Status = models.OutboundStatus2
// 	}
// 	model.Remark = s.Remark
// 	if s.UpdateBy != 0 {
// 		model.UpdateBy = s.UpdateBy
// 	}
// 	if s.UpdateByName != "" {
// 		model.UpdateByName = s.UpdateByName
// 	}
// }

type StockOutboundPartCancelForCsOrderReq struct {
	CsNo            string                                         `json:"csNo" comment:"售后单号"`
	OrderId         string                                         `json:"orderId" comment:"订单号"`
	CsOrderProducts []StockOutboundProductsPartCancelForCsOrderReq `json:"csOrderProducts"`
	Remark          string                                         `json:"remark" comment:"备注"`
	common.ControlBy
}

// 处理出库单状态
func (s *StockOutboundPartCancelForCsOrderReq) GenerateOutBound(model *models.StockOutbound) {
	currTime := time.Now()
	if model.Status == "1" { // 未出库-作废
		model.Status = models.OutboundStatus0
	}
	if model.Status == "3" { // 部分出库-已发货
		model.Status = models.OutboundStatus2
		model.OutboundEndTime = currTime
	}
	model.Remark = s.Remark
	if s.UpdateBy != 0 {
		model.UpdateBy = s.UpdateBy
	}
	if s.UpdateByName != "" {
		model.UpdateByName = s.UpdateByName
	}
}

// 处理商品状态
func (s *StockOutboundPartCancelForCsOrderReq) GenerateProduct(model *models.StockOutboundProducts) {
	if s.UpdateBy != 0 {
		model.UpdateBy = s.UpdateBy
	}
	if s.UpdateByName != "" {
		model.UpdateByName = s.UpdateByName
	}
}

type StockOutboundUpdateReq struct {
	Id                 int       `uri:"id" comment:"id"` // id
	OutboundCode       string    `json:"outboundCode" comment:"出库单编码"`
	Type               string    `json:"type" comment:"出库类型:  0 大货出库  1 订单出库  2 其他"`
	Status             string    `json:"status" comment:"状态:0-已作废 1-创建 2-已完成"`
	SourceCode         string    `json:"sourceCode" comment:"来源单据code"`
	Remark             string    `json:"remark" comment:"备注"`
	OutboundTime       time.Time `json:"outboundTime" comment:"出库时间"`
	WarehouseCode      string    `json:"warehouseCode" comment:"实体仓code"`
	LogicWarehouseCode string    `json:"logicWarehouseCode" comment:"逻辑仓code"`
	CreateByName       string    `json:"createByName" comment:""`
	UpdateByName       string    `json:"updateByName" comment:""`
	common.ControlBy
}

func (s *StockOutboundUpdateReq) Generate(model *models.StockOutbound) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.OutboundCode = s.OutboundCode
	model.Type = s.Type
	model.Status = s.Status
	model.SourceCode = s.SourceCode
	model.Remark = s.Remark
	model.OutboundTime = s.OutboundTime
	model.WarehouseCode = s.WarehouseCode
	model.LogicWarehouseCode = s.LogicWarehouseCode
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
	model.CreateByName = s.CreateByName
	model.UpdateByName = s.UpdateByName
}

func (s *StockOutboundUpdateReq) GetId() interface{} {
	return s.Id
}

// StockOutboundGetReq 功能获取请求参数
type StockOutboundGetReq struct {
	Id   int    `uri:"id"`
	Type string `form:"type"`
}

func (s *StockOutboundGetReq) GetId() interface{} {
	return s.Id
}

type StockOutboundCommonResp struct {
	models.StockOutbound
	TypeName           string `json:"typeName"`
	SourceTypeName     string `json:"sourceTypeName"`
	StatusName         string `json:"statusName"`
	WarehouseName      string `json:"warehouseName"`
	LogicWarehouseName string `json:"logicWarehouseName"`
	AddressFullName    string `json:"addressFullName"`
	Address            string `json:"address"`
	District           int    `json:"district"`
	City               int    `json:"city"`
	Province           int    `json:"province"`
	DistrictName       string `json:"districtName"`
	CityName           string `json:"cityName"`
	ProvinceName       string `json:"provinceName"`
	ExternalOrderNo    string `json:"externalOrderNo"`
	ReceiveAddrInfo    struct {
		DistrictName string `json:"districtName"`
		CityName     string `json:"cityName"`
		ProvinceName string `json:"provinceName"`

		AddressFullName string `json:"addressFullName"`
		Address         string `json:"address"`
		District        int    `json:"district"`
		City            int    `json:"city"`
		Province        int    `json:"province"`
		Mobile          string `json:"mobile"`
		Linkman         string `json:"linkman"`
		// 领用单特有， 下单人信息
		UserName        string `json:"userName"`
		UserCompanyName string `json:"userCompanyName"`
		UserPhone       string `json:"userPhone"`
		UserDepartment  string `json:"userDepartment"`
	} `json:"receiveAddrInfo"`
}

type StockOutboundNoLocationResp struct {
	*StockOutboundCommonResp
	StockOutboundProducts []StockOutboundProductsNoLocationResp `json:"stockOutboundProducts" copier:"-"`
}

func (s *StockOutboundNoLocationResp) InitProductData(tx *gorm.DB) error {
	stockOutboundProducts := models.StockOutboundProducts{}
	products, err := stockOutboundProducts.GetList(tx, s.OutboundCode)
	if err != nil {
		return err
	}
	// 获取goods product 信息通过pc
	goodsIdSlice := lo.Uniq(lo.Map(products, func(item models.StockOutboundProducts, _ int) int {
		return item.GoodsId
	}))
	goodsProductMap := models.GetGoodsProductInfoFromPc(tx, goodsIdSlice)
	// 获取vendor map
	vendorIdSlice := lo.Uniq(lo.Map(products, func(item models.StockOutboundProducts, _ int) int {
		return item.VendorId
	}))
	vendorIdMap := models.GetVendorsMapByIds(tx, vendorIdSlice)

	for index, item := range products {
		productGetResp := StockOutboundProductsNoLocationResp{}
		if err := utils.CopyDeep(&productGetResp, &item); err != nil {
			return err
		}
		productGetResp.Number = index + 1
		productGetResp.VendorName = vendorIdMap[productGetResp.VendorId]
		productGetResp.FillProductGoodsRespData(productGetResp.GoodsId, goodsProductMap)
		s.StockOutboundProducts = append(s.StockOutboundProducts, productGetResp)
	}
	return nil
}

type StockOutboundProductsNoLocationResp struct {
	models.StockOutboundProducts
	Number int `json:"number"`
	StockDto.ProductGoodsResp
}

type StockOutboundGetResp struct {
	*StockOutboundCommonResp
	StockOutboundProducts []*StockOutboundProductsGetResp `json:"stockOutboundProducts" copier:"-"`
}

type StockOutboundProductsGetResp struct {
	models.OutboundProductSubCustom
	LocationCode string `json:"locationCode"`
	DiffNum      int    `json:"diffNum"`
	Number       int    `json:"number"`
	StockDto.ProductGoodsResp
}

func (s *StockOutboundGetResp) InitProductData(tx *gorm.DB, Type string) error {
	OutboundProductSubCustom := models.OutboundProductSubCustom{}
	subProducts, err := OutboundProductSubCustom.GetList(tx, []string{s.OutboundCode})
	if err != nil {
		return err
	}

	// 获取goods product 信息通过pc
	goodsIdSlice := lo.Uniq(lo.Map(subProducts, func(item models.OutboundProductSubCustom, _ int) int {
		return item.GoodsId
	}))
	goodsProductMap := models.GetGoodsProductInfoFromPc(tx, goodsIdSlice)

	// 获取vendor map
	vendorIdSlice := lo.Uniq(lo.Map(subProducts, func(item models.OutboundProductSubCustom, _ int) int {
		return item.VendorId
	}))
	vendorIdMap := models.GetVendorsMapByIds(tx, vendorIdSlice)

	// 查询库位部分出库历史记录
	subIds := lo.Map(subProducts, func(item models.OutboundProductSubCustom, _ int) int {
		return item.SubId
	})
	subLog := &models.StockOutboundProductsSubLog{}
	subList, err := subLog.ListbySubIds(tx, subIds)
	if err != nil {
		return err
	}
	subListMap := map[int][]*models.StockOutboundProductsSubLog{}
	for _, item := range subList {
		subListMap[item.SubId] = append(subListMap[item.SubId], item)
	}

	for index, item := range subProducts {
		productGetResp := StockOutboundProductsGetResp{}
		if err := utils.CopyDeep(&productGetResp, &item); err != nil {
			return err
		}
		// if Type == "confirm" && s.IsOrderOutbound() { // 旧逻辑-订单出库必须全部出库
		// 	productGetResp.ActQuantity = productGetResp.Quantity
		// 	productGetResp.LocationActQuantity = productGetResp.LocationQuantity
		// }
		productGetResp.LocationCode = productGetResp.StockLocation.LocationCode
		productGetResp.DiffNum = productGetResp.LocationQuantity - productGetResp.LocationActQuantity
		productGetResp.Number = index + 1
		productGetResp.VendorName = vendorIdMap[productGetResp.VendorId]
		productGetResp.FillProductGoodsRespData(productGetResp.GoodsId, goodsProductMap)

		// 追加库位出库日志
		productGetResp.SubLog = subListMap[item.SubId]
		// 兼容旧数据
		if s.Status == "2" && len(productGetResp.SubLog) == 0 {
			oldSubLog := []*models.StockOutboundProductsSubLog{
				{
					SubId:               item.SubId,
					OutboundTime:        s.OutboundTime,
					LocationActQuantity: item.LocationActQuantity,
				},
			}
			productGetResp.SubLog = oldSubLog
		}

		s.StockOutboundProducts = append(s.StockOutboundProducts, &productGetResp)
	}
	return nil
}

// StockOutboundDeleteReq 功能删除请求参数
type StockOutboundDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *StockOutboundDeleteReq) GetId() interface{} {
	return s.Ids
}

// 部分出库请求
type StockPartOutboundReq struct {
	Id                    int                                `json:"id" vd:"@:$>0;msg:'出库单ID必填'"`                               // 出库单ID
	OutboundType          int                                `json:"outboundType" vd:"@:in($,1,2);msg:'outboundType默认必填[1,2]'"` // 1-部分出库 2-全部出库
	StockOutboundProducts []*StockOutboundProductsConfirmReq `json:"stockOutboundProducts" vd:"@:len($)>0;msg:'出库单商品信息必填'"`     // 出库单商品信息
	common.ControlBy
}

// 确认出库请求
type StockOutboundConfirmReq struct {
	Id                    int                               `json:"id" vd:"@:$>0;msg:'出库单商品ID必填'"`
	StockOutboundProducts []StockOutboundProductsConfirmReq `json:"stockOutboundProducts" vd:"@:len($)>0;msg:'出库单商品信息必填'"`
	common.ControlBy
}

func (s *StockOutboundConfirmReq) GetId() interface{} {
	return s.Id
}

func (s *StockOutboundConfirmReq) Generate(model *models.StockOutbound) {
	model.Status = models.OutboundStatus2
	model.OutboundTime = time.Now()
	model.UpdateBy = s.UpdateBy
	model.UpdateByName = s.UpdateByName
	for index := range model.StockOutboundProducts {
		model.StockOutboundProducts[index].UpdateBy = s.UpdateBy
		model.StockOutboundProducts[index].UpdateByName = s.UpdateByName

	}
}

type StockOutboundPrintReq struct {
	Id int `uri:"id"`
}

func (s *StockOutboundPrintReq) GetId() interface{} {
	return s.Id
}

type StockOutboundPrintResp struct {
	Html string `json:"html"`
	Data any    `json:"data"`
}
