package dto

import (
	"time"

	"go-admin/app/oc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type StatementsGetPageReq struct {
	dto.Pagination `search:"-"`
	CompanyId      string `form:"companyId"  search:"type:exact;column:company_id;table:statements"`
	Start          string `form:"startTime" search:"type:gte;column:created_at;table:statements"`
	End            string `form:"endTime" search:"type:lte;column:created_at;table:statements"`
	StatementsOrder
}

type StatementsOrder struct {
	Id        string `form:"idOrder"  search:"type:order;column:id;table:statements"`
	CompanyId string `form:"companyIdOrder"  search:"type:order;column:company_id;table:statements"`
	CreateAt  string `form:"createdAtOrder"  search:"type:order;column:created_at;table:statements"`
}

func (m *StatementsGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type StatementsInsertReq struct {
	Id           int       `json:"-" comment:"自增ID"` // 自增ID
	StatementsNo string    `json:"statementsNo" comment:"对账单编号"`
	CompanyId    int       `json:"companyId" comment:"公司"`
	TotalAmount  float64   `json:"totalAmount" comment:"账单总金额"`
	StartTime    time.Time `json:"startTime" comment:"账单开始日期"`
	EndTime      time.Time `json:"endTime" comment:"账单结束日期"`
	OrderCount   int       `json:"orderCount" comment:"订单总数"`
	ProductCount int       `json:"productCount" comment:"商品总数"`
	CreatedAt    time.Time `json:"createdAt" gorm:"comment:创建时间"`
}

func (s *StatementsInsertReq) Generate(model *models.Statements) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.StatementsNo = s.StatementsNo
	model.CompanyId = s.CompanyId
	model.CreatedAt = s.CreatedAt
	model.TotalAmount = s.TotalAmount
	model.StartTime = s.StartTime
	model.EndTime = s.EndTime
	model.OrderCount = s.OrderCount
	model.ProductCount = s.ProductCount
}

func (s *StatementsInsertReq) GetId() interface{} {
	return s.Id
}

type StatementsUpdateReq struct {
	Id           int       `uri:"id" comment:"自增ID"` // 自增ID
	StatementsNo string    `json:"statementsNo" comment:"对账单编号"`
	CompanyId    int       `json:"companyId" comment:"公司"`
	TotalAmount  float64   `json:"totalAmount" comment:"账单总金额"`
	StartTime    time.Time `json:"startTime" comment:"账单开始日期"`
	EndTime      time.Time `json:"endTime" comment:"账单结束日期"`
	OrderCount   int       `json:"orderCount" comment:"订单总数"`
	ProductCount int       `json:"productCount" comment:"商品总数"`
	CreatedAt    time.Time `json:"createdAt" gorm:"comment:创建时间"`
}

func (s *StatementsUpdateReq) Generate(model *models.Statements) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.StatementsNo = s.StatementsNo
	model.CompanyId = s.CompanyId
	model.CreatedAt = s.CreatedAt
	model.TotalAmount = s.TotalAmount
	model.StartTime = s.StartTime
	model.EndTime = s.EndTime
	model.OrderCount = s.OrderCount
	model.ProductCount = s.ProductCount
}

func (s *StatementsUpdateReq) GetId() interface{} {
	return s.Id
}

// StatementsGetReq 功能获取请求参数
type StatementsGetReq struct {
	Id int `uri:"id"`
}

func (s *StatementsGetReq) GetId() interface{} {
	return s.Id
}

// StatementsDeleteReq 功能删除请求参数
type StatementsDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *StatementsDeleteReq) GetId() interface{} {
	return s.Ids
}

// OrderToStatementsGetReq 功能获取请求参数
type OrderToStatementsGetReq struct {
	dto.Pagination `search:"-"`
	Id             int `uri:"id"`
}

// ExportOrderToStatementsList 导出结构体
type ExportOrderToStatementsList struct {
	Key                   int     `json:"key" gorm:"-"`
	OrderId               string  `json:"orderId"`
	FinalTotalAmount      float64 `json:"finalTotalAmount"`
	SkuCode               string  `json:"skuCode"`
	ProductName           string  `json:"productName"`
	BrandName             string  `json:"brandName"`
	ProductModel          string  `json:"productModel"`
	ProductNo             string  `json:"productNo"`
	VendorName            string  `json:"vendorName"`
	SupplierSkuCode       string  `json:"supplierSkuCode"`
	SalePrice             float64 `json:"salePrice"`
	FinalQuantity         int     `json:"finalQuantity"`
	FinalSubTotalAmount   float64 `json:"finalSubTotalAmount"`
	UserName              string  `json:"userName"`
	UserId                int     `json:"userId"`
	UserPhone             string  `json:"userPhone"`
	FullCompanyName       string  `json:"fullCompanyName"`
	CreateFrom            string  `json:"createFrom"`
	SkuClassificationName string  `json:"skuClassificationName"`
	PayAccount            string  `json:"payAccount"`
}

type OrderToStatementsListResp struct {
	Key                         int     `json:"key" gorm:"-"`
	OrderId                     string  `json:"orderId"`
	FinalTotalAmount            float64 `json:"finalTotalAmount"`
	SkuCode                     string  `json:"skuCode"`
	ProductName                 string  `json:"productName"`
	BrandName                   string  `json:"brandName"`
	ProductModel                string  `json:"productModel"`
	ProductNo                   string  `json:"productNo"`
	VendorName                  string  `json:"vendorName"`
	SupplierSkuCode             string  `json:"supplierSkuCode"`
	SalePrice                   float64 `json:"salePrice"`
	FinalQuantity               int     `json:"finalQuantity"`
	FinalSubTotalAmount         float64 `json:"finalSubTotalAmount"`
	UserName                    string  `json:"userName"`
	UserId                      int     `json:"userId"`
	UserPhone                   string  `json:"userPhone"`
	FullCompanyName             string  `json:"fullCompanyName"`
	CreateFrom                  string  `json:"createFrom"`
	ParentCompanyDepartmentName string  `json:"parentCompanyDepartmentName"`
	CompanyDepartmentName       string  `json:"companyDepartmentName"`
	SkuClassificationSwitch     string  `json:"skuClassificationSwitch"`
	SkuClassificationName       string  `json:"skuClassificationName"`
	PayAccount                  string  `json:"payAccount"`
}

type CompanySwitch struct {
	Keyword      string `json:"keyword"`
	SwitchStatus string `json:"switchStatus"`
}

type CompanySwitchPageBaseResp struct {
	CompanySwitch CompanySwitch
	Count         int64 `json:"count"`
	PageIndex     int   `json:"pageIndex"`
	PageSize      int   `json:"pageSize"`
}

type OrderToStatementsListPageResp struct {
	CompanySwitchPageBaseResp
	List []OrderToStatementsListResp `json:"list"`
}

type StatementsGetPageResp struct {
	models.Statements
	CompanyName string `json:"companyName" gorm:"-"`
}
