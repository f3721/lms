package models

import (
	"go-admin/common/global"
	"go-admin/common/models"
	"gorm.io/gorm"
)

type SysDictData struct {
	DictCode  int    `json:"dictCode" gorm:"primaryKey;column:dict_code;autoIncrement;comment:主键编码"`
	DictSort  int    `json:"dictSort" gorm:"size:20;comment:DictSort"`
	DictLabel string `json:"dictLabel" gorm:"size:128;comment:DictLabel"`
	DictValue string `json:"dictValue" gorm:"size:255;comment:DictValue"`
	DictType  string `json:"dictType" gorm:"size:64;comment:DictType"`
	CssClass  string `json:"cssClass" gorm:"size:128;comment:CssClass"`
	ListClass string `json:"listClass" gorm:"size:128;comment:ListClass"`
	IsDefault string `json:"isDefault" gorm:"size:8;comment:IsDefault"`
	Status    int    `json:"status" gorm:"size:4;comment:Status"`
	Default   string `json:"default" gorm:"size:8;comment:Default"`
	Remark    string `json:"remark" gorm:"size:255;comment:Remark"`
	models.ControlBy
	models.ModelTime
}

func (SysDictData) TableName() string {
	return "sys_dict_data"
}

func (e *SysDictData) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *SysDictData) GetId() interface{} {
	return e.DictCode
}

func (e *SysDictData) GetDictDataByTypes (tx *gorm.DB, dictTypes []string) (res map[string]map[string]string) {
	adminPrefix := global.GetTenantAdminDBNameWithDB(tx)
	res = make(map[string]map[string]string)
	var list []SysDictData
	err := tx.Table(adminPrefix + "." + e.TableName()).Select("dict_label, dict_value, dict_type").Where("dict_type in ?", dictTypes).Order("dict_type asc, dict_sort desc").Find(&list).Error
	if err != nil {
		return
	}
	for _, data := range list {
		if _, ok := res[data.DictType]; !ok {
			res[data.DictType] = make(map[string]string)
		}
		res[data.DictType][data.DictValue] = data.DictLabel
	}
	return
}
