package dto

import (
	"go-admin/app/uc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type UserInfoGetPageReq struct {
	dto.Pagination      `search:"-"`
	Id                  string `form:"id"  search:"type:exact;column:id;table:user_info" comment:"id"`
	UserEmail           string `form:"userEmail"  search:"type:exact;column:user_email;table:user_info" comment:"用户邮箱"`
	LoginName           string `form:"loginName"  search:"type:exact;column:login_name;table:user_info" comment:"用户登录名"`
	UserPhone           string `form:"userPhone"  search:"type:exact;column:user_phone;table:user_info" comment:"手机号码"`
	UserName            string `form:"userName"  search:"type:exact;column:user_name;table:user_info" comment:"用户名称"`
	UserStatus          string `form:"userStatus"  search:"type:exact;column:user_status;table:user_info" comment:"用户状态（1可用，0不可用）"`
	CompanyId           int    `form:"companyId"  search:"type:exact;column:company_id;table:user_info" comment:"公司ID"`
	IsAdminShow         string `form:"isAdminShow"  search:"type:exact;column:is_admin_show;table:user_info" comment:"用户状态（1是，0否）"`
	CanLogin            int    `form:"canLogin"  search:"type:exact;column:can_login;table:user_info" comment:"是否可登陆"`
	CompanyDepartmentId int    `form:"companyDepartmentId"  search:"type:exact;column:company_department_id;table:user_info" comment:"用户所属部门"`

	FilterLoginName           string `form:"filterLoginName" search:"type:contains;column:login_name;table:user_info"`
	FilterUserPhone           string `form:"filterUserPhone" search:"type:contains;column:user_phone;table:user_info"`
	FilterUserName            string `form:"filterUserName" search:"type:contains;column:user_name;table:user_info"`
	FilterUserEmail           string `form:"filterUserEmail" search:"type:contains;column:user_email;table:user_info"`
	FilterCompanyDepartmentId int    `form:"filterCompanyDepartmentId" search:"-"`
	RoleId                    int    `form:"roleId"  search:"-"`
	UserInfoOrder
}

type UserInfoOrder struct {
	Id                   string `form:"idOrder"  search:"type:order;column:id;table:user_info"`
	UserEmail            string `form:"userEmailOrder"  search:"type:order;column:user_email;table:user_info"`
	LoginName            string `form:"loginNameOrder"  search:"type:order;column:login_name;table:user_info"`
	UserPhone            string `form:"userPhoneOrder"  search:"type:order;column:user_phone;table:user_info"`
	LoginPassword        string `form:"loginPasswordOrder"  search:"type:order;column:login_password;table:user_info"`
	UserName             string `form:"userNameOrder"  search:"type:order;column:user_name;table:user_info"`
	Gender               string `form:"genderOrder"  search:"type:order;column:gender;table:user_info"`
	HeadPortrait         string `form:"headPortraitOrder"  search:"type:order;column:head_portrait;table:user_info"`
	RegisterMode         string `form:"registerModeOrder"  search:"type:order;column:register_mode;table:user_info"`
	PhoneType            string `form:"phoneTypeOrder"  search:"type:order;column:phone_type;table:user_info"`
	DeviceToken          string `form:"deviceTokenOrder"  search:"type:order;column:device_token;table:user_info"`
	UserSource           string `form:"userSourceOrder"  search:"type:order;column:user_source;table:user_info"`
	UserStatus           string `form:"userStatusOrder"  search:"type:order;column:user_status;table:user_info"`
	CompanyId            string `form:"companyIdOrder"  search:"type:order;column:company_id;table:user_info"`
	EmailApproveCronExpr string `form:"emailApproveCronExprOrder"  search:"type:order;column:email_approve_cron_expr;table:user_info"`
	IsOpen               string `form:"isOpenOrder"  search:"type:order;column:is_open;table:user_info"`
	UserType             string `form:"userTypeOrder"  search:"type:order;column:user_type;table:user_info"`
	Position             string `form:"positionOrder"  search:"type:order;column:position;table:user_info"`
	RegisterSource       string `form:"registerSourceOrder"  search:"type:order;column:register_source;table:user_info"`
	DeliveryMailset      string `form:"deliveryMailsetOrder"  search:"type:order;column:delivery_mailset;table:user_info"`
	InvoiceMailset       string `form:"invoiceMailsetOrder"  search:"type:order;column:invoice_mailset;table:user_info"`
	EasReceiveMsgSet     string `form:"easReceiveMsgSetOrder"  search:"type:order;column:eas_receive_msg_set;table:user_info"`
	IsAdminShow          string `form:"isAdminShowOrder"  search:"type:order;column:is_admin_show;table:user_info"`
	Telephone            string `form:"telephoneOrder"  search:"type:order;column:telephone;table:user_info"`
	EmailVerifyStatus    string `form:"emailVerifyStatusOrder"  search:"type:order;column:email_verify_status;table:user_info"`
	CanLogin             string `form:"canLoginOrder"  search:"type:order;column:can_login;table:user_info"`
	CompanyDepartmentId  string `form:"companyDepartmentIdOrder"  search:"type:order;column:company_department_id;table:user_info"`
	CreatedAt            string `form:"createdAtOrder"  search:"type:order;column:created_at;table:user_info"`
	UpdatedAt            string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:user_info"`
	CreateByName         string `form:"createByNameOrder"  search:"type:order;column:create_by_name;table:user_info"`
	UpdateByName         string `form:"updateByNameOrder"  search:"type:order;column:update_by_name;table:user_info"`
	DeletedAt            string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:user_info"`
}

func (m *UserInfoGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type UserInfoGetListPageRes struct {
	Id                    int    `json:"id"`
	UserEmail             string `json:"userEmail"  search:"type:exact;column:user_email;table:user_info" comment:"用户邮箱"`             // 用户邮箱
	LoginName             string `json:"loginName"  search:"type:exact;column:login_name;table:user_info" comment:"用户登录名"`            // 用户登录名
	UserPhone             string `json:"userPhone"  search:"type:exact;column:user_phone;table:user_info" comment:"手机号码"`             //手机号码
	UserName              string `json:"userName"  search:"type:exact;column:user_name;table:user_info" comment:"用户名称"`               // 用户名称
	UserStatus            int    `json:"userStatus"  search:"type:exact;column:user_status;table:user_info" comment:"用户状态（1可用，0不可用）"` //用户状态
	CompanyId             int    `json:"companyId"  search:"type:exact;column:company_id;table:user_info" comment:"公司ID"`             // 公司ID
	IsAdminShow           string `json:"isAdminShow"  search:"type:exact;column:is_admin_show;table:user_info" comment:"用户状态（1是，0否）"` //用户状态（1是，0否）
	EmailVerifyStatus     int    `json:"emailVerifyStatus" comment:"邮箱验证状态（0：未验证 1：已验证 2验证失败）"`                                       //邮箱验证状态（0：未验证 1：已验证 2验证失败）
	CompanyDepartmentId   int    `json:"companyDepartmentId" comment:"用户所属部门"`                                                        //用户所属部门id
	CompanyDepartmentText string `json:"companyDepartmentText" gorm:"-" comment:"用户所属部门"`
	UserRoleText          string `json:"userRoleText" gorm:"-"`

	UserRole    []*models.UserRole `json:"userRole" gorm:"foreignKey:user_id"`
	CompanyInfo *CompanyInfo       `json:"companyInfo" gorm:"foreignKey:company_id;table:company_info"`
	common.ModelTime
}
type CompanyInfo struct {
	Id          int    `json:"id"`
	CompanyName string `json:"companyName" comment:"用户所属公司"` // 用户公司名称
}

type UserInfoGetSelectListRes struct {
	Id                    int    `json:"id"`
	UserEmail             string `json:"userEmail"  search:"type:exact;column:user_email;table:user_info" comment:"用户邮箱"`             // 用户邮箱
	LoginName             string `json:"loginName"  search:"type:exact;column:login_name;table:user_info" comment:"用户登录名"`            // 用户登录名
	UserPhone             string `json:"userPhone"  search:"type:exact;column:user_phone;table:user_info" comment:"手机号码"`             //手机号码
	UserName              string `json:"userName"  search:"type:exact;column:user_name;table:user_info" comment:"用户名称"`               // 用户名称
	UserStatus            int    `json:"userStatus"  search:"type:exact;column:user_status;table:user_info" comment:"用户状态（1可用，0不可用）"` //用户状态
	CompanyId             int    `json:"companyId"  search:"type:exact;column:company_id;table:user_info" comment:"公司ID"`             // 公司ID
	IsAdminShow           string `json:"isAdminShow"  search:"type:exact;column:is_admin_show;table:user_info" comment:"用户状态（1是，0否）"` //用户状态（1是，0否）
	EmailVerifyStatus     int    `json:"emailVerifyStatus" comment:"邮箱验证状态（0：未验证 1：已验证 2验证失败）"`                                       //邮箱验证状态（0：未验证 1：已验证 2验证失败）
	CompanyDepartmentId   int    `json:"companyDepartmentId" comment:"用户所属部门"`                                                        //用户所属部门id
	CompanyName           string `json:"companyName" gorm:"-" comment:"用户所属公司"`                                                       // 用户公司名称
	CompanyDepartmentText string `json:"companyDepartmentText" gorm:"-" comment:"用户所属部门"`
}

type UserInfoInsertReq struct {
	Id                   int    `json:"-" comment:""` //
	UserEmail            string `json:"userEmail" comment:"用户邮箱"`
	LoginName            string `json:"loginName" comment:"用户登录名"`
	UserPhone            string `json:"userPhone" comment:"手机号码"`
	LoginPassword        string `json:"loginPassword" comment:"登录密码" vd:"len($)>0" `
	UserName             string `json:"userName" comment:"用户名称" vd:"len($)>0" `
	Gender               string `json:"gender" comment:"性别 (0:女 1:男)"`
	HeadPortrait         string `json:"headPortrait" comment:"头像地址"`
	UserStatus           int    `json:"userStatus" comment:"用户状态（1可用，0不可用）"`
	CompanyId            int    `json:"companyId" comment:"公司ID" vd:"$>0" `
	Telephone            string `json:"telephone" comment:"用户座机号"`
	CompanyDepartmentId  int    `json:"companyDepartmentId" comment:"用户所属部门"`
	CompanyDepartmentIds []int  `json:"companyDepartmentIds" comment:"用户所属部门"`
	UserRole             []int  `json:"userRole"` // 用户角色
	EmailVerifyStatus    int    `json:"-" comment:"邮箱验证状态（0：未验证 1：已验证 2验证失败）"`
	RegisterMode         int    `json:"-" comment:"注册方式 (1 手机注册，2邮箱注册,3 微信注册)"`
	PhoneType            int    `json:"-" comment:"手机类型(1 IOS 2android)"`
	DeviceToken          string `json:"-" comment:"手机标识码"`
	UserSource           int    `json:"-" comment:"用户来源（1www,2opc,3android,4IOS,5M,6微信,7微信小程序，8脉信，9CRM）"`
	EmailApproveCronExpr string `json:"-" comment:"邮件定时规则表达式"`
	IsOpen               int    `json:"-" comment:"此用户的邮件审批是否开启"`
	UserType             string `json:"-" comment:"用户类型 1.下单人2联系人3审批人4通用"`
	Position             string `json:"-" comment:""`
	RegisterSource       string `json:"-" comment:"注册用户渠道来源来源"`
	DeliveryMailset      string `json:"-" comment:"发货邮件通知设置，英文逗号分割{1,2,3=下单人,收货人,销售}"`
	InvoiceMailset       string `json:"-" comment:"发票邮件通知设置，英文逗号分割{1,2,3,4=下单人,收货人,销售,收票人}"`
	EasReceiveMsgSet     string `json:"-" comment:"EAS审批短信通知设置，英文逗号分割{1=采购人（审批驳回）;2=审批人（采购审批）}"`
	IsAdminShow          string `json:"-" comment:"用户状态（1是，0否）"`
	CanLogin             int    `json:"-" comment:"是否可登陆"`

	common.ControlBy `json:"-"`
}

type CompanyDepartmentIds struct {
	level1Id int
	level2Id int
}

func (s *UserInfoInsertReq) Generate(model *models.UserInfo) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.UserEmail = s.UserEmail
	model.EmailVerifyStatus = 0
	if s.UserEmail != "" {
		model.EmailVerifyStatus = 1
	}
	model.LoginName = s.LoginName
	model.UserPhone = s.UserPhone
	model.LoginPassword = s.LoginPassword
	model.UserName = s.UserName
	model.Gender = s.Gender
	model.HeadPortrait = s.HeadPortrait
	model.RegisterMode = s.RegisterMode
	model.PhoneType = s.PhoneType
	model.DeviceToken = s.DeviceToken
	model.UserSource = s.UserSource
	model.UserStatus = s.UserStatus
	model.CompanyId = s.CompanyId
	model.EmailApproveCronExpr = s.EmailApproveCronExpr
	model.IsOpen = s.IsOpen
	model.UserType = s.UserType
	model.Position = s.Position
	model.RegisterSource = s.RegisterSource
	model.DeliveryMailset = s.DeliveryMailset
	model.InvoiceMailset = s.InvoiceMailset
	model.EasReceiveMsgSet = s.EasReceiveMsgSet
	model.IsAdminShow = s.IsAdminShow
	model.Telephone = s.Telephone
	model.CanLogin = s.CanLogin
	model.CompanyDepartmentId = s.CompanyDepartmentId
	model.CreateByName = s.CreateByName
}

func (s *UserInfoInsertReq) GetId() interface{} {
	return s.Id
}

type UserInfoUpdateReq struct {
	Id                   int    `uri:"id" comment:""` //
	UserEmail            string `json:"userEmail" comment:"用户邮箱"  `
	LoginName            string `json:"loginName" comment:"用户登录名"`
	UserPhone            string `json:"userPhone" comment:"手机号码"  `
	LoginPassword        string `json:"loginPassword" comment:"登录密码"`
	UserName             string `json:"userName" comment:"用户名称" vd:"len($)>0" `
	Gender               string `json:"gender" comment:"性别 (0:女 1:男)"`
	HeadPortrait         string `json:"headPortrait" comment:"头像地址"`
	UserStatus           int    `json:"userStatus" comment:"用户状态（1可用，0不可用）"`
	CompanyId            int    `json:"companyId" comment:"公司ID" vd:"$>0" `
	Telephone            string `json:"telephone" comment:"用户座机号"`
	CompanyDepartmentId  int    `json:"companyDepartmentId" comment:"用户所属部门"`
	UserRole             []int  `json:"userRole"` // 用户角色
	EmailVerifyStatus    int    `json:"-" comment:"邮箱验证状态（0：未验证 1：已验证 2验证失败）"`
	RegisterMode         int    `json:"-" comment:"注册方式 (1 手机注册，2邮箱注册,3 微信注册)"`
	PhoneType            int    `json:"-" comment:"手机类型(1 IOS 2android)"`
	DeviceToken          string `json:"-" comment:"手机标识码"`
	UserSource           int    `json:"-" comment:"用户来源（1www,2opc,3android,4IOS,5M,6微信,7微信小程序，8脉信，9CRM）"`
	EmailApproveCronExpr string `json:"-" comment:"邮件定时规则表达式"`
	IsOpen               int    `json:"-" comment:"此用户的邮件审批是否开启"`
	UserType             string `json:"-" comment:"用户类型 1.下单人2联系人3审批人4通用"`
	Position             string `json:"-" comment:""`
	RegisterSource       string `json:"-" comment:"注册用户渠道来源来源"`
	DeliveryMailset      string `json:"-" comment:"发货邮件通知设置，英文逗号分割{1,2,3=下单人,收货人,销售}"`
	InvoiceMailset       string `json:"-" comment:"发票邮件通知设置，英文逗号分割{1,2,3,4=下单人,收货人,销售,收票人}"`
	EasReceiveMsgSet     string `json:"-" comment:"EAS审批短信通知设置，英文逗号分割{1=采购人（审批驳回）;2=审批人（采购审批）}"`
	IsAdminShow          string `json:"-" comment:"用户状态（1是，0否）"`
	CanLogin             int    `json:"-" comment:"是否可登陆"`

	common.ControlBy `json:"-"`
}

func (s *UserInfoUpdateReq) Generate(model *models.UserInfo) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.UserEmail = s.UserEmail
	model.EmailVerifyStatus = 0
	if s.UserEmail != "" {
		model.EmailVerifyStatus = 1
	}
	model.LoginName = s.LoginName
	model.UserPhone = s.UserPhone
	if s.LoginPassword != "" {
		model.LoginPassword = s.LoginPassword
	}
	model.UserName = s.UserName
	model.Gender = s.Gender
	model.HeadPortrait = s.HeadPortrait
	model.RegisterMode = s.RegisterMode
	model.PhoneType = s.PhoneType
	model.DeviceToken = s.DeviceToken
	model.UserSource = s.UserSource
	model.UserStatus = s.UserStatus
	model.CompanyId = s.CompanyId
	model.EmailApproveCronExpr = s.EmailApproveCronExpr
	model.IsOpen = s.IsOpen
	model.UserType = s.UserType
	model.Position = s.Position
	model.RegisterSource = s.RegisterSource
	model.DeliveryMailset = s.DeliveryMailset
	model.InvoiceMailset = s.InvoiceMailset
	model.EasReceiveMsgSet = s.EasReceiveMsgSet
	model.IsAdminShow = s.IsAdminShow
	model.Telephone = s.Telephone
	model.CanLogin = s.CanLogin
	model.CompanyDepartmentId = s.CompanyDepartmentId
	model.UpdateByName = s.UpdateByName // 添加这而，需要记录是被谁更新的
}

func (s *UserInfoUpdateReq) GetId() interface{} {
	return s.Id
}

// UserInfoGetReq 功能获取请求参数
type UserInfoGetReq struct {
	Id int `uri:"id"`
}

func (s *UserInfoGetReq) GetId() interface{} {
	return s.Id
}

type UserInfoGetRes struct {
	models.UserInfo
	UserDepartments UserInfoGetDepartmentNames `json:"userDepartments" gorm:"-"`
}

type UserInfoGetDepartmentNames struct {
	Names []string `json:"names"`
	Ids   []int    `json:"ids"`
}

// UserInfoDeleteReq 功能删除请求参数
type UserInfoDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *UserInfoDeleteReq) GetId() interface{} {
	return s.Ids
}

type UserInfoUpdatePassword struct {
	Id            int    `json:"id" comment:"" vd:"$>0" ` //
	LoginPassword string `json:"loginPassword" comment:"登录密码" vd:"len($)>0" `

	common.ControlBy `json:"-"`
}

type UserInfoUpdateUserPhone struct {
	Id        int    `json:"id" comment:"" vd:"$>0" ` //
	UserPhone string `json:"userPhone" vd:"len($)>0" `

	common.ControlBy `json:"-"`
}

type UserInfoProxyLoginReq struct {
	SellPeople       string `form:"sellPeople"`
	UserID           int    `form:"userId"`
	Website          string `form:"website"`
	common.ControlBy `json:"-"`
}

type UserInfoProxyLoginRes struct {
	RedirectURL string `json:"redirectURL"`
	DevModel    string `json:"devModel"`
}

type UserInfoGetUserCompanyInfo struct {
	CompanyId   int    `json:"companyId"`
	CompanyName string `json:"companyName"`
}

type UserInfoImportData struct {
	UserName        string `json:"userName"`
	UserPhone       string `json:"userPhone"`
	Telephone       string `json:"telephone"`
	LoginName       string `json:"loginName"`
	UserEmail       string `json:"userEmail"`
	LoginPassword   string `json:"loginPassword"`
	UserRole        string `json:"userRole"`
	CompanyName     string `json:"companyName"`
	DepartmentName  string `json:"departmentName"`
	FDepartmentName string `json:"FDepartmentName"`

	CompanyId    int   `json:"-"`
	DepartmentId int   `json:"-"`
	UserRoleList []int `json:"-"`
}
