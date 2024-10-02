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

type CategoryLog struct {
	service.Service
}

// GetPage 获取CategoryLog列表
func (e *CategoryLog) GetPage(c *dto.CategoryLogGetPageReq, p *actions.DataPermission, list *[]models.CategoryLog, count *int64) error {
	var err error
	var data models.CategoryLog

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("CategoryLogService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取CategoryLog对象
func (e *CategoryLog) Get(d *dto.CategoryLogGetReq, p *actions.DataPermission, detailResp *utils.OperateLogDetailResp) error {
	var model models.CategoryLog

	err := e.Orm.
		Scopes(
			actions.Permission(model.TableName(), p),
		).
		First(&model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetCategoryLog error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	if err := json.Unmarshal([]byte(fmt.Sprintf(`%s`, model.DiffData)), detailResp); err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}
