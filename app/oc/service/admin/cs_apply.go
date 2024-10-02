package admin

import (
	"errors"
	modelsPc "go-admin/app/pc/models"
	dtoPc "go-admin/app/pc/service/admin/dto"
	modelsUc "go-admin/app/uc/models"
	modelsWc "go-admin/app/wc/models"
	serviceWc "go-admin/app/wc/service/admin"
	dtoWc "go-admin/app/wc/service/admin/dto"
	pcClient "go-admin/common/client/pc"
	ucClient "go-admin/common/client/uc"
	wcClient "go-admin/common/client/wc"
	"go-admin/common/global"
	cModels "go-admin/common/models"
	"go-admin/common/utils"
	"strconv"
	"strings"
	"time"

	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"github.com/jinzhu/copier"
	"github.com/prometheus/common/log"

	"github.com/gin-gonic/gin"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/oc/models"
	"go-admin/app/oc/service/admin/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type CsApply struct {
	service.Service
}

// GetPage 获取CsApply列表
func (e *CsApply) GetPage(c *dto.CsApplyGetPageReq, p *actions.DataPermission, list *[]models.CsApply, count *int64) error {
	var err error
	var data models.CsApply

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("CsApplyService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取CsApply对象
func (e *CsApply) Get(c *dto.CsApplyGetReq, p *actions.DataPermission, model *models.CsApply) error {
	var data models.CsApply

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		Where("cs_no =?", c.CsNo).
		First(model).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetCsApply error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// GetPage 获取CsApply列表
func (e *CsApply) GetListPage(c *dto.CsApplyGetPageReq, p *actions.DataPermission, list *[]*dto.CsApplyListData, count *int64) error {
	var err error
	var data models.CsApply
	query := e.Orm
	var queryCsNoList []string
	if c.FilterSkuCode != "" || c.FilterProductNo != "" || c.FilterProductName != "" {
		if c.FilterSkuCode != "" {
			query = query.Or("sku_code IN (?)", strings.Split(c.FilterSkuCode, ","))
		}
		if c.FilterProductNo != "" {
			query = query.Or("product_no IN (?)", strings.Split(c.FilterProductNo, ","))
		}
		if c.FilterProductName != "" {
			query = query.Or("product_name LIKE ?", "%"+c.FilterProductName+"%")
		}
		csApplyDetailList := []*models.CsApplyDetail{}
		// 构建查询语句
		query = query.
			Select("cs_no").
			Find(&csApplyDetailList)
		if len(csApplyDetailList) > 0 {
			for _, detail := range csApplyDetailList {
				queryCsNoList = append(queryCsNoList, detail.CsNo)
			}
		} else {
			queryCsNoList = append(queryCsNoList, "-1")
		}
	}
	err = e.Orm.Model(&data).
		Scopes(
			func(db *gorm.DB) *gorm.DB {
				if len(queryCsNoList) > 0 {
					if queryCsNoList[0] == "-1" {
						db = db.Where("cs_no = ?", queryCsNoList[0])
					} else {
						db = db.Where("cs_no in ?", queryCsNoList)
					}
					if p.AuthorityCompanyId != "" {
						//actions.SysUserPermission(data.TableName(), p, 1),
						db.Where("o.user_company_id in ?", utils.Split(p.AuthorityCompanyId))
					}
				}
				return db
			},
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
			actions.SysUserPermission(data.TableName(), p, 2),
			actions.SysUserPermission(data.TableName(), p, 4),
		).
		Select("cs_apply.*,o.user_company_name company_name,o.total_amount total_amount,o.user_name customer_name").
		Joins("left join order_info o on o.order_id = cs_apply.order_id").
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("CsApplyService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取CsApply对象
func (e *CsApply) GetInfo(c *dto.CsApplyGetReq, p *actions.DataPermission, model *dto.CsApplyInfoData) error {
	var data models.CsApply

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		Where("cs_no =?", c.CsNo).
		Select("cs_apply.*,o.user_company_name company_name,o.total_amount total_amount,o.user_name customer_name").
		Joins("left join order_info o on o.order_id = cs_apply.order_id").
		First(model).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetCsApply error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}

	if v, ok := models.CsApplyStatusText[model.CsStatus]; ok {
		model.CsStatusText = v
	}
	if v, ok := models.CsApplyTypeText[model.CsStatus]; ok {
		model.CsTypeText = v
	}

	return nil
}

// Get 获取CsApply对象
func (e *CsApply) GetByCsNo(csNo string) (*models.CsApply, error) {
	var data models.CsApply

	ocPrefix := global.GetTenantOcDBNameWithDB(e.Orm)

	err := e.Orm.Table(ocPrefix+"."+data.TableName()).
		Where("cs_no", csNo).
		First(&data).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("售后单不存在")
		return nil, err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return nil, err
	}
	return &data, nil
}

// GetInfoByCsNo 获取CsApply对象以及售后商品信息
func (e *CsApply) GetInfoByCsNo(csNo string) (*models.CsApply, error) {
	data, err := e.GetByCsNo(csNo)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (e *CsApply) GetApplyDetailList(csNo string, csType int) (csApplyDetailList []*models.CsApplyDetail, err error) {
	csApplyDetailService := CsApplyDetail{e.Service}
	err = csApplyDetailService.GetAll(&dto.CsApplyDetailGetPageReq{
		CsType: csType,
		CsNo:   csNo,
	}, &csApplyDetailList)
	if err != nil {
		return
	}
	return
}

// Insert 创建CsApply对象
func (e *CsApply) Insert(c *dto.CsApplyInsertReq) error {
	var err error
	var data models.CsApply
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("CsApplyService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改CsApply对象
func (e *CsApply) Update(c *dto.CsApplyUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.CsApply{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("CsApplyService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除CsApply
func (e *CsApply) Remove(d *dto.CsApplyDeleteReq, p *actions.DataPermission) error {
	var data models.CsApply

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveCsApply error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}

// Cancel 取消售后
func (e *CsApply) Cancel(c *dto.CsApplyCancelReq, p *actions.DataPermission) (err error) {
	var data = models.CsApply{}
	e.Orm.Where("cs_no = ?", c.CsNo).First(&data)

	if data.Id == 0 {
		return errors.New("没有这个售后单")
	}
	// 售后状态：0 => '待审核', 1 => '已确认', 2 => '已作废', 3 => '已取消', 10 => '待退货入库',11 => '待取消出库',99 => '已完结'
	// 待审核 待退货入库 待取消出库可以进行售后取消0, 10, 11
	var CsApplyCancelOkStatus = map[int]int{
		0:  0,
		10: 10,
		11: 11,
	}
	// 待审核 待退货入库 待取消出库可以进行售后取消0, 10, 11
	if _, ok := CsApplyCancelOkStatus[data.CsStatus]; !ok {
		return errors.New("售后单状态不可关闭")
	}

	// 判断是否有退货入库仓库的权限

	// 判断用户是否有货主的权限

	// 开始事务
	tx := e.Orm.Begin()
	// 如果是退货的 待退货入库状态下取消 =>  通知入库单售后已取消 因为退货生产的入库单作废
	if data.CsType == 0 && data.CsStatus == 10 {
		// 调用入库单作废方法
		err = serviceWc.CancelEntryForCancelCsOrder(tx, &dtoWc.StockEntryCancelForCancelCsOrderReq{
			CsNo:   c.CsNo,
			Remark: "售后单关闭",
			ControlBy: cModels.ControlBy{
				UpdateBy:     user.GetUserId(e.Orm.Statement.Context.(*gin.Context)),
				UpdateByName: user.GetUserName(e.Orm.Statement.Context.(*gin.Context)),
			},
		})
		if err != nil {
			// 遇到错误时回滚事务
			tx.Rollback()
			return errors.New("售后单取消失败" + err.Error())
		}
	}

	// 更新售后单状态为 已取消
	err = e.UpdateCsStatus(e.Orm, data.Id, data.CsNo, 3, &data, c.AuditReason)
	if err != nil {
		// 遇到错误时回滚事务
		tx.Rollback()
		return err
	}

	// 通知订单售后状态改为售后完结
	err = e.OrderAfterSales(e.Orm, data.OrderId, 99)
	if err != nil {
		// 遇到错误时回滚事务
		tx.Rollback()
		return errors.New("售后单取消失败")
	}

	// 否则，提交事务
	tx.Commit()
	return nil
}

// Undo 作废售后单
func (e *CsApply) Undo(c *dto.CsApplyCancelReq, p *actions.DataPermission) (err error) {
	var data = models.CsApply{}
	e.Orm.Where("cs_no = ?", c.CsNo).First(&data)

	if data.Id == 0 {
		return errors.New("没有这个售后单")
	}
	// 售后状态：0 => '待审核', 1 => '已确认', 2 => '已作废', 3 => '已取消', 10 => '待退货入库',11 => '待取消出库',99 => '已完结'
	if data.CsStatus != 0 {
		return errors.New("售后单状态不可拒绝")
	}

	// 判断是否有退货入库仓库的权限

	// 开始事务
	tx := e.Orm.Begin()

	// 更新售后单状态为 已拒绝
	err = e.UpdateCsStatus(e.Orm, data.Id, data.CsNo, 2, &data, c.AuditReason)
	if err != nil {
		// 遇到错误时回滚事务
		tx.Rollback()
		return err
	}

	// 通知订单售后状态改为售后完结
	err = e.OrderAfterSales(e.Orm, data.OrderId, 99)
	if err != nil {
		// 遇到错误时回滚事务
		tx.Rollback()
		return errors.New("售后单拒绝失败")
	}

	// 否则，提交事务
	tx.Commit()
	return nil
}

// Confirm 确认售后单(通过操作)
func (e *CsApply) Confirm(db *gorm.DB, csNo string, p *actions.DataPermission) (err error, callStockChangeEventErr bool) {
	if db == nil {
		db = e.Orm
	}
	csApplyM := models.CsApply{}
	data, err := csApplyM.GetByCsNo(db, csNo)
	if err != nil {
		return
	}
	if data.CsStatus != 0 {
		err = errors.New("售后单状态不可操作")
		return
	}

	csApplyDetailList, err := csApplyM.GetApplyDetailList(db, csNo, 0)
	if err != nil {
		return err, false
	}
	confirmProduct, err := csApplyM.ConfirmProduct(db, data.OrderId, csApplyDetailList, data.CsType)
	if err != nil {
		log.Info(err)
		err = errors.New("售后单状态不可操作.")
		return
	}

	// 售后类型：0-退货、1-换货、2-退款、3-发票问题、4-技术及资料支持、5-技术及资料支持、6-缺货少配件、7-售后维修、8-其他、9-订单取消、10-补发货（仅限ERP迁移老数据使用）
	// 驿站只有退货和取消

	// 事务
	// 开始事务
	tx := db.Begin()

	// 根据退货 取消不同 做订单和库存的处理

	productIsEnd := 0
	if confirmProduct.IsAll || confirmProduct.IsEnd {
		productIsEnd = 1
	}
	var orderIsEnd bool
	if confirmProduct.IsOrderEnd {
		orderIsEnd = true
	}

	// 退货
	if data.CsType == 0 {
		// 原售后状态变更逻辑：
		// 审核通过 暂时状态改为"待退货入库" => “等待退货入库”状态下，通知库存中心生成入库单，进行入库流程；
		// 后续不更新售后单状态  确认入库后，入库单调用售后完结方法售后单更新状态“已完结”；
		// 售后状态优化后：
		// 1.确认审核售后 =》 售后状态变更为已确认
		// 2.调用库存中心后(通知库存中心生成入库单，进行入库流程；)  =》 售后状态变更为  "待退货入库" （后续不更新售后单状态  确认入库后，入库单调用售后完结方法售后单更新状态“已完结”）

		// 1.确认审核售后 =》 售后状态变更为已确认
		err = e.UpdateCsStatus(tx, data.Id, data.CsNo, 1, nil, "")
		if err != nil {
			tx.Rollback()
			return err, false
		}

		// 调用库存中心
		err = e.callStockChangeEvent(tx, data.CsType, productIsEnd, &dto.CallStockChangeEventReq{
			OrderId:            data.OrderId,
			CsNo:               data.CsNo,
			WarehouseCode:      data.WarehouseCode,
			LogicWarehouseCode: data.LogicWarehouseCode,
			ProductList:        csApplyDetailList,
			VendorsId:          data.VendorId,
		})
		if err != nil {
			tx.Rollback()
			return
		}

		//2.调用库存中心后 =》 售后状态变更为  "待退货入库" （后续不更新售后单状态  确认入库后，入库单调用售后完结方法售后单更新状态“已完结”）
		err = e.UpdateCsStatus(tx, data.Id, data.CsNo, 10, nil, "")
		if err != nil {
			tx.Rollback()
			return err, false
		}
	} else if data.CsType == 9 { // 取消
		// 原售后状态变更逻辑：
		// 审核通过 暂时状态改为"等待取消出库" => 通知库存中心更新原出库单
		// 调用库存方法:成功=>后售后单更新状态“已完结” 失败售后单状态不更新
		// 售后状态优化后：
		// 1.确认审核售后 =》 售后状态变更为已确认
		// 2.调用库存中心后 成功=>后售后单更新状态“已完结”进行后续逻辑操作    失败 =》 2023-09-16 修改逻辑 取消”待取消入库“ 状态 修改为如果调用库存中心后失败后直接回滚

		var order models.OrderInfo
		err = db.Where("order_id = ?", data.OrderId).First(&order).Error
		// 订单未确认和缺货状态都是 未出库状态 这时候申请售后审核时不需要调用库存中心
		if order.OrderStatus != 5 && order.OrderStatus != 6 {
			// 调用库存中心
			err = e.callStockChangeEvent(db, data.CsType, productIsEnd, &dto.CallStockChangeEventReq{
				OrderId:            data.OrderId,
				CsNo:               data.CsNo,
				WarehouseCode:      data.WarehouseCode,
				LogicWarehouseCode: data.LogicWarehouseCode,
				ProductList:        csApplyDetailList,
				VendorsId:          data.VendorId,
			})
		}
		//err = nil
		if err != nil {
			tx.Rollback()

			e.Log.Errorf("售后单库存中心出库入库变更失败 error:%s", err)
			//售后单库存中心出库入库变更失败
			callStockChangeEventErr = true

			//// 2.调用库存中心后 失败 =》 售后状态变更为  "待取消入库"（后续无代码逻辑 直接return）(2023-09-16 这段代码废弃)
			//err = e.UpdateCsStatus(e.Orm, data.Id, data.CsNo, 11, nil, "")
			//if err != nil {
			//	return
			//}

			// 2023-09-16 修改逻辑 取消 ”待取消入库“ 状态 修改为如果 调用库存中心后失败后直接回滚
			return
		} else {
			// 2.调用库存中心后 成功=>后售后单更新状态“已完结”
			err = e.applyEnd(tx, data.Id, data.CsNo, csApplyDetailList, data)
			if err != nil {
				tx.Rollback()
				return
			}

			// 订单售后取消的操作（修改状态更新商品已取消数量等）
			err = e.OrderCancel(tx, data.OrderId, productIsEnd, csApplyDetailList, data, p, orderIsEnd)
			if err != nil {
				tx.Rollback()
				return err, false
			}
		}
	} else {
		err = errors.New("售后类型错误")
		return
	}

	tx.Commit()
	return nil, callStockChangeEventErr
}

/**
 * 售后退货完成
 *
 * @param csNo string
 * @return void
 */
func (e *CsApply) ReturnCompleted(tx *gorm.DB, csNo string) (err error) {

	data, err := e.GetByCsNo(csNo)
	if err != nil {
		return
	}

	if data.CsStatus != 10 {
		err = errors.New("售后单状态不可操作")
		return
	}
	applyDetailList, err := e.GetApplyDetailList(csNo, 0)
	if err != nil {
		return
	}
	orderId := data.OrderId

	// 确认商品是否正确 有没有购买数量等
	_, err = data.ConfirmProduct(tx, orderId, applyDetailList, data.CsType)
	if err != nil {
		return
	}

	// 售后类型：0-退货、1-换货、2-退款、3-发票问题、4-技术及资料支持、5-技术及资料支持、6-缺货少配件、7-售后维修、8-其他、9-订单取消、10-补发货（仅限ERP迁移老数据使用）
	// sxyz只有退货和取消

	//productIsEnd := 0
	//if confirmProduct.IsAll && confirmProduct.IsEnd {
	//	productIsEnd = 1
	//}

	// 订单修改退货状态
	err = e.OrderReturn(tx, orderId, applyDetailList)
	if err != nil {
		return err
	}

	err = e.applyEnd(tx, data.Id, csNo, applyDetailList, data)
	if err != nil {
		return
	}

	return

}

// OrderReturn 售后退货 操作订单状态
func (e *CsApply) OrderReturn(db *gorm.DB, orderId string, csApplyDetailList []*models.CsApplyDetail) (err error) {
	csApplyM := models.CsApply{}
	// 如果有在售后中的就不用改订单状态 直到没有一个在售后中的
	isOrderInAfterSales, err := csApplyM.IsOrderInAfterSales(db, orderId)
	log.Info("isOrderInAfterSales", isOrderInAfterSales)
	if !isOrderInAfterSales {
		err = e.UpdateOrder(db, orderId, 99)
		if err != nil {
			return
		}
	}

	// 订单商品表处理
	for _, detail := range csApplyDetailList {
		err = e.UpdateReturnOrderDetail(e.Orm, orderId, detail.GoodsId, detail.Quantity)
		if err != nil {
			return err
		}
	}

	return
}

// OrderCancel 售后取消 操作订单状态
func (e *CsApply) OrderCancel(db *gorm.DB, orderId string, isAll int, csApplyDetailList []*models.CsApplyDetail, data *models.CsApply, p *actions.DataPermission, orderIsEnd bool) (err error) {
	csApplyM := models.CsApply{}
	orderInfo := models.OrderInfo{}
	db.Model(&models.OrderInfo{}).Where("order_id = ?", orderId).First(&orderInfo)

	// 订单是否未确认
	var isOrderUnconfirmed bool
	if orderInfo.OrderStatus == 5 || orderInfo.OrderStatus == 6 {
		isOrderUnconfirmed = true
	}

	orderDetails := []*models.OrderDetail{}
	err = db.Model(models.OrderDetail{}).Where("order_id = ?", orderId).Find(&orderDetails).Error
	if err != nil {
		return
	}
	orderDetailsMap := make(map[string]*models.OrderDetail)
	for _, detail := range orderDetails {
		orderDetailsMap[detail.SkuCode] = detail
	}

	csApplyDetailListMap := map[string]*models.CsApplyDetail{}
	// 订单商品表处理
	for _, detail := range csApplyDetailList {
		err = e.UpdateCancelOrderDetail(db, orderId, detail.GoodsId, detail.Quantity)
		if err != nil {
			// 遇到错误时回滚事务
			return
		}
		csApplyDetailListMap[detail.SkuCode] = detail
		if isOrderUnconfirmed {
			// 缺货和待确认订单 售后取消需要释放库存
			// 申请售后的商品 重置锁定库存 重置商品数量
			if orderDetail, ok := orderDetailsMap[detail.SkuCode]; ok {
				// 订单原数量 - 申请售后的数量 - 原锁库数量
				unlock := orderDetail.Quantity - orderDetail.CancelQuantity - detail.Quantity - orderDetail.LockStock
				// 减少数量 恢复库存
				if unlock < 0 {
					err = modelsWc.UnLockStockInfoForOrder(db, -unlock, detail.GoodsId, data.WarehouseCode, data.OrderId, "后台订单修改商品数量恢复库存")
					if err != nil {
						return err
					}
					//	这里锁库数量要相关 LockStock - 解锁数量
					db.Model(orderDetail).Where("id = ?", orderDetail.Id).UpdateColumn("lock_stock", gorm.Expr("lock_stock - ?", -unlock))
				}
			}
		}
	}

	// 如果有在售后中的就不用改订单状态 直到没有一个在售后中的
	isOrderInAfterSales, err := csApplyM.IsOrderInAfterSales(db, orderId)
	log.Info("isOrderInAfterSales", isOrderInAfterSales)

	isSaveOrderInfo := false
	saveOrderInfo := models.OrderInfo{}

	// 缺货和待确认订单处理
	if isOrderUnconfirmed {
		// 部分取消并且 订单处于缺货状态下判断
		if isAll == 0 && orderInfo.OrderStatus == 6 {
			// 查询订单是否还是缺货状态
			isOutOfStock, err := e.GetOrderProductsIsOutOfStock(db, orderInfo.OrderId)
			if err != nil {
				return err
			}

			// 如果订单不是缺货状态的话 缺货订单订单状态改为5待确认
			if !isOutOfStock {
				isSaveOrderInfo = true
				saveOrderInfo = orderInfo
				saveOrderInfo.OrderStatus = 5
			}
		}
	}

	if !isOrderInAfterSales {
		isSaveOrderInfo = true
		if saveOrderInfo.Id == 0 {
			saveOrderInfo = orderInfo
		}
		saveOrderInfo.RmaStatus = 99

		// 如果是 全部取消和最后一次部分取消 就直接把订单状态改成取消
		if isAll == 1 {
			saveOrderInfo.OrderStatus = 9

			// 因为新增部分发货逻辑 新增判断 如果是部分发货的订单 全部取消后订单状态改为 1已发货
			if orderInfo.OrderStatus == 11 {
				saveOrderInfo.OrderStatus = 1
				// 订单已经全部售后完成 改成 已关闭状态
				if orderIsEnd {
					saveOrderInfo.OrderStatus = 10
				}
			}
			saveOrderInfo.CancelByName = user.GetUserName(e.Orm.Statement.Context.(*gin.Context))
			saveOrderInfo.CancelBy = user.GetUserId(e.Orm.Statement.Context.(*gin.Context))
		}

	}
	if isSaveOrderInfo {
		err = csApplyM.OrderUpdate(db, orderId, &orderInfo, &saveOrderInfo, 99)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetOrderProductsIsOutOfStock 检查订单内商品是否缺货 如果商品数量-取消数量 == 锁库数量就是不缺货了 反之还是缺货状态让他走补货脚本
func (e *CsApply) GetOrderProductsIsOutOfStock(tx *gorm.DB, orderId string) (isOutOfStock bool, err error) {
	order := &models.OrderInfo{}
	err = tx.Model(models.OrderInfo{}).Where("order_id = ?", orderId).First(&order).Error
	if err != nil {
		return
	}
	products := []*models.OrderDetail{}
	err = tx.Model(models.OrderDetail{}).Where("order_id = ?", orderId).Find(&products).Error
	if err != nil {
		return
	}

	// 全部满足：订单状态=待确认 /  不完全满足，订单状态=缺货
	for _, product := range products {
		if (product.Quantity - product.CancelQuantity) != product.LockStock {
			isOutOfStock = true
			break
		}
	}
	return
}

func (e *CsApply) UpdateOrder(db *gorm.DB, orderID string, rmaStatus int) error {
	result := db.Table(global.GetTenantOcDBNameWithDB(e.Orm)+".order_info").
		Where("order_id = ?", orderID).
		Updates(map[string]interface{}{
			"rma_status": rmaStatus,
		})

	if result.Error != nil {
		return result.Error
	}

	//订单更新日志..

	return nil
}

// UpdateCancelOrderDetail 订单商品表取消处理
func (e *CsApply) UpdateCancelOrderDetail(db *gorm.DB, orderID string, goodsID int, quantity int) error {
	err := db.Table(global.GetTenantOcDBNameWithDB(e.Orm)+".order_detail").
		Where("order_id = ? AND goods_id = ?", orderID, goodsID).
		Updates(map[string]interface{}{
			"final_quantity":  gorm.Expr("final_quantity - ?", quantity),
			"cancel_quantity": gorm.Expr("cancel_quantity + ?", quantity),
		}).Error
	return err
}

// UpdateReturnOrderDetail 订单商品表退货处理
func (e *CsApply) UpdateReturnOrderDetail(db *gorm.DB, orderID string, goodsID int, quantity int) error {
	err := db.Table(global.GetTenantOcDBNameWithDB(e.Orm)+".order_detail").
		Where("order_id = ? AND goods_id = ?", orderID, goodsID).
		Updates(map[string]interface{}{
			"final_quantity":  gorm.Expr("final_quantity - ?", quantity),
			"cancel_quantity": gorm.Expr("cancel_quantity + ?", quantity),
		}).Error
	return err
}

func (e *CsApply) applyEnd(db *gorm.DB, id int, csNo string, csApplyDetailList []*models.CsApplyDetail, data *models.CsApply) error {
	err := e.UpdateCsStatus(db, id, csNo, 99, nil, "")
	if err != nil {
		return err
	}
	var csApplyDetailAdds []*models.CsApplyDetail
	// 售后商品表添加实际售后的记录
	for i, detail := range csApplyDetailList {
		csApplyDetailAdds = append(csApplyDetailAdds, detail)
		csApplyDetailAdds[i].Id = 0
		csApplyDetailAdds[i].CsType = 1
	}

	err = db.Model(&models.CsApplyDetail{}).Create(csApplyDetailAdds).Error
	if err != nil {
		return err
	}

	if data.CsType == 0 || data.CsType == 9 {
		var order models.OrderInfo
		err = db.Where("order_id = ?", data.OrderId).First(&order).Error
		if err != nil {
			return err
		}
		if order.CreatedAt.Year() == time.Now().Year() && order.CreatedAt.Month() == time.Now().Month() {
			// 售后 退货取消完成后 调用预算记录更新方法
			// 释放预算
			var departmentBudgetM = modelsUc.DepartmentBudget{}
			log.Info(order.UserId, -data.ApplyPrice, time.Now().Format("200601"))
			err = departmentBudgetM.UpdateBudget(db, order.UserId, -data.ApplyPrice, time.Now().Format("200601"))
			//$this->model_department_budget->updateBudget($data['order_user_id'], $data['apply_price'] * (-1), date('Ym'));
		}
	}

	return nil
}
func (e *CsApply) callStockChangeEvent(db *gorm.DB, csType int, isAll int, callStockChangeEventReq *dto.CallStockChangeEventReq) (err error) {
	//stockOutboundService := serviceWc.StockOutbound{e.Service}
	tx := db
	// 退货
	if csType == 0 {
		// 退货入库单
		StockEntryProducts := []dtoWc.StockEntryProductsReq{}
		copier.Copy(&StockEntryProducts, callStockChangeEventReq.ProductList)
		_, err = serviceWc.CreateEntryForCsOrder(tx, &dtoWc.StockEntryInsertReq{
			SourceCode:         callStockChangeEventReq.CsNo,
			Remark:             "",
			WarehouseCode:      callStockChangeEventReq.WarehouseCode,
			LogicWarehouseCode: "",
			VendorId:           callStockChangeEventReq.VendorsId,
			ControlBy: cModels.ControlBy{
				UpdateBy:     user.GetUserId(e.Orm.Statement.Context.(*gin.Context)),
				UpdateByName: user.GetUserName(e.Orm.Statement.Context.(*gin.Context)),
			},
			StockEntryProducts: StockEntryProducts,
		})

	} else if csType == 9 { // 取消
		if isAll == 1 {
			csOrderProducts := []dtoWc.StockOutboundProductsPartCancelForCsOrderReq{}
			copier.Copy(&csOrderProducts, callStockChangeEventReq.ProductList)
			log.Info(&dtoWc.StockOutboundPartCancelForCsOrderReq{
				CsNo:            callStockChangeEventReq.CsNo,
				OrderId:         callStockChangeEventReq.OrderId,
				CsOrderProducts: csOrderProducts,
				Remark:          "",
				ControlBy: cModels.ControlBy{
					UpdateBy:     user.GetUserId(e.Orm.Statement.Context.(*gin.Context)),
					UpdateByName: user.GetUserName(e.Orm.Statement.Context.(*gin.Context)),
				},
			})

			// 订单取消全部时库存变更
			err = serviceWc.CancelOutboundForCsOrder(tx, &dtoWc.StockOutboundPartCancelForCsOrderReq{
				CsNo:            callStockChangeEventReq.CsNo,
				OrderId:         callStockChangeEventReq.OrderId,
				CsOrderProducts: csOrderProducts,
				Remark:          "",
				ControlBy: cModels.ControlBy{
					UpdateBy:     user.GetUserId(e.Orm.Statement.Context.(*gin.Context)),
					UpdateByName: user.GetUserName(e.Orm.Statement.Context.(*gin.Context)),
				},
			})
			//stockOutboundService.CancelOutboundForCsOrder
			//$status = $this->model_lion_post_stock_entry->createEntryForCancelAll($data);
		} else {
			csOrderProducts := []dtoWc.StockOutboundProductsPartCancelForCsOrderReq{}
			copier.Copy(&csOrderProducts, callStockChangeEventReq.ProductList)
			err = serviceWc.PartCancelOutboundForCsOrder(tx, &dtoWc.StockOutboundPartCancelForCsOrderReq{
				CsNo:            callStockChangeEventReq.CsNo,
				OrderId:         callStockChangeEventReq.OrderId,
				CsOrderProducts: csOrderProducts,
				ControlBy: cModels.ControlBy{
					UpdateBy:     user.GetUserId(e.Orm.Statement.Context.(*gin.Context)),
					UpdateByName: user.GetUserName(e.Orm.Statement.Context.(*gin.Context)),
				},
			})
			// 订单取消时库存变更
			//$status = $this->model_lion_post_stock_entry->createEntryForCancel($data);
		}
	}
	if err != nil {
		e.Log.Errorf("[售后-出入库]:%s \r\n", err)
		return err
	}
	return nil
}

// 订单售后状态改变 status 售后状态：0-无售后、1-售后处理中、2-售后已确认、99-售后处理完成
func (e *CsApply) OrderAfterSales(db *gorm.DB, orderId string, status int) error {
	var data models.CsApply
	// 如果有在售后中的就不用改订单状态 直到没有一个在售后中的
	err := db.Where("order_id = ?", orderId).Where("cs_status = 0").First(&data).Error
	if data.Id > 0 {
		// 如果未找到记录，则售后单不存在不用改变订单状态
		return nil
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// 如果出现其他错误，则返回错误信息
		return err
	}

	var order models.OrderInfo
	err = db.Where("order_id = ?", orderId).First(&order).Error
	if err != nil {
		return err
	}

	order.RmaStatus = status

	// 调用修改订单售后状态的方法
	err = data.OrderUpdate(db, orderId, nil, &order, status)
	if err != nil {
		return err
	}
	return nil
}

// IsOrderInAfterPendingReview 通过订单号查询 该订单是 售后待审核状态的
func (e *CsApply) IsOrderInAfterPendingReview(orderID string) (bool, error) {
	var count int64
	err := e.Orm.Model(&models.CsApply{}).
		Where("order_id = ?", orderID).
		Where("cs_status = 0").
		Limit(1).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// UpdateCsStatus 更新售后单状态
func (e *CsApply) UpdateCsStatus(db *gorm.DB, id int, csNo string, csStatus int, data *models.CsApply, auditReason string) (err error) {
	ocPrefix := global.GetTenantOcDBNameWithDB(db)

	csApplyCsStatusText, ok := dto.CsApplyCsStatus[csStatus]
	if !ok {
		return errors.New("错误的状态修改")
	}
	if data == nil {
		data = &models.CsApply{}
		db.First(&data, id)
	}
	data.CsStatus = csStatus
	data.AuditReason = auditReason // 审核原因
	data.UpdateBy = user.GetUserId(e.Orm.Statement.Context.(*gin.Context))
	data.UpdateByName = user.GetUserName(e.Orm.Statement.Context.(*gin.Context))

	saveDb := db.Table(ocPrefix+"."+data.TableName()).Where("id = ?", id).Where("cs_status in (?)", dto.CsApplyStatusUpdateOkList).Save(&data)
	if err = saveDb.Error; err != nil {
		e.Log.Errorf("CsApplyService Save error:%s \r\n", err)
		return err
	}
	if saveDb.RowsAffected == 0 {

		return errors.New("无权更新该数据")
	}

	// 需要调用售后日志Service
	var csApplyLogService CsApplyLog
	csApplyLogService.Service = e.Service

	err = csApplyLogService.AddLog(db, csNo, "售后单状态变更为:"+csApplyCsStatusText)
	if err != nil {
		return err
	}

	return nil
}

func (e *CsApply) GetSaleProducts(orderId string, csType int) (res *dto.CsApplyGetSaleProductsRes, err error) {
	orderInfo := models.OrderInfo{}
	e.Orm.Where("order_id", orderId).First(&orderInfo)
	if orderInfo.Id == 0 {
		err = errors.New("未查询到该订单，请确认订单号输入无误")
		return
	}
	//IS_FORMAL_ORDER 0 是否是正式订单(1:是, 0:否)
	if orderInfo.IsFormalOrder != 1 {
		err = errors.New("预订单不可创建售后申请单，请确认订单号输入无误！")
		return
	}

	products, err := e.GetOrderProducts(orderId)
	if err != nil {
		return
	}
	e.ProductsInit(orderId, csType, &products)

	warehouse := e.GetWcGetWarehousesByCode(e.Orm, orderInfo.WarehouseCode)
	res = &dto.CsApplyGetSaleProductsRes{
		Products: products,
		Warehouse: &dto.CsApplyWarehouseData{
			WarehouseCode: warehouse.WarehouseCode,
			WarehouseName: warehouse.WarehouseName,
		},
		IsStatements: e.GetIsStatements(orderId), // 订单是否已经对账
	}
	return
}

func (e *CsApply) ProductsInit(orderId string, csType int, products *[]*dto.CsApplyGetOrderProducts) {
	if products == nil || len(*products) == 0 {
		return
	}
	orderInfo := models.OrderInfo{}
	csApplyModel := models.CsApply{}
	e.Orm.Where("order_id", orderId).First(&orderInfo)
	allowCancel := e.VerifyOrderCancelStatusForAfter(orderInfo.OrderStatus)
	allowReturn := e.VerifyOrderReturnStatus(orderInfo.OrderStatus)
	returnProducts := csApplyModel.GetAfterReturnProductsBySaleId(e.Orm, orderId)

	for _, product := range *products {
		returnQantity := 0
		quantity := 0
		if csType == 0 { // 退货
			if allowReturn {
				if v, ok := returnProducts[product.SkuCode]; ok {
					returnQantity = v.Quantity
				}
				//$quantity = $SEND_QUANTITY - $return_quantity - $after_num_0;
				quantity = product.SendQuantity - returnQantity - product.AfterTypeNum[0]
				log.Info("product.SendQuantity, returnQantity, product.AfterTypeNum[0]")
				log.Info(product.SendQuantity, returnQantity, product.AfterTypeNum[0])
			}
		} else if csType == 9 { // 取消
			if allowCancel {
				//$quantity = $ORIGINAL_QUANTITY - $SEND_QUANTITY - $CANCEL_QUANTITY - $after_num_9; // - $return_quantity
				quantity = product.OriginalQuantity - product.SendQuantity - product.CancelQuantity - product.AfterTypeNum[9]
				log.Info(product.OriginalQuantity, product.SendQuantity, product.CancelQuantity, product.AfterTypeNum[9])
			}
		} else { // 其他类型:申请数量小于订单数量
			quantity = product.OriginalQuantity
		}

		product.AllowQuantity = quantity
		product.Quantity = 0
		log.Info("ProductsInit-end")
	}
	return
}

// VerifyOrderCancelStatusForAfter 校验订单状态是否可以申请售后取消 0-未发货、11-部分发货、1-已发货、2-部分收货、3-待评价、4-已评价、5-待确认、6、缺货、7-已签收、9-已取消、10-已关闭 11 部分发货
func (e *CsApply) VerifyOrderCancelStatusForAfter(orderStatus int) (verifyStatus bool) {
	if _, ok := models.OrderStatusCsApplyCancelOk[orderStatus]; ok {
		verifyStatus = true
	}
	return
}

// VerifyOrderReturnStatus 校验订单状态是否可以申请售后退货 0-未发货、11-部分发货、1-已发货、2-部分收货、3-待评价、4-已评价、5-待确认、6、缺货、7-已签收、9-已取消、10-已关闭 11 部分发货
func (e *CsApply) VerifyOrderReturnStatus(orderStatus int) (verifyStatus bool) {
	if _, ok := models.OrderStatusCsApplyReturnOk[orderStatus]; ok {
		verifyStatus = true
	}
	return
}

func (e *CsApply) GetOrderProducts(orderId string) (orderProducts []*dto.CsApplyGetOrderProducts, err error) {
	csApplyModel := models.CsApply{}
	//only_valid
	//根据用户是否是否登录以及是否绑定 company_id 来判断筛选条件
	orderInfo := models.OrderInfo{}
	e.Orm.Where("order_id", orderId).First(&orderInfo)
	//companyId := 0
	//if orderInfo.UserCompanyId != 0 {
	//	companyId = orderInfo.UserCompanyId
	//}

	query := e.Orm.Table("order_detail od").
		Select("od.*, oi.create_from,oi.order_status").
		Joins("LEFT JOIN order_info oi ON oi.order_id = od.order_id").
		Where("od.order_id = ?", orderId)

	query.Order("od.id").Find(&orderProducts)

	goodsIds := []int{}
	vendorIds := []int{}
	skuCodes := []string{}
	for _, products := range orderProducts {
		skuCodes = append(skuCodes, products.SkuCode)
		goodsIds = append(goodsIds, products.GoodsId)
		vendorIds = append(vendorIds, products.VendorId)
	}
	//// 获取商品信息
	//pcGoodsMap := e.GetPcGetGoodsBySku(e.Orm, goodsIds)
	// 获取货主信息
	wcVendorsMap := e.GetWcGetVendorsById(e.Orm, vendorIds)

	// 售后数
	afterApplyProducts := csApplyModel.GetAfterApplyProductsByTypeBySaleId(e.Orm, orderId)
	// 已退货数
	returnProducts := csApplyModel.GetAfterReturnProductsBySaleId(e.Orm, orderId)
	listSort := 0
	for _, product := range orderProducts {
		if _, ok := wcVendorsMap[product.VendorId]; !ok {
			err = errors.New("货主信息不存在")
			return
		}
		product.VendorsCode = wcVendorsMap[product.VendorId].Code
		product.VendorsName = wcVendorsMap[product.VendorId].NameZh
		//product.ProductNo = pcGoodsMap[product.GoodsId].ProductNo
		log.Info("_, product := range orderProducts")
		// 实际缺货数量
		product.ActualStock = product.Quantity - product.LockStock
		// 毛利率
		if product.SubTotalAmount > 0 {
			product.SkuProfit = ((product.SubTotalAmount - product.PurchasePrice*float64(product.Quantity)) / product.SubTotalAmount) * 100
		} else {
			product.SkuProfit = 0
		}

		product.AfterTypeNum = make(map[int]int)
		product.AfterTypeNum[0] = 0
		product.AfterTypeNum[9] = 0
		if v, ok := afterApplyProducts[product.SkuCode][0]; ok {
			product.AfterTypeNum[0] = v.Quantity
		}
		if v, ok := afterApplyProducts[product.SkuCode][9]; ok {
			product.AfterTypeNum[9] = v.Quantity
		}
		product.AfterNum = product.AfterTypeNum[0] + product.AfterTypeNum[9]
		if v, ok := returnProducts[product.SkuCode]; ok {
			product.ReturnNum = v.Quantity
		}

		listSort = listSort + 1
		product.ListSort = listSort
	}

	return
}

func (e *CsApply) GetIsStatements(orderId string) bool {
	statements := models.OrderToStatements{}
	e.Orm.Where("order_id = ?", orderId).First(&statements)
	if statements.Id == 0 {
		return false
	}
	return true
}

func (e *CsApply) GetPcGetProductBySku(db *gorm.DB, skuSlice []string) (productMap map[string]*dtoPc.InnerGetProductBySkuResp) {
	result := pcClient.ApiByDbContext(db).GetProductBySku(skuSlice)
	resultInfo := &struct {
		response.Response
		Data []dtoPc.InnerGetProductBySkuResp
	}{}
	result.Scan(resultInfo)

	productMap = make(map[string]*dtoPc.InnerGetProductBySkuResp)
	for _, product := range resultInfo.Data {
		productMap[product.SkuCode] = &product
	}

	return
}

func (e *CsApply) GetPcGetGoodsBySku(db *gorm.DB, goodIds []int) (goodsMap map[int]*modelsPc.Goods) {
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

func (e *CsApply) GetWcGetVendorsById(db *gorm.DB, ids []int) (dataMap map[int]*modelsWc.Vendors) {
	strSlice := make([]string, len(ids))
	for i, id := range ids {
		strSlice[i] = strconv.Itoa(id)
	}
	// 获取商品信息
	req := dtoWc.InnerVendorsGetListReq{
		Ids: strings.Join(strSlice, ","),
	}
	result := wcClient.ApiByDbContext(db).GetVendorList(req)

	resultInfo := &struct {
		response.Response
		Data []modelsWc.Vendors
	}{}
	result.Scan(resultInfo)
	dataMap = make(map[int]*modelsWc.Vendors)
	for _, data := range resultInfo.Data {
		copyData := data
		dataMap[copyData.Id] = &copyData
	}
	return
}

func (e *CsApply) GetWcGetWarehousesByCode(db *gorm.DB, code string) (data dtoWc.WarehouseGetPageResp) {
	// 获取商品信息
	req := dtoWc.InnerWarehouseGetListReq{
		WarehouseCode: code,
	}
	result := wcClient.ApiByDbContext(db).GetWarehouseList(req)

	resultInfo := &struct {
		response.Response
		Data []dtoWc.WarehouseGetPageResp
	}{}
	result.Scan(resultInfo)

	if len(resultInfo.Data) == 0 {
		return dtoWc.WarehouseGetPageResp{}
	}
	return resultInfo.Data[0]
}

/**
 * 查询订单是否可以售后（未发货的可以取消）（已发货的可以退货）
 *
 * @param orderID int         订单ID
 * @param orderData map[string]interface{} 订单数据
 * @param typ int             售后类型：0-退货、1-换货、2-退款、3-发票问题、4-技术及资料支持、5-技术及资料支持、6-缺货少配件、7-售后维修、8-其他、9-订单取消、10-补发货（仅限ERP迁移老数据使用）
 * @return map[string]interface{} orderData
 */
func (e *CsApply) OrderIsCanSale(orderID string, orderData *models.OrderInfo, applyType int) *models.OrderInfo {
	if orderData == nil {
		// Replace the line below with your own code to fetch order data
		// orderData = FetchOrderData(orderID)
		e.Orm.Where("order_id", orderID).First(&orderData)
	}

	if orderData.OrderStatus == 9 || orderData.OrderStatus == 10 {
		return nil
	}

	// 订单状态：0-未发货、1-已发货、2-部分收货、3-待评价、4-已评价、5-待确认、6、缺货、7-已签收、9-已取消

	// 取消 未发货的可以 和 部分发货可以申请取消
	if applyType == 9 && e.VerifyOrderCancelStatusForAfter(orderData.OrderStatus) {
		return orderData
	}

	// 退货 已发货 已签收 部分发货 可以申请退货
	if applyType == 0 && e.VerifyOrderReturnStatus(orderData.OrderStatus) {
		return orderData
	}

	return nil
}

func (e *CsApply) Add(req *dto.CsApplyInsertRequest, p *actions.DataPermission) (err error) {
	csApplyModel := models.CsApply{}
	//actualDetailList := make([]map[string]interface{}, 0)
	orderData := e.OrderIsCanSale(req.OrderId, nil, req.CsType)
	if orderData == nil {
		err = errors.New("该订单状态不能申请售后")
		return
	}
	req.IsStatements = 0
	if e.GetIsStatements(req.OrderId) {
		req.IsStatements = 1
	}
	warehouseData := e.GetWcGetWarehousesByCode(e.Orm, req.WarehouseCode)
	log.Info("p.AuthorityWarehouseId")
	log.Info(p.AuthorityWarehouseId)
	if p.AuthorityWarehouseId != "all" && !contains(strings.Split(p.AuthorityWarehouseId, ","), string(req.WarehouseCode)) {
		err = errors.New("没有退货入库仓库的权限，无法创建售后单")
		return
	}

	goodIds := []int{}
	//addDataList 根据货主分成多个订单插入
	addDataList := make(map[int]*dto.CsApplyInsertRequest)
	// 应该是遍历商品信息
	for _, product := range req.Products {
		// 数量大于0才说明有售后
		if product.Quantity > 0 {
			goodIds = append(goodIds, product.GoodsId)
			if _, ok := addDataList[product.VendorId]; !ok {
				addDataList[product.VendorId] = &dto.CsApplyInsertRequest{
					CsType:             req.CsType,
					OrderId:            req.OrderId,
					WarehouseCode:      req.WarehouseCode,
					WarehouseName:      warehouseData.WarehouseName,
					LogicWarehouseCode: warehouseData.WarehouseCode,
					IsStatements:       req.IsStatements,
					Products:           nil,
					CsNo:               csApplyModel.GenerateCsNo(),
					VendorsId:          product.VendorId,
					VendorsName:        product.VendorName,
					VendorsSkuCode:     product.VendorSkuCode,
					ApplyPrice:         0,
					ApplyQuantity:      0,
				}
			}
			product.CsNo = addDataList[product.VendorId].CsNo

			addDataList[product.VendorId].Products = append(addDataList[product.VendorId].Products, product)
			pric := utils.MulFloat64AndInt(product.SalePrice, product.Quantity)
			addDataList[product.VendorId].ApplyPrice = utils.AddFloat64(addDataList[product.VendorId].ApplyPrice, pric)
			addDataList[product.VendorId].ApplyQuantity = addDataList[product.VendorId].ApplyQuantity + product.Quantity
		}
	}
	// 确认商品信息是否正确 售后商品数量是否超出
	actualDetailList := []*models.CsApplyDetail{}
	copier.Copy(&actualDetailList, req.Products)
	status, err := csApplyModel.ConfirmProduct(e.Orm, req.OrderId, actualDetailList, req.CsType)
	if status == nil {
		return
	}

	if len(actualDetailList) == 0 {
		err = errors.New("请选择要售后的商品")
		return
	}
	// 查询订单有没有在售后中的商品 如果有就不能提交
	saleGoods, isSaleGoods := e.IsSaleGoods(goodIds, req.OrderId)
	if isSaleGoods {
		err = errors.New("SKU: " + strings.Join(saleGoods, ",") + "已存在售后，请先结束其他售后再申请")
		return
	}

	isAutoAudit := false
	if orderData.UserCompanyId > 0 {
		isAutoAudit = e.GetCompanyIsAfterAutoAudit(orderData.UserCompanyId, req.CsType)
	}

	tx := e.Orm.Begin()
	for _, add := range addDataList {

		// 添加售后时订单售后状态改成售后中
		err = e.OrderAfterSales(tx, add.OrderId, 1)
		if err != nil {
			tx.Rollback()
			return
		}

		// 插入售后表
		err = tx.Create(&models.CsApply{
			CsNo:               add.CsNo,
			CsType:             add.CsType,
			OrderId:            orderData.OrderId,
			CsStatus:           0,
			Telephone:          orderData.Telephone,
			UserId:             orderData.UserId,
			UserName:           orderData.UserName,
			WarehouseCode:      add.WarehouseCode,
			WarehouseName:      add.WarehouseName,
			LogicWarehouseCode: add.LogicWarehouseCode,
			CsSource:           "LMS",
			VendorId:           add.VendorsId,
			VendorName:         add.VendorsName,
			VendorSkuCode:      add.VendorsSkuCode,
			IsStatements:       add.IsStatements,
			ApplyPrice:         add.ApplyPrice,
			ApplyQuantity:      add.ApplyQuantity,
			ControlBy:          cModels.ControlBy{},
		}).Error
		if err != nil {
			tx.Rollback()
			return
		}

		// 插入售后详细商品表
		addDetail := []*models.CsApplyDetail{}
		copier.Copy(&addDetail, add.Products)
		log.Info(addDetail[0])
		err = tx.Create(&addDetail).Error
		if err != nil {
			tx.Rollback()
			return
		}
		// 添加log
		csApplyLogService := CsApplyLog{e.Service}
		csStatusText, _ := models.CsApplyStatusText[0]
		csApplyLogService.AddLog(e.Orm, add.CsNo, "售后单状态变更为:"+csStatusText)

	}
	tx.Commit()
	// 自动审核通过
	if isAutoAudit {
		for _, add := range addDataList {
			err, _ = e.Confirm(e.Orm, add.CsNo, p)
		}
	}
	return err
}

func (e *CsApply) IsSaleGoods(goodIds []int, orderId string) ([]string, bool) {
	if len(goodIds) == 0 {
		return nil, false
	}

	var details []models.CsApplyDetail
	err := e.Orm.Table("cs_apply a").
		Select("cs_apply_detail.sku_code").
		Joins("JOIN cs_apply_detail cs_apply_detail ON cs_apply_detail.cs_no = a.cs_no").
		Where("a.cs_status NOT IN (2, 3, 99)").
		Where("cs_apply_detail.goods_id IN ?", goodIds).
		Where("a.order_id = ?", orderId).
		Find(&details).Error

	if err != nil {
		return nil, false
	}

	skuCodes := make([]string, len(details))
	for i, detail := range details {
		skuCodes[i] = detail.SkuCode
	}

	return skuCodes, len(skuCodes) > 0
}

func (e *CsApply) GetCompanyIsAfterAutoAudit(id int, afterType int) bool {
	companyInfo := e.GetUcGetCompanyById(id)
	if companyInfo.AfterAutoAudit == "" {
		return false
	}

	afterAutoAuditSlice := strings.Split(companyInfo.AfterAutoAudit, ",")
	for _, value := range afterAutoAuditSlice {
		intValue, err := strconv.Atoi(value)
		if err == nil && intValue == afterType {
			return true
		}
	}
	return false
}

// GetPage 获取CsApply列表
func (e *CsApply) GetOrderInfoByCsNo(c *dto.GetOrderInfoByCsReq) (*models.OrderInfo, error) {
	var model models.CsApply
	var data models.OrderInfo

	e.Orm.Model(&model).
		Select("o.*").
		Where("a.cs_no = ?", c.CsNo).
		Joins("left join order_info o on o.order_id = cs_apply.order_id").
		First(&data)

	return &data, nil
}

func contains(slice []string, element string) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}

func (e *CsApply) GetUcGetCompanyById(id int) modelsUc.CompanyInfo {
	// 公司名称
	companyResult := ucClient.ApiByDbContext(e.Orm).GetCompanyInfoById(id)
	companyResultInfo := &struct {
		response.Response
		Data modelsUc.CompanyInfo
	}{}
	companyResult.Scan(companyResultInfo)

	return companyResultInfo.Data
}
