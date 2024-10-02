package models

import (
	"time"

	"go-admin/common/models"
)

type OrderImage struct {
    models.Model
    
    Name string `json:"name" gorm:"type:varchar(45);comment:图片名称"` 
    OrderId string `json:"orderId" gorm:"type:varchar(30);comment:订单编号"` 
    Url string `json:"url" gorm:"type:varchar(255);comment:图片地址"` 
    Type int `json:"type" gorm:"type:tinyint(1);comment:类型: 0-订单相关文件 1-订单签收文件"`
	CreatedAt time.Time      `json:"createdAt" gorm:"comment:创建时间"`
}

func (OrderImage) TableName() string {
    return "order_image"
}

func (e *OrderImage) GetId() interface{} {
	return e.Id
}
