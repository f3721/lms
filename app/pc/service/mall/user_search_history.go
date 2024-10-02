package mall

import (
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"go-admin/app/pc/models"
	"go-admin/app/pc/service/mall/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type UserSearchHistory struct {
	service.Service
}

// GetPage 获取UserSearchHistory列表
func (e *UserSearchHistory) GetPage(c *dto.UserSearchHistoryGetPageReq, p *actions.DataPermission, list *[]dto.UserSearchHistoryGetPageResp, count *int64) error {
	var err error
	var data models.UserSearchHistory

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			dto.MakeUserSearchHistoryReqCondition(c),
		).
		Distinct("keyword").
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("UserSearchHistoryService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Insert 创建UserSearchHistory对象
func (e *UserSearchHistory) Insert(c *dto.UserSearchHistoryInsertReq) error {
	var err error
	var data models.UserSearchHistory
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("UserSearchHistoryService Insert error:%s \r\n", err)
		return err
	}
	return nil
}
