package admin

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/uc/models"
	service "go-admin/app/uc/service/admin"
	"go-admin/app/uc/service/admin/dto"
	"go-admin/common/actions"
)

type RoleInfo struct {
	api.Api
}

// GetPage 获取角色信息表列表
// @Summary 获取角色信息表列表
// @Description 获取角色信息表列表
// @Tags 角色信息表
// @Param roleName query string false "角色名称"
// @Param roleStatus query int false "状态"
// @Param manageCompany query string false "判断公司是否可以管理该权限(1xx:punchout管理,x1x:EAS管理,xx1:普通管理)"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.RoleInfo}} "{"code": 200, "data": [...]}"
// @Router /api/v1/uc/admin/role [get]
// @Security Bearer
func (e RoleInfo) GetPage(c *gin.Context) {
	req := dto.RoleInfoGetPageReq{}
	s := service.RoleInfo{}
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
	list := make([]models.RoleInfo, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取角色信息表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取角色信息表
// @Summary 获取角色信息表
// @Description 获取角色信息表
// @Tags 角色信息表
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.RoleInfo} "{"code": 200, "data": [...]}"
// @Router /api/v1/uc/role/{id} [get]
// @Security Bearer
func (e RoleInfo) Get(c *gin.Context) {
	req := dto.RoleInfoGetReq{}
	s := service.RoleInfo{}
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
	var object models.RoleInfo

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取角色信息表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建角色信息表
// @Summary 创建角色信息表
// @Description 创建角色信息表
// @Tags 角色信息表
// @Accept application/json
// @Product application/json
// @Param data body dto.RoleInfoInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/uc/role [post]
// @Security Bearer
func (e RoleInfo) Insert(c *gin.Context) {
	req := dto.RoleInfoInsertReq{}
	s := service.RoleInfo{}
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
	err = s.Insert(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建角色信息表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改角色信息表
// @Summary 修改角色信息表
// @Description 修改角色信息表
// @Tags 角色信息表
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.RoleInfoUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/uc/role/{id} [put]
// @Security Bearer
func (e RoleInfo) Update(c *gin.Context) {
	req := dto.RoleInfoUpdateReq{}
	s := service.RoleInfo{}
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

	err = s.Update(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改角色信息表失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除角色信息表
// @Summary 删除角色信息表
// @Description 删除角色信息表
// @Tags 角色信息表
// @Param data body dto.RoleInfoDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/uc/role [delete]
// @Security Bearer
func (e RoleInfo) Delete(c *gin.Context) {
	s := service.RoleInfo{}
	req := dto.RoleInfoDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除角色信息表失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}
