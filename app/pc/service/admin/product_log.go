package admin

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-admin/common/utils"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/pc/models"
	"go-admin/app/pc/service/admin/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type ProductLog struct {
	service.Service
}

// GetPage 获取ProductLog列表
func (e *ProductLog) GetPage(c *dto.ProductLogGetPageReq, p *actions.DataPermission, list *[]models.ProductLog, count *int64) error {
	var err error
	var data models.ProductLog

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("ProductLogService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取ProductLog对象
func (e *ProductLog) Get(d *dto.ProductLogGetReq, p *actions.DataPermission, detailResp *utils.OperateLogDetailResp) error {
	var data models.ProductLog

	err := e.Orm.
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(&data, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetProductLog error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	if err := json.Unmarshal([]byte(fmt.Sprintf(`%s`, data.DiffData)), detailResp); err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建ProductLog对象
func (e *ProductLog) Insert(c *dto.ProductLogInsertReq) error {
	var err error
	var data models.ProductLog
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("ProductLogService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改ProductLog对象
func (e *ProductLog) Update(c *dto.ProductLogUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.ProductLog{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("ProductLogService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除ProductLog
func (e *ProductLog) Remove(d *dto.ProductLogDeleteReq, p *actions.DataPermission) error {
	var data models.ProductLog

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveProductLog error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
