package models

import (
	"go-admin/common/global"
	"go-admin/common/utils"
	"time"

	"gorm.io/gorm"

	"go-admin/common/models"
)

type OrderDetail struct {
	models.Model

	OrderId             string    `json:"orderId" gorm:"type:varchar(30);comment:订单编号"`
	SendQuantity        int       `json:"sendQuantity" gorm:"type:int unsigned;comment:已发货商品数量"`
	UserId              int       `json:"userId" gorm:"type:int unsigned;comment:用户ID"`
	UserName            string    `json:"userName" gorm:"type:varchar(200);comment:用户名称"`
	GoodsId             int       `json:"goodsId" gorm:"type:int unsigned;comment:goods表ID"`
	WarehouseCode       string    `json:"warehouseCode" gorm:"type:varchar(10);comment:发货仓"`
	CatId               int       `json:"catId" gorm:"type:int unsigned;comment:一级产线ID"`
	CatName             string    `json:"catName" gorm:"type:varchar(50);comment:一级产线名称"`
	CatId2              int       `json:"catId2" gorm:"type:int unsigned;comment:二级产线ID"`
	CatName2            string    `json:"catName2" gorm:"type:varchar(50);comment:二级产线名称"`
	CatId3              int       `json:"catId3" gorm:"type:int unsigned;comment:三级产线ID"`
	CatName3            string    `json:"catName3" gorm:"type:varchar(50);comment:三级产线名称"`
	CatId4              int       `json:"catId4" gorm:"type:int unsigned;comment:四级产线ID"`
	CatName4            string    `json:"catName4" gorm:"type:varchar(50);comment:四级产线名称"`
	SkuCode             string    `json:"skuCode" gorm:"type:varchar(20);comment:sku"`
	ProductId           int       `json:"productId" gorm:"type:int unsigned;comment:商品ID"`
	ProductName         string    `json:"productName" gorm:"type:varchar(512);comment:商品名称"`
	ProductPic          string    `json:"productPic" gorm:"type:varchar(200);comment:商品图片"`
	ProductModel        string    `json:"productModel" gorm:"type:varchar(255);comment:商品型号"`
	BrandId             int       `json:"brandId" gorm:"type:int unsigned;comment:商品品牌ID"`
	BrandName           string    `json:"brandName" gorm:"type:varchar(100);comment:商品品牌名字"`
	BrandEname          string    `json:"brandEname" gorm:"type:varchar(100);comment:商品品牌英文"`
	Unit                string    `json:"unit" gorm:"type:varchar(10);comment:商品单位"`
	SalePrice           float64   `json:"salePrice" gorm:"type:int unsigned;comment:商品销售价格"`
	PurchasePrice       float64   `json:"purchasePrice" gorm:"type:int unsigned;comment:商品采购价"`
	Quantity            int       `json:"quantity" gorm:"type:int unsigned;comment:商品数量"`
	CancelQuantity      int       `json:"cancelQuantity" gorm:"type:int unsigned;comment:已取消数量"`
	LockStock           int       `json:"lockStock" gorm:"type:int unsigned;comment:已锁库存"`
	VendorId            int       `json:"vendorId" gorm:"type:int unsigned;comment:供应商信息ID"`
	SubTotalAmount      float64   `json:"subTotalAmount" gorm:"type:int unsigned;comment:行项目小计金额"`
	OriginalItemsMount  float64   `json:"originalItemsMount" gorm:"type:int unsigned;comment:原商品总金额"`
	OriginalQuantity    int       `json:"originalQuantity" gorm:"type:int unsigned;comment:原商品数量"`
	BatchGroup          int       `json:"batchGroup" gorm:"type:tinyint unsigned;comment:发货分组"`
	SupplierWarehouse   string    `json:"supplierWarehouse" gorm:"type:varchar(10);comment:供应商仓库编号"`
	DeliveryWarehouse   string    `json:"deliveryWarehouse" gorm:"type:varchar(10);comment:西域发货仓"`
	MarketPrice         float64   `json:"marketPrice" gorm:"type:int unsigned;comment:系统价格"`
	Moq                 int       `json:"moq" gorm:"type:int unsigned;comment:最小起订量"`
	UserProductName     string    `json:"userProductName" gorm:"type:varchar(255);comment:客户商品物料描述"`
	UserProductRemark   string    `json:"userProductRemark" gorm:"type:varchar(255);comment:客户商品物料采购备注"`
	FinalSubTotalAmount float64   `json:"finalSubTotalAmount" gorm:"type:int unsigned;comment:最终行项目小计金额"`
	FinalQuantity       int       `json:"finalQuantity" gorm:"type:int unsigned;comment:最终数量（总数量-已取消-已退货）"`
	ProductNo           string    `json:"productNo" gorm:"type:varchar(30);comment:物料编码"`
	CreatedAt           time.Time `json:"createdAt" gorm:"comment:创建时间"`
	UpdatedAt           time.Time `json:"updatedAt" gorm:"comment:最后更新时间"`
}

func (OrderDetail) TableName() string {
	return "order_detail"
}

//func (e *OrderDetail) Generate() models.ActiveRecord {
//	o := *e
//	return &o
//}

func (e *OrderDetail) GetId() interface{} {
	return e.Id
}

func (e *OrderDetail) SetFinalQuantityAndAmountAndGetFinalTotalAmount(tx *gorm.DB, orderId string) (totalAmount float64) {
	var details []OrderDetail
	tx.Table(e.TableName()).Where("order_id = ?", orderId).Find(&details)
	if len(details) == 0 {
		return
	}
	var csApply CsApply
	returnList := csApply.GetAfterReturnProductsBySaleId(tx, orderId)
	for _, detail := range details {
		returnQuantity := 0
		if _, ok := returnList[detail.SkuCode]; ok {
			returnQuantity = returnList[detail.SkuCode].Quantity
		}
		detail.FinalQuantity = detail.Quantity - detail.CancelQuantity - returnQuantity
		detail.FinalSubTotalAmount = utils.MulFloat64AndInt(detail.SalePrice, detail.FinalQuantity)
		tx.Save(&detail)
		totalAmount = totalAmount + detail.FinalSubTotalAmount
	}
	var order OrderInfo
	tx.Model(&order).Where("order_id = ?", orderId).First(&order)
	order.FinalTotalAmount = totalAmount
	tx.Save(&order)
	return
}

// 查询goodsId 订单缺货数
type OrderLackStock struct {
	GoodsId   int `json:"goodsId"`
	LackStock int `json:"lackStock"`
}

func (e *OrderDetail) GetOrderLackStock(tx *gorm.DB, goodsIds []int) (res map[int]int) {
	ocPrefix := global.GetTenantOcDBNameWithDB(tx)

	res = make(map[int]int)
	var list []OrderLackStock
	tx.Debug().Table(ocPrefix+"."+e.TableName()+" od").
		Joins("left join "+ocPrefix+".order_info oi on oi.order_id = od.order_id").
		Select("od.goods_id, sum(od.quantity-od.lock_stock-od.cancel_quantity) as lack_stock").
		Where("oi.order_status = ?", 6).
		Where(" od.goods_id in ?", goodsIds).
		Group("od.goods_id").
		Scan(&list)
	for _, v := range list {
		res[v.GoodsId] = v.LackStock
	}
	return
}

type ReturnQuantiryStruct struct {
	OrderId        string `json:"orderId"`
	ReturnQuantiry int    `json:"returnQuantiry"`
}

// 获取订单最终商品数量
func (e *OrderDetail) GetReturnQuantiryMap(tx *gorm.DB, orderIds []string) (res map[string]ReturnQuantiryStruct) {
	var data []ReturnQuantiryStruct
	err := tx.Raw(`
			SELECT order_id, sum(final_quantity) return_quantiry
			FROM order_detail as t
			WHERE order_id IN ? GROUP BY order_id
		`, orderIds).Scan(&data).Error
	if err != nil {
		return nil
	}
	res = make(map[string]ReturnQuantiryStruct)
	_ = utils.StructColumn(&res, data, "", "OrderId")
	return
}
