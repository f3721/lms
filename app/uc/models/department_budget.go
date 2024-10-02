package models

import (
	"database/sql"
	"errors"
	"go-admin/common/global"
	"go-admin/common/models"
	"go-admin/common/utils"
	"gorm.io/gorm"
	"time"
)

type DepartmentBudget struct {
	models.Model

	Month         string    `json:"month" gorm:"type:varchar(10);comment:年月"`
	DepId         int       `json:"depId" gorm:"type:int;comment:部门ID"`
	InitialAmount float64   `json:"initialAmount" gorm:"type:decimal(10,2);comment:预算金额"`
	Spent         float64   `json:"spent" gorm:"type:decimal(10,2);comment:已使用金额"`
	Balance       float64   `json:"balance" gorm:"type:decimal(10,2);comment:剩余金额"`
	ExcessAmount  float64   `json:"excessAmount" gorm:"type:decimal(10,2);comment:超出金额"`
	CreatedAt     time.Time `json:"createdAt" gorm:"comment:创建时间"`
	UpdatedAt     time.Time `json:"updatedAt" gorm:"comment:最后更新时间"`
}

func (DepartmentBudget) TableName() string {
	return "department_budget"
}

func (e *DepartmentBudget) GetId() interface{} {
	return e.Id
}

// Generate 生成预算
// userID大于0是仅处理单条用户的预算
// mode 1=>初次新增或预算调整, 2=>部门预算和人均预算 模式相互转换
func (e *DepartmentBudget) Generate(tx *gorm.DB, depId, userId, mode int) (err error) {
	users := make([]UserInfo, 0)
	if userId > 0 {
		err = tx.Table("user_info").Where("id = ? and user_status = 1", userId).Find(&users).Error
		if len(users) == 1 {
			depId = users[0].CompanyDepartmentId
		}
	} else {
		err = tx.Table("user_info").Where("company_department_id = ? and user_status = 1", depId).Find(&users).Error
	}
	if err != nil {
		return
	}
	department := CompanyDepartment{}
	err = tx.Table("company_department").First(&department, depId).Error
	if err != nil {
		return
	}

	month := time.Now().Format("200601")
	tx = tx.Debug().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	// 1 使用部门预算
	if department.DepartmentBudget.Valid == true && department.DepartmentBudget.Float64 >= 0 {
		departmentBudget := DepartmentBudget{}
		err = tx.Table("department_budget").Where("month = ? and dep_id = ?", month, depId).First(&departmentBudget).Error
		if err != nil {
			// 1.1 不存在当月部门预算记录 则新增
			if errors.Is(err, gorm.ErrRecordNotFound) {
				departmentBudget.Month = month
				departmentBudget.DepId = depId
				departmentBudget.InitialAmount = department.DepartmentBudget.Float64
				departmentBudget.Balance = department.DepartmentBudget.Float64
				err = tx.Create(&departmentBudget).Error
				if err != nil {
					tx.Rollback()
					return
				}
			} else {
				tx.Rollback()
				return
			}
		} else {
			// 1.2 已存在当月部门预算
			// 1.2.1 部门预算调整金额 || 1.2.2 人均预算改为部门预算 => 根据新的部门预算金额及使用金额重新计算
			if mode == 1 && departmentBudget.InitialAmount != department.DepartmentBudget.Float64 || mode == 2 {
				balance := 0.00
				excessAmount := 0.00
				if department.DepartmentBudget.Float64 > departmentBudget.Spent {
					balance = utils.SubFloat64(department.DepartmentBudget.Float64, departmentBudget.Spent)
				}
				if balance <= 0 {
					excessAmount = utils.SubFloat64(departmentBudget.Spent, department.DepartmentBudget.Float64)
				}
				departmentBudget.InitialAmount = department.DepartmentBudget.Float64
				departmentBudget.Balance = balance
				departmentBudget.ExcessAmount = excessAmount
				err = tx.Save(&departmentBudget).Error
				if err != nil {
					tx.Rollback()
					return
				}
			}
		}
		// 1.3 生成或改为 只记录‘使用金额’的用户预算记录
		if len(users) > 0 {
			for _, user := range users {
				userBudget := UserBudget{}
				err = tx.Table("user_budget").Where("user_id = ? and month = ?", user.Id, month).First(&userBudget).Error
				if err != nil {
					// 不存在当月用户预算
					if errors.Is(err, gorm.ErrRecordNotFound) {
						userBudget.Month = month
						userBudget.UserId = user.Id
						err = tx.Create(&userBudget).Error
						if err != nil {
							tx.Rollback()
							return
						}
					} else {
						tx.Rollback()
						return
					}
				} else if mode == 2 {
					// 已存在当月用户预算
					userBudget.InitialAmount = 0
					userBudget.Balance = 0
					userBudget.ExcessAmount = 0
					err = tx.Save(&userBudget).Error
					if err != nil {
						tx.Rollback()
						return
					}
				}
			}
		}
	} else if department.PersonalBudget.Valid == true && department.PersonalBudget.Float64 >= 0 {
		// 2 使用人均预算
		// 2.1 初次新增或预算调整 || 2.2 部门预算改为人均预算 => 生成或修改用户预算里的 预算金额，并根据已使用金额计算出剩余及超出金额
		num := len(users)
		if num > 0 && (mode == 1 || mode == 2) {
			userSpentTotal := 0.00
			for _, user := range users {
				userBudget := UserBudget{}
				err = tx.Table("user_budget").Where("user_id = ? and month = ?", user.Id, month).First(&userBudget).Error
				if err != nil {
					// 不存在当月用户预算
					if errors.Is(err, gorm.ErrRecordNotFound) {
						userBudget.Month = month
						userBudget.UserId = user.Id
						userBudget.InitialAmount = department.PersonalBudget.Float64
						userBudget.Balance = department.PersonalBudget.Float64
						err = tx.Create(&userBudget).Error
						if err != nil {
							tx.Rollback()
							return
						}
					} else {
						tx.Rollback()
						return
					}
				} else {
					// 已存在当月用户预算
					userSpentTotal = utils.AddFloat64(userSpentTotal, userBudget.Spent)
					// 人均预算修改，重置用户预算
					if userBudget.InitialAmount != department.PersonalBudget.Float64 {
						balance := 0.00
						excessAmount := 0.00
						if department.PersonalBudget.Float64 > userBudget.Spent {
							balance = utils.SubFloat64(department.PersonalBudget.Float64, userBudget.Spent)
						}
						if balance <= 0 {
							excessAmount = utils.SubFloat64(userBudget.Spent, department.PersonalBudget.Float64)
						}
						userBudget.InitialAmount = department.PersonalBudget.Float64
						userBudget.Balance = balance
						userBudget.ExcessAmount = excessAmount
						err = tx.Save(&userBudget).Error
						if err != nil {
							tx.Rollback()
							return
						}
					}
				}
			}
			departmentBudget := DepartmentBudget{}
			err = tx.Table("department_budget").Where("month = ? and dep_id = ?", month, depId).First(&departmentBudget).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					balance := 0.00
					excessAmount := 0.00
					initialAmount := utils.MulFloat64AndInt(department.PersonalBudget.Float64, num)
					if initialAmount > userSpentTotal {
						balance = utils.SubFloat64(initialAmount, userSpentTotal)
					}
					if balance <= 0 {
						excessAmount = utils.SubFloat64(userSpentTotal, initialAmount)
					}
					departmentBudget.Month = month
					departmentBudget.DepId = depId
					departmentBudget.InitialAmount = initialAmount
					departmentBudget.Spent = userSpentTotal
					departmentBudget.Balance = balance
					departmentBudget.ExcessAmount = excessAmount
					err = tx.Create(&departmentBudget).Error
					if err != nil {
						tx.Rollback()
						return
					}
				} else {
					tx.Rollback()
					return
				}
			} else {
				initialAmount := 0.00
				if userId > 0 {
					initialAmount = utils.AddFloat64(departmentBudget.InitialAmount, department.PersonalBudget.Float64)
					userSpentTotal = userSpentTotal + departmentBudget.Spent
				} else {
					initialAmount = utils.MulFloat64AndInt(department.PersonalBudget.Float64, num)
				}
				balance := 0.00
				excessAmount := 0.00
				if initialAmount > userSpentTotal {
					balance = utils.SubFloat64(initialAmount, userSpentTotal)
				}
				if balance <= 0 {
					excessAmount = utils.SubFloat64(userSpentTotal, initialAmount)
				}
				departmentBudget.InitialAmount = initialAmount
				departmentBudget.Spent = userSpentTotal
				departmentBudget.Balance = balance
				departmentBudget.ExcessAmount = excessAmount
				err = tx.Save(&departmentBudget).Error
				if err != nil {
					tx.Rollback()
					return
				}
			}
		}
	}

	tx.Commit()
	return
}

// CheckIsOverBudget 是否超预算
// userId 用户id
// money 金额
// month 月份
func (e *DepartmentBudget) CheckIsOverBudget(tx *gorm.DB, userId int, money float64, month string) bool {
	ucPrefix := global.GetTenantUcDBNameWithDB(tx)
	var dep struct {
		Id               int
		DepartmentBudget sql.NullFloat64
		PersonalBudget   sql.NullFloat64
	}

	err := tx.Table(ucPrefix+".company_department t").
		Select("t.id, t.department_budget, t.personal_budget").
		Joins("left join "+ucPrefix+".user_info ui on ui.company_department_id = t.id").
		Where("ui.id = ?", userId).
		Scan(&dep).Error

	if err != nil {
		return false
	}
	if dep.DepartmentBudget.Valid == false && dep.PersonalBudget.Valid == false {
		return false
	}
	// 部门预算
	if dep.DepartmentBudget.Float64 >= 0 {
		departmentBudget := DepartmentBudget{}
		err = tx.Table(ucPrefix+".department_budget").Where("month = ? and dep_id = ?", month, dep.Id).First(&departmentBudget).Error
		if err == nil {
			return departmentBudget.Balance < money
		}
	} else if dep.PersonalBudget.Float64 >= 0 {
		userBudget := UserBudget{}
		err = tx.Table(ucPrefix+".user_budget").Where("user_id = ? and month = ?", userId, month).First(&userBudget).Error
		if err == nil {
			return userBudget.Balance < money
		}
	}
	return false
}

// UpdateBudget 修改预算
// userId 用户id
// money 金额
// month 月份
func (e *DepartmentBudget) UpdateBudget(tx *gorm.DB, userId int, money float64, month string) (err error) {
	ucPrefix := global.GetTenantUcDBNameWithDB(tx)
	var dep struct {
		Id               int
		DepartmentBudget sql.NullFloat64
		PersonalBudget   sql.NullFloat64
	}
	err = tx.Table(ucPrefix+".company_department t").
		Select("t.id, t.department_budget, t.personal_budget").
		Joins("left join "+ucPrefix+".user_info ui on ui.company_department_id = t.id").
		Where("ui.id = ?", userId).
		Scan(&dep).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) || (dep.DepartmentBudget.Valid == false && dep.PersonalBudget.Valid == false) {
		return nil
	}

	userBudget := UserBudget{}
	err = tx.Table(ucPrefix+".user_budget").Where("user_id = ? and month = ?", userId, month).First(&userBudget).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	// 更新预算
	if dep.DepartmentBudget.Float64 >= 0 || dep.PersonalBudget.Float64 >= 0 {
		// 用户预算
		userBudget.Spent = utils.AddFloat64(userBudget.Spent, money)
		if dep.DepartmentBudget.Valid == false && dep.PersonalBudget.Float64 >= 0 {
			if userBudget.InitialAmount > userBudget.Spent {
				userBudget.Balance = utils.SubFloat64(userBudget.InitialAmount, userBudget.Spent)
			} else {
				userBudget.Balance = 0
			}
			if userBudget.Spent > userBudget.InitialAmount {
				userBudget.ExcessAmount = utils.SubFloat64(userBudget.Spent, userBudget.InitialAmount)
			} else {
				userBudget.ExcessAmount = 0
			}
		}
		err = tx.Table(ucPrefix + ".user_budget").Save(&userBudget).Error
		if err != nil {
			return
		}

		// 部门预算
		departmentBudget := DepartmentBudget{}
		err = tx.Table(ucPrefix+".department_budget").Where("month = ? and dep_id = ?", month, dep.Id).First(&departmentBudget).Error
		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		departmentBudget.Spent = utils.AddFloat64(departmentBudget.Spent, money)

		if departmentBudget.InitialAmount > departmentBudget.Spent {
			departmentBudget.Balance = utils.SubFloat64(departmentBudget.InitialAmount, departmentBudget.Spent)
		} else {
			departmentBudget.Balance = 0
		}
		if departmentBudget.Spent > departmentBudget.InitialAmount {
			departmentBudget.ExcessAmount = utils.SubFloat64(departmentBudget.Spent, departmentBudget.InitialAmount)
		} else {
			departmentBudget.ExcessAmount = 0
		}
		err = tx.Table(ucPrefix + ".department_budget").Save(&departmentBudget).Error
		if err != nil {
			return
		}
	}
	return
}
