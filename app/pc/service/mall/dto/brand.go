package dto

import (
	"go-admin/common/dto"
)

type BrandGetPageReq struct {
	dto.Pagination `search:"-"`
	BrandZh        string `form:"brandZh"  search:"type:contains;column:brand_zh;table:brand"`
	BrandEn        string `form:"brandEn"  search:"type:contains;column:brand_en;table:brand"`
	Status         int    `form:"status"  search:"type:exact;column:status;table:brand"`
	BrandOrder
}

type BrandOrder struct {
	Id               string `form:"idOrder"  search:"type:order;column:id;table:brand"`
	BrandZh          string `form:"brandZhOrder"  search:"type:order;column:brand_zh;table:brand"`
	BrandEn          string `form:"brandEnOrder"  search:"type:order;column:brand_en;table:brand"`
	BrandDescription string `form:"brandDesciptionOrder"  search:"type:order;column:brand_desciption;table:brand"`
	Status           string `form:"statusOrder"  search:"type:order;column:status;table:brand"`
	CreatedAt        string `form:"createdAtOrder"  search:"type:order;column:created_at;table:brand"`
	UpdatedAt        string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:brand"`
	DeletedAt        string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:brand"`
	CreateBy         string `form:"createByOrder"  search:"type:order;column:create_by;table:brand"`
	CreateByName     string `form:"createByNameOrder"  search:"type:order;column:create_by_name;table:brand"`
	UpdateBy         string `form:"updateByOrder"  search:"type:order;column:update_by;table:brand"`
	UpdateByName     string `form:"updateByNameOrder"  search:"type:order;column:update_by_name;table:brand"`
}

func (m *BrandGetPageReq) GetNeedSearch() interface{} {
	return *m
}

// BrandGetReq 功能获取请求参数
type BrandGetReq struct {
	Id int `uri:"id"`
}

func (s *BrandGetReq) GetId() interface{} {
	return s.Id
}

type BrandInfo struct {
	BrandId int    `json:"brandId"`
	BrandZh string `json:"brandZh"`
}
