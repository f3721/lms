package admin

import (
	"fmt"
	service "go-admin/app/oc/service/admin"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/oc/models"
	"go-admin/app/oc/service/admin/dto"
	"go-admin/common/actions"
)

type CsApplyLog struct {
	api.Api
}

// GetPage 获取售后日志列表
// @Summary 获取售后日志列表
// @Description 获取售后日志列表
// @Tags 售后中心
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.CsApplyLog}} "{"code": 200, "data": [...]}"
// @Router /api/v1/oc/admin/cs-apply-log [get]
// @Security Bearer
func (e CsApplyLog) GetPage(c *gin.Context) {
	req := dto.CsApplyLogGetPageReq{}
	s := service.CsApplyLog{}
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
	list := make([]models.CsApplyLog, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取CsApplyLog失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

//// Get 获取CsApplyLog
//// @Summary 获取CsApplyLog
//// @Description 获取CsApplyLog
//// @Tags CsApplyLog
//// @Param id path int false "id"
//// @Success 200 {object} response.Response{data=models.CsApplyLog} "{"code": 200, "data": [...]}"
//// @Router /api/v1/oc/admin/cs-apply-log/{id} [get]
//// @Security Bearer
//func (e CsApplyLog) Get(c *gin.Context) {
//	req := dto.CsApplyLogGetReq{}
//	s := service.CsApplyLog{}
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
//	var object models.CsApplyLog
//
//	p := actions.GetPermissionFromContext(c)
//	err = s.Get(&req, p, &object)
//	if err != nil {
//		e.Error(500, err, fmt.Sprintf("获取CsApplyLog失败，\r\n失败信息 %s", err.Error()))
//        return
//	}
//
//	e.OK( object, "查询成功")
//}
//
//// Insert 创建CsApplyLog
//// @Summary 创建CsApplyLog
//// @Description 创建CsApplyLog
//// @Tags CsApplyLog
//// @Accept application/json
//// @Product application/json
//// @Param data body dto.CsApplyLogInsertReq true "data"
//// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
//// @Router /api/v1/oc/admin/cs-apply-log [post]
//// @Security Bearer
//func (e CsApplyLog) Insert(c *gin.Context) {
//    req := dto.CsApplyLogInsertReq{}
//    s := service.CsApplyLog{}
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
//		e.Error(500, err, fmt.Sprintf("创建CsApplyLog失败，\r\n失败信息 %s", err.Error()))
//        return
//	}
//
//	e.OK(req.GetId(), "创建成功")
//}
//
//// Update 修改CsApplyLog
//// @Summary 修改CsApplyLog
//// @Description 修改CsApplyLog
//// @Tags CsApplyLog
//// @Accept application/json
//// @Product application/json
//// @Param id path int true "id"
//// @Param data body dto.CsApplyLogUpdateReq true "body"
//// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
//// @Router /api/v1/oc/admin/cs-apply-log/{id} [put]
//// @Security Bearer
//func (e CsApplyLog) Update(c *gin.Context) {
//    req := dto.CsApplyLogUpdateReq{}
//    s := service.CsApplyLog{}
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
//		e.Error(500, err, fmt.Sprintf("修改CsApplyLog失败，\r\n失败信息 %s", err.Error()))
//        return
//	}
//	e.OK( req.GetId(), "修改成功")
//}
//
//// Delete 删除CsApplyLog
//// @Summary 删除CsApplyLog
//// @Description 删除CsApplyLog
//// @Tags CsApplyLog
//// @Param data body dto.CsApplyLogDeleteReq true "body"
//// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
//// @Router /api/v1/oc/admin/cs-apply-log [delete]
//// @Security Bearer
//func (e CsApplyLog) Delete(c *gin.Context) {
//    s := service.CsApplyLog{}
//    req := dto.CsApplyLogDeleteReq{}
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
//		e.Error(500, err, fmt.Sprintf("删除CsApplyLog失败，\r\n失败信息 %s", err.Error()))
//        return
//	}
//	e.OK( req.GetId(), "删除成功")
//}
