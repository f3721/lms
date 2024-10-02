package mall_handler

import (
	"errors"
	"go-admin/app/admin/models"
	"go-admin/common"
	"net/http"

	"go-admin/common/global"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"github.com/mssola/user_agent"
)

func PayloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(map[string]interface{}); ok {
		u, _ := v["user"].(UserInfo)
		r, _ := v["role"].(RoleInfo)
		return jwt.MapClaims{
			jwt.IdentityKey:  u.Id,
			jwt.RoleIdKey:    r.Id,
			jwt.RoleKey:      r.RoleKey,
			jwt.NiceKey:      u.UserName,
			jwt.DataScopeKey: r.DataScope,
			jwt.RoleNameKey:  r.RoleName,
		}
	}
	return jwt.MapClaims{}
}

func IdentityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	return map[string]interface{}{
		"IdentityKey": claims["identity"],
		"UserName":    claims["nice"],
		"RoleKey":     claims["rolekey"],
		"UserId":      claims["identity"],
		"RoleIds":     claims["roleid"],
		"DataScope":   claims["datascope"],
	}
}

// Authenticator 获取token
// @Summary 商城登陆
// @Description 获取token
// @Description LoginHandler can be used by clients to get a jwt token.
// @Description Payload needs to be json in the form of {"username": "USERNAME", "password": "PASSWORD"}.
// @Description Reply will be of the form {"token": "TOKEN"}.
// @Description dev mode：It should be noted that all fields cannot be empty, and a value of 0 can be passed in addition to the account password
// @Description 注意：开发模式：需要注意全部字段不能为空，账号密码外可以传入0值
// @Tags Mall登陆
// @Accept  application/json
// @Product application/json
// @Param account body Login  true "account"
// @Success 200 {string} string "{"code": 200, "expire": "2019-08-07T12:45:48+08:00", "token": ".eyJleHAiOjE1NjUxNTMxNDgsImlkIjoiYWRtaW4iLCJvcmlnX2lhdCI6MTU2NTE0OTU0OH0.-zvzHvbg0A" }"
// @Router /api/v1/uc/mall/login [post]
func Authenticator(c *gin.Context) (interface{}, error) {
	log := api.GetRequestLogger(c)
	db, err := pkg.GetOrm(c)
	if err != nil {
		log.Errorf("get db error, %s", err.Error())
		response.Error(c, 500, err, "数据库连接获取失败")
		return nil, errors.New("系统繁忙，请稍后再试")
	}

	var loginVals Login
	var status = "2"
	var msg = "登录成功"
	var username = ""
	var userId int
	defer func() {
		LoginLogToDB(c, status, msg, username, userId)
	}()

	if err = c.ShouldBind(&loginVals); err != nil {
		username = loginVals.Username
		msg = "数据解析失败"
		status = "1"

		return nil, errors.New("缺少用户名或密码")
	}
	username = loginVals.Username

	// 登录查询
	user, role, e := loginVals.GetUser(db)
	if e != nil {
		msg = "登录失败"
		status = "1"
		return nil, e
	}
	userId = user.Id

	// 默认仓库选中
	err = InitWarehouse(c, userId)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"user": user, "role": role}, nil
}

// LoginLogToDB Write log to database
func LoginLogToDB(c *gin.Context, status string, msg string, username string, userId int) {
	//if !config.LoggerConfig.EnabledDB {
	//	return
	//}
	// if config.ApplicationConfig.Mode == "dev" {
	// 	return
	// }
	log := api.GetRequestLogger(c)
	l := make(map[string]interface{})

	dbPrefix := global.GetTenantDBPrefixWithContext(c)
	ua := user_agent.New(c.Request.UserAgent())
	l["db_prefix"] = dbPrefix
	l["userId"] = userId
	l["ipaddr"] = common.GetClientIP(c)
	l["loginLocation"] = "" // pkg.GetLocation(common.GetClientIP(c),gaConfig.ExtConfig.AMap.Key)
	l["loginTime"] = pkg.GetCurrentTime()
	l["status"] = status
	l["remark"] = c.Request.UserAgent()
	browserName, browserVersion := ua.Browser()
	l["browser"] = browserName + " " + browserVersion
	l["os"] = ua.OS()
	l["platform"] = ua.Platform()
	l["username"] = username
	l["msg"] = msg

	q := sdk.Runtime.GetMemoryQueue(c.Request.Host)
	message, err := sdk.Runtime.GetStreamMessage("", global.MallLoginLog, l)
	if err != nil {
		log.Errorf("GetStreamMessage error, %s", err.Error())
		//日志报错错误，不中断请求
	} else {
		err = q.Append(message)
		if err != nil {
			log.Errorf("Append message error, %s", err.Error())
		}
	}
}

// LogOut
// @Summary 退出登录
// @Description 获取token
// @Tags Mall登陆
// LoginHandler can be used by clients to get a jwt token.
// Reply will be of the form {"token": "TOKEN"}.
// @Accept  application/json
// @Product application/json
// @Success 200 {string} string "{"code": 200, "msg": "成功退出系统" }"
// @Router /api/v1/uc/mall/logout [post]
// @Security Bearer
func LogOut(c *gin.Context) {
	LoginLogToDB(c, "2", "退出成功", user.GetUserName(c), user.GetUserId(c))
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "退出成功",
	})

}

func Authorizator(data interface{}, c *gin.Context) bool {

	if v, ok := data.(map[string]interface{}); ok {
		u, _ := v["user"].(models.SysUser)
		r, _ := v["role"].(models.SysRole)
		c.Set("role", r.RoleName)
		c.Set("roleIds", r.RoleId)
		c.Set("userId", u.UserId)
		c.Set("userName", u.Username)
		c.Set("dataScope", r.DataScope)
		return true
	}
	return false
}

func Unauthorized(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  message,
	})
}
