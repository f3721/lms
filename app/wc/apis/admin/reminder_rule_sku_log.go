package admin

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/wc/models"
	service "go-admin/app/wc/service/admin"
	"go-admin/app/wc/service/admin/dto"
	"go-admin/common/actions"
)

type ReminderRuleSkuLog struct {
	api.Api
}

// GetPage 获取补货提醒规则sku单独配置子表的log表列表
// @Summary 获取补货提醒规则sku单独配置子表的log表列表
// @Description 获取补货提醒规则sku单独配置子表的log表列表
// @Tags 补货提醒规则sku单独配置子表的log表
// @Param reminderRuleSkuId query int64 false "sku补货提醒规则表id"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.ReminderRuleSkuLog}} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/reminder-rule-sku-log [get]
// @Security Bearer
func (e ReminderRuleSkuLog) GetPage(c *gin.Context) {
	req := dto.ReminderRuleSkuLogGetPageReq{}
	s := service.ReminderRuleSkuLog{}
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
	list := make([]models.ReminderRuleSkuLog, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取补货提醒规则sku单独配置子表的log表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

//
//// Get 获取补货提醒规则sku单独配置子表的log表
//// @Summary 获取补货提醒规则sku单独配置子表的log表
//// @Description 获取补货提醒规则sku单独配置子表的log表
//// @Tags 补货提醒规则sku单独配置子表的log表
//// @Param id path int false "id"
//// @Success 200 {object} response.Response{data=models.ReminderRuleSkuLog} "{"code": 200, "data": [...]}"
//// @Router /api/v1/wc/admin/reminder-rule-sku-log/{id} [get]
//// @Security Bearer
//func (e ReminderRuleSkuLog) Get(c *gin.Context) {
//	req := dto.ReminderRuleSkuLogGetReq{}
//	s := service.ReminderRuleSkuLog{}
//    err := e.MakeContext(c).
//		MakeOrm().
//		Bind(&req).
//		MakeService(&s.Service).
//		Errors
//	if err != nil {
//		e.Logger.Error(err)
//		e.Error(500, err, err.Error())
//		return
//	}
//	var object models.ReminderRuleSkuLog
//
//	p := actions.GetPermissionFromContext(c)
//	err = s.Get(&req, p, &object)
//	if err != nil {
//		e.Error(500, err, fmt.Sprintf("获取补货提醒规则sku单独配置子表的log表失败，\r\n失败信息 %s", err.Error()))
//        return
//	}
//
//	e.OK( object, "查询成功")
//}
//
//// Insert 创建补货提醒规则sku单独配置子表的log表
//// @Summary 创建补货提醒规则sku单独配置子表的log表
//// @Description 创建补货提醒规则sku单独配置子表的log表
//// @Tags 补货提醒规则sku单独配置子表的log表
//// @Accept application/json
//// @Product application/json
//// @Param data body dto.ReminderRuleSkuLogInsertReq true "data"
//// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
//// @Router /api/v1/wc/admin/reminder-rule-sku-log [post]
//// @Security Bearer
//func (e ReminderRuleSkuLog) Insert(c *gin.Context) {
//    req := dto.ReminderRuleSkuLogInsertReq{}
//    s := service.ReminderRuleSkuLog{}
//    err := e.MakeContext(c).
//        MakeOrm().
//        Bind(&req).
//        MakeService(&s.Service).
//        Errors
//    if err != nil {
//        e.Logger.Error(err)
//        e.Error(500, err, err.Error())
//        return
//    }
//	// 设置创建人
//	req.SetCreateBy(user.GetUserId(c))
//
//	err = s.Insert(&req)
//	if err != nil {
//		e.Error(500, err, fmt.Sprintf("创建补货提醒规则sku单独配置子表的log表失败，\r\n失败信息 %s", err.Error()))
//        return
//	}
//
//	e.OK(req.GetId(), "创建成功")
//}
//
//// Update 修改补货提醒规则sku单独配置子表的log表
//// @Summary 修改补货提醒规则sku单独配置子表的log表
//// @Description 修改补货提醒规则sku单独配置子表的log表
//// @Tags 补货提醒规则sku单独配置子表的log表
//// @Accept application/json
//// @Product application/json
//// @Param id path int true "id"
//// @Param data body dto.ReminderRuleSkuLogUpdateReq true "body"
//// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
//// @Router /api/v1/wc/admin/reminder-rule-sku-log/{id} [put]
//// @Security Bearer
//func (e ReminderRuleSkuLog) Update(c *gin.Context) {
//    req := dto.ReminderRuleSkuLogUpdateReq{}
//    s := service.ReminderRuleSkuLog{}
//    err := e.MakeContext(c).
//        MakeOrm().
//        Bind(&req).
//        MakeService(&s.Service).
//        Errors
//    if err != nil {
//        e.Logger.Error(err)
//        e.Error(500, err, err.Error())
//        return
//    }
//	req.SetUpdateBy(user.GetUserId(c))
//	p := actions.GetPermissionFromContext(c)
//
//	err = s.Update(&req, p)
//	if err != nil {
//		e.Error(500, err, fmt.Sprintf("修改补货提醒规则sku单独配置子表的log表失败，\r\n失败信息 %s", err.Error()))
//        return
//	}
//	e.OK( req.GetId(), "修改成功")
//}
//
//// Delete 删除补货提醒规则sku单独配置子表的log表
//// @Summary 删除补货提醒规则sku单独配置子表的log表
//// @Description 删除补货提醒规则sku单独配置子表的log表
//// @Tags 补货提醒规则sku单独配置子表的log表
//// @Param data body dto.ReminderRuleSkuLogDeleteReq true "body"
//// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
//// @Router /api/v1/wc/admin/reminder-rule-sku-log [delete]
//// @Security Bearer
//func (e ReminderRuleSkuLog) Delete(c *gin.Context) {
//    s := service.ReminderRuleSkuLog{}
//    req := dto.ReminderRuleSkuLogDeleteReq{}
//    err := e.MakeContext(c).
//        MakeOrm().
//        Bind(&req).
//        MakeService(&s.Service).
//        Errors
//    if err != nil {
//        e.Logger.Error(err)
//        e.Error(500, err, err.Error())
//        return
//    }
//
//	// req.SetUpdateBy(user.GetUserId(c))
//	p := actions.GetPermissionFromContext(c)
//
//	err = s.Remove(&req, p)
//	if err != nil {
//		e.Error(500, err, fmt.Sprintf("删除补货提醒规则sku单独配置子表的log表失败，\r\n失败信息 %s", err.Error()))
//        return
//	}
//	e.OK( req.GetId(), "删除成功")
//}
