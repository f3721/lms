package database

import (
	"encoding/json"
	"fmt"
	log "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk"
	"github.com/go-admin-team/go-admin-core/sdk/config"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"github.com/monaco-io/request"
	"go-admin/common/global"
	"go-admin/common/middleware/admin_handler"
	"go-admin/common/models"
	appConfig "go-admin/config"
	"time"
)

func dynamicDBInit() {
	go func() {
		for {
			time.Sleep(time.Second * 1)
			dynamicDBRun()
			time.Sleep(time.Second * 30)
		}
	}()
}

func dynamicDBRun() {
	//获取db
	result, err := getTenants()
	if err != nil {
		log.Info(pkg.Red(fmt.Sprintf("动态获取租户数据库失败!!!error:%v", err)))
		return
	}
	//sdk.Runtime.SetDb("ok", *go)
	for _, tenant := range *result {
		if tenant.Status == 1 {
			continue
		}
		global.SetTenant(tenant)
		source, err := tenant.GetDBSourceWithAppName(config.ApplicationConfig.Name)
		if err != nil {
			log.Info(pkg.Red(fmt.Sprintf("租户数据库链接构建失败!error:%v", err)))
			break
		}
		driver := "mysql"
		key := tenant.TenantDBPrefix()
		conf := &config.Database{
			driver,
			source,
			0,
			0,
			0,
			0,
			nil,
		}
		db := sdk.Runtime.GetDbByKey(key)
		if db == nil {
			SetupSimpleDatabase(key, conf, true)
		}
	}
}

// 获取鲁班token
func GetLubanToken() (token string) {
	c := request.Client{
		URL:    appConfig.ExtConfig.LubanHost + "/admin-api/system/oauth2/token?username=sxyz&client_secret=VS7yjxjRVLWYkmcdDAtbQKoD26UAM7s8&password=ehsy2023&client_id=sxyz&grant_type=password",
		Method: "POST",
		Header: models.HS{
			"tenant-id":    "1",
		},
	}
	var resp struct{
		 models.LubanResp
		 Data models.LubanToken
	}
	result := c.Send()
	result.Scan(&resp)
	if resp.Code == 0 && resp.Data.AccessToken != "" {
		token = resp.Data.AccessToken
	}

	return
}

// 获取租户列表
func GetInitOkTenantList(token string) (data []models.SystemTenant) {
	tenants := GetTenantList(token)
	for _, datum := range tenants {
		if datum.DatabaseInitOk == 1 {
			data = append(data, datum)
		}
	}

	return data
}

// 获取租户列表
func GetTenantList(token string) (data []models.SystemTenant) {
	c := request.Client{
		URL:    appConfig.ExtConfig.LubanHost + "/admin-api/system/tenant/list",
		Method: "GET",
		Header: models.HS{
			"Authorization":    "Bearer " + token,
			"tenant-id":    "1",
		},
		Query: models.HS{
			"systemName":    "狮行驿站",
		},
	}

	var resp struct{
		models.LubanResp
		Data []models.SystemTenant
	}
	result := c.Send()
	result.Scan(&resp)
	// 防止鲁班接口失效： 使用缓存存储
	key := "lubanTenantList"
	if resp.Data == nil {
		value, _ := admin_handler.GetCache(key)
		_ = json.Unmarshal([]byte(value), &data)
	} else {
		if diffBytes, err := json.Marshal(resp.Data); err == nil {
			_ = admin_handler.SetCache(key, string(diffBytes), 86400*30)
		}
		data = resp.Data
	}

	return
}

func getTenants() (*[]models.SystemTenant, error) {
	token := GetLubanToken()
	result := GetInitOkTenantList(token)
	return &result, nil
	//httplib.Get("https://www.baidu.com").String()


//	text := `{
//"system_tenant": [
//	{
//		"id" : 1,
//		"name" : "lms",
//		"contact_user_id" : 120,
//		"contact_name" : "LMS",
//		"contact_mobile" : "13775677995",
//		"status" : 0,
//		"domain" : "https:\/\/pms-master-test.ehsy.com",
//		"package_id" : 111,
//		"expire_time" : "2099-01-01 00:00:00",
//		"account_count" : 100,
//		"creator" : "1",
//		"create_time" : "2022-10-11 11:19:49",
//		"updater" : "1",
//		"update_time" : "2022-12-09 11:02:26",
//		"deleted" : 1,
//
//		"db_user" : "root",
//		"db_password" : "root",
//		"db_host" : "localhost",
//		"db_port" : "3306",
//		"db_prefix" : "lms",
//		"db_init_ok" : true
//	}
//]}`
//
//	var result []models.SystemTenant
//	err := json.Unmarshal([]byte(text), &result)
//	if err != nil {
//		return nil, err
//	}
//	return &result, nil
}
