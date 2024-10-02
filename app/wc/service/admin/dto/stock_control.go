package dto

import (
	"github.com/samber/lo"
	StockDto "go-admin/common/dto/stock/dto"
	"go-admin/common/utils"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"

	"go-admin/app/wc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type StockControlGetPageReq struct {
	dto.Pagination     `search:"-"`
	ControlCode        string    `form:"controlCode"  search:"type:exact;column:control_code;table:stock_control" comment:"调整单编码"`
	Type               string    `form:"type"  search:"type:exact;column:type;table:stock_control" comment:"调整类型: 0 调增 1 调减"`
	Status             string    `form:"status"  search:"type:exact;column:status;table:stock_control" comment:"状态:0-已作废 1-创建 2-已完成 99未提交"`
	WarehouseCode      string    `form:"warehouseCode"  search:"type:exact;column:warehouse_code;table:stock_control" comment:"实体仓code"`
	LogicWarehouseCode string    `form:"logicWarehouseCode"  search:"type:exact;column:logic_warehouse_code;table:stock_control" comment:"逻辑仓code"`
	VerifyStatus       string    `form:"verifyStatus"  search:"type:exact;column:verify_status;table:stock_control" comment:"审核状态 0 待审核 1 审核通过 2 审核驳回 99初始化"`
	VendorId           int       `form:"vendorId"  search:"type:exact;column:vendor_id;table:stock_control" comment:"货主id"`
	CreatedAtStart     time.Time `form:"createdAtStart"  search:"-" comment:"" time_format:"2006-01-02 15:04:05"`
	CreatedAtEnd       time.Time `form:"createdAtEnd"  search:"-" comment:"" time_format:"2006-01-02 15:04:05"`
	StockDto.ProductSearch
}

func (m *StockControlGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type StockControlGetPageResp struct {
	models.StockControl
	VendorName         string                      `json:"vendorName"`
	TypeName           string                      `json:"typeName"`
	StatusName         string                      `json:"statusName"`
	WarehouseName      string                      `json:"warehouseName"`
	LogicWarehouseName string                      `json:"logicWarehouseName"`
	VerifyStatusName   string                      `json:"verifyStatusName"`
	ControlTotalNum    int                         `json:"controlTotalNum"`
	OptionRulesControl StockDto.OptionRulesControl `json:"optionRules" gorm:"-"`
}

func (m *StockControlGetPageResp) SetControlRulesByStatus() {
	// casbin 设置 todo
	m.OptionRulesControl.Update = true
	m.OptionRulesControl.Delete = true
	m.OptionRulesControl.Detail = true
	m.OptionRulesControl.Audit = true
	switch m.Status {
	case models.ControlStatus0:
		m.OptionRulesControl.Update.False()
		m.OptionRulesControl.Delete.False()
		m.OptionRulesControl.Detail.True()
		m.OptionRulesControl.Audit.False()
	case models.ControlStatus1:
		m.OptionRulesControl.Update.False()
		m.OptionRulesControl.Delete.False()
		m.OptionRulesControl.Detail.True()
		m.OptionRulesControl.Audit.True()
	case models.ControlStatus2:
		m.OptionRulesControl.Update.False()
		m.OptionRulesControl.Delete.False()
		m.OptionRulesControl.Detail.True()
		m.OptionRulesControl.Audit.False()
	case models.ControlStatus99:
		m.OptionRulesControl.Update.True()
		m.OptionRulesControl.Delete.True()
		m.OptionRulesControl.Detail.True()
		m.OptionRulesControl.Audit.False()
	}
}

func (s *StockControlGetPageResp) InitData() {
	s.SetControlRulesByStatus()
	s.VendorName = s.Vendor.NameZh
	s.WarehouseName = s.Warehouse.WarehouseName
	s.LogicWarehouseName = s.LogicWarehouse.LogicWarehouseName
	s.TypeName = utils.GetTFromMap(s.Type, models.ControlTypeMap)
	s.StatusName = utils.GetTFromMap(s.Status, models.ControlStatusMap)
	s.VerifyStatusName = utils.GetTFromMap(s.VerifyStatus, models.ControlVerifyStatusMap)
}

type StockControlGetResp struct {
	models.StockControl
	VendorName           string                        `json:"vendorName"`
	TypeName             string                        `json:"typeName"`
	StatusName           string                        `json:"statusName"`
	VerifyStatusName     string                        `json:"verifyStatusName"`
	ControlTotalNum      int                           `json:"controlTotalNum"`
	StockControlProducts []StockControlProductsGetResp `json:"stockControlProducts" copier:"-"`
}

func (s *StockControlGetResp) InitData(tx *gorm.DB, model *models.StockControl, Type string) error {
	var sum int
	// 获取goods product 信息通过pc
	goodsIdSlice := lo.Uniq(lo.Map(model.StockControlProducts, func(item models.StockControlProducts, _ int) int {
		return item.GoodsId
	}))

	goodsProductMap := models.GetGoodsProductInfoFromPc(tx, goodsIdSlice)
	for index, item := range model.StockControlProducts {
		productGetResp := StockControlProductsGetResp{}
		if err := utils.CopyDeep(&productGetResp, &item); err != nil {
			return err
		}
		// 获取库位列表
		if Type == "edit" {
			stockLocations, err := models.GetStockLocationsHasId(tx, item.LogicWarehouseCode, item.StockLocationId)
			if err == nil {
				stockControlProductsLocationGetResp := StockControlProductsLocationGetResp{}
				for _, item := range *stockLocations {
					stockControlProductsLocationGetResp.Regenerate(item)
					productGetResp.StockLocation = append(productGetResp.StockLocation, stockControlProductsLocationGetResp)
				}
			}
		}
		productGetResp.LocationCode = item.StockLocation.LocationCode
		productGetResp.Number = index + 1
		productGetResp.VendorName = model.Vendor.NameZh
		productGetResp.FillProductGoodsRespData(productGetResp.GoodsId, goodsProductMap)
		s.StockControlProducts = append(s.StockControlProducts, productGetResp)
	}
	s.VendorName = s.Vendor.NameZh
	s.TypeName = utils.GetTFromMap(s.Type, models.ControlTypeMap)
	s.StatusName = utils.GetTFromMap(s.Status, models.ControlStatusMap)
	s.VerifyStatusName = utils.GetTFromMap(s.VerifyStatus, models.ControlVerifyStatusMap)
	for _, item := range s.StockControlProducts {
		sum += item.Quantity
	}
	s.ControlTotalNum = sum
	return nil
}

type StockControlProductsGetResp struct {
	models.StockControlProducts
	StockDto.ProductGoodsResp
	StockLocation []StockControlProductsLocationGetResp `json:"stockLocations" copier:"-"`
	LocationCode  string                                `json:"locationCode"`
	Number        int                                   `json:"number"`
}

type StockControlProductsLocationGetResp struct {
	Id           int    `json:"id"`
	LocationCode string `json:"locationCode"`
	Stock        int    `json:"stock"`
	LockStock    int    `json:"lockStock"`
}

func (s *StockControlProductsLocationGetResp) Regenerate(location models.StockLocation) {
	s.Id = location.Id
	s.LocationCode = location.LocationCode
}

type StockControlInsertReq struct {
	Id                 int    `json:"-" comment:"id"` // id
	Type               string `json:"type" comment:"调整类型: 0 调增 1 调减" default:"0"`
	WarehouseCode      string `json:"warehouseCode" comment:"实体仓code" vd:"@:len($)>0; msg:'实体仓库不能为空'"`
	LogicWarehouseCode string `json:"logicWarehouseCode" comment:"逻辑仓code" vd:"@:len($)>0; msg:'逻辑仓库不能为空'"`
	Remark             string `json:"remark" comment:"备注"`
	VendorId           int    `json:"vendorId" comment:"货主id"  vd:"$>0; msg:'货主不能为空'"`
	common.ControlBy
	StockControlProducts []StockControlProductsReq `json:"stockControlProducts"`
}

func (s *StockControlInsertReq) Generate(model *models.StockControl) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.Type = s.Type
	model.WarehouseCode = s.WarehouseCode
	model.LogicWarehouseCode = s.LogicWarehouseCode
	model.Remark = s.Remark
	model.VendorId = s.VendorId
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.CreateByName = s.CreateByName
	for _, item := range s.StockControlProducts {
		item.SetCreateBy(s.CreateBy)
		item.SetCreateByName(s.CreateByName)
		modelStockControlProducts := models.StockControlProducts{}
		item.Generate(&modelStockControlProducts)
		model.StockControlProducts = append(model.StockControlProducts, modelStockControlProducts)
	}
}

func (s *StockControlInsertReq) GetId() interface{} {
	return s.Id
}

type StockControlUpdateReq struct {
	Id                 int    `uri:"id" comment:"id"` // id
	WarehouseCode      string `json:"warehouseCode" comment:"实体仓code" vd:"@:len($)>0; msg:'实体仓库不能为空'"`
	LogicWarehouseCode string `json:"logicWarehouseCode" comment:"逻辑仓code" vd:"@:len($)>0; msg:'逻辑仓库不能为空'"`
	Remark             string `json:"remark" comment:"备注"`
	common.ControlBy
	StockControlProducts []StockControlProductsReq `json:"stockControlProducts"`
}

func (s *StockControlUpdateReq) Generate(model *models.StockControl) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.WarehouseCode = s.WarehouseCode
	model.LogicWarehouseCode = s.LogicWarehouseCode
	model.Remark = s.Remark
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
	model.UpdateByName = s.UpdateByName
	for _, item := range s.StockControlProducts {
		item.SetUpdateBy(s.UpdateBy)
		item.SetUpdateByName(s.UpdateByName)
		modelStockControlProducts := models.StockControlProducts{}
		item.GenerateForUpdate(&modelStockControlProducts)
		model.StockControlProducts = append(model.StockControlProducts, modelStockControlProducts)
	}
}

func (s *StockControlUpdateReq) GetId() interface{} {
	return s.Id
}

// StockControlGetReq 功能获取请求参数
type StockControlGetReq struct {
	Id   int    `uri:"id"`
	Type string `form:"type"`
}

func (s *StockControlGetReq) GetId() interface{} {
	return s.Id
}

// StockControlDeleteReq 功能删除请求参数
type StockControlDeleteReq struct {
	Id int `json:"id"`
	common.ControlBy
}

func (s *StockControlDeleteReq) GetId() interface{} {
	return s.Id
}

type StockControlAuditReq struct {
	Id           int    `json:"id"`
	VerifyRemark string `json:"verifyRemark"`
	VerifyStatus string `json:"verifyStatus" vd:"$=='1' || $=='2'; msg:'verifyStatus为1、2'"`
	common.ControlBy
}

func (s *StockControlAuditReq) Generate(model *models.StockControl) {
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

func (s *StockControlAuditReq) GetId() interface{} {
	return s.Id
}

// StockControlValidateSkusReq 调拨单验证skus
type StockControlValidateSkusReq struct {
	SkuCodes           string `form:"skuCodes" vd:"@:len($)>0; msg:'SKU不能为空'"`
	CurrentQuantity    int    `form:"currentQuantity" vd:"@:$>=0; msg:'盘后数量不能小于0'"`
	StockLocationId    int    `form:"stockLocationId" vd:"?"`
	VendorId           int    `form:"vendorId" vd:"@:$>0; msg:'商品所属货主不能为空'"`
	WarehouseCode      string `form:"warehouseCode" vd:"@:len($)>0; msg:'实体仓不能为空'"`
	LogicWarehouseCode string `form:"logicWarehouseCode" vd:"@:len($)>0; msg:'逻辑仓不能为空'"`
}

type StockControlValidateSkusResp struct {
	StockDto.ProductGoodsResp
	SkuCode         string                                `json:"skuCode"`
	ErrorMsg        string                                `json:"errorMsg"`
	Stock           int                                   `json:"stock" gorm:"type:int;comment:可用库存"`
	LockStock       int                                   `json:"lockStock" gorm:"type:int;comment:占用库存"`
	TotalStock      int                                   `json:"totalStock" gorm:"type:int;comment:在库库存"`
	CurrentQuantity int                                   `json:"currentQuantity" gorm:"type:int;comment:盘后数量"`
	Type            string                                `json:"type" gorm:"type:tinyint;comment:盘点结果: 0 盘盈 1 盘亏 2 无差异"`
	Quantity        int                                   `json:"quantity" gorm:"type:int unsigned;comment:差异数量"`
	StockLocationId int                                   `json:"stockLocationId"`
	StockLocation   []StockControlProductsLocationGetResp `json:"stockLocations"`
}

type ImportStockControl struct {
	Key                int    `json:"key"`
	VendorName         string `json:"vendorName" gorm:"type:string;comment:商品所属货主" vd:"@:len($)>0; msg:'商品所属货主不能为空'"`
	WarehouseName      string `json:"warehouseName" gorm:"type:string;comment:实体仓" vd:"@:len($)>0; msg:'实体仓不能为空'"`
	LogicWarehouseName string `json:"logicWarehouseName" gorm:"type:string;comment:逻辑仓" vd:"@:len($)>0; msg:'逻辑仓不能为空'"`
	SkuCode            string `json:"skuCode" gorm:"type:string;comment:SKU" vd:"@:len($)>0; msg:'SKU不能为空'"`
	CurrentQuantity    int    `json:"currentQuantity" gorm:"type:int;comment:盘后数量" vd:"@:$>=0; msg:'盘后数量不能小于0'"`
	LocationCode       string `json:"locationCode" vd:"@:len($)>0; msg:'库位编号不能为空'"`
}

type ImportReq struct {
	Data []ImportStockControl `json:"data"`
}

func (e *ImportStockControl) Trim(data map[string]interface{}) {
	e.VendorName = strings.TrimSpace(data["vendorName"].(string))
	e.WarehouseName = strings.TrimSpace(data["warehouseName"].(string))
	e.LogicWarehouseName = strings.TrimSpace(data["logicWarehouseName"].(string))
	e.SkuCode = strings.ToUpper(strings.TrimSpace(data["skuCode"].(string)))
	currentQuantity, _ := strconv.Atoi(strings.TrimSpace(data["currentQuantity"].(string)))
	e.CurrentQuantity = currentQuantity
	e.LocationCode = strings.ToUpper(strings.TrimSpace(data["locationCode"].(string)))
}

type GoodsCurrentQuantity struct {
	Keys            []int  `json:"keys"`
	GoodsId         int    `json:"goodsId"`
	LocationId      int    `json:"locationId"`
	SkuCode         string `json:"skuCode"`
	CurrentQuantity int    `json:"currentQuantity"`
	Stock           int    `json:"stock"`
	LockStock       int    `json:"lockStock"`
	Quantity        int    `json:"quantity"`
}

type BaseStockControl struct {
	BaseVendorId           int                                   `json:"baseVendorId"`
	BaseVendorName         string                                `json:"baseVendorName"`
	BaseWarehouseCode      string                                `json:"baseWarehouseCode"`
	BaseWarehouseName      string                                `json:"baseWarehouseName"`
	BaseLogicWarehouseCode string                                `json:"baseLogicWarehouseCode"`
	BaseLogicWarehouseName string                                `json:"baseLogicWarehouseName"`
	StockLocation          []StockControlProductsLocationGetResp `json:"stockLocations"`
}
