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

type ReminderListSku struct {
	api.Api
}

// GetPage 获取补货清单SKU明细列表
// @Summary 获取补货清单SKU明细列表
// @Description 获取补货清单SKU明细列表
// @Tags 补货清单SKU明细
// @Param reminderListId query int64 false "sxyz_reminder_list 补货清单表id"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.ReminderListSku}} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/reminder-list-sku [get]
// @Security Bearer
func (e ReminderListSku) GetPage(c *gin.Context) {
	req := dto.ReminderListSkuGetPageReq{}
	s := service.ReminderListSku{}
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
	list := make([]models.ReminderListSku, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取补货清单SKU明细失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

//
//// Get 获取补货清单SKU明细
//// @Summary 获取补货清单SKU明细
//// @Description 获取补货清单SKU明细
//// @Tags 补货清单SKU明细
//// @Param id path int false "id"
//// @Success 200 {object} response.Response{data=models.ReminderListSku} "{"code": 200, "data": [...]}"
//// @Router /api/v1/wc/admin/reminder-list-sku/{id} [get]
//// @Security Bearer
//func (e ReminderListSku) Get(c *gin.Context) {
//	req := dto.ReminderListSkuGetReq{}
//	s := service.ReminderListSku{}
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
//	var object models.ReminderListSku
//
//	p := actions.GetPermissionFromContext(c)
//	err = s.Get(&req, p, &object)
//	if err != nil {
//		e.Error(500, err, fmt.Sprintf("获取补货清单SKU明细失败，\r\n失败信息 %s", err.Error()))
//		return
//	}
//
//	e.OK(object, "查询成功")
//}
//
//// Insert 创建补货清单SKU明细
//// @Summary 创建补货清单SKU明细
//// @Description 创建补货清单SKU明细
//// @Tags 补货清单SKU明细
//// @Accept application/json
//// @Product application/json
//// @Param data body dto.ReminderListSkuInsertReq true "data"
//// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
//// @Router /api/v1/wc/admin/reminder-list-sku [post]
//// @Security Bearer
//func (e ReminderListSku) Insert(c *gin.Context) {
//	req := dto.ReminderListSkuInsertReq{}
//	s := service.ReminderListSku{}
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
//	err = s.Insert(&req)
//	if err != nil {
//		e.Error(500, err, fmt.Sprintf("创建补货清单SKU明细失败，\r\n失败信息 %s", err.Error()))
//		return
//	}
//
//	e.OK(req.GetId(), "创建成功")
//}
//
//// Update 修改补货清单SKU明细
//// @Summary 修改补货清单SKU明细
//// @Description 修改补货清单SKU明细
//// @Tags 补货清单SKU明细
//// @Accept application/json
//// @Product application/json
//// @Param id path int true "id"
//// @Param data body dto.ReminderListSkuUpdateReq true "body"
//// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
//// @Router /api/v1/wc/admin/reminder-list-sku/{id} [put]
//// @Security Bearer
//func (e ReminderListSku) Update(c *gin.Context) {
//	req := dto.ReminderListSkuUpdateReq{}
//	s := service.ReminderListSku{}
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
//	req.SetUpdateBy(user.GetUserId(c))
//	p := actions.GetPermissionFromContext(c)
//
//	err = s.Update(&req, p)
//	if err != nil {
//		e.Error(500, err, fmt.Sprintf("修改补货清单SKU明细失败，\r\n失败信息 %s", err.Error()))
//		return
//	}
//	e.OK(req.GetId(), "修改成功")
//}
//
//// Delete 删除补货清单SKU明细
//// @Summary 删除补货清单SKU明细
//// @Description 删除补货清单SKU明细
//// @Tags 补货清单SKU明细
//// @Param data body dto.ReminderListSkuDeleteReq true "body"
//// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
//// @Router /api/v1/wc/admin/reminder-list-sku [delete]
//// @Security Bearer
//func (e ReminderListSku) Delete(c *gin.Context) {
//	s := service.ReminderListSku{}
//	req := dto.ReminderListSkuDeleteReq{}
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
//		e.Error(500, err, fmt.Sprintf("删除补货清单SKU明细失败，\r\n失败信息 %s", err.Error()))
//		return
//	}
//	e.OK(req.GetId(), "删除成功")
//}
