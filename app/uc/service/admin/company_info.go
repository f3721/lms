package admin

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	adminClient "go-admin/common/client/admin"
	ocClient "go-admin/common/client/oc"
	"go-admin/common/utils"
	"gorm.io/gorm"
	"regexp"
	"strconv"
	"strings"

	"go-admin/app/uc/models"
	"go-admin/app/uc/service/admin/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
	commonModels "go-admin/common/models"
)

type CompanyInfo struct {
	service.Service
}

// GetPage 获取CompanyInfo列表
func (e *CompanyInfo) GetPage(c *dto.CompanyInfoGetPageReq, p *actions.DataPermission, list *[]models.CompanyInfo, count *int64) error {
	var err error
	var data models.CompanyInfo

	err = e.Orm.Model(&data).
		Scopes(
			func(db *gorm.DB) *gorm.DB {
				if c.IgnoreUserPermission != true {
					//actions.SysUserPermission(data.TableName(), p, 1),
					db.Where("id in ?", utils.Split(p.AuthorityCompanyId))
				}
				if c.QueryCompanyNames != "" {
					db = db.Where("company_name IN (?)", strings.Split(c.QueryCompanyNames, ","))
				}
				if c.QueryCompanyIds != "" {
					db = db.Where("id IN (?)", strings.Split(c.QueryCompanyIds, ","))
				}
				return db
			},
			cDto.MakeCondition(c.GetNeedSearch()),
			func(db *gorm.DB) *gorm.DB {
				if c.PageSize != -1 {
					cDto.Paginate(c.GetPageSize(), c.GetPageIndex())
				}

				return db
			},

			actions.Permission(data.TableName(), p),
		).
		Order("pid asc, id desc").
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("CompanyInfoService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取CompanyInfo对象
func (e *CompanyInfo) Get(d *dto.CompanyInfoGetReq, p *actions.DataPermission, model *models.CompanyInfo) error {
	var data models.CompanyInfo

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetCompanyInfo error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Get 获取CompanyInfo对象
func (e *CompanyInfo) GetInfo(d *dto.CompanyInfoGetReq, p *actions.DataPermission) (data *dto.CompanyInfoGetRes, err error) {
	var companyInfo models.CompanyInfo
	err = e.Get(d, p, &companyInfo)
	if err != nil {
		return
	}

	data = &dto.CompanyInfoGetRes{companyInfo, ""}
	if companyInfo.ParentId > 0 {
		var parentCompanyInfo models.CompanyInfo
		err = e.Get(&dto.CompanyInfoGetReq{Id: companyInfo.ParentId}, p, &parentCompanyInfo)
		if err != nil {
			return
		}
		data.ParentName = parentCompanyInfo.CompanyName
	}

	return
}

// Insert 创建CompanyInfo对象
func (e *CompanyInfo) Insert(c *dto.CompanyInfoInsertReq) error {
	var err error
	var data models.CompanyInfo
	var saveData models.CompanyInfo
	c.Generate(&saveData)

	err = e.SaveSubVerifyGenerate(&data, &saveData)
	if err != nil {
		return err
	}

	err = e.Orm.Create(&saveData).Error
	if err != nil {
		e.Log.Errorf("CompanyInfoService Insert error:%s \r\n", err)
		return err
	}

	// 创建公司成功后给创建公司的用户 该公司权限
	err = e.addCompanyUserAuth(saveData.Id)
	if err != nil {
		return err
	}

	//记录操作日志
	DataStr, _ := json.Marshal(&saveData)
	opLog := commonModels.OperateLogs{
		DataId:       strconv.Itoa(saveData.Id),
		ModelName:    "companyInfo",
		Type:         "create",
		DoStatus:     "",
		Before:       "",
		Data:         string(DataStr),
		After:        string(DataStr),
		OperatorId:   c.CreateBy,
		OperatorName: c.CreateByName,
	}
	_ = opLog.InsertItem(e.Orm)
	return nil
}

func (e *CompanyInfo) SaveSubVerifyGenerate(data *models.CompanyInfo, saveData *models.CompanyInfo) (err error) {
	topData, err := e.GetTopInfo()
	if err != nil {
		return
	}

	if data.Id == 0 {
		saveData.CompanyLevels = 2
		saveData.ParentId = topData.Id
		saveData.Pid = topData.Id
	}

	//// 使用正则表达式检查是否包含中文括号或英文括号
	//match := regexp.MustCompile(`[\p{Han}()\[\]]`).MatchString(data.CompanyName)
	//if match {
	//	return errors.New("")
	//}

	if data.CompanyName != saveData.CompanyName {

		isExists, err := e.CompanyIsExistsByName(saveData.CompanyName)
		if err != nil {
			return err
		}
		if isExists {
			return errors.New("公司名称已存在！")
		}
		// 公司名括号强转
		saveData.CompanyName = strings.ReplaceAll(saveData.CompanyName, "(", "（")
		saveData.CompanyName = strings.ReplaceAll(saveData.CompanyName, ")", "）")
		saveData.CompanyName = regexp.MustCompile(`[\s　]+`).ReplaceAllString(saveData.CompanyName, "")
	}
	if data.Id > 1 && data.CompanyStatus == 1 && saveData.CompanyStatus == 0 {
		checkExistUnCompletedOrder := e.OcClientClCheckExistUnCompletedOrder(e.Orm, data.Id)
		if !checkExistUnCompletedOrder {
			return errors.New("该公司有未完结的订单，不可禁用！")
		}
	}

	return
}

// OcClientClCheckExistUnCompletedOrder 判断是否有未完结订单
func (e *CompanyInfo) OcClientClCheckExistUnCompletedOrder(tx *gorm.DB, companyId int) bool {
	result := ocClient.ApiByDbContext(tx).CheckExistUnCompletedOrder(strconv.Itoa(companyId))
	resultInfo := &struct {
		response.Response
		Data int
	}{}
	result.Scan(resultInfo)
	if resultInfo.Data == 400 {
		return false
	}
	return true
}

// addCompanyUserAuth 创建公司成功后给创建公司的用户 该公司权限
func (e *CompanyInfo) addCompanyUserAuth(companyId int) error {
	result := adminClient.Api(e.Orm.Statement.Context.(*gin.Context)).
		UpdatePermission(adminClient.ApiClientPermissionUpdateRequest{
			UserID:                    user.GetUserId(e.Orm.Statement.Context.(*gin.Context)),
			AuthorityCompanyID:        strconv.Itoa(companyId),
			AuthorityWarehouseID:      "",
			AuthorityWarehouseAllocID: "",
			AuthorityVendorID:         "",
			UpdateBy:                  user.GetUserId(e.Orm.Statement.Context.(*gin.Context)),
			UpdateByName:              user.GetUserName(e.Orm.Statement.Context.(*gin.Context)),
		})
	if !result.OK() {
		return result.Error()
	}
	return nil
}

func (e *CompanyInfo) GetTopInfo() (*models.CompanyInfo, error) {
	var err error
	var data models.CompanyInfo
	err = e.Orm.Model(&data).
		Where("company_levels = 1").
		First(&data).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetCompanyInfo error:%s \r\n", err)
		return nil, err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return nil, err
	}
	return &data, nil
}

// CompanyIsAvailable 公司是否可用
func (e *CompanyInfo) CompanyIsAvailable(id int) (bool, error) {
	var err error
	var data models.CompanyInfo
	err = e.Orm.Model(&data).
		First(&data, id).Error
	if err != nil {
		return false, err
	}
	if data.CompanyStatus == 1 {
		return true, nil
	}
	return false, nil
}

func (e *CompanyInfo) CompanyIsExists(id int) (bool, error) {
	var err error
	var data models.CompanyInfo
	err = e.Orm.Model(&data).
		First(&data, id).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return false, err
	}
	return true, nil
}

func (e *CompanyInfo) CompanyIsExistsByName(name string) (bool, error) {
	var err error
	var data models.CompanyInfo
	err = e.Orm.Model(&data).
		Where("company_name = ?", name).
		First(&data).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return false, err
	}
	return true, nil
}

func (e *CompanyInfo) CompanyByName(name string) (*models.CompanyInfo, error) {
	var err error
	var data models.CompanyInfo
	err = e.Orm.Model(&data).
		Where("company_name = ?", name).
		First(&data).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return &data, nil
}

// Update 修改CompanyInfo对象
func (e *CompanyInfo) Update(c *dto.CompanyInfoUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.CompanyInfo{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	saveData := data
	c.Generate(&saveData)

	err = e.SaveSubVerifyGenerate(&data, &saveData)
	if err != nil {
		return err
	}

	db := e.Orm.Save(&saveData)
	if err = db.Error; err != nil {
		e.Log.Errorf("CompanyInfoService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}

	//记录操作日志
	oldDataStr, _ := json.Marshal(&data)
	dataStr, _ := json.Marshal(&saveData)
	opLog := commonModels.OperateLogs{
		DataId:       strconv.Itoa(saveData.Id),
		ModelName:    "companyInfo",
		Type:         "update",
		DoStatus:     "",
		Before:       string(oldDataStr),
		Data:         string(dataStr),
		After:        string(dataStr),
		OperatorId:   c.UpdateBy,
		OperatorName: c.UpdateByName,
	}
	_ = opLog.InsertItem(e.Orm)
	return nil
}

// Remove 删除CompanyInfo
func (e *CompanyInfo) Remove(d *dto.CompanyInfoDeleteReq, p *actions.DataPermission) error {
	var data models.CompanyInfo

	err := e.Orm.First(&data, d.GetId()).Error
	if err != nil {
		return err
	}
	if data.CompanyLevels == 1 {
		return errors.New("主体公司信息无法删除")
	}
	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveCompanyInfo error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
