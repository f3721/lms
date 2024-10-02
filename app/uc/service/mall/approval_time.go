package mall

import (
	"errors"
	"go-admin/app/uc/models"
	"go-admin/app/uc/service/mall/dto"
	"go-admin/common/actions"
	"go-admin/common/utils"
	"sort"
	"strings"

	"github.com/go-admin-team/go-admin-core/sdk/service"
)

type ApprovalTime struct {
	service.Service
}

// 格式化时间
func (s *ApprovalTime) fmtTimes(d *dto.ApprovalTimeInsertReq) string {
	// weeks 先排序
	sort.Strings(d.Weeks)
	group := []string{d.Hour, d.Min, "*", "*", strings.Join(d.Weeks, ",")}
	return strings.Join(group, " ")
}

// 新建审批时间
func (s *ApprovalTime) Insert(d *dto.ApprovalTimeInsertReq, userId int) error {
	// 校验星期范围
	for _, val := range d.Weeks {
		if val < "0" || val > "6" {
			return errors.New("星期范围[0-6]")
		}
	}

	// 查询用户信息
	user := models.UserInfo{}
	err := s.Orm.Where("id = ?", userId).First(&user).Error
	if err != nil {
		return err
	}

	// 组合新时间
	var times []string
	newTime := s.fmtTimes(d)
	if user.EmailApproveCronExpr != "" { // 兼容首次创建
		times = strings.Split(user.EmailApproveCronExpr, "@")
		if len(times) >= 3 {
			return errors.New("最多可设置3条提醒")
		}

		// 校验是否有重复项
		if utils.InArrayString(newTime, times) {
			return errors.New("已存在相同的时间设置！")
		}
	}

	times = append(times, newTime)
	expr := strings.Join(times, "@")

	// 存储数据
	err = s.Orm.Model(&user).Where("id = ?", userId).UpdateColumn("email_approve_cron_expr", expr).Error
	return err
}

// 更新审批时间
func (s *ApprovalTime) Update(d *dto.ApprovalTimeUpdateReq, userId int, p *actions.DataPermission) error {
	// 校验星期范围
	for _, val := range d.Weeks {
		if val < "0" || val > "6" {
			return errors.New("星期范围[0-6]")
		}
	}

	// 查询用户信息
	user := models.UserInfo{}
	err := s.Orm.Where("id = ?", userId).First(&user).Error
	if err != nil {
		return err
	}

	// 查询旧数据
	if user.EmailApproveCronExpr == "" { // 兼容首次创建
		return errors.New("该用户还未创建审批时间，请先创建")
	}
	times := strings.Split(user.EmailApproveCronExpr, "@")

	// 校验重复项
	fileds := &dto.ApprovalTimeInsertReq{Hour: d.Hour, Min: d.Min, Weeks: d.Weeks}
	newTime := s.fmtTimes(fileds)
	if utils.InArrayString(newTime, times) {
		return errors.New("已存在相同的时间设置！")
	}

	// 更新时间
	if d.Id >= len(times) {
		return errors.New("传参范围错误")
	}
	times[d.Id] = newTime
	expr := strings.Join(times, "@")

	// 存储数据
	err = s.Orm.Model(&user).Where("id = ?", userId).UpdateColumn("email_approve_cron_expr", expr).Error
	return err
}

// 查询列表
func (s ApprovalTime) GetPage(userId int, p *actions.DataPermission) ([]dto.ApprovalTimeInsertReq, error) {
	res := []dto.ApprovalTimeInsertReq{}

	// 查询数据
	find := &models.UserInfo{}
	err := s.Orm.Select("email_approve_cron_expr").
		Scopes(actions.Permission(find.TableName(), p)).
		Where("id = ?", userId).First(find).Error
	if err != nil {
		return res, err
	}
	if find.EmailApproveCronExpr == "" {
		return res, err
	}

	// 解析数据
	data := strings.Split(find.EmailApproveCronExpr, "@")
	for _, val := range data {
		itemArr := strings.Split(val, " ")
		itemData := dto.ApprovalTimeInsertReq{
			Hour:   itemArr[0],
			Min:    itemArr[1],
			Weeks:  strings.Split(itemArr[4], ","),
			Repeat: "1",
		}
		res = append(res, itemData)
	}

	return res, nil
}

// 查询数据
func (s ApprovalTime) Get(d *dto.ApprovalTimeGetReq, userId int, p *actions.DataPermission) (any, error) {

	// 查询所有
	all, err := s.GetPage(userId, p)
	if err != nil {
		return all, err
	}

	// 参数校验
	if d.Id >= len(all) {
		return "", errors.New("参数范围错误")
	}

	// 返回指定
	return all[d.Id], nil
}

// 删除数据
func (s ApprovalTime) Remove(d *dto.ApprovalTimeDeleteReq, userId int, p *actions.DataPermission) error {
	// 查出所有
	find := &models.UserInfo{}
	err := s.Orm.Select("email_approve_cron_expr").
		Scopes(actions.Permission(find.TableName(), p)).
		Where("id = ?", userId).First(find).Error
	if err != nil {
		return err
	}
	if find.EmailApproveCronExpr == "" {
		return errors.New("当前用户没有任何设置")
	}

	// 删除指定
	arr := strings.Split(find.EmailApproveCronExpr, "@")

	if d.Id >= len(arr) {
		return errors.New("传参错误")
	}
	newArr := utils.UnsetArray(arr, d.Id)
	exprStr := strings.Join(newArr, "@")

	// 更新
	err = s.Orm.Model(&find).Where("id = ?", userId).UpdateColumn("email_approve_cron_expr", exprStr).Error
	return err
}
