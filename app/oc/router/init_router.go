package router

import (
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk"
	common "go-admin/common/middleware"
)

// InitRouter 路由初始化，不要怀疑，这里用到了
func InitRouter() {
	var r *gin.Engine
	h := sdk.Runtime.GetEngine()
	if h == nil {
		log.Fatal("not found engine...")
		os.Exit(-1)
	}
	switch h.(type) {
	case *gin.Engine:
		r = h.(*gin.Engine)
	default:
		log.Fatal("not support other engine")
		os.Exit(-1)
	}

	// the admin jwt middleware
	adminAuthMiddleware, err := common.AdminAuthInit()
	if err != nil {
		log.Fatalf("admin JWT Init Error, %s", err.Error())
	}
	InitAdminRouter(r, adminAuthMiddleware)

	// the mall jwt middleware
	mallAuthMiddleware, err := common.MallAuthInit()
	if err != nil {
		log.Fatalf("mall JWT Init Error, %s", err.Error())
	}
	InitMallRouter(r, mallAuthMiddleware)
}
