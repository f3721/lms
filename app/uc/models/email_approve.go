package models

import (
	"go-admin/common/models"
)

const (
	EmailApproveStatus0 = 0
	EmailApproveStatus1 = 1

	EmailApproveType9 = 9

	EmailApproveLogModelName   = "emailApprove"
	EmailApproveLogModelInsert = "insert"
	EmailApproveLogModelUpdate = "update"
	EmailApproveLogModelDelete = "delete"
)

type EmailApprove struct {
	models.Model

	ApproveStatus       int                  `json:"approveStatus" gorm:"type:tinyint(1);comment:审批流状态(1可用0不可用)"`
	CompanyId           int                  `json:"companyId" gorm:"type:int unsigned;comment:公司ID"`
	ApproveType         int                  `json:"approveType" gorm:"type:tinyint(1);comment:审批流类型(1系统审批,9用户审批流)"`
	EmailApproveDetails []EmailApproveDetail `json:"emailApproveDetails" gorm:"foreignKey:ApproveId;"`
	UserApproves        []UserApprove        `json:"userApproves" gorm:"foreignKey:ApproveId;"`
	models.ModelTime
	models.ControlBy
}

func (EmailApprove) TableName() string {
	return "email_approve"
}

func (e *EmailApprove) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *EmailApprove) GetId() interface{} {
	return e.Id
}

/*type EmailApproveProcessCustomer struct {
	models.Model
	ApproveStatus        int    `json:"approveStatus"`
	ApproveRankType      int    `json:"approveRankType"`
	ProcessUserId        string `json:"processUserId"`
	ProcessUserName      string `json:"processUserName"`
	ProcessUserLoginName string `json:"processUserLoginName"`
	ProcessUserPhone     string `json:"processUserPhone"`
	ProcessUserEmail     string `json:"processUserEmail"`
}

func (e *EmailApproveProcessCustomer) GetProcessListByCompanyId(tx *gorm.DB, companyId int) (*[]EmailApproveProcessCustomer, error) {
	var emailApproveCustomer = &[]EmailApproveProcessCustomer{}
	if err := tx.Model(e).
		Joins("LEFT JOIN email_approve_detail detail ON detail.approve_id = email_approve.id").
		Joins("LEFT JOIN user_info ON detail.user_id = user_info.id").
		Where("email_approve.company_id = ?", companyId).
		Where("email_approve_detail.approve_detail_status = ?", EmailApproveDetailStatus1).Group("detail.priority,email_approve.id").
		Select(" email_approve.id,detail.priority,detail.approve_rank_type,email_approve.approve_status," +
			"GROUP_CONCAT(user_info.id SEPARATOR ',' ) process_user_id" +
			",GROUP_CONCAT(user_info.user_name SEPARATOR ',' ) process_user_name" +
			",GROUP_CONCAT(user_info.login_name SEPARATOR ',' ) process_login_name" +
			",GROUP_CONCAT(user_info.user_phone SEPARATOR ',' ) process_user_phone" +
			",GROUP_CONCAT(user_info.user_email SEPARATOR ',' ) process_user_email").
		Find(emailApproveCustomer).Error; err != nil {
		return nil, err
	}
	sdk.Runtime.GetCacheAdapter()
	for index, item := range *emailApproveCustomer {
		if item.ApproveRankType == EmailApproveApproveRankType2 {
			(*emailApproveCustomer)[index].ProcessUserId = strings.ReplaceAll(item.ProcessUserId, ",", "、")
		}
		if item.ApproveRankType == EmailApproveApproveRankType3 {
			(*emailApproveCustomer)[index].ProcessUserId = strings.ReplaceAll(item.ProcessUserId, ",", "/")
		}
	}
	return emailApproveCustomer, nil
}*/

type EmailApproveListCustomer struct {
	models.Model
	ApproveStatus int    `json:"approveStatus"`
	Priority      int    `json:"priority"`
	Process       string `json:"process"`
	Person        string `json:"person"`
}

func (EmailApproveListCustomer) TableName() string {
	return "email_approve"
}
