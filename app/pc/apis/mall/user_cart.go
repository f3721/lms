package mall

import (
	"fmt"
	"go-admin/common/middleware/mall_handler"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	service "go-admin/app/pc/service/mall"
	"go-admin/app/pc/service/mall/dto"
	"go-admin/common/actions"
)

type UserCart struct {
	api.Api
}

// GetPage 获取购物车信息表列表
// @Summary 获取购物车信息表列表
// @Description 获取购物车信息表列表
// @Tags 购物车信息表
// @Success 200 {object} response.Response{data=dto.UserCartGetProductForCartPageResp} "{"code": 200, "data": [...]}"
// @Router /api/v1/pc/mall/user-cart [get]
// @Security Bearer
func (e UserCart) GetPage(c *gin.Context) {
	req := dto.UserCartGetProductForCartPageReq{}
	s := service.UserCart{}
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

	req.UserId = user.GetUserId(c)
	req.WarehouseCode = mall_handler.GetUserCurrentWarehouseCode(c)
	p := actions.GetPermissionFromContext(c)
	data := dto.UserCartGetProductForCartPageResp{}

	err = s.GetCartDataPage(&req, p, &data)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取购物车信息表失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(data, "查询成功")
}

// Insert 创建购物车信息表
// @Summary 创建购物车信息表
// @Description 创建购物车信息表
// @Tags 购物车信息表
// @Accept application/json
// @Product application/json
// @Param data body dto.UserCartInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/pc/mall/user-cart [post]
// @Security Bearer
func (e UserCart) Insert(c *gin.Context) {
	req := dto.UserCartInsertReq{}
	s := service.UserCart{}
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
	req.WarehouseCode = mall_handler.GetUserCurrentWarehouseCode(c)
	err = s.Insert(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建购物车信息表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改购物车信息表
// @Summary 修改购物车信息表
// @Description 修改购物车信息表
// @Tags 购物车信息表
// @Accept application/json
// @Product application/json
// @Param data body dto.UserCartUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/pc/mall/user-cart [put]
// @Security Bearer
func (e UserCart) Update(c *gin.Context) {
	req := dto.UserCartUpdateReq{}
	s := service.UserCart{}
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
	req.UserId = user.GetUserId(c)
	req.WarehouseCode = mall_handler.GetUserCurrentWarehouseCode(c)
	p := actions.GetPermissionFromContext(c)

	err = s.Update(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改购物车信息表失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(0, "修改成功")
}

// Delete 删除购物车信息表
// @Summary 删除购物车信息表
// @Description 删除购物车信息表
// @Tags 购物车信息表
// @Param goodsId query int false "goodsId"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/pc/mall/user-cart [delete]
// @Security Bearer
func (e UserCart) Delete(c *gin.Context) {
	s := service.UserCart{}
	req := dto.UserCartDeleteReq{}
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
	req.UserId = user.GetUserId(c)
	req.WarehouseCode = mall_handler.GetUserCurrentWarehouseCode(c)
	p := actions.GetPermissionFromContext(c)

	err = s.Remove(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("删除购物车信息表失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(0, "删除成功")
}

// SelectOne 购物车SelectOne
// @Summary 购物车SelectOne
// @Description 购物车SelectOne
// @Tags 购物车信息表
// @Param data body dto.UserCartSelectOneReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "选择成功"}"
// @Router /api/v1/pc/mall/user-cart/select-one [post]
// @Security Bearer
func (e UserCart) SelectOne(c *gin.Context) {
	s := service.UserCart{}
	req := dto.UserCartSelectOneReq{}
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
	req.UserId = user.GetUserId(c)
	req.WarehouseCode = mall_handler.GetUserCurrentWarehouseCode(c)

	p := actions.GetPermissionFromContext(c)

	err = s.SelectOne(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("选择购物车信息表失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(0, "选择成功")
}

// SelectAll 购物车SelectAll
// @Summary 购物车SelectAll
// @Description 购物车SelectAll
// @Tags 购物车信息表
// @Success 200 {object} response.Response	"{"code": 200, "message": "选择成功"}"
// @Router /api/v1/pc/mall/user-cart/select-all [post]
// @Security Bearer
func (e UserCart) SelectAll(c *gin.Context) {
	s := service.UserCart{}
	req := dto.UserCartSelectAllReq{}
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
	req.UserId = user.GetUserId(c)
	req.WarehouseCode = mall_handler.GetUserCurrentWarehouseCode(c)
	p := actions.GetPermissionFromContext(c)

	err = s.SelectAll(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("选择购物车信息表失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(0, "选择成功")
}

// UnSelectAll 购物车UnSelectAll
// @Summary 购物车UnSelectAll
// @Description 购物车UnSelectAll
// @Tags 购物车信息表
// @Success 200 {object} response.Response	"{"code": 200, "message": "选择成功"}"
// @Router /api/v1/pc/mall/user-cart/unselect-all [post]
// @Security Bearer
func (e UserCart) UnSelectAll(c *gin.Context) {
	s := service.UserCart{}
	req := dto.UserCartUnSelectAllReq{}
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
	req.UserId = user.GetUserId(c)
	req.WarehouseCode = mall_handler.GetUserCurrentWarehouseCode(c)
	p := actions.GetPermissionFromContext(c)

	err = s.UnSelectAll(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("选择购物车信息表失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(0, "选择成功")
}

// ClearSelect 购物车ClearSelect
// @Summary 购物车ClearSelect
// @Description 购物车ClearSelect
// @Tags 购物车信息表
// @Success 200 {object} response.Response	"{"code": 200, "message": "成功"}"
// @Router /api/v1/pc/mall/user-cart/clear-select [post]
// @Security Bearer
func (e UserCart) ClearSelect(c *gin.Context) {
	s := service.UserCart{}
	req := dto.UserCartClearSelectReq{}
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
	req.UserId = user.GetUserId(c)
	req.WarehouseCode = mall_handler.GetUserCurrentWarehouseCode(c)
	p := actions.GetPermissionFromContext(c)

	err = s.ClearSelect(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("ClearSelect购物车信息表失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(0, "成功")
}

// ClearSelect 购物车ClearInvalid
// @Summary 购物车ClearInvalid
// @Description 购物车ClearInvalid
// @Tags 购物车信息表
// @Success 200 {object} response.Response	"{"code": 200, "message": "成功"}"
// @Router /api/v1/pc/mall/user-cart/clear-invalid [post]
// @Security Bearer
func (e UserCart) ClearInvalid(c *gin.Context) {
	s := service.UserCart{}
	req := dto.UserCartClearSelectReq{}
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
	req.UserId = user.GetUserId(c)
	req.WarehouseCode = mall_handler.GetUserCurrentWarehouseCode(c)
	p := actions.GetPermissionFromContext(c)

	err = s.ClearInvalid(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("ClearSelect购物车信息表失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(0, "成功")
}

// VerifyCart 购物车VerifyCart
// @Summary 购物车VerifyCart
// @Description 购物车VerifyCart
// @Tags 购物车信息表
// @Success 200 {object} response.Response	"{"code": 200, "message": "成功"}"
// @Router /api/v1/pc/mall/user-cart/verify-cart [get]
// @Security Bearer
func (e UserCart) VerifyCart(c *gin.Context) {
	s := service.UserCart{}
	req := dto.UserCartVerifyCartReq{}
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

	req.UserId = user.GetUserId(c)
	req.WarehouseCode = mall_handler.GetUserCurrentWarehouseCode(c)

	p := actions.GetPermissionFromContext(c)

	err = s.VerifyCart(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("VerifyCart购物车信息表失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(0, "成功")
}

// VerifyBuyNow 购物车VerifyBuyNow
// @Summary 购物车VerifyBuyNow
// @Description 购物车VerifyBuyNow
// @Tags 购物车信息表
// @Param goodsId query int false "goodsId"
// @Param quantity query int false "quantity"
// @Success 200 {object} response.Response	"{"code": 200, "message": "成功"}"
// @Router /api/v1/pc/mall/user-cart/verify-buy-now [get]
// @Security Bearer
func (e UserCart) VerifyBuyNow(c *gin.Context) {
	s := service.UserCart{}
	req := dto.UserCartVerifyBuyNowReq{}
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

	req.UserId = user.GetUserId(c)
	req.WarehouseCode = mall_handler.GetUserCurrentWarehouseCode(c)

	p := actions.GetPermissionFromContext(c)

	err = s.VerifyBuyNow(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("VerifyBuyNow购物车信息表失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(0, "成功")
}

// BatchAdd 批量添加购物车信息表
// @Summary 批量添加购物车信息表
// @Description 批量添加购物车信息表
// @Tags 购物车信息表
// @Accept application/json
// @Product application/json
// @Param data body dto.UserCartBatchAddReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "批量添加成功"}"
// @Router /api/v1/pc/mall/user-cart/batch-add [post]
// @Security Bearer
func (e UserCart) BatchAdd(c *gin.Context) {
	req := dto.UserCartBatchAddReq{}
	s := service.UserCart{}
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
	req.WarehouseCode = mall_handler.GetUserCurrentWarehouseCode(c)
	err = s.BatchAdd(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("批量添加购物车信息表，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(0, "创建成功")
}

// GetPageForOrder 获取领用单（购物车）商品列表
// @Summary 获取领用单（购物车）商品列表
// @Description 获取领用单（购物车）商品列表
// @Tags 购物车信息表
// @Success 200 {object} response.Response{data=dto.UserCartGetProductForOrderPageResp} "{"code": 200, "data": [...]}"
// @Router /api/v1/pc/mall/user-cart/order-cart-products [get]
// @Security Bearer
func (e UserCart) GetPageForOrder(c *gin.Context) {
	req := dto.UserCartGetProductForOrderPageReq{}
	s := service.UserCart{}
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

	req.UserId = user.GetUserId(c)
	req.WarehouseCode = mall_handler.GetUserCurrentWarehouseCode(c)
	p := actions.GetPermissionFromContext(c)
	data := dto.UserCartGetProductForOrderPageResp{}

	err = s.GetOrderProductsPage(&req, p, &data)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取领用单（购物车）商品列表，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(data, "查询成功")
}

// GetProductForOrderOnBuyNow 获取领用单（立即购买）商品列表
// @Summary 获取领用单（立即购买）商品列表
// @Description 获取领用单（立即购买）商品列表
// @Tags 购物车信息表
// @Param goodsId query int false "goodsId"
// @Param quantity query int false "quantity"
// @Success 200 {object} response.Response{data=dto.UserCartGetProductForOrderBuyNowResp} "{"code": 200, "data": [...]}"
// @Router /api/v1/pc/mall/user-cart/order-buy-now-product [get]
// @Security Bearer
func (e UserCart) GetProductForOrderOnBuyNow(c *gin.Context) {
	req := dto.UserCartGetProductForOrderBuyNowReq{}
	s := service.UserCart{}
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

	req.UserId = user.GetUserId(c)
	req.WarehouseCode = mall_handler.GetUserCurrentWarehouseCode(c)
	p := actions.GetPermissionFromContext(c)
	data := dto.UserCartGetProductForOrderBuyNowResp{}

	err = s.GetOrderProductForBuyNow(&req, p, &data)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取领用单（购物车）商品列表，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(data, "查询成功")
}

// SaleMoq 购物车SaleMoq
// @Summary 购物车SaleMoq
// @Description 购物车SaleMoq
// @Tags 购物车信息表
// @Param data body dto.UserCartGetProductForSaleMoqReq true "data"
// @Success 200 {object} response.Response{data=dto.UserCartGetProductForSaleMoqResp} "{"code": 200, "data": [...]}"
// @Router /api/v1/pc/mall/user-cart/sale-moq [post]
// @Security Bearer
func (e UserCart) SaleMoq(c *gin.Context) {
	s := service.UserCart{}
	req := dto.UserCartGetProductForSaleMoqReq{}
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
	req.WarehouseCode = mall_handler.GetUserCurrentWarehouseCode(c)
	p := actions.GetPermissionFromContext(c)
	data := dto.UserCartGetProductForSaleMoqResp{}

	err = s.SaleMoq(&req, p, &data)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("购物车SaleMoq失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(data, "查询成功")
}
