package models

import (
     
     

	"go-admin/common/models"

)

type DataRemark struct {
    models.Model
    
    Type string `json:"type" gorm:"type:varchar(50);comment:类型{sale_order:销售售后，purchase_order:采购售后}"`// 类型{sale_order:销售售后，purchase_order:采购售后} 
    DataId string `json:"dataId" gorm:"type:varchar(30);comment:DataId"`//  
    Remark string `json:"remark" gorm:"type:text;comment:备注"`// 备注 
    Usertype int `json:"usertype" gorm:"type:tinyint(1);comment:操作人员类型(0:opc,1:spc,2:other)"`// 操作人员类型(0:opc,1:spc,2:other) 
    models.ModelTime
    models.ControlBy
}

func (DataRemark) TableName() string {
    return "data_remark"
}

func (e *DataRemark) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *DataRemark) GetId() interface{} {
	return e.Id
}
