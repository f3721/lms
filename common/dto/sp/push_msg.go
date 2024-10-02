package sp

// GetMsgReq 消息列表

type GetMsgReq struct {
	Token string `json:"token"`
	Type  int64  `json:"type"`
}

type GetMsgResultResult struct {
	SkuId      string `json:"skuId"`
	State      int64  `json:"state"`
	PageNum    int64  `json:"page_num"`
	TenantId   string `json:"tenantId"`
	TenantName string `json:"tenantName"`
}

type GetMsgResult struct {
	Id     string             `json:"id"`
	Time   string             `json:"time"`
	Type   int                `json:"type"`
	Result GetMsgResultResult `json:"result"`
}

type GetMsgResp struct {
	Result []GetMsgResult `json:"result"`
	Response
}
