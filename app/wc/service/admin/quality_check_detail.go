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

type QualityCheckDetail struct {
	service.Service
}

// GetPage 获取QualityCheckDetail列表
func (e *QualityCheckDetail) GetPage(c *dto.QualityCheckGetPageReq, p *actions.DataPermission, list *[]models.QualityCheckDetail, count *int64) error {
	var err error
	var data models.QualityCheckDetail

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("QualityCheckDetailService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取QualityCheckDetail对象
func (e *QualityCheckDetail) Get(d *dto.QualityCheckDetailGetReq, p *actions.DataPermission, model *models.QualityCheckDetail) error {
	var data models.QualityCheckDetail
	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetQualityCheckDetail error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}
