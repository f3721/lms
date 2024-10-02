package dto

type UserConfigSelectWarehouseReq struct {
	WarehouseID   int    `form:"warehouseId" comment:"仓库ID"  vd:"@:$>0;msg:'仓库warehouseId不能为空'"`
	WarehouseCode string `form:"warehouseCode" comment:"仓库的CODE"   vd:"@:len($)>0;msg:'仓库warehouseCode不能为空'"`
	WarehouseName string `form:"warehouseName" comment:"仓库的名称"   vd:"@:len($)>0;msg:'仓库warehouseName不能为空'"`
}
