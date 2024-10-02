package dto

import (
	"go-admin/app/uc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type UserCollectGetPageReq struct {
	dto.Pagination `search:"-"`
	SkuCode        string `form:"skuCode"  search:"type:exact;column:sku_code;table:user_collect" comment:"产品SKU"`              //产品SKU
	UserId         int    `form:"userId"  search:"type:exact;column:user_id;table:user_collect" comment:"用户ID"`                 //用户ID
	GoodsId        int    `form:"goodsId"  search:"type:exact;column:goods_id;table:user_collect" comment:"商品表id"`              //商品表id
	WarehouseCode  string `form:"warehouseCode"  search:"type:exact;column:warehouse_code;table:user_collect" comment:"仓库code"` //仓库code
	CompanyId      int    `form:"-" search:"-" comment:"公司id"`

	FilterLevel1Catid int    `form:"filterLevel1Catid" search:"-"`
	FilterLevel2Catid int    `form:"filterLevel2Catid" search:"-"`
	FilterLevel3Catid int    `form:"filterLevel3Catid" search:"-"`
	FilterLevel4Catid int    `form:"filterLevel4Catid" search:"-"`
	FilterKeyword     string `form:"filterKeyword" search:"-"`
	IsShowStock       int    `form:"isShowStock" search:"-"` // 是否显示库存
	UserCollectOrder
}

type UserCollectOrder struct {
	Id            string `form:"idOrder"  search:"type:order;column:id;table:user_collect"`
	SkuCode       string `form:"skuCodeOrder"  search:"type:order;column:sku_code;table:user_collect"`
	UserId        string `form:"userIdOrder"  search:"type:order;column:user_id;table:user_collect"`
	GoodsId       string `form:"goodsIdOrder"  search:"type:order;column:goods_id;table:user_collect"`
	WarehouseCode string `form:"warehouseCodeOrder"  search:"type:order;column:warehouse_code;table:user_collect"`
	CreatedAt     string `form:"createdAtOrder"  search:"type:order;column:created_at;table:user_collect"`
	UpdatedAt     string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:user_collect"`
	CreateByName  string `form:"createByNameOrder"  search:"type:order;column:create_by_name;table:user_collect"`
	UpdateByName  string `form:"updateByNameOrder"  search:"type:order;column:update_by_name;table:user_collect"`
	DeletedAt     string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:user_collect"`
}

func (m *UserCollectGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type UserCollectGetListPageRes struct {
	models.UserCollect
	UserCollectGetListPageProduct
}

type UserCollectGetListPageProduct struct {
	ProductName        string  `json:"productName"`        // 商品名
	ProductImage       string  `json:"productImage"`       // 商品图片
	ProductSalesMoq    int     `json:"productSalesMoq"`    // 商品最小起订量
	ProductMarketPrice float64 `json:"productMarketPrice"` // 销售价
	ProductBrandZh     string  `json:"productBrandZh"`     //品牌名称
	ProductMfgModel    string  `json:"productMfgModel"`    // 产品型号
	ProductVendorName  string  `json:"productVendorName"`  //货主名称
	ProductStock       int     `json:"productStock"`       // 商品库存
}

type UserCollectInsertReq struct {
	GoodsId       int    `json:"goodsId" comment:"商品表id"` // 商品表id
	GoodsIds      []int  `json:"goodsIds"`                //商品表id 和 GoodsId任选一个
	WarehouseCode string `json:"-" comment:"仓库code"`      // 仓库code
	UserId        int    `json:"-" comment:"用户ID"`        // 用户ID
}

func (s *UserCollectInsertReq) Generate(model *models.UserCollect) {
	model.UserId = s.UserId
	model.GoodsId = s.GoodsId
	model.WarehouseCode = s.WarehouseCode
}

type UserCollectUpdateReq struct {
	Id            int    `uri:"id" comment:"编号"`                 // 编号
	SkuCode       string `json:"skuCode" comment:"产品SKU"`        // 产品SKU
	UserId        int    `json:"userId" comment:"用户ID"`          // 用户ID
	GoodsId       int    `json:"goodsId" comment:"商品表id"`        // 商品表id
	WarehouseCode string `json:"warehouseCode" comment:"仓库code"` // 仓库code
	common.ControlBy
}

func (s *UserCollectUpdateReq) Generate(model *models.UserCollect) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.SkuCode = s.SkuCode
	model.UserId = s.UserId
	model.GoodsId = s.GoodsId
	model.WarehouseCode = s.WarehouseCode
	model.UpdateByName = s.UpdateByName // 添加这而，需要记录是被谁更新的
}

func (s *UserCollectUpdateReq) GetId() interface{} {
	return s.Id
}

// UserCollectGetReq 功能获取请求参数
type UserCollectGetReq struct {
	Id int `uri:"id"`
}

func (s *UserCollectGetReq) GetId() interface{} {
	return s.Id
}

// UserCollectDeleteReq 功能删除请求参数
type UserCollectDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *UserCollectDeleteReq) GetId() interface{} {
	return s.Ids
}

// UserCollectDeleteReq 功能删除请求参数
type UserCollectDeleteGoodsIdsReq struct {
	GoodsIds []int `json:"goodsIds"`
}

// UserCollectGetReq 功能获取请求参数
type UserCollectGetGoodsIsCollected struct {
	GoodsIds string `form:"goodsIds" vd:"len($)>0" `
	UserId   int    `form:"userId" vd:"$>0" `
}

type UserCollectGetIsUserCollectResData struct {
	GoodsId   int  `json:"goodsId"`
	IsCollect bool `json:"isCollect"`
}
