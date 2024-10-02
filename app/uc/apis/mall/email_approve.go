package mall

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"go-admin/common/middleware/mall_handler"

	"go-admin/app/uc/models"
	service "go-admin/app/uc/service/mall"
	"go-admin/app/uc/service/mall/dto"
	"go-admin/common/actions"
)

type EmailApprove struct {
	api.Api
}

// GetPage 获取审批流表列表
// @Summary 获取审批流表列表
// @Description 获取审批流表列表
// @Tags 审批流表
// @Param filterProcess query string false "审批流程过滤"
// @Param filterPerson query string false "领用人过滤"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.EmailApproveListCustomer}} "{"code": 200, "data": [...]}"
// @Router /api/v1/uc/mall/email-approve [get]
// @Security Bearer
func (e EmailApprove) GetPage(c *gin.Context) {
	req := dto.EmailApproveGetPageReq{}
	s := service.EmailApprove{}
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

	companyId, err := mall_handler.GetUserCompanyID(e.Context)
	if err != nil {
		e.Error(500, err, err.Error())
	}
	req.CompanyId = companyId

	p := actions.GetPermissionFromContext(c)
	list := make([]models.EmailApproveListCustomer, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取审批流表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取审批流表
// @Summary 获取审批流表
// @Description 获取审批流表
// @Tags 审批流表
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=dto.EmailApproveGetResp} "{"code": 200, "data": [...]}"
// @Router /api/v1/uc/mall/email-approve/{id} [get]
// @Security Bearer
func (e EmailApprove) Get(c *gin.Context) {
	req := dto.EmailApproveGetReq{}
	s := service.EmailApprove{}
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
	var object dto.EmailApproveGetResp

	p := actions.GetPermissionFromContext(c)
	companyId, err := mall_handler.GetUserCompanyID(e.Context)
	if err != nil {
		e.Error(500, err, err.Error())
	}
	req.CompanyId = companyId

	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取审批流表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Get 获取审批人和领用人
// @Summary 获取审批人和领用人
// @Description 获取审批人和领用人
// @Tags 审批流表
// @Success 200 {object} response.Response{data=dto.EmailApproveGetApproveAndRecipient} "{"code": 200, "data": [...]}"
// @Router /api/v1/uc/mall/email-approve/approve-recipient-users [get]
// @Security Bearer
func (e EmailApprove) GetApproveAndRecipient(c *gin.Context) {
	s := service.EmailApprove{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	var object dto.EmailApproveGetApproveAndRecipient

	p := actions.GetPermissionFromContext(c)
	companyId, err := mall_handler.GetUserCompanyID(e.Context)
	if err != nil {
		e.Error(500, err, err.Error())
	}

	err = s.GetAllApproveAndRecipientUsers(companyId, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取审批人和领用人失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建审批流表
// @Summary 创建审批流表
// @Description 创建审批流表
// @Tags 审批流表
// @Accept application/json
// @Product application/json
// @Param data body dto.EmailApproveInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/uc/mall/email-approve [post]
// @Security Bearer
func (e EmailApprove) Insert(c *gin.Context) {
	req := dto.EmailApproveInsertReq{}
	s := service.EmailApprove{}
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
	req.SetCreateByName(user.GetUserName(c))

	companyId, err := mall_handler.GetUserCompanyID(e.Context)
	if err != nil {
		e.Error(500, err, err.Error())
	}
	req.CompanyId = companyId

	err = s.Insert(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建审批流表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改审批流表
// @Summary 修改审批流表
// @Description 修改审批流表
// @Tags 审批流表
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.EmailApproveUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/uc/mall/email-approve/{id} [put]
// @Security Bearer
func (e EmailApprove) Update(c *gin.Context) {
	req := dto.EmailApproveUpdateReq{}
	s := service.EmailApprove{}
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
	req.SetUpdateByName(user.GetUserName(c))
	p := actions.GetPermissionFromContext(c)
	companyId, err := mall_handler.GetUserCompanyID(e.Context)
	if err != nil {
		e.Error(500, err, err.Error())
	}
	req.CompanyId = companyId

	err = s.Update(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改审批流表失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除审批流表
// @Summary 删除审批流表
// @Description 删除审批流表
// @Tags 审批流表
// @Param data body dto.EmailApproveDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/uc/mall/email-approve [delete]
// @Security Bearer
func (e EmailApprove) Delete(c *gin.Context) {
	s := service.EmailApprove{}
	req := dto.EmailApproveDeleteReq{}
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
	req.SetUpdateBy(user.GetUserId(c))
	req.SetUpdateByName(user.GetUserName(c))
	companyId, err := mall_handler.GetUserCompanyID(e.Context)
	if err != nil {
		e.Error(500, err, err.Error())
	}
	req.CompanyId = companyId

	err = s.Remove(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("删除审批流表失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}

// Get 获取审批流程
// @Summary 获取审批流程
// @Description 获取审批流程
// @Tags 审批流表
// @Success 200 {object} response.Response{data=dto.EmailApproveWorkflow} "{"code": 200, "data": [...]}"
// @Router /api/v1/uc/mall/email-approve/workflow [get]
// @Security Bearer
func (e EmailApprove) Workflow(c *gin.Context) {
	s := service.EmailApprove{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	var object []dto.EmailApproveWorkflow
	p := actions.GetPermissionFromContext(c)

	err = s.Workflow(user.GetUserId(c), p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取审批流程失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}
