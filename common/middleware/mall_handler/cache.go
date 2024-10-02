package mall_handler

import (
	"github.com/go-admin-team/go-admin-core/sdk"
	"github.com/go-admin-team/go-admin-core/sdk/config"
)

func GetCache(key string) (data string, err error)  {
	data, err = sdk.Runtime.GetCacheAdapter().Get(GetMallCacheKeyPrefix()+key)
	return
}

func SetCache(key string, val interface{}, expire int) (err error)  {
	err = sdk.Runtime.GetCacheAdapter().Set(GetMallCacheKeyPrefix()+key, val, expire)
	return
}

func DelCache(key string) (err error)  {
	err = sdk.Runtime.GetCacheAdapter().Del(GetMallCacheKeyPrefix()+key)
	return
}

func GetMallCacheKeyPrefix() string {
	return "SXYZ"+ config.ApplicationConfig.Mode + "MALL-"
}