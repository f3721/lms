package models

import (
	"go-admin/common/models"
)

type OperateLogTypes struct {
	models.Model

	ModelName string `json:"modelName" gorm:"type:varchar(50);comment:模型name"`
	Type      string `json:"type" gorm:"type:varchar(20);comment:eg: create"`
	Mapping   string `json:"mapping" gorm:"type:text;comment:字段映射"`
	Name      string `json:"name" gorm:"type:varchar(50);comment:操作名称"`
}

func (OperateLogTypes) TableName() string {
	return "operate_log_types"
}

func (e *OperateLogTypes) GetId() interface{} {
	return e.Id
}
