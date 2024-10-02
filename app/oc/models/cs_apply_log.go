package models

import (
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	"go-admin/common/global"
	"gorm.io/gorm"
	"time"

	"go-admin/common/models"
)

type CsApplyLog struct {
	models.Model

	CsNo         string    `json:"csNo" gorm:"type:varchar(30);comment:售后申请编号(对应于cs_apply.cs_no)"`
	UserId       string    `json:"userId" gorm:"type:varchar(32);comment:备注人id"`
	HandlerLog   string    `json:"handlerLog" gorm:"type:text;comment:处理记录"`
	UserName     string    `json:"userName" gorm:"type:varchar(128);comment:UserName"`
	CreatedAt    time.Time `json:"createdAt" gorm:"comment:创建时间"`
	UpdatedAt    time.Time `json:"updatedAt" gorm:"comment:最后更新时间"`
	CreateBy     int       `json:"createBy" gorm:"index;comment:创建者"`
	CreateByName string    `json:"createByName" gorm:"index;comment:创建者姓名"`
}

func (CsApplyLog) TableName() string {
	return "cs_apply_log"
}

func (e *CsApplyLog) Generate() *CsApplyLog {
	o := *e
	return &o
}

func (e *CsApplyLog) GetId() interface{} {
	return e.Id
}

// Insert 创建CsApplyLog对象
func (e *CsApplyLog) AddLog(db *gorm.DB, csNo string, text string) error {
	ocPrefix := global.GetTenantOcDBNameWithDB(db)

	var err error
	data := CsApplyLog{
		CsNo:         csNo,
		UserId:       "",
		HandlerLog:   text,
		UserName:     "",
		CreateBy:     user.GetUserId(db.Statement.Context.(*gin.Context)),
		CreateByName: user.GetUserName(db.Statement.Context.(*gin.Context)),
	}

	err = db.Table(ocPrefix + "." + e.TableName()).Create(&data).Error
	if err != nil {
		return err
	}
	return nil
}
