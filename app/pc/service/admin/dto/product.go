package dto

import (
	"go-admin/app/pc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
	"go-admin/common/utils"
	"gorm.io/gorm"
	"strconv"
)

type ProductGetPageReq struct {
	dto.Pagination  `search:"-"`
	Ids             []int  `form:"ids[]"  search:"-"`
	SkuCode         string `form:"skuCode"  search:"-"`
	NameZh          string `form:"nameZh"  search:"-"`
	BrandZh         string `form:"brandZh"  search:"-"`
	MfgModel        string `form:"mfgModel"  search:"-"`
	Status          int    `form:"status"  search:"type:exact;column:status;table:product"`
	Level1Catid     int    `form:"level1Catid" search:"-"`
	Level2Catid     int    `form:"level2Catid" search:"-"`
	Level3Catid     int    `form:"level3Catid" search:"-"`
	Level4Catid     int    `form:"level4Catid" search:"-"`
	VendorId        int    `form:"vendorId"  search:"type:exact;column:vendor_id;table:product"`
	SupplierSkuCode string `form:"supplierSkuCode"  search:"type:exact;column:supplier_sku_code;table:product"`
	CreateDateStart string `form:"createDateStart" search:"-"`
	CreateDateEnd   string `form:"createDateEnd" search:"-"`
	ProductOrder
}

type ProductOrder struct {
	Id                  string `form:"idOrder"  search:"type:order;column:id;table:product"`
	VendorId            string `form:"vendorIdOrder"  search:"type:order;column:vendor_id;table:product"`
	SupplierSkuCode     string `form:"supplierSkuCodeOrder"  search:"type:order;column:supplier_sku_code;table:product"`
	SkuCode             string `form:"skuCodeOrder"  search:"type:order;column:sku_code;table:product"`
	NameZh              string `form:"nameZhOrder"  search:"type:order;column:name_zh;table:product"`
	NameEn              string `form:"nameEnOrder"  search:"type:order;column:name_en;table:product"`
	Title               string `form:"titleOrder"  search:"type:order;column:title;table:product"`
	BriefDesc           string `form:"briefDescOrder"  search:"type:order;column:brief_desc;table:product"`
	MfgModel            string `form:"mfgModelOrder"  search:"type:order;column:mfg_model;table:product"`
	BrandId             string `form:"brandIdOrder"  search:"type:order;column:brand_id;table:product"`
	SalesUom            string `form:"salesUomOrder"  search:"type:order;column:sales_uom;table:product"`
	PhysicalUom         string `form:"physicalUomOrder"  search:"type:order;column:physical_uom;table:product"`
	SalesPhysicalFactor string `form:"salesPhysicalFactorOrder"  search:"type:order;column:sales_physical_factor;table:product"`
	SalesMoq            string `form:"salesMoqOrder"  search:"type:order;column:sales_moq;table:product"`
	PackLength          string `form:"packLengthOrder"  search:"type:order;column:pack_length;table:product"`
	PackWidth           string `form:"packWidthOrder"  search:"type:order;column:pack_width;table:product"`
	PackHeight          string `form:"packHeightOrder"  search:"type:order;column:pack_height;table:product"`
	PackWeight          string `form:"packWeightOrder"  search:"type:order;column:pack_weight;table:product"`
	Accessories         string `form:"accessoriesOrder"  search:"type:order;column:accessories;table:product"`
	ProductAlias        string `form:"productAliasOrder"  search:"type:order;column:product_alias;table:product"`
	FragileFlag         string `form:"fragileFlagOrder"  search:"type:order;column:fragile_flag;table:product"`
	HazardFlag          string `form:"hazardFlagOrder"  search:"type:order;column:hazard_flag;table:product"`
	BulkyFlag           string `form:"bulkyFlagOrder"  search:"type:order;column:bulky_flag;table:product"`
	AssembleFlag        string `form:"assembleFlagOrder"  search:"type:order;column:assemble_flag;table:product"`
	IsValuables         string `form:"isValuablesOrder"  search:"type:order;column:is_valuables;table:product"`
	IsFluid             string `form:"isFluidOrder"  search:"type:order;column:is_fluid;table:product"`
	Status              string `form:"statusOrder"  search:"type:order;column:status;table:product"`
	ConsumptiveFlag     string `form:"consumptiveFlagOrder"  search:"type:order;column:consumptive_flag;table:product"`
	StorageFlag         string `form:"storageFlagOrder"  search:"type:order;column:storage_flag;table:product"`
	StorageTime         string `form:"storageTimeOrder"  search:"type:order;column:storage_time;table:product"`
	CustomMadeFlag      string `form:"customMadeFlagOrder"  search:"type:order;column:custom_made_flag;table:product"`
	OtherLabel          string `form:"otherLabelOrder"  search:"type:order;column:other_label;table:product"`
	Tax                 string `form:"taxOrder"  search:"type:order;column:tax;table:product"`
	CreateBy            string `form:"createByOrder"  search:"type:order;column:create_by;table:product"`
	CreateByName        string `form:"createByNameOrder"  search:"type:order;column:create_by_name;table:product"`
	UpdateBy            string `form:"updateByOrder"  search:"type:order;column:update_by;table:product"`
	UpdateByName        string `form:"updateByNameOrder"  search:"type:order;column:update_by_name;table:product"`
	CreatedAt           string `form:"createdAtOrder"  search:"type:order;column:created_at;table:product"`
	UpdatedAt           string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:product"`
	DeletedAt           string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:product"`
}

func (m *ProductGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type ProductCategory struct {
	CategoryId   int    `json:"categoryId" comment:"分类ID" vd:"@:$>0; msg:'分类ID不能为空'"`
	CategoryName string `json:"categoryName" comment:"分类名称" vd:"@:len($)>0; msg:'分类名称不能为空'"`
}

type ProductInsertReq struct {
	Id                     int                      `json:"-" comment:"产品id"` // 产品id
	VendorId               int                      `json:"vendorId" comment:"货主ID" vd:"@:$>0; msg:'货主必填'"`
	SupplierSkuCode        string                   `json:"supplierSkuCode" comment:"货主SKU" vd:"@:len($)>0; msg:'货主SKU必填'"`
	ProductCategory        []ProductCategory        `json:"productCategory" comment:"产品目录"`
	AttrList               []ProductExtAttribute    `json:"attrList" comment:"产品扩展属性"`
	ProductCategoryPrimary int                      `json:"productCategoryPrimary" comment:"商品主分类" vd:"@:$>0; msg:'请选择商品主分类'"`
	SkuCode                string                   `json:"skuCode" comment:"产品sku"`
	NameZh                 string                   `json:"nameZh" comment:"中文名" vd:"@:len($)>0; msg:'商品名称必填'"`
	NameEn                 string                   `json:"nameEn" comment:"英文名"`
	Title                  string                   `json:"title" comment:"标题"`
	BriefDesc              string                   `json:"briefDesc" comment:"产品描述" vd:"@:len($)<700000; msg:'产品描述字段超出长度限制'"`
	MfgModel               string                   `json:"mfgModel" comment:"制造厂型号" vd:"@:len($)>0; msg:'制造厂型号必填'"`
	BrandId                int                      `json:"brandId" comment:"品牌id" vd:"@:$>0; msg:'品牌必填'"`
	SalesUom               string                   `json:"salesUom" comment:"售卖包装单位" vd:"@:len($)>0; msg:'售卖包装单位必填'"`
	PhysicalUom            string                   `json:"physicalUom" comment:"物理单位"`
	SalesPhysicalFactor    string                   `json:"salesPhysicalFactor" comment:"销售包装单位中含物理单位数量"`
	SalesMoq               int                      `json:"salesMoq" comment:"销售最小起订量" vd:"@:$>0; msg:'销售最小起订量至少为1'"`
	PackLength             string                   `json:"packLength" comment:"售卖包装长(mm)"`
	PackWidth              string                   `json:"packWidth" comment:"售卖包装宽(mm)"`
	PackHeight             string                   `json:"packHeight" comment:"售卖包装高(mm)"`
	PackWeight             string                   `json:"packWeight" comment:"售卖包装质量(kg)"`
	Accessories            string                   `json:"accessories" comment:"包装清单"`
	ProductAlias           string                   `json:"productAlias" comment:"产品别名"`
	FragileFlag            int                      `json:"fragileFlag" comment:"易碎标志(0:否;1:是)" vd:"@:in($,0,1); msg:'易碎标志只能为0或1'"`
	HazardFlag             int                      `json:"hazardFlag" comment:"危险品标志(0:否;1:是)" vd:"@:in($,0,1); msg:'危险品标志只能为0或1'"`
	HazardClass            int                      `json:"hazardClass" comment:"危险品等级"`
	BulkyFlag              int                      `json:"bulkyFlag" comment:"抛货标志(0:否; 1:是)" vd:"@:in($,0,1); msg:'抛货标志只能为0或1'"`
	AssembleFlag           int                      `json:"assembleFlag" comment:"拼装件标志(0:否; 1:是)" vd:"@:in($,0,1); msg:'拼装件标志只能为0或1'"`
	IsValuables            int                      `json:"isValuables" comment:"是否贵重品(0:否; 1:是)" vd:"@:in($,0,1); msg:'是否贵重品只能为0或1'"`
	IsFluid                int                      `json:"isFluid" comment:"是否液体(0:否; 1:是)" vd:"@:in($,0,1); msg:'是否液体只能为0或1'"`
	Status                 int                      `json:"status" comment:"产品状态" vd:"@:in($,1,2,3); msg:'产品状态只能为1或2或3'"`
	ConsumptiveFlag        int                      `json:"consumptiveFlag" comment:"耗材标志" vd:"@:in($,0,1); msg:'耗材标志只能为0或1'"`
	StorageFlag            int                      `json:"storageFlag" comment:"保存期标志" vd:"@:in($,0,1); msg:'保存期标志只能为0或1'"`
	StorageTime            int                      `json:"storageTime" comment:"保存期限(月)"`
	CustomMadeFlag         int                      `json:"customMadeFlag" comment:"定制品标志(0:否,1:是)" vd:"@:in($,0,1); msg:'定制品标志只能为0或1'"`
	RefundFlag             int                      `json:"refundFlag" comment:"可退换货标志(0:否; 1:是)" vd:"@:in($,0,1); msg:'可退换货标志只能为0或1'"`
	OtherLabel             string                   `json:"otherLabel" comment:"其他备注"`
	Barcode                string                   `json:"barcode" comment:"商品条形码"`
	Tax                    string                   `json:"tax" comment:"产品税率：默认空 值有（13%，9%，6%）" vd:"in($,'0.13','0.06','0.09'); msg:'产线税率的值有（0.13,0.06,0.09）'"`
	CertificateImage       string                   `json:"certificateImage" comment:"合格证"`
	ProductImage           []MediaInstanceInsertReq `json:"productImage" comment:"商品图片"`
	common.ControlBy
}

func (s *ProductInsertReq) Generate(model *models.Product) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.VendorId = s.VendorId
	model.SupplierSkuCode = s.SupplierSkuCode
	model.SkuCode = s.SkuCode
	model.NameZh = s.NameZh
	model.NameEn = s.NameEn
	model.Title = s.Title
	model.BriefDesc = s.BriefDesc
	model.MfgModel = s.MfgModel
	model.BrandId = s.BrandId
	model.SalesUom = s.SalesUom
	model.PhysicalUom = s.PhysicalUom
	model.SalesPhysicalFactor = s.SalesPhysicalFactor
	model.SalesMoq = s.SalesMoq
	model.PackLength = s.PackLength
	model.PackWidth = s.PackWidth
	model.PackHeight = s.PackHeight
	model.PackWeight = s.PackWeight
	model.Accessories = s.Accessories
	model.ProductAlias = s.ProductAlias
	model.FragileFlag = s.FragileFlag
	model.HazardFlag = s.HazardFlag
	model.HazardClass = s.HazardClass
	model.BulkyFlag = s.BulkyFlag
	model.AssembleFlag = s.AssembleFlag
	model.IsValuables = s.IsValuables
	model.IsFluid = s.IsFluid
	model.Status = s.Status
	model.ConsumptiveFlag = s.ConsumptiveFlag
	model.StorageFlag = s.StorageFlag
	model.StorageTime = s.StorageTime
	model.CustomMadeFlag = s.CustomMadeFlag
	model.RefundFlag = s.RefundFlag
	model.OtherLabel = s.OtherLabel
	model.Barcode = s.Barcode
	model.Tax = s.Tax
	model.CertificateImage = s.CertificateImage
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.CreateByName = s.CreateByName
}

func (s *ProductInsertReq) GetId() interface{} {
	return s.Id
}

type ProductImportReq struct {
	Data []ProductImportData
}

type ProductImportData struct {
	Id                     int               `json:"-" comment:"产品id"` // 产品id
	SkuCode                string            `json:"-" comment:"产品sku"`
	ProductCategory        []ProductCategory `json:"-" comment:"产品目录"`
	ProductCategoryPrimary int               `json:"-" comment:"商品主分类"`
	Level1CatName          string            `json:"level1CatName" comment:"一级目录" vd:"@:len($)>0; msg:'一级目录必填'"`
	Level2CatName          string            `json:"level2CatName" comment:"二级目录" vd:"@:len($)>0; msg:'二级目录必填'"`
	Level3CatName          string            `json:"level3CatName" comment:"三级目录" vd:"@:len($)>0; msg:'三级目录必填'"`
	Level4CatName          string            `json:"level4CatName" comment:"四级目录"`
	NameZh                 string            `json:"nameZh" comment:"中文名" vd:"@:len($)>0; msg:'商品名称必填'"`
	NameEn                 string            `json:"nameEn" comment:"英文名"`
	VendorId               int               `json:"-" comment:"货主ID"`
	VendorName             string            `json:"vendorName" comment:"货主名称" vd:"@:len($)>0; msg:'货主必填'"`
	SupplierSkuCode        string            `json:"supplierSkuCode" comment:"货主SKU" vd:"regexp('^[a-zA-Z0-9]{1,20}$'); msg:'货主SKU只能英文、数字，长度在20以内'"`
	BrandId                int               `json:"-" comment:"品牌ID"`
	BrandZh                string            `json:"brandZh" comment:"品牌(中文)" vd:"@:len($)>0; msg:'品牌中文必填'"`
	BrandEn                string            `json:"brandEn" comment:"品牌(英文)"`
	MfgModel               string            `json:"mfgModel" comment:"制造厂型号" vd:"@:len($)>0; msg:'制造厂型号必填'"`
	SalesUom               string            `json:"salesUom" comment:"售卖包装单位" vd:"@:len($)>0; msg:'售卖包装单位必填'"`
	PhysicalUom            string            `json:"physicalUom" comment:"物理单位"`
	SalesPhysicalFactor    string            `json:"salesPhysicalFactor" comment:"销售包装单位中含物理单位数量"`
	SalesMoq               int               `json:"salesMoq" comment:"销售最小起订量" vd:"@:$>0; msg:'销售最小起订量至少为1'"`
	PackWeight             string            `json:"packWeight" comment:"重量(kg)"`
	PackLength             string            `json:"packLength" comment:"售卖包装长(mm)"`
	PackWidth              string            `json:"packWidth" comment:"售卖包装宽(mm)"`
	PackHeight             string            `json:"packHeight" comment:"售卖包装高(mm)"`
	FragileFlag            int               `json:"fragileFlag" comment:"易碎标志(0:否;1:是)" vd:"@:in($,0,1); msg:'易碎标志只能为0或1'"`
	HazardFlag             int               `json:"hazardFlag" comment:"危险品标志(0:否;1:是)" vd:"@:in($,0,1); msg:'危险品标志只能为0或1'"`
	HazardClass            int               `json:"hazardClass" comment:"危险品等级"`
	BulkyFlag              int               `json:"bulkyFlag" comment:"抛货标志(0:否; 1:是)" vd:"@:in($,0,1); msg:'抛货标志只能为0或1'"`
	AssembleFlag           int               `json:"assembleFlag" comment:"拼装件标志(0:否; 1:是)" vd:"@:in($,0,1); msg:'拼装件标志只能为0或1'"`
	IsValuables            int               `json:"isValuables" comment:"是否贵重品(0:否; 1:是)" vd:"@:in($,0,1); msg:'是否贵重品只能为0或1'"`
	IsFluid                int               `json:"isFluid" comment:"是否液体(0:否; 1:是)" vd:"@:in($,0,1); msg:'是否液体只能为0或1'"`
	ConsumptiveFlag        int               `json:"consumptiveFlag" comment:"耗材标志" vd:"@:in($,0,1); msg:'耗材标志只能为0或1'"`
	StorageFlag            int               `json:"storageFlag" comment:"保存期标志" vd:"@:in($,0,1); msg:'保存期标志只能为0或1'"`
	StorageTime            int               `json:"storageTime" comment:"保存期限(月)"`
	CustomMadeFlag         int               `json:"customMadeFlag" comment:"定制品标志(0:否,1:是)" vd:"@:in($,0,1); msg:'定制品标志只能为0或1'"`
	RefundFlag             int               `json:"refundFlag" comment:"可退换货标志(0:否; 1:是)" vd:"@:in($,0,1); msg:'可退换货标志只能为0或1'"`
	Tax                    string            `json:"tax" comment:"产品税率：默认空 值有（0.13，0.09，0.06）" vd:"in($,'0.13','0.06','0.09'); msg:'产线税率的值有（0.13,0.06,0.09）'"`
	common.ControlBy
}

func (product *ProductImportData) Generate(k string, v any) {
	switch k {
	case "level1CatName":
		product.Level1CatName = v.(string)
	case "level2CatName":
		product.Level2CatName = v.(string)
	case "level3CatName":
		product.Level3CatName = v.(string)
	case "level4CatName":
		product.Level4CatName = v.(string)
	case "nameZh":
		product.NameZh = v.(string)
	case "nameEn":
		product.NameEn = v.(string)
	case "vendorName":
		product.VendorName = v.(string)
	case "supplierSkuCode":
		product.SupplierSkuCode = v.(string)
	case "brandZh":
		product.BrandZh = v.(string)
	case "brandEn":
		product.BrandEn = v.(string)
	case "mfgModel":
		product.MfgModel = v.(string)
	case "salesUom":
		product.SalesUom = v.(string)
	case "physicalUom":
		product.PhysicalUom = v.(string)
	case "salesPhysicalFactor":
		product.SalesPhysicalFactor = v.(string)
	case "salesMoq":
		if v == "" {
			product.SalesMoq = 1
		} else {
			product.SalesMoq, _ = strconv.Atoi(v.(string))
		}
	case "packWeight":
		product.PackWeight = v.(string)
	case "packLength":
		product.PackLength = v.(string)
	case "packWidth":
		product.PackWidth = v.(string)
	case "packHeight":
		product.PackHeight = v.(string)
	case "fragileFlag":
		if v == "" {
			product.FragileFlag = 0
		} else {
			product.FragileFlag, _ = strconv.Atoi(v.(string))
		}
	case "hazardFlag":
		if v == "" {
			product.HazardFlag = 0
		} else {
			product.HazardFlag, _ = strconv.Atoi(v.(string))
		}
	case "hazardClass":
		if v == "" {
			product.HazardClass = 0
		} else {
			product.HazardClass, _ = strconv.Atoi(v.(string))
		}
	case "bulkyFlag":
		if v == "" {
			product.BulkyFlag = 0
		} else {
			product.BulkyFlag, _ = strconv.Atoi(v.(string))
		}
	case "assembleFlag":
		if v == "" {
			product.AssembleFlag = 0
		} else {
			product.AssembleFlag, _ = strconv.Atoi(v.(string))
		}
	case "isValuables":
		if v == "" {
			product.IsValuables = 0
		} else {
			product.IsValuables, _ = strconv.Atoi(v.(string))
		}
	case "isFluid":
		if v == "" {
			product.IsFluid = 0
		} else {
			product.IsFluid, _ = strconv.Atoi(v.(string))
		}
	case "consumptiveFlag":
		if v == "" {
			product.ConsumptiveFlag = 0
		} else {
			product.ConsumptiveFlag, _ = strconv.Atoi(v.(string))
		}
	case "storageFlag":
		if v == "" {
			product.StorageFlag = 0
		} else {
			product.StorageFlag, _ = strconv.Atoi(v.(string))
		}
	case "storageTime":
		if v == "" {
			product.StorageTime = 0
		} else {
			product.StorageTime, _ = strconv.Atoi(v.(string))
		}
	case "customMadeFlag":
		if v == "" {
			product.CustomMadeFlag = 0
		} else {
			product.CustomMadeFlag, _ = strconv.Atoi(v.(string))
		}
	case "refundFlag":
		if v == "" {
			product.RefundFlag = 0
		} else {
			product.RefundFlag, _ = strconv.Atoi(v.(string))
		}
	case "tax":
		product.Tax = v.(string)
	}
}

func (s *ProductImportData) GetId() interface{} {
	return s.Id
}

type ProductImportUpdateReq struct {
	Id                     int               `json:"-" comment:"产品id"` // 产品id
	SkuCode                string            `json:"skuCode" comment:"产品sku" vd:"@:len($)>0; msg:'产品sku必填'"`
	ProductCategory        []ProductCategory `json:"-" comment:"产品目录"`
	ProductCategoryPrimary int               `json:"-" comment:"商品主分类"`
	Level1CatName          string            `json:"level1CatName" comment:"一级目录" vd:"@:len($)>0; msg:'一级目录必填'"`
	Level2CatName          string            `json:"level2CatName" comment:"二级目录" vd:"@:len($)>0; msg:'二级目录必填'"`
	Level3CatName          string            `json:"level3CatName" comment:"三级目录" vd:"@:len($)>0; msg:'三级目录必填'"`
	Level4CatName          string            `json:"level4CatName" comment:"四级目录"`
	NameZh                 string            `json:"nameZh" comment:"产品名称(中文)" vd:"@:len($)>0; msg:'产品名称(中文)必填'"`
	NameEn                 string            `json:"nameEn" comment:"产品名称(英文)"`
	VendorName             string            `json:"vendorName" comment:"货主ID" vd:"@:len($)>0; msg:'货主必填'"`
	SupplierSkuCode        string            `json:"supplierSkuCode" comment:"货主SKU" vd:"regexp('^[a-zA-Z0-9]{1,20}$'); msg:'货主SKU只能英文、数字，长度在20以内'"`
	BrandId                int               `json:"-" comment:"品牌ID"`
	BrandZh                string            `json:"brandZh" comment:"品牌(中文)" vd:"@:len($)>0; msg:'品牌中文必填'"`
	BrandEn                string            `json:"brandEn" comment:"品牌(英文)"`
	MfgModel               string            `json:"mfgModel" comment:"制造厂型号" vd:"@:len($)>0; msg:'制造厂型号必填'"`
	SalesUom               string            `json:"salesUom" comment:"售卖包装单位" vd:"@:len($)>0; msg:'售卖包装单位必填'"`
	PhysicalUom            string            `json:"physicalUom" comment:"物理单位"`
	SalesPhysicalFactor    string            `json:"salesPhysicalFactor" comment:"销售包装单位中含物理单位数量"`
	SalesMoq               int               `json:"salesMoq" comment:"销售最小起订量" vd:"@:$>0; msg:'销售最小起订量至少为1'"`
	PackWeight             string            `json:"packWeight" comment:"重量(kg)"`
	PackLength             string            `json:"packLength" comment:"售卖包装长(mm)"`
	PackWidth              string            `json:"packWidth" comment:"售卖包装宽(mm)"`
	PackHeight             string            `json:"packHeight" comment:"售卖包装高(mm)"`
	FragileFlag            int               `json:"fragileFlag" comment:"易碎标志(0:否;1:是)" vd:"@:in($,0,1); msg:'易碎标志只能为0或1'"`
	HazardFlag             int               `json:"hazardFlag" comment:"危险品标志(0:否;1:是)" vd:"@:in($,0,1); msg:'危险品标志只能为0或1'"`
	HazardClass            int               `json:"hazardClass" comment:"危险品等级"`
	BulkyFlag              int               `json:"bulkyFlag" comment:"抛货标志(0:否; 1:是)" vd:"@:in($,0,1); msg:'抛货标志只能为0或1'"`
	AssembleFlag           int               `json:"assembleFlag" comment:"拼装件标志(0:否; 1:是)" vd:"@:in($,0,1); msg:'拼装件标志只能为0或1'"`
	IsValuables            int               `json:"isValuables" comment:"是否贵重品(0:否; 1:是)" vd:"@:in($,0,1); msg:'是否贵重品只能为0或1'"`
	IsFluid                int               `json:"isFluid" comment:"是否液体(0:否; 1:是)" vd:"@:in($,0,1); msg:'是否液体只能为0或1'"`
	ConsumptiveFlag        int               `json:"consumptiveFlag" comment:"耗材标志" vd:"@:in($,0,1); msg:'耗材标志只能为0或1'"`
	StorageFlag            int               `json:"storageFlag" comment:"保存期标志" vd:"@:in($,0,1); msg:'保存期标志只能为0或1'"`
	StorageTime            int               `json:"storageTime" comment:"保存期限(月)"`
	CustomMadeFlag         int               `json:"customMadeFlag" comment:"定制品标志(0:否,1:是)" vd:"@:in($,0,1); msg:'定制品标志只能为0或1'"`
	RefundFlag             int               `json:"refundFlag" comment:"可退换货标志(0:否; 1:是)" vd:"@:in($,0,1); msg:'可退换货标志只能为0或1'"`
	common.ControlBy
}

func (model *ProductImportUpdateReq) MapToStruct(k string, v any) {
	switch k {
	case "skuCode":
		model.SkuCode = v.(string)
	case "level1CatName":
		model.Level1CatName = v.(string)
	case "level2CatName":
		model.Level2CatName = v.(string)
	case "level3CatName":
		model.Level3CatName = v.(string)
	case "level4CatName":
		model.Level4CatName = v.(string)
	case "nameZh":
		model.NameZh = v.(string)
	case "nameEn":
		model.NameEn = v.(string)
	case "vendorName":
		model.VendorName = v.(string)
	case "supplierSkuCode":
		model.SupplierSkuCode = v.(string)
	case "brandZh":
		model.BrandZh = v.(string)
	case "brandEn":
		model.BrandEn = v.(string)
	case "mfgModel":
		model.MfgModel = v.(string)
	case "salesUom":
		model.SalesUom = v.(string)
	case "physicalUom":
		model.PhysicalUom = v.(string)
	case "salesPhysicalFactor":
		model.SalesPhysicalFactor = v.(string)
	case "salesMoq":
		if v == "" {
			model.SalesMoq = 1
		} else {
			model.SalesMoq, _ = strconv.Atoi(v.(string))
		}
	case "packWeight":
		model.PackWeight = v.(string)
	case "packLength":
		model.PackLength = v.(string)
	case "packWidth":
		model.PackWidth = v.(string)
	case "packHeight":
		model.PackHeight = v.(string)
	case "fragileFlag":
		if v == "" {
			model.FragileFlag = 0
		} else {
			model.FragileFlag, _ = strconv.Atoi(v.(string))
		}
	case "hazardFlag":
		if v == "" {
			model.HazardFlag = 0
		} else {
			model.HazardFlag, _ = strconv.Atoi(v.(string))
		}
	case "hazardClass":
		if v == "" {
			model.HazardClass = 0
		} else {
			model.HazardClass, _ = strconv.Atoi(v.(string))
		}
	case "bulkyFlag":
		if v == "" {
			model.BulkyFlag = 0
		} else {
			model.BulkyFlag, _ = strconv.Atoi(v.(string))
		}
	case "assembleFlag":
		if v == "" {
			model.AssembleFlag = 0
		} else {
			model.AssembleFlag, _ = strconv.Atoi(v.(string))
		}
	case "isValuables":
		if v == "" {
			model.IsValuables = 0
		} else {
			model.IsValuables, _ = strconv.Atoi(v.(string))
		}
	case "isFluid":
		if v == "" {
			model.IsFluid = 0
		} else {
			model.IsFluid, _ = strconv.Atoi(v.(string))
		}
	case "consumptiveFlag":
		if v == "" {
			model.ConsumptiveFlag = 0
		} else {
			model.ConsumptiveFlag, _ = strconv.Atoi(v.(string))
		}
	case "storageFlag":
		if v == "" {
			model.StorageFlag = 0
		} else {
			model.StorageFlag, _ = strconv.Atoi(v.(string))
		}
	case "storageTime":
		if v == "" {
			model.StorageTime = 0
		} else {
			model.StorageTime, _ = strconv.Atoi(v.(string))
		}
	case "customMadeFlag":
		if v == "" {
			model.CustomMadeFlag = 0
		} else {
			model.CustomMadeFlag, _ = strconv.Atoi(v.(string))
		}
	case "refundFlag":
		if v == "" {
			model.RefundFlag = 0
		} else {
			model.RefundFlag, _ = strconv.Atoi(v.(string))
		}
	}
}

func (s *ProductImportUpdateReq) Generate(model *models.Product) {
	model.SupplierSkuCode = s.SupplierSkuCode
	model.SkuCode = s.SkuCode
	model.NameZh = s.NameZh
	model.NameEn = s.NameEn
	model.MfgModel = s.MfgModel
	model.BrandId = s.BrandId
	model.SalesUom = s.SalesUom
	model.PhysicalUom = s.PhysicalUom
	model.SalesPhysicalFactor = s.SalesPhysicalFactor
	model.SalesMoq = s.SalesMoq
	model.PackLength = s.PackLength
	model.PackWidth = s.PackWidth
	model.PackHeight = s.PackHeight
	model.PackWeight = s.PackWeight
	model.FragileFlag = s.FragileFlag
	model.HazardFlag = s.HazardFlag
	model.HazardClass = s.HazardClass
	model.BulkyFlag = s.BulkyFlag
	model.AssembleFlag = s.AssembleFlag
	model.IsValuables = s.IsValuables
	model.IsFluid = s.IsFluid
	model.ConsumptiveFlag = s.ConsumptiveFlag
	model.StorageFlag = s.StorageFlag
	model.StorageTime = s.StorageTime
	model.CustomMadeFlag = s.CustomMadeFlag
	model.RefundFlag = s.RefundFlag
	model.UpdateBy = s.UpdateBy
	model.UpdateByName = s.UpdateByName // 添加这而，需要记录是被谁更新的
}

func (s *ProductImportUpdateReq) GetId() interface{} {
	return s.Id
}

type ProductUpdateReq struct {
	Id                     int                      `uri:"id" comment:"产品id"` // 产品id
	VendorId               int                      `json:"vendorId" comment:"货主ID" vd:"@:$>0; msg:'货主必填'"`
	SupplierSkuCode        string                   `json:"supplierSkuCode" comment:"货主SKU" vd:"@:len($)>0; msg:'货主SKU必填'"`
	ProductCategory        []ProductCategory        `json:"productCategory" comment:"产品目录"`
	AttrList               []ProductExtAttribute    `json:"attrList" comment:"产品扩展属性"`
	ProductCategoryPrimary int                      `json:"productCategoryPrimary" comment:"商品主分类" vd:"@:$>0; msg:'请选择商品主分类'"`
	SkuCode                string                   `json:"skuCode" comment:"产品sku"`
	NameZh                 string                   `json:"nameZh" comment:"中文名" vd:"@:len($)>0; msg:'商品名称必填'"`
	NameEn                 string                   `json:"nameEn" comment:"英文名"`
	Title                  string                   `json:"title" comment:"标题"`
	BriefDesc              string                   `json:"briefDesc" comment:"产品描述" vd:"@:len($)<700000; msg:'产品描述字段超出长度限制'"`
	MfgModel               string                   `json:"mfgModel" comment:"制造厂型号" vd:"@:len($)>0; msg:'制造厂型号必填'"`
	BrandId                int                      `json:"brandId" comment:"品牌id" vd:"@:$>0; msg:'品牌必填'"`
	SalesUom               string                   `json:"salesUom" comment:"售卖包装单位" vd:"@:len($)>0; msg:'售卖包装单位必填'"`
	PhysicalUom            string                   `json:"physicalUom" comment:"物理单位"`
	SalesPhysicalFactor    string                   `json:"salesPhysicalFactor" comment:"销售包装单位中含物理单位数量"`
	SalesMoq               int                      `json:"salesMoq" comment:"销售最小起订量" vd:"@:$>0; msg:'销售最小起订量至少为1'"`
	PackLength             string                   `json:"packLength" comment:"售卖包装长(mm)"`
	PackWidth              string                   `json:"packWidth" comment:"售卖包装宽(mm)"`
	PackHeight             string                   `json:"packHeight" comment:"售卖包装高(mm)"`
	PackWeight             string                   `json:"packWeight" comment:"售卖包装质量(kg)"`
	Accessories            string                   `json:"accessories" comment:"包装清单"`
	ProductAlias           string                   `json:"productAlias" comment:"产品别名"`
	FragileFlag            int                      `json:"fragileFlag" comment:"易碎标志(0:否;1:是)" vd:"@:in($,0,1); msg:'易碎标志只能为0或1'"`
	HazardFlag             int                      `json:"hazardFlag" comment:"危险品标志(0:否;1:是)" vd:"@:in($,0,1); msg:'危险品标志只能为0或1'"`
	HazardClass            int                      `json:"hazardClass" comment:"危险品等级"`
	BulkyFlag              int                      `json:"bulkyFlag" comment:"抛货标志(0:否; 1:是)" vd:"@:in($,0,1); msg:'抛货标志只能为0或1'"`
	AssembleFlag           int                      `json:"assembleFlag" comment:"拼装件标志(0:否; 1:是)" vd:"@:in($,0,1); msg:'拼装件标志只能为0或1'"`
	IsValuables            int                      `json:"isValuables" comment:"是否贵重品(0:否; 1:是)" vd:"@:in($,0,1); msg:'是否贵重品只能为0或1'"`
	IsFluid                int                      `json:"isFluid" comment:"是否液体(0:否; 1:是)" vd:"@:in($,0,1); msg:'是否液体只能为0或1'"`
	Status                 int                      `json:"status" comment:"产品状态" vd:"@:in($,1,2,3); msg:'产品状态只能为1或2或3'"`
	ConsumptiveFlag        int                      `json:"consumptiveFlag" comment:"耗材标志" vd:"@:in($,0,1); msg:'耗材标志只能为0或1'"`
	StorageFlag            int                      `json:"storageFlag" comment:"保存期标志" vd:"@:in($,0,1); msg:'保存期标志只能为0或1'"`
	StorageTime            int                      `json:"storageTime" comment:"保存期限(月)"`
	CustomMadeFlag         int                      `json:"customMadeFlag" comment:"定制品标志(0:否,1:是)" vd:"@:in($,0,1); msg:'定制品标志只能为0或1'"`
	RefundFlag             int                      `json:"refundFlag" comment:"可退换货标志(0:否; 1:是)" vd:"@:in($,0,1); msg:'可退换货标志只能为0或1'"`
	Barcode                string                   `json:"barcode" comment:"商品条形码"`
	OtherLabel             string                   `json:"otherLabel" comment:"其他备注"`
	Tax                    string                   `json:"tax" comment:"产品税率：默认空 值有（13%，9%，6%）" vd:"in($,'0.13','0.06','0.09'); msg:'产线税率的值有（0.13,0.06,0.09）'"`
	CertificateImage       string                   `json:"certificateImage" comment:"合格证"`
	ProductImage           []MediaInstanceInsertReq `json:"productImage" comment:"商品图片"`
	common.ControlBy
}

func (s *ProductUpdateReq) Generate(model *models.Product) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.VendorId = s.VendorId
	model.SupplierSkuCode = s.SupplierSkuCode
	model.SkuCode = s.SkuCode
	model.NameZh = s.NameZh
	model.NameEn = s.NameEn
	model.Title = s.Title
	model.BriefDesc = s.BriefDesc
	model.MfgModel = s.MfgModel
	model.BrandId = s.BrandId
	model.SalesUom = s.SalesUom
	model.PhysicalUom = s.PhysicalUom
	model.SalesPhysicalFactor = s.SalesPhysicalFactor
	model.SalesMoq = s.SalesMoq
	model.PackLength = s.PackLength
	model.PackWidth = s.PackWidth
	model.PackHeight = s.PackHeight
	model.PackWeight = s.PackWeight
	model.Accessories = s.Accessories
	model.ProductAlias = s.ProductAlias
	model.FragileFlag = s.FragileFlag
	model.HazardFlag = s.HazardFlag
	model.HazardClass = s.HazardClass
	model.BulkyFlag = s.BulkyFlag
	model.AssembleFlag = s.AssembleFlag
	model.IsValuables = s.IsValuables
	model.IsFluid = s.IsFluid
	model.Status = s.Status
	model.ConsumptiveFlag = s.ConsumptiveFlag
	model.StorageFlag = s.StorageFlag
	model.StorageTime = s.StorageTime
	model.CustomMadeFlag = s.CustomMadeFlag
	model.RefundFlag = s.RefundFlag
	model.OtherLabel = s.OtherLabel
	model.Barcode = s.Barcode
	model.Tax = s.Tax
	model.CertificateImage = s.CertificateImage
	model.UpdateBy = s.UpdateBy
	model.UpdateByName = s.UpdateByName // 添加这而，需要记录是被谁更新的
}

func (s *ProductUpdateReq) GetId() interface{} {
	return s.Id
}

type ProductExportResp struct {
	SkuCode             string `json:"skuCode" comment:"产品SKU"`
	NameZh              string `json:"nameZh" comment:"产品名称(中文)"`
	NameEn              string `json:"nameEn" comment:"产品名称(英文)"`
	BrandZh             string `json:"brandZh" comment:"品牌(中文)"`
	BrandEn             string `json:"brandEn" comment:"品牌(英文)"`
	VendorName          string `json:"vendorName" comment:"货主名称"`
	SupplierSkuCode     string `json:"supplierSkuCode" comment:"货主SKU"`
	Level1CatName       string `json:"level1CatName" comment:"一级目录"`
	Level2CatName       string `json:"level2CatName" comment:"二级目录"`
	Level3CatName       string `json:"level3CatName" comment:"三级目录"`
	Level4CatName       string `json:"level4CatName" comment:"四级目录"`
	MfgModel            string `json:"mfgModel" comment:"制造厂型号"`
	PhysicalUom         string `json:"physicalUom" comment:"物理单位"`
	SalesUom            string `json:"salesUom" comment:"售卖包装单位"`
	SalesPhysicalFactor string `json:"salesPhysicalFactor" comment:"销售包装单位中含物理单位数量"`
	PackWeight          string `json:"packWeight" comment:"重量(kg)"`
	PackLength          string `json:"packLength" comment:"售卖包装长(mm)"`
	PackWidth           string `json:"packWidth" comment:"售卖包装宽(mm)"`
	PackHeight          string `json:"packHeight" comment:"售卖包装高(mm)"`
	FragileFlag         int    `json:"fragileFlag" comment:"易碎标志(0:否;1:是)"`
	HazardFlag          int    `json:"hazardFlag" comment:"危险品标志(0:否;1:是)"`
	HazardClass         string `json:"hazardClass" comment:"危险品等级"`
	CustomMadeFlag      int    `json:"customMadeFlag" comment:"定制品标志(0:否,1:是)"`
	BulkyFlag           int    `json:"bulkyFlag" comment:"抛货标志(0:否; 1:是)"`
	AssembleFlag        int    `json:"assembleFlag" comment:"拼装件标志(0:否; 1:是)"`
	IsValuables         int    `json:"isValuables" comment:"是否贵重品(0:否; 1:是)"`
	IsFluid             int    `json:"isFluid" comment:"是否液体(0:否; 1:是)"`
	ConsumptiveFlag     int    `json:"consumptiveFlag" comment:"耗材标志(0:否; 1:是)"`
	StorageFlag         int    `json:"storageFlag" comment:"保存期标志(0:否; 1:是)"`
	StorageTime         int    `json:"storageTime" comment:"保存期限(月)"`
	RefundFlag          int    `json:"refundFlag" comment:"可退换货标志(0:否; 1:是)"`
}

type BatchProApprovalReq struct {
	Selected []int `json:"selected" comment:"选中行"`
	common.ControlBy
}

type ProApprovalReq struct {
	Status    int `json:"status" comment:"审核状态" vd:"@:in($,2,3); msg:'产品状态只能为2或3！'"`
	ProductId int `json:"productId" comment:"产品状态" vd:"@:$ > 0; msg:'产品ID不能为空！'"`
	common.ControlBy
}

// ProductGetReq 功能获取请求参数
type ProductGetReq struct {
	Id int `uri:"id"`
}

func (s *ProductGetReq) GetId() interface{} {
	return s.Id
}

type GetProductBySkuCodeReq struct {
	SkuCode string `uri:"skuCode"`
}

type GetProductBySkuCodeResp struct {
	NameZh          string `json:"nameZh" gorm:"name_zh"`
	VendorId        int    `json:"vendorId" gorm:"vendor_id"`
	VendorName      string `json:"vendorName" gorm:"-"`
	SupplierSkuCode string `json:"supplierSkuCode" gorm:"supplier_sku_code"`
}

// ProductDeleteReq 功能删除请求参数
type ProductDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *ProductDeleteReq) GetId() interface{} {
	return s.Ids
}

type ProductAttr struct {
	NameZh     string
	CategoryId int
}

func MakeCondition(configure map[string]interface{}, id int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if id != 0 {
			db.Where("id <> ?", id)
		}
		db.Where(configure)
		return db
	}
}

type ExportReq struct {
	Ids []int `form:"ids"`
}

func MakeSearchCondition(c *ProductGetPageReq, tx *gorm.DB) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if c.SkuCode != "" {
			skuCodes := utils.Split(c.SkuCode)
			db.Where("product.sku_code in ?", skuCodes)
		}
		if len(c.Ids) > 0 {
			db.Where("product.id in ?", c.Ids)
		}
		categoryPath := models.CategoryPath{}
		if c.Level1Catid != 0 && c.Level2Catid == 0 {
			skuCode := categoryPath.GetSkuCodeByCateId(tx, c.Level1Catid)
			if len(skuCode) > 0 {
				db.Where("product.sku_code in ?", skuCode)
			} else {
				db.Where("product.sku_code = ''")
			}
		} else if c.Level2Catid != 0 && c.Level3Catid == 0 {
			skuCode := categoryPath.GetSkuCodeByCateId(tx, c.Level2Catid)
			if len(skuCode) > 0 {
				db.Where("product.sku_code in ?", skuCode)
			} else {
				db.Where("product.sku_code = ''")
			}
		} else if c.Level3Catid != 0 && c.Level4Catid == 0 {
			skuCode := categoryPath.GetSkuCodeByCateId(tx, c.Level3Catid)
			if len(skuCode) > 0 {
				db.Where("product.sku_code in ?", skuCode)
			} else {
				db.Where("product.sku_code = ''")
			}
		} else if c.Level4Catid != 0 {
			skuCode := categoryPath.GetSkuCodeByCateId(tx, c.Level4Catid)
			if len(skuCode) > 0 {
				db.Where("product.sku_code in ?", skuCode)
			} else {
				db.Where("product.sku_code = ''")
			}
		}
		if c.NameZh != "" {
			db.Where("product.name_zh LIKE ?", "%"+c.NameZh+"%")
		}
		if c.MfgModel != "" {
			db.Where("product.mfg_model LIKE ?", "%"+c.MfgModel+"%")
		}
		if c.BrandZh != "" {
			db.Where("brand.brand_zh LIKE ?", "%"+c.BrandZh+"%")
		}
		if c.CreateDateStart != "" {
			db.Where("product.created_at > ?", c.CreateDateStart)
		}
		if c.CreateDateEnd != "" {
			db.Where("product.created_at < ?", c.CreateDateEnd)
		}
		return db
	}
}

type ProductExportAttrResp struct {
	SkuCode        string
	NameZh         string
	MfgModel       string
	CategoryId     int
	AttributeId    int
	AttributeName  string
	AttributeValue string
}

type InnerGetProductBySkuReq struct {
	SkuCode []string `json:"skuCode"`
}

type InnerGetProductBySkuResp struct {
	Id              int          `json:"id"`
	SkuCode         string       `json:"skuCode"`
	NameZh          string       `json:"nameZH"`
	MfgModel        string       `json:"mfgModel"`
	SalesUom        string       `json:"salesUom"`
	SalesMoq        int          `json:"salesMoq"`
	VendorId        int          `json:"vendorId"`
	SupplierSkuCode string       `json:"supplierSkuCode"`
	BrandId         int          `json:"brandId"`
	Brand           models.Brand `json:"brand"`
}

// ProductAttributeImportReq 产品属性维护导入
type ProductAttributeImportReq struct {
	CategoryLevel1 string `json:"categoryLevel1" comment:"一级目录"`
	CategoryLevel2 string `json:"categoryLevel2" comment:"二级目录"`
	CategoryLevel3 string `json:"categoryLevel3" comment:"三级目录"`
	CategoryLevel4 string `json:"categoryLevel4" comment:"四级目录"`
	CategoryId     int    `json:"categoryId" comment:"产线ID[终级目录]"`
	SkuCode        string `json:"skuCode" comment:"产品SKU" vd:"@:len($)>0; msg:'产品SKU必填'"`
	NameZh         string `json:"nameZH"  comment:"产品名称(中文)" vd:"@:len($)>0; msg:'产品名称(中文)必填'"`
	MfgModel       string `json:"mfgModel"  comment:"制造商型号" vd:"@:len($)>0; msg:'制造商型号必填'"`
}

func (model *ProductAttributeImportReq) Generate(k string, v any) {
	switch k {
	case "categoryLevel1":
		model.CategoryLevel1 = v.(string)
	case "categoryLevel2":
		model.CategoryLevel2 = v.(string)
	case "categoryLevel3":
		model.CategoryLevel3 = v.(string)
	case "categoryLevel4":
		model.CategoryLevel4 = v.(string)
	case "categoryId":
		if v == "" {
			model.CategoryId = 0
		} else {
			model.CategoryId, _ = strconv.Atoi(v.(string))
		}
	case "skuCode":
		model.SkuCode = v.(string)
	case "nameZh":
		model.NameZh = v.(string)
	case "mfgModel":
		model.MfgModel = v.(string)
	}
}

type GetInfoResp struct {
	models.Product
	ProductCategory        []CategoryPath  `json:"productCategory" gorm:"-"`
	ProductCategoryPrimary int             `json:"productCategoryPrimary" gorm:"-"`
	AttrList               *[]AttrsKeyName `json:"attrList" gorm:"-"`
}

type BatchUploadProductImageReq struct {
	ProductImage []MediaInstanceInsertReq `json:"productImage" comment:"商品图片"`
}

type InnerGetProductCategoryBySkuResp struct {
	SkuCode         string            `json:"skuCode"`
	ProductCategory []models.Category `json:"productCategory"`
}

type FindSkuIsBindReq struct {
	SkuCode         string
	VendorId        int
	SupplierSkuCode string
}

func FindSkuIsBind(c *FindSkuIsBindReq) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db.Where("sku_code = ?", c.SkuCode)
		db.Where("vendor_id = ?", c.VendorId)
		db.Where("supplier_sku_code = ?", c.SupplierSkuCode)
		return db
	}
}

// ProductSort 批量排序
type ProductSort struct {
	ProductId int `json:"productId"`
	Seq       int `json:"seq"`
}

type ProductSortReq struct {
	Sort []ProductSort `json:"sort"`
	common.ControlBy
}

type ProductUpdater struct {
	common.ControlBy
}
