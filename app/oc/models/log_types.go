package models

import (

	"go-admin/common/models"

)

type LogTypes struct {
    models.Model
    
    ModelName string `json:"model_name" gorm:"type:varchar(30);comment:ModelName"` 
    Type string `json:"type" gorm:"type:varchar(10);comment:Type"` 
    Mapping string `json:"mapping" gorm:"type:json;comment:字段映射"` 
    Title string `json:"title" gorm:"type:varchar(50);comment:操作名称"`
}

func (LogTypes) TableName() string {
    return "log_types"
}

func (e *LogTypes) GetId() interface{} {
	return e.Id
}
