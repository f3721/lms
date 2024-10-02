package admin

import (
	"errors"
	log "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/oc/models"
	"go-admin/app/oc/service/admin/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type CsApplyDetail struct {
	service.Service
}

// GetPage 获取CsApplyDetail列表
func (e *CsApplyDetail) GetPage(c *dto.CsApplyDetailGetPageReq, p *actions.DataPermission, list *[]models.CsApplyDetail, count *int64) error {
	var err error
	var data models.CsApplyDetail

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("CsApplyDetailService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// GetAll 获取CsApplyDetail列表
func (e *CsApplyDetail) GetAll(c *dto.CsApplyDetailGetPageReq, list *[]*models.CsApplyDetail) error {
	var err error
	var data models.CsApplyDetail

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
		).
		Find(list).Error
	if err != nil {
		e.Log.Errorf("CsApplyDetailService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取CsApplyDetail对象
func (e *CsApplyDetail) Get(d *dto.CsApplyDetailGetReq, p *actions.DataPermission, model *models.CsApplyDetail) error {
	var data models.CsApplyDetail

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetCsApplyDetail error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Get 获取CsApplyDetail对象
func (e *CsApplyDetail) GetCsApplyDetail(req *dto.CsApplyDetailGetReq, p *actions.DataPermission, model *dto.CsApplyDetailGetRes) error {

	model.ApplicationList, _ = e.GetApplicationList(req.CsNo)
	log.Info(model.ApplicationList)
	model.AfterSalesList, _ = e.GetAfterSalesList(req.CsNo)
	return nil
}

func (e *CsApplyDetail) GetApplicationList(csNo string) (*[]models.CsApplyDetail, error) {
	var err error
	var list []models.CsApplyDetail
	err = e.Orm.Debug().
		Where("cs_type = 0").
		Where("cs_no = ?", csNo).
		Find(&list).Error

	if err != nil {
		return nil, err
	}
	return &list, nil
}

func (e *CsApplyDetail) GetAfterSalesList(csNo string) (*[]models.CsApplyDetail, error) {
	var err error
	var list []models.CsApplyDetail
	err = e.Orm.Debug().
		Where("cs_type = 1").
		Where("cs_no = ?", csNo).
		Find(&list).Error
	if err != nil {
		return nil, err
	}
	return &list, nil
}

// Insert 创建CsApplyDetail对象
func (e *CsApplyDetail) Insert(c *dto.CsApplyDetailInsertReq) error {
	var err error
	var data models.CsApplyDetail
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("CsApplyDetailService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改CsApplyDetail对象
func (e *CsApplyDetail) Update(c *dto.CsApplyDetailUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.CsApplyDetail{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("CsApplyDetailService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除CsApplyDetail
func (e *CsApplyDetail) Remove(d *dto.CsApplyDetailDeleteReq, p *actions.DataPermission) error {
	var data models.CsApplyDetail

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveCsApplyDetail error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}

// GetAfterReturnProductsBySaleId 售后需要算的退货数
func (e *CsApplyDetail) GetAfterReturnProductsBySaleId(orderId string) (dto.AfterReturnProductsQuantity, error) {

	var csApplyDetail []*models.CsApplyDetail
	err := e.Orm.Table(models.CsApplyDetail{}.TableName()+" ad").
		Select("ad.sku_code,sum(ad.quantity) quantity").
		Where("ad.cs_type = 1 and a.cs_type=0 and a.cs_status in(1,99)").
		Where("a.order_id = ?", orderId).
		Joins(models.CsApply{}.TableName()+" a", "a.cs_no = ad.cs_no").
		Group("ad.sku_code").
		Find(&csApplyDetail).Error
	if err != nil {
		return nil, err
	}

	afterReturnProductsQuantity := make(dto.AfterReturnProductsQuantity)
	for _, detail := range csApplyDetail {
		afterReturnProductsQuantity[detail.SkuCode] = struct {
			SkuCode  string
			Quantity int
		}{
			SkuCode:  detail.SkuCode,
			Quantity: detail.Quantity,
		}
	}
	return afterReturnProductsQuantity, nil
}
