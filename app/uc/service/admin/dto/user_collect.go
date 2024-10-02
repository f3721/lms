package dto

import (
	"go-admin/app/uc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type UserCollectGetPageReq struct {
	dto.Pagination `search:"-"`
	SkuCode        string `form:"skuCode"  search:"type:exact;column:sku_code;table:user_collect" comment:"产品SKU"`
	UserId         int    `form:"userId"  search:"type:exact;column:user_id;table:user_collect" comment:"用户ID"`
	GoodsId        int    `form:"goodsId"  search:"type:exact;column:goods_id;table:user_collect" comment:"商品表id"`
	WarehouseCode  string `form:"warehouseCode"  search:"type:exact;column:warehouse_code;table:user_collect" comment:"仓库code"`
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
	ProductName  string `json:"productName"`  // 商品名
	ProductImage string `json:"productImage"` // 商品图片
}

type UserCollectInsertReq struct {
	Id            int    `json:"-" comment:"编号"` // 编号
	SkuCode       string `json:"skuCode" comment:"产品SKU"`
	UserId        int    `json:"userId" comment:"用户ID"`
	GoodsId       int    `json:"goodsId" comment:"商品表id"`
	WarehouseCode string `json:"warehouseCode" comment:"仓库code"`
	common.ControlBy
}

func (s *UserCollectInsertReq) Generate(model *models.UserCollect) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.SkuCode = s.SkuCode
	model.UserId = s.UserId
	model.GoodsId = s.GoodsId
	model.WarehouseCode = s.WarehouseCode
	model.CreateByName = s.CreateByName
}

func (s *UserCollectInsertReq) GetId() interface{} {
	return s.Id
}

type UserCollectUpdateReq struct {
	Id            int    `uri:"id" comment:"编号"` // 编号
	SkuCode       string `json:"skuCode" comment:"产品SKU"`
	UserId        int    `json:"userId" comment:"用户ID"`
	GoodsId       int    `json:"goodsId" comment:"商品表id"`
	WarehouseCode string `json:"warehouseCode" comment:"仓库code"`
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
