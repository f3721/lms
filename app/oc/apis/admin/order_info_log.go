package admin

import (
	"fmt"
	"go-admin/common/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/oc/models"
	service "go-admin/app/oc/service/admin"
	"go-admin/app/oc/service/admin/dto"
	"go-admin/common/actions"
)

type OrderInfoLog struct {
	api.Api
}

// GetPage 获取订单日志列表
// @Summary 获取订单日志列表
// @Description 获取订单日志列表
// @Tags 订单日志
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.OrderInfoLog}} "{"code": 200, "data": [...]}"
// @Router /api/v1/oc/order-info-log [get]
// @Security Bearer
func (e OrderInfoLog) GetPage(c *gin.Context) {
    req := dto.OrderInfoLogGetPageReq{}
    s := service.OrderInfoLog{}
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
	list := make([]models.OrderInfoLog, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取订单日志失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取订单日志
// @Summary 获取订单日志
// @Description 获取订单日志
// @Tags 订单日志
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.OrderInfoLog} "{"code": 200, "data": [...]}"
// @Router /api/v1/oc/order-info-log/{id} [get]
// @Security Bearer
func (e OrderInfoLog) Get(c *gin.Context) {
	req := dto.OrderInfoLogGetReq{}
	s := service.OrderInfoLog{}
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
	var object utils.OperateLogDetailResp

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取订单日志失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK( object, "查询成功")
}
