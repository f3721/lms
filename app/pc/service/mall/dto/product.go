package dto

import (

	"go-admin/app/pc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type ProductGetPageReq struct {
	dto.Pagination     `search:"-"`
    ProductOrder
}

type ProductOrder struct {
    Id string `form:"idOrder"  search:"type:order;column:id;table:product"`
    VendorId string `form:"vendorIdOrder"  search:"type:order;column:vendor_id;table:product"`
    SupplierSkuCode string `form:"supplierSkuCodeOrder"  search:"type:order;column:supplier_sku_code;table:product"`
    SkuCode string `form:"skuCodeOrder"  search:"type:order;column:sku_code;table:product"`
    NameZh string `form:"nameZhOrder"  search:"type:order;column:name_zh;table:product"`
    NameEn string `form:"nameEnOrder"  search:"type:order;column:name_en;table:product"`
    Title string `form:"titleOrder"  search:"type:order;column:title;table:product"`
    BriefDesc string `form:"briefDescOrder"  search:"type:order;column:brief_desc;table:product"`
    MfgModel string `form:"mfgModelOrder"  search:"type:order;column:mfg_model;table:product"`
    BrandId string `form:"brandIdOrder"  search:"type:order;column:brand_id;table:product"`
    SalesUom string `form:"salesUomOrder"  search:"type:order;column:sales_uom;table:product"`
    PhysicalUom string `form:"physicalUomOrder"  search:"type:order;column:physical_uom;table:product"`
    SalesPhysicalFactor string `form:"salesPhysicalFactorOrder"  search:"type:order;column:sales_physical_factor;table:product"`
    SalesMoq string `form:"salesMoqOrder"  search:"type:order;column:sales_moq;table:product"`
    PackLength string `form:"packLengthOrder"  search:"type:order;column:pack_length;table:product"`
    PackWidth string `form:"packWidthOrder"  search:"type:order;column:pack_width;table:product"`
    PackHeight string `form:"packHeightOrder"  search:"type:order;column:pack_height;table:product"`
    PackWeight string `form:"packWeightOrder"  search:"type:order;column:pack_weight;table:product"`
    Accessories string `form:"accessoriesOrder"  search:"type:order;column:accessories;table:product"`
    ProductAlias string `form:"productAliasOrder"  search:"type:order;column:product_alias;table:product"`
    FragileFlag string `form:"fragileFlagOrder"  search:"type:order;column:fragile_flag;table:product"`
    HazardFlag string `form:"hazardFlagOrder"  search:"type:order;column:hazard_flag;table:product"`
    BulkyFlag string `form:"bulkyFlagOrder"  search:"type:order;column:bulky_flag;table:product"`
    AssembleFlag string `form:"assembleFlagOrder"  search:"type:order;column:assemble_flag;table:product"`
    IsValuables string `form:"isValuablesOrder"  search:"type:order;column:is_valuables;table:product"`
    IsFluid string `form:"isFluidOrder"  search:"type:order;column:is_fluid;table:product"`
    Status string `form:"statusOrder"  search:"type:order;column:status;table:product"`
    ConsumptiveFlag string `form:"consumptiveFlagOrder"  search:"type:order;column:consumptive_flag;table:product"`
    StorageFlag string `form:"storageFlagOrder"  search:"type:order;column:storage_flag;table:product"`
    StorageTime string `form:"storageTimeOrder"  search:"type:order;column:storage_time;table:product"`
    CustomMadeFlag string `form:"customMadeFlagOrder"  search:"type:order;column:custom_made_flag;table:product"`
    OtherLabel string `form:"otherLabelOrder"  search:"type:order;column:other_label;table:product"`
    Tax string `form:"taxOrder"  search:"type:order;column:tax;table:product"`
    CreateBy string `form:"createByOrder"  search:"type:order;column:create_by;table:product"`
    CreateByName string `form:"createByNameOrder"  search:"type:order;column:create_by_name;table:product"`
    UpdateBy string `form:"updateByOrder"  search:"type:order;column:update_by;table:product"`
    UpdateByName string `form:"updateByNameOrder"  search:"type:order;column:update_by_name;table:product"`
    CreatedAt string `form:"createdAtOrder"  search:"type:order;column:created_at;table:product"`
    UpdatedAt string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:product"`
    DeletedAt string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:product"`
    
}

func (m *ProductGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type ProductInsertReq struct {
    Id int `json:"-" comment:"产品id"` // 产品id
    VendorId int `json:"vendorId" comment:"货主ID"`
    SupplierSkuCode string `json:"supplierSkuCode" comment:"货主SKU"`
    SkuCode string `json:"skuCode" comment:"产品sku"`
    NameZh string `json:"nameZh" comment:"中文名"`
    NameEn string `json:"nameEn" comment:"英文名"`
    Title string `json:"title" comment:"标题"`
    BriefDesc string `json:"briefDesc" comment:"产品描述"`
    MfgModel string `json:"mfgModel" comment:"制造厂型号"`
    BrandId int `json:"brandId" comment:"品牌id"`
    SalesUom string `json:"salesUom" comment:"售卖包装单位"`
    PhysicalUom string `json:"physicalUom" comment:"物理单位"`
    SalesPhysicalFactor string `json:"salesPhysicalFactor" comment:"销售包装单位中含物理单位数量"`
    SalesMoq int `json:"salesMoq" comment:"销售最小起订量"`
    PackLength string `json:"packLength" comment:"售卖包装长(mm)"`
    PackWidth string `json:"packWidth" comment:"售卖包装宽(mm)"`
    PackHeight string `json:"packHeight" comment:"售卖包装高(mm)"`
    PackWeight string `json:"packWeight" comment:"售卖包装质量(kg)"`
    Accessories string `json:"accessories" comment:"包装清单"`
    ProductAlias string `json:"productAlias" comment:"产品别名"`
    FragileFlag int `json:"fragileFlag" comment:"易碎标志(0:否;1:是)"`
    HazardFlag int `json:"hazardFlag" comment:"危险品标志(0:否;1:是)"`
    BulkyFlag int `json:"bulkyFlag" comment:"抛货标志(0:否; 1:是)"`
    AssembleFlag int `json:"assembleFlag" comment:"拼装件标志(0:否; 1:是)"`
    IsValuables int `json:"isValuables" comment:"是否贵重品(0:否; 1:是)"`
    IsFluid int `json:"isFluid" comment:"是否液体(0:否; 1:是)"`
    Status int `json:"status" comment:"产品状态"`
    ConsumptiveFlag int `json:"consumptiveFlag" comment:"耗材标志"`
    StorageFlag int `json:"storageFlag" comment:"保存期标志"`
    StorageTime int `json:"storageTime" comment:"保存期限(月)"`
    CustomMadeFlag int `json:"customMadeFlag" comment:"定制品标志(0:否,1:是)"`
    OtherLabel string `json:"otherLabel" comment:"其他备注"`
    Tax string `json:"tax" comment:"产品税率：默认空 值有（13%，9%，6%）"`
    common.ControlBy
}

func (s *ProductInsertReq) Generate(model *models.Product)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
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
    model.BulkyFlag = s.BulkyFlag
    model.AssembleFlag = s.AssembleFlag
    model.IsValuables = s.IsValuables
    model.IsFluid = s.IsFluid
    model.Status = s.Status
    model.ConsumptiveFlag = s.ConsumptiveFlag
    model.StorageFlag = s.StorageFlag
    model.StorageTime = s.StorageTime
    model.CustomMadeFlag = s.CustomMadeFlag
    model.OtherLabel = s.OtherLabel
    model.Tax = s.Tax
    model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
    model.CreateByName = s.CreateByName
}

func (s *ProductInsertReq) GetId() interface{} {
	return s.Id
}

type ProductUpdateReq struct {
    Id int `uri:"id" comment:"产品id"` // 产品id
    VendorId int `json:"vendorId" comment:"货主ID"`
    SupplierSkuCode string `json:"supplierSkuCode" comment:"货主SKU"`
    SkuCode string `json:"skuCode" comment:"产品sku"`
    NameZh string `json:"nameZh" comment:"中文名"`
    NameEn string `json:"nameEn" comment:"英文名"`
    Title string `json:"title" comment:"标题"`
    BriefDesc string `json:"briefDesc" comment:"产品描述"`
    MfgModel string `json:"mfgModel" comment:"制造厂型号"`
    BrandId int `json:"brandId" comment:"品牌id"`
    SalesUom string `json:"salesUom" comment:"售卖包装单位"`
    PhysicalUom string `json:"physicalUom" comment:"物理单位"`
    SalesPhysicalFactor string `json:"salesPhysicalFactor" comment:"销售包装单位中含物理单位数量"`
    SalesMoq int `json:"salesMoq" comment:"销售最小起订量"`
    PackLength string `json:"packLength" comment:"售卖包装长(mm)"`
    PackWidth string `json:"packWidth" comment:"售卖包装宽(mm)"`
    PackHeight string `json:"packHeight" comment:"售卖包装高(mm)"`
    PackWeight string `json:"packWeight" comment:"售卖包装质量(kg)"`
    Accessories string `json:"accessories" comment:"包装清单"`
    ProductAlias string `json:"productAlias" comment:"产品别名"`
    FragileFlag int `json:"fragileFlag" comment:"易碎标志(0:否;1:是)"`
    HazardFlag int `json:"hazardFlag" comment:"危险品标志(0:否;1:是)"`
    BulkyFlag int `json:"bulkyFlag" comment:"抛货标志(0:否; 1:是)"`
    AssembleFlag int `json:"assembleFlag" comment:"拼装件标志(0:否; 1:是)"`
    IsValuables int `json:"isValuables" comment:"是否贵重品(0:否; 1:是)"`
    IsFluid int `json:"isFluid" comment:"是否液体(0:否; 1:是)"`
    Status int `json:"status" comment:"产品状态"`
    ConsumptiveFlag int `json:"consumptiveFlag" comment:"耗材标志"`
    StorageFlag int `json:"storageFlag" comment:"保存期标志"`
    StorageTime int `json:"storageTime" comment:"保存期限(月)"`
    CustomMadeFlag int `json:"customMadeFlag" comment:"定制品标志(0:否,1:是)"`
    OtherLabel string `json:"otherLabel" comment:"其他备注"`
    Tax string `json:"tax" comment:"产品税率：默认空 值有（13%，9%，6%）"`
    common.ControlBy
}

func (s *ProductUpdateReq) Generate(model *models.Product)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
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
    model.BulkyFlag = s.BulkyFlag
    model.AssembleFlag = s.AssembleFlag
    model.IsValuables = s.IsValuables
    model.IsFluid = s.IsFluid
    model.Status = s.Status
    model.ConsumptiveFlag = s.ConsumptiveFlag
    model.StorageFlag = s.StorageFlag
    model.StorageTime = s.StorageTime
    model.CustomMadeFlag = s.CustomMadeFlag
    model.OtherLabel = s.OtherLabel
    model.Tax = s.Tax
    model.UpdateBy = s.UpdateBy
    model.UpdateByName = s.UpdateByName // 添加这而，需要记录是被谁更新的
}

func (s *ProductUpdateReq) GetId() interface{} {
	return s.Id
}

// ProductGetReq 功能获取请求参数
type ProductGetReq struct {
     Id int `uri:"id"`
}
func (s *ProductGetReq) GetId() interface{} {
	return s.Id
}

// ProductDeleteReq 功能删除请求参数
type ProductDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *ProductDeleteReq) GetId() interface{} {
	return s.Ids
}
