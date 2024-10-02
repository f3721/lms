package global

import (
	"errors"
	"go-admin/common/models"
	"go-admin/common/utils"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var tenants = make(map[string]*models.SystemTenant)

var lock = new(sync.RWMutex)

func GetTenant(key string) (*models.SystemTenant, error) {
	lock.RLock()
	defer lock.RUnlock()
	tenant, ok := tenants[key]
	if !ok {
		return nil, errors.New("对应tenantId的租户不存在")
	}
	return tenant, nil
}

func GetTenantWithContext(c *gin.Context) (*models.SystemTenant, error) {
	tenantId := c.GetHeader("tenant-id")
	if tenantId == "" {
		//response.Error(c, 500, nil, "header中tenant-id未填写")
		return nil, errors.New("header中tenant-id未填写")
	}
	//id, err := pkg.StringToInt(tenantId)
	//if err != nil {
	//	//response.Error(c, 500, nil, "tenant-id非数值型")
	//	return nil, errors.New("tenant-id非数值型")
	//}

	return GetTenant(tenantId)
}

func GetTenantDBPrefixWithContext(c *gin.Context) string {
	tenant, err := GetTenantWithContext(c)
	if err != nil {
		return "base"
	}
	return tenant.TenantDBPrefix()
}

func GetTenantAdminDBNameWithContext(c *gin.Context) string {
	prefix := GetTenantDBPrefixWithContext(c)
	return prefix + "_admin"
}

func GetTenantUcDBNameWithContext(c *gin.Context) string {
	prefix := GetTenantDBPrefixWithContext(c)
	return prefix + "_uc"
}

func GetTenantPcDBNameWithContext(c *gin.Context) string {
	prefix := GetTenantDBPrefixWithContext(c)
	return prefix + "_pc"
}

func GetTenantWcDBNameWithContext(c *gin.Context) string {
	prefix := GetTenantDBPrefixWithContext(c)
	return prefix + "_wc"
}

func GetTenantOcDBNameWithContext(c *gin.Context) string {
	prefix := GetTenantDBPrefixWithContext(c)
	return prefix + "_oc"
}

func GetTenantUcDBNameWithDB(db *gorm.DB) string {
	ctx := db.Statement.Context
	switch ctx.(type) {
	case *gin.Context:
		return GetTenantUcDBNameWithContext(ctx.(*gin.Context))
	default:
		return "199_databasesxyz004_uc"
	}
}

func GetTenantPcDBNameWithDB(db *gorm.DB) string {
	ctx := db.Statement.Context
	switch ctx.(type) {
	case *gin.Context:
		return GetTenantPcDBNameWithContext(ctx.(*gin.Context))
	default:
		return "199_databasesxyz004_pc"
	}
}

func GetTenantWcDBNameWithDB(db *gorm.DB) string {
	ctx := db.Statement.Context
	switch ctx.(type) {
	case *gin.Context:
		return GetTenantWcDBNameWithContext(ctx.(*gin.Context))
	default:
		return "199_databasesxyz004_wc"
	}
}

func GetTenantOcDBNameWithDB(db *gorm.DB) string {
	ctx := db.Statement.Context
	switch ctx.(type) {
	case *gin.Context:
		return GetTenantOcDBNameWithContext(ctx.(*gin.Context))
	default:
		return "199_databasesxyz004_oc"
	}
}

func GetTenantDBPrefixWithDB(db *gorm.DB) (string, error) {
	ctx := db.Statement.Context
	switch ctx.(type) {
	case *gin.Context:
		return GetTenantDBPrefixWithContext(ctx.(*gin.Context)), nil
	default:
		return "", errors.New("db.Statement.Context未非gin.Context,无法获取gin.Context中的租户相关信息")
	}
}

func GetTenantAdminDBNameWithDB(db *gorm.DB) string {
	ctx := db.Statement.Context
	switch ctx.(type) {
	case *gin.Context:
		return GetTenantAdminDBNameWithContext(ctx.(*gin.Context))
	default:
		return "199_databasesxyz004_admin"
	}
}

func GetTenants() map[string]*models.SystemTenant {
	return tenants
}

func SetTenant(tenant models.SystemTenant) error {
	lock.Lock()
	defer lock.Unlock()
	tenants[EncryptTenantId(tenant.ID)] = &tenant
	return nil
}

func DelTenant(tenant *models.SystemTenant) {
	lock.Lock()
	defer lock.Unlock()
	delete(tenants, EncryptTenantId(tenant.ID))
}

func DelTenantWithId(tid int) {
	lock.Lock()
	defer lock.Unlock()
	delete(tenants, EncryptTenantId(tid))
}

func EncryptTenantId(tid int) string {
	salt := "23530lms1452"
	return utils.EncryptMd5(string(rune(tid)), salt)
}

func GetSyncTenantDBNameWithDB(db *gorm.DB, tenantId int) string {
	ctx := db.Statement.Context
	switch ctx.(type) {
	case *gin.Context:
		tenant, err := GetTenant(EncryptTenantId(tenantId))
		if err != nil {
			return "199_databasesxyz004_pc"
		}
		return tenant.TenantDBPrefix()
	default:
		return "199_databasesxyz004_pc"
	}
}

// 跨库DB操作，指定表名
func TenantTable(prefixName, tableName string, alias ...string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// 获取库名
		dbPrefix := GetPrefix(db, prefixName)
		if dbPrefix == "" {
			return db
		}

		// 指定表名
		fullTableName := dbPrefix + "." + tableName
		if len(alias) > 0 {
			fullTableName += " AS " + alias[0]
		}
		return db.Table(fullTableName)
	}
}

// GetPrefix 查询跨库库名
func GetPrefix(db *gorm.DB, prefixName string) (prefix string) {
	switch prefixName {
	case "pc":
		prefix = GetTenantPcDBNameWithDB(db)
	case "oc":
		prefix = GetTenantOcDBNameWithDB(db)
	case "wc":
		prefix = GetTenantWcDBNameWithDB(db)
	case "uc":
		prefix = GetTenantUcDBNameWithDB(db)
	case "admin":
		prefix = GetTenantAdminDBNameWithDB(db)
	default:
		return
	}
	return
}

// 跨库DB操作，获取表名
func TenantTableName(db *gorm.DB, prefixName, tableName string, alias ...string) string {
	// 获取库名
	dbPrefix := GetPrefix(db, prefixName)

	// 指定表名
	fullTableName := dbPrefix + "." + tableName
	if len(alias) > 0 {
		fullTableName += " AS " + alias[0]
	}
	return fullTableName
}
