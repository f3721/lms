package admin

import (
	"errors"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/wc/models"
	"go-admin/app/wc/service/admin/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type ReminderRuleSkuLog struct {
	service.Service
}

// GetPage 获取ReminderRuleSkuLog列表
func (e *ReminderRuleSkuLog) GetPage(c *dto.ReminderRuleSkuLogGetPageReq, p *actions.DataPermission, list *[]models.ReminderRuleSkuLog, count *int64) error {
	var err error
	var data models.ReminderRuleSkuLog

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("ReminderRuleSkuLogService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取ReminderRuleSkuLog对象
func (e *ReminderRuleSkuLog) Get(d *dto.ReminderRuleSkuLogGetReq, p *actions.DataPermission, model *models.ReminderRuleSkuLog) error {
	var data models.ReminderRuleSkuLog

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetReminderRuleSkuLog error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建ReminderRuleSkuLog对象
func (e *ReminderRuleSkuLog) Insert(c *dto.ReminderRuleSkuLogInsertReq) error {
	var err error
	var data models.ReminderRuleSkuLog
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("ReminderRuleSkuLogService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改ReminderRuleSkuLog对象
func (e *ReminderRuleSkuLog) Update(c *dto.ReminderRuleSkuLogUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.ReminderRuleSkuLog{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("ReminderRuleSkuLogService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除ReminderRuleSkuLog
func (e *ReminderRuleSkuLog) Remove(d *dto.ReminderRuleSkuLogDeleteReq, p *actions.DataPermission) error {
	var data models.ReminderRuleSkuLog

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveReminderRuleSkuLog error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
