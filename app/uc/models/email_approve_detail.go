package models

import (
	"go-admin/common/models"
	"gorm.io/gorm"
)

const (
	EmailApproveDetailStatus1 = 1

	EmailApproveRankType1 = 1
	EmailApproveRankType2 = 2
	EmailApproveRankType3 = 3
)

type EmailApproveDetail struct {
	models.Model

	ApproveId           int      `json:"approveId" gorm:"type:int unsigned;comment:审批流ID"`
	UserId              int      `json:"userId" gorm:"type:int unsigned;comment:用户ID"`
	Priority            int      `json:"priority" gorm:"type:int unsigned;comment:审批者审批时所在的位置"`
	TotalPriceLimit     string   `json:"totalPriceLimit" gorm:"type:decimal(10,2);comment:采购限额"`
	LimitType           int      `json:"limitType" gorm:"type:tinyint(1);comment:采购限额类型(1:订单金额,2:商品金额)"`
	ApproveDetailStatus int      `json:"approveDetailStatus" gorm:"type:tinyint(1);comment:字审批流状态"`
	ApproveRankType     int      `json:"approveRankType" gorm:"type:tinyint(1);comment:审批层级类型：1：普签；2：会签；3：或签；默认为：1；"`
	User                UserInfo `json:"user" gorm:"foreignKey:UserId;"`
	models.ModelTime
	models.ControlBy
}

func (EmailApproveDetail) TableName() string {
	return "email_approve_detail"
}

func (e *EmailApproveDetail) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *EmailApproveDetail) GetId() interface{} {
	return e.Id
}

func (e *EmailApproveDetail) DeleteByApproveId(tx *gorm.DB, approveId int) error {
	return tx.Where("approve_id = ?", approveId).Delete(e).Error
}

func GetByApproveId(tx *gorm.DB, approveId int) (*[]EmailApproveDetail, error) {
	var list = &[]EmailApproveDetail{}
	err := tx.Where("approve_id = ?", approveId).Find(list).Error
	return list, err
}

func GetEmailApproveListByApproveId(tx *gorm.DB, approveId int) (*[]EmailApproveDetail, error) {
	var list = &[]EmailApproveDetail{}
	err := tx.
		Where("approve_id = ?", approveId).Where("approve_detail_status = ?", EmailApproveDetailStatus1).Group("priority").Find(list).Error
	return list, err
}

func GetEmailApproveListByApproveIdAndPriority(tx *gorm.DB, approveId, priority int) (*[]EmailApproveDetail, error) {
	var list = &[]EmailApproveDetail{}
	err := tx.Preload("User").Where("priority = ?", priority).
		Where("approve_id = ?", approveId).Where("approve_detail_status = ?", EmailApproveDetailStatus1).Order("id ASC").Find(list).Error
	return list, err
}
