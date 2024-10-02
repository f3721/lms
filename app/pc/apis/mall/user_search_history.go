package mall

import (
	"fmt"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	"go-admin/common/middleware/mall_handler"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	service "go-admin/app/pc/service/mall"
	"go-admin/app/pc/service/mall/dto"
	"go-admin/common/actions"
)

type UserSearchHistory struct {
	api.Api
}

// GetPage 获取用户搜索历史记录表列表
// @Summary 获取用户搜索历史记录表列表
// @Description 获取用户搜索历史记录表列表
// @Tags 用户搜索历史记录表
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.UserSearchHistory}} "{"code": 200, "data": [...]}"
// @Router /api/v1/pc/user-search-history [get]
// @Security Bearer
func (e UserSearchHistory) GetPage(c *gin.Context) {
	req := dto.UserSearchHistoryGetPageReq{}
	s := service.UserSearchHistory{}
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
	list := make([]dto.UserSearchHistoryGetPageResp, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取用户搜索历史记录表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}
