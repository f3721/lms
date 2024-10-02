package dto

import (

	"go-admin/app/pc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type ProductCategoryGetPageReq struct {
	dto.Pagination     `search:"-"`
    ProductCategoryOrder
}

type ProductCategoryOrder struct {
    Id string `form:"idOrder"  search:"type:order;column:id;table:product_category"`
    SkuCode string `form:"skuCodeOrder"  search:"type:order;column:sku_code;table:product_category"`
    CategoryId string `form:"categoryIdOrder"  search:"type:order;column:category_id;table:product_category"`
    MainCateFlag string `form:"mainCateFlagOrder"  search:"type:order;column:main_cate_flag;table:product_category"`
    CreatedAt string `form:"createdAtOrder"  search:"type:order;column:created_at;table:product_category"`
    UpdatedAt string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:product_category"`
    DeletedAt string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:product_category"`
    CreateBy string `form:"createByOrder"  search:"type:order;column:create_by;table:product_category"`
    CreateByName string `form:"createByNameOrder"  search:"type:order;column:create_by_name;table:product_category"`
    UpdateBy string `form:"updateByOrder"  search:"type:order;column:update_by;table:product_category"`
    UpdateByName string `form:"updateByNameOrder"  search:"type:order;column:update_by_name;table:product_category"`
    
}

func (m *ProductCategoryGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type ProductCategoryInsertReq struct {
    Id int `json:"-" comment:"id"` // id
    SkuCode string `json:"skuCode" comment:"sku"`
    CategoryId int `json:"categoryId" comment:"产线id"`
    MainCateFlag int `json:"mainCateFlag" comment:"主产线标志0否1是"`
    common.ControlBy
}

func (s *ProductCategoryInsertReq) Generate(model *models.ProductCategory)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.SkuCode = s.SkuCode
    model.CategoryId = s.CategoryId
    model.MainCateFlag = s.MainCateFlag
    model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
    model.CreateByName = s.CreateByName
}

func (s *ProductCategoryInsertReq) GetId() interface{} {
	return s.Id
}

type ProductCategoryUpdateReq struct {
    Id int `uri:"id" comment:"id"` // id
    SkuCode string `json:"skuCode" comment:"sku"`
    CategoryId int `json:"categoryId" comment:"产线id"`
    MainCateFlag int `json:"mainCateFlag" comment:"主产线标志0否1是"`
    common.ControlBy
}

func (s *ProductCategoryUpdateReq) Generate(model *models.ProductCategory)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.SkuCode = s.SkuCode
    model.CategoryId = s.CategoryId
    model.MainCateFlag = s.MainCateFlag
    model.UpdateBy = s.UpdateBy
    model.UpdateByName = s.UpdateByName // 添加这而，需要记录是被谁更新的
}

func (s *ProductCategoryUpdateReq) GetId() interface{} {
	return s.Id
}

// ProductCategoryGetReq 功能获取请求参数
type ProductCategoryGetReq struct {
     Id int `uri:"id"`
}
func (s *ProductCategoryGetReq) GetId() interface{} {
	return s.Id
}

// ProductCategoryDeleteReq 功能删除请求参数
type ProductCategoryDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *ProductCategoryDeleteReq) GetId() interface{} {
	return s.Ids
}
