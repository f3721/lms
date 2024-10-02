package models

import (
	"errors"
	"fmt"
	"go-admin/common/global"
	"go-admin/common/models"

	"gorm.io/gorm"
)

const (
	LwhType0 = "0"
	LwhType1 = "1"

	LogicWarehouseModelName   = "logicWarehouse"
	LogicWarehouseModeInsert  = "insert"
	LogicWarehouseModeUpdate  = "update"
	LogicWarehouseModeStatus0 = "0"
	LogicWarehouseModeStatus1 = "1"
)

var LwhTypeMap = map[string]string{
	LwhType0: "正品仓",
	LwhType1: "次品仓",
}

var LogicWarehouseModeStatus = map[string]string{
	LogicWarehouseModeStatus0: "无效",
	LogicWarehouseModeStatus1: "有效",
}

type LogicWarehouse struct {
	models.Model

	LogicWarehouseCode string    `json:"logicWarehouseCode" gorm:"type:varchar(20);comment:逻辑仓库编码"`
	LogicWarehouseName string    `json:"logicWarehouseName" gorm:"type:varchar(20);comment:逻辑仓库名称"`
	WarehouseCode      string    `json:"warehouseCode" gorm:"type:varchar(20);comment:逻辑仓库对应实体仓code"`
	Mobile             string    `json:"mobile" gorm:"type:varchar(20);comment:Mobile"`
	Linkman            string    `json:"linkman" gorm:"type:varchar(50);comment:联系人"`
	Email              string    `json:"email" gorm:"type:varchar(50);comment:邮箱"`
	Type               string    `json:"type" gorm:"type:tinyint(1);comment:0 正品仓 1次品仓"`
	Status             string    `json:"status" gorm:"type:tinyint(1) unsigned;comment:是否使用 0-否，1-是"`
	Remark             string    `json:"remark" gorm:"type:text;comment:Remark"`
	Warehouse          Warehouse `json:"-" gorm:"foreignKey:WarehouseCode;references:WarehouseCode"`
	models.ModelTime
	models.ControlBy
}

func (LogicWarehouse) TableName() string {
	return "logic_warehouse"
}

func (e *LogicWarehouse) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *LogicWarehouse) GetId() interface{} {
	return e.Id
}

func (e *LogicWarehouse) GenerateLwhCode(tx *gorm.DB) (string, error) {
	var result struct {
		MaxId int
	}
	if err := tx.Model(e).Select("max(id) as MaxId").Scan(&result).Error; err != nil {
		return "", err
	}
	code := fmt.Sprintf("%05d", result.MaxId+1)
	e.LogicWarehouseCode = code
	return code, nil
}

func (e *LogicWarehouse) CheckLwhName(tx *gorm.DB, name string, id int) (bool, error) {
	var result struct {
		Id int
	}
	tx = tx.Model(&LogicWarehouse{}).Select("id")
	if id != 0 {
		tx.Where("id <> ?", id)
	}
	if err := tx.Where("logic_warehouse_name = ?", name).Scan(&result).Error; err != nil {
		return false, err
	}
	if result.Id != 0 {
		return false, nil
	}
	return true, nil
}

func (e *LogicWarehouse) CheckLwhExist(tx *gorm.DB, whCode, Type string, id int) (bool, error) {
	var result struct {
		Id int
	}
	tx = tx.Model(&LogicWarehouse{}).Select("id")
	if id != 0 {
		tx.Where("id <> ?", id)
	}
	if err := tx.Where("warehouse_code = ?", whCode).Where("type = ?", Type).Where("status = 1").Scan(&result).Error; err != nil {
		return false, err
	}
	if result.Id != 0 {
		return false, nil
	}
	return true, nil
}

func (e *LogicWarehouse) GetLogicWharehouseNameByCodes(tx *gorm.DB, codes []string) ([]map[string]string, error) {
	var result = []map[string]string{}

	tx = tx.Model(&LogicWarehouse{}).Select("logic_warehouse_code,logic_warehouse_name")

	if err := tx.Where("logic_warehouse_code in ?", codes).Scan(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (e *LogicWarehouse) GetLogicWarehouseByCode(tx *gorm.DB, lwhCode string) error {
	err := tx.Model(&LogicWarehouse{}).Preload("Warehouse").Where("logic_warehouse_code = ?", lwhCode).First(e).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("逻辑仓不存在")
	}
	return err
}

func (e *LogicWarehouse) GetWhAndLwhInfo(tx *gorm.DB, whCode, lwhCode string) error {
	if err := e.GetLogicWarehouseByCode(tx, lwhCode); err != nil {
		return err
	}
	if e.Warehouse.WarehouseCode == "" || e.Warehouse.WarehouseCode != whCode {
		return errors.New("实体仓和逻辑仓不匹配")
	}
	if e.Warehouse.Status != "1" {
		return errors.New("实体仓状态不正确")
	}
	if e.Status != "1" {
		return errors.New("逻辑仓状态不正确")
	}
	//仓库权限检查 todo
	return nil
}

// 正品逻辑仓 -> 次品仓逻辑仓  次品仓逻辑仓 ->正品逻辑仓
func (e *LogicWarehouse) GetDefectiveOrPassedLogicWarehouse(tx *gorm.DB, lwhCode, Type string) error {
	oriLogicWarehouse := &LogicWarehouse{}
	if err := oriLogicWarehouse.GetLogicWarehouseByCode(tx, lwhCode); err != nil {
		return err
	}
	if err := tx.Where("warehouse_code = ?", oriLogicWarehouse.WarehouseCode).Where("type = ?", Type).Where("status = ?", LogicWarehouseModeStatus1).Take(e).Error; err != nil {
		return err
	}
	return nil
}

func (e *LogicWarehouse) GetLogicWarehouseByName(tx *gorm.DB, lwhName string, ruleWarehouseCodes []string) error {
	tx = tx.Model(LogicWarehouse{}).Where("logic_warehouse_name = ?", lwhName)
	if ruleWarehouseCodes != nil {
		tx = tx.Where("warehouse_code in ?", ruleWarehouseCodes)
	}
	err := tx.First(e).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("未找到权限内的实体仓对应的逻辑仓库")
	}
	return err
}

func (e *LogicWarehouse) GetPassLogicWarehouseByWhCode(tx *gorm.DB, whCode string) error {
	wcPrefix := global.GetTenantWcDBNameWithDB(tx)
	err := tx.Table(wcPrefix+"."+e.TableName()).Where("status = ?", LogicWarehouseModeStatus1).Where("type = ?", LwhType0).Where("warehouse_code = ?", whCode).First(e).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("正品仓不存在")
	}
	return err
}

func (e *LogicWarehouse) GetLogicWarehouseByWarehouseName(tx *gorm.DB, logicWarehouseName string) (logicWarehouse *LogicWarehouse, err error) {
	err = tx.Where("logic_warehouse_name = ?", logicWarehouseName).First(&logicWarehouse).Error
	if logicWarehouse.Id != 0 {
		return logicWarehouse, nil
	}
	return logicWarehouse, err
}
