package admin

import (
	"encoding/json"
	"errors"
	"go-admin/app/uc/models"
	"go-admin/app/uc/service/admin/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
	commonModels "go-admin/common/models"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"
)

type CompanyDepartment struct {
	service.Service
}

// GetPage 获取CompanyDepartment列表
func (e *CompanyDepartment) GetPage(c *dto.CompanyDepartmentGetPageReq, p *actions.DataPermission, list *[]models.CompanyDepartment, count *int64) error {
	var err error
	var data models.CompanyDepartment

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Order("company_id asc,level asc").
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("CompanyDepartmentService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取CompanyDepartment对象
func (e *CompanyDepartment) Get(d *dto.CompanyDepartmentGetReq, p *actions.DataPermission, model *models.CompanyDepartment) error {
	var data models.CompanyDepartment

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetCompanyDepartment error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// GetListPage 获取部门列表（带有上级部门信息和公司信息的列表）
func (e *CompanyDepartment) GetListPage(c *dto.CompanyDepartmentGetPageReq, p *actions.DataPermission, list *[]dto.CompanyDepartmentGetPageListData, count *int64) error {
	var err error
	var data models.CompanyDepartment
	err = e.Orm.Model(&data).Debug().
		Scopes(
			func(db *gorm.DB) *gorm.DB {
				if c.QueryFid > 0 {
					db.Where("company_department.id = ? or company_department.f_id = ?", c.QueryFid, c.QueryFid)
				}
				if c.Ids != "" {
					db.Where("company_department.id in ?", strings.Split(c.Ids, ","))
				}
				return db
			},
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSizeNegative(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
			actions.SysUserPermission(data.TableName(), p, 1),
		).
		Select("company_department.*, company_info.company_name as company_name, f_company_department.name as f_department_name").
		Joins("LEFT JOIN company_info ON company_department.company_id = company_info.id").
		Joins("LEFT JOIN company_department as f_company_department ON company_department.f_id = f_company_department.id").
		Order("company_department.company_id asc,company_department.level asc").
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("CompanyDepartmentService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// GetInfo 获取CompanyDepartment对象
func (e *CompanyDepartment) GetInfo(d *dto.CompanyDepartmentGetReq, p *actions.DataPermission, model *dto.CompanyDepartmentGetRes) error {
	var data models.CompanyDepartment

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetCompanyDepartment error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}

	if *model.FId == 0 {
		model.FId = nil
	}
	return nil
}

// GetById 根据id获取部门信息
func (e *CompanyDepartment) GetById(department *models.CompanyDepartment, id int) (err error) {
	err = e.Orm.Model(&department).First(department, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 处理记录不存在的情况
			err = errors.New("没有这个部门！")
		} else {
			// 处理其他错误情况
			return err
		}
	}
	return err
}

// GetByNameLevel 根据部门名称和等级获取部门信息
func (e *CompanyDepartment) GetByNameLevel(department *models.CompanyDepartment, companyId int, name string, fId int, level int) (err error) {
	db := e.Orm.Model(&department).Where("company_id = ?", companyId).Where("name = ?", name)
	if level > 0 {
		db = db.Where("level = ?", level)
	}
	if fId > 0 {
		db = db.Where("f_id = ?", fId)
	}
	err = db.First(&department).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			// 处理其他错误情况
			return err
		}
	}
	return nil
}

// GetAll 获取CompanyDepartment列表
func (e *CompanyDepartment) GetAll(c *dto.CompanyDepartmentGetPageReq, p *actions.DataPermission, list *[]models.CompanyDepartment) error {
	var err error
	var data models.CompanyDepartment

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Error
	if err != nil {
		e.Log.Errorf("CompanyDepartmentService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Insert 创建CompanyDepartment对象
func (e *CompanyDepartment) Insert(c *dto.CompanyDepartmentInsertReq) (id int, err error) {
	var data models.CompanyDepartment
	var saveData models.CompanyDepartment

	c.Generate(&saveData)

	err = e.SaveVerifyGenerate(&data, &saveData)
	if err != nil {
		return
	}

	err = e.Orm.Create(&saveData).Error
	if err != nil {
		e.Log.Errorf("CompanyDepartmentService Insert error:%s \r\n", err)
		return
	}
	tx := e.Orm.Begin()
	tx.Model(saveData).
		Where("level = ? AND top_id = ?", 1, 0).
		Updates(map[string]interface{}{"top_id": gorm.Expr("id")})

	_ = e.AddLog(saveData.Id, data, saveData, models.CompanyDepartmentOperationCreate, c.CreateBy, c.CreateByName)

	// 如果修改了预算
	budgetType, _ := e.GetSaveBudgetType(&data, &saveData)
	if budgetType > 0 {
		//* 生成预算方法
		//x @param int $dep_id
		//* @param int $type 1=>初次新增或预算调整，2=>部门预算和人均预算 模式相互转换
		var departmentBudgetM = models.DepartmentBudget{}
		err = departmentBudgetM.Generate(e.Orm, saveData.Id, 0, budgetType)
		if err != nil {
			tx.Rollback()
			return
		}
	}

	tx.Commit()
	return saveData.Id, nil
}

// SaveVerifyGenerate 创建CompanyDepartment对象
func (e *CompanyDepartment) SaveVerifyGenerate(data *models.CompanyDepartment, saveData *models.CompanyDepartment) (err error) {
	//上级部门的构造方法
	err = e.SaveFidGenerate(data, saveData)
	if err != nil {
		return err
	}
	if data.CompanyId != saveData.CompanyId {
		company := models.CompanyInfo{}
		e.Orm.First(&company, saveData.CompanyId)
		if company.Id == 0 {
			err = errors.New("公司不存在！")
			return
		}
	}
	if data.Name != saveData.Name {
		// 使用正则表达式匹配部门名称格式
		regex := regexp.MustCompile("^[a-zA-Z0-9_\u4e00-\u9fa5]+$")
		regexStatus := regex.MatchString(saveData.Name)
		if !regexStatus {
			// 正则表达式匹配出错，视为部门名称格式不正确
			err = errors.New("部门名称格式错误")
		}
		// 检查部门名称长度是否超过限制
		if utf8.RuneCountInString(saveData.Name) > 50 {
			err = errors.New("部门名称不能超出50字符")
		}
		var dataName models.CompanyDepartment
		err = e.GetByNameLevel(&dataName, saveData.CompanyId, saveData.Name, saveData.FId, saveData.Level)
		if err != nil {
			return err
		}
		e.Log.Info(dataName.Id)
		if dataName.Id > 0 {
			err = errors.New("部门名称不能重复！")
			return err
		}
	}

	// 如果修改了预算
	budgetType, _ := e.GetSaveBudgetType(data, saveData)
	if budgetType > 0 {
		err = e.SaveBudgetVerify(saveData)
		if err != nil {
			return err
		}
	}

	return nil
}

// Update 修改CompanyDepartment对象
func (e *CompanyDepartment) Update(c *dto.CompanyDepartmentUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.CompanyDepartment{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())

	saveData := data
	c.Generate(&saveData)

	err = e.SaveVerifyGenerate(&data, &saveData)
	if err != nil {
		return err
	}
	tx := e.Orm.Begin()

	db := tx.Save(&saveData)
	if err = db.Error; err != nil {
		tx.Rollback()
		e.Log.Errorf("CompanyDepartmentService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}

	_ = e.AddLog(saveData.Id, data, saveData, models.CompanyDepartmentOperationUpdate, c.UpdateBy, c.UpdateByName)

	tx.Commit()
	budgetType, _ := e.GetSaveBudgetType(&data, &saveData)
	if budgetType > 0 {
		//* 生成预算方法
		//x @param int $dep_id
		//* @param int $type 1=>初次新增或预算调整，2=>部门预算和人均预算 模式相互转换
		var departmentBudgetM = models.DepartmentBudget{}
		_ = departmentBudgetM.Generate(e.Orm, saveData.Id, 0, budgetType)
	}

	return nil
}

// SaveFidGenerate 涉及到修改上级部门的一些数据操作（修改Level， TopId等）
func (e *CompanyDepartment) SaveFidGenerate(data *models.CompanyDepartment, saveData *models.CompanyDepartment) (err error) {
	// 如果修改了上级部门
	if saveData.FId == 0 {
		saveData.Level = 1
		saveData.TopId = saveData.Id
	}

	if data.FId != saveData.FId {
		if saveData.FId > 0 {
			if saveData.FId == saveData.Id {
				return errors.New("上级部门不能是自己！")
			}
			if data.Id > 0 && data.FId == 0 {
				subDepartment := models.CompanyDepartment{}
				e.Orm.Model(&models.CompanyDepartment{}).Where("f_id = ?", saveData.Id).First(&subDepartment)
				if subDepartment.Id > 0 {
					return errors.New("该部门已有下级部门不能修改父级部门！")
				}
			}

			fDepartment := models.CompanyDepartment{}
			err = e.GetById(&fDepartment, saveData.FId)
			if err != nil {
				return err
			}

			if fDepartment.Level > 1 {
				return errors.New("上级部门不是一级部门不能添加！")
			}
			if fDepartment.TopId == 0 {
				return errors.New("上级部门数据异常！")
			}
			if fDepartment.CompanyId != saveData.CompanyId {
				return errors.New("上级部门不属于该公司！")
			}

			saveData.Level = fDepartment.Level + 1
			saveData.TopId = fDepartment.TopId
		}
	}
	return err
}

// SaveBudgetVerify 新增/修改预算前的校验
func (e *CompanyDepartment) SaveBudgetVerify(department *models.CompanyDepartment) (err error) {
	if department.DepartmentBudget.Valid && department.PersonalBudget.Valid {
		err = errors.New("部门预算和人均预算只能设置一个")
		return
	}
	if (department.DepartmentBudget.Valid && department.DepartmentBudget.Float64 < 0) ||
		(department.PersonalBudget.Valid && department.PersonalBudget.Float64 < 0) {
		err = errors.New("预算不能是负数")
		return
	}
	if (department.DepartmentBudget.Valid && department.DepartmentBudget.Float64 > 100000000) ||
		(department.PersonalBudget.Valid && department.PersonalBudget.Float64 > 100000000) {
		err = errors.New("预算不能>=100000000")
		return
	}
	return
}

// GetSaveBudgetType 获取预算调整类型 0 没调整 1 初次新增或预算调整，2 部门预算和人均预算模式相互转换
func (e *CompanyDepartment) GetSaveBudgetType(department *models.CompanyDepartment, saveDepartment *models.CompanyDepartment) (budgetType int, err error) {
	// 初次新增或预算调整budgetType
	// 0 没调整不掉哟joe预算重算方法
	// 1 初次新增或预算调整，2 部门预算和人均预算模式相互转换
	budgetType = 0

	// 第一次设置预算
	if department.Id == 0 && !department.PersonalBudget.Valid && department.DepartmentBudget.Valid {
		budgetType = 1
	}

	// 预算调整
	if (department.PersonalBudget.Float64 != saveDepartment.PersonalBudget.Float64) ||
		(department.DepartmentBudget.Float64 != saveDepartment.DepartmentBudget.Float64) {
		budgetType = 1
	}

	// 部门预算和人均预算模式相互转换
	if (department.DepartmentBudget.Valid && saveDepartment.PersonalBudget.Valid && !department.PersonalBudget.Valid) ||
		(department.PersonalBudget.Valid && saveDepartment.DepartmentBudget.Valid && !department.DepartmentBudget.Valid) {
		budgetType = 2
	}

	return
}

// UpdateBudget 修改预算方法
func (e *CompanyDepartment) UpdateBudget(c *dto.CompanyDepartmentUpdateBudgetReq, p *actions.DataPermission, db *gorm.DB) (err error) {
	if db != nil {
		db = e.Orm
	}
	var department = models.CompanyDepartment{}
	var saveDepartment = models.CompanyDepartment{}
	e.Orm.Scopes(
		actions.Permission(department.TableName(), p),
	).First(&department, c.GetId())
	c.Generate(&saveDepartment)
	err = e.SaveBudgetVerify(&saveDepartment)
	if err != nil {
		return err
	}

	// 初次新增或预算调整
	// 0 没调整不掉哟joe预算重算方法
	// 1 初次新增或预算调整，2 部门预算和人均预算模式相互转换
	budgetType, _ := e.GetSaveBudgetType(&department, &saveDepartment)

	tx := db.Begin()
	err = tx.Save(&saveDepartment).Error
	if err != nil {
		tx.Rollback()
		return
	}

	if budgetType > 0 {
		// * 生成预算方法
		// x @param int $dep_id
		// * @param int $type 1=>初次新增或预算调整，2=>部门预算和人均预算 模式相互转换
		var departmentBudgetM = models.DepartmentBudget{}
		err = departmentBudgetM.Generate(tx, department.Id, 0, budgetType)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return
}

// RemoveVerify 删除前的校验 校验是否可以删除
func (e *CompanyDepartment) RemoveVerify(d *dto.CompanyDepartmentDeleteReq, p *actions.DataPermission) error {
	// Check if there are second-level departments under the current department
	var count int64
	err := e.GetPage(&dto.CompanyDepartmentGetPageReq{
		Pagination: cDto.Pagination{
			PageIndex: 1,
			PageSize:  1,
		},
		TopId: d.GetId(),
	}, p, &[]models.CompanyDepartment{}, &count)

	if count > 1 {
		err = errors.New("部门下存在二级部门不能删除")
		return err
	}

	e.Orm.Model(&models.UserInfo{}).Where("company_department_id = ?", d.GetId()).Limit(1).Count(&count)
	if count > 0 {
		err = errors.New("部门下存在用户不能删除！")
		return err
	}

	return nil
}

// Remove 删除CompanyDepartment
func (e *CompanyDepartment) Remove(d *dto.CompanyDepartmentDeleteReq, p *actions.DataPermission) (err error) {
	var data models.CompanyDepartment
	err = e.RemoveVerify(d, p)
	if err != nil {
		return err
	}
	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err = db.Error; err != nil {
		e.Log.Errorf("Service RemoveCompanyDepartment error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}

// Insert 创建CompanyDepartment对象
func (e *CompanyDepartment) Import(req *dto.CompanyDepartmentImportData, p *actions.DataPermission) (errs []error) {
	var companyId int
	companyId, _ = strconv.Atoi(req.CompanyId)
	department := models.CompanyDepartment{}
	fId := 0
	level := 1
	if req.Name == "" {
		errs = append(errs, errors.New("部门名称不能为空"))
	} else {
		// 使用正则表达式匹配部门名称格式
		regex := regexp.MustCompile("^[a-zA-Z0-9_\u4e00-\u9fa5]+$")
		regexStatus := regex.MatchString(req.Name)
		if !regexStatus {
			// 正则表达式匹配出错，视为部门名称格式不正确
			errs = append(errs, errors.New("部门名称格式错误"))
		}
		// 检查部门名称长度是否超过限制
		if utf8.RuneCountInString(req.Name) > 50 {
			errs = append(errs, errors.New("部门名称不能超出50字符"))
		}
	}
	if companyId == 0 {
		errs = append(errs, errors.New("公司信息错误"))
	} else {
		company := models.CompanyInfo{}
		e.Orm.First(&company, companyId)
		if company.Id == 0 {
			errs = append(errs, errors.New("公司不存在"))
		} else if req.CompanyName != company.CompanyName {
			errs = append(errs, errors.New("公司信息错误"))
		}
	}

	if req.FDepartmentName != "" {
		// 校验父级部门
		_ = e.GetByNameLevel(&department, companyId, req.FDepartmentName, 0, 0)
		if department.Id == 0 {
			errs = append(errs, errors.New("父级部门不存在"))
			return
		}
		if department.Level != 1 {
			errs = append(errs, errors.New("父级上级部门不是一级部门！"))
			return
		}
		level = 2
		fId = department.Id
	}

	if errs != nil {
		return
	}

	// 根据部门名称和父级部门id获取部门信息
	department = models.CompanyDepartment{}
	_ = e.GetByNameLevel(&department, companyId, req.Name, fId, level)
	// 没有部门调用插入
	var err error

	if department.Id == 0 {
		var personalBudget *float64 = nil
		if req.PersonalBudget != "" {
			value, _ := strconv.ParseFloat(req.PersonalBudget, 64)
			personalBudget = &value
		}

		var departmentBudget *float64 = nil
		if req.DepartmentBudget != "" {
			value, _ := strconv.ParseFloat(req.DepartmentBudget, 64)
			departmentBudget = &value
		}

		_, err = e.Insert(&dto.CompanyDepartmentInsertReq{
			Name:             req.Name,
			FId:              fId,
			CompanyId:        companyId,
			PersonalBudget:   personalBudget,
			DepartmentBudget: departmentBudget,
			ControlBy: commonModels.ControlBy{
				CreateBy:     user.GetUserId(e.Orm.Statement.Context.(*gin.Context)),
				CreateByName: user.GetUserName(e.Orm.Statement.Context.(*gin.Context)),
			},
		})
	} else {
		var personalBudget *float64 = nil
		if req.PersonalBudget != "" {
			value, _ := strconv.ParseFloat(req.PersonalBudget, 64)
			personalBudget = &value
		}

		var departmentBudget *float64 = nil
		if req.DepartmentBudget != "" {
			value, _ := strconv.ParseFloat(req.DepartmentBudget, 64)
			departmentBudget = &value
		}
		err = e.Update(&dto.CompanyDepartmentUpdateReq{
			Id:               department.Id,
			FId:              department.FId,
			PersonalBudget:   personalBudget,
			DepartmentBudget: departmentBudget,
			ControlBy: commonModels.ControlBy{
				UpdateBy:     user.GetUserId(e.Orm.Statement.Context.(*gin.Context)),
				UpdateByName: user.GetUserName(e.Orm.Statement.Context.(*gin.Context)),
			},
		}, p)
	}
	if err != nil {
		return append(errs, err)
	}
	return
}

func (e *CompanyDepartment) AddLog(dataId int, oldData, saveData interface{}, modelType string, operatorId int, operatorName string) error {
	//记录操作日志
	oldDataStr := ""
	if saveData != nil {
		oldDataJson, _ := json.Marshal(&oldData)
		oldDataStr = string(oldDataJson)
	}
	dataStr, _ := json.Marshal(&saveData)
	opLog := commonModels.OperateLogs{
		DataId:       strconv.Itoa(dataId),
		ModelName:    models.CompanyDepartmentOperationModel,
		Type:         modelType,
		DoStatus:     "",
		Before:       oldDataStr,
		Data:         string(dataStr),
		After:        string(dataStr),
		OperatorId:   operatorId,
		OperatorName: operatorName,
	}
	_ = opLog.InsertItem(e.Orm)
	return nil
}
