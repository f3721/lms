package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-admin/app/admin/models"
	"go-admin/app/admin/service/dto"
	modelsWc "go-admin/app/wc/models"
	dtoWc "go-admin/app/wc/service/admin/dto"
	wcClient "go-admin/common/client/wc"
	"go-admin/common/global"
	"go-admin/common/middleware"
	"go-admin/common/middleware/admin_handler"
	modelsCommon "go-admin/common/models"
	"go-admin/common/utils"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jinzhu/copier"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/config"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	log "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type SysUser struct {
	service.Service
}

// GetPage 获取SysUser列表
func (e *SysUser) GetPage(c *dto.SysUserGetPageReq, p *actions.DataPermission, list *[]dto.SysUserGetPageResp, count *int64) error {
	var err error
	var data models.SysUser

	err = e.Orm.Debug().Preload("Dept").
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}

	// 货主map
	vendorResult := wcClient.ApiByDbContext(e.Orm).GetVendorList(dtoWc.InnerVendorsGetListReq{
		Status: "1",
	})
	vendorResultInfo := &struct {
		response.Response
		Data []modelsWc.Vendors
	}{}
	vendorResult.Scan(vendorResultInfo)
	vendorMap := make(map[int]string, len(vendorResultInfo.Data))
	for _, vendor := range vendorResultInfo.Data {
		vendorMap[vendor.Id] = vendor.NameZh
	}

	// 仓库map
	warehouseResult := wcClient.ApiByDbContext(e.Orm).GetWarehouseList(dtoWc.InnerWarehouseGetListReq{
		Status: "1",
	})
	warehouseResultInfo := &struct {
		response.Response
		Data []modelsWc.Warehouse
	}{}
	warehouseResult.Scan(warehouseResultInfo)
	warehouseMap := make(map[string]string, len(warehouseResultInfo.Data))
	for _, warehouse := range warehouseResultInfo.Data {
		warehouseMap[warehouse.WarehouseCode] = warehouse.WarehouseName
	}

	tmpList := *list
	for i, v := range tmpList {
		temWareHouseCode := utils.Split(v.AuthorityWarehouseId)
		for _, v2 := range temWareHouseCode {
			if _, ok := warehouseMap[v2]; ok {
				tmpList[i].AuthorityWarehouseName = tmpList[i].AuthorityWarehouseName + warehouseMap[v2] + ","
			}
		}
		tmpList[i].AuthorityWarehouseName = strings.TrimRight(tmpList[i].AuthorityWarehouseName, ",")
		temVendorId := utils.SplitToInt(v.AuthorityVendorId)
		for _, v2 := range temVendorId {
			if _, ok := vendorMap[v2]; ok {
				tmpList[i].AuthorityVendorName = tmpList[i].AuthorityVendorName + vendorMap[v2] + ","
			}
		}
		tmpList[i].AuthorityVendorName = strings.TrimRight(tmpList[i].AuthorityVendorName, ",")
	}
	*list = tmpList

	return nil
}

// Get 获取SysUser对象
func (e *SysUser) Get(d *dto.SysUserById, p *actions.DataPermission, model *dto.SysUserGetResp) error {
	var data models.SysUser

	err := e.Orm.Model(&data).Debug().
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("db error: %s", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	// 仓库map
	warehouseResult := wcClient.ApiByDbContext(e.Orm).GetWarehouseList(dtoWc.InnerWarehouseGetListReq{
		Status: "1",
	})
	warehouseResultInfo := &struct {
		response.Response
		Data []modelsWc.Warehouse
	}{}
	warehouseResult.Scan(warehouseResultInfo)
	warehouseMap := make(map[string][]string, len(warehouseResultInfo.Data))
	for _, warehouse := range warehouseResultInfo.Data {
		warehouseMap[warehouse.WarehouseCode] = []string{strconv.Itoa(warehouse.CompanyId), warehouse.WarehouseCode}
	}

	// 仓库调拨 下拉选中项
	model.WarehouseAllocateSelected = [][]string{}
	temAllocateWarehouseCodeArr := utils.Split(model.AuthorityWarehouseAllocateId)
	for _, v2 := range temAllocateWarehouseCodeArr {
		if _, ok := warehouseMap[v2]; ok {
			model.WarehouseAllocateSelected = append(model.WarehouseAllocateSelected, warehouseMap[v2])
		}
	}

	// 公司及仓库 下拉选中项
	model.CompanyWarehouseSelected = [][]string{}
	temCompanyIDArr := utils.Split(model.AuthorityCompanyId)
	temWarehouseCodeArr := utils.Split(model.AuthorityWarehouseId)
	for _, v2 := range temWarehouseCodeArr {
		if _, ok := warehouseMap[v2]; ok {
			model.CompanyWarehouseSelected = append(model.CompanyWarehouseSelected, warehouseMap[v2])
			temCompanyIDArr = utils.RemoveIfExist(temCompanyIDArr, warehouseMap[v2][0])
		}
	}
	for _, v3 := range temCompanyIDArr {
		model.CompanyWarehouseSelected = append(model.CompanyWarehouseSelected, []string{v3})
	}

	return nil
}

func (e *SysUser) UploadSignImg(d *dto.UploadSignImg, p *actions.DataPermission, model *dto.SysUserGetResp) error {
	var data models.SysUser
	err := e.Orm.Model(&data).Debug().
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("db error: %s", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}

	//上传电子签
	err = e.Orm.Model(&data).Where("id", d.Id).Update("sign_img", d.SignImage).Error
	if err != nil {
		return err
	}
	return nil
}

// Insert 创建SysUser对象
func (e *SysUser) Insert(c *dto.SysUserInsertReq) error {
	var err error
	var data models.SysUser
	var i int64
	err = e.Orm.Model(&data).Where("username = ? or phone = ?", c.Username, c.Phone).Count(&i).Error
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	if i > 0 {
		err := errors.New("用户名或手机号已存在！")
		e.Log.Errorf("db error: %s", err)
		return err
	}
	c.Generate(&data)

	// 使用base权限DB
	dsn := mysql.Open(config.DatabasesConfig["base"].Source)
	baseDb, err := gorm.Open(dsn, &gorm.Config{})
	if err != nil {
		return errors.New("数据库连接失败")
	}

	// 事务落库
	err = baseDb.Transaction(func(tx *gorm.DB) error {
		// 保存信息
		adminPrefix := global.GetTenantAdminDBNameWithDB(e.Orm)
		err = tx.Table(adminPrefix + ".sys_user").Create(&data).Error
		if err != nil {
			e.Log.Errorf("db error1: %s", err)
			return err
		}

		// 同步宽表
		adminUser := &modelsCommon.AdminUsers{}
		err = copier.Copy(adminUser, &data)
		adminUser.Password = data.Password
		ctx := e.Orm.Statement.Context.(*gin.Context)
		adminUser.TenantId = ctx.GetHeader("tenant-id")
		if err != nil {
			return err
		}

		err = tx.Table("common.admin_users").Create(adminUser).Error
		if err != nil {
			e.Log.Errorf("db error2: %s", err)
			return err
		}

		return nil
	})
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}

	// 生成日志
	e.createLog(c, models.SysUser{}, data, global.LogTypeCreate)
	return nil
}

// Update 修改SysUser对象
func (e *SysUser) Update(c *dto.SysUserUpdateReq, p *actions.DataPermission) error {
	var err error
	var model models.SysUser
	db := e.Orm.Scopes(
		actions.Permission(model.TableName(), p),
	).First(&model, c.GetId())
	if err = db.Error; err != nil {
		e.Log.Errorf("Service UpdateSysUser error: %s", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")

	}
	oldData := model
	c.Generate(&model)

	// 保证用户名和手机号的唯一性
	if c.Username != "" && oldData.Username != c.Username {
		usernameExit, err := model.UsernameExist(e.Orm, c.Username)
		if err != nil {
			return err
		}
		if usernameExit {
			return fmt.Errorf("用户名[%v]已经存在,请更换", c.Username)
		}
	}
	if c.Phone != "" && oldData.Phone != c.Phone {
		phoneExit, err := model.PhoneExist(e.Orm, c.Phone)
		if err != nil {
			return err
		}
		if phoneExit {
			return fmt.Errorf("手机号[%v]已经存在,请更换", c.Phone)
		}
	}

	// 使用base权限DB
	dsn := mysql.Open(config.DatabasesConfig["base"].Source)
	baseDb, err := gorm.Open(dsn, &gorm.Config{})
	if err != nil {
		return errors.New("数据库连接失败")
	}

	// 查询宽表数据
	ctx := e.Orm.Statement.Context.(*gin.Context)
	tenantId := ctx.GetHeader("tenant-id")
	adminUser := &modelsCommon.AdminUsers{}
	err = baseDb.Table("common.admin_users").Where("tenant_id", tenantId).Where("username", oldData.Username).First(adminUser).Error
	if err != nil {
		return err
	}

	// 事务落库
	baseDb.Transaction(func(tx *gorm.DB) error {
		// 更新用户信息
		adminPrefix := global.GetTenantAdminDBNameWithDB(e.Orm)
		update := tx.Table(adminPrefix+".sys_user").Where("user_id = ?", &model.UserId).Omit("password", "salt").Save(&model)
		if err = update.Error; err != nil {
			e.Log.Errorf("db error: %s", err)
			return err
		}
		if update.RowsAffected == 0 {
			err = errors.New("update userinfo error")
			log.Warnf("db update error")
			return err
		}

		// 更新用户宽表信息
		err = copier.Copy(adminUser, &model)
		if err != nil {
			return err
		}
		err := tx.Table("common.admin_users").Where("id", adminUser.Id).Omit("password", "tenantId").Updates(adminUser).Error
		if err != nil {
			return err
		}
		return nil
	})

	e.createLog(c, oldData, model, global.LogTypeUpdate)
	return nil
}

// UpdateAvatar 更新用户头像
func (e *SysUser) UpdateAvatar(c *dto.UpdateSysUserAvatarReq, p *actions.DataPermission) error {
	var err error
	var model models.SysUser
	db := e.Orm.Scopes(
		actions.Permission(model.TableName(), p),
	).First(&model, c.GetId())
	if err = db.Error; err != nil {
		e.Log.Errorf("Service UpdateSysUser error: %s", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")

	}
	oldData := model
	err = e.Orm.Table(model.TableName()).Where("user_id =? ", c.UserId).Updates(c).Error
	if err != nil {
		e.Log.Errorf("Service UpdateSysUser error: %s", err)
		return err
	}
	e.createLog(c, oldData, model, global.LogTypeUpdate)
	return nil
}

// UpdateStatus 更新用户状态
func (e *SysUser) UpdateStatus(c *dto.UpdateSysUserStatusReq, p *actions.DataPermission) error {
	var err error
	var model models.SysUser
	db := e.Orm.Scopes(
		actions.Permission(model.TableName(), p),
	).First(&model, c.GetId())
	if err = db.Error; err != nil {
		e.Log.Errorf("Service UpdateSysUser error: %s", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")

	}
	oldData := model

	// 使用base权限DB
	dsn := mysql.Open(config.DatabasesConfig["base"].Source)
	baseDb, err := gorm.Open(dsn, &gorm.Config{})
	if err != nil {
		return errors.New("数据库连接失败")
	}

	// 查询宽表数据
	ctx := e.Orm.Statement.Context.(*gin.Context)
	tenantId := ctx.GetHeader("tenant-id")
	adminUser := &modelsCommon.AdminUsers{}
	err = baseDb.Table("common.admin_users").Where("tenant_id", tenantId).Where("username", model.Username).First(adminUser).Error
	if err != nil {
		return err
	}

	// 事务落库
	baseDb.Transaction(func(tx *gorm.DB) error {
		// 更新用户状态
		adminPrefix := global.GetTenantAdminDBNameWithDB(e.Orm)
		err := tx.Debug().Table(adminPrefix+".sys_user").Where("user_id =? ", c.UserId).Updates(c).Error
		if err != nil {
			e.Log.Errorf("db error: %s", err)
			return err
		}

		// 更新用户宽表状态
		err = tx.Table("common.admin_users").Where("id", adminUser.Id).Update("status", c.Status).Error
		if err != nil {
			return err
		}
		return nil
	})

	e.createLog(c, oldData, model, global.LogTypeUpdate)
	return nil
}

// ResetPwd 重置用户密码
func (e *SysUser) ResetPwd(c *dto.ResetSysUserPwdReq, p *actions.DataPermission) error {
	var err error
	var model models.SysUser
	db := e.Orm.Scopes(
		actions.Permission(model.TableName(), p),
	).First(&model, c.GetId())
	if err = db.Error; err != nil {
		e.Log.Errorf("At Service ResetSysUserPwd error: %s", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	oldData := model
	c.Generate(&model)

	// 使用base权限DB
	dsn := mysql.Open(config.DatabasesConfig["base"].Source)
	baseDb, err := gorm.Open(dsn, &gorm.Config{})
	if err != nil {
		return errors.New("数据库连接失败")
	}

	// 查询宽表数据
	ctx := e.Orm.Statement.Context.(*gin.Context)
	tenantId := ctx.GetHeader("tenant-id")
	adminUser := &modelsCommon.AdminUsers{}
	err = baseDb.Table("common.admin_users").Where("tenant_id", tenantId).Where("username", model.Username).First(adminUser).Error
	if err != nil {
		return errors.New("查询宽表数据失败")
	}

	// 事务落库
	baseDb.Transaction(func(tx *gorm.DB) error {
		// 更新用户密码
		adminPrefix := global.GetTenantAdminDBNameWithDB(e.Orm)
		update := tx.Debug().Table(adminPrefix+".sys_user").Where("user_id = ?", &model.UserId).Omit("username", "nick_name", "phone", "role_id", "avatar", "sex").Save(&model)
		if err = update.Error; err != nil {
			e.Log.Errorf("db error: %s", err)
			return err
		}
		if update.RowsAffected == 0 {
			return err
		}

		// 更新用户宽表密码
		err = copier.Copy(adminUser, &model)
		if err != nil {
			return err
		}
		err := tx.Table("common.admin_users").Where("id", adminUser.Id).Omit("username", "phone").Updates(adminUser).Error
		if err != nil {
			return err
		}
		return nil
	})

	e.createLog(c, oldData, model, global.LogTypeUpdate)
	return nil
}

// Remove 删除SysUser | 备注：目前还没有做同步宽表，页面没有删除操作
func (e *SysUser) Remove(c *dto.SysUserById, p *actions.DataPermission) error {
	var err error
	var data models.SysUser

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, c.GetId())
	if err = db.Error; err != nil {
		e.Log.Errorf("Error found in  RemoveSysUser : %s", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}

// UpdatePwd 修改SysUser对象密码
func (e *SysUser) UpdatePwd(id int, oldPassword, newPassword string, p *actions.DataPermission) error {
	var err error

	if newPassword == "" {
		return nil
	}
	c := &models.SysUser{}

	err = e.Orm.Model(c).
		Scopes(
			actions.Permission(c.TableName(), p),
		).Select("UserId", "Username", "Password", "Salt").
		First(c, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("无权更新该数据")
		}
		e.Log.Errorf("db error: %s", err)
		return err
	}
	oldData := c
	var ok bool
	ok, err = pkg.CompareHashAndPassword(c.Password, oldPassword)
	if err != nil {
		e.Log.Errorf("CompareHashAndPassword error, %s", err.Error())
		return err
	}
	if !ok {
		err = errors.New("incorrect Password")
		e.Log.Warnf("user[%d] %s", id, err.Error())
		return err
	}
	c.Password = newPassword

	// 使用base权限DB
	dsn := mysql.Open(config.DatabasesConfig["base"].Source)
	baseDb, err := gorm.Open(dsn, &gorm.Config{})
	if err != nil {
		return errors.New("数据库连接失败")
	}

	// 查询宽表数据
	ctx := e.Orm.Statement.Context.(*gin.Context)
	tenantId := ctx.GetHeader("tenant-id")
	adminUser := &modelsCommon.AdminUsers{}
	err = baseDb.Table("common.admin_users").Where("tenant_id", tenantId).Where("username", c.Username).First(adminUser).Error
	if err != nil {
		return err
	}

	// 事务落库
	baseDb.Transaction(func(tx *gorm.DB) error {
		// 更新用户密码
		adminPrefix := global.GetTenantAdminDBNameWithDB(e.Orm)
		update := tx.Debug().Table(adminPrefix+".sys_user").Where("user_id = ?", id).Select("Password", "Salt").Updates(c)
		if err = update.Error; err != nil {
			e.Log.Errorf("db error: %s", err)
			return err
		}
		if update.RowsAffected == 0 {
			return err
		}

		// 更新用户宽表密码
		err = copier.Copy(adminUser, &c)
		if err != nil {
			return err
		}
		err := tx.Table("common.admin_users").Where("id", adminUser.Id).Omit("username", "phone").Updates(adminUser).Error
		if err != nil {
			return err
		}
		return nil
	})

	e.createLog(map[string]interface{}{
		"id":          id,
		"oldPassword": oldPassword,
		"newPassword": newPassword,
	}, *oldData, *c, global.LogTypeUpdate)
	return nil
}

func (e *SysUser) GetProfile(c *dto.SysUserById, user *models.SysUser, roles *[]models.SysRole, posts *[]models.SysPost) error {
	err := e.Orm.Preload("Dept").First(user, c.GetId()).Error
	if err != nil {
		return err
	}
	err = e.Orm.Find(roles, user.RoleId).Error
	if err != nil {
		return err
	}
	err = e.Orm.Find(posts, user.PostIds).Error
	if err != nil {
		return err
	}

	return nil
}

// 生成日志
func (e *SysUser) createLog(req interface{}, beforeModel models.SysUser, afterModel models.SysUser, logType string) {
	dataLog, _ := json.Marshal(&req)
	beforeDataStr := []byte("")
	if !reflect.DeepEqual(beforeModel, models.SysUser{}) {
		beforeDataStr, _ = json.Marshal(&beforeModel)
	}
	afterDataStr, _ := json.Marshal(&afterModel)
	sysUserLog := models.SysUserLog{
		DataId:       afterModel.UserId,
		Type:         logType,
		Data:         string(dataLog),
		BeforeData:   string(beforeDataStr),
		AfterData:    string(afterDataStr),
		CreateBy:     user.GetUserId(e.Orm.Statement.Context.(*gin.Context)),
		CreateByName: user.GetUserName(e.Orm.Statement.Context.(*gin.Context)),
	}
	_ = sysUserLog.CreateLog("sysUser", e.Orm)
}

// UpdatePermission 修改SysUser权限
func (e *SysUser) UpdatePermission(c *dto.SysUserUpdatePermissionReq) error {
	var err error
	var model models.SysUser
	db := e.Orm.First(&model, c.GetId())
	if err = db.Error; err != nil {
		e.Log.Errorf("Service UpdateSysUser error: %s", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")

	}
	oldData := model
	c.Generate(&model)
	update := e.Orm.Model(&model).Where("user_id = ?", &model.UserId).Select("AuthorityCompanyId", "AuthorityWarehouseId", "AuthorityWarehouseAllocateId", "AuthorityVendorId", "UpdateBy", "UpdateByName").Updates(&model)
	if err = update.Error; err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	if update.RowsAffected == 0 {
		err = errors.New("update user permission error")
		log.Warnf("db update error")
		return err
	}
	e.createLog(c, oldData, model, global.LogTypeUpdate)
	return nil
}

// 后台代登录
func (e *SysUser) ProxyLogin(d *dto.SysUserProxyLoginReq, c *gin.Context) (any, error) {
	// 验证ID是否正确
	user := admin_handler.SysUser{}
	err := e.Orm.Where("user_id = ?", d.UserId).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return "", errors.New("传入的ID有误，请检查")
	}
	if err != nil {
		return "", nil
	}

	// 查询角色信息
	role := admin_handler.SysRole{}
	err = e.Orm.Where("role_id = ? ", user.RoleId).First(&role).Error
	if err == gorm.ErrRecordNotFound {
		return "", errors.New("该账户还未设置角色")
	}
	if err != nil {
		return "", nil
	}
	data := map[string]interface{}{"user": user, "role": role}

	// 初始化中间件
	mw, err := middleware.AdminAuthInit()
	if err != nil {
		return "", errors.New("系统繁忙，请稍后再试！")
	}

	// 创建Token
	token := jwt.New(jwt.GetSigningMethod(mw.SigningAlgorithm))
	claims := token.Claims.(jwt.MapClaims)
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
	if mw.SendCookie {
		maxage := int(expire.Unix() - time.Now().Unix())
		c.SetCookie(
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

	return res, nil
}

// GetAuthorityWarehouseUser 根据 warehouse code 获取相应的仓管员
func (e *SysUser) GetAuthorityWarehouseUser(tx *gorm.DB, warehouseCode string) (sysUsers *[]models.SysUser, err error) {
	sysUsers = &[]models.SysUser{}
	adminPrefix := global.GetTenantAdminDBNameWithDB(tx)
	err = tx.Table(adminPrefix+".sys_user").
		Where("send_email_status = ?", "1").
		Where("authority_warehouse_id like ?", "%"+warehouseCode+"%").Find(sysUsers).Error

	return sysUsers, err
}
