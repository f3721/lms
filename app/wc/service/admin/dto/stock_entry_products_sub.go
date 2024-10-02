package dto

import (
	"go-admin/app/wc/models"
	common "go-admin/common/models"
	"time"
)

type StockEntryProductsSubReq struct {
	Id               int       `json:"id"`
	EntryCode        string    `json:"entryCode" comment:"入库单code"`
	EntryProductId   int       `json:"entryProductId" comment:"出库单产品ID"`
	ShouldQuantity   int       `json:"shouldQuantity" comment:"应入数量"`
	StockLocationId  int       `json:"stockLocationId" comment:"库位ID"`
	ActQuantity      int       `json:"actQuantity" comment:"实际入库数量"`
	StashLocationId  int       `json:"stashLocationId" comment:"暂存库位id"`
	StashActQuantity int       `json:"stashActQuantity" comment:"暂存数量"`
	EntryTime        time.Time `json:"entryTime" comment:"入库时间"`
	common.ControlBy
}

func (s *StockEntryProductsSubReq) Generate(model *models.StockEntryProductsSub) {
	model.EntryCode = s.EntryCode
	model.EntryProductId = s.EntryProductId
	model.ShouldQuantity = s.ShouldQuantity
	model.StockLocationId = s.StockLocationId
	model.ActQuantity = s.ActQuantity
	model.StashLocationId = s.StashLocationId
	model.StashActQuantity = s.StashActQuantity
}
