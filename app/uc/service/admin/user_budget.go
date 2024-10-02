package admin

import (
	"errors"

    "github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/uc/models"
	"go-admin/app/uc/service/admin/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type UserBudget struct {
	service.Service
}

// GetPage 获取UserBudget列表
func (e *UserBudget) GetPage(c *dto.UserBudgetGetPageReq, p *actions.DataPermission, list *[]dto.UserBudgetGetPageResp, count *int64) error {
	var err error
	var data models.UserBudget

	err = e.Orm.Model(&data).
		Select("user_budget.*, ui.user_name").
		Joins("left join user_info ui on user_budget.user_id = ui.id").
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
			dto.UserBudgetGetPageMakeCondition(c),
			actions.SysUserPermission("ui", p, 1),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("UserBudgetService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取UserBudget对象
func (e *UserBudget) Get(d *dto.UserBudgetGetReq, p *actions.DataPermission, model *models.UserBudget) error {
	var data models.UserBudget

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetUserBudget error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建UserBudget对象
func (e *UserBudget) Insert(c *dto.UserBudgetInsertReq) error {
    var err error
    var data models.UserBudget
    c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("UserBudgetService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改UserBudget对象
func (e *UserBudget) Update(c *dto.UserBudgetUpdateReq, p *actions.DataPermission) error {
    var err error
    var data = models.UserBudget{}
    e.Orm.Scopes(
            actions.Permission(data.TableName(), p),
        ).First(&data, c.GetId())
    c.Generate(&data)

    db := e.Orm.Save(&data)
    if err = db.Error; err != nil {
        e.Log.Errorf("UserBudgetService Save error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权更新该数据")
    }
    return nil
}

// Remove 删除UserBudget
func (e *UserBudget) Remove(d *dto.UserBudgetDeleteReq, p *actions.DataPermission) error {
	var data models.UserBudget

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
        e.Log.Errorf("Service RemoveUserBudget error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权删除该数据")
    }
	return nil
}
