package models

import (
	"go-admin/common/models"
	"gorm.io/gorm"
	"time"
)

type UserPasswordChangeLog struct {
	models.Model

	UserId     int       `json:"userId" gorm:"type:int unsigned;comment:用户ID"`
	RecordTime time.Time `json:"recordTime" gorm:"type:datetime;comment:修改时间时间"`
	RecordName string    `json:"recordName" gorm:"type:varchar(100);comment:修改人"`
}

func (UserPasswordChangeLog) TableName() string {
	return "user_password_change_log"
}

func (e *UserPasswordChangeLog) Generate() *UserPasswordChangeLog {
	o := *e
	return &o
}

func (e *UserPasswordChangeLog) GetId() interface{} {
	return e.Id
}

// AddLog 修改密码日志
func (e *UserPasswordChangeLog) AddLog(db *gorm.DB, userId int, recordName string) error {
	err := db.Create(&UserPasswordChangeLog{
		UserId:     userId,
		RecordTime: time.Now(),
		RecordName: recordName,
	}).Error
	if err != nil {
		return err
	}
	return nil
}
