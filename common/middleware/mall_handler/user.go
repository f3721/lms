package mall_handler

import (
	"errors"
	"fmt"
	"go-admin/common/global"
	"go-admin/common/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
)

func GetUserCompanyID(c *gin.Context) (int, error) {
	db, err := pkg.GetOrm(c)
	if err != nil {
		return 0, fmt.Errorf("获取数据库连接失败:%v", err)
	}
	userId := user.GetUserId(c)
	ucDbName := global.GetTenantUcDBNameWithContext(c)
	var companyId int
	db.Raw("SELECT company_id FROM "+ucDbName+".user_info WHERE id = ?", userId).Scan(&companyId)
	return companyId, nil
}

func GetUserId(c *gin.Context) int {
	return user.GetUserId(c)
}

func GetUserInfo(c *gin.Context) (*UserInfo, error) {
	db, err := pkg.GetOrm(c)
	if err != nil {
		return nil, fmt.Errorf("获取数据库连接失败:%v", err)
	}
	userId := user.GetUserId(c)
	ucDbName := global.GetTenantUcDBNameWithContext(c)
	var userInfo UserInfo
	err = db.Table(ucDbName+".user_info").Where("id=?", userId).First(&userInfo).Error
	if err != nil {
		return nil, err
	}
	return &userInfo, nil
}

func GetUserCompanyInfo(c *gin.Context) (*CompanyInfo, error) {
	db, err := pkg.GetOrm(c)
	if err != nil {
		return nil, fmt.Errorf("获取数据库连接失败:%v", err)
	}
	userId := user.GetUserId(c)
	ucDbName := global.GetTenantUcDBNameWithContext(c)
	var companyId int
	db.Raw("SELECT company_id FROM "+ucDbName+".user_info WHERE id = ?", userId).Scan(&companyId)
	if companyId == 0 {
		return nil, errors.New("未查询到对应用户的公司")
	}
	var companyInfo CompanyInfo
	err = db.Table(ucDbName+".company_info").Where("id=?", companyId).First(&companyInfo).Error
	if err != nil {
		return nil, err
	}
	return &companyInfo, nil
}

type UserInfo struct {
	models.Model

	UserEmail            string    `json:"userEmail" gorm:"type:varchar(50);comment:用户邮箱"`
	LoginName            string    `json:"loginName" gorm:"type:varchar(100);comment:用户登录名"`
	UserPhone            string    `json:"userPhone" gorm:"type:varchar(16);comment:手机号码"`
	LoginPassword        string    `json:"loginPassword" gorm:"type:varchar(100);comment:登录密码"`
	UserName             string    `json:"userName" gorm:"type:varchar(50);comment:用户名称"`
	Gender               string    `json:"gender" gorm:"type:varchar(2);comment:性别 (0:女 1:男)"`
	HeadPortrait         string    `json:"headPortrait" gorm:"type:varchar(100);comment:头像地址"`
	RegisterMode         int       `json:"registerMode" gorm:"type:tinyint unsigned;comment:注册方式 (1 手机注册，2邮箱注册,3 微信注册)"`
	PhoneType            int       `json:"phoneType" gorm:"type:tinyint unsigned;comment:手机类型(1 IOS 2android)"`
	DeviceToken          string    `json:"deviceToken" gorm:"type:varchar(64);comment:手机标识码"`
	UserSource           int       `json:"userSource" gorm:"type:tinyint unsigned;comment:用户来源（1www,2opc,3android,4IOS,5M,6微信,7微信小程序，8脉信，9CRM）"`
	UserStatus           int       `json:"userStatus" gorm:"type:tinyint(1);comment:用户状态（1可用，0不可用）"`
	CompanyId            int       `json:"companyId" gorm:"type:int unsigned;comment:公司ID"`
	EmailApproveCronExpr string    `json:"emailApproveCronExpr" gorm:"type:varchar(200);comment:邮件定时规则表达式"`
	IsOpen               int       `json:"isOpen" gorm:"type:int;comment:此用户的邮件审批是否开启"`
	UserType             string    `json:"userType" gorm:"type:varchar(2);comment:用户类型 1.下单人2联系人3审批人4通用"`
	Position             string    `json:"position" gorm:"type:varchar(100);comment:Position"`
	RegisterSource       string    `json:"registerSource" gorm:"type:varchar(100);comment:注册用户渠道来源来源"`
	DeliveryMailset      string    `json:"deliveryMailset" gorm:"type:varchar(10);comment:发货邮件通知设置，英文逗号分割{1,2,3=下单人,收货人,销售}"`
	InvoiceMailset       string    `json:"invoiceMailset" gorm:"type:varchar(10);comment:发票邮件通知设置，英文逗号分割{1,2,3,4=下单人,收货人,销售,收票人}"`
	EasReceiveMsgSet     string    `json:"easReceiveMsgSet" gorm:"type:varchar(10);comment:EAS审批短信通知设置，英文逗号分割{1=采购人（审批驳回）;2=审批人（采购审批）}"`
	IsAdminShow          string    `json:"isAdminShow" gorm:"type:varchar(2);comment:用户状态（1是，0否）"`
	Telephone            string    `json:"telephone" gorm:"type:varchar(50);comment:用户座机号"`
	EmailVerifyStatus    int       `json:"emailVerifyStatus" gorm:"type:tinyint(1);comment:邮箱验证状态（0：未验证 1：已验证 2验证失败）"`
	CanLogin             int       `json:"canLogin" gorm:"type:tinyint(1);comment:是否可登陆"`
	CompanyDepartmentId  int       `json:"companyDepartmentId" gorm:"type:int unsigned;comment:用户所属部门"`
	CreatedAt            time.Time `json:"createdAt" gorm:"comment:创建时间"`
	UpdatedAt            time.Time `json:"updatedAt" gorm:"comment:最后更新时间"`
	models.ControlBy
}

func (UserInfo) TableName() string {
	return "user_info"
}

func (e *UserInfo) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *UserInfo) GetId() interface{} {
	return e.Id
}

type CompanyInfo struct {
	models.Model

	CompanyStatus     int    `json:"companyStatus" gorm:"type:tinyint(1);comment:公司状态（1可用 0不可用）"` // 公司状态（1可用 0不可用）
	CompanyName       string `json:"companyName" gorm:"type:varchar(50);comment:公司名称"`
	CompanyNature     int    `json:"companyNature" gorm:"type:tinyint(1);comment:公司性质 （2终端，3分销）"`                     // 公司性质 （2终端，3分销）
	CompanyType       int    `json:"companyType" gorm:"type:tinyint(1);comment:公司类型（1:KA 2:SME 3:DS 4:整包 5:央企 6:其他）"` // （1:KA 2:SME 3:DS 4:整包 5:央企 6:其他）
	IsPunchout        int    `json:"isPunchout" gorm:"type:tinyint(1);comment:punchout服务(1是 0否)"`                     // punchout服务(1是 0否)
	IsEas             int    `json:"isEas" gorm:"type:tinyint(1);comment:eas服务(1是 0否)"`                               // 领用单是否审批(1是 0否)
	IsEis             int    `json:"isEis" gorm:"type:tinyint(1);comment:eis服务(1是 0否)"`                               // eis服务(1是 0否)
	ParentId          int    `json:"parentId" gorm:"type:int;comment:父级节点"`                                           // 父公司id
	Address           string `json:"address" gorm:"type:varchar(200);comment:地址"`                                     //地址
	TaxNo             string `json:"taxNo" gorm:"type:varchar(50);comment:税号"`                                        //税号
	BankName          string `json:"bankName" gorm:"type:varchar(100);comment:开户行"`                                   //开户行
	BankAccount       string `json:"bankAccount" gorm:"type:varchar(20);comment:银行帐号"`                                //银行帐号
	CompanyLogo       string `json:"companyLogo" gorm:"type:varchar(200);comment:公司logo"`                             // 公司logo
	CompanyMediumLogo string `json:"companyMediumLogo" gorm:"type:varchar(200);comment:中等大小logo"`                     // 中等大小logo
	CompanySmallLogo  string `json:"companySmallLogo" gorm:"type:varchar(200);comment:小logo"`                         // 小logo
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
	models.ModelTime
	models.ControlBy
}

type Warehouse struct {
	models.Model

	WarehouseCode string `json:"warehouseCode" gorm:"type:varchar(20);comment:仓库编码"`
	WarehouseName string `json:"warehouseName" gorm:"type:varchar(20);comment:仓库名称"`
	CompanyId     int    `json:"companyId" gorm:"type:int(10);comment:仓库对应公司d"`
	Mobile        string `json:"mobile" gorm:"type:varchar(20);comment:Mobile"`
	Linkman       string `json:"linkman" gorm:"type:varchar(50);comment:联系人"`
	Email         string `json:"email" gorm:"type:varchar(50);comment:邮箱"`
	Status        string `json:"status" gorm:"type:tinyint(1) unsigned;comment:是否使用 0-否，1-是"`
	IsVirtual     string `json:"isVirtual" gorm:"type:tinyint(1);comment:是否为虚拟仓 0-否，1-是"`
	PostCode      string `json:"postCode" gorm:"type:varchar(50);comment:仓库所在地址邮编"`
	Province      int    `json:"province" gorm:"type:int unsigned;comment:省"`
	City          int    `json:"city" gorm:"type:int unsigned;comment:市"`
	District      int    `json:"district" gorm:"type:int unsigned;comment:区"`
	Address       string `json:"address" gorm:"type:varchar(100);comment:地址"`
	Remark        string `json:"remark" gorm:"type:varchar(255);comment:Remark"`
	models.RegionName
	models.ModelTime
	models.ControlBy
}

// 生成Token用户信息
type TokenUserInfo struct {
	User UserInfo
	Role RoleInfo
}
