package admin

import (
	"encoding/json"
	"errors"
	"fmt"
	vd "github.com/bytedance/go-tagexpr/v2/validator"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"github.com/jinzhu/copier"
	"github.com/samber/lo"
	dtoUc "go-admin/app/uc/service/admin/dto"
	modelsWc "go-admin/app/wc/models"
	dtoWc "go-admin/app/wc/service/admin/dto"
	ucClient "go-admin/common/client/uc"
	wcClient "go-admin/common/client/wc"
	"go-admin/common/excel"
	"go-admin/common/global"
	"go-admin/common/utils"
	"mime/multipart"
	"regexp"
	"strings"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/pc/models"
	"go-admin/app/pc/service/admin/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
	cModel "go-admin/common/models"
)

type Goods struct {
	service.Service
}

var goodsStatusTxt = []string{"停用", "启用"}
var onlineStatusTxt = []string{"未上架", "上架", "下架"}

// GetPage 获取Goods列表
func (e *Goods) GetPage(c *dto.GoodsGetPageReq, p *actions.DataPermission, list *[]dto.GoodsGetPageResp, count *int64) error {
	var err error
	var data models.Goods

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
			dto.GoodsMakeCondition(c, p),
			actions.SysUserPermission(data.TableName(), p, 4), //货主权限
		).
		Joins("INNER JOIN product on goods.sku_code = product.sku_code").
		Preload("Product.Brand").
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("GoodsService GetPage error:%s \r\n", err)
		return err
	}
	if *count > 0 {
		var vendorIds []int
		for _, resp := range *list {
			vendorIds = append(vendorIds, resp.VendorId)
		}
		var warehouseCode []string
		for _, resp := range *list {
			warehouseCode = append(warehouseCode, resp.WarehouseCode)
		}
		if len(vendorIds) > 0 {
			// 货主map
			vendorResult := e.apiGetVendorInfoById(vendorIds)
			// 仓库map
			warehouseResult := e.apiGetWarehouseInfoByCode(warehouseCode)
			tmpList := *list
			for k, goods := range tmpList {
				if len(vendorResult) > 0 {
					tmpList[k].VendorName = vendorResult[goods.VendorId].NameZh
					tmpList[k].VendorShortName = vendorResult[goods.VendorId].ShortName
				}
				if len(warehouseResult) > 0 {
					tmpList[k].WarehouseName = warehouseResult[goods.WarehouseCode].WarehouseName
					tmpList[k].CompanyName = warehouseResult[goods.WarehouseCode].CompanyName
				}
			}
			list = &tmpList
		}
	}
	return nil
}

// Get 获取Goods对象
func (e *Goods) Get(d *dto.GoodsGetReq, p *actions.DataPermission, model *dto.GoodsGetResp) error {
	var data models.Goods

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		Preload("Product").
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetGoods error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}

	if model.WarehouseCode != "" {
		// 仓库名称
		warehouseMap := e.apiGetWarehouseInfoByCode([]string{model.WarehouseCode})
		// 仓库map
		warehouseResult := e.apiGetWarehouseInfoByCode([]string{model.WarehouseCode})
		if len(warehouseResult) > 0 {
			model.CompanyName = warehouseResult[model.WarehouseCode].CompanyName
		}
		if value, ok := warehouseMap[model.WarehouseCode]; ok {
			model.WarehouseName = value.WarehouseName
		}
	}

	return nil
}

// Insert 创建Goods对象
func (e *Goods) Insert(c *dto.GoodsInsertReq, p *actions.DataPermission) error {
	var err error
	var data models.Goods
	c.Generate(&data)

	// 获取用户权限仓库
	authWarehouseCodes := utils.Split(p.AuthorityWarehouseId)
	// api获取仓库信息
	warehouseMap := e.apiGetWarehouseInfoByCode([]string{data.WarehouseCode})
	// api获取货主信息
	vendorMap := e.apiGetVendorInfoById([]int{data.VendorId})

	status, msg := e._validate(c, warehouseMap, vendorMap, authWarehouseCodes)
	if status {
		err = e.Orm.Create(&data).Error
		if err != nil {
			e.Log.Errorf("GoodsService Insert error:%s \r\n", err)
			return err
		}

		// 生成日志
		dataLog, _ := json.Marshal(&c)
		goodsLog := models.GoodsLog{
			DataId:     data.Id,
			Type:       global.LogTypeCreate,
			Data:       string(dataLog),
			BeforeData: "",
			AfterData:  string(dataLog),
			ControlBy:  c.ControlBy,
		}
		_ = goodsLog.CreateLog("goods", e.Orm)

		return nil
	}

	return errors.New(strings.Join(msg, ","))
}

// Update 修改Goods对象
func (e *Goods) Update(c *dto.GoodsUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.Goods{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())

	before := data
	c.Generate(&data)
	// 如果停用
	if data.Status == 0 && before.Status == 1 {
		// 如果是停用
		if before.OnlineStatus == 1 {
			return errors.New("请先下架该商品！")
		}
		// 校验 SKU+仓库+货主 库存为0时 可以停用，否则不可停用
		goodsStockMap := e.getStockMap([]int{c.Id}, c.WarehouseCode)
		if goodsStockMap[before.Id] > 0 {
			return errors.New("SKU+仓库+货主 库存 > 0,不可停用！")
		}
	} else {
		c.Status = 1
		var goods models.Goods
		findSameWhere := dto.FindSameWhere{
			SkuCode:       c.SkuCode,
			WarehouseCode: c.WarehouseCode,
			VendorId:      c.VendorId,
			Status:        c.Status,
			ApproveStatus: c.ApproveStatus,
		}
		if before.Status == 0 && e.ExistsSkuSupplierWarehouseUniquee(&findSameWhere, &goods) {
			return errors.New("SKU+货主+仓库不唯一，同一个仓库下存在相同SKU！")
		}
		findSkuSames := dto.FindSkuSames{
			SkuCode:       c.SkuCode,
			WarehouseCode: c.WarehouseCode,
		}
		if before.Status == 0 && e.ExistsSkuWarehouseUniquee(&findSkuSames) {
			return errors.New("SKU+仓库不唯一，同一个仓库下存在可用的相同SKU！")
		}
		// 如果变更物料编码或销售价或启用进入待审核状态
		e.doingApprove(&before, &data)
	}

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("GoodsService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}

	// 生成日志
	dataLog, _ := json.Marshal(&c)
	beforeDataStr, _ := json.Marshal(&before)
	afterDataStr, _ := json.Marshal(&data)
	goodsLog := models.GoodsLog{
		DataId:     data.Id,
		Type:       global.LogTypeUpdate,
		Data:       string(dataLog),
		BeforeData: string(beforeDataStr),
		AfterData:  string(afterDataStr),
		ControlBy:  c.ControlBy,
	}
	_ = goodsLog.CreateLog("goods", e.Orm)
	return nil
}

// Remove 删除Goods
func (e *Goods) Remove(d *dto.GoodsDeleteReq, p *actions.DataPermission) error {
	var data models.Goods

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveGoods error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}

// Approve 审核
func (e *Goods) Approve(c *dto.GoodsApproveReq, p *actions.DataPermission) (string, error) {
	if c.ApproveStatus == 0 {
		msg, err := e.OnlineOffline(&dto.OlineOfflineReq{
			Ids:        c.Ids,
			ActionType: 1,
			ControlBy:  c.ControlBy,
		}, p)
		return msg, err
	} else if c.ApproveStatus == 1 {
		data := map[string]interface{}{
			"approve_status": c.ApproveStatus,
			"approve_remark": "",
		}
		err := e._update(&dto.GoodsUpdater{ControlBy: c.ControlBy}, c.Ids, data)
		if err != nil {
			return "", err
		}
	} else {
		if c.ApproveRemark == "" {
			return "", errors.New("请填写审核不通过原因！")
		}
		data := map[string]interface{}{
			"approve_status": c.ApproveStatus,
			"approve_remark": c.ApproveRemark,
		}
		err := e._update(&dto.GoodsUpdater{ControlBy: c.ControlBy}, c.Ids, data)
		if err != nil {
			return "", err
		}
	}
	return "", nil
}

// OnlineOffline 批量上下架
func (e *Goods) OnlineOffline(c *dto.OlineOfflineReq, p *actions.DataPermission) (string, error) {
	if len(c.Ids) <= 0 {
		return "", errors.New("请先选择数据！")
	}
	data := make([]models.Goods, 0)
	err := e.Orm.Find(&data, c.Ids).Error
	if err != nil {
		return "", err
	}
	todoOnline := e.filterData(data, []interface{}{0, 2}) //0 未上架 2 已下架
	todoOffline := e.filterData(data, []interface{}{1})   // 已上架
	todoStatusDisable := e.filterDisableData(data)
	// 判断是上架还是下架
	if c.ActionType == 1 {
		tips := []string{
			" 上架成功", " 已上架不可重复上架！", " 状态停用不可以上架！", " 货主无效", " 仓库无效",
		}
		return e.actionOnlineOffline(c, todoOnline, todoOffline, todoStatusDisable, c.ActionType, tips), nil
	} else {
		tips := []string{
			" 下架成功", " 已下架或未上架不可操作下架！", " 状态停用不可以下架！",
		}
		return e.actionOnlineOffline(c, todoOffline, todoOnline, todoStatusDisable, c.ActionType, tips), nil
	}
}

func (e *Goods) actionOnlineOffline(c *dto.OlineOfflineReq, todoData []models.Goods, noHandle []models.Goods, disabledData []models.Goods, actionType int, msgSuffix []string) string {
	var green []models.Goods
	var msg []string
	// 可以上下架
	if len(todoData) > 0 {
		// 供应商失效
		var disabledVendor []models.Goods
		// 仓库失效
		var disabledWarehouse []models.Goods
		for _, goods := range todoData {
			if actionType == 1 && !e.supplierHasDisabled(goods.VendorId) {
				disabledVendor = append(disabledVendor, goods)
			} else if actionType == 1 && !e.warehouseHasDisabled(goods.WarehouseCode) {
				disabledWarehouse = append(disabledWarehouse, goods)
			} else {
				green = append(green, goods)
				// 更新
				data := map[string]interface{}{
					"online_status":  actionType,
					"approve_status": 1,
				}
				e._update(&dto.GoodsUpdater{ControlBy: c.ControlBy}, []int{goods.Id}, data)
			}
		}
		if len(disabledVendor) > 0 {
			txt := strings.Join(e.splitJoin(disabledVendor), "，") + msgSuffix[3]
			msg = append(msg, e.setMsgColor(txt, "#909399"))
		}
		if len(disabledWarehouse) > 0 {
			txt := strings.Join(e.splitJoin(disabledWarehouse), "，") + msgSuffix[4]
			msg = append(msg, e.setMsgColor(txt, "#F56C6C"))
		}
		if len(green) > 0 {
			txt := strings.Join(e.splitJoin(green), "，") + msgSuffix[0]
			msg = append(msg, e.setMsgColor(txt, "green"))
		}
	}
	// 不满足上下架条件
	if len(noHandle) > 0 {
		txt := strings.Join(e.splitJoin(noHandle), "，") + msgSuffix[1]
		msg = append(msg, e.setMsgColor(txt, "green"))
		if len(msg) == 1 {
			return msg[0]
		}
	}

	// 停用状态不可以上下架
	if len(disabledData) > 0 {
		txt := strings.Join(e.splitJoin(disabledData), "，") + msgSuffix[2]
		msg = append(msg, e.setMsgColor(txt, "red"))
		if len(msg) == 1 {
			return msg[0]
		}
	}
	if len(msg) > 0 {
		return strings.Join(msg, "，")
	}
	return ""
}

// GoodsImport 批量新增维护导入
func (e *Goods) GoodsImport(file *multipart.FileHeader, p *actions.DataPermission) (err error, errTitleList []map[string]string, errData []map[string]interface{}) {
	excelApp := excel.NewExcel()
	/**
	templateFilePath := "static/exceltpl/stock_location_import.xlsx"
	fieldsCorrect := excelApp.ValidImportFieldsCorrect(file, templateFilePath)
	if fieldsCorrect == false {
		err = errors.New("导入的excel字段与模板字段不一致，请重试")
	}
	**/

	err, datas, titleList := excelApp.GetExcelData(file)
	if err != nil {
		return
	}
	if len(datas) == 0 {
		err = errors.New("导入数据不能为空！")
		return
	}
	if len(datas) > 0 && len(datas) <= 20000 {
		errMsgList := make(map[int]string, 0)
		successCount := 0
		// 赋值到结构体
		importReq := dto.ImportReq{}
		for i, goods := range datas {
			errStr := ""
			var importGoods dto.ImportGoods
			importGoods.Trim(goods)
			// 基本数据校验
			err = vd.Validate(importGoods)
			if err != nil {
				errStr = err.Error()
			} else {
				importReq.Data = append(importReq.Data, importGoods)
			}
			if errStr == "" {
				successCount += 1
			}
			errMsgList[i] = errStr
		}
		if len(datas) > successCount {
			// 有校验错误的行时调用 导出保存excel
			err = errors.New("有错误！")
			errTitleList, errData = excelApp.MergeErrMsgColumn(titleList, datas, errMsgList)
		} else {
			var vendorName []string
			_ = utils.StructColumn(&vendorName, importReq.Data, "VendorName", "")
			var warehouseName []string
			_ = utils.StructColumn(&warehouseName, importReq.Data, "WarehouseName", "")
			var companyName []string
			_ = utils.StructColumn(&companyName, importReq.Data, "CompanyName", "")
			// 批量查询货主信息
			vendorList := e.apiGetVendorInfoByName(vendorName)
			warehouseList := e.apiGetWarehouseInfoByName(warehouseName)
			companyList := e.apiGetCompanyList(companyName)
			// 批量查询仓库信息
			successCount = 0
			for k, importGoods := range importReq.Data {
				var goods models.Goods
				errStr := ""
				// 数据校验
				err = e.importValidate(&importGoods, &goods, vendorList, warehouseList, companyList, p)
				if err != nil {
					errStr = err.Error()
				} else {
					if goods.Id > 0 {
						var before models.Goods
						_ = e.findGoodsById(goods.Id, &before)
						old := before
						// 如果变更物料编码或销售价或启用进入待审核状态
						e.doingApprove(&before, &goods)
						// 导入更新只允许更新销售价，价格调整备注，物料编码，停、启用
						before.MarketPrice = goods.MarketPrice
						before.PriceModifyReason = goods.PriceModifyReason
						before.ProductNo = goods.ProductNo
						before.Status = goods.Status
						err = e.Orm.Save(&before).Error
						if err != nil {
							e.Log.Errorf("GoodsService GoodsImport Save error:%s \r\n", err)
							errStr = err.Error()
						} else {
							// 生成日志----------------------------
							dataLog, _ := json.Marshal(&goods)
							beforeDataStr, _ := json.Marshal(&old)
							afterDataStr, _ := json.Marshal(&before)
							goodsLog := models.GoodsLog{
								DataId:     goods.Id,
								Type:       global.LogTypeUpdate,
								Data:       string(dataLog),
								BeforeData: string(beforeDataStr),
								AfterData:  string(afterDataStr),
								ControlBy: cModel.ControlBy{
									CreateBy:     user.GetUserId(e.Orm.Statement.Context.(*gin.Context)),
									CreateByName: user.GetUserName(e.Orm.Statement.Context.(*gin.Context)),
								},
							}
							_ = goodsLog.CreateLog("goods", e.Orm)
						}

					} else {
						err = e.Orm.Create(&goods).Error
						if err != nil {
							e.Log.Errorf("GoodsService GoodsImport Create error:%s \r\n", err)
							errStr = err.Error()
						} else {
							// 生成日志----------------------------
							dataLog, _ := json.Marshal(&goods)
							goodsLog := models.GoodsLog{
								DataId:     goods.Id,
								Type:       global.LogTypeCreate,
								Data:       string(dataLog),
								BeforeData: "",
								AfterData:  string(dataLog),
								ControlBy: cModel.ControlBy{
									CreateBy:     user.GetUserId(e.Orm.Statement.Context.(*gin.Context)),
									CreateByName: user.GetUserName(e.Orm.Statement.Context.(*gin.Context)),
								},
							}
							_ = goodsLog.CreateLog("goods", e.Orm)
						}
					}
				}
				if errStr == "" {
					successCount += 1
				}
				errMsgList[k] = errStr
			}
			if len(importReq.Data) > successCount {
				err = errors.New("有错误！")
				errTitleList, errData = excelApp.MergeErrMsgColumn(titleList, datas, errMsgList)
			}
		}
		return
	} else {
		err = errors.New("导入数据不能超过2万行！")
		return
	}
}

func (e *Goods) Export(c *dto.GoodsGetPageReq, p *actions.DataPermission) ([]interface{}, error) {
	var err error
	var data models.Goods
	var count int64
	var list []dto.GoodsGetPageResp
	db := e.Orm.Model(&data).
		Joins("Product").
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			dto.GoodsMakeCondition(c, p),
			actions.Permission(data.TableName(), p),
			actions.SysUserPermission(data.TableName(), p, 2), //仓库权限
			actions.SysUserPermission(data.TableName(), p, 4), //货主权限
		).
		Preload("Product.Brand").
		Find(&list)
	err = db.Count(&count).Error

	if err != nil {
		e.Log.Errorf("GoodsService GetPage error:%s \r\n", err)
		return nil, err
	}
	if count > 0 {
		var vendorIds []int
		for _, resp := range list {
			vendorIds = append(vendorIds, resp.VendorId)
		}
		var warehouseCode []string
		for _, resp := range list {
			warehouseCode = append(warehouseCode, resp.WarehouseCode)
		}
		if len(vendorIds) > 0 {
			// 货主map
			vendorResult := e.apiGetVendorInfoById(vendorIds)
			// 仓库map
			warehouseResult := e.apiGetWarehouseInfoByCode(warehouseCode)
			tmpList := list
			for k, goods := range tmpList {
				if len(vendorResult) > 0 {
					tmpList[k].VendorName = vendorResult[goods.VendorId].NameZh
					tmpList[k].VendorShortName = vendorResult[goods.VendorId].ShortName
				}
				if len(warehouseResult) > 0 {
					tmpList[k].WarehouseName = warehouseResult[goods.WarehouseCode].WarehouseName
					tmpList[k].CompanyName = warehouseResult[goods.WarehouseCode].CompanyName
				}
			}
			list = tmpList
		}

		exportData := make([]interface{}, 0)
		for _, goodsInfo := range list {
			goodsExportResp := dto.GoodsExportResp{
				SkuCode:           goodsInfo.SkuCode,
				GoodsName:         goodsInfo.Product.NameZh,
				BrandZh:           goodsInfo.Product.Brand.BrandZh,
				MfgModel:          goodsInfo.Product.MfgModel,
				CompanyName:       goodsInfo.CompanyName,
				WarehouseName:     goodsInfo.WarehouseName,
				VendorName:        goodsInfo.VendorName,
				SupplierSkuCode:   goodsInfo.Product.SupplierSkuCode,
				MarketPrice:       goodsInfo.MarketPrice,
				PriceModifyReason: goodsInfo.PriceModifyReason,
				ProductNo:         goodsInfo.ProductNo,
				OnlineStatus:      onlineStatusTxt[goodsInfo.OnlineStatus],
				Status:            goodsStatusTxt[goodsInfo.Status],
			}
			exportData = append(exportData, goodsExportResp)
		}
		return exportData, nil
	}
	return []interface{}{}, nil
}

func (e *Goods) _update(c *dto.GoodsUpdater, ids []int, data map[string]interface{}) error {
	// 日志记录
	var err error
	for _, id := range ids {
		var before models.Goods
		_ = e.findGoodsById(id, &before)
		if before.Id > 0 {
			beforeDataStr, _ := json.Marshal(&before)
			for k, v := range data {
				switch k {
				case "approve_status":
					before.ApproveStatus = v.(int)
				case "approve_remark":
					before.ApproveRemark = v.(string)
				case "online_status":
					before.OnlineStatus = v.(int)
				}
			}
			err = e.Orm.Save(&before).Error
			if err != nil {
				break
			}
			dataLog, _ := json.Marshal(&data)

			afterDataStr, _ := json.Marshal(&before)
			goodsLog := models.GoodsLog{
				DataId:     id,
				Type:       global.LogTypeUpdate,
				Data:       string(dataLog),
				BeforeData: string(beforeDataStr),
				AfterData:  string(afterDataStr),
				ControlBy: cModel.ControlBy{
					CreateBy:     c.UpdateBy,
					CreateByName: c.UpdateByName,
				},
			}
			_ = goodsLog.CreateLog("goods", e.Orm)
		}
	}

	return err
}

func (e *Goods) _validate(c *dto.GoodsInsertReq, warehouseMap map[string]dtoWc.WarehouseGetPageResp, vendorMap map[int]modelsWc.Vendors, authWarehouse []string) (bool, []string) {
	product := Product{e.Service}
	var msg []string
	status := true
	if product.IsApprove(c.SkuCode) {
		status = false
		msg = append(msg, "商品图文审核状态未审核或未通过，创建失败！")
	}
	var goods models.Goods
	findSameWhere := dto.FindSameWhere{
		SkuCode:       c.SkuCode,
		WarehouseCode: c.WarehouseCode,
		VendorId:      c.VendorId,
		Status:        c.Status,
		ApproveStatus: c.ApproveStatus,
	}
	if e.ExistsSkuSupplierWarehouseUniquee(&findSameWhere, &goods) {
		status = false
		msg = append(msg, "SKU+货主+仓库不唯一，同一个仓库下存在相同SKU！")
	}
	findSkuSames := dto.FindSkuSames{
		SkuCode:       c.SkuCode,
		WarehouseCode: c.WarehouseCode,
	}
	if e.ExistsSkuWarehouseUniquee(&findSkuSames) {
		status = false
		msg = append(msg, "SKU+仓库不唯一，同一个仓库下存在可用的相同SKU！")
	}
	if _, ok := warehouseMap[c.WarehouseCode]; !ok {
		msg = append(msg, fmt.Sprintf("仓库[%s}]为虚拟仓或该实体仓下无逻辑仓!", c.WarehouseCode))
	}
	if _, ok := vendorMap[c.VendorId]; !ok {
		msg = append(msg, "货主无效！")
	}
	if !utils.InArrayString(c.WarehouseCode, authWarehouse) {
		msg = append(msg, "仓库无权限！")
	}
	return status, msg
}

func (e *Goods) importValidate(c *dto.ImportGoods, model *models.Goods, vendorList map[int]modelsWc.Vendors, warehouseList map[string]dtoWc.WarehouseGetPageResp, companyList map[int]string, p *actions.DataPermission) error {
	c.Generate(model)
	// 判断公司是否存在
	companyId := e.getCompanyId(companyList, c.CompanyName)
	if companyId == 0 {
		return errors.New("公司不存在或公司无效！")
	}
	// 判断仓库是否为当前公司的仓库
	warehouseCode := e.getWarehouseCode(warehouseList, c.WarehouseName, c.CompanyName)
	if warehouseCode == "" {
		return errors.New("仓库不存在或仓库未启用！")
	}
	model.WarehouseCode = warehouseCode
	// 判断货主是否存在
	vendorId := e.getVendorId(vendorList, c.VendorName)
	if vendorId == 0 {
		return errors.New("货主不存在！")
	}
	model.VendorId = vendorId
	if c.ProductNo != "" {
		reg := regexp.MustCompile(`^[a-zA-Z0-9]{1,30}$`)
		res := reg.FindAllString(c.ProductNo, -1)
		if len(res) == 0 {
			return errors.New("物料编码只能英文、数字，长度在30以内，请检查！")
		}
	}
	if !e.FindSkuIsBind(model) {
		return errors.New("(SKU+货主+货主SKU)与驿站SKU无绑定关系！")
	}

	// 判断是新增还是更新（sku+货主+仓库唯一）
	var goods models.Goods
	findSameWhere := dto.FindSameWhere{
		SkuCode:       model.SkuCode,
		WarehouseCode: model.WarehouseCode,
		VendorId:      model.VendorId,
		Status:        model.Status,
		ApproveStatus: model.ApproveStatus,
	}

	if e.ExistsSkuSupplierWarehouseUniquee(&findSameWhere, &goods) {
		model.Id = goods.Id
		// 更新 继续更新校验
		return e.importUpdateValidate(model)
	} else {
		// 新增 继续新增校验
		// 获取用户权限仓库
		authWarehouseCodes := utils.Split(p.AuthorityWarehouseId)
		// api获取仓库信息
		warehouseMap := e.apiGetWarehouseInfoByCode([]string{model.WarehouseCode})
		// api获取货主信息
		vendorMap := e.apiGetVendorInfoById([]int{model.VendorId})
		var goodsinsertReq dto.GoodsInsertReq
		copier.Copy(&goodsinsertReq, &model)
		status, msg := e._validate(&goodsinsertReq, warehouseMap, vendorMap, authWarehouseCodes)
		if status {
			return nil
		}
		return errors.New(strings.Join(msg, ","))
	}
}

func (e *Goods) importUpdateValidate(model *models.Goods) error {
	var before models.Goods
	err := e.findGoods(model, &before)
	if err != nil {
		return err
	}
	if len(model.PriceModifyReason) > 200 {
		return errors.New("价格修改备注长度不能大于200！")
	}
	// 如果停用
	if model.Status == 0 && before.Status == 1 {
		// 如果是停用
		if before.OnlineStatus == 1 {
			return errors.New("请先下架该商品！")
		}
		// 校验 SKU+仓库+货主 库存为0时 可以停用，否则不可停用
		goodsStockMap := e.getStockMap([]int{before.Id}, model.WarehouseCode)
		if goodsStockMap[before.Id] > 0 {
			return errors.New("SKU+仓库+货主 库存 > 0,不可停用！")
		}
	} else {
		model.Status = 1
		var goods models.Goods
		findSameWhere := dto.FindSameWhere{
			SkuCode:       model.SkuCode,
			WarehouseCode: model.WarehouseCode,
			VendorId:      model.VendorId,
			Status:        model.Status,
			ApproveStatus: model.ApproveStatus,
		}
		if before.Status == 0 && e.ExistsSkuSupplierWarehouseUniquee(&findSameWhere, &goods) {
			return errors.New("SKU+货主+仓库不唯一，同一个仓库下存在相同SKU！")
		}
		findSkuSames := dto.FindSkuSames{
			SkuCode:       model.SkuCode,
			WarehouseCode: model.WarehouseCode,
		}
		if before.Status == 0 && e.ExistsSkuWarehouseUniquee(&findSkuSames) {
			return errors.New("SKU+仓库不唯一，同一个仓库下存在可用的相同SKU！")
		}
	}
	return nil
}

// 获取公司名称
func (e *Goods) getCompanyId(companyList map[int]string, companyName string) int {
	companyId := 0
	for k, v := range companyList {
		if v == companyName {
			companyId = k
		}
	}
	return companyId
}

// 获取仓库Code
func (e *Goods) getWarehouseCode(warehouselist map[string]dtoWc.WarehouseGetPageResp, warehouseName string, companyName string) string {
	warehouseCode := ""
	for k, warehouse := range warehouselist {
		if warehouse.WarehouseName == warehouseName && warehouse.CompanyName == companyName {
			warehouseCode = k
		}
	}
	return warehouseCode
}

// 获取货主ID
func (e *Goods) getVendorId(vendorList map[int]modelsWc.Vendors, vendorName string) int {
	vendorId := 0
	for k, vendors := range vendorList {
		if vendors.NameZh == vendorName {
			vendorId = k
		}
	}
	return vendorId
}

// ExistsSkuSupplierWarehouseUniquee SKU+货主+仓库唯一
func (e *Goods) ExistsSkuSupplierWarehouseUniquee(c *dto.FindSameWhere, data *models.Goods) bool {
	result := e.Orm.Scopes(dto.FindSame(c)).First(&data)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false
	}
	if result.Error != nil {
		return false
	}
	return true
}

// FindSkuIsBind (SKU+货主+货主SKU)与驿站SKU是否存在绑定关系
func (e *Goods) FindSkuIsBind(model *models.Goods) bool {
	productService := Product{e.Service}
	return productService.FindSkuIsBind(model)
}

// ExistsSkuWarehouseUniquee 仓库 + 可用sku唯一
func (e *Goods) ExistsSkuWarehouseUniquee(c *dto.FindSkuSames) bool {
	var data models.Goods
	result := e.Orm.Scopes(dto.FindSkuSame(c)).First(&data)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false
	}
	if result.Error != nil {
		return false
	}
	return true
}

func (e *Goods) doingApprove(before *models.Goods, data *models.Goods) {
	// 未上架或下架状态 只要是启用启用 或 更改‘物料编码’或‘销售价’重新进入审核，
	if before.OnlineStatus == 0 || before.OnlineStatus == 2 {
		if (data.Status == 1 && before.Status == 0) || (data.Status == 1 && before.Status == 1 && (data.ProductNo != before.ProductNo || data.MarketPrice != before.MarketPrice)) {
			data.ApproveStatus = 0
			data.ApproveRemark = ""
		}
	} else if before.OnlineStatus == 1 {
		if data.ProductNo != before.ProductNo || data.MarketPrice != before.MarketPrice {
			data.OnlineStatus = 2
			data.ApproveStatus = 0
			data.ApproveRemark = ""
			before.OnlineStatus = 2
			before.ApproveStatus = 0
		}
	}
}

func (e *Goods) filterData(data []models.Goods, filter []interface{}) (filterData []models.Goods) {
	for _, goods := range data {
		if goods.Status == 1 && utils.InArray(goods.OnlineStatus, filter) {
			filterData = append(filterData, goods)
		}
	}
	return
}

func (e *Goods) filterDisableData(data []models.Goods) (filterData []models.Goods) {
	for _, goods := range data {
		if goods.Status == 0 {
			filterData = append(filterData, goods)
		}
	}
	return
}

func (e *Goods) supplierHasDisabled(vendorId int) bool {
	vendorInfo := e.apiGetVendorInfoById([]int{vendorId})
	if _, ok := vendorInfo[vendorId]; ok {
		return true
	}
	return false
}

func (e *Goods) warehouseHasDisabled(warehouseCode string) bool {
	warehouseInfo := e.apiGetWarehouseInfoByCode([]string{warehouseCode})
	if _, ok := warehouseInfo[warehouseCode]; ok {
		return true
	}
	return false
}

func (e *Goods) splitJoin(data []models.Goods) (msg []string) {
	for _, goods := range data {
		msg = append(msg, goods.WarehouseCode+"-"+goods.SkuCode)
	}
	return
}

func (e *Goods) setMsgColor(msg string, color string) string {
	return fmt.Sprintf("<span style=\"color:%s\">%s</span>", color, msg)
}

func (e *Goods) apiGetVendorInfoById(vendorIds []int) map[int]modelsWc.Vendors {
	// 货主map
	vendorResult := wcClient.ApiByDbContext(e.Orm).GetVendorList(dtoWc.InnerVendorsGetListReq{
		Ids: strings.Trim(strings.Join(strings.Fields(fmt.Sprint(lo.Uniq(vendorIds))), ","), "[]"),
	})
	vendorResultInfo := &struct {
		response.Response
		Data []modelsWc.Vendors
	}{}
	vendorResult.Scan(vendorResultInfo)
	vendorMap := make(map[int]modelsWc.Vendors, len(vendorResultInfo.Data))
	for _, vendor := range vendorResultInfo.Data {
		vendorMap[vendor.Id] = vendor
	}
	return vendorMap
}

func (e *Goods) apiGetVendorInfoByName(vendorName []string) map[int]modelsWc.Vendors {
	// 货主map
	vendorResult := wcClient.ApiByDbContext(e.Orm).GetVendorList(dtoWc.InnerVendorsGetListReq{
		NameZh: strings.Join(lo.Uniq(vendorName), ","),
	})
	vendorResultInfo := &struct {
		response.Response
		Data []modelsWc.Vendors
	}{}
	vendorResult.Scan(vendorResultInfo)
	vendorMap := make(map[int]modelsWc.Vendors, len(vendorResultInfo.Data))
	for _, vendor := range vendorResultInfo.Data {
		vendorMap[vendor.Id] = vendor
	}
	return vendorMap
}

func (e *Goods) apiGetWarehouseInfoByCode(warehouseCodes []string) map[string]dtoWc.WarehouseGetPageResp {
	// 货主map
	warehouseResult := wcClient.ApiByDbContext(e.Orm).GetWarehouseList(dtoWc.InnerWarehouseGetListReq{
		WarehouseCode: strings.Join(lo.Uniq(warehouseCodes), ","),
	})
	warehouseResultInfo := &struct {
		response.Response
		Data []dtoWc.WarehouseGetPageResp
	}{}
	warehouseResult.Scan(warehouseResultInfo)
	warehouseMap := make(map[string]dtoWc.WarehouseGetPageResp, len(warehouseResultInfo.Data))
	for _, warehouse := range warehouseResultInfo.Data {
		warehouseMap[warehouse.WarehouseCode] = warehouse
	}
	return warehouseMap
}

func (e *Goods) apiGetWarehouseInfoByName(warehouseNames []string) map[string]dtoWc.WarehouseGetPageResp {
	// 货主map
	warehouseResult := wcClient.ApiByDbContext(e.Orm).GetWarehouseList(dtoWc.InnerWarehouseGetListReq{
		WarehouseName: strings.Join(lo.Uniq(warehouseNames), ","),
	})
	warehouseResultInfo := &struct {
		response.Response
		Data []dtoWc.WarehouseGetPageResp
	}{}
	warehouseResult.Scan(warehouseResultInfo)
	warehouseMap := make(map[string]dtoWc.WarehouseGetPageResp, len(warehouseResultInfo.Data))
	for _, warehouse := range warehouseResultInfo.Data {
		warehouseMap[warehouse.WarehouseCode] = warehouse
	}
	return warehouseMap
}

func (e *Goods) apiGetCompanyList(companyName []string) map[int]string {
	companyResult := ucClient.ApiByDbContext(e.Orm).GetCompanyByName(strings.Join(companyName, ","))
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

// findGoods (SKU+货主+货主SKU)
func (e *Goods) findGoods(c *models.Goods, model *models.Goods) error {
	var data models.Goods
	err := e.Orm.Model(&data).Scopes(dto.FindGoods(&dto.FindGoodsReq{
		SkuCode:         c.SkuCode,
		WarehouseCode:   c.WarehouseCode,
		SupplierSkuCode: c.SupplierSkuCode,
	})).First(&model).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		return err
	}
	if err != nil {
		return err
	}
	return nil
}

func (e *Goods) findGoodsById(goodsId int, data *models.Goods) error {
	return e.Orm.Model(&data).First(&data, goodsId).Error
}

func (e *Goods) findGoodsByIds(goodsId []int, data *[]models.Goods) error {
	return e.Orm.Model(&data).Find(data, goodsId).Error
}

//------------------------------------------------------INNER----------------------------------------------------------------

func (e *Goods) GetGoodsInfo(c *dto.GoodsInfoReq, data *[]models.Goods) error {
	var model models.Goods
	db := e.Orm.Model(&model)
	for _, info := range c.Query {
		db.Or(info)
	}
	err := db.Find(&data).Error
	if err != nil {
		return err
	}
	return nil
}

func (e *Goods) GetGoodsById(c *dto.GetGoodsByIdReq, data *[]models.Goods) error {
	var model models.Goods
	err := e.Orm.Model(&model).Preload("Product.Brand").
		Preload("Product.MediaRelationship", func(db *gorm.DB) *gorm.DB {
			return db.Order("media_relationship.seq ASC")
		}).
		Preload("Product.MediaRelationship.MediaInstant").
		Find(&data, c.Ids).Error
	if err != nil {
		return err
	}
	return nil
}

func (e *Goods) GetGoodsBySkuCodeReq(c *dto.GetGoodsBySkuCodeReq, data *[]models.Goods) error {
	var model models.Goods
	tx := e.Orm.Model(&model).
		Preload("Product.MediaRelationship", func(db *gorm.DB) *gorm.DB {
			return db.Order("media_relationship.seq ASC")
		}).
		Preload("Product.MediaRelationship.MediaInstant").
		Preload("Product.Brand")
	tx.Where("sku_code in ?", c.SkuCode)
	if c.WarehouseCode != "" {
		tx.Where("warehouse_code = ?", c.WarehouseCode)
	}
	if c.Status > -1 {
		tx.Where("status = ?", c.Status)
	}
	if c.OnlineStatus > -1 {
		tx.Where("online_status = ?", c.OnlineStatus)
	}
	err := tx.Find(&data).Error
	if err != nil {
		return err
	}
	return nil
}

// getStockMap 获取库存map goodsid => stock
func (e *Goods) getStockMap(goodsIds []int, warehouseCode string) (stockMap map[int]int) {
	stockResult := wcClient.ApiByDbContext(e.Orm).GetStockListByGoodsIdAndWarehouseCode(dtoWc.InnerStockInfoGetByGoodsIdAndWarehouseCodeReq{
		GoodsIds:      goodsIds,
		WarehouseCode: warehouseCode,
	})
	stockResultInfo := &struct {
		response.Response
		Data []modelsWc.StockInfo
	}{}
	stockResult.Scan(stockResultInfo)
	stockMap = make(map[int]int)
	for _, stock := range stockResultInfo.Data {
		stockMap[stock.GoodsId] = stock.Stock
	}
	return
}
