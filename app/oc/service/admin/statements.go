package admin

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	modelsUc "go-admin/app/uc/models"
	dtoUc "go-admin/app/uc/service/admin/dto"
	ucClient "go-admin/common/client/uc"
	"go-admin/common/global"
	"strconv"
	"time"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/oc/models"
	"go-admin/app/oc/service/admin/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type Statements struct {
	service.Service
}

type OrderToStatements struct {
	service.Service
}

// GetPage 获取Statements列表
func (e *Statements) GetPage(c *dto.StatementsGetPageReq, p *actions.DataPermission, list *[]dto.StatementsGetPageResp, count *int64) error {
	var err error
	var data models.Statements

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
			actions.SysUserPermission(data.TableName(), p, 1),
		).
		Scan(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("StatementsService GetPage error:%s \r\n", err)
		return err
	}

	// 公司名称
	companyResult := ucClient.ApiByDbContext(e.Orm).GetCompanyByIds("")
	companyResultInfo := &struct {
		response.Response
		Data struct {
			response.Page
			List []dtoUc.CompanyInfoGetSelectPageData
		}
	}{}
	companyResult.Scan(companyResultInfo)
	companyMap := make(map[int]string, len(companyResultInfo.Data.List))
	for _, company := range companyResultInfo.Data.List {
		companyMap[company.Id] = company.CompanyName
	}

	tmpList := *list
	for i, v := range tmpList {
		if _, ok := companyMap[v.CompanyId]; ok {
			tmpList[i].CompanyName = companyMap[v.CompanyId]
		}
	}
	*list = tmpList

	return nil
}

// Get 获取Statements对象
func (e *Statements) Get(d *dto.StatementsGetReq, p *actions.DataPermission, model *dto.StatementsGetPageResp) error {
	var data models.Statements

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetStatements error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}

	// 公司名称
	companyResult := ucClient.ApiByDbContext(e.Orm).GetCompanyByIds(strconv.Itoa(model.CompanyId))
	companyResultInfo := &struct {
		response.Response
		Data struct {
			response.Page
			List []dtoUc.CompanyInfoGetSelectPageData
		}
	}{}
	companyResult.Scan(companyResultInfo)
	if len(companyResultInfo.Data.List) == 1 {
		model.CompanyName = companyResultInfo.Data.List[0].CompanyName
	}
	return nil
}

// GetDetails 获取Statements对象
func (e *OrderToStatements) GetDetails(d *dto.OrderToStatementsGetReq, p *actions.DataPermission, list *[]dto.OrderToStatementsListResp, count *int64, switchStatus string, exportStatus bool) error {
	var data models.OrderToStatements
	pcPrefix := global.GetTenantPcDBNameWithDB(e.Orm)
	wcPrefix := global.GetTenantWcDBNameWithDB(e.Orm)
	ucPrefix := global.GetTenantUcDBNameWithDB(e.Orm)
	query := e.Orm.Model(&data).
		Select("order_to_statements.id, oi.order_id, oi.final_total_amount, u_i.company_department_id, p_c_d.id parent_company_department_id,"+
			" c_d.name company_department_name, p_c_d.name parent_company_department_name, oi.create_from, oi.user_id, u_i.user_name, u_i.user_phone,"+
			" od.sku_code, od.product_name, od.brand_name, od.product_model,od.sale_price, od.final_quantity, od.final_sub_total_amount, v.name_zh vendor_name,"+
			" g.supplier_sku_code, g.product_no, if('"+switchStatus+"'=1,c.name,null) as sku_classification_name,  if('"+switchStatus+"'=1,u_p_a.pay_account,null) as pay_account").
		Joins("left join order_info oi on order_to_statements.order_id = oi.order_id").
		Joins("left join order_detail od on od.order_id = oi.order_id").
		Joins("left join statements s on order_to_statements.statements_id = s.id").
		Joins("left join "+pcPrefix+".goods g on od.goods_id = g.id").
		Joins("left join "+wcPrefix+".vendors v on od.vendor_id = v.id").
		Joins("left join "+ucPrefix+".user_info u_i on oi.user_id = u_i.id").
		Joins("left join "+ucPrefix+".company_department c_d on c_d.id = u_i.company_department_id").
		Joins("left join "+ucPrefix+".company_department p_c_d on p_c_d.id = c_d.f_id").
		Joins("left join "+ucPrefix+".sku_classification s_c on s_c.company_id = c_d.company_id and s_c.sku_code = od.sku_code and s_c.status = 1").
		Joins("left join "+ucPrefix+".classification c on c.company_id = c_d.company_id and s_c.classification_id = c.id").
		Joins("left join "+ucPrefix+".user_pay_account u_p_a on u_p_a.company_id = c_d.company_id and u_p_a.user_id = oi.user_id and (u_p_a.classification_id = c.id or u_p_a.classification_id = -1)").
		Where("s.id = ?", d.Id)

	var err error
	if exportStatus == false {
		err = query.Scopes(
			cDto.Paginate(d.GetPageSize(), d.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).Find(list).Limit(-1).Offset(-1).Count(count).Error
	} else {
		err = query.Scopes(
			actions.Permission(data.TableName(), p),
		).Find(list).Error
	}

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetStatements error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}

	return nil
}

func (e *OrderToStatements) GetCompanySwitchInfo(statementsId int, companySkuSwitch string) (companyIndividualitySwitch modelsUc.CompanyIndividualitySwitch, err error) {
	var statements models.Statements
	statementInfo := new(models.Statements)
	err = e.Orm.Model(&statements).
		First(statementInfo, statementsId).
		Error
	err = companyIndividualitySwitch.GetRowByCompanyId(e.Orm, companySkuSwitch, statementInfo.CompanyId)
	if err != nil {
		return
	}

	return
}

// InitStatements 根据对账日 生成对账单
// 条件：订单最后修改时间在账单时间范围内、订单已签收、无正在进行中的售后、未对账
func (e *Statements) InitStatements(c *gin.Context) (echoMsg string, err error) {
	tenants := global.GetTenants()
	if len(tenants) <= 0 {
		err = errors.New("暂无租户")
		return
	}

	now := time.Now()
	day := now.Day()
	startTime := time.Date(now.Year(), now.Month()-1, day, 0, 4, 59, 0, now.Location())
	endTime := time.Date(now.Year(), now.Month(), day, 0, 5, 0, 0, now.Location())
	echoMsg = "生成对账单脚本执行开始："
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
		companyList := modelCompany.GetRowsByCondition(tenantDb, "company_status = 1 and reconciliation_day = ?", day)
		if len(companyList) == 0 {
			echoMsg = fmt.Sprintf("%s 无%v日对账的公司", echoMsg, day)
			continue
		}
		type tmpStruct struct {
			OrderId string
			CaId    int
		}

		for _, company := range companyList {
			var list []tmpStruct
			err = tenantDb.Table("order_info t").
				Select("t.order_id, ca.id ca_id").
				Joins("left join cs_apply ca on t.order_id = ca.order_id").
				Joins("left join order_to_statements ots on ots.order_id = t.order_id").
				Where("t.order_status = 7 and (ca.id is null or ca.cs_status in (3,99)) and ots.order_id is null").
				//Where("t.updated_at BETWEEN ? AND ?", startTime, endTime).
				Where("t.confirm_order_receipt_time < ?", endTime).
				Where("t.user_company_id = ?", company.Id).
				Group("t.order_id").
				Find(&list).Error
			if err != nil || len(list) == 0 {
				continue
			}
			var orderIds []string
			var modelOrderDetail models.OrderDetail
			var data models.Statements
			totalAmount := 0.00
			for _, v := range list {
				// 订单的所有售后单已完成 才算售后完成
				var count1, count2 int64
				tenantDb.Model(&models.CsApply{}).Where("order_id = ?", v.OrderId).Count(&count1)
				tenantDb.Model(&models.CsApply{}).Where("order_id = ? and cs_status in (2, 3, 99)", v.OrderId).Count(&count2)
				if count1 != count2 {
					continue
				}
				orderIds = append(orderIds, v.OrderId)
				totalAmount = totalAmount + modelOrderDetail.SetFinalQuantityAndAmountAndGetFinalTotalAmount(tenantDb, v.OrderId)
			}
			data.CompanyId = company.Id
			data.TotalAmount = totalAmount
			data.StartTime = startTime
			data.EndTime = endTime
			data.OrderCount = len(orderIds)
			var productCount int64
			tenantDb.Model(&models.OrderDetail{}).Where("order_id in ?", orderIds).Count(&productCount)
			data.ProductCount = int(productCount)
			err = tenantDb.Create(&data).Error
			if err != nil {
				echoMsg = fmt.Sprintf("%s 公司[%v]fail:%s | ", echoMsg, company.Id, err.Error())
				continue
			}
			data.StatementsNo = "AR" + now.Format("20060102") + fmt.Sprintf("%04d", data.Id)
			err = tenantDb.Save(&data).Error
			if err != nil {
				echoMsg = fmt.Sprintf("%s 公司[%v]fail:%s | ", echoMsg, company.Id, err.Error())
				continue
			}
			for _, orderId := range orderIds {
				ots := models.OrderToStatements{
					OrderId:      orderId,
					StatementsId: data.Id,
				}
				tenantDb.Save(&ots)
			}
			echoMsg = fmt.Sprintf("%s 公司[%v]success | ", echoMsg, company.Id)
			echoMsg = fmt.Sprintf("%s\r\n", echoMsg)
		}
	}
	echoMsg = fmt.Sprintf("%s\r\n 脚本执行完毕", echoMsg)

	return echoMsg, nil
}
