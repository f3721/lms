package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-admin/common/utils"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/admin/models"
	"go-admin/app/admin/service/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type SysUserLog struct {
	service.Service
}

// GetPage 获取SysUserLog列表
func (e *SysUserLog) GetPage(c *dto.SysUserLogGetPageReq, p *actions.DataPermission, list *[]models.SysUserLog, count *int64) error {
	var err error
	var data models.SysUserLog

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("SysUserLogService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取SysUserLog对象
func (e *SysUserLog) Get(d *dto.SysUserLogGetReq, p *actions.DataPermission, detailResp *utils.OperateLogDetailResp) error {
	var data models.SysUserLog

	err := e.Orm.
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(&data, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetSysUserLog error:%s \r\n", err)
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
