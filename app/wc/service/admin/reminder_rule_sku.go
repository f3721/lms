package admin

import (
	"errors"
	modelsPc "go-admin/app/pc/models"
	dtoPc "go-admin/app/pc/service/admin/dto"
	"go-admin/app/wc/models"
	"go-admin/app/wc/service/admin/dto"
	"go-admin/common/actions"
	pcClient "go-admin/common/client/pc"
	cDto "go-admin/common/dto"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type ReminderRuleSku struct {
	service.Service
}

// GetPage 获取ReminderRuleSku列表
func (e *ReminderRuleSku) GetPage(c *dto.ReminderRuleSkuGetPageReq, p *actions.DataPermission, list *[]models.ReminderRuleSku, count *int64) error {
	var err error
	var data models.ReminderRuleSku

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSizeNegative(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(&list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("ReminderRuleSkuService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取ReminderRuleSku对象
func (e *ReminderRuleSku) Get(d *dto.ReminderRuleSkuGetReq, p *actions.DataPermission, model *models.ReminderRuleSku) error {
	var data models.ReminderRuleSku

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetReminderRuleSku error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建ReminderRuleSku对象
func (e *ReminderRuleSku) Insert(c *dto.ReminderRuleSkuInsertReq) error {
	var err error
	var data models.ReminderRuleSku
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("ReminderRuleSkuService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Inserts 批量创建ReminderRuleSku对象
func (e *ReminderRuleSku) Inserts(db *gorm.DB, skus []*dto.ReminderRuleSkuInsertReq, id int, warehouseCode string, gc *gin.Context) error {
	var err error

	if id == 0 {
		return errors.New("数据错误无法插入sku设置")
	}
	err = e.CheckSkuList(skus, warehouseCode)
	if err != nil {
		return err
	}

	skuAdds := []models.ReminderRuleSku{}
	for _, req := range skus {
		data := models.ReminderRuleSku{}
		req.Generate(&data)
		data.ReminderRuleId = id
		skuAdds = append(skuAdds, data)
	}

	err = e.Orm.Create(skuAdds).Error
	if err != nil {
		e.Log.Errorf("ReminderRuleSkuService Insert error:%s \r\n", err)
		return err
	}

	//添加log
	e.AddLog(db, skuAdds, 1, gc)
	return nil
}

func (e *ReminderRuleSku) CheckSkuList(skus []*dto.ReminderRuleSkuInsertReq, warehouseCode string) error {
	visited := make(map[string]bool)
	var skuList = []string{}
	for _, req := range skus {
		if visited[req.SkuCode] {
			return errors.New(req.SkuCode + "SKU不能重复配置")
		}
		visited[req.SkuCode] = true
		skuList = append(skuList, req.SkuCode)
	}
	// 校验sku真实性
	goodsMap := e.GetPcGetGoodsBySku(e.Orm, skuList, warehouseCode)
	skuDiff := make([]string, 0)
	// 遍历待比较的 sku_code 列表，如果不在已有列表中，则添加到差异列表中
	for _, sku := range skuList {
		if _, exists := goodsMap[sku]; !exists {
			skuDiff = append(skuDiff, sku)
		}
	}

	if len(skuDiff) > 0 {
		return errors.New("[" + strings.Join(skuDiff, ",") + "] SKU不存在或已被禁用，不能添加！")
	}
	return nil
}

// Update 修改ReminderRuleSku对象
func (e *ReminderRuleSku) Update(db *gorm.DB, c *dto.ReminderRuleSkuUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.ReminderRuleSku{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	err = db.Save(&data).Error
	if err != nil {
		e.Log.Errorf("ReminderRuleSkuService Save error:%s \r\n", err)
		return err
	}

	addLogData := []models.ReminderRuleSku{}
	addLogData = append(addLogData, data)

	//添加log
	e.AddLog(db, addLogData, 2, nil)
	return nil
}

func (e *ReminderRuleSku) AddLog(db *gorm.DB, adds []models.ReminderRuleSku, logType int, gc *gin.Context) error {
	logAdds := []*models.ReminderRuleSkuLog{}

	for _, add := range adds {
		addData := &models.ReminderRuleSkuLog{}
		copier.Copy(addData, add)
		addData.LogType = logType
		addData.CreatedAt = time.Now()
		addData.Id = 0
		addData.ReminderRuleSkuId = add.Id
		//addData.CreateBy = user.GetUserId(gc)
		//addData.CreateByName = user.GetUserName(gc)

		logAdds = append(logAdds, addData)
	}

	err := db.Create(logAdds).Error
	if err != nil {
		return err
	}

	return nil
}

// Remove 删除ReminderRuleSku
func (e *ReminderRuleSku) Remove(d *dto.ReminderRuleSkuDeleteReq, p *actions.DataPermission) error {
	var data models.ReminderRuleSku

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveReminderRuleSku error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}

func (e *ReminderRuleSku) GetPcGetGoodsBySku(db *gorm.DB, skuSlice []string, warehouseCode string) (productMap map[string]*modelsPc.Goods) {
	result := pcClient.ApiByDbContext(db).GetGoodsBySkuCodeReq(dtoPc.GetGoodsBySkuCodeReq{
		SkuCode:       skuSlice,
		Status:        1,
		WarehouseCode: warehouseCode,
		OnlineStatus:  -1,
	})
	resultInfo := &struct {
		response.Response
		Data []modelsPc.Goods
	}{}
	result.Scan(resultInfo)

	productMap = make(map[string]*modelsPc.Goods)
	for _, goods := range resultInfo.Data {
		productMap[goods.SkuCode] = &goods
	}

	return
}
