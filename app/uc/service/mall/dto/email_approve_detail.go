package dto

import (
	"go-admin/app/uc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type EmailApproveDetailGetPageReq struct {
	dto.Pagination `search:"-"`
	EmailApproveDetailOrder
}

type EmailApproveDetailOrder struct {
	Id                  string `form:"idOrder"  search:"type:order;column:id;table:email_approve_detail"`
	ApproveId           string `form:"approveIdOrder"  search:"type:order;column:approve_id;table:email_approve_detail"`
	UserId              string `form:"userIdOrder"  search:"type:order;column:user_id;table:email_approve_detail"`
	Priority            string `form:"priorityOrder"  search:"type:order;column:priority;table:email_approve_detail"`
	TotalPriceLimit     string `form:"totalPriceLimitOrder"  search:"type:order;column:total_price_limit;table:email_approve_detail"`
	LimitType           string `form:"limitTypeOrder"  search:"type:order;column:limit_type;table:email_approve_detail"`
	ApproveDetailStatus string `form:"approveDetailStatusOrder"  search:"type:order;column:approve_detail_status;table:email_approve_detail"`
	ApproveRankType     string `form:"approveRankTypeOrder"  search:"type:order;column:approve_rank_type;table:email_approve_detail"`
	CreatedAt           string `form:"createdAtOrder"  search:"type:order;column:created_at;table:email_approve_detail"`
	CreateBy            string `form:"createByOrder"  search:"type:order;column:create_by;table:email_approve_detail"`
	UpdateBy            string `form:"updateByOrder"  search:"type:order;column:update_by;table:email_approve_detail"`
	UpdatedAt           string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:email_approve_detail"`
	DeletedAt           string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:email_approve_detail"`
	CreateByName        string `form:"createByNameOrder"  search:"type:order;column:create_by_name;table:email_approve_detail"`
	UpdateByName        string `form:"updateByNameOrder"  search:"type:order;column:update_by_name;table:email_approve_detail"`
}

func (m *EmailApproveDetailGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type EmailApproveDetailInsertReq struct {
	ApproveUsers    []int  `json:"approveUsers" comment:"用户ID" vd:"@:len($)>0; msg:'审批人不能为空'"`
	TotalPriceLimit string `json:"totalPriceLimit" comment:"采购限额"`
	LimitType       int    `json:"limitType" comment:"采购限额类型(1:订单金额,2:商品金额)" vd:"$==1 || $==2; msg:'limitType为1或2'"`
	ApproveRankType int    `json:"approveRankType" comment:"审批层级类型：1：普签；2：会签；3：或签；默认为：1；" vd:"$==1 || $==2 || $==3; msg:'approveRankType为1或2或3'"`
	common.ControlBy
}

func (s *EmailApproveDetailInsertReq) Generate(model *models.EmailApproveDetail, index int) {
	model.TotalPriceLimit = s.TotalPriceLimit
	model.LimitType = s.LimitType
	model.ApproveDetailStatus = models.EmailApproveDetailStatus1
	model.ApproveRankType = s.ApproveRankType
	model.Priority = index + 1
	model.CreateBy = s.CreateBy
	model.CreateByName = s.CreateByName
}

type EmailApproveDetailUpdateReq struct {
	Id                  int    `uri:"id" comment:"id"` // id
	ApproveId           int    `json:"approveId" comment:"审批流ID"`
	UserId              int    `json:"userId" comment:"用户ID"`
	Priority            string `json:"priority" comment:"审批者审批时所在的位置"`
	TotalPriceLimit     string `json:"totalPriceLimit" comment:"采购限额"`
	LimitType           int    `json:"limitType" comment:"采购限额类型(1:订单金额,2:商品金额)"`
	ApproveDetailStatus int    `json:"approveDetailStatus" comment:"字审批流状态"`
	ApproveRankType     int    `json:"approveRankType" comment:"审批层级类型：1：普签；2：会签；3：或签；默认为：1；"`
	CreateByName        string `json:"createByName" comment:""`
	UpdateByName        string `json:"updateByName" comment:""`
	common.ControlBy
}

func (s *EmailApproveDetailUpdateReq) Generate(model *models.EmailApproveDetail) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.ApproveId = s.ApproveId
	model.UserId = s.UserId
	model.TotalPriceLimit = s.TotalPriceLimit
	model.LimitType = s.LimitType
	model.ApproveDetailStatus = s.ApproveDetailStatus
	model.ApproveRankType = s.ApproveRankType
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
	model.CreateByName = s.CreateByName
	model.UpdateByName = s.UpdateByName
}

func (s *EmailApproveDetailUpdateReq) GetId() interface{} {
	return s.Id
}

// EmailApproveDetailGetReq 功能获取请求参数
type EmailApproveDetailGetReq struct {
	Id int `uri:"id"`
}

func (s *EmailApproveDetailGetReq) GetId() interface{} {
	return s.Id
}

// EmailApproveDetailDeleteReq 功能删除请求参数
type EmailApproveDetailDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *EmailApproveDetailDeleteReq) GetId() interface{} {
	return s.Ids
}
