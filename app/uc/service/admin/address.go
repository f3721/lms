package admin

import (
	"errors"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"go-admin/app/uc/models"
	"go-admin/app/uc/service/admin/dto"
	modelsWc "go-admin/app/wc/models"
	"go-admin/common/actions"
	wcClient "go-admin/common/client/wc"
	cDto "go-admin/common/dto"
	"gorm.io/gorm"
	"strconv"
	"strings"
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

// GetPage 获取Address列表
func (e *Address) GetListPage(c *dto.AddressGetPageReq, p *actions.DataPermission, resList *[]dto.AddressGetPageRes, count *int64) error {
	list := &[]models.Address{}
	err := e.GetPage(c, p, list, count)
	if err != nil {
		return err
	}
	for _, address := range *list {
		*resList = append(*resList, dto.AddressGetPageRes{Address: address, AddressTypeText: models.AddressType[address.AddressType]})
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
	var err error
	var data models.Address
	var saveData models.Address
	c.Generate(&saveData)

	err = e.SaveGenerate(&data, &saveData)
	if err != nil {
		return err
	}
	tx := e.Orm.Begin()

	// 设置默认的情况下清空其他默认
	if c.IsDefault == 1 {
		err = tx.Model(&models.Address{}).
			Where("address_type = ?", c.AddressType).
			Where("user_id = ?", c.UserId).
			Where("is_default = ?", 1).
			Update("is_default", 0).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Create(&saveData).Error
	if err != nil {
		e.Log.Errorf("AddressService Insert error:%s \r\n", err)
		tx.Rollback()
		return err
	}
	//日志
	operateLogsService := OperateLogs{e.Service}
	operateLogsService.AddLog(saveData.UserId, data, saveData, models.AddressModelName, models.AddressOperationCreate, c.CreateBy, c.CreateByName)

	// 如果发票地址没有添加收货地址的时候同步给发票地址
	if saveData.AddressType == 1 && !data.IsAddressTypeExists(tx, 2, saveData.UserId) {
		saveData.Id = 0
		saveData.AddressType = 2
		tx.Create(&saveData)
	}

	tx.Commit()
	return nil
}

// Update 修改Address对象
func (e *Address) Update(c *dto.AddressUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.Address{}
	var saveData = models.Address{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())

	saveData = data
	c.Generate(&saveData)
	e.SaveGenerate(&data, &saveData)

	db := e.Orm.Save(&saveData)
	if err = db.Error; err != nil {
		e.Log.Errorf("AddressService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}

	//日志
	operateLogsService := OperateLogs{e.Service}
	operateLogsService.AddLog(saveData.UserId, data, saveData, models.AddressModelName, models.AddressOperationUpdate, c.UpdateBy, c.UpdateByName)
	return nil
}

// Update 修改Address对象
func (e *Address) SaveGenerate(data *models.Address, saveData *models.Address) error {
	if saveData.Telephone == "" && saveData.CellPhone == "" {
		return errors.New("手机号/座机（二选一）不能为空！")
	}
	if data.Id == 0 {
		if data.AddressType == 0 {
			data.AddressType = 1
		}

		//	获取用户的公司信息
		userInfoService := UserInfo{e.Service}
		userCompany, _ := userInfoService.GetUserCompanyByUserId(saveData.UserId)

		saveData.CompanyId = userCompany.CompanyId
		saveData.CompanyName = userCompany.CompanyName
	}
	if saveData.CountryId == 1 {
		if saveData.ProvinceId == 0 || saveData.CityId == 0 || saveData.AreaId == 0 || saveData.TownId == 0 {
			return errors.New("省、市、区、镇/街道必选")
		}

		// 获取省市区中文名
		regionIdList := []string{
			strconv.Itoa(saveData.CountryId),
			strconv.Itoa(saveData.ProvinceId),
			strconv.Itoa(saveData.CityId),
			strconv.Itoa(saveData.AreaId),
			strconv.Itoa(saveData.TownId),
		}
		regionIds := strings.Join(regionIdList, ",")
		regionResult := wcClient.ApiByDbContext(e.Orm).GetRegionByIds(regionIds)
		regionResultInfo := &struct {
			response.Response
			Data []modelsWc.Region
		}{}
		regionResult.Scan(regionResultInfo)
		for _, region := range regionResultInfo.Data {
			if region.Id == saveData.CountryId {
				saveData.CountryName = region.Name
			}
			if region.Id == saveData.ProvinceId {
				saveData.ProvinceName = region.Name
			}
			if region.Id == saveData.CityId {
				saveData.CityName = region.Name
			}
			if region.Id == saveData.AreaId {
				saveData.AreaName = region.Name
			}
			if region.Id == saveData.TownId {
				saveData.TownName = region.Name
			}
		}
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
