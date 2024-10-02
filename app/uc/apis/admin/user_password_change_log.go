package admin

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/uc/models"
	service "go-admin/app/uc/service/admin"
	"go-admin/app/uc/service/admin/dto"
	"go-admin/common/actions"
)

type UserPasswordChangeLog struct {
	api.Api
}

// GetPage 获取用户密码修改记录列表
// @Summary 获取用户密码修改记录列表
// @Description 获取用户密码修改记录列表
// @Tags 用户密码修改记录
// @Param userId query int false "用户ID"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.UserPasswordChangeLog}} "{"code": 200, "data": [...]}"
// @Router /api/v1/uc/admin/user-password-change-log [get]
// @Security Bearer
func (e UserPasswordChangeLog) GetPage(c *gin.Context) {
	req := dto.UserPasswordChangeLogGetPageReq{}
	s := service.UserPasswordChangeLog{}
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
	list := make([]models.UserPasswordChangeLog, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取用户密码修改记录失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

//
//// Get 获取用户密码修改记录
//// @Summary 获取用户密码修改记录
//// @Description 获取用户密码修改记录
//// @Tags 用户密码修改记录
//// @Param id path int false "id"
//// @Success 200 {object} response.Response{data=models.UserPasswordChangeLog} "{"code": 200, "data": [...]}"
//// @Router /api/v1/uc/admin/user-password-change-log/{id} [get]
//// @Security Bearer
//func (e UserPasswordChangeLog) Get(c *gin.Context) {
//	req := dto.UserPasswordChangeLogGetReq{}
//	s := service.UserPasswordChangeLog{}
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
//	var object models.UserPasswordChangeLog
//
//	p := actions.GetPermissionFromContext(c)
//	err = s.Get(&req, p, &object)
//	if err != nil {
//		e.Error(500, err, fmt.Sprintf("获取用户密码修改记录失败，\r\n失败信息 %s", err.Error()))
//		return
//	}
//
//	e.OK(object, "查询成功")
//}
//
//// Insert 创建用户密码修改记录
//// @Summary 创建用户密码修改记录
//// @Description 创建用户密码修改记录
//// @Tags 用户密码修改记录
//// @Accept application/json
//// @Product application/json
//// @Param data body dto.UserPasswordChangeLogInsertReq true "data"
//// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
//// @Router /api/v1/uc/admin/user-password-change-log [post]
//// @Security Bearer
//func (e UserPasswordChangeLog) Insert(c *gin.Context) {
//	req := dto.UserPasswordChangeLogInsertReq{}
//	s := service.UserPasswordChangeLog{}
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
//		e.Error(500, err, fmt.Sprintf("创建用户密码修改记录失败，\r\n失败信息 %s", err.Error()))
//		return
//	}
//
//	e.OK(req.GetId(), "创建成功")
//}
//
//// Update 修改用户密码修改记录
//// @Summary 修改用户密码修改记录
//// @Description 修改用户密码修改记录
//// @Tags 用户密码修改记录
//// @Accept application/json
//// @Product application/json
//// @Param id path int true "id"
//// @Param data body dto.UserPasswordChangeLogUpdateReq true "body"
//// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
//// @Router /api/v1/uc/admin/user-password-change-log/{id} [put]
//// @Security Bearer
//func (e UserPasswordChangeLog) Update(c *gin.Context) {
//	req := dto.UserPasswordChangeLogUpdateReq{}
//	s := service.UserPasswordChangeLog{}
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
//	p := actions.GetPermissionFromContext(c)
//
//	err = s.Update(&req, p)
//	if err != nil {
//		e.Error(500, err, fmt.Sprintf("修改用户密码修改记录失败，\r\n失败信息 %s", err.Error()))
//		return
//	}
//	e.OK(req.GetId(), "修改成功")
//}
//
//// Delete 删除用户密码修改记录
//// @Summary 删除用户密码修改记录
//// @Description 删除用户密码修改记录
//// @Tags 用户密码修改记录
//// @Param data body dto.UserPasswordChangeLogDeleteReq true "body"
//// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
//// @Router /api/v1/uc/admin/user-password-change-log [delete]
//// @Security Bearer
//func (e UserPasswordChangeLog) Delete(c *gin.Context) {
//	s := service.UserPasswordChangeLog{}
//	req := dto.UserPasswordChangeLogDeleteReq{}
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
//		e.Error(500, err, fmt.Sprintf("删除用户密码修改记录失败，\r\n失败信息 %s", err.Error()))
//		return
//	}
//	e.OK(req.GetId(), "删除成功")
//}
