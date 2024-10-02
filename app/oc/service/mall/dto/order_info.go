package dto

import (
	"errors"
	"fmt"
	modelsPc "go-admin/app/pc/models"
	"go-admin/common/utils"
	"math/rand"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"go-admin/app/oc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type OrderInfoGetPageReq struct {
	dto.Pagination `search:"-"`
	Keyword        string `form:"keyword"  search:"-"`
	OrderType      int    `form:"orderType"  search:"-"`
	Start          string `form:"startTime" search:"type:gte;column:created_at;table:order_info"`
	End            string `form:"endTime" search:"type:lte;column:created_at;table:order_info"`
	OrderInfoOrder
}

type OrderInfoOrder struct {
	Id string `form:"idOrder"  search:"type:order;column:id;table:order_info"`
}

func (m *OrderInfoGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type OrderInfoUpdatePoReq struct {
	Id         string `uri:"id" comment:""` //
	ContractNo string `json:"contractNo" comment:"客户合同编号"`
}

func (s *OrderInfoUpdatePoReq) Generate(model *models.OrderInfo) {
	model.ContractNo = s.ContractNo
}

func (s *OrderInfoUpdatePoReq) GetId() interface{} {
	return s.Id
}

type OrderInfoInsertReq struct {
	Id              int                             `json:"-" comment:""` //
	UserId          int                             `json:"userId" comment:"客户编号"`
	UserName        string                          `json:"userName" comment:"客户用户名"`
	UserCompanyId   int                             `json:"userCompanyId" comment:"客户公司ID"`
	UserCompanyName string                          `json:"userCompanyName" comment:"客户公司名称"`
	ContractNo      string                          `json:"contractNo" comment:"客户合同编号"`
	Remark          string                          `json:"remark" comment:"客户留言"`
	OrderFileRemark string                          `json:"orderFileRemark" comment:"订单附件备注"`
	OrderImages     []OrderInfoOrderImagesInsertReq `json:"orderImages" comment:"订单相关文件"`
	Products        []OrderInfoProductInsertReq     `json:"products" comment:"商品"`

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

	CreateBy      int    `json:"createBy" gorm:"index;comment:创建者"`
	CreateByName  string `json:"createByName" gorm:"index;comment:创建者姓名"`
	WarehouseCode string `json:"warehouseCode" comment:"发货仓"`

	//Id int `json:"-" comment:""` //
	//OrderId string `json:"orderId" comment:"订单编号"`
	//OrderStatus int `json:"orderStatus" comment:"订单状态：0-未发货、1-已发货、2-部分收货、3-待评价、4-已评价、5-待确认、6、缺货、7-已签收、9-已取消"`
	//ItemsAmount int `json:"itemsAmount" comment:"订单商品总金额(单位:分)"`
	//TotalAmount int `json:"totalAmount" comment:"订单总金额(单位:分)"`
	//CompanyName string `json:"companyName" comment:"收货人公司名称"`
	//Remark string `json:"remark" comment:"客户留言"`
	//ContractNo string `json:"contractNo" comment:"客户合同编号"`
	//UserId int `json:"userId" comment:"客户编号"`
	//UserName string `json:"userName" comment:"客户用户名"`
	//UserCompanyId int `json:"userCompanyId" comment:"客户公司ID"`
	//UserCompanyName string `json:"userCompanyName" comment:"客户公司名称"`
	//Ip string `json:"ip" comment:"客户IP地址"`
	//SaleUserId string `json:"saleUserId" comment:"用户对应的销售用户ID"`
	//SaleUserName string `json:"saleUserName" comment:"用户对应的销售名字"`
	//ProxyUserId string `json:"proxyUserId" comment:"代理下单人"`
	//ProxyUserName string `json:"proxyUserName" comment:"代理下单人姓名"`
	//SendStatus int `json:"sendStatus" comment:"发货状态：0-未发货、1-部分发货、2-全部发货"`
	//ValidFlag int `json:"validFlag" comment:"有效标识位：0-无效、1-有效"`
	//CreateFrom string `json:"createFrom" comment:"订单来源：LMS/MALL/XCX"`
	//ClassifyQuantity int `json:"classifyQuantity" comment:"包含商品种类"`
	//ProductQuantity int `json:"productQuantity" comment:"包含商品数量"`
	//ContactEmail string `json:"contactEmail" comment:"联系人邮箱"`
	//OriginalItemsAmount int `json:"originalItemsAmount" comment:"订单商品初始总金额(单位:分)"`
	//OriginalTotalAmount int `json:"originalTotalAmount" comment:"订单初始总金额(单位:分)"`
	//CountryId int `json:"countryId" comment:"客户国家ID"`
	//CountryName string `json:"countryName" comment:"客户国家名称"`
	//ProvinceId int `json:"provinceId" comment:"客户省份ID"`
	//ProvinceName string `json:"provinceName" comment:"客户省份名称"`
	//DistrictId int `json:"districtId" comment:"客户区县ID"`
	//DistrictName string `json:"districtName" comment:"客户区县名称"`
	//TownId int `json:"townId" comment:"镇/街道 ID"`
	//TownName string `json:"townName" comment:"镇/街道名称"`
	//ConfirmTime time.Time `json:"confirmTime" comment:"订单确认时间"`
	//PartialShipments int `json:"partialShipments" comment:"订单是否分批发货,1-分批发货，2-整单发货"`
	//ActualOrderingCompany string `json:"actualOrderingCompany" comment:"实际下单公司"`
	//ActualOrderingCompanyId int `json:"actualOrderingCompanyId" comment:"实际下单公司"`
	//ConfirmOrderReceiptTime time.Time `json:"confirmOrderReceiptTime" comment:"订单确认签收时间（客户发起签收）"`
	//FinalTotalAmount string `json:"finalTotalAmount" comment:"最终订单总金额（对账单生成时各明细FINAL_SUB_TOTAL_AMOUNT之和）"`
	//IsOverBudget int `json:"isOverBudget" comment:"是否超出预算"`
	//OrderFileRemark string `json:"orderFileRemark" comment:"订单附件备注"`
	//LogisticalRemark string `json:"logisticalRemark" comment:"物流备注"`
	//CancelBy int `json:"cancelBy" comment:"取消人id"`
	//CancelByName string `json:"cancelByName" comment:"取消人姓名"`
	//CancelByType int `json:"cancelByType" comment:"取消人类型：0-客户、1-销售"`
	//common.ControlBy
}

type OrderInfoOrderImagesInsertReq struct {
	Name string `json:"name" comment:"图片名称"`
	Url  string `json:"url" comment:"url"`
}

type OrderInfoProductInsertReq struct {
	SkuCode        string  `json:"skuCode" comment:"sku"`
	ProductName    string  `json:"productName" comment:"商品名称"`
	ProductModel   string  `json:"productModel" comment:"商品型号"`
	BrandId        int     `json:"brandId" comment:"商品品牌ID"`
	BrandName      string  `json:"brandName" comment:"商品品牌名字"`
	BrandEname     string  `json:"brandEname" comment:"商品品牌英文"`
	ProductId      int     `json:"productId" comment:"comment:商品ID"`
	Quantity       int     `json:"quantity" comment:"商品数量"`
	Moq            int     `json:"moq" comment:"销售最小起订量"`
	Unit           string  `json:"unit" comment:"商品单位"`
	SalePrice      float64 `json:"salePrice" comment:"商品销售价格"`
	VendorId       int     `json:"vendorId" comment:"供应商信息ID"`
	VendorName     string  `json:"vendorName" comment:"货主名称"`
	VendorSkuCode  string  `json:"vendorSkuCode" comment:"货主sku"`
	ProductNo      string  `json:"productNo" comment:"物料编码"`
	GoodsId        int     `json:"goodsId"comment:"goods表ID"`
	Stock          int     `json:"stock" comment:"可用库存"`
	Tax            string  `json:"tax" comment:"税率"`
	SubTotalAmount float64 `json:"subTotalAmount" comment:"行项目小计金额"`

	//models.Model
	//
	//OrderId string `json:"orderId" gorm:"type:varchar(30);comment:订单编号"`
	//SendQuantity int `json:"sendQuantity" gorm:"type:int unsigned;comment:已发货商品数量"`
	//UserId int `json:"userId" gorm:"type:int unsigned;comment:用户ID"`
	//UserName string `json:"userName" gorm:"type:varchar(200);comment:用户名称"`
	//CatId int `json:"catId" gorm:"type:int unsigned;comment:一级产线ID"`
	//CatName string `json:"catName" gorm:"type:varchar(50);comment:一级产线名称"`
	//CatId2 int `json:"catId2" gorm:"type:int unsigned;comment:二级产线ID"`
	//CatName2 string `json:"catName2" gorm:"type:varchar(50);comment:二级产线名称"`
	//CatId3 int `json:"catId3" gorm:"type:int unsigned;comment:三级产线ID"`
	//CatName3 string `json:"catName3" gorm:"type:varchar(50);comment:三级产线名称"`
	//CatId4 int `json:"catId4" gorm:"type:int unsigned;comment:四级产线ID"`
	//CatName4 string `json:"catName4" gorm:"type:varchar(50);comment:四级产线名称"`
	//SkuCode string `json:"skuCode" gorm:"type:varchar(20);comment:sku"`
	//ProductId int `json:"productId" gorm:"type:int unsigned;comment:商品ID"`
	//ProductPic string `json:"productPic" gorm:"type:varchar(200);comment:商品名称"`
	//ProductModel string `json:"productModel" gorm:"type:varchar(255);comment:商品型号"`
	//PurchasePrice int `json:"purchasePrice" gorm:"type:int unsigned;comment:商品采购价(单位:分)"`
	//LockStock int `json:"lockStock" gorm:"type:int unsigned;comment:已锁库存"`
	//SubTotalAmount int `json:"subTotalAmount" gorm:"type:int unsigned;comment:行项目小计金额(单位:分)"`
	//OriginalItemsMount int `json:"originalItemsMount" gorm:"type:int unsigned;comment:原商品总金额(单位:分)"`
	//OriginalQuantity int `json:"originalQuantity" gorm:"type:int unsigned;comment:原商品数量(单位:分)"`
	//BatchGroup int `json:"batchGroup" gorm:"type:tinyint unsigned;comment:发货分组"`
	//SupplierWarehouse string `json:"supplierWarehouse" gorm:"type:varchar(10);comment:供应商仓库编号"`
	//DeliveryWarehouse string `json:"deliveryWarehouse" gorm:"type:varchar(10);comment:西域发货仓"`
	//MarketPrice int `json:"marketPrice" gorm:"type:int unsigned;comment:系统价格(单位:分)"`
	//Moq int `json:"moq" gorm:"type:int unsigned;comment:最小起订量"`
	//UserProductName string `json:"userProductName" gorm:"type:varchar(255);comment:客户商品物料描述"`
	//UserProductRemark string `json:"userProductRemark" gorm:"type:varchar(255);comment:客户商品物料采购备注"`
	//FinalSubTotalAmount int `json:"finalSubTotalAmount" gorm:"type:int unsigned;comment:最终行项目小计金额(单位:分)"`
	//FinalQuantity int `json:"finalQuantity" gorm:"type:int unsigned;comment:最终数量（总数量-已取消-已退货）"`
	//CreatedAt time.Time      `json:"createdAt" gorm:"comment:创建时间"`
	//UpdatedAt time.Time      `json:"updatedAt" gorm:"comment:最后更新时间"`
}

func (s *OrderInfoInsertReq) Generate(model *models.OrderInfo) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.OrderId = GenerateOrderId()
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

	//model.OrderId = s.OrderId
	//model.OrderStatus = s.OrderStatus
	//model.ItemsAmount = s.ItemsAmount
	//model.TotalAmount = s.TotalAmount
	//model.AreaId = s.AreaId
	//model.AreaName = s.AreaName
	//model.CompanyName = s.CompanyName
	//model.Remark = s.Remark
	//model.ContractNo = s.ContractNo
	//model.Ip = s.Ip
	//model.SaleUserId = s.SaleUserId
	//model.SaleUserName = s.SaleUserName
	//model.ProxyUserId = s.ProxyUserId
	//model.ProxyUserName = s.ProxyUserName
	//model.SendStatus = s.SendStatus
	//model.ValidFlag = s.ValidFlag
	//model.CreateFrom = s.CreateFrom
	//model.ClassifyQuantity = s.ClassifyQuantity
	//model.ProductQuantity = s.ProductQuantity
	//model.ContactEmail = s.ContactEmail
	//model.OriginalItemsAmount = s.OriginalItemsAmount
	//model.OriginalTotalAmount = s.OriginalTotalAmount
	//model.CountryId = s.CountryId
	//model.CountryName = s.CountryName
	//model.ProvinceId = s.ProvinceId
	//model.ProvinceName = s.ProvinceName
	//model.DistrictId = s.DistrictId
	//model.DistrictName = s.DistrictName
	//model.TownId = s.TownId
	//model.TownName = s.TownName
	//model.ConfirmTime = s.ConfirmTime
	//model.PartialShipments = s.PartialShipments
	//model.ActualOrderingCompany = s.ActualOrderingCompany
	//model.ActualOrderingCompanyId = s.ActualOrderingCompanyId
	//model.ConfirmOrderReceiptTime = s.ConfirmOrderReceiptTime
	//model.FinalTotalAmount = s.FinalTotalAmount
	//model.IsOverBudget = s.IsOverBudget
	//model.OrderFileRemark = s.OrderFileRemark
	//model.LogisticalRemark = s.LogisticalRemark
	//model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	//model.CreateByName = s.CreateByName
	//model.CancelBy = s.CancelBy
	//model.CancelByName = s.CancelByName
	//model.CancelByType = s.CancelByType
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
	var limitSaleMoqSkus []string
	for _, product := range s.Products {
		if product.Quantity < product.Moq {
			limitSaleMoqSkus = append(limitSaleMoqSkus, product.SkuCode)
		}
	}
	if len(limitSaleMoqSkus) > 0 {
		return errors.New("产品[" + strings.Join(limitSaleMoqSkus, ",") + "]未达到最小起订量!")
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

func GenerateOrderId() string {
	return "RO" + time.Now().Format("20060102150405") + fmt.Sprintf("%04d", (1+rand.Intn(9999)))
}

func GetPageMakeCondition(c *OrderInfoGetPageReq, newDb *gorm.DB) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db.Where("order_info.valid_flag = 1 and order_info.is_formal_order = 1 and order_info.user_id = ?", user.GetUserId(db.Statement.Context.(*gin.Context)))
		//db.Where("order_info.valid_flag = 1 and order_info.is_formal_order = 1 ")
		if c.Keyword != "" {
			db.Where("order_info.order_id like ? or order_info.external_order_no like ? or od.sku_code like ? or od.product_name like ? or od.product_no like ?", "%"+c.Keyword+"%", "%"+c.Keyword+"%", "%"+c.Keyword+"%", "%"+c.Keyword+"%", "%"+c.Keyword+"%")
		}
		if c.OrderType > 0 {
			switch c.OrderType {
			case 5: // 小程序-仓库处理中 mall-待发货
				db.Where("order_info.order_status in (0,5,6)")
				break
			case 6: // mall-待发货(这个状态前端已经不用了 统一都用5)
				db.Where("order_info.order_status = 0")
				break
			case 1: // 待收货
				db.Where("order_info.order_status in (1,11)")
				break
			case 7: // 已完成
				db.Where("order_info.order_status = 7")
				break
			case 9: // 已取消
				db.Where("order_info.order_status in (9, 10)")
				break
			case 99: //全部
				break
			}
		}
		return db
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

	return
}

type OrderInfoReceiptReq struct {
	Id string `uri:"id" comment:""` //
	//ReceiptImages []OrderInfoOrderImagesInsertReq `json:"receiptImages" comment:"订单签收文件"`
}

func (s *OrderInfoReceiptReq) GetId() interface{} {
	return s.Id
}

func (s *OrderInfoReceiptReq) ReceiptValid(orderInfo *models.OrderInfo) (err error) {
	if orderInfo.Id <= 0 {
		return errors.New("此订单不存在！")
	}
	if orderInfo.OrderStatus == 1 && (orderInfo.RmaStatus == 1 || orderInfo.RmaStatus == 2) {
		return errors.New("此订单正处于售后中！")
	}

	return
}

type OrderInfoCancelReq struct {
	Id           string `uri:"id" comment:""` //
	Remark       string `json:"remark" comment:"取消原因"`
	CancelBy     int    `json:"cancelBy" comment:"取消人id"`
	CancelByName string `json:"cancelByName" comment:"取消人姓名"`
}

func (s *OrderInfoCancelReq) Generate(model *models.OrderInfo) {
	model.Remark = s.Remark
	model.CancelBy = s.CancelBy
	model.CancelByName = s.CancelByName
}

func (s *OrderInfoCancelReq) GetId() interface{} {
	return s.Id
}

func (s *OrderInfoCancelReq) CancelValid(tx *gorm.DB) (err error) {
	var object models.OrderInfo
	orderInfo := object.GetRowByOrderId(tx, s.Id)
	if orderInfo.OrderStatus != 5 && orderInfo.OrderStatus != 6 {
		return errors.New("此订单不可取消")
	}
	if orderInfo.RmaStatus != 0 && orderInfo.RmaStatus != 99 {
		return errors.New("订单正在售后中，不可取消")
	}

	return
}

type OrderInfoBuyAgainReq struct {
	Id string `uri:"id" comment:""` //
}

// OrderInfoGetReq 功能获取请求参数
type OrderInfoGetReq struct {
	Id string `uri:"id"`
}

func (s *OrderInfoGetReq) GetId() interface{} {
	return s.Id
}

type OrderInfoDeleteReq struct {
	Id string `uri:"id" comment:""` //
}

func (s *OrderInfoDeleteReq) GetId() interface{} {
	return s.Id
}

func (s *OrderInfoDeleteReq) Valid(orderInfo models.OrderInfo) (err error) {
	if orderInfo.OrderStatus != 9 {
		return errors.New("该订单不允许删除")
	}
	if orderInfo.ValidFlag == 0 {
		return errors.New("此订单已被删除")
	}

	return
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

	var object models.OrderInfo
	orderInfo := object.GetRow(tx, s.Id)
	if orderInfo.OrderStatus != 5 && orderInfo.OrderStatus != 6 {
		return errors.New("该订单不可编辑")
	}

	return
}

type OrderInfoByOrderIdReq struct {
	OrderIds []string `form:"orderIds"  search:"type:in;column:order_id;table:order_info" vd:"@:len($[0])>0 ; msg:'订单id必填'"`
}

func (m *OrderInfoByOrderIdReq) GetNeedSearch() interface{} {
	return *m
}

type OrderInfoAddProductReq struct {
	SkuCode       string `json:"skuCode" comment:"SKU（多个逗号分割）"`
	WarehouseCode string `json:"warehouseCode" comment:"发货仓库code"`
}

type OrderInfoAddProductResp struct {
	SkuCode       string `json:"skuCode" comment:"SKU（多个逗号分割）"`
	WarehouseCode string `json:"warehouseCode" comment:"发货仓库code"`
}

func (s *OrderInfoAddProductReq) Valid(tx *gorm.DB, products *[]modelsPc.Goods) (err error) {
	if s.WarehouseCode == "" {
		return errors.New("仓库不存在")
	}
	skus := utils.Split(s.SkuCode)
	productMap := make(map[string]modelsPc.Goods, len(*products))
	_ = utils.StructColumn(&productMap, *products, "", "SkuCode")
	var hasProductSkus []string
	_ = utils.StructColumn(&hasProductSkus, *products, "SkuCode", "")

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
	OrderDetails    []OrderInfoGetDetailResp `json:"orderDetails" gorm:"foreignkey:OrderId;references:OrderId"`
	OrderImages     []models.OrderImage      `json:"orderImages" gorm:"foreignkey:OrderId;references:OrderId"`
	ReceiptImages   []models.OrderImage      `json:"receiptImages" gorm:"foreignkey:OrderId;references:OrderId"`
	AddressFullName string                   `json:"addressFullName" comment:"完整地址" gorm:"-"`
	ButtonList      []int                    `json:"buttonList" comment:"展示列表按钮" gorm:"-"`
}

type OrderInfoGetDetailResp struct {
	models.OrderDetail
	ActualStock int    `json:"actualStock" comment:"实缺数量"`
	ReturnNum   int    `json:"returnNum" comment:"已退货数量"`
	AfterNum    int    `json:"afterNum" comment:"售后中数量"`
	VendorName  string `json:"vendorName" comment:"货主"`
}

type OrderInfoGetPageResp struct {
	models.OrderInfo
	WarehouseName string                            `json:"warehouseName" comment:"发货仓"`
	OrderDetails  []OrderInfoGetPageRespOrderDetail `json:"orderDetails" gorm:"foreignkey:OrderId;references:OrderId"`
	ButtonList    []int                             `json:"buttonList" comment:"展示列表按钮" gorm:"-"`
}

type OrderInfoGetPageRespOrderDetail struct {
	models.OrderDetail
	UnTaxPrice float64 `json:"unTaxPrice"`
	VendorName string  `json:"vendorName" comment:"货主"`
}

type OrderIdsReq struct {
	OrderIds []string `json:"orderIds" comment:"订单编号"`
}

type OrderInfoGetExportData struct {
	OrderId           string  `json:"orderId" comment:"订单编号"`
	ExternalOrderNo   string  `json:"externalOrderNo" comment:"审批单号"`
	ContractNo        string  `json:"contractNo" comment:"客户合同编号"`
	UserCompanyName   string  `json:"userCompanyName" comment:"客户公司名称"`
	UserName          string  `json:"userName" comment:"客户用户名"`
	CreatedTime       string  `json:"createdTime" comment:"下单时间"`
	WarehouseName     string  `json:"warehouseName" comment:"发货仓"`
	Consignee         string  `json:"consignee" comment:"收货人姓名"`
	SendStatusText    string  `json:"sendStatusText"`
	ReceiveStatusText string  `json:"receiveStatusText"`
	ProductName       string  `json:"productName" comment:"商品名称"`
	UserProductRemark string  `json:"userProductRemark" comment:"客户商品物料采购备注"`
	SkuCode           string  `json:"skuCode" comment:"sku"`
	VendorName        string  `json:"vendorName" comment:"货主名称"`
	SupplierSkuCode   string  `json:"supplierSkuCode" comment:"货主sku"`
	ProductNo         string  `json:"productNo" comment:"物料编码"`
	BrandName         string  `json:"brandName" comment:"商品品牌名字"`
	ProductModel      string  `json:"productModel" comment:"商品型号"`
	CatName           string  `json:"catName" comment:"一级产线名称"`
	CatName2          string  `json:"catName2" comment:"二级产线名称"`
	Unit              string  `json:"unit" comment:"商品单位"`
	SalePrice         float64 `json:"salePrice" comment:"含税单价"`
	Quantity          int
	TotalPrice        float64 `json:"totalPrice" comment:"含税总单价"`
	CancelQuantity    float64 `json:"cancelQuantity"`
	ReturnQuantity    int     `json:"returnQuantity"`
}
