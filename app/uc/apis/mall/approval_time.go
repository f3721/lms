package mall

import (
	"fmt"
	service "go-admin/app/uc/service/mall"
	"go-admin/app/uc/service/mall/dto"
	"go-admin/common/actions"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
)

type ApprovalTime struct {
	api.Api
}

// GetPage 获取审批时间列表
// @Summary 获取审批时间列表
// @Description 获取审批时间列表
// @Tags 审批时间
// @Success 200 {object} response.Response{data=response.Page{list=[]dto.ApprovalTimeInsertReq}} "{"code": 200, "data": [...]}"
// @Router /api/v1/uc/mall/approval-time [get]
// @Security Bearer
func (e ApprovalTime) GetPage(c *gin.Context) {
	s := service.ApprovalTime{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	// 获取用户ID
	userId := user.GetUserId(c)
	p := actions.GetPermissionFromContext(c)

	res, err := s.GetPage(userId, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取审批时间失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(res, "查询成功")
}

// Get 获取审批时间
// @Summary 获取审批时间
// @Description 获取审批时间
// @Tags 审批时间
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=dto.ApprovalTimeInsertReq} "{"code": 200, "data": [...]}"
// @Router /api/v1/uc/mall/approval-time/{id} [get]
// @Security Bearer
func (e ApprovalTime) Get(c *gin.Context) {
	req := dto.ApprovalTimeGetReq{}
	s := service.ApprovalTime{}
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

	// 获取用户ID
	userId := user.GetUserId(c)
	p := actions.GetPermissionFromContext(c)
	res, err := s.Get(&req, userId, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取审批时间失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(res, "查询成功")
}

// Insert 创建审批时间
// @Summary 创建审批时间
// @Description 创建审批时间
// @Tags 审批时间
// @Accept application/json
// @Product application/json
// @Param data body dto.ApprovalTimeInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/uc/mall/approval-time [post]
// @Security Bearer
func (e ApprovalTime) Insert(c *gin.Context) {
	req := dto.ApprovalTimeInsertReq{}
	s := service.ApprovalTime{}
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

	// 获取用户ID
	userId := user.GetUserId(c)

	// 逻辑处理
	err = s.Insert(&req, userId)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建审批时间失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK("", "创建成功")
}

// Update 修改审批时间
// @Summary 修改审批时间
// @Description 修改审批时间
// @Tags 审批时间
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.ApprovalTimeUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/uc/mall/approval-time/{id} [put]
// @Security Bearer
func (e ApprovalTime) Update(c *gin.Context) {
	req := dto.ApprovalTimeUpdateReq{}
	s := service.ApprovalTime{}
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

	// 增加更新字段
	userId := user.GetUserId(c)
	req.SetUpdateBy(userId)
	p := actions.GetPermissionFromContext(c)

	err = s.Update(&req, userId, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改审批时间失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK("", "修改成功")
}

// Delete 删除审批时间
// @Summary 删除审批时间
// @Description 删除审批时间
// @Tags 审批时间
// @Param data body dto.ApprovalTimeDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/uc/mall/approval-time [delete]
// @Security Bearer
func (e ApprovalTime) Delete(c *gin.Context) {
	s := service.ApprovalTime{}
	req := dto.ApprovalTimeDeleteReq{}
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

	userId := user.GetUserId(c)
	p := actions.GetPermissionFromContext(c)
	err = s.Remove(&req, userId, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("删除审批时间失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK("", "删除成功")
}
