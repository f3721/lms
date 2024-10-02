package admin_handler

import (
	"github.com/go-admin-team/go-admin-core/sdk"
	"github.com/go-admin-team/go-admin-core/sdk/config"
)

func GetCache(key string) (data string, err error)  {
	data, err = sdk.Runtime.GetCacheAdapter().Get(GetAdminCacheKeyPrefix()+key)
	return
}

func SetCache(key string, val interface{}, expire int) (err error)  {
	err = sdk.Runtime.GetCacheAdapter().Set(GetAdminCacheKeyPrefix()+key, val, expire)
	return
}

func DelCache(key string) (err error)  {
	err = sdk.Runtime.GetCacheAdapter().Del(GetAdminCacheKeyPrefix()+key)
	return
}

func GetAdminCacheKeyPrefix() string {
	return "SXYZ"+ config.ApplicationConfig.Mode + "LMS-"
}