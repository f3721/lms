package admin

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"github.com/jinzhu/copier"
	"github.com/prometheus/common/log"

	service "go-admin/app/wc/service/admin"
	"go-admin/app/wc/service/admin/dto"
	"go-admin/common/actions"
)

type ReminderRule struct {
	api.Api
}

// GetPage 获取补货提醒规则列表
// @Summary 获取补货提醒规则列表
// @Description 获取补货提醒规则列表
// @Tags 补货提醒规则
// @Param companyId query int64 false "公司id"
// @Param warehouseCode query string false "仓库code"
// @Param status query string false "状态 1启用 0未启用"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.ReminderRule}} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/reminder-rule [get]
// @Security Bearer
func (e ReminderRule) GetPage(c *gin.Context) {
	req := dto.ReminderRuleGetPageReq{}
	s := service.ReminderRule{}
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
	list := make([]*dto.ReminderRuleData, 0)
	var count int64
	log.Info(req.Status)
	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取补货提醒规则失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取补货提醒规则
// @Summary 获取补货提醒规则
// @Description 获取补货提醒规则
// @Tags 补货提醒规则
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=dto.ReminderRuleGetRes} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/reminder-rule/{id} [get]
// @Security Bearer
func (e ReminderRule) Get(c *gin.Context) {
	req := dto.ReminderRuleGetReq{}
	s := service.ReminderRule{}
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
	var object dto.ReminderRuleData

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取补货提醒规则失败，\r\n失败信息 %s", err.Error()))
		return
	}

	res := dto.ReminderRuleGetRes{}
	copier.Copy(&res, object)

	e.OK(res, "查询成功")
}

// Insert 创建补货提醒规则
// @Summary 创建补货提醒规则
// @Description 创建补货提醒规则
// @Tags 补货提醒规则
// @Accept application/json
// @Product application/json
// @Param data body dto.ReminderRuleInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/wc/admin/reminder-rule [post]
// @Security Bearer
func (e ReminderRule) Insert(c *gin.Context) {
	req := dto.ReminderRuleInsertReq{}
	s := service.ReminderRule{}

	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	for _, insertReq := range req.SkuList {
		err = e.Bind(&insertReq, binding.JSON).Errors
		if err != nil {
			e.Logger.Error(err)
			e.Error(500, err, err.Error())
			return
		}
	}
	log.Info("ser.GetUserId(c)")
	log.Info(user.GetUserId(c))
	// 设置创建人k
	req.SetCreateBy(user.GetUserId(c))
	req.SetCreateByName(user.GetUserName(c))
	log.Info(req)
	p := actions.GetPermissionFromContext(c)

	_, err = s.Insert(&req, c, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建补货提醒规则失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改补货提醒规则
// @Summary 修改补货提醒规则
// @Description 修改补货提醒规则
// @Tags 补货提醒规则
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.ReminderRuleUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/wc/admin/reminder-rule/{id} [put]
// @Security Bearer
func (e ReminderRule) Update(c *gin.Context) {
	req := dto.ReminderRuleUpdateReq{}
	s := service.ReminderRule{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	for _, updateReq := range req.SkuList {
		err = e.Bind(&updateReq, binding.JSON).Errors
		if err != nil {
			e.Logger.Error(err)
			e.Error(500, err, err.Error())
			return
		}
	}
	req.SetUpdateBy(user.GetUserId(c))
	req.SetUpdateByName(user.GetUserName(c))
	p := actions.GetPermissionFromContext(c)

	err = s.Update(&req, p, c)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改补货提醒规则失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "修改成功")
}

//// Delete 删除补货提醒规则
//// @Summary 删除补货提醒规则
//// @Description 删除补货提醒规则
//// @Tags 补货提醒规则
//// @Param data body dto.ReminderRuleDeleteReq true "body"
//// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
//// @Router /api/v1/wc/admin/reminder-rule [delete]
//// @Security Bearer
//func (e ReminderRule) Delete(c *gin.Context) {
//	s := service.ReminderRule{}
//	req := dto.ReminderRuleDeleteReq{}
//	err := e.MakeContext(c).
//		MakeOrm().
//		Bind(&req).
//		MakeService(&s.Service).
//		Errors
//	if err != nil {
//		e.Logger.Error(err)
//		e.Error(500, err, err.Error())
//		return
//	}
//
//	// req.SetUpdateBy(user.GetUserId(c))
//	p := actions.GetPermissionFromContext(c)
//
//	err = s.Remove(&req, p)
//	if err != nil {
//		e.Error(500, err, fmt.Sprintf("删除补货提醒规则失败，\r\n失败信息 %s", err.Error()))
//		return
//	}
//	e.OK(req.GetId(), "删除成功")
//}
