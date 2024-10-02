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

type Brand struct {
	api.Api
}

// GetPage 获取品牌表列表
// @Summary 获取品牌表列表
// @Description 获取品牌表列表
// @Tags 品牌表
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.Brand}} "{"code": 200, "data": [...]}"
// @Router /api/v1/pc/brand [get]
// @Security Bearer
func (e Brand) GetPage(c *gin.Context) {
	req := dto.BrandGetPageReq{}
	s := service.Brand{}
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
	list := make([]models.Brand, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取品牌表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取品牌表
// @Summary 获取品牌表
// @Description 获取品牌表
// @Tags 品牌表
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.Brand} "{"code": 200, "data": [...]}"
// @Router /api/v1/pc/brand/{id} [get]
// @Security Bearer
func (e Brand) Get(c *gin.Context) {
	req := dto.BrandGetReq{}
	s := service.Brand{}
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
	var object models.Brand

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取品牌表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}
