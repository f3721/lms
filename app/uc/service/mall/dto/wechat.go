package dto

type UserLoginRequest struct {
	JsCode string `json:"jsCode,optional"` // 登录时获取的 code，可通过wx.login获取
}

type UserLoginResponse struct {
	SessionKey string `json:"sessionKey"` // 会话密钥
	UnionId    string `json:"unionId"`    // 用户在开放平台的唯一标识符，若当前小程序已绑定到微信开放平台帐号下会返回
	OpenId     string `json:"openId"`     // 用户唯一标识
	ErrMsg     string `json:"errMsg"`     // 错误信息
	ErrCode    int64  `json:"errCode"`    // 错误码
}

type PhoneNumberRequest struct {
	Code string `json:"code"` // 手机号获取凭证
}

type Watermark struct {
	Timestamp int32  `json:"timestamp"` // 用户获取手机号操作的时间戳
	Appid     string `json:"appid"`     // string
}

type PhoneInfo struct {
	PhoneNumber     string     `json:"phoneNumber"`     // 用户绑定的手机号（国外手机号会有区号）
	PurePhoneNumber string     `json:"purePhoneNumber"` // 没有区号的手机号
	CountryCode     string     `json:"countryCode"`     // 区号
	Watermark       *Watermark `json:"watermark"`       // 数据水印
}

type PhoneNumberResponse struct {
	ErrCode   int32      `json:"errCode"`   // 错误码
	ErrMsg    string     `json:"errMsg"`    // 错误信息
	PhoneInfo *PhoneInfo `json:"phoneInfo"` // 用户手机号信息
}

type CommonResp struct {
	Code    int64       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type CommonIdPathReq struct {
	Id int64 `path:"id"`
}
