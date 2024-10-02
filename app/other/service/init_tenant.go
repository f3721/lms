package service

import (
	"errors"
	"fmt"
	modelsAdmin "go-admin/app/admin/models"
	"go-admin/app/other/service/dto"
	modelsUc "go-admin/app/uc/models"
	"go-admin/common/database"
	"go-admin/common/global"
	"go-admin/common/models"
	"go-admin/common/utils"
	"golang.org/x/crypto/bcrypt"
	"reflect"
	"strings"
	"time"

	"github.com/go-admin-team/go-admin-core/sdk/config"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"github.com/jinzhu/copier"
	"github.com/samber/lo"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type InitTenant struct {
	service.Service
}

// Insert 创建SysDept对象
func (e *InitTenant) InitTenant(c *dto.InitTenantReq) error {
	//tenant, err := global.GetTenant(global.EncryptTenantId(c.Id))
	// 明文密码加密存储
	var err error
	var hash []byte
	if hash, err = bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost); err != nil {
		return err
	} else {
		c.Password = string(hash)
	}
	if err != nil {
		return err
	}
	var tenant models.SystemTenant
	tenants := database.GetTenantList(database.GetLubanToken())
	if len(tenants) <= 0 {
		return errors.New("暂无驿站租户")
	}
	for _, systemTenant := range tenants {
		if systemTenant.ID == c.Id {
			tenant = systemTenant
			break
		}
	}
	if tenant.Name == "" {
		return errors.New("租户不存在")
	}
	if tenant.DatabaseInitOk == 1 {
		return errors.New("租户数据库已初始化，请勿重复操作")
	}

	baseDbPrefix := "base"
	//syncSource := "ea_yzwl_create:JvuY!0Z3k9618p7@tcp(652050ca7b6b4ed08b104bd3fec4eb13in01.internal.cn-east-3.mysql.rds.myhuaweicloud.com:3306)/base_admin?charset=utf8mb4&parseTime=True&loc=Local&timeout=1000ms"
	baseDb, err := gorm.Open(mysql.Open(config.DatabasesConfig[baseDbPrefix].Source), &gorm.Config{})
	//baseDb, err := gorm.Open(mysql.Open(syncSource), &gorm.Config{})
	if err != nil {
		return errors.New("数据库连接失败")
	}
	//newDbSource := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v_%v?charset=utf8mb4&parseTime=True&loc=Local&timeout=1000ms",
	//	tenant.DatabaseUsername,
	//	tenant.DatabasePassword,
	//	tenant.DatabaseHost,
	//	tenant.DatabasePort,
	//	tenant.TenantDBPrefix(),
	//	"admin")
	//_, err = gorm.Open(mysql.Open(newDbSource), &gorm.Config{})
	//if err != nil {
	//	return errors.New("租户数据库连接失败")
	//}

	// 建表
	dbSuffix := []string{"admin", "oc", "uc", "pc", "wc"}
	for _, suffix := range dbSuffix {
		// 从base数据库中获取所有表的表名
		var tables []string
		err = baseDb.Raw("SHOW TABLES from " + baseDbPrefix + "_" + suffix).Scan(&tables).Error
		if err != nil {
			return errors.New(baseDbPrefix + "_" + suffix + "数据库获取表失败:" + err.Error())
		}

		baseDb.Exec(fmt.Sprintf("DROP TABLE IF Exists %s_%s.%s", tenant.DatabaseName, suffix, "sys_casbin_rule"))

		for _, table := range tables {
			err = baseDb.Exec(fmt.Sprintf("CREATE TABLE %s_%s.%s LIKE %s_%s.%s", tenant.DatabaseName, suffix, table, baseDbPrefix, suffix, table)).Error
			if err != nil {
				return errors.New(baseDbPrefix + "_" + suffix + "数据库创建表失败:" + err.Error())
			}

			//err = baseDb.Exec(fmt.Sprintf("INSERT INTO %s_%s.%s SELECT * FROM %s_%s.%s", tenant.DatabaseName, suffix, table, baseDbPrefix, suffix, table)).Error
			//if err != nil {
			//	return errors.New(baseDbPrefix + "_" + suffix + "." + table + "数据同步失败:" + err.Error())
			//}
		}
	}

	// 初始化数据
	copyTableMap := map[string][]string{
		"admin": {"log_types", "sys_api", "sys_config", "sys_dict_type", "sys_dict_data", "sys_menu", "sys_menu_api_rule", "sys_role_menu"},
		"oc":    {"log_types"},
		"pc":    {"log_types", "attribute_config", "uommaster"},
		"uc":    {"operate_log_types", "manage_menu", "role_info"},
		"wc":    {"operate_log_types", "region"},
	}
	for copySuffix, copyTables := range copyTableMap {
		for _, copyTable := range copyTables {
			err = baseDb.Exec(fmt.Sprintf("INSERT INTO %s_%s.%s SELECT * FROM %s_%s.%s", tenant.DatabaseName, copySuffix, copyTable, baseDbPrefix, copySuffix, copyTable)).Error
			if err != nil {
				return errors.New(baseDbPrefix + "_" + copySuffix + "." + copyTable + "数据同步失败:" + err.Error())
			}
		}
	}

	copySqls := []string{
		"INSERT INTO " + tenant.DatabaseName + "_admin.sys_role (role_id, role_name, status, role_key, role_sort, flag, remark, admin, data_scope, create_by, update_by, created_at, updated_at, deleted_at, create_by_name, update_by_name) VALUES(1, '系统管理员', '2', 'admin', 1, '', '', 1, '', 1, 1, '2021-05-13 19:56:37.913000000', '2021-05-13 19:56:37.913000000', NULL, '', '')",
		"INSERT INTO " + tenant.DatabaseName + "_admin.sys_user (user_id, username, password, nick_name, nick_name_en, phone, role_id, salt, avatar, sex, email, dept_id, post_id, telephone, fax, remark, status, authority_company_id, authority_warehouse_id, authority_warehouse_allocate_id, authority_vendor_id, create_by, update_by, created_at, updated_at, deleted_at, create_by_name, update_by_name) VALUES(1, '" + c.UserName + "', '" + c.Password + "', '系统管理员', '', '13818888888', 1, '', '', '1', '1@qq.com', 1, 1, '', '', '', '2', '1', '', '','', 1, 1, '2021-05-13 19:56:37.914000000', '2023-07-12 18:21:49.416000000', NULL, '', '')",
		"INSERT INTO common.admin_users (tenant_id, username, password, phone, status, create_by, update_by, created_at, updated_at, deleted_at, create_by_name, update_by_name) VALUES ('" + global.EncryptTenantId(tenant.TenantId) + "', '" + c.UserName + "', '" + c.Password + "', '13818888888', '2', 1, 1, '2021-05-13 19:56:37', '2023-09-28 17:16:10', null, '', '');",
		"INSERT INTO " + tenant.DatabaseName + "_uc.company_info (id, company_status, company_name, company_nature, company_type, is_punchout, is_eas, is_eis, parent_id, address, tax_no, bank_name, bank_account, company_logo, company_medium_logo, company_small_logo, industry_id, is_send_email, show_easflow, theme, `domain`, payment_methods, site_title, site_css, company_levels, login_type, pid, preso_term, country_id, province_id, city_id, area_id, pdf_type, is_tax_price, service_fee, order_eamil_set, approve_email_type, reconciliation_day, after_auto_audit, order_auto_confirm, created_at, updated_at, deleted_at, create_by, create_by_name, update_by, update_by_name) VALUES(1, 1, '" + tenant.Name + "', 2, 1, 0, 0, 0, 0, '', '', '', '', '', '', '', 0, 0, 0, 'orange', '', '', '', '', 1, 0, 0, 15, 1, 0, 0, 0, 0, 0, 0.00, '', 1, 0, ' ', 0, '2023-05-15 16:26:48', '2023-07-04 15:47:33', NULL, 1, '', 1, 'admin')",
		// 隐藏菜单： 开发工具、系统工具、定时任务、接口管理、字典数据、参数管理、字典管理、岗位管理、部门管理、菜单管理
		"UPDATE " + tenant.DatabaseName + "_admin.sys_menu set visible = 1 where menu_id in (60, 537, 459, 528, 59, 62, 58, 57, 56, 51) ",
	}
	for _, copySql := range copySqls {
		err = baseDb.Exec(copySql).Error
		if err != nil {
			return errors.New("初始数据同步失败:" + err.Error())
		}
	}
	return nil
}

// 去除租户名查询tenant-id
func (e *InitTenant) LoginNew(req *dto.Login) (any, error) {

	// 查询所有租户信息
	tenants := global.GetTenants()
	if len(tenants) <= 0 {
		return "", errors.New("暂无租户")
	}

	// 生成迁移SQL | 待发布前删除
	if req.Gen == "240" {
		sql, err := e.MigrateUsers(tenants)
		if err != nil {
			return "", err
		}
		return sql, nil
	}

	// 后台登录
	if req.Type == "0" {
		// 查询数据
		var query []*models.AdminUsers
		unionSql := `
		select * from common.admin_users where phone = ? and status = '2' and deleted_at is null
		union 
		select * from common.admin_users where username = ? and status = '2' and deleted_at is null
		`
		err := e.Orm.Raw(unionSql, req.Username, req.Username).Find(&query).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误，请检查！")
		}
		if err != nil {
			return nil, err
		}

		// 校验密码
		var infos []*models.AdminUsers = lo.Filter(query, func(item *models.AdminUsers, _ int) bool {
			match, err := pkg.CompareHashAndPassword(item.Password, req.Password)
			if err != nil {
				return false
			}
			return match
		})
		if len(infos) == 0 {
			return nil, errors.New("用户名或密码错误，请检查！")
		}

		// 补全&去除密码
		for _, item := range infos {
			item.Password = ""
			item.TenantName = tenants[item.TenantId].Name
		}

		return infos, nil
	}

	// 商城登录
	if req.Type == "1" {
		// 查询数据 | 这里去掉了用户名称
		var query []*models.MallUsers
		unionSql := `
		select * from common.mall_users where user_phone = ? and user_status = '1' and deleted_at is null
		union 
		select * from common.mall_users where user_email = ? and user_status = '1' and deleted_at is null
		union 
		select * from common.mall_users where login_name = ? and user_status = '1' and deleted_at is null
		`
		err := e.Orm.Raw(unionSql, req.Username, req.Username, req.Username).Find(&query).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误，请检查！")
		}
		if err != nil {
			return nil, err
		}

		// 校验密码
		var infos []*models.MallUsers = lo.Filter(query, func(item *models.MallUsers, _ int) bool {
			inputMd5Pwd := utils.Md5Uc(req.Password)
			return inputMd5Pwd == item.LoginPassword
		})
		if len(infos) == 0 {
			return nil, errors.New("用户名或密码错误，请检查！")
		}

		// 补全&去除密码
		for _, item := range infos {
			item.LoginPassword = ""
			item.TenantName = tenants[item.TenantId].Name
		}

		return infos, nil
	}

	return nil, errors.New("Type传参错误")
}

// 去除租户名查询tenant-id
func (e *InitTenant) LoginNewMini(req *dto.LoginMini) (any, error) {

	// 查询所有租户信息
	tenants := global.GetTenants()
	if len(tenants) <= 0 {
		return "", errors.New("暂无租户")
	}

	// 商城登录
	var query []*models.MallUsers
	unionSql := `
	select * from common.mall_users where user_phone = ? and user_status = '1' and deleted_at is null
	`
	err := e.Orm.Raw(unionSql, req.Username).Find(&query).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("该用户不存在，请检查！")
	}
	if err != nil {
		return nil, err
	}

	// 校验密码
	infos := query

	// 补全&去除密码
	for _, item := range infos {
		item.LoginPassword = ""
		item.TenantName = tenants[item.TenantId].Name
	}

	return infos, nil
}

// 生成执行SQL
func (e *InitTenant) MigrateUsers(tenants map[string]*models.SystemTenant) (any, error) {
	// 链接DB
	baseDbPrefix := "base"
	db, err := gorm.Open(mysql.Open(config.DatabasesConfig[baseDbPrefix].Source), &gorm.Config{})
	if err != nil {
		return "", err
	}

	// 生成迁移SQL
	insertSqls := []string{}

	// -------admin----------
	for tenantId, item := range tenants {
		var sysUser []modelsAdmin.SysUser
		err = db.Debug().Table(item.DatabaseName + "_admin.sys_user").Find(&sysUser).Error
		if err != nil {
			return "", err
		}

		for _, user := range sysUser {
			// 赋值单条数据
			adminUsers := models.AdminUsers{}
			err := copier.Copy(&adminUsers, &user)
			if err != nil {
				return "", err
			}
			adminUsers.TenantId = tenantId

			// 反射生成SQL
			sql := e.genInsertSQL(adminUsers)
			insertSqls = append(insertSqls, sql)
		}
	}

	// -------mall----------
	for tenantId, item := range tenants {
		var mallUser []modelsUc.UserInfo
		err = db.Debug().Table(item.DatabaseName + "_uc.user_info").Find(&mallUser).Error
		if err != nil {
			return "", err
		}

		for _, user := range mallUser {
			// 赋值单条数据
			mallUsers := models.MallUsers{}
			err := copier.Copy(&mallUsers, &user)
			if err != nil {
				return "", err
			}
			mallUsers.TenantId = tenantId

			// 反射生成SQL
			sql := e.genInsertSQL(mallUsers)
			insertSqls = append(insertSqls, sql)
		}
	}

	resSql := strings.Join(insertSqls, ";")
	return resSql, nil
}

// 生成迁移Insert
func (e *InitTenant) genInsertSQL(data interface{}) string {
	// 获取结构体类型信息
	val := reflect.ValueOf(data)
	typ := val.Type()

	// 准备要插入的列和值
	var columns []string
	var values []string

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		column := fieldType.Tag.Get("db")
		if column == "id" || column == "-" {
			continue
		}
		// 转换字段值为字符串，如果是零值则使用空字符串
		var value string
		switch field.Interface().(type) {
		case int:
			value = fmt.Sprintf("%d", field.Interface())
		case int64:
			value = fmt.Sprintf("%d", field.Interface())
		case time.Time:
			if !field.Interface().(time.Time).IsZero() {
				value = fmt.Sprintf("'%s'", field.Interface().(time.Time).Format("2006-01-02 15:04:05"))
			} else {
				value = "null"
			}
		default:
			// 处理其他类型
			value = fmt.Sprintf("'%s'", field.Interface())
		}

		columns = append(columns, column)
		values = append(values, value)
	}

	// 表名
	tableName := ""
	tableNameVale := val.MethodByName("TableName")
	if !tableNameVale.IsValid() {
		fmt.Println("TableName方法不存在")
	}
	tableName = tableNameVale.Call([]reflect.Value{})[0].String()

	return fmt.Sprintf("INSERT INTO common.%s (%s) VALUES (%s)", tableName, strings.Join(columns, ", "), strings.Join(values, ", "))
}
