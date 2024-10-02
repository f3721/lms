package dto

import (
	"go-admin/app/oc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type CsApplyDetailGetPageReq struct {
	dto.Pagination `search:"-"`
	CsNo           string `form:"csNo"  search:"type:exact;column:cs_no;table:cs_apply_detail" comment:"售后申请编号"`
	SkuCode        string `form:"skuCode"  search:"type:exact;column:sku_code;table:cs_apply_detail" comment:"商品sku"`
	CsType         int    `form:"csType"  search:"type:exact;column:cs_type;table:cs_apply_detail" comment:"类型：0-用户申请、1-实际售后"`
	ProductPic     string `form:"productPic"  search:"type:exact;column:product_pic;table:cs_apply_detail" comment:"商品图片url"`
	DeliveryTime   string `form:"deliveryTime"  search:"type:exact;column:delivery_time;table:cs_apply_detail" comment:"货期"`
	WarehouseCode  string `form:"warehouseCode"  search:"type:exact;column:warehouse_code;table:cs_apply_detail" comment:"到货仓库"`
	GoodsId        int    `form:"goodsId"  search:"type:exact;column:goods_id;table:cs_apply_detail" comment:""`
	VendorSkuCode  string `form:"vendorSkuCode"  search:"type:exact;column:vendor_sku_code;table:cs_apply_detail" comment:"货主sku"`
	ProductNo      string `form:"productNo"  search:"type:exact;column:product_no;table:cs_apply_detail" comment:"物料编码"`
	IsDefective    int    `form:"isDefective"  search:"type:exact;column:is_defective;table:cs_apply_detail" comment:"是否次品 1是 0否"`
	CsApplyDetailOrder
}

type CsApplyDetailOrder struct {
	Id            string `form:"idOrder"  search:"type:order;column:id;table:cs_apply_detail"`
	CsNo          string `form:"csNoOrder"  search:"type:order;column:cs_no;table:cs_apply_detail"`
	ProductModel  string `form:"productModelOrder"  search:"type:order;column:product_model;table:cs_apply_detail"`
	SkuCode       string `form:"skuCodeOrder"  search:"type:order;column:sku_code;table:cs_apply_detail"`
	ProductName   string `form:"productNameOrder"  search:"type:order;column:product_name;table:cs_apply_detail"`
	Quantity      string `form:"quantityOrder"  search:"type:order;column:quantity;table:cs_apply_detail"`
	Unit          string `form:"unitOrder"  search:"type:order;column:unit;table:cs_apply_detail"`
	RefundAmt     string `form:"refundAmtOrder"  search:"type:order;column:refund_amt;table:cs_apply_detail"`
	ReparationAmt string `form:"reparationAmtOrder"  search:"type:order;column:reparation_amt;table:cs_apply_detail"`
	Remark        string `form:"remarkOrder"  search:"type:order;column:remark;table:cs_apply_detail"`
	CreatedAt     string `form:"createdAtOrder"  search:"type:order;column:created_at;table:cs_apply_detail"`
	UpdatedAt     string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:cs_apply_detail"`
	CsType        string `form:"csTypeOrder"  search:"type:order;column:cs_type;table:cs_apply_detail"`
	SalePrice     string `form:"salePriceOrder"  search:"type:order;column:sale_price;table:cs_apply_detail"`
	ProductPic    string `form:"productPicOrder"  search:"type:order;column:product_pic;table:cs_apply_detail"`
	BrandName     string `form:"brandNameOrder"  search:"type:order;column:brand_name;table:cs_apply_detail"`
	DeliveryTime  string `form:"deliveryTimeOrder"  search:"type:order;column:delivery_time;table:cs_apply_detail"`
	WarehouseCode string `form:"warehouseCodeOrder"  search:"type:order;column:warehouse_code;table:cs_apply_detail"`
	GoodsId       string `form:"goodsIdOrder"  search:"type:order;column:goods_id;table:cs_apply_detail"`
	VendorName    string `form:"vendorNameOrder"  search:"type:order;column:vendor_name;table:cs_apply_detail"`
	VendorSkuCode string `form:"vendorSkuCodeOrder"  search:"type:order;column:vendor_sku_code;table:cs_apply_detail"`
	ProductNo     string `form:"productNoOrder"  search:"type:order;column:product_no;table:cs_apply_detail"`
	IsDefective   string `form:"isDefectiveOrder"  search:"type:order;column:is_defective;table:cs_apply_detail"`
	Pics          string `form:"picsOrder"  search:"type:order;column:pics;table:cs_apply_detail"`
}

func (m *CsApplyDetailGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type CsApplyDetailInsertReq struct {
	Id            int     `json:"-" comment:""` //
	CsNo          string  `json:"csNo" comment:"售后申请编号"`
	ProductModel  string  `json:"productModel" comment:"售后申请产品名原始型号"`
	SkuCode       string  `json:"skuCode" comment:"商品sku"`
	ProductName   string  `json:"productName" comment:"商品名称"`
	Quantity      int     `json:"quantity" comment:"数量"`
	Unit          string  `json:"unit" comment:"单位"`
	RefundAmt     string  `json:"refundAmt" comment:"退款金额"`
	ReparationAmt string  `json:"reparationAmt" comment:"赔款金额"`
	Remark        string  `json:"remark" comment:"备注"`
	CsType        int     `json:"csType" comment:"类型：0-用户申请、1-实际售后"`
	SalePrice     float64 `json:"salePrice" comment:"销售价"`
	ProductPic    string  `json:"productPic" comment:"商品图片url"`
	BrandName     string  `json:"brandName" comment:"产品的品牌"`
	DeliveryTime  string  `json:"deliveryTime" comment:"货期"`
	WarehouseCode string  `json:"warehouseCode" comment:"到货仓库"`
	GoodsId       int     `json:"goodsId" comment:""`
	VendorName    string  `json:"vendorName" comment:"货主"`
	VendorSkuCode string  `json:"vendorSkuCode" comment:"货主sku"`
	ProductNo     string  `json:"productNo" comment:"物料编码"`
	IsDefective   int     `json:"isDefective" comment:"是否次品 1是 0否"`
	Pics          string  `json:"pics" comment:"售后申请的图片"`
	common.ControlBy
}

func (s *CsApplyDetailInsertReq) Generate(model *models.CsApplyDetail) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.CsNo = s.CsNo
	model.ProductModel = s.ProductModel
	model.SkuCode = s.SkuCode
	model.ProductName = s.ProductName
	model.Quantity = s.Quantity
	model.Unit = s.Unit
	model.RefundAmt = s.RefundAmt
	model.ReparationAmt = s.ReparationAmt
	model.Remark = s.Remark
	model.CsType = s.CsType
	model.SalePrice = s.SalePrice
	model.ProductPic = s.ProductPic
	model.BrandName = s.BrandName
	model.DeliveryTime = s.DeliveryTime
	model.WarehouseCode = s.WarehouseCode
	model.GoodsId = s.GoodsId
	model.VendorName = s.VendorName
	model.VendorSkuCode = s.VendorSkuCode
	model.ProductNo = s.ProductNo
	model.IsDefective = s.IsDefective
	model.Pics = s.Pics
}

func (s *CsApplyDetailInsertReq) GetId() interface{} {
	return s.Id
}

type CsApplyDetailUpdateReq struct {
	Id            int     `uri:"id" comment:""` //
	CsNo          string  `json:"csNo" comment:"售后申请编号"`
	ProductModel  string  `json:"productModel" comment:"售后申请产品名原始型号"`
	SkuCode       string  `json:"skuCode" comment:"商品sku"`
	ProductName   string  `json:"productName" comment:"商品名称"`
	Quantity      int     `json:"quantity" comment:"数量"`
	Unit          string  `json:"unit" comment:"单位"`
	RefundAmt     string  `json:"refundAmt" comment:"退款金额"`
	ReparationAmt string  `json:"reparationAmt" comment:"赔款金额"`
	Remark        string  `json:"remark" comment:"备注"`
	CsType        int     `json:"csType" comment:"类型：0-用户申请、1-实际售后"`
	SalePrice     float64 `json:"salePrice" comment:"销售价"`
	ProductPic    string  `json:"productPic" comment:"商品图片url"`
	BrandName     string  `json:"brandName" comment:"产品的品牌"`
	DeliveryTime  string  `json:"deliveryTime" comment:"货期"`
	WarehouseCode string  `json:"warehouseCode" comment:"到货仓库"`
	GoodsId       int     `json:"goodsId" comment:""`
	VendorName    string  `json:"vendorName" comment:"货主"`
	VendorSkuCode string  `json:"vendorSkuCode" comment:"货主sku"`
	ProductNo     string  `json:"productNo" comment:"物料编码"`
	IsDefective   int     `json:"isDefective" comment:"是否次品 1是 0否"`
	Pics          string  `json:"pics" comment:"售后申请的图片"`
	common.ControlBy
}

func (s *CsApplyDetailUpdateReq) Generate(model *models.CsApplyDetail) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.CsNo = s.CsNo
	model.ProductModel = s.ProductModel
	model.SkuCode = s.SkuCode
	model.ProductName = s.ProductName
	model.Quantity = s.Quantity
	model.Unit = s.Unit
	model.RefundAmt = s.RefundAmt
	model.ReparationAmt = s.ReparationAmt
	model.Remark = s.Remark
	model.CsType = s.CsType
	model.SalePrice = s.SalePrice
	model.ProductPic = s.ProductPic
	model.BrandName = s.BrandName
	model.DeliveryTime = s.DeliveryTime
	model.WarehouseCode = s.WarehouseCode
	model.GoodsId = s.GoodsId
	model.VendorName = s.VendorName
	model.VendorSkuCode = s.VendorSkuCode
	model.ProductNo = s.ProductNo
	model.IsDefective = s.IsDefective
	model.Pics = s.Pics
}

func (s *CsApplyDetailUpdateReq) GetId() interface{} {
	return s.Id
}

// CsApplyDetailGetReq 功能获取请求参数
type CsApplyDetailGetReq struct {
	CsNo string `uri:"csNo"`
}

// CsApplyDetailGetReq 功能获取请求参数
type CsApplyDetailGetRes struct {
	ApplicationList *[]models.CsApplyDetail `json:"applicationList"` // 用户申请的
	AfterSalesList  *[]models.CsApplyDetail `json:"afterSalesList"`  // 实际售后的数据
}

// CsApplyDetailDeleteReq 功能删除请求参数
type CsApplyDetailDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *CsApplyDetailDeleteReq) GetId() interface{} {
	return s.Ids
}

type AfterReturnProductsQuantity map[string]struct {
	SkuCode  string
	Quantity int
}
