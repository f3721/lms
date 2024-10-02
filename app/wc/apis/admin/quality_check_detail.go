package admin

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/wc/models"
	service "go-admin/app/wc/service/admin"
	"go-admin/app/wc/service/admin/dto"
	"go-admin/common/actions"
)

type QualityCheckDetail struct {
	api.Api
}

func (e QualityCheckDetail) GetPage(c *gin.Context) {
	req := dto.QualityCheckGetPageReq{}
	s := service.QualityCheckDetail{}
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
	list := make([]models.QualityCheckDetail, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取质检内容列表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

func (e QualityCheckDetail) Get(c *gin.Context) {
	req := dto.QualityCheckDetailGetReq{}
	s := service.QualityCheckDetail{}
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
	var object models.QualityCheckDetail

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取质检详情信息失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}
