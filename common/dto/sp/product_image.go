package sp

// GetProductImageReq 产品图片
type GetProductImageReq struct {
	Timestamp string   `json:"timestamp"`
	Token     string   `json:"token"`
	Sku       []string `json:"sku"`
}

type SkuPic struct {
	IsPrimary int    `json:"isPrimary"`
	OrderSort int    `json:"orderSort"`
	Path      string `json:"path"`
}

type ProductImageResult struct {
	SkuPic []SkuPic `json:"skuPic"`
	Sku    string   `json:"sku"`
}

type GetProductImageResp struct {
	Result []ProductImageResult `json:"result"`
	Response
}
