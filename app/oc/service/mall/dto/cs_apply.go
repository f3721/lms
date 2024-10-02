package dto

import (
	"go-admin/app/oc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type CsApplyGetPageReq struct {
	dto.Pagination     `search:"-"`
	CsNo               string `form:"csNo"  search:"type:exact;column:cs_no;table:cs_apply" comment:"售后申请编号"`                                                                                             //售后申请编号
	CsType             int    `form:"csType"  search:"type:exact;column:cs_type;table:cs_apply" comment:"售后类型：0-退货、1-换货、2-退款、3-发票问题、4-技术及资料支持、5-技术及资料支持、6-缺货少配件、7-售后维修、8-其他、9-订单取消、10-补发货（仅限erp迁移老数据使用）"` //售后类型：0-退货、1-换货、2-退款、3-发票问题、4-技术及资料支持、5-技术及资料支持、6-缺货少配件、7-售后维修、8-其他、9-订单取消、10-补发货（仅限erp迁移老数据使用）
	OrderId            string `form:"orderId"  search:"type:exact;column:order_id;table:cs_apply" comment:"销售订单号"`                                                                                        //销售订单号
	CsStatus           string `form:"csStatus"  search:"-" comment:"售后状态：0-待处理、1-已确认、2-已驳回、99-完结"`                                                                                                        //提交人id
	WarehouseCode      string `form:"warehouseCode"  search:"type:exact;column:warehouse_code;table:cs_apply" comment:"退货实体仓库code"`                                                                       //退货实体仓库code
	LogicWarehouseCode string `form:"logicWarehouseCode"  search:"type:exact;column:logic_warehouse_code;table:cs_apply" comment:"退货逻辑仓code"`                                                             //退货逻辑仓code
	CsSource           string `form:"csSource"  search:"type:exact;column:cs_source;table:cs_apply" comment:"mall ,sxyz"`                                                                                 //mall ,sxyz
	VendorId           int    `form:"vendorId"  search:"type:exact;column:vendor_id;table:cs_apply" comment:"售后单所属货主id"`                                                                                  //售后单所属货主id
	VendorSkuCode      string `form:"vendorSkuCode"  search:"type:exact;column:vendor_sku_code;table:cs_apply" comment:"售后单所属货主sku"`                                                                      //售后单所属货主sku
	IsStatements       int    `form:"isStatements"  search:"type:exact;column:is_statements;table:cs_apply" comment:"订单是否存在对账单 0否 1是"`                                                                    //订单是否存在对账单 0否 1是
	FilterKeyword      string `form:"filterKeyword" search:"-"`
	UserId             int    `form:"userId"  search:"-" comment:""`
	CsApplyOrder
}

type CsApplyOrder struct {
	Id                 string `form:"idOrder"  search:"type:order;column:id;table:cs_apply"`
	CsNo               string `form:"csNoOrder"  search:"type:order;column:cs_no;table:cs_apply"`
	CsType             string `form:"csTypeOrder"  search:"type:order;column:cs_type;table:cs_apply"`
	CsDescription      string `form:"csDescriptionOrder"  search:"type:order;column:cs_description;table:cs_apply"`
	OrderId            string `form:"orderIdOrder"  search:"type:order;column:order_id;table:cs_apply"`
	CsStatus           string `form:"csStatusOrder"  search:"type:order;column:cs_status;table:cs_apply"`
	CreatedAt          string `form:"createdAtOrder"  search:"type:order;column:created_at;table:cs_apply"`
	UpdatedAt          string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:cs_apply"`
	Telephone          string `form:"telephoneOrder"  search:"type:order;column:telephone;table:cs_apply"`
	Pics               string `form:"picsOrder"  search:"type:order;column:pics;table:cs_apply"`
	UserId             string `form:"userIdOrder"  search:"type:order;column:user_id;table:cs_apply"`
	UserName           string `form:"userNameOrder"  search:"type:order;column:user_name;table:cs_apply"`
	RefundAmt          string `form:"refundAmtOrder"  search:"type:order;column:refund_amt;table:cs_apply"`
	ReparationAmt      string `form:"reparationAmtOrder"  search:"type:order;column:reparation_amt;table:cs_apply"`
	TransferAmount     string `form:"transferAmountOrder"  search:"type:order;column:transfer_amount;table:cs_apply"`
	CsReason           string `form:"csReasonOrder"  search:"type:order;column:cs_reason;table:cs_apply"`
	CsIssueDetail      string `form:"csIssueDetailOrder"  search:"type:order;column:cs_issue_detail;table:cs_apply"`
	WarehouseCode      string `form:"warehouseCodeOrder"  search:"type:order;column:warehouse_code;table:cs_apply"`
	WarehouseName      string `form:"warehouseNameOrder"  search:"type:order;column:warehouse_name;table:cs_apply"`
	LogicWarehouseCode string `form:"logicWarehouseCodeOrder"  search:"type:order;column:logic_warehouse_code;table:cs_apply"`
	CsSource           string `form:"csSourceOrder"  search:"type:order;column:cs_source;table:cs_apply"`
	VendorId           string `form:"vendorIdOrder"  search:"type:order;column:vendor_id;table:cs_apply"`
	VendorName         string `form:"vendorNameOrder"  search:"type:order;column:vendor_name;table:cs_apply"`
	VendorSkuCode      string `form:"vendorSkuCodeOrder"  search:"type:order;column:vendor_sku_code;table:cs_apply"`
	AuditReason        string `form:"auditReasonOrder"  search:"type:order;column:audit_reason;table:cs_apply"`
	IsStatements       string `form:"isStatementsOrder"  search:"type:order;column:is_statements;table:cs_apply"`
	ApplyPrice         string `form:"applyPriceOrder"  search:"type:order;column:apply_price;table:cs_apply"`
	ApplyQuantity      string `form:"applyQuantityOrder"  search:"type:order;column:apply_quantity;table:cs_apply"`
	CreateBy           string `form:"createByOrder"  search:"type:order;column:create_by;table:cs_apply"`
	UpdateBy           string `form:"updateByOrder"  search:"type:order;column:update_by;table:cs_apply"`
}

func (m *CsApplyGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type CsApplyGetPageGroupByOrder struct {
	OrderId    string                               `json:"orderId"`
	CreatedAt  string                               `json:"createdAt"`
	ContractNo string                               `json:"contractNo"`
	CsApply    []*CsApplyGetPageGroupByOrderCsApply `json:"csApply" gorm:"foreignKey:OrderId;references:OrderId"`
}

type CsApplyGetPageGroupByOrderCsApply struct {
	models.CsApply
	CsApplyDetail []*CsApplyGetPageGroupByOrderCsApplyDetail `json:"csApplyDetail" gorm:"foreignKey:CsNo;references:CsNo"`
}

type CsApplyGetPageGroupByOrderCsApplyDetail struct {
	models.CsApplyDetail
	Moq           int `json:"moq"`
	ApplyQuantity int `json:"applyQuantity"`
	FinalQuantity int `json:"finalQuantity"`
}

type CsApplyInsertReq struct {
	Id                 int     `json:"-" comment:""`                                                                                                    //
	CsNo               string  `json:"csNo" comment:"售后申请编号"`                                                                                           // 售后申请编号
	CsType             int     `json:"csType" comment:"售后类型：0-退货、1-换货、2-退款、3-发票问题、4-技术及资料支持、5-技术及资料支持、6-缺货少配件、7-售后维修、8-其他、9-订单取消、10-补发货（仅限erp迁移老数据使用）"` // 售后类型：0-退货、1-换货、2-退款、3-发票问题、4-技术及资料支持、5-技术及资料支持、6-缺货少配件、7-售后维修、8-其他、9-订单取消、10-补发货（仅限erp迁移老数据使用）
	CsDescription      string  `json:"csDescription" comment:"售后申请描述"`                                                                                  // 售后申请描述
	OrderId            string  `json:"orderId" comment:"销售订单号"`                                                                                         // 销售订单号
	CsStatus           int     `json:"csStatus" comment:"售后状态：0-待处理、1-已确认、2-已驳回、99-完结"`                                                                 // 售后状态：0-待处理、1-已确认、2-已驳回、99-完结
	Telephone          string  `json:"telephone" comment:"联系电话"`                                                                                        // 联系电话
	Pics               string  `json:"pics" comment:"售后申请图片"`                                                                                           // 售后申请图片
	UserId             int     `json:"userId" comment:"提交人id"`                                                                                          // 提交人id
	UserName           string  `json:"userName" comment:"提交人名称"`                                                                                        // 提交人名称
	RefundAmt          string  `json:"refundAmt" comment:"退款金额"`                                                                                        // 退款金额
	ReparationAmt      string  `json:"reparationAmt" comment:"赔款金额"`                                                                                    // 赔款金额
	TransferAmount     string  `json:"transferAmount" comment:"转款金额"`                                                                                   // 转款金额
	CsReason           int     `json:"csReason" comment:"售后原因id"`                                                                                       // 售后原因id
	CsIssueDetail      string  `json:"csIssueDetail" comment:"产品质量问题投诉必填信息"`                                                                            // 产品质量问题投诉必填信息
	WarehouseCode      string  `json:"warehouseCode" comment:"退货实体仓库code"`                                                                              // 退货实体仓库code
	WarehouseName      string  `json:"warehouseName" comment:"退货实体仓库name"`                                                                              // 退货实体仓库name
	LogicWarehouseCode string  `json:"logicWarehouseCode" comment:"退货逻辑仓code"`                                                                          // 退货逻辑仓code
	CsSource           string  `json:"csSource" comment:"mall ,sxyz"`                                                                                   // mall ,sxyz
	VendorId           int     `json:"vendorId" comment:"售后单所属货主id"`                                                                                    // 售后单所属货主id
	VendorName         string  `json:"vendorName" comment:"售后单所属货主名"`                                                                                   // 售后单所属货主名
	VendorSkuCode      string  `json:"vendorSkuCode" comment:"售后单所属货主sku"`                                                                              // 售后单所属货主sku
	AuditReason        string  `json:"auditReason" comment:"审核原因"`                                                                                      // 审核原因
	IsStatements       int     `json:"isStatements" comment:"订单是否存在对账单 0否 1是"`                                                                          // 订单是否存在对账单 0否 1是
	ApplyPrice         float64 `json:"applyPrice" comment:"申请金额"`                                                                                       // 申请金额
	ApplyQuantity      int     `json:"applyQuantity" comment:"申请总数量"`                                                                                   // 申请总数量
	common.ControlBy
}

func (s *CsApplyInsertReq) Generate(model *models.CsApply) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.CsNo = s.CsNo
	model.CsType = s.CsType
	model.CsDescription = s.CsDescription
	model.OrderId = s.OrderId
	model.CsStatus = s.CsStatus
	model.Telephone = s.Telephone
	model.Pics = s.Pics
	model.UserId = s.UserId
	model.UserName = s.UserName
	model.RefundAmt = s.RefundAmt
	model.ReparationAmt = s.ReparationAmt
	model.TransferAmount = s.TransferAmount
	model.CsReason = s.CsReason
	model.CsIssueDetail = s.CsIssueDetail
	model.WarehouseCode = s.WarehouseCode
	model.WarehouseName = s.WarehouseName
	model.LogicWarehouseCode = s.LogicWarehouseCode
	model.CsSource = s.CsSource
	model.VendorId = s.VendorId
	model.VendorName = s.VendorName
	model.VendorSkuCode = s.VendorSkuCode
	model.AuditReason = s.AuditReason
	model.IsStatements = s.IsStatements
	model.ApplyPrice = s.ApplyPrice
	model.ApplyQuantity = s.ApplyQuantity
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
}

func (s *CsApplyInsertReq) GetId() interface{} {
	return s.Id
}

type CsApplyUpdateReq struct {
	Id                 int     `uri:"id" comment:""`                                                                                                    //
	CsNo               string  `json:"csNo" comment:"售后申请编号"`                                                                                           // 售后申请编号
	CsType             int     `json:"csType" comment:"售后类型：0-退货、1-换货、2-退款、3-发票问题、4-技术及资料支持、5-技术及资料支持、6-缺货少配件、7-售后维修、8-其他、9-订单取消、10-补发货（仅限erp迁移老数据使用）"` // 售后类型：0-退货、1-换货、2-退款、3-发票问题、4-技术及资料支持、5-技术及资料支持、6-缺货少配件、7-售后维修、8-其他、9-订单取消、10-补发货（仅限erp迁移老数据使用）
	CsDescription      string  `json:"csDescription" comment:"售后申请描述"`                                                                                  // 售后申请描述
	OrderId            string  `json:"orderId" comment:"销售订单号"`                                                                                         // 销售订单号
	CsStatus           int     `json:"csStatus" comment:"售后状态：0-待处理、1-已确认、2-已驳回、99-完结"`                                                                 // 售后状态：0-待处理、1-已确认、2-已驳回、99-完结
	Telephone          string  `json:"telephone" comment:"联系电话"`                                                                                        // 联系电话
	Pics               string  `json:"pics" comment:"售后申请图片"`                                                                                           // 售后申请图片
	UserId             int     `json:"userId" comment:"提交人id"`                                                                                          // 提交人id
	UserName           string  `json:"userName" comment:"提交人名称"`                                                                                        // 提交人名称
	RefundAmt          string  `json:"refundAmt" comment:"退款金额"`                                                                                        // 退款金额
	ReparationAmt      string  `json:"reparationAmt" comment:"赔款金额"`                                                                                    // 赔款金额
	TransferAmount     string  `json:"transferAmount" comment:"转款金额"`                                                                                   // 转款金额
	CsReason           int     `json:"csReason" comment:"售后原因id"`                                                                                       // 售后原因id
	CsIssueDetail      string  `json:"csIssueDetail" comment:"产品质量问题投诉必填信息"`                                                                            // 产品质量问题投诉必填信息
	WarehouseCode      string  `json:"warehouseCode" comment:"退货实体仓库code"`                                                                              // 退货实体仓库code
	WarehouseName      string  `json:"warehouseName" comment:"退货实体仓库name"`                                                                              // 退货实体仓库name
	LogicWarehouseCode string  `json:"logicWarehouseCode" comment:"退货逻辑仓code"`                                                                          // 退货逻辑仓code
	CsSource           string  `json:"csSource" comment:"mall ,sxyz"`                                                                                   // mall ,sxyz
	VendorId           int     `json:"vendorId" comment:"售后单所属货主id"`                                                                                    // 售后单所属货主id
	VendorName         string  `json:"vendorName" comment:"售后单所属货主名"`                                                                                   // 售后单所属货主名
	VendorSkuCode      string  `json:"vendorSkuCode" comment:"售后单所属货主sku"`                                                                              // 售后单所属货主sku
	AuditReason        string  `json:"auditReason" comment:"审核原因"`                                                                                      // 审核原因
	IsStatements       int     `json:"isStatements" comment:"订单是否存在对账单 0否 1是"`                                                                          // 订单是否存在对账单 0否 1是
	ApplyPrice         float64 `json:"applyPrice" comment:"申请金额"`                                                                                       // 申请金额
	ApplyQuantity      int     `json:"applyQuantity" comment:"申请总数量"`                                                                                   // 申请总数量
	common.ControlBy
}

func (s *CsApplyUpdateReq) Generate(model *models.CsApply) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.CsNo = s.CsNo
	model.CsType = s.CsType
	model.CsDescription = s.CsDescription
	model.OrderId = s.OrderId
	model.CsStatus = s.CsStatus
	model.Telephone = s.Telephone
	model.Pics = s.Pics
	model.UserId = s.UserId
	model.UserName = s.UserName
	model.RefundAmt = s.RefundAmt
	model.ReparationAmt = s.ReparationAmt
	model.TransferAmount = s.TransferAmount
	model.CsReason = s.CsReason
	model.CsIssueDetail = s.CsIssueDetail
	model.WarehouseCode = s.WarehouseCode
	model.WarehouseName = s.WarehouseName
	model.LogicWarehouseCode = s.LogicWarehouseCode
	model.CsSource = s.CsSource
	model.VendorId = s.VendorId
	model.VendorName = s.VendorName
	model.VendorSkuCode = s.VendorSkuCode
	model.AuditReason = s.AuditReason
	model.IsStatements = s.IsStatements
	model.ApplyPrice = s.ApplyPrice
	model.ApplyQuantity = s.ApplyQuantity
	model.UpdateBy = s.UpdateBy
}

func (s *CsApplyUpdateReq) GetId() interface{} {
	return s.Id
}

// CsApplyGetReq 功能获取请求参数
type CsApplyGetReq struct {
	Id int `uri:"id"`
}

func (s *CsApplyGetReq) GetId() interface{} {
	return s.Id
}

// CsApplyDeleteReq 功能删除请求参数
type CsApplyDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *CsApplyDeleteReq) GetId() interface{} {
	return s.Ids
}

type CsApplyCancelReq struct {
	CsNo        string `json:"csNo" `                                   // 售后单号
	AuditReason string `json:"auditReason" vd:"len($)>0;msg:'售后原因必填'" ` // 审核原因
	UserId      int    `json:"-"`
}

type CsApplyGoodsQuantity struct {
	CsNo     string `json:"csNo"`
	GoodsId  int    `json:"goodsId"`
	Quantity int    `json:"quantity"`
}
