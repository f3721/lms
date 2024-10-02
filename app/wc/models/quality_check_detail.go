package models

import (
	"go-admin/common/models"
)

type QualityCheckDetail struct {
	models.Model
	QualityCheckTaskId int    `json:"qualityCheckTaskId" gorm:"type:int;comment:质检任务id"`
	QualityCheckOption string `json:"qualityCheckOption" gorm:"type:varchar(100);comment:质检内容"`
	QualityBy          int    `json:"qualityBy" gorm:"type:int;comment:质检人"`
	QualityByName      string `json:"qualityByName" gorm:"type:varchar(20);comment:质检人名称"`
	QualityRes         int    `json:"qualityRes" gorm:"type:int;comment:质检结果:1合格，2不合格"`
	Remark             string `json:"remark" gorm:"type:varchar(100);comment:备注"`
	models.ModelTime
	models.ControlBy
}

func (QualityCheckDetail) TableName() string {
	return "quality_check_task_detail"
}

func (e *QualityCheckDetail) GetId() interface{} {
	return e.Id
}
