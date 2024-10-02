package mall

import (
	"errors"
	"go-admin/common/global"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/pc/models"
	"go-admin/app/pc/service/mall/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type MediaRelationship struct {
	service.Service
}

// GetPage 获取MediaRelationship列表
func (e *MediaRelationship) GetPage(c *dto.MediaRelationshipGetPageReq, p *actions.DataPermission, list *[]models.MediaRelationship, count *int64) error {
	var err error
	var data models.MediaRelationship

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("MediaRelationshipService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取MediaRelationship对象
func (e *MediaRelationship) Get(d *dto.MediaRelationshipGetReq, p *actions.DataPermission, model *models.MediaRelationship) error {
	var data models.MediaRelationship

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetMediaRelationship error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

func (e *MediaRelationship) getMediaList(bizId []string, list *[]models.MediaRelationship) error {
	var data models.MediaRelationship
	pcPrefix := global.GetTenantPcDBNameWithDB(e.Orm)
	err := e.Orm.Table(pcPrefix+"."+data.TableName()).Preload("MediaInstant", func(db *gorm.DB) *gorm.DB {
		return db.Table(pcPrefix + ".media_instance")
	}).Where("busz_id in ?", bizId).Order("seq ASC").Find(list).Error
	if err != nil {
		return err
	}
	return nil
}

func mediaRelation(data []models.MediaRelationship) map[string][]models.MediaRelationship {
	res := make(map[string][]models.MediaRelationship, 0)
	for _, mediaRelation := range data {
		res[mediaRelation.BuszId] = append(res[mediaRelation.BuszId], mediaRelation)
	}
	return res
}
