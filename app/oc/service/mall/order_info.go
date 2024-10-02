package mall

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"github.com/jinzhu/copier"
	modelsAdmin "go-admin/app/admin/models"
	serviceAdmin "go-admin/app/admin/service"
	modelsPc "go-admin/app/pc/models"
	dtoPc "go-admin/app/pc/service/admin/dto"
	servicePc "go-admin/app/pc/service/mall"
	dtoMallPc "go-admin/app/pc/service/mall/dto"
	modelsUc "go-admin/app/uc/models"
	modelsWc "go-admin/app/wc/models"
	dtoWc "go-admin/app/wc/service/admin/dto"
	"go-admin/common"
	pcClient "go-admin/common/client/pc"
	wcClient "go-admin/common/client/wc"
	"go-admin/common/global"
	"go-admin/common/middleware/mall_handler"
	"go-admin/common/msg/email"
	"go-admin/common/utils"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/oc/models"
	"go-admin/app/oc/service/mall/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type OrderInfo struct {
	service.Service
}

// GetPage 获取OrderInfo列表
func (e *OrderInfo) GetPage(c *dto.OrderInfoGetPageReq, p *actions.DataPermission, list *[]dto.OrderInfoGetPageResp, count *int64) error {
	var err error
	var data models.OrderInfo

	err = e.Orm.Model(&data).Preload("OrderDetails").
		//Select("order_info.*").
		Joins("left join order_detail od on order_info.order_id = od.order_id").
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
			dto.GetPageMakeCondition(c, e.Orm),
		).
		Group("order_info.order_id").
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("OrderInfoService GetPage error:%s \r\n", err)
		return err
	}

	tmpList := *list
	var vendorIds []string
	var skus []string
	var warehouseCodes []string
	for _, orderInfo := range tmpList {
		warehouseCodes = append(warehouseCodes, orderInfo.WarehouseCode)
		for _, detail := range orderInfo.OrderDetails {
			vendorIds = append(vendorIds, strconv.Itoa(detail.VendorId))
			skus = append(skus, detail.SkuCode)
		}
	}
	// 仓库
	warehouseMap := getWarehouseMap(e.Orm, warehouseCodes)
	// 税率
	var modelProduct modelsPc.Product
	taxMap := modelProduct.GetTaxBySku(e.Orm, skus)
	currentTax := utils.GetCurrentTaxRate()
	// 货主map
	vendorMap := getVendorMap(e.Orm, vendorIds)
	for i, orderInfo := range tmpList {
		if _, ok := warehouseMap[orderInfo.WarehouseCode]; ok {
			tmpList[i].WarehouseName = warehouseMap[orderInfo.WarehouseCode].WarehouseName
		}
		tmpList[i].ButtonList = e.getButtonList(e.Orm, &orderInfo)
		for i2, detail := range orderInfo.OrderDetails {
			tax := currentTax
			if _, ok := vendorMap[detail.VendorId]; ok {
				tmpList[i].OrderDetails[i2].VendorName = vendorMap[detail.VendorId].NameZh
			}
			// 含税单价
			tmpTax, ok3 := taxMap[detail.SkuCode]
			if ok3 {
				tax, _ = strconv.ParseFloat(tmpTax, 64)
			}
			tmpList[i].OrderDetails[i2].UnTaxPrice, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", utils.DivFloat64(detail.SalePrice, utils.AddFloat64(tax, 1.00))), 64)
		}
	}

	*list = tmpList
	return nil
}

func (e *OrderInfo) GetPageCount(c *dto.OrderInfoGetPageReq, p *actions.DataPermission, list *[]dto.OrderInfoGetPageResp) error {
	var err error
	var data models.OrderInfo

	err = e.Orm.Model(&data).Preload("OrderDetails").
		Joins("left join order_detail od on order_info.order_id = od.order_id").
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			actions.Permission(data.TableName(), p),
			dto.GetPageMakeCondition(c, e.Orm),
		).
		Group("order_info.order_id").
		Find(list).Error
	if err != nil {
		e.Log.Errorf("OrderInfoService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取OrderInfo对象
func (e *OrderInfo) Get(d *dto.OrderInfoGetReq, p *actions.DataPermission, model *dto.OrderInfoGetResp) error {
	var data models.OrderInfo

	err := e.Orm.Preload("OrderDetails").Preload("OrderImages", "type=0").Preload("ReceiptImages", models.OrderImage{Type: 1}).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		Where("order_info.order_id = ?", d.GetId()).
		First(model).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetOrderInfo error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	var vendorIds []string
	for _, detail := range model.OrderDetails {
		vendorIds = append(vendorIds, strconv.Itoa(detail.VendorId))
	}
	// 货主map
	vendorMap := getVendorMap(e.Orm, vendorIds)
	for i, detail := range model.OrderDetails {
		if vendor, ok := vendorMap[detail.VendorId]; ok {
			model.OrderDetails[i].VendorName = vendor.NameZh
		}
	}
	//model.AddressFullName = model.Consignee + ", " + model.Mobile + ", " + model.Address
	model.AddressFullName = model.Consignee + ", " + model.Mobile
	if model.Telephone != "" {
		model.AddressFullName = model.AddressFullName + ", " + model.Telephone
	}
	model.AddressFullName = model.AddressFullName + ", " + model.ProvinceName + " " + model.CityName + " " + model.AreaName + ", " + model.Address
	model.ButtonList = e.getViewButtonList(e.Orm, model)

	return nil
}

// UpdatePo 签收领用单
func (e *OrderInfo) UpdatePo(c *dto.OrderInfoUpdatePoReq, p *actions.DataPermission) error {
	var err error
	var order = models.OrderInfo{}
	e.Orm.Scopes(
		actions.Permission(order.TableName(), p),
	).Where("order_info.order_id = ?", c.GetId()).First(&order)
	oldData := order

	c.Generate(&order)
	db := e.Orm.Save(&order)
	if err = db.Error; err != nil {
		e.Log.Errorf("OrderInfoService Save error:%s \r\n", err)
		return err
	}

	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	e.createLog(c, oldData, order, models.LogTypeOrderInfoUpdateInfo)
	return nil
}

// Cancel 取消领用单
func (e *OrderInfo) Cancel(c *dto.OrderInfoCancelReq, p *actions.DataPermission) error {
	var err error
	var order = models.OrderInfo{}
	if c.Id == "" {
		return errors.New("参数有误")
	}
	e.Orm.Preload("OrderDetails").Scopes(
		actions.Permission(order.TableName(), p),
	).Where("order_info.order_id = ?", c.GetId()).First(&order)
	err = c.CancelValid(e.Orm)
	if err != nil {
		return err
	}
	oldData := order

	tx := e.Orm.Debug().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	c.Generate(&order)
	order.OrderStatus = 9
	order.CancelByType = 0
	db := tx.Updates(&order)
	if err = db.Error; err != nil {
		e.Log.Errorf("OrderInfoService Save error:%s \r\n", err)
		tx.Rollback()
		return err
	}
	var goodsIds []int
	for _, detail := range order.OrderDetails {
		goodsIds = append(goodsIds, detail.GoodsId)
	}
	stockMap := getStockMap(e.Orm, goodsIds, order.WarehouseCode)
	// 释放库存
	for _, detail := range order.OrderDetails {
		if _, ok := stockMap[detail.GoodsId]; !ok || detail.LockStock <= 0 {
			continue
		}
		err = modelsWc.UnLockStockInfoForOrder(e.Orm, detail.LockStock, detail.GoodsId, order.WarehouseCode, order.OrderId, "后台取消订单释放锁库")
		if err != nil {
			tx.Rollback()
			return err
		}
		detail.LockStock = 0
		err = tx.Updates(&detail).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	// 释放预算
	var departmentBudgetM = modelsUc.DepartmentBudget{}
	err = departmentBudgetM.UpdateBudget(e.Orm, order.UserId, -order.TotalAmount, order.CreatedAt.Format("200601"))
	if err != nil {
		tx.Rollback()
		return errors.New(fmt.Sprintf("预算修改失败，\r\n失败信息 %s", err.Error()))
	}

	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	tx.Commit()
	e.createLog(c, oldData, order, models.LogTypeOrderInfoCancel)
	return nil
}

// BuyAgain 再次购买
func (e *OrderInfo) BuyAgain(c *dto.OrderInfoBuyAgainReq, p *actions.DataPermission) (errList []string, err error) {
	var order = models.OrderInfo{}
	if c.Id == "" {
		return nil, errors.New("参数有误")
	}
	e.Orm.Preload("OrderDetails").Scopes(
		actions.Permission(order.TableName(), p),
	).Where("order_info.order_id = ?", c.Id).First(&order)

	if len(order.OrderDetails) <= 0 {
		return nil, errors.New("订单无商品")
	}
	ctx := e.Orm.Statement.Context.(*gin.Context)
	if order.WarehouseCode != mall_handler.GetUserCurrentWarehouseCode(ctx) {
		return nil, errors.New("该订单所属仓库不在当前所选仓库下，无法再次购买")
	}
	userCartService := servicePc.UserCart{Service: e.Service}
	// 添加到购物车
	for _, detail := range order.OrderDetails {
		data := dtoMallPc.UserCartInsertReq{
			GoodsId:       detail.GoodsId,
			Quantity:      detail.Quantity,
			WarehouseCode: order.WarehouseCode,
		}
		data.CreateBy = user.GetUserId(ctx)
		data.CreateByName = user.GetUserName(ctx)
		err2 := userCartService.Insert(&data)
		if err2 != nil {
			errList = append(errList, fmt.Sprintf("[%s] %s", detail.SkuCode, err2.Error()))
		}
	}
	if len(errList) > 0 {
		return errList, nil
	}

	return nil, nil
}

// GetExport 获取导出数据
func (e *OrderInfo) GetExport(c *dto.OrderInfoGetPageReq, p *actions.DataPermission, list *[]dto.OrderInfoGetExportData) error {
	var err error
	var data models.OrderInfo

	var resList []struct {
		models.OrderInfo
		models.OrderDetail
	}

	err = e.Orm.Table("order_info").
		Select("order_info.*,od.*").
		Joins("left join order_detail od on order_info.order_id = od.order_id").
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			actions.Permission(data.TableName(), p),
			dto.GetPageMakeCondition(c, e.Orm),
		).
		Find(&resList).Error
	if err != nil {
		return err
	}
	tmpList := *list
	var skus []string
	var warehouseCodes []string
	var vendorIds []string
	for _, v := range resList {
		warehouseCodes = append(warehouseCodes, v.OrderInfo.WarehouseCode)
		skus = append(skus, v.SkuCode)
		vendorIds = append(vendorIds, strconv.Itoa(v.VendorId))
	}
	//  字段值转中文
	var modelDictData modelsAdmin.SysDictData
	dictDataList := modelDictData.GetDictDataByTypes(e.Orm, []string{"order_info_receive_status", "order_info_send_status"})
	// 税率
	//var modelProduct modelsPc.Product
	//taxMap := modelProduct.GetTaxBySku(e.Orm, skus)
	//tax := utils.GetCurrentTaxRate()
	// 仓库
	warehouseMap := getWarehouseMap(e.Orm, warehouseCodes)
	// 货主
	vendorMap := getVendorMap(e.Orm, vendorIds)
	// 商品map
	productMap := getProductMap(e.Orm, skus)
	for _, v := range resList {
		var tmp dto.OrderInfoGetExportData
		_ = copier.Copy(&tmp, &v)
		tmp.OrderId = v.OrderInfo.OrderId
		tmp.UserName = v.OrderInfo.UserName

		if vendor, ok := vendorMap[v.VendorId]; ok {
			tmp.VendorName = vendor.NameZh
		}
		if product, ok := productMap[v.SkuCode]; ok {
			tmp.SupplierSkuCode = product.SupplierSkuCode
		}
		// 仓库名称
		if _, ok := warehouseMap[v.OrderInfo.WarehouseCode]; ok {
			tmp.WarehouseName = warehouseMap[v.OrderInfo.WarehouseCode].WarehouseName
		}
		//发货状态
		if text, ok := dictDataList["order_info_send_status"][strconv.Itoa(v.SendStatus)]; ok {
			tmp.SendStatusText = text
		}
		//收货状态
		if text, ok := dictDataList["order_info_receive_status"][strconv.Itoa(v.ReceiveStatus)]; ok {
			tmp.ReceiveStatusText = text
		}
		tmp.CreatedTime = utils.TimeFormat(v.OrderInfo.CreatedAt)

		var csApply models.CsApply
		returnListMap := map[string]map[string]*models.GetAfterReturnProductsBySaleIdResult{}
		//已退货数
		returnList, ok2 := returnListMap[v.OrderInfo.OrderId]
		if !ok2 {
			returnList = csApply.GetAfterReturnProductsBySaleId(e.Orm, v.OrderInfo.OrderId)
			returnListMap[v.OrderInfo.OrderId] = returnList

		}
		if _, okk2 := returnList[v.SkuCode]; okk2 {
			tmp.ReturnQuantity = returnList[v.SkuCode].Quantity
		}
		// 含税单价
		//tmpTax, ok3 :=  taxMap[v.SkuCode]
		//if ok3 {
		//	tax, _ = strconv.ParseFloat(tmpTax, 64)
		//}
		tmp.SalePrice = v.SalePrice
		tmp.TotalPrice, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", utils.MulFloat64AndInt(v.SalePrice, v.Quantity)), 64)
		tmpList = append(tmpList, tmp)
	}
	*list = tmpList
	return nil
}

// Delete
func (e *OrderInfo) Delete(c *dto.OrderInfoDeleteReq, p *actions.DataPermission) error {
	var err error
	var data = models.OrderInfo{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).Where("order_info.order_id = ?", c.Id).First(&data)
	err = c.Valid(data)
	if err != nil {
		return err
	}
	oldData := data

	data.ValidFlag = 0
	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("OrderInfoService Save error:%s \r\n", err)
		return err
	}

	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	e.createLog(c, oldData, data, models.LogTypeOrderInfoDelete)
	return nil
}

// Insert 创建OrderInfo对象
func (e *OrderInfo) Insert(c *dto.PresoSubmitApprovalReq, order *models.OrderInfo) error {
	var err error
	ctx := e.Orm.Statement.Context.(*gin.Context)
	userId := user.GetUserId(ctx)
	userName := user.GetUserName(ctx)
	warehouseCode := mall_handler.GetUserConfig(ctx).SelectedWarehouseCode

	//var order models.OrderInfo
	var shippingAddress modelsUc.Address
	// 购物车
	userCartService := servicePc.UserCart{Service: e.Service}
	var productList []modelsPc.UserCartGoodsProduct
	if c.BuyNow == 1 {
		productList, err = userCartService.GetBuyNowProductForOrder(userId, c.GoodsId, c.Quantity, warehouseCode)
	} else {
		productList, err = userCartService.GetCartProductForOrder(userId, warehouseCode, false)
	}

	var companyInfoM modelsUc.CompanyInfo
	companyInfo := companyInfoM.GetRowByUserId(e.Orm, user.GetUserId(ctx))
	if companyInfo.IsEas != 0 {
		return errors.New("用户所属公司需提交审批，请刷新页面重新尝试")
	}
	err = c.Valid(e.Orm, &shippingAddress, productList, companyInfo)
	if err != nil {
		return err
	}

	order.CreateFrom = c.CreateFrom
	order.ContractNo = c.ContractNo
	order.Remark = c.Remark
	order.OrderId = dto.GenerateOrderId()
	order.Ip = common.GetClientIP(ctx)
	order.UserId = userId
	order.UserName = userName
	order.UserCompanyId = companyInfo.Id
	order.UserCompanyName = companyInfo.CompanyName
	order.CreateBy = userId
	order.CreateByName = userName
	order.WarehouseCode = warehouseCode
	//地址相关
	order.DeliverId = shippingAddress.Id
	order.CountryId = shippingAddress.CountryId
	order.CountryName = shippingAddress.CountryName
	order.ProvinceId = shippingAddress.ProvinceId
	order.ProvinceName = shippingAddress.ProvinceName
	order.CityId = shippingAddress.CityId
	order.CityName = shippingAddress.CityName
	order.AreaId = shippingAddress.AreaId
	order.AreaName = shippingAddress.AreaName
	order.TownId = shippingAddress.TownId
	order.TownName = shippingAddress.TownName
	order.Mobile = shippingAddress.CellPhone
	order.Telephone = shippingAddress.Telephone
	order.Address = shippingAddress.DetailAddress
	order.CompanyName = shippingAddress.CompanyName
	order.Consignee = shippingAddress.ReceiverName
	order.ContactEmail = shippingAddress.Email

	order.IsOverBudget = -1
	order.OrderStatus = 5

	var goodsIds []int
	var skus []string
	productQuantity := 0
	totalAmount := 0.00
	for _, product := range productList {
		goodsIds = append(goodsIds, product.GoodsId)
		skus = append(skus, product.SkuCode)
		productQuantity = productQuantity + product.Quantity
		totalAmount = utils.AddFloat64(totalAmount, utils.MulFloat64AndInt(product.MarketPrice, product.Quantity))
	}
	// 遍历商品 判断库存是否完全满足
	// 全部满足：订单状态=待确认 /  不完全满足，订单状态=缺货
	stockMap := getStockMap(e.Orm, goodsIds, warehouseCode)

	for _, product := range productList {
		stock, ok := stockMap[product.GoodsId]
		if !ok || stock < product.Quantity {
			order.OrderStatus = 6
			break
		}
	}

	order.IsFormalOrder = 1
	order.ValidFlag = 1
	order.ClassifyQuantity = len(productList) // 商品校验了必须有产线，且每个商品主产线只有1个，故 商品种类数量为sku个数
	order.ProductQuantity = productQuantity
	// 无运费
	order.ItemsAmount = totalAmount
	order.TotalAmount = totalAmount
	order.OriginalItemsAmount = totalAmount
	order.OriginalTotalAmount = totalAmount

	tx := e.Orm.Debug().Begin()
	err = tx.Create(&order).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// 预订单明细
	var userProductRemark map[string]dto.PresoUserProductRemarkReq
	if len(c.UserProductRemark) > 0 {
		_ = utils.StructColumn(&userProductRemark, c.UserProductRemark, "SkuCode", "")
	}
	categoryMap := getCategoryMap(e.Orm, skus)
	productMap := getProductMap(e.Orm, skus)
	for _, detail := range productList {
		var orderDetail models.OrderDetail
		orderDetail.OrderId = order.OrderId
		orderDetail.UserId = order.UserId
		orderDetail.UserName = order.UserName
		orderDetail.WarehouseCode = order.WarehouseCode
		orderDetail.SkuCode = detail.SkuCode
		orderDetail.VendorId = detail.VendorId
		orderDetail.GoodsId = detail.GoodsId
		orderDetail.Quantity = detail.Quantity
		orderDetail.MarketPrice = detail.MarketPrice
		orderDetail.SalePrice = detail.MarketPrice
		orderDetail.ProductNo = detail.ProductNo
		orderDetail.ProductName = detail.NameZh
		orderDetail.ProductPic = detail.Image.MediaDir
		if _, ok := userProductRemark[detail.SkuCode]; ok {
			orderDetail.UserProductRemark = userProductRemark[detail.SkuCode].Remark
		}
		orderDetail.FinalQuantity = orderDetail.Quantity
		orderDetail.OriginalQuantity = orderDetail.Quantity
		orderDetail.OriginalItemsMount = utils.MulFloat64AndInt(orderDetail.SalePrice, orderDetail.Quantity)
		orderDetail.SubTotalAmount = orderDetail.OriginalItemsMount
		// 商品
		setDetailProduct(&orderDetail, productMap)
		// 产线
		setDetailCategory(&orderDetail, categoryMap)

		// 锁库
		lockStock := 0
		if _, ok := stockMap[orderDetail.GoodsId]; ok && stockMap[orderDetail.GoodsId] > 0 {
			lockStock = orderDetail.Quantity
			if stockMap[orderDetail.GoodsId] < orderDetail.Quantity {
				lockStock = stockMap[orderDetail.GoodsId]
			}
			err = modelsWc.LockStockInfoForOrder(e.Orm, lockStock, orderDetail.GoodsId, order.WarehouseCode, order.OrderId, "商城创建订单锁库")
			if err != nil {
				tx.Rollback()
				return err
			}
		}
		orderDetail.LockStock = lockStock
		err = tx.Create(&orderDetail).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		// 预算
		var departmentBudgetM = modelsUc.DepartmentBudget{}
		err = departmentBudgetM.UpdateBudget(e.Orm, orderDetail.UserId, orderDetail.SubTotalAmount, time.Now().Format("200601"))
		if err != nil {
			tx.Rollback()
			return errors.New(fmt.Sprintf("预算修改失败，\r\n失败信息 %s", err.Error()))
		}

		order.OrderDetails = append(order.OrderDetails, orderDetail)
	}

	// 清空购物车
	if c.BuyNow != 1 {
		data := dtoMallPc.UserCartClearSelectReq{
			WarehouseCode: warehouseCode,
			UserId:        userId,
		}
		err = userCartService.ClearSelect(&data, nil)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	e.createLog(c, models.OrderInfo{}, *order, global.LogTypeCreate)
	// 生成正式领用单邮件通知仓管员
	_ = e.NoticeSystemUser(*order)

	return nil
}

// Receipt 签收领用单
func (e *OrderInfo) Receipt(c *dto.OrderInfoReceiptReq, p *actions.DataPermission) error {
	var err error
	var order = models.OrderInfo{}
	e.Orm.Scopes(
		actions.Permission(order.TableName(), p),
	).Where("order_info.order_id = ?", c.Id).First(&order)
	err = c.ReceiptValid(&order)
	if err != nil {
		return err
	}
	oldData := order

	order.OrderStatus = 7
	order.ReceiveStatus = 2
	order.ConfirmOrderReceiptTime = time.Now()
	db := e.Orm.Updates(&order)
	if err = db.Error; err != nil {
		e.Log.Errorf("OrderInfoService Save error:%s \r\n", err)
		return err
	}

	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	e.createLog(c, oldData, order, models.LogTypeOrderInfoReceipt)
	return nil
}

// 获取我的订单列表中可展示的按钮
func (e *OrderInfo) getButtonList(tx *gorm.DB, order *dto.OrderInfoGetPageResp) (res []int) {
	res = append(res, 1) // 查看详情
	var modelCsApply models.CsApply
	if csApplys := modelCsApply.GetRowsByOrderId(tx, order.OrderId); len(csApplys) > 0 {
		res = append(res, 2) // 查看售后
	}
	res = append(res, 3) // 再次购买
	res = append(res, 5) // 填写po单号
	switch order.OrderStatus {
	case 9: //已取消
		res = append(res, 4) // 删除订单
		break
	case 5, 6:
		res = append(res, 6) // 申请取消
		break
	case 1:
		var modelOrderDetail models.OrderDetail
		tmpMap := modelOrderDetail.GetReturnQuantiryMap(tx, []string{order.OrderId})
		if tmpStruct, ok := tmpMap[order.OrderId]; ok && tmpStruct.ReturnQuantiry > 0 && (order.RmaStatus == 0 || order.RmaStatus == 99) {
			res = append(res, 7) // 签收
		}
		break
	}
	return
}

// 获取我的订单详情中可展示的按钮
func (e *OrderInfo) getViewButtonList(tx *gorm.DB, order *dto.OrderInfoGetResp) (res []int) {
	res = append(res, 3) // 再次购买
	switch order.OrderStatus {
	case 5, 6:
		res = append(res, 6) // 申请取消
		break
	case 1:
		var modelOrderDetail models.OrderDetail
		tmpMap := modelOrderDetail.GetReturnQuantiryMap(tx, []string{order.OrderId})
		if tmpStruct, ok := tmpMap[order.OrderId]; ok && tmpStruct.ReturnQuantiry > 0 && (order.RmaStatus == 0 || order.RmaStatus == 99) {
			res = append(res, 7) // 签收
		}
		break
	}
	return
}

// 获取库存map goodsid => stock
func getStockMap(tx *gorm.DB, goodsIds []int, warehouseCode string) (stockMap map[int]int) {
	stockResult := wcClient.ApiByDbContext(tx).GetStockListByGoodsIdAndWarehouseCode(dtoWc.InnerStockInfoGetByGoodsIdAndWarehouseCodeReq{
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

// 获取库存map sku => category
func getCategoryMap(tx *gorm.DB, skus []string) (categoryMap map[string][]modelsPc.Category) {
	categoryResult := pcClient.ApiByDbContext(tx).GetCategoryBySku(skus)
	categoryResultInfo := &struct {
		response.Response
		Data []dtoPc.InnerGetProductCategoryBySkuResp
	}{}
	categoryResult.Scan(categoryResultInfo)
	categoryMap = make(map[string][]modelsPc.Category)
	for _, category := range categoryResultInfo.Data {
		categoryMap[category.SkuCode] = category.ProductCategory
	}
	return
}

// 货主map
func getVendorMap(tx *gorm.DB, vendorIds []string) (vendorMap map[int]modelsWc.Vendors) {
	vendorResult := wcClient.ApiByDbContext(tx).GetVendorList(dtoWc.InnerVendorsGetListReq{
		Ids:    strings.Join(vendorIds, ","),
		Status: "1",
	})
	vendorResultInfo := &struct {
		response.Response
		Data []modelsWc.Vendors
	}{}
	vendorResult.Scan(vendorResultInfo)
	vendorMap = make(map[int]modelsWc.Vendors)
	for _, vendor := range vendorResultInfo.Data {
		vendorMap[vendor.Id] = vendor
	}
	return
}

// 仓库map
func getWarehouseMap(tx *gorm.DB, warehouseCodes []string) (warehouseMap map[string]modelsWc.Warehouse) {
	warehouseResult := wcClient.ApiByDbContext(tx).GetWarehouseList(dtoWc.InnerWarehouseGetListReq{
		Status:        "1",
		WarehouseCode: strings.Join(warehouseCodes, ","),
	})
	warehouseResultInfo := &struct {
		response.Response
		Data []modelsWc.Warehouse
	}{}
	warehouseResult.Scan(warehouseResultInfo)
	warehouseMap = make(map[string]modelsWc.Warehouse, len(warehouseResultInfo.Data))
	for _, warehouse := range warehouseResultInfo.Data {
		warehouseMap[warehouse.WarehouseCode] = warehouse
	}
	return
}

// 生成日志
func (e *OrderInfo) createLog(req interface{}, beforeModel models.OrderInfo, afterModel models.OrderInfo, logType string) {
	dataLog, _ := json.Marshal(&req)
	beforeDataStr := []byte("")
	if !reflect.DeepEqual(beforeModel, models.OrderInfo{}) {
		beforeDataStr, _ = json.Marshal(&beforeModel)
	}
	afterDataStr, _ := json.Marshal(&afterModel)
	sysUserLog := models.OrderInfoLog{
		DataId:     afterModel.Id,
		Type:       logType,
		Data:       string(dataLog),
		BeforeData: string(beforeDataStr),
		AfterData:  string(afterDataStr),
	}
	_ = sysUserLog.CreateLog("orderInfo", e.Orm)
}

func (e *OrderInfo) NoticeSystemUser(order models.OrderInfo) error {
	adminService := serviceAdmin.SysUser{Service: e.Service}
	sysUsers, err := adminService.GetAuthorityWarehouseUser(e.Orm, order.WarehouseCode)
	if err != nil {
		return err
	}
	warehouse, err2 := e.GetWarehouseByCode(order.WarehouseCode)
	if err2 != nil {
		return err2
	}

	logo := "https://image-c.ehsy.com/uploadfile/sxyz/img/2022/09/16/20220916100258386.png"
	// 获取后台订单详情地址
	host := utils.GetHostUrl(1)
	approveUrl := host + "/sale/order/index?orderId=" + order.OrderId

	// 获取模板路径
	tmpPath := "static/viewtpl/oc/ensure_approve.html"

	// 获取邮件标题
	subject := fmt.Sprintf("【狮行驿站】领用单待确认通知(%s)", order.OrderId)

	body := ""
	var to []string
	data := models.EnsureApprovalEmail{}
	for _, item := range *sysUsers {
		// 组合内容
		data = models.EnsureApprovalEmail{
			UserName:      item.Username,
			OrderId:       order.OrderId,
			ApproveUrl:    approveUrl,
			WarehouseName: warehouse.WarehouseName,
			LogoUrl:       logo,
			OrderDetails:  order.OrderDetails,
		}
		// 模板赋值
		body, err = utils.View(tmpPath, data)
		if err != nil {
			return err
		}
		// 发送邮件
		to = []string{item.Email}
		err = email.AsyncSendEmails(to, subject, body)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *OrderInfo) GetWarehouseByCode(warehouseCode string) (warehouse *modelsWc.Warehouse, err error) {
	warehouse = &modelsWc.Warehouse{}
	wcPrefix := global.GetTenantWcDBNameWithDB(e.Orm)
	err = e.Orm.Table(wcPrefix+".warehouse").
		Where("warehouse_code = ?", warehouseCode).
		Find(warehouse).Error

	return warehouse, err
}
