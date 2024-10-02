package dto

import (
	"go-admin/common/utils"
	"gorm.io/gorm"
	"strings"
)

type WarehousesSearch struct {
	QueryWarehouseCode string `form:"queryWarehouseCode"  search:"-"`
}

func GenWarehousesSearch(queryWarehouseCode, tableAlias string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		queryWarehouseCode = strings.Trim(queryWarehouseCode, " ")
		if queryWarehouseCode != "" {
			db.Where(tableAlias+".warehouse_code IN ?", utils.Split(strings.ReplaceAll(queryWarehouseCode, "，", ",")))
		}
		return db
	}
}
