package dto

import (
	"go-admin/app/pc/models"
	"go-admin/common/dto"
	"gorm.io/gorm"
	"regexp"
	"strings"
)

type GoodsGetPageReq struct {
	dto.Pagination `search:"-"`
	CategoryId     int    `form:"categoryId"  search:"-"`
	BrandId        []int  `form:"brandId[]"  search:"-"`
	VendorId       []int  `form:"vendorId[]"  search:"-"`
	Keyword        string `form:"keyword"  search:"-"`
	UserId         int    `form:"-"  search:"-"`
	UserName       string `form:"-"  search:"-"`
	WarehouseCode  string `form:"-"  search:"-"`
	HasStock       int    `form:"hasStock"  search:"-"`
	GoodsOrder
}

type GoodsOrder struct {
	MarketPrice string `form:"marketPriceOrder"  search:"type:order;column:market_price;table:goods"`
}

func (m *GoodsGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type GoodsGetPageResp struct {
	GoodsId        int                  `json:"goodsId" comment:"商品ID"`
	NameZh         string               `json:"nameZh" comment:"商品名称"`
	MarketPrice    float64              `json:"marketPrice" comment:"含税价"`
	NakedSalePrice float64              `json:"nakedSalePrice" gorm:"-" comment:"未税价"`
	BrandId        int                  `json:"brandId" comment:"品牌ID"`
	BrandZh        string               `json:"brandZh" comment:"品牌中文"`
	BrandEn        string               `json:"brandEn" comment:"品牌英文"`
	FirstLetter    string               `json:"brandFirstLetter" comment:"-"`
	VendorId       int                  `json:"vendorId" comment:"货主ID"`
	VendorName     string               `json:"vendorName" gorm:"-" comment:"货主名称"`
	SkuCode        string               `json:"skuCode" comment:"SKU"`
	Stock          int                  `json:"stock" comment:"库存" gorm:"-"`
	ProductNo      string               `json:"productNo" comment:"货号"`
	MfgModel       string               `json:"mfgModel" comment:"型号"`
	SalesUom       string               `json:"salesUom" comment:"售卖包装单位"`
	Tax            string               `json:"tax" comment:"税率"`
	SalesMoq       int                  `json:"salesMoq" comment:"销售最小起订量"`
	Image          models.MediaInstance `json:"image" gorm:"-" comment:"图片"`
	IsCollect      bool                 `json:"isCollect" gorm:"-" comment:"是否收藏"`
	Quantity       int                  `json:"quantity" gorm:"-" comment:"加入购物车数量"`
}

type GoodsGetInfo struct {
	models.Goods
	VendorName         string  `json:"vendorName" gorm:"-" comment:"货主名称"`
	SalesMoq           int     `json:"salesMoq" gorm:"-" comment:"最小起订量"`
	Stock              int     `json:"stock" gorm:"-" comment:"库存"`
	NakedSalePrice     float64 `json:"nakedSalePrice" gorm:"-" comment:"未税价"`
	TaxStr             string  `json:"tax" comment:"税率" gorm:"-" comment:"未税价"`
	IsCollect          bool    `json:"isCollect" gorm:"-" comment:"是否收藏"`
	Quantity           int     `json:"quantity" gorm:"-" comment:"加入购物车数量"`
	IsCurrentWarehouse bool    `json:"isCurrentWarehouse" gorm:"-" comment:"是否当前仓库"`
	CheckStockStatus   int     `json:"checkStockStatus" gorm:"-" comment:"库存不足是否允许下单 0否 1是"`
}

// GoodsGetReq 功能获取请求参数
type GoodsGetReq struct {
	Id            int  `uri:"id"`
	View          bool `form:"view"`
	UserId        int
	WarehouseCode string
}

func (s *GoodsGetReq) GetId() interface{} {
	return s.Id
}

func MakeSearchCondition(c *GoodsGetPageReq) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db.Where("goods.warehouse_code = ?", c.WarehouseCode)
		db.Where("goods.online_status = 1")
		// 正则匹配关键词是否是SKU
		c.Keyword = strings.TrimSpace(c.Keyword)
		if c.Keyword != "" {
			reg := regexp.MustCompile(`^[A-Za-z]{3}\d{3}$`)
			res := reg.FindAllString(c.Keyword, -1)
			if len(res) > 0 {
				db.Where("goods.sku_code = ?", c.Keyword)
			} else {
				db.Where("product.name_zh LIKE ? OR product.mfg_model LIKE ? OR brand.brand_zh LIKE ? OR goods.product_no = ? OR goods.supplier_sku_code = ?", "%"+c.Keyword+"%", "%"+c.Keyword+"%", "%"+c.Keyword+"%", c.Keyword, c.Keyword)
				db.Where("product.status = 2")
			}
		}
		// 品牌搜素
		if len(c.BrandId) > 0 {
			db.Where("product.brand_id in ?", c.BrandId)
		}
		// 货主搜索
		if len(c.VendorId) > 0 {
			db.Where("product.vendor_id in ?", c.VendorId)
		}
		return db
	}
}

func MakeWarehouseCondition(c *GoodsGetReq) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if !c.View {
			db.Where("goods.warehouse_code = ?", c.WarehouseCode)
		}
		return db
	}
}

func MakeFilterCategoryIdCondition(c *GoodsGetPageReq, tx *gorm.DB) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// 分类搜索
		if c.CategoryId != 0 {
			categoryPath := models.CategoryPath{}
			skuCode := categoryPath.GetSkuCodeByCateId(tx, c.CategoryId)
			if len(skuCode) > 0 {
				db.Where("goods.sku_code in ?", skuCode)
			} else {
				db.Where("goods.sku_code = ''")
			}
		}
		return db
	}
}

type List struct {
	Product     []GoodsGetPageResp       `json:"product"`
	BrandAll    []map[string]any         `json:"brandAll"`
	BrandFilter []map[string][]BrandInfo `json:"brandFilter"`
	Category    []map[string]any         `json:"category"`
	CategoryNav [][]CategoryNav          `json:"categoryNav"`
	BrandNav    []string                 `json:"brandNav"`
	Vendor      []map[string]any         `json:"vendor"`
}

type GoodsGetResp struct {
	Info                GoodsGetInfo     `json:"info"`
	Navigation          []map[string]any `json:"navigation"`
	ProductExtAttribute any              `json:"productExtAttribute"`
}

type NoWhereResult struct {
	SkuCode string
	BrandId int
	BrandZh string
}

type MiniProgramHomeFilter struct {
	BrandAll []map[string]any `json:"brandAll"`
	Vendor   []map[string]any `json:"vendor"`
}

type AllBrand struct {
	BrandId     int
	BrandZh     string
	FirstLetter string
}

type AllVendor struct {
	VendorId int
}
