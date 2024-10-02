package mall

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	service "go-admin/app/uc/service/mall"
	"go-admin/app/uc/service/mall/dto"
)

type Wechat struct {
	api.Api
}

// Login 微信小程序登录
// @Summary 微信小程序登录
// @Description 微信小程序登录
// @Tags 商城-微信登录api
// @Param data body dto.UserLoginRequest true "data"
// @Success 200 {object} response.Response{data=dto.UserLoginResponse} "{"code": 200, "data": [...]}"
// @Router /api/v1/uc/mall/wechat/login [post]
// @Security Bearer
func (e Wechat) Login(c *gin.Context) {
	req := dto.UserLoginRequest{}
	s := service.Wechat{}
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

	userInfo, err := s.Login(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取用户中心失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(userInfo, "查询成功")
}

// PhoneNumber 微信小程序获取手机号
// @Summary 微信小程序获取手机号
// @Description 获取手机号
// @Tags 商城-微信登录api
// @Param data body dto.PhoneNumberRequest true "data"
// @Success 200 {object} response.Response{data=dto.PhoneNumberResponse} "{"code": 200, "data": [...]}"
// @Router /api/v1/uc/mall/wechat/phone-number [post]
// @Security Bearer
func (e Wechat) PhoneNumber(c *gin.Context) {
	req := dto.PhoneNumberRequest{}
	s := service.Wechat{}
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

	res, err := s.PhoneNumber(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取用户中心失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(res, "查询成功")
}
