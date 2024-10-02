package sp

type GetAccessTokenReq struct {
	GrantType    string `json:"grant_type"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Timestamp    string `json:"timestamp"`
	Username     string `json:"username"`
	Password     string `json:"password"`
}

type GetAccessTokenResult struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Time         string `json:"time"`
	ExpiresIn    int    `json:"expires_in"`
}

type GetAccessTokenResp struct {
	Result GetAccessTokenResult `json:"result"`
	Response
}

type RefreshTokenReq struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResp struct {
	Result GetAccessTokenResult `json:"result"`
	Response
}
