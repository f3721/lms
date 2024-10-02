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

type ReminderListSku struct {
	service.Service
}

// GetPage 获取ReminderListSku列表
func (e *ReminderListSku) GetPage(c *dto.ReminderListSkuGetPageReq, p *actions.DataPermission, list *[]models.ReminderListSku, count *int64) error {
	var err error
	var data models.ReminderListSku

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("ReminderListSkuService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// GetExportAll 获取ReminderListSku列表
func (e *ReminderListSku) GetExportAll(c *dto.ReminderListSkuGetPageReq, p *actions.DataPermission, list *[]models.ReminderListSku) error {
	var err error
	var data models.ReminderListSku

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Error
	if err != nil {
		e.Log.Errorf("ReminderListSkuService GetAll error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取ReminderListSku对象
func (e *ReminderListSku) Get(d *dto.ReminderListSkuGetReq, p *actions.DataPermission, model *models.ReminderListSku) error {
	var data models.ReminderListSku

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetReminderListSku error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建ReminderListSku对象
func (e *ReminderListSku) Insert(c *dto.ReminderListSkuInsertReq) error {
	var err error
	var data models.ReminderListSku
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("ReminderListSkuService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Inserts 批量创建ReminderListSku对象
func (e *ReminderListSku) Inserts(skus []*dto.ReminderListSkuInsertReq, id int) error {
	var err error

	if id == 0 {
		return errors.New("数据错误无法插入sku设置")
	}

	skuAdds := []models.ReminderListSku{}
	for _, req := range skus {
		data := models.ReminderListSku{}
		req.Generate(&data)
		data.ReminderListId = id
		skuAdds = append(skuAdds, data)
	}

	err = e.Orm.Create(skuAdds).Error
	if err != nil {
		e.Log.Errorf("ReminderRuleSkuService Insert error:%s \r\n", err)
		return err
	}

	return nil
}

// Update 修改ReminderListSku对象
func (e *ReminderListSku) Update(c *dto.ReminderListSkuUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.ReminderListSku{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("ReminderListSkuService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除ReminderListSku
func (e *ReminderListSku) Remove(d *dto.ReminderListSkuDeleteReq, p *actions.DataPermission) error {
	var data models.ReminderListSku

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveReminderListSku error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
