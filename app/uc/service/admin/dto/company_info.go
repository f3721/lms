package dto

import (
	"go-admin/app/uc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type CompanyInfoGetPageReq struct {
	dto.Pagination       `search:"-"`
	CompanyStatus        *int   `form:"companyStatus"  search:"type:exact;column:company_status;table:company_info" comment:"公司状态（1可用 0不可用）"`
	CompanyName          string `form:"companyName"  search:"type:contains;column:company_name;table:company_info" comment:"公司名称"`
	CompanyNature        int    `form:"companyNature"  search:"type:exact;column:company_nature;table:company_info" comment:"公司性质 （2终端，3分销）"`
	CompanyType          int    `form:"companyType"  search:"type:exact;column:company_type;table:company_info" comment:"公司类型（1:KA 2:SME 3:DS 4:整包 5:央企 6:其他）"`
	ParentId             int    `form:"parentId"  search:"type:exact;column:parent_id;table:company_info" comment:"父级节点"`
	CompanyLevels        int    `form:"companyLevels"  search:"type:exact;column:company_levels;table:company_info" comment:"新建公司等级1-5"`
	Pid                  int    `form:"pid"  search:"type:exact;column:pid;table:company_info" comment:"最高级公司"`
	QueryCompanyIds      string `form:"queryCompanyIds" search:"-"`
	QueryCompanyNames    string `form:"queryCompanyNames" search:"-"` // 公司名批量查询
	IgnoreUserPermission bool   `json:"-" search:"-"`
	CompanyInfoOrder
}

type CompanyInfoOrder struct {
	Id                string `form:"idOrder"  search:"type:order;column:id;table:company_info"`
	CompanyStatus     string `form:"companyStatusOrder"  search:"type:order;column:company_status;table:company_info"`
	CompanyName       string `form:"companyNameOrder"  search:"type:order;column:company_name;table:company_info"`
	CompanyNature     string `form:"companyNatureOrder"  search:"type:order;column:company_nature;table:company_info"`
	CompanyType       string `form:"companyTypeOrder"  search:"type:order;column:company_type;table:company_info"`
	IsPunchout        string `form:"isPunchoutOrder"  search:"type:order;column:is_punchout;table:company_info"`
	IsEas             string `form:"isEasOrder"  search:"type:order;column:is_eas;table:company_info"`
	IsEis             string `form:"isEisOrder"  search:"type:order;column:is_eis;table:company_info"`
	ParentId          string `form:"parentIdOrder"  search:"type:order;column:parent_id;table:company_info"`
	Address           string `form:"addressOrder"  search:"type:order;column:address;table:company_info"`
	TaxNo             string `form:"taxNoOrder"  search:"type:order;column:tax_no;table:company_info"`
	BankName          string `form:"bankNameOrder"  search:"type:order;column:bank_name;table:company_info"`
	BankAccount       string `form:"bankAccountOrder"  search:"type:order;column:bank_account;table:company_info"`
	CompanyLogo       string `form:"companyLogoOrder"  search:"type:order;column:company_logo;table:company_info"`
	CompanyMediumLogo string `form:"companyMediumLogoOrder"  search:"type:order;column:company_medium_logo;table:company_info"`
	CompanySmallLogo  string `form:"companySmallLogoOrder"  search:"type:order;column:company_small_logo;table:company_info"`
	IndustryId        string `form:"industryIdOrder"  search:"type:order;column:industry_id;table:company_info"`
	IsSendEmail       string `form:"isSendEmailOrder"  search:"type:order;column:is_send_email;table:company_info"`
	ShowEasflow       string `form:"showEasflowOrder"  search:"type:order;column:show_easflow;table:company_info"`
	Theme             string `form:"themeOrder"  search:"type:order;column:theme;table:company_info"`
	Domain            string `form:"domainOrder"  search:"type:order;column:domain;table:company_info"`
	PaymentMethods    string `form:"paymentMethodsOrder"  search:"type:order;column:payment_methods;table:company_info"`
	SiteTitle         string `form:"siteTitleOrder"  search:"type:order;column:site_title;table:company_info"`
	SiteCss           string `form:"siteCssOrder"  search:"type:order;column:site_css;table:company_info"`
	CompanyLevels     string `form:"companyLevelsOrder"  search:"type:order;column:company_levels;table:company_info"`
	LoginType         string `form:"loginTypeOrder"  search:"type:order;column:login_type;table:company_info"`
	Pid               string `form:"pidOrder"  search:"type:order;column:pid;table:company_info"`
	PresoTerm         string `form:"presoTermOrder"  search:"type:order;column:preso_term;table:company_info"`
	CountryId         string `form:"countryIdOrder"  search:"type:order;column:country_id;table:company_info"`
	ProvinceId        string `form:"provinceIdOrder"  search:"type:order;column:province_id;table:company_info"`
	CityId            string `form:"cityIdOrder"  search:"type:order;column:city_id;table:company_info"`
	AreaId            string `form:"areaIdOrder"  search:"type:order;column:area_id;table:company_info"`
	PdfType           string `form:"pdfTypeOrder"  search:"type:order;column:pdf_type;table:company_info"`
	IsTaxPrice        string `form:"isTaxPriceOrder"  search:"type:order;column:is_tax_price;table:company_info"`
	ServiceFee        string `form:"serviceFeeOrder"  search:"type:order;column:service_fee;table:company_info"`
	OrderEamilSet     string `form:"orderEamilSetOrder"  search:"type:order;column:order_eamil_set;table:company_info"`
	ApproveEmailType  string `form:"approveEmailTypeOrder"  search:"type:order;column:approve_email_type;table:company_info"`
	ReconciliationDay string `form:"reconciliationDayOrder"  search:"type:order;column:reconciliation_day;table:company_info"`
	AfterAutoAudit    string `form:"afterAutoAuditOrder"  search:"type:order;column:after_auto_audit;table:company_info"`
	OrderAutoConfirm  string `form:"orderAutoConfirmOrder"  search:"type:order;column:order_auto_confirm;table:company_info"`
	CreatedAt         string `form:"createdAtOrder"  search:"type:order;column:created_at;table:company_info"`
	UpdatedAt         string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:company_info"`
	DeletedAt         string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:company_info"`
	CreateBy          string `form:"createByOrder"  search:"type:order;column:create_by;table:company_info"`
	CreateByName      string `form:"createByNameOrder"  search:"type:order;column:create_by_name;table:company_info"`
	UpdateBy          string `form:"updateByOrder"  search:"type:order;column:update_by;table:company_info"`
	UpdateByName      string `form:"updateByNameOrder"  search:"type:order;column:update_by_name;table:company_info"`
}

func (m *CompanyInfoGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type CompanyInfoGetPageRes struct {
	List []*CompanyInfoGetPageData `json:"list"`
}

type CompanyInfoGetPageData struct {
	Id               int    `json:"id" comment:"主键"`                                                     // 主键
	CompanyStatus    int    `json:"companyStatus" comment:"公司状态（1可用 0不可用）"`                              // 公司状态（1可用 0不可用）
	CompanyName      string `json:"companyName" comment:"公司名称"`                                          // 公司名称
	CompanyNature    int    `json:"companyNature" comment:"公司性质 （2终端，3分销）"`                              // 公司性质 （2终端，3分销）
	CompanyType      int    `json:"companyType" comment:"公司类型（1:KA 2:SME 3:DS 4:整包 5:央企 6:其他）" vd:"$>0"` // 公司类型（1:KA 2:SME 3:DS 4:整包 5:央企 6:其他）
	CompanyLevels    int    `json:"companyLevels" comment:"新建公司等级1-5"`                                   //公司层级
	OrderAutoSignFor int    `json:"orderAutoSignFor"`                                                    //订单自动签收 0否 1是
	common.ModelTime
	common.ControlBy
	Children []*CompanyInfoGetPageData `json:"children"`
}

type CompanyInfoGetSelectPageData struct {
	Id                int    `json:"id" comment:"主键"`                                                     // 主键
	CompanyStatus     int    `json:"companyStatus" comment:"公司状态（1可用 0不可用）"`                              // 公司状态（1可用 0不可用）
	CompanyName       string `json:"companyName" comment:"公司名称"`                                          // 公司名称
	CompanyNature     int    `json:"companyNature" comment:"公司性质 （2终端，3分销）"`                              // 公司性质 （2终端，3分销）
	CompanyType       int    `json:"companyType" comment:"公司类型（1:KA 2:SME 3:DS 4:整包 5:央企 6:其他）" vd:"$>0"` // 公司类型（1:KA 2:SME 3:DS 4:整包 5:央企 6:其他）
	CompanyLevels     int    `json:"companyLevels" comment:"新建公司等级1-5"`                                   //公司层级
	ReconciliationDay int    `json:"reconciliationDay" comment:"当月对账日"`                                   // 当月对账日
	OrderAutoSignFor  int    `json:"orderAutoSignFor"`

	common.ModelTime
	common.ControlBy
}

type CompanyInfoInsertReq struct {
	Id                int    `json:"-" comment:"主键"` // 主键
	CompanyStatus     int    `json:"companyStatus" comment:"公司状态（1可用 0不可用）"`
	CompanyName       string `json:"companyName" comment:"公司名称" vd:"len($)>1,len($)<=100" `
	CompanyNature     int    `json:"companyNature" comment:"公司性质 （2终端，3分销）" vd:"$>0" `
	CompanyType       int    `json:"companyType" comment:"公司类型（1:KA 2:SME 3:DS 4:整包 5:央企 6:其他）" vd:"$>0" `
	Address           string `json:"address"`                                                                          // 地址
	IsEas             int    `json:"isEas" comment:"eas服务(1是 0否)"`                                                     // 领用单是否审批
	CheckStockStatus  int    `json:"checkStockStatus" comment:"库存不足允许下单：1、是；2、否" vd:"@:in($,0,1); msg:'库存不足允许下单状态错误'"` // 领用单是否审批
	CompanyLogo       string `json:"companyLogo" comment:"公司logo"`
	CompanyMediumLogo string `json:"companyMediumLogo" comment:"中等大小logo"`
	CompanySmallLogo  string `json:"companySmallLogo" comment:"小logo"`
	Domain            string `json:"-" comment:"公司自定义域名"`
	SiteTitle         string `json:"-" comment:"自定义域名后，网站定制标题"`
	IsPunchout        int    `json:"-" comment:"punchout服务(1是 0否)"`
	ParentId          int    `json:"-" comment:"父级节点"`
	SiteCss           string `json:"-" comment:"公司定义的css样式文件名称，如common-sanxia.css"`
	CompanyLevels     int    `json:"-" comment:"新建公司等级1-5"`
	LoginType         int    `json:"-" comment:"登录页类型"`
	ReconciliationDay int    `json:"-" comment:"当月对账日"`
	AfterAutoAudit    string `json:"-" comment:"售后自动审核 空不自动 自动的售后type逗号分隔的字符串"`
	OrderAutoConfirm  int    `json:"-" comment:"订单自动确认 0否 1自动确认"`
	OrderAutoSignFor  int    `json:"orderAutoSignFor"` //订单自动签收 0否 1是
	common.ControlBy  `json:"-"`
}

func (s *CompanyInfoInsertReq) Generate(model *models.CompanyInfo) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.CompanyStatus = s.CompanyStatus
	model.CompanyName = s.CompanyName
	model.CompanyNature = s.CompanyNature
	model.CompanyType = s.CompanyType
	model.IsPunchout = s.IsPunchout
	model.IsEas = s.IsEas
	model.CheckStockStatus = s.CheckStockStatus
	model.ParentId = s.ParentId
	model.Address = s.Address
	//model.TaxNo = s.TaxNo
	//model.BankName = s.BankName
	//model.BankAccount = s.BankAccount
	//model.CompanyLogo = s.CompanyLogo
	//model.CompanyMediumLogo = s.CompanyMediumLogo
	//model.CompanySmallLogo = s.CompanySmallLogo
	//model.IndustryId = s.IndustryId
	//model.IsSendEmail = s.IsSendEmail
	//model.ShowEasflow = s.ShowEasflow
	//model.Theme = s.Theme
	model.Domain = s.Domain
	//model.PaymentMethods = s.PaymentMethods
	model.SiteTitle = s.SiteTitle
	model.SiteCss = s.SiteCss
	model.CompanyLevels = s.CompanyLevels
	model.LoginType = s.LoginType
	//model.Pid = s.Pid
	//model.PresoTerm = s.PresoTerm
	//model.CountryId = s.CountryId
	//model.ProvinceId = s.ProvinceId
	//model.CityId = s.CityId
	//model.AreaId = s.AreaId
	//model.PdfType = s.PdfType
	//model.IsTaxPrice = s.IsTaxPrice
	//model.ServiceFee = s.ServiceFee
	//model.ApproveEmailType = s.ApproveEmailType
	model.ReconciliationDay = s.ReconciliationDay
	model.AfterAutoAudit = s.AfterAutoAudit
	model.OrderAutoConfirm = s.OrderAutoConfirm
	model.OrderAutoSignFor = s.OrderAutoSignFor

	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.CreateByName = s.CreateByName
}

func (s *CompanyInfoInsertReq) GetId() interface{} {
	return s.Id
}

type CompanyInfoUpdateReq struct {
	Id                int    `uri:"id" comment:"主键" vd:"$>0" `                                                         // 主键
	CompanyName       string `json:"companyName" comment:"公司名称" vd:"len($)>0" `                                        // 公司名称
	CompanyNature     int    `json:"companyNature" comment:"公司性质 （2终端，3分销）" vd:"$>0" `                                 //公司性质
	CompanyType       int    `json:"companyType" comment:"公司类型（1:KA 2:SME 3:DS 4:整包 5:央企 6:其他）" vd:"$>0" `             //公司类型
	Address           string `json:"address" comment:"地址"`                                                             //地址
	IsEas             int    `json:"isEas" comment:"eas服务(1是 0否)"`                                                     // 领用单是否审批(1是 0否)
	CheckStockStatus  int    `json:"checkStockStatus" comment:"库存不足允许下单：1、是；2、否" vd:"@:in($,0,1); msg:'库存不足允许下单状态错误'"` // 领用单是否审批
	ReconciliationDay int    `json:"reconciliationDay" comment:"当月对账日"`                                                //当月对账日
	CompanyStatus     int    `json:"companyStatus" comment:"公司状态（1可用 0不可用）"`                                           //公司状态
	CompanyLogo       string `json:"companyLogo" comment:"公司logo"`                                                     //公司logo
	CompanyMediumLogo string `json:"companyMediumLogo" comment:"中等大小logo"`                                             //中等大小logo
	CompanySmallLogo  string `json:"companySmallLogo" comment:"小logo"`                                                 //小logo
	SiteTitle         string `json:"siteTitle" comment:"自定义域名后，网站定制标题"`                                                // 自定义域名后，网站定制标题
	Domain            string `json:"domain" comment:"公司自定义域名"`                                                         // 公司自定义域名

	OrderAutoConfirm int `json:"orderAutoConfirm" comment:"订单自动确认 0否 1自动确认"`
	OrderAutoSignFor int `json:"orderAutoSignFor"` //订单自动签收 0否 1是
	common.ControlBy `json:"-"`
}

func (s *CompanyInfoUpdateReq) Generate(model *models.CompanyInfo) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.CompanyStatus = s.CompanyStatus
	model.CompanyName = s.CompanyName
	model.CompanyNature = s.CompanyNature
	model.CompanyType = s.CompanyType
	model.IsEas = s.IsEas
	model.CheckStockStatus = s.CheckStockStatus
	model.Address = s.Address
	model.CompanyLogo = s.CompanyLogo
	model.CompanyMediumLogo = s.CompanyMediumLogo
	model.CompanySmallLogo = s.CompanySmallLogo
	model.Domain = s.Domain
	model.SiteTitle = s.SiteTitle
	model.ReconciliationDay = s.ReconciliationDay
	model.OrderAutoConfirm = s.OrderAutoConfirm
	model.OrderAutoSignFor = s.OrderAutoSignFor
	model.UpdateBy = s.UpdateBy
	model.UpdateByName = s.UpdateByName // 添加这而，需要记录是被谁更新的
}

func (s *CompanyInfoUpdateReq) GetId() interface{} {
	return s.Id
}

// CompanyInfoGetReq 功能获取请求参数
type CompanyInfoGetReq struct {
	Id int `uri:"id"`
}

func (s *CompanyInfoGetReq) GetId() interface{} {
	return s.Id
}

// CompanyInfoGetRes 功能获取请求参数
type CompanyInfoGetRes struct {
	models.CompanyInfo
	ParentName string `json:"parentName"`
}

// CompanyInfoGetReq 功能获取请求参数
type CompanyInfoGetByNameReq struct {
	CompanyName string `form:"companyName" vd:"len($)>0" `
}

// CompanyInfoDeleteReq 功能删除请求参数
type CompanyInfoDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *CompanyInfoDeleteReq) GetId() interface{} {
	return s.Ids
}

// CompanyInfoGetIdReq 功能获取请求参数
type CompanyInfoGetIdReq struct {
	Id int `uri:"id"`
}

type CompanyInfoIsAvailableRes struct {
	IsAvailable bool `json:"isAvailable"` //是否可用
}
