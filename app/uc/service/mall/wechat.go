package mall

import (
	"errors"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"github.com/prometheus/common/log"
	"github.com/silenceper/wechat/v2/miniprogram/auth"
	"go-admin/app/uc/service/mall/dto"
	common "go-admin/common/wechat"
)

type Wechat struct {
	service.Service
}

func (l *Wechat) Login(req *dto.UserLoginRequest) (data *auth.ResCode2Session, err error) {

	wc := common.GetMiniProgram()
	log.Info(req.JsCode)
	userInfo, err := wc.GetAuth().Code2Session(req.JsCode)
	if err != nil {
		log.Info(err.Error())
		return nil, errors.New("参数错误")
	}

	return &userInfo, nil
}

func (l *Wechat) PhoneNumber(req *dto.PhoneNumberRequest) (response *auth.GetPhoneNumberResponse, err error) {
	wc := common.GetMiniProgram()
	log.Info(req.Code)
	phoneNumberInfo, err := wc.GetAuth().GetPhoneNumber(req.Code)
	if err != nil {
		log.Info(err)
		return nil, errors.New("参数错误")
	}

	return phoneNumberInfo, nil
}
