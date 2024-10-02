package router

import (
	apis "go-admin/app/uc/apis/mall"
	"go-admin/common/actions"
	"go-admin/common/middleware/mall_handler"

	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
)

func init() {
	routerMallNoCheckRole = append(routerMallNoCheckRole, registerMallNoAuthApiRouter)
	routerMallCheckRole = append(routerMallCheckRole, registerMallApiRouter)
}

// registerMallNoAuthApiRouter
func registerMallNoAuthApiRouter(v1 *gin.RouterGroup) {
	// 找回密码相关
	api := apis.UserInfo{}
	r := v1.Group("user")
	{
		r.GET("/captcha", api.Captcha)
		r.POST("/check-phone", api.CheckPhone)
		r.POST("/verify-code", api.VerifyCode)
		r.POST("/comit-verify-code", api.ComitVerifyCode)
		r.POST("/send-email", api.SendEmail)
		r.POST("/change-pwd", api.ChangePwd)
		r.POST("/proxy-login", api.ProxyLogin)
		r.POST("/mini-login", api.MiniLogin)
	}

	// 找回密码相关
	apiWechat := apis.Wechat{}
	r = v1.Group("wechat")
	{
		r.POST("/login", apiWechat.Login)
		r.POST("/phone-number", apiWechat.PhoneNumber)
	}
}

// registerMallApiRouter
func registerMallApiRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	// 登录
	v1.POST("/login", authMiddleware.LoginHandler)
	// Refresh time can be longer than token timeout
	v1.GET("/refresh_token", authMiddleware.RefreshHandler)
	v1.GET("/logout", mall_handler.LogOut)

	// 用户相关接口
	api := apis.UserInfo{}
	r := v1.Group("user").Use(authMiddleware.MiddlewareFunc())
	{
		r.POST("test", func(context *gin.Context) {
			context.JSON(200, gin.H{"status": "user validated, test success!"})
		})
		r.POST("select-warehouse", apis.UserConfigAPI{}.SelectWarehouse)
		r.GET("/info", actions.PermissionAction(), api.GetUserInfo)
		r.PUT("change-basic-info", actions.PermissionAction(), api.ChangeBasicInfo)
		r.PUT("change-password", actions.PermissionAction(), api.ChangePassword)
		r.POST("send-sms", actions.PermissionAction(), api.SendSms)
		r.POST("send-email-for-change-email", actions.PermissionAction(), api.SendEmailForChangeEmail)
		r.POST("send-email-for-check-email", actions.PermissionAction(), api.SendEmailForCheckEmail)
		r.POST("check-send-email-token", actions.PermissionAction(), api.CheckUserTokenAvailable)
		r.PUT("change-user-email", actions.PermissionAction(), api.ChangeUserEmail)
		r.PUT("change-phone", actions.PermissionAction(), api.ChangePhone)

		r.GET("", actions.PermissionAction(), api.GetPage)
		r.GET("/:id", actions.PermissionAction(), api.Get)
		r.POST("", api.Insert)
		r.PUT("/:id", actions.PermissionAction(), api.Update)
		r.DELETE("", api.Delete)
	}

	// 用户收藏
	apiUserCollect := apis.UserCollect{}
	r = v1.Group("/user-collect").Use(authMiddleware.MiddlewareFunc())
	{
		r.GET("", actions.PermissionAction(), apiUserCollect.GetPage)
		r.GET("/:id", actions.PermissionAction(), apiUserCollect.Get)
		r.POST("", apiUserCollect.Insert)
		//r.PUT("/:id", actions.PermissionAction(), apiUserCollect.Update)
		r.DELETE("", apiUserCollect.Delete)
	}

	// 公司收藏
	r = v1.Group("/company-collect").Use(authMiddleware.MiddlewareFunc())
	{
		r.GET("", actions.PermissionAction(), apiUserCollect.GetCompanyPage)
	}

	// 管理中心菜单相关
	apiMenu := apis.ManageMenu{}
	r = v1.Group("/manage-menu").Use(authMiddleware.MiddlewareFunc())
	{
		r.GET("", actions.PermissionAction(), apiMenu.GetPage)
		r.GET("/:id", actions.PermissionAction(), apiMenu.Get)
		r.POST("", apiMenu.Insert)
		r.PUT("/:id", actions.PermissionAction(), apiMenu.Update)
		r.DELETE("", apiMenu.Delete)
	}

	// 用户地址相关接口
	apiAddress := apis.Address{}
	r = v1.Group("address").Use(authMiddleware.MiddlewareFunc())
	{
		r.GET("", actions.PermissionAction(), apiAddress.GetPage)
		r.GET("/:id", actions.PermissionAction(), apiAddress.Get)
		r.POST("", apiAddress.Insert)
		r.PUT("/:id", actions.PermissionAction(), apiAddress.Update)
		r.PUT("/default/:id", actions.PermissionAction(), apiAddress.Default)
		r.DELETE("", apiAddress.Delete)
	}

	// 审批时间设置
	apiApprovalTime := apis.ApprovalTime{}
	r = v1.Group("approval-time").Use(authMiddleware.MiddlewareFunc())
	{
		r.GET("", actions.PermissionAction(), apiApprovalTime.GetPage)
		r.GET("/:id", actions.PermissionAction(), apiApprovalTime.Get)
		r.POST("", apiApprovalTime.Insert)
		r.PUT("/:id", actions.PermissionAction(), apiApprovalTime.Update)
		r.DELETE("", apiApprovalTime.Delete)
	}

	// 用户权限处理 todo
	rEmailApprove := v1.Group("/email-approve").Use(authMiddleware.MiddlewareFunc())
	apiEmailApprove := apis.EmailApprove{}
	{
		rEmailApprove.GET("", actions.PermissionAction(), apiEmailApprove.GetPage)
		rEmailApprove.GET("approve-recipient-users", actions.PermissionAction(), apiEmailApprove.GetApproveAndRecipient)
		rEmailApprove.GET("/:id", actions.PermissionAction(), apiEmailApprove.Get)
		rEmailApprove.POST("", apiEmailApprove.Insert)
		rEmailApprove.PUT("/:id", actions.PermissionAction(), apiEmailApprove.Update)
		rEmailApprove.DELETE("", apiEmailApprove.Delete)
		rEmailApprove.GET("workflow", actions.PermissionAction(), apiEmailApprove.Workflow)
	}

	apiUserFootprint := apis.UserFootprint{}
	r = v1.Group("/user-footprint").Use(authMiddleware.MiddlewareFunc())
	{
		r.GET("", actions.PermissionAction(), apiUserFootprint.GetPage)
		r.GET("/:id", actions.PermissionAction(), apiUserFootprint.Get)
		r.POST("", apiUserFootprint.Insert)
		r.PUT("/:id", actions.PermissionAction(), apiUserFootprint.Update)
		r.DELETE("", apiUserFootprint.Delete)
	}
}
