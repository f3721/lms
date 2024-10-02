package admin

import (
	"errors"
	"fmt"
	"strings"

	modelsPc "go-admin/app/pc/models"
	modelsUc "go-admin/app/uc/models"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"github.com/samber/lo"
	"gorm.io/gorm"

	"go-admin/app/wc/models"
	"go-admin/app/wc/service/admin/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
	"go-admin/common/global"
	"go-admin/common/utils"
)

type QualityCheckConfig struct {
	service.Service
}

// GetPage 获取QualityCheckConfig列表
func (e *QualityCheckConfig) GetPage(c *dto.QualityCheckConfigGetPageReq, p *actions.DataPermission, list *[]*dto.QualityCheckConfigGetPageResp, count *int64) error {

	// 搜索条件
	searchScope := func(db *gorm.DB) *gorm.DB {
		if c.CompanyIds != "" {
			db.Where("qd.company_id in ?", utils.Split(c.CompanyIds))
		}
		if c.WarehouseCodes != "" {
			db.Where("qd.warehouse_code in ?", utils.Split(c.WarehouseCodes))
		}
		if c.Type != "" {
			db.Where("quality_check_config.type = ?", c.Type)
		}
		return db
	}

	// DB查询
	var data models.QualityCheckConfig
	err := e.Orm.Model(&data).
		Scopes(
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
			searchScope,
		).
		Joins("join quality_check_config_detail qd on quality_check_config.id = qd.config_id").
		Preload("QualityCheckConfigDetail").
		Group("quality_check_config.id").
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("QualityCheckConfigService GetPage error:%s \r\n", err)
		return err
	}

	// 返回结果处理
	for _, item := range *list {
		// 去重组合公司名称和仓库名称
		companyNameMap := make(map[string]bool)
		warehouseNameMap := make(map[string]bool)
		for _, v := range item.QualityCheckConfigDetail {
			companyNameMap[v.CompanyName] = true
			warehouseNameMap[v.WarehouseName] = true
		}
		item.CompanyNames = strings.Join(lo.Keys[string](companyNameMap), ",")
		item.WarehouseNames = strings.Join(lo.Keys[string](warehouseNameMap), ",")

		// 质检类型中文
		typeMap := map[int]string{
			0: "全检",
			1: "抽检",
		}
		item.TypeName = typeMap[item.Type]
	}

	return nil
}

// Get 获取QualityCheckConfig对象
func (e *QualityCheckConfig) Get(d *dto.QualityCheckConfigGetReq, p *actions.DataPermission, model *models.QualityCheckConfig) error {
	var data models.QualityCheckConfig

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetQualityCheckConfig error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建QualityCheckConfig对象
func (e *QualityCheckConfig) Insert(c *dto.QualityCheckConfigInsertReq) error {
	// 查询仓库列表
	warehouseModel := &models.Warehouse{}
	warehouseList := &[]models.Warehouse{}
	err := e.Orm.Model(warehouseModel).
		Where("company_id in (?)", c.CompanyIds).
		Where("warehouse_code in (?)", c.WarehouseCodes).
		Find(&warehouseList).Error
	if err != nil {
		return err
	}

	// 查询公司Map
	companyModel := &modelsUc.CompanyInfo{}
	companyMap, err := companyModel.MapByIds(e.Orm, c.CompanyIds)
	if err != nil {
		return err
	}

	// 组合子表数据
	details := []*models.QualityCheckConfigDetail{}
	for _, item := range *warehouseList {
		detail := &models.QualityCheckConfigDetail{
			CompanyId:     item.CompanyId,
			CompanyName:   companyMap[item.CompanyId].CompanyName,
			WarehouseCode: item.WarehouseCode,
			WarehouseName: item.WarehouseName,
			OrderType:     c.OrderType,
		}
		if c.OrderType == 0 { // 全部类型
			detail.OrderType = 1 // 采购类型
			details = append(details, detail)

			detail2 := &models.QualityCheckConfigDetail{}
			*detail2 = *detail
			detail2.OrderType = 2 // 大货类型
			details = append(details, detail2)
		} else {
			detail.OrderType = c.OrderType
			details = append(details, detail)
		}
	}
	c.QualityCheckConfigDetail = details

	// 参数校验
	err = e.verification(c, warehouseList, companyMap)
	if err != nil {
		return err
	}

	// 新建数据
	var data models.QualityCheckConfig
	c.Generate(&data)
	err = e.Orm.Save(&data).Error
	if err != nil {
		e.Log.Errorf("QualityCheckConfigService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改QualityCheckConfig对象
func (e *QualityCheckConfig) Update(c *dto.QualityCheckConfigInsertReq, p *actions.DataPermission) error {

	// 查询数据
	var data = models.QualityCheckConfig{}
	err := e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId()).Error
	if err != nil {
		return err
	}

	// 查询仓库列表
	warehouseModel := &models.Warehouse{}
	warehouseList := &[]models.Warehouse{}
	err = e.Orm.Model(warehouseModel).
		Where("company_id in (?)", c.CompanyIds).
		Where("warehouse_code in (?)", c.WarehouseCodes).
		Find(&warehouseList).Error
	if err != nil {
		return err
	}

	// 查询公司Map
	companyModel := &modelsUc.CompanyInfo{}
	companyMap, err := companyModel.MapByIds(e.Orm, c.CompanyIds)
	if err != nil {
		return err
	}

	// 组合子表数据
	details := []*models.QualityCheckConfigDetail{}
	for _, item := range *warehouseList {
		detail := &models.QualityCheckConfigDetail{
			CompanyId:     item.CompanyId,
			CompanyName:   companyMap[item.CompanyId].CompanyName,
			WarehouseCode: item.WarehouseCode,
			WarehouseName: item.WarehouseName,
			OrderType:     c.OrderType,
		}
		if c.OrderType == 0 { // 全部类型
			detail.OrderType = 1 // 采购类型
			details = append(details, detail)

			detail2 := &models.QualityCheckConfigDetail{}
			*detail2 = *detail
			detail2.OrderType = 2 // 大货类型
			details = append(details, detail2)
		} else {
			detail.OrderType = c.OrderType
			details = append(details, detail)
		}
	}
	c.QualityCheckConfigDetail = details

	// 参数校验
	err = e.verification(c, warehouseList, companyMap)
	if err != nil {
		return err
	}

	// 更新处理
	c.Generate(&data)
	err = e.Orm.Transaction(func(tx *gorm.DB) error {
		// 清除旧子数据
		err := tx.Where("config_id = ?", c.Id).Unscoped().Delete(&models.QualityCheckConfigDetail{}).Error
		if err != nil {
			return err
		}

		// 更新数据
		err = tx.Omit("create_by", "create_by_name").Save(&data).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// Remove 删除QualityCheckConfig
func (e *QualityCheckConfig) Remove(d *dto.QualityCheckConfigDeleteReq, p *actions.DataPermission) error {
	var data models.QualityCheckConfig

	// 校验数据
	err := e.Orm.Preload("QualityCheckConfigDetail").First(&data, d.GetId()).Error
	if err != nil {
		return errors.New("数据不存在,请检查！")
	}

	// 联动删除
	err = e.Orm.Transaction(func(tx *gorm.DB) error {
		// 删除子数据
		err := tx.Where("config_id = ?", d.GetId()).Unscoped().Delete(&models.QualityCheckConfigDetail{}).Error
		if err != nil {
			return err
		}

		// 删除数据
		err = tx.Scopes(actions.Permission(data.TableName(), p)).Delete(&data).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		e.Log.Errorf("Service RemoveQualityCheckConfig error:%s \r\n", err)
		return err
	}

	return nil
}

// 初始化接口
func (e *QualityCheckConfig) Init() (any, error) {

	// 字段选项
	skuFieldsOptionsMap := []map[string]any{}
	categoryFieldsOptionsMap := []map[string]any{}
	valuesOptionsMap := map[string]any{}
	suggestOperator := map[string]string{}
	for key, constraint := range models.ConstraintMap {
		// 排序
		orderArr := make([]string, len(constraint))
		for filed, item := range constraint {
			orderArr[item["order"].(int)-1] = filed

			// 建议运算符
			suggestOperator[filed] = item["operator"].(string)
		}

		// 取字段选项
		for _, filed := range orderArr {
			// 字符选项
			itemOptions := map[string]any{"label": constraint[filed]["name"], "value": filed}
			if key == "skuConstraint" {
				skuFieldsOptionsMap = append(skuFieldsOptionsMap, itemOptions)
			} else {
				categoryFieldsOptionsMap = append(categoryFieldsOptionsMap, itemOptions)
			}

			// 字段对应的值
			targetValue := constraint[filed]["targetValue"].([]string)
			if len(targetValue) == 0 {
				continue
			}
			valueArr := []map[string]string{}
			for _, v := range targetValue {
				itemKvMap := map[string]string{"label": v, "value": v}
				valueArr = append(valueArr, itemKvMap)
			}
			valuesOptionsMap[filed] = valueArr
		}
	}

	// 初始化数据
	res := map[string]any{
		"typeOption": []map[string]string{
			{
				"label": "全部",
				"value": "0",
			},
			{
				"label": "采购入库",
				"value": "1",
			},
			{
				"label": "大货入库",
				"value": "2",
			},
		},
		"relationOption": []map[string]string{
			{
				"label": "并且",
				"value": "and",
			},
			{
				"label": "或者",
				"value": "or",
			},
		},
		"operatorOption": []map[string]string{
			{
				"label": "包含",
				"value": "like",
			},
			{
				"label": "等于",
				"value": "=",
			},
			{
				"label": "其中之一",
				"value": "in",
			},
		},
		"skuFieldsOptions":      skuFieldsOptionsMap,
		"categoryFieldsOptions": categoryFieldsOptionsMap,
		"valuesOptionsMap":      valuesOptionsMap,
		"suggestOperator":       suggestOperator,
	}

	return res, nil
}

// QualityCheckNum 查询需要质检的SKU
func (e *QualityCheckConfig) QualityCheckNum(req *dto.QualityCheckNumReq) ([]*dto.QualityCheckNumProducts, error) {
	//return req.QualityCheckNumProducts, nil
	res := []*dto.QualityCheckNumProducts{}

	// 0. 参数校验
	if lo.IndexOf([]string{"0", "3"}, req.OrderType) == -1 {
		return nil, errors.New("传参错误: OrderType范围[0,3]")
	}

	// 1. 定位配置项
	detail := &models.QualityCheckConfigDetail{}
	orderTypeMap := map[string]string{
		"3": "1", // 采购入库
		"0": "2", // 大货入库
	}
	err := e.Orm.Where("order_type = ?", orderTypeMap[req.OrderType]).Where("warehouse_code = ?", req.WarehouseCode).First(detail).Error
	if err != nil {
		return res, nil
	}
	data := &models.QualityCheckConfig{}
	err = e.Orm.Where("id = ?", detail.ConfigId).First(data).Error
	if err != nil {
		return res, nil
	}

	// 2. 汇总SKU
	skuArr := []string{}
	skuMap := map[string]int{}
	for _, item := range req.QualityCheckNumProducts {
		skuArr = append(skuArr, item.SkuCode)
		skuMap[item.SkuCode] = item.Quantity
	}

	// 3. 组合SQL
	whereArr := []string{}
	orArr := []string{}
	for _, item := range data.SkuConstraint {

		// 处理值
		finelVal := e.SqlVal(item.Operator, item.TargetValue)
		sql := fmt.Sprintf("p.%v %v %v", item.Field, item.Operator, finelVal)

		// 分流
		if item.Relation == "and" {
			whereArr = append(whereArr, sql)
		} else {
			orArr = append(orArr, sql)
		}
	}
	for _, item := range data.CategoryConstraint {

		// 处理值
		finelVal := e.SqlVal(item.Operator, item.TargetValue)
		sql := fmt.Sprintf("c.%v %v %v", item.Field, item.Operator, finelVal)

		// 分流
		if item.Relation == "and" {
			whereArr = append(whereArr, sql)
		} else {
			orArr = append(orArr, sql)
		}
	}

	// 4.1 查询符合条件的SKU
	products := []modelsPc.Product{}
	producCategoryTable := global.TenantTableName(e.Orm, "pc", "product_category", "pc")
	categoryTable := global.TenantTableName(e.Orm, "pc", "category", "c")
	queryDb := e.Orm.Scopes(global.TenantTable("pc", "product", "p")).
		Joins("join "+producCategoryTable+" on p.sku_code = pc.sku_code and main_cate_flag = 1 and pc.deleted_at IS NULL").
		Joins("join "+categoryTable+" on pc.category_id = c.id and c.deleted_at IS NULL").
		Where("p.sku_code in (?)", skuArr)

	// 4.1 拼接自定义条件
	if len(whereArr) > 0 {
		queryDb.Where(strings.Join(whereArr, " AND "))
	}
	if len(orArr) > 0 {
		queryDb.Where(strings.Join(orArr, " OR "))
	}

	// 4.2 最终查询
	err = queryDb.Find(&products).Group("p.sku_code").Error
	if err != nil {
		return nil, err
	}

	// 5. 过滤SKU&指定质检数量
	if len(products) == 0 {
		return res, err
	}
	for _, item := range products {
		itemRes := &dto.QualityCheckNumProducts{}
		if data.Type == 0 || skuMap[item.SkuCode] < data.SamplingNum { // 全检或入库数量小于抽奖数量
			itemRes.SkuCode = item.SkuCode
			itemRes.Quantity = skuMap[item.SkuCode]
			itemRes.Type = data.Type
			itemRes.Unqualified = data.Unqualified
		} else { // 抽检
			itemRes.SkuCode = item.SkuCode
			itemRes.Quantity = data.SamplingNum
			itemRes.Type = data.Type
			itemRes.Unqualified = data.Unqualified
		}
		res = append(res, itemRes)
	}
	return res, nil
}

// -------------------------私有函数开始-------------------------

// verification 新增和更新入参校验
func (e *QualityCheckConfig) verification(c *dto.QualityCheckConfigInsertReq, warehouseList *[]models.Warehouse, companyMap map[int]*modelsUc.CompanyInfo) error {

	// 校验约束条件Key
	allConstraint := []*models.Constraint{}
	allConstraint = append(allConstraint, c.SkuConstraint...)
	allConstraint = append(allConstraint, c.CategoryConstraint...)
	SkuConstraintKeys := lo.Map(allConstraint, func(item *models.Constraint, _ int) string {
		return item.Field
	})
	skuRepeat := lo.FindDuplicates(SkuConstraintKeys)
	if len(skuRepeat) > 0 {
		return fmt.Errorf("传参错误：字段[%v]重复出现,请检查", skuRepeat[0])
	}

	// 校验约束条件详情
	for _, item := range c.SkuConstraint {
		files := []string{}
		for key := range models.ConstraintMap["skuConstraint"] {
			files = append(files, key)
		}

		if lo.IndexOf[string](files, item.Field) == -1 {
			return fmt.Errorf("传参错误：字段[%v]不符合预期,请检查", item.Field)
		}

		values := models.ConstraintMap["skuConstraint"][item.Field]["targetValue"].([]string)
		if len(values) > 0 {
			// 包含
			if item.Operator == "in" {
				targetValues := strings.Split(item.TargetValue, ",")
				diff, _ := lo.Difference(targetValues, values)
				if len(diff) > 0 {
					return fmt.Errorf("传参错误：字段[%v]值[%v]不符合预期,请检查", item.Field, diff[0])
				}
			}

			// 相等
			if item.Operator == "=" {
				if lo.IndexOf(values, item.TargetValue) == -1 {
					return fmt.Errorf("传参错误：字段[%v]值[%v]不符合预期,请检查", item.Field, item.TargetValue)
				}
			}
		}

	}
	for _, item := range c.CategoryConstraint {
		files := []string{}
		for key := range models.ConstraintMap["categoryConstraint"] {
			files = append(files, key)
		}
		if lo.IndexOf(files, item.Field) == -1 {
			return fmt.Errorf("传参错误：字段[%v]不符合预期,请检查", item.Field)
		}

		values := models.ConstraintMap["categoryConstraint"][item.Field]["targetValue"].([]string)
		if len(values) > 0 {
			// 包含
			if item.Operator == "in" {
				targetValues := strings.Split(item.TargetValue, ",")
				diff, _ := lo.Difference(targetValues, values)
				if len(diff) > 0 {
					return fmt.Errorf("传参错误：字段[%v]值[%v]不符合预期,请检查", item.Field, diff[0])
				}
			}

			// 相等
			if item.Operator == "=" {
				if lo.IndexOf(values, item.TargetValue) == -1 {
					return fmt.Errorf("传参错误：字段[%v]值[%v]不符合预期,请检查", item.Field, item.TargetValue)
				}
			}
		}

		// 验证通过后将产线搜索字段修正
		item.Field = "cate_level"
	}

	// 区分新增和编辑
	db := e.Orm
	if c.Id != 0 {
		db = e.Orm.Where("config_id <> ?", c.Id)
	}

	// 校验公司+仓库不能重复
	whereCondition1 := []string{}
	for _, item := range *warehouseList {
		whereCondition1 = append(whereCondition1, fmt.Sprintf("(%v,%s)", item.CompanyId, "'"+item.WarehouseCode+"'"))
	}
	check1 := &models.QualityCheckConfigDetail{}
	err := db.Where(fmt.Sprintf("(company_id, warehouse_code) in (%s)", strings.Join(whereCondition1, ","))).First(check1).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if check1.Id != 0 {
		return fmt.Errorf("公司[%v]+仓库[%v]不能重复", check1.CompanyName, check1.WarehouseName)
	}

	// 校验订单类型+仓库不能重复
	whereCondition2 := []string{}
	for _, item := range c.QualityCheckConfigDetail {
		whereCondition2 = append(whereCondition2, fmt.Sprintf("(%v,%s)", item.OrderType, "'"+item.WarehouseCode+"'"))
	}
	check2 := &models.QualityCheckConfigDetail{}
	err = db.Where(fmt.Sprintf("(order_type, warehouse_code) in (%s)", strings.Join(whereCondition2, ","))).First(check2).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if check2.Id != 0 {
		return fmt.Errorf("订单类型[%v]+仓库[%v]不能重复", check2.OrderType, check2.WarehouseName)
	}

	return nil
}

// 约束Sql值的处理
func (c *QualityCheckConfig) SqlVal(operator, targetValue string) string {
	if operator == "in" {
		arr := strings.Split(targetValue, ",")
		transformed := make([]string, len(arr))
		for i, item := range arr {
			// 如果有映射的值做处理
			mapVal, ok := models.ConstraintValMap[item]
			if ok {
				item = mapVal
			}
			transformed[i] = "'" + item + "'"
		}
		targetValue = "(" + strings.Join(transformed, ", ") + ")"
		return targetValue
	}

	if operator == "like" {
		targetValue = "'%" + targetValue + "%'"
	}

	// 如果有映射的值做处理
	mapVal, ok := models.ConstraintValMap[targetValue]
	if ok {
		targetValue = mapVal
	}
	return targetValue
}
