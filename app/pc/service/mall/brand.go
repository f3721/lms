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

type Brand struct {
	service.Service
}

// GetPage 获取Brand列表
func (e *Brand) GetPage(c *dto.BrandGetPageReq, p *actions.DataPermission, list *[]models.Brand, count *int64) error {
	var err error
	var data models.Brand

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("BrandService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取Brand对象
func (e *Brand) Get(d *dto.BrandGetReq, p *actions.DataPermission, model *models.Brand) error {
	var data models.Brand

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetBrand error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

func (e *Brand) GetBrandList(ids []int, list *[]dto.BrandInfo) error {
	var data models.Brand
	err := e.Orm.Model(&data).Select([]string{"id as brand_id", "brand_zh"}).Find(list, ids).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		return err
	}
	if err != nil {
		return err
	}
	return nil
}
