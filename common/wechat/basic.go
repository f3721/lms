package commonWechat

import (
	"context"
	sdkConfig "github.com/go-admin-team/go-admin-core/sdk/config"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/miniprogram"
	miniConfig "github.com/silenceper/wechat/v2/miniprogram/config"
	"go-admin/config"
)

func InitWechat() (*wechat.Wechat) {
	wc := wechat.NewWechat()
	redisOpts := &cache.RedisOpts{
		Host:    sdkConfig.CacheConfig.Redis.Addr,
		Password: sdkConfig.CacheConfig.Redis.Password,
	}
	redisCache := cache.NewRedis(context.Background(), redisOpts)
	wc.SetCache(redisCache)
	return wc
}

func GetMiniProgram() *miniprogram.MiniProgram {
	wc := InitWechat()
	memory := cache.NewMemory()
	cfg := &miniConfig.Config{
		AppID:     config.ExtConfig.MiniProgramConfig.AppID,
		AppSecret: config.ExtConfig.MiniProgramConfig.AppSecret,
		Cache:     memory,
	}
	miniProgram := wc.GetMiniProgram(cfg)
	return miniProgram
}
