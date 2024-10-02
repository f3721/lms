package dto

import (
	"go-admin/app/oc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
	"strconv"
)

type CsApplyGetPageReq struct {
	dto.Pagination     `search:"-"`
	CsNo               string `form:"csNo"  search:"type:exact;column:cs_no;table:cs_apply" comment:"售后申请编号"`
	CsType             string `form:"csType"  search:"type:exact;column:cs_type;table:cs_apply" comment:"售后类型：0-退货、1-换货、2-退款、3-发票问题、4-技术及资料支持、5-技术及资料支持、6-缺货少配件、7-售后维修、8-其他、9-订单取消、10-补发货（仅限erp迁移老数据使用）"`
	OrderId            string `form:"orderId"  search:"type:exact;column:order_id;table:cs_apply" comment:"领用订单号"`
	CsStatus           string `form:"csStatus"  search:"type:exact;column:cs_status;table:cs_apply" comment:"售后状态：0-待处理、1-已确认、2-已驳回、99-完结"`
	UserId             int    `form:"userId"  search:"type:exact;column:user_id;table:cs_apply" comment:"提交人id"`
	WarehouseCode      string `form:"warehouseCode"  search:"type:exact;column:warehouse_code;table:cs_apply" comment:"退货实体仓库code"`
	LogicWarehouseCode string `form:"logicWarehouseCode"  search:"type:exact;column:logic_warehouse_code;table:cs_apply" comment:"退货逻辑仓code"`
	CsSource           string `form:"csSource"  search:"type:exact;column:cs_source;table:cs_apply" comment:"mall ,sxyz"`
	VendorId           int    `form:"vendorId"  search:"type:exact;column:vendor_id;table:cs_apply" comment:"售后单所属货主id"`
	VendorSkuCode      string `form:"vendorSkuCode"  search:"type:exact;column:vendor_sku_code;table:cs_apply" comment:"售后单所属货主sku"`
	IsStatements       string `form:"isStatements"  search:"type:exact;column:is_statements;table:cs_apply" comment:"订单是否存在对账单 0否 1是"`
	CompanyId          int    `form:"companyId" search:"type:in;column:user_company_id;table:o"`
	CreateAtStart      string `form:"createAtStart"  search:"type:gte;column:created_at;table:cs_apply" comment:"订单是否存在对账单 0否 1是"`
	CreateAtEnd        string `form:"createAtEnd"  search:"type:lte;column:created_at;table:cs_apply" comment:"订单是否存在对账单 0否 1是"`
	FilterSkuCode      string `form:"skuCode" search:"-"`
	FilterProductNo    string `form:"productNo" search:"-"`
	FilterProductName  string `form:"productName" search:"-"`
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

type CsApplyData struct {
	models.CsApply
	CompanyName  string `json:"companyName"`  // 公司名
	TotalAmount  string `json:"totalAmount"`  // 订单支付金额
	CustomerName string `json:"customerName"` // 客户名称
}

type CsApplyListData struct {
	CsApplyData
}

type CsApplyInfoData struct {
	CsApplyData
	CsStatusText string `json:"csStatusText" gorm:"-"`
	CsTypeText   string `json:"csTypeText" gorm:"-"`
}

type CsApplyInsertReq struct {
	Id                 int    `json:"-" comment:""` //
	CsNo               string `json:"csNo" comment:"售后申请编号"`
	CsType             int    `json:"csType" comment:"售后类型：0-退货、1-换货、2-退款、3-发票问题、4-技术及资料支持、5-技术及资料支持、6-缺货少配件、7-售后维修、8-其他、9-订单取消、10-补发货（仅限erp迁移老数据使用）"`
	CsDescription      string `json:"csDescription" comment:"售后申请描述"`
	OrderId            string `json:"orderId" comment:"销售订单号"`
	Pics               string `json:"pics" comment:"售后申请图片"`
	UserId             int    `json:"userId" comment:"提交人id"`
	UserName           string `json:"userName" comment:"提交人名称"`
	RefundAmt          string `json:"refundAmt" comment:"退款金额"`
	ReparationAmt      string `json:"reparationAmt" comment:"赔款金额"`
	TransferAmount     string `json:"transferAmount" comment:"转款金额"`
	CsReason           int    `json:"csReason" comment:"售后原因id"`
	CsIssueDetail      string `json:"csIssueDetail" comment:"产品质量问题投诉必填信息"`
	WarehouseCode      string `json:"warehouseCode" comment:"退货实体仓库code"`
	WarehouseName      string `json:"warehouseName" comment:"退货实体仓库name"`
	LogicWarehouseCode string `json:"logicWarehouseCode" comment:"退货逻辑仓code"`
	CsSource           string `json:"csSource" comment:"mall ,sxyz"`
	VendorId           int    `json:"vendorId" comment:"售后单所属货主id"`
	VendorName         string `json:"vendorName" comment:"售后单所属货主名"`
	VendorSkuCode      string `json:"vendorSkuCode" comment:"售后单所属货主sku"`
	AuditReason        string `json:"auditReason" comment:"审核原因"`
	IsStatements       int    `json:"isStatements" comment:"订单是否存在对账单 0否 1是"`
	ApplyPrice         string `json:"applyPrice" comment:"申请金额"`
	ApplyQuantity      int    `json:"applyQuantity" comment:"申请总数量"`
	common.ControlBy   `json:"-"`
}

func (s *CsApplyInsertReq) Generate(model *models.CsApply) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.CsNo = s.CsNo
	model.CsType = s.CsType
	model.CsDescription = s.CsDescription
	model.OrderId = s.OrderId
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
	// 尝试将 s.ApplyPrice 转换为 float64 类型
	applyPrice, _ := strconv.ParseFloat(s.ApplyPrice, 64)
	model.ApplyPrice = applyPrice
	model.ApplyQuantity = s.ApplyQuantity
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
}

func (s *CsApplyInsertReq) GetId() interface{} {
	return s.Id
}

type CsApplyUpdateReq struct {
	Id                 int    `uri:"id" comment:""` //
	CsNo               string `json:"csNo" comment:"售后申请编号"`
	CsType             int    `json:"csType" comment:"售后类型：0-退货、1-换货、2-退款、3-发票问题、4-技术及资料支持、5-技术及资料支持、6-缺货少配件、7-售后维修、8-其他、9-订单取消、10-补发货（仅限erp迁移老数据使用）"`
	CsDescription      string `json:"csDescription" comment:"售后申请描述"`
	OrderId            string `json:"orderId" comment:"销售订单号"`
	CsStatus           int    `json:"csStatus" comment:"售后状态：0-待处理、1-已确认、2-已驳回、99-完结"`
	Telephone          string `json:"telephone" comment:"联系电话"`
	Pics               string `json:"pics" comment:"售后申请图片"`
	UserId             int    `json:"userId" comment:"提交人id"`
	UserName           string `json:"userName" comment:"提交人名称"`
	RefundAmt          string `json:"refundAmt" comment:"退款金额"`
	ReparationAmt      string `json:"reparationAmt" comment:"赔款金额"`
	TransferAmount     string `json:"transferAmount" comment:"转款金额"`
	CsReason           int    `json:"csReason" comment:"售后原因id"`
	CsIssueDetail      string `json:"csIssueDetail" comment:"产品质量问题投诉必填信息"`
	WarehouseCode      string `json:"warehouseCode" comment:"退货实体仓库code"`
	WarehouseName      string `json:"warehouseName" comment:"退货实体仓库name"`
	LogicWarehouseCode string `json:"logicWarehouseCode" comment:"退货逻辑仓code"`
	CsSource           string `json:"csSource" comment:"mall ,sxyz"`
	VendorId           int    `json:"vendorId" comment:"售后单所属货主id"`
	VendorName         string `json:"vendorName" comment:"售后单所属货主名"`
	VendorSkuCode      string `json:"vendorSkuCode" comment:"售后单所属货主sku"`
	AuditReason        string `json:"auditReason" comment:"审核原因"`
	IsStatements       int    `json:"isStatements" comment:"订单是否存在对账单 0否 1是"`
	ApplyPrice         string `json:"applyPrice" comment:"申请金额"`
	ApplyQuantity      int    `json:"applyQuantity" comment:"申请总数量"`
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
	// 尝试将 s.ApplyPrice 转换为 float64 类型
	applyPrice, _ := strconv.ParseFloat(s.ApplyPrice, 64)
	model.ApplyPrice = applyPrice
	model.ApplyQuantity = s.ApplyQuantity
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
}

func (s *CsApplyUpdateReq) GetId() interface{} {
	return s.Id
}

// CsApplyGetReq 功能获取请求参数
type CsApplyGetReq struct {
	CsNo string `uri:"csNo"`
}

// CsApplyGetRes 功能获取返回参数
type CsApplyGetRes struct {
	models.CsApply
	ApplyDetailList  []*models.CsApplyDetail //申请售后列表
	ActualDetailList []*models.CsApplyDetail //实际售后列表
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
}

type CsApplyConfirmReq struct {
	CsNo string `json:"csNo"` // 售后单号
}

type CsApplyAllAuditReq struct {
	OperationType string   `json:"operationType"` // 售后审核类型  undo作废 confirm确认
	CsNos         []string `json:"csNos"`         //售后单号
	AuditReason   string   `json:"auditReason"`   // 审核原因
}

type CsApplyAllAuditRes struct {
	Msg     string `json:"msg"`     // 操作后提示
	MsgType string `json:"msgType"` // 提示弹窗的类型 success 或者 warning
}

type ConfirmProductStatus struct {
	IsAll bool
	IsEnd bool
}

// 售后单状态
var CsApplyCsStatus = map[int]string{
	0:  "待审核",
	1:  "已确认",
	2:  "已拒绝",
	3:  "已关闭",
	10: "待退货入库",
	11: "待取消出库",
	99: "已完结",
}

// CsApplyStatus 可以进行售后状态修改的类型0, 10, 11
var CsApplyStatusUpdateOk = map[int]int{
	0:  0,
	1:  1,
	10: 10,
	11: 11,
}

// CsApplyStatus 可以进行售后状态修改的类型0, 10, 11
var CsApplyStatusUpdateOkList = []int{
	0,
	1,
	10,
	11,
}

// AfterSalesOrderLogType 订单表显示的售后状态只有这两个
var AfterSalesOrderLogType = map[int]string{
	1:  "after_sales_apply",
	99: "after_sales_complete",
}

// CsApplyGetReq 功能获取请求参数
type CsApplyGetSaleProductsReq struct {
	OrderId string `uri:"orderId" vd:"len($)>0" `
	CsType  int    `form:"csType"`
}

type CsApplyGetSaleProductsRes struct {
	Products     []*CsApplyGetOrderProducts `json:"products"`
	Warehouse    *CsApplyWarehouseData      `json:"warehouse"`
	IsStatements bool                       `json:"isStatements"`
}

type CsApplyWarehouseData struct {
	WarehouseCode string `json:"warehouseCode"`
	WarehouseName string `json:"warehouseName"`
}

type CsApplyGetOrderProducts struct {
	models.OrderDetail
	// 添加其他字段...
	VendorsName string `json:"vendorsName" gorm:"-"` // 货主名
	VendorsCode string `json:"vendorsCode" gorm:"-"` // 货主code
	AreaId      int    `json:"areaId" gorm:"-"`
	OrderStatus string `json:"orderStatus" gorm:"-"` // 订单状态

	ActualStock   int         `json:"actualStock" gorm:"-"`
	SkuProfit     float64     `json:"skuProfit" gorm:"-"`
	SkuProfitText string      `json:"skuProfitText" gorm:"-"`
	AfterTypeNum  map[int]int `json:"afterNum" gorm:"-"`
	AfterNum      int         `json:"afterNum" gorm:"-"`
	ReturnNum     int         `json:"returnNum" gorm:"-"`
	ListSort      int         `json:"listSort" gorm:"-"`
	AllowQuantity int         `json:"allowQuantity" gorm:"-"` // 允许申请数量
	Quantity      int         `json:"quantity" gorm:"-"`
}

type CsApplyInsertRequest struct {
	CsType        int                     `json:"csType"`        // CS类型
	OrderId       string                  `json:"orderId"`       // 订单ID
	WarehouseCode string                  `json:"warehouseCode"` // 仓库代码
	IsStatements  int                     `json:"isStatements"`  // 是否结算
	Products      []*CsApplyInsertProduct `json:"products"`

	CsNo               string  `json:"-"`
	SupplierInfoID     string  `json:"-"`
	SupplierName       string  `json:"-"`
	SupplierSKUCode    string  `json:"-"`
	ApplyPrice         float64 `json:"-"`
	ApplyQuantity      int     `json:"-"`
	VendorsId          int     `json:"-"` // 货主名称
	VendorsName        string  `json:"-"` // 货主名称
	VendorsSkuCode     string  `json:"-"` // 货主SKU编码
	WarehouseName      string  `json:"-"`
	LogicWarehouseCode string  `json:"-"`

	common.ControlBy `json:"-"`
}

type CsApplyInsertProduct struct {
	GoodsId          int     `json:"goodsId"`          // 商品ID
	SkuCode          string  `json:"skuCode"`          // SKU编码
	Quantity         int     `json:"quantity"`         // 数量
	BrandName        string  `json:"brandName"`        // 品牌名称
	ProductName      string  `json:"productName"`      // 产品名称
	Unit             string  `json:"unit"`             // 单位
	ProductPic       string  `json:"productPic"`       // 产品图片
	ProductModel     string  `json:"productModel"`     // 产品型号
	SalePrice        float64 `json:"salePrice"`        // 销售价格
	RefundAmt        float64 `json:"refundAmt"`        // 退款金额
	ReparationAmt    float64 `json:"reparationAmt"`    // 赔偿金额
	WarehouseCode    string  `json:"warehouseCode"`    // 仓库编码
	Worthless        int     `json:"worthless"`        // 无效标志
	TransferAmount   float64 `json:"transferAmount"`   // 转账金额
	ReturnToSupplier string  `json:"returnToSupplier"` // 退回供应商
	CSType           int     `json:"csType"`           // 客服类型
	VendorId         int     `json:"vendorId"`         // 货主名称
	VendorName       string  `json:"vendorsName"`      // 货主名称
	VendorSkuCode    string  `json:"vendorsCode"`      // 货主SKU编码
	ProductNo        string  `json:"productNo"`        // 产品编号
	Remark           string  `json:"remark"`           // 备注
	IsDefective      int     `json:"isDefective"`      // 是否有质量问题
	Pics             string  `json:"pics"`             // 图片列表

	CsNo string `json:"-"`
}

type CallStockChangeEventReq struct {
	OrderId            string
	CsNo               string
	WarehouseCode      string
	LogicWarehouseCode string
	ProductList        interface{}
	VendorsId          int
}

type GetOrderInfoByCsReq struct {
	CsNo string `uri:"csNo" vd:"len($)>0" `
}

type CsApplyIsOrderInAfterPendingReview struct {
	OrderId string `uri:"orderId" vd:"len($)>0" `
}

type CsApplyIsOrderInAfterPendingReviewRes struct {
	IsPendingReview bool `json:"isPendingReview"`
}
