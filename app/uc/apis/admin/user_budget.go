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

type UserBudget struct {
	api.Api
}

// GetPage 获取用户预算列表
// @Summary 获取用户预算列表
// @Description 获取用户预算列表
// @Tags 用户预算
// @Param userId query int false "用户id"
// @Param userName query string false "用户姓名"
// @Param month query string false "年月"
// @Param depId query int false "部门id"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]dto.UserBudgetGetPageResp}} "{"code": 200, "data": [...]}"
// @Router /api/v1/uc/user-budget [get]
// @Security Bearer
func (e UserBudget) GetPage(c *gin.Context) {
    req := dto.UserBudgetGetPageReq{}
    s := service.UserBudget{}
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
	list := make([]dto.UserBudgetGetPageResp, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取用户预算失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取用户预算
// @Summary 获取用户预算
// @Description 获取用户预算
// @Tags 用户预算
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.UserBudget} "{"code": 200, "data": [...]}"
// @Router /api/v1/uc/user-budget/{id} [get]
// @Security Bearer
func (e UserBudget) Get(c *gin.Context) {
	req := dto.UserBudgetGetReq{}
	s := service.UserBudget{}
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
	var object models.UserBudget

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取用户预算失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK( object, "查询成功")
}

//// Insert 创建用户预算
//// @Summary 创建用户预算
//// @Description 创建用户预算
//// @Tags 用户预算
//// @Accept application/json
//// @Product application/json
//// @Param data body dto.UserBudgetInsertReq true "data"
//// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
//// @Router /api/v1/uc/user-budget [post]
//// @Security Bearer
//func (e UserBudget) Insert(c *gin.Context) {
//    req := dto.UserBudgetInsertReq{}
//    s := service.UserBudget{}
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
//    req.SetCreateByName(user.GetUserName(c))
//	err = s.Insert(&req)
//	if err != nil {
//		e.Error(500, err, fmt.Sprintf("创建用户预算失败，\r\n失败信息 %s", err.Error()))
//        return
//	}
//
//	e.OK(req.GetId(), "创建成功")
//}
//
//// Update 修改用户预算
//// @Summary 修改用户预算
//// @Description 修改用户预算
//// @Tags 用户预算
//// @Accept application/json
//// @Product application/json
//// @Param id path int true "id"
//// @Param data body dto.UserBudgetUpdateReq true "body"
//// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
//// @Router /api/v1/uc/user-budget/{id} [put]
//// @Security Bearer
//func (e UserBudget) Update(c *gin.Context) {
//    req := dto.UserBudgetUpdateReq{}
//    s := service.UserBudget{}
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
//		e.Error(500, err, fmt.Sprintf("修改用户预算失败，\r\n失败信息 %s", err.Error()))
//        return
//	}
//	e.OK( req.GetId(), "修改成功")
//}
//
//// Delete 删除用户预算
//// @Summary 删除用户预算
//// @Description 删除用户预算
//// @Tags 用户预算
//// @Param data body dto.UserBudgetDeleteReq true "body"
//// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
//// @Router /api/v1/uc/user-budget [delete]
//// @Security Bearer
//func (e UserBudget) Delete(c *gin.Context) {
//    s := service.UserBudget{}
//    req := dto.UserBudgetDeleteReq{}
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
//		e.Error(500, err, fmt.Sprintf("删除用户预算失败，\r\n失败信息 %s", err.Error()))
//        return
//	}
//	e.OK( req.GetId(), "删除成功")
//}
