package models

import (
	"go-admin/common/models"
)

type MediaInstance struct {
	models.Model
	MediaDir  string `json:"mediaDir" gorm:"type:varchar(255);comment:文件路径"`
	MediaName string `json:"mediaName" gorm:"type:varchar(255);comment:文件名"`
}

func (MediaInstance) TableName() string {
	return "media_instance"
}

func (e *MediaInstance) GetId() interface{} {
	return e.Id
}
