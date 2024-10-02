package admin

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	adminClient "go-admin/common/client/admin"
	"go-admin/common/utils"
	"gorm.io/gorm"
	"strconv"
	"strings"

	"go-admin/app/wc/models"
	"go-admin/app/wc/service/admin/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type Vendors struct {
	service.Service
}

// 货主数据权限

func VendorPermission(p *actions.DataPermission, noCheckPermission int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if noCheckPermission != 1 {
			db.Where("vendors.id in ?", utils.SplitToInt(p.AuthorityVendorId))
		}
		return db
	}
}

// GetPage 获取Vendors列表

func (e *Vendors) GetPage(c *dto.VendorsGetPageReq, p *actions.DataPermission, list *[]models.Vendors, count *int64) error {
	var err error
	var data models.Vendors

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSizeNegative(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
			VendorPermission(p, c.NoCheckPermission),
		).Order("id DESC").
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("VendorsService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取Vendors对象
func (e *Vendors) Get(d *dto.VendorsGetReq, p *actions.DataPermission, model *models.Vendors) error {
	var data models.Vendors

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
			VendorPermission(p, 0),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetVendors error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建Vendors对象
func (e *Vendors) Insert(c *dto.VendorsInsertReq) error {
	var err error
	var data models.Vendors
	if resFlag, _ := data.CheckNameZh(e.Orm, c.NameZh, 0); !resFlag {
		return errors.New("货主中文名已存在")
	}
	if resFlag, _ := data.CheckCode(e.Orm, c.Code, 0); !resFlag {
		return errors.New("货主编码已存在")
	}
	c.Generate(&data)
	// code大写
	data.Code = strings.ToUpper(data.Code)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("VendorsService Insert error:%s \r\n", err)
		return err
	}
	//记录操作日志
	DataStr, _ := json.Marshal(&data)
	opLog := models.OperateLogs{
		DataId:       data.Code,
		ModelName:    models.VendorsModelName,
		Type:         models.VendorsModelInsert,
		DoStatus:     models.VendorsModelStatus[data.Status],
		Before:       "",
		Data:         string(DataStr),
		After:        string(DataStr),
		OperatorId:   c.CreateBy,
		OperatorName: c.CreateByName,
	}
	_ = opLog.InsertItem(e.Orm)
	_ = adminClient.ApiByDbContext(e.Orm).
		UpdatePermission(adminClient.ApiClientPermissionUpdateRequest{
			UserID:            user.GetUserId(e.Orm.Statement.Context.(*gin.Context)),
			AuthorityVendorID: strconv.Itoa(data.Id),
		})
	return nil
}

// Update 修改Vendors对象
func (e *Vendors) Update(c *dto.VendorsUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.Vendors{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
		VendorPermission(p, 0),
	).First(&data, c.GetId())
	if data.Id == 0 {
		return errors.New("货主不存在或没有权限")
	}
	oldData := data
	if resFlag, _ := data.CheckNameZh(e.Orm, c.NameZh, data.Id); !resFlag {
		return errors.New("货主中文名已存在")
	}
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("VendorsService Save error:%s \r\n", err)
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
		ModelName:    models.VendorsModelName,
		Type:         models.VendorsModelUpdate,
		DoStatus:     models.VendorsModelStatus[data.Status],
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
func (e *Vendors) Remove(d *dto.VendorsDeleteReq, p *actions.DataPermission) error {
	var data models.Vendors

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveVendors error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}

func (e *Vendors) Select(c *dto.VendorsSelectReq, p *actions.DataPermission, outData *[]dto.VendorsSelectResp, count *int64) error {
	var err error
	var data models.Vendors
	var list = &[]models.Vendors{}
	var outItem = &dto.VendorsSelectResp{}

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSizeNegative(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
			VendorPermission(p, c.NoCheckPermission),
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

func (e *Vendors) InnerGetList(c *dto.InnerVendorsGetListReq, list *[]models.Vendors) error {
	var err error
	var data models.Vendors
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
		e.Log.Errorf("VendorsService InnerGetList error:%s \r\n", err)
		return err
	}
	return nil
}
