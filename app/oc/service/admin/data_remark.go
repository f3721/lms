package admin

import (
	"errors"

    "github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/oc/models"
	"go-admin/app/oc/service/admin/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type DataRemark struct {
	service.Service
}

// GetPage 获取DataRemark列表
func (e *DataRemark) GetPage(c *dto.DataRemarkGetPageReq, p *actions.DataPermission, list *[]models.DataRemark, count *int64) error {
	var err error
	var data models.DataRemark

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("DataRemarkService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取DataRemark对象
func (e *DataRemark) Get(d *dto.DataRemarkGetReq, p *actions.DataPermission, model *models.DataRemark) error {
	var data models.DataRemark

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetDataRemark error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建DataRemark对象
func (e *DataRemark) Insert(c *dto.DataRemarkInsertReq) error {
    var err error
    var data models.DataRemark
    c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("DataRemarkService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改DataRemark对象
func (e *DataRemark) Update(c *dto.DataRemarkUpdateReq, p *actions.DataPermission) error {
    var err error
    var data = models.DataRemark{}
    e.Orm.Scopes(
            actions.Permission(data.TableName(), p),
        ).First(&data, c.GetId())
    c.Generate(&data)

    db := e.Orm.Save(&data)
    if err = db.Error; err != nil {
        e.Log.Errorf("DataRemarkService Save error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权更新该数据")
    }
    return nil
}

// Remove 删除DataRemark
func (e *DataRemark) Remove(d *dto.DataRemarkDeleteReq, p *actions.DataPermission) error {
	var data models.DataRemark

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
        e.Log.Errorf("Service RemoveDataRemark error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权删除该数据")
    }
	return nil
}
