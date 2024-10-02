package sms

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"go-admin/common/httplib"
	"go-admin/config"
	"strconv"
	"strings"
	"time"
)

type MessageIdType string

const (
	// MessageIdTypeSendVerifyCode 发送验证码，两个参数，第一个参数 验证码 第二个参数 有效时间 分钟
	MessageIdTypeSendVerifyCode MessageIdType = "70032645" //[狮行驿站]验证码: %P%，该验证码%P%分钟内有效。 如非您本人操作，请忽略本短信。
	// MessageIdTypeApplyNewOrder 领用申请单，三个参数：用户 单号 审批链接
	MessageIdTypeApplyNewOrder MessageIdType = "70032641" //[狮行驿站]尊敬的%P%，您收到一条新的领用申请单，单号:%P%，审批请点击 %P%
	// MessageIdTypeRevocationOrder 领用人撤销申请单，三个参数：用户 单号 详情链接
	MessageIdTypeRevocationOrder MessageIdType = "70032643" //[狮行驿站] 尊敬的%P%，您收到的领用申请单已经被领用人撤回，单号: %P%，如有疑问，请联系领用人。
	// MessageIdTypeTurnDownUserOrder 驳回申请单，三个参数：用户 单号 详情链接
	MessageIdTypeTurnDownUserOrder MessageIdType = "70032644" //[狮行驿站]尊敬的%P%，您的领用申请有产品被驳回，单号:%P%，详情请点击%P%
)

type SmsResponse struct {
	Code int    `json:"code"` //code=0，正常
	Msg  string `json:"msg"`  //消息提示
}

// https://www.shlianlu.com/console/document/api_4_4
func SendSMS(userPhones []string, messageIdType MessageIdType, replaceParams []string) (result *SmsResponse, err error) {

	// 读取配置
	config := config.ExtConfig.Sms
	smsUrl := config.SmsUrl
	enterpriseId := config.EnterpriseId
	AppId := config.AppId
	AppKey := config.AppKey

	req := httplib.Post(smsUrl)

	// 设置Header
	req.Header("Accept", "application/json")
	req.Header("Content-Type", "application/json;charset=utf-8")

	// 设置传参
	timeStamp := strconv.Itoa(int(time.Now().Unix()))
	signatureStr := "AppId=" + AppId + "&MchId=" + enterpriseId + "&SignType=" + "MD5" + "&TemplateId=" + string(messageIdType) +
		"&TimeStamp=" + timeStamp + "&Type=3&Version=1.1.0" + "&key=" + AppKey
	signatureMd5 := md5.Sum([]byte(signatureStr))
	signature := strings.ToUpper(hex.EncodeToString(signatureMd5[:]))
	body := map[string]any{
		"MchId":            enterpriseId,
		"AppId":            AppId,
		"Version":          "1.1.0",
		"Type":             "3",
		"PhoneNumberSet":   userPhones,
		"TemplateId":       string(messageIdType),
		"TemplateParamSet": replaceParams,
		"TimeStamp":        timeStamp,
		"SignType":         "MD5",
		"Signature":        signature,
	}
	req.JSONBody(body)

	// 发送请求
	data, err := req.String()
	if err != nil {
		return nil, err
	}

	// 接口返回代码兼容
	newRes := struct {
		TaskId    string `json:"taskId"`
		Message   string `json:"message"`
		Timestamp int    `json:"timestamp"`
		Status    string `json:"status"`
	}{}
	err = json.Unmarshal([]byte(data), &newRes)
	if err != nil {
		return nil, err
	}
	result = &SmsResponse{
		Code: 0,
		Msg:  newRes.Message,
	}
	if newRes.Status != "00" {
		result.Code = 300
		fmt.Println("[SMS]", newRes)
	}

	return result, nil
}
