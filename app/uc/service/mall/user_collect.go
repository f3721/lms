package mall

import (
	"errors"
	"fmt"
	modelsPc "go-admin/app/pc/models"
	modelsWc "go-admin/app/wc/models"
	dtoWc "go-admin/app/wc/service/admin/dto"
	pcClient "go-admin/common/client/pc"
	wcClient "go-admin/common/client/wc"
	"go-admin/common/global"
	"strconv"
	"strings"

	"github.com/samber/lo"

	"github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"github.com/jinzhu/copier"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/uc/models"
	"go-admin/app/uc/service/mall/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type UserCollect struct {
	service.Service
}

// GetPage 获取UserCollect列表
func (e *UserCollect) GetPage(c *dto.UserCollectGetPageReq, p *actions.DataPermission, list *[]models.UserCollect, count *int64) error {
	var err error
	var data models.UserCollect

	query := e.Orm.Table(data.TableName()).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		)

	pcPrefix := global.GetTenantPcDBNameWithDB(e.Orm)
	query = query.Joins("left JOIN " + pcPrefix + ".product as p ON p.sku_code = user_collect.sku_code")
	if c.FilterLevel1Catid > 0 {
		query = query.Where("user_collect.sku_code IN (?)", e.Orm.Table(pcPrefix+".category_path cp").
			Select("pc.sku_code").
			Joins("INNER JOIN "+pcPrefix+".product_category AS pc ON (cp.category_id = pc.category_id)").
			Where("path_id = ?", c.FilterLevel1Catid).
			Group("pc.sku_code"))
	}
	if c.FilterLevel2Catid > 0 {
		query = query.Where("user_collect.sku_code IN (?)", e.Orm.Table(pcPrefix+".category_path cp").
			Select("pc.sku_code").
			Joins("INNER JOIN "+pcPrefix+".product_category AS pc ON (cp.category_id = pc.category_id)").
			Where("path_id = ?", c.FilterLevel2Catid).
			Group("pc.sku_code"))
	}
	if c.FilterLevel3Catid > 0 {
		query = query.Where("user_collect.sku_code IN (?)", e.Orm.Table(pcPrefix+".category_path cp").
			Select("pc.sku_code").
			Joins("INNER JOIN "+pcPrefix+".product_category AS pc ON (cp.category_id = pc.category_id)").
			Where("path_id = ?", c.FilterLevel3Catid).
			Group("pc.sku_code"))
	}
	if c.FilterLevel4Catid > 0 {
		query = query.Where("user_collect.sku_code IN (?)", e.Orm.Table(pcPrefix+".category_path cp").
			Select("pc.sku_code").
			Joins("INNER JOIN "+pcPrefix+".product_category AS pc ON (cp.category_id = pc.category_id)").
			Where("path_id = ?", c.FilterLevel4Catid).
			Group("pc.sku_code"))
	}
	if c.FilterKeyword != "" {
		query = query.Where("(p.sku_code LIKE ? OR p.name_zh LIKE ?)", "%"+c.FilterKeyword+"%", "%"+c.FilterKeyword+"%")
	}
	if c.CompanyId != 0 {
		query = query.Joins("INNER JOIN user_info as ui ON ui.id = user_collect.user_id").
			Where("ui.company_id = ?", c.CompanyId).Group("goods_id")
	}

	err = query.Order("user_collect.id DESC").
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error

	if err != nil {
		e.Log.Errorf("UserCollectService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// GetListPage  获取UserCollect列表(包含商品名 商品图片)
func (e *UserCollect) GetListPage(c *dto.UserCollectGetPageReq, p *actions.DataPermission, resList *[]dto.UserCollectGetListPageRes, count *int64) error {
	var err error
	var list []models.UserCollect

	err = e.GetPage(c, p, &list, count)
	if err != nil {
		return err
	}
	var goodIds []int
	goodIdsMapWarehouseCode := make(map[string][]int)
	for _, collect := range list {
		goodIds = append(goodIds, collect.GoodsId)
		goodIdsMapWarehouseCode[collect.WarehouseCode] = append(goodIdsMapWarehouseCode[collect.WarehouseCode], collect.GoodsId)
	}

	stockMap := make(map[int]int)
	e.Log.Info(c.IsShowStock)
	if c.IsShowStock == 1 {
		for s, ints := range goodIdsMapWarehouseCode {
			// 获取库存信息
			e.Log.Info(ints, s)
			stockList := e.GetStockList(ints, s)
			for _, info := range stockList {
				stockMap[info.GoodsId] = info.Stock
			}
		}
	}

	// 获取商品信息
	goodsResult := pcClient.ApiByDbContext(e.Orm).GetGoodsById(goodIds)

	goodsResultInfo := &struct {
		response.Response
		Data []modelsPc.Goods
	}{}
	vendorId := []int{}
	goodsResult.Scan(goodsResultInfo)
	goodsMap := make(map[int]modelsPc.Goods)
	for _, goods := range goodsResultInfo.Data {
		goodsMap[goods.Id] = goods
		vendorId = append(vendorId, goods.VendorId)
	}

	vendorMap := e.apiGetVendorInfoById(vendorId)

	var resListCopy []dto.UserCollectGetListPageRes
	for _, collect := range list {
		goods, ok := goodsMap[collect.GoodsId]
		productData := dto.UserCollectGetListPageProduct{}
		if ok {
			productData.ProductName = goods.Product.NameZh
			productData.ProductSalesMoq = goods.Product.SalesMoq
			productData.ProductMarketPrice = goods.MarketPrice
			productData.ProductBrandZh = goods.Product.Brand.BrandZh
			productData.ProductMfgModel = goods.Product.MfgModel
			productData.ProductVendorName = goods.Product.VendorName
			if goods.Product.MediaRelationship != nil && len(*goods.Product.MediaRelationship) > 0 {
				mediaRelationship := (*goods.Product.MediaRelationship)[0]
				productData.ProductImage = mediaRelationship.MediaInstant.MediaDir
			}
			if _, ok1 := stockMap[collect.GoodsId]; ok1 {
				productData.ProductStock = stockMap[collect.GoodsId]
			}

			if _, ok1 := vendorMap[goods.VendorId]; ok1 {
				productData.ProductVendorName = vendorMap[goods.VendorId].NameZh
			}

		}
		resListCopy = append(resListCopy, dto.UserCollectGetListPageRes{
			UserCollect:                   collect,
			UserCollectGetListPageProduct: productData,
		})
	}
	copier.Copy(&resList, resListCopy)
	return nil
}

// GetGoodsIsCollected 批量获取商品收藏状态
func (e *UserCollect) GetGoodsIsCollected(d *dto.UserCollectGetGoodsIsCollected) (*map[int]*dto.UserCollectGetIsUserCollectResData, error) {
	if d.UserId == 0 || d.GoodsIds == "" {
		return nil, errors.New("参数错误")
	}

	var list []*models.UserCollect
	res := make(map[int]*dto.UserCollectGetIsUserCollectResData)

	goodsIdList := strings.Split(d.GoodsIds, ",")
	err := e.Orm.
		Where("user_id = ? ", d.UserId).
		Where("goods_id in ? ", goodsIdList).
		Find(&list).Error

	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return nil, err
	}

	for _, goodsIds := range goodsIdList {
		goodsId, _ := strconv.Atoi(goodsIds)
		res[goodsId] = &dto.UserCollectGetIsUserCollectResData{
			GoodsId:   goodsId,
			IsCollect: false,
		}
	}

	for _, collect := range list {
		if _, ok := res[collect.GoodsId]; ok {
			res[collect.GoodsId].IsCollect = true
		}
	}
	return &res, nil
}

// Get 获取UserCollect对象
func (e *UserCollect) Get(d *dto.UserCollectGetReq, p *actions.DataPermission, model *models.UserCollect) error {
	var data models.UserCollect

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetUserCollect error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建UserCollect对象
func (e *UserCollect) Insert(c *dto.UserCollectInsertReq) error {
	var err error
	var data models.UserCollect

	c.Generate(&data)

	goods := e.GetPcGetGoodsByGoodsId(e.Orm, c.GoodsId)
	if goods == nil {
		return errors.New("商品不存在")
	}
	data.WarehouseCode = goods.WarehouseCode
	data.SkuCode = goods.SkuCode

	err = e.Orm.Where("user_id = ?", c.UserId).Where("goods_id = ?", c.GoodsId).FirstOrCreate(&data).Error
	if err != nil {
		e.Log.Errorf("UserCollectService Insert error:%s \r\n", err)
		return err
	}

	return nil
}

// Update 修改UserCollect对象
func (e *UserCollect) Update(c *dto.UserCollectUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.UserCollect{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("UserCollectService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除UserCollect
func (e *UserCollect) Remove(d *dto.UserCollectDeleteReq, p *actions.DataPermission) error {
	var data models.UserCollect

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveUserCollect error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}

// Remove 删除UserCollect
func (e *UserCollect) RemoveByGoodsIds(d *dto.UserCollectDeleteGoodsIdsReq, p *actions.DataPermission) error {
	var data models.UserCollect

	db := e.Orm.Model(&data).
		Where("goods_id in ?", d.GoodsIds).
		Delete(&data)
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveUserCollect error:%s \r\n", err)
		return err
	}

	return nil
}

func (e *UserCollect) isExists(userId int, goodsId int) (bool, error) {
	var count int64
	err := e.Orm.Model(&models.UserCollect{}).
		Where("user_id = ?", userId).
		Where("goods_id = ?", goodsId).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (e *UserCollect) GetPcGetGoodsByGoodsIds(db *gorm.DB, goodIds []int) (goodsMap map[int]*modelsPc.Goods) {
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

func (e *UserCollect) GetPcGetGoodsByGoodsId(db *gorm.DB, goodsId int) *modelsPc.Goods {
	// 获取商品信息
	goodIds := []int{goodsId}
	goodsMap := e.GetPcGetGoodsByGoodsIds(db, goodIds)

	if _, ok := goodsMap[goodsId]; !ok {
		return nil
	}
	return goodsMap[goodsId]
}

// GetStockList 获取库存list
func (e *UserCollect) GetStockList(goodsIds []int, warehouseCode string) (stock []modelsWc.StockInfo) {
	stockResult := wcClient.ApiByDbContext(e.Orm).GetStockListByGoodsIdAndWarehouseCode(dtoWc.InnerStockInfoGetByGoodsIdAndWarehouseCodeReq{
		GoodsIds:      goodsIds,
		WarehouseCode: warehouseCode,
	})
	stockResultInfo := &struct {
		response.Response
		Data *[]modelsWc.StockInfo
	}{}
	stockResult.Scan(stockResultInfo)
	if stockResultInfo.Data != nil {
		return *stockResultInfo.Data
	}
	return []modelsWc.StockInfo{}
}

func (e *UserCollect) apiGetVendorInfoById(vendorIds []int) map[int]modelsWc.Vendors {
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
