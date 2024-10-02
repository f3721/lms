package mall

import (
	"errors"
	"fmt"
	"go-admin/common/excel"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	service "go-admin/app/oc/service/mall"
	"go-admin/app/oc/service/mall/dto"
	"go-admin/common/actions"
)

type Preso struct {
	api.Api
}

// SubmitApproval 提交审批
// @Summary 提交审批
// @Description 提交审批
// @Tags 领用审批
// @Accept application/json
// @Product application/json
// @Param data body dto.PresoSubmitApprovalReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "提交成功"}"
// @Router /api/v1/oc/mall/preso/submit-approval [post]
// @Security Bearer
func (e Preso) SubmitApproval(c *gin.Context) {
	req := dto.PresoSubmitApprovalReq{}
	s := service.Preso{}
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

	respData := make(map[string]string)
	err = s.SubmitApproval(&req, &respData)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("提交审批失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(respData, "提交成功")
}

// GetPage 领用申请+领用审批列表
// @Summary 领用申请+领用审批列表
// @Description 领用申请+领用审批列表
// @Tags 领用审批
// @Param keyword query string false "关键词"
// @Param approveStatus query int false "审批状态"
// @Param type query int false "类型: 1-领用申请 2-领用审批"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]dto.PresoGetPageResp}} "{"code": 200, "data": [...]}"
// @Router /api/v1/oc/mall/preso [get]
// @Security Bearer
func (e Preso) GetPage(c *gin.Context) {
    req := dto.PresoGetPageReq{}
    s := service.Preso{}
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
	list := make([]dto.PresoGetPageResp, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取列表失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// GetPageCount 领用申请tab数量
// @Summary tab数量
// @Description tab数量
// @Tags 领用审批
// @Success 200 {object} response.Response{data=map[string]int} "{"code": 200, "data": [...]}"
// @Router /api/v1/oc/mall/preso/count [get]
// @Security Bearer
func (e Preso) GetPageCount(c *gin.Context) {
	req := dto.PresoGetPageReq{}
	s := service.Preso{}
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
	list := make([]dto.PresoGetPageResp, 0)
	req.Type = 1

	err = s.GetPageCount(&req, p, &list)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取列表统计失败，\r\n失败信息 %s", err.Error()))
		return
	}
	resp := map[string]int{
		"approveStatus10": 0,
		"approveStatus1": 0,
		"approveStatus-1": 0,
		"approveStatus-2": 0,
		"approveStatus-3": 0,
	}
	for _, v := range list {
		if v.ApproveStatus == 10 || v.ApproveStatus == 0 {
			resp["approveStatus10"]++
		} else if v.ApproveStatus == 1 {
			resp["approveStatus1"]++
		} else if v.ApproveStatus == -1 {
			resp["approveStatus-1"]++
		} else if v.ApproveStatus == -2 {
			resp["approveStatus-2"]++
		} else if v.ApproveStatus == -3 {
			resp["approveStatus-3"]++
		}
	}
	e.OK(resp, "查询成功")
}


// GetApprovePageCount 领用审批tab数量
// @Summary tab数量
// @Description tab数量
// @Tags 领用审批
// @Success 200 {object} response.Response{data=map[string]int} "{"code": 200, "data": [...]}"
// @Router /api/v1/oc/mall/preso/count [get]
// @Security Bearer
func (e Preso) GetApprovePageCount(c *gin.Context) {
	req := dto.PresoGetPageReq{}
	s := service.Preso{}
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
	list := make([]dto.PresoGetPageResp, 0)
	req.Type = 2

	err = s.GetPageCount(&req, p, &list)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取列表统计失败，\r\n失败信息 %s", err.Error()))
		return
	}
	resp := map[string]int{
		"approveStatus10": 0,
		"approveStatus1": 0,
		"approveStatus-1": 0,
		"approveStatus-2": 0,
	}
	// 驳回、超时
	for _, v := range list {
		 if v.ApproveStatus == -1 {
			resp["approveStatus-1"]++
		} else if v.ApproveStatus == -2 {
			resp["approveStatus-2"]++
		}
	}
	// 已审批
	list1 := make([]dto.PresoGetPageResp, 0)
	req.ApproveStatus = 1
	err = s.GetPageCount(&req, p, &list1)
	resp["approveStatus1"] = len(list1)
	// 待审批
	list10 := make([]dto.PresoGetPageResp, 0)
	req.ApproveStatus = 10
	err = s.GetPageCount(&req, p, &list10)
	resp["approveStatus10"] = len(list10)

	e.OK(resp, "查询成功")
}

// Get 获取预订单
// @Summary 获取预订单
// @Description 获取预订单
// @Tags 领用审批
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=dto.PresoGetResp} "{"code": 200, "data": [...]}"
// @Router /api/v1/oc/mall/preso/{id} [get]
// @Security Bearer
func (e Preso) Get(c *gin.Context) {
	req := dto.PresoGetReq{}
	s := service.Preso{}
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
	var object dto.PresoGetResp

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取预订单失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK( object, "查询成功")
}

// FinishApproval 完成审批
// @Summary 完成审批
// @Description 完成审批
// @Tags 领用审批
// @Accept application/json
// @Product application/json
// @Param data body dto.PresoFinishApprovalReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "审批成功"}"
// @Router /api/v1/oc/mall/preso/finish-approval [post]
// @Security Bearer
func (e Preso) FinishApproval(c *gin.Context) {
	req := dto.PresoFinishApprovalReq{}
	s := service.Preso{}
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
	//// 设置创建人
	//req.CreateBy = user.GetUserId(c)
	//req.CreateByName = user.GetUserName(c)

	err = s.FinishApproval(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("审批失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(nil, "审批成功")
}

// BatchApproval 批量审批
// @Summary 批量审批
// @Description 批量审批
// @Tags 领用审批
// @Accept application/json
// @Product application/json
// @Param data body dto.PresoBatchApprovalReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "审批成功"}"
// @Router /api/v1/oc/mall/preso/batch-approval [post]
// @Security Bearer
func (e Preso) BatchApproval(c *gin.Context) {
	req := dto.PresoBatchApprovalReq{}
	s := service.Preso{}
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

	err, errMsgList := s.BatchApproval(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("审批失败，\r\n失败信息 %s", err.Error()))
		return
	} else if len(errMsgList) > 0 {
		err = errors.New(strings.Join(errMsgList, "\r\n"))
		e.OK("error", "部分审批单审批失败："+err.Error())
	} else {
		e.OK("success", "审批成功")
	}
}

// Withdraw 撤回审批单
// @Summary 撤回审批单
// @Description 撤回审批单
// @Tags 领用审批
// @Accept application/json
// @Product application/json
// @Param data body dto.PresoWithdrawReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "撤回成功"}"
// @Router /api/v1/oc/mall/preso/withdraw [post]
// @Security Bearer
func (e Preso) Withdraw(c *gin.Context) {
	req := dto.PresoWithdrawReq{}
	s := service.Preso{}
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

	err = s.Withdraw(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("撤回失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(nil, "审批成功")
}

// SaveFile 上传文档
// @Summary 上传文档
// @Description 上传文档
// @Tags 领用审批
// @Accept application/json
// @Product application/json
// @Param data body dto.PresoSaveFileReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "保存成功"}"
// @Router /api/v1/oc/mall/preso/save-file [post]
// @Security Bearer
func (e Preso) SaveFile(c *gin.Context) {
	req := dto.PresoSaveFileReq{}
	s := service.Preso{}
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

	err = s.SaveFile(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("保存失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(nil, "保存成功")
}

// Export 导出审批明细
// @Summary 导出审批明细
// @Description 导出审批明细
// @Tags 领用审批
// @Success 200 {object} response.Response{data=response.Page{list=[]dto.OrderInfoGetPageResp}} "{"code": 200, "data": [...]}"
// @Router /api/v1/oc/mall/preso/export/{id} [get]
// @Security Bearer
func (e Preso) Export(c *gin.Context) {
	req := dto.PresoGetReq{}
	s := service.Preso{}
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
	var object []dto.PresoGetExportResp

	p := actions.GetPermissionFromContext(c)
	err = s.GetExport(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("导出审批明细失败，\r\n失败信息 %s", err.Error()))
		return
	}

	var exportData []interface{}

	for _, resp := range object {
		exportData = append(exportData, resp)
	}

	title := []map[string]string{
		{"skuCode": "订货号"},
		{"productName": "商品名称"},
		{"vendorName": "货主"},
		{"supplierSkuCode": "货主SKU号"},
		{"userProductRemark": "SKU备注"},
		{"quantity": "数量"},
		{"unit": "单位"},
		{"nakedUnitPrice": "未税单价"},
		{"salePrice": "含税单价"},
		{"tax": "税额"},
		{"untaxedTotal": "未税总计"},
		{"taxedTotal": "含税总计"},
		{"approveRemark": "备注"},
	}
	excelApp := excel.NewExcel()
	err = excelApp.ExportExcelByStruct(c, title, exportData, "审批单-" + req.Id, "Sheet1")
	if err != nil {
		e.Error(500, err, fmt.Sprintf("订单明细导出失败： %s", err.Error()))
	}
}

// BuyAgain 再次购买
// @Summary 再次购买
// @Description 再次购买
// @Tags 领用审批
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.OrderInfoBuyAgainReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/oc/mall/preso/buy-again/{id} [put]
// @Security Bearer
func (e Preso) BuyAgain(c *gin.Context) {
	req := dto.PresoBuyAgainReq{}
	s := service.Preso{}
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

	errList, err := s.BuyAgain(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("再次购买失败，\r\n失败信息 %s", err.Error()))
		return
	} else if len(errList) > 0 {
		e.OK("warning", "部分商品再次购买失败："+ strings.Join(errList, "\r\n"))
	} else {
		e.OK("success", "再次购买成功")
	}
}

// Expire 审批单过期脚本：(15分钟一次)
func (e Preso) Expire(c *gin.Context) {
	s := service.Preso{}
	err := e.MakeContext(c).
		//MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	echoMsg, err := s.Expire()
	if err != nil {
		c.String(500, err.Error())
	} else {
		c.String(200, echoMsg)
	}
}

// CronApprove 定时任务审批提醒：(5分钟一次)
func (e Preso) CronApprove(c *gin.Context) {
	s := service.Preso{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	echoMsg, err := s.CronApprove(c)
	if err != nil {
		c.String(500, err.Error())
	} else {
		c.String(200, echoMsg)
	}
}

// DeleteFile 删除文件
// @Summary 删除文件
// @Description 删除文件
// @Tags 领用审批
// @Param data body dto.CsApplyDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/oc/mall/preso/file [delete]
// @Security Bearer
func (e Preso) DeleteFile(c *gin.Context) {
	req := dto.PresoDeleteFleReq{}
	s := service.Preso{}
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

	err = s.DeleteFile(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("删除审批单附件失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}





//// Insert 创建预订单
//// @Summary 创建预订单
//// @Description 创建预订单
//// @Tags 预订单
//// @Accept application/json
//// @Product application/json
//// @Param data body dto.PresoInsertReq true "data"
//// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
//// @Router /api/v1/oc/mall/preso [post]
//// @Security Bearer
//func (e Preso) Insert(c *gin.Context) {
//    req := dto.PresoInsertReq{}
//    s := service.Preso{}
//    err := e.MakeContext(c).
//        MakeOrm().
//        Bind(&req).
//        MakeService(&s.Service).
//        Errors
//    if err != nil {
//        e.Logger.Error(err)
//        e.Error(500, err, err.Error())
//        return
//    }
//	// 设置创建人
//	req.SetCreateBy(user.GetUserId(c))
//    req.SetCreateByName(user.GetUserName(c))
//	err = s.Insert(&req)
//	if err != nil {
//		e.Error(500, err, fmt.Sprintf("创建预订单失败，\r\n失败信息 %s", err.Error()))
//        return
//	}
//
//	e.OK(req.GetId(), "创建成功")
//}
//
//// Update 修改预订单
//// @Summary 修改预订单
//// @Description 修改预订单
//// @Tags 预订单
//// @Accept application/json
//// @Product application/json
//// @Param id path int true "id"
//// @Param data body dto.PresoUpdateReq true "body"
//// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
//// @Router /api/v1/oc/mall/preso/{id} [put]
//// @Security Bearer
//func (e Preso) Update(c *gin.Context) {
//    req := dto.PresoUpdateReq{}
//    s := service.Preso{}
//    err := e.MakeContext(c).
//        MakeOrm().
//        Bind(&req).
//        MakeService(&s.Service).
//        Errors
//    if err != nil {
//        e.Logger.Error(err)
//        e.Error(500, err, err.Error())
//        return
//    }
//	req.SetUpdateBy(user.GetUserId(c))
//	p := actions.GetPermissionFromContext(c)
//
//	err = s.Update(&req, p)
//	if err != nil {
//		e.Error(500, err, fmt.Sprintf("修改预订单失败，\r\n失败信息 %s", err.Error()))
//        return
//	}
//	e.OK( req.GetId(), "修改成功")
//}
//
//// Delete 删除预订单
//// @Summary 删除预订单
//// @Description 删除预订单
//// @Tags 预订单
//// @Param data body dto.PresoDeleteReq true "body"
//// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
//// @Router /api/v1/oc/mall/preso [delete]
//// @Security Bearer
//func (e Preso) Delete(c *gin.Context) {
//    s := service.Preso{}
//    req := dto.PresoDeleteReq{}
//    err := e.MakeContext(c).
//        MakeOrm().
//        Bind(&req).
//        MakeService(&s.Service).
//        Errors
//    if err != nil {
//        e.Logger.Error(err)
//        e.Error(500, err, err.Error())
//        return
//    }
//
//	// req.SetUpdateBy(user.GetUserId(c))
//	p := actions.GetPermissionFromContext(c)
//
//	err = s.Remove(&req, p)
//	if err != nil {
//		e.Error(500, err, fmt.Sprintf("删除预订单失败，\r\n失败信息 %s", err.Error()))
//        return
//	}
//	e.OK( req.GetId(), "删除成功")
//}
