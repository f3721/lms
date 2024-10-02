package dto

type InitTenantReq struct {
	Id       int    `json:"id" comment:"租户id" vd:"@:$>0; msg:'租户id必填'"`
	UserName string `json:"username" comment:"用户名" vd:"@:len($)>0; msg:'用户名必填'"`
	Password string `json:"password" comment:"密码" vd:"@:len($)>0; msg:'密码必填'"`
}

type Login struct {
	Username string `json:"username" vd:"@:len($)>0;msg:'登录名必填'"`
	Password string `json:"password" vd:"@:len($)>0;msg:'密码必填'"`
	Type     string `json:"type" vd:"@:in($,'0','1');msg:'type必填,且范围是[0,1]'"`
	Gen      string `json:"gen"`
}

type LoginMini struct {
	Username string `json:"username" vd:"@:len($)>0;msg:'登录名必填'"`
}
