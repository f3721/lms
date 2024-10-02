package admin

import (
	"errors"
	"fmt"
	"github.com/go-admin-team/go-admin-core/sdk"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"go-admin/app/uc/models"
	"go-admin/app/uc/service/admin/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
	"go-admin/common/global"
	"gorm.io/gorm"
	"strconv"
)

type DepartmentBudget struct {
	service.Service
}

// GetPage 获取DepartmentBudget列表
func (e *DepartmentBudget) GetPage(c *dto.DepartmentBudgetGetPageReq, p *actions.DataPermission, list *[]dto.DepartmentBudgetListResp, count *int64) error {
	var err error
	var data models.DepartmentBudget

	err = e.Orm.Model(&data).
		Select("department_budget.*, cd.name dep_name, cd.company_id, ci.company_name, cd2.name parent_dep_name, cd.personal_budget").
		Joins("left join company_department cd on department_budget.dep_id = cd.id").
		Joins("left join company_info ci on cd.company_id = ci.id").
		Joins("left join company_department cd2 on cd.f_id = cd2.id").
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
			dto.GetPageMakeCondition(c),
			actions.SysUserPermission("cd", p, 1),
		).
		Scan(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("DepartmentBudgetService GetPage error:%s \r\n", err)
		return err
	}
	tmpList := *list
	for i, v := range tmpList {
		if v.PersonalBudget.Valid == true {
			tmpList[i].IsUserBudget = 1
		}
	}
	*list = tmpList
	return nil
}

// Get 获取DepartmentBudget对象
func (e *DepartmentBudget) Get(d *dto.DepartmentBudgetGetReq, p *actions.DataPermission, model *dto.DepartmentBudgetGetResp) error {
	var data models.DepartmentBudget

	err := e.Orm.Model(&data).
		Select("department_budget.*, cd.name dep_name, ci.company_name").
		Joins("left join company_department cd on department_budget.dep_id = cd.id").
		Joins("left join company_info ci on cd.company_id = ci.id").
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetDepartmentBudget error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建DepartmentBudget对象
func (e *DepartmentBudget) Insert(c *dto.DepartmentBudgetInsertReq) error {
    var err error
    var data models.DepartmentBudget
    c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("DepartmentBudgetService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改DepartmentBudget对象
func (e *DepartmentBudget) Update(c *dto.DepartmentBudgetUpdateReq, p *actions.DataPermission) error {
    var err error
    var data = models.DepartmentBudget{}
    e.Orm.Scopes(
            actions.Permission(data.TableName(), p),
        ).First(&data, c.GetId())
    c.Generate(&data)

    db := e.Orm.Save(&data)
    if err = db.Error; err != nil {
        e.Log.Errorf("DepartmentBudgetService Save error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权更新该数据")
    }
    return nil
}

// Remove 删除DepartmentBudget
func (e *DepartmentBudget) Remove(d *dto.DepartmentBudgetDeleteReq, p *actions.DataPermission) error {
	var data models.DepartmentBudget

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
        e.Log.Errorf("Service RemoveDepartmentBudget error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权删除该数据")
    }
	return nil
}

func (e *DepartmentBudget) Init() (echoMsg string, err error) {
	tenants := global.GetTenants()
	if len(tenants) <= 0 {
		err = errors.New("暂无租户")
		return
	}
	echoMsg = "初始化部门预算脚本执行开始："
	for _, tenant := range tenants {
		//tenantDbSource := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v_%v?charset=utf8mb4&parseTime=True&loc=Local&timeout=1000ms",
		//	tenant.DatabaseUsername,
		//	tenant.DatabasePassword,
		//	tenant.DatabaseHost,
		//	tenant.DatabasePort,
		//	tenant.TenantDBPrefix(),
		//	"oc")
		//tenantDb, err := gorm.Open(mysql.Open(tenantDbSource), &gorm.Config{})
		//if err != nil {
		//	echoMsg = fmt.Sprintf("%s 数据库连接失败", echoMsg)
		//	continue
		//}

		tenantDBPrefix := tenant.TenantDBPrefix()
		tenantDb := sdk.Runtime.GetDbByKey(tenantDBPrefix)
		echoMsg = fmt.Sprintf("%s\r\n租户%s[%s]:", echoMsg, tenant.Name, tenant.DatabaseName)
		if tenantDb == nil {
			echoMsg = fmt.Sprintf("%s 数据库连接失败", echoMsg)
			continue
		}

		var model models.DepartmentBudget
		var list []models.CompanyDepartment
		tenantDb.Table("company_department").Joins("left join company_info ci on company_department.company_id = ci.id").Where("ci.company_status = 1").Find(&list)
		if len(list) == 0 {
			echoMsg = fmt.Sprintf("%s 无部门需要初始化预算", echoMsg)
			continue
		}

		echoMsg = fmt.Sprintf("%s\r\n", echoMsg)
		for _, data := range list {
			err = model.Generate(tenantDb, data.Id, 0, 1)
			if err != nil {
				echoMsg = fmt.Sprintf("%s部门[%s]fail | ", echoMsg, strconv.Itoa(data.Id))
			} else {
				echoMsg = fmt.Sprintf("%s部门[%s]success | ", echoMsg, strconv.Itoa(data.Id))
			}
		}
		echoMsg = fmt.Sprintf("%s\r\n", echoMsg)
	}
	echoMsg = fmt.Sprintf("%s\r\n 脚本执行完毕", echoMsg)

	return echoMsg, nil
}
