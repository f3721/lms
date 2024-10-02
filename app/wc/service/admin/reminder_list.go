package admin

import (
	"errors"
	"fmt"
	modelsOc "go-admin/app/oc/models"
	dtoUc "go-admin/app/uc/service/admin/dto"
	dtoWc "go-admin/app/wc/service/admin/dto"
	ucClient "go-admin/common/client/uc"
	wcClient "go-admin/common/client/wc"
	"go-admin/common/global"
	cModel "go-admin/common/models"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"github.com/prometheus/common/log"
	"gorm.io/driver/mysql"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/wc/models"
	"go-admin/app/wc/service/admin/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type ReminderList struct {
	service.Service
}

// GetPage 获取ReminderList列表
func (e *ReminderList) GetPage(c *dto.ReminderListGetPageReq, p *actions.DataPermission, list *[]*dto.ReminderListData, count *int64) error {
	var err error
	var data models.ReminderList

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
			actions.SysUserPermission(data.TableName(), p, 1),
			actions.SysUserPermission(data.TableName(), p, 2),
			actions.SysUserPermission(data.TableName(), p, 4),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("ReminderListService GetPage error:%s \r\n", err)
		return err
	}

	companyByIds := []int{}
	warehouseCodes := []string{}
	vendorIds := []int{}
	for _, ruleData := range *list {
		companyByIds = append(companyByIds, int(ruleData.CompanyId))
		warehouseCodes = append(warehouseCodes, ruleData.WarehouseCode)
		vendorIds = append(vendorIds, ruleData.VendorId)
	}
	companyNameMap := e.GetUcGetCompanyByIds(companyByIds)
	warehousesByCodeMap := e.GetWcGetWarehousesByCodeMap(e.Orm, warehouseCodes)
	vendorsByCodeMap := e.GetWcGetVendorsByCodeMap(e.Orm, vendorIds)
	for _, listData := range *list {
		if v, ok := companyNameMap[listData.CompanyId]; ok {
			listData.CompanyName = v
		}
		if v, ok := warehousesByCodeMap[listData.WarehouseCode]; ok {
			listData.WarehouseName = v
		}
		if v, ok := vendorsByCodeMap[listData.VendorId]; ok {
			listData.VendorName = v
		}
	}

	return nil
}

// Get 获取ReminderList对象
func (e *ReminderList) Get(d *dto.ReminderListGetReq, p *actions.DataPermission, model *dto.ReminderListData) error {
	var data models.ReminderList

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetReminderList error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}

	companyByIds := []int{model.CompanyId}
	warehouseCodes := []string{model.WarehouseCode}
	vendorIds := []int{model.VendorId}

	companyNameMap := e.GetUcGetCompanyByIds(companyByIds)
	warehousesByCodeMap := e.GetWcGetWarehousesByCodeMap(e.Orm, warehouseCodes)
	vendorsByCodeMap := e.GetWcGetVendorsByCodeMap(e.Orm, vendorIds)
	if v, ok := companyNameMap[model.CompanyId]; ok {
		model.CompanyName = v
	}
	if v, ok := warehousesByCodeMap[model.WarehouseCode]; ok {
		model.WarehouseName = v
	}
	if v, ok := vendorsByCodeMap[model.VendorId]; ok {
		model.VendorName = v
	}
	return nil
}

// Insert 创建ReminderList对象
func (e *ReminderList) Insert(c *dto.ReminderListInsertReq) error {
	var err error
	var data models.ReminderList
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("ReminderListService Insert error:%s \r\n", err)
		return err
	}

	c.Id = data.Id
	return nil
}

// Update 修改ReminderList对象
func (e *ReminderList) Update(c *dto.ReminderListUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.ReminderList{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("ReminderListService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除ReminderList
func (e *ReminderList) Remove(d *dto.ReminderListDeleteReq, p *actions.DataPermission) error {
	var data models.ReminderList

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveReminderList error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}

func (e *ReminderList) GetUcGetCompanyByIds(ids []int) map[int]string {
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

func (e *ReminderList) GetWcGetWarehousesByCode(db *gorm.DB, code []string) []dtoWc.WarehouseGetPageResp {
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

func (e *ReminderList) GetWcGetWarehousesByCodeMap(db *gorm.DB, code []string) (dataMap map[string]string) {
	list := e.GetWcGetWarehousesByCode(db, code)
	dataMap = make(map[string]string)
	for _, data := range list {
		dataMap[data.WarehouseCode] = data.WarehouseName
	}

	return dataMap
}

func (e *ReminderList) GetWcGetVendorsByCodeMap(db *gorm.DB, ids []int) map[int]string {
	strNumbers := make([]string, len(ids))
	for i, num := range ids {
		strNumbers[i] = fmt.Sprintf("%d", num)
	}

	// 获取货主名称
	list := []models.Vendors{}
	vendorsService := Vendors{e.Service}
	vendorsService.InnerGetList(&dto.InnerVendorsGetListReq{Ids: strings.Join(strNumbers, ",")}, &list)
	dataMap := make(map[int]string)
	for _, data := range list {
		dataMap[data.Id] = data.NameZh
	}

	return dataMap
}

// GetPage 获取LogicWarehouse列表
func (e *ReminderList) LogicWarehouseGetAll(c *dto.LogicWarehouseGetPageReq, p *actions.DataPermission, list *[]models.LogicWarehouse) error {
	var err error
	var data models.LogicWarehouse
	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			actions.Permission(data.TableName(), p),
		).Order("id DESC").
		Find(list).Error
	if err != nil {
		e.Log.Errorf("LogicWarehouseService GetPage error:%s \r\n", err)
		return err
	}

	return nil
}

func (e *ReminderList) TenantsCreate(c *gin.Context, p *actions.DataPermission) (echoMsg string, err error) {
	tenants := global.GetTenants()
	if len(tenants) <= 0 {
		err = errors.New("暂无租户")
		return
	}
	echoMsg = "补货清单脚本执行开始："
	for tenantKey, tenant := range tenants {
		tenantDbSource := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v_%v?charset=utf8mb4&parseTime=True&loc=Local&timeout=1000ms",
			tenant.DatabaseUsername,
			tenant.DatabasePassword,
			tenant.DatabaseHost,
			tenant.DatabasePort,
			tenant.TenantDBPrefix(),
			"wc")
		tenantDb, err := gorm.Open(mysql.Open(tenantDbSource), &gorm.Config{})
		echoMsg = fmt.Sprintf("%s\r\n租户%s[%s]:", echoMsg, tenant.Name, tenant.DatabaseName)

		c.Request.Header.Set("tenant-id", tenantKey)
		tenantDb = tenantDb.Session(&gorm.Session{
			Context: c,
		})
		e.Orm = tenantDb
		echoMsg = fmt.Sprintf("%s\r\n租户%s[%s]:", echoMsg, tenant.Name, tenant.DatabaseName)
		if err != nil {
			echoMsg = fmt.Sprintf("%s 数据库连接失败", echoMsg)
			continue
		}
		err = e.Create(p)
		if err != nil {
			echoMsg = fmt.Sprintf("%s[%s]fail | ", echoMsg, tenant.DatabaseUsername)
		} else {
			echoMsg = fmt.Sprintf("%s[%s]success | ", echoMsg, tenant.DatabaseUsername)
		}
	}
	echoMsg = fmt.Sprintf("%s\r\n 脚本执行完毕", echoMsg)

	return echoMsg, nil
}

func (e *ReminderList) Create(p *actions.DataPermission) (err error) {
	pcPrefix := global.GetTenantPcDBNameWithDB(e.Orm)

	queryStatus := 1
	req := dto.ReminderRuleGetPageReq{Status: strconv.Itoa(queryStatus)}

	list := make([]models.ReminderRule, 0)

	reminderRuleService := ReminderRule{e.Service}
	// 获取设置的补货规则list
	err = reminderRuleService.GetAll(&req, p, &list)
	if err != nil {
		return
	}

	//循环查询单独配置sku
	//Service
	skuService := ReminderRuleSku{e.Service}
	reminderListService := ReminderList{e.Service}
	reminderListSkuService := ReminderListSku{e.Service}

	for _, rule := range list {
		skuList := make([]models.ReminderRuleSku, 0)
		var count int64
		skuReq := dto.ReminderRuleSkuGetPageReq{
			Pagination: cDto.Pagination{
				PageIndex: -1,
				PageSize:  -1,
			},
			ReminderRuleId: int(rule.Id),
			Status:         "1",
		}

		err = skuService.GetPage(&skuReq, nil, &skuList, &count)
		if err != nil {
			return
		}

		// 根据仓库id查询sku库存
		physicalWarehouseCode := rule.WarehouseCode //physical_warehouse_id 逻辑仓库对应实体仓id

		// 根据实体仓id查出逻辑仓信息（多条 正品次品仓）
		//$sql = "select id,type,warehouse_code from pc.sxyz_logic_warehouse
		//where physical_warehouse_id=" . $value['warehouse_id'];
		physicalWarehouseList := []models.LogicWarehouse{}
		err = e.LogicWarehouseGetAll(&dto.LogicWarehouseGetPageReq{
			WarehouseCode: physicalWarehouseCode,
		}, p, &physicalWarehouseList)
		if err != nil {
			return
		}
		log.Info(physicalWarehouseList)
		var physicalCodez []string // 正品仓code
		var physicalCodec []string // 次品仓code
		for _, warehouse := range physicalWarehouseList {
			// 0 正品仓 1次品仓
			if warehouse.Type == "0" {
				physicalCodez = append(physicalCodez, warehouse.LogicWarehouseCode)
			} else {
				physicalCodec = append(physicalCodec, warehouse.LogicWarehouseCode)
			}
		}

		// 根据逻辑仓code查库存 查询正品仓sku库存是否到通用达预警值
		genuineArr := []dto.ReminderListGenuine{}
		query := e.Orm.Table(pcPrefix+".goods g").
			Select("stock_info.*, IFNULL(stock_info.stock, 0) AS stock, IFNULL(stock_info.lock_stock, 0) AS lock_stock,"+
				"g.supplier_sku_code vendor_sku, g.sku_code, g.id,g.vendor_id").
			Joins("LEFT JOIN stock_info ON g.sku_code = stock_info.sku_code AND g.id = stock_info.goods_id").
			Where("IFNULL(stock_info.logic_warehouse_code, 0) IN (?)", append(physicalCodez, "0")).
			Where("IFNULL(stock_info.stock, 0) <= ?", rule.WarningValue).
			Where("g.warehouse_code = ?", rule.WarehouseCode)

		if skuList != nil && len(skuList) > 0 {
			skuCodes := []string{}
			for _, sku := range skuList {
				skuCodes = append(skuCodes, sku.SkuCode)
			}
			query = query.Where("ifnull(stock_info.sku_code, '0') not IN (?)", skuCodes)
		}

		query.Find(&genuineArr)

		var insertSupplier map[int][]*dto.ReminderListSkuInsertReq
		insertSupplier = make(map[int][]*dto.ReminderListSkuInsertReq)

		goodsIds := []int{}
		for _, genuineData := range genuineArr {
			goodsIds = append(goodsIds, genuineData.Id)
		}
		goodsOrderLackStockMap := e.GetOrderLackStock(goodsIds)
		// 循环这些到达库存预警值的sku
		for _, genuineData := range genuineArr {
			// 次品数量
			var defectiveStock int = 0
			// 当前次品仓在库数量
			var defectiveAllStock int = 0

			// 根据skucode 查询次品仓数量
			if physicalCodec != nil && len(physicalCodec) > 0 {
				//$sql = "SELECT ifnull(sinfo.stock,0) stock
				//FROM pc.goods g
				//left join pc.sxyz_stock_info sinfo on g.sku_code = sinfo.sku_code and g.id = sinfo.goods_id
				//WHERE sinfo.logic_warehouse_code in('" . $physicalCodeCpStr . "')"
				//. " AND sinfo.sku_code = '" . $genuineData['sku_code']."'"
				//. " AND g.warehouse_code = '" . $warehouseCode."'";
				//if($skuList){
				//	$skuStr = implode("','",array_column($skuList, 'sku_code'));
				//	$sql .=. " AND sinfo.sku_code not in('" . $skuStr  ."')";
				//}
				//$query = $this->db->query($sql);

				var defectiveData int64
				query = e.Orm.Table(pcPrefix+".goods g").
					Select("IFNULL(stock_info.stock, 0) AS stock").
					Joins("LEFT JOIN stock_info ON g.sku_code = stock_info.sku_code AND g.id = stock_info.goods_id").
					Where("logic_warehouse_code IN (?)", physicalCodec).
					Where("stock_info.sku_code = ?", genuineData.SkuCode).
					Where("g.warehouse_code = ?", rule.WarehouseCode)

				if skuList != nil && len(skuList) > 0 {
					skuCodes := []string{}
					for _, sku := range skuList {
						skuCodes = append(skuCodes, sku.SkuCode)
					}
					query = query.Where("ifnull(stock_info.sku_code) IN (?)", skuCodes)
				}

				query.First(&defectiveData)

				// 当前次品仓在库数量 可用数量+锁定库存
				defectiveAllStock = defectiveStock
			}
			// 添加的时候根据货主分
			if _, ok := insertSupplier[genuineData.VendorId]; !ok {
				insertSupplier[genuineData.VendorId] = []*dto.ReminderListSkuInsertReq{}
			}
			orderLackStock := 0
			if v, ok := goodsOrderLackStockMap[genuineData.Id]; ok {
				orderLackStock = v
			}
			addSkuData := GetAddSkuData(&genuineData, defectiveAllStock, rule, orderLackStock)
			insertSupplier[genuineData.VendorId] = append(insertSupplier[genuineData.VendorId], &addSkuData)
			log.Info("genuineData.VendorId", genuineData.VendorId, addSkuData)

		}

		// 根据逻辑仓code查库存 查询正品仓单的设置了预警值的sku库存是否到达预警值
		for _, sku := range skuList {
			// 根据逻辑仓code查库存 查询正品仓单的设置了预警值的sku库存是否到达预警值
			var genuineData dto.ReminderListGenuine
			query = e.Orm.Table(pcPrefix+".goods g").
				Select("stock_info.*, g.supplier_sku_code vendor_sku, g.sku_code, g.id,g.vendor_id").
				Joins("LEFT JOIN stock_info ON g.sku_code = stock_info.sku_code AND g.id = stock_info.goods_id").
				Where("IFNULL(stock_info.logic_warehouse_code, 0) IN (?)", append(physicalCodez, "0")).
				Where("IFNULL(stock_info.stock, 0) <= ?", sku.WarningValue).
				Where("stock_info.sku_code = ?", sku.SkuCode).
				Where("g.warehouse_code = ?", rule.WarehouseCode).
				First(&genuineData)

			// 如果存在就说明到预警值
			if genuineData.Id > 0 {
				// 次品数量
				var defectiveStock int = 0
				// 当前次品仓在库数量
				var defectiveAllStock int = 0

				if physicalCodec != nil && len(physicalCodec) > 0 {
					// 根据skucode 查询次品仓数量
					//$sql = "SELECT sinfo.stock FROM pc.sxyz_stock_info  sinfo
					//join pc.goods g on g.sku_code = sinfo.sku_code and g.id = sinfo.goods_id
					//WHERE sinfo.logic_warehouse_code in('" . $physicalCodeCpStr . "')"
					//. " AND sinfo.sku_code = '" . $skuV['sku_code']."'"
					//. " AND g.warehouse_code = '" . $warehouseCode."'";
					//$query = $this->db->query($sql);
					//$defectiveData = $query->rows[0];

					// 这里需要一个库存查询的方法
					e.Orm.Table("stock_info").
						Select("stock_info.stock").
						Joins("JOIN "+pcPrefix+".goods as g ON g.sku_code = stock_info.sku_code AND g.id = stock_info.goods_id").
						Where("logic_warehouse_code IN (?)", physicalCodec).
						Where("stock_info.sku_code = ?", sku.SkuCode).
						Where("g.warehouse_code = ?", rule.WarehouseCode).
						First(&defectiveStock)

					// 使用 stock 进行后续操作

					// 当前次品仓在库数量 可用数量+锁定库存
					defectiveAllStock = defectiveStock

				}

				// 添加的时候根据货主分
				if _, ok := insertSupplier[genuineData.VendorId]; !ok {
					insertSupplier[genuineData.VendorId] = []*dto.ReminderListSkuInsertReq{}
				}
				goodsIds = []int{}
				goodsIds = append(goodsIds, genuineData.Id)
				goodsOrderLackStockMap = e.GetOrderLackStock(goodsIds)
				orderLackStock := 0
				if v, ok := goodsOrderLackStockMap[genuineData.Id]; ok {
					orderLackStock = v
				}
				skuStatus, _ := strconv.Atoi(sku.Status)
				addSkuData := GetAddSkuData(&genuineData, defectiveAllStock, models.ReminderRule{
					Model:              cModel.Model{Id: sku.Id},
					CompanyId:          rule.CompanyId,
					WarehouseCode:      rule.WarehouseCode,
					WarningValue:       sku.WarningValue,
					ReplenishmentValue: sku.ReplenishmentValue,
					Status:             skuStatus,
				}, orderLackStock)
				insertSupplier[genuineData.VendorId] = append(insertSupplier[genuineData.VendorId], &addSkuData)
				log.Info("genuineData.VendorId", genuineData.VendorId, addSkuData)

			}

		}
		for VendorId, skus := range insertSupplier {
			log.Info("insertSupplier-VendorId", VendorId)
			insertReminderList := &dto.ReminderListInsertReq{
				Id:             0,
				ReminderRuleId: int(rule.Id),
				CompanyId:      rule.CompanyId,
				WarehouseCode:  rule.WarehouseCode,
				VendorId:       VendorId,
				SkuCount:       int(len(skus)),
			}
			// 添加补货清单
			err = reminderListService.Insert(insertReminderList)
			if err != nil {

			}

			// 添加补货清单sku记录
			if err == nil {
				err = reminderListSkuService.Inserts(skus, insertReminderList.Id)
				if err != nil {
					return
				}
			}
		}

	}
	return
}

func GetAddSkuData(genuineData *dto.ReminderListGenuine, defectiveAllStock int, rule models.ReminderRule, orderLackStock int) dto.ReminderListSkuInsertReq {
	var genuineStock, genuineAllStock int
	// 当前正品可用数量
	genuineStock = genuineData.Stock
	// 当前正品仓在库数量 可用数量+锁定库存
	genuineAllStock = genuineStock + int(genuineData.LockStock)

	// 订单实际缺货数量：补货SKU 在缺货订单中实际缺货数量之和；
	orderLackStock = orderLackStock

	AddSkuData := dto.ReminderListSkuInsertReq{
		WarningValue: rule.WarningValue,
		// 补货量
		ReplenishmentValue: rule.ReplenishmentValue,
		// 建议补货量：建议补货量=备货量-当前正品用数量+订单实际缺货数量；
		RecommendReplenishmentValue: rule.ReplenishmentValue - int(genuineStock) + int(orderLackStock),
		SkuCode:                     genuineData.SkuCode,
		VendorId:                    genuineData.VendorId,
		VendorSku:                   genuineData.VendorSku,
		// 当前正品可用数量：SKU在当前实体仓的正品仓下的可用数量
		GenuineStock: genuineStock,
		// 当前在库数量：SKU在当前实体仓的正品仓在库数量+次品仓在库数量
		AllStock: genuineAllStock + defectiveAllStock,
		// 当前占用数量：SKU在当前实体仓的正品仓占用数量
		OccupyStock: int(genuineData.LockStock),
		// 订单缺货数量
		OrderLackStock: orderLackStock,
		CreatedAt:      time.Now(),
	}
	return AddSkuData
}

func (e *ReminderList) GetOrderLackStock(ids []int) map[int]int {
	orderDetailM := modelsOc.OrderDetail{}
	return orderDetailM.GetOrderLackStock(e.Orm, ids)

}
