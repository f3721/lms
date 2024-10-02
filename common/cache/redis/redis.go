package redis

import (
	"github.com/go-admin-team/go-admin-core/sdk"
)

func GetCache(key string) (data string, err error) {
	data, err = sdk.Runtime.GetCacheAdapter().Get(key)
	return
}

func SetCache(key string, val interface{}, expire int) (err error) {
	err = sdk.Runtime.GetCacheAdapter().Set(key, val, expire)
	return
}

func DelCache(key string) (err error) {
	err = sdk.Runtime.GetCacheAdapter().Del(key)
	return
}
