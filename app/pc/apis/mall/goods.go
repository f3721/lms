package mall

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"go-admin/common/middleware/mall_handler"

	service "go-admin/app/pc/service/mall"
	"go-admin/app/pc/service/mall/dto"
	"go-admin/common/actions"
)

type Goods struct {
	api.Api
}

// GetPage 获取商品管理表列表
// @Summary 获取商品管理表列表
// @Description 获取商品管理表列表
// @Tags 商品管理表
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Param categoryId query int false "分类ID"
// @Param brandId query []int false "品牌ID"
// @Param keyword query string false "关键词"
// @Param marketPriceOrder query string false "价格排序"
// @Success 200 {object} response.Response{data=response.Page{list=dto.List}} "{"code": 200, "data": [...]}"
// @Router /api/v1/pc/mall/product [get]
// @Security Bearer
func (e Goods) GetPage(c *gin.Context) {
	req := dto.GoodsGetPageReq{}
	s := service.Goods{}
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

	// 当前登录的仓
	req.UserId = user.GetUserId(c)
	req.UserName = user.GetUserName(c)
	req.WarehouseCode = mall_handler.GetUserConfig(c).SelectedWarehouseCode

	p := actions.GetPermissionFromContext(c)
	list := dto.List{
		Product:     nil,
		BrandAll:    nil,
		BrandFilter: nil,
		Category:    nil,
		CategoryNav: nil,
		BrandNav:    nil,
	}
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取商品管理表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取商品管理表
// @Summary 获取商品管理表
// @Description 获取商品管理表
// @Tags 商品管理表
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=dto.GoodsGetResp} "{"code": 200, "data": [...]}"
// @Router /api/v1/pc/mall/product/{id} [get]
// @Security Bearer
func (e Goods) Get(c *gin.Context) {
	req := dto.GoodsGetReq{}
	s := service.Goods{}
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
	var object dto.GoodsGetResp

	req.UserId = user.GetUserId(c)
	// 当前登录的仓
	req.WarehouseCode = mall_handler.GetUserConfig(c).SelectedWarehouseCode

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取商品管理表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// GetMiniProgramHomeFilter 获取小程序端筛选项
// @Summary 获取小程序端筛选项
// @Description 获取小程序端筛选项
// @Tags 商品管理表
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=dto.GoodsGetResp} "{"code": 200, "data": [...]}"
// @Router /api/v1/pc/mall/product/filter-items [get]
// @Security Bearer
func (e Goods) GetMiniProgramHomeFilter(c *gin.Context) {
	req := dto.GoodsGetReq{}
	s := service.Goods{}
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
	var object dto.MiniProgramHomeFilter

	req.UserId = user.GetUserId(c)
	// 当前登录的仓
	req.WarehouseCode = mall_handler.GetUserConfig(c).SelectedWarehouseCode

	err = s.GetMiniProgramHomeFilter(&req, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取商品管理表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}
