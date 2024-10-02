package admin

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/wc/models"
	service "go-admin/app/wc/service/admin"
	"go-admin/app/wc/service/admin/dto"
	"go-admin/common/actions"
)

type ReminderRuleSku struct {
	api.Api
}

// GetPage 获取补货提醒规则sku单独配置子表列表
// @Summary 获取补货提醒规则sku单独配置子表列表
// @Description 获取补货提醒规则sku单独配置子表列表
// @Tags 补货提醒规则sku单独配置子表
// @Param reminderRuleId query int64 false "sxyz_reminder_rule 补货提醒规则表id"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.ReminderRuleSku}} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/reminder-rule-sku [get]
// @Security Bearer
func (e ReminderRuleSku) GetPage(c *gin.Context) {
	req := dto.ReminderRuleSkuGetPageReq{}
	s := service.ReminderRuleSku{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	p := actions.GetPermissionFromContext(c)
	list := make([]models.ReminderRuleSku, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取补货提醒规则sku单独配置子表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取补货提醒规则sku单独配置子表
// @Summary 获取补货提醒规则sku单独配置子表
// @Description 获取补货提醒规则sku单独配置子表
// @Tags 补货提醒规则sku单独配置子表
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.ReminderRuleSku} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/reminder-rule-sku/{id} [get]
// @Security Bearer
func (e ReminderRuleSku) Get(c *gin.Context) {
	req := dto.ReminderRuleSkuGetReq{}
	s := service.ReminderRuleSku{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	var object models.ReminderRuleSku

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取补货提醒规则sku单独配置子表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建补货提醒规则sku单独配置子表
// @Summary 创建补货提醒规则sku单独配置子表
// @Description 创建补货提醒规则sku单独配置子表
// @Tags 补货提醒规则sku单独配置子表
// @Accept application/json
// @Product application/json
// @Param data body dto.ReminderRuleSkuInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/wc/admin/reminder-rule-sku [post]
// @Security Bearer
func (e ReminderRuleSku) Insert(c *gin.Context) {
	req := dto.ReminderRuleSkuInsertReq{}
	s := service.ReminderRuleSku{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	// 设置创建人
	req.SetCreateBy(user.GetUserId(c))

	err = s.Insert(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建补货提醒规则sku单独配置子表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改补货提醒规则sku单独配置子表
// @Summary 修改补货提醒规则sku单独配置子表
// @Description 修改补货提醒规则sku单独配置子表
// @Tags 补货提醒规则sku单独配置子表
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.ReminderRuleSkuUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/wc/admin/reminder-rule-sku/{id} [put]
// @Security Bearer
func (e ReminderRuleSku) Update(c *gin.Context) {
	req := dto.ReminderRuleSkuUpdateReq{}
	s := service.ReminderRuleSku{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	req.SetUpdateBy(user.GetUserId(c))
	p := actions.GetPermissionFromContext(c)

	err = s.Update(nil, &req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改补货提醒规则sku单独配置子表失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除补货提醒规则sku单独配置子表
// @Summary 删除补货提醒规则sku单独配置子表
// @Description 删除补货提醒规则sku单独配置子表
// @Tags 补货提醒规则sku单独配置子表
// @Param data body dto.ReminderRuleSkuDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/wc/admin/reminder-rule-sku [delete]
// @Security Bearer
func (e ReminderRuleSku) Delete(c *gin.Context) {
	s := service.ReminderRuleSku{}
	req := dto.ReminderRuleSkuDeleteReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	// req.SetUpdateBy(user.GetUserId(c))
	p := actions.GetPermissionFromContext(c)

	err = s.Remove(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("删除补货提醒规则sku单独配置子表失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}
