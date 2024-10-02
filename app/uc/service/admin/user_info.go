package admin

import (
	"errors"
	"fmt"
	"go-admin/common/global"
	commonGlobal "go-admin/common/global"
	cModels "go-admin/common/models"
	modelsCommon "go-admin/common/models"
	"go-admin/common/utils"
	"html"
	"strings"

	"github.com/go-admin-team/go-admin-core/sdk/config"
	"github.com/jinzhu/copier"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"github.com/prometheus/common/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"go-admin/app/uc/models"
	"go-admin/app/uc/service/admin/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
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
			cDto.Paginate(c.GetPageSizeNegative(), c.GetPageIndex()),
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

// GetListPage 获取UserInfo列表
func (e *UserInfo) GetListPage(c *dto.UserInfoGetPageReq, p *actions.DataPermission, list *[]*dto.UserInfoGetListPageRes, count *int64) error {
	var err error
	var data models.UserInfo

	roleInfoService := RoleInfo{e.Service}

	roleMap, err := roleInfoService.GetStatusOkMapList()
	log.Info(roleMap)
	err = e.Orm.Model(&data).Debug().
		Scopes(
			func(db *gorm.DB) *gorm.DB {
				if c.RoleId > 0 {
					db.Where("id in (select user_id from user_role where role_id = ? and `deleted_at` IS NULL ) ", c.RoleId)
				}
				if c.FilterCompanyDepartmentId > 0 {
					db.Where("company_department_id in (select id from company_department where (id = ? || f_id = ?) and `deleted_at` IS NULL ) ", c.FilterCompanyDepartmentId, c.FilterCompanyDepartmentId)
				}
				return db
			},
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSizeNegative(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
			// 公司权限筛选
			actions.SysUserPermission(data.TableName(), p, 1),
		).
		Preload("UserRole").
		Preload("CompanyInfo", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, company_name")
		}).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	var companyDepartmentIDs []int
	for _, res := range *list {
		companyDepartmentIDs = append(companyDepartmentIDs, res.CompanyDepartmentId)
	}
	companyDepartmentNames := e.GetDepartmentNamesByIDs(companyDepartmentIDs)

	for _, res := range *list {
		if res.CompanyDepartmentId > 0 {
			res.CompanyDepartmentText = strings.Join(companyDepartmentNames[res.CompanyDepartmentId].Names, ">")
		}
		//	查询用户角色
		for _, userRoleData := range res.UserRole {
			if v, ok := roleMap[userRoleData.RoleId]; ok {
				log.Info(v)
				res.UserRoleText = res.UserRoleText + " " + v.RoleName
			}
		}
	}
	if err != nil {
		e.Log.Errorf("UserInfoService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

func (e *UserInfo) GetDepartmentNamesByIDs(companyDepartmentIDs []int) (names map[int]dto.UserInfoGetDepartmentNames) {
	names = make(map[int]dto.UserInfoGetDepartmentNames)
	log.Info("companyDepartmentIDs")
	log.Info(companyDepartmentIDs)
	if len(companyDepartmentIDs) == 0 {
		return names
	}

	var departments []models.CompanyDepartment
	err := e.Orm.Table("company_department").Where("id IN ?", companyDepartmentIDs).Find(&departments).Error
	if err != nil {
		return names
	}

	var tipIds []int
	for _, department := range departments {
		tipIds = append(tipIds, department.TopId)
	}
	var allDepartments []models.CompanyDepartment
	err = e.Orm.Table("company_department").Where("top_id IN ?", tipIds).Order("top_id asc,level desc").Find(&allDepartments).Error
	if err != nil {
		return names
	}

	allDepartmentsMap := make(map[int]models.CompanyDepartment)
	for _, department := range allDepartments {
		allDepartmentsMap[department.Id] = department
	}
	//fid := -1
	var nameArr []string
	var nameArrIds []int
	forUserDepartments := models.CompanyDepartment{}
	for _, userDepartmentId := range companyDepartmentIDs {
		nameArr = []string{}
		nameArrIds = []int{}
		userDepartments := allDepartmentsMap[userDepartmentId]
		level := userDepartments.Level
		for i := level; i != 0; i-- {
			if i == level {
				forUserDepartments = userDepartments
			} else {
				forUserDepartments = allDepartmentsMap[userDepartments.FId]
			}
			nameArr = append(nameArr, forUserDepartments.Name)
			nameArrIds = append(nameArrIds, forUserDepartments.Id)
		}

		utils.ReverseSlice(nameArr)
		utils.ReverseSlice(nameArrIds)
		names[userDepartmentId] = dto.UserInfoGetDepartmentNames{
			Names: nameArr,
			Ids:   nameArrIds,
		}
	}

	return names
}

// GetInfo 获取UserInfo对象
func (e *UserInfo) GetInfo(d *dto.UserInfoGetReq, p *actions.DataPermission, model *dto.UserInfoGetRes) error {
	var data models.UserInfo

	err := e.Orm.Model(&data).Debug().
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		Preload("UserRole").
		Find(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetUserInfo error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	companyDepartmentNames := e.GetDepartmentNamesByIDs([]int{
		model.CompanyDepartmentId,
	})

	model.UserDepartments = companyDepartmentNames[model.CompanyDepartmentId]
	model.LoginPassword = ""
	return nil
}

// Get 获取UserInfo对象
func (e *UserInfo) Get(d *dto.UserInfoGetReq, p *actions.DataPermission, model *models.UserInfo) error {
	var data models.UserInfo

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
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

// GetUserCompanyByUserId 获取用户所属公司名
func (e *UserInfo) GetUserCompanyByUserId(userId int) (data *dto.UserInfoGetUserCompanyInfo, err error) {
	data = &dto.UserInfoGetUserCompanyInfo{}
	err = e.Orm.Model(&models.UserInfo{}).
		Where("user_info.id = ?", userId).
		Joins("join company_info on company_info.id = user_info.company_id").
		Select("company_info.id company_id,company_info.company_name").
		First(&data).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		e.Log.Errorf("db error:%s", err)
		return
	}
	return
}

// Insert 创建UserInfo对象
func (e *UserInfo) Insert(c *dto.UserInfoInsertReq, isImport bool) (errs []error) {
	var err error
	var data models.UserInfo
	var saveData models.UserInfo
	c.Generate(&saveData)

	errs = e.SaveValidate(&data, &saveData, isImport)
	if errs != nil {
		return errs
	}

	// 使用base权限DB
	dsn := mysql.Open(config.DatabasesConfig["base"].Source)
	baseDb, err := gorm.Open(dsn, &gorm.Config{})
	if err != nil {
		return append(errs, errors.New("数据库连接失败"))
	}
	tx := baseDb.Begin()

	// 保存用户
	ucPrefix := global.GetTenantUcDBNameWithDB(e.Orm)
	err = tx.Table(ucPrefix + ".user_info").Create(&saveData).Error
	if err != nil {
		tx.Rollback()
		e.Log.Errorf("UserInfoService Insert error:%s \r\n", err)
		return append(errs, err)
	}

	if len(c.UserRole) > 0 {
		userRoleService := UserRole{e.Service}
		err = userRoleService.UserRoleSave(e.Orm, c.UserRole, saveData.Id)
		if err != nil {
			tx.Rollback()
			return append(errs, err)
		}
	}

	operateLogsService := OperateLogs{e.Service}
	_ = operateLogsService.AddLog(saveData.Id, data, saveData, models.UserInfoModel, commonGlobal.LogTypeCreate, c.CreateBy, c.CreateByName)

	// 同步宽表
	mallUser := &modelsCommon.MallUsers{}
	err = copier.Copy(mallUser, &saveData)
	mallUser.LoginPassword = saveData.LoginPassword
	ctx := e.Orm.Statement.Context.(*gin.Context)
	mallUser.TenantId = ctx.GetHeader("tenant-id")
	if err != nil {
		tx.Rollback()
		return append(errs, err)
	}

	err = tx.Table("common.mall_users").Create(mallUser).Error
	if err != nil {
		tx.Rollback()
		e.Log.Errorf("新增同步宽表错误:%s \r\n", err)
		return append(errs, err)
	}

	tx.Commit()

	//* 生成预算方法
	//x @param int $dep_id
	//* @param int $type 1=>初次新增或预算调整，2=>部门预算和人均预算 模式相互转换
	var departmentBudgetM = models.DepartmentBudget{}
	err = departmentBudgetM.Generate(e.Orm, 0, saveData.Id, 1)
	return nil
}

// Update 修改UserInfo对象
func (e *UserInfo) Update(c *dto.UserInfoUpdateReq, p *actions.DataPermission, isImport bool) (errs []error) {
	var err error
	var data = models.UserInfo{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())

	saveData := data
	c.Generate(&saveData)
	errs = e.SaveValidate(&data, &saveData, isImport)
	if errs != nil {
		return errs
	}

	// 使用base权限DB
	dsn := mysql.Open(config.DatabasesConfig["base"].Source)
	baseDb, err := gorm.Open(dsn, &gorm.Config{})
	if err != nil {
		return append(errs, errors.New("数据库连接失败"))
	}

	// 查宽表数据
	ctx := e.Orm.Statement.Context.(*gin.Context)
	tenantId := ctx.GetHeader("tenant-id")
	mallUser := &modelsCommon.MallUsers{}
	err = baseDb.Table("common.mall_users").Where("tenant_id", tenantId).Where("user_phone", data.UserPhone).First(mallUser).Error
	if err != nil {
		return append(errs, errors.New("系统繁忙,请稍后再试[1]"))
	}

	tx := baseDb.Begin()

	// 编辑用户信息
	ucPrefix := global.GetTenantUcDBNameWithDB(e.Orm)
	dbRes := tx.Table(ucPrefix + ".user_info").Save(&saveData)
	if err = dbRes.Error; err != nil {
		tx.Rollback()
		e.Log.Errorf("UserInfoService Save error:%s \r\n", err)
		return append(errs, errors.New("222"))
	}
	if dbRes.RowsAffected == 0 {
		tx.Rollback()
		return append(errs, errors.New("无权更新该数据"))
	}

	if len(c.UserRole) > 0 {
		userRoleService := UserRole{e.Service}
		err = userRoleService.UserRoleSave(e.Orm, c.UserRole, saveData.Id)
		if err != nil {
			tx.Rollback()
			return append(errs, errors.New("csc"))
		}
	}

	if saveData.LoginPassword != "" {
		// 修改密码日志
		_ = e.UpdatePasswordLog(saveData.Id, user.GetUserName(e.Orm.Statement.Context.(*gin.Context)))
	}
	operateLogsService := OperateLogs{e.Service}
	_ = operateLogsService.AddLog(saveData.Id, data, saveData, models.UserInfoModel, commonGlobal.LogTypeUpdate, c.UpdateBy, c.UpdateByName)

	// 同步宽表
	err = copier.Copy(mallUser, &saveData)
	if err != nil {
		tx.Rollback()
		return append(errs, err)
	}

	err = tx.Table("common.mall_users").Where("id = ?", mallUser.Id).Save(mallUser).Error
	if err != nil {
		tx.Rollback()
		e.Log.Errorf("更新同步宽表错误:%s \r\n", err)
		return append(errs, err)
	}

	tx.Commit()

	if data.CompanyDepartmentId != saveData.CompanyDepartmentId {
		//* 生成预算方法
		//x @param int $dep_id
		//* @param int $type 1=>初次新增或预算调整，2=>部门预算和人均预算 模式相互转换
		var departmentBudgetM = models.DepartmentBudget{}
		err = departmentBudgetM.Generate(e.Orm, 0, saveData.Id, 1)
	}
	return nil
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
		if !utils.ValidatePassword(saveData.LoginPassword) {
			errs = append(errs, errors.New("密码格式错误"))
			if !isAllErr {
				return
			}
		}
		// 密码md5处理
		decodedPassword := html.UnescapeString(saveData.LoginPassword)
		newPassword := utils.Md5Uc(decodedPassword)
		saveData.LoginPassword = newPassword
	}

	if data.CompanyId != saveData.CompanyId {
		//公司校验
		companyInfoService := CompanyInfo{e.Service}
		isExists, _ := companyInfoService.CompanyIsExists(saveData.CompanyId)
		if !isExists {
			errs = append(errs, errors.New("关联公司不存在"))
			if !isAllErr {
				return
			}
		}
	}
	if data.LoginName != saveData.LoginName {
		//loginName的验证逻辑
		if !utils.ValidateLoginName(saveData.LoginName) {
			errs = append(errs, errors.New("登录账户要求:8-20位,由大写字母、小写字母、数字和英文字符(除空格外)至少二种组成"))
			if !isAllErr {
				return
			}
		}
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
	if saveData.UserEmail != "" && !strings.EqualFold(data.UserEmail, saveData.UserEmail) {
		if !utils.ValidateEmailFormat(saveData.UserEmail) {
			errs = append(errs, errors.New("邮箱格式不正确"))
			if !isAllErr {
				return
			}
		}
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

// UpdatePassword 修改密码
func (e *UserInfo) UpdatePassword(c *dto.UserInfoUpdatePassword, p *actions.DataPermission) error {
	var err error
	var data = models.UserInfo{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.Id)

	// 使用base权限DB
	dsn := mysql.Open(config.DatabasesConfig["base"].Source)
	baseDb, err := gorm.Open(dsn, &gorm.Config{})
	if err != nil {
		return err
	}

	// 查宽表数据
	ctx := e.Orm.Statement.Context.(*gin.Context)
	tenantId := ctx.GetHeader("tenant-id")
	mallUser := &modelsCommon.MallUsers{}
	err = baseDb.Table("common.mall_users").Where("tenant_id", tenantId).Where("user_phone", data.UserPhone).First(mallUser).Error
	if err != nil {
		return err
	}

	tx := baseDb.Begin()

	// 密码md5处理
	decodedPassword := html.UnescapeString(c.LoginPassword)
	newPassword := utils.Md5Uc(decodedPassword)

	data.LoginPassword = newPassword

	ucPrefix := global.GetTenantUcDBNameWithDB(e.Orm)
	dbRes := tx.Table(ucPrefix + ".user_info").Save(&data)
	if err = dbRes.Error; err != nil {
		tx.Rollback()
		e.Log.Errorf("UserInfoService Save error:%s \r\n", err)
		return err
	}
	if dbRes.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("无权更新该数据")
	}

	// 修改密码日志
	err = e.UpdatePasswordLog(data.Id, user.GetUserName(e.Orm.Statement.Context.(*gin.Context)))
	if err != nil {
		tx.Rollback()
		return err
	}

	// 同步宽表
	err = tx.Table("common.mall_users").Where("id = ?", mallUser.Id).Update("login_password", newPassword).Error
	if err != nil {
		tx.Rollback()
		e.Log.Errorf("更新同步宽表错误:%s \r\n", err)
		return err
	}

	tx.Commit()
	return nil
}

// UpdateUserPhone 修改手机号
func (e *UserInfo) UpdateUserPhone(c *dto.UserInfoUpdateUserPhone, p *actions.DataPermission) error {
	var err error
	var data = models.UserInfo{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.Id)
	saveData := data

	saveData.UserPhone = c.UserPhone

	errs := e.SaveValidate(&data, &saveData, false)
	if errs != nil {
		return errs[0]
	}

	// 使用base权限DB
	dsn := mysql.Open(config.DatabasesConfig["base"].Source)
	baseDb, err := gorm.Open(dsn, &gorm.Config{})
	if err != nil {
		return err
	}

	// 查宽表数据
	ctx := e.Orm.Statement.Context.(*gin.Context)
	tenantId := ctx.GetHeader("tenant-id")
	mallUser := &modelsCommon.MallUsers{}
	err = baseDb.Table("common.mall_users").Where("tenant_id", tenantId).Where("user_phone", data.UserPhone).First(mallUser).Error
	if err != nil {
		return err
	}

	tx := baseDb.Begin()

	// 更新手机号
	ucPrefix := global.GetTenantUcDBNameWithDB(e.Orm)
	dbRes := tx.Table(ucPrefix + ".user_info").Save(&saveData)
	if err = dbRes.Error; err != nil {
		e.Log.Errorf("UserInfoService Save error:%s \r\n", err)
		tx.Rollback()
		return err
	}
	if dbRes.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}

	operateLogsService := OperateLogs{e.Service}
	_ = operateLogsService.AddLog(data.Id, data, saveData, models.UserInfoModel, models.UserInfoOperationUpdateUserPhone, c.UpdateBy, c.UpdateByName)

	// 更新宽表手机号
	err = tx.Table("common.mall_users").Where("id = ?", mallUser.Id).Update("user_phone", c.UserPhone).Error
	if err != nil {
		tx.Rollback()
		e.Log.Errorf("更新同步宽表错误:%s \r\n", err)
		return err
	}

	tx.Commit()
	return nil
}

// UpdatePasswordLog 修改密码日志
func (e *UserInfo) UpdatePasswordLog(userId int, recordName string) error {
	userPasswordChangeLogModels := models.UserPasswordChangeLog{}
	err := userPasswordChangeLogModels.AddLog(e.Orm, userId, recordName)
	if err != nil {
		return err
	}
	return nil
}

// ProxyLogin 代登录
func (e *UserInfo) ProxyLogin(c *dto.UserInfoProxyLoginReq, p *actions.DataPermission) (res *dto.UserInfoProxyLoginRes, err error) {
	//// 实现代登录逻辑
	//sellPeople := c.SellPeople
	//if len(sellPeople) == 0 {
	//	err = errors.New("没有绑定销售，绑定之后才可登录")
	//	return
	//}

	if c.UserID == 0 {
		err = errors.New("参数错误")
		return
	}

	userData := models.UserInfo{}
	err = e.Get(&dto.UserInfoGetReq{Id: c.UserID}, p, &userData)
	if err != nil {
		return
	}

	CompanyInfoModels := models.CompanyInfo{}
	companyDomain, err := CompanyInfoModels.GetPCompanyDomainByID(e.Orm, userData.CompanyId)
	if err != nil {
		return
	}

	// 生产要登录的token
	token := "EHSY20171227$@#&*(123"
	systemId := user.GetUserName(e.Orm.Statement.Context.(*gin.Context))

	tenantId := e.Orm.Statement.Context.(*gin.Context).GetHeader("tenant-id")

	accessToken := userData.GenerateAccessToken(userData.Id, tenantId, userData.LoginPassword, systemId, token)

	customerWWWLoginURL := ""
	//customerWWWLoginURLs := map[string]string{
	//	"dev":     "http://eis-sxyz-t.ehsy.com",
	//	"test":    "https://mall-test.lionmile.com",
	//	"staging": "http://eis-staging.ehsy.com",
	//	"pro":     "https://mall.lionmile.com",
	//}
	//
	devModel := config.ApplicationConfig.Mode // 假设存在DEV_MODEL变量
	//if url, ok := customerWWWLoginURLs[devModel]; ok {
	//	customerWWWLoginURL = url
	//}

	if len(companyDomain) > 0 {
		customerWWWLoginURL = "http://" + companyDomain
	} else {
		customerWWWLoginURL = utils.GetHostUrl(0)
	}

	redirectURL := fmt.Sprintf("%s/#/login?code=%d&deCode=%s&accessToken=%s&tenantId=%s",
		customerWWWLoginURL, userData.Id, systemId, accessToken, tenantId)

	res = &dto.UserInfoProxyLoginRes{
		RedirectURL: redirectURL,
		DevModel:    devModel,
	}
	return
}

// ImportValidate 导入校验
func (e *UserInfo) ImportValidate(req *dto.UserInfoImportData) (errs []error) {

	companyInfoService := CompanyInfo{e.Service}
	companyDepartmentService := CompanyDepartment{e.Service}

	fId := 0
	level := 1
	if req.UserName == "" {
		errs = append(errs, errors.New("用户姓名不能为空"))
	}
	if req.UserPhone == "" {
		errs = append(errs, errors.New("手机号不能为空"))
	}
	if req.LoginPassword == "" {
		errs = append(errs, errors.New("密码不能为空"))
	}
	if req.CompanyName == "" {
		errs = append(errs, errors.New("公司不能为空"))
	} else {
		companyInfo, _ := companyInfoService.CompanyByName(req.CompanyName)
		if companyInfo != nil && companyInfo.Id == 0 {
			errs = append(errs, errors.New("公司不存在"))
		}
		req.CompanyId = companyInfo.Id
	}

	if req.FDepartmentName != "" {
		// 校验父级部门
		department := models.CompanyDepartment{}
		_ = companyDepartmentService.GetByNameLevel(&department, req.CompanyId, req.FDepartmentName, 0, 0)
		if department.Id == 0 {
			errs = append(errs, errors.New("父级部门不存在"))
			return
		}
		if department.Level != 1 {
			errs = append(errs, errors.New("父级上级部门不是一级部门！"))
			return
		}
		level = 2
		fId = department.Id
	}

	// 根据部门名称和父级部门id获取部门信息
	if req.DepartmentName != "" {
		department := models.CompanyDepartment{}
		_ = companyDepartmentService.GetByNameLevel(&department, req.CompanyId, req.DepartmentName, fId, level)
		req.DepartmentId = department.Id
	}

	if errs != nil {
		return
	}

	return
}

// Import 导入
func (e *UserInfo) Import(req *dto.UserInfoImportData, roleList []*models.RoleInfo) []error {
	if roleList == nil {
		roleInfoService := RoleInfo{e.Service}
		roleList, _ = roleInfoService.GetStatusOkList()
	}

	resErrs := e.ImportValidate(req)
	// 判断角色是否存在
	userRoles := strings.Split(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(req.UserRole, "，", ","), "|", ","), ";", ","), "、", ","), ",")
	// 遍历判断角色是否存在
	var userRoleIds []int
	for _, role := range userRoles {
		roleValidate := false
		for _, info := range roleList {
			if role == info.RoleName {
				roleValidate = true
				userRoleIds = append(userRoleIds, info.Id)
			}
		}
		if !roleValidate {
			resErrs = append(resErrs, errors.New("角色【"+role+"】错误"))
		}
	}
	req.UserRoleList = userRoleIds
	if len(resErrs) == 0 {
		errs := e.Insert(&dto.UserInfoInsertReq{
			UserEmail:           req.UserEmail,
			LoginName:           req.LoginName,
			UserPhone:           req.UserPhone,
			LoginPassword:       req.LoginPassword,
			UserName:            req.UserName,
			UserStatus:          1,
			CompanyId:           req.CompanyId,
			Telephone:           req.Telephone,
			CompanyDepartmentId: req.DepartmentId,
			UserRole:            req.UserRoleList,
			CanLogin:            1,
			ControlBy: cModels.ControlBy{
				CreateBy:     user.GetUserId(e.Orm.Statement.Context.(*gin.Context)),
				CreateByName: user.GetUserName(e.Orm.Statement.Context.(*gin.Context)),
			},
		}, true)
		if len(errs) > 0 {
			// 创建一个新的切片，容量为两个切片的长度之和
			resultErrs := make([]error, len(resErrs)+len(errs))
			// 复制slice1的元素到result
			copy(resultErrs, resErrs)
			// 复制slice2的元素到result
			copy(resultErrs[len(resErrs):], errs)
			resErrs = resultErrs
		}
	}
	if len(resErrs) > 0 {
		return resErrs
	}
	return resErrs
}
