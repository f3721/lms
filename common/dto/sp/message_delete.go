package sp

type MessageDelete struct {
	Token string `json:"token"`
	Id    string `json:"id"`
}

type MessageDeleteResp struct {
	Result bool `json:"result"`
	Response
}
