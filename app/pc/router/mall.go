package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	apis "go-admin/app/pc/apis/mall"
	"go-admin/common/actions"
)

func init() {
	routerMallCheckRole = append(routerMallCheckRole, registerMallApiRouter)
}

// registerMallApiRouter
func registerMallApiRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	//用户验证
	//r := v1.Group("product").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	//{
	//	r.POST("test", func(context *gin.Context) {
	//		userId := user.GetUserIdStr(context)
	//		context.JSON(200, gin.H{"status": "user validated, test success! userId is " + userId})
	//	})
	//}

	category := apis.Category{}
	r := v1.Group("/category").Use(authMiddleware.MiddlewareFunc()) //.Use(middleware.AuthCheckRole())
	{
		r.GET("", actions.PermissionAction(), category.GetPage)
		r.GET("/:id", actions.PermissionAction(), category.Get)
		r.GET("/list/:id", actions.PermissionAction(), category.GetList)
	}

	brand := apis.Brand{}
	r = v1.Group("/brand").Use(authMiddleware.MiddlewareFunc()) //.Use(middleware.AuthCheckRole())
	{
		r.GET("", actions.PermissionAction(), brand.GetPage)
		r.GET("/:id", actions.PermissionAction(), brand.Get)
	}

	goods := apis.Goods{}
	r = v1.Group("/product").Use(authMiddleware.MiddlewareFunc()) //.Use(middleware.AuthCheckRole())
	{
		r.GET("", actions.PermissionAction(), goods.GetPage)
		r.GET("/:id", actions.PermissionAction(), goods.Get)
		r.GET("/filter-items", actions.PermissionAction(), goods.GetMiniProgramHomeFilter)
	}

	apiCart := apis.UserCart{}
	rCart := v1.Group("/user-cart").Use(authMiddleware.MiddlewareFunc()) /*.Use(middleware.AuthCheckRole())*/
	{
		rCart.GET("", actions.PermissionAction(), apiCart.GetPage)
		rCart.GET("order-cart-products", actions.PermissionAction(), apiCart.GetPageForOrder)
		rCart.GET("order-buy-now-product", actions.PermissionAction(), apiCart.GetProductForOrderOnBuyNow)
		rCart.POST("", apiCart.Insert)
		rCart.PUT("", actions.PermissionAction(), apiCart.Update)
		rCart.POST("delete", apiCart.Delete)
		rCart.POST("batch-add", apiCart.BatchAdd)
		rCart.POST("select-one", apiCart.SelectOne)
		rCart.POST("select-all", apiCart.SelectAll)
		rCart.POST("unselect-all", apiCart.UnSelectAll)
		rCart.POST("clear-select", apiCart.ClearSelect)
		rCart.POST("clear-invalid", apiCart.ClearInvalid)
		rCart.GET("verify-cart", apiCart.VerifyCart)
		rCart.GET("verify-buy-now", apiCart.VerifyBuyNow)
		rCart.POST("sale-moq", apiCart.SaleMoq)
	}

	userSearchHistory := apis.UserSearchHistory{}
	r = v1.Group("/user-search-history").Use(authMiddleware.MiddlewareFunc()) //.Use(middleware.AuthCheckRole())
	{
		r.GET("", actions.PermissionAction(), userSearchHistory.GetPage)
	}
}
