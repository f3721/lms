package admin

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	service "go-admin/app/oc/service/admin"
	"go-admin/app/oc/service/admin/dto"
	"go-admin/common/actions"
)

type CsApplyDetail struct {
	api.Api
}

// GetByCsNo 获取售后相关商品表
// @Summary 获取售后相关商品表
// @Description 获取售后相关商品表
// @Tags 售后相关商品表
// @Param csNo path string false "csNo"
// @Success 200 {object} response.Response{data=dto.CsApplyDetailGetRes} "{"code": 200, "data": [...]}"
// @Router /api/v1/oc/admin/cs-apply-detail/cs-apply-info/{csNo} [get]
// @Security Bearer
func (e CsApplyDetail) GetByCsNo(c *gin.Context) {
	req := dto.CsApplyDetailGetReq{}
	s := service.CsApplyDetail{}
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
	var object dto.CsApplyDetailGetRes

	p := actions.GetPermissionFromContext(c)
	err = s.GetCsApplyDetail(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取售后相关商品表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}
