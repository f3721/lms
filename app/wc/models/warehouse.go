package models

import (
	"errors"
	"fmt"
	"go-admin/common/global"
	"go-admin/common/models"

	"github.com/samber/lo"
	"gorm.io/gorm"
)

const (
	WarehouseModelName      = "warehouse"
	WarehouseModeInsert     = "insert"
	WarehouseModeUpdate     = "update"
	WarehouseModeStatus0    = "0"
	WarehouseModeStatus1    = "1"
	WarehouseModeIsVirtual0 = "0"
	WarehouseModeIsVirtual1 = "1"
)

var WarehouseModeStatus = map[string]string{
	WarehouseModeStatus0: "无效",
	WarehouseModeStatus1: "有效",
}

var WarehouseModeIsVirtualMap = map[string]string{
	WarehouseModeIsVirtual0: "否",
	WarehouseModeIsVirtual1: "是",
}

type Warehouse struct {
	models.Model

	WarehouseCode string `json:"warehouseCode" gorm:"type:varchar(20);comment:仓库编码"`
	WarehouseName string `json:"warehouseName" gorm:"type:varchar(20);comment:仓库名称"`
	CompanyId     int    `json:"companyId" gorm:"type:int(10);comment:仓库对应公司d"`
	Mobile        string `json:"mobile" gorm:"type:varchar(20);comment:Mobile"`
	Linkman       string `json:"linkman" gorm:"type:varchar(50);comment:联系人"`
	Email         string `json:"email" gorm:"type:varchar(50);comment:邮箱"`
	Status        string `json:"status" gorm:"type:tinyint(1) unsigned;comment:是否使用 0-否，1-是"`
	IsVirtual     string `json:"isVirtual" gorm:"type:tinyint(1);comment:是否为虚拟仓 0-否，1-是"`
	PostCode      string `json:"postCode" gorm:"type:varchar(50);comment:仓库所在地址邮编"`
	Province      int    `json:"province" gorm:"type:int unsigned;comment:省"`
	City          int    `json:"city" gorm:"type:int unsigned;comment:市"`
	District      int    `json:"district" gorm:"type:int unsigned;comment:区"`
	Address       string `json:"address" gorm:"type:varchar(100);comment:地址"`
	Remark        string `json:"remark" gorm:"type:varchar(255);comment:Remark"`
	models.RegionName
	models.ModelTime
	models.ControlBy
}

func (Warehouse) TableName() string {
	return "warehouse"
}

func (e *Warehouse) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *Warehouse) GetId() interface{} {
	return e.Id
}

//自增生成实体仓code

func (e *Warehouse) GenerateWhCode(tx *gorm.DB) (string, error) {
	var result struct {
		MaxId int
	}
	if err := tx.Model(e).Select("max(id) as MaxId").Scan(&result).Error; err != nil {
		return "", err
	}
	code := "WH" + fmt.Sprintf("%04d", result.MaxId+1)
	e.WarehouseCode = code
	return code, nil
}

//查重实体仓name

func (e *Warehouse) CheckWhName(tx *gorm.DB, name string, id int) (bool, error) {
	var result struct {
		Id int
	}
	tx = tx.Model(&Warehouse{}).Select("id")
	if id != 0 {
		tx.Where("id <> ?", id)
	}
	if err := tx.Where("warehouse_name = ?", name).Scan(&result).Error; err != nil {
		return false, err
	}
	if result.Id != 0 {
		return false, nil
	}
	return true, nil
}

func (e *Warehouse) GetWharehouseNameByCodes(tx *gorm.DB, codes []string) ([]map[string]string, error) {
	var result = []map[string]string{}

	tx = tx.Model(&Warehouse{}).Select("warehouse_code,warehouse_name")

	if err := tx.Where("warehouse_code in ?", codes).Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

// 根据code获取实体仓信息
func (e *Warehouse) GetByWarehouseCode(tx *gorm.DB, code string) error {
	wcPrefix := global.GetTenantWcDBNameWithDB(tx)
	err := tx.Table(wcPrefix+"."+e.TableName()).Where("warehouse_code = ?", code).Take(e).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("实体仓不存在")
	}
	return err
}

// 判断实体仓是否为虚拟仓
func (e *Warehouse) CheckIsVirtual() bool {
	return e.IsVirtual == WarehouseModeIsVirtual1
}

// 根据name获取实体仓信息
func (e *Warehouse) GetByWarehouseName(tx *gorm.DB, name string, ruleWarehouseCodes []string) error {
	tx = tx.Where("warehouse_name = ?", name)
	if ruleWarehouseCodes != nil {
		tx = tx.Where("warehouse_code in ?", ruleWarehouseCodes)
	}
	err := tx.Take(e).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("未找到权限内的实体仓库")
	}
	return err
}

func (e *Warehouse) GetWarehoseByName(tx *gorm.DB, warehouseName string) (*Warehouse, error) {
	var warehouse = &Warehouse{}
	err := tx.Where("warehouse_name = ?", warehouseName).First(&warehouse).Error
	if warehouse.Id != 0 {
		return warehouse, nil
	}
	return warehouse, err
}

// 根据Codes查询列表
func (e *Warehouse) ListByCodes(tx *gorm.DB, warehouseCodes []string) ([]*Warehouse, error) {
	list := []*Warehouse{}
	err := tx.Where("warehouse_code in (?)", warehouseCodes).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, err
}

// 根据Codes查询Map
func (e *Warehouse) MapByCodes(tx *gorm.DB, warehouseCodes []string) (map[string]*Warehouse, error) {
	// 查列表
	list, err := e.ListByCodes(tx, warehouseCodes)
	if err != nil {
		return nil, err
	}

	// 组map
	res := lo.Associate(list, func(item *Warehouse) (string, *Warehouse) {
		return item.WarehouseCode, item
	})

	return res, nil
}
