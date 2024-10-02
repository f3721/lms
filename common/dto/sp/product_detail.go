package sp

// GetProductDetailReq 产品明细

type ProductParams struct {
	AttrCode      string `json:"attrCode"`
	AttrName      string `json:"attrName"`
	AttrValue     string `json:"attrValue"`
	AttrSeq       string `json:"attrSeq"`
	AttrGroupCode string `json:"attrGroupCode"`
	AttrGroupName string `json:"attrGroupName"`
	AttrGroupSeq  string `json:"attrGroupSeq"`
}

type ProductDetail struct {
	Sku          string          `json:"sku"`
	Weight       string          `json:"weight"`
	ImagePath    string          `json:"imagePath"`
	BrandName    string          `json:"brandName"`
	BrandPic     string          `json:"brandPic"`
	Name         string          `json:"name"`
	SaleUnit     string          `json:"saleUnit"`
	Category     []string        `json:"category"`
	Moq          int             `json:"moq"`
	MfgSku       string          `json:"mfgSku"`
	DeliveryTime int             `json:"deliveryTime"`
	IsReturn     int             `json:"isReturn"`
	Introduction string          `json:"introduction"`
	WareQD       string          `json:"wareQD"`
	SettleUnit   string          `json:"settleUnit"`
	WareNum      string          `json:"wareNum"`
	TaxCode      string          `json:"taxCode"`
	Param        []ProductParams `json:"param"`
}

type GetProductDetailReq struct {
	Timestamp string `json:"timestamp"`
	Token     string `json:"token"`
	Sku       string `json:"sku"`
}

type GetProductDetailResp struct {
	Result ProductDetail `json:"result"`
	Response
}
