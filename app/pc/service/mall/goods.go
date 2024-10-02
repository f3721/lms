package mall

import (
	"errors"
	"fmt"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"github.com/samber/lo"
	modelsUc "go-admin/app/uc/models"
	dtoUc "go-admin/app/uc/service/mall/dto"
	modelsWc "go-admin/app/wc/models"
	dtoWc "go-admin/app/wc/service/admin/dto"
	ucClient "go-admin/common/client/uc"
	wcClient "go-admin/common/client/wc"
	"go-admin/common/global"
	cModel "go-admin/common/models"
	"go-admin/common/utils"
	"sort"
	"strconv"
	"strings"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/pc/models"
	"go-admin/app/pc/service/mall/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type Goods struct {
	service.Service
}

// GetPage 获取Goods列表
func (e *Goods) GetPage(c *dto.GoodsGetPageReq, p *actions.DataPermission, list *dto.List, count *int64) error {
	var err error

	goodsResult := make([]dto.GoodsGetPageResp, 0)
	err = e.getDb(c).Scopes(cDto.Paginate(c.GetPageSize(), c.GetPageIndex())).
		Select([]string{
			"brand.id as brand_id",
			"brand.brand_zh",
			"brand.brand_en",
			"brand.first_letter",
			"goods.id as GoodsId",
			"goods.sku_code",
			"goods.market_price",
			"goods.product_no",
			"goods.vendor_id",
			"goods.warehouse_code",
			"product.name_zh",
			"product.tax",
			"product.mfg_model",
			"product.sales_moq",
			"product.sales_uom",
		}).
		Find(&goodsResult).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("GoodsService GetPage error:%s \r\n", err)
		return err
	}

	// 品牌导航
	list.BrandNav = e.brandNav(c)

	// 搜素历史记录
	if c.Keyword != "" {
		err = e.addHistoryKeyword(c)
		if err != nil {
			return err
		}
	}

	if *count > 0 {
		var skuCodes []string
		err = utils.StructColumn(&skuCodes, goodsResult, "SkuCode", "")
		if err != nil {
			return err
		}

		skuCategory := make(map[int]dto.CategoryList)

		allResult := e.getAllSku(c)
		var allSku []string
		_ = utils.StructColumn(&allSku, allResult, "SkuCode", "")
		// 分类导航
		categoryService := Category{e.Service}
		if c.CategoryId > 0 {
			categoryPath := categoryService.GetCategoryPath(c.CategoryId)
			categoryPathModel := models.CategoryPath{}
			allCategorySku := categoryPathModel.GetSkuCodeByCateId(e.Service.Orm, categoryPath[0].Id)
			err = categoryService.GetCategoryBySku(allCategorySku, &skuCategory)
			if err != nil {
				return err
			}
			categoryNav := categoryService.GetParentCategoryList(c.CategoryId, allSkuCategoryToSlice(skuCategory))
			list.CategoryNav = categoryNav
		} else {
			err = categoryService.GetCategoryBySku(allSku, &skuCategory)
			if err != nil {
				return err
			}
			categoryNav := categoryService.GetParentCategoryList(c.CategoryId, allSkuCategoryToSlice(skuCategory))
			list.CategoryNav = categoryNav

		}
		// 商品列表所有品牌
		allBrand := e.getAllBrand(c)
		list.BrandAll = e.brandAll(&allBrand)

		//商品列表所有货主
		allVendor := e.getAllVendor(c)
		list.Vendor = e.vendorAll(&allVendor)

		// 所有品牌字母排序
		list.BrandFilter = e.brandFilterList(&allBrand)

		// 数据组装
		err = e.AssembleGoods(c, &goodsResult)
		if err != nil {
			return err
		}
		// 综合排序
		if c.MarketPrice == "" {
			e.comprehensiveSorting(skuCodes)
		}
		list.Product = goodsResult
		list.Category = e.goodsListCategoryList(c, allSkuCategoryToSlice(skuCategory))
	}
	return nil
}

// Get 获取Goods对象
func (e *Goods) Get(d *dto.GoodsGetReq, p *actions.DataPermission, resp *dto.GoodsGetResp) error {
	var data models.Goods

	var goodsInfo dto.GoodsGetInfo
	err := e.Orm.Model(&data).
		Preload("Product.Brand").
		Preload("Product.MediaRelationship", func(db *gorm.DB) *gorm.DB {
			return db.Order("media_relationship.seq ASC")
		}).
		Preload("Product.MediaRelationship.MediaInstant").
		Scopes(
			dto.MakeWarehouseCondition(d),
		).
		First(&goodsInfo, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		return err
	}
	if err != nil {
		return err
	}

	// 是否当前仓库
	if goodsInfo.WarehouseCode == d.WarehouseCode {
		goodsInfo.IsCurrentWarehouse = true
	}

	// 分类导航
	categoryService := Category{e.Service}
	result := make(map[int]dto.CategoryList)
	err = categoryService.GetCategoryBySku([]string{goodsInfo.SkuCode}, &result)
	if err != nil {
		return err
	}
	navigation := getKV(result)
	resp.Navigation = navigation

	// 扩展属性
	var productExtAttr any
	productCategory := ProductCategory{e.Service}
	productExtAttribute := ProductExtAttribute{e.Service}
	productCategoryId := productCategory.GetCategoryBySkuMainCategoryId(goodsInfo.SkuCode)
	if productCategoryId != 0 {
		productExtAttr, _ = productExtAttribute.GetAttrs(&dto.GetProductExtAttributeReq{
			SkuCode:    goodsInfo.SkuCode,
			CategoryId: productCategoryId,
		})
	}
	resp.ProductExtAttribute = productExtAttr

	// 货主
	vendorResult := e.ApiGetVendorInfoById([]int{goodsInfo.VendorId})

	// 库存
	goodsStockMap := e.GetStockMap([]int{goodsInfo.Id}, d.WarehouseCode)

	// 是否收藏
	userCollectMap := e.ApiGetUserCollectList(d.UserId, []int{goodsInfo.Id})
	goodsInfo.IsCollect = userCollectMap[goodsInfo.Id]

	//用户购物车商品map
	userCartProductList := e.UserCartProductList(&dto.GetProductListReq{
		WarehouseCode: d.WarehouseCode,
		UserId:        d.UserId,
		GoodsId:       []int{goodsInfo.Id},
	})

	var companyInfoM modelsUc.CompanyInfo
	companyInfo := companyInfoM.GetRowByUserId(e.Orm, d.UserId)

	// 数据组装
	goodsInfo.VendorName = vendorResult[goodsInfo.VendorId].NameZh
	goodsInfo.SalesMoq = goodsInfo.Product.SalesMoq
	goodsInfo.Stock = goodsStockMap[goodsInfo.Id]
	score, _ := strconv.ParseFloat(goodsInfo.Product.Tax, 64)
	goodsInfo.NakedSalePrice = utils.DivFloat64Fixed(goodsInfo.MarketPrice, 1+score, true)
	goodsInfo.TaxStr = strconv.FormatFloat(score*100, 'f', 0, 64) + "%"
	goodsInfo.Quantity = userCartProductList[goodsInfo.Id]
	goodsInfo.CheckStockStatus = companyInfo.CheckStockStatus
	resp.Info = goodsInfo

	return nil
}

func (e *Goods) ApiGetVendorInfoById(vendorIds []int) map[int]modelsWc.Vendors {
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

func (e *Goods) getDb(c *dto.GoodsGetPageReq) (db *gorm.DB) {
	var data models.Goods
	db = e.Orm.Model(&data).
		Joins("INNER JOIN product ON goods.sku_code = product.sku_code").
		Joins("INNER JOIN brand ON product.brand_id = brand.id").
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			dto.MakeSearchCondition(c),
			dto.MakeFilterCategoryIdCondition(c, e.Orm),
		)
	if c.HasStock == 1 {
		pcPrefix := global.GetTenantWcDBNameWithDB(e.Orm)
		db.Joins("INNER JOIN " + pcPrefix + ".stock_info si ON (goods.id = si.goods_id AND goods.warehouse_code = si.warehouse_code)").
			Where("si.stock > 0")
	}

	if c.MarketPrice == "" {
		db.Order("product.seq ASC").Order("goods.id ASC")
	}
	return
}

func (e *Goods) getMiniProgramDb(c *dto.GoodsGetReq) (db *gorm.DB) {
	var data models.Goods
	db = e.Orm.Model(&data).
		Joins("INNER JOIN product ON goods.sku_code = product.sku_code").
		Joins("INNER JOIN brand ON product.brand_id = brand.id").
		Scopes(
			dto.MakeWarehouseCondition(c),
		)
	return
}

func (e *Goods) GetMiniProgramHomeFilter(c *dto.GoodsGetReq, list *dto.MiniProgramHomeFilter) error {
	var err error
	// 仓库所有品牌
	allBrand := e.getMiniProgramAllBrand(c)
	list.BrandAll = e.brandAll(&allBrand)

	// 商品列表所有货主
	allVendor := e.getMiniProgramAllVendor(c)
	list.Vendor = e.vendorAll(&allVendor)

	return err
}

func (e *Goods) getMiniProgramAllBrand(c *dto.GoodsGetReq) []dto.AllBrand {
	skuList := make([]dto.AllBrand, 0)
	e.getMiniProgramDb(c).Select([]string{
		"goods.sku_code",
		"brand.id as brand_id",
		"brand.brand_zh",
	}).Group("brand.id").Order("brand.id ASC").Find(&skuList)
	return skuList
}

func (e *Goods) getAllSku(c *dto.GoodsGetPageReq) []dto.NoWhereResult {
	skuList := make([]dto.NoWhereResult, 0)
	e.getDb(c).Select([]string{
		"goods.sku_code",
		"brand.id as brand_id",
		"brand.brand_zh",
	}).Order("brand.id ASC").Find(&skuList)
	return skuList
}

func (e *Goods) getAllBrand(c *dto.GoodsGetPageReq) []dto.AllBrand {
	skuList := make([]dto.AllBrand, 0)
	e.getDb(c).Select([]string{
		"goods.sku_code",
		"brand.id as brand_id",
		"brand.brand_zh",
		"brand.first_letter",
	}).Group("brand.id").Order("brand.id ASC").Find(&skuList)
	return skuList
}

func (e *Goods) getMiniProgramAllVendor(c *dto.GoodsGetReq) []dto.AllVendor {
	vendorList := make([]dto.AllVendor, 0)
	e.getMiniProgramDb(c).Select([]string{
		"product.vendor_id",
	}).Group("product.vendor_id").Order("product.vendor_id ASC").Find(&vendorList)
	return vendorList
}

func (e *Goods) getAllVendor(c *dto.GoodsGetPageReq) []dto.AllVendor {
	vendorList := make([]dto.AllVendor, 0)
	e.getDb(c).Select([]string{
		"product.vendor_id",
	}).Group("product.vendor_id").Order("product.vendor_id ASC").Find(&vendorList)
	return vendorList
}

// GetStockMap 获取库存map goodsIds => stock
func (e *Goods) GetStockMap(goodsIds []int, warehouseCode string) (stockMap map[int]int) {
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

func getKV(categoryList map[int]dto.CategoryList) []map[string]any {
	var categoryMap []map[string]any
	for _, list := range categoryList {
		if list.Id1 != 0 {
			categoryMap = append(categoryMap, map[string]any{
				"id":   list.Id1,
				"name": list.NameZh1,
			})
		}
		if list.Id2 != 0 {
			categoryMap = append(categoryMap, map[string]any{
				"id":   list.Id2,
				"name": list.NameZh2,
			})
		}
		if list.Id3 != 0 {
			categoryMap = append(categoryMap, map[string]any{
				"id":   list.Id3,
				"name": list.NameZh3,
			})
		}
		if list.Id4 != 0 {
			categoryMap = append(categoryMap, map[string]any{
				"id":   list.Id4,
				"name": list.NameZh4,
			})
		}
	}
	return categoryMap
}

func (e *Goods) ApiGetUserCollectList(userId int, goodsIds []int) map[int]bool {
	userCollectResult := ucClient.ApiByDbContext(e.Orm).GetUserCollect(&dtoUc.UserCollectGetGoodsIsCollected{
		GoodsIds: strings.Trim(strings.Join(strings.Fields(fmt.Sprint(lo.Uniq(goodsIds))), ","), "[]"),
		UserId:   userId,
	})
	userCollectInfo := &struct {
		response.Response
		Data map[int]*dtoUc.UserCollectGetIsUserCollectResData
	}{}
	userCollectResult.Scan(userCollectInfo)
	result := make(map[int]bool, 0)
	for _, data := range userCollectInfo.Data {
		result[data.GoodsId] = data.IsCollect
	}
	return result
}

func (e *Goods) UserCartProductList(c *dto.GetProductListReq) map[int]int {
	userCart := make([]models.UserCart, 0)
	userCartService := UserCart{e.Service}
	_ = userCartService.GetProductListById(c, &userCart)

	var userProductMap = make(map[int]int, 0)
	for _, cart := range userCart {
		userProductMap[cart.GoodsId] = cart.Quantity
	}
	return userProductMap
}

// ComprehensiveSorting 综合排序
func (e *Goods) comprehensiveSorting(skuCodes []string) {
	// 获取商品分类
	skuCategoryInfo := make([]dto.CategoryInfo, 0)
	productCategory := ProductCategory{e.Service}
	_ = productCategory.GetCategoryBySkuMainCate(skuCodes, &skuCategoryInfo)

	// 分类父级关系映射
	var categoryMap map[int]int
	_ = utils.StructColumn(&categoryMap, skuCategoryInfo, "ParentId", "CategoryId")

	categoryGroup := map[int][]dto.CategoryInfo{}

	// 末级分类一样的放一起
	for _, info := range skuCategoryInfo {
		categoryGroup[info.CategoryId] = append(categoryGroup[info.CategoryId], info)
	}
	// 父级分类一样的放一起
	for i := range categoryGroup {
		for k, categoryInfos := range categoryGroup {
			if i == k {
				continue
			}
			// 父级分类相同
			if categoryMap[i] == categoryMap[k] {
				for _, info := range categoryInfos {
					categoryGroup[i] = append(categoryGroup[i], info)
				}
				delete(categoryGroup, k)
			}
		}
	}

	// 结果集重新排序
	categorySort := make([]dto.CategoryInfo, 0)
	for _, infos := range categoryGroup {
		for _, info := range infos {
			categorySort = append(categorySort, info)
		}
	}
	//fmt.Println(categorySort)
}

// AssembleGoods 组装
func (e *Goods) AssembleGoods(c *dto.GoodsGetPageReq, goodsResult *[]dto.GoodsGetPageResp) error {
	var skuCodes []string
	err := utils.StructColumn(&skuCodes, *goodsResult, "SkuCode", "")
	if err != nil {
		return err
	}
	// 产品图片
	mediaRelationship := make([]models.MediaRelationship, 0)
	mediaService := MediaRelationship{e.Service}
	err = mediaService.getMediaList(skuCodes, &mediaRelationship)
	if err != nil {
		return err
	}
	skuMedia := mediaRelation(mediaRelationship)

	// 货主
	var vendorIds []int
	err = utils.StructColumn(&vendorIds, *goodsResult, "VendorId", "")
	if err != nil {
		return err
	}
	vendorResult := e.ApiGetVendorInfoById(vendorIds)

	// 库存
	var goodsIds []int
	err = utils.StructColumn(&goodsIds, *goodsResult, "GoodsId", "")
	if err != nil {
		return err
	}
	goodsStockMap := e.GetStockMap(goodsIds, c.WarehouseCode)

	// 是否收藏
	userCollectMap := e.ApiGetUserCollectList(c.UserId, goodsIds)

	//用户购物车商品map
	userCartProductList := e.UserCartProductList(&dto.GetProductListReq{
		WarehouseCode: c.WarehouseCode,
		UserId:        c.UserId,
		GoodsId:       goodsIds,
	})

	temp := *goodsResult
	// 数据组装
	for k, goods := range temp {
		if len(skuMedia[goods.SkuCode]) > 0 {
			temp[k].Image = skuMedia[goods.SkuCode][0].MediaInstant
		}
		temp[k].VendorName = vendorResult[goods.VendorId].NameZh
		temp[k].Stock = goodsStockMap[goods.GoodsId]
		score, _ := strconv.ParseFloat(goods.Tax, 64)
		temp[k].NakedSalePrice = utils.DivFloat64Fixed(goods.MarketPrice, 1+score, true)
		temp[k].IsCollect = userCollectMap[goods.GoodsId]
		temp[k].Quantity = userCartProductList[goods.GoodsId]
	}
	goodsResult = &temp

	return nil
}

func (e *Goods) goodsListCategoryList(c *dto.GoodsGetPageReq, filterCategory []int) []map[string]any {
	categoryService := Category{e.Service}
	categoryMapArr := make([]map[string]any, 0)
	categoryList := make([]models.Category, 0)
	err := categoryService.GetCategoryListByParentId(c.CategoryId, &categoryList)
	if err != nil {
		return categoryMapArr
	}

	// 过滤出页面显示的分类
	for _, category := range categoryList {
		if lo.Contains[int](filterCategory, category.Id) {
			categoryMapArr = append(categoryMapArr, map[string]any{
				"id":   category.Id,
				"name": category.NameZh,
				"icon": category.MediaRelationship.MediaInstant.MediaDir,
			})
		}
	}

	return categoryMapArr
}

func (e *Goods) brandFilterList(allBrand *[]dto.AllBrand) []map[string][]dto.BrandInfo {
	brandList := make([]map[string][]dto.BrandInfo, 0)
	// 所有品牌字母排序
	var word []string
	err := utils.StructColumn(&word, *allBrand, "FirstLetter", "")
	if err != nil {
		return brandList
	}
	wordList := lo.Uniq[string](word)
	sort.Slice(wordList, func(i, j int) bool {
		return wordList[i] < wordList[j]
	})
	var brandMap map[int]dto.AllBrand
	err = utils.StructColumn(&brandMap, *allBrand, "", "BrandId")
	if err != nil {
		return brandList
	}

	for _, v := range wordList {
		mapV := map[string][]dto.BrandInfo{}
		mapB := make([]dto.BrandInfo, 0)
		for _, brand := range brandMap {
			if v == brand.FirstLetter {
				temp := dto.BrandInfo{
					BrandId: brand.BrandId,
					BrandZh: brand.BrandZh,
				}
				mapB = append(mapB, temp)
			}
		}
		mapV[v] = mapB
		brandList = append(brandList, mapV)
	}
	return brandList
}

func (e *Goods) brandAll(goodsResult *[]dto.AllBrand) []map[string]any {
	brandAll := make([]map[string]any, 0)
	for _, v := range *goodsResult {
		brandAll = append(brandAll, map[string]any{
			"id":          v.BrandId,
			"name":        v.BrandZh,
			"firstLetter": v.FirstLetter,
		})
	}
	return brandAll
}

func (e *Goods) vendorAll(allVendor *[]dto.AllVendor) []map[string]any {
	vendorAll := make([]map[string]any, 0)
	var vendorIds []int
	err := utils.StructColumn(&vendorIds, *allVendor, "VendorId", "")
	if err != nil {
		return vendorAll
	}
	vendorResult := e.ApiGetVendorInfoById(vendorIds)

	for _, v := range *allVendor {
		vendorAll = append(vendorAll, map[string]any{
			"id":   v.VendorId,
			"name": vendorResult[v.VendorId].NameZh,
		})
	}
	return vendorAll
}

func (e *Goods) brandNav(c *dto.GoodsGetPageReq) []string {
	brandNav := make([]string, 0)
	if len(c.BrandId) > 0 {
		brandList := make([]dto.BrandInfo, 0)
		brandService := Brand{e.Service}
		_ = brandService.GetBrandList(c.BrandId, &brandList)
		if len(brandList) > 0 {
			var brandName []string
			_ = utils.StructColumn(&brandName, brandList, "BrandZh", "")
			brandNav = append(brandNav, strings.Join(brandName, ","))
		}
	}
	return brandNav
}

func (e *Goods) addHistoryKeyword(c *dto.GoodsGetPageReq) error {
	userSearchHistory := UserSearchHistory{e.Service}
	err := userSearchHistory.Insert(&dto.UserSearchHistoryInsertReq{
		UserId:        c.UserId,
		WarehouseCode: c.WarehouseCode,
		Keyword:       c.Keyword,
		ControlBy: cModel.ControlBy{
			CreateBy:     c.UserId,
			CreateByName: c.UserName,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func allSkuCategoryToSlice(data map[int]dto.CategoryList) (uniqValues []int) {
	result := make([]int, 0)
	for _, list := range data {
		result = append(result, list.Id1)
		result = append(result, list.Id2)
		result = append(result, list.Id3)
		result = append(result, list.Id4)
	}
	uniqValues = lo.Uniq[int](result)
	return
}
