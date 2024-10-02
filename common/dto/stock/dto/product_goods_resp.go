package dto

import (
	modelsPc "go-admin/app/pc/models"
)

type ProductGoodsResp struct {
	GoodId        int    `json:"goodId"`
	ProductName   string `json:"productName"`
	MfgModel      string `json:"mfgModel"`
	BrandName     string `json:"brandName"`
	SalesUom      string `json:"salesUom"`
	VendorName    string `json:"vendorName"`
	ProductNo     string `json:"productNo"`
	VendorSkuCode string `json:"vendorSkuCode"`
}

func (p *ProductGoodsResp) FillProductGoodsRespData(goodsId int, goodsProductMap map[int]modelsPc.Goods) {
	if data, ok := goodsProductMap[goodsId]; ok {
		p.ProductName = data.Product.NameZh
		p.MfgModel = data.Product.MfgModel
		p.BrandName = data.Product.Brand.BrandZh
		p.SalesUom = data.Product.SalesUom
		p.ProductNo = data.ProductNo
		p.VendorSkuCode = data.Product.SupplierSkuCode
		p.GoodId = data.Id
	}
}
