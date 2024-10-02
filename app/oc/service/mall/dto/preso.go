package dto

import (
	"errors"
	"fmt"
	"go-admin/app/oc/models"
	modelsPc "go-admin/app/pc/models"
	modelsUc "go-admin/app/uc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
	"go-admin/common/utils"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	"gorm.io/gorm"
)

type PresoImagesInsertReq struct {
	Name string `json:"name" comment:"图片名称"`
	Url  string `json:"url" comment:"url"`
}
type PresoUserProductRemarkReq struct {
	SkuCode string `json:"skuCode" comment:"sku"`
	Remark  string `json:"remark" comment:"备注"`
}

func (s *PresoImagesInsertReq) Generate(model *models.PresoImage) {
	model.Name = s.Name
	model.Url = s.Url
}

// 提交审批
type PresoSubmitApprovalReq struct {
	ShippingAddressId int                         `json:"shippingAddressId" comment:"收货地址id"`
	ApproveflowId     int                         `json:"approveflowId" comment:"审批流id"`
	ContractNo        string                      `json:"contractNo" comment:"Po单号"`
	Remark            string                      `json:"remark" comment:"备注"`
	CreateFrom        string                      `json:"createFrom" comment:"订单来源"`
	BuyNow            int                         `json:"buyNow" comment:"是否立即领用"`
	GoodsId           int                         `json:"goodsId" comment:"goods表ID"`
	Quantity          int                         `json:"quantity" comment:"商品数量"`
	ApproveRemark     string                      `json:"approveRemark" comment:"审批说明"`
	PresoImage        []PresoImagesInsertReq      `json:"presoImage" comment:"预订单附件"`
	UserProductRemark []PresoUserProductRemarkReq `json:"userProductRemark" comment:"SKU备注"`
}

func (s *PresoSubmitApprovalReq) Generate(model *models.Preso) {
	model.PresoNo = generatePresoNo()
	model.ApproveflowId = s.ApproveflowId
	model.Remark = s.Remark
	model.ApproveRemark = s.ApproveRemark
	model.ContractNo = s.ContractNo
	model.DeliverId = s.ShippingAddressId
	model.CreateFrom = s.CreateFrom
}

func generatePresoNo() string {
	return "PR" + time.Now().Format("20060102150405") + fmt.Sprintf("%04d", (1+rand.Intn(9999)))
}

func (s *PresoSubmitApprovalReq) Valid(tx *gorm.DB, shippingAddress *modelsUc.Address, productList []modelsPc.UserCartGoodsProduct, companyInfo modelsUc.CompanyInfo) (err error) {
	if s.ShippingAddressId <= 0 {
		return errors.New("请选择收货地址！")
	}
	if companyInfo.Id <= 0 {
		return errors.New("用户所属公司不存在")
	}
	err = shippingAddress.Get(tx, s.ShippingAddressId)
	if err != nil {
		return errors.New("请选择收货地址！")
	}
	if shippingAddress.ProvinceId <= 0 || shippingAddress.CityId <= 0 || shippingAddress.AreaId <= 0 {
		return errors.New("收货地址省市区异常，请核对后下单！")
	}
	var salesMoqLimitSkus []string
	var unSelectedSkus []string
	var unSaleSkus []string
	totalAmount := 0.00
	for _, product := range productList {
		if companyInfo.CheckStockStatus == 0 && product.Stock < product.Quantity {
			return errors.New("库存不足，不允许下单！")
		}
		if product.Quantity < product.SalesMoq {
			salesMoqLimitSkus = append(salesMoqLimitSkus, product.SkuCode)
		}
		// 购物车中是否被选中
		if product.Selected != 1 {
			unSelectedSkus = append(unSelectedSkus, product.SkuCode)
		} else if product.SaleStatus != 1 {
			unSaleSkus = append(unSaleSkus, product.SkuCode)
		}
		totalAmount = utils.AddFloat64(totalAmount, utils.MulFloat64AndInt(product.MarketPrice, product.Quantity))
	}
	if len(salesMoqLimitSkus) > 0 {
		return errors.New("存在" + strconv.Itoa(len(salesMoqLimitSkus)) + "类商品[" + strings.Join(salesMoqLimitSkus, ",") + "]小于最小订货量")
	}
	if len(unSelectedSkus) > 0 {
		return errors.New("请检查购物车商品")
	}
	if len(unSaleSkus) > 0 {
		return errors.New("请检查购物车是否还存在可售商品")
	}

	// 订单总金额不能小于0
	if totalAmount <= 0 {
		return errors.New("订单应付金额异常")
	}

	return
}

type PresoGetPageReq struct {
	dto.Pagination `search:"-"`
	Keyword        string `form:"keyword"  search:"-"`
	ApproveStatus  int    `form:"approveStatus"  search:"-"`
	Type           int    `form:"type"  search:"-"`
	Client         string `form:"client"  search:"-"`
	PresoOrder
}

type PresoOrder struct {
	Id string `form:"idOrder"  search:"type:order;column:id;table:preso"`
}

func (m *PresoGetPageReq) GetNeedSearch() interface{} {
	return *m
}

func PresoGetPageMakeCondition(c *PresoGetPageReq, newDb *gorm.DB) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		userId := user.GetUserId(db.Statement.Context.(*gin.Context))
		if c.Keyword != "" {
			db.Where("preso.preso_no like ? or preso.user_name like ? or pd.sku_code like ? or pd.product_name like ? or pd.product_no like ?", "%"+c.Keyword+"%", "%"+c.Keyword+"%", "%"+c.Keyword+"%", "%"+c.Keyword+"%", "%"+c.Keyword+"%")
		}
		if c.ApproveStatus != 0 {
			switch c.ApproveStatus {
			case 10:
				db.Where("preso.approve_status in (0, 10)")
				break
			case -1, -2, -3:
				db.Where("preso.approve_status = ?", c.ApproveStatus)
				break
			}
		}
		// 领用申请
		if c.Type == 1 {
			db.Where("preso.user_id = ?", userId)
			if c.ApproveStatus == 1 {
				db.Where("preso.approve_status = 1")
			}
		} else if c.Type == 2 { //领用审批
			db.Where("pl.user_id = ?", userId)
			// 待审批 只展示下个节点是当前登陆人的
			if c.ApproveStatus == 10 {
				//// 待审批 下个节点是当前登陆人的排前面
				//db.Order("CASE WHEN preso.step+1 = pl.step THEN 1 ELSE 2 END")
				db.Where("preso.step+1 = pl.step")
			} else if c.ApproveStatus == 1 {
				// 已审批 当前登陆人审批过的
				db.Where("preso.approve_status in (0, 10, 1)")
				db.Where("pl.approve_status = 1 and not exists (select t.id from preso t left join preso_log pl on t.preso_no = pl.preso_no where pl.user_id = ? and preso.step+1 = pl.step)", userId)
			} else if c.ApproveStatus == 0 {
				// 已审批 当前登陆人审批过的
				db.Where("(preso.approve_status in (0, 10) and (preso.step+1 = pl.step or pl.approve_status = 1)) or preso.approve_status in (-2, -1, 1)")
			}
		}
		return db
	}
}

type WorkFlowNodes struct {
	Step   int    `json:"step" comment:"步骤"`
	Text   string `json:"text" comment:"内容"`
	Status int    `json:"status" comment:"状态 0-未审核 1-已审核"`
}

type PresoGetPageResp struct {
	models.Preso
	PresoDetails        []PresoGetPageRespPresoDetail `json:"presoDetails" gorm:"foreignkey:PresoNo;references:PresoNo"`
	TotalAmount         float64                       `json:"totalAmount" comment:"含税总金额"`
	TotalQuantity       int                           `json:"totalQuantity" comment:"总数量"`
	UnRejectTotalAmount float64                       `json:"unRejectTotalAmount" comment:"未驳回的含税总金额"`
	TaxTotalAmount      float64                       `json:"taxTotalAmount" comment:"税额"`
	ExpireTimeText      string                        `json:"expireTimeText" comment:"过期时间描述"`
	WorkFollowNodes     []models.WorkFlowNodes        `json:"workFollowNodes" gorm:"-"`
	OrderId             string                        `json:"orderId" comment:"订单编号" gorm:"-"`
	ButtonList          []int                         `json:"buttonList" comment:"展示列表按钮" gorm:"-"`
	ApproveUser         string                        `json:"approveUser" comment:"订单编号" gorm:"-"`
}

type PresoGetPageRespPresoDetail struct {
	models.PresoDetail
	NakedUnitPrice float64 `json:"nakedUnitPrice" comment:"未税价"`
	BrandName      string  `json:"brandName" comment:"品牌名称"`
	MfgModel       string  `json:"mfgModel" comment:"型号"`
	VendorName     string  `json:"vendorName" comment:"货主"`
}

type PresoBuyAgainReq struct {
	Id string `uri:"id" comment:""` //
}

type PresoDeleteFleReq struct {
	Id int `json:"id"`
}

func (s *PresoDeleteFleReq) GetId() interface{} {
	return s.Id
}

type PresoGetRespPresoDetail struct {
	models.PresoDetail
	NakedUnitPrice float64 `json:"nakedUnitPrice" comment:"未税价"`
	UntaxedTotal   float64 `json:"untaxedTotal" comment:"未税总计"`
	TaxedTotal     float64 `json:"taxedTotal" comment:"含税总计"`
	Tax            float64 `json:"tax" comment:"税额"`
	OperUser       string  `json:"operUser" comment:"驳回人"`
	ApproveRemark  string  `json:"approveRemark" comment:"审批说明"`
	Unit           string  `json:"unit" comment:"单位"`
	BrandName      string  `json:"brandName" comment:"品牌名称"`
	MfgModel       string  `json:"mfgModel" comment:"型号"`
	VendorName     string  `json:"vendorName" comment:"货主"`
	VendorSkuCode  string  `json:"vendorSkuCode" comment:"货主SkuCode"`
}

type PresoGetResp struct {
	models.Preso
	PresoDetails          []PresoGetRespPresoDetail `json:"presoDetails" gorm:"foreignkey:PresoNo;references:PresoNo"`
	TotalAmount           float64                   `json:"totalAmount" comment:"含税总金额" gorm:"-"`
	ExpireTimeText        string                    `json:"expireTimeText" comment:"过期时间描述" gorm:"-"`
	AddressFullName       string                    `json:"addressFullName" comment:"完整地址" gorm:"-"`
	WorkFollowNodes       []models.WorkFlowNodes    `json:"workFollowNodes" gorm:"-"`
	ApproveRemarkText     string                    `json:"approveRemarkText" comment:"审批描述" gorm:"-"`
	IsCurrentNodeApprover int                       `json:"isCurrentNodeApprover" comment:"是否当前节点的审批人" gorm:"-"`
	OrderId               string                    `json:"orderId" comment:"订单编号" gorm:"-"`
}

type PresoGetExportResp struct {
	SkuCode           string  `json:"skuCode" comment:"sku"`
	ProductName       string  `json:"productName" comment:"商品名称"`
	VendorName        string  `json:"vendorName" comment:"货主"`
	SupplierSkuCode   string  `json:"supplierSkuCode" comment:"货主sku"`
	UserProductRemark string  `json:"userProductRemark" comment:"SKU备注"`
	Quantity          int     `json:"quantity" comment:":商品数量"`
	Unit              string  `json:"unit" comment:"单位"`
	NakedUnitPrice    float64 `json:"nakedUnitPrice" comment:"未税价"`
	SalePrice         float64 `json:"salePrice" comment:"商品销售价格"`
	Tax               float64 `json:"tax" comment:"税额"`
	UntaxedTotal      float64 `json:"untaxedTotal" comment:"未税总计"`
	TaxedTotal        float64 `json:"taxedTotal" comment:"含税总计"`
	ApproveRemark     string  `json:"approveRemark" comment:"备注"`
}

type PresoFinishApprovalProductReq struct {
	SkuCode          string  `json:"skuCode" comment:"sku"`
	Approved         int     `json:"approved" comment:"状态 0-待审核 1-已通过 -1'-已驳回"`
	Quantity         int     `json:"quantity" comment:":商品数量"`
	ApprovedQuantity int     `json:"approvedQuantity" comment:":通过数量"`
	Price            float64 `json:"price" comment:"价格"`
}

type PresoFinishApprovalOperContent struct {
	PassTotal    int                             `json:"passTotal" comment:"通过数量"`
	RejectTotal  int                             `json:"rejectTotal" comment:":驳回数量"`
	ApproveItems []PresoFinishApprovalProductReq `json:"approveItems"`
}

// 完成审批
type PresoFinishApprovalReq struct {
	PresoNo       string `json:"presoNo" comment:"审批单编号"`
	Step          int    `json:"step" comment:"审批第几级"`
	ApproveRemark string `json:"approveRemark" comment:"审批意见"`
	ContractNo    string `json:"contractNo" comment:"PO单号"`
	Remark        string `json:"remark" comment:"采购备注"`
	PresoFinishApprovalOperContent
}

func (s *PresoFinishApprovalReq) Valid(tx *gorm.DB, preso *models.Preso, companyInfo modelsUc.CompanyInfo) (err error) {
	if s.PresoNo == "" {
		return errors.New("审批单号有误！")
	}
	err = tx.Model(&preso).Preload("PresoDetails").Where("preso_no = ?", s.PresoNo).First(&preso).Error
	if err != nil {
		return errors.New("审批单不存在！")
	}
	if preso.ApproveStatus != 0 && preso.ApproveStatus != 10 {
		return errors.New("此审批单已审批完成！")
	}
	if preso.Step >= s.Step {
		return errors.New("此流程已审批！")
	}
	if len(s.ApproveItems) <= 0 {
		return errors.New("审批的商品不能为空！")
	}

	// TODO 审批校验逻辑

	for _, approveItem := range s.ApproveItems {
		if approveItem.Approved != 1 && approveItem.Approved != -1 {
			return errors.New("商品中有审批未完成！")
		}
	}

	return
}

// 批量审批
type PresoBatchApprovalReq struct {
	PresoNos      []string `json:"presoNos" comment:"审批单编号"`
	ApproveStatus int      `json:"approveStatus" comment:"审批状态 -1 审批不通过  1 审批通过"`
	ApproveRemark string   `json:"approveRemark" comment:"审批意见"`
}

func (s *PresoBatchApprovalReq) Valid(tx *gorm.DB, presos *[]models.Preso) (err error) {
	if len(s.PresoNos) <= 0 {
		err = errors.New("请选择审批单！")
		return
	}
	if s.ApproveStatus != -1 && s.ApproveStatus != 1 {
		err = errors.New("审批状态有误！")
		return
	}
	err = tx.Preload("PresoDetails").Where("preso_no in ?", s.PresoNos).Find(&presos).Error
	if err != nil {
		err = errors.New("审批单不存在！")
		return
	}

	return
}

// 撤回审批单
type PresoWithdrawReq struct {
	PresoNo string `json:"presoNo" comment:"审批单编号"`
}

// 保存上传文档
type PresoSaveFileReq struct {
	PresoNo string `json:"presoNo" comment:"审批单编号"`
	Name    string `json:"name" comment:"图片名称"`
	Url     string `json:"url" comment:"url"`
}

type PresoInsertReq struct {
	Id              int       `json:"-" comment:""` //
	PresoNo         string    `json:"presoNo" comment:"审批单编号"`
	WarehouseCode   string    `json:"warehouseCode" comment:"发货仓"`
	ApproveflowId   int       `json:"approveflowId" comment:"审批流id"`
	ApproveUsers    string    `json:"approveUsers" comment:"审批流用户"`
	UserId          int       `json:"userId" comment:"客户编号"`
	UserName        string    `json:"userName" comment:"客户用户名"`
	UserCompanyId   int       `json:"userCompanyId" comment:"客户公司ID"`
	UserCompanyName string    `json:"userCompanyName" comment:"客户公司名称"`
	DeliverId       int       `json:"deliverId" comment:"收货地址ID"`
	Consignee       string    `json:"consignee" comment:"收货人姓名"`
	CountryId       int       `json:"countryId" comment:"客户国家ID"`
	CountryName     string    `json:"countryName" comment:"客户国家名称"`
	ProvinceId      int       `json:"provinceId" comment:"客户省份ID"`
	ProvinceName    string    `json:"provinceName" comment:"客户省份名称"`
	CityId          int       `json:"cityId" comment:"收货人城市编号"`
	CityName        string    `json:"cityName" comment:"收货人城市名称"`
	AreaId          int       `json:"areaId" comment:"客户区县ID"`
	AreaName        string    `json:"areaName" comment:"客户区县名称"`
	TownId          int       `json:"townId" comment:"镇/街道 ID"`
	TownName        string    `json:"townName" comment:"镇/街道名称"`
	CompanyName     string    `json:"companyName" comment:"收货人公司名称"`
	Address         string    `json:"address" comment:"收货人详细地址"`
	Mobile          string    `json:"mobile" comment:"收货人手机号"`
	Telephone       string    `json:"telephone" comment:"收货人座机号"`
	ContractNo      string    `json:"contractNo" comment:"客户合同编号"`
	CreateFrom      string    `json:"createFrom" comment:"订单来源：LMS/MALL/XCX"`
	Remark          string    `json:"remark" comment:"客户留言"`
	Ip              string    `json:"ip" comment:"客户IP地址"`
	ApproveStatus   int       `json:"approveStatus" comment:"审批状态 -1 审批不通过  1 审批通过 0 初始提交审批 10 审批中 -2 超时 -3 撤回"`
	Step            int       `json:"step" comment:"审批到第几级"`
	ExpireTime      time.Time `json:"expireTime" comment:"过期时间"`
	common.ControlBy
}

func (s *PresoInsertReq) Generate(model *models.Preso) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.PresoNo = s.PresoNo
	model.WarehouseCode = s.WarehouseCode
	model.ApproveflowId = s.ApproveflowId
	model.ApproveUsers = s.ApproveUsers
	model.UserId = s.UserId
	model.UserName = s.UserName
	model.UserCompanyId = s.UserCompanyId
	model.UserCompanyName = s.UserCompanyName
	model.DeliverId = s.DeliverId
	model.Consignee = s.Consignee
	model.CountryId = s.CountryId
	model.CountryName = s.CountryName
	model.ProvinceId = s.ProvinceId
	model.ProvinceName = s.ProvinceName
	model.AreaId = s.AreaId
	model.AreaName = s.AreaName
	model.CityId = s.CityId
	model.CityName = s.CityName
	model.TownId = s.TownId
	model.TownName = s.TownName
	model.CompanyName = s.CompanyName
	model.Address = s.Address
	model.Mobile = s.Mobile
	model.Telephone = s.Telephone
	model.ContractNo = s.ContractNo
	model.CreateFrom = s.CreateFrom
	model.Remark = s.Remark
	model.Ip = s.Ip
	model.ApproveStatus = s.ApproveStatus
	model.Step = s.Step
	model.ExpireTime = s.ExpireTime
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.CreateByName = s.CreateByName
}

func (s *PresoInsertReq) GetId() interface{} {
	return s.Id
}

type PresoUpdateReq struct {
	Id              int       `uri:"id" comment:""` //
	PresoNo         string    `json:"presoNo" comment:"审批单编号"`
	WarehouseCode   string    `json:"warehouseCode" comment:"发货仓"`
	ApproveflowId   int       `json:"approveflowId" comment:"审批流id"`
	ApproveUsers    string    `json:"approveUsers" comment:"审批流用户"`
	UserId          int       `json:"userId" comment:"客户编号"`
	UserName        string    `json:"userName" comment:"客户用户名"`
	UserCompanyId   int       `json:"userCompanyId" comment:"客户公司ID"`
	UserCompanyName string    `json:"userCompanyName" comment:"客户公司名称"`
	DeliverId       int       `json:"deliverId" comment:"收货地址ID"`
	Consignee       string    `json:"consignee" comment:"收货人姓名"`
	CountryId       int       `json:"countryId" comment:"客户国家ID"`
	CountryName     string    `json:"countryName" comment:"客户国家名称"`
	ProvinceId      int       `json:"provinceId" comment:"客户省份ID"`
	ProvinceName    string    `json:"provinceName" comment:"客户省份名称"`
	CityId          int       `json:"cityId" comment:"收货人城市编号"`
	CityName        string    `json:"cityName" comment:"收货人城市名称"`
	AreaId          int       `json:"areaId" comment:"客户区县ID"`
	AreaName        string    `json:"areaName" comment:"客户区县名称"`
	TownId          int       `json:"townId" comment:"镇/街道 ID"`
	TownName        string    `json:"townName" comment:"镇/街道名称"`
	CompanyName     string    `json:"companyName" comment:"收货人公司名称"`
	Address         string    `json:"address" comment:"收货人详细地址"`
	Mobile          string    `json:"mobile" comment:"收货人手机号"`
	Telephone       string    `json:"telephone" comment:"收货人座机号"`
	ContractNo      string    `json:"contractNo" comment:"客户合同编号"`
	CreateFrom      string    `json:"createFrom" comment:"订单来源：LMS/MALL/XCX"`
	Remark          string    `json:"remark" comment:"客户留言"`
	Ip              string    `json:"ip" comment:"客户IP地址"`
	ApproveStatus   int       `json:"approveStatus" comment:"审批状态 -1 审批不通过  1 审批通过 0 初始提交审批 10 审批中 -2 超时 -3 撤回"`
	Step            int       `json:"step" comment:"审批到第几级"`
	ExpireTime      time.Time `json:"expireTime" comment:"过期时间"`
	common.ControlBy
}

func (s *PresoUpdateReq) Generate(model *models.Preso) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.PresoNo = s.PresoNo
	model.WarehouseCode = s.WarehouseCode
	model.ApproveflowId = s.ApproveflowId
	model.ApproveUsers = s.ApproveUsers
	model.UserId = s.UserId
	model.UserName = s.UserName
	model.UserCompanyId = s.UserCompanyId
	model.UserCompanyName = s.UserCompanyName
	model.DeliverId = s.DeliverId
	model.Consignee = s.Consignee
	model.CountryId = s.CountryId
	model.CountryName = s.CountryName
	model.ProvinceId = s.ProvinceId
	model.ProvinceName = s.ProvinceName
	model.AreaId = s.AreaId
	model.AreaName = s.AreaName
	model.CityId = s.CityId
	model.CityName = s.CityName
	model.TownId = s.TownId
	model.TownName = s.TownName
	model.CompanyName = s.CompanyName
	model.Address = s.Address
	model.Mobile = s.Mobile
	model.Telephone = s.Telephone
	model.ContractNo = s.ContractNo
	model.CreateFrom = s.CreateFrom
	model.Remark = s.Remark
	model.Ip = s.Ip
	model.ApproveStatus = s.ApproveStatus
	model.Step = s.Step
	model.ExpireTime = s.ExpireTime
}

func (s *PresoUpdateReq) GetId() interface{} {
	return s.Id
}

// PresoGetReq 功能获取请求参数
type PresoGetReq struct {
	Id string `uri:"id"`
}

func (s *PresoGetReq) GetId() interface{} {
	return s.Id
}

// PresoDeleteReq 功能删除请求参数
type PresoDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *PresoDeleteReq) GetId() interface{} {
	return s.Ids
}
