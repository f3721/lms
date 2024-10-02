package dto

import (
	"go-admin/app/pc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
	"gorm.io/gorm"
)

type UserCartGetPageReq struct {
	dto.Pagination `search:"-"`
	GoodsId        int    `form:"goodsId"  search:"type:exact;column:goods_id;table:user_cart" comment:"goods表 主键id"`
	UserId         int    `form:"userId"  search:"type:exact;column:user_id;table:user_cart" comment:"用户编号"`
	WarehouseCode  string `form:"warehouseCode"  search:"type:exact;column:warehouse_code;table:user_cart" comment:""`
	SkuCode        string `form:"skuCode"  search:"type:exact;column:sku_code;table:user_cart" comment:"商品订货号"`
	Selected       int    `form:"selected"  search:"type:exact;column:selected;table:user_cart" comment:"选中标记"`
}

func (m *UserCartGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type UserCartInsertReq struct {
	Id            int    `json:"-" comment:"主键"` // 主键
	GoodsId       int    `json:"goodsId" comment:"goods表 主键id" vd:"$>0; msg:'GoodsId不能为空'"`
	Quantity      int    `json:"quantity" comment:"商品数量" vd:"$>0; msg:'Quantity不能为空'"`
	WarehouseCode string `json:"warehouseCode" comment:"仓库code"`
	common.ControlBy
}

func (s *UserCartInsertReq) Generate(model *models.UserCart) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.GoodsId = s.GoodsId
	model.Selected = models.UserCartSelected1
	model.WarehouseCode = s.WarehouseCode
	model.UserId = s.CreateBy
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.CreateByName = s.CreateByName
}

func (s *UserCartInsertReq) GetId() interface{} {
	return s.Id
}

type UserCartBatchAddReq struct {
	Data          []UserCartBatchAddItemReq `json:"data"`
	WarehouseCode string                    `json:"warehouseCode" comment:"仓库code"`
	common.ControlBy
}

type UserCartBatchAddItemReq struct {
	SkuCode  string `json:"skuCode" vd:"len($)>0; msg:'skuCode不能为空'"`
	Quantity int    `json:"quantity" comment:"商品数量" vd:"$>0; msg:'Quantity不能为空'"`
}

type UserCartUpdateReq struct {
	GoodsId       int    `json:"goodsId" comment:"goods表 主键id" vd:"$>0; msg:'GoodsId不能为空'"`
	Quantity      int    `json:"quantity" comment:"商品数量" vd:"$>0; msg:'Quantity不能为空'"`
	WarehouseCode string `json:"warehouseCode" comment:"仓库code"`
	UserId        int    `json:"userId"`
}

func (s *UserCartUpdateReq) Generate(model *models.UserCart) {
	model.GoodsId = s.GoodsId
	model.WarehouseCode = s.WarehouseCode
	model.UserId = s.UserId
}

// UserCartGetReq 功能获取请求参数
type UserCartGetReq struct {
	Id int `uri:"id"`
}

func (s *UserCartGetReq) GetId() interface{} {
	return s.Id
}

// UserCartDeleteReq 功能删除请求参数
type UserCartDeleteReq struct {
	GoodsId       int    `json:"goodsId" comment:"goods表 主键id" vd:"$>0; msg:'GoodsId不能为空'"`
	WarehouseCode string `json:"warehouseCode" comment:"仓库code"`
	UserId        int    `json:"userId"`
}

func (s *UserCartDeleteReq) Generate(model *models.UserCart) {
	model.GoodsId = s.GoodsId
	model.WarehouseCode = s.WarehouseCode
	model.UserId = s.UserId
}

// UserCartSelectOneReq
type UserCartSelectOneReq struct {
	GoodsId       int    `json:"goodsId" comment:"goods表 主键id" vd:"$>0; msg:'GoodsId不能为空'"`
	WarehouseCode string `json:"warehouseCode" comment:"仓库code"`
	UserId        int    `json:"userId"`
}

func (s *UserCartSelectOneReq) Generate(model *models.UserCart) {
	model.GoodsId = s.GoodsId
	model.WarehouseCode = s.WarehouseCode
	model.UserId = s.UserId
}

type UserCartSelectAllReq struct {
	WarehouseCode string `json:"warehouseCode" comment:"仓库code"`
	UserId        int    `json:"userId"`
}

type UserCartUnSelectAllReq struct {
	WarehouseCode string `json:"warehouseCode" comment:"仓库code"`
	UserId        int    `json:"userId"`
}

type UserCartClearSelectReq struct {
	WarehouseCode string `json:"warehouseCode" comment:"仓库code"`
	UserId        int    `json:"userId"`
}

type UserCartVerifyCartReq struct {
	WarehouseCode string `form:"warehouseCode" comment:"仓库code"`
	UserId        int    `form:"userId"`
}

type UserCartVerifyBuyNowReq struct {
	WarehouseCode string `form:"warehouseCode" comment:"仓库code"`
	UserId        int    `form:"userId"`
	GoodsId       int    `form:"goodsId" vd:"$>0; msg:'GoodsId不能为空'"`
	Quantity      int    `form:"quantity" vd:"$>0; msg:'Quantity不能为空'"`
}

type UserCartGetProductForCartPageReq struct {
	WarehouseCode string `form:"warehouseCode" comment:"仓库code"`
	UserId        int    `form:"userId"`
}

type UserCartGetProductForCartPageResp struct {
	Products         []models.UserCartPageGoodsProduct `json:"products"`
	IsSelectAll      int                               `json:"isSelectAll"`
	IsExistNotSale   int                               `json:"isExistNotSale"`
	CartProductNum   int                               `json:"cartProductNum"`
	TotalProductNum  int                               `json:"totalProductNum"`
	TotalAmount      float64                           `json:"totalAmount"`
	TotalNakedAmount float64                           `json:"totalNakedAmount"`
	CheckStockStatus int                               `json:"checkStockStatus"`
}

type UserCartGetProductForOrderPageReq struct {
	WarehouseCode string `form:"warehouseCode" comment:"仓库code"`
	UserId        int    `form:"userId"`
}

type UserCartGetProductForOrderPageResp struct {
	Products         []models.UserCartGoodsProduct `json:"products"`
	TotalQuantity    int                           `json:"totalQuantity"`
	TotalAmount      float64                       `json:"totalAmount"`
	TotalNakedAmount float64                       `json:"totalNakedAmount"`
	CheckStockStatus int                           `json:"checkStockStatus"`
	AllowOrder       int                           `json:"allowOrder" default:"1"`
	AllowOrderMsg    string                        `json:"allowOrderMsg"`
}

type UserCartGetProductForOrderBuyNowReq struct {
	WarehouseCode string `form:"warehouseCode" comment:"仓库code"`
	UserId        int    `form:"userId"`
	GoodsId       int    `form:"goodsId" vd:"$>0; msg:'GoodsId不能为空'"`
	Quantity      int    `form:"quantity" vd:"$>0; msg:'Quantity不能为空'"`
}

type UserCartGetProductForOrderBuyNowResp struct {
	Products           []models.UserCartGoodsProduct `json:"products"`
	ProductNum         int                           `json:"productNum"`
	ProductCategoryNum int                           `json:"productCategoryNum"`
	TotalAmount        float64                       `json:"totalAmount"`
	TotalNakedAmount   float64                       `json:"totalNakedAmount"`
}

type UserCartGetProductForSaleMoqReq struct {
	SkuCodes      []string `json:"skuCodes" comment:"SkuCodes" vd:"len($)>0; msg:'skuCodes不能为空'"`
	WarehouseCode string   `json:"warehouseCode" comment:"仓库code"`
}

type UserCartGetProductForSaleMoqResp struct {
	SaleMoq map[string]int `json:"saleMoq" comment:"SaleMoq"`
}

type GetProductListReq struct {
	WarehouseCode string `form:"warehouseCode" comment:"仓库code"`
	UserId        int    `form:"userId"`
	GoodsId       []int  `form:"goodsId"`
}

func MakeProductListCondition(c *GetProductListReq) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db.Where("user_id = ?", c.UserId)
		db.Where("warehouse_code = ?", c.WarehouseCode)
		db.Where("goods_id in ?", c.GoodsId)
		return db
	}
}
