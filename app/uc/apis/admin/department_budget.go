package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	service "go-admin/app/uc/service/admin"
	"go-admin/app/uc/service/admin/dto"
	"go-admin/common/actions"
)

type DepartmentBudget struct {
	api.Api
}

// GetPage 获取部门预算列表
// @Summary 获取部门预算列表
// @Description 获取部门预算列表
// @Tags 部门预算
// @Param companyId query int false "公司id"
// @Param parentDepId query int false "父部门id"
// @Param depId query int false "父部门id"
// @Param startMonth query string false "开始月份"
// @Param endMonth query string false "结束月份"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]dto.DepartmentBudgetListResp}} "{"code": 200, "data": [...]}"
// @Router /api/v1/uc/department-budget [get]
// @Security Bearer
func (e DepartmentBudget) GetPage(c *gin.Context) {
    req := dto.DepartmentBudgetGetPageReq{}
    s := service.DepartmentBudget{}
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
	list := make([]dto.DepartmentBudgetListResp, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取部门预算失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取部门预算
// @Summary 获取部门预算
// @Description 获取部门预算
// @Tags 部门预算
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=dto.DepartmentBudgetGetResp} "{"code": 200, "data": [...]}"
// @Router /api/v1/uc/department-budget/{id} [get]
// @Security Bearer
func (e DepartmentBudget) Get(c *gin.Context) {
	req := dto.DepartmentBudgetGetReq{}
	s := service.DepartmentBudget{}
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
	var object dto.DepartmentBudgetGetResp

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取部门预算失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK( object, "查询成功")
}

// Init 初始化部门预算 每月初
func (e DepartmentBudget) Init(c *gin.Context) {
	s := service.DepartmentBudget{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	echoMsg, err := s.Init()
	if err != nil {
		c.String(500, err.Error())
	} else {
		c.String(200, echoMsg)
	}
}

//// Insert 创建部门预算
//// @Summary 创建部门预算
//// @Description 创建部门预算
//// @Tags 部门预算
//// @Accept application/json
//// @Product application/json
//// @Param data body dto.DepartmentBudgetInsertReq true "data"
//// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
//// @Router /api/v1/uc/department-budget [post]
//// @Security Bearer
//func (e DepartmentBudget) Insert(c *gin.Context) {
//    req := dto.DepartmentBudgetInsertReq{}
//    s := service.DepartmentBudget{}
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
//		e.Error(500, err, fmt.Sprintf("创建部门预算失败，\r\n失败信息 %s", err.Error()))
//        return
//	}
//
//	e.OK(req.GetId(), "创建成功")
//}
//
//// Update 修改部门预算
//// @Summary 修改部门预算
//// @Description 修改部门预算
//// @Tags 部门预算
//// @Accept application/json
//// @Product application/json
//// @Param id path int true "id"
//// @Param data body dto.DepartmentBudgetUpdateReq true "body"
//// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
//// @Router /api/v1/uc/department-budget/{id} [put]
//// @Security Bearer
//func (e DepartmentBudget) Update(c *gin.Context) {
//    req := dto.DepartmentBudgetUpdateReq{}
//    s := service.DepartmentBudget{}
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
//		e.Error(500, err, fmt.Sprintf("修改部门预算失败，\r\n失败信息 %s", err.Error()))
//        return
//	}
//	e.OK( req.GetId(), "修改成功")
//}
//
//// Delete 删除部门预算
//// @Summary 删除部门预算
//// @Description 删除部门预算
//// @Tags 部门预算
//// @Param data body dto.DepartmentBudgetDeleteReq true "body"
//// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
//// @Router /api/v1/uc/department-budget [delete]
//// @Security Bearer
//func (e DepartmentBudget) Delete(c *gin.Context) {
//    s := service.DepartmentBudget{}
//    req := dto.DepartmentBudgetDeleteReq{}
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
//		e.Error(500, err, fmt.Sprintf("删除部门预算失败，\r\n失败信息 %s", err.Error()))
//        return
//	}
//	e.OK( req.GetId(), "删除成功")
//}
