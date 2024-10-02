package models

import (

	"go-admin/common/models"
	"time"
)

type PresoImage struct {
    models.Model
    
    Name string `json:"name" gorm:"type:varchar(45);comment:文件名称"`
	PresoNo string `json:"presoNo" gorm:"type:varchar(30);comment:审批单编号"`
    Url string `json:"url" gorm:"type:varchar(255);comment:文件地址"`
    Type int `json:"type" gorm:"type:tinyint(1);comment:类型: 0-审批单文件"`
	CreatedAt time.Time      `json:"createdAt" gorm:"comment:创建时间"`
}

func (PresoImage) TableName() string {
    return "preso_image"
}

func (e *PresoImage) GetId() interface{} {
	return e.Id
}
