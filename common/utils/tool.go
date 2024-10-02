package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	goStore "go-admin/common/client/go_store"
	"math"
	"time"

	"gorm.io/gorm"
)

func GetTFromMap[T any](key string, mapping map[string]T) (r T) {
	res, ok := mapping[key]
	if !ok {
		return r
	}
	return res
}

//获取当天时间

func GetTodayTime() (start, end string) {
	todayDate := time.Now().Format("2006-01-02")
	start = todayDate + " 00:00:00"
	end = todayDate + " 23:59:59"
	return
}

func EncryptTenantId(tid int) string {
	salt := "23520lms1452"
	data := []byte(string(tid) + salt)
	md5New := md5.New()
	md5New.Write(data)
	return hex.EncodeToString(md5New.Sum(nil))
}

func GetShortLink(tx *gorm.DB, url string) string {
	fmt.Println("GetShortLink[1]:", url)
	result := goStore.ApiByDbContext(tx).GetShortLink(url)
	resultInfo := &struct {
		Code    int               `json:"code"`
		Message string            `json:"message"`
		Data    map[string]string `json:"data"`
	}{}
	result.Scan(resultInfo)
	return resultInfo.Data["url"]
}

func AbsToInt(val int) int {
	return int(math.Abs(float64(val)))
}
