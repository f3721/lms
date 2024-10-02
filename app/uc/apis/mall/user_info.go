package mall

import (
	"fmt"
	"go-admin/common/middleware/mall_handler"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/captcha"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/uc/models"
	service "go-admin/app/uc/service/mall"
	"go-admin/app/uc/service/mall/dto"
	"go-admin/common/actions"
)

type UserInfo struct {
	api.Api
}

// GetPage 获取用户中心列表
// @Summary 获取用户中心列表
// @Description 获取用户中心列表
// @Tags 商城-个人中心
// @Param userEmail query string false "用户邮箱"
// @Param loginName query string false "用户登录名"
// @Param userPhone query string false "手机号码"
// @Param userName query string false "用户名称"
// @Param userStatus query int false "用户状态（1可用，0不可用）"
// @Param companyId query int false "公司ID"
// @Param isAdminShow query string false "用户状态（1是，0否）"
// @Param canLogin query int false "是否可登陆"
// @Param companyDepartmentId query int false "用户所属部门"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.UserInfo}} "{"code": 200, "data": [...]}"
// @Router /api/v1/uc/mall/user-info [get]
// @Security Bearer
func (e UserInfo) GetPage(c *gin.Context) {
	req := dto.UserInfoGetPageReq{}
	s := service.UserInfo{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	p := actions.GetPermissionFromContext(c)
	list := make([]models.UserInfo, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取用户中心失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// GetUserInfo 获取当前用户信息
// @Summary 获取当前用户信息
// @Description 获取当前用户信息
// @Tags 商城-个人中心
// @Success 200 {object} response.Response{data=models.UserInfo} "{"code": 200, "data": [...]}"
// @Router /api/v1/uc/mall/user/info [get]
// @Security Bearer
func (e UserInfo) GetUserInfo(c *gin.Context) {
	req := dto.UserInfoGetReq{}
	s := service.UserInfo{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	req.Id = user.GetUserId(c)
	var object models.UserInfo

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取用户中心失败，\r\n失败信息 %s", err.Error()))
		return
	}

	object.UserConfig = mall_handler.GetUserConfig(c)
	e.OK(object, "查询成功")
}

// Get 获取用户中心
// @Summary 获取用户中心
// @Description 获取用户中心
// @Tags 商城-个人中心
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.UserInfo} "{"code": 200, "data": [...]}"
// @Router /api/v1/uc/mall/user-info/{id} [get]
// @Security Bearer
func (e UserInfo) Get(c *gin.Context) {
	req := dto.UserInfoGetReq{}
	s := service.UserInfo{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	var object models.UserInfo

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取用户中心失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建用户中心
// @Summary 创建用户中心
// @Description 创建用户中心
// @Tags 商城-个人中心
// @Accept application/json
// @Product application/json
// @Param data body dto.UserInfoInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/uc/mall/user-info [post]
// @Security Bearer
func (e UserInfo) Insert(c *gin.Context) {
	req := dto.UserInfoInsertReq{}
	s := service.UserInfo{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	// 设置创建人
	req.SetCreateBy(user.GetUserId(c))
	req.SetCreateByName(user.GetUserName(c))
	err = s.Insert(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建用户中心失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改用户中心
// @Summary 修改用户中心
// @Description 修改用户中心
// @Tags 商城-个人中心
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.UserInfoUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/uc/mall/user-info/{id} [put]
// @Security Bearer
func (e UserInfo) Update(c *gin.Context) {
	req := dto.UserInfoUpdateReq{}
	s := service.UserInfo{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	req.SetUpdateBy(user.GetUserId(c))
	p := actions.GetPermissionFromContext(c)

	err = s.Update(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改用户中心失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除用户中心
// @Summary 删除用户中心
// @Description 删除用户中心
// @Tags 商城-个人中心
// @Param data body dto.UserInfoDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/uc/mall/user-info [delete]
// @Security Bearer
func (e UserInfo) Delete(c *gin.Context) {
	s := service.UserInfo{}
	req := dto.UserInfoDeleteReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	// req.SetUpdateBy(user.GetUserId(c))
	p := actions.GetPermissionFromContext(c)

	err = s.Remove(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("删除用户中心失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}

// ChangeBasicInfo 用户基础资料修改
// @Summary 用户基础资料修改
// @Description 用户基础资料修改
// @Tags 商城-个人中心
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.UserInfoChangeBasicInfoReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/uc/mall/user/change-basic-info [put]
// @Security Bearer
func (e UserInfo) ChangeBasicInfo(c *gin.Context) {
	req := dto.UserInfoChangeBasicInfoReq{}
	s := service.UserInfo{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, "11")
		return
	}

	req.UserId = user.GetUserId(c)
	err = s.ChangeBasicInfo(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改资料失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(nil, "修改成功")
}

// ChangePassword 用户修改密码
// @Summary 用户修改密码
// @Description 用户修改密码
// @Tags 商城-个人中心
// @Accept application/json
// @Product application/json
// @Param data body dto.UserInfoChangePasswordReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/uc/mall/user/change-password [put]
// @Security Bearer
func (e UserInfo) ChangePassword(c *gin.Context) {
	req := dto.UserInfoChangePasswordReq{}
	s := service.UserInfo{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	//req.SetUpdateBy(user.GetUserId(c))
	//p := actions.GetPermissionFromContext(c)
	req.UserId = user.GetUserId(c)
	err = s.ChangePassword(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改密码失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(nil, "修改成功")
}

// ChangeUserEmail 修改邮箱
// @Summary 修改邮箱
// @Description 修改邮箱
// @Tags 商城-个人中心
// @Accept application/json
// @Product application/json
// @Param data body dto.UserInfoChangeUserEmailReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/uc/mall/user/change-user-email [put]
// @Security Bearer
func (e UserInfo) ChangeUserEmail(c *gin.Context) {
	req := dto.UserInfoChangeUserEmailReq{}
	s := service.UserInfo{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	//req.SetUpdateBy(user.GetUserId(c))
	//p := actions.GetPermissionFromContext(c)
	req.UserId = user.GetUserId(c)
	err = s.ChangeUserEmail(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改邮箱失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(nil, "修改成功")
}

// ChangePhone 修改手机号
// @Summary 修改手机号
// @Description 修改手机号
// @Tags 商城-个人中心
// @Accept application/json
// @Product application/json
// @Param data body dto.UserInfoChangePhoneReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/uc/mall/user/change-phone [put]
// @Security Bearer
func (e UserInfo) ChangePhone(c *gin.Context) {
	req := dto.UserInfoChangePhoneReq{}
	s := service.UserInfo{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	//req.SetUpdateBy(user.GetUserId(c))
	//p := actions.GetPermissionFromContext(c)
	req.UserId = user.GetUserId(c)
	err = s.ChangePhone(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改手机号失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(nil, "修改成功")
}

// SendEmailForCheckEmail 修改邮箱发送验证邮件
// @Summary 修改邮箱发送验证邮件
// @Description 修改邮箱发送验证邮件
// @Tags 商城-个人中心
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.SendEmailForCheckEmailReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/uc/mall/user/send-email-for-check-email [post]
// @Security Bearer
func (e UserInfo) SendEmailForCheckEmail(c *gin.Context) {
	req := dto.SendEmailForCheckEmailReq{}
	s := service.UserInfo{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	//req.SetUpdateBy(user.GetUserId(c))
	//p := actions.GetPermissionFromContext(c)
	req.UserId = user.GetUserId(c)
	err = s.SendEmailForCheckEmail(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("发送邮件失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(nil, "修改成功")
}

// SendSms 发送手机验证码(不校验手机号是否存在)
// @Summary 发送手机验证码(不校验手机号是否存在)
// @Description 发送手机验证码(不校验手机号是否存在)
// @Tags 商城-个人中心
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.VerifyCodeReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/uc/mall/user/send-sms [post]
// @Security Bearer
func (e UserInfo) SendSms(c *gin.Context) {
	req := dto.VerifyCodeReq{}
	s := service.UserInfo{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	//req.SetUpdateBy(user.GetUserId(c))
	//p := actions.GetPermissionFromContext(c)
	_, err = s.VerifyCode(&req, false)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改用户中心失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(nil, "验证码发送成功")
}

// SendEmailForChangeEmail 修改邮箱发送修改邮件
// @Summary 修改邮箱发送修改邮件
// @Description 修改邮箱发送修改邮件
// @Tags 商城-个人中心
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.SendEmailForChangeEmailReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/uc/mall/user/send-email-for-change-email [post]
// @Security Bearer
func (e UserInfo) SendEmailForChangeEmail(c *gin.Context) {
	req := dto.SendEmailForChangeEmailReq{}
	s := service.UserInfo{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	//req.SetUpdateBy(user.GetUserId(c))
	//p := actions.GetPermissionFromContext(c)
	req.UserId = user.GetUserId(c)
	err = s.SendEmailForChangeEmail(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("发送邮件失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(nil, "修改成功")
}

// CheckUserTokenAvailable 验证修改邮箱的Token
// @Summary 验证修改邮箱的Token
// @Description 验证修改邮箱的Token
// @Tags 商城-个人中心
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.CheckUserTokenAvailableReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/uc/mall/user/check-send-email-token [post]
// @Security Bearer
func (e UserInfo) CheckUserTokenAvailable(c *gin.Context) {
	req := dto.CheckUserTokenAvailableReq{}
	s := service.UserInfo{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	//req.SetUpdateBy(user.GetUserId(c))
	//p := actions.GetPermissionFromContext(c)
	req.UserId = user.GetUserId(c)
	tokenAvailable, _ := s.CheckUserEmailTokenAvailable(&req)

	e.OK(tokenAvailable, "修改成功")
}

// Captcha 找回密码验证码
// @Summary 找回密码验证码
// @Description 找回密码验证码
// @Tags 找回密码
// @Success 200 {object} response.Response	"{"code": 200, "data": {"data":data,"id":id}, "msg": "success"}"
// @Router /api/v1/uc/mall/user/captcha [get]
// @Security Bearer
func (e UserInfo) Captcha(c *gin.Context) {
	err := e.MakeContext(c).Errors
	if err != nil {
		e.Error(500, err, "服务初始化失败！")
		return
	}
	id, b64s, err := captcha.DriverDigitFunc()
	if err != nil {
		e.Logger.Errorf("DriverDigitFunc error, %s", err.Error())
		e.Error(500, err, "验证码获取失败")
		return
	}
	res := map[string]any{"id": id, "data": b64s}
	e.OK(res, "获取图片验证码成功")
}

// checkPhone 检查手机号是否注册
// @Summary 检查手机号是否注册
// @Description 检查手机号是否注册
// @Tags 找回密码
// @Param data body dto.VerifyCodeReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "data": true, "msg": "success"}"
// @Router /api/v1/uc/mall/user/check-phone [post]
// @Security Bearer
func (e UserInfo) CheckPhone(c *gin.Context) {
	s := service.UserInfo{}
	req := dto.VerifyCodeReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	err = s.CheckPhone(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("检查手机号是否注册失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(true, "检查手机号是否注册成功")
}

// VerifyCode 获取手机验证码
// @Summary 获取手机验证码
// @Description 获取手机验证码
// @Tags 找回密码
// @Param data body dto.VerifyCodeReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "data": "data", "msg": "success"}"
// @Router /api/v1/uc/mall/user/verify-code [post]
// @Security Bearer
func (e UserInfo) VerifyCode(c *gin.Context) {
	s := service.UserInfo{}
	req := dto.VerifyCodeReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	_, err = s.VerifyCode(&req, true)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取手机验证码失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(nil, "获取手机验证码成功")
}

// ComitVerifyCode 提交手机验证码
// @Summary 提交手机验证码
// @Description 提交手机验证码
// @Tags 找回密码
// @Param data body dto.ComitVerifyCodeReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "data": "data", "msg": "success"}"
// @Router /api/v1/uc/mall/user/comit-verify-code [post]
// @Security Bearer
func (e UserInfo) ComitVerifyCode(c *gin.Context) {
	s := service.UserInfo{}
	req := dto.ComitVerifyCodeReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	token, err := s.ComitVerifyCode(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("提交手机验证码失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(token, "提交手机验证码成功")
}

// SendEmail 发送至邮箱
// @Summary 修改密码
// @Description 修改密码
// @Tags 找回密码
// @Param data body dto.SendEmailReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "data": "data", "msg": "success"}"
// @Router /api/v1/uc/mall/user/send-email [post]
// @Security Bearer
func (e UserInfo) SendEmail(c *gin.Context) {
	s := service.UserInfo{}
	req := dto.SendEmailReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	err = s.SendEmail(c, &req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改密码失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK("", "已发送至邮箱，请查收！")
}

// ChangePwd 修改密码
// @Summary 修改密码
// @Description 修改密码
// @Tags 找回密码
// @Param data body dto.ChangePwdReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "data": "data", "msg": "success"}"
// @Router /api/v1/uc/mall/user/change-pwd [post]
// @Security Bearer
func (e UserInfo) ChangePwd(c *gin.Context) {
	s := service.UserInfo{}
	req := dto.ChangePwdReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	err = s.ChangePwd(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改密码失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(true, "修改密码成功")
}

// proxy-login 代登录[前端部分]
// @Summary 代登录[前端部分]
// @Description 代登录[前端部分]
// @Tags 商城-个人中心
// @Param data body dto.ChangePwdReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "data": "data", "msg": "success"}"
// @Router /api/v1/uc/mall/user/proxy-login [post]
// @Security Bearer
func (e UserInfo) ProxyLogin(c *gin.Context) {
	s := service.UserInfo{}
	req := dto.ProxyLoginReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	res, err := s.ProxyLogin(&req, c)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("代登录失败，\r\n失败信息 %s", err.Error()))
		return
	}

	// 设置cookie
	e.OK(res, "代登录成功")
}

// mini-login 小程序登录
// @Summary 小程序登录
// @Description 小程序登录
// @Tags 商城-个人中心
// @Param data body dto.ChangePwdReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "data": "data", "msg": "success"}"
// @Router /api/v1/uc/mall/user/proxy-login [post]
// @Security Bearer
func (e UserInfo) MiniLogin(c *gin.Context) {
	s := service.UserInfo{}
	req := dto.MiniLoginReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	res, err := s.MiniLogin(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("小程序登录失败，\r\n失败信息 %s", err.Error()))
		return
	}

	// 设置cookie
	e.OK(res, "小程序登录成功")
}
