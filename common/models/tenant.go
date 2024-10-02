package models

import (
	"errors"
	"fmt"
)

type HS map[string]string

type LubanResp struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
}
type LubanToken struct {
	Scope string `json:"scope"`
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokeType string `json:"toke_type"`
	ExpiresIn string `json:"expires_in"`
}

type SystemTenantResponse struct {
	Data []SystemTenant `json:"data"`
}

type SystemTenant struct {
	ID            	  int    `json:"id"`
	TenantId      	  int    `json:"tenantId"`
	Name          	  string `json:"name"`
	SystemName    	  string `json:"systemName"`
	ContactName   	  string `json:"contactName"`
	ContactMobile 	  string `json:"contactMobile"`
	Status        	  int    `json:"status" comment:"0-开启 1-关闭"`
	Domain        	  string `json:"domain"`
	PackageID     	  int    `json:"packageId"`
	ExpireTime    	  string `json:"expireTime"`
	AccountCount  	  int    `json:"accountCount"`
	CreateTime    	  string `json:"createTime"`
	UpdateTime    	  string `json:"updateTime"`

	DatabaseName     string `json:"databaseName"`
	DatabaseUsername string `json:"databaseUsername"`
	DatabasePassword string `json:"databasePassword"`
	DatabaseHost     string `json:"databaseHost"`
	DatabasePort     int 	`json:"databasePort"`
	DatabaseInitOk   int    `json:"databaseInitOk"`
}

func (s *SystemTenant) GetDBSourceWithAppName(appName string) (string, error) {
	if s.DatabaseInitOk == 0 {
		return "", errors.New("租户数据库未初始化")
	}
	if s.DatabaseUsername == "" {
		return "", errors.New("租户数据库未填写用户")
	}
	if s.DatabasePassword == "" {
		return "", errors.New("租户数据库未填写密码")
	}
	if s.DatabaseHost == "" {
		return "", errors.New("租户数据库未填写域名")
	}
	if s.DatabasePort == 0 {
		return "", errors.New("租户数据库未填写端口")
	}
	if s.DatabaseName == "" {
		return "", errors.New("租户数据库未填写前缀")
	}
	source := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v_%v?charset=utf8mb4&parseTime=True&loc=Local&timeout=1000ms",
		s.DatabaseUsername,
		s.DatabasePassword,
		s.DatabaseHost,
		s.DatabasePort,
		s.TenantDBPrefix(),
		appName)
	return source, nil
}

func (s *SystemTenant) TenantDBPrefix() string {
	return fmt.Sprintf("%v", s.DatabaseName)
}
