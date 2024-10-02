package admin

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"github.com/samber/lo"
	dtoUc "go-admin/app/uc/service/admin/dto"
	adminClient "go-admin/common/client/admin"
	ucClient "go-admin/common/client/uc"
	"go-admin/common/utils"
	"gorm.io/gorm"
	"strconv"
	"strings"

	"go-admin/app/wc/models"
	"go-admin/app/wc/service/admin/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type Warehouse struct {
	service.Service
}

// GetPage 获取Warehouse列表
func (e *Warehouse) GetPage(c *dto.WarehouseGetPageReq, p *actions.DataPermission, outData *[]dto.WarehouseGetPageResp, count *int64) error {
	var err error
	var data models.Warehouse
	var list = &[]models.Warehouse{}
	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSizeNegative(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
			actions.SysUserPermission(data.TableName(), p, 2),
		).Order("id DESC").
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("WarehouseService GetPage error:%s \r\n", err)
		return err
	}
	if err := e.FormatPageOutput(list, outData); err != nil {
		e.Log.Errorf("WarehouseService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

func (e *Warehouse) FormatPageOutput(list *[]models.Warehouse, outData *[]dto.WarehouseGetPageResp) error {
	companyIds := lo.Uniq(lo.Map(*list, func(item models.Warehouse, index int) string {
		return strconv.Itoa(item.CompanyId)
	}))
	companyIdsStr := strings.Join(companyIds, ",")
	companyMap := e.GetCompanyMapByClient(companyIdsStr)
	for _, item := range *list {
		dataItem := dto.WarehouseGetPageResp{}
		if err := utils.CopyDeep(&dataItem, &item); err != nil {
			return err
		}
		dataItem.IsVirtualName = utils.GetTFromMap(dataItem.IsVirtual, models.WarehouseModeIsVirtualMap)
		dataItem.CompanyName = companyMap[dataItem.CompanyId]
		*outData = append(*outData, dataItem)
	}
	return nil
}

func (e *Warehouse) GetCompanyMapByClient(companyIds string) map[int]string {
	result := ucClient.ApiByDbContext(e.Orm).GetCompanyByIds(companyIds)
	resultInfo := &struct {
		response.Response
		Data struct {
			response.Page
			List []dtoUc.CompanyInfoGetSelectPageData
		}
	}{}
	result.Scan(resultInfo)

	return lo.Associate(resultInfo.Data.List, func(f dtoUc.CompanyInfoGetSelectPageData) (int, string) {
		return f.Id, f.CompanyName
	})
}

// Get 获取Warehouse对象
func (e *Warehouse) Get(d *dto.WarehouseGetReq, p *actions.DataPermission, outData *dto.WarehouseGetResp) error {
	var data models.Warehouse

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
			actions.SysUserPermission(data.TableName(), p, 2),
		).Where("id= ?", d.GetId()).
		Scan(outData).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetWarehouse error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}

	companyMap := e.GetCompanyMapByClient(strconv.Itoa(outData.CompanyId))
	outData.CompanyName = companyMap[outData.CompanyId]
	return nil
}

// Insert 创建Warehouse对象
func (e *Warehouse) Insert(c *dto.WarehouseInsertReq) error {
	var err error
	var data models.Warehouse

	if sameFlag, _ := data.CheckWhName(e.Orm, c.WarehouseName, 0); !sameFlag {
		return errors.New("实体仓名称重复")
	}
	//检查公司
	/*if !e.CheckCompany(c.CompanyId) {
		return errors.New("公司信息异常")
	}*/
	c.Generate(e.Orm, &data)
	if _, err = data.GenerateWhCode(e.Orm); err != nil {
		return err
	}

	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("WarehouseService Insert error:%s \r\n", err)
		return err
	}

	//记录操作日志
	DataStr, _ := json.Marshal(&data)
	opLog := models.OperateLogs{
		DataId:       data.WarehouseCode,
		ModelName:    models.WarehouseModelName,
		Type:         models.WarehouseModeInsert,
		DoStatus:     models.WarehouseModeStatus[data.Status],
		Before:       "",
		Data:         string(DataStr),
		After:        string(DataStr),
		OperatorId:   c.CreateBy,
		OperatorName: c.CreateByName,
	}
	_ = opLog.InsertItem(e.Orm)

	_ = adminClient.ApiByDbContext(e.Orm).
		UpdatePermission(adminClient.ApiClientPermissionUpdateRequest{
			UserID:                    user.GetUserId(e.Orm.Statement.Context.(*gin.Context)),
			AuthorityWarehouseID:      data.WarehouseCode,
			AuthorityWarehouseAllocID: data.WarehouseCode,
		})

	return nil
}

// 根据公司id检查公司

func (e *Warehouse) CheckCompany(companyId int) bool {
	companyAccess := &struct {
		response.Response
		Data dtoUc.CompanyInfoIsAvailableRes `json:"data,omitempty"`
	}{}
	result := ucClient.ApiByDbContext(e.Orm).CompanyIsAvailable(companyId)
	result.Scan(companyAccess)
	return companyAccess.Data.IsAvailable
}

// Update 修改Warehouse对象
func (e *Warehouse) Update(c *dto.WarehouseUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.Warehouse{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
		actions.SysUserPermission(data.TableName(), p, 2),
	).First(&data, c.GetId())
	if data.Id == 0 {
		return errors.New("实体仓不存在或没有权限")
	}
	oldData := data
	if sameFlag, _ := data.CheckWhName(e.Orm, c.WarehouseName, data.Id); !sameFlag {
		return errors.New("实体仓名称重复")
	}

	c.Generate(e.Orm, &data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("WarehouseService Save error:%s \r\n", err)
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
		ModelName:    models.WarehouseModelName,
		Type:         models.WarehouseModeUpdate,
		DoStatus:     models.WarehouseModeStatus[data.Status],
		Before:       string(oldDataStr),
		Data:         string(DataStr),
		After:        string(DataStr),
		OperatorId:   c.UpdateBy,
		OperatorName: c.UpdateByName,
	}
	_ = opLog.InsertItem(e.Orm)
	return nil
}

// Remove 删除Warehouse
func (e *Warehouse) Remove(d *dto.WarehouseDeleteReq, p *actions.DataPermission) error {
	var data models.Warehouse

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveWarehouse error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}

func (e *Warehouse) Select(c *dto.WarehouseSelectReq, p *actions.DataPermission, resData *[]dto.WarehouseSelectResp, count *int64) error {
	var err error
	var data models.Warehouse
	var list = &[]models.Warehouse{}
	var outItem = &dto.WarehouseSelectResp{}
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
		e.Log.Errorf("WarehouseService Select error:%s \r\n", err)
		return err
	}
	for _, item := range *list {
		outItem.ReGenerate(&item)
		*resData = append(*resData, *outItem)
	}
	return nil
}

func (e *Warehouse) GetByWarehouseCode(d *dto.WarehouseGetByCodeReq, p *actions.DataPermission, model *models.Warehouse) error {
	var data models.Warehouse

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		Where("warehouse_code = ?", d.GetCode()).
		First(model).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetByWarehouseCode error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

func (e *Warehouse) InnerGetListByNameAndCompanyId(c *dto.InnerWarehouseGetListByNameAndCompanyIdReq, list *[]models.Warehouse) error {
	var err error
	var data models.Warehouse
	searchSlice := [][]interface{}{}
	if len(c.Query) == 0 {
		return errors.New("参数不能为空")
	}
	for _, item := range c.Query {
		searchSlice = append(searchSlice, []interface{}{
			item.WarehouseName,
			item.CompanyId,
		})
	}
	err = e.Orm.Model(&data).
		Where("status = ?", models.WarehouseModeStatus1).
		Scopes(
			func(db *gorm.DB) *gorm.DB {
				if len(c.Query) > 0 {
					db.Where("(warehouse_name,company_id) IN ?", searchSlice)
				}
				return db
			},
		).
		Find(list).Error
	if err != nil {
		e.Log.Errorf("WarehouseService InnerGetList error:%s \r\n", err)
		return err
	}
	return nil
}

func (e *Warehouse) InnerGetList(c *dto.InnerWarehouseGetListReq, outData *[]dto.WarehouseGetPageResp) error {
	var err error
	var data models.Warehouse
	var list = &[]models.Warehouse{}
	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			func(db *gorm.DB) *gorm.DB {
				if c.WarehouseCode != "" {
					db.Where("warehouse_code IN ?", lo.Uniq(utils.Split(c.WarehouseCode)))
				}
				if c.WarehouseName != "" {
					db.Where("warehouse_name IN ?", lo.Uniq(utils.Split(c.WarehouseName)))
				}
				return db
			},
		).
		Find(list).Error
	if err != nil {
		e.Log.Errorf("WarehouseService InnerGetList error:%s \r\n", err)
		return err
	}
	if err := e.FormatPageOutput(list, outData); err != nil {
		e.Log.Errorf("WarehouseService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// GetCompanyWarehouseTree 获取Warehouse对象
func (e *Warehouse) GetCompanyWarehouseTree() (outData []dto.CompanyWarehouseTreeResp, err error) {
	result := ucClient.ApiByDbContext(e.Orm).GetCompanyByIds("")
	resultInfo := &struct {
		response.Response
		Data struct {
			response.Page
			List []dtoUc.CompanyInfoGetSelectPageData
		}
	}{}
	result.Scan(resultInfo)
	dataLen := len(resultInfo.Data.List)
	if dataLen <= 0 {
		return
	}
	outData = make([]dto.CompanyWarehouseTreeResp, dataLen)
	var companyIds []int
	for _, v := range resultInfo.Data.List {
		companyIds = append(companyIds, v.Id)
	}

	var warehouseList []struct {
		CompanyId     int    `json:"companyId"`
		WarehouseCode string `json:"warehouseCode"`
		WarehouseName string `json:"warehouseName"`
	}

	var data models.Warehouse
	err = e.Orm.Model(&data).Where("company_id in ?", companyIds).
		Scan(&warehouseList).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}

	warehouseMap := make(map[int][]dto.CompanyWarehouseTreeChildren, len(warehouseList))
	for _, warehouse := range warehouseList {
		warehouseMap[warehouse.CompanyId] = append(warehouseMap[warehouse.CompanyId], dto.CompanyWarehouseTreeChildren{
			Value: warehouse.WarehouseCode,
			Label: warehouse.WarehouseName,
		})
	}

	for i, v := range resultInfo.Data.List {
		tmp := dto.CompanyWarehouseTreeResp{
			Value: v.Id,
			Label: v.CompanyName,
		}
		if _, ok := warehouseMap[v.Id]; ok {
			tmp.Children = warehouseMap[v.Id]
		}
		outData[i] = tmp
	}
	return
}
