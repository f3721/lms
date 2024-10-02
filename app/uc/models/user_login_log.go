package models

import (
	"encoding/json"
	"errors"
	log "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk"
	"github.com/go-admin-team/go-admin-core/storage"
	"go-admin/common/models"
	"time"
)

type UserLoginLog struct {
	models.Model
	// 主键编码
	UserId        int       `json:"userId" gorm:"type:int unsigned;comment:用户ID"`       // 用户ID
	Username      string    `json:"username" gorm:"type:varchar(128);comment:用户名"`      // 用户名
	Status        string    `json:"status" gorm:"type:varchar(4);comment:状态"`           // 状态
	Ipaddr        string    `json:"ipaddr" gorm:"type:varchar(255);comment:ip地址"`       // ip地址
	LoginLocation string    `json:"loginLocation" gorm:"type:varchar(255);comment:归属地"` // 归属地
	Browser       string    `json:"browser" gorm:"type:varchar(255);comment:浏览器"`       // 浏览器
	Os            string    `json:"os" gorm:"type:varchar(255);comment:系统"`             // 系统
	Platform      string    `json:"platform" gorm:"type:varchar(255);comment:固件"`       // 固件
	LoginTime     time.Time `json:"loginTime" gorm:"type:timestamp;comment:登录时间"`       // 登录时间
	Remark        string    `json:"remark" gorm:"type:varchar(255);comment:备注"`         // 备注
	Msg           string    `json:"msg" gorm:"type:varchar(255);comment:信息"`            // 信息
	models.ModelTime
	models.ControlBy
}

func (UserLoginLog) TableName() string {
	return "user_login_log"
}

func (e *UserLoginLog) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *UserLoginLog) GetId() interface{} {
	return e.Id
}

// MallSaveLoginLog 从队列中获取登录日志
func MallSaveLoginLog(message storage.Messager) (err error) {
	values := message.GetValues()
	//准备db
	db := sdk.Runtime.GetDbByKey(values["db_prefix"].(string)) //默认base，有租户取租户ID+前缀
	if db == nil {
		err = errors.New("db not exist")
		log.Errorf("host[%s]'s %s", message.GetPrefix(), err.Error())
		return err
	}
	var rb []byte
	rb, err = json.Marshal(values)
	if err != nil {
		log.Errorf("json Marshal error, %s", err.Error())
		return err
	}
	var l UserLoginLog
	err = json.Unmarshal(rb, &l)
	if err != nil {
		log.Errorf("json Unmarshal error, %s", err.Error())
		return err
	}
	err = db.Create(&l).Error
	if err != nil {
		log.Errorf("db create error, %s", err.Error())
		return err
	}
	return nil
}
