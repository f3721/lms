package router

import (
	apis "go-admin/app/uc/apis/admin"
	"go-admin/common/actions"
	"go-admin/common/middleware"

	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerSysAdminApiRouter)
	routerNoCheckRole = append(routerNoCheckRole, registerNoCheckRoleApiRouter)
}

// registerSysApiRouter
func registerSysAdminApiRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	// 部门预算
	apiDeptBudget := apis.DepartmentBudget{}
	r := v1.Group("/department-budget").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", actions.PermissionAction(), apiDeptBudget.GetPage)
		r.GET("/:id", actions.PermissionAction(), apiDeptBudget.Get)
	}

	// 用户预算
	apiUserBudget := apis.UserBudget{}
	r = v1.Group("/user-budget").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", actions.PermissionAction(), apiUserBudget.GetPage)
		r.GET("/:id", actions.PermissionAction(), apiUserBudget.Get)
	}

	// 公司
	apiCompanyInfo := apis.CompanyInfo{}
	r = v1.Group("/admin/company-info").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", actions.PermissionAction(), apiCompanyInfo.GetPage)
		r.GET("/:id", actions.PermissionAction(), apiCompanyInfo.Get)
		r.POST("", apiCompanyInfo.Insert)
		r.PUT("/:id", actions.PermissionAction(), apiCompanyInfo.Update)
		//r.DELETE("", api.Delete)

		r.GET("select-list", actions.PermissionAction(), apiCompanyInfo.GetSelectList)
		r.GET("is-available/:id", actions.PermissionAction(), apiCompanyInfo.IsAvailable)
		r.GET("parameters", actions.PermissionAction(), apiCompanyInfo.Parameters)
	}

	// 部门
	apiCompanyDepartment := apis.CompanyDepartment{}
	r = v1.Group("/admin/company-department").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", actions.PermissionAction(), apiCompanyDepartment.GetPage)
		r.GET("/:id", actions.PermissionAction(), apiCompanyDepartment.Get)
		r.POST("", apiCompanyDepartment.Insert)
		r.PUT("/:id", actions.PermissionAction(), apiCompanyDepartment.Update)
		r.DELETE("", apiCompanyDepartment.Delete)

		r.GET("/select-list", actions.PermissionAction(), apiCompanyDepartment.GetSelectList)
		r.PUT("/update-budget/:id", actions.PermissionAction(), apiCompanyDepartment.UpdateBudget)
		r.GET("/is-cane-delete/:id", actions.PermissionAction(), apiCompanyDepartment.IsCaneDelete)

		r.POST("/import", actions.PermissionAction(), apiCompanyDepartment.Import)
		r.GET("/export", actions.PermissionAction(), apiCompanyDepartment.Export)
	}

	// 用户
	apiUserInfo := apis.UserInfo{}
	r = v1.Group("/admin/user-info").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", actions.PermissionAction(), apiUserInfo.GetPage)
		r.GET("/:id", actions.PermissionAction(), apiUserInfo.Get)
		r.POST("", apiUserInfo.Insert)
		r.PUT("/:id", actions.PermissionAction(), apiUserInfo.Update)

		r.GET("select-list", actions.PermissionAction(), apiUserInfo.GetSelectList)
		// 修改密码
		r.PUT("update-password", actions.PermissionAction(), apiUserInfo.UpdatePassword)
		// 修改密码
		r.PUT("update-userphone", actions.PermissionAction(), apiUserInfo.UpdateUserPhone)
		r.GET("proxy-login", actions.PermissionAction(), apiUserInfo.ProxyLogin)
		// 导入
		r.POST("import", actions.PermissionAction(), apiUserInfo.Import)
	}

	// 用户地址
	apiAddress := apis.Address{}
	r = v1.Group("/admin/user/address").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", actions.PermissionAction(), apiAddress.GetPage)
		r.GET("/:id", actions.PermissionAction(), apiAddress.Get)
		r.POST("", apiAddress.Insert)
		r.PUT("/:id", actions.PermissionAction(), apiAddress.Update)
		r.DELETE("", apiAddress.Delete)
	}

	// 用户收藏
	apiUserCollect := apis.UserCollect{}
	r = v1.Group("/admin/user-collect").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", actions.PermissionAction(), apiUserCollect.GetPage)
		//r.GET("/:id", actions.PermissionAction(), api.Get)
		//r.POST("", api.Insert)
		//r.PUT("/:id", actions.PermissionAction(), api.Update)
		//r.DELETE("", api.Delete)
	}

	// 用户登录日志
	apiUserLoginLog := apis.UserLoginLog{}
	r = v1.Group("/admin/user/login-log").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", actions.PermissionAction(), apiUserLoginLog.GetPage)
	}

	// 用户密码修改日志
	apiUserPasswordChangeLog := apis.UserPasswordChangeLog{}
	r = v1.Group("/admin/user-password-change-log").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", actions.PermissionAction(), apiUserPasswordChangeLog.GetPage)
	}

	// 角色信息
	apiRoleInfo := apis.RoleInfo{}
	r = v1.Group("/admin/role").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", actions.PermissionAction(), apiRoleInfo.GetPage)
		r.GET("/:id", actions.PermissionAction(), apiRoleInfo.Get)
		r.POST("", apiRoleInfo.Insert)
		r.PUT("/:id", actions.PermissionAction(), apiRoleInfo.Update)
		r.DELETE("", apiRoleInfo.Delete)
	}

	//操作日志
	apiOperateLogs := apis.OperateLogs{}
	rOperateLogs := v1.Group("/admin/operate-logs").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		rOperateLogs.GET("", apiOperateLogs.GetPage)
		rOperateLogs.GET("/:id", apiOperateLogs.Get)
	}

	// 客户分类
	apiClassification := apis.Classification{}
	r = v1.Group("/admin/classification").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", actions.PermissionAction(), apiClassification.GetPage)
		r.GET("/:id", actions.PermissionAction(), apiClassification.Get)
		r.POST("", apiClassification.Insert)
		r.PUT("/:id", actions.PermissionAction(), apiClassification.Update)
		//r.DELETE("", apiClassification.Delete)
	}

	// 公司开关
	apiCompanyIndividualitySwitch := apis.CompanyIndividualitySwitch{}
	r = v1.Group("/admin/company-individuality-switch").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", actions.PermissionAction(), apiCompanyIndividualitySwitch.GetPage)
		r.GET("/:id", actions.PermissionAction(), apiCompanyIndividualitySwitch.Get)
		r.POST("", apiCompanyIndividualitySwitch.Insert)
		r.PUT("/:id", actions.PermissionAction(), apiCompanyIndividualitySwitch.Update)
		//r.DELETE("", apiCompanyIndividualitySwitch.Delete)
	}

	// SKU对应客户分类
	apiSkuClassification := apis.SkuClassification{}
	r = v1.Group("/uc/admin/sku-classification").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", actions.PermissionAction(), apiSkuClassification.GetPage)
		r.GET("/:id", actions.PermissionAction(), apiSkuClassification.Get)
		r.POST("", apiSkuClassification.Insert)
		r.PUT("/:id", actions.PermissionAction(), apiSkuClassification.Update)
		//r.DELETE("", apiSkuClassification.Delete)
	}

	// 客户分类对应支付账户
	apiUserPayAccount := apis.UserPayAccount{}
	r = v1.Group("/uc/admin/user-pay-account").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", actions.PermissionAction(), apiUserPayAccount.GetPage)
		r.GET("/:id", actions.PermissionAction(), apiUserPayAccount.Get)
		r.POST("", apiUserPayAccount.Insert)
		r.PUT("/:id", actions.PermissionAction(), apiUserPayAccount.Update)
		//r.DELETE("", apiUserPayAccount.Delete)
	}
}

func registerNoCheckRoleApiRouter(v1 *gin.RouterGroup) {
	// 部门预算
	apiDeptBudget := apis.DepartmentBudget{}
	r := v1.Group("/department-budget")
	{
		r.GET("/init", apiDeptBudget.Init)
	}
}
