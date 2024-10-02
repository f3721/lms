package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	service "go-admin/app/pc/service/admin"
	"go-admin/app/pc/service/admin/dto"
	"go-admin/common/actions"
)

type Category struct {
	api.Api
}

// GetPage 获取产线表列表
// @Summary 获取产线表列表
// @Description 获取产线表列表
// @Tags 产线表
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.Category}} "{"code": 200, "data": [...]}"
// @Router /api/v1/pc/category [get]
// @Security Bearer
func (e Category) GetPage(c *gin.Context) {
	req := dto.CategoryGetPageReq{}
	s := service.Category{}
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
	list := make([]dto.CategoryGetPageResp, 0)

	var count int
	err = s.GetPage(&req, p, &list, &count).Error
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取产线表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取产线表
// @Summary 获取产线表
// @Description 获取产线表
// @Tags 产线表
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.Category} "{"code": 200, "data": [...]}"
// @Router /api/v1/pc/category/{id} [get]
// @Security Bearer
func (e Category) Get(c *gin.Context) {
	req := dto.CategoryGetReq{}
	s := service.Category{}
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
	var object dto.CategoryGetResp

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取产线表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// GetList 获取子产线列表
// @Summary 获取子产线列表
// @Description 获取子产线列表
// @Tags 获取子产线列表
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=dto.CategoryChildList} "{"code": 200, "data": [...]}"
// @Router /api/v1/pc/category/{id} [get]
// @Security Bearer
func (e Category) GetList(c *gin.Context) {
	req := dto.CategoryGetReq{}
	s := service.Category{}
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
	object := make([]dto.CategoryGetPageResp, 0)

	p := actions.GetPermissionFromContext(c)
	err = s.GetList(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取产线表失败，\r\n失败信息 %s", err.Error()))
		return
	}
	for i, _ := range object {
		object[i].Addchild = true
	}

	e.OK(object, "查询成功")
}

// Insert 创建产线表
// @Summary 创建产线表
// @Description 创建产线表
// @Tags 产线表
// @Accept application/json
// @Product application/json
// @Param data body dto.CategoryInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/pc/category [post]
// @Security Bearer
func (e Category) Insert(c *gin.Context) {
	req := dto.CategoryInsertReq{}
	s := service.Category{}
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
		e.Error(500, err, fmt.Sprintf("创建产线表失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "创建成功")
}

// Update 修改产线表
// @Summary 修改产线表
// @Description 修改产线表
// @Tags 产线表
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.CategoryUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/pc/category/{id} [put]
// @Security Bearer
func (e Category) Update(c *gin.Context) {
	req := dto.CategoryUpdateReq{}
	s := service.Category{}
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

	err = s.Update(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改产线表失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Sort 排序
// @Summary 排序
// @Description 排序
// @Tags 排序
// @Accept application/json
// @Product application/json
// @Param data body dto.SortReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/pc/category/sort [post]
// @Security Bearer
func (e Category) Sort(c *gin.Context) {
	req := dto.SortReq{}
	s := service.Category{}
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

	err = s.Sort(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("排序失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK("", "修改成功")
}

// Delete 删除产线表
// @Summary 删除产线表
// @Description 删除产线表
// @Tags 产线表
// @Param data body dto.CategoryDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/pc/category [delete]
// @Security Bearer
func (e Category) Delete(c *gin.Context) {
	s := service.Category{}
	req := dto.CategoryDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除产线表失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}

// GetCategoryPath 父级获取产线
// @Summary 父级获取产线
// @Description 父级获取产线
// @Tags 产线表
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.Category} "{"code": 200, "data": [...]}"
// @Router /api/v1/pc/category-path/{id} [get]
// @Security Bearer
func (e Category) GetCategoryPath(c *gin.Context) {
	req := dto.CategoryPathReq{}
	s := service.Category{}
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

	result := s.GetCategoryPath(req.Id)

	e.OK(result, "查询成功")
}
