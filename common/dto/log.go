package dto

import "go-admin/common/models"

type LogCreateReq struct {
	DataId     int    `json:"dataId" gorm:"type:int;comment:关联的主键ID"`
	Type       string `json:"type" gorm:"type:varchar(10);comment:变更类型"`
	Data       []byte `json:"data" gorm:"type:json;comment:源数据"`
	BeforeData []byte `json:"beforeData" gorm:"type:json;comment:变更前数据"`
	AfterData  []byte `json:"afterData" gorm:"type:json;comment:变更后数据"`
	models.ControlBy
}
