package dto

import (
	"go-admin/app/pc/models"
	"go-admin/common/actions"
	"go-admin/common/dto"
	common "go-admin/common/models"
	"go-admin/common/utils"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type GoodsGetPageReq struct {
	dto.Pagination    `search:"-"`
	Ids               []int    `form:"ids[]"  search:"-"`
	NameZh            string   `form:"nameZh"  search:"-"`
	VendorShortName   string   `json:"vendorShortName" search:"-"`
	SkuCode           string   `form:"skuCode"  search:"type:exact;column:sku_code;table:goods"`
	SupplierSkuCode   string   `form:"supplierSkuCode"  search:"type:exact;column:supplier_sku_code;table:goods"`
	OnlineStatus      int      `form:"onlineStatus"  search:"-"`
	ApproveStatus     int      `form:"approveStatus"  search:"-"`
	Status            int      `form:"status"  search:"-"`
	ProductNo         string   `form:"productNo"  search:"type:exact;column:product_no;table:goods"`
	VendorId          int      `form:"vendorId"  search:"type:exact;column:vendor_id;table:goods"`
	WarehouseCode     string   `form:"warehouseCode"  search:"type:exact;column:warehouse_code;table:goods"`
	WarehouseCodeList []string `form:"warehouseCodeList[]"  search:"-"`
	GoodsOrder
}

type GoodsGetPageResp struct {
	models.Goods
	VendorName      string `json:"vendorName" gorm:"-"`
	VendorShortName string `json:"vendorShortName" gorm:"-"`
	WarehouseName   string `json:"warehouseName" gorm:"-"`
	CompanyName     string `json:"companyName" gorm:"-"`
}

type GoodsOrder struct {
	Id                string `form:"idOrder"  search:"type:order;column:id;table:goods"`
	SkuCode           string `form:"skuCodeOrder"  search:"type:order;column:sku_code;table:goods"`
	SupplierSkuCode   string `form:"supplierSkuCodeOrder"  search:"type:order;column:supplier_sku_code;table:goods"`
	WarehouseCode     string `form:"warehouseCodeOrder"  search:"type:order;column:warehouse_code;table:goods"`
	ProductNo         string `form:"productNoOrder"  search:"type:order;column:product_no;table:goods"`
	MarketPrice       string `form:"marketPriceOrder"  search:"type:order;column:market_price;table:goods"`
	PriceModifyReason string `form:"priceModifyReasonOrder"  search:"type:order;column:price_modify_reason;table:goods"`
	ApproveStatus     string `form:"approveStatusOrder"  search:"type:order;column:approve_status;table:goods"`
	ApproveRemark     string `form:"approveRemarkOrder"  search:"type:order;column:approve_remark;table:goods"`
	Status            string `form:"statusOrder"  search:"type:order;column:status;table:goods"`
	OnlineStatus      string `form:"onlineStatusOrder"  search:"type:order;column:online_status;table:goods"`
	VendorId          string `form:"vendorIdOrder"  search:"type:order;column:vendor_id;table:goods"`
	CreatedAt         string `form:"createdAtOrder"  search:"type:order;column:created_at;table:goods"`
	UpdatedAt         string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:goods"`
	DeletedAt         string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:goods"`
	CreateBy          string `form:"createByOrder"  search:"type:order;column:create_by;table:goods"`
	CreateByName      string `form:"createByNameOrder"  search:"type:order;column:create_by_name;table:goods"`
	UpdateBy          string `form:"updateByOrder"  search:"type:order;column:update_by;table:goods"`
	UpdateByName      string `form:"updateByNameOrder"  search:"type:order;column:update_by_name;table:goods"`
}

func (m *GoodsGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type GoodsInsertReq struct {
	Id                int     `json:"-" comment:""` //
	SkuCode           string  `json:"skuCode" comment:"商品SKU" vd:"@:len($)>0&&mblen($)<10; msg:'SKU长度在1-10之间'"`
	SupplierSkuCode   string  `json:"supplierSkuCode" comment:"货主SKU" vd:"@:len($)>0&&mblen($)<=20; msg:'货主SKU长度在1-20之间'"`
	WarehouseCode     string  `json:"warehouseCode" comment:"仓库ID" vd:"@:len($)>0&&mblen($)<10; msg:'仓库必填'"`
	ProductNo         string  `json:"productNo" comment:"物料编码"`
	MarketPrice       float64 `json:"marketPrice" comment:"价格" vd:"@:$>0&&$<9999999; msg:'价格必填'"`
	PriceModifyReason string  `json:"priceModifyReason" comment:"价格调整备注"`
	ApproveStatus     int     `json:"approveStatus" comment:"审核状态 0 待审核  1 审核通过  2 审核失败" vd:"@:in($,0,1,2); msg:'审核状态错误'"`
	ApproveRemark     string  `json:"approveRemark" comment:"审核备注"`
	Status            int     `json:"status" default:"1" comment:"商品状态 0 禁用  1启用" vd:"@:in($,0,1); msg:'状态只能为0或1'"`
	OnlineStatus      int     `json:"onlineStatus" comment:"上下架状态  0未上架 1上架  2下架" vd:"@:in($,0,1,2); msg:'上下架状态错误'"`
	VendorId          int     `json:"vendorId" comment:"供应商ID" vd:"@:$>0; msg:'货主必填'"`
	common.ControlBy
}

func (s *GoodsInsertReq) Generate(model *models.Goods) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.SkuCode = strings.ToUpper(s.SkuCode)
	model.SupplierSkuCode = s.SupplierSkuCode
	model.WarehouseCode = s.WarehouseCode
	model.ProductNo = s.ProductNo
	model.MarketPrice = s.MarketPrice
	model.PriceModifyReason = s.PriceModifyReason
	model.ApproveStatus = s.ApproveStatus
	model.ApproveRemark = s.ApproveRemark
	model.Status = s.Status
	model.OnlineStatus = s.OnlineStatus
	model.VendorId = s.VendorId
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.CreateByName = s.CreateByName
}

func (s *GoodsInsertReq) GetId() interface{} {
	return s.Id
}

type GoodsUpdateReq struct {
	Id                int     `uri:"id" comment:""` //
	SkuCode           string  `json:"skuCode" comment:"商品SKU" vd:"@:len($)>0&&mblen($)<10; msg:'SKU长度在1-10之间'"`
	SupplierSkuCode   string  `json:"supplierSkuCode" comment:"货主SKU" vd:"@:len($)>0&&mblen($)<=20; msg:'货主SKU长度在1-20之间'"`
	WarehouseCode     string  `json:"warehouseCode" comment:"仓库ID" vd:"@:len($)>0&&mblen($)<10; msg:'仓库必填'"`
	ProductNo         string  `json:"productNo" comment:"物料编码"`
	MarketPrice       float64 `json:"marketPrice" comment:"价格" vd:"@:$>0&&$<9999999; msg:'价格必填'"`
	PriceModifyReason string  `json:"priceModifyReason" comment:"价格调整备注"`
	ApproveStatus     int     `json:"approveStatus" comment:"审核状态 0 待审核  1 审核通过  2 审核失败" vd:"@:in($,0,1,2); msg:'审核状态错误'"`
	ApproveRemark     string  `json:"approveRemark" comment:"审核备注"`
	Status            int     `json:"status" comment:"商品状态 0 禁用  1启用" vd:"@:in($,0,1); msg:'状态只能为0或1'"`
	OnlineStatus      int     `json:"onlineStatus" comment:"上下架状态  0未上架 1上架  2下架" vd:"@:in($,0,1,2); msg:'上下架状态错误'"`
	VendorId          int     `json:"vendorId" comment:"供应商ID" vd:"@:$>0; msg:'货主必填'"`
	common.ControlBy
}

func (s *GoodsUpdateReq) Generate(model *models.Goods) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.SkuCode = strings.ToUpper(s.SkuCode)
	model.SupplierSkuCode = s.SupplierSkuCode
	model.WarehouseCode = s.WarehouseCode
	model.ProductNo = s.ProductNo
	model.MarketPrice = s.MarketPrice
	model.PriceModifyReason = s.PriceModifyReason
	model.ApproveStatus = s.ApproveStatus
	model.ApproveRemark = s.ApproveRemark
	model.Status = s.Status
	model.OnlineStatus = s.OnlineStatus
	model.VendorId = s.VendorId
	model.UpdateBy = s.UpdateBy
	model.UpdateByName = s.UpdateByName // 添加这而，需要记录是被谁更新的
}

func (s *GoodsUpdateReq) GetId() interface{} {
	return s.Id
}

// GoodsGetReq 功能获取请求参数
type GoodsGetReq struct {
	Id int `uri:"id"`
}

func (s *GoodsGetReq) GetId() interface{} {
	return s.Id
}

type GoodsGetResp struct {
	models.Goods
	CompanyName   string `json:"companyName" gorm:"-"`
	WarehouseName string `json:"warehouseName" gorm:"-"`
}

// GoodsDeleteReq 功能删除请求参数
type GoodsDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *GoodsDeleteReq) GetId() interface{} {
	return s.Ids
}

// GoodsApproveReq 商品审核请求参数
type GoodsApproveReq struct {
	Ids           []int  `json:"ids"`
	ApproveStatus int    `json:"approveStatus" vd:"@:in($,0,1,2); msg:'审核状态错误！'"`
	ApproveRemark string `json:"approveRemark"`
	common.ControlBy
}

type GoodsUpdater struct {
	common.ControlBy
}

// GoodsExportResp 商品导出
type GoodsExportResp struct {
	SkuCode           string  `json:"skuCode"`
	GoodsName         string  `json:"goodsName"`
	BrandZh           string  `json:"brandZh"`
	MfgModel          string  `json:"mfgModel"`
	CompanyName       string  `json:"companyName"`
	WarehouseName     string  `json:"warehouseName"`
	VendorName        string  `json:"vendorName"`
	SupplierSkuCode   string  `json:"supplierSkuCode"`
	MarketPrice       float64 `json:"marketPrice"`
	PriceModifyReason string  `json:"priceModifyReason"`
	ProductNo         string  `json:"productNo"`
	OnlineStatus      string  `json:"onlineStatus"`
	Status            string  `json:"status"`
}

type FindSameWhere struct {
	SkuCode       string `json:"skuCode"`
	WarehouseCode string `json:"warehouseCode"`
	VendorId      int    `json:"vendorId"`
	Status        int    `json:"status"`
	ApproveStatus int    `json:"approveStatus"`
}

func FindSame(c *FindSameWhere) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db.Where("sku_code = ?", c.SkuCode)
		db.Where("warehouse_code = ?", c.WarehouseCode)
		db.Where("vendor_id = ?", c.VendorId)
		if c.Status == 1 {
			db.Where("status = ?", c.Status)
		}
		if c.ApproveStatus == 1 {
			db.Where("approve_status = ?", c.ApproveStatus)
		}
		return db
	}
}

type FindSkuSames struct {
	SkuCode       string
	WarehouseCode string
}

func FindSkuSame(c *FindSkuSames) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db.Where("sku_code = ?", c.SkuCode)
		db.Where("warehouse_code = ?", c.WarehouseCode)
		db.Where("status = 1")
		return db
	}
}

type OlineOfflineReq struct {
	Ids        []int `json:"ids"`
	ActionType int   `json:"actionType" comment:"类型" vd:"@:in($,1,2); msg:'操作类型只能为1或2！'"`
	common.ControlBy
}

type ImportGoods struct {
	SkuCode           string  `json:"skuCode" comment:"商品SKU" vd:"@:len($)>0&&mblen($)<7; msg:'SKU长度在1-7之间'"`
	CompanyName       string  `json:"companyName" comment:"公司名称" vd:"@:len($)>0&&mblen($)<=30; msg:'公司名称长度在1-30之间'"`
	WarehouseName     string  `json:"warehouseName" comment:"仓库名称" vd:"@:len($)>0&&mblen($)<=20; msg:'仓库名称长度在1-20之间'"`
	VendorName        string  `json:"vendorName" comment:"货主名称" vd:"@:len($)>0&&mblen($)<=50; msg:'货主长度在1-50之间'"`
	SupplierSkuCode   string  `json:"supplierSkuCode" comment:"货主SKU" vd:"@:len($)>0&&mblen($)<=20; msg:'货主SKU长度在1-20之间'"`
	MarketPrice       float64 `json:"marketPrice" comment:"价格" vd:"@:$>0 && $<=9999999; msg:'价格不为空且大于0！'"`
	ProductNo         string  `json:"productNo" comment:"物料编码"`
	PriceModifyReason string  `json:"priceModifyReason" comment:"价格调整备注"`
	Status            int     `json:"status" comment:"商品状态 0禁用 1启用" vd:"@:in($,0,1); msg:'状态只能为0或1'"  default:"1"`
	common.ControlBy
}

func (e *ImportGoods) Trim(data map[string]interface{}) {
	e.SkuCode = strings.ToUpper(strings.TrimSpace(data["skuCode"].(string)))
	e.CompanyName = strings.TrimSpace(data["companyName"].(string))
	e.WarehouseName = strings.TrimSpace(data["warehouseName"].(string))
	e.VendorName = strings.TrimSpace(data["vendorName"].(string))
	e.SupplierSkuCode = strings.TrimSpace(data["supplierSkuCode"].(string))
	e.ProductNo = strings.TrimSpace(data["productNo"].(string))
	e.PriceModifyReason = strings.TrimSpace(data["priceModifyReason"].(string))
	e.MarketPrice, _ = strconv.ParseFloat(data["marketPrice"].(string), 64)
	if data["status"] == "" {
		e.Status = 1
	} else {
		e.Status, _ = strconv.Atoi(data["status"].(string))
	}
}

func (e *ImportGoods) Generate(model *models.Goods) {
	model.SkuCode = e.SkuCode
	model.SupplierSkuCode = e.SupplierSkuCode
	model.MarketPrice = e.MarketPrice
	model.ProductNo = e.ProductNo
	model.PriceModifyReason = e.PriceModifyReason
	model.Status = e.Status
}

type ImportReq struct {
	Data []ImportGoods `json:"data"`
}

type GoodsInfo struct {
	FindSameWhere
}

type GoodsInfoReq struct {
	Query []GoodsInfo `json:"query"`
}

// GetGoodsByIdReq ID查询商品信息
type GetGoodsByIdReq struct {
	Ids []int `json:"ids"`
}

func (s *GetGoodsByIdReq) GetId() interface{} {
	return s.Ids
}

// GetGoodsBySkuCodeReq 仓库Code + sku查商品信息
type GetGoodsBySkuCodeReq struct {
	SkuCode       []string `json:"skuCode"`
	WarehouseCode string   `json:"warehouseCode"`
	Status        int      `json:"status"`
	OnlineStatus  int      `json:"onlineStatus"`
}

type FindGoodsReq struct {
	SkuCode         string
	WarehouseCode   string
	SupplierSkuCode string
}

func FindGoods(c *FindGoodsReq) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db.Where("sku_code = ?", c.SkuCode)
		db.Where("warehouse_code = ?", c.WarehouseCode)
		db.Where("supplier_sku_code = ?", c.SupplierSkuCode)
		return db
	}
}

func GoodsMakeCondition(c *GoodsGetPageReq, p *actions.DataPermission) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(c.Ids) > 0 {
			db.Where("goods.id in ?", c.Ids)
		}
		if c.NameZh != "" {
			db.Where("product.name_zh LIKE ?", "%"+c.NameZh+"%")
		}
		if c.OnlineStatus >= 0 {
			db.Where("goods.online_status = ?", c.OnlineStatus)
		}
		if c.ApproveStatus >= 0 {
			db.Where("goods.approve_status = ?", c.ApproveStatus)
		}
		if c.Status >= 0 {
			db.Where("goods.status = ?", c.Status)
		}
		if len(c.WarehouseCodeList) > 0 {
			authorityWarehouseId := utils.Split(p.AuthorityWarehouseId)
			warehouseCode := utils.Intersection(authorityWarehouseId, c.WarehouseCodeList)
			if len(warehouseCode) > 0 {
				db.Where("goods.warehouse_code in ?", warehouseCode)
			} else {
				db.Where("goods.warehouse_code = '-1'")
			}
		} else {
			db.Where("goods.warehouse_code in ?", utils.Split(p.AuthorityWarehouseId))
		}
		return db
	}
}
