package models

import (
	"go-admin/common/models"
	"time"
)

type Preso struct {
	models.Model

	PresoNo         string        `json:"presoNo" gorm:"type:varchar(30);comment:审批单编号"`
	WarehouseCode   string        `json:"warehouseCode" gorm:"type:varchar(10);comment:发货仓"`
	ApproveflowId   int           `json:"approveflowId" gorm:"type:int unsigned;comment:审批流id"`
	ApproveUsers    string        `json:"approveUsers" gorm:"type:varchar(200);comment:审批流用户"`
	UserId          int           `json:"userId" gorm:"type:int unsigned;comment:客户编号"`
	UserName        string        `json:"userName" gorm:"type:varchar(30);comment:客户用户名"`
	UserCompanyId   int           `json:"userCompanyId" gorm:"type:int unsigned;comment:客户公司ID"`
	UserCompanyName string        `json:"userCompanyName" gorm:"type:varchar(50);comment:客户公司名称"`
	DeliverId       int           `json:"deliverId" gorm:"type:int unsigned;comment:收货地址ID"`
	Consignee       string        `json:"consignee" gorm:"type:varchar(20);comment:收货人姓名"`
	CountryId       int           `json:"countryId" gorm:"type:int unsigned;comment:客户国家ID"`
	CountryName     string        `json:"countryName" gorm:"type:varchar(50);comment:客户国家名称"`
	ProvinceId      int           `json:"provinceId" gorm:"type:int unsigned;comment:客户省份ID"`
	ProvinceName    string        `json:"provinceName" gorm:"type:varchar(50);comment:客户省份名称"`
	CityId          int           `json:"cityId" gorm:"type:int unsigned;comment:收货人城市编号"`
	CityName        string        `json:"cityName" gorm:"type:varchar(20);comment:收货人城市名称"`
	AreaId          int           `json:"areaId" gorm:"type:int unsigned;comment:客户区县ID"`
	AreaName        string        `json:"areaName" gorm:"type:varchar(50);comment:客户区县名称"`
	TownId          int           `json:"townId" gorm:"type:int unsigned;comment:镇/街道 ID"`
	TownName        string        `json:"townName" gorm:"type:varchar(50);comment:镇/街道名称"`
	CompanyName     string        `json:"companyName" gorm:"type:varchar(50);comment:收货人公司名称"`
	Address         string        `json:"address" gorm:"type:varchar(500);comment:收货人详细地址"`
	Mobile          string        `json:"mobile" gorm:"type:varchar(20);comment:收货人手机号"`
	Telephone       string        `json:"telephone" gorm:"type:varchar(20);comment:收货人座机号"`
	ContactEmail    string        `json:"contactEmail" gorm:"type:varchar(50);comment:联系人邮箱"`
	ContractNo      string        `json:"contractNo" gorm:"type:varchar(100);comment:客户合同编号"`
	CreateFrom      string        `json:"createFrom" gorm:"type:varchar(30);comment:订单来源：LMS/MALL/XCX"`
	Remark          string        `json:"remark" gorm:"type:varchar(500);comment:客户留言"`
	Ip              string        `json:"ip" gorm:"type:varchar(20);comment:客户IP地址"`
	ApproveStatus   int           `json:"approveStatus" gorm:"type:tinyint;comment:审批状态 -1 审批不通过  1 审批通过 0 初始提交审批 10 审批中 -2 超时 -3 撤回"`
	ApproveRemark   string        `json:"approveRemark" gorm:"type:varchar(400);comment:审批备注"`
	Step            int           `json:"step" gorm:"type:tinyint unsigned;comment:审批到第几级"`
	ExpireTime      time.Time     `json:"expireTime" gorm:"type:datetime;comment:过期时间"`
	CreatedAt       time.Time     `json:"createdAt" gorm:"comment:创建时间"`
	UpdatedAt       time.Time     `json:"updatedAt" gorm:"comment:最后更新时间"`
	CreateBy        int           `json:"createBy" gorm:"index;comment:创建者"`
	CreateByName    string        `json:"createByName" gorm:"index;comment:创建者姓名"`
	PresoDetails    []PresoDetail `json:"presoDetails" gorm:"foreignkey:PresoNo;references:PresoNo"`
	Files           []PresoImage  `json:"files" gorm:"foreignkey:PresoNo;references:PresoNo"`
}

func (Preso) TableName() string {
	return "preso"
}

func (e *Preso) GetId() interface{} {
	return e.Id
}

// 审批邮件相关模板
type ApprovalEmail struct {
	UserName     string        `json:"userName" comment:"审批人名称"`
	PresoNo      string        `json:"presoNo" comment:"预订单号"`
	ApproveUrl   string        `json:"approveUrl" comment:"审批URL地址"`
	ProductsTab  string        `json:"productsTab" comment:""`
	LogoUrl      string        `json:"logoUrl" comment:"logo地址"`
	ContractNo   string        `json:"contractNo" comment:"客户PO号"`
	PresoDetails []PresoDetail `json:"presoDetails" comment:"产品信息"`
}

// 定时任务审批邮件模板
type CornApprovalEmail struct {
	UserName   string `json:"userName" comment:"审批人名称"`
	LogoUrl    string `json:"logoUrl" comment:"logo地址"`
	ApproveUrl string `json:"approveUrl" comment:"待审批URL地址"`
	Num        int    `json:"num" comment:"待审批条数"`
}
