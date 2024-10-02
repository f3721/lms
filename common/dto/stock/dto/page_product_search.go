package dto

import (
	"go-admin/common/utils"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ProductSearch struct {
	Sku         string `form:"sku"  search:"-"`
	ProductNo   string `form:"productNo"  search:"-"`
	ProductName string `form:"productName"  search:"-"`
}

func GenProductSearch(sku, productNo, productName, tableAlias string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		sku = strings.Trim(sku, " ")
		if sku != "" {
			db.Where(tableAlias+".sku_code IN ?", utils.Split(sku))
		}
		if productNo != "" {
			db.Where(" goods.product_no IN ?", utils.Split(productNo))
		}
		if productName != "" {
			//nameZh := utils.Split(productName)
			//if len(nameZh) > 1 {
			//	db.Where(" product.name_zh IN ?", nameZh)
			//} else {
			//	db.Where(" product.name_zh LIKE ?", "%"+productName+"%")
			//}
			db.Where(" product.name_zh LIKE ?", "%"+productName+"%")
		}
		return db
	}
}

func GenRecipientSearch(userName string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if userName != "" {
			db.Where("stock_outbound.type = 1")
			db.Where("oi.user_name LIKE ?", "%"+userName+"%")
		}
		return db
	}
}

func GenEntryRecipientSearch(userName string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if userName != "" {
			db.Where("stock_entry.type = 1")
			db.Where("oi.user_name LIKE ?", "%"+userName+"%")
		}
		return db
	}
}

func GenIdsSearch(ids []int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(ids) > 0 {
			db.Where("stock_outbound.id in ?", ids)
		}
		return db
	}
}

func GenEntryIdsSearch(ids []int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(ids) > 0 {
			db.Where("stock_entry.id in ?", ids)
		}
		return db
	}
}

func GenCreatedAtTimeSearch(createdAtStart, createdAtEnd time.Time, tableAlias string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if !createdAtStart.IsZero() {
			db.Where(tableAlias+".created_at >= ?", createdAtStart)
		}
		if !createdAtEnd.IsZero() {
			db.Where(tableAlias+".created_at <= ?", createdAtEnd)
		}
		return db
	}
}
