package admin

import (
	"encoding/json"
	"errors"
	"go-admin/common/utils"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/wc/models"
	"go-admin/app/wc/service/admin/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type LogicWarehouse struct {
	service.Service
}

// GetPage 获取LogicWarehouse列表
func (e *LogicWarehouse) GetPage(c *dto.LogicWarehouseGetPageReq, p *actions.DataPermission, outData *[]dto.LogicWarehouseGetPageResp, count *int64) error {
	var err error
	var data models.LogicWarehouse
	var list = &[]models.LogicWarehouse{}
	err = e.Orm.Model(&data).Preload("Warehouse").
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSizeNegative(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
			actions.SysUserPermission(data.TableName(), p, 2),
		).Order("id DESC").
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("LogicWarehouseService GetPage error:%s \r\n", err)
		return err
	}
	if err := e.FormatPageOutput(list, outData); err != nil {
		e.Log.Errorf("LogicWarehouseService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

func (e *LogicWarehouse) FormatPageOutput(list *[]models.LogicWarehouse, outData *[]dto.LogicWarehouseGetPageResp) error {
	for _, item := range *list {
		dataItem := dto.LogicWarehouseGetPageResp{}
		if err := utils.CopyDeep(&dataItem, &item); err != nil {
			return err
		}
		dataItem.WarehouseName = dataItem.Warehouse.WarehouseName
		dataItem.TypeName = utils.GetTFromMap(dataItem.Type, models.LwhTypeMap)
		*outData = append(*outData, dataItem)
	}
	return nil
}

// Get 获取LogicWarehouse对象
func (e *LogicWarehouse) Get(d *dto.LogicWarehouseGetReq, p *actions.DataPermission, outData *dto.LogicWarehouseGetResp) error {
	var data models.LogicWarehouse
	var model = &models.LogicWarehouse{}
	err := e.Orm.Model(&data).Preload("Warehouse").
		Scopes(
			actions.Permission(data.TableName(), p),
			actions.SysUserPermission(data.TableName(), p, 2),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetLogicWarehouse error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	if err := utils.CopyDeep(outData, model); err != nil {
		return err
	}
	outData.WarehouseName = model.Warehouse.WarehouseName
	outData.TypeName = utils.GetTFromMap(outData.Type, models.LwhTypeMap)

	return nil
}

// Insert 创建LogicWarehouse对象
func (e *LogicWarehouse) Insert(c *dto.LogicWarehouseInsertReq) error {
	var err error
	var warehouse models.Warehouse
	var data models.LogicWarehouse
	if err := warehouse.GetByWarehouseCode(e.Orm, c.WarehouseCode); err != nil {
		return err
	}
	if resFlag, _ := data.CheckLwhName(e.Orm, c.LogicWarehouseName, 0); !resFlag {
		return errors.New("逻辑仓名称重复")
	}
	if resFlag, _ := data.CheckLwhExist(e.Orm, c.WarehouseCode, c.Type, 0); !resFlag {
		return errors.New("当前实体仓库下已存在" + models.LwhTypeMap[c.Type])
	}

	c.Generate(&data)
	if _, err = data.GenerateLwhCode(e.Orm); err != nil {
		return err
	}

	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("LogicWarehouseService Insert error:%s \r\n", err)
		return err
	}

	//记录操作日志
	DataStr, _ := json.Marshal(&data)
	opLog := models.OperateLogs{
		DataId:       data.LogicWarehouseCode,
		ModelName:    models.LogicWarehouseModelName,
		Type:         models.LogicWarehouseModeInsert,
		DoStatus:     models.LogicWarehouseModeStatus[data.Status],
		Before:       "",
		Data:         string(DataStr),
		After:        string(DataStr),
		OperatorId:   c.CreateBy,
		OperatorName: c.CreateByName,
	}
	_ = opLog.InsertItem(e.Orm)
	return nil
}

// Update 修改LogicWarehouse对象
func (e *LogicWarehouse) Update(c *dto.LogicWarehouseUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.LogicWarehouse{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
		actions.SysUserPermission(data.TableName(), p, 2),
	).First(&data, c.GetId())

	if data.Id == 0 {
		return errors.New("逻辑仓不存在或没有权限")
	}
	if resFlag, _ := data.CheckLwhName(e.Orm, c.LogicWarehouseName, data.Id); !resFlag {
		return errors.New("逻辑仓名称重复")
	}
	oldData := data
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("LogicWarehouseService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}

	//记录操作日志
	oldDataStr, _ := json.Marshal(&oldData)
	DataStr, _ := json.Marshal(&data)
	opLog := models.OperateLogs{
		DataId:       data.WarehouseCode,
		ModelName:    models.LogicWarehouseModelName,
		Type:         models.LogicWarehouseModeUpdate,
		DoStatus:     models.LogicWarehouseModeStatus[data.Status],
		Before:       string(oldDataStr),
		Data:         string(DataStr),
		After:        string(DataStr),
		OperatorId:   c.UpdateBy,
		OperatorName: c.UpdateByName,
	}
	_ = opLog.InsertItem(e.Orm)

	return nil
}

// Remove 删除LogicWarehouse
func (e *LogicWarehouse) Remove(d *dto.LogicWarehouseDeleteReq, p *actions.DataPermission) error {
	var data models.LogicWarehouse

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveLogicWarehouse error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}

func (e *LogicWarehouse) Select(c *dto.LogicWarehouseSelectReq, p *actions.DataPermission, resData *[]dto.LogicWarehouseSelectResp, count *int64) error {
	var err error
	var data models.LogicWarehouse
	var list = &[]models.LogicWarehouse{}
	var outItem = &dto.LogicWarehouseSelectResp{}
	var sysUserPermission = actions.SysUserPermission(data.TableName(), p, 2)

	if c.IsTransfer == "1" {
		sysUserPermission = actions.SysUserPermission(data.TableName(), p, 3)
	}
	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSizeNegative(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
			sysUserPermission,
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("LogicWarehouseService Select error:%s \r\n", err)
		return err
	}
	for _, item := range *list {
		outItem.ReGenerate(&item)
		*resData = append(*resData, *outItem)
	}
	return nil
}
