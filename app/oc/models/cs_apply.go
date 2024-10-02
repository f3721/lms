package models

import (
	"encoding/json"
	"errors"
	"fmt"
	modelsUc "go-admin/app/uc/models"
	"go-admin/common/global"
	"go-admin/common/models"
	"math/rand"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	"gorm.io/gorm"
)

// 售后类型
var CsApplyTypeText = map[int]string{
	0: "退货",
	9: "取消",
}

// 售后状态
var CsApplyStatusText = map[int]string{
	0:  "待审核",
	1:  "已确认",
	2:  "已拒绝",
	3:  "已关闭",
	10: "待退货入库",
	11: "待取消出库",
	99: "已完结",
}

// 订单是否存在对账单 0否 1是
var CsApplyStatementsText = map[int]string{
	0: "否",
	1: "是",
}

// 售后来源
var CsApplySource = map[string]string{
	"mall": "mall",
	"sxyz": "sxyz",
}

var AfterSalesOrderLogType = map[int]string{
	1:  "after_sales_apply",
	99: "after_sales_complete",
}

// 售后状态：0 => '待审核', 1 => '已确认', 2 => '已作废', 3 => '已取消', 10 => '待退货入库',11 => '待取消出库',99 => '已完结'
// 待审核 待退货入库 待取消出库可以进行售后取消0, 10, 11
var CsApplyCancelOkStatus = map[int]int{
	0:  0,
	10: 10,
	11: 11,
}

var CsApplyAuditOkStatus = map[int]int{
	0: 0,
}

// CsApplyStatus 可以进行售后状态修改的类型0, 10, 11
var CsApplyStatusUpdateOkList = []int{
	0:  0,
	10: 10,
	11: 11,
}

// OrderStatusCsApplyCancelOk  订单状态中可以申请售后取消的
// 订单状态 0-未发货、11-部分发货、1-已发货、2-部分收货、3-待评价、4-已评价、5-待确认、6、缺货、7-已签收、9-已取消、10-已关闭 11 部分发货
var OrderStatusCsApplyCancelOk = map[int]int{
	0:  0,  // 0-未发货
	5:  5,  // 5-待确认
	6:  6,  // 6、缺货
	11: 11, // 11-部分发货
}

// OrderStatusCsApplyReturnOk  订单状态中可以申请售后退货的 (订单状态现在使用到的 应该只有 已发货 已签收 部分发货 待评价之类的状态现在应该没有在用  这里是以防万一)
// 订单状态 0-未发货、11-部分发货、1-已发货、2-部分收货、3-待评价、4-已评价、5-待确认、6、缺货、7-已签收、9-已取消、10-已关闭 11 部分发货
var OrderStatusCsApplyReturnOk = map[int]int{
	1:  1,  // 1-已发货
	2:  2,  // 2-部分收货
	3:  3,  // 3-待评价
	4:  4,  // 4-已评价
	7:  7,  // 7-已签收
	11: 11, // 11-部分发货
}

type CsApply struct {
	models.Model

	CsNo               string  `json:"csNo" gorm:"type:varchar(30);comment:售后申请编号"`                                                                                                // 售后申请编号
	CsType             int     `json:"csType" gorm:"type:tinyint unsigned;comment:售后类型：0-退货、1-换货、2-退款、3-发票问题、4-技术及资料支持、5-技术及资料支持、6-缺货少配件、7-售后维修、8-其他、9-订单取消、10-补发货（仅限erp迁移老数据使用）"` // 售后类型：0-退货、1-换货、2-退款、3-发票问题、4-技术及资料支持、5-技术及资料支持、6-缺货少配件、7-售后维修、8-其他、9-订单取消、10-补发货（仅限erp迁移老数据使用）
	CsDescription      string  `json:"csDescription" gorm:"type:mediumtext;comment:售后申请描述"`                                                                                        // 售后申请描述
	OrderId            string  `json:"orderId" gorm:"type:varchar(30);comment:销售订单号"`                                                                                              // 销售订单号
	CsStatus           int     `json:"csStatus" gorm:"type:tinyint(1);comment:售后状态：0-待处理、1-已确认、2-已驳回、99-完结"`                                                                       // 售后状态：0-待处理、1-已确认、2-已驳回、99-完结
	Telephone          string  `json:"telephone" gorm:"type:varchar(30);comment:联系电话"`                                                                                             // 联系电话
	Pics               string  `json:"pics" gorm:"type:mediumtext;comment:售后申请图片"`                                                                                                 // 售后申请图片
	UserId             int     `json:"userId" gorm:"type:tinyint unsigned;comment:提交人id"`                                                                                          // 提交人id
	UserName           string  `json:"userName" gorm:"type:varchar(50);comment:提交人名称"`                                                                                             // 提交人名称
	RefundAmt          string  `json:"refundAmt" gorm:"type:decimal(10,2);comment:退款金额"`                                                                                           // 退款金额
	ReparationAmt      string  `json:"reparationAmt" gorm:"type:decimal(10,2);comment:赔款金额"`                                                                                       // 赔款金额
	TransferAmount     string  `json:"transferAmount" gorm:"type:decimal(10,2);comment:转款金额"`                                                                                      // 转款金额
	CsReason           int     `json:"csReason" gorm:"type:tinyint unsigned;comment:售后原因id"`                                                                                       // 售后原因id
	CsIssueDetail      string  `json:"csIssueDetail" gorm:"type:mediumtext;comment:产品质量问题投诉必填信息"`                                                                                  // 产品质量问题投诉必填信息
	WarehouseCode      string  `json:"warehouseCode" gorm:"type:varchar(64);comment:退货实体仓库code"`                                                                                   // 退货实体仓库code
	WarehouseName      string  `json:"warehouseName" gorm:"type:varchar(100);comment:退货实体仓库name"`                                                                                  // 退货实体仓库name
	LogicWarehouseCode string  `json:"logicWarehouseCode" gorm:"type:varchar(64);comment:退货逻辑仓code"`                                                                               // 退货逻辑仓code
	CsSource           string  `json:"csSource" gorm:"type:varchar(50);comment:mall ,sxyz"`                                                                                        // mall ,sxyz
	VendorId           int     `json:"vendorId" gorm:"type:int unsigned;comment:售后单所属货主id"`                                                                                        // 售后单所属货主id
	VendorName         string  `json:"vendorName" gorm:"type:varchar(255);comment:售后单所属货主名"`                                                                                       // 售后单所属货主名
	VendorSkuCode      string  `json:"vendorSkuCode" gorm:"type:varchar(255);comment:售后单所属货主sku"`                                                                                  // 售后单所属货主sku
	AuditReason        string  `json:"auditReason" gorm:"type:mediumtext;comment:审核原因"`                                                                                            // 审核原因
	IsStatements       int     `json:"isStatements" gorm:"type:tinyint unsigned;comment:订单是否存在对账单 0否 1是"`                                                                          // 订单是否存在对账单 0否 1是
	ApplyPrice         float64 `json:"applyPrice" gorm:"type:decimal(10,2);comment:申请金额"`                                                                                          // 申请金额
	ApplyQuantity      int     `json:"applyQuantity" gorm:"type:int unsigned;comment:申请总数量"`                                                                                       // 申请总数量
	models.ModelTime
	models.ControlBy
}

func (CsApply) TableName() string {
	return "cs_apply"
}

func (e *CsApply) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *CsApply) GetId() interface{} {
	return e.Id
}

// GenerateCsNo 生成订单号
func (e *CsApply) GenerateCsNo() string {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	randomNum := rand.Intn(99999-1+1) + 1
	return fmt.Sprintf("CS%013d%05d", timestamp, randomNum)
}

type GetAfterApplyProductsByTypeBySaleIdResult struct {
	SkuCode  string
	CsType   int
	Quantity int
}

// 售后数
func (e *CsApply) GetAfterApplyProductsByTypeBySaleId(db *gorm.DB, orderID string) (result map[string]map[int]*GetAfterApplyProductsByTypeBySaleIdResult) {
	result = make(map[string]map[int]*GetAfterApplyProductsByTypeBySaleIdResult)
	resultList := []*GetAfterApplyProductsByTypeBySaleIdResult{}

	ocPrefix := global.GetTenantOcDBNameWithDB(db)

	db.Table(ocPrefix+".cs_apply_detail as ad").
		Select("ad.sku_code, a.cs_type, sum(ad.quantity) as quantity").
		Joins("INNER JOIN "+ocPrefix+".cs_apply a ON (ad.cs_no = a.cs_no)").
		Where("a.order_id = ? AND ad.cs_type = 0 AND a.cs_status = 0 ", orderID).
		Group("ad.sku_code, a.cs_type").
		Scan(&resultList)

	for _, idResult := range resultList {
		result[idResult.SkuCode] = make(map[int]*GetAfterApplyProductsByTypeBySaleIdResult)
		result[idResult.SkuCode][idResult.CsType] = idResult
	}

	return result
}

type GetAfterReturnProductsBySaleIdResult struct {
	SkuCode  string
	Quantity int
}

// 已退货数
func (e *CsApply) GetAfterReturnProductsBySaleId(db *gorm.DB, orderID string) (result map[string]*GetAfterReturnProductsBySaleIdResult) {
	result = make(map[string]*GetAfterReturnProductsBySaleIdResult)
	resultList := []*GetAfterReturnProductsBySaleIdResult{}

	ocPrefix := global.GetTenantOcDBNameWithDB(db)

	db.Table(ocPrefix+".cs_apply_detail as ad").Debug().
		Select("ad.sku_code, sum(ad.quantity) as quantity").
		Joins("INNER JOIN "+ocPrefix+".cs_apply a ON (ad.cs_no = a.cs_no)").
		Where("a.order_id = ? AND ad.cs_type = 1 AND a.cs_type = 0 AND a.cs_status = 99", orderID).
		Group("ad.sku_code").
		Scan(&resultList)
	for _, idResult := range resultList {
		result[idResult.SkuCode] = idResult
	}
	return result
}

func (e *CsApply) GetRowsByOrderId(db *gorm.DB, orderID string) (result []CsApply) {
	db.Table("cs_apply").Debug().
		Where("order_id = ?", orderID).
		Find(&result)
	return
}

// 订单售后状态改变 status 售后状态：0-无售后、1-售后处理中、2-售后已确认、99-售后处理完成
func (e *CsApply) OrderAfterSales(db *gorm.DB, orderId string, status int) error {
	var data CsApply
	// 如果有在售后中的就不用改订单状态 直到没有一个在售后中的
	err := db.Where("order_id = ?", orderId).Where("cs_status = 0").First(&data).Error
	if data.Id > 0 {
		// 如果未找到记录，则售后单不存在不用改变订单状态
		return nil
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// 如果出现其他错误，则返回错误信息
		return err
	}

	var order OrderInfo
	err = db.Where("order_id = ?", orderId).First(&order).Error
	if err != nil {
		return err
	}

	order.RmaStatus = status

	// 调用修改订单售后状态的方法
	err = data.OrderUpdate(db, orderId, nil, &order, status)
	if err != nil {
		return err
	}
	return nil
}

/**
 * 售后退货完成
 *
 * @param csNo string
 * @return void
 */
func (e *CsApply) ReturnCompleted(db *gorm.DB, csNo string) (err error) {

	data, err := e.GetByCsNo(db, csNo)
	if err != nil {
		return
	}

	if data.CsStatus != 10 {
		err = errors.New("售后单状态不可操作")
		return
	}
	applyDetailList, err := e.GetApplyDetailList(db, csNo, 0)
	if err != nil {
		return
	}
	orderId := data.OrderId

	// 确认商品是否正确 有没有购买数量等
	confirmProduct, err := e.ConfirmProduct(db, orderId, applyDetailList, data.CsType)
	if err != nil {
		return
	}

	// 售后类型：0-退货、1-换货、2-退款、3-发票问题、4-技术及资料支持、5-技术及资料支持、6-缺货少配件、7-售后维修、8-其他、9-订单取消、10-补发货（仅限ERP迁移老数据使用）
	// sxyz只有退货和取消

	productIsEnd := 0
	if confirmProduct.IsAll || confirmProduct.IsEnd {
		productIsEnd = 1
	}

	// 订单修改退货状态
	err = e.OrderReturn(db, orderId, applyDetailList, productIsEnd)
	if err != nil {
		return err
	}

	err = e.applyEnd(db, data.Id, csNo, applyDetailList, data)
	if err != nil {
		return
	}

	return

}

// Get 获取CsApply对象
func (e *CsApply) GetByCsNo(db *gorm.DB, csNo string) (*CsApply, error) {
	var data CsApply

	ocPrefix := global.GetTenantOcDBNameWithDB(db)

	err := db.Table(ocPrefix+"."+data.TableName()).
		Where("cs_no = ?", csNo).
		First(&data).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("售后单不存在")
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (e *CsApply) GetApplyDetailList(db *gorm.DB, csNo string, csType int) (csApplyDetailList []*CsApplyDetail, err error) {
	ocPrefix := global.GetTenantOcDBNameWithDB(db)
	err = db.Table(ocPrefix + "." + CsApplyDetail{}.TableName()).Where(&CsApplyDetail{
		CsType: csType,
		CsNo:   csNo,
	}).Find(&csApplyDetailList).Error
	if err != nil {
		return
	}
	return
}

type ConfirmProductStatus struct {
	IsAll      bool
	IsEnd      bool
	IsOrderEnd bool
}

// ConfirmProduct 确认商品是否正确 有没有购买数量等
func (e *CsApply) ConfirmProduct(db *gorm.DB, orderId string, shProduct []*CsApplyDetail, csType int) (*ConfirmProductStatus, error) {
	//var orderDetail []*OrderDetail{}
	ocPrefix := global.GetTenantOcDBNameWithDB(db)

	// 获取订单商品信息
	var orderDetails []*OrderDetail
	db.Table(ocPrefix+".order_detail od").Where("od.order_id = ?", orderId).Order("od.id").Find(&orderDetails)

	afterReturnProdutsQuantity := e.GetAfterReturnProductsBySaleId(db, orderId)

	// 是否全部商品退货或者取消
	isAll := true
	// 是否最后一次商品退货或者取消
	isEnd := true
	// 是否订单全部商品已经结束无法售后
	isOrderEnd := true

	okGoods := make(map[string]bool)
	for _, product := range orderDetails {
		// 超过最大库存
		quantity := product.Quantity - product.CancelQuantity
		//quantity
		productAllQuantity := product.Quantity

		// 已退货数量
		returnQuantity := 0
		if returnProdutQuantity, ok := afterReturnProdutsQuantity[product.SkuCode]; ok && returnProdutQuantity.Quantity > 0 {
			returnQuantity = returnProdutQuantity.Quantity
		}
		if csType == 0 { // 退货
			// 能申请的最大库存 = 已发货商品数量-已退货数量
			quantity = product.SendQuantity - returnQuantity

			productAllQuantity = product.SendQuantity - returnQuantity
		} else if csType == 9 { // 取消
			// 能申请的最大库存 = 原商品数量 - 已发货商品数量(存在部分发货逻辑 部分发货时可以申请取消未发货的商品 所以申请数量要扣减已发货数量) - 已取消的数量
			quantity = product.OriginalQuantity - product.SendQuantity - product.CancelQuantity

			productAllQuantity = product.OriginalQuantity - product.CancelQuantity - returnQuantity
		}
		for _, applyDetail := range shProduct {
			if applyDetail.SkuCode == product.SkuCode && product.GoodsId == 0 {
				return nil, errors.New("商品数据错误1")
			}
			if applyDetail.SkuCode == product.SkuCode && applyDetail.GoodsId == 0 {
				return nil, errors.New("商品数据错误2")
			}
			if product.SkuCode == applyDetail.SkuCode && product.GoodsId == applyDetail.GoodsId {
				//// 没有仓库信息不能退货
				//if product.SupplierWarehouse == "" {
				//	return nil, errors.New("商品仓库信息不存在无法申请退货")
				//}

				if applyDetail.Quantity > quantity {
					return nil, errors.New("申请售后商品数量大于可申请数量！")
				}

				if isAll == true && applyDetail.Quantity == product.OriginalQuantity {
					isAll = true
				} else {
					isAll = false
				}

				if isEnd == true && applyDetail.Quantity == quantity {
					isEnd = true
				} else {
					isEnd = false
				}

				if isOrderEnd == true && applyDetail.Quantity == productAllQuantity {
					isOrderEnd = true
				} else {
					isOrderEnd = false
				}

				okGoods[product.SkuCode] = true
			}
		}
		if _, ok := okGoods[product.SkuCode]; !ok {
			isAll = false
			if quantity > 0 {
				isEnd = false
			}
			if productAllQuantity > 0 {
				isOrderEnd = false
			}
		}
	}

	if len(okGoods) == 0 {
		return nil, errors.New("商品数据错误")
	}

	return &ConfirmProductStatus{
		IsAll:      isAll,
		IsEnd:      isEnd,
		IsOrderEnd: isOrderEnd,
	}, nil

}

// OrderReturn 售后退货 操作订单状态
func (e *CsApply) OrderReturn(db *gorm.DB, orderId string, csApplyDetailList []*CsApplyDetail, productIsEnd int) (err error) {
	// 如果有在售后中的就不用改订单状态 直到没有一个在售后中的
	isOrderInAfterSales, err := e.IsOrderInAfterSales(db, orderId)
	if !isOrderInAfterSales {
		orderInfo := OrderInfo{}
		db.Table(global.GetTenantOcDBNameWithDB(db)+".order_info").Where("order_id = ?", orderId).First(&orderInfo)
		saveOrderInfo := orderInfo
		saveOrderInfo.RmaStatus = 99
		// 如果全部退货了 并且订单状态是已发货的   订单状态改成已关闭
		if productIsEnd == 1 && orderInfo.OrderStatus == 1 {
			saveOrderInfo.OrderStatus = 10
		}
		err = e.OrderUpdate(db, orderId, &orderInfo, &saveOrderInfo, 99)
		if err != nil {
			return
		}
	}

	// 订单商品表处理
	for _, detail := range csApplyDetailList {
		err = e.UpdateReturnOrderDetail(db, orderId, detail.GoodsId, detail.Quantity)
		if err != nil {
			return err
		}
	}

	return
}

// IsOrderInAfterSales 通过订单号查询 该订单是否在售后中
func (e *CsApply) IsOrderInAfterSales(db *gorm.DB, orderID string) (bool, error) {
	var count int64

	ocPrefix := global.GetTenantOcDBNameWithDB(db)

	err := db.Table(ocPrefix+"."+e.TableName()).
		Where("order_id = ?", orderID).
		Where("cs_status = 0").
		Limit(1).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (e *CsApply) UpdateOrder(db *gorm.DB, orderID string, rmaStatus int) error {
	result := db.Table(global.GetTenantOcDBNameWithDB(db)+".order_info").
		Where("ORDER_ID = ?", orderID).
		Updates(map[string]interface{}{
			"rma_status": rmaStatus,
		})

	if result.Error != nil {
		return result.Error
	}

	//订单更新日志..

	return nil
}

// UpdateReturnOrderDetail 订单商品表退货处理
func (e *CsApply) UpdateReturnOrderDetail(db *gorm.DB, orderID string, goodsID int, quantity int) error {
	err := db.Table(global.GetTenantOcDBNameWithDB(db)+".order_detail").
		Where("order_id = ? AND goods_id = ?", orderID, goodsID).
		Updates(map[string]interface{}{
			"final_quantity": gorm.Expr("final_quantity - ?", quantity),
		}).Error
	return err
}

func (e *CsApply) applyEnd(db *gorm.DB, id int, csNo string, csApplyDetailList []*CsApplyDetail, data *CsApply) error {
	err := e.UpdateCsStatus(db, id, csNo, 99, nil, "")
	if err != nil {
		return err
	}
	var csApplyDetailAdds []*CsApplyDetail
	// 售后商品表添加实际售后的记录
	for i, detail := range csApplyDetailList {
		csApplyDetailAdds = append(csApplyDetailAdds, detail)
		csApplyDetailAdds[i].Id = 0
		csApplyDetailAdds[i].CsType = 1
	}

	ocPrefix := global.GetTenantOcDBNameWithDB(db)
	err = db.Table(ocPrefix + "." + CsApplyDetail{}.TableName()).Create(csApplyDetailAdds).Error
	if err != nil {
		return err
	}

	if data.CsType == 0 || data.CsType == 9 {
		var order OrderInfo
		err = db.Table(ocPrefix+"."+OrderInfo{}.TableName()).Where("order_id = ?", data.OrderId).First(&order).Error
		if err != nil {
			return err
		}
		if order.CreatedAt.Year() == time.Now().Year() && order.CreatedAt.Month() == time.Now().Month() {
			// 售后 退货取消完成后 调用预算记录更新方法
			// 释放预算
			var departmentBudgetM = modelsUc.DepartmentBudget{}
			err = departmentBudgetM.UpdateBudget(db, order.UserId, -data.ApplyPrice, time.Now().Format("200601"))
			//$this->model_department_budget->updateBudget($data['order_user_id'], $data['apply_price'] * (-1), date('Ym'));
		}
	}

	return nil
}

// UpdateCsStatus 更新售后单状态
func (e *CsApply) UpdateCsStatus(db *gorm.DB, id int, csNo string, csStatus int, data *CsApply, auditReason string) (err error) {
	ocPrefix := global.GetTenantOcDBNameWithDB(db)

	csApplyCsStatusText, ok := CsApplyStatusText[csStatus]
	if !ok {
		return errors.New("错误的状态修改")
	}
	if data == nil {
		data = &CsApply{}
		db.Table(ocPrefix+"."+data.TableName()).First(&data, id)
	}
	data.CsStatus = csStatus
	data.AuditReason = auditReason // 审核原因
	data.UpdateBy = user.GetUserId(db.Statement.Context.(*gin.Context))
	data.UpdateByName = user.GetUserName(db.Statement.Context.(*gin.Context))

	saveDb := db.Table(ocPrefix+"."+data.TableName()).Where("id = ?", id).Where("cs_status in (?)", CsApplyStatusUpdateOkList).Save(&data)
	if err = saveDb.Error; err != nil {
		return err
	}
	if saveDb.RowsAffected == 0 {

		return errors.New("无权更新该数据")
	}

	// 需要调用售后日志Service
	var csApplyLogM CsApplyLog
	err = csApplyLogM.AddLog(db, csNo, "售后单状态变更为:"+csApplyCsStatusText)
	if err != nil {
		return err
	}

	return nil
}

// OrderUpdate 订单表售后状态更新
func (e *CsApply) OrderUpdate(db *gorm.DB, orderId string, orderData *OrderInfo, updateOrder *OrderInfo, rmaStatus int) error {
	ocPrefix := global.GetTenantOcDBNameWithDB(db)

	// 添加订单日志
	if orderData == nil {
		orderData = &OrderInfo{}
		err := db.Table(ocPrefix+"."+orderData.TableName()).Where("order_id = ?", orderId).First(orderData).Error
		if err != nil {
			return err
		}
	}
	// 订单表更新
	result := db.Table(ocPrefix+".order_info").Where("order_id = ?", orderId).Save(updateOrder)
	if result.Error != nil {
		return result.Error
	}
	// 添加订单日志
	logType, _ := AfterSalesOrderLogType[rmaStatus]
	e.CreateOrderInfoLog(db, "", *orderData, *updateOrder, logType)
	return nil
}

// 生成日志
func (e *CsApply) CreateOrderInfoLog(db *gorm.DB, req interface{}, beforeModel OrderInfo, afterModel OrderInfo, logType string) {
	dataLog, _ := json.Marshal(&req)
	beforeDataStr := []byte("")
	if !reflect.DeepEqual(beforeModel, OrderInfo{}) {
		beforeDataStr, _ = json.Marshal(&beforeModel)
	}
	afterDataStr, _ := json.Marshal(&afterModel)
	log := OrderInfoLog{
		DataId:       afterModel.Id,
		Type:         logType,
		Data:         string(dataLog),
		BeforeData:   string(beforeDataStr),
		AfterData:    string(afterDataStr),
		CreateBy:     user.GetUserId(db.Statement.Context.(*gin.Context)),
		CreateByName: user.GetUserName(db.Statement.Context.(*gin.Context)),
	}
	_ = log.CreateLog("orderInfo", db)
}
