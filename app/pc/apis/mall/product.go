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

type Product struct {
	api.Api
}

// GetPage 获取商品档案列表
// @Summary 获取商品档案列表
// @Description 获取商品档案列表
// @Tags 商品档案
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.Product}} "{"code": 200, "data": [...]}"
// @Router /api/v1/pc/product [get]
// @Security Bearer
func (e Product) GetPage(c *gin.Context) {
	req := dto.ProductGetPageReq{}
	s := service.Product{}
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
	list := make([]models.Product, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取商品档案失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取商品档案
// @Summary 获取商品档案
// @Description 获取商品档案
// @Tags 商品档案
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.Product} "{"code": 200, "data": [...]}"
// @Router /api/v1/pc/product/{id} [get]
// @Security Bearer
func (e Product) Get(c *gin.Context) {
	req := dto.ProductGetReq{}
	s := service.Product{}
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
	var object models.Product

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取商品档案失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}
