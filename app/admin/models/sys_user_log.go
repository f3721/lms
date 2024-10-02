package models

import (
	"encoding/json"
	"go-admin/common/utils"
	"gorm.io/gorm"
	"time"

	"go-admin/common/models"

)

type SysUserLog struct {
    models.Model
    
    DataId int `json:"dataId" gorm:"type:int;comment:关联的主键ID"` 
    Type string `json:"type" gorm:"type:varchar(10);comment:变更类型"` 
    Data string `json:"data" gorm:"type:json;comment:源数据"` 
    BeforeData string `json:"beforeData" gorm:"type:json;comment:变更前数据"` 
    AfterData string `json:"afterData" gorm:"type:json;comment:变更后数据"` 
    DiffData string `json:"diffData" gorm:"type:json;comment:差异数据"` 
	CreatedAt time.Time      `json:"createdAt" gorm:"comment:创建时间"`
	CreateBy     int    `json:"createBy" gorm:"index;comment:创建者"`
	CreateByName string `json:"createByName" gorm:"type:varchar(20);comment:创建人姓名"`
}

func (SysUserLog) TableName() string {
    return "sys_user_log"
}

func (e *SysUserLog) GetId() interface{} {
	return e.Id
}

func (e *SysUserLog) BeforeCreate(tx *gorm.DB) (err error) {
	if e.AfterData == "" {
		e.AfterData = "{}"
	}
	if e.Data == "" {
		e.Data = "{}"
	}
	if e.BeforeData == "" {
		e.BeforeData = "{}"
	}
	if e.DiffData == "" {
		e.DiffData = "[]"
	}
	return
}

func (e *SysUserLog) CreateLog(modelName string, tx *gorm.DB) error {
	logType := LogTypes{}
	beforeMap := map[string]interface{}{}
	afterMap := map[string]interface{}{}
	mappingMap := map[string]interface{}{}
	diffSlice := make(utils.OperateLogDetailResp, 0)

	err := tx.Where("model_name = ?", modelName).Where("type = ?", e.Type).Take(&logType).Error
	if err == nil {
		if e.AfterData != "" {
			if errAf := json.Unmarshal([]byte(e.AfterData), &afterMap); errAf == nil {
				if e.BeforeData != "" {
					_ = json.Unmarshal([]byte(e.BeforeData), &beforeMap)
				}
				if logType.Mapping != "" {
					_ = json.Unmarshal([]byte(logType.Mapping), &mappingMap)
				}
				diffData := utils.CompareDiff(beforeMap, afterMap)
				utils.FormatDiff(diffData, mappingMap, &diffSlice, "")
			}
		}
	}
	if diffBytes, err := json.Marshal(diffSlice); err == nil {
		e.DiffData = string(diffBytes)
	}
	//数据库操作
	if err := tx.Create(e).Error; err != nil {
		return err
	}
	return nil
}