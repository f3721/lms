package dto

import (
	"errors"
	"go-admin/app/wc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
	"go-admin/common/utils"
	"gorm.io/gorm"
	"regexp"
	"strings"
)

type StockLocationGetPageReq struct {
	dto.Pagination     `search:"-"`
	LocationCode string `search:"-" form:"locationCode"`
	WarehouseCode 	string `form:"warehouseCode"  search:"type:exact;column:warehouse_code;table:stock_location"`
	LogicWarehouseCode 	string `form:"logicWarehouseCode"  search:"type:exact;column:logic_warehouse_code;table:stock_location"`
	Status 	string `form:"status"  search:"type:exact;column:status;table:stock_location"`
	StockLocationOrder
}

type StockLocationOrder struct {
	Id string `form:"idOrder"  search:"type:order;column:id;table:stock_location"`

}

func (m *StockLocationGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type StockLocationInsertReq struct {
	Id int `json:"-" comment:"id"` // id
	LocationCode string `json:"locationCode" comment:"库位编码"`
	WarehouseCode string `json:"warehouseCode" comment:"实体仓库编码"`
	LogicWarehouseCode string `json:"logicWarehouseCode" comment:"逻辑仓库编码"`
	Status string `json:"status" comment:"是否启用"`
	//SizeHeight string `json:"sizeHeight" comment:"高"`
	//SizeLength string `json:"sizeLength" comment:"长"`
	//SizeWidth string `json:"sizeWidth" comment:"宽"`
	//Capacity string `json:"capacity" comment:"容量"`
	Remark string `json:"remark" comment:"备注"`
	common.ControlBy
}

func (s *StockLocationInsertReq) InsertValid(tx *gorm.DB) (err error) {
	if s.WarehouseCode  == "" {
		return errors.New("实体仓必填")
	}

	if s. LogicWarehouseCode == "" {
		return errors.New("逻辑仓必填")
	}

	var logicWarehouse models.LogicWarehouse
	err = logicWarehouse.GetLogicWarehouseByCode(tx, s.LogicWarehouseCode)
	if err != nil {
		return
	}
	if logicWarehouse.WarehouseCode != s.WarehouseCode {
		return errors.New("实体仓和逻辑仓不匹配")
	}

	if s.LocationCode == "" {
		return errors.New("库位编号必填")
	} else {
		reg := regexp.MustCompile(`^[A-Za-z]{2}[\d-]{1,18}$`)
		if !reg.MatchString(s.LocationCode) {
			return errors.New("库位编号格式校验失败")
		}

		s.LocationCode = strings.ToUpper(s.LocationCode)
		var object models.StockLocation

		if flag := object.CheckLocationCodeExist(tx, s.LocationCode); flag == true {
			err = errors.New("库位编号已存在")
		}
	}

	return
}

func (s *StockLocationUpdateReq) UpdateValid(tx *gorm.DB) (err error) {
	var object models.StockLocationGoods
	if s.Status == "0" {
		s.LocationCode = strings.ToUpper(s.LocationCode)
		if total := object.GetTotalStockByLocationCode(tx, s.LocationCode); total > 0 {
			err = errors.New("库位下无商品才可停用")
		}
	}
	return
}

func (s *StockLocationInsertReq) Generate(model *models.StockLocation)  {
	if s.Id == 0 {
		model.Model = common.Model{ Id: s.Id }
	}
	model.LocationCode = s.LocationCode
	model.WarehouseCode = s.WarehouseCode
	model.LogicWarehouseCode = s.LogicWarehouseCode
	model.Status = s.Status
	//model.SizeHeight = s.SizeHeight
	//model.SizeLength = s.SizeLength
	//model.SizeWidth = s.SizeWidth
	//model.Capacity = s.Capacity
	model.Remark = s.Remark
}

func (s *StockLocationInsertReq) GetId() interface{} {
	return s.Id
}

type StockLocationUpdateReq struct {
	Id int `uri:"id" comment:"id"` // id
	LocationCode string `json:"locationCode" comment:"库位编码"`
	WarehouseCode string `json:"warehouseCode" comment:"实体仓库编码"`
	LogicWarehouseCode string `json:"logicWarehouseCode" comment:"逻辑仓库编码"`
	Status string `json:"status" comment:"是否启用"`
	SizeHeight string `json:"sizeHeight" comment:"高"`
	SizeLength string `json:"sizeLength" comment:"长"`
	SizeWidth string `json:"sizeWidth" comment:"宽"`
	Capacity string `json:"capacity" comment:"容量"`
	Remark string `json:"remark" comment:"备注"`
	common.ControlBy
}

func (s *StockLocationUpdateReq) Generate(model *models.StockLocation)  {
	if s.Id == 0 {
		model.Model = common.Model{ Id: s.Id }
	}
	model.LocationCode = s.LocationCode
	model.WarehouseCode = s.WarehouseCode
	model.LogicWarehouseCode = s.LogicWarehouseCode
	model.Status = s.Status
	//model.SizeHeight = s.SizeHeight
	//model.SizeLength = s.SizeLength
	//model.SizeWidth = s.SizeWidth
	//model.Capacity = s.Capacity
	model.Remark = s.Remark
}

func (s *StockLocationUpdateReq) GetId() interface{} {
	return s.Id
}

// StockLocationGetReq 功能获取请求参数
type StockLocationGetReq struct {
	Id int `uri:"id"`
}
func (s *StockLocationGetReq) GetId() interface{} {
	return s.Id
}

// StockLocationDeleteReq 功能删除请求参数
type StockLocationDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *StockLocationDeleteReq) GetId() interface{} {
	return s.Ids
}

type StockLocationResp struct {
	Id int `json:"id" comment:"id"` // id
	LocationCode string `json:"locationCode" comment:"库位编码"`
	WarehouseCode string `json:"warehouseCode" comment:"实体仓库编码"`
	LogicWarehouseCode string `json:"logicWarehouseCode" comment:"逻辑仓库编码"`
	WarehouseName string `json:"warehouseName" comment:"实体仓名称"`
	LogicWarehouseName string `json:"logicWarehouseName" comment:"逻辑仓名称"`
	Status string `json:"status" comment:"是否启用"`
	Remark string `json:"remark" comment:"备注"`
}

func StockLocationGetPageMakeCondition(c *StockLocationGetPageReq) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if c.LocationCode != "" {
			locationCode := utils.Split(c.LocationCode)
			length := len(locationCode)
			if length > 1 {
				db.Where("location_code in ?", locationCode)
			} else if length == 1 {
				db.Where("location_code like ?", "%"+c.LocationCode+"%")
			}
		}
		return db
	}
}
