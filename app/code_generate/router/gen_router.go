package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	"go-admin/app/code_generate/apis"
)

func init() {
	routerCheckRole = append(routerCheckRole, codeGenRouter)
}
func codeGenRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	r1 := v1.Group("/code")
	{
		gen := apis.Gen{}
		r1.GET("/gen", gen.CustomGen)
	}
}
