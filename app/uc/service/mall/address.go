package mall

import (
	"errors"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/uc/models"
	"go-admin/app/uc/service/mall/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type Address struct {
	service.Service
}

// GetPage 获取Address列表
func (e *Address) GetPage(c *dto.AddressGetPageReq, p *actions.DataPermission, list *[]models.Address, count *int64) error {
	var err error
	var data models.Address

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("AddressService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取Address对象
func (e *Address) Get(d *dto.AddressGetReq, p *actions.DataPermission, model *models.Address) error {
	var data models.Address

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetAddress error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建Address对象
func (e *Address) Insert(c *dto.AddressInsertReq) error {
	// 手机或固定电话必须有一个
	if len(c.CellPhone) == 0 && len(c.Telephone) == 0 {
		return errors.New("手机号码、固定电话至少填一个")
	}

	var err error
	var data models.Address
	c.Generate(&data)
	err = e.Orm.Transaction(func(tx *gorm.DB) error {
		// 设置默认的情况下清空其他默认
		if c.IsDefault == 1 {
			err = tx.Model(&models.Address{}).
				Where("address_type = ?", c.AddressType).
				Where("user_id = ?", c.UserId).
				Where("is_default = ?", 1).
				Update("is_default", 0).Error
			if err != nil {
				return err
			}
		}

		// 创建新数据
		err = tx.Create(&data).Error

		// 如果发票地址没有添加收货地址的时候同步给发票地址
		if data.AddressType == 1 && !data.IsAddressTypeExists(e.Orm, 2, data.UserId) {
			data.Id = 0
			data.AddressType = 2
			err = tx.Create(&data).Error
		}
		return err
	})
	if err != nil {
		e.Log.Errorf("AddressService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改Address对象
func (e *Address) Update(c *dto.AddressUpdateReq, p *actions.DataPermission) error {
	// 手机或固定电话必须有一个
	if len(c.CellPhone) == 0 && len(c.Telephone) == 0 {
		return errors.New("手机号码、固定电话至少填一个")
	}

	var err error
	var data = models.Address{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm
	err = db.Transaction(func(tx *gorm.DB) error {
		// 设置默认的情况下清空其他默认
		if c.IsDefault == 1 {
			err = tx.Model(&models.Address{}).
				Where("address_type = ?", c.AddressType).
				Where("user_id = ?", c.UserId).
				Where("is_default = ?", 1).
				Update("is_default", 0).Error
			if err != nil {
				return err
			}
		}

		// 正常更新
		err = tx.Save(&data).Error
		return err
	})
	if err = db.Error; err != nil {
		e.Log.Errorf("AddressService Save error:%s \r\n", err)
		return err
	}
	return nil
}

// Remove 删除Address
func (e *Address) Remove(d *dto.AddressDeleteReq, p *actions.DataPermission) error {
	var data models.Address

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveAddress error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}

// Default 设置默认地址
func (e *Address) Default(d *dto.AddressDefaultReq, p *actions.DataPermission) error {
	db := e.Orm
	var data models.Address
	// 校验这条数据是否存在
	err := db.
		Where("address_type = ?", d.AddressType).
		Where("user_id = ?", d.UserId).
		Where("id = ?", d.Id).First(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("此条地址不存在，请检查")
		} else {
			return err
		}
	}

	db.Transaction(func(tx *gorm.DB) error {
		// 清空之前的默认
		err := tx.Model(&models.Address{}).
			Where("address_type = ?", d.AddressType).
			Where("user_id = ?", d.UserId).
			Where("is_default = ?", 1).
			Scopes(
				actions.Permission(data.TableName(), p),
			).
			Update("is_default", 0).Error
		if err != nil {
			return err
		}

		// 取消默认
		if d.Cancel == 1 {
			return nil
		}

		// 设置新的默认
		err = tx.Model(&data).
			Where("address_type = ?", d.AddressType).
			Where("user_id = ?", d.UserId).
			Where("id = ?", d.Id).
			Scopes(
				actions.Permission(data.TableName(), p),
			).
			Update("is_default", 1).Error
		if err != nil {
			return err
		}
		return nil
	})

	return nil
}
