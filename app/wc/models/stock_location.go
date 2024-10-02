package models

import (
	"errors"
	"go-admin/common/models"

	"gorm.io/gorm"
)

type StockLocation struct {
	models.Model

	LocationCode       string `json:"locationCode" gorm:"type:varchar(20);comment:库位编码"`
	WarehouseCode      string `json:"warehouseCode" gorm:"type:varchar(64);comment:实体仓库编码"`
	LogicWarehouseCode string `json:"logicWarehouseCode" gorm:"type:varchar(64);comment:逻辑仓库编码"`
	Status             string `json:"status" gorm:"type:tinyint(1) unsigned;comment:是否启用"`
	SizeHeight         string `json:"sizeHeight" gorm:"type:int(11) unsigned;comment:高"`
	SizeLength         string `json:"sizeLength" gorm:"type:int(11) unsigned;comment:长"`
	SizeWidth          string `json:"sizeWidth" gorm:"type:int(11) unsigned;comment:宽"`
	Capacity           string `json:"capacity" gorm:"type:int(11) unsigned;comment:容量"`
	Remark             string `json:"remark" gorm:"type:varchar(500);comment:备注"`
	models.ModelTime
	models.ControlBy
}

func (StockLocation) TableName() string {
	return "stock_location"
}

func (e *StockLocation) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *StockLocation) GetId() interface{} {
	return e.Id
}

func (e *StockLocation) CheckLocationCodeExist(tx *gorm.DB, locationCode string) bool {
	var data StockLocation

	err := tx.Table(e.TableName()).Where("location_code = ?", locationCode).Where("status = 1").First(&data).Error

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}

func (e *StockLocation) GetByLocationCode(tx *gorm.DB, locationCode string) error {
	return tx.Table(e.TableName()).Where("location_code = ?", locationCode).Where("status = 1").First(e).Error
}

// 根据逻辑仓编码获取库位列表
func (e *StockLocation) StockLocationListByLwsCode(tx *gorm.DB, LwsCode string) (*[]*StockLocation, error) {
	list := &[]*StockLocation{}
	err := tx.Table(e.TableName()).Where("logic_warehouse_code = ?", LwsCode).Where("status = 1").Find(list).Error
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (e *StockLocation) Find(tx *gorm.DB) []StockLocation {
	var list []StockLocation
	err := tx.Table(e.TableName()).Find(&list).Error
	if err != nil {
		return nil
	}
	return list
}
