package admin

import (
	"errors"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"go-admin/app/pc/models"
	"go-admin/app/pc/service/admin/dto"
	"go-admin/common/actions"
)

type MediaRelationship struct {
	service.Service
}

// Insert 创建MediaRelationship对象
func (e *MediaRelationship) Insert(c *dto.MediaRelationshipInsertReq, data *models.MediaRelationship) error {
	var err error
	c.Generate(data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("MediaRelationshipService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改MediaRelationship对象
func (e *MediaRelationship) Update(c *dto.MediaRelationshipUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.MediaRelationship{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("MediaRelationshipService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除MediaRelationship
func (e *MediaRelationship) Remove(mediaTypeId int, buszId interface{}) error {
	var data []models.MediaRelationship
	err := e.Orm.Model(&models.MediaRelationship{}).Where("media_type_id = ?", mediaTypeId).Where("busz_id = ?", buszId).Find(&data).Error
	if err != nil {
		return err
	}
	if len(data) > 0 {
		//var ids []int
		//if err := utils.StructColumn(&ids, data, "MediaInstantId", ""); err != nil {
		//	return err
		//}
		db := e.Orm.Model(&models.MediaRelationship{}).Unscoped().Delete(&data)
		if err := db.Error; err != nil {
			return err
		}
		//mediaInstance := MediaInstance{e.Service}
		//err = mediaInstance.Remove(&dto.MediaInstanceDeleteReq{
		//	Ids: ids,
		//})
		//if err != nil {
		//	return err
		//}
	}
	return nil
}
