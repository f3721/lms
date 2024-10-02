package sp

type SyncFailed struct {
	Token          string           `json:"token"`
	RejectionState string           `json:"rejectionState"`
	Skus           []SyncFailedSkus `json:"skus"`
}

type SyncFailedSkus struct {
	SkuId           string `json:"skuId"`
	RejectionReason string `json:"rejectionReason"`
	TenantName      string `json:"tenantName"`
}

type SyncFailedPushResp struct {
	Result bool `json:"result"`
	Response
}

type SyncFailedPush struct {
	SkuId           string
	RejectionReason string
	TenantName      string
	MessageId       string
}
