package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	service "go-admin/app/uc/service/admin"
	"go-admin/app/uc/service/admin/dto"
	"go-admin/common/actions"
)

type UserCollect struct {
	api.Api
}

// GetPage 获取用户收藏列表
// @Summary 获取用户收藏列表
// @Description 获取用户收藏列表
// @Tags 用户收藏
// @Param skuCode query string false "产品SKU"
// @Param userId query int false "用户ID"
// @Param goodsId query int false "商品表id"
// @Param warehouseCode query string false "仓库code"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.UserCollect}} "{"code": 200, "data": [...]}"
// @Router /api/v1/uc/admin/user-collect [get]
// @Security Bearer
func (e UserCollect) GetPage(c *gin.Context) {
	req := dto.UserCollectGetPageReq{}
	s := service.UserCollect{}
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
	list := make([]dto.UserCollectGetListPageRes, 0)
	var count int64

	err = s.GetListPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取用户收藏失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

//// Get 获取用户收藏
//// @Summary 获取用户收藏
//// @Description 获取用户收藏
//// @Tags 用户收藏
//// @Param id path int false "id"
//// @Success 200 {object} response.Response{data=models.UserCollect} "{"code": 200, "data": [...]}"
//// @Router /api/v1/uc/admin/user-collect/{id} [get]
//// @Security Bearer
//func (e UserCollect) Get(c *gin.Context) {
//	req := dto.UserCollectGetReq{}
//	s := service.UserCollect{}
//	err := e.MakeContext(c).
//		MakeOrm().
//		Bind(&req).
//		MakeService(&s.Service).
//		Errors
//	if err != nil {
//		e.Logger.Error(err)
//		e.Error(500, err, err.Error())
//		return
//	}
//	var object models.UserCollect
//
//	p := actions.GetPermissionFromContext(c)
//	err = s.Get(&req, p, &object)
//	if err != nil {
//		e.Error(500, err, fmt.Sprintf("获取用户收藏失败，\r\n失败信息 %s", err.Error()))
//		return
//	}
//
//	e.OK(object, "查询成功")
//}
//
//// Insert 创建用户收藏
//// @Summary 创建用户收藏
//// @Description 创建用户收藏
//// @Tags 用户收藏
//// @Accept application/json
//// @Product application/json
//// @Param data body dto.UserCollectInsertReq true "data"
//// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
//// @Router /api/v1/uc/admin/user-collect [post]
//// @Security Bearer
//func (e UserCollect) Insert(c *gin.Context) {
//	req := dto.UserCollectInsertReq{}
//	s := service.UserCollect{}
//	err := e.MakeContext(c).
//		MakeOrm().
//		Bind(&req).
//		MakeService(&s.Service).
//		Errors
//	if err != nil {
//		e.Logger.Error(err)
//		e.Error(500, err, err.Error())
//		return
//	}
//	// 设置创建人
//	req.SetCreateBy(user.GetUserId(c))
//	req.SetCreateByName(user.GetUserName(c))
//	err = s.Insert(&req)
//	if err != nil {
//		e.Error(500, err, fmt.Sprintf("创建用户收藏失败，\r\n失败信息 %s", err.Error()))
//		return
//	}
//
//	e.OK(req.GetId(), "创建成功")
//}
//
//// Update 修改用户收藏
//// @Summary 修改用户收藏
//// @Description 修改用户收藏
//// @Tags 用户收藏
//// @Accept application/json
//// @Product application/json
//// @Param id path int true "id"
//// @Param data body dto.UserCollectUpdateReq true "body"
//// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
//// @Router /api/v1/uc/admin/user-collect/{id} [put]
//// @Security Bearer
//func (e UserCollect) Update(c *gin.Context) {
//	req := dto.UserCollectUpdateReq{}
//	s := service.UserCollect{}
//	err := e.MakeContext(c).
//		MakeOrm().
//		Bind(&req).
//		MakeService(&s.Service).
//		Errors
//	if err != nil {
//		e.Logger.Error(err)
//		e.Error(500, err, err.Error())
//		return
//	}
//	req.SetUpdateBy(user.GetUserId(c))
//	p := actions.GetPermissionFromContext(c)
//
//	err = s.Update(&req, p)
//	if err != nil {
//		e.Error(500, err, fmt.Sprintf("修改用户收藏失败，\r\n失败信息 %s", err.Error()))
//		return
//	}
//	e.OK(req.GetId(), "修改成功")
//}
//
//// Delete 删除用户收藏
//// @Summary 删除用户收藏
//// @Description 删除用户收藏
//// @Tags 用户收藏
//// @Param data body dto.UserCollectDeleteReq true "body"
//// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
//// @Router /api/v1/uc/admin/user-collect [delete]
//// @Security Bearer
//func (e UserCollect) Delete(c *gin.Context) {
//	s := service.UserCollect{}
//	req := dto.UserCollectDeleteReq{}
//	err := e.MakeContext(c).
//		MakeOrm().
//		Bind(&req).
//		MakeService(&s.Service).
//		Errors
//	if err != nil {
//		e.Logger.Error(err)
//		e.Error(500, err, err.Error())
//		return
//	}
//
//	// req.SetUpdateBy(user.GetUserId(c))
//	p := actions.GetPermissionFromContext(c)
//
//	err = s.Remove(&req, p)
//	if err != nil {
//		e.Error(500, err, fmt.Sprintf("删除用户收藏失败，\r\n失败信息 %s", err.Error()))
//		return
//	}
//	e.OK(req.GetId(), "删除成功")
//}
