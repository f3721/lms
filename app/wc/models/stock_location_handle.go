package models

import (
	"github.com/samber/lo"
	"gorm.io/gorm"
)

const (
	StockLocationStatus0   = "0"
	StockLocationStatus1   = "1"
	StockLocationDataLimit = 100
)

// 根据id获取库位
func (e *StockLocation) GetStockLocationById(tx *gorm.DB, id int) error {
	return tx.Take(e, id).Error
}

// 根据ids获取库位

func GetStockLocationByIds(tx *gorm.DB, ids []int) (*[]StockLocation, error) {
	stockLocations := &[]StockLocation{}
	ids = lo.Uniq(ids)
	if err := tx.Find(stockLocations, ids).Error; err != nil {
		return nil, err
	}
	return stockLocations, nil
}

func GetStockLocationMapByIds(tx *gorm.DB, ids []int) (map[int]string, error) {
	stockLocations, err := GetStockLocationByIds(tx, ids)
	if err != nil {
		return nil, err
	}
	aMap := lo.Associate(*stockLocations, func(s StockLocation) (int, string) {
		return s.Id, s.LocationCode
	})
	return aMap, nil
}

// 根据逻辑仓code获取库位limit
func GetStockLocationByLwhCodeWithLimit(tx *gorm.DB, lwhCode string) (*[]StockLocation, error) {
	return GetStockLocationByLwhCodeWithOptions(tx, lwhCode, func(tx *gorm.DB) *gorm.DB {
		return tx.Limit(StockLocationDataLimit)
	})
}

// 根据逻辑仓code获取库位 库位code模糊查询
func GetStockLocationByLwhCodeWithLocationCode(tx *gorm.DB, lwhCode, locationCode string) (*[]StockLocation, error) {
	stockLocations := &[]StockLocation{}
	query := tx.Where("logic_warehouse_code = ?", lwhCode).Debug().Where("status = ?", StockLocationStatus1)

	if locationCode != "" {
		query.Where("location_code LIKE ?", "%"+locationCode+"%")
	}

	err := query.Find(stockLocations).Error
	if err != nil {
		return nil, err
	}
	return stockLocations, nil
}

// 根据逻辑仓code获取库位
func GetStockLocationByLwhCodeWithOptions(tx *gorm.DB, lwhCode string, option func(*gorm.DB) *gorm.DB) (*[]StockLocation, error) {
	stockLocations := &[]StockLocation{}
	tx = option(tx)
	if err := tx.Where("logic_warehouse_code = ?", lwhCode).
		Where("status = ?", StockLocationStatus1).
		Find(stockLocations).Error; err != nil {
		return nil, err
	}
	return stockLocations, nil
}

// 从库位切片中根据id获取库位
func getLocationsFromStockLocationsById(stockLocations *[]StockLocation, id int) *StockLocation {
	for index, item := range *stockLocations {
		if item.Id == id {
			return &(*stockLocations)[index]
		}
	}
	return &StockLocation{}
}

// 根据id和逻辑仓code检查库位是否有效
func (e *StockLocation) CheckStockLocation(tx *gorm.DB, id int, logicWarehouseCode string) bool {
	if err := e.GetStockLocationById(tx, id); err != nil {
		return false
	}
	if e.Status != StockLocationStatus1 {
		return false
	}
	if e.LogicWarehouseCode != logicWarehouseCode {
		return false
	}
	return true
}

// 根据goodsId、逻辑仓code获取入库库位(有指定id)
func GetStockLocationsForEntryHasId(tx *gorm.DB, lwhCode string, goodsId, id int) (*[]StockLocation, bool, error) {
	if id == 0 {
		return GetStockLocationsForEntry(tx, lwhCode, goodsId)
	}
	stockLocations, err := GetStockLocationsHasId(tx, lwhCode, id)
	return stockLocations, false, err
}

// 根据goodsId、逻辑仓code获取库位(有指定id)
func GetStockLocationsHasId(tx *gorm.DB, lwhCode string, id int) (*[]StockLocation, error) {
	stockLocation := StockLocation{}
	if err := tx.Take(&stockLocation, id).Error; err != nil {
		return nil, err
	}
	stockLocations, err := GetStockLocationByLwhCodeWithOptions(tx, lwhCode, func(tx *gorm.DB) *gorm.DB {
		return tx.Where("id != ?", id).Limit(StockLocationDataLimit)
	})
	if err != nil {
		return nil, err
	}
	*stockLocations = append([]StockLocation{stockLocation}, *stockLocations...)
	return stockLocations, err
}

// 根据goodsId、逻辑仓code获取入库库位
func GetStockLocationsForEntry(tx *gorm.DB, lwhCode string, goodsId int) (*[]StockLocation, bool, error) {
	// 找到当前逻辑仓下所有有效库位
	stockLocations, err := GetStockLocationByLwhCodeWithLimit(tx, lwhCode)
	if err != nil {
		return nil, false, err
	}
	stockLocationGoods, err := GetLocationGoodsByGoodsIdAndLwhCode(tx, goodsId, lwhCode)
	if err != nil {
		return nil, false, err
	}
	//置顶库位
	topLocation := &StockLocation{}
	//大于0库位
	gtZeroLocations := []int{}
	//等于0库位
	zeroLocations := []int{}
	for _, item := range *stockLocationGoods {
		if item.Stock > 0 {
			gtZeroLocations = append(gtZeroLocations, item.LocationId)
		}
		if item.Stock == 0 {
			zeroLocations = append(zeroLocations, item.LocationId)
		}
	}
	if len(gtZeroLocations) > 0 {
		topLocation = getLocationsFromStockLocationsById(stockLocations, gtZeroLocations[0])
	} else if len(zeroLocations) > 0 {
		topLocation = getLocationsFromStockLocationsById(stockLocations, zeroLocations[0])
	}
	if topLocation.Id != 0 {
		filterStockLocations := lo.Filter(*stockLocations, func(item StockLocation, index int) bool {
			return item.Id != topLocation.Id
		})
		outDData := append([]StockLocation{*topLocation}, filterStockLocations...)
		return &outDData, true, nil
	}

	return stockLocations, false, nil
}

// 通过指定goodsId和库位信息获取相关库位库存信息

func GetLocationGoodsByGoodsIdAndLocationIdSlice(tx *gorm.DB, goodsId int, locationIds []int) (*[]StockLocationGoods, error) {
	stockLocationGoods := &[]StockLocationGoods{}
	if err := tx.Where("goods_id = ?", goodsId).Where("location_id in ?", locationIds).Order("stock asc,updated_at desc").Find(stockLocationGoods).Error; err != nil {
		return nil, err
	}
	return stockLocationGoods, nil
}

// 通过指定goodsId和逻辑仓code获取相关库位库存信息

func GetLocationGoodsByGoodsIdAndLwhCode(tx *gorm.DB, goodsId int, lwhCode string) (*[]StockLocationGoods, error) {
	stockLocationGoods := &[]StockLocationGoods{}

	if err := tx.Joins("JOIN stock_location location ON location.id = stock_location_goods.location_id").
		Select("stock_location_goods.*").
		Where("stock_location_goods.goods_id = ?", goodsId).
		Where("location.logic_warehouse_code = ?", lwhCode).
		Where("location.status = ?", StockLocationStatus1).
		Order("stock_location_goods.stock asc,stock_location_goods.updated_at desc").
		Find(stockLocationGoods).Error; err != nil {
		return nil, err
	}
	return stockLocationGoods, nil
}
