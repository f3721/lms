package sp

import (
	"encoding/json"
	"errors"
	"go-admin/common/cache/redis"
	"go-admin/common/dto/sp"
	"go-admin/common/httplib"
	"go-admin/common/utils"
)

var clientId string
var clientSecret string
var username string
var password string
var tokenKey string
var accessToken string

type AccessToken struct {
	GrantType    string `json:"grant_type"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Timestamp    string `json:"timestamp"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	SpApiUrl     string `json:"sp_api_url"`
}

func New(clientId string, clientSecret string, username string, password string, spApiUrl string) *AccessToken {
	spClient := AccessToken{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		Username:     username,
		Password:     password,
		SpApiUrl:     spApiUrl,
	}
	tokenKey = utils.EncryptMd5(clientId, username)
	spClient.GetToken()
	return &spClient
}

func (e *AccessToken) GetToken() error {
	val, _ := redis.GetCache(tokenKey)
	if val != "" {
		accessToken = val
		return nil
	} else {
		_, err := e.GetAccessToken()
		if err != nil {
			return err
		}
		return nil
	}
}

func (e *AccessToken) GetAccessToken() (result *sp.GetAccessTokenResp, err error) {
	req := sp.GetAccessTokenReq{
		ClientId:     e.ClientId,
		ClientSecret: e.ClientSecret,
		Timestamp:    utils.GetDateTimeString(),
		Username:     e.Username,
		Password:     e.Password,
	}
	data, err := post(e.SpApiUrl+"/accessToken", req)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(data), &result)
	if err != nil {
		return nil, err
	}
	if result.ResultCode == "0000" {
		accessToken = result.Result.AccessToken
		_ = redis.SetCache(tokenKey, result.Result.AccessToken, result.Result.ExpiresIn-100)
		return nil, nil
	}
	return nil, errors.New(result.ResultMessage)
}

func (e *AccessToken) RefreshAccessToken(in *sp.RefreshTokenReq) (result *sp.RefreshTokenResp, err error) {
	req := httplib.Post(e.SpApiUrl + "/refreshToken")
	req.Param("client_id", in.ClientId)
	req.Param("client_secret", in.ClientSecret)
	req.Param("refresh_token", in.RefreshToken)
	data, err := req.String()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(data), &result)
	if err != nil {
		return nil, err
	}
	return
}

// GetPushMsg 消息列表
func (e *AccessToken) GetPushMsg(in *sp.GetMsgReq) (result *sp.GetMsgResp, err error) {
	in.Token = accessToken
	data, err := post(e.SpApiUrl+"/mall/get", in)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(data), &result)
	if err != nil {
		return nil, err
	}
	if result.ResultCode == "2007" {
		_, err = e.GetAccessToken()
		if err != nil {
			return nil, err
		}
		e.GetPushMsg(in)
	}
	return
}

// GetProductDetail 商品详情
func (e *AccessToken) GetProductDetail(in *sp.GetProductDetailReq) (result *sp.GetProductDetailResp, err error) {
	in.Token = accessToken
	in.Timestamp = utils.GetDateTimeString()
	data, err := post(e.SpApiUrl+"/mall/getDetail", in)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(data), &result)
	if err != nil {
		return nil, err
	}
	if result.ResultCode == "2007" {
		_, err = e.GetAccessToken()
		if err != nil {
			return nil, err
		}
		e.GetProductDetail(in)
	}
	return
}

// GetProductImage 商品详情
func (e *AccessToken) GetProductImage(in *sp.GetProductImageReq) (result *sp.GetProductImageResp, err error) {
	in.Token = accessToken
	in.Timestamp = utils.GetDateTimeString()
	data, err := post(e.SpApiUrl+"/mall/skuImage", in)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(data), &result)
	if err != nil {
		return nil, err
	}
	if result.ResultCode == "2007" {
		_, err = e.GetAccessToken()
		if err != nil {
			return nil, err
		}
		e.GetProductImage(in)
	}
	return
}

// SyncFailedPush 商品详情
func (e *AccessToken) SyncFailedPush(in *sp.SyncFailed) (result *sp.SyncFailedPushResp, err error) {
	in.Token = accessToken
	data, err := post(e.SpApiUrl+"/mall/syncRejectionReason", in)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(data), &result)
	if err != nil {
		return nil, err
	}
	if result.ResultCode == "2007" {
		_, err = e.GetAccessToken()
		if err != nil {
			return nil, err
		}
		e.SyncFailedPush(in)
	}
	return
}

// MessageDelete 商品详情
func (e *AccessToken) MessageDelete(in *sp.MessageDelete) (result *sp.MessageDeleteResp, err error) {
	in.Token = accessToken
	data, err := post(e.SpApiUrl+"/mall/delete", in)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(data), &result)
	if err != nil {
		return nil, err
	}
	if result.ResultCode == "2007" {
		_, err = e.GetAccessToken()
		if err != nil {
			return nil, err
		}
		e.MessageDelete(in)
	}
	return
}

func post(url string, in any) (data string, err error) {
	req := httplib.Post(url)
	datas, err := json.Marshal(in)
	req.Body(datas)
	req.Header("Content-Type", "application/json")
	data, err = req.String()
	if err != nil {
		return "", err
	}
	return
}
