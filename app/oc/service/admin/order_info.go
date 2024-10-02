package admin

import (
	"encoding/json"
	"errors"
	"fmt"
	modelsPc "go-admin/app/pc/models"
	dtoPc "go-admin/app/pc/service/admin/dto"
	modelsUc "go-admin/app/uc/models"
	dtoUc "go-admin/app/uc/service/admin/dto"
	modelsWc "go-admin/app/wc/models"
	serviceWc "go-admin/app/wc/service/admin"
	dtoWc "go-admin/app/wc/service/admin/dto"
	"go-admin/common"
	pcClient "go-admin/common/client/pc"
	ucClient "go-admin/common/client/uc"
	wcClient "go-admin/common/client/wc"
	"go-admin/common/global"
	"go-admin/common/utils"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/oc/models"
	"go-admin/app/oc/service/admin/dto"
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

	wcPrefix := global.GetTenantWcDBNameWithDB(e.Orm)
	err = e.Orm.Model(&data).
		Select("order_info.*, w.warehouse_name").
		Joins("left join "+wcPrefix+".warehouse w on order_info.warehouse_code = w.warehouse_code").
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
			dto.GetPageMakeCondition(c, e.Orm),
			actions.SysUserPermission(data.TableName(), p, 2),          // 仓库权限
			dto.OrderInfoGetPageCompanyPermission(data.TableName(), p), // 公司权限
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("OrderInfoService GetPage error:%s \r\n", err)
		return err
	}
	var orderIds []string
	var companyByIds []string
	for _, v := range *list {
		orderIds = append(orderIds, v.OrderId)
		companyByIds = append(companyByIds, strconv.Itoa(v.UserCompanyId))
	}

	// 仓库map
	warehouseResult := wcClient.ApiByDbContext(e.Orm).GetWarehouseList(dtoWc.InnerWarehouseGetListReq{
		Status: "1",
	})
	warehouseResultInfo := &struct {
		response.Response
		Data []modelsWc.Warehouse
	}{}
	warehouseResult.Scan(warehouseResultInfo)
	warehouseMap := make(map[string]string, len(warehouseResultInfo.Data))
	for _, warehouse := range warehouseResultInfo.Data {
		warehouseMap[warehouse.WarehouseCode] = warehouse.WarehouseName
	}

	if len(orderIds) > 0 {
		// 获取公司信息
		companyMap := e.GetCompanyByIds(companyByIds)

		type tmpStruct struct {
			OrderId        string `json:"orderId"`
			ReturnQuantiry int    `json:"returnQuantiry"`
		}
		var tmp []tmpStruct
		err = e.Orm.Raw(`
			SELECT order_id, sum(final_quantity) return_quantiry
			FROM order_detail as t
			WHERE order_id IN ? GROUP BY order_id
		`, orderIds).Scan(&tmp).Error
		if err != nil {
			return err
		}
		var tmpMap map[string]tmpStruct
		_ = utils.StructColumn(&tmpMap, tmp, "", "OrderId")
		tmpList := *list
		for i, v := range tmpList {
			if _, ok := tmpMap[v.OrderId]; ok {
				tmpList[i].ReturnQuantiry = tmpMap[v.OrderId].ReturnQuantiry
			}
			if _, ok := companyMap[v.UserCompanyId]; ok {
				if companyMap[v.UserCompanyId].OrderAutoSignFor == 1 {
					tmpList[i].IsUploadReceipt = 1
				}
			}
			orderIds = append(orderIds, v.OrderId)
		}
		*list = tmpList
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
		First(model, d.GetId()).Error
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
	var skus []string
	var goodsIds []int
	for _, detail := range model.OrderDetails {
		vendorIds = append(vendorIds, strconv.Itoa(detail.VendorId))
		skus = append(skus, detail.SkuCode)
		goodsIds = append(goodsIds, detail.GoodsId)
	}
	// 税率
	var modelProduct modelsPc.Product
	taxMap := modelProduct.GetTaxBySku(e.Orm, skus)
	currentTax := utils.GetCurrentTaxRate()
	// 货主map
	vendorMap := getVendorMap(e.Orm, vendorIds)
	// 库存map
	stockMap := getStockMap(e.Orm, goodsIds, model.WarehouseCode)
	// 商品map
	productMap := getProductMap(e.Orm, skus, model.WarehouseCode)
	var csApply models.CsApply
	afterList := csApply.GetAfterApplyProductsByTypeBySaleId(e.Orm, model.OrderId)
	returnList := csApply.GetAfterReturnProductsBySaleId(e.Orm, model.OrderId)
	for i, detail := range model.OrderDetails {
		model.OrderDetails[i].ActualStock = detail.Quantity - detail.CancelQuantity - detail.LockStock
		//售后数
		if _, ok := afterList[detail.SkuCode]; ok {
			afterNum := 0
			for i2, v2 := range afterList[detail.SkuCode] {
				if i2 == 0 || i2 == 9 {
					afterNum = afterNum + v2.Quantity
				}
			}
			model.OrderDetails[i].AfterNum = afterNum
		}

		//已退货数
		if _, ok := returnList[detail.SkuCode]; ok {
			model.OrderDetails[i].ReturnNum = returnList[detail.SkuCode].Quantity
		}
		if vendor, ok := vendorMap[detail.VendorId]; ok {
			model.OrderDetails[i].VendorName = vendor.NameZh
			model.OrderDetails[i].VendorSkuCode = vendor.Code
		}
		tax := currentTax
		tmpTax, ok3 := taxMap[detail.SkuCode]
		if ok3 {
			tax, _ = strconv.ParseFloat(tmpTax, 64)
		}
		model.OrderDetails[i].Tax = tax
		if stock, ok := stockMap[detail.GoodsId]; ok {
			model.OrderDetails[i].Stock = stock
		}
		if product, ok := productMap[detail.SkuCode]; ok {
			model.OrderDetails[i].Moq = product.Product.SalesMoq
		}
	}

	if model.OrderStatus == 5 || model.OrderStatus == 6 {
		for i, detail := range model.OrderDetails {
			if product, ok := productMap[detail.SkuCode]; ok {
				model.OrderDetails[i].ProductNo = product.ProductNo
			}
		}
	}

	return nil
}

// Insert 创建OrderInfo对象
func (e *OrderInfo) Insert(c *dto.OrderInfoInsertReq) error {
	var err error
	var order models.OrderInfo
	err = c.InsertValid(e.Orm)
	if err != nil {
		return err
	}

	c.Generate(&order)

	var goodsIds []int
	var skus []string
	productQuantity := 0
	totalAmount := 0.00
	for _, product := range c.Products {
		goodsIds = append(goodsIds, product.GoodsId)
		skus = append(skus, product.SkuCode)
		productQuantity = productQuantity + product.Quantity
		totalAmount = utils.AddFloat64(totalAmount, utils.MulFloat64AndInt(product.SalePrice, product.Quantity))
	}

	// 遍历商品 判断库存是否完全满足
	// 全部满足：订单状态=待确认 /  不完全满足，订单状态=缺货
	stockMap := getStockMap(e.Orm, goodsIds, c.WarehouseCode)

	for _, product := range c.Products {
		stock, ok := stockMap[product.GoodsId]
		if !ok || stock < product.Quantity {
			order.OrderStatus = 6
			break
		}
	}
	order.CreateFrom = "LMS"
	order.Ip = common.GetClientIP(e.Orm.Statement.Context.(*gin.Context))
	order.IsFormalOrder = 1
	order.ValidFlag = 1
	order.ClassifyQuantity = len(c.Products) // 商品校验了必须有产线，且每个商品主产线只有1个，故 商品种类数量为sku个数
	order.ProductQuantity = productQuantity
	// 无运费
	order.ItemsAmount = totalAmount
	order.TotalAmount = totalAmount
	order.OriginalItemsAmount = totalAmount
	order.OriginalTotalAmount = totalAmount

	tx := e.Orm.Debug().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	err = tx.Create(&order).Error
	if err != nil {
		tx.Rollback()
		e.Log.Errorf("OrderInfoService Insert error:%s \r\n", err)
		return err
	}

	categoryMap := getCategoryMap(e.Orm, skus)
	// 商品
	for _, product := range c.Products {
		var detail models.OrderDetail
		product.Generate(&detail)
		detail.OrderId = order.OrderId
		detail.UserId = order.UserId
		detail.UserName = order.UserName
		detail.WarehouseCode = order.WarehouseCode
		// 产线
		setDetailCategory(&detail, categoryMap)
		// 锁库
		lockStock := 0
		if _, ok := stockMap[detail.GoodsId]; ok && stockMap[detail.GoodsId] > 0 {
			lockStock = detail.Quantity
			if stockMap[detail.GoodsId] < detail.Quantity {
				lockStock = stockMap[detail.GoodsId]
			}
			err = modelsWc.LockStockInfoForOrder(e.Orm, lockStock, detail.GoodsId, order.WarehouseCode, order.OrderId, "后台创建订单锁库")
			if err != nil {
				tx.Rollback()
				return err
			}
		}
		detail.LockStock = lockStock
		err = tx.Create(&detail).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		// 预算
		var departmentBudgetM = modelsUc.DepartmentBudget{}
		err = departmentBudgetM.UpdateBudget(e.Orm, detail.UserId, detail.SubTotalAmount, time.Now().Format("200601"))
		if err != nil {
			tx.Rollback()
			return errors.New(fmt.Sprintf("预算修改失败，\r\n失败信息 %s", err.Error()))
		}
	}

	// 相关文件
	for _, orderImg := range c.OrderImages {
		var orderImage models.OrderImage
		orderImg.Generate(&orderImage)
		orderImage.OrderId = order.OrderId
		err = tx.Create(&orderImage).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	e.createLog(c, models.OrderInfo{}, order, global.LogTypeCreate)
	return nil
}

// Update 修改OrderInfo对象
func (e *OrderInfo) Update(c *dto.OrderInfoUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.OrderInfo{}
	e.Orm.Preload("OrderImages", "type=0").
		Scopes(
			actions.Permission(data.TableName(), p),
		).First(&data, c.GetId())
	err = c.UpdateValid(e.Orm)
	if err != nil {
		return err
	}
	oldData := data

	tx := e.Orm.Debug().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	c.Generate(&data)
	db := tx.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("OrderInfoService Save error:%s \r\n", err)
		tx.Rollback()
		return err
	}

	// 相关文件
	var imgs []string
	var newImgs []string
	_ = utils.StructColumn(&imgs, data.OrderImages, "Url", "")
	_ = utils.StructColumn(&newImgs, c.OrderImages, "Url", "")
	for _, img := range data.OrderImages {
		if !utils.InArrayString(img.Url, newImgs) {
			err = tx.Unscoped().Delete(&img).Error
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	for _, image := range c.OrderImages {
		// 新图片
		if !utils.InArrayString(image.Url, imgs) {
			var orderImage models.OrderImage
			image.Generate(&orderImage)
			orderImage.OrderId = data.OrderId
			err = tx.Create(&orderImage).Error
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	tx.Commit()
	e.createLog(c, oldData, data, models.LogTypeOrderInfoUpdateInfo)
	return nil
}

// Receipt 签收领用单
func (e *OrderInfo) Receipt(c *dto.OrderInfoReceiptReq, p *actions.DataPermission) error {
	var err error
	var order = models.OrderInfo{}
	e.Orm.Scopes(
		actions.Permission(order.TableName(), p),
	).First(&order, c.GetId())
	err = c.ReceiptValid(e.Orm)
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
	order.OrderStatus = 7
	order.ReceiveStatus = 2
	order.ConfirmOrderReceiptTime = time.Now()
	db := tx.Updates(&order)
	if err = db.Error; err != nil {
		e.Log.Errorf("OrderInfoService Save error:%s \r\n", err)
		tx.Rollback()
		return err
	}

	// 签收文件
	for _, receiptImage := range c.ReceiptImages {
		var orderImage models.OrderImage
		receiptImage.Generate(&orderImage)
		orderImage.OrderId = order.OrderId
		orderImage.Type = 1
		err = tx.Create(&orderImage).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	tx.Commit()
	logType := models.LogTypeOrderInfoReceipt
	if c.IsAuto == 1 {
		logType = models.LogTypeOrderInfoAutoReceipt
	}
	e.createLog(c, oldData, order, logType)
	return nil
}

// GetReceiptImage 签收图片
func (e *OrderInfo) GetReceiptImage(c *dto.OrderInfoGetReceiptImageReq, res *dto.OrderInfoGetReceiptImageRes) error {

	var err error
	var data models.OrderImage

	res.List = &[]*models.OrderImage{}

	err = e.Orm.Model(&data).
		Where("order_id = ?", c.OrderId).
		Where("type = 2").
		Find(res.List).Limit(-1).Offset(-1).Error
	if err != nil {
		return err
	}

	return nil
}

// SaveReceiptImage 签收图片修改
func (e *OrderInfo) SaveReceiptImage(c *dto.OrderInfoSaveReceiptImageReq, p *actions.DataPermission) error {
	var err error
	var order = models.OrderInfo{}
	e.Orm.Scopes(
		actions.Permission(order.TableName(), p),
	).Where("order_id", c.OrderId).First(&order)
	if err != nil {
		return err
	}
	if len(c.ReceiptImages) > 6 {
		return errors.New("回单图片不能超过6张！")
	}
	if order.OrderStatus != 1 {
		return errors.New("订单状态不在已发货状态 不能上传回单！")
	}

	var data models.OrderImage
	list := &[]*models.OrderImage{}
	err = e.Orm.Model(&data).
		Where("order_id = ?", c.OrderId).
		Where("type = 2").
		Find(list).Error

	tx := e.Orm.Debug().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	orderImageIdMap := make(map[int]bool)
	// 签收文件
	for _, receiptImage := range c.ReceiptImages {
		if receiptImage.Id == 0 {
			var orderImage models.OrderImage
			receiptImage.Generate(&orderImage)
			orderImage.OrderId = order.OrderId
			orderImage.Type = 2
			err = tx.Save(&orderImage).Error
			if err != nil {
				tx.Rollback()
				return err
			}
		}
		orderImageIdMap[receiptImage.Id] = true
	}

	// 把不需要的删除掉
	listDelIds := []int{}
	for _, listData := range *list {
		if _, ok := orderImageIdMap[listData.Id]; !ok {
			listDelIds = append(listDelIds, listData.Id)
		}
	}
	if listDelIds != nil && len(listDelIds) > 0 {
		err = tx.Delete(data, listDelIds).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	e.createLog(c, models.OrderInfo{}, models.OrderInfo{}, models.LogTypeOrderInfoReceiptImageSave)
	return nil
}

// Cancel 取消领用单
func (e *OrderInfo) Cancel(c *dto.OrderInfoCancelReq, p *actions.DataPermission) error {
	var err error
	var order = models.OrderInfo{}
	e.Orm.Preload("OrderDetails").Scopes(
		actions.Permission(order.TableName(), p),
	).First(&order, c.GetId())
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
	order.CancelByType = 2
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

// Remove 删除OrderInfo
func (e *OrderInfo) Remove(d *dto.OrderInfoDeleteReq, p *actions.DataPermission) error {
	var data models.OrderInfo

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveOrderInfo error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}

// UpdateShipping 修改OrderInfo对象
func (e *OrderInfo) UpdateShipping(c *dto.OrderInfoUpdateShippingReq, p *actions.DataPermission) error {
	var err error
	var data = models.OrderInfo{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	err = c.UpdateValid(e.Orm)
	if err != nil {
		return err
	}
	oldData := data

	tx := e.Orm.Debug().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	c.Generate(&data)
	db := tx.Updates(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("OrderInfoService Save error:%s \r\n", err)
		tx.Rollback()
		return err
	}

	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	tx.Commit()
	e.createLog(c, oldData, data, models.LogTypeOrderInfoUpdateShipping)
	return nil
}

// UpdateProduct 修改OrderInfo对象
func (e *OrderInfo) UpdateProduct(c *dto.OrderInfoUpdateProductReq, p *actions.DataPermission) error {
	var err error
	var data = models.OrderInfo{}
	e.Orm.Preload("OrderDetails").Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	err = c.UpdateValid(e.Orm)
	if err != nil {
		return err
	}

	oldData := data

	tx := e.Orm.Debug().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var goodsIds []int
	productQuantity := 0
	totalAmount := 0.00
	for _, product := range c.Products {
		goodsIds = append(goodsIds, product.GoodsId)
		productQuantity = productQuantity + product.Quantity
		totalAmount = utils.AddFloat64(totalAmount, utils.MulFloat64AndInt(product.SalePrice, product.Quantity))
	}
	for _, detail := range data.OrderDetails {
		goodsIds = append(goodsIds, detail.GoodsId)
	}
	stockMap := getStockMap(e.Orm, goodsIds, c.WarehouseCode)

	// 商品
	var skus []string
	var newSkus []string
	var newDetailMap map[string]dto.OrderInfoProductInsertReq
	_ = utils.StructColumn(&skus, data.OrderDetails, "SkuCode", "")
	_ = utils.StructColumn(&newSkus, c.Products, "SkuCode", "")
	_ = utils.StructColumn(&newDetailMap, c.Products, "", "SkuCode")
	for i, detail := range data.OrderDetails {
		oriSubTotalAmount := detail.SubTotalAmount
		budgetMoney := 0.00
		// 商品被移除
		if !utils.InArrayString(detail.SkuCode, newSkus) {
			budgetMoney = 0 - oriSubTotalAmount
			err = tx.Unscoped().Delete(&detail).Error
			if err != nil {
				tx.Rollback()
				return err
			}
			// 恢复库存
			if detail.LockStock > 0 {
				err = modelsWc.UnLockStockInfoForOrder(e.Orm, detail.LockStock, detail.GoodsId, data.WarehouseCode, data.OrderId, "后台订单删除商品恢复库存")
				if err != nil {
					tx.Rollback()
					return err
				}
			}

			data.OrderDetails = append(data.OrderDetails[:i], data.OrderDetails[i+1:]...)
		} else {
			// 已存在商品 重置锁定库存 重置商品数量
			if newDetail, ok := newDetailMap[detail.SkuCode]; ok {
				unlock := newDetail.Quantity - detail.CancelQuantity - detail.LockStock
				// 减少数量 恢复库存
				if unlock < 0 {
					err = modelsWc.UnLockStockInfoForOrder(e.Orm, -unlock, detail.GoodsId, data.WarehouseCode, data.OrderId, "后台订单修改商品数量恢复库存")
					if err != nil {
						tx.Rollback()
						return err
					}
					detail.LockStock = detail.LockStock + unlock
				} else if unlock > 0 {
					// 锁库
					lockStock := unlock
					if _, ok := stockMap[detail.GoodsId]; ok && stockMap[detail.GoodsId] > 0 {
						if stockMap[detail.GoodsId] < unlock {
							lockStock = stockMap[detail.GoodsId]
						}
						err = modelsWc.LockStockInfoForOrder(e.Orm, lockStock, detail.GoodsId, data.WarehouseCode, data.OrderId, "后台订单修改商品锁库")
						if err != nil {
							tx.Rollback()
							return err
						}
						detail.LockStock = detail.LockStock + lockStock
					}
				}

				detail.OriginalQuantity = newDetail.Quantity
				detail.Quantity = newDetail.Quantity
				detail.FinalQuantity = newDetail.Quantity
				detail.SubTotalAmount = utils.MulFloat64AndInt(detail.SalePrice, detail.Quantity)

				budgetMoney = detail.SubTotalAmount - oriSubTotalAmount
				err = tx.Updates(&detail).Error
				if err != nil {
					tx.Rollback()
					return err
				}
			}
		}

		// 预算
		var departmentBudgetM = modelsUc.DepartmentBudget{}
		if budgetMoney != 0 {
			err = departmentBudgetM.UpdateBudget(e.Orm, detail.UserId, budgetMoney, time.Now().Format("200601"))
			if err != nil {
				tx.Rollback()
				return errors.New(fmt.Sprintf("预算修改失败，\r\n失败信息 %s", err.Error()))
			}
		}
	}

	categoryMap := getCategoryMap(e.Orm, newSkus)
	for _, product := range c.Products {
		// 新商品
		if !utils.InArrayString(product.SkuCode, skus) {
			var orderDetail models.OrderDetail
			product.Generate(&orderDetail)
			orderDetail.OrderId = data.OrderId
			orderDetail.UserId = data.UserId
			orderDetail.UserName = data.UserName
			orderDetail.WarehouseCode = data.WarehouseCode
			// 产线
			setDetailCategory(&orderDetail, categoryMap)
			// 锁库
			lockStock := 0
			if _, ok := stockMap[orderDetail.GoodsId]; ok && stockMap[orderDetail.GoodsId] > 0 {
				lockStock = orderDetail.Quantity
				if stockMap[orderDetail.GoodsId] < orderDetail.Quantity {
					lockStock = stockMap[orderDetail.GoodsId]
				}
				err = modelsWc.LockStockInfoForOrder(e.Orm, lockStock, orderDetail.GoodsId, data.WarehouseCode, data.OrderId, "后台修改订单添加商品锁库")
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
		}
	}

	data.OrderStatus = 5
	// 全部满足：订单状态=待确认 /  不完全满足，订单状态=缺货
	for _, product := range c.Products {
		stock, ok := stockMap[product.GoodsId]
		if !ok || stock < product.Quantity {
			data.OrderStatus = 6
			break
		}
	}

	c.Generate(&data)
	data.ClassifyQuantity = len(c.Products) // 商品校验了必须有产线，且每个商品主产线只有1个，故 商品种类数量为sku个数
	data.ProductQuantity = productQuantity
	// 无运费
	data.ItemsAmount = totalAmount
	data.TotalAmount = totalAmount
	data.OriginalItemsAmount = totalAmount
	data.OriginalTotalAmount = totalAmount
	db := tx.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("OrderInfoService Save error:%s \r\n", err)
		tx.Rollback()
		return err
	}

	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	tx.Commit()
	e.createLog(c, oldData, data, models.LogTypeOrderInfoUpdateProduct)
	return nil
}

// GetByOrderId 获取OrderInfo对象
func (e *OrderInfo) GetByOrderId(d *dto.OrderInfoByOrderIdReq, list *[]models.OrderInfo) error {

	err := e.Orm.Preload("OrderDetails").Preload("OrderImages", "type=0").Preload("ReceiptImages", models.OrderImage{Type: 1}).
		Scopes(
			cDto.MakeCondition(d.GetNeedSearch()),
		).
		Find(list).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetOrderInfo error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// CheckExistUnCompletedOrder
func (e *OrderInfo) CheckExistUnCompletedOrder(d *dto.OrderInfoCheckExistUnCompletedOrderReq, list *[]models.OrderInfo) error {
	err := e.Orm.Where("user_company_id = ? and order_status in (0, 1, 2 ,99)", d.Id).Find(list).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// AddProduct 添加商品
func (e *OrderInfo) AddProduct(c *dto.OrderInfoAddProductReq) (data []dto.OrderInfoProductInsertReq, err error) {
	skus := utils.Split(c.SkuCode)
	productMap := getProductMap(e.Orm, skus, c.WarehouseCode)

	err = c.Valid(e.Orm, productMap)
	if err != nil {
		return nil, err
	}

	var goodsIds []int
	var vendorIds []string
	for _, product := range productMap {
		goodsIds = append(goodsIds, product.Id)
		vendorIds = append(vendorIds, strconv.Itoa(product.VendorId))
	}
	// 货主map
	vendorMap := getVendorMap(e.Orm, vendorIds)
	// 库存map
	stockMap := getStockMap(e.Orm, goodsIds, c.WarehouseCode)

	for _, datum := range productMap {
		tmp := dto.OrderInfoProductInsertReq{
			GoodsId:        datum.Id,
			SkuCode:        datum.SkuCode,
			ProductName:    datum.Product.NameZh,
			ProductModel:   datum.Product.MfgModel,
			ProductId:      datum.Product.Id,
			BrandId:        datum.Product.BrandId,
			BrandName:      datum.Product.Brand.BrandZh,
			BrandEname:     datum.Product.Brand.BrandEn,
			Moq:            datum.Product.SalesMoq,
			Unit:           datum.Product.SalesUom,
			SalePrice:      datum.MarketPrice,
			VendorId:       datum.VendorId,
			VendorSkuCode:  datum.SupplierSkuCode,
			ProductNo:      datum.ProductNo,
			SubTotalAmount: datum.MarketPrice,
			Quantity:       1,
		}

		tmp.Tax, _ = strconv.ParseFloat(datum.Product.Tax, 64)
		if vendor, ok := vendorMap[tmp.VendorId]; ok {
			tmp.VendorName = vendor.NameZh
		}
		if stock, ok := stockMap[tmp.GoodsId]; ok {
			tmp.Stock = stock
		}
		mediaRelationship := *datum.Product.MediaRelationship
		if len(mediaRelationship) > 0 {
			tmp.ProductPic = mediaRelationship[0].MediaInstant.MediaDir
		}
		data = append(data, tmp)
	}

	return
}

// CheckIsOverBudget 查看是否超出预算
func (e *OrderInfo) CheckIsOverBudget(c *dto.OrderIdsReq) (orderIds []string, err error) {
	var list []models.OrderInfo
	var data models.OrderInfo

	err = e.Orm.Model(&data).Where("order_id in ?", c.OrderIds).Find(&list).Error
	if err != nil {
		return nil, err
	}

	for _, order := range list {
		var departmentBudgetM = modelsUc.DepartmentBudget{}
		if isOverBudget := departmentBudgetM.CheckIsOverBudget(e.Orm, order.UserId, order.TotalAmount, time.Now().Format("200601")); isOverBudget {
			orderIds = append(orderIds, order.OrderId)
		}
	}

	return
}

// Confirm 确认订单
func (e *OrderInfo) Confirm(c *dto.OrderIdsReq) (orderIds []string, err error) {
	var list []models.OrderInfo
	var data models.OrderInfo

	err = e.Orm.Model(&data).Preload("OrderDetails").Where("order_id in ?", c.OrderIds).Find(&list).Error
	if err != nil {
		return nil, err
	}

	resList, err := e.CheckIsOverBudget(c)
	for _, order := range list {
		if order.OrderStatus != 5 {
			orderIds = append(orderIds, "订单["+order.OrderId+"]状态不允许确认")
			continue
		}
		if order.RmaStatus != 0 && order.RmaStatus != 99 {
			orderIds = append(orderIds, "订单["+order.OrderId+"]正在售后中，不可确认")
			continue
		}
		oldData := order

		order.OrderStatus = 0
		order.ConfirmTime = time.Now()
		order.IsFormalOrder = 1
		// 订单是否超预算
		if utils.InArrayString(order.OrderId, resList) {
			order.IsOverBudget = 1
		} else {
			order.IsOverBudget = 0
		}

		tx := e.Orm.Debug().Begin()
		// 订单变更为未发货
		err = tx.Save(&order).Error
		if err != nil {
			tx.Rollback()
			orderIds = append(orderIds, "订单["+order.OrderId+"]更新失败，"+err.Error())
			continue
		}
		var stockOutbound dtoWc.StockOutboundInsertForOrderReq
		stockOutbound.SourceCode = order.OrderId
		stockOutbound.Remark = order.Remark
		stockOutbound.WarehouseCode = order.WarehouseCode
		createStockOutboundFlag := true

		var skus []string
		for _, detail := range order.OrderDetails {
			skus = append(skus, detail.SkuCode)
		}
		productMap := getProductMap(e.Orm, skus, order.WarehouseCode)
		for _, detail := range order.OrderDetails {
			if product, ok := productMap[detail.SkuCode]; ok {
				detail.ProductNo = product.ProductNo
			}
			detail.LockStock = 0
			err = tx.Save(&detail).Error
			if err != nil {
				tx.Rollback()
				orderIds = append(orderIds, "订单["+order.OrderId+"]库存更新失败，"+err.Error())
				createStockOutboundFlag = false
				break
			}

			// 商品数量-取消数量>0的商品才生成出库单 商品数量-取消数量=0说明该商品行以及全部取消
			if detail.Quantity-detail.CancelQuantity > 0 {
				stockOutbound.StockOutboundProducts = append(stockOutbound.StockOutboundProducts, dtoWc.StockOutboundProductsForOrderReq{
					SkuCode:  detail.SkuCode,
					Quantity: detail.Quantity - detail.CancelQuantity,
					GoodsId:  detail.GoodsId,
					VendorId: detail.VendorId,
				})
			}
		}
		if createStockOutboundFlag == false {
			continue
		}
		// 生成出货单
		_, err = serviceWc.CreateOutboundForOrder(tx, &stockOutbound)
		if err != nil {
			tx.Rollback()
			orderIds = append(orderIds, "订单["+order.OrderId+"]库存更新失败，"+err.Error())
			continue
		}
		//// 扣减预算
		//var departmentBudgetM = modelsUc.DepartmentBudget{}
		//err = departmentBudgetM.UpdateBudget(e.Orm, order.UserId, -order.TotalAmount, time.Now().Format("200601"))
		//if err != nil {
		//	tx.Rollback()
		//	orderIds = append(orderIds, "订单[" + order.OrderId + "]预算修改失败，" + err.Error())
		//	break
		//}
		e.createLog(c, oldData, order, models.LogTypeOrderInfoConfirm)
		tx.Commit()
	}

	return orderIds, nil
}

func (e *OrderInfo) AutoConfirm(c *gin.Context) (echoMsg string, err error) {
	tenants := global.GetTenants()
	if len(tenants) <= 0 {
		err = errors.New("暂无租户")
		return
	}
	echoMsg = "订单自动确认脚本执行开始："
	for tenantKey, tenant := range tenants {
		//tenantDbSource := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v_%v?charset=utf8mb4&parseTime=True&loc=Local&timeout=1000ms",
		//	tenant.DatabaseUsername,
		//	tenant.DatabasePassword,
		//	tenant.DatabaseHost,
		//	tenant.DatabasePort,
		//	tenant.TenantDBPrefix(),
		//	"oc")
		//tenantDb, err := gorm.Open(mysql.Open(tenantDbSource), &gorm.Config{})
		//if err != nil {
		//	echoMsg = fmt.Sprintf("%s 数据库连接失败", echoMsg)
		//	continue
		//}

		// oc出的接口，获取的tenantDb就是 oc的db连接
		tenantDBPrefix := tenant.TenantDBPrefix()
		tenantDb := sdk.Runtime.GetDbByKey(tenantDBPrefix)
		echoMsg = fmt.Sprintf("%s\r\n租户%s[%s]:", echoMsg, tenant.Name, tenant.DatabaseName)
		if tenantDb == nil {
			echoMsg = fmt.Sprintf("%s 数据库连接失败", echoMsg)
			continue
		}

		// 将tenant-id放到db.statements.context中，用于接口或方法中 获取数据库前缀
		c.Request.Header.Set("tenant-id", tenantKey)
		tenantDb = tenantDb.Session(&gorm.Session{
			Context: c,
		})

		var modelCompany modelsUc.CompanyInfo
		companyList := modelCompany.GetRowsByCondition(tenantDb, "company_status = 1 and order_auto_confirm = ?", 1)
		if len(companyList) == 0 {
			echoMsg = fmt.Sprintf("%s 无需要自动确认订单的公司", echoMsg)
			continue
		}
		var companyIds []int
		for _, info := range companyList {
			companyIds = append(companyIds, info.Id)
		}
		_ = utils.StructColumn(&companyIds, companyList, "Id", "")

		var orderInfos []models.OrderInfo
		tenantDb.Model(&models.OrderInfo{}).Select("order_id").Where("order_status = 5 and user_company_id in ?", companyIds).Find(&orderInfos)
		var orderIdsReq dto.OrderIdsReq
		_ = utils.StructColumn(&orderIdsReq.OrderIds, orderInfos, "OrderId", "")
		if len(orderIdsReq.OrderIds) == 0 {
			echoMsg = fmt.Sprintf("%s 无需要自动确认的订单", echoMsg)
			continue
		}
		fmt.Printf("#自动确认订单:%s\n", utils.TimeFormat(time.Now()))
		e.Orm = tenantDb
		errOrderMsgs, _ := e.Confirm(&orderIdsReq)
		if len(errOrderMsgs) > 0 {
			echoMsg = fmt.Sprintf("%s\r\n部分订单确认失败：%s", echoMsg, strings.Join(errOrderMsgs, " | "))
		}
		echoMsg = fmt.Sprintf("%s\r\n", echoMsg)
	}
	echoMsg = fmt.Sprintf("%s\r\n 脚本执行完毕", echoMsg)

	return echoMsg, nil
}

func (e *OrderInfo) AutoSignFor(c *gin.Context, p *actions.DataPermission) (echoMsg string, err error) {
	tenants := global.GetTenants()
	if len(tenants) <= 0 {
		err = errors.New("暂无租户")
		return
	}
	echoMsg = "订单自动签收脚本执行开始："
	for tenantKey, tenant := range tenants {

		// oc出的接口，获取的tenantDb就是 oc的db连接
		tenantDBPrefix := tenant.TenantDBPrefix()

		tenantDb := sdk.Runtime.GetDbByKey(tenantDBPrefix)
		echoMsg = fmt.Sprintf("%s\r\n租户%s[%s]:", echoMsg, tenant.Name, tenant.DatabaseName)
		if tenantDb == nil {
			echoMsg = fmt.Sprintf("%s 数据库连接失败", echoMsg)
			continue
		}

		// 将tenant-id放到db.statements.context中，用于接口或方法中 获取数据库前缀
		c.Request.Header.Set("tenant-id", tenantKey)
		tenantDb = tenantDb.Session(&gorm.Session{
			Context: c,
		})

		var modelCompany modelsUc.CompanyInfo
		companyList := modelCompany.GetRowsByCondition(tenantDb, "company_status = 1 and order_auto_sign_for = ?", 1)
		if len(companyList) == 0 {
			echoMsg = fmt.Sprintf("%s 无需要自动签收订单的公司", echoMsg)
			continue
		}

		var companyIds []int
		for _, info := range companyList {
			companyIds = append(companyIds, info.Id)
		}
		_ = utils.StructColumn(&companyIds, companyList, "Id", "")

		// 查询订单状态=2已发货 并且无售后或售后完成 并且订单商品都已经发货完毕并且出库时间>48小时 并且上传过回单图片的
		orders, err := e.GetAutoAcceptOrders(tenantDb, companyIds)
		if err != nil {
			echoMsg = fmt.Sprintf("%s 租户%s自动签收订单查询失败", echoMsg, tenant.Name)
			continue
		}
		if len(orders) == 0 {
			echoMsg = fmt.Sprintf("%s 无需要自动签收的订单", echoMsg)
			continue
		}

		fmt.Printf("#自动签收开始:%s\n", utils.TimeFormat(time.Now()))

		okNum := 0
		errNum := 0
		for _, order := range orders {
			e.Orm = tenantDb

			req := &dto.OrderInfoReceiptReq{
				Id:            order.Id,
				ReceiptImages: []dto.OrderInfoOrderImagesInsertReq{},
				IsAuto:        1,
			}

			e.Orm.Model(&models.OrderImage{}).
				Where("order_id = ?", order.OrderID).
				Where("type = 2").
				Find(&req.ReceiptImages)

			receiptErr := e.Receipt(req, p)
			if receiptErr != nil {
				echoMsg += fmt.Sprintf("%s\r\n订单签收失败：%s", order.OrderID, receiptErr.Error())
				errNum++
			} else {
				okNum++
			}
		}

		echoMsg += fmt.Sprintf("\r\n 执行完毕租户%s:%d成功 %d失败", tenant.Name, okNum, errNum)
	}
	echoMsg += fmt.Sprintf("\r\n 脚本执行完毕")

	return echoMsg, nil
}

// 定义订单信息模型
type AutoAcceptOrders struct {
	Id          int      `gorm:"column:id;primaryKey"`
	OrderID     string   `gorm:"column:order_id"`
	OrderStatus int      `gorm:"column:order_status"`
	RMAStatus   int      `gorm:"column:rma_status"`
	Outbound    Outbound // 关联到 Outbound 模型
}

// 定义出库记录模型
type Outbound struct {
	SourceCode   string    `gorm:"column:source_code"`
	Type         int       `gorm:"column:type"`
	Status       int       `gorm:"column:status"`
	OutboundTime time.Time `gorm:"column:outbound_time"`
}

// GetAutoAcceptOrders 获取满足自动签收条件的订单
func (e *OrderInfo) GetAutoAcceptOrders(db *gorm.DB, companyIds []int) ([]AutoAcceptOrders, error) {
	var orders []AutoAcceptOrders

	wcPrefix := global.GetTenantWcDBNameWithDB(db)

	// 执行查询
	if err := db.Debug().
		Model(&models.OrderInfo{}).
		Joins("JOIN "+wcPrefix+".stock_outbound outbound ON outbound.source_code = order_info.order_id AND outbound.type = 1 AND outbound.status = 2").
		Where("order_info.order_status = ? AND TIMESTAMPDIFF(HOUR, outbound.outbound_time, NOW()) > 48 AND rma_status IN (?, ?)", 1, 0, 99).
		Where("user_company_id in ?", companyIds).
		Where("order_id in (select order_info.order_id from order_info join order_image oi on order_info.order_id = oi.order_id where order_status = 1 and oi.type = 2 group by order_info.order_id)").
		Find(&orders).Error; err != nil {
		return nil, err
	}

	return orders, nil
}

func (e *OrderInfo) GetCompanyByIds(companyIds []string) map[int]dtoUc.CompanyInfoGetSelectPageData {
	// 公司名称
	companyResult := ucClient.ApiByDbContext(e.Orm).GetCompanyByIds(strings.Join(companyIds, ","))
	companyResultInfo := &struct {
		response.Response
		Data struct {
			response.Page
			List []dtoUc.CompanyInfoGetSelectPageData
		}
	}{}
	companyResult.Scan(companyResultInfo)

	companyMap := make(map[int]dtoUc.CompanyInfoGetSelectPageData)

	for _, data := range companyResultInfo.Data.List {
		companyMap[data.Id] = data
	}

	return companyMap
}

func (e *OrderInfo) AutOfStock(c *gin.Context) (echoMsg string, err error) {

	tenants := global.GetTenants()
	if len(tenants) <= 0 {
		err = errors.New("暂无租户")
		return
	}
	echoMsg = "缺货订单补货脚本执行开始："
	for tenantKey, tenant := range tenants {
		//tenantDbSource := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v_%v?charset=utf8mb4&parseTime=True&loc=Local&timeout=1000ms",
		//	tenant.DatabaseUsername,
		//	tenant.DatabasePassword,
		//	tenant.DatabaseHost,
		//	tenant.DatabasePort,
		//	tenant.TenantDBPrefix(),
		//	"oc")
		//tenantDb, err := gorm.Open(mysql.Open(tenantDbSource), &gorm.Config{})
		//if err != nil {
		//	echoMsg = fmt.Sprintf("%s 数据库连接失败", echoMsg)
		//	continue
		//}

		// oc出的接口，获取的tenantDb就是 oc的db连接
		tenantDBPrefix := tenant.TenantDBPrefix()
		tenantDb := sdk.Runtime.GetDbByKey(tenantDBPrefix)
		echoMsg = fmt.Sprintf("%s\r\n租户%s[%s]:", echoMsg, tenant.Name, tenant.DatabaseName)
		if tenantDb == nil {
			echoMsg = fmt.Sprintf("%s 数据库连接失败", echoMsg)
			continue
		}

		// 将tenant-id放到db.statements.context中，用于接口或方法中 获取数据库前缀
		c.Request.Header.Set("tenant-id", tenantKey)
		tenantDb = tenantDb.Session(&gorm.Session{
			Context: c,
		})

		var orderDetails []models.OrderDetail
		tenantDb.Table("order_detail t").Select("t.*").
			Joins("left join order_info oi on t.order_id = oi.order_id").
			Where("oi.order_status = 6 and t.warehouse_code != ''").
			Where("oi.rma_status in (0,99)"). // 新增逻辑 缺货订单可以申请售后 所以在补货这里做拦截 售后中的订单就先不补货 直到售后完结
			Find(&orderDetails)
		if len(orderDetails) == 0 {
			echoMsg = fmt.Sprintf("%s 无订单缺货", echoMsg)
			continue
		}
		detailMap := make(map[string][]models.OrderDetail)
		for _, orderDetail := range orderDetails {
			detailMap[orderDetail.OrderId] = append(detailMap[orderDetail.OrderId], orderDetail)
		}
		for orderId, details := range detailMap {
			echoMsg = fmt.Sprintf("%s\r\n订单[%s]开始：", echoMsg, orderId)
			orderStatusConfirm := true
			var goodsIds []int
			_ = utils.StructColumn(&goodsIds, details, "GoodsId", "")
			warehouseCode := details[0].WarehouseCode
			// 可用库存map
			stockMap := getStockMap(tenantDb, goodsIds, warehouseCode)
			tx := tenantDb.Debug().Begin()
			for _, detail := range details {
				// 商品数量 因为加了缺货状态可以售后取消 所以 实际上的商品数量是Quantity-CancelQuantity
				quantity := detail.Quantity - detail.CancelQuantity
				if quantity == detail.LockStock {
					continue
				}
				unLockStock := quantity - detail.LockStock // 未锁库存
				// 可用库存
				stock, ok := stockMap[detail.GoodsId]
				if !ok {
					stock = 0
				}

				if stock >= unLockStock {
					detail.LockStock = quantity
					err = tx.Save(&detail).Error
					if err != nil {
						tx.Rollback()
						echoMsg = fmt.Sprintf("%s商品[%s]锁库数量更新失败%s | ", echoMsg, detail.SkuCode, err.Error())
						continue
					}
					// 锁库
					err = modelsWc.LockStockInfoForOrder(tenantDb, unLockStock, detail.GoodsId, warehouseCode, orderId, "补货定时任务锁库")
					if err != nil {
						tx.Rollback()
						echoMsg = fmt.Sprintf("%s商品[%s]锁库失败%s | ", echoMsg, detail.SkuCode, err.Error())
						continue
					}
				} else {
					orderStatusConfirm = false
					if stock > 0 && unLockStock > stock {
						detail.LockStock = detail.LockStock + stock
						err = tx.Save(&detail).Error
						if err != nil {
							tx.Rollback()
							echoMsg = fmt.Sprintf("%s商品[%s]锁库数量更新失败%s | ", echoMsg, detail.SkuCode, err.Error())
							continue
						}
						// 锁库
						err = modelsWc.LockStockInfoForOrder(tenantDb, stock, detail.GoodsId, warehouseCode, orderId, "补货定时任务锁库")
						if err != nil {
							tx.Rollback()
							echoMsg = fmt.Sprintf("%s商品[%s]锁库失败%s | ", echoMsg, detail.SkuCode, err.Error())
							continue
						}
					}
				}

			}
			// 订单所有缺口商品均已补货完成 修改订单状态未待确认
			if orderStatusConfirm == true {
				var orderInfo models.OrderInfo
				tenantDb.Model(&orderInfo).Where("order_id = ?", orderId).First(&orderInfo)
				orderInfo.OrderStatus = 5
				err = tx.Save(&orderInfo).Error
				if err != nil {
					tx.Rollback()
					echoMsg = fmt.Sprintf("%s更新订单状态失败%s | ", echoMsg, err.Error())
					continue
				}
			}
			tx.Commit()
			echoMsg = fmt.Sprintf("%s\r\n订单[%s]结束", echoMsg, orderId)
		}
		echoMsg = fmt.Sprintf("%s\r\n", echoMsg)

	}
	echoMsg = fmt.Sprintf("%s\r\n 脚本执行完毕", echoMsg)

	return echoMsg, nil
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

// 获取库存map sku => goods
func getProductMap(tx *gorm.DB, skus []string, warehouseCode string) (productMap map[string]modelsPc.Goods) {
	productResult := pcClient.ApiByDbContext(tx).GetGoodsBySkuCodeReq(dtoPc.GetGoodsBySkuCodeReq{
		SkuCode:       skus,
		WarehouseCode: warehouseCode,
		OnlineStatus:  1,
		Status:        1,
	})
	productResultInfo := &struct {
		response.Response
		Data []modelsPc.Goods
	}{}
	productResult.Scan(productResultInfo)
	productMap = make(map[string]modelsPc.Goods)
	for _, product := range productResultInfo.Data {
		productMap[product.SkuCode] = product
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
	log := models.OrderInfoLog{
		DataId:       afterModel.Id,
		Type:         logType,
		Data:         string(dataLog),
		BeforeData:   string(beforeDataStr),
		AfterData:    string(afterDataStr),
		CreateBy:     user.GetUserId(e.Orm.Statement.Context.(*gin.Context)),
		CreateByName: user.GetUserName(e.Orm.Statement.Context.(*gin.Context)),
	}
	_ = log.CreateLog("orderInfo", e.Orm)
}

// 给order_detail 设置4级产线
func setDetailCategory(detail *models.OrderDetail, categoryMap map[string][]modelsPc.Category) {
	if categoryList, ok := categoryMap[detail.SkuCode]; ok {
		categoryLen := len(categoryList)
		if categoryLen > 0 {
			detail.CatId = categoryList[0].Id
			detail.CatName = categoryList[0].NameZh
		}
		if categoryLen > 1 {
			detail.CatId2 = categoryList[1].Id
			detail.CatName2 = categoryList[1].NameZh
		}
		if categoryLen > 2 {
			detail.CatId3 = categoryList[2].Id
			detail.CatName3 = categoryList[2].NameZh
		}
		if categoryLen > 3 {
			detail.CatId4 = categoryList[3].Id
			detail.CatName4 = categoryList[3].NameZh
		}
	}
}

// ------------------------------------------------------INNER----------------------------------------------------------

// GetOrderListByUserName 获取OrderInfo对象
func (e *OrderInfo) GetOrderListByUserName(in *dto.OrderListReq, p *actions.DataPermission, list *[]dto.OrderListResp) error {
	var model models.OrderInfo
	err := e.Orm.Model(&model).Where("user_name LIKE ?", "%"+in.UserName+"%").Find(list).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetOrderInfo error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}
