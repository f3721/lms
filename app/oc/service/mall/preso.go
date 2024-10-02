package mall

import (
	"encoding/json"
	"errors"
	"fmt"
	modelsPc "go-admin/app/pc/models"
	dtoPc "go-admin/app/pc/service/admin/dto"
	servicePc "go-admin/app/pc/service/mall"
	dtoMallPc "go-admin/app/pc/service/mall/dto"
	modelsUc "go-admin/app/uc/models"
	"go-admin/app/uc/service/admin"
	modelsWc "go-admin/app/wc/models"
	"go-admin/common"
	pcClient "go-admin/common/client/pc"
	"go-admin/common/global"
	"go-admin/common/middleware/mall_handler"
	commonModels "go-admin/common/models"
	"go-admin/common/msg/email"
	"go-admin/common/msg/sms"
	"go-admin/common/utils"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/samber/lo"

	"github.com/go-admin-team/go-admin-core/sdk"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"go-admin/app/oc/models"
	"go-admin/app/oc/service/mall/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type Preso struct {
	service.Service
}

// SubmitApproval 提交审批
func (e *Preso) SubmitApproval(c *dto.PresoSubmitApprovalReq, respData *map[string]string) error {
	var err error
	ctx := e.Orm.Statement.Context.(*gin.Context)
	userId := user.GetUserId(ctx)
	userName := user.GetUserName(ctx)
	warehouseCode := mall_handler.GetUserConfig(ctx).SelectedWarehouseCode

	var preso models.Preso
	var shippingAddress modelsUc.Address
	// 购物车
	userCartService := servicePc.UserCart{Service: e.Service}
	var productList []modelsPc.UserCartGoodsProduct
	if c.BuyNow == 1 {
		productList, err = userCartService.GetBuyNowProductForOrder(userId, c.GoodsId, c.Quantity, warehouseCode)
	} else {
		productList, err = userCartService.GetCartProductForOrder(userId, warehouseCode, true)
	}
	if err != nil {
		return err
	}
	var companyInfoM modelsUc.CompanyInfo
	companyInfo := companyInfoM.GetRowByUserId(e.Orm, user.GetUserId(ctx))
	if companyInfo.IsEas != 1 {
		return errors.New("用户所属公司无需提交审批，请刷新页面重新尝试")
	}
	err = c.Valid(e.Orm, &shippingAddress, productList, companyInfo)
	if err != nil {
		return err
	}
	if c.ApproveflowId <= 0 {
		return errors.New("请选择审批流")
	}

	c.Generate(&preso)
	preso.Ip = common.GetClientIP(ctx)
	preso.CreateBy = userId
	preso.CreateByName = userName
	preso.ApproveStatus = 0
	preso.Step = 0
	// preso.ApproveUsers =
	preso.UserId = userId
	preso.UserName = userName
	preso.UserCompanyId = companyInfo.Id
	preso.UserCompanyName = companyInfo.CompanyName
	presoTerm := 15
	if companyInfo.PresoTerm > 0 {
		presoTerm = companyInfo.PresoTerm
	}
	now := time.Now()
	preso.ExpireTime = time.Date(now.Year(), now.Month(), now.Day()+presoTerm, now.Hour(), now.Minute(), now.Second(), 0, now.Location())
	preso.WarehouseCode = warehouseCode
	//地址相关
	preso.CountryId = shippingAddress.CountryId
	preso.CountryName = shippingAddress.CountryName
	preso.ProvinceId = shippingAddress.ProvinceId
	preso.ProvinceName = shippingAddress.ProvinceName
	preso.CityId = shippingAddress.CityId
	preso.CityName = shippingAddress.CityName
	preso.AreaId = shippingAddress.AreaId
	preso.AreaName = shippingAddress.AreaName
	preso.TownId = shippingAddress.TownId
	preso.TownName = shippingAddress.TownName
	preso.Mobile = shippingAddress.CellPhone
	preso.Telephone = shippingAddress.Telephone
	preso.Address = shippingAddress.DetailAddress
	preso.CompanyName = shippingAddress.CompanyName
	preso.Consignee = shippingAddress.ReceiverName
	preso.ContactEmail = shippingAddress.Email

	tx := e.Orm.Debug().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// 预订单
	err = tx.Save(&preso).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	var userProductRemark map[string]dto.PresoUserProductRemarkReq
	if len(c.UserProductRemark) > 0 {
		_ = utils.StructColumn(&userProductRemark, c.UserProductRemark, "", "SkuCode")
	}

	totalAmount := 0.00
	// 预订单明细
	for index, detail := range productList {
		var presoDetail models.PresoDetail
		presoDetail.PresoNo = preso.PresoNo
		presoDetail.UserId = preso.UserId
		presoDetail.UserName = preso.UserName
		presoDetail.CreateByName = preso.UserName
		presoDetail.WarehouseCode = preso.WarehouseCode
		presoDetail.CreateBy = preso.CreateBy
		presoDetail.SkuCode = detail.SkuCode
		presoDetail.VendorId = detail.VendorId
		presoDetail.GoodsId = detail.GoodsId
		presoDetail.Quantity = detail.Quantity
		presoDetail.MarketPrice = detail.MarketPrice
		presoDetail.SalePrice = detail.MarketPrice
		presoDetail.ProductNo = detail.ProductNo
		presoDetail.ProductName = detail.NameZh
		presoDetail.ProductPic = detail.Image.MediaDir
		presoDetail.ApproveQuantity = 0
		presoDetail.Step = 0
		if _, ok := userProductRemark[detail.SkuCode]; ok {
			presoDetail.UserProductRemark = userProductRemark[detail.SkuCode].Remark
		}
		totalAmount = utils.AddFloat64(totalAmount, utils.MulFloat64AndInt(presoDetail.MarketPrice, presoDetail.Quantity))
		err = tx.Save(&presoDetail).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		// 组合发邮件时候需要的产品信息|最多取5条
		if index <= 5 {
			preso.PresoDetails = append(preso.PresoDetails, presoDetail)
		}
	}

	// 审批文档
	for _, img := range c.PresoImage {
		var presoImage models.PresoImage
		img.Generate(&presoImage)
		presoImage.PresoNo = preso.PresoNo
		err = tx.Create(&presoImage).Error
		if err != nil {
			tx.Rollback()
			return err
		}
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

	// 自动审批流程
	// 通过c.ApproveflowId获取审批流及明细,并生成preso_log
	var modelUserApprove modelsUc.UserApprove
	workflows, err := modelUserApprove.Workflows(e.Orm, userId, []int{c.ApproveflowId})
	if err != nil {
		tx.Rollback()
		return errors.New("用户未设置审批流：" + err.Error())
	}
	for _, workflow := range workflows {
		var presoLog models.PresoLog
		presoLog.PresoNo = preso.PresoNo
		presoLog.Type = 2
		presoLog.ApproveflowId = c.ApproveflowId
		presoLog.ApproveflowItemId = workflow.ApproveID
		presoLog.UserId = workflow.UserID
		presoLog.ApproveStatus = 0
		presoLog.Step = workflow.Priority
		presoLog.ApproveRankType = workflow.ApproveRankType
		presoLog.TotalPriceLimit = workflow.TotalPriceLimit
		presoLog.LimitType = workflow.LimitType
		presoLog.CreateBy = userId
		presoLog.CreateByName = userName
		err = tx.Create(&presoLog).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	// 自动审批
	autoApproveFlag, err := e.autoApproveNextStep(preso.PresoNo, 1)
	if err != nil {
		return err
	}
	if !autoApproveFlag {
		approveUsers := e.GetStepApprovers(preso, preso.Step+1)
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			e.PendingApprovalEmail(preso, approveUsers) // 邮件
		}()
		go func() {
			defer wg.Done()
			e.PendingApprovalMsg(preso, approveUsers) // 短信
		}()
		wg.Wait()
	}

	tmp := *respData
	tmp["preso"] = preso.PresoNo
	tmp["totalAmount"] = strconv.FormatFloat(totalAmount, 'f', 2, 64)
	*respData = tmp
	return nil
}

// GetPage 获取Preso列表
func (e *Preso) GetPage(c *dto.PresoGetPageReq, p *actions.DataPermission, list *[]dto.PresoGetPageResp, count *int64) error {
	var err error
	var data models.Preso

	if c.Type != 1 && c.Type != 2 {
		return errors.New("参数类型有误")
	}

	err = e.Orm.Model(&data).Preload("PresoDetails").Preload("Files", "type=0").
		Joins("left join preso_detail pd on preso.preso_no = pd.preso_no").
		Joins("left join preso_log pl on preso.preso_no = pl.preso_no").
		Scopes(
			dto.PresoGetPageMakeCondition(c, e.Orm),
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Group("preso.preso_no").
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("PresoService GetPage error:%s \r\n", err)
		return err
	}
	if *count <= 0 {
		return nil
	}
	tmpList := *list

	var skus []string
	var vendorIds []string
	var presoNos []string
	for _, preso := range tmpList {
		presoNos = append(presoNos, preso.PresoNo)
		for _, detail := range preso.PresoDetails {
			skus = append(skus, detail.SkuCode)
			vendorIds = append(vendorIds, strconv.Itoa(detail.VendorId))
		}
	}

	// 税率
	var modelProduct modelsPc.Product
	taxMap := modelProduct.GetTaxBySku(e.Orm, skus)
	currentTax := utils.GetCurrentTaxRate()

	// 审批流节点
	var modelPresoLog models.PresoLog
	workFlowNodesMap := modelPresoLog.GetWorkFlowNodes(e.Orm, presoNos)

	var orders []models.OrderInfo
	e.Orm.Table("order_info").Where("external_order_no in ?", presoNos).Find(&orders)
	orderMap := make(map[string]string)
	for _, order := range orders {
		orderMap[order.ExternalOrderNo] = order.OrderId
	}

	// 商品map
	productMap := getProductMap(e.Orm, skus)
	// 货主map
	vendorMap := getVendorMap(e.Orm, vendorIds)

	for i, preso := range tmpList {
		totalAmount := 0.00
		totalQuantity := 0
		taxTotalAmount := 0.00
		unRejectTotalAmount := 0.00
		for i2, detail := range preso.PresoDetails {
			tax := currentTax
			if product, ok := productMap[detail.SkuCode]; ok {
				preso.PresoDetails[i2].BrandName = product.Brand.BrandZh
				preso.PresoDetails[i2].MfgModel = product.MfgModel
			}
			if _, ok := vendorMap[detail.VendorId]; ok {
				preso.PresoDetails[i2].VendorName = vendorMap[detail.VendorId].NameZh
			}
			// 含税单价
			tmpTax, ok3 := taxMap[detail.SkuCode]
			if ok3 {
				tax, _ = strconv.ParseFloat(tmpTax, 64)
			}
			tmpList[i].PresoDetails[i2].NakedUnitPrice, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", utils.DivFloat64(detail.SalePrice, utils.AddFloat64(tax, 1.00))), 64)
			detailAmount := utils.MulFloat64AndInt(detail.SalePrice, detail.Quantity) // 行小计
			totalAmount = utils.AddFloat64(totalAmount, detailAmount)                 // 总金额
			totalQuantity = totalQuantity + detail.Quantity                           // 总数量
			if detail.ApproveStatus != -1 {
				unRejectTotalAmount = utils.AddFloat64(unRejectTotalAmount, detailAmount) // 未驳回的总金额
			}
			taxTotalAmount = utils.AddFloat64(taxTotalAmount, utils.SubFloat64(detailAmount, utils.MulFloat64AndInt(tmpList[i].PresoDetails[i2].NakedUnitPrice, detail.Quantity))) // 税额
		}
		tmpList[i].TotalAmount = totalAmount
		tmpList[i].TotalQuantity = totalQuantity
		tmpList[i].UnRejectTotalAmount = unRejectTotalAmount
		//tmpList[i].TaxTotalAmount, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", utils.SubFloat64(totalAmount, utils.DivFloat64(totalAmount, utils.AddFloat64(1.00, tax)))), 64)
		tmpList[i].TaxTotalAmount = taxTotalAmount
		tmpList[i].ButtonList = e.getButtonList(&preso)
		if preso.ApproveStatus == 0 || preso.ApproveStatus == 10 {
			tmpList[i].ExpireTimeText = utils.GetDiffTextByTime(preso.ExpireTime, time.Now())
		}
		if workFollowNodes, ok := workFlowNodesMap[preso.PresoNo]; ok {
			tmpList[i].WorkFollowNodes = workFollowNodes
			// 审批人 （小程序使用列表） 未完成审批展示下一个审批节点的审批人||审批完成展示最后一个节点的审批人|| 审批驳回 展示最后驳回节点的审批人
			for _, node := range workFollowNodes {
				if (preso.ApproveStatus == 0 || preso.ApproveStatus == 10) && node.Step == preso.Step+1 {
					tmpList[i].ApproveUser = node.ApproveUser
				} else if node.Step == preso.Step {
					tmpList[i].ApproveUser = node.ApproveUser
				}
				if node.IsAutoApprove == 1 && c.Client != "xcx" {
					tmpList[i].ApproveUser = tmpList[i].ApproveUser + "(自动跳过)"
				}
			}
			if preso.ApproveStatus == -2 && c.Client != "xcx" {
				tmpList[i].ApproveUser = tmpList[i].ApproveUser + "(已超时)"
			}

		}
		if orderId, ok := orderMap[preso.PresoNo]; ok {
			tmpList[i].OrderId = orderId
		}
	}

	*list = tmpList

	return nil
}

func (e *Preso) GetPageCount(c *dto.PresoGetPageReq, p *actions.DataPermission, list *[]dto.PresoGetPageResp) error {
	var err error
	var data models.Preso

	if c.Type != 1 && c.Type != 2 {
		return errors.New("参数类型有误")
	}

	err = e.Orm.Model(&data).
		Joins("left join preso_detail pd on preso.preso_no = pd.preso_no").
		Joins("left join preso_log pl on preso.preso_no = pl.preso_no").
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			actions.Permission(data.TableName(), p),
			dto.PresoGetPageMakeCondition(c, e.Orm),
		).
		Group("preso.preso_no").
		Find(list).Error
	if err != nil {
		e.Log.Errorf("PresoService GetPage error:%s \r\n", err)
		return err
	}

	return nil
}

// 获取领用申请列表中可展示的按钮
func (e *Preso) getButtonList(preso *dto.PresoGetPageResp) (res []int) {
	res = append(res, 1) // 查看详情
	res = append(res, 2) // 加入购物车
	if preso.ApproveStatus == 0 {
		res = append(res, 3) // 撤回审批单
	}
	res = append(res, 4) // 上传文档
	res = append(res, 5) // 下载审批明细
	return
}

// Get 获取Preso对象
func (e *Preso) Get(d *dto.PresoGetReq, p *actions.DataPermission, model *dto.PresoGetResp) error {
	var data models.Preso

	err := e.Orm.Model(&data).Preload("PresoDetails").Preload("Files", "type=0").
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		Where("preso_no = ?", d.GetId()).
		First(model).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetPreso error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}

	model.ExpireTimeText = utils.GetDiffTextByTime(model.ExpireTime, time.Now())
	model.AddressFullName = model.Consignee + " " + model.Mobile
	if model.Telephone != "" {
		model.AddressFullName = model.AddressFullName + " " + model.Telephone
	}
	model.AddressFullName = model.AddressFullName + " " + model.ProvinceName + " " + model.CityName + " " + model.AreaName + " " + model.Address

	var skus []string
	var vendorIds []string
	for _, detail := range model.PresoDetails {
		skus = append(skus, detail.SkuCode)
		vendorIds = append(vendorIds, strconv.Itoa(detail.VendorId))
	}
	// 商品map
	productMap := getProductMap(e.Orm, skus)
	// 货主map
	vendorMap := getVendorMap(e.Orm, vendorIds)
	// 税率
	var modelProduct modelsPc.Product
	taxMap := modelProduct.GetTaxBySku(e.Orm, skus)
	currentTax := utils.GetCurrentTaxRate()

	totalAmount := 0.00
	for i, detail := range model.PresoDetails {
		// 含税单价
		tax := currentTax
		tmpTax, ok3 := taxMap[detail.SkuCode]
		if ok3 {
			tax, _ = strconv.ParseFloat(tmpTax, 64)
		}
		model.PresoDetails[i].NakedUnitPrice, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", utils.DivFloat64(detail.SalePrice, utils.AddFloat64(tax, 1.00))), 64)
		model.PresoDetails[i].UntaxedTotal = utils.MulFloat64AndInt(model.PresoDetails[i].NakedUnitPrice, detail.Quantity)
		model.PresoDetails[i].TaxedTotal = utils.MulFloat64AndInt(detail.SalePrice, detail.Quantity)
		model.PresoDetails[i].Tax = utils.SubFloat64(model.PresoDetails[i].TaxedTotal, model.PresoDetails[i].UntaxedTotal)
		if detail.ApproveStatus != -1 {
			totalAmount = utils.AddFloat64(totalAmount, utils.MulFloat64AndInt(detail.SalePrice, detail.Quantity)) // 总金额
		}

		if product, ok := productMap[detail.SkuCode]; ok {
			model.PresoDetails[i].Unit = product.SalesUom
			model.PresoDetails[i].BrandName = product.Brand.BrandZh
			model.PresoDetails[i].MfgModel = product.MfgModel
		}
		if _, ok := vendorMap[detail.VendorId]; ok {
			model.PresoDetails[i].VendorName = vendorMap[detail.VendorId].NameZh
			model.PresoDetails[i].VendorSkuCode = vendorMap[detail.VendorId].Code
		}
	}
	model.TotalAmount = totalAmount

	model.ApproveRemarkText = model.ApproveRemark
	// 审批流节点
	var presoLogM models.PresoLog
	workFlowNodesMap := presoLogM.GetWorkFlowNodes(e.Orm, []string{model.PresoNo})
	if workFollowNodes, ok := workFlowNodesMap[model.PresoNo]; ok {
		model.WorkFollowNodes = workFollowNodes
		for _, node := range workFollowNodes {
			if len(node.Children) == 0 && node.Remark != "" {
				model.ApproveRemarkText = model.ApproveRemarkText + " " + node.Remark
			} else if len(node.Children) > 0 {
				for _, child := range node.Children {
					if child.Remark != "" {
						model.ApproveRemarkText = model.ApproveRemarkText + " " + child.Remark
					}
				}
			}
		}
	}
	model.ApproveRemarkText = strings.TrimSpace(model.ApproveRemarkText)
	// 是否当前节点审批人
	var currentLog models.PresoLog
	err2 := e.Orm.Model(currentLog).Where("preso_no = ? and step = ? and user_id = ?", model.PresoNo, model.Step+1, user.GetUserId(e.Orm.Statement.Context.(*gin.Context))).First(&currentLog).Error
	if err2 == nil {
		model.IsCurrentNodeApprover = 1
	}
	var order models.OrderInfo
	err = e.Orm.Table("order_info").Where("external_order_no = ?", model.PresoNo).First(&order).Error
	if err == nil {
		model.OrderId = order.OrderId
	}

	return nil
}

func (e *Preso) GetExport(d *dto.PresoGetReq, p *actions.DataPermission, res *[]dto.PresoGetExportResp) error {
	var data models.Preso

	var model dto.PresoGetResp
	err := e.Orm.Model(&data).Preload("PresoDetails").
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		Where("preso_no = ?", d.GetId()).
		First(&model).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetPreso error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}

	var skus []string
	var vendorIds []string
	for _, detail := range model.PresoDetails {
		skus = append(skus, detail.SkuCode)
		vendorIds = append(vendorIds, strconv.Itoa(detail.VendorId))
	}
	// 商品map
	productMap := getProductMap(e.Orm, skus)
	// 货主map
	vendorMap := getVendorMap(e.Orm, vendorIds)
	// 税率
	var modelProduct modelsPc.Product
	taxMap := modelProduct.GetTaxBySku(e.Orm, skus)
	currentTax := utils.GetCurrentTaxRate()

	tmpList := *res
	for _, detail := range model.PresoDetails {
		var tmp dto.PresoGetExportResp
		_ = copier.Copy(&tmp, &detail)
		// 含税单价
		tax := currentTax
		tmpTax, ok3 := taxMap[detail.SkuCode]
		if ok3 {
			tax, _ = strconv.ParseFloat(tmpTax, 64)
		}
		tmp.NakedUnitPrice, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", utils.DivFloat64(detail.SalePrice, utils.AddFloat64(tax, 1.00))), 64)
		tmp.UntaxedTotal = utils.MulFloat64AndInt(tmp.NakedUnitPrice, detail.Quantity)
		tmp.TaxedTotal = utils.MulFloat64AndInt(detail.SalePrice, detail.Quantity)
		//tmp.Tax = utils.SubFloat64(tmp.TaxedTotal, tmp.UntaxedTotal)
		tmp.Tax, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", utils.SubFloat64(tmp.TaxedTotal, utils.DivFloat64(tmp.TaxedTotal, 1.13))), 64)

		if product, ok := productMap[detail.SkuCode]; ok {
			tmp.Unit = product.SalesUom
			tmp.SupplierSkuCode = product.SupplierSkuCode
		}
		if _, ok := vendorMap[detail.VendorId]; ok {
			tmp.VendorName = vendorMap[detail.VendorId].NameZh
		}
		tmp.ApproveRemark = model.ApproveRemark
		tmpList = append(tmpList, tmp)
	}
	*res = tmpList

	return nil
}

// FinishApproval 提交审批
func (e *Preso) FinishApproval(c *dto.PresoFinishApprovalReq) error {
	var err error

	ctx := e.Orm.Statement.Context.(*gin.Context)
	userId := user.GetUserId(ctx)
	userName := user.GetUserName(ctx)

	var companyInfoM modelsUc.CompanyInfo
	companyInfo := companyInfoM.GetRowByUserId(e.Orm, user.GetUserId(ctx))

	var preso models.Preso
	err = c.Valid(e.Orm, &preso, companyInfo)
	if err != nil {
		return err
	}

	var presoLogM models.PresoLog
	presoLogsMap := presoLogM.GetPresoLogsMap(e.Orm, c.PresoNo)
	if len(presoLogsMap[c.Step]) <= 0 {
		return errors.New("此流程不存在！")
	}
	_, existNextStep := presoLogsMap[c.Step+1]

	var presoLog models.PresoLog
	// 普签 || 只有一个审批人的或签
	if len(presoLogsMap[c.Step]) == 1 {
		if presoLogsMap[c.Step][0].UserId == userId {
			presoLog = presoLogsMap[c.Step][0]
		}
		if presoLog.ApproveStatus != 0 && presoLog.ApproveStatus != 10 {
			return errors.New("此流程已审批！")
		}
	} else if len(presoLogsMap[c.Step]) > 1 && presoLogsMap[c.Step][0].ApproveRankType == 3 {
		// 或签
		for _, log2 := range presoLogsMap[c.Step] {
			if log2.ApproveStatus != 0 && log2.ApproveStatus != 10 {
				return errors.New("此流程已审批！")
			}
			if log2.UserId == userId {
				presoLog = log2
			}
		}
	}
	if presoLog.Id <= 0 {
		return errors.New("当前审批节点不存在！")
	}
	presoLog.OperUser = userName
	presoLog.UpdateBy = userId
	presoLog.UpdateByName = userName
	presoLog.ApproveRemark = c.ApproveRemark
	if operContentBytes, err := json.Marshal(c.PresoFinishApprovalOperContent); err == nil {
		presoLog.OperContent = string(operContentBytes)
	}
	if c.PassTotal == 0 && c.RejectTotal > 0 {
		presoLog.ApproveStatus = -1
		// 驳回 短信+邮件
		existNextStep = false
		procurementer := e.GetProcurementer(preso)
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			_ = e.RejectEmail(preso, procurementer)
		}()
		go func() {
			defer wg.Done()
			_ = e.RejectMsg(preso, procurementer)
		}()
		wg.Wait()
	} else if c.PassTotal > 0 {
		presoLog.ApproveStatus = 1
	}

	tx := e.Orm.Debug().Begin()

	err = tx.Save(&presoLog).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	// 更新预订单步骤
	preso.Step = c.Step
	preso.ContractNo = c.ContractNo
	//preso.Remark = c.Remark
	//preso.ApproveRemark = c.ApproveRemark
	// 审批状态
	if existNextStep {
		preso.ApproveStatus = 10
	} else {
		if c.PassTotal == 0 && c.RejectTotal > 0 {
			preso.ApproveStatus = -1
		} else if c.PassTotal > 0 {
			preso.ApproveStatus = 1
		}
	}
	err = tx.Save(&preso).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	// 更新预订单商品
	approveItemMap := make(map[string]dto.PresoFinishApprovalProductReq)
	for _, item := range c.ApproveItems {
		approveItemMap[item.SkuCode] = item
	}
	for _, detail := range preso.PresoDetails {
		if _, ok := approveItemMap[detail.SkuCode]; ok {
			if approveItemMap[detail.SkuCode].Approved == 1 {
				detail.Step = c.Step
				if existNextStep == true {
					detail.ApproveStatus = 10
				} else {
					detail.ApproveStatus = 1
				}
			} else if approveItemMap[detail.SkuCode].Approved == -1 {
				detail.Step = c.Step
				detail.ApproveStatus = -1
				detail.RejectByName = userName
				detail.ApproveRemark = c.ApproveRemark
			}
			err = tx.Save(&detail).Error
			if err != nil {
				tx.Rollback()
				return err
			}
		}

	}
	tx.Commit()

	// 驳回 短信+邮件
	if preso.ApproveStatus == -1 {
		procurementer := e.GetProcurementer(preso)
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			_ = e.RejectEmail(preso, procurementer)
		}()
		go func() {
			defer wg.Done()
			_ = e.RejectMsg(preso, procurementer)
		}()
		wg.Wait()
	}

	// 存在下一个节点
	if existNextStep {
		autoApproveFlag, err := e.autoApproveNextStep(preso.PresoNo, c.Step+1)
		if !autoApproveFlag {
			approveUsers := e.GetStepApprovers(preso, preso.Step+1)
			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				defer wg.Done()
				e.PendingApprovalEmail(preso, approveUsers) // 邮件
			}()
			go func() {
				defer wg.Done()
				e.PendingApprovalMsg(preso, approveUsers) // 短信
			}()
			wg.Wait()
			return nil
		}
		if err != nil {
			return err
		}
	} else {
		// 不存在下一个节点 审批完成，预订单转订单 发邮件
		err = e.presoToOrderInfo(preso.PresoNo)
		if err != nil {
			return err
		}
	}

	return nil
}

// 自动处理下一个节点
func (e *Preso) autoApproveNextStep(presoNo string, step int) (autoApproveFlag bool, err error) {

	var preso models.Preso
	err = e.Orm.Model(&preso).Preload("PresoDetails").Where("preso_no = ?", presoNo).First(&preso).Error
	if err != nil {
		err = errors.New("审批单不存在！")
		return
	}

	var presoLogM models.PresoLog
	presoLogsMap := presoLogM.GetPresoLogsMap(e.Orm, presoNo)

	autoApproveFlag = false
	if len(presoLogsMap[step]) <= 0 {
		err = errors.New("此流程不存在！")
		return
	}
	// 是否满足自动审批
	if presoLogsMap[step][0].TotalPriceLimit > 0 {
		// 订单金额
		if presoLogsMap[step][0].LimitType == 1 {
			totalAmount := 0.00
			for _, detail := range preso.PresoDetails {
				totalAmount = utils.AddFloat64(totalAmount, utils.MulFloat64AndInt(detail.SalePrice, detail.Quantity))
			}
			if totalAmount <= presoLogsMap[step][0].TotalPriceLimit {
				autoApproveFlag = true
			}
		} else if presoLogsMap[step][0].LimitType == 2 {
			// 商品金额 需每个商品单价都满足
			count := 0
			for _, detail := range preso.PresoDetails {
				if detail.SalePrice <= presoLogsMap[step][0].TotalPriceLimit {
					count++
				}
			}
			if count == len(preso.PresoDetails) {
				autoApproveFlag = true
			}
		}
	}
	if !autoApproveFlag {
		return
	}
	_, existNextStep := presoLogsMap[step+1]

	tx := e.Orm.Debug().Begin()
	for _, presoLog := range presoLogsMap[step] {
		if presoLog.ApproveStatus != 0 && presoLog.ApproveStatus != 10 {
			err = errors.New("此流程已审批！")
			return
		}

		presoLog.IsAutoApprove = 1
		presoLog.ApproveStatus = 1
		presoLog.ApproveRemark = "系统自动审批"
		err = tx.Save(&presoLog).Error
		if err != nil {
			tx.Rollback()
			return
		}
	}
	// 更新预订单步骤
	preso.Step = step
	//preso.ContractNo = c.ContractNo
	//preso.Remark = c.Remark
	// 审批状态
	if existNextStep == true {
		preso.ApproveStatus = 10
	} else {
		preso.ApproveStatus = 1
	}
	err = tx.Save(&preso).Error
	if err != nil {
		tx.Rollback()
		return
	}

	// 更新预订单商品
	for _, detail := range preso.PresoDetails {
		detail.Step = step
		if existNextStep == true {
			detail.ApproveStatus = 10
		} else {
			detail.ApproveStatus = 1
		}
		err = tx.Save(&detail).Error
		if err != nil {
			tx.Rollback()
			return
		}
	}
	tx.Commit()

	// 存在下一个节点 发邮件
	if existNextStep {

		autoApproveFlagNew, err := e.autoApproveNextStep(preso.PresoNo, step+1)
		if autoApproveFlagNew == false {
			// 发邮件
			return autoApproveFlagNew, nil
		}
		if err != nil {
			return autoApproveFlagNew, err
		}
	} else {
		// 不存在下一个节点 审批完成，预订单转订单 发邮件
		err = e.presoToOrderInfo(preso.PresoNo)
		if err != nil {
			return
		}
	}
	return
}

// 预订单转正式订单
func (e *Preso) presoToOrderInfo(presoNo string) (err error) {
	var preso models.Preso
	err = e.Orm.Model(&preso).Preload("PresoDetails").Where("preso_no = ?", presoNo).First(&preso).Error
	if err != nil {
		return errors.New("审批单不存在！")
	}

	if preso.ApproveStatus == 1 {
		tx := e.Orm.Debug().Begin()
		var approveProducts []models.PresoDetail
		var goodsIds []int
		var skus []string
		productQuantity := 0
		totalAmount := 0.00
		for _, detail := range preso.PresoDetails {
			if detail.ApproveStatus == 1 {
				approveProducts = append(approveProducts, detail)
				goodsIds = append(goodsIds, detail.GoodsId)
				skus = append(skus, detail.SkuCode)
				productQuantity = productQuantity + detail.Quantity
				totalAmount = utils.AddFloat64(totalAmount, utils.MulFloat64AndInt(detail.SalePrice, detail.Quantity))
			}
		}
		if len(approveProducts) > 0 {
			// 订单
			var order models.OrderInfo
			_ = copier.Copy(&order, &preso)
			order.Model = commonModels.Model{Id: 0}
			order.CreatedAt = time.Now()
			order.OrderId = dto.GenerateOrderId()
			order.CreateFrom = "MALL"
			order.IsFormalOrder = 1
			order.IsTransform = 1
			order.ValidFlag = 1
			order.ExternalOrderNo = preso.PresoNo
			order.ClassifyQuantity = len(approveProducts) // 商品校验了必须有产线，且每个商品主产线只有1个，故 商品种类数量为sku个数
			order.ProductQuantity = productQuantity
			order.ItemsAmount = totalAmount
			order.TotalAmount = totalAmount
			order.OriginalItemsAmount = totalAmount
			order.OriginalTotalAmount = totalAmount
			order.OrderStatus = 5

			// 遍历商品 判断库存是否完全满足
			// 全部满足：订单状态=待确认 /  不完全满足，订单状态=缺货
			stockMap := getStockMap(e.Orm, goodsIds, preso.WarehouseCode)

			for _, product := range approveProducts {
				stock, ok := stockMap[product.GoodsId]
				if !ok || stock < product.Quantity {
					order.OrderStatus = 6
					break
				}
			}

			err = tx.Save(&order).Error
			if err != nil {
				tx.Rollback()
				return err
			}

			// 订单明细
			categoryMap := getCategoryMap(e.Orm, skus)
			productMap := getProductMap(e.Orm, skus)
			for _, detail := range approveProducts {
				var orderDetail models.OrderDetail
				_ = copier.Copy(&orderDetail, &detail)
				orderDetail.OrderId = order.OrderId
				orderDetail.Model = commonModels.Model{Id: 0}
				// 商品
				setDetailProduct(&orderDetail, productMap)
				// 产线
				setDetailCategory(&orderDetail, categoryMap)

				orderDetail.FinalQuantity = orderDetail.Quantity
				orderDetail.OriginalQuantity = orderDetail.Quantity
				orderDetail.OriginalItemsMount = utils.MulFloat64AndInt(orderDetail.SalePrice, orderDetail.Quantity)
				orderDetail.SubTotalAmount = orderDetail.OriginalItemsMount

				// 锁库
				lockStock := 0
				if _, ok := stockMap[orderDetail.GoodsId]; ok && stockMap[orderDetail.GoodsId] > 0 {
					lockStock = orderDetail.Quantity
					if stockMap[orderDetail.GoodsId] < orderDetail.Quantity {
						lockStock = stockMap[orderDetail.GoodsId]
					}
					err = modelsWc.LockStockInfoForOrder(e.Orm, lockStock, orderDetail.GoodsId, order.WarehouseCode, order.OrderId, "商城预订单转订单锁库")
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

			// 生成正式领用单邮件通知仓管员
			orderInfoService := OrderInfo{Service: e.Service}
			_ = orderInfoService.NoticeSystemUser(order)
		}
		tx.Commit()

		// 邮件
		procurementer := e.GetProcurementer(preso)
		approveUsers := e.GetStepApprovers(preso, 1)
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = e.ApprovedEmail(preso, procurementer, approveUsers)
		}()
		wg.Wait()
	}
	return
}

// BatchApproval 批量审批
func (e *Preso) BatchApproval(c *dto.PresoBatchApprovalReq) (err error, errMsgList []string) {
	ctx := e.Orm.Statement.Context.(*gin.Context)
	userId := user.GetUserId(ctx)
	userName := user.GetUserName(ctx)
	var presos []models.Preso
	err = c.Valid(e.Orm, &presos)
	if err != nil {
		return
	}

	var wg sync.WaitGroup

	for _, preso := range presos {
		if preso.ApproveStatus != 0 && preso.ApproveStatus != 10 {
			errMsgList = append(errMsgList, "审批单["+preso.PresoNo+"]已审批完成")
			continue
		}

		var presoLogM models.PresoLog
		presoLogsMap := presoLogM.GetPresoLogsMap(e.Orm, preso.PresoNo)
		step := preso.Step + 1
		if len(presoLogsMap[step]) <= 0 {
			errMsgList = append(errMsgList, "审批单["+preso.PresoNo+"]当前审批流程不存在！")
			continue
		}
		_, existNextStep := presoLogsMap[step+1]

		var presoLog models.PresoLog
		// 普签 || 只有一个审批人的或签
		if len(presoLogsMap[step]) == 1 {
			if presoLogsMap[step][0].UserId == userId {
				presoLog = presoLogsMap[step][0]
			}
			if presoLog.ApproveStatus != 0 && presoLog.ApproveStatus != 10 {
				errMsgList = append(errMsgList, "审批单["+preso.PresoNo+"]当前流程已审批！")
				continue
			}
		} else if len(presoLogsMap[step]) > 1 && presoLogsMap[step][0].ApproveRankType == 3 {
			// 或签
			tmpErrMsg := ""
			for _, log2 := range presoLogsMap[step] {
				if log2.ApproveStatus != 0 && log2.ApproveStatus != 10 {
					tmpErrMsg = "审批单[" + preso.PresoNo + "]当前流程已审批！"
					break
				}
				if log2.UserId == userId {
					presoLog = log2
				}
			}
			if tmpErrMsg != "" {
				errMsgList = append(errMsgList, tmpErrMsg)
				continue
			}
		}
		if presoLog.Id <= 0 {
			errMsgList = append(errMsgList, "审批单["+preso.PresoNo+"您无权审批当前节点！")
			continue
		}
		presoLog.OperUser = userName
		presoLog.UpdateBy = userId
		presoLog.UpdateByName = userName
		presoLog.ApproveRemark = c.ApproveRemark

		var operContent dto.PresoFinishApprovalOperContent
		for _, detail := range preso.PresoDetails {
			operContent.ApproveItems = append(operContent.ApproveItems, dto.PresoFinishApprovalProductReq{
				SkuCode:          detail.SkuCode,
				Approved:         c.ApproveStatus,
				Quantity:         detail.Quantity,
				ApprovedQuantity: detail.Quantity,
				Price:            detail.SalePrice,
			})
		}
		if c.ApproveStatus == 1 {
			operContent.PassTotal = len(preso.PresoDetails)
		} else if c.ApproveStatus == -1 {
			operContent.RejectTotal = len(preso.PresoDetails)
		}

		if operContentBytes, err := json.Marshal(operContent); err == nil {
			presoLog.OperContent = string(operContentBytes)
		}
		presoLog.ApproveStatus = c.ApproveStatus

		tx := e.Orm.Debug().Begin()

		err = tx.Save(&presoLog).Error
		if err != nil {
			tx.Rollback()
			errMsgList = append(errMsgList, "审批单["+preso.PresoNo+"]审批失败！"+err.Error())
			continue
		}
		// 更新预订单步骤
		preso.Step = step

		// 驳回 短信+邮件
		if preso.ApproveStatus == -1 {
			existNextStep = false
			procurementer := e.GetProcurementer(preso)
			wg.Add(2)
			go func(preso models.Preso) {
				defer wg.Done()
				_ = e.RejectEmail(preso, procurementer)
			}(preso)
			go func(preso models.Preso) {
				defer wg.Done()
				_ = e.RejectMsg(preso, procurementer)
			}(preso)
		}

		// 审批状态
		if existNextStep {
			preso.ApproveStatus = 10
		} else {
			preso.ApproveStatus = c.ApproveStatus
		}
		err = tx.Save(&preso).Error
		if err != nil {
			tx.Rollback()
			errMsgList = append(errMsgList, "审批单["+preso.PresoNo+"]审批失败！"+err.Error())
			continue
		}
		// 更新预订单商品
		tmpErrMsg2 := ""
		for _, detail := range preso.PresoDetails {
			detail.Step = step
			if existNextStep {
				detail.ApproveStatus = 10
			} else {
				detail.ApproveStatus = c.ApproveStatus
			}
			if c.ApproveStatus == -1 {
				detail.RejectByName = userName
				detail.ApproveRemark = c.ApproveRemark
			}
			err = tx.Save(&detail).Error
			if err != nil {
				tx.Rollback()
				tmpErrMsg2 = "审批单[" + preso.PresoNo + "]sku[" + detail.SkuCode + "]审批失败！" + err.Error()
				break
			}
		}
		if tmpErrMsg2 != "" {
			errMsgList = append(errMsgList, tmpErrMsg2)
			continue
		}
		tx.Commit()

		// 存在下一个节点
		if existNextStep {
			autoApproveFlag, err := e.autoApproveNextStep(preso.PresoNo, step+1)
			if err != nil {
				errMsgList = append(errMsgList, "审批单["+preso.PresoNo+"]审批失败！"+err.Error())
				continue
			}
			if !autoApproveFlag {
				approveUsers := e.GetStepApprovers(preso, preso.Step+1)
				wg.Add(2)
				go func(preso models.Preso) { // 邮件
					defer wg.Done()
					_ = e.PendingApprovalEmail(preso, approveUsers)
				}(preso)
				go func(preso models.Preso) { // 短信
					defer wg.Done()
					_ = e.PendingApprovalMsg(preso, approveUsers)
				}(preso)
			}
		} else {
			// 不存在下一个节点 审批完成，预订单转订单 发邮件
			err = e.presoToOrderInfo(preso.PresoNo)
			if err != nil {
				errMsgList = append(errMsgList, "审批单["+preso.PresoNo+"]审批失败！"+err.Error())
				continue
			}
		}

	}

	wg.Wait()

	return
}

// Withdraw 撤回审批单
func (e *Preso) Withdraw(c *dto.PresoWithdrawReq) (err error) {
	var preso models.Preso
	e.Orm.Where("preso_no = ?", c.PresoNo).First(&preso)
	if preso.Id <= 0 {
		return errors.New("审批单不存在")
	}
	if preso.ApproveStatus != 0 && preso.ApproveStatus != 10 {
		return errors.New("审批单当前状态不能被撤回")
	}
	preso.ApproveStatus = -3

	err = e.Orm.Save(&preso).Error
	if err != nil {
		return err
	}

	// 邮件 + 短信
	approveUsers := e.GetStepApprovers(preso, preso.Step+1)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		_ = e.RevokeEmail(preso, approveUsers)
	}()
	go func() {
		defer wg.Done()
		_ = e.RevokeMsg(preso, approveUsers)
	}()
	wg.Wait()
	return nil
}

// SaveFile 保存上传的文档
func (e *Preso) SaveFile(c *dto.PresoSaveFileReq) (err error) {
	var data models.PresoImage
	if c.PresoNo == "" {
		return errors.New("审批单号必填")
	}
	var preso models.Preso
	e.Orm.Where("preso_no = ?", c.PresoNo).First(&preso)
	if preso.Id <= 0 {
		return errors.New("审批单不存在")
	}

	var count int64
	e.Orm.Model(&data).Where("preso_no = ?", c.PresoNo).Count(&count)
	if count >= 3 {
		return errors.New("审批附件最多上传三个")
	}

	_ = copier.Copy(&data, &c)
	err = e.Orm.Create(&data).Error
	if err != nil {
		return
	}

	return
}

// BuyAgain 再次购买
func (e *Preso) BuyAgain(c *dto.PresoBuyAgainReq, p *actions.DataPermission) (errList []string, err error) {
	var preso = models.Preso{}
	if c.Id == "" {
		return nil, errors.New("参数有误")
	}
	e.Orm.Preload("PresoDetails").Scopes(
		actions.Permission(preso.TableName(), p),
	).Where("preso.preso_no = ?", c.Id).First(&preso)

	if len(preso.PresoDetails) <= 0 {
		return nil, errors.New("审批单无商品")
	}
	ctx := e.Orm.Statement.Context.(*gin.Context)
	if preso.WarehouseCode != mall_handler.GetUserCurrentWarehouseCode(ctx) {
		return nil, errors.New("该审批单所属仓库不在当前所选仓库下，无法再次购买")
	}
	userCartService := servicePc.UserCart{Service: e.Service}
	// 添加到购物车
	for _, detail := range preso.PresoDetails {
		data := dtoMallPc.UserCartInsertReq{
			GoodsId:       detail.GoodsId,
			Quantity:      detail.Quantity,
			WarehouseCode: preso.WarehouseCode,
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

func (e *Preso) CronApprove(c *gin.Context) (echoMsg string, err error) {
	tenants := global.GetTenants()
	if len(tenants) <= 0 {
		err = errors.New("暂无租户")
		return
	}
	echoMsg = "定时审批提醒脚本执行开始："
	for tenantKey, tenant := range tenants {
		echoMsg = fmt.Sprintf("%s\r\n租户%s[%s]:", echoMsg, tenant.Name, tenant.DatabaseName)

		tenantDBPrefix := tenant.TenantDBPrefix()
		tenantDb := sdk.Runtime.GetDbByKey(tenantDBPrefix)
		if tenantDb == nil {
			echoMsg = fmt.Sprintf("%s 数据库连接失败", echoMsg)
			continue
		}
		now := time.Now()
		hour := now.Hour()
		minute := now.Minute() - now.Minute()%15
		prefix := fmt.Sprintf("%v %v * * ", hour, minute)
		var users []modelsUc.UserInfo
		tenantDb.Table(tenant.DatabaseName+"_uc.user_info").Where("user_status = 1 and email_approve_cron_expr like ?", "%"+prefix+"%").Find(&users)
		if len(users) == 0 {
			echoMsg = fmt.Sprintf("%s 无设置当前时间定时审批的用户", echoMsg)
			continue
		}
		var matchUsers []modelsUc.UserInfo
		for _, user := range users {
			cronExprs := strings.Split(user.EmailApproveCronExpr, "@")
			for _, expr := range cronExprs {
				if cronExpr := strings.Split(expr, prefix); len(cronExpr) == 2 {
					weekday := int(now.Weekday()) + 1
					if lo.Contains(strings.Split(cronExpr[1], ","), strconv.Itoa(weekday)) {
						matchUsers = append(matchUsers, user)
						break
					}
				}
			}

		}
		if len(matchUsers) == 0 {
			echoMsg = fmt.Sprintf("%s 无设置当前时间定时审批的用户", echoMsg)
			continue
		}
		// 将tenant-id放到db.statements.context中，用于接口或方法中 获取数据库前缀
		c.Request.Header.Set("tenant-id", tenantKey)
		tenantDb = tenantDb.Session(&gorm.Session{
			Context: c,
		})
		e.Orm = tenantDb
		for _, matchUser := range matchUsers {
			var count int64
			tenantDb.Table("preso").Joins("left join preso_log pl on preso.preso_no = pl.preso_no and preso.step+1 = pl.step").Where("preso.approve_status in (0, 10) and pl.approve_status in (0, 10) and pl.user_id = ?", matchUser.Id).Group("preso.preso_no").Count(&count)
			if count == 0 {
				echoMsg = fmt.Sprintf("%s用户id[%v]无待审批单 | ", echoMsg, matchUser.Id)
				continue
			}
			err = e.CornApprovalEmail(matchUser, int(count))
			if err != nil {
				echoMsg = fmt.Sprintf("%s用户id[%v]fail,msg:%s | ", echoMsg, matchUser.Id, err.Error())
			} else {
				echoMsg = fmt.Sprintf("%s用户id[%v]success | ", echoMsg, matchUser.Id)
			}
		}
		echoMsg = fmt.Sprintf("%s\r\n", echoMsg)
	}
	echoMsg = fmt.Sprintf("%s\r\n 脚本执行完毕", echoMsg)

	return echoMsg, nil
}

func (e *Preso) Expire() (echoMsg string, err error) {
	tenants := global.GetTenants()
	if len(tenants) <= 0 {
		err = errors.New("暂无租户")
		return
	}
	echoMsg = "审批单过期脚本执行开始："
	for _, tenant := range tenants {
		//tenantDbSource := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v_%v?charset=utf8mb4&parseTime=True&loc=Local&timeout=1000ms",
		//	tenant.DatabaseUsername,
		//	tenant.DatabasePassword,
		//	tenant.DatabaseHost,
		//	tenant.DatabasePort,
		//	tenant.TenantDBPrefix(),
		//	"oc")
		//tenantDb, err := gorm.Open(mysql.Open(tenantDbSource), &gorm.Config{})
		echoMsg = fmt.Sprintf("%s\r\n租户%s[%s]:", echoMsg, tenant.Name, tenant.DatabaseName)
		//if err != nil {
		//	echoMsg = fmt.Sprintf("%s 数据库连接失败", echoMsg)
		//	continue
		//}

		// oc出的接口，获取的tenantDb就是 oc的db连接
		tenantDBPrefix := tenant.TenantDBPrefix()
		tenantDb := sdk.Runtime.GetDbByKey(tenantDBPrefix)
		if tenantDb == nil {
			echoMsg = fmt.Sprintf("%s 数据库连接失败", echoMsg)
			continue
		}
		var presos []models.Preso
		tenantDb.Model(&models.Preso{}).Where("approve_status in (0, 10) and expire_time <= ?", time.Now()).Find(&presos)
		if len(presos) == 0 {
			echoMsg = fmt.Sprintf("%s 无过期的审批单", echoMsg)
			continue
		}
		echoMsg = fmt.Sprintf("%s\r\n", echoMsg)
		for _, preso := range presos {
			preso.ApproveStatus = -2
			err = tenantDb.Save(&preso).Error
			if err != nil {
				echoMsg = fmt.Sprintf("%s[%s]fail | ", echoMsg, preso.PresoNo)
			} else {
				echoMsg = fmt.Sprintf("%s[%s]success | ", echoMsg, preso.PresoNo)
			}
		}
	}
	echoMsg = fmt.Sprintf("%s\r\n 脚本执行完毕", echoMsg)

	return echoMsg, nil
}

// DeleteFile 删除审批单文件
func (e *Preso) DeleteFile(d *dto.PresoDeleteFleReq, p *actions.DataPermission) error {
	var data models.PresoImage

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Unscoped().Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemovePresoImage error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}

func getProductMap(db *gorm.DB, skus []string) (productMap map[string]dtoPc.InnerGetProductBySkuResp) {
	result := pcClient.ApiByDbContext(db).GetProductBySku(skus)
	resultInfo := &struct {
		response.Response
		Data []dtoPc.InnerGetProductBySkuResp
	}{}
	result.Scan(resultInfo)

	productMap = make(map[string]dtoPc.InnerGetProductBySkuResp)
	for _, product := range resultInfo.Data {
		productMap[product.SkuCode] = product
	}

	return
}

// Insert 创建Preso对象
func (e *Preso) Insert(c *dto.PresoInsertReq) error {
	var err error
	var data models.Preso
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("PresoService Insert error:%s \r\n", err)
		return err
	}
	return nil
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

// Update 修改Preso对象
func (e *Preso) Update(c *dto.PresoUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.Preso{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("PresoService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除Preso
func (e *Preso) Remove(d *dto.PresoDeleteReq, p *actions.DataPermission) error {
	var data models.Preso

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemovePreso error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}

// 给order_detail 设置商品相关字段
func setDetailProduct(detail *models.OrderDetail, productMap map[string]dtoPc.InnerGetProductBySkuResp) {
	if product, ok := productMap[detail.SkuCode]; ok {
		detail.ProductId = product.Id
		detail.ProductModel = product.MfgModel
		detail.BrandId = product.Brand.Id
		detail.BrandName = product.Brand.BrandZh
		detail.BrandEname = product.Brand.BrandEn
		detail.Unit = product.SalesUom
		detail.VendorId = product.VendorId
		detail.Moq = product.SalesMoq
	}
}

// --------------- 邮件+短信 ---------------

func (e *Preso) GetStepApprovers(preso models.Preso, step int) (approvers []modelsUc.UserInfo) {
	var logs []models.PresoLog
	e.Orm.Where("preso_no = ? and step = ?", preso.PresoNo, step).Find(&logs)
	var userIds []int
	for _, log := range logs {
		userIds = append(userIds, log.UserId)
	}
	ucPrefix := global.GetTenantUcDBNameWithDB(e.Orm)
	if len(userIds) > 0 {
		e.Orm.Table(ucPrefix+".user_info").Where("id in ?", userIds).Find(&approvers)
	}
	return
}

func (e *Preso) GetProcurementer(preso models.Preso) (procurementer modelsUc.UserInfo) {
	ucPrefix := global.GetTenantUcDBNameWithDB(e.Orm)
	e.Orm.Table(ucPrefix+".user_info").Where("id = ?", preso.UserId).First(&procurementer)
	return
}

func (e *Preso) GetUserCompany(user modelsUc.UserInfo) (company modelsUc.CompanyInfo) {
	ucPrefix := global.GetTenantUcDBNameWithDB(e.Orm)
	e.Orm.Table(ucPrefix+".company_info").Where("id = ?", user.CompanyId).First(&company)
	return
}

// 获取审批地址
func (e *Preso) approveUrl(presoNo string, user modelsUc.UserInfo) string {
	// 获取代登录链接地址
	userService := new(admin.UserInfo)
	userService.Orm = e.Orm
	proxyLoginUrl, err := user.ProxyLogin(e.Orm, user.Id)

	e.Log.Info("代登录url地址:", proxyLoginUrl)

	// 获取审批详情地址
	host := utils.GetHostUrl(0)
	redirect := "/approval/detail?presoNo=" + presoNo

	res := ""
	if err == nil {
		// 组合最后的URL
		redirect = url.QueryEscape(redirect)
		res = proxyLoginUrl + "&redirect=" + redirect
	} else {
		res = host + "/#" + redirect
	}

	// 追加多租户ID
	ctx := e.Orm.Statement.Context.(*gin.Context)
	tenantId := ctx.GetHeader("tenant-id")
	res = res + "&tenantId=" + tenantId

	return res
}

// 待审批发邮件|预订单，当前步骤审批人
func (e *Preso) PendingApprovalEmail(preso models.Preso, approver []modelsUc.UserInfo) error {

	// 如果过期取消推送消息
	if preso.ExpireTime.Unix() <= time.Now().Unix() {
		return nil
	}

	for _, user := range approver {
		userCompany := e.GetUserCompany(user)
		// 获取logo
		logo := userCompany.CompanyLogo
		if len(logo) == 0 {
			logo = "https://image-c.ehsy.com/uploadfile/sxyz/img/2022/09/16/20220916100258386.png"
		}

		// 获取审批详情地址
		approveUrl := e.approveUrl(preso.PresoNo, user)

		// 获取模板路径
		path := "static"
		tmpPath := path + "/viewtpl/oc/remind_approve.html"
		if userCompany.ApproveEmailType == 2 {
			tmpPath = path + "/viewtpl/oc/remind_approve_japan.html"
		}

		// 获取邮件标题
		subject := fmt.Sprintf("【狮行驿站】领用审批提醒(%s)", preso.PresoNo)

		// 组合内容
		data := models.ApprovalEmail{
			UserName:     user.UserName,
			PresoNo:      preso.PresoNo,
			ApproveUrl:   approveUrl,
			ProductsTab:  "",
			LogoUrl:      logo,
			ContractNo:   preso.ContractNo,
			PresoDetails: preso.PresoDetails,
		}

		// 模板赋值
		body, err := utils.View(tmpPath, data)
		if err != nil {
			return err
		}

		// 发送邮件
		to := []string{user.UserEmail}
		err = email.AsyncSendEmails(to, subject, body)
		if err != nil {
			return err
		}
	}

	return nil
}

// 待审批发短信|预订单，当前步骤审批人
func (e *Preso) PendingApprovalMsg(preso models.Preso, approver []modelsUc.UserInfo) error {
	// 如果过期取消推送消息
	if preso.ExpireTime.Unix() <= time.Now().Unix() {
		return nil
	}

	for _, user := range approver {
		// 获取审批详情地址
		approveUrl := e.approveUrl(preso.PresoNo, user)

		// 发送短信
		idType := sms.MessageIdTypeApplyNewOrder
		userPhones := []string{user.UserPhone}
		replaceParams := []string{user.UserName, preso.PresoNo, approveUrl}
		sendCode, err := sms.SendSMS(userPhones, idType, replaceParams)
		if err != nil || sendCode.Code != 0 {
			return err
		}
	}

	return nil
}

// 审批通过给一级审批人发邮件
func (e *Preso) sendCompanyEmail(preso models.Preso, approver []modelsUc.UserInfo) error {
	// 获取订单附件
	attachUrls := []string{}
	err := e.Orm.Model(&models.PresoImage{}).Select("url").Where("preso_no = ?", preso.PresoNo).Find(&attachUrls).Error
	if err != nil {
		return err
	}

	for _, user := range approver {
		// 基本信息
		to := []string{user.UserEmail}
		subject := "【狮行驿站】您的领用申请已审批通过"
		body := "亲爱的用户, 您在狮行驿站提交的EAS订单已审批完成。"

		// 发送邮件
		err := email.SendEmailWithAttach(to, subject, body, attachUrls)
		if err != nil {
			return err
		}
	}
	return nil
}

// 审批通过发邮件|预订单，采购人，一级审批人
func (e *Preso) ApprovedEmail(preso models.Preso, procurementer modelsUc.UserInfo, approver []modelsUc.UserInfo) error {
	procurementerCompany := e.GetUserCompany(procurementer)
	// 获取邮件标题
	subject := "【狮行驿站】您的领用申请已审批通过"
	// 获取审批详情地址
	approveUrl := e.approveUrl(preso.PresoNo, procurementer)
	// 获取模板路径
	path := "static"
	tmpPath := path + "/viewtpl/oc/remind_approve_finish.html"
	// 获取logo
	logo := procurementerCompany.CompanyLogo
	if len(logo) == 0 {
		logo = "https://image-c.ehsy.com/uploadfile/sxyz/img/2022/09/16/20220916100258386.png"
	}
	// 组合内容
	data := models.ApprovalEmail{
		UserName:    procurementer.UserName,
		PresoNo:     preso.PresoNo,
		ApproveUrl:  approveUrl,
		ProductsTab: "",
		LogoUrl:     logo,
		ContractNo:  preso.ContractNo,
	}
	// 模板赋值
	body, err := utils.View(tmpPath, data)
	if err != nil {
		return err
	}
	// 发送邮件
	to := []string{procurementer.UserEmail}
	err = email.AsyncSendEmails(to, subject, body)
	if err != nil {
		return err
	}

	// 给一级审批人发送邮件
	if procurementerCompany.IsSendEmail == 1 {
		err = e.sendCompanyEmail(preso, approver)
		if err != nil {
			return err
		}
	}

	return err
}

// 审批驳回发邮件|预订单，采购人
func (e *Preso) RejectEmail(preso models.Preso, procurementer modelsUc.UserInfo) error {
	procurementerCompany := e.GetUserCompany(procurementer)
	// 获取邮件标题
	subject := "【狮行驿站】领用审批产品驳回提醒"
	// 获取审批详情地址
	approveUrl := e.approveUrl(preso.PresoNo, procurementer)
	// 获取模板路径
	path := "static"
	tmpPath := path + "/viewtpl/oc/remind_purchaser.html"
	// 获取logo
	logo := procurementerCompany.CompanyLogo
	if len(logo) == 0 {
		logo = "https://image-c.ehsy.com/uploadfile/sxyz/img/2022/09/16/20220916100258386.png"
	}
	// 组合内容
	data := models.ApprovalEmail{
		UserName:    procurementer.UserName,
		PresoNo:     preso.PresoNo,
		ApproveUrl:  approveUrl,
		ProductsTab: "",
		LogoUrl:     logo,
		ContractNo:  preso.ContractNo,
	}
	// 模板赋值
	body, err := utils.View(tmpPath, data)
	if err != nil {
		return err
	}
	// 发送邮件
	to := []string{procurementer.UserEmail}
	err = email.AsyncSendEmails(to, subject, body)
	if err != nil {
		return err
	}
	return nil
}

// 审批驳回发短信
func (e *Preso) RejectMsg(preso models.Preso, procurementer modelsUc.UserInfo) error {
	// 获取审批详情地址
	approveUrl := e.approveUrl(preso.PresoNo, procurementer)

	// 发送短信
	idType := sms.MessageIdTypeTurnDownUserOrder
	userPhones := []string{procurementer.UserPhone}
	replaceParams := []string{procurementer.UserName, preso.PresoNo, approveUrl}
	sendCode, err := sms.SendSMS(userPhones, idType, replaceParams)
	if err != nil || sendCode.Code != 0 {
		return err
	}

	return nil
}

// 撤回发邮件
func (e *Preso) RevokeEmail(preso models.Preso, approver []modelsUc.UserInfo) error {
	// 如果过期取消推送消息
	if preso.ExpireTime.Unix() <= time.Now().Unix() {
		return nil
	}

	for _, user := range approver {
		userCompany := e.GetUserCompany(user)
		// 获取logo
		logo := userCompany.CompanyLogo
		if len(logo) == 0 {
			logo = "https://image-c.ehsy.com/uploadfile/sxyz/img/2022/09/16/20220916100258386.png"
		}

		// 获取审批详情地址
		approveUrl := e.approveUrl(preso.PresoNo, user)

		// 获取模板路径
		path := "static"
		tmpPath := path + "/viewtpl/oc/withdraw_approve.html"

		// 获取邮件标题
		subject := fmt.Sprintf("【狮行驿站】领用撤回审批提醒(%s)", preso.PresoNo)

		// 组合内容
		data := models.ApprovalEmail{
			UserName:    user.UserName,
			PresoNo:     preso.PresoNo,
			ApproveUrl:  approveUrl,
			ProductsTab: "",
			LogoUrl:     logo,
			ContractNo:  preso.ContractNo,
		}

		// 模板赋值
		body, err := utils.View(tmpPath, data)
		if err != nil {
			return err
		}

		// 发送邮件
		to := []string{user.UserEmail}
		err = email.AsyncSendEmails(to, subject, body)
		if err != nil {
			return err
		}
	}

	return nil
}

// 撤回发短信
func (e *Preso) RevokeMsg(preso models.Preso, approver []modelsUc.UserInfo) error {

	// 发送短信
	for _, user := range approver {
		// 获取审批详情地址
		approveUrl := e.approveUrl(preso.PresoNo, user)

		// 发送短信
		idType := sms.MessageIdTypeRevocationOrder
		userPhones := []string{user.UserPhone}
		replaceParams := []string{user.UserName, preso.PresoNo, approveUrl}
		sendCode, err := sms.SendSMS(userPhones, idType, replaceParams)
		if err != nil || sendCode.Code != 0 {
			return err
		}
	}

	return nil
}

// 获取审批列表地址
func (e *Preso) approveListUrl(user modelsUc.UserInfo) string {
	// 获取代登录链接地址
	userService := new(admin.UserInfo)
	userService.Orm = e.Orm
	proxyLoginUrl, err := user.ProxyLogin(e.Orm, user.Id)

	e.Log.Info("代登录url地址:", proxyLoginUrl)

	// 获取待审批地址
	host := utils.GetHostUrl(0)
	redirect := "Account/approval?activeName=approval"

	res := ""
	if err == nil {
		// 组合最后的URL
		redirect = url.QueryEscape(redirect)
		res = proxyLoginUrl + "&redirect=" + redirect
	} else {
		res = host + "/#" + redirect
	}

	// 追加多租户ID
	ctx := e.Orm.Statement.Context.(*gin.Context)
	tenantId := ctx.GetHeader("tenant-id")
	res = res + "&tenantId=" + tenantId

	return res
}

// 定时任务审批提醒 | 用户信息，待审批条数
func (e *Preso) CornApprovalEmail(user modelsUc.UserInfo, num int) error {

	// 获取logo
	userCompany := e.GetUserCompany(user)
	logo := userCompany.CompanyLogo
	if len(logo) == 0 {
		logo = "https://image-c.ehsy.com/uploadfile/sxyz/img/2022/09/16/20220916100258386.png"
	}

	// 获取审批详情地址
	approveUrl := e.approveListUrl(user)

	// 获取模板路径
	tmpPath := "static/viewtpl/oc/cron_approve.html"

	// 组合模板内容
	data := models.CornApprovalEmail{
		UserName:   user.UserName,
		ApproveUrl: approveUrl,
		Num:        num,
		LogoUrl:    logo,
	}

	// 模板赋值
	body, err := utils.View(tmpPath, data)
	if err != nil {
		return err
	}

	// 发送邮件
	to := []string{user.UserEmail}
	subject := "【狮行驿站】定时领用审批提醒"
	err = email.AsyncSendEmails(to, subject, body)
	if err != nil {
		return err
	}

	return nil
}
