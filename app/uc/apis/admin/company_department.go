package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"github.com/jinzhu/copier"
	"github.com/mitchellh/mapstructure"
	"go-admin/app/uc/models"
	service "go-admin/app/uc/service/admin"
	"go-admin/app/uc/service/admin/dto"
	"go-admin/common/actions"
	commonApis "go-admin/common/apis"
	"go-admin/common/excel"
	"go-admin/config"
	"strconv"
	"strings"
)

type CompanyDepartment struct {
	api.Api
	Api2 commonApis.Api
}

// GetPage 获取公司部门列表
// @Summary 获取公司部门列表
// @Description 获取公司部门列表
// @Tags 部门管理
// @Param name query string false "部门名称"
// @Param level query int false "部门层级"
// @Param fId query int false "父级部门id"
// @Param topId query int false "一级部门ID"
// @Param companyId query int false "公司id"
// @Param personalBudget query string false "人均预算 null不设置无限预算 >=0 有预算限制 "
// @Param departmentBudget query string false "部门预算 null不设置无限预算 >=0 有预算限制 "
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]dto.CompanyDepartmentGetPageListData}} "{"code": 200, "data": [...]}"
// @Router /api/v1/uc/admin/company-department [get]
// @Security Bearer
func (e CompanyDepartment) GetPage(c *gin.Context) {
	req := dto.CompanyDepartmentGetPageReq{}
	s := service.CompanyDepartment{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	p := actions.GetPermissionFromContext(c)
	list := make([]dto.CompanyDepartmentGetPageListData, 0)
	var count int64
	// 部门列表
	err = s.GetListPage(&req, p, &list, &count)

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// GetSelectList 获取部门列表根据等级
// @Summary 获取公司部门列表
// @Description 获取公司部门列表
// @Tags 部门管理
// @Param companyId query int true "公司id"
// // @Param topId query int false "一级部门ID"
// @Param name query string false "部门名称"
// @Param level query int false "部门层级"
// @Param fId query int false "父级部门id"
// @Param pageSize query int false "页条数" 默认1000
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]dto.CompanyDepartmentGetPageListData}} "{"code": 200, "data": [...]}"
// @Router /api/v1/uc/admin/company-department/select-list [get]
// @Security Bearer
func (e CompanyDepartment) GetSelectList(c *gin.Context) {
	req := dto.CompanyDepartmentGetPageReq{}
	s := service.CompanyDepartment{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	p := actions.GetPermissionFromContext(c)
	list := make([]models.CompanyDepartment, 0)
	var count int64

	if req.PageSize == 0 {
		req.PageSize = 1000
	}
	// 部门列表
	err = s.GetPage(&req, p, &list, &count)

	var resList []*dto.CompanyDepartmentGetPageData
	for _, department := range list {
		resData := &dto.CompanyDepartmentGetPageData{}
		_ = copier.Copy(resData, department)
		if !department.DepartmentBudget.Valid {
			resData.DepartmentBudget = nil
		}
		if !department.PersonalBudget.Valid {
			resData.PersonalBudget = nil
		}

		resList = append(resList, resData)
	}
	e.PageOK(resList, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取部门信息
// @Summary 获取部门信息
// @Description 获取部门信息
// @Tags 部门管理
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.CompanyDepartment} "{"code": 200, "data": [...]}"
// @Router /api/v1/uc/admin/company-department/{id} [get]
// @Security Bearer
func (e CompanyDepartment) Get(c *gin.Context) {
	req := dto.CompanyDepartmentGetReq{}
	s := service.CompanyDepartment{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	var object dto.CompanyDepartmentGetRes

	p := actions.GetPermissionFromContext(c)
	err = s.GetInfo(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取公司部门失败，\r\\\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建部门
// @Summary 创建部门
// @Description 创建部门
// @Tags 部门管理
// @Accept application/json
// @Product application/json
// @Param data body dto.CompanyDepartmentInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/uc/admin/company-department [post]
// @Security Bearer
func (e CompanyDepartment) Insert(c *gin.Context) {
	req := dto.CompanyDepartmentInsertReq{}
	s := service.CompanyDepartment{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	// 设置创建人
	req.SetCreateBy(user.GetUserId(c))
	req.SetCreateByName(user.GetUserName(c))
	id, err := s.Insert(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建公司部门失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(id, "创建成功")
}

// Update 修改公司部门
// @Summary 修改公司部门
// @Description 修改公司部门
// @Tags 部门管理
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.CompanyDepartmentUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/uc/admin/company-department/{id} [put]
// @Security Bearer
func (e CompanyDepartment) Update(c *gin.Context) {
	req := dto.CompanyDepartmentUpdateReq{}
	s := service.CompanyDepartment{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	req.SetUpdateBy(user.GetUserId(c))
	req.SetUpdateByName(user.GetUserName(c))
	p := actions.GetPermissionFromContext(c)

	err = s.Update(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改公司部门失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// UpdateBudget 修改预算
// @Summary 修改预算
// @Description 修改预算
// @Tags 部门管理
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.CompanyDepartmentUpdateBudgetReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/uc/admin/company-department/update-budget/{id} [put]
// @Security Bearer
func (e CompanyDepartment) UpdateBudget(c *gin.Context) {
	req := dto.CompanyDepartmentUpdateBudgetReq{}
	s := service.CompanyDepartment{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	req.SetUpdateBy(user.GetUserId(c))
	p := actions.GetPermissionFromContext(c)

	err = s.UpdateBudget(&req, p, nil)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改公司部门失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除公司部门
// @Summary 删除公司部门
// @Description 删除公司部门
// @Tags 部门管理
// @Param data body dto.CompanyDepartmentDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/uc/admin/company-department [delete]
// @Security Bearer
func (e CompanyDepartment) Delete(c *gin.Context) {
	s := service.CompanyDepartment{}
	req := dto.CompanyDepartmentDeleteReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	// req.SetUpdateBy(user.GetUserId(c))
	p := actions.GetPermissionFromContext(c)

	err = s.Remove(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("删除公司部门失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}

// IsCaneDelete 是否可以删除删除部门
// @Summary 是否可以删除删除部门
// @Description 是否可以删除删除部门
// @Tags 部门管理
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=dto.CompanyDepartmentIsCanDeleteRes}	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/uc/admin/company-department/is-cane-delete/{id} [get]
// @Security Bearer
func (e CompanyDepartment) IsCaneDelete(c *gin.Context) {
	s := service.CompanyDepartment{}
	req := dto.CompanyDepartmentGetReq{}
	res := dto.CompanyDepartmentIsCanDeleteRes{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	// req.SetUpdateBy(user.GetUserId(c))
	p := actions.GetPermissionFromContext(c)

	err = s.RemoveVerify(&dto.CompanyDepartmentDeleteReq{Id: req.Id}, p)
	if err != nil {
		res.IsCanDelete = false
		e.OK(res, err.Error())
		return
	}
	res.IsCanDelete = true
	e.OK(res, "该操作将删除选中的部门，是否继续?")
}

// Import 导入
// @Summary 导入
// @Description 导入
// @Tags 部门管理
// @Accept application/json
// @Product application/json
// @Param data body dto.ImportReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/uc/admin/company-department/import [post]
// @Security Bearer
func (e CompanyDepartment) Import(c *gin.Context) {
	s := service.CompanyDepartment{}

	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	p := actions.GetPermissionFromContext(c)

	importFile, err := c.FormFile("file")
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	excelApp := excel.NewExcel()

	templateFilePath := config.ExtConfig.ApiHost + "/static/exceltpl/company_department_budget_import.xlsx"
	fieldsCorrect := excelApp.ValidImportFieldsCorrect(importFile, templateFilePath)
	if fieldsCorrect == false {
		fmt.Println("导入的excel字段与模板字段不一致，请重试")
		e.Error(500, nil, "导入的excel字段与模板字段不一致，请重试")
		return
	}

	// 读取上传的excel 并返回data
	_, list, titleList := excelApp.GetExcelData(importFile)
	// 校验data 获取 errMsgList
	errMsgList := map[int]string{}

	for i := 0; i < len(list); i++ {
		data := list[i]
		req := dto.CompanyDepartmentImportData{}
		_ = mapstructure.Decode(data, &req)
		fmt.Println("req")
		fmt.Println(req)
		errs := s.Import(&req, p)
		errText := ""
		if len(errs) > 0 {
			errStrings := make([]string, len(errs))
			for errI, err := range errs {
				errStrings[errI] = err.Error()
			}
			errText = strings.Join(errStrings, ",")
			errMsgList[i] = errText
		}
	}

	if len(errMsgList) > 0 {
		//e.OK(errMsgList, "部门导入部分失败")
		// 有校验错误的行时调用 导出保存excel
		titleList2, data2 := excelApp.MergeErrMsgColumn(titleList, list, errMsgList)
		_ = excelApp.ExportExcelByMap(c, titleList2, data2, "部门导入部分失败", "sheet1")
	} else {
		//没有校验错误则进行数据操作并接口返回
		e.OK(nil, "导入成功")
	}

	return
}

// Export 部门导出
// @Summary 部门导出
// @Description 部门导出
// @Tags 部门管理
// @Param ids query string false "ids"
// @Param name query string false "部门名称"
// @Param level query int false "部门层级"
// @Param fId query int false "父级部门id"
// @Param topId query int false "一级部门ID"
// @Param companyId query int false "公司id"
// @Param personalBudget query string false "人均预算 null不设置无限预算 >=0 有预算限制 "
// @Param departmentBudget query string false "部门预算 null不设置无限预算 >=0 有预算限制 "
// @Param tenant-id query string false "tenant-id"
// @Success 200 {object} response.Response{data=response.Page{list=[]dto.CompanyDepartmentGetPageListData}} "{"code": 200, "data": [...]}"
// @Router /api/v1/uc/admin/company-department/export [get]
// @Security Bearer
func (e CompanyDepartment) Export(c *gin.Context) {
	req := dto.CompanyDepartmentGetPageReq{}
	s := service.CompanyDepartment{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	p := actions.GetPermissionFromContext(c)
	var list []dto.CompanyDepartmentGetPageListData
	var count int64
	// 部门列表
	req.PageSize = -1
	err = s.GetListPage(&req, p, &list, &count)

	var exportList []interface{}

	// 遍历公司部门列表，并将每个元素转换为 interface{} 类型后添加到切片中
	for _, department := range list {
		exportData := dto.CompanyDepartmentExportData{}
		_ = copier.Copy(&exportData, department)

		if department.DepartmentBudget != nil {
			exportData.DepartmentBudget = strconv.FormatFloat(*department.DepartmentBudget, 'f', -1, 64)
		}
		if department.PersonalBudget != nil {
			exportData.PersonalBudget = strconv.FormatFloat(*department.PersonalBudget, 'f', -1, 64)
		}

		exportList = append(exportList, exportData)
	}

	title := []map[string]string{
		{"companyId": "公司ID"},
		{"companyName": "公司名称"},
		{"name": "部门名称"},
		{"fDepartmentBudget": "上级部门"},
		{"departmentBudget": "部门预算"},
		{"personalBudget": "人均预算"},
	}
	excelApp := excel.NewExcel()
	err = excelApp.ExportExcelByStruct(c, title, exportList, "company-department", "sheet1")
	if err != nil {
		e.Error(500, err, fmt.Sprintf("导出部门失败： %s", err.Error()))
		return
	}
}

//// GetPageBf GetPage报废的 获取公司部门列表
//// @Summary 获取公司部门列表
//// @Description 获取公司部门列表
//// @Tags 部门管理
//// @Param name query string false "部门名称"
//// @Param level query int false "部门层级"
//// @Param fId query int false "父级部门id"
//// @Param topId query int false "一级部门ID"
//// @Param companyId query int false "公司id"
//// @Param personalBudget query string false "人均预算 null不设置无限预算 >=0 有预算限制 "
//// @Param departmentBudget query string false "部门预算 null不设置无限预算 >=0 有预算限制 "
//// @Param pageSize query int false "页条数"
//// @Param pageIndex query int false "页码"
//// @Success 200 {object} response.Response{data=response.Page{list=[]dto.CompanyDepartmentGetPageCompanyData}} "{"code": 200, "data": [...]}"
//// @Router /api/v1/uc/admin/company-department [get]
//// @Security Bearer
//func (e CompanyDepartment) GetPageBf(c *gin.Context) {
//	req := dto.CompanyDepartmentGetPageReq{}
//	s := service.CompanyDepartment{}
//	companyInfoService := service.CompanyInfo{}
//	err := e.MakeContext(c).
//		MakeOrm().
//		Bind(&req).
//		MakeService(&s.Service).
//		MakeService(&companyInfoService.Service).
//		Errors
//	if err != nil {
//		e.Logger.Error(err)
//		e.Error(500, err, err.Error())
//		return
//	}
//
//	p := actions.GetPermissionFromContext(c)
//	list := make([]models.CompanyDepartment, 0)
//	var count int64
//	// 部门列表
//	err = s.GetPage(&req, p, &list, &count)
//	if err != nil {
//		e.Error(500, err, fmt.Sprintf("获取公司部门失败，\r\n失败信息 %s", err.Error()))
//		return
//	}
//
//	var queryCompanyIds []int
//	for _, department := range list {
//		queryCompanyIds = append(queryCompanyIds, department.CompanyId)
//	}
//
//	// 获取公司信息 部门列表是以公司分组展示的
//	companyList := make([]models.CompanyInfo, 0)
//	if len(queryCompanyIds) > 0 {
//		err = companyInfoService.GetPage(&dto.CompanyInfoGetPageReq{
//			QueryCompanyIds: queryCompanyIds,
//		}, p, &companyList, &count)
//	}
//
//	//departmentMap := make(map[int]*dto.GetPageDepartmentData)
//	departmentNameMap := make(map[int]string)
//	res := &dto.CompanyDepartmentGetPageRes{}
//
//	for _, info := range companyList {
//		data := &dto.CompanyDepartmentGetPageCompanyData{}
//		copier.Copy(&data, info)
//
//		for _, department := range list {
//			departmentResData := &dto.CompanyDepartmentGetPageDepartmentData{}
//			copier.Copy(departmentResData, department)
//			if department.FId > 0 {
//				departmentResData.Name = departmentNameMap[department.Id]
//			}
//
//			departmentNameMap[department.Id] = department.Name
//
//			if info.Id == department.CompanyId {
//				data.DepartmentList = append(data.DepartmentList, departmentResData)
//			}
//		}
//		res.List = append(res.List, data)
//	}
//
//	e.PageOK(res.List, len(res.List), req.GetPageIndex(), req.GetPageSize(), "查询成功")
//}
