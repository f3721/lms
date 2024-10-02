package admin

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/wc/models"
	"go-admin/app/wc/service/admin/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type QualityCheck struct {
	service.Service
}

// GetPage 获取QualityCheck列表
func (e *QualityCheck) GetPage(c *dto.QualityCheckGetPageReq, p *actions.DataPermission, outData *[]dto.QualityCheckRes, count *int64) error {
	var err error
	var data models.QualityCheck
	db := e.Orm.Model(&data).Preload("QualityCheckDetail").
		Scopes(
			//cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			//actions.Permission(data.TableName(), p),
			dto.GenQualityCheckSearch(c, "quality_check_task"),
		).
		Find(&outData)
	err = db.Error
	if err != nil {
		e.Log.Errorf("QualityCheck GetPage error:%s \r\n", err)
		return err
	}

	for k, od := range *outData {
		(*outData)[k].StatusName = models.QsStatusName[od.Status]
		(*outData)[k].TypeName = models.QsTypeName[od.Type]
		(*outData)[k].QualityStatusName = models.QualityStatusName[od.QualityStatus]
		(*outData)[k].QualityResName = models.QualityResName[od.QualityRes]
	}

	err = e.Orm.Table("(?) as u", db.Limit(-1).Offset(-1)).Count(count).Error
	if err != nil {
		e.Log.Errorf("QualityCheckService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取QualityCheck对象
func (e *QualityCheck) Get(d *dto.QualityCheckGetReq, p *actions.DataPermission, outData *dto.QualityCheckRes) error {
	var data models.QualityCheck
	err := e.Orm.Model(&data).Preload("QualityCheckDetail").
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(outData, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetQualityCheck error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	outData.StatusName = models.QsStatusName[outData.Status]
	outData.TypeName = models.QsTypeName[outData.Type]
	outData.QualityStatusName = models.QualityStatusName[outData.QualityStatus]
	outData.QualityResName = models.QualityResName[outData.QualityRes]
	return nil
}

// 获取配置项
func (e *QualityCheck) getInfo(d *dto.QualityCheckGetReq, p *actions.DataPermission) (*models.QualityCheckConfig, error) {
	var data models.QualityCheckConfig
	var model = &models.QualityCheckConfig{}

	err := e.Orm.Model(&data).
		Select("quality_check_config.*").
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetStockOutbound error:%s \r\n", err)
		return nil, err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return nil, err
	}
	return model, nil
}

// 上传质检明细
func (e *QualityCheck) UploadQualityRes(d *dto.QualityCheckUpdateReq, p *actions.DataPermission) error {
	var dataDetail models.QualityCheckDetail
	model := models.QualityCheck{}
	err := e.Orm.Model(&models.QualityCheck{}).Preload("QualityCheckDetail").
		Scopes(
		//actions.Permission(data.TableName(), p),
		).
		First(&model, d.Id).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetQualityCheckDetail error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}

	//1.先更新主信息，2.再更新配置结果，3增加事务处理，4需要增加子任务，按数量匹配
	var data models.QualityCheck
	curTime := time.Now()
	data.QuantityNum = d.QuantityNum
	if model.StayQualityNum > d.QuantityNum {
		data.StayQualityNum = model.StayQualityNum - d.QuantityNum
	}
	data.QualityRes = d.QualityRes
	data.QualityTime = curTime
	data.UpdateBy = model.CreateBy
	data.UpdateByName = model.CreateByName
	err = e.Orm.Model(&models.QualityCheck{}).Where("id = ?", d.Id).Updates(data).Error
	if err != nil {
		return err
	}
	//如果质检还有质检需要处理，继续新增
	if model.StayQualityNum > d.QuantityNum { //新的子表，待质检就是当前需要质检的数量
		model.StayQualityNum = d.QuantityNum
		model.QuantityNum = 0
		err = e.Orm.Model(&models.QualityCheck{}).Omit("id").Create(model).Error
		if err != nil {
			return err
		}
	}

	var configDetail []models.QualityCheckDetail
	for _, option := range d.QualityOption {
		item := models.QualityCheckDetail{
			QualityCheckTaskId: model.Id,
			QualityCheckOption: option.QualityCheckOption,
			QualityBy:          option.QualityBy,
			QualityRes:         option.QualityRes,
			Remark:             option.Remark,
		}
		item.CreateBy = model.CreateBy
		item.CreateByName = model.CreateByName
		item.QualityByName = model.CreateByName
		configDetail = append(configDetail, item)
	}
	err = e.Orm.Model(&dataDetail).Create(&configDetail).Error
	if err != nil {
		return err
	}
	return nil
}

// 配置信息
func (e *QualityCheck) getQualityConfigDetailInfo(id int, model *models.QualityCheckConfigDetail) error {
	var data models.QualityCheckConfigDetail
	err := e.Orm.Model(&data).First(model, id).Error
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

func (e *QualityCheck) Export(d *dto.QualityCheckGetPageReq, p *actions.DataPermission) (outData []dto.QualityExportRes, err error) {
	var data models.QualityCheck
	QualityCheckRes := []dto.QualityCheckRes{}
	db := e.Orm.Model(&data).Preload("QualityCheckDetail").
		Scopes(
			//cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(d.GetPageSize(), d.GetPageIndex()),
			//actions.Permission(data.TableName(), p),
			dto.GenQualityCheckSearch(d, "quality_check_task"),
		).
		Find(&QualityCheckRes)
	err = db.Error
	if err != nil {
		e.Log.Errorf("QualityCheck GetPage error:%s \r\n", err)
		return nil, err
	}

	var count int64
	err = e.Orm.Table("(?) as u", db.Limit(-1).Offset(-1)).Count(&count).Error
	if err != nil {
		e.Log.Errorf("QualityCheckService GetPage error:%s \r\n", err)
		return nil, err
	}

	for _, i2 := range QualityCheckRes {
		it := dto.QualityExportRes{
			ID:                 fmt.Sprintf("%v", i2.Id),
			QualityCheckCode:   i2.QualityCheckCode,
			SourceCode:         i2.SourceCode,
			EntryCode:          i2.EntryCode,
			WarehouseCode:      i2.WarehouseCode,
			LogicWarehouseCode: i2.LogicWarehouseCode,
			SourceName:         i2.SourceName,
			Status:             models.QsTypeName[i2.Status],
			Type:               models.QsTypeName[i2.Type],
			SkuCode:            i2.SkuCode,
			QualityStatus:      models.QualityStatusName[i2.QualityStatus],
			StayQualityNum:     fmt.Sprintf("%v", i2.StayQualityNum),
			QuantityNum:        fmt.Sprintf("%v", i2.QuantityNum),
			QualityRes:         fmt.Sprintf("%v", i2.QualityRes),
			QualityTime:        i2.QualityTime.Format("2006-01-02 15:04:05"),
		}
		outData = append(outData, it)
	}
	return outData, nil
}

// 打印
func (e *QualityCheck) Print() {

}

func (e *QualityCheck) IsQualityNumOk(req *dto.StockEntryPartReq, stockEntry *models.StockEntry) error {
	//先查配置开关，再看质检校验
	skuMap := map[string]int{}
	qualityRes := []models.QualityCheck{}
	err := e.Orm.Model(&models.QualityCheck{}).
		Where("warehouse_code = ?", stockEntry.WarehouseCode).
		Where("type = ?", stockEntry.Type).
		Find(&qualityRes).Error
	if err != nil {
		return err
	}
	if len(qualityRes) > 0 {
		for _, items := range qualityRes {
			if items.QualityRes == 2 {
				return errors.New("还有质检任务，请前往质检单查看！")
			}
			skuMap[items.SkuCode] += items.QuantityNum
		}
	}

	//对数据进行比对
	skuNumMap := map[string]int{}
	for _, product := range stockEntry.StockEntryProducts {
		skuNumMap[product.SkuCode] = product.ActQuantity
	}
	for _, todoProduct := range req.StockEntryProducts {
		skuNumMap[todoProduct.SkuCode] += todoProduct.ActQuantityTotal
	}

	//进行比对
	for sk, num := range skuNumMap {
		q_num, ok := skuMap[sk]
		if !ok {
			return errors.New("找不到sku")
		}

		if num < q_num {
			return errors.New("还有质检任务，请前往质检单查看！")
		}
	}
	return nil
}
