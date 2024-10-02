package admin

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"github.com/jinzhu/copier"
	"github.com/prometheus/common/log"
	dtoUc "go-admin/app/uc/service/admin/dto"
	dtoWc "go-admin/app/wc/service/admin/dto"
	ucClient "go-admin/common/client/uc"
	wcClient "go-admin/common/client/wc"
	"strconv"
	"strings"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/wc/models"
	"go-admin/app/wc/service/admin/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type ReminderRule struct {
	service.Service
}

// GetPage 获取ReminderRule列表
func (e *ReminderRule) GetPage(c *dto.ReminderRuleGetPageReq, p *actions.DataPermission, list *[]*dto.ReminderRuleData, count *int64) error {
	var err error
	var data models.ReminderRule

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
			actions.SysUserPermission(data.TableName(), p, 1),
			actions.SysUserPermission(data.TableName(), p, 2),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("ReminderRuleService GetPage error:%s \r\n", err)
		return err
	}
	companyByIds := []int{}
	warehouseCodes := []string{}
	for _, ruleData := range *list {
		companyByIds = append(companyByIds, int(ruleData.CompanyId))
		warehouseCodes = append(warehouseCodes, ruleData.WarehouseCode)
	}
	companyNameMap := e.GetOcGetCompanyByIds(companyByIds)
	warehousesByCodeMap := e.GetWcGetWarehousesByCodeMap(e.Orm, warehouseCodes)
	log.Info(companyNameMap)
	for _, ruleData := range *list {
		if v, ok := companyNameMap[ruleData.CompanyId]; ok {
			ruleData.CompanyName = v
		}
		if v, ok := warehousesByCodeMap[ruleData.WarehouseCode]; ok {
			ruleData.WarehouseName = v
		}
	}
	return nil
}

// GetAll 获取ReminderRule全部列表
func (e *ReminderRule) GetAll(c *dto.ReminderRuleGetPageReq, p *actions.DataPermission, list *[]models.ReminderRule) error {
	var err error
	var data models.ReminderRule

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).Error
	if err != nil {
		e.Log.Errorf("ReminderRuleService GetPage error:%s \r\n", err)
		return err
	}

	return nil
}

// Get 获取ReminderRule对象
func (e *ReminderRule) Get(d *dto.ReminderRuleGetReq, p *actions.DataPermission, model *dto.ReminderRuleData) error {
	var data models.ReminderRule

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetReminderRule error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}

	companyByIds := []int{model.CompanyId}
	warehouseCodes := []string{model.WarehouseCode}

	companyNameMap := e.GetOcGetCompanyByIds(companyByIds)
	warehousesByCodeMap := e.GetWcGetWarehousesByCodeMap(e.Orm, warehouseCodes)
	if v, ok := companyNameMap[model.CompanyId]; ok {
		model.CompanyName = v
	}
	if v, ok := warehousesByCodeMap[model.WarehouseCode]; ok {
		model.WarehouseName = v
	}
	return nil
}

// Insert 创建ReminderRule对象
func (e *ReminderRule) Insert(c *dto.ReminderRuleInsertReq, gc *gin.Context, p *actions.DataPermission) (int, error) {
	var err error
	var data models.ReminderRule
	c.Generate(&data)
	isExists, err := e.ExistsByWarehouseCode(data.WarehouseCode, data.CompanyId)
	if isExists == true {
		return 0, errors.New("已存在该公司仓库的补货规则")
	}
	if !contains(strings.Split(p.AuthorityWarehouseId, ","), c.WarehouseCode) {
		return 0, errors.New("没有仓库权限无法创建")
	}
	if !contains(strings.Split(p.AuthorityCompanyId, ","), strconv.Itoa(c.CompanyId)) {
		return 0, errors.New("没有公司权限无法创建")
	}

	tx := e.Orm.Begin()
	err = tx.Create(&data).Error
	if err != nil {
		tx.Rollback()
		e.Log.Errorf("ReminderRuleService Insert error:%s \r\n", err)
		return 0, err
	}

	if c.SkuList != nil && len(c.SkuList) > 0 {
		ReminderRuleSkuS := ReminderRuleSku{}
		ReminderRuleSkuS.Service = e.Service
		ReminderRuleSkuS.Orm = tx
		err = ReminderRuleSkuS.Inserts(tx, c.SkuList, int(data.Id), c.WarehouseCode, gc)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	tx.Commit()
	return int(data.Id), nil
}

func contains(slice []string, element string) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}

// Update 修改ReminderRule对象
func (e *ReminderRule) Update(c *dto.ReminderRuleUpdateReq, p *actions.DataPermission, gc *gin.Context) error {
	var err error
	var data = models.ReminderRule{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	tx := e.Orm.Begin()

	db := tx.Save(&data)
	if err = db.Error; err != nil {
		tx.Rollback()
		e.Log.Errorf("ReminderRuleService Save error:%s \r\n", err)
		return err
	}

	skuS := ReminderRuleSku{e.Service}

	checkSkuListReq := []*dto.ReminderRuleSkuInsertReq{}
	copier.Copy(&checkSkuListReq, c.SkuList)
	err = skuS.CheckSkuList(checkSkuListReq, c.WarehouseCode)
	if err != nil {
		tx.Rollback()
		return err
	}

	insertArr := []*dto.ReminderRuleSkuInsertReq{}
	for _, updateReq := range c.SkuList {
		if updateReq.Id == 0 {
			insertData := &dto.ReminderRuleSkuInsertReq{}
			copier.Copy(insertData, updateReq)
			insertArr = append(insertArr, insertData)
		} else {
			err = skuS.Update(tx, updateReq, p)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	if len(insertArr) > 0 {
		err = skuS.Inserts(tx, insertArr, c.Id, c.WarehouseCode, gc)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

// Remove 删除ReminderRule
func (e *ReminderRule) Remove(d *dto.ReminderRuleDeleteReq, p *actions.DataPermission) error {
	var data models.ReminderRule

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveReminderRule error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}

// ExistsByWarehouseCode 判断是否已存在公司id仓库id
func (e *ReminderRule) ExistsByWarehouseCode(warehouseCode string, companyId int) (bool, error) {
	var err error
	var data models.ReminderRule

	err = e.Orm.Where("warehouse_code", warehouseCode).Where("company_id", companyId).Find(&data).Error
	if err != nil {
		e.Log.Errorf("ReminderRuleService ExistsByWarehouseCode error:%s \r\n", err)
		return false, err
	}

	if data.Id == 0 {
		return false, nil
	}
	return true, nil
}

func (e *ReminderRule) GetOcGetCompanyByIds(ids []int) map[int]string {
	strNumbers := make([]string, len(ids))
	for i, num := range ids {
		strNumbers[i] = fmt.Sprintf("%d", num)
	}

	// 公司名称
	companyResult := ucClient.ApiByDbContext(e.Orm).GetCompanyByIds(strings.Join(strNumbers, ","))
	companyResultInfo := &struct {
		response.Response
		Data struct {
			response.Page
			List []dtoUc.CompanyInfoGetSelectPageData
		}
	}{}
	companyResult.Scan(companyResultInfo)
	companyMap := make(map[int]string, len(companyResultInfo.Data.List))
	for _, company := range companyResultInfo.Data.List {
		companyMap[company.Id] = company.CompanyName
	}

	return companyMap
}

func (e *ReminderRule) GetWcGetWarehousesByCode(db *gorm.DB, code []string) []dtoWc.WarehouseGetPageResp {
	// 获取商品信息
	req := dtoWc.InnerWarehouseGetListReq{
		WarehouseCode: strings.Join(code, ","),
	}
	result := wcClient.ApiByDbContext(db).GetWarehouseList(req)

	resultInfo := &struct {
		response.Response
		Data []dtoWc.WarehouseGetPageResp
	}{}
	result.Scan(resultInfo)

	return resultInfo.Data
}

func (e *ReminderRule) GetWcGetWarehousesByCodeMap(db *gorm.DB, code []string) (dataMap map[string]string) {
	list := e.GetWcGetWarehousesByCode(db, code)
	dataMap = make(map[string]string)
	for _, data := range list {
		dataMap[data.WarehouseCode] = data.WarehouseName
	}

	return dataMap
}
