package middleware

import (
	"time"

	"go-admin/common/middleware/admin_handler"
	"go-admin/common/middleware/mall_handler"

	"github.com/go-admin-team/go-admin-core/sdk/config"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
)

// AdminAuthInit AuthInit jwt验证new
func AdminAuthInit() (*jwt.GinJWTMiddleware, error) {
	timeout := time.Hour
	if config.ApplicationConfig.Mode == "dev" {
		timeout = time.Duration(876010) * time.Hour
	} else {
		if config.JwtConfig.Timeout != 0 {
			timeout = time.Duration(config.JwtConfig.Timeout) * time.Second
		}
	}
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:           "test zone",
		Key:             []byte(config.JwtConfig.Secret),
		Timeout:         timeout,
		MaxRefresh:      time.Hour,
		PayloadFunc:     admin_handler.PayloadFunc,
		IdentityHandler: admin_handler.IdentityHandler,
		Authenticator:   admin_handler.Authenticator,
		Authorizator:    admin_handler.Authorizator,
		Unauthorized:    admin_handler.Unauthorized,
		TokenLookup:     "header: Authorization, query: token, cookie: jwt",
		TokenHeadName:   "Bearer",
		TimeFunc:        time.Now,
	})

}

// MallAuthInit AuthInit jwt验证new
func MallAuthInit() (*jwt.GinJWTMiddleware, error) {
	timeout := time.Hour
	if config.ApplicationConfig.Mode == "dev" {
		timeout = time.Duration(876010) * time.Hour
	} else {
		if config.JwtConfig.Timeout != 0 {
			timeout = time.Duration(config.JwtConfig.Timeout) * time.Second
		}
	}
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:           "test zone",
		Key:             []byte(config.JwtConfig.Secret+"mall2374"),
		Timeout:         timeout,
		MaxRefresh:      time.Hour,
		PayloadFunc:     mall_handler.PayloadFunc,
		IdentityHandler: mall_handler.IdentityHandler,
		Authenticator:   mall_handler.Authenticator,
		Authorizator:    mall_handler.Authorizator,
		Unauthorized:    mall_handler.Unauthorized,
		TokenLookup:     "header: Authorization, query: token, cookie: jwt",
		TokenHeadName:   "Bearer",
		TimeFunc:        time.Now,
	})

}
