package models

import (
	"go-admin/common/global"
	"go-admin/common/models"

	"github.com/samber/lo"
	"gorm.io/gorm"
)

type UserApprove struct {
	models.Model

	UserId    int      `json:"userId" gorm:"type:int unsigned;comment:用户ID"`
	ApproveId int      `json:"approveId" gorm:"type:int unsigned;comment:审批流ID"`
	Status    int      `json:"status" gorm:"type:tinyint(1);comment:状态"`
	User      UserInfo `json:"user" gorm:"foreignKey:UserId;references:Id"`
	models.ModelTime
	models.ControlBy
}

const (
	UserApproveStatus1 = 1
)

func (UserApprove) TableName() string {
	return "user_approve"
}

func (e *UserApprove) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *UserApprove) GetId() interface{} {
	return e.Id
}

func (e *UserApprove) DeleteByApproveId(tx *gorm.DB, approveId int) error {
	return tx.Where("approve_id = ?", approveId).Delete(e).Error
}

func GetUserApproveListById(tx *gorm.DB, approveId int) (*[]UserApprove, error) {
	var list = &[]UserApprove{}
	err := tx.Model(&UserApprove{}).Preload("User").
		Where("approve_id = ?", approveId).Where("status = ?", UserApproveStatus1).Find(list).Error
	return list, err
}

func (e *UserApprove) GetApproveListByUserId(tx *gorm.DB, userId int) (*[]EmailApproveDetail, error) {
	var list = &[]EmailApproveDetail{}
	err := tx.Model(e).
		Joins("INNER JOIN email_approve ON user_approve.approve_id = email_approve.id").
		Joins("INNER JOIN email_approve_detail ON email_approve_detail.approve_id = email_approve.id").
		Where("user_approve.user_id = ?", userId).Where("user_approve.status = ?", UserApproveStatus1).
		Where("email_approve.approve_type = ?", EmailApproveType9).Where("email_approve.approve_status = ?", EmailApproveStatus1).
		Where("email_approve_detail.approve_detail_status = ?", EmailApproveDetailStatus1).
		Group("user_approve.approve_id,email_approve_detail.priority").
		Order("user_approve.approve_id ASC,email_approve_detail.priority ASC").
		Select("email_approve_detail.*").
		Find(list).Error
	return list, err
}

func (e *UserApprove) GetApproveIdsByUserId(tx *gorm.DB, approveId int) ([]int, error) {
	var list = &[]UserApprove{}
	if err := tx.Model(e).Where("user_id = ?", approveId).Where("status = ?", UserApproveStatus1).Find(list).Error; err != nil {
		return []int{}, err
	}
	return lo.Map(*list, func(item UserApprove, _ int) int {
		return item.ApproveId
	}), nil
}

// 获取审批流 通过用户ID 以及 审批流ID
type Workflows struct {
	AdId            int     `json:"adId" comment:"详情ID"`
	Priority        int     `json:"priority" comment:"审批者审批时所在的位置"`
	UserID          int     `json:"userId" comment:"用户ID"`
	TotalPriceLimit float64 `json:"totalPriceLimit" comment:"采购限额"`
	LimitType       int  `json:"limitType" comment:"采购限额类型(1:订单金额,2:商品金额)"`
	ApproveRankType int  `json:"approveRankType" comment:"审批层级类型:1:普签;2:会签;3:或签；默认为:1"`
	ApproveID       int     `json:"approveID" comment:"审批流ID"`
	LoginName       string  `json:"loginName" comment:"用户登录名称"`
	UserName        string  `json:"userName" comment:"用户姓名"`
	UserEmail       string  `json:"userEmail" comment:"用户邮箱"`
	// Department      string `json:"department" comment:"用户部门"`
	Position string `json:"position" comment:"用户Position"`
}

func (e *UserApprove) Workflows(tx *gorm.DB, userId int, approveId []int) ([]Workflows, error) {
	ucPrefix := global.GetTenantUcDBNameWithDB(tx)
	var workflows []Workflows
	db := tx.Table(ucPrefix+"."+e.TableName()+" a").
		Select("c.id as ad_id, c.priority, c.user_id, c.total_price_limit, c.limit_type, c.approve_rank_type, a.approve_id, d.login_name, d.user_name, d.user_email, d.position").
		Joins("INNER JOIN "+ucPrefix+".email_approve b ON (a.approve_id = b.id)").
		Joins("INNER JOIN "+ucPrefix+".email_approve_detail c ON (b.id = c.approve_id)").
		Joins("INNER JOIN "+ucPrefix+".user_info d ON (c.user_id = d.id)").
		Where("a.user_id = ?", userId).
		Where("a.status = ?", 1).
		Where("a.deleted_at is null and c.deleted_at is null").
		Where("b.approve_type = ? AND b.approve_status = ?", 9, 1).
		Where("c.approve_detail_status = ?", 1)
	if len(approveId) > 0 {
		db.Where("b.id in ?", approveId)
	}
	err := db.Order("c.id, a.approve_id ASC, c.priority ASC").Find(&workflows).Error
	if err != nil {
		return nil, err
	}

	return workflows, nil
}
