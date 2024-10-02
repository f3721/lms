package models

import (
	"go-admin/common/models"
)

type MediaRelationship struct {
	models.Model

	MediaTypeId    int           `json:"mediaTypeId" gorm:"type:int unsigned;comment:媒体类型id"`
	BuszId         string        `json:"buszId" gorm:"type:varchar(11);comment:业务id"`
	MediaInstantId int           `json:"mediaInstantId" gorm:"type:int unsigned;comment:MediaInstantId"`
	Seq            int           `json:"seq" gorm:"type:smallint unsigned;comment:排序号"`
	ImgFeatureId   string        `json:"imgFeatureId" gorm:"type:varchar(50);comment:图片的特征id"`
	Watermark      int           `json:"watermark" gorm:"type:tinyint(1);comment:图片是否有水印,0无水印,1,图片仅供参考,2,产品仅为附件"`
	MediaInstant   MediaInstance `json:"mediaInstance"`
	models.ModelTime
	models.ControlBy
}

func (MediaRelationship) TableName() string {
	return "media_relationship"
}

func (e *MediaRelationship) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *MediaRelationship) GetId() interface{} {
	return e.Id
}
