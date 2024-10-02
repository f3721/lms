package models

import (
	"database/sql/driver"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ControlBy struct {
	CreateBy     int    `json:"createBy" gorm:"index;comment:创建者"`
	CreateByName string `json:"createByName" gorm:"index;comment:创建者姓名"`
	UpdateBy     int    `json:"updateBy" gorm:"index;comment:更新者姓名"`
	UpdateByName string `json:"updateByName" gorm:"index;comment:更新者姓名"`
}

// SetCreateBy 设置创建人id
func (e *ControlBy) SetCreateBy(createBy int) {
	e.CreateBy = createBy
}

// SetCreateByName 设置创建人姓名
func (e *ControlBy) SetCreateByName(createByName string) {
	e.CreateByName = createByName
}

// SetUpdateBy 设置修改人id
func (e *ControlBy) SetUpdateBy(updateBy int) {
	e.UpdateBy = updateBy
}

// SetUpdateByName 设置修改人姓名
func (e *ControlBy) SetUpdateByName(updateByName string) {
	e.UpdateByName = updateByName
}

type Model struct {
	Id int `json:"id" gorm:"primaryKey;autoIncrement;comment:主键编码"`
}

type ModelTime struct {
	CreatedAt time.Time      `json:"createdAt" gorm:"comment:创建时间"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"comment:最后更新时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
}

type RegionName struct {
	CityName     string `json:"cityName" gorm:"comment:市名"`
	ProvinceName string `json:"provinceName" gorm:"comment:省名"`
	DistrictName string `json:"districtName" gorm:"comment:区名"`
}

func (e *RegionName) GenerateRegionName(provinceId, cityId, district int, regionMap map[int]string) {
	e.ProvinceName = regionMap[provinceId]
	e.CityName = regionMap[cityId]
	e.DistrictName = regionMap[district]
}

// 逗号隔开的字符串
type Strs []string

func (m *Strs) Scan(val interface{}) error {
	s := val.([]uint8)
	ss := strings.Split(string(s), ",")
	*m = ss
	return nil
}

func (m Strs) Value() (driver.Value, error) {
	str := strings.Join(m, ",")
	return str, nil
}
