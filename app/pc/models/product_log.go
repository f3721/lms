package models

import (
	"encoding/json"
	"go-admin/common/models"
	"go-admin/common/utils"
	"gorm.io/gorm"
)

type ProductLog struct {
	models.Model

	DataId     int    `json:"dataId" gorm:"type:int;comment:关联的主键ID"`
	Type       string `json:"type" gorm:"type:varchar(10);comment:变更类型"`
	Data       string `json:"data" gorm:"type:json;comment:源数据"`
	BeforeData string `json:"beforeData" gorm:"type:json;comment:变更前数据"`
	AfterData  string `json:"afterData" gorm:"type:json;comment:变更后数据"`
	DiffData   string `json:"diffData" gorm:"type:json;comment:差异数据"`
	models.ModelTime
	models.ControlBy
}

func (ProductLog) TableName() string {
	return "product_log"
}

func (e *ProductLog) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *ProductLog) GetId() interface{} {
	return e.Id
}

func (e *ProductLog) BeforeCreate(tx *gorm.DB) (err error) {
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

func (e *ProductLog) CreateLog(modelName string, tx *gorm.DB) error {
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
	//更新没有变更字段不插入日志
	if e.Type == "update" && len(diffSlice) == 0 {
		return nil
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
