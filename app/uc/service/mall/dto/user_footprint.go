package dto

import (
	"go-admin/app/uc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type UserFootprintGetPageReq struct {
	dto.Pagination `search:"-"`
	GoodsId        int    `form:"goodsId"  search:"type:exact;column:goods_id;table:user_footprint" comment:""` //
	SkuCode        string `form:"skuCode"  search:"type:exact;column:sku_code;table:user_footprint" comment:""` //
	UserId         int    `form:"-"  search:"type:exact;column:user_id;table:user_footprint" comment:""`        //
	UserFootprintOrder
}

type UserFootprintOrder struct {
	Id           string `form:"idOrder"  search:"type:order;column:id;table:user_footprint"`
	GoodsId      string `form:"goodsIdOrder"  search:"type:order;column:goods_id;table:user_footprint"`
	SkuCode      string `form:"skuCodeOrder"  search:"type:order;column:sku_code;table:user_footprint"`
	UserId       string `form:"userIdOrder"  search:"type:order;column:user_id;table:user_footprint"`
	CreatedAt    string `form:"createdAtOrder"  search:"type:order;column:created_at;table:user_footprint"`
	UpdatedAt    string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:user_footprint"`
	CreateByName string `form:"createByNameOrder"  search:"type:order;column:create_by_name;table:user_footprint"`
	UpdateByName string `form:"updateByNameOrder"  search:"type:order;column:update_by_name;table:user_footprint"`
	DeletedAt    string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:user_footprint"`
}

func (m *UserFootprintGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type UserFootprintGetListPageRes struct {
	models.UserFootprint
	UserFootprintGetListPageProduct
}

type UserFootprintGetListPageProduct struct {
	ProductName        string  `json:"productName"`        // 商品名
	ProductImage       string  `json:"productImage"`       // 商品图片
	ProductSalesMoq    int     `json:"productSalesMoq"`    // 商品最小起订量
	ProductMarketPrice float64 `json:"productMarketPrice"` // 销售价
}

type UserFootprintInsertReq struct {
	GoodsId int `json:"goodsId" comment:""` //
	UserId  int `json:"-" comment:""`       //
	common.ControlBy
}

func (s *UserFootprintInsertReq) Generate(model *models.UserFootprint) {
	model.GoodsId = s.GoodsId
	model.UserId = s.UserId
}

type UserFootprintUpdateReq struct {
	Id      int    `uri:"id" comment:""`       //
	GoodsId int    `json:"goodsId" comment:""` //
	SkuCode string `json:"skuCode" comment:""` //
	UserId  int    `json:"userId" comment:""`  //
	common.ControlBy
}

func (s *UserFootprintUpdateReq) Generate(model *models.UserFootprint) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.GoodsId = s.GoodsId
	model.SkuCode = s.SkuCode
	model.UserId = s.UserId
}

func (s *UserFootprintUpdateReq) GetId() interface{} {
	return s.Id
}

// UserFootprintGetReq 功能获取请求参数
type UserFootprintGetReq struct {
	Id int `uri:"id"`
}

func (s *UserFootprintGetReq) GetId() interface{} {
	return s.Id
}

// UserFootprintDeleteReq 功能删除请求参数
type UserFootprintDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *UserFootprintDeleteReq) GetId() interface{} {
	return s.Ids
}
