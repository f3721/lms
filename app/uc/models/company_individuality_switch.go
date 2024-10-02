package models

import (
	"go-admin/common/global"
	"go-admin/common/models"
	"gorm.io/gorm"
)

type CompanyIndividualitySwitch struct {
	models.Model

	CompanyId    int    `json:"companyId" gorm:"type:int;comment:公司ID"`
	Keyword      string `json:"keyword" gorm:"type:varchar(45);comment:关键字"`
	SwitchStatus string `json:"switchStatus" gorm:"type:varchar(45);comment:开关状态，每个keyword的状态（是否打开或者下拉选择某个状态）"`
	Status       int    `json:"status" gorm:"type:tinyint;comment:状态 0关闭 1启用 默认1"`
	Desc         string `json:"desc" gorm:"type:varchar(100);comment:字段描述"`
	Sort         int    `json:"sort" gorm:"type:smallint;comment:排序 值越高 排名越前"`
	IsDel        int    `json:"isDel" gorm:"type:tinyint;comment:是否删除 0.否 1.是  默认0"`
	Remark       string `json:"remark" gorm:"type:varchar(255);comment:备注"`
	models.ControlBy
}

func (CompanyIndividualitySwitch) TableName() string {
	return "company_individuality_switch"
}

func (e *CompanyIndividualitySwitch) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *CompanyIndividualitySwitch) GetId() interface{} {
	return e.Id
}

func (e *CompanyIndividualitySwitch) GetRowByCompanyId(tx *gorm.DB, keyword string, companyId int) (err error) {
	ucPrefix := global.GetTenantUcDBNameWithDB(tx)
	err = tx.Table(ucPrefix+"."+e.TableName()).
		Where("keyword = ?", keyword).
		Where("company_id = ?", companyId).
		Where("status = ?", 1).
		First(e).Error
	return
}
