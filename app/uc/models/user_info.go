package models

import (
	"fmt"
	"go-admin/common/global"
	"go-admin/common/middleware/mall_handler"
	"go-admin/common/utils"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"gorm.io/gorm"

	"crypto/md5"
	"encoding/hex"
	"go-admin/common/models"
	"strconv"
)

const UserInfoModel = "userInfo"
const UserInfoOperationUpdateUserPhone = "updateUserPhone"
const UserInfoProxyLoginToken = "EHSY20171227$@#&*(123"

type UserInfo struct {
	models.Model

	UserEmail            string `json:"userEmail" gorm:"type:varchar(50);comment:用户邮箱"`   //用户邮箱
	LoginName            string `json:"loginName" gorm:"type:varchar(100);comment:用户登录名"` //用户登录名
	UserPhone            string `json:"userPhone" gorm:"type:varchar(16);comment:手机号码"`   //手机号码
	LoginPassword        string `json:"loginPassword" gorm:"type:varchar(100);comment:登录密码"`
	UserName             string `json:"userName" gorm:"type:varchar(50);comment:用户名称"`      // 用户名称
	Gender               string `json:"gender" gorm:"type:varchar(2);comment:性别 (0:女 1:男)"` //性别 (0:女 1:男)
	HeadPortrait         string `json:"headPortrait" gorm:"type:varchar(100);comment:头像地址"` //头像地址
	RegisterMode         int    `json:"registerMode" gorm:"type:tinyint unsigned;comment:注册方式 (1 手机注册，2邮箱注册,3 微信注册)"`
	PhoneType            int    `json:"phoneType" gorm:"type:tinyint unsigned;comment:手机类型(1 IOS 2android)"`
	DeviceToken          string `json:"deviceToken" gorm:"type:varchar(64);comment:手机标识码"`
	UserSource           int    `json:"userSource" gorm:"type:tinyint unsigned;comment:用户来源（1www,2opc,3android,4IOS,5M,6微信,7微信小程序，8脉信，9CRM）"`
	UserStatus           int    `json:"userStatus" gorm:"type:tinyint(1);comment:用户状态（1可用，0不可用）"` //用户状态（1可用，0不可用）
	CompanyId            int    `json:"companyId" gorm:"type:int unsigned;comment:公司ID"`          //公司ID
	EmailApproveCronExpr string `json:"emailApproveCronExpr" gorm:"type:varchar(200);comment:邮件定时规则表达式"`
	IsOpen               int    `json:"isOpen" gorm:"type:int;comment:此用户的邮件审批是否开启"`
	UserType             string `json:"userType" gorm:"type:varchar(2);comment:用户类型 1.下单人2联系人3审批人4通用"` //用户类型 1.下单人2联系人3审批人4通用
	Position             string `json:"position" gorm:"type:varchar(100);comment:Position"`
	RegisterSource       string `json:"registerSource" gorm:"type:varchar(100);comment:注册用户渠道来源来源"`
	DeliveryMailset      string `json:"deliveryMailset" gorm:"type:varchar(10);comment:发货邮件通知设置，英文逗号分割{1,2,3=下单人,收货人,销售}"`
	InvoiceMailset       string `json:"invoiceMailset" gorm:"type:varchar(10);comment:发票邮件通知设置，英文逗号分割{1,2,3,4=下单人,收货人,销售,收票人}"`
	EasReceiveMsgSet     string `json:"easReceiveMsgSet" gorm:"type:varchar(10);comment:EAS审批短信通知设置，英文逗号分割{1=采购人（审批驳回）;2=审批人（采购审批）}"`
	IsAdminShow          string `json:"isAdminShow" gorm:"type:varchar(2);comment:用户状态（1是，0否）"`
	Telephone            string `json:"telephone" gorm:"type:varchar(50);comment:用户座机号"`                            //用户座机号
	EmailVerifyStatus    int    `json:"emailVerifyStatus" gorm:"type:tinyint(1);comment:邮箱验证状态（0：未验证 1：已验证 2验证失败）"` //邮箱验证状态（0：未验证 1：已验证 2验证失败）
	CanLogin             int    `json:"canLogin" gorm:"type:tinyint(1);comment:是否可登陆"`                              //是否可登陆
	CompanyDepartmentId  int    `json:"companyDepartmentId" gorm:"type:int unsigned;comment:用户所属部门"`                //用户所属部门
	models.ModelTime
	models.ControlBy

	CompanyInfo *CompanyInfo             `json:"companyInfo" gorm:"foreignKey:company_id"`
	UserRole    []*UserRole              `json:"userRole" gorm:"foreignKey:user_id"`
	RoleInfo    []RoleInfo               `json:"roleInfo" gorm:"-"`
	SelectMenu  []ManageMenu             `json:"selectMenu" gorm:"-"`
	RouterMenu  []RouterMenu             `json:"routerMenu" gorm:"-"`
	Department  CompanyDepartment        `json:"department" gorm:"foreignKey:company_department_id"`
	UserConfig  *mall_handler.UserConfig `json:"userConfig" gorm:"-"`
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

func (e *UserInfo) GenerateAccessToken(id int, tenantId, loginPassword, systemId, token string) string {
	data := strconv.Itoa(id) + tenantId + loginPassword + systemId + token
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

// 根据id批量获取用户

func (e *UserInfo) GetUsersByIds(tx *gorm.DB, userIds []int) (users *[]UserInfo, err error) {
	ucPrefix := global.GetTenantUcDBNameWithDB(tx)
	users = &[]UserInfo{}
	err = tx.Table(ucPrefix+".user_info").Find(users, userIds).Error
	return
}

func (e *UserInfo) GetFormatNameUsersByIds(tx *gorm.DB, userIds []int) (map[int]string, error) {
	outData := map[int]string{}
	users, err := e.GetUsersByIds(tx, userIds)
	if err != nil {
		return outData, err
	}
	outData = lo.Associate(*users, func(item UserInfo) (int, string) {
		return item.Id, item.FormatUserName()
	})
	return outData, nil
}

func (e *UserInfo) FormatUserName() string {
	if e.UserName != "" {
		return e.UserName
	}
	if e.LoginName != "" {
		return e.LoginName
	}
	if e.UserPhone != "" {
		return e.UserPhone
	}
	if e.UserEmail != "" {
		return e.UserEmail
	}
	return ""
}

func GetApproveUserByCompanyId(tx *gorm.DB, companyId int) (*[]UserInfo, error) {
	return GetApproveUserByCompanyIdAndRoleKey(tx, companyId, "approver")
}

func GetApproveUserByCompanyIdAndRoleKey(tx *gorm.DB, companyId int, roleKey string) (*[]UserInfo, error) {
	users := &[]UserInfo{}
	err := tx.Model(&UserInfo{}).
		Joins("LEFT JOIN user_role ON user_role.user_id = user_info.id").
		Joins("LEFT JOIN role_info ON user_role.role_id = role_info.id").
		Select("user_info.*").
		Where("role_info.role_key = ?", roleKey).
		Where("user_info.company_id = ?", companyId).
		Group("user_info.id").
		Find(users).Error
	return users, err
}

func GetRecipientUserByCompanyId(tx *gorm.DB, companyId int) (*[]UserInfo, error) {
	return GetApproveUserByCompanyIdAndRoleKey(tx, companyId, "recipient")
}

// 用户代登录URL地址
func (e *UserInfo) ProxyLogin(tx *gorm.DB, userId int) (string, error) {
	// 查询公司信息
	CompanyInfoModels := CompanyInfo{}
	companyDomain, err := CompanyInfoModels.GetPCompanyDomainByID(tx, e.CompanyId)
	if err != nil {
		return "", err
	}

	// 获取多租户ID
	ctx := tx.Statement.Context.(*gin.Context)
	tenantId := ctx.GetHeader("tenant-id")

	// 组合代登录地址
	token := UserInfoProxyLoginToken
	systemId := e.LoginName + "system" // 兼容登录账户有空值情况
	accessToken := e.GenerateAccessToken(userId, tenantId, e.LoginPassword, systemId, token)

	customerWWWLoginURL := ""
	if len(companyDomain) > 0 {
		customerWWWLoginURL = "http://" + companyDomain
	} else {
		customerWWWLoginURL = utils.GetHostUrl(0)
	}

	redirectURL := fmt.Sprintf("%s/#/login?code=%d&deCode=%s&accessToken=%s", customerWWWLoginURL, userId, systemId, accessToken)

	return redirectURL, nil
}
