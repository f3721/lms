package mall_handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-admin/common/global"
	"strconv"

	"github.com/go-redis/redis/v7"

	"github.com/gin-gonic/gin"
	log "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	"gorm.io/gorm"
)

func getUserConfigCacheKey(userId string) string {
	return fmt.Sprintf("UserConfigWithID:%s", userId)
}

type UserConfig struct {
	UserID                int    `json:"userId"`
	SelectedWarehouseId   int    `json:"selectedWarehouseId"`
	SelectedWarehouseCode string `json:"selectedWarehouseCode"`
	SelectedWarehouseName string `json:"selectedWarehouseName"`
}

// 获取用户配置信息|仓库
func GetUserConfig(c *gin.Context) *UserConfig {
	userId := user.GetUserId(c)
	tenantId := c.GetHeader("tenant-id")
	key := getUserConfigCacheKey(tenantId + "_" + strconv.Itoa(userId))
	data, err := GetCache(key)
	if err != nil {
		return &UserConfig{UserID: userId}
	}
	var config UserConfig
	err = json.Unmarshal([]byte(data), &config)
	if err != nil {
		return &UserConfig{UserID: userId}
	}
	return &config
}

// 设置用户信息
func SetUserConfig(c *gin.Context, config *UserConfig) error {
	userId := user.GetUserId(c)
	tenantId := c.GetHeader("tenant-id")
	key := getUserConfigCacheKey(tenantId + "_" + strconv.Itoa(userId))
	data, _ := json.Marshal(config)
	err := SetCache(key, data, 3600*24*30)
	if err != nil {
		return err
	}
	return nil
}

// 获取用户当前仓库code

func GetUserCurrentWarehouseCode(c *gin.Context) string {
	//return "WH0037" //lily1仓库
	return GetUserConfig(c).SelectedWarehouseCode
}

// 登录后默认仓库
func InitWarehouse(c *gin.Context, userId int) error {
	// 查询用户是否已经设置仓库
	tenantId := c.GetHeader("tenant-id")
	key := getUserConfigCacheKey(tenantId + "_" + strconv.Itoa(userId))
	data, err := GetCache(key)
	if err != nil && err != redis.Nil {
		return err
	}
	if data != "" {
		// 初始化失败,继续重设
		old := UserConfig{}
		_ = json.Unmarshal([]byte(data), &old)
		if old.SelectedWarehouseCode != "" {
			return nil
		}
	}

	// 初始化DB
	db, err := pkg.GetOrm(c)
	if err != nil {
		log.Errorf("get db error, %s", err.Error())
		return err
	}

	// 查询公司ID
	ucDbName := global.GetTenantUcDBNameWithContext(c)
	var companyId int
	db.Raw("SELECT company_id FROM "+ucDbName+".user_info WHERE id = ?", userId).Scan(&companyId)
	if companyId == 0 {
		return errors.New("未查询到对应用户的公司")
	}

	// 公司第一个仓库
	var warehouse Warehouse
	wcDbName := global.GetTenantWcDBNameWithContext(c)
	err = db.Raw("SELECT * FROM "+wcDbName+".warehouse WHERE company_id = ? and is_virtual = ?", companyId, 0).Scan(&warehouse).Error
	// err = db.Where("companyId = ?", companyId).Where("isVirtual = ?", 0).Order("id ASC").First(&warehouse).Error
	if err == gorm.ErrRecordNotFound {
		return errors.New("未查询到该公司对应实体仓")
	}
	if err != nil {
		return err
	}

	// 设置默认
	config := UserConfig{
		UserID:                userId,
		SelectedWarehouseId:   warehouse.Id,
		SelectedWarehouseCode: warehouse.WarehouseCode,
		SelectedWarehouseName: warehouse.WarehouseName,
	}
	con, _ := json.Marshal(config)
	err = SetCache(key, con, 3600*24*30)
	if err != nil {
		return err
	}

	return nil
}
