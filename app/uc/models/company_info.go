package models

import (
	"errors"
	"go-admin/common/global"
	"go-admin/common/models"

	"github.com/samber/lo"
	"gorm.io/gorm"
)

// CompanyNature 公司性质 （2终端，3分销）
var CompanyNature = map[int64]string{
	2: "终端",
	3: "分销",
}

// CompanyType 公司类型（1:KA 2:SME 3:DS 4:整包 5:央企 6:其他）
var CompanyType = map[int64]string{
	1: "KA",
	2: "SME",
	3: "DS",
	//4: "整包",
	5: "央企",
	7: "地方国企",
	6: "其他",
}

var CompanyTypeGroup = map[int64]map[int64]string{
	//终端
	2: {
		1: "KA",
		2: "SME",
		5: "央企",
		7: "地方国企",
		6: "其他",
	},
	// 分销
	3: {
		2: "SME",
		3: "DS",
	},
}

// ReconciliationDay 对账日
var ReconciliationDay = map[int64]string{
	0:  "未选择",
	1:  "每月1日",
	5:  "每月5日",
	10: "每月10日",
	15: "每月15日",
	20: "每月20日",
	25: "每月25日",
}

var CompanyStatus = map[int64]string{
	0: "停用",
	1: "启用",
}

var CompanyTypeData = struct {
	CompanyStatus     interface{} `json:"companyStatus"`
	BusinessType      interface{} `json:"businessType"`
	CompanyNature     interface{} `json:"companyNature"`
	CompanyTypeGroup  interface{} `json:"companyTypeGroup"`
	SiteCss           interface{} `json:"siteCss"`
	ReconciliationDay interface{} `json:"reconciliationDay"`
}{
	CompanyStatus:     CompanyStatus,
	CompanyNature:     CompanyNature,
	CompanyTypeGroup:  CompanyTypeGroup,
	ReconciliationDay: ReconciliationDay,
}

type CompanyInfo struct {
	models.Model

	CompanyStatus     int    `json:"companyStatus" gorm:"type:tinyint(1);comment:公司状态（1可用 0不可用）"` // 公司状态（1可用 0不可用）
	CompanyName       string `json:"companyName" gorm:"type:varchar(50);comment:公司名称"`
	CompanyNature     int    `json:"companyNature" gorm:"type:tinyint(1);comment:公司性质 （2终端，3分销）"`                     // 公司性质 （2终端，3分销）
	CompanyType       int    `json:"companyType" gorm:"type:tinyint(1);comment:公司类型（1:KA 2:SME 3:DS 4:整包 5:央企 6:其他）"` // （1:KA 2:SME 3:DS 4:整包 5:央企 6:其他）
	IsPunchout        int    `json:"isPunchout" gorm:"type:tinyint(1);comment:punchout服务(1是 0否)"`                     // punchout服务(1是 0否)
	IsEas             int    `json:"isEas" gorm:"type:tinyint(1);comment:eas服务(1是 0否)"`                               // 领用单是否审批(1是 0否)
	CheckStockStatus  int    `json:"checkStockStatus" gorm:"type: tinyint(1);comment:库存不足允许下单：1、是；2、否 默认为1"`
	IsEis             int    `json:"isEis" gorm:"type:tinyint(1);comment:eis服务(1是 0否)"`           // eis服务(1是 0否)
	ParentId          int    `json:"parentId" gorm:"type:int;comment:父级节点"`                       // 父公司id
	Address           string `json:"address" gorm:"type:varchar(200);comment:地址"`                 //地址
	TaxNo             string `json:"taxNo" gorm:"type:varchar(50);comment:税号"`                    //税号
	BankName          string `json:"bankName" gorm:"type:varchar(100);comment:开户行"`               //开户行
	BankAccount       string `json:"bankAccount" gorm:"type:varchar(20);comment:银行帐号"`            //银行帐号
	CompanyLogo       string `json:"companyLogo" gorm:"type:varchar(200);comment:公司logo"`         // 公司logo
	CompanyMediumLogo string `json:"companyMediumLogo" gorm:"type:varchar(200);comment:中等大小logo"` // 中等大小logo
	CompanySmallLogo  string `json:"companySmallLogo" gorm:"type:varchar(200);comment:小logo"`     // 小logo
	IndustryId        int    `json:"industryId" gorm:"type:int;comment:所属行业ID"`
	IsSendEmail       int    `json:"isSendEmail" gorm:"type:tinyint(1);comment:下单时是否邮件通知 0:否  1:是"`
	ShowEasflow       int    `json:"showEasflow" gorm:"type:tinyint(1);comment:订单PDF是否显示EAS审批流信息{1：显示，0：不显示}"`
	Theme             string `json:"theme" gorm:"type:varchar(20);comment:Theme"`
	Domain            string `json:"domain" gorm:"type:varchar(45);comment:公司自定义域名"` // 公司自定义域名
	PaymentMethods    string `json:"paymentMethods" gorm:"type:varchar(45);comment:公司所支持的付款方式：1支付宝 2微信 0线下转账 99账期付款。拼接方式：1|2|0|99。不填默认所有。"`
	SiteTitle         string `json:"siteTitle" gorm:"type:varchar(45);comment:自定义域名后，网站定制标题"` // 自定义域名后，网站定制标题
	SiteCss           string `json:"siteCss" gorm:"type:varchar(45);comment:公司定义的css样式文件名称，如common-sanxia.css"`
	CompanyLevels     int    `json:"companyLevels" gorm:"type:tinyint(1);comment:新建公司等级1-5"` // 公司等级
	LoginType         int    `json:"loginType" gorm:"type:tinyint;comment:登录页类型"`
	Pid               int    `json:"pid" gorm:"type:int;comment:最高级公司"`
	PresoTerm         int    `json:"presoTerm" gorm:"type:int;comment:审批单有效期默认15天"`
	CountryId         int    `json:"countryId" gorm:"type:tinyint;comment:国家，1中国，0国外"`
	ProvinceId        int    `json:"provinceId" gorm:"type:int;comment:省份"`
	CityId            int    `json:"cityId" gorm:"type:int;comment:城市"`
	AreaId            int    `json:"areaId" gorm:"type:int;comment:区域"`
	PdfType           int    `json:"pdfType" gorm:"type:tinyint;comment:eis打印pdf格式"`
	IsTaxPrice        int    `json:"isTaxPrice" gorm:"type:tinyint;comment:未税价(1是0否)"`
	ServiceFee        string `json:"serviceFee" gorm:"type:decimal(10,2);comment:服务费，值为百分号前面的数"`
	OrderEamilSet     string `json:"orderEamilSet" gorm:"type:varchar(255);comment:下单邮箱设置"`
	ApproveEmailType  int    `json:"approveEmailType" gorm:"type:tinyint(1);comment:审批提醒邮件类型：1：中英文；2：中日文；默认：1；"`
	ReconciliationDay int    `json:"reconciliationDay" gorm:"type:smallint;comment:当月对账日"`                          // 当月对账日
	AfterAutoAudit    string `json:"afterAutoAudit" gorm:"type:varchar(100);comment:售后自动审核 空不自动 自动的售后type逗号分隔的字符串"` //售后自动审核 空不自动 自动的售后type逗号分隔的字符串
	OrderAutoConfirm  int    `json:"orderAutoConfirm" gorm:"type:tinyint;comment:订单自动确认 0否 1自动确认"`                  //订单自动确认 0否 1自动确认
	OrderAutoSignFor  int    `json:"orderAutoSignFor" gorm:"type:tinyint;comment:订单自动签收 0否 1是"`                     //订单自动签收 0否 1是
	models.ModelTime
	models.ControlBy
}

func (CompanyInfo) TableName() string {
	return "company_info"
}

func (e *CompanyInfo) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *CompanyInfo) GetId() interface{} {
	return e.Id
}

// GetPCompanyDomainByID 获取要跳转的域名 （如果有父公司以夫公司的为准）
func (e *CompanyInfo) GetPCompanyDomainByID(tx *gorm.DB, companyId int) (domain string, err error) {
	if companyId == 0 {
		err = errors.New("公司id不能为空")
		return
	}

	companyInfo := CompanyInfo{}
	// 执行查询
	prefix := global.GetTenantUcDBNameWithDB(tx)
	err = tx.Table(prefix+".company_info").First(&companyInfo, companyId).Error
	if err != nil {
		return
	}
	pId := companyInfo.Pid

	if pId != 0 {
		companyInfo = CompanyInfo{}
		err = tx.Table(prefix+".company_info").First(&companyInfo, pId).Error
		if err != nil {
			return
		}
	}
	domain = companyInfo.Domain

	return
}

func (e *CompanyInfo) GetRowsByCondition(tx *gorm.DB, query interface{}, args ...interface{}) (data []CompanyInfo) {
	ucPrefix := global.GetTenantUcDBNameWithDB(tx)
	tx.Table(ucPrefix+"."+e.TableName()).Where(query, args).Find(&data)

	return
}

// 根据用户id获取用户所属公司
func (e *CompanyInfo) GetRowByUserId(tx *gorm.DB, userId int) (data CompanyInfo) {
	ucPrefix := global.GetTenantUcDBNameWithDB(tx)
	tx.Table(ucPrefix+"."+e.TableName()).Joins("left join "+ucPrefix+".user_info ui on "+e.TableName()+".id = ui.company_id").Where("ui.id = ?", userId).Find(&data)

	return
}

// 根据公司IDs查询公司列表
func (e *CompanyInfo) ListByIds(tx *gorm.DB, ids []int) ([]*CompanyInfo, error) {
	list := []*CompanyInfo{}
	err := tx.Scopes(global.TenantTable("uc", e.TableName())).Where("id in (?)", ids).Find(&list).Error
	if err != nil {
		return nil, err
	}

	return list, nil
}

// 根据公司IDs查询公司Map
func (e *CompanyInfo) MapByIds(tx *gorm.DB, ids []int) (map[int]*CompanyInfo, error) {
	// 查列表
	list, err := e.ListByIds(tx, ids)
	if err != nil {
		return nil, err
	}

	// 组map
	res := lo.Associate(list, func(item *CompanyInfo) (int, *CompanyInfo) {
		return item.Id, item
	})
	return res, nil
}
