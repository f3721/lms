package admin

import (
	"encoding/json"
	"errors"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"go-admin/common/utils"
	"gorm.io/gorm"
	"strings"

	"go-admin/app/wc/models"
	"go-admin/app/wc/service/admin/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type Supplier struct {
	service.Service
}

// GetPage 获取supplier列表
func (e *Supplier) GetPage(c *dto.SupplierGetPageReq, p *actions.DataPermission, list *[]models.Supplier, count *int64) error {
	var err error
	var data models.Supplier

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSizeNegative(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).Order("id DESC").
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("SupplierService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取Supplier对象
func (e *Supplier) Get(d *dto.SupplierGetReq, p *actions.DataPermission, model *models.Supplier) error {
	var data models.Supplier

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetSupplier error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建Supplie对象
func (e *Supplier) Insert(c *dto.SupplierInsertReq) error {
	var err error
	var data models.Supplier
	if resFlag, _ := data.CheckNameZh(e.Orm, c.NameZh, 0); !resFlag {
		return errors.New("供应商中文名已存在")
	}
	if resFlag, _ := data.CheckCode(e.Orm, c.Code, 0); !resFlag {
		return errors.New("供应商编码已存在")
	}
	c.Generate(&data)
	// code大写
	data.Code = strings.ToUpper(data.Code)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("SupplieService Insert error:%s \r\n", err)
		return err
	}
	//记录操作日志
	DataStr, _ := json.Marshal(&data)
	opLog := models.OperateLogs{
		DataId:       data.Code,
		ModelName:    models.SupplierModelName,
		Type:         models.SupplierModelInsert,
		DoStatus:     models.SupplierModelStatus[data.Status],
		Before:       "",
		Data:         string(DataStr),
		After:        string(DataStr),
		OperatorId:   c.CreateBy,
		OperatorName: c.CreateByName,
	}
	_ = opLog.InsertItem(e.Orm)
	//_ = adminClient.ApiByDbContext(e.Orm).
	//	UpdatePermission(adminClient.ApiClientPermissionUpdateRequest{
	//		UserID:            user.GetUserId(e.Orm.Statement.Context.(*gin.Context)),
	//		AuthorityVendorID: strconv.Itoa(data.Id),
	//	})
	return nil
}

// Update 修改Supplier对象
func (e *Supplier) Update(c *dto.SupplierUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.Supplier{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	if data.Id == 0 {
		return errors.New("供应商不存在")
	}
	oldData := data
	if resFlag, _ := data.CheckNameZh(e.Orm, c.NameZh, data.Id); !resFlag {
		return errors.New("供应商中文名已存在")
	}
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("SupplierService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	//记录操作日志
	oldDataStr, _ := json.Marshal(&oldData)
	DataStr, _ := json.Marshal(&data)
	opLog := models.OperateLogs{
		DataId:       data.Code,
		ModelName:    models.SupplierModelName,
		Type:         models.SupplierModelUpdate,
		DoStatus:     models.SupplierModelStatus[data.Status],
		Before:       string(oldDataStr),
		Data:         string(DataStr),
		After:        string(DataStr),
		OperatorId:   c.UpdateBy,
		OperatorName: c.UpdateByName,
	}
	_ = opLog.InsertItem(e.Orm)
	return nil
}

// Remove 删除Vendors
//func (e *Supplier) Remove(d *dto.VendorsDeleteReq, p *actions.DataPermission) error {
//	var data models.Vendors
//
//	db := e.Orm.Model(&data).
//		Scopes(
//			actions.Permission(data.TableName(), p),
//		).Delete(&data, d.GetId())
//	if err := db.Error; err != nil {
//		e.Log.Errorf("Service RemoveVendors error:%s \r\n", err)
//		return err
//	}
//	if db.RowsAffected == 0 {
//		return errors.New("无权删除该数据")
//	}
//	return nil
//}

func (e *Supplier) Select(c *dto.SupplierSelectReq, p *actions.DataPermission, outData *[]dto.SupplierSelectResp, count *int64) error {
	var err error
	var data models.Supplier
	var list = &[]models.Supplier{}
	var outItem = &dto.SupplierSelectResp{}

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSizeNegative(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("VendorsService GetPage error:%s \r\n", err)
		return err
	}
	for _, item := range *list {
		outItem.ReGenerate(&item)
		*outData = append(*outData, *outItem)
	}
	return nil
}

func (e *Supplier) InnerGetList(c *dto.InnerSupplierGetListReq, list *[]models.Supplier) error {
	var err error
	var data models.Supplier
	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			func(db *gorm.DB) *gorm.DB {
				if c.NameZh != "" {
					db.Where("name_zh in ?", utils.Split(c.NameZh))
				}
				if c.Code != "" {
					db.Where("code in ?", utils.Split(c.Code))
				}
				if c.Ids != "" {
					db.Where("id in ?", utils.SplitToInt(c.Ids))
				}
				return db
			},
		).
		Find(list).Error
	if err != nil {
		e.Log.Errorf("SupplierService InnerGetList error:%s \r\n", err)
		return err
	}
	return nil
}
