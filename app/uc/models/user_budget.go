package models

import (

	"go-admin/common/models"
	"time"
)

type UserBudget struct {
    models.Model
    
    Month string `json:"month" gorm:"type:varchar(10);comment:年月"` 
    UserId int `json:"userId" gorm:"type:int;comment:用户ID"` 
    InitialAmount float64 `json:"initialAmount" gorm:"type:decimal(10,2);comment:预算金额"`
    Spent float64 `json:"spent" gorm:"type:decimal(10,2);comment:已使用金额"`
    Balance float64 `json:"balance" gorm:"type:decimal(10,2);comment:剩余金额"`
    ExcessAmount float64 `json:"excessAmount" gorm:"type:decimal(10,2);comment:超出金额"`
	CreatedAt time.Time      `json:"createdAt" gorm:"comment:创建时间"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"comment:最后更新时间"`
}

func (UserBudget) TableName() string {
    return "user_budget"
}

func (e *UserBudget) GetId() interface{} {
	return e.Id
}
