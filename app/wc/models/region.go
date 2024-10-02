package models

import (
	"github.com/samber/lo"
	"go-admin/common/models"
	"gorm.io/gorm"
)

type Region struct {
	models.Model

	ParentId   string `json:"parentId" gorm:"type:int(11);comment:ParentId"`
	Name       string `json:"name" gorm:"type:varchar(50);comment:Name"`
	Level      string `json:"level" gorm:"type:int(11);comment:Level"`
	PostalCode string `json:"postalCode" gorm:"type:varchar(50);comment:PostalCode"`
	Latitude   string `json:"latitude" gorm:"type:varchar(50);comment:纬度"`
	Longitude  string `json:"longitude" gorm:"type:varchar(50);comment:经度"`
	Adcode     string `json:"adcode" gorm:"type:int(11);comment:Adcode"`
}

func (Region) TableName() string {
	return "region"
}

func (e *Region) GetId() interface{} {
	return e.Id
}

func (e *Region) GetByIds(tx *gorm.DB, ids []int) *[]Region {
	regions := &[]Region{}
	tx.Find(regions, ids)
	return regions
}

func GeRegionMapByIds(tx *gorm.DB, regionIds []int) map[int]string {
	region := &Region{}
	regionIds = lo.Uniq(regionIds)
	regions := region.GetByIds(tx, regionIds)
	return lo.Associate(*regions, func(s Region) (int, string) {
		return s.Id, s.Name
	})
}
