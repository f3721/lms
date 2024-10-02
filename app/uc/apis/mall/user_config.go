package mall

import (
	"go-admin/app/uc/service/mall/dto"
	"go-admin/common/middleware/mall_handler"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"
)

type UserConfigAPI struct {
	api.Api
}

// SelectWarehouse 切换当前用户的仓库
// @Summary 切换当前用户的仓库
// @Description 切换当前用户的仓库
// @Tags Mall登陆
// @Accept application/json
// @Product application/json
// @Param data body dto.UserConfigSelectWarehouseReq true "data"
// @Success 200 {object} response.Response	"{"requestId":"40180c07-5a91-42fb-960a-aabec7deaca1","code":200,"msg":"仓库选择成功!","data":{"userId":1,"selectedWarehouseId":1,"selectedWarehouseCode":"121212"}}"
// @Router /api/v1/uc/mall/user/select-warehouse [post]
// @Security Bearer
func (e UserConfigAPI) SelectWarehouse(c *gin.Context) {
	req := dto.UserConfigSelectWarehouseReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	config := mall_handler.GetUserConfig(c)
	config.SelectedWarehouseId = req.WarehouseID
	config.SelectedWarehouseCode = req.WarehouseCode
	config.SelectedWarehouseName = req.WarehouseName
	mall_handler.SetUserConfig(c, config)
	e.OK(config, "仓库选择成功!")
}
