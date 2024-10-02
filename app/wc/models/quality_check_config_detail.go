package models

import "go-admin/common/models"

type QualityCheckConfigDetail struct {
	models.Model
	ConfigId      int    `json:"configId"`      //  配置表id
	CompanyId     int    `json:"companyId"`     //  公司id
	CompanyName   string `json:"companyName"`   //  公司名称
	WarehouseCode string `json:"warehouseCode"` //  仓库code
	WarehouseName string `json:"warehouseName"` //  实体仓名称
	OrderType     int    `json:"orderType"`     //  1-采购入库 2-大货入库
	models.ModelTime
}

func (QualityCheckConfigDetail) TableName() string {
	return "quality_check_config_detail"
}

func (e *QualityCheckConfigDetail) GetId() interface{} {
	return e.Id
}
