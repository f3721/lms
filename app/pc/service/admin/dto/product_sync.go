package dto

import common "go-admin/common/models"

type ProductSyncReq struct {
	VendorId               int                      `json:"vendorId" comment:"货主ID" vd:"@:$>0; msg:'货主必填'"`
	SupplierSkuCode        string                   `json:"supplierSkuCode" comment:"货主SKU" vd:"@:len($)>0; msg:'货主SKU必填'"`
	ProductCategory        []ProductCategory        `json:"productCategory" comment:"产品目录"`
	AttrList               []ProductExtAttribute    `json:"attrList" comment:"产品扩展属性"`
	ProductCategoryPrimary int                      `json:"productCategoryPrimary" comment:"商品主分类" vd:"@:$>0; msg:'请选择商品主分类'"`
	NameZh                 string                   `json:"nameZh" comment:"中文名" vd:"@:len($)>0; msg:'商品名称必填'"`
	BriefDesc              string                   `json:"briefDesc" comment:"产品描述" vd:"@:len($)<700000; msg:'产品描述字段超出长度限制'"`
	MfgModel               string                   `json:"mfgModel" comment:"制造厂型号" vd:"@:len($)>0; msg:'制造厂型号必填'"`
	BrandId                int                      `json:"brandId" comment:"品牌id" vd:"@:$>0; msg:'品牌必填'"`
	SalesUom               string                   `json:"salesUom" comment:"售卖包装单位" vd:"@:len($)>0; msg:'售卖包装单位必填'"`
	SalesPhysicalFactor    string                   `json:"salesPhysicalFactor" comment:"销售包装单位中含物理单位数量"`
	SalesMoq               int                      `json:"salesMoq" comment:"销售最小起订量" vd:"@:$>0; msg:'销售最小起订量至少为1'"`
	PackWeight             string                   `json:"packWeight" comment:"售卖包装质量(kg)"`
	Tax                    string                   `json:"tax" comment:"产品税率：默认空 值有（13%，9%，6%）" vd:"in($,'0.13','0.06','0.09'); msg:'产线税率的值有（0.13,0.06,0.09）'"`
	Status                 int                      `json:"status" comment:"产品状态" vd:"@:in($,1,2,3); msg:'产品状态只能为1或2或3'"`
	RefundFlag             int                      `json:"refundFlag" comment:"可退换货标志(0:否; 1:是)" vd:"@:in($,0,1); msg:'可退换货标志只能为0或1'"`
	ProductImage           []MediaInstanceInsertReq `json:"productImage" comment:"商品图片"`
	common.ControlBy
}
