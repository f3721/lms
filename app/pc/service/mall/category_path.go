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

type CategoryPath struct {
	service.Service
}

// GetPage 获取CategoryPath列表
func (e *CategoryPath) GetPage(c *dto.CategoryPathGetPageReq, p *actions.DataPermission, list *[]models.CategoryPath, count *int64) error {
	var err error
	var data models.CategoryPath

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("CategoryPathService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取CategoryPath对象
func (e *CategoryPath) Get(d *dto.CategoryPathGetReq, p *actions.DataPermission, model *models.CategoryPath) error {
	var data models.CategoryPath

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetCategoryPath error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}
