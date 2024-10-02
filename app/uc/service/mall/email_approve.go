package mall

import (
	"errors"
	"github.com/jinzhu/copier"
	"github.com/samber/lo"
	modelsCommon "go-admin/common/models"
	"strconv"
	"strings"
	"text/template"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/uc/models"
	"go-admin/app/uc/service/mall/dto"
	"go-admin/common/actions"
)

type EmailApprove struct {
	service.Service
}

func EmailApprovePermission(companyId int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db.Where("email_approve.company_id = ?", companyId)
		return db
	}
}

// GetPage 获取EmailApprove列表

func (e *EmailApprove) GetPage(c *dto.EmailApproveGetPageReq, p *actions.DataPermission, list *[]models.EmailApproveListCustomer, count *int64) error {
	var err error
	err = e.GetList(c, list, count)
	if err != nil {
		e.Log.Errorf("EmailApproveService GetPage error:%s \r\n", err)
		return err
	}
	FormatGetPage(e.Orm, list)
	return nil
}

func FormatGetPage(tx *gorm.DB, list *[]models.EmailApproveListCustomer) {
	// 格式化 用户名称
	userIds := []int{}
	for _, item := range *list {
		itemProcessIds := SplitUserIdsBySep(item.Process)
		itemPersonIds := SplitUserIdsBySep(item.Person)
		userIds = append(userIds, itemProcessIds...)
		userIds = append(userIds, itemPersonIds...)
	}
	userIds = lo.Uniq(userIds)
	userInfo := &models.UserInfo{}
	userNameMap, _ := userInfo.GetFormatNameUsersByIds(tx, userIds)
	for index, item := range *list {
		(*list)[index].Process = JoinProcessForApproveList(item.Process, userNameMap)
		(*list)[index].Person = JoinPersonForApproveList(item.Person, userNameMap)
	}
}

func JoinPersonForApproveList(person string, userNameMap map[int]string) string {
	outStr := ""
	personList := strings.Split(person, ",")
	for _, item := range personList {
		val, _ := strconv.Atoi(item)
		outStr = userNameMap[val] + " "
	}
	return strings.TrimRight(outStr, " ")
}

func JoinProcessForApproveList(process string, userNameMap map[int]string) string {
	outStr := ""
	processList := strings.Split(process, ",")
	for _, item := range processList {
		itemUserIds := []int{}
		itemUserIdsStr := ""
		if strings.Contains(item, "、") {
			itemUserIds = lo.Map(strings.Split(item, "、"), func(str string, _ int) int {
				val, _ := strconv.Atoi(str)
				return val
			})
			for _, userId := range itemUserIds {
				itemUserIdsStr += userNameMap[userId] + "、"
			}
			itemUserIdsStr = strings.TrimRight(itemUserIdsStr, "、")
		} else if strings.Contains(item, "/") {
			itemUserIds = lo.Map(strings.Split(item, "/"), func(str string, _ int) int {
				val, _ := strconv.Atoi(str)
				return val
			})
			for _, userId := range itemUserIds {
				itemUserIdsStr += userNameMap[userId] + "/"
			}
			itemUserIdsStr = strings.TrimRight(itemUserIdsStr, "/")
		} else {
			val, _ := strconv.Atoi(item)
			itemUserIdsStr = userNameMap[val]
		}
		outStr += itemUserIdsStr + " > "
	}
	return strings.TrimRight(outStr, " > ")
}

// 根据",、/"分割符 分割用户id

func SplitUserIdsBySep(userIds string) []int {
	splitUserIds := strings.FieldsFunc(userIds, func(r rune) bool {
		if r == ',' || r == '、' || r == '/' {
			return true
		} else {
			return false
		}
	})
	return lo.Map(splitUserIds, func(item string, _ int) int {
		val, _ := strconv.Atoi(item)
		return val
	})
}

func (e *EmailApprove) GetList(c *dto.EmailApproveGetPageReq, list *[]models.EmailApproveListCustomer, count *int64) error {
	whereStr := " WHERE ea.company_id = " + strconv.Itoa(c.CompanyId)

	processSql := `SELECT GROUP_CONCAT(t.PROCESS SEPARATOR ',') process,
                            t.priority,
                            t.id,
                            t.approve_status,
                            t.process_user_name,
                            t.process_login_name,
                            t.process_user_phone,
                            t.process_user_email
                            FROM 
                            (SELECT ea.id,ead.priority,ea.approve_status,
                          CASE
                                ead.approve_rank_type
                                WHEN 1 THEN
                                GROUP_CONCAT(ui.id SEPARATOR ',' )
                                WHEN 2 THEN
                                GROUP_CONCAT(ui.id ORDER BY ead.id SEPARATOR '、' ) 
                                ELSE GROUP_CONCAT( ui.id ORDER BY ead.id SEPARATOR '/' )
                          END PROCESS,
                          CASE
                                ead.approve_rank_type
                                WHEN 1 THEN
                                GROUP_CONCAT(ui.user_name SEPARATOR ',' )
                                WHEN 2 THEN
                                GROUP_CONCAT(ui.user_name SEPARATOR '、' ) 
                                ELSE GROUP_CONCAT( ui.user_name SEPARATOR '/' )
                          END process_user_name,
                          CASE
                                ead.approve_rank_type
                                WHEN 1 THEN
                                GROUP_CONCAT(ui.login_name SEPARATOR ',' )
                                WHEN 2 THEN
                                GROUP_CONCAT(ui.login_name SEPARATOR '、' ) 
                                ELSE GROUP_CONCAT( ui.login_name SEPARATOR '/' )
                          END process_login_name,
                          CASE
                                ead.approve_rank_type
                                WHEN 1 THEN
                                GROUP_CONCAT(ui.user_phone SEPARATOR ',' )
                                WHEN 2 THEN
                                GROUP_CONCAT(ui.user_phone SEPARATOR '、' ) 
                                ELSE GROUP_CONCAT( ui.user_phone SEPARATOR '/' )
                          END process_user_phone,
                          CASE
                                ead.approve_rank_type
                                WHEN 1 THEN
                                GROUP_CONCAT(ui.user_email SEPARATOR ',' )
                                WHEN 2 THEN
                                GROUP_CONCAT(ui.user_email SEPARATOR '、' ) 
                                ELSE GROUP_CONCAT( ui.user_email SEPARATOR '/' )
                          END process_user_email
                          FROM email_approve ea
                          LEFT JOIN email_approve_detail ead ON ea.id = ead.approve_id
                          LEFT JOIN user_info ui on ead.user_id=ui.id` + whereStr + ` AND ead.approve_detail_status = 1 AND ead.deleted_at IS NULL 
                          GROUP BY ead.priority,ea.id ORDER BY ead.priority ASC ) t GROUP BY t.id`

	personSql := `SELECT ea.id,GROUP_CONCAT(ui.id) person,
                         GROUP_CONCAT(ui.user_name ) person_user_name,
                         GROUP_CONCAT(ui.login_name ) person_login_name,
                         GROUP_CONCAT(ui.user_phone ) person_user_phone,
                         GROUP_CONCAT(ui.user_email ) person_user_email
                         FROM email_approve ea
                         LEFT JOIN user_approve ua ON ea.id = ua.approve_id
                         LEFT JOIN user_info ui on ua.user_id=ui.id` + whereStr + ` AND ua.status=1 AND ua.deleted_at IS NULL 
                         GROUP BY ea.id ORDER BY ea.id ASC`

	sql := `SELECT ea.id, aa.process, bb.person,aa.priority,aa.approve_status FROM email_approve ea
	INNER JOIN (` + processSql + `) aa ON ea.id = aa.id
	INNER JOIN (` + personSql + `) bb ON ea.id = bb.id  WHERE ea.approve_status = 1 `

	sqlCount := `SELECT count(ea.id) as count FROM email_approve ea
	INNER JOIN (` + processSql + `) aa ON ea.id = aa.id
	INNER JOIN (` + personSql + `) bb ON ea.id = bb.id  WHERE ea.approve_status = 1 `

	if c.FilterProcess != "" {
		filterProcess := template.HTMLEscapeString(c.FilterProcess)
		sqlFilterProcess := " And (aa.process LIKE '%?%' OR aa.process_user_name LIKE '%?%' OR aa.process_login_name LIKE '%?%' OR aa.process_user_phone LIKE '%?%' OR aa.process_user_email LIKE '%?%') "
		sql += strings.ReplaceAll(sqlFilterProcess, "?", filterProcess)
		sqlCount += strings.ReplaceAll(sqlFilterProcess, "?", filterProcess)
	}

	if c.FilterPerson != "" {
		filterPerson := template.HTMLEscapeString(c.FilterPerson)
		sqlFilterPerson := " And (bb.person LIKE '%?%' OR bb.person_user_name LIKE '%?%' OR bb.person_login_name LIKE '%?%' OR bb.person_user_phone LIKE '%?%' OR bb.person_user_email LIKE '%?%') "
		sql += strings.ReplaceAll(sqlFilterPerson, "?", filterPerson)
		sqlCount += strings.ReplaceAll(sqlFilterPerson, "?", filterPerson)
	}
	err := e.Orm.Raw(sqlCount).Scan(count).Error
	if err != nil {
		return err
	}
	sql += " ORDER BY ea.id DESC "
	pageSize := c.GetPageSize()
	offset := (c.GetPageIndex() - 1) * pageSize
	sql += " Limit " + strconv.Itoa(offset) + "," + strconv.Itoa(pageSize)
	return e.Orm.Raw(sql).Find(list).Error
}

// Get 获取EmailApprove对象
func (e *EmailApprove) Get(d *dto.EmailApproveGetReq, p *actions.DataPermission, outData *dto.EmailApproveGetResp) error {
	var data models.EmailApprove
	var model = &models.EmailApprove{}

	err := e.Orm.Model(&data).
		Scopes(
			EmailApprovePermission(d.CompanyId),
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetEmailApprove error:%s \r\n", err)
		return err
	}

	emailApproveDetails, _ := models.GetEmailApproveListByApproveId(e.Orm, model.Id)
	for _, item := range *emailApproveDetails {
		emailApproveDetailGetResp := dto.EmailApproveDetailGetResp{}
		if err := copier.Copy(&emailApproveDetailGetResp, &item); err != nil {
			return err
		}
		priorityEmailApproveDetails, _ := models.GetEmailApproveListByApproveIdAndPriority(e.Orm, model.Id, item.Priority)
		for _, priorityItem := range *priorityEmailApproveDetails {
			emailApproveDetailGetResp.ApproveUsers = append(emailApproveDetailGetResp.ApproveUsers, dto.EmailApproveGetUserResp{
				UserId:   priorityItem.UserId,
				UserName: priorityItem.User.FormatUserName(),
			})
		}
		outData.EmailApproveDetails = append(outData.EmailApproveDetails, emailApproveDetailGetResp)
	}

	userApproves, _ := models.GetUserApproveListById(e.Orm, model.Id)
	for _, item := range *userApproves {
		outData.RecipientUsers = append(outData.RecipientUsers, dto.EmailApproveGetUserResp{
			UserId:   item.UserId,
			UserName: item.User.FormatUserName(),
		})
	}
	approveUser, _ := models.GetApproveUserByCompanyId(e.Orm, d.CompanyId)
	for _, item := range *approveUser {
		outData.AllApproveUsers = append(outData.AllApproveUsers, dto.EmailApproveGetUserResp{
			UserId:   item.Id,
			UserName: item.FormatUserName(),
		})
	}

	recipientUser, _ := models.GetRecipientUserByCompanyId(e.Orm, d.CompanyId)
	for _, item := range *recipientUser {
		outData.AllRecipientUsers = append(outData.AllRecipientUsers, dto.EmailApproveGetUserResp{
			UserId:   item.Id,
			UserName: item.FormatUserName(),
		})
	}

	return nil
}

func (e *EmailApprove) GetAllApproveAndRecipientUsers(companyId int, p *actions.DataPermission, outData *dto.EmailApproveGetApproveAndRecipient) error {
	approveUser, _ := models.GetApproveUserByCompanyId(e.Orm, companyId)
	approveUserList := []dto.EmailApproveGetUserResp{}
	for _, item := range *approveUser {
		approveUserList = append(approveUserList, dto.EmailApproveGetUserResp{
			UserId:   item.Id,
			UserName: item.FormatUserName(),
		})
	}
	recipientUser, _ := models.GetRecipientUserByCompanyId(e.Orm, companyId)
	recipientUserList := []dto.EmailApproveGetUserResp{}
	for _, item := range *recipientUser {
		recipientUserList = append(recipientUserList, dto.EmailApproveGetUserResp{
			UserId:   item.Id,
			UserName: item.FormatUserName(),
		})
	}
	outData.AllApproveUsers = approveUserList
	outData.AllRecipientUsers = recipientUserList
	return nil
}

// Insert 创建EmailApprove对象
func (e *EmailApprove) Insert(c *dto.EmailApproveInsertReq) error {
	var err error
	var data models.EmailApprove
	if err := c.Generate(&data); err != nil {
		return err
	}
	err = e.Orm.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&data).Error
	})
	if err != nil {
		e.Log.Errorf("EmailApproveService Insert error:%s \r\n", err)
		return err
	}
	_ = modelsCommon.AddOperateLog(e.Orm, data.Id, "", c, &data, models.EmailApproveLogModelName, models.EmailApproveLogModelInsert, 1, c.CreateBy, c.CreateByName)
	return nil
}

// Update 修改EmailApprove对象
func (e *EmailApprove) Update(c *dto.EmailApproveUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.EmailApprove{}
	e.Orm.Scopes(
		EmailApprovePermission(c.CompanyId),
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())

	if err := e.CheckApproveModelData(&data); err != nil {
		return err
	}

	oldData := data
	if err := c.Generate(&data); err != nil {
		return err
	}

	err = e.Orm.Transaction(func(tx *gorm.DB) error {
		// 去除原始的detail信息
		emailApproveDetail := &models.EmailApproveDetail{}
		if err := emailApproveDetail.DeleteByApproveId(tx, data.Id); err != nil {
			return err
		}
		userApprove := &models.UserApprove{}
		if err := userApprove.DeleteByApproveId(tx, data.Id); err != nil {
			return err
		}
		return tx.Save(&data).Error
	})
	if err != nil {
		e.Log.Errorf("EmailApproveService Save error:%s \r\n", err)
		return err
	}
	_ = modelsCommon.AddOperateLog(e.Orm, data.Id, &oldData, c, &data, models.EmailApproveLogModelName, models.EmailApproveLogModelUpdate, 1, c.UpdateBy, c.UpdateByName)
	return nil
}

func (e *EmailApprove) CheckApproveModelData(data *models.EmailApprove) error {
	if data.Id == 0 {
		return errors.New("审批流不存在或没有权限")
	}
	if data.ApproveStatus != models.EmailApproveStatus1 {
		return errors.New("审批流状态不正确")
	}
	return nil
}

// Remove 删除EmailApprove
func (e *EmailApprove) Remove(d *dto.EmailApproveDeleteReq, p *actions.DataPermission) error {
	var data models.EmailApprove

	e.Orm.Model(&data).
		Scopes(
			EmailApprovePermission(d.CompanyId),
			actions.Permission(data.TableName(), p),
		).First(&data, d.GetId())

	if err := e.CheckApproveModelData(&data); err != nil {
		return err
	}
	oldData := data

	data.ApproveStatus = models.EmailApproveStatus0
	data.UpdateBy = d.UpdateBy
	data.UpdateByName = d.UpdateByName

	db := e.Orm.Save(&data)

	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveEmailApprove error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	_ = modelsCommon.AddOperateLog(e.Orm, data.Id, &oldData, d, &data, models.EmailApproveLogModelName, models.EmailApproveLogModelDelete, 1, d.UpdateBy, d.UpdateByName)
	return nil
}

func (e *EmailApprove) Workflow(userId int, p *actions.DataPermission, outData *[]dto.EmailApproveWorkflow) error {
	userApprove := models.UserApprove{}
	approveDetailsMap := map[int][]string{}
	emailApproveDetails, _ := userApprove.GetApproveListByUserId(e.Orm, userId)
	for _, item := range *emailApproveDetails {
		ApproveUserNameSlice := []string{}
		ApproveUserNameStr := ""
		priorityDetails, _ := models.GetEmailApproveListByApproveIdAndPriority(e.Orm, item.ApproveId, item.Priority)
		for _, priorityItem := range *priorityDetails {
			ApproveUserNameSlice = append(ApproveUserNameSlice, priorityItem.User.FormatUserName())
		}
		if item.ApproveRankType == models.EmailApproveRankType2 {
			ApproveUserNameStr = strings.Join(ApproveUserNameSlice, "、")
		} else if item.ApproveRankType == models.EmailApproveRankType3 {
			ApproveUserNameStr = strings.Join(ApproveUserNameSlice, "/")
		} else {
			ApproveUserNameStr = strings.Join(ApproveUserNameSlice, ",")
		}
		approveDetailsMap[item.ApproveId] = append(approveDetailsMap[item.ApproveId], ApproveUserNameStr)
	}
	for approveId, nameSlice := range approveDetailsMap {
		*outData = append(*outData, dto.EmailApproveWorkflow{
			Id:      approveId,
			Process: strconv.Itoa(approveId) + "-" + strings.Join(nameSlice, " > "),
		})
	}
	return nil
}
