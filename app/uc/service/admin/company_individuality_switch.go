package admin

import (
	"errors"
	"go-admin/app/uc/service/admin/dto"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/uc/models"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type CompanyIndividualitySwitch struct {
	service.Service
}

// GetPage 获取CompanyIndividualitySwitch列表
func (e *CompanyIndividualitySwitch) GetPage(c *dto.CompanyIndividualitySwitchGetPageReq, p *actions.DataPermission, list *[]models.CompanyIndividualitySwitch, count *int64) error {
	var err error
	var data models.CompanyIndividualitySwitch

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("CompanyIndividualitySwitchService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取CompanyIndividualitySwitch对象
func (e *CompanyIndividualitySwitch) Get(d *dto.CompanyIndividualitySwitchGetReq, p *actions.DataPermission, model *models.CompanyIndividualitySwitch) error {
	var data models.CompanyIndividualitySwitch

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetCompanyIndividualitySwitch error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建CompanyIndividualitySwitch对象
func (e *CompanyIndividualitySwitch) Insert(c *dto.CompanyIndividualitySwitchInsertReq) error {
	var err error
	var data models.CompanyIndividualitySwitch
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("CompanyIndividualitySwitchService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改CompanyIndividualitySwitch对象
func (e *CompanyIndividualitySwitch) Update(c *dto.CompanyIndividualitySwitchUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.CompanyIndividualitySwitch{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("CompanyIndividualitySwitchService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除CompanyIndividualitySwitch
func (e *CompanyIndividualitySwitch) Remove(d *dto.CompanyIndividualitySwitchDeleteReq, p *actions.DataPermission) error {
	var data models.CompanyIndividualitySwitch

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveCompanyIndividualitySwitch error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
