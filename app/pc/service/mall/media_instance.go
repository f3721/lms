package mall

import (
	"errors"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/pc/models"
	"go-admin/app/pc/service/mall/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type MediaInstance struct {
	service.Service
}

// GetPage 获取MediaInstance列表
func (e *MediaInstance) GetPage(c *dto.MediaInstanceGetPageReq, p *actions.DataPermission, list *[]models.MediaInstance, count *int64) error {
	var err error
	var data models.MediaInstance

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("MediaInstanceService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取MediaInstance对象
func (e *MediaInstance) Get(d *dto.MediaInstanceGetReq, p *actions.DataPermission, model *models.MediaInstance) error {
	var data models.MediaInstance

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetMediaInstance error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}
