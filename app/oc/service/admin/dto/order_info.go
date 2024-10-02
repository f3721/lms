package dto

import (
	"errors"
	"fmt"
	modelsPc "go-admin/app/pc/models"
	"go-admin/common/actions"
	"go-admin/common/utils"
	"math/rand"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"go-admin/app/oc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type OrderInfoGetPageReq struct {
	dto.Pagination `search:"-"`
	//OrderId string     `search:"type:in;column:order_id;table:order_info"`
	UserId        string `form:"userId"  search:"type:exact;column:user_id;table:order_info"`
	OrderStatus   []int  `form:"orderStatus[]"  search:"type:in;column:order_status;table:order_info"`
	RmaStatus     []int  `form:"rmaStatus[]"  search:"type:in;column:rma_status;table:order_info"`
	CreateFrom    string `form:"createFrom"  search:"type:exact;column:create_from;table:order_info"`
	OrderIds      string `search:"-" form:"orderIds"`
	Start         string `form:"startTime" search:"type:gte;column:created_at;table:order_info"`
	End           string `form:"endTime" search:"type:lte;column:created_at;table:order_info"`
	IsOverBudget  string `form:"isOverBudget"  search:"type:exact;column:is_over_budget;table:order_info"`
	WarehouseCode string `form:"warehouseCode"  search:"type:exact;column:warehouse_code;table:order_info"`
	UserCompanyId string `form:"userCompanyId"  search:"type:exact;column:user_company_id;table:order_info"`
	SkuCode       string `form:"skuCode"  search:"-"`
	ProductName   string `form:"productName"  search:"-"`
	ProductNo     string `form:"productNo"  search:"-"`
	OrderInfoOrder
}

type OrderInfoOrder struct {
	Id string `form:"idOrder"  search:"type:order;column:id;table:order_info"`
	//CreatedAt string `form:"createdAtOrder"  search:"type:order;column:created_at;table:order_info"`
	//UpdatedAt string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:order_info"`

}

func (m *OrderInfoGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type OrderInfoInsertReq struct {
	Id               int                             `json:"-" comment:""` //
	UserId           int                             `json:"userId" comment:"客户编号"`
	UserName         string                          `json:"userName" comment:"客户用户名"`
	UserCompanyId    int                             `json:"userCompanyId" comment:"客户公司ID"`
	UserCompanyName  string                          `json:"userCompanyName" comment:"客户公司名称"`
	ContractNo       string                          `json:"contractNo" comment:"客户合同编号"`
	Remark           string                          `json:"remark" comment:"客户留言"`
	OrderFileRemark  string                          `json:"orderFileRemark" comment:"订单附件备注"`
	OrderImages      []OrderInfoOrderImagesInsertReq `json:"orderImages" comment:"订单相关文件"`
	Products         []OrderInfoProductInsertReq     `json:"products" comment:"商品"`
	DeliverId        int                             `json:"deliverId" comment:"收货地址ID"`
	Address          string                          `json:"address" comment:"收货人详细地址"`
	Consignee        string                          `json:"consignee" comment:"收货人姓名"`
	Mobile           string                          `json:"mobile" comment:"收货人手机号"`
	Telephone        string                          `json:"telephone" comment:"收货人座机号"`
	Email            string                          `json:"email" comment:"邮箱"`
	CountryId        int                             `json:"countryId" comment:"客户国家ID"`
	CountryName      string                          `json:"countryName" comment:"客户国家名称"`
	ProvinceId       int                             `json:"provinceId" comment:"客户省份ID"`
	ProvinceName     string                          `json:"provinceName" comment:"客户省份名称"`
	CityId           int                             `json:"cityId" comment:"收货人城市编号"`
	CityName         string                          `json:"cityName" comment:"收货人城市名称"`
	AreaId           int                             `json:"areaId" comment:"客户区县ID"`
	AreaName         string                          `json:"areaName" comment:"客户区县名称"`
	TownId           int                             `json:"townId" comment:"镇/街道 ID"`
	TownName         string                          `json:"townName" comment:"镇/街道名称"`
	LogisticalRemark string                          `json:"logisticalRemark" comment:"物流备注"`
	CreateBy         int                             `json:"createBy" gorm:"index;comment:创建者"`
	CreateByName     string                          `json:"createByName" gorm:"index;comment:创建者姓名"`
	WarehouseCode    string                          `json:"warehouseCode" comment:"发货仓"`
}

type OrderInfoOrderImagesInsertReq struct {
	Name string `json:"name" comment:"图片名称"`
	Url  string `json:"url" comment:"url"`
}

type OrderInfoProductInsertReq struct {
	SkuCode          string  `json:"skuCode" comment:"sku"`
	ProductName      string  `json:"productName" comment:"商品名称"`
	ProductModel     string  `json:"productModel" comment:"商品型号"`
	BrandId          int     `json:"brandId" comment:"商品品牌ID"`
	BrandName        string  `json:"brandName" comment:"商品品牌名字"`
	BrandEname       string  `json:"brandEname" comment:"商品品牌英文"`
	ProductId        int     `json:"productId" comment:"comment:商品ID"`
	Quantity         int     `json:"quantity" comment:"商品数量"`
	Moq              int     `json:"moq" comment:"销售最小起订量"`
	Unit             string  `json:"unit" comment:"商品单位"`
	SalePrice        float64 `json:"salePrice" comment:"商品销售价格(单位:分)"`
	VendorId         int     `json:"vendorId" comment:"供应商信息ID"`
	VendorName       string  `json:"vendorName" comment:"货主名称"`
	VendorSkuCode    string  `json:"vendorSkuCode" comment:"货主sku"`
	ProductNo        string  `json:"productNo" comment:"物料编码"`
	GoodsId          int     `json:"goodsId" comment:"goods表ID"`
	Stock            int     `json:"stock" comment:"可用库存"`
	Tax              float64 `json:"tax" comment:"税率"`
	SubTotalAmount   float64 `json:"subTotalAmount" comment:"行项目小计金额"`
	ProductPic       string  `json:"productPic" comment:"商品图片"`
	CompanyId        int     `json:"companyId" comment:"公司ID"`
	CheckStockStatus int     `json:"checkStockStatus" comment:"库存不足是否允许下单 0否 1是"`
}

func (s *OrderInfoInsertReq) Generate(model *models.OrderInfo) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.OrderId = generateOrderId()
	model.OrderStatus = 5
	model.WarehouseCode = s.WarehouseCode
	model.DeliverId = s.DeliverId
	model.Consignee = s.Consignee
	model.Address = s.Address
	model.Mobile = s.Mobile
	model.Telephone = s.Telephone
	model.ContactEmail = s.Email
	model.Remark = s.Remark
	model.ContractNo = s.ContractNo
	model.UserId = s.UserId
	model.UserName = s.UserName
	model.UserCompanyId = s.UserCompanyId
	model.UserCompanyName = s.UserCompanyName
	model.OrderFileRemark = s.OrderFileRemark
	model.LogisticalRemark = s.LogisticalRemark
	model.CountryId = s.CountryId
	model.CountryName = s.CountryName
	model.ProvinceId = s.ProvinceId
	model.ProvinceName = s.ProvinceName
	model.CityId = s.CityId
	model.CityName = s.CityName
	model.AreaId = s.AreaId
	model.AreaName = s.AreaName
	model.TownId = s.TownId
	model.TownName = s.TownName
	model.IsOverBudget = -1
}

func (s *OrderInfoInsertReq) GetId() interface{} {
	return s.Id
}

func (s *OrderInfoInsertReq) InsertValid(tx *gorm.DB) (err error) {
	// 基本信息
	if s.UserCompanyId <= 0 {
		return errors.New("请选择公司")
	}
	if s.UserId <= 0 {
		return errors.New("请选择客户")
	}
	if s.ContractNo == "" {
		return errors.New("请输入外部单号")
	}
	if s.Remark != "" && utf8.RuneCountInString(s.Remark) > 200 {
		return errors.New("下单备注过长(最多200个字符)！")
	}
	if s.OrderFileRemark != "" && utf8.RuneCountInString(s.OrderFileRemark) > 200 {
		return errors.New("订单附件备注过长(最多200个字符)！")
	}
	if len(s.OrderImages) <= 0 {
		return errors.New("请上传订单相关文件")
	}

	// 配送
	if s.DeliverId <= 0 || s.CountryId <= 0 || s.ProvinceId <= 0 || s.CityId <= 0 || s.AreaId <= 0 || s.TownId <= 0 || s.Address == "" || s.Consignee == "" {
		return errors.New("选择的地址信息不完整！")
	}
	if s.Mobile == "" && s.Telephone == "" {
		return errors.New("收货人电话及手机至少有一个！")
	}
	if s.LogisticalRemark != "" && utf8.RuneCountInString(s.LogisticalRemark) > 200 {
		return errors.New("物流备注备注过长(最多200个字符)！")
	}

	// 产品
	if s.WarehouseCode == "" {
		return errors.New("发货仓必填")
	}
	if len(s.Products) <= 0 {
		return errors.New("请添加产品")
	}
	err = checkProduct(s.Products)
	if err != nil {
		return err
	}

	return
}

func (s *OrderInfoProductInsertReq) Generate(model *models.OrderDetail) {

	_ = copier.Copy(model, s)
	model.MarketPrice = model.SalePrice
	model.FinalQuantity = model.Quantity
	model.OriginalQuantity = model.Quantity
	model.OriginalItemsMount = utils.MulFloat64AndInt(model.SalePrice, model.Quantity)
	model.SubTotalAmount = model.OriginalItemsMount
	//model.ProductId = s.ProductId
}

func (s *OrderInfoOrderImagesInsertReq) Generate(model *models.OrderImage) {
	model.Name = s.Name
	model.Url = s.Url
}

func generateOrderId() string {
	return "RO" + time.Now().Format("20060102150405") + fmt.Sprintf("%04d", (1+rand.Intn(9999)))
}

func GetPageMakeCondition(c *OrderInfoGetPageReq, newDb *gorm.DB) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {

		var orderIds [][]string
		if c.OrderIds != "" {
			orderIds = append(orderIds, utils.Split(c.OrderIds))
		}

		var orderDetailM models.OrderDetail
		var ids []string
		if c.SkuCode != "" {
			newDb.Model(&orderDetailM).Select("order_id").Where("sku_code = ?", c.SkuCode).Find(&ids)
			orderIds = append(orderIds, ids)
		}
		if c.ProductName != "" {
			newDb.Model(&orderDetailM).Select("order_id").Where("product_name like ?", "%"+c.ProductName+"%").Find(&ids)
			orderIds = append(orderIds, ids)
		}
		if c.ProductNo != "" {
			newDb.Model(&orderDetailM).Select("order_id").Where("product_no in ?", utils.Split(c.ProductNo)).Find(&ids)
			orderIds = append(orderIds, ids)
		}

		if len(orderIds) > 0 {
			a := utils.IntersectionString(orderIds)
			db.Where("order_id in ?", a)
		}
		return db
	}
}

func OrderInfoGetPageCompanyPermission(tableName string, p *actions.DataPermission) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(tableName+".user_company_id in ?", utils.SplitToInt(p.AuthorityCompanyId))
	}
}

type OrderInfoUpdateReq struct {
	Id              int                             `uri:"id" comment:""` //
	ContractNo      string                          `json:"contractNo" comment:"客户合同编号"`
	Remark          string                          `json:"remark" comment:"客户留言"`
	OrderFileRemark string                          `json:"orderFileRemark" comment:"订单附件备注"`
	OrderImages     []OrderInfoOrderImagesInsertReq `json:"orderImages" comment:"订单相关文件"`
}

func (s *OrderInfoUpdateReq) Generate(model *models.OrderInfo) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}

	model.Remark = s.Remark
	model.ContractNo = s.ContractNo
	model.OrderFileRemark = s.OrderFileRemark
}

func (s *OrderInfoUpdateReq) GetId() interface{} {
	return s.Id
}

func (s *OrderInfoUpdateReq) UpdateValid(tx *gorm.DB) (err error) {
	if len(s.OrderImages) <= 0 {
		return errors.New("请上传订单相关文件")
	}
	if s.ContractNo == "" {
		return errors.New("请输入外部单号")
	}

	var object models.OrderInfo
	orderInfo := object.GetRow(tx, s.Id)
	if orderInfo.OrderStatus != 5 && orderInfo.OrderStatus != 6 {
		return errors.New("该订单不可编辑")
	}
	if _, ok := models.OrderInfoEditableOrderSources[orderInfo.CreateFrom]; !ok {
		return errors.New("该订单不可编辑")
	}

	return
}

type OrderInfoReceiptReq struct {
	Id            int                             `uri:"id" comment:""` //id
	ReceiptImages []OrderInfoOrderImagesInsertReq `json:"receiptImages" comment:"订单签收文件"`
	IsAuto        int                             `json:"-"` // 是否自动签收 默认0 只有自动签收脚本时传1
}

func (s *OrderInfoReceiptReq) Generate(model *models.OrderInfo) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
}

func (s *OrderInfoReceiptReq) GetId() interface{} {
	return s.Id
}

func (s *OrderInfoReceiptReq) ReceiptValid(tx *gorm.DB) (err error) {
	if len(s.ReceiptImages) <= 0 {
		return errors.New("请上传订单签收文件")
	}
	if len(s.ReceiptImages) > 6 {
		return errors.New("订单签收文件最多上传6张")
	}

	var object models.OrderInfo
	orderInfo := object.GetRow(tx, s.Id)
	if orderInfo.OrderStatus == 1 && (orderInfo.RmaStatus == 1 || orderInfo.RmaStatus == 2) {
		return errors.New("此订单正处于售后中！")
	}

	return
}

type OrderInfoGetReceiptImageReq struct {
	OrderId string `uri:"orderId" comment:""` //OrderId
}

type OrderInfoGetReceiptImageRes struct {
	List *[]*models.OrderImage `json:"list"`
}

type OrderInfoSaveReceiptImageReq struct {
	OrderId       string                          `uri:"orderId" comment:""` //OrderId
	ReceiptImages []OrderInfoOrderImagesUpdateReq `json:"receiptImages" comment:"订单签收文件"`
}

type OrderInfoOrderImagesUpdateReq struct {
	Id   int    `json:"id"` // 如果是原来的图片就有id 如果是新的图片没有id可以传0或者不传这个字段
	Name string `json:"name" comment:"图片名称"`
	Url  string `json:"url" comment:"url"`
}

func (s *OrderInfoOrderImagesUpdateReq) Generate(model *models.OrderImage) {
	if s.Id > 0 {
		model.Id = s.Id
	}
	model.Name = s.Name
	model.Url = s.Url
}

type OrderInfoCancelReq struct {
	Id           int    `uri:"id" comment:""` //
	Remark       string `json:"remark" comment:"取消原因"`
	CancelBy     int    `json:"cancelBy" comment:"取消人id"`
	CancelByName string `json:"cancelByName" comment:"取消人姓名"`
}

func (s *OrderInfoCancelReq) Generate(model *models.OrderInfo) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.Remark = s.Remark
	model.CancelBy = s.CancelBy
	model.CancelByName = s.CancelByName
}

func (s *OrderInfoCancelReq) GetId() interface{} {
	return s.Id
}

func (s *OrderInfoCancelReq) CancelValid(tx *gorm.DB) (err error) {
	var object models.OrderInfo
	orderInfo := object.GetRow(tx, s.Id)
	if orderInfo.OrderStatus != 5 && orderInfo.OrderStatus != 6 {
		return errors.New("此订单不可取消")
	}
	if orderInfo.RmaStatus != 0 && orderInfo.RmaStatus != 99 {
		return errors.New("订单正在售后中，不可取消")
	}

	return
}

// OrderInfoGetReq 功能获取请求参数
type OrderInfoGetReq struct {
	Id int `uri:"id"`
}

func (s *OrderInfoGetReq) GetId() interface{} {
	return s.Id
}

// OrderInfoDeleteReq 功能删除请求参数
type OrderInfoDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *OrderInfoDeleteReq) GetId() interface{} {
	return s.Ids
}

type OrderInfoUpdateShippingReq struct {
	Id               int    `uri:"id" comment:""` //
	DeliverId        int    `json:"deliverId" comment:"收货地址ID"`
	Address          string `json:"address" comment:"收货人详细地址"`
	Consignee        string `json:"consignee" comment:"收货人姓名"`
	Mobile           string `json:"mobile" comment:"收货人手机号"`
	Telephone        string `json:"telephone" comment:"收货人座机号"`
	Email            string `json:"email" comment:"邮箱"`
	CountryId        int    `json:"countryId" comment:"客户国家ID"`
	CountryName      string `json:"countryName" comment:"客户国家名称"`
	ProvinceId       int    `json:"provinceId" comment:"客户省份ID"`
	ProvinceName     string `json:"provinceName" comment:"客户省份名称"`
	CityId           int    `json:"cityId" comment:"收货人城市编号"`
	CityName         string `json:"cityName" comment:"收货人城市名称"`
	AreaId           int    `json:"areaId" comment:"客户区县ID"`
	AreaName         string `json:"areaName" comment:"客户区县名称"`
	TownId           int    `json:"townId" comment:"镇/街道 ID"`
	TownName         string `json:"townName" comment:"镇/街道名称"`
	LogisticalRemark string `json:"logisticalRemark" comment:"物流备注"`
}

func (s *OrderInfoUpdateShippingReq) Generate(model *models.OrderInfo) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}

	model.CountryId = s.CountryId
	model.CountryName = s.CountryName
	model.ProvinceId = s.ProvinceId
	model.ProvinceName = s.ProvinceName
	model.CityId = s.CityId
	model.CityName = s.CityName
	model.AreaId = s.AreaId
	model.AreaName = s.AreaName
	model.TownId = s.TownId
	model.TownName = s.TownName
	model.LogisticalRemark = s.LogisticalRemark
}

func (s *OrderInfoUpdateShippingReq) GetId() interface{} {
	return s.Id
}

func (s *OrderInfoUpdateShippingReq) UpdateValid(tx *gorm.DB) (err error) {
	if s.DeliverId <= 0 || s.CountryId <= 0 || s.ProvinceId <= 0 || s.CityId <= 0 || s.AreaId <= 0 || s.TownId <= 0 || s.Address == "" || s.Consignee == "" {
		return errors.New("选择的地址信息不完整！")
	}
	if s.Mobile == "" && s.Telephone == "" {
		return errors.New("收货人电话及手机至少有一个！")
	}
	if s.LogisticalRemark != "" && utf8.RuneCountInString(s.LogisticalRemark) > 200 {
		return errors.New("物流备注备注过长(最多200个字符)！")
	}
	var object models.OrderInfo
	orderInfo := object.GetRow(tx, s.Id)
	if orderInfo.OrderStatus != 5 && orderInfo.OrderStatus != 6 {
		return errors.New("该订单不可编辑")
	}

	return
}

type OrderInfoUpdateProductReq struct {
	Id            int                         `uri:"id" comment:""` //
	WarehouseCode string                      `json:"warehouseCode" comment:"发货仓"`
	Products      []OrderInfoProductInsertReq `json:"products" comment:"商品"`
}

func (s *OrderInfoUpdateProductReq) Generate(model *models.OrderInfo) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
}

func (s *OrderInfoUpdateProductReq) GetId() interface{} {
	return s.Id
}

func (s *OrderInfoUpdateProductReq) UpdateValid(tx *gorm.DB) (err error) {
	if len(s.Products) <= 0 {
		return errors.New("请添加产品")
	}

	// 产品
	if s.WarehouseCode == "" {
		return errors.New("发货仓必填")
	}
	if len(s.Products) <= 0 {
		return errors.New("请添加产品")
	}

	err = checkProduct(s.Products)
	if err != nil {
		return err
	}

	var object models.OrderInfo
	orderInfo := object.GetRow(tx, s.Id)
	if orderInfo.OrderStatus != 5 && orderInfo.OrderStatus != 6 {
		return errors.New("该订单不可编辑")
	}
	if orderInfo.RmaStatus != 0 && orderInfo.RmaStatus != 99 {
		return errors.New("订单正在售后中，不可编辑")
	}
	if _, ok := models.OrderInfoEditableOrderSources[orderInfo.CreateFrom]; !ok {
		return errors.New("该订单不可编辑")
	}

	return
}

func checkProduct(s []OrderInfoProductInsertReq) (err error) {
	var limitSaleMoqSkus []string
	for _, product := range s {
		if product.Quantity < product.Moq {
			limitSaleMoqSkus = append(limitSaleMoqSkus, product.SkuCode)
		}
	}
	if len(limitSaleMoqSkus) > 0 {
		return errors.New("产品[" + strings.Join(limitSaleMoqSkus, ",") + "]未达到最小起订量!")
	}

	return
}

type OrderInfoByOrderIdReq struct {
	OrderIds []string `form:"orderIds"  search:"type:in;column:order_id;table:order_info" vd:"@:len($[0])>0 ; msg:'订单id必填'"`
}

func (m *OrderInfoByOrderIdReq) GetNeedSearch() interface{} {
	return *m
}

type OrderInfoCheckExistUnCompletedOrderReq struct {
	Id int `uri:"id" comment:""` //
}

type OrderInfoAddProductReq struct {
	SkuCode       string `json:"skuCode" comment:"SKU（多个逗号分割）"`
	WarehouseCode string `json:"warehouseCode" comment:"发货仓库code"`
}

type OrderInfoAddProductResp struct {
	SkuCode       string `json:"skuCode" comment:"SKU（多个逗号分割）"`
	WarehouseCode string `json:"warehouseCode" comment:"发货仓库code"`
}

func (s *OrderInfoAddProductReq) Valid(tx *gorm.DB, productMap map[string]modelsPc.Goods) (err error) {
	if s.WarehouseCode == "" {
		return errors.New("仓库不存在")
	}
	skus := utils.Split(s.SkuCode)
	var hasProductSkus []string
	for sku, _ := range productMap {
		hasProductSkus = append(hasProductSkus, sku)
	}

	var noProductSkus []string
	var anomalyPriceSkus []string
	for _, sku := range skus {
		if !utils.InArrayString(sku, hasProductSkus) {
			noProductSkus = append(noProductSkus, sku)
		} else if productMap[sku].MarketPrice <= 0 {
			anomalyPriceSkus = append(anomalyPriceSkus, sku)
		}
	}
	if len(noProductSkus) > 0 {
		return errors.New("产品[" + strings.Join(noProductSkus, ",") + "]不存在或已下架，请先联系相关人员处理!")
	}
	if len(anomalyPriceSkus) > 0 {
		return errors.New("产品[" + strings.Join(anomalyPriceSkus, ",") + "]销售价格异常，请先联系相关人员处理!")
	}

	var categoryObj modelsPc.Category
	hasCategorySku := categoryObj.GetHasCategorySku(tx, skus)
	var noCategorySku []string
	for _, sku := range skus {
		if !utils.InArrayString(sku, hasCategorySku) {
			noCategorySku = append(noCategorySku, sku)
		}
	}
	if len(noCategorySku) > 0 {
		return errors.New("产品[" + strings.Join(noCategorySku, ",") + "]没有分配产线，请先联系相关人员处理!")
	}

	// 税率都相同 在此不做校验

	return
}

type OrderInfoGetResp struct {
	models.OrderInfo
	OrderDetails  []OrderInfoGetDetailResp `json:"orderDetails" gorm:"foreignkey:OrderId;references:OrderId"`
	OrderImages   []models.OrderImage      `json:"orderImages" gorm:"foreignkey:OrderId;references:OrderId"`
	ReceiptImages []models.OrderImage      `json:"receiptImages" gorm:"foreignkey:OrderId;references:OrderId"`
}

type OrderInfoGetDetailResp struct {
	models.OrderDetail
	ActualStock   int     `json:"actualStock" comment:"实缺数量"`
	ReturnNum     int     `json:"returnNum" comment:"已退货数量"`
	AfterNum      int     `json:"afterNum" comment:"售后中数量"`
	VendorName    string  `json:"vendorName" comment:"货主"`
	VendorSkuCode string  `json:"vendorSkuCode" comment:"货主sku"`
	Tax           float64 `json:"tax" comment:"税率"`
	Stock         int     `json:"stock" comment:"可用库存"`
}

type OrderInfoGetPageResp struct {
	models.OrderInfo
	WarehouseName   string `json:"warehouseName"`
	ReturnQuantiry  int    `json:"returnQuantiry" comment:"订单总最终数量" gorm:"-"`
	IsUploadReceipt int    `json:"isUploadReceipt"` //是否回单上传 1是 0否 等于1时签收按钮替换成回单按钮
}

type OrderIdsReq struct {
	OrderIds []string `json:"orderIds" comment:"订单编号"`
}

type OrderListReq struct {
	UserName string `form:"userName" vd:"@:len($)>0; msg:'领用人必填'" comment:"领用人"`
}

type OrderListResp struct {
	UserName string `json:"userName"`
	OrderId  string `json:"orderId"`
}
