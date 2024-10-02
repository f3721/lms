package models

import (
	"go-admin/common/global"
	"time"

	"gorm.io/gorm"

	"go-admin/common/models"
)

// OrderInfoEditableOrderSources 包含可以编辑的订单来源
var OrderInfoEditableOrderSources = map[string]bool{
	"LMS": true,
}

type OrderInfo struct {
	models.Model

	OrderId                 string        `json:"orderId" gorm:"type:varchar(30);comment:订单编号"`
	OrderStatus             int           `json:"orderStatus" gorm:"type:tinyint(1);comment:订单状态：0-未发货、11-部分发货、1-已发货、2-部分收货、3-待评价、4-已评价、5-待确认、6、缺货、7-已签收、9-已取消、10-已关闭"`
	IsFormalOrder           int           `json:"isFormalOrder" gorm:"type:tinyint(1);comment:是否是正式订单(1:是, 0:否)"`
	IsTransform             int           `json:"isTransform" gorm:"type:tinyint(1);comment:是否转正式订单(预订单专用 0:未转正式订单, 1 已转正式订单)"`
	ItemsAmount             float64       `json:"itemsAmount" gorm:"type:int unsigned;comment:订单商品总金额"`
	TotalAmount             float64       `json:"totalAmount" gorm:"type:int unsigned;comment:订单总金额"`
	WarehouseCode           string        `json:"warehouseCode" gorm:"type:varchar(10);comment:发货仓"`
	Consignee               string        `json:"consignee" gorm:"type:varchar(20);comment:收货人姓名"`
	CompanyName             string        `json:"companyName" gorm:"type:varchar(50);comment:收货人公司名称"`
	Address                 string        `json:"address" gorm:"type:varchar(500);comment:收货人详细地址"`
	Mobile                  string        `json:"mobile" gorm:"type:varchar(20);comment:收货人手机号"`
	Telephone               string        `json:"telephone" gorm:"type:varchar(20);comment:收货人座机号"`
	Remark                  string        `json:"remark" gorm:"type:varchar(500);comment:客户留言"`
	ContractNo              string        `json:"contractNo" gorm:"type:varchar(100);comment:客户合同编号"`
	UserId                  int           `json:"userId" gorm:"type:int unsigned;comment:客户编号"`
	UserName                string        `json:"userName" gorm:"type:varchar(30);comment:客户用户名"`
	UserCompanyId           int           `json:"userCompanyId" gorm:"type:int unsigned;comment:客户公司ID"`
	UserCompanyName         string        `json:"userCompanyName" gorm:"type:varchar(50);comment:客户公司名称"`
	Ip                      string        `json:"ip" gorm:"type:varchar(20);comment:客户IP地址"`
	SendStatus              int           `json:"sendStatus" gorm:"type:tinyint unsigned;comment:发货状态：0-未发货、1-部分发货、2-全部发货"`
	SendTime                time.Time     `json:"sendTime" gorm:"type:datetime;comment:发货时间"`
	ValidFlag               int           `json:"validFlag" gorm:"type:tinyint unsigned;comment:有效标识位：0-无效、1-有效"`
	CreateFrom              string        `json:"createFrom" gorm:"type:varchar(30);comment:订单来源：LMS/MALL/XCX"`
	ClassifyQuantity        int           `json:"classifyQuantity" gorm:"type:int unsigned;comment:包含商品种类"`
	ProductQuantity         int           `json:"productQuantity" gorm:"type:int unsigned;comment:包含商品数量"`
	DeliverId               int           `json:"deliverId" gorm:"type:int unsigned;comment:收货地址ID"`
	ContactEmail            string        `json:"contactEmail" gorm:"type:varchar(50);comment:联系人邮箱"`
	OriginalItemsAmount     float64       `json:"originalItemsAmount" gorm:"type:int unsigned;comment:订单商品初始总金额"`
	OriginalTotalAmount     float64       `json:"originalTotalAmount" gorm:"type:int unsigned;comment:订单初始总金额"`
	CountryId               int           `json:"countryId" gorm:"type:int unsigned;comment:客户国家ID"`
	CountryName             string        `json:"countryName" gorm:"type:varchar(50);comment:客户国家名称"`
	ProvinceId              int           `json:"provinceId" gorm:"type:int unsigned;comment:客户省份ID"`
	ProvinceName            string        `json:"provinceName" gorm:"type:varchar(50);comment:客户省份名称"`
	CityId                  int           `json:"cityId" gorm:"type:int unsigned;comment:客户城市ID"`
	CityName                string        `json:"cityName" gorm:"type:varchar(20);comment:客户城市名称"`
	AreaId                  int           `json:"areaId" gorm:"type:int unsigned;comment:客户区县ID"`
	AreaName                string        `json:"areaName" gorm:"type:varchar(50);comment:客户区县名称"`
	TownId                  int           `json:"townId" gorm:"type:int unsigned;comment:镇/街道 ID"`
	TownName                string        `json:"townName" gorm:"type:varchar(50);comment:镇/街道名称"`
	ConfirmTime             time.Time     `json:"confirmTime" gorm:"type:datetime;comment:订单确认时间"`
	ConfirmOrderReceiptTime time.Time     `json:"confirmOrderReceiptTime" gorm:"type:datetime;comment:订单确认签收时间（客户发起签收）"`
	FinalTotalAmount        float64       `json:"finalTotalAmount" gorm:"type:decimal(10,2);comment:最终订单总金额（对账单生成时各明细FINAL_SUB_TOTAL_AMOUNT之和）"`
	IsOverBudget            int           `json:"isOverBudget" gorm:"type:tinyint(1);comment:是否超出预算 -1-未计算 1-是 0-否"`
	OrderFileRemark         string        `json:"orderFileRemark" gorm:"type:varchar(500);comment:订单附件备注"`
	LogisticalRemark        string        `json:"logisticalRemark" gorm:"type:varchar(500);comment:物流备注"`
	RmaStatus               int           `json:"rmaStatus" gorm:"type:tinyint(1);comment:售后状态：0-无售后、1-售后处理中、2-售后已确认、99-售后处理完成"`
	ReceiveStatus           int           `json:"receiveStatus" gorm:"type:tinyint(1);comment:收货状态：0-未收货、1-部分收货、2-全部收货"`
	ExternalOrderNo         string        `json:"externalOrderNo" gorm:"type:varchar(100);comment:审批单号"`
	CreatedAt               time.Time     `json:"createdAt" gorm:"comment:创建时间"`
	UpdatedAt               time.Time     `json:"updatedAt" gorm:"comment:最后更新时间"`
	CreateBy                int           `json:"createBy" gorm:"index;comment:创建者"`
	CreateByName            string        `json:"createByName" gorm:"index;comment:创建者姓名"`
	CancelBy                int           `json:"cancelBy" gorm:"type:int unsigned;comment:取消人id"`
	CancelByName            string        `json:"cancelByName" gorm:"type:varchar(20);comment:取消人姓名"`
	CancelByType            int           `json:"cancelByType" gorm:"type:tinyint unsigned;comment:取消人类型：0-客户、1-销售"`
	OrderDetails            []OrderDetail `json:"orderDetails" gorm:"foreignkey:OrderId;references:OrderId"`
	OrderImages             []OrderImage  `json:"orderImages" gorm:"foreignkey:OrderId;references:OrderId"`
	ReceiptImages           []OrderImage  `json:"receiptImages" gorm:"foreignkey:OrderId;references:OrderId"`
}

func (OrderInfo) TableName() string {
	return "order_info"
}

func (e *OrderInfo) GetId() interface{} {
	return e.Id
}

func (e *OrderInfo) GetRow(tx *gorm.DB, id int) (row OrderInfo) {
	tx.Table(e.TableName()).First(&row, id)
	return
}

// GetOrderMapByOrderIds 获取以orderId 为key 的订单map
func (e *OrderInfo) GetOrderMapByOrderIds(tx *gorm.DB, orderIds []string) (res map[string]OrderInfo) {
	var orders []OrderInfo
	tx.Table(e.TableName()).Where("order_id in ?", orderIds).First(&orders)
	res = make(map[string]OrderInfo)
	for _, order := range orders {
		res[order.OrderId] = order
	}
	return
}

func (e *OrderInfo) GetRowByOrderId(tx *gorm.DB, orderId string) (row OrderInfo) {
	tx.Table(e.TableName()).Where("order_id = ?", orderId).First(&row)
	return
}

// ShipmentUpdateOrder 出货修改订单 - 逐步废弃
func (e *OrderInfo) ShipmentUpdateOrder(tx *gorm.DB, orderId string, quantityList map[int]int) (err error) {
	ocPrefix := global.GetTenantOcDBNameWithDB(tx)
	err = tx.Table(ocPrefix+"."+e.TableName()).Where("order_id = ?", orderId).First(e).Error
	if err != nil {
		return
	}
	e.OrderStatus = 1
	e.SendStatus = 2
	e.SendTime = time.Now()
	err = tx.Table(ocPrefix + "." + e.TableName()).Save(e).Error
	if err != nil {
		return
	}

	var orderDetails []OrderDetail
	err = tx.Table(ocPrefix+".order_detail").Where("order_id = ?", orderId).Find(&orderDetails).Error
	if err != nil {
		return
	}

	for _, detail := range orderDetails {
		if quantity, ok := quantityList[detail.GoodsId]; ok {
			detail.SendQuantity = quantity
			err = tx.Table(ocPrefix + ".order_detail").Save(&detail).Error
			if err != nil {
				return
			}
		}
	}

	return
}

// ShipmentUpdateOrder 部分发货修改订单
//
// db,订单ID,是否全部出库,部分出库商品map
func (e *OrderInfo) OutboundUpdateOrder(tx *gorm.DB, orderId string, allOutbound bool, quantityList map[int]int) (err error) {
	ocPrefix := global.GetTenantOcDBNameWithDB(tx)
	err = tx.Table(ocPrefix+"."+e.TableName()).Where("order_id = ?", orderId).First(e).Error
	if err != nil {
		return
	}

	if allOutbound { // 全部出库
		e.OrderStatus = 1 // 已发货
		e.SendStatus = 2  // 全部发货
	} else {
		e.OrderStatus = 11 // 部分发货
		e.SendStatus = 1   // 部分发货
	}
	e.SendTime = time.Now()
	err = tx.Table(ocPrefix + "." + e.TableName()).Save(e).Error
	if err != nil {
		return
	}

	var orderDetails []OrderDetail
	err = tx.Table(ocPrefix+".order_detail").Debug().Where("order_id = ?", orderId).Find(&orderDetails).Error
	if err != nil {
		return
	}

	for _, detail := range orderDetails {
		if quantity, ok := quantityList[detail.GoodsId]; ok {
			detail.SendQuantity = quantity
			detailInrc := map[string]any{
				"send_quantity": gorm.Expr("send_quantity + ?", quantity),
			}
			// err = tx.Model(&detail).Updates(detailInrc).Error
			err = tx.Table(ocPrefix+".order_detail").Where("id = ?", detail.Id).Updates(detailInrc).Error
			if err != nil {
				return
			}
		}
	}

	return
}

type EnsureApprovalEmail struct {
	UserName      string        `json:"userName" comment:"审批人名称"`
	OrderId       string        `json:"orderId" comment:"订单号"`
	ApproveUrl    string        `json:"approveUrl" comment:"审批URL地址"`
	WarehouseName string        `json:"warehouseName" comment:""`
	LogoUrl       string        `json:"logoUrl" comment:"logo地址"`
	OrderDetails  []OrderDetail `json:"orderDetails" comment:"产品信息"`
}
