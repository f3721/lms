package admin

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/wc/models"
	service "go-admin/app/wc/service/admin"
	"go-admin/app/wc/service/admin/dto"
	"go-admin/common/actions"
)

type Region struct {
	api.Api
}

// GetPage 获取Region列表
// @Summary 获取Region列表
// @Description 获取Region列表
// @Tags Region
// @Param parentId query string false "parentId"
// @Param name query string false "name"
// @Param level query string false "level"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.Region}} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/region [get]

// @Security Bearer
func (e Region) GetPage(c *gin.Context) {
	req := dto.RegionGetPageReq{}
	s := service.Region{}
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
	list := make([]models.Region, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取Region失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取Region
// @Summary 获取Region
// @Description 获取Region
// @Tags Region
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.Region} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/region/{id} [get]

// @Security Bearer
func (e Region) Get(c *gin.Context) {
	req := dto.RegionGetReq{}
	s := service.Region{}
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
	var object models.Region

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取Region失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建Region
// @Summary 创建Region
// @Description 创建Region
// @Tags Region
// @Accept application/json
// @Product application/json
// @Param data body dto.RegionInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/region [post]
// @Security Bearer
func (e Region) Insert(c *gin.Context) {
	req := dto.RegionInsertReq{}
	s := service.Region{}
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

	err = s.Insert(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建Region失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改Region
// @Summary 修改Region
// @Description 修改Region
// @Tags Region
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.RegionUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/region/{id} [put]
// @Security Bearer
func (e Region) Update(c *gin.Context) {
	req := dto.RegionUpdateReq{}
	s := service.Region{}
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
		e.Error(500, err, fmt.Sprintf("修改Region失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除Region
// @Summary 删除Region
// @Description 删除Region
// @Tags Region
// @Param data body dto.RegionDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/region [delete]
// @Security Bearer
func (e Region) Delete(c *gin.Context) {
	s := service.Region{}
	req := dto.RegionDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除Region失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}

// InnerGetByIds Inner通过ids获取Region信息
// @Summary Inner通过ids获取Region信息
// @Description Inner通过ids获取Region信息
// @Tags InnerRegion
// @Param ids query string false "ids 多个用,分割"
// @Success 200 {object} response.Response{data=[]models.Region} "{"code": 200, "data": [...]}"
// @Router /inner/wc/admin/region/get-info [get]
// @Security Bearer
func (e Region) InnerGetByIds(c *gin.Context) {
	req := dto.InnerRegionGetByIdsReq{}
	s := service.Region{}
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

	list := make([]models.Region, 0)

	err = s.InnerGetByIds(&req, &list)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("Inner通过ids获取Region信息失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(list, "成功")
}
