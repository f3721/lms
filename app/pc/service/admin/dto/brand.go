package dto

import (
	"go-admin/app/pc/models"
	"go-admin/common/actions"
	"go-admin/common/dto"
	common "go-admin/common/models"
	"gorm.io/gorm"
)

type BrandGetPageReq struct {
	dto.Pagination `search:"-"`
	Id             int    `form:"id"  search:"type:exact;column:id;table:brand"`
	BrandZh        string `form:"brandZh"  search:"type:contains;column:brand_zh;table:brand"`
	BrandEn        string `form:"brandEn"  search:"type:contains;column:brand_en;table:brand"`
	Status         int    `form:"status"  search:"-"`
	BrandOrder
}

type BrandOrder struct {
	Id               string `form:"idOrder"  search:"type:order;column:id;table:brand"`
	BrandZh          string `form:"brandZhOrder"  search:"type:order;column:brand_zh;table:brand"`
	BrandEn          string `form:"brandEnOrder"  search:"type:order;column:brand_en;table:brand"`
	BrandDescription string `form:"brandDescriptionOrder"  search:"type:order;column:brand_description;table:brand"`
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

type BrandInsertReq struct {
	Id               int    `json:"-" comment:"品牌id"` // 品牌id
	BrandZh          string `json:"brandZh" comment:"品牌中文" vd:"@:len($)>0&&mblen($)<=100; msg:'品牌名称长度在1-100之间！'"`
	BrandEn          string `json:"brandEn" comment:"品牌英文"`
	FirstLetter      string `json:"-" comment:"首字母"`
	BrandDescription string `json:"brandDescription" comment:"品牌描述"`
	Status           int    `json:"status" comment:"激活状态(0:不使用1:使用)"`
	Confirm          int    `json:"confirm" comment:"二次确认"`
	common.ControlBy
}

func (s *BrandInsertReq) Generate(model *models.Brand) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.BrandZh = s.BrandZh
	model.BrandEn = s.BrandEn
	model.FirstLetter = s.FirstLetter
	model.BrandDescription = s.BrandDescription
	model.Status = s.Status
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.CreateByName = s.CreateByName
}

func (s *BrandInsertReq) GetId() interface{} {
	return s.Id
}

type BrandUpdateReq struct {
	Id               int    `uri:"id" comment:"品牌id"` // 品牌id
	BrandZh          string `json:"brandZh" comment:"品牌中文"`
	BrandEn          string `json:"brandEn" comment:"品牌英文"`
	FirstLetter      string `json:"-" comment:"首字母"`
	BrandDescription string `json:"brandDescription" comment:"品牌描述"`
	Status           int    `json:"status" comment:"激活状态(0:不使用1:使用)"`
	Confirm          int    `json:"confirm" comment:"二次确认"`
	common.ControlBy
}

func (s *BrandUpdateReq) Generate(model *models.Brand) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.BrandZh = s.BrandZh
	model.BrandEn = s.BrandEn
	model.FirstLetter = s.FirstLetter
	model.BrandDescription = s.BrandDescription
	model.Status = s.Status
	model.UpdateBy = s.UpdateBy
	model.UpdateByName = s.UpdateByName // 添加这而，需要记录是被谁更新的
}

func (s *BrandUpdateReq) GetId() interface{} {
	return s.Id
}

// BrandGetReq 功能获取请求参数
type BrandGetReq struct {
	Id int `uri:"id"`
}

func (s *BrandGetReq) GetId() interface{} {
	return s.Id
}

// BrandDeleteReq 功能删除请求参数
type BrandDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *BrandDeleteReq) GetId() interface{} {
	return s.Ids
}

func BrandMakeCondition(c *BrandGetPageReq, p *actions.DataPermission) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if c.Status >= 0 {
			db.Where("status = ?", c.Status)
		}
		return db
	}
}
