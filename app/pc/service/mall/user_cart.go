package mall

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	modelsUc "go-admin/app/uc/models"
	"go-admin/common/utils"
	"strconv"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"github.com/samber/lo"
	"gorm.io/gorm"

	"go-admin/app/pc/models"
	"go-admin/app/pc/service/mall/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type UserCart struct {
	service.Service
}

// GetPage 获取UserCart列表
func (e *UserCart) GetPage(c *dto.UserCartGetPageReq, p *actions.DataPermission, list *[]models.UserCart, count *int64) error {
	var err error
	var data models.UserCart

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("UserCartService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取UserCart对象
func (e *UserCart) Get(d *dto.UserCartGetReq, p *actions.DataPermission, model *models.UserCart) error {
	var data models.UserCart

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetUserCart error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建UserCart对象
func (e *UserCart) Insert(c *dto.UserCartInsertReq) error {
	var err error
	var data models.UserCart
	c.Generate(&data)
	err = data.Add(e.Orm, c.Quantity)
	if err != nil {
		e.Log.Errorf("UserCartService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改UserCart对象
func (e *UserCart) Update(c *dto.UserCartUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.UserCart{}
	c.Generate(&data)
	err = data.Edit(e.Orm, c.Quantity)
	if err != nil {
		e.Log.Errorf("UserCartService Save error:%s \r\n", err)
		return err
	}
	return nil
}

// Remove 删除UserCart
func (e *UserCart) Remove(d *dto.UserCartDeleteReq, p *actions.DataPermission) error {
	var err error
	var data models.UserCart
	d.Generate(&data)
	err = data.Remove(e.Orm)
	if err != nil {
		e.Log.Errorf("UserCartService Remove error:%s \r\n", err)
		return err
	}
	return nil
}

// SelectOne
func (e *UserCart) SelectOne(d *dto.UserCartSelectOneReq, p *actions.DataPermission) error {
	var err error
	var data models.UserCart
	d.Generate(&data)
	err = data.SelectOne(e.Orm)
	if err != nil {
		e.Log.Errorf("UserCartService SelectOne error:%s \r\n", err)
		return err
	}
	return nil
}

func (e *UserCart) SelectAll(d *dto.UserCartSelectAllReq, p *actions.DataPermission) error {
	var err error
	var data models.UserCart
	if e.IsUserCartSelectAll(d.UserId, d.WarehouseCode) {
		err = data.UnSelectAll(e.Orm, d.UserId, d.WarehouseCode)
	} else {
		err = data.SelectAll(e.Orm, d.UserId, d.WarehouseCode)
	}
	if err != nil {
		e.Log.Errorf("UserCartService SelectAll error:%s \r\n", err)
		return err
	}
	return nil
}

func (e *UserCart) UnSelectAll(d *dto.UserCartUnSelectAllReq, p *actions.DataPermission) error {
	var err error
	var data models.UserCart
	err = data.UnSelectAll(e.Orm, d.UserId, d.WarehouseCode)
	if err != nil {
		e.Log.Errorf("UserCartService UnSelectAll error:%s \r\n", err)
		return err
	}
	return nil
}

func (e *UserCart) ClearSelect(d *dto.UserCartClearSelectReq, p *actions.DataPermission) error {
	var data models.UserCart
	goodsIdsSaleStatus0, err := e.GetInvalidSaleStatusGoods(d.UserId, d.WarehouseCode)
	if err != nil {
		e.Log.Errorf("UserCartService ClearSelect error:%s \r\n", err)
		return err
	}
	err = data.ClearSelect(e.Orm, d.UserId, d.WarehouseCode, goodsIdsSaleStatus0)
	if err != nil {
		e.Log.Errorf("UserCartService ClearSelect error:%s \r\n", err)
		return err
	}
	return nil
}

// 清空失效商品

func (e *UserCart) ClearInvalid(d *dto.UserCartClearSelectReq, p *actions.DataPermission) error {
	goodsIdsSaleStatus0, err := e.GetInvalidSaleStatusGoods(d.UserId, d.WarehouseCode)
	if err != nil {
		e.Log.Errorf("UserCartService ClearSelect error:%s \r\n", err)
		return err
	}
	err = models.UserCartRemoveByGoodsIds(e.Orm, d.UserId, d.WarehouseCode, goodsIdsSaleStatus0)
	if err != nil {
		e.Log.Errorf("UserCartService ClearSelect error:%s \r\n", err)
		return err
	}
	return nil
}

// 获取失效商品

func (e *UserCart) GetInvalidSaleStatusGoods(userId int, warehouseCode string) ([]int, error) {
	var err error
	var data models.UserCart
	goodsIdsSaleStatus0 := []int{}
	userCarts, err := data.GetUserCartAll(e.Orm, userId, warehouseCode)
	if err != nil {
		return goodsIdsSaleStatus0, err
	}
	goodsResult := e.GetProducts(userCarts)
	for _, item := range *goodsResult {
		if item.ShowCartProduct != 1 {
			goodsIdsSaleStatus0 = append(goodsIdsSaleStatus0, item.GoodsId)
		}
	}
	return goodsIdsSaleStatus0, nil
}

func (e *UserCart) GetCartDataPage(c *dto.UserCartGetProductForCartPageReq, p *actions.DataPermission, outData *dto.UserCartGetProductForCartPageResp) error {
	var err error
	var data models.UserCart
	userCarts, err := data.GetUserCartAll(e.Orm, c.UserId, c.WarehouseCode)
	if err != nil {
		e.Log.Errorf("UserCartService GetProductForCart error:%s \r\n", err)
		return err
	}

	var companyInfoM modelsUc.CompanyInfo
	companyInfo := companyInfoM.GetRowByUserId(e.Orm, c.UserId)
	outData.CheckStockStatus = companyInfo.CheckStockStatus

	goodsResult := e.GetProducts(userCarts)
	e.FormatDataCartPage(goodsResult, outData)
	// 调用接口获取收藏信息
	goodsIds := lo.Map(*goodsResult, func(item models.UserCartGoodsProduct, index int) int {
		return item.GoodsId
	})
	goodsService := Goods{e.Service}
	collectMap := goodsService.ApiGetUserCollectList(c.UserId, goodsIds)
	pageGoodsResult := []models.UserCartPageGoodsProduct{}
	for _, item := range *goodsResult {
		pageGoodsResultItem := models.UserCartPageGoodsProduct{}
		_ = utils.CopyDeep(&pageGoodsResultItem, &item)
		pageGoodsResultItem.Collected = collectMap[item.GoodsId]
		pageGoodsResult = append(pageGoodsResult, pageGoodsResultItem)
	}
	outData.Products = pageGoodsResult
	return nil
}

func (e *UserCart) IsUserCartSelectAll(userId int, warehouseCode string) bool {
	var data models.UserCart
	userCarts, err := data.GetUserCartAll(e.Orm, userId, warehouseCode)
	if err != nil {
		return false
	}
	goodsResult := e.GetProducts(userCarts)
	isAll := true
	for _, item := range *goodsResult {
		if item.Selected == 0 && item.ShowCartProduct == 1 {
			isAll = false
			break
		}
	}
	return isAll
}

func (e *UserCart) FormatDataCartPage(cartProducts *[]models.UserCartGoodsProduct, outData *dto.UserCartGetProductForCartPageResp) {
	isSelectAll := 1
	for _, item := range *cartProducts {
		if item.ShowCartProduct == 1 && item.Selected == 0 {
			isSelectAll = 0
			break
		}
	}
	outData.IsSelectAll = isSelectAll

	isExistNotSale := 0
	for _, item := range *cartProducts {
		if item.ShowCartProduct == 0 {
			isExistNotSale = 1
			break
		}
	}
	outData.IsExistNotSale = isExistNotSale

	cartProductNum := 0
	totalProductNum := 0
	totalAmount := float64(0)
	totalNakedAmount := float64(0)

	for _, item := range *cartProducts {
		totalProductNum += item.Quantity
		if item.Selected == 1 && item.ShowCartProduct == 1 {
			cartProductNum += item.Quantity
			totalAmount = utils.AddFloat64Fixed(totalAmount, item.TotalMarketPrice, true)
			totalNakedAmount = utils.AddFloat64Fixed(totalNakedAmount, item.TotalNakedSalePrice, true)
		}
	}
	outData.TotalProductNum = totalProductNum
	outData.CartProductNum = cartProductNum
	outData.TotalAmount = totalAmount
	outData.TotalNakedAmount = totalNakedAmount
}

func (e *UserCart) VerifyCart(c *dto.UserCartVerifyCartReq, p *actions.DataPermission) error {
	var err error
	var data models.UserCart
	userCarts, err := data.GetUserCartSeleted(e.Orm, c.UserId, c.WarehouseCode)
	if err != nil {
		e.Log.Errorf("UserCartService VerifyCart error:%s \r\n", err)
		return err
	}
	goodsResult := e.GetProducts(userCarts)
	if len(*goodsResult) == 0 {
		return errors.New("购物车中未选择商品")
	}

	var companyInfoM modelsUc.CompanyInfo
	companyInfo := companyInfoM.GetRowByUserId(e.Orm, c.UserId)

	saleStatus1Slice := []int{}
	for _, item := range *goodsResult {
		if companyInfo.CheckStockStatus == 0 && item.Stock < item.Quantity {
			return errors.New("库存不足，不允许下单！")
		}
		if item.Quantity < item.SalesMoq {
			return errors.New(item.SkuCode + ",商品小于最小订货量")
		}
		if item.ShowCartProduct == 1 {
			saleStatus1Slice = append(saleStatus1Slice, item.GoodsId)
		}
	}
	if len(saleStatus1Slice) == 0 {
		return errors.New("请检查购物车是否还存在可售商品")
	}
	return nil
}

func (e *UserCart) VerifyBuyNow(c *dto.UserCartVerifyBuyNowReq, p *actions.DataPermission) error {
	var data = models.UserCart{
		GoodsId:       c.GoodsId,
		UserId:        c.UserId,
		WarehouseCode: c.WarehouseCode,
		Quantity:      c.Quantity,
		Selected:      1,
	}
	userCarts := &[]models.UserCart{data}
	goodsResult := e.GetProducts(userCarts)
	if len(*goodsResult) == 0 {
		return errors.New("未选择商品")
	}
	item := (*goodsResult)[0]
	if item.Quantity < item.SalesMoq {
		return errors.New(item.SkuCode + ",商品小于最小订货量")
	}
	if item.ShowCartProduct == 0 {
		return errors.New("不可售商品")
	}
	return nil
}

func (e *UserCart) GetProducts(userCarts *[]models.UserCart) *[]models.UserCartGoodsProduct {
	goodsIds := lo.Map(*userCarts, func(item models.UserCart, index int) int {
		return item.GoodsId
	})
	userCartsMap := lo.Associate(*userCarts, func(item models.UserCart) (int, models.UserCart) {
		return item.GoodsId, item
	})

	userCartGoodsProduct := &models.UserCartGoodsProduct{}
	goodsResult, _ := userCartGoodsProduct.GetByGoodsIds(e.Orm, goodsIds, "")
	_ = e.AssembleProducts(goodsResult)
	goodsResultMap := lo.Associate(*goodsResult, func(item models.UserCartGoodsProduct) (int, models.UserCartGoodsProduct) {
		return item.GoodsId, item
	})
	cartProductSlice := []models.UserCartGoodsProduct{}
	for _, item := range *userCarts {
		cartProduct := goodsResultMap[item.GoodsId]
		cartProduct.Quantity = item.Quantity
		cartProduct.TotalNakedSalePrice = utils.MulFloat64AndIntFixed(cartProduct.NakedSalePrice, cartProduct.Quantity, true)
		cartProduct.TotalMarketPrice = utils.MulFloat64AndIntFixed(cartProduct.MarketPrice, cartProduct.Quantity, true)
		cartProduct.TotalTaxPrice = utils.SubFloat64Fixed(cartProduct.TotalMarketPrice, cartProduct.TotalNakedSalePrice, true)
		cartProduct.Selected = item.Selected
		cartProduct.ShowCartProduct = 1
		cartProduct.SaleStatus = 1
		cartProductSlice = append(cartProductSlice, cartProduct)
	}
	e.VerifyCartProducts(&cartProductSlice, userCartsMap)
	return &cartProductSlice
}

func (e *UserCart) VerifyCartProducts(cartProduct *[]models.UserCartGoodsProduct, userCartsMap map[int]models.UserCart) {
	ctx := e.Orm.Statement.Context.(*gin.Context)
	var companyInfoM modelsUc.CompanyInfo
	companyInfo := companyInfoM.GetRowByUserId(e.Orm, user.GetUserId(ctx))

	for index, item := range *cartProduct {
		if item.WarehouseCode != userCartsMap[item.GoodsId].WarehouseCode {
			(*cartProduct)[index].ShowCartProduct = 0
			(*cartProduct)[index].SaleStatus = 0
			(*cartProduct)[index].SaleStatusRemark = "该商品不可用"
		}
		if item.Status != 1 {
			(*cartProduct)[index].ShowCartProduct = 0
			(*cartProduct)[index].SaleStatus = 0
			(*cartProduct)[index].SaleStatusRemark = "该商品不可用"
		}
		if item.OnlineStatus != 1 {
			(*cartProduct)[index].ShowCartProduct = 0
			(*cartProduct)[index].SaleStatus = 0
			(*cartProduct)[index].SaleStatusRemark = "该商品未上架"
		}
		if item.MarketPrice <= 0 {
			(*cartProduct)[index].ShowCartProduct = 0
			(*cartProduct)[index].SaleStatus = 0
			(*cartProduct)[index].SaleStatusRemark = "销售价格异常"
		}
		if companyInfo.CheckStockStatus == 0 && item.Stock < item.Quantity {
			(*cartProduct)[index].ShowCartProduct = 1
			(*cartProduct)[index].SaleStatus = 0
			(*cartProduct)[index].SaleStatusRemark = "该商品库存不足"
		}
	}
	tempSlice := []models.UserCartGoodsProduct{}

	for _, item := range *cartProduct {
		if item.ShowCartProduct == 0 {
			tempSlice = append(tempSlice, item)
		}
	}
	filterSlice := lo.Filter(*cartProduct, func(x models.UserCartGoodsProduct, _ int) bool {
		if x.ShowCartProduct == 0 {
			return false
		}
		return true
	})
	*cartProduct = append(filterSlice, tempSlice...)
}

func (e *UserCart) AssembleProducts(goodsResult *[]models.UserCartGoodsProduct) error {
	var skuCodes []string
	var err error
	_ = utils.StructColumn(&skuCodes, *goodsResult, "SkuCode", "")
	// 产品图片
	var mediaRelationship []models.MediaRelationship
	mediaService := MediaRelationship{e.Service}
	_ = mediaService.getMediaList(skuCodes, &mediaRelationship)

	skuMedia := mediaRelation(mediaRelationship)

	// 货主
	var vendorIds []int
	_ = utils.StructColumn(&vendorIds, *goodsResult, "VendorId", "")
	goodsService := Goods{e.Service}
	vendorResult := goodsService.ApiGetVendorInfoById(vendorIds)

	// 数据组装
	for k, goods := range *goodsResult {
		if len(skuMedia[goods.SkuCode]) > 0 {
			(*goodsResult)[k].Image = skuMedia[goods.SkuCode][0].MediaInstant
		}
		(*goodsResult)[k].VendorName = vendorResult[goods.VendorId].NameZh
		score, _ := strconv.ParseFloat(goods.Tax, 64)
		(*goodsResult)[k].NakedSalePrice = utils.DivFloat64Fixed(goods.MarketPrice, 1+score, true)
	}
	return err
}

// GetCartProductForOrder 下单获取购物车商品 for joe
func (e *UserCart) GetCartProductForOrder(userId int, warehouseCode string, checkStock bool) ([]models.UserCartGoodsProduct, error) {
	var err error
	var data models.UserCart
	var cartProduct = []models.UserCartGoodsProduct{}
	userCarts, err := data.GetUserCartSeleted(e.Orm, userId, warehouseCode)
	if err != nil {
		return cartProduct, err
	}
	goodsResult := e.GetProducts(userCarts)
	// 需要检查库存
	if checkStock == true {
		for _, gr := range *goodsResult {
			// 库存不足
			if gr.ShowCartProduct == 1 && gr.SaleStatus == 0 {
				err = errors.New("SKU：" + gr.SkuCode + " 库存不足，不允许下单！")
				return cartProduct, err
			}
		}
	}
	goodsResultOutput := lo.Filter(*goodsResult, func(x models.UserCartGoodsProduct, _ int) bool {
		if x.ShowCartProduct == 0 {
			return false
		}
		// 需要检查库存
		if checkStock == true {
			// 库存不足
			if x.ShowCartProduct == 1 && x.SaleStatus == 0 {
				return false
			}
		}
		return true
	})

	return goodsResultOutput, nil
}

// 下单获取立即购买商品 for joe
func (e *UserCart) GetBuyNowProductForOrder(userId, goodsId, quantity int, warehouseCode string) ([]models.UserCartGoodsProduct, error) {
	var data = models.UserCart{
		GoodsId:       goodsId,
		UserId:        userId,
		WarehouseCode: warehouseCode,
		Quantity:      quantity,
		Selected:      1,
	}
	userCarts := &[]models.UserCart{data}
	goodsResult := e.GetProducts(userCarts)
	return *goodsResult, nil
}

func (e *UserCart) GetProductListById(c *dto.GetProductListReq, data *[]models.UserCart) error {
	var model models.UserCart
	err := e.Orm.Model(&model).Scopes(dto.MakeProductListCondition(c)).Find(data).Error
	if err != nil {
		return err
	}
	return nil
}

// BatchAdd 创建UserCart对象
func (e *UserCart) BatchAdd(c *dto.UserCartBatchAddReq) error {
	var err error
	err = e.Orm.Transaction(func(tx *gorm.DB) error {
		for _, item := range c.Data {
			data := models.UserCart{
				UserId:        c.CreateBy,
				WarehouseCode: c.WarehouseCode,
				SkuCode:       item.SkuCode,
			}
			data.CreateBy = c.CreateBy
			data.CreateByName = c.CreateByName
			if err := data.AddBySku(tx, item.Quantity); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		e.Log.Errorf("UserCartService BatchAdd error:%s \r\n", err)
		return err
	}
	return nil
}

func (e *UserCart) GetOrderProductsPage(c *dto.UserCartGetProductForOrderPageReq, p *actions.DataPermission, outData *dto.UserCartGetProductForOrderPageResp) error {
	goodsResult, err := e.GetCartProductForOrder(c.UserId, c.WarehouseCode, false)
	if err != nil {
		return err
	}
	var companyInfoM modelsUc.CompanyInfo
	companyInfo := companyInfoM.GetRowByUserId(e.Orm, c.UserId)
	outData.CheckStockStatus = companyInfo.CheckStockStatus
	e.FormatProductOrderPage(&goodsResult, outData)
	outData.Products = goodsResult
	return nil
}

func (e *UserCart) FormatProductOrderPage(cartProducts *[]models.UserCartGoodsProduct, outData *dto.UserCartGetProductForOrderPageResp) {
	totalAmount := float64(0)
	totalNakedAmount := float64(0)
	totalQuantity := 0

	outData.AllowOrder = 1
	for _, item := range *cartProducts {
		// 库存不足 不允许下单
		if outData.CheckStockStatus == 0 {
			// 库存不足
			if item.ShowCartProduct == 1 && item.SaleStatus == 0 {
				outData.AllowOrder = 0
				outData.AllowOrderMsg = "库存不足"
			}
		}
		if item.Selected == 1 && item.ShowCartProduct == 1 {
			totalAmount = utils.AddFloat64Fixed(totalAmount, item.TotalMarketPrice, true)
			totalNakedAmount = utils.AddFloat64Fixed(totalNakedAmount, item.TotalNakedSalePrice, true)
			totalQuantity = totalQuantity + item.Quantity
		}
	}

	outData.TotalQuantity = totalQuantity
	outData.TotalAmount = totalAmount
	outData.TotalNakedAmount = totalNakedAmount
}

func (e *UserCart) GetOrderProductForBuyNow(c *dto.UserCartGetProductForOrderBuyNowReq, p *actions.DataPermission, outData *dto.UserCartGetProductForOrderBuyNowResp) error {
	goodsResult, err := e.GetBuyNowProductForOrder(c.UserId, c.GoodsId, c.Quantity, c.WarehouseCode)
	if err != nil {
		return err
	}
	e.FormatProductOrderBuyNow(e.Orm, &goodsResult, outData)
	outData.Products = goodsResult
	return nil
}

func (e *UserCart) FormatProductOrderBuyNow(tx *gorm.DB, cartProducts *[]models.UserCartGoodsProduct, outData *dto.UserCartGetProductForOrderBuyNowResp) {
	totalAmount := float64(0)
	totalNakedAmount := float64(0)
	skuCodes := []string{}
	productNum := 0

	for _, item := range *cartProducts {
		productNum += item.Quantity
		totalAmount = utils.AddFloat64Fixed(totalAmount, item.TotalMarketPrice, true)
		totalNakedAmount = utils.AddFloat64Fixed(totalNakedAmount, item.TotalNakedSalePrice, true)
		skuCodes = append(skuCodes, item.SkuCode)
	}
	productCategory := models.ProductCategory{}
	outData.ProductCategoryNum = productCategory.GetMainCategoryCountBySkus(tx, skuCodes)
	outData.ProductNum = productNum
	outData.TotalAmount = totalAmount
	outData.TotalNakedAmount = totalNakedAmount
}

func (e *UserCart) SaleMoq(c *dto.UserCartGetProductForSaleMoqReq, p *actions.DataPermission, outData *dto.UserCartGetProductForSaleMoqResp) error {
	outData.SaleMoq = map[string]int{}
	for _, item := range c.SkuCodes {
		saleMoq, err := e.GetCartSaleMoqBySku(item, c.WarehouseCode)
		if err != nil {
			return err
		}
		outData.SaleMoq[item] = saleMoq
	}
	return nil
}

func (e *UserCart) GetCartSaleMoqBySku(skuCode, warehouseCode string) (int, error) {
	goods := models.CheckBySkuAndWarehouseCode(e.Orm, skuCode, warehouseCode)
	if goods.Id == 0 {
		return 0, errors.New(skuCode + ",商品在当前仓库中，已下架或不存在")
	}
	userCartGoodsProduct := &models.UserCartGoodsProduct{}
	goodsResult, _ := userCartGoodsProduct.GetByGoodsIds(e.Orm, []int{goods.Id}, "")
	if len(*goodsResult) == 0 {
		return 0, errors.New("未找到该商品")
	}
	return (*goodsResult)[0].SalesMoq, nil
}
