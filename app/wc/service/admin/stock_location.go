package admin

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"go-admin/common/excel"
	"go-admin/common/utils"
	"go-admin/config"
	"gorm.io/gorm"
	"mime/multipart"
	"regexp"

	"go-admin/app/wc/models"
	"go-admin/app/wc/service/admin/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type StockLocation struct {
	service.Service
}

// GetPage 获取StockLocation列表
func (e *StockLocation) GetPage(c *dto.StockLocationGetPageReq, p *actions.DataPermission, list *[]dto.StockLocationResp, count *int64) error {
	var err error
	var data models.StockLocation

	err = e.Orm.Model(&data).
		Select("stock_location.*, w.warehouse_name, lw.logic_warehouse_name").
		Joins("left join warehouse w on stock_location.warehouse_code = w.warehouse_code").
		Joins("left join logic_warehouse lw on stock_location.logic_warehouse_code = lw.logic_warehouse_code").
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
			dto.StockLocationGetPageMakeCondition(c),
			actions.SysUserPermission(data.TableName(), p, 2), //仓库权限
		).
		Scan(&list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("StockLocationService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取StockLocation对象
func (e *StockLocation) Get(d *dto.StockLocationGetReq, p *actions.DataPermission, model *models.StockLocation) error {
	var data models.StockLocation

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetStockLocation error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建StockLocation对象
func (e *StockLocation) Insert(c *dto.StockLocationInsertReq) error {
	var err error
	var data models.StockLocation
	err = c.InsertValid(e.Orm)
	if err != nil {
		return err
	}
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("StockLocationService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改StockLocation对象
func (e *StockLocation) Update(c *dto.StockLocationUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.StockLocation{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())

	err = c.UpdateValid(e.Orm)
	if err != nil {
		return err
	}
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("StockLocationService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除StockLocation
func (e *StockLocation) Remove(d *dto.StockLocationDeleteReq, p *actions.DataPermission) error {
	var data models.StockLocation

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveStockLocation error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}

// Import 导入
func (e *StockLocation) Import(file *multipart.FileHeader) (err error, errTitleList []map[string]string, errData []map[string]interface{}) {
	excelApp := excel.NewExcel()

	templateFilePath := config.ExtConfig.ApiHost + "/static/exceltpl/stock_location_import.xlsx"
	fieldsCorrect := excelApp.ValidImportFieldsCorrect(file, templateFilePath)
	if fieldsCorrect == false {
		err = errors.New("导入的excel字段与模板字段不一致，请重试")
		return
	}

	// 读取上传的excel 并返回data
	_, data, titleList := excelApp.GetExcelData(file)

	if len(data) == 0 {
		err = errors.New("导入数据不能为空！")
		return
	}
	data = utils.TrimMapSpace(data)
	// 校验data 获取 errMsgList
	errMsgList := map[int]string{}
	var successData []map[string]interface{}
	var errDataList []map[string]interface{}
	errKey := 0
	for _, datum := range data {
		if validErr := e.importValid(datum); validErr != nil {
			errMsgList[errKey] = validErr.Error()
			errDataList = append(errDataList, datum)
			errKey++
		} else {
			successData = append(successData, datum)
		}
	}
	if len(errMsgList) > 0 {
		// 有校验错误的行时调用 导出保存excel
		errTitleList, errData = excelApp.MergeErrMsgColumn(titleList, errDataList, errMsgList)
	}
	for _, datum := range successData {
		//if errMsg, ok := errMsgList[i]; ok && errMsg != "" {
		//	continue
		//}
		var stockLocation models.StockLocation
		err = stockLocation.GetByLocationCode(e.Orm, datum["location_code"].(string))
		stockLocation.Status = datum["status"].(string)
		stockLocation.Remark = datum["remark"].(string)
		// 默认启用
		if stockLocation.Status == "" {
			stockLocation.Status = "1"
		}
		if err == nil {
			e.Orm.Save(&stockLocation)
		} else if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			stockLocation.LocationCode = datum["location_code"].(string)
			var warehouse models.Warehouse
			if err = warehouse.GetByWarehouseName(e.Orm, datum["warehouse_name"].(string), nil); err == nil {
				stockLocation.WarehouseCode = warehouse.WarehouseCode
			}
			var logicWarehouse models.LogicWarehouse
			if err = logicWarehouse.GetLogicWarehouseByName(e.Orm, datum["logic_warehouse_name"].(string), nil); err == nil {
				stockLocation.LogicWarehouseCode = logicWarehouse.LogicWarehouseCode
			}
			e.Orm.Create(&stockLocation)
		}
	}
	return
}

func (e *StockLocation) importValid(data map[string]interface{}) (err error) {
	if data["location_code"] == nil || data["location_code"] == "" {
		return errors.New("库位编码必填")
	}
	reg := regexp.MustCompile(`^[A-Za-z]{2}[\d-]{1,18}$`)
	if !reg.MatchString(data["location_code"].(string)) {
		return errors.New("库位编号格式校验失败")
	}
	var stockLocation models.StockLocation
	if flag := stockLocation.CheckLocationCodeExist(e.Orm, data["location_code"].(string)); flag == false {
		if data["warehouse_name"] == nil || data["warehouse_name"] == "" {
			return errors.New("实体仓库必填")
		}

		p := actions.GetPermissionFromContext(e.Orm.Statement.Context.(*gin.Context))
		ruleWarehouseCodes := utils.Split(p.AuthorityWarehouseId)
		var warehouse models.Warehouse
		if err = warehouse.GetByWarehouseName(e.Orm, data["warehouse_name"].(string), ruleWarehouseCodes); err != nil {
			return err
		}
		if data["logic_warehouse_name"] == nil || data["logic_warehouse_name"] == "" {
			return errors.New("逻辑仓库必填")
		}
		var logicWarehouse models.LogicWarehouse
		if err = logicWarehouse.GetLogicWarehouseByName(e.Orm, data["logic_warehouse_name"].(string), ruleWarehouseCodes); err != nil {
			return err
		}
		if warehouse.WarehouseCode != logicWarehouse.WarehouseCode {
			return errors.New("实体仓和逻辑仓不匹配")
		}
	} else {
		var stockLocationGoods models.StockLocationGoods
		if total := stockLocationGoods.GetTotalStockByLocationCode(e.Orm, data["location_code"].(string)); data["status"].(string) == "1" && total > 0 {
			return errors.New("库位下无商品才可停用")
		}
	}

	if data["remark"] != "" && len([]rune(data["remark"].(string))) > 200 {
		return errors.New("备注长度不能超过200")
	}

	return
}
