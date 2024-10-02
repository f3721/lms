package dto

import (
	"errors"
	"github.com/samber/lo"
	"go-admin/app/uc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
	"strconv"
)

type EmailApproveGetPageReq struct {
	dto.Pagination `search:"-"`
	CompanyId      int
	FilterProcess  string `form:"filterProcess"  search:"-" comment:"审批流程过滤"`
	FilterPerson   string `form:"filterPerson"  search:"-" comment:" 领用人过滤"`
}

func (m *EmailApproveGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type EmailApproveInsertReq struct {
	Id                  int                           `json:"-" comment:"id"` // id
	CompanyId           int                           `json:"-" comment:"CompanyId"`
	EmailApproveDetails []EmailApproveDetailInsertReq `json:"emailApproveDetails" comment:"EmailApproveDetails" vd:"@:len($)>0; msg:'审批层级不能为空'"`
	RecipientUsers      []int                         `json:"recipientUsers" comment:"RecipientUsers" vd:"@:len($)>0; msg:'领用人不能为空'"`
	common.ControlBy
}

func (s *EmailApproveInsertReq) Generate(model *models.EmailApprove) error {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.CompanyId = s.CompanyId
	model.ApproveStatus = models.EmailApproveStatus1
	model.ApproveType = models.EmailApproveType9
	model.CreateBy = s.CreateBy
	model.CreateByName = s.CreateByName
	for index, item := range s.EmailApproveDetails {
		if item.ApproveRankType == models.EmailApproveRankType1 {
			if len(item.ApproveUsers) > 1 {
				return errors.New("审批流程" + strconv.Itoa(index+0) + ",普签审批人数量最多为一人")
			}
		} else {
			if len(item.ApproveUsers) != len(lo.Uniq(item.ApproveUsers)) {
				return errors.New("审批流程" + strconv.Itoa(index+0) + ",审批人重复")
			}
		}
		item.SetCreateBy(s.CreateBy)
		item.SetCreateByName(s.CreateByName)
		emailApproveDetail := models.EmailApproveDetail{}
		item.Generate(&emailApproveDetail, index)
		for _, userId := range item.ApproveUsers {
			emailApproveDetail.UserId = userId
			model.EmailApproveDetails = append(model.EmailApproveDetails, emailApproveDetail)
		}
	}
	if len(s.RecipientUsers) != len(lo.Uniq(s.RecipientUsers)) {
		return errors.New("领用人重复")
	}
	for _, personId := range s.RecipientUsers {
		userApprove := models.UserApprove{}
		userApprove.UserId = personId
		userApprove.Status = models.UserApproveStatus1
		userApprove.CreateBy = s.CreateBy
		userApprove.CreateByName = s.CreateByName
		model.UserApproves = append(model.UserApproves, userApprove)
	}
	return nil
}

func (s *EmailApproveInsertReq) GetId() interface{} {
	return s.Id
}

type EmailApproveUpdateReq struct {
	Id                  int                           `uri:"id" comment:"id"` // id
	CompanyId           int                           `json:"-" comment:"CompanyId"`
	EmailApproveDetails []EmailApproveDetailInsertReq `json:"emailApproveDetails" comment:"EmailApproveDetails" vd:"@:len($)>0; msg:'审批层级不能为空'"`
	RecipientUsers      []int                         `json:"recipientUsers" comment:"RecipientUsers" vd:"@:len($)>0; msg:'领用人不能为空'"`
	common.ControlBy
}

func (s *EmailApproveUpdateReq) Generate(model *models.EmailApprove) error {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.UpdateBy = s.UpdateBy
	model.UpdateByName = s.UpdateByName
	for index, item := range s.EmailApproveDetails {
		if item.ApproveRankType == models.EmailApproveRankType1 {
			if len(item.ApproveUsers) > 1 {
				return errors.New("审批流程" + strconv.Itoa(index+0) + ",普签审批人数量最多为一人")
			}
		} else {
			if len(item.ApproveUsers) != len(lo.Uniq(item.ApproveUsers)) {
				return errors.New("审批流程" + strconv.Itoa(index+0) + ",审批人重复")
			}
		}
		item.SetCreateBy(s.UpdateBy)
		item.SetCreateByName(s.UpdateByName)
		emailApproveDetail := models.EmailApproveDetail{}
		item.Generate(&emailApproveDetail, index)
		for _, userId := range item.ApproveUsers {
			emailApproveDetail.UserId = userId
			model.EmailApproveDetails = append(model.EmailApproveDetails, emailApproveDetail)
		}
	}
	if len(s.RecipientUsers) != len(lo.Uniq(s.RecipientUsers)) {
		return errors.New("领用人重复")
	}
	for _, personId := range s.RecipientUsers {
		userApprove := models.UserApprove{}
		userApprove.UserId = personId
		userApprove.Status = models.UserApproveStatus1
		userApprove.CreateBy = s.UpdateBy
		userApprove.CreateByName = s.UpdateByName
		model.UserApproves = append(model.UserApproves, userApprove)
	}
	return nil
}

func (s *EmailApproveUpdateReq) GetId() interface{} {
	return s.Id
}

// EmailApproveGetReq 功能获取请求参数
type EmailApproveGetReq struct {
	Id        int `uri:"id"`
	CompanyId int `json:"-" comment:"CompanyId"`
}

func (s *EmailApproveGetReq) GetId() interface{} {
	return s.Id
}

type EmailApproveGetResp struct {
	EmailApproveDetails []EmailApproveDetailGetResp `json:"emailApproveDetails" comment:"EmailApproveDetails"`
	RecipientUsers      []EmailApproveGetUserResp   `json:"recipientUsers"`
	AllApproveUsers     []EmailApproveGetUserResp   `json:"allApproveUsers"`
	AllRecipientUsers   []EmailApproveGetUserResp   `json:"allRecipientUsers"`
}

type EmailApproveDetailGetResp struct {
	ApproveId           int                       `json:"approveId"`
	Priority            int                       `json:"priority"`
	TotalPriceLimit     string                    `json:"totalPriceLimit"`
	LimitType           int                       `json:"limitType"`
	ApproveDetailStatus int                       `json:"approveDetailStatus"`
	ApproveRankType     int                       `json:"approveRankType"`
	ApproveUsers        []EmailApproveGetUserResp `json:"approveUsers"`
}

type EmailApproveGetUserResp struct {
	UserId   int    `json:"userId" comment:"UserId"`
	UserName string `json:"userName" comment:"UserName"`
}

type EmailApproveGetApproveAndRecipient struct {
	AllApproveUsers   []EmailApproveGetUserResp `json:"allApproveUsers"`
	AllRecipientUsers []EmailApproveGetUserResp `json:"allRecipientUsers"`
}

// EmailApproveDeleteReq 功能删除请求参数
type EmailApproveDeleteReq struct {
	Id        int `json:"id"`
	CompanyId int `json:"-" comment:"CompanyId"`
	common.ControlBy
}

func (s *EmailApproveDeleteReq) GetId() interface{} {
	return s.Id
}

type EmailApproveWorkflow struct {
	Id      int    `json:"id"`
	Process string `json:"process"`
}
