package actions

import (
	"errors"
	"github.com/go-admin-team/go-admin-core/sdk/config"
	"go-admin/common/global"
	"go-admin/common/utils"

	"github.com/gin-gonic/gin"
	log "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"gorm.io/gorm"
)

type DataPermission struct {
	DataScope string
	UserId    int
	DeptId    int
	RoleId    int
	AuthorityCompanyId string
	AuthorityWarehouseId string
	AuthorityWarehouseAllocateId string
	AuthorityVendorId string
}

func PermissionAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		db, err := pkg.GetOrm(c)
		if err != nil {
			log.Error(err)
			return
		}

		msgID := pkg.GenerateMsgIDFromContext(c)
		var p = new(DataPermission)
		if userId := user.GetUserIdStr(c); userId != "" {
			p, err = newDataPermission(db, userId)
			if err != nil {
				log.Errorf("MsgID[%s] PermissionAction error: %s", msgID, err)
				response.Error(c, 500, err, "权限范围鉴定错误")
				c.Abort()
				return
			}
		}
		c.Set(PermissionKey, p)
		c.Next()
	}
}

func newDataPermission(tx *gorm.DB, userId interface{}) (*DataPermission, error) {
	var err error
	p := &DataPermission{}
	adminDBName := global.GetTenantAdminDBNameWithDB(tx)
	err = tx.Table(adminDBName+".sys_user").
		Select("sys_user.user_id", "sys_role.role_id", "sys_user.dept_id", "sys_role.data_scope", "sys_user.authority_company_id", "sys_user.authority_warehouse_id", "sys_user.authority_warehouse_allocate_id", "sys_user.authority_vendor_id").
		Joins("left join "+adminDBName+".sys_role on sys_role.role_id = sys_user.role_id").
		Where("sys_user.user_id = ?", userId).
		Scan(p).Error
	if err != nil {
		err = errors.New("获取用户数据出错 msg:" + err.Error())
		return nil, err
	}
	return p, nil
}

func Permission(tableName string, p *DataPermission) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if !config.ApplicationConfig.EnableDP {
			return db
		}
		adminDbName := global.GetTenantAdminDBNameWithDB(db)
		switch p.DataScope {
		case "2":
			return db.Where(tableName+".create_by in (select sys_user.user_id from "+adminDbName+".sys_role_dept left join "+adminDbName+".sys_user on sys_user.dept_id=sys_role_dept.dept_id where sys_role_dept.role_id = ?)", p.RoleId)
		case "3":
			return db.Where(tableName+".create_by in (SELECT user_id from "+adminDbName+".sys_user where dept_id = ? )", p.DeptId)
		case "4":
			return db.Where(tableName+".create_by in (SELECT user_id from "+adminDbName+".sys_user where sys_user.dept_id in(select dept_id from "+adminDbName+".sys_dept where dept_path like ? ))", "%/"+pkg.IntToString(p.DeptId)+"/%")
		case "5":
			return db.Where(tableName+".create_by = ?", p.UserId)
		default:
			return db
		}
	}
}

func getPermissionFromContext(c *gin.Context) *DataPermission {
	p := new(DataPermission)
	if pm, ok := c.Get(PermissionKey); ok {
		switch pm.(type) {
		case *DataPermission:
			p = pm.(*DataPermission)
		}
	}
	return p
}

// GetPermissionFromContext 提供非action写法数据范围约束
func GetPermissionFromContext(c *gin.Context) *DataPermission {
	return getPermissionFromContext(c)
}

//SysUserPermission 系统用户 公司、仓库、仓库调拨、货主权限校验
func SysUserPermission(tableName string, p *DataPermission, mode int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		//if !config.ApplicationConfig.EnableDP {
		//	return db
		//}
		//adminDbName := global.GetTenantAdminDBNameWithDB(db)
		switch mode {
		// 公司
		case 1:
			return db.Where(tableName+".company_id in ?", utils.SplitToInt(p.AuthorityCompanyId))
		// 仓库
		case 2:
			return db.Where(tableName+".warehouse_code in ?", utils.Split(p.AuthorityWarehouseId))
		// 仓库调拨
		case 3:
			return db.Where(tableName+".warehouse_code in ?", utils.Split(p.AuthorityWarehouseAllocateId))
		// 货主
		case 4:
			return db.Where(tableName+".vendor_id in ?", utils.SplitToInt(p.AuthorityVendorId))
		default:
			return db
		}
	}
}
