package sp

type Response struct {
	ResultCode    string `json:"resultCode"`
	ResultMessage string `json:"resultMessage"`
	Success       bool   `json:"success"`
}
