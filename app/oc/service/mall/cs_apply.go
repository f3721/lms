package mall

import (
	"errors"
	modelsPc "go-admin/app/pc/models"
	serviceWc "go-admin/app/wc/service/admin"
	dtoWc "go-admin/app/wc/service/admin/dto"
	pcClient "go-admin/common/client/pc"
	cModels "go-admin/common/models"
	"go-admin/common/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/oc/models"
	"go-admin/app/oc/service/mall/dto"
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

// GetPage 获取CsApply列表
func (e *CsApply) GetPageDetails(c *dto.CsApplyGetPageReq, p *actions.DataPermission, list *[]*dto.CsApplyGetPageGroupByOrderCsApply, count *int64) error {
	var err error
	var data models.CsApply

	err = e.Orm.Model(&data).
		//Select("cs_apply.*").
		Where("order_info.user_id = ?", c.UserId).
		Scopes(
			func(db *gorm.DB) *gorm.DB {
				if c.CsStatus != "" {
					db = db.Where("cs_apply.cs_status in ?", utils.Split(c.CsStatus))
				}
				if c.FilterKeyword != "" {
					db = db.Where("cs_apply.order_id LIKE ? OR cs_apply_detail.sku_code = ? OR cs_apply_detail.product_name LIKE ? OR cs_apply_detail.product_no LIKE ?",
						"%"+c.FilterKeyword+"%", "%"+c.FilterKeyword+"%", c.FilterKeyword, "%"+c.FilterKeyword+"%")
				}
				return db
			},
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
		).
		Joins("left JOIN order_info ON cs_apply.order_id = order_info.order_id").
		Joins("left JOIN cs_apply_detail ON cs_apply_detail.cs_no = cs_apply.cs_no").
		Preload("CsApplyDetail", "cs_type = 0").
		Group("cs_apply.cs_no").
		Find(&list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("CsApplyService GetPage error:%s \r\n", err)
		return err
	}
	csNoList := []string{}
	for _, apply := range *list {
		csNoList = append(csNoList, apply.CsNo)
	}

	finalQuantityMap := e.GetGoodsQuantityMap(csNoList, 1)
	for _, apply := range *list {
		for _, detail := range apply.CsApplyDetail {
			detail.ApplyQuantity = detail.Quantity
			if _, ok := finalQuantityMap[apply.CsNo]; ok {
				if finalQuantity, ok := finalQuantityMap[apply.CsNo][detail.GoodsId]; ok {
					detail.FinalQuantity = finalQuantity
				}
			}
		}
	}
	return nil
}

func (e *CsApply) GetGoodsQuantityMap(csNos []string, csType int) map[string]map[int]int {
	goodsQuantity := []dto.CsApplyGoodsQuantity{}
	e.Orm.Model(models.CsApplyDetail{}).
		Where("cs_no in ? and cs_type = ?", csNos, csType).
		Select("quantity,cs_no,goods_id").
		Find(&goodsQuantity)
	quantityMap := make(map[string]map[int]int)
	for _, quantity := range goodsQuantity {
		if _, ok := quantityMap[quantity.CsNo]; !ok {
			quantityMap[quantity.CsNo] = make(map[int]int)
		}
		quantityMap[quantity.CsNo][quantity.GoodsId] = quantity.Quantity
	}
	return quantityMap
}

// GetPage 获取CsApply列表
func (e *CsApply) GetPageGroupByOrder(c *dto.CsApplyGetPageReq, p *actions.DataPermission, list *[]*dto.CsApplyGetPageGroupByOrder, count *int64) error {
	var err error
	var data models.OrderInfo

	err = e.Orm.Model(&data).
		Select("order_info.order_id,order_info.created_at,order_info.contract_no").
		Where("order_info.user_id = ?", c.UserId).
		Scopes(
			func(db *gorm.DB) *gorm.DB {
				if c.CsStatus != "" {
					db = db.Where("cs_apply.cs_status in ?", utils.Split(c.CsStatus))
				}
				if c.FilterKeyword != "" {
					db = db.Where("order_info.order_id LIKE ? OR order_info.contract_no LIKE ? OR cs_apply_detail.sku_code = ? OR cs_apply_detail.product_name LIKE ? OR cs_apply_detail.product_no LIKE ?",
						"%"+c.FilterKeyword+"%", "%"+c.FilterKeyword+"%", c.FilterKeyword, "%"+c.FilterKeyword+"%", "%"+c.FilterKeyword+"%")
				}

				return db
			},
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
		).
		Joins("JOIN cs_apply ON cs_apply.order_id = order_info.order_id").
		Joins("JOIN cs_apply_detail ON cs_apply_detail.cs_no = cs_apply.cs_no").
		Preload("CsApply", func(db *gorm.DB) *gorm.DB {
			if c.CsStatus != "" {
				return db.Where("cs_status in ?", utils.Split(c.CsStatus))
			}
			return db
		}).
		Preload("CsApply.CsApplyDetail", "cs_type = 0").
		Group("order_info.order_id").
		Find(&list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("CsApplyService GetPage error:%s \r\n", err)
		return err
	}

	csNoList := []string{}
	for _, order := range *list {
		for _, apply := range order.CsApply {
			csNoList = append(csNoList, apply.CsNo)
		}
	}
	finalQuantityMap := e.GetGoodsQuantityMap(csNoList, 1)

	for _, order := range *list {
		for _, apply := range order.CsApply {
			for _, detail := range apply.CsApplyDetail {
				detail.ApplyQuantity = detail.Quantity
				if _, ok := finalQuantityMap[apply.CsNo]; ok {
					if finalQuantity, ok := finalQuantityMap[apply.CsNo][detail.GoodsId]; ok {
						detail.FinalQuantity = finalQuantity
					}
				}
			}
		}
	}
	return nil
}

func (e *CsApply) GetAfterOrders(c *dto.CsApplyGetPageReq, list *[]*dto.CsApplyGetPageGroupByOrder, count *int64) error {
	var data models.OrderInfo

	err := e.Orm.Model(&data).
		Select("order_id,create_at,contract_no").
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
		).
		Find(&list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("CsApplyService GetPage error:%s \r\n", err)
		return err
	}

	return nil
}

// Cancel 取消售后
func (e *CsApply) Cancel(c *dto.CsApplyCancelReq, p *actions.DataPermission) (err error) {
	var data = models.CsApply{}
	e.Orm.Where("cs_no = ?", c.CsNo).Where("user_id = ?", c.UserId).First(&data)

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
	err = data.UpdateCsStatus(e.Orm, data.Id, data.CsNo, 3, &data, c.AuditReason)
	if err != nil {
		// 遇到错误时回滚事务
		tx.Rollback()
		return err
	}

	// 通知订单售后状态改为售后完结
	err = data.OrderAfterSales(e.Orm, data.OrderId, 99)
	if err != nil {
		// 遇到错误时回滚事务
		tx.Rollback()
		return errors.New("售后单取消失败")
	}

	// 否则，提交事务
	tx.Commit()
	return nil
}

// Get 获取CsApply对象
func (e *CsApply) Get(d *dto.CsApplyGetReq, p *actions.DataPermission, apply *dto.CsApplyGetPageGroupByOrderCsApply) error {
	var data models.CsApply

	err := e.Orm.Debug().Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		Preload("CsApplyDetail", "cs_type = 0").
		First(apply, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetCsApply error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}

	csNoList := []string{}
	csNoList = append(csNoList, apply.CsNo)

	var goodsIds []int
	for _, detail := range apply.CsApplyDetail {
		goodsIds = append(goodsIds, detail.GoodsId)

	}

	goodsMap := e.getGoodsMap(goodsIds)

	finalQuantityMap := e.GetGoodsQuantityMap(csNoList, 1)
	for _, detail := range apply.CsApplyDetail {
		detail.ApplyQuantity = detail.Quantity
		detail.Moq = 1
		if goods, ok := goodsMap[detail.GoodsId]; ok {
			detail.Moq = goods.Product.SalesMoq
		}
		if _, ok := finalQuantityMap[apply.CsNo]; ok {
			if finalQuantity, ok := finalQuantityMap[apply.CsNo][detail.GoodsId]; ok {
				detail.FinalQuantity = finalQuantity
			}
		}
	}
	return nil
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

func (e *CsApply) getGoodsMap(goodIds []int) map[int]modelsPc.Goods {
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

	return goodsMap
}
