package admin

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/oc/models"
	"go-admin/app/oc/service/admin/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type CsApplyLog struct {
	service.Service
}

// GetPage 获取CsApplyLog列表
func (e *CsApplyLog) GetPage(c *dto.CsApplyLogGetPageReq, p *actions.DataPermission, list *[]models.CsApplyLog, count *int64) error {
	var err error
	var data models.CsApplyLog

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("CsApplyLogService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取CsApplyLog对象
func (e *CsApplyLog) Get(d *dto.CsApplyLogGetReq, p *actions.DataPermission, model *models.CsApplyLog) error {
	var data models.CsApplyLog

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetCsApplyLog error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建CsApplyLog对象
func (e *CsApplyLog) Insert(c *dto.CsApplyLogInsertReq) error {
	var err error
	var data models.CsApplyLog
	c.Generate(&data)
	
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("CsApplyLogService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改CsApplyLog对象
func (e *CsApplyLog) Update(c *dto.CsApplyLogUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.CsApplyLog{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("CsApplyLogService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除CsApplyLog
func (e *CsApplyLog) Remove(d *dto.CsApplyLogDeleteReq, p *actions.DataPermission) error {
	var data models.CsApplyLog

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveCsApplyLog error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}

// Insert 创建CsApplyLog对象
func (e *CsApplyLog) AddLog(db *gorm.DB, csNo string, text string) error {
	var err error
	data := &models.CsApplyLog{
		CsNo:         csNo,
		UserId:       "",
		HandlerLog:   text,
		UserName:     "",
		CreateBy:     user.GetUserId(db.Statement.Context.(*gin.Context)),
		CreateByName: user.GetUserName(db.Statement.Context.(*gin.Context)),
	}

	err = db.Create(&data).Error
	if err != nil {
		e.Log.Errorf("CsApplyLogService Insert error:%s \r\n", err)
		return err
	}
	return nil
}
