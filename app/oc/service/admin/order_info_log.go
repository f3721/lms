package admin

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-admin/common/utils"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/oc/models"
	"go-admin/app/oc/service/admin/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type OrderInfoLog struct {
	service.Service
}

// GetPage 获取OrderInfoLog列表
func (e *OrderInfoLog) GetPage(c *dto.OrderInfoLogGetPageReq, p *actions.DataPermission, list *[]models.OrderInfoLog, count *int64) error {
	var err error
	var data models.OrderInfoLog

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("OrderInfoLogService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取OrderInfoLog对象
func (e *OrderInfoLog) Get(d *dto.OrderInfoLogGetReq, p *actions.DataPermission, detailResp *utils.OperateLogDetailResp) error {
	var data models.OrderInfoLog

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(&data, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetOrderInfoLog error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	if err = json.Unmarshal([]byte(fmt.Sprintf(`%s`, data.DiffData)), detailResp); err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建OrderInfoLog对象
func (e *OrderInfoLog) Insert(c *dto.OrderInfoLogInsertReq) error {
    var err error
    var data models.OrderInfoLog
    c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("OrderInfoLogService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改OrderInfoLog对象
func (e *OrderInfoLog) Update(c *dto.OrderInfoLogUpdateReq, p *actions.DataPermission) error {
    var err error
    var data = models.OrderInfoLog{}
    e.Orm.Scopes(
            actions.Permission(data.TableName(), p),
        ).First(&data, c.GetId())
    c.Generate(&data)

    db := e.Orm.Save(&data)
    if err = db.Error; err != nil {
        e.Log.Errorf("OrderInfoLogService Save error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权更新该数据")
    }
    return nil
}

// Remove 删除OrderInfoLog
func (e *OrderInfoLog) Remove(d *dto.OrderInfoLogDeleteReq, p *actions.DataPermission) error {
	var data models.OrderInfoLog

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
        e.Log.Errorf("Service RemoveOrderInfoLog error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权删除该数据")
    }
	return nil
}
