package mall

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/pc/models"
	service "go-admin/app/pc/service/mall"
	"go-admin/app/pc/service/mall/dto"
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
	list := make([]models.Category, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
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
	var object models.Category

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
