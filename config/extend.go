package config

var ExtConfig Extend

// Extend 扩展配置
//
//	extend:
//	  demo:
//	    name: demo-name
//
// 使用方法： config.ExtConfig......即可！！
type Extend struct {
	AMap              AMap // 这里配置对应配置文件的结构即可
	Module            string
	LmsHost           string
	MallHost          string
	ApiHost           string
	StaticPath        string
	NoNetwork         bool
	Wc                Wc
	LubanHost         string
	MiniProgramConfig MiniProgramConfig
	Sms               Sms
}

type AMap struct {
	Key string
}

type Wc struct {
	PrintSkuPage string
}

// MiniProgramConfig 小程序相关配置
type MiniProgramConfig struct {
	AppID     string `yaml:"appID"`
	AppSecret string `yaml:"appSecret"`
}

// 短信服务商配置
type Sms struct {
	SmsUrl       string `yaml:"smsUrl"`
	EnterpriseId string `yaml:"enterpriseId"`
	AppId        string `yaml:"appId"`
	AppKey       string `yaml:"appKey"`
}
