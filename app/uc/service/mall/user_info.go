package mall

import (
	"errors"
	"fmt"
	"go-admin/common/global"
	"go-admin/common/msg/sms"
	"html"
	"strconv"
	"strings"
	"time"

	"github.com/go-admin-team/go-admin-core/sdk/config"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	"github.com/jinzhu/copier"

	"go-admin/app/uc/models"
	"go-admin/app/uc/service/mall/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
	"go-admin/common/middleware"
	"go-admin/common/middleware/mall_handler"
	modelsCommon "go-admin/common/models"
	"go-admin/common/msg/email"
	"go-admin/common/utils"

	"github.com/go-admin-team/go-admin-core/sdk/pkg/captcha"

	"github.com/gin-gonic/gin"

	"github.com/golang-jwt/jwt/v4"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type UserInfo struct {
	service.Service
}

// GetPage 获取UserInfo列表
func (e *UserInfo) GetPage(c *dto.UserInfoGetPageReq, p *actions.DataPermission, list *[]models.UserInfo, count *int64) error {
	var err error
	var data models.UserInfo

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("UserInfoService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取UserInfo对象
func (e *UserInfo) Get(d *dto.UserInfoGetReq, p *actions.DataPermission, model *models.UserInfo) error {
	var data models.UserInfo

	// 查询用户信息
	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		Preload("CompanyInfo").
		Preload("UserRole").
		Preload("Department").
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetUserInfo error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	model.LoginPassword = ""

	// 查询顶级公司信息
	topCompany := models.CompanyInfo{}
	err = e.Orm.Where("pid = ?", 0).First(&topCompany).Error
	if err != nil {
		return err
	}

	// 使用顶级公司的logo
	if topCompany.CompanyLogo == "" {
		topCompany.CompanyLogo = "https://image-c.ehsy.com/uploadfile/sxyz/img/2022/09/16/20220916100258386.png"
	}
	model.CompanyInfo.CompanyLogo = topCompany.CompanyLogo

	// 查询角色信息
	var ids []int
	for _, item := range model.UserRole {
		ids = append(ids, item.RoleId)
	}
	roles := []models.RoleInfo{}
	err = e.Orm.Where("id in ?", ids).Find(&roles).Error
	if err != nil {
		return nil
	}
	model.RoleInfo = roles
	model.UserRole = nil

	// 查询角色拥有的权限
	var menusIds []string
	for _, item := range roles {
		menu := strings.Split(item.Menus, ",")
		menusIds = append(menusIds, menu...)
	}
	menusIds = utils.ArrayUnique(menusIds)
	menus := []models.ManageMenu{}
	err = e.Orm.Where("id in ?", menusIds).Order("order_by ASC, id ASC").Find(&menus).Error
	if err != nil {
		return err
	}

	// 菜单遍历
	var keyMenus = make(map[string][]models.ManageMenu)
	var selectMenus []models.ManageMenu
	var keys []string
	for i := 0; i < len(menus); i++ {
		// 左侧菜单
		if menus[i].Type == 1 {
			keys = append(keys, menus[i].GroupName)
			keyMenus[menus[i].GroupName] = append(keyMenus[menus[i].GroupName], menus[i])
		}
		// 下拉菜单
		if menus[i].Type == 2 {
			selectMenus = append(selectMenus, menus[i])
		}
	}
	keys = utils.ArrayUnique(keys)
	res := []models.RouterMenu{}
	for _, key := range keys {
		res = append(res, models.RouterMenu{
			Title:    keyMenus[key][0].GroupName,
			Children: keyMenus[key],
		})
	}
	model.RouterMenu = res
	model.SelectMenu = selectMenus

	return nil
}

// Insert 创建UserInfo对象
func (e *UserInfo) Insert(c *dto.UserInfoInsertReq) error {
	var err error
	var data models.UserInfo
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("UserInfoService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改UserInfo对象
func (e *UserInfo) Update(c *dto.UserInfoUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.UserInfo{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("UserInfoService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除UserInfo
func (e *UserInfo) Remove(d *dto.UserInfoDeleteReq, p *actions.DataPermission) error {
	var data models.UserInfo

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveUserInfo error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}

func (e *UserInfo) ChangePassword(d *dto.UserInfoChangePasswordReq) (err error) {
	// 图形验证码校验
	checkCaptcha := captcha.Verify(d.CodeId, d.Code, true)
	if !checkCaptcha {
		return errors.New("图形验证码输入不正确，请重试")
	}

	var userInfo models.UserInfo
	err = e.Orm.Debug().Table("user_info").Where("id=?", d.UserId).First(&userInfo).Error
	if err != nil {
		return err
	}

	if utils.Md5Uc(d.Password) != userInfo.LoginPassword {
		return errors.New("原始密码不正确！")
	}
	if d.NewPassword != d.NewPasswordCheck {
		return errors.New("新密码与确认密码不一样！")
	}
	// 校验密码|大小写,英文，数字，符号
	checkPwd := utils.ValidatePassword(d.NewPassword)
	if !checkPwd {
		return errors.New("密码要求:8-20位,由大写字母、小写字母、数字和英文字符(除空格外)至少三种组成")
	}

	// 使用base权限DB
	dsn := mysql.Open(config.DatabasesConfig["base"].Source)
	baseDb, err := gorm.Open(dsn, &gorm.Config{})
	if err != nil {
		return errors.New("数据库连接失败")
	}

	// 查宽表数据
	ctx := e.Orm.Statement.Context.(*gin.Context)
	tenantId := ctx.GetHeader("tenant-id")
	mallUser := &modelsCommon.MallUsers{}
	err = baseDb.Table("common.mall_users").Where("tenant_id", tenantId).Where("user_phone", userInfo.UserPhone).First(mallUser).Error
	if err != nil {
		return err
	}

	// 事务落库
	err = baseDb.Transaction(func(tx *gorm.DB) error {
		// 用户表更新
		newPassword := utils.Md5Uc(d.NewPassword)
		userInfo.LoginPassword = newPassword
		ucPrefix := global.GetTenantUcDBNameWithDB(e.Orm)
		err = tx.Table(ucPrefix + ".user_info").Save(&userInfo).Error
		if err != nil {
			return err
		}

		// 宽表更新
		err = tx.Table("common.mall_users").Where("id = ?", mallUser.Id).Update("login_password", newPassword).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	userPasswordChangeLogModels := models.UserPasswordChangeLog{}
	userPasswordChangeLogModels.AddLog(e.Orm, d.UserId, user.GetUserName(e.Orm.Statement.Context.(*gin.Context)))
	return
}
func (e *UserInfo) ChangeBasicInfo(d *dto.UserInfoChangeBasicInfoReq) (err error) {
	var userInfo models.UserInfo
	err = e.Orm.Table("user_info").Where("id=?", d.UserId).First(&userInfo).Error
	if err != nil {
		return err
	}

	oldUserInfo := userInfo

	userInfo.LoginName = d.LoginName
	userInfo.UserName = d.UserName
	userInfo.Gender = strconv.Itoa(d.Gender)
	userInfo.HeadPortrait = d.HeadPortrait
	userInfo.LoginPassword = ""
	errs := e.SaveValidate(&oldUserInfo, &userInfo, false)
	if errs != nil {
		return errs[0]
	}

	// 使用base权限DB
	dsn := mysql.Open(config.DatabasesConfig["base"].Source)
	baseDb, err := gorm.Open(dsn, &gorm.Config{})
	if err != nil {
		return errors.New("数据库连接失败")
	}

	// 查宽表数据
	ctx := e.Orm.Statement.Context.(*gin.Context)
	tenantId := ctx.GetHeader("tenant-id")
	mallUser := &modelsCommon.MallUsers{}
	err = baseDb.Table("common.mall_users").Where("tenant_id", tenantId).Where("user_phone", oldUserInfo.UserPhone).First(mallUser).Error
	if err != nil {
		return err
	}

	// err = e.Orm.Save(&userInfo).Error
	// if err != nil {
	// 	return err
	// }

	// 事务落库
	err = baseDb.Transaction(func(tx *gorm.DB) error {
		// 用户表更新
		ucPrefix := global.GetTenantUcDBNameWithDB(e.Orm)
		err = tx.Table(ucPrefix + ".user_info").Omit("login_password").Save(&userInfo).Error
		if err != nil {
			return err
		}

		err = copier.Copy(mallUser, &userInfo)
		if err != nil {
			return err
		}

		err = tx.Table("common.mall_users").Where("id = ?", mallUser.Id).Omit("login_password").Save(mallUser).Error
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return
}

func (e *UserInfo) ChangeUserEmail(d *dto.UserInfoChangeUserEmailReq) (err error) {

	var userInfo models.UserInfo
	err = e.Orm.Table("user_info").Where("id=?", d.UserId).First(&userInfo).Error
	if err != nil {
		return err
	}

	checkUserEmailTokenAvailable, err := e.CheckUserEmailTokenAvailable(&dto.CheckUserTokenAvailableReq{
		Token:         d.Token,
		UserId:        d.UserId,
		UserEmail:     d.Email,
		LoginPassword: userInfo.LoginPassword,
	})
	if err != nil {
		return err
	}
	if !checkUserEmailTokenAvailable {
		return errors.New("邮箱验证失败，请重试")
	}
	oldUserInfo := userInfo

	userInfo.UserEmail = d.Email
	userInfo.EmailVerifyStatus = 1
	errs := e.SaveValidate(&oldUserInfo, &userInfo, false)
	if errs != nil {
		return errs[0]
	}

	// 使用base权限DB
	dsn := mysql.Open(config.DatabasesConfig["base"].Source)
	baseDb, err := gorm.Open(dsn, &gorm.Config{})
	if err != nil {
		return errors.New("数据库连接失败")
	}

	// 查宽表数据
	ctx := e.Orm.Statement.Context.(*gin.Context)
	tenantId := ctx.GetHeader("tenant-id")
	mallUser := &modelsCommon.MallUsers{}
	err = baseDb.Table("common.mall_users").Where("tenant_id", tenantId).Where("user_phone", oldUserInfo.UserPhone).First(mallUser).Error
	if err != nil {
		return err
	}

	// 事务落库
	err = baseDb.Transaction(func(tx *gorm.DB) error {
		// 用户表更新
		ucPrefix := global.GetTenantUcDBNameWithDB(e.Orm)
		err = tx.Table(ucPrefix + ".user_info").Save(&userInfo).Error
		if err != nil {
			return err
		}

		// 宽表更新
		err = tx.Table("common.mall_users").Where("id = ?", mallUser.Id).Update("user_email", userInfo.UserEmail).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return
}

func (e *UserInfo) ChangePhone(d *dto.UserInfoChangePhoneReq) (err error) {
	var userInfo models.UserInfo
	err = e.Orm.Table("user_info").Where("id=?", d.UserId).First(&userInfo).Error
	if err != nil {
		return err
	}

	// 手机号+验证码是否正确
	cacheCoke, err := mall_handler.GetCache(d.UserPhone)
	if err != nil {
		return err
	}
	if cacheCoke != d.UserPhoneCode {
		return errors.New("手机验证码失效或不正确，请重试")
	}
	oldUserInfo := userInfo

	userInfo.UserPhone = d.UserPhone
	errs := e.SaveValidate(&oldUserInfo, &userInfo, false)
	if errs != nil {
		return errs[0]
	}

	// 使用base权限DB
	dsn := mysql.Open(config.DatabasesConfig["base"].Source)
	baseDb, err := gorm.Open(dsn, &gorm.Config{})
	if err != nil {
		return errors.New("数据库连接失败")
	}

	// 查宽表数据
	ctx := e.Orm.Statement.Context.(*gin.Context)
	tenantId := ctx.GetHeader("tenant-id")
	mallUser := &modelsCommon.MallUsers{}
	err = baseDb.Table("common.mall_users").Where("tenant_id", tenantId).Where("user_phone", oldUserInfo.UserPhone).First(mallUser).Error
	if err != nil {
		return err
	}

	// 事务落库
	err = baseDb.Transaction(func(tx *gorm.DB) error {
		// 用户表更新
		ucPrefix := global.GetTenantUcDBNameWithDB(e.Orm)
		err = tx.Table(ucPrefix + ".user_info").Save(&userInfo).Error
		if err != nil {
			return err
		}

		// 宽表更新
		err = tx.Table("common.mall_users").Where("id = ?", mallUser.Id).Update("user_phone", userInfo.UserPhone).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return
}

func (e *UserInfo) SaveValidate(data *models.UserInfo, saveData *models.UserInfo, isAllErr bool) (errs []error) {
	if saveData.UserPhone == "" {
		errs = append(errs, errors.New("手机号为必填项"))
		if !isAllErr {
			return
		}
	}
	if saveData.UserName == "" {
		errs = append(errs, errors.New("用户姓名必填"))
		if !isAllErr {
			return
		}
	}

	//密码不为空 或 当前操作为新增用户 则验证密码规则
	if data.Id == 0 || (saveData.LoginPassword != "" && data.LoginPassword != saveData.LoginPassword) {
		// 校验密码|大小写,英文，数字，符号
		checkPwd := utils.ValidatePassword(saveData.LoginPassword)
		if !checkPwd {
			errs = append(errs, errors.New("密码要求:8-20位,由大写字母、小写字母、数字和英文字符(除空格外)至少三种组成"))
			if !isAllErr {
				return
			}
		}
		// 密码md5处理
		decodedPassword := html.UnescapeString(saveData.LoginPassword)
		newPassword := utils.Md5Uc(decodedPassword)
		saveData.LoginPassword = newPassword
	}

	if data.LoginName != saveData.LoginName {
		if strings.Contains(saveData.LoginName, "@") || utils.IsAllDigits(saveData.LoginName) {
			errs = append(errs, errors.New("登录账号不能纯数字、且不能包含“@”"))
			if !isAllErr {
				return
			}
		} else {
			userInfo, err := e.GetByLoginName(saveData.LoginName)
			if err != nil {
				errs = append(errs, errors.New(err.Error()))
				if !isAllErr {
					return
				}
			}
			if userInfo.Id > 0 {
				errs = append(errs, errors.New("登录名已被使用，请修改"))
				if !isAllErr {
					return
				}
			}
		}
	}

	if data.UserPhone != saveData.UserPhone {
		if !utils.IsPhoneNumber(saveData.UserPhone) {
			errs = append(errs, errors.New("手机号格式不正确"))
			if !isAllErr {
				return
			}
		}
		userInfo, err := e.GetByUserPhone(saveData.UserPhone)
		if err != nil {
			errs = append(errs, errors.New(err.Error()))
			if !isAllErr {
				return
			}
		} else if userInfo.Id > 0 {
			errs = append(errs, errors.New("手机号已被使用，请修改"))
			if !isAllErr {
				return
			}
		}
	}
	if saveData.UserEmail != "" && data.UserEmail != saveData.UserEmail {
		userInfo, err := e.GetByUserEmail(saveData.UserEmail)
		if err != nil {
			errs = append(errs, errors.New(err.Error()))
			if !isAllErr {
				return
			}
		}
		if userInfo.Id > 0 {
			errs = append(errs, errors.New("邮箱已被使用，请修改"))
			if !isAllErr {
				return
			}
		}
	}
	if isAllErr {
		return
	}
	return nil
}

// GetByLoginName 根据登录名获取用户信息
func (e *UserInfo) GetByLoginName(loginName string) (data *models.UserInfo, err error) {
	data = &models.UserInfo{}

	err = e.Orm.Model(data).
		Where("login_name = ?", loginName).
		First(data).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		e.Log.Errorf("db error:%s", err)
		return nil, err
	}
	return data, nil
}

// GetByUserEmail 根据邮箱获用户信息
func (e *UserInfo) GetByUserEmail(userEmail string) (data *models.UserInfo, err error) {
	data = &models.UserInfo{}

	err = e.Orm.Model(&data).
		Where("user_email = ?", userEmail).
		First(data).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		e.Log.Errorf("db error:%s", err)
		return
	}
	return data, nil
}

// GetByUserPhone 根据手机号获用户信息
func (e *UserInfo) GetByUserPhone(userPhone string) (data *models.UserInfo, err error) {
	data = &models.UserInfo{}

	err = e.Orm.Model(&data).
		Where("user_phone = ?", userPhone).
		First(data).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		e.Log.Errorf("db error:%s", err)
		return nil, err
	}
	return data, nil
}

// CheckPhone 实时检查手机号是否注册
func (e *UserInfo) CheckPhone(d *dto.VerifyCodeReq) error {
	// 手机号存在与否
	find := models.UserInfo{}
	err := e.Orm.Where("user_phone = ?", d.UserPhone).First(&find).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("该手机号未注册！")
		}
		return err
	}
	return nil
}

// VerifyCode 获取手机验证码
func (e *UserInfo) VerifyCode(d *dto.VerifyCodeReq, checkPhoneExist bool) (code string, err error) {
	// 校验Header中的租户ID是否存在
	ctx := e.Orm.Statement.Context.(*gin.Context)
	tenantId := ctx.GetHeader("tenant-id")
	if tenantId == "" {
		return "", errors.New("租户ID不存在，请检查！")
	}

	// 图形验证码校验
	checkCaptcha := captcha.Verify(d.Id, d.Code, true)
	if !checkCaptcha {
		return "", errors.New("图形验证码输入不正确，请重试")
	}

	// 检查手机号存在的话就不发短信
	if d.CheckPhoneExistNoSms {
		find := models.UserInfo{}
		err = e.Orm.Debug().Where("user_phone = ?", d.UserPhone).First(&find).Error
		if find.Id > 0 {
			return "", errors.New("手机号已被占用")
		}
		if err != nil {
			return "", err
		}
	}

	// 检查手机号存在与否
	if checkPhoneExist {
		find := models.UserInfo{}
		err = e.Orm.Debug().Where("user_phone = ?", d.UserPhone).First(&find).Error
		if err == gorm.ErrRecordNotFound {
			return "", errors.New("手机号不存在，请检查")
		}
		if err != nil {
			return "", err
		}
	}

	// 缓存验证码 | 3分钟
	verifyCode := strconv.Itoa(utils.GenRandNum(6))
	err = mall_handler.SetCache(d.UserPhone, verifyCode, 60*3)
	if err != nil {
		return "", nil
	}

	// 发送手机验证码
	idType := sms.MessageIdTypeSendVerifyCode
	userPhones := []string{d.UserPhone}
	replaceParams := []string{verifyCode, "1"}
	sendCode, err := sms.SendSMS(userPhones, idType, replaceParams)
	if err != nil || sendCode.Code != 0 {
		return "", err
	}

	return "", nil
}

// ComitVerifyCode 提交手机验证码
func (e *UserInfo) ComitVerifyCode(d *dto.ComitVerifyCodeReq) (any, error) {
	// 校验Header中的租户ID是否存在
	ctx := e.Orm.Statement.Context.(*gin.Context)
	tenantId := ctx.GetHeader("tenant-id")
	if tenantId == "" {
		return "", errors.New("租户ID不存在，请检查！")
	}

	// 手机号存在与否
	var token string
	find := models.UserInfo{}
	err := e.Orm.Model(&find).
		Select("SHA1(CONCAT(user_phone,login_password)) token").
		Where("user_phone = ?", d.UserPhone).First(&token).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", errors.New("手机号不存在，请检查")
		}
		return "", err
	}

	// 手机号+验证码是否正确
	cacheCoke, err := mall_handler.GetCache(d.UserPhone)
	if err != nil {
		return cacheCoke, err
	}
	if cacheCoke != d.VerifyCode {
		return "", errors.New("手机验证码失效或不正确，请重试")
	}

	// 生成加密token | phone +
	res := map[string]string{"token": token}

	return res, nil
}

// SendEmail 发送至邮箱
func (e *UserInfo) SendEmail(c *gin.Context, d *dto.SendEmailReq) error {
	// 校验Header中的租户ID是否存在
	ctx := e.Orm.Statement.Context.(*gin.Context)
	tenantId := ctx.GetHeader("tenant-id")
	if tenantId == "" {
		return errors.New("租户ID不存在，请检查！")
	}

	// 图形验证码校验
	checkCaptcha := captcha.Verify(d.Id, d.Code, true)
	if !checkCaptcha {
		return errors.New("图形验证码输入不正确，请重试")
	}

	// 校验邮箱是否存在
	var token string
	find := models.UserInfo{}
	err := e.Orm.Model(&find).
		Select("SHA1(CONCAT(user_email,login_password)) token").
		Where("user_email = ?", d.Email).First(&token).Error
	if err == gorm.ErrRecordNotFound {
		return errors.New("邮箱不存在，请检查")
	}
	if err != nil {
		return err
	}

	// 缓存Token过期时间|2小时
	err = mall_handler.SetCache(token, "1", 60*60*2)
	if err != nil {
		return err
	}

	// 读取前台地址
	timestamp := time.Now().Unix()
	timestampInt := int(timestamp)
	host := utils.GetHostUrl(0)
	url := fmt.Sprintf("%s/#/forgetPassword?token=%s&timeStamp=%s&tenantId=%s", host, token, strconv.Itoa(timestampInt), tenantId)

	// 发送邮件
	recipients := []string{d.Email}
	subject := "找回密码邮件-狮行驿站"
	body := fmt.Sprintf(`请点击以下链接完成重置密码<br/> <a href="%s">%s</a> <br/> 如果以上链接无法打开，请将上面网页地址复制到浏览器地址栏中打开（该链接2小时有效）`, url, url)
	err = email.SendEmails(recipients, subject, body)

	return err
}

// SendEmailForCheckEmail 发送至邮箱 验证邮箱
func (e *UserInfo) SendEmailForCheckEmail(d *dto.SendEmailForCheckEmailReq) error {
	// 图形验证码校验
	checkCaptcha := captcha.Verify(d.Id, d.Code, true)
	if !checkCaptcha {
		return errors.New("图形验证码输入不正确，请重试")
	}

	userInfo := &models.UserInfo{}
	err := e.Orm.Table("user_info").Where("id=?", d.UserId).First(userInfo).Error
	if err != nil {
		return err
	}

	// 生成加密Token
	token := utils.Sha1(userInfo.UserEmail + userInfo.LoginPassword)

	err = mall_handler.SetCache(token, userInfo.UserEmail, 60*60*2)
	if err != nil {
		return err
	}

	timestamp := time.Now().Unix()
	url := fmt.Sprintf(d.RedirectUri, token, strconv.FormatInt(timestamp, 10))

	// 发送邮件
	recipients := []string{d.Email}
	subject := "验证邮箱邮件-狮行驿站"
	body := fmt.Sprintf(`请点击以下链接完成重置邮箱<br/> <a href="%s">%s</a> <br/> 如果以上链接无法打开，请将上面网页地址复制到浏览器地址栏中打开（该链接2小时有效）`, url, url)
	err = email.SendEmails(recipients, subject, body)

	return err
}

func (e *UserInfo) CheckUserEmailTokenAvailable(req *dto.CheckUserTokenAvailableReq) (bool, error) {
	value, err := mall_handler.GetCache(req.Token)
	if err != nil {
		return false, errors.New("链接已失效")
	}
	if value == "" {
		return false, nil
	}

	// 生成加密Token
	userToken := utils.Sha1(req.UserEmail + req.LoginPassword)
	if userToken != req.Token {
		return false, nil
	}
	return true, nil
}

// SendEmailForChangeEmail 发送至邮箱 修改邮箱
func (e *UserInfo) SendEmailForChangeEmail(d *dto.SendEmailForChangeEmailReq) error {
	// 图形验证码校验
	checkCaptcha := captcha.Verify(d.Id, d.Code, true)
	if !checkCaptcha {
		return errors.New("图形验证码输入不正确，请重试")
	}

	// 校验邮箱是否存在
	userInfo, err := e.GetByUserEmail(d.Email)
	if err != nil {
		return err
	}
	if userInfo.Id > 0 {
		return errors.New("邮箱已被使用，请修改")
	}

	userInfo = &models.UserInfo{}
	err = e.Orm.Table("user_info").Where("id=?", d.UserId).First(userInfo).Error
	if err != nil {
		return err
	}
	userInfo.EmailVerifyStatus = 1
	e.Orm.Save(userInfo)
	if userInfo.UserEmail != "" {
		checkUserEmailTokenAvailable, err := e.CheckUserEmailTokenAvailable(&dto.CheckUserTokenAvailableReq{
			Token:         d.Token,
			UserId:        userInfo.Id,
			UserEmail:     userInfo.UserEmail,
			LoginPassword: userInfo.LoginPassword,
		})
		if err != nil {
			return err
		}
		if !checkUserEmailTokenAvailable {
			return errors.New("邮箱验证链接已失效")
		}
	}

	// 生成加密Token
	token := utils.Sha1(d.Email + userInfo.LoginPassword)

	// 缓存Token过期时间|2小时
	err = mall_handler.SetCache(token, d.Email, 60*60*2)
	if err != nil {
		return err
	}

	//url := utils.GetHostUrl(e.Orm.Statement.Context.(*gin.Context), 0)
	userEmail := d.Email

	timestamp := time.Now().Unix()
	url := fmt.Sprintf(d.RedirectUri, token, strconv.FormatInt(timestamp, 10))

	ymd := time.Now().Format("2006年01月02日")
	his := time.Now().Format("15时04分05秒")

	title := "新邮箱验证邮件-狮行驿站"
	content := `<table border='0' cellpadding='0' cellspacing='0' width='100%'>
		<tbody><tr>
			<td style='padding: 10px 0 30px 0;'>
				<table align='center' border='0' cellpadding='0' cellspacing='0' width='700' style='border: 1px solid #cccccc; border-collapse: collapse;'>
					<tbody><tr bgcolor='#007454'>
						<td width='100'>
							<img src='https://static-c.ehsy.com/content/email_image/加盟.jpg' alt='狮行驿站' style='display: block'>
						</td>
					</tr>
					<tr>
						<td style='font-size: 18px;padding: 50px 50px 15px 50px;' align='left' height='20'>
							尊敬的用户 :<a href='mailto:` + userEmail + `'>` + userEmail + `</a>
						</td>
					</tr>
					<tr>
						<td style='font-size: 14px;padding: 5px 50px 0 50px;' align='left' height='20'>
								您好，
						</td>
					</tr>
					<tr>
						<td style='font-size: 14px;padding: 5px 50px 25px 50px;' align='left' height='20'>
								您于` + ymd + " " + his + ` 申请验证邮箱，点击以下按钮，即可完成验证：
						</td>
					</tr>
					<tr>
						<td width='100' style='padding: 0 50px 20px 50px'>
							<table border='0' cellspacing='0' cellpadding='0' width='100%'>
								<tbody><tr>
									<td height='44' width='220' style='color:#ffffff;background: #01b382;font-size: 14px;border-radius:4px;' align='center'>
										<a style='text-decoration: none;color:#ffffff;font-size: 16px;' target='_blank' href='` + url + `'>验证邮箱</a>
									</td>
									<td></td>
									<td></td>
								</tr>
							</tbody></table>
						</td>
					</tr>
					<tr>
						<td style='font-size: 14px;padding: 5px 50px 20px 50px;' align='left' height='20'>
								为保障您的账号安全，请在2小时内点击该链接
						</td>
					</tr>
				</tbody></table>
			</td>
		</tr>
	</tbody></table>`

	fmt.Println(content)

	// 发送邮件
	recipients := []string{d.Email}
	err = email.SendEmails(recipients, title, content)

	return err
}

// ChangePwd 修改密码
func (e *UserInfo) ChangePwd(d *dto.ChangePwdReq) error {
	// 校验密码是否确认
	if d.Pwd != d.PwdCheck {
		return errors.New("确认密码和密码不一致，请检查")
	}

	// 校验密码|大小写,英文，数字，符号
	checkPwd := utils.ValidatePassword(d.Pwd)
	if !checkPwd {
		return errors.New("密码要求:8-20位,由大写字母、小写字母、数字和英文字符(除空格外)至少三种组成")
	}

	// Email校验是否过期|2小时
	if d.Type == "email" {
		timeCheck, err := mall_handler.GetCache(d.Token)
		if err != nil {
			return nil
		}
		if timeCheck == "" {
			return errors.New("已过期,请重新申请")
		}
	}

	// 加密密码
	d.Pwd = utils.Md5Uc(d.Pwd)

	// 查询用户数据
	var err error
	find := models.UserInfo{}
	// 根据手机号查ID
	if d.Type == "phone" {
		err = e.Orm.Where("SHA1(CONCAT(user_phone,login_password)) = ?", d.Token).First(&find).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.New("参数有被篡改，请检查")
			}
			return err
		}
	}
	// 根据邮箱查ID
	if d.Type == "email" {
		err = e.Orm.Where("SHA1(CONCAT(user_email,login_password)) = ?", d.Token).First(&find).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.New("参数有被篡改，请检查")
			}
			return err
		}
	}

	// 使用base权限DB
	dsn := mysql.Open(config.DatabasesConfig["base"].Source)
	baseDb, err := gorm.Open(dsn, &gorm.Config{})
	if err != nil {
		return errors.New("数据库连接失败")
	}

	// 查宽表数据
	ctx := e.Orm.Statement.Context.(*gin.Context)
	tenantId := ctx.GetHeader("tenant-id")
	mallUser := &modelsCommon.MallUsers{}
	err = baseDb.Table("common.mall_users").Where("tenant_id", tenantId).Where("user_phone", find.UserPhone).First(mallUser).Error
	if err != nil {
		return err
	}

	// 事务落库
	ucPrefix := global.GetTenantUcDBNameWithDB(e.Orm)
	err = baseDb.Transaction(func(tx *gorm.DB) error {
		// 修改密码
		err = tx.Table(ucPrefix+".user_info").Where("id", find.Id).Update("login_password", d.Pwd).Error
		if err != nil {
			return err
		}

		// 宽表更新
		err = tx.Table("common.mall_users").Where("id = ?", mallUser.Id).Update("login_password", d.Pwd).Error
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

// 根据用户信息生成Token
func (e *UserInfo) genToken(user mall_handler.UserInfo, role mall_handler.RoleInfo) (any, error) {

	// 初始化中间件
	mw, err := middleware.MallAuthInit()
	if err != nil {
		return "", errors.New("系统繁忙，请稍后再试！")
	}

	// 创建Token
	token := jwt.New(jwt.GetSigningMethod(mw.SigningAlgorithm))
	claims := token.Claims.(jwt.MapClaims)
	data := map[string]interface{}{"user": user, "role": role}
	if mw.PayloadFunc != nil {
		for key, value := range mw.PayloadFunc(data) {
			claims[key] = value
		}
	}

	// 设置过期时间
	expire := mw.TimeFunc().Add(mw.Timeout)
	claims["exp"] = expire.Unix()
	claims["orig_iat"] = mw.TimeFunc().Unix()

	// 备注：这里需要约定只走配置里的Key
	tokenString, err := token.SignedString(mw.Key)
	if err != nil {
		return "", err
	}

	// 设置Cookie
	ctx := e.Orm.Statement.Context.(*gin.Context)
	if mw.SendCookie {
		maxage := int(expire.Unix() - time.Now().Unix())
		ctx.SetCookie(
			mw.CookieName,
			tokenString,
			maxage,
			"/",
			mw.CookieDomain,
			mw.SecureCookie,
			mw.CookieHTTPOnly,
		)
	}

	// 返回Token和过期时间
	res := map[string]any{
		"token":  tokenString,
		"expire": expire,
	}

	return res, err
}

// ProxyLogin 代登录前端
func (e *UserInfo) ProxyLogin(d *dto.ProxyLoginReq, c *gin.Context) (any, error) {
	// 多租户处理
	prefix := global.GetTenantUcDBNameWithDB(e.Orm)

	// 验证用户
	find := mall_handler.UserInfo{}
	err := e.Orm.Table(prefix+".user_info").Where("id = ?", d.Code).First(&find).Error
	if err == gorm.ErrRecordNotFound {
		return "", errors.New("参数验证错误[1]")
	}
	if err != nil {
		return "", err
	}

	// 参数校验
	salt := models.UserInfoProxyLoginToken
	user := &models.UserInfo{}
	tenantId := c.GetHeader("tenant-id")
	checkToken := user.GenerateAccessToken(find.Id, tenantId, find.LoginPassword, d.DeCode, salt)
	if checkToken != d.AccessToken {
		return "", errors.New("参数验证错误[2]")
	}

	// 查询角色ID
	userRole := mall_handler.UserRole{}
	err = e.Orm.Table(prefix+".user_role").Where("user_id = ? ", find.Id).First(&userRole).Error
	if err == gorm.ErrRecordNotFound {
		return "", errors.New("该用户还没有角色，请配置[1]")
	}
	if err != nil {
		return "", err
	}

	// 查询角色
	role := mall_handler.RoleInfo{}
	err = e.Orm.Table(prefix+".role_info").Where("id = ? ", userRole.RoleId).Order("id ASC").First(&role).Error
	if err == gorm.ErrRecordNotFound {
		return "", errors.New("该用户还没有角色，请配置[2]")
	}
	if err != nil {
		return "", err
	}

	// 生成token
	res, err := e.genToken(find, role)
	if err != nil {
		return "", err
	}

	// 默认选择仓库
	err = mall_handler.InitWarehouse(c, find.Id)
	if err != nil {
		return nil, err
	}

	// 记录登录日志
	mall_handler.LoginLogToDB(c, "2", "代登录成功", find.UserName, find.Id)

	return res, nil
}

// MiniLogin 小程序
func (e *UserInfo) MiniLogin(d *dto.MiniLoginReq) (any, error) {

	// 校验Header中的租户ID是否存在
	ctx := e.Orm.Statement.Context.(*gin.Context)
	tenantId := ctx.GetHeader("tenant-id")
	if tenantId == "" {
		return "", errors.New("租户ID不存在，请检查！")
	}

	// 验证用户
	find := mall_handler.UserInfo{}
	err := e.Orm.Where("user_phone = ?", d.UserPhone).First(&find).Error
	if err == gorm.ErrRecordNotFound {
		return "", errors.New("该手机号未绑定用户")
	}
	if err != nil {
		return "", err
	}

	// 查询角色ID
	userRole := mall_handler.UserRole{}
	err = e.Orm.Where("user_id = ? ", find.Id).First(&userRole).Error
	if err == gorm.ErrRecordNotFound {
		return "", errors.New("该用户还没有角色，请配置[1]")
	}
	if err != nil {
		return "", err
	}

	// 查询角色
	role := mall_handler.RoleInfo{}
	err = e.Orm.Where("id = ? ", userRole.RoleId).Order("id ASC").First(&role).Error
	if err == gorm.ErrRecordNotFound {
		return "", errors.New("该用户还没有角色，请配置[2]")
	}
	if err != nil {
		return "", err
	}

	// 生成token
	res, err := e.genToken(find, role)
	if err != nil {
		return "", err
	}

	// 默认选择仓库
	c := e.Orm.Statement.Context.(*gin.Context)
	err = mall_handler.InitWarehouse(c, find.Id)
	if err != nil {
		return nil, err
	}

	// 记录登录日志
	mall_handler.LoginLogToDB(c, "2", "小程序登录成功", find.UserName, find.Id)
	return res, nil
}
