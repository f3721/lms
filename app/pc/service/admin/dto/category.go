package dto

import (
	"go-admin/app/pc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
	"gorm.io/gorm"
)

type CategoryGetPageReq struct {
	dto.Pagination  `search:"-"`
	Id              int    `form:"id"  search:"type:exact;column:id;table:category"`
	NameZh          string `form:"nameZh"  search:"type:contains;column:name_zh;table:category"`
	Tax             string `form:"tax"  search:"type:contains;column:tax;table:category"`
	CategoryTaxCode string `form:"categoryTaxCode"  search:"type:contains;column:category_tax_code;table:category"`
	Status          int    `form:"status"  search:"-"`
	ParentId        int    `form:"parentId"  search:"type:exact;column:parent_id;table:category"`
	CateLevel       int    `form:"cateLevel"  search:"type:exact;column:cate_level;table:category"`
	CategoryOrder
}

type CategoryOrder struct {
	Id              string `form:"idOrder"  search:"type:order;column:id;table:category"`
	CateLevel       string `form:"cateLevelOrder"  search:"type:order;column:cate_level;table:category"`
	Seq             string `form:"seqOrder"  search:"type:order;column:seq;table:category"`
	NameZh          string `form:"nameZhOrder"  search:"type:order;column:name_zh;table:category"`
	NameEn          string `form:"nameEnOrder"  search:"type:order;column:name_en;table:category"`
	ParentId        string `form:"parentIdOrder"  search:"type:order;column:parent_id;table:category"`
	Description     string `form:"descriptionOrder"  search:"type:order;column:description;table:category"`
	Status          string `form:"statusOrder"  search:"type:order;column:status;table:category"`
	KeyWords        string `form:"keyWordsOrder"  search:"type:order;column:key_words;table:category"`
	Tax             string `form:"taxOrder"  search:"type:order;column:tax;table:category"`
	CategoryTaxCode string `form:"categoryTaxCodeOrder"  search:"type:order;column:category_tax_code;table:category"`
	CreatedAt       string `form:"createdAtOrder"  search:"type:order;column:created_at;table:category"`
	UpdatedAt       string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:category"`
	DeletedAt       string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:category"`
	CreateBy        string `form:"createByOrder"  search:"type:order;column:create_by;table:category"`
	CreateByName    string `form:"createByNameOrder"  search:"type:order;column:create_by_name;table:category"`
	UpdateBy        string `form:"updateByOrder"  search:"type:order;column:update_by;table:category"`
	UpdateByName    string `form:"updateByNameOrder"  search:"type:order;column:update_by_name;table:category"`
}

func (m *CategoryGetPageReq) GetNeedSearch() interface{} {
	return *m
}

// CategoryGetReq 功能获取请求参数
type CategoryGetReq struct {
	Id              int    `uri:"id" search:"-"`
	CategoryId      int    `form:"categoryId" search:"type:exact;column:id;table:category"`
	NameZh          string `form:"nameZh"  search:"type:contains;column:name_zh;table:category"`
	Tax             string `form:"tax"  search:"type:contains;column:tax;table:category"`
	CategoryTaxCode string `form:"categoryTaxCode"  search:"type:contains;column:category_tax_code;table:category"`
	Status          string `form:"status"  search:"type:exact;column:status;table:category"`
	ParentId        string `form:"parentId"  search:"type:exact;column:parent_id;table:category"`
	CateLevel       int    `form:"cateLevel"  search:"type:exact;column:cate_level;table:category"`
	CategoryOrder
}

func (s *CategoryGetReq) GetId() interface{} {
	return s.Id
}

func (m *CategoryGetReq) GetNeedSearch() interface{} {
	return *m
}

type CategoryInsertReq struct {
	Id              int                    `json:"-" comment:"id"` // id
	CateLevel       int                    `json:"-"`
	Level1Catid     int                    `json:"level1Catid" comment:"一级分类"`
	Level2Catid     int                    `json:"level2Catid" comment:"二级分类"`
	Level3Catid     int                    `json:"level3Catid" comment:"三级分类"`
	Level4Catid     int                    `json:"level4Catid" comment:"四级分类"`
	Seq             int                    `json:"seq" comment:"序列"`
	NameZh          string                 `json:"nameZh" comment:"中文名" vd:"@:len($)>0&&mblen($)<=255; msg:'中文名长度在0-255之间'"`
	NameEn          string                 `json:"nameEn" comment:"英文名" vd:"@:mblen($)<=255; msg:'英文名长度小于255'"`
	ParentId        int                    `json:"-" comment:"父类id"`
	Description     string                 `json:"description" comment:"描述" vd:"@:mblen($)<=255; msg:'描述长度小于255'"`
	Status          int                    `json:"status" comment:"产线状态" vd:"@:in($,0,1); msg:'状态只能为0或1'"`
	KeyWords        string                 `json:"keyWords" comment:"关键字"`
	Tax             string                 `json:"tax" comment:"产线税率：默认空 值有（0.13,0.06,0.09）"`
	CategoryTaxCode string                 `json:"categoryTaxCode" comment:"产线税号" vd:"@:mblen($)<=50; msg:'产线税号长度小于50'"`
	MediaInstance   MediaInstanceInsertReq `json:"mediaInstance"`
	common.ControlBy
}

func (s *CategoryInsertReq) InsertGenerate() {
	if s.Level3Catid > 0 {
		s.CateLevel = 4
		s.ParentId = s.Level3Catid
	} else if s.Level2Catid > 0 {
		s.CateLevel = 3
		s.ParentId = s.Level2Catid
	} else if s.Level1Catid > 0 {
		s.CateLevel = 2
		s.ParentId = s.Level1Catid
	} else {
		s.CateLevel = 1
		s.ParentId = 0
	}
}

func (s *CategoryInsertReq) Generate(model *models.Category) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.CateLevel = s.CateLevel
	model.Seq = s.Seq
	model.NameZh = s.NameZh
	model.NameEn = s.NameEn
	model.ParentId = s.ParentId
	model.Description = s.Description
	model.Status = s.Status
	model.KeyWords = s.KeyWords
	model.Tax = s.Tax
	model.CategoryTaxCode = s.CategoryTaxCode
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.CreateByName = s.CreateByName
}

func (s *CategoryInsertReq) GetId() interface{} {
	return s.Id
}

type CategoryUpdateReq struct {
	Id              int                    `uri:"id" comment:"id"` // id
	CateLevel       int                    `json:"-"`
	Level1Catid     int                    `json:"level1Catid" comment:"一级分类"`
	Level2Catid     int                    `json:"level2Catid" comment:"二级分类"`
	Level3Catid     int                    `json:"level3Catid" comment:"三级分类"`
	Level4Catid     int                    `json:"level4Catid" comment:"四级分类"`
	Seq             int                    `json:"seq" comment:"序列"`
	NameZh          string                 `json:"nameZh" comment:"中文名"`
	NameEn          string                 `json:"nameEn" comment:"英文名"`
	ParentId        int                    `json:"parentId" comment:"父类id"`
	Description     string                 `json:"description" comment:"描述"`
	Status          int                    `json:"status" comment:"产线状态"`
	KeyWords        string                 `json:"keyWords" comment:"关键字"`
	Tax             string                 `json:"tax" comment:"产线税率：默认空 值有（0.13,0.06.0.09）"`
	CategoryTaxCode string                 `json:"categoryTaxCode" comment:"产线税号"`
	MediaInstance   MediaInstanceInsertReq `json:"mediaInstance"`
	common.ControlBy
}

func (s *CategoryUpdateReq) UpdateGenerate() {
	if s.Level3Catid > 0 {
		s.CateLevel = 4
		s.ParentId = s.Level3Catid
	} else if s.Level2Catid > 0 {
		s.CateLevel = 3
		s.ParentId = s.Level2Catid
	} else if s.Level1Catid > 0 {
		s.CateLevel = 2
		s.ParentId = s.Level1Catid
	} else {
		s.CateLevel = 1
		s.ParentId = 0
	}
}

func (s *CategoryUpdateReq) Generate(model *models.Category) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.CateLevel = s.CateLevel
	model.Seq = s.Seq
	model.NameZh = s.NameZh
	model.NameEn = s.NameEn
	model.ParentId = s.ParentId
	model.Description = s.Description
	model.Status = s.Status
	model.KeyWords = s.KeyWords
	model.Tax = s.Tax
	model.CategoryTaxCode = s.CategoryTaxCode
	model.UpdateBy = s.UpdateBy
	model.UpdateByName = s.UpdateByName // 添加这而，需要记录是被谁更新的
}

func (s *CategoryUpdateReq) GetId() interface{} {
	return s.Id
}

// CategoryDeleteReq 功能删除请求参数
type CategoryDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *CategoryDeleteReq) GetId() interface{} {
	return s.Ids
}

type CategoryChildList struct {
	Id     int    `json:"id"`
	NameZh string `json:"name_zh"`
}

type CategoryList struct {
	Id1     int
	Id2     int
	Id3     int
	Id4     int
	NameZh1 string
	NameZh2 string
	NameZh3 string
	NameZh4 string
}

type CategoryLevel struct {
	Level1CatId int `json:"level1CatId" gorm:"-" comment:"一级分类"`
	Level2CatId int `json:"level2CatId" gorm:"-" comment:"二级分类"`
	Level3CatId int `json:"level3CatId" gorm:"-" comment:"三级分类"`
	Level4CatId int `json:"level4CatId" gorm:"-" comment:"四级分类"`
}

var StatusMap = []string{"禁用", "启用"}

type CategoryGetPageResp struct {
	Id              int                   `json:"id" gorm:"type:int(11) unsigned;"`
	CateLevel       int                   `json:"cateLevel" gorm:"type:int(11) unsigned;comment:层级"`
	Seq             int                   `json:"seq" gorm:"type:smallint(5) unsigned;comment:序列"`
	NameZh          string                `json:"nameZh" gorm:"type:varchar(255);comment:中文名"`
	NameEn          string                `json:"nameEn" gorm:"type:varchar(255);comment:英文名"`
	ParentId        int                   `json:"parentId" gorm:"type:int(11) unsigned;comment:父类id"`
	Status          int                   `json:"status" gorm:"type:tinyint(1);comment:产线状态"`
	Tax             string                `json:"tax" gorm:"type:varchar(4);comment:产线税率：默认空 值有（0.13,0.06.0.09）"`
	CategoryTaxCode string                `json:"categoryTaxCode" gorm:"type:varchar(50);comment:产线税号"`
	StatusTxt       string                `gorm:"-" json:"statusTxt"`
	TaxTxt          string                `gorm:"-" json:"taxTxt"`
	Addchild        bool                  `gorm:"-" json:"addchild"`
	Haschild        bool                  `json:"haschild"`
	Children        []CategoryGetPageResp `gorm:"-" json:"children"`
}

type CategoryGetResp struct {
	models.Category
	CategoryLevel
	Addchild bool `gorm:"-" json:"addchild"`
}

type CategoryPathReq struct {
	Id int `uri:"id"`
}

func (s *CategoryPathReq) GetId() interface{} {
	return s.Id
}

type CategoryPath struct {
	Id         int    `json:"id"`
	CategoryId int    `json:"categoryId"`
	Name       string `json:"name"`
	Value      int    `json:"value"`
	PathName   string `json:"pathName"`
}

// CategorySort 批量排序
type CategorySort struct {
	CategoryId int `json:"categoryId"`
	Seq        int `json:"seq"`
}

type SortReq struct {
	Sort []CategorySort `json:"sort"`
	common.ControlBy
}

type Ids struct {
	Id int
}

func MakeCategoryCondition(c *CategoryGetPageReq) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if c.Status >= 0 {
			db.Where("category.status = ?", c.Status)
		}
		return db
	}
}

func MakeCategoryGetReqCondition(c *CategoryGetReq) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		//if c.CategoryId <= 0 && c.Status == "" && c.CategoryTaxCode == "" && c.NameZh == "" && c.Tax == "" {
		db.Where("category.parent_id = ?", c.GetId())
		//}
		return db
	}
}
