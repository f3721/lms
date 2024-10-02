package admin

import (
	"fmt"
	"go-admin/common/actions"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	service "go-admin/app/wc/service/admin"
	"go-admin/app/wc/service/admin/dto"
)

type PrintsJson struct {
	api.Api
}

// CommonPrints 获取打印数据接口
// @Summary 获取打印数据接口
// @Description 共用的获取打印数据接口-json
// @Tags WC公用接口
// @Param id path int false "id"
// @Param print-type path string false "打印类型:outboundprint>打印出库单 pickingprint>打印拣货单 entryprint>入库单打印"
// @Success 200 {object} response.Response{data=dto.CommonPrintsResp} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/common-prints/:print-type/:id [get]
// @Security Bearer
func (p PrintsJson) CommonPrints(c *gin.Context) {
	req := dto.CommonPrintsReq{}
	s := service.PrintsJson{}
	err := p.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		p.Logger.Error(err)
		p.Error(500, err, err.Error())
		return
	}
	var object dto.CommonPrintsResp

	Per := actions.GetPermissionFromContext(c)
	switch {
	case req.PrintType == "outboundprint":
		outboundSer := service.StockOutbound{}
		p.MakeService(&outboundSer.Service)
		err = s.OutboundPrint(&req, Per, &outboundSer, &object)
	case req.PrintType == "pickingprint":
		outboundSer := service.StockOutbound{}
		p.MakeService(&outboundSer.Service)
		err = s.StockPrintPicking(&req, Per, &outboundSer, &object)
	case req.PrintType == "entryprint":
		entrySer := service.StockEntry{}
		p.MakeService(&entrySer.Service)
		err = s.StockEntryPrints(&req, Per, &entrySer, &object)
	case req.PrintType == "qualityreportprint":
		qualitySer := service.QualityCheck{}
		p.MakeService(&qualitySer.Service)
		err = s.QualityReportPrint(&req, Per, &qualitySer, &object)
	default:
		p.Error(404, err, "page not found")
		return

	}
	if err != nil {
		p.Error(500, err, fmt.Sprintf("打印失败，\r\n失败信息 %s", err.Error()))
		return
	}
	p.OK(object, "成功")
}
