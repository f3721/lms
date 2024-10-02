package dto

import (
	"go-admin/app/wc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type QualityCheckConfigGetPageReq struct {
	dto.Pagination      `search:"-"`
	Id                  string `form:"idOrder"  search:"type:order;column:id;table:quality_check_config"`
	Status              string `form:"statusOrder"  search:"type:order;column:status;table:quality_check_config"`
	CompanyIds          string `form:"companyIds"  search:"type:order;column:company_ids;table:quality_check_config"`
	WarehouseCodes      string `form:"warehouseCodes"  search:"type:order;column:warehouse_codes;table:quality_check_config"`
	Type                string `form:"type"  search:"type:exact;column:type;table:quality_check_config"`
	SamplingNum         string `form:"samplingNumOrder"  search:"type:order;column:sampling_num;table:quality_check_config"`
	OrderType           string `form:"orderTypeOrder"  search:"type:order;column:order_type;table:quality_check_config"`
	SkuConstraint       string `form:"skuConstraintOrder"  search:"type:order;column:sku_constraint;table:quality_check_config"`
	CategoryConstraint  string `form:"categoryConstraintOrder"  search:"type:order;column:category_constraint;table:quality_check_config"`
	QualityCheckRoles   string `form:"qualityCheckRolesOrder"  search:"type:order;column:quality_check_roles;table:quality_check_config"`
	QualityCheckOptions string `form:"qualityCheckOptionsOrder"  search:"type:order;column:quality_check_options;table:quality_check_config"`
	CreatedAt           string `form:"createdAtOrder"  search:"type:order;column:created_at;table:quality_check_config"`
	UpdatedAt           string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:quality_check_config"`
	DeletedAt           string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:quality_check_config"`
	CreateBy            string `form:"createByOrder"  search:"type:order;column:create_by;table:quality_check_config"`
	UpdateBy            string `form:"updateByOrder"  search:"type:order;column:update_by;table:quality_check_config"`
	CreateByName        string `form:"createByNameOrder"  search:"type:order;column:create_by_name;table:quality_check_config"`
	UpdateByName        string `form:"updateByNameOrder"  search:"type:order;column:update_by_name;table:quality_check_config"`
}

type QualityCheckConfigGetPageResp struct {
	models.QualityCheckConfig
	CompanyNames   string `json:"companyNames" comment:"适用公司"`
	WarehouseNames string `json:"warehouseNames" comment:"适用仓库"`
	TypeName       string `json:"typeName" comment:"质检类型中文"`
}

func (m *QualityCheckConfigGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type QualityCheckConfigInsertReq struct {
	Id                       int                                `uri:"id" comment:"Primary Key"` // Primary Key
	Status                   int                                `json:"status" comment:"质检开关 0-关 1-开" vd:"@:in($,0,1); msg:'质检开关范围[0,1]'"`
	CompanyIds               []int                              `json:"companyIds" comment:"适用公司IDS" vd:"@:len($)>0; msg:'适用公司必填'"`
	WarehouseCodes           []string                           `json:"warehouseCodes" comment:"适用仓库IDs" vd:"@:len($)>0; msg:'适用仓库必填'"`
	Type                     int                                `json:"type" comment:"质检类型 0-全检  1-抽检" vd:"@:in($,0,1); msg:'质检类型范围[0,1]'"`
	SamplingNum              int                                `json:"samplingNum" comment:"抽检数量"`
	OrderType                int                                `json:"orderType" comment:"质检订单类型:0-全部 1-采购入库 2-大货" vd:"@:in($,0,1,2); msg:'质检订单类型范围[0,1,2]'"`
	SkuConstraint            []*models.Constraint               `json:"skuConstraint" comment:"sku约束"`
	CategoryConstraint       []*models.Constraint               `json:"categoryConstraint" comment:"产线约束"`
	QualityCheckRoles        string                             `json:"qualityCheckRoles" comment:"质检角色IDs" vd:"@:len($)>0; msg:'质检角色必填'"`
	QualityCheckOptions      string                             `json:"qualityCheckOptions" comment:"质检内容" vd:"@:len($)>0; msg:'质检内容必填'"`
	Unqualified              int                                `json:"unqualified" comment:"不合格处理办法" vd:"@:in($,0,1); msg:'不合格处理办法范围[0,1]'"`
	QualityCheckConfigDetail []*models.QualityCheckConfigDetail `json:"-"  comment:"子表数据"`
	common.ControlBy
}

func (s *QualityCheckConfigInsertReq) Generate(model *models.QualityCheckConfig) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.Status = s.Status
	model.Type = s.Type
	model.SamplingNum = s.SamplingNum
	model.OrderType = s.OrderType
	model.SkuConstraint = s.SkuConstraint
	model.CategoryConstraint = s.CategoryConstraint
	model.QualityCheckRoles = s.QualityCheckRoles
	model.QualityCheckOptions = s.QualityCheckOptions
	model.Unqualified = s.Unqualified
	model.QualityCheckConfigDetail = s.QualityCheckConfigDetail
	model.CreateBy = s.CreateBy
	model.UpdateBy = s.UpdateBy
	model.CreateByName = s.ControlBy.CreateByName
	model.UpdateByName = s.ControlBy.UpdateByName
}

func (s *QualityCheckConfigInsertReq) GetId() interface{} {
	return s.Id
}

// QualityCheckConfigGetReq 功能获取请求参数
type QualityCheckConfigGetReq struct {
	Id int `uri:"id"`
}

func (s *QualityCheckConfigGetReq) GetId() interface{} {
	return s.Id
}

// QualityCheckConfigDeleteReq 功能删除请求参数
type QualityCheckConfigDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *QualityCheckConfigDeleteReq) GetId() interface{} {
	return s.Ids
}

// 查询质检任务请求参数
type QualityCheckNumReq struct {
	WarehouseCode           string                     `json:"warehouseCode"`
	OrderType               string                     `json:"orderType"`
	QualityCheckNumProducts []*QualityCheckNumProducts `json:"qualityCheckNumProducts"`
}

type QualityCheckNumProducts struct {
	SkuCode     string `json:"skuCode" comment:"sku"`
	Quantity    int    `json:"quantity" comment:"调拨数量"`
	Type        int    `json:"type" comment:"质检类型 0-全检  1-抽检"`
	Unqualified int    `json:"unqualified" comment:"不合格处理办法: 0-拒收 1-异常填报"`
}
