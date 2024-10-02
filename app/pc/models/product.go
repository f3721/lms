package models

import (
	"go-admin/common/global"
	"go-admin/common/models"
	"gorm.io/gorm"
)

type Product struct {
	models.Model
	VendorId            int                  `json:"vendorId" gorm:"type:int unsigned;comment:货主ID"`
	VendorName          string               `json:"vendorName" gorm:"-" comment:"货主名称"`
	SupplierSkuCode     string               `json:"supplierSkuCode" gorm:"type:varchar(50);comment:货主SKU"`
	SkuCode             string               `json:"skuCode" gorm:"type:varchar(10);comment:产品sku"`
	NameZh              string               `json:"nameZh" gorm:"type:varchar(512);comment:中文名"`
	NameEn              string               `json:"nameEn" gorm:"type:varchar(512);comment:英文名"`
	Title               string               `json:"title" gorm:"type:varchar(1000);comment:标题"`
	BriefDesc           string               `json:"briefDesc" gorm:"type:varchar(4000);comment:产品描述"`
	MfgModel            string               `json:"mfgModel" gorm:"type:varchar(255);comment:制造厂型号"`
	BrandId             int                  `json:"brandId" gorm:"type:int unsigned;comment:品牌id"`
	SalesUom            string               `json:"salesUom" gorm:"type:varchar(255);comment:售卖包装单位"`
	PhysicalUom         string               `json:"physicalUom" gorm:"type:varchar(255);comment:物理单位"`
	SalesPhysicalFactor string               `json:"salesPhysicalFactor" gorm:"type:decimal(10,2);comment:销售包装单位中含物理单位数量"`
	SalesMoq            int                  `json:"salesMoq" gorm:"type:int unsigned;comment:销售最小起订量"`
	PackLength          string               `json:"packLength" gorm:"type:decimal(10,2);comment:售卖包装长(mm)"`
	PackWidth           string               `json:"packWidth" gorm:"type:decimal(10,2);comment:售卖包装宽(mm)"`
	PackHeight          string               `json:"packHeight" gorm:"type:decimal(10,2);comment:售卖包装高(mm)"`
	PackWeight          string               `json:"packWeight" gorm:"type:decimal(10,3);comment:售卖包装质量(kg)"`
	Accessories         string               `json:"accessories" gorm:"type:varchar(500);comment:包装清单"`
	ProductAlias        string               `json:"productAlias" gorm:"type:varchar(255);comment:产品别名"`
	FragileFlag         int                  `json:"fragileFlag" gorm:"type:tinyint(1);comment:易碎标志(0:否;1:是)"`
	HazardFlag          int                  `json:"hazardFlag" gorm:"type:tinyint(1);comment:危险品标志(0:否;1:是)"`
	HazardClass         int                  `json:"hazardClass" gorm:"type:tinyint;comment:危险品等级"`
	BulkyFlag           int                  `json:"bulkyFlag" gorm:"type:tinyint(1);comment:抛货标志(0:否; 1:是)"`
	AssembleFlag        int                  `json:"assembleFlag" gorm:"type:tinyint(1);comment:拼装件标志(0:否; 1:是)"`
	IsValuables         int                  `json:"isValuables" gorm:"type:tinyint(1);comment:是否贵重品(0:否; 1:是)"`
	IsFluid             int                  `json:"isFluid" gorm:"type:tinyint(1);comment:是否液体(0:否; 1:是)"`
	Status              int                  `json:"status" gorm:"type:tinyint(1);comment:产品状态"`
	ConsumptiveFlag     int                  `json:"consumptiveFlag" gorm:"type:tinyint(1);comment:耗材标志"`
	StorageFlag         int                  `json:"storageFlag" gorm:"type:tinyint(1);comment:保存期标志"`
	StorageTime         int                  `json:"storageTime" gorm:"type:tinyint;comment:保存期限(月)"`
	CustomMadeFlag      int                  `json:"customMadeFlag" gorm:"type:tinyint(1);comment:定制品标志(0:否,1:是)"`
	RefundFlag          int                  `json:"refundFlag" gorm:"type:tinyint(1);comment:可退换货(0:否,1:是)"`
	Barcode             string               `json:"barcode" gorm:"type:varchar(50);comment:商品条形码"`
	OtherLabel          string               `json:"otherLabel" gorm:"type:varchar(50);comment:其他备注"`
	Tax                 string               `json:"tax" gorm:"type:varchar(4);comment:产品税率：默认空 值有（13%，9%，6%）"`
	CertificateImage    string               `json:"certificateImage" gorm:"type:varchar(200);comment:合格证"`
	Seq                 int                  `json:"seq" gorm:"type:smallint unsigned;comment:排序"`
	Brand               Brand                `json:"brand" gorm:"<-:false;"`
	MediaRelationship   *[]MediaRelationship `json:"mediaRelationship" gorm:"foreignKey:BuszId;references:SkuCode"`
	models.ModelTime
	models.ControlBy
}

func (Product) TableName() string {
	return "product"
}

func (e *Product) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *Product) GetId() interface{} {
	return e.Id
}

func (e *Product) SearchSkuByKeyword(tx *gorm.DB, keyword string, result *[]map[string]string) error {
	var model Product
	err := tx.Model(&model).
		Joins("LEFT JOIN brand ON product.brand_id = brand.id").
		Where("product.name_zh LIKE ? OR product.mfg_model LIKE ? OR brand.brand_zh LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%").
		Where("product.status = 2").
		Select("product.sku_code").
		Find(result).Error
	return err
}

func (e *Product) GetTaxBySku(tx *gorm.DB, skus []string) (res map[string]string) {
	pcPrefix := global.GetTenantPcDBNameWithDB(tx)
	var list []Product
	tx.Table(pcPrefix+"."+e.TableName()).Select("sku_code, tax").Where("sku_code in ?", skus).Find(&list)
	res = make(map[string]string)
	for _, product := range list {
		res[product.SkuCode] = product.Tax
	}
	return
}
