package admin

import (
	"errors"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/uc/models"
	"go-admin/app/uc/service/admin/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type RoleInfo struct {
	service.Service
}

// GetPage 获取RoleInfo列表
func (e *RoleInfo) GetPage(c *dto.RoleInfoGetPageReq, p *actions.DataPermission, list *[]models.RoleInfo, count *int64) error {
	var err error
	var data models.RoleInfo

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSizeNegative(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Where("role_status = 1").
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("RoleInfoService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取RoleInfo对象
func (e *RoleInfo) Get(d *dto.RoleInfoGetReq, p *actions.DataPermission, model *models.RoleInfo) error {
	var data models.RoleInfo

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetRoleInfo error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Get 获取RoleInfo对象
func (e *RoleInfo) GetByName(name string) (model *models.RoleInfo, err error) {
	var data models.RoleInfo

	err = e.Orm.Model(&data).
		Where("name = ?", name).
		First(model).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		e.Log.Errorf("db error:%s", err)
		return
	}
	return
}

// Get 获取RoleInfo对象
func (e *RoleInfo) GetByNames(name []string) (model *models.RoleInfo, err error) {
	var data models.RoleInfo

	err = e.Orm.Model(&data).
		Where("name in ?", name).
		First(model).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		e.Log.Errorf("db error:%s", err)
		return
	}
	return
}

// Get 获取RoleInfo对象
func (e *RoleInfo) GetStatusOkList() (list []*models.RoleInfo, err error) {
	var data models.RoleInfo

	data.RoleStatus = 1
	e.Orm.Model(&data).Find(&list)

	return
}

// Get 获取RoleInfo对象
func (e *RoleInfo) GetStatusOkMapList() (mapList map[int]*models.RoleInfo, err error) {
	var data models.RoleInfo
	var list []*models.RoleInfo
	mapList = make(map[int]*models.RoleInfo)
	data.RoleStatus = 1
	e.Orm.Model(&data).Find(&list)
	for _, info := range list {
		mapList[info.Id] = info
	}
	return
}

// Insert 创建RoleInfo对象
func (e *RoleInfo) Insert(c *dto.RoleInfoInsertReq) error {
	var err error
	var data models.RoleInfo
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("RoleInfoService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改RoleInfo对象
func (e *RoleInfo) Update(c *dto.RoleInfoUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.RoleInfo{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("RoleInfoService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除RoleInfo
func (e *RoleInfo) Remove(d *dto.RoleInfoDeleteReq, p *actions.DataPermission) error {
	var data models.RoleInfo

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveRoleInfo error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
