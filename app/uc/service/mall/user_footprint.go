package mall

import (
	"errors"
	modelsPc "go-admin/app/pc/models"
	"go-admin/app/uc/models"
	"go-admin/app/uc/service/mall/dto"
	"go-admin/common/actions"
	pcClient "go-admin/common/client/pc"
	cDto "go-admin/common/dto"
	commonModels "go-admin/common/models"
	"time"

	"github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type UserFootprint struct {
	service.Service
}

// GetPage 获取UserFootprint列表
func (e *UserFootprint) GetPage(c *dto.UserFootprintGetPageReq, p *actions.DataPermission, list *[]models.UserFootprint, count *int64) error {
	var err error
	var data models.UserFootprint

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("UserFootprintService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// GetListPage  获取UserCollect列表(包含商品名 商品图片)
func (e *UserFootprint) GetListPage(c *dto.UserFootprintGetPageReq, p *actions.DataPermission, resList *[]*dto.UserFootprintGetListPageRes, count *int64) error {
	var err error
	var list []models.UserFootprint

	err = e.GetPage(c, p, &list, count)
	if err != nil {
		return err
	}
	var goodIds []int
	for _, collect := range list {
		goodIds = append(goodIds, collect.GoodsId)
	}

	// 获取商品信息
	goodsResult := pcClient.ApiByDbContext(e.Orm).GetGoodsById(goodIds)

	goodsResultInfo := &struct {
		response.Response
		Data []modelsPc.Goods
	}{}
	goodsResult.Scan(goodsResultInfo)
	goodsMap := make(map[int]modelsPc.Goods)
	for _, goods := range goodsResultInfo.Data {
		goodsMap[goods.Id] = goods
	}

	var resListCopy []dto.UserFootprintGetListPageRes
	for _, collect := range list {
		goods, ok := goodsMap[collect.GoodsId]
		productData := dto.UserFootprintGetListPageProduct{}
		if ok {
			productData.ProductName = goods.Product.NameZh
			productData.ProductSalesMoq = goods.Product.SalesMoq
			productData.ProductMarketPrice = goods.MarketPrice
			if goods.Product.MediaRelationship != nil && len(*goods.Product.MediaRelationship) > 0 {
				mediaRelationship := (*goods.Product.MediaRelationship)[0]
				productData.ProductImage = mediaRelationship.MediaInstant.MediaDir
			}
		}
		resListCopy = append(resListCopy, dto.UserFootprintGetListPageRes{
			UserFootprint:                   collect,
			UserFootprintGetListPageProduct: productData,
		})
	}
	copier.Copy(&resList, resListCopy)
	return nil
}

// Get 获取UserFootprint对象
func (e *UserFootprint) Get(d *dto.UserFootprintGetReq, p *actions.DataPermission, model *models.UserFootprint) error {
	var data models.UserFootprint

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetUserFootprint error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建UserFootprint对象
func (e *UserFootprint) Insert(c *dto.UserFootprintInsertReq) error {
	var err error
	var data models.UserFootprint
	c.Generate(&data)

	goods := e.GetPcGetGoodsByGoodsId(e.Orm, c.GoodsId)
	if goods == nil {
		return errors.New("商品不存在")
	}
	data.WarehouseCode = goods.WarehouseCode
	data.SkuCode = goods.SkuCode

	err = e.Orm.Where("user_id = ?", c.UserId).
		Where("goods_id = ?", c.GoodsId).
		Assign(models.UserFootprint{
			ModelTime: commonModels.ModelTime{
				UpdatedAt: time.Now(),
			},
		}).
		FirstOrCreate(&data).Error
	if err != nil {
		e.Log.Errorf("UserCollectService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

func (e *UserFootprint) GetPcGetGoodsByGoodsId(db *gorm.DB, goodsId int) *modelsPc.Goods {
	// 获取商品信息
	goodIds := []int{goodsId}
	goodsMap := e.GetPcGetGoodsByGoodsIds(db, goodIds)

	if _, ok := goodsMap[goodsId]; !ok {
		return nil
	}
	return goodsMap[goodsId]
}

func (e *UserFootprint) GetPcGetGoodsByGoodsIds(db *gorm.DB, goodIds []int) (goodsMap map[int]*modelsPc.Goods) {
	// 获取商品信息
	goodsResult := pcClient.ApiByDbContext(db).GetGoodsById(goodIds)

	goodsResultInfo := &struct {
		response.Response
		Data []modelsPc.Goods
	}{}
	goodsResult.Scan(goodsResultInfo)
	goodsMap = make(map[int]*modelsPc.Goods)
	for _, goods := range goodsResultInfo.Data {
		goodsMap[goods.Id] = &goods
	}
	return
}

// Update 修改UserFootprint对象
func (e *UserFootprint) Update(c *dto.UserFootprintUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.UserFootprint{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("UserFootprintService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除UserFootprint
func (e *UserFootprint) Remove(d *dto.UserFootprintDeleteReq, p *actions.DataPermission) error {
	var data models.UserFootprint

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveUserFootprint error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
