package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"github.com/jinzhu/copier"
	"github.com/mitchellh/mapstructure"
	"go-admin/common/excel"
	"go-admin/config"
	"strings"

	"go-admin/app/uc/models"
	service "go-admin/app/uc/service/admin"
	"go-admin/app/uc/service/admin/dto"
	"go-admin/common/actions"
)

type UserInfo struct {
	api.Api
}

// GetPage 获取用户中心列表
// @Summary 获取用户中心列表
// @Description 获取用户中心列表
// @Tags 用户中心
// @Param userEmail query string false "用户邮箱"
// @Param loginName query string false "用户登录名"
// @Param userPhone query string false "手机号码"
// @Param userName query string false "用户名称"
// @Param userStatus query int false "用户状态（1可用，0不可用）"
// @Param companyId query int false "公司ID"
// @Param companyDepartmentId query int false "用户所属部门"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.UserInfo}} "{"code": 200, "data": [...]}"
// @Router /api/v1/uc/admin/user-info [get]
// @Security Bearer
func (e UserInfo) GetPage(c *gin.Context) {
	req := dto.UserInfoGetPageReq{}
	s := service.UserInfo{}
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
	list := make([]*dto.UserInfoGetListPageRes, 0)
	var count int64

	err = s.GetListPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取用户中心失败，\r\n失败信息 %s", err.Error()))
		return
	}

	//copier.Copy()
	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// GetSelectList 获取用户SelectList
// @Summary 获取用户SelectList
// @Description 获取用户SelectList
// @Tags 用户中心
// @Param userEmail query string false "用户邮箱"
// @Param loginName query string false "用户登录名"
// @Param userPhone query string false "手机号码"
// @Param userName query string false "用户名称"
// @Param userStatus query int false "用户状态（1可用，0不可用）"
// @Param companyId query int false "公司ID"
// @Param companyDepartmentId query int false "用户所属部门"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.UserInfo}} "{"code": 200, "data": [...]}"
// @Router /api/v1/uc/admin/user-info/select-list [get]
// @Security Bearer
func (e UserInfo) GetSelectList(c *gin.Context) {
	req := dto.UserInfoGetPageReq{}
	var res []*dto.UserInfoGetSelectListRes
	s := service.UserInfo{}
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
	list := make([]models.UserInfo, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取用户中心失败，\r\n失败信息 %s", err.Error()))
		return
	}

	_ = copier.Copy(&res, list)
	//copier.Copy()
	e.PageOK(res, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取用户中心
// @Summary 获取用户中心
// @Description 获取用户中心
// @Tags 用户中心
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.UserInfo} "{"code": 200, "data": [...]}"
// @Router /api/v1/uc/admin/user-info/{id} [get]
// @Security Bearer
func (e UserInfo) Get(c *gin.Context) {
	req := dto.UserInfoGetReq{}
	s := service.UserInfo{}
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
	var object dto.UserInfoGetRes

	p := actions.GetPermissionFromContext(c)
	err = s.GetInfo(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取用户中心失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建用户中心
// @Summary 创建用户中心
// @Description 创建用户中心
// @Tags 用户中心
// @Accept application/json
// @Product application/json
// @Param data body dto.UserInfoInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/uc/admin/user-info [post]
// @Security Bearer
func (e UserInfo) Insert(c *gin.Context) {
	req := dto.UserInfoInsertReq{}
	s := service.UserInfo{}
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
	errs := s.Insert(&req, false)
	if len(errs) > 0 {
		e.Error(500, err, fmt.Sprintf("创建用户中心失败，\r\n失败信息 %s", errs[0].Error()))
		return
	}

	e.OK(req.GetId(), "创建成功2")
}

// Update 修改用户中心
// @Summary 修改用户中心
// @Description 修改用户中心
// @Tags 用户中心
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.UserInfoUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/uc/admin/user-info/{id} [put]
// @Security Bearer
func (e UserInfo) Update(c *gin.Context) {
	req := dto.UserInfoUpdateReq{}
	s := service.UserInfo{}
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
	// 设置修改
	req.SetUpdateBy(user.GetUserId(c))
	req.SetUpdateByName(user.GetUserName(c))
	p := actions.GetPermissionFromContext(c)

	errs := s.Update(&req, p, false)
	if len(errs) > 0 {
		e.Error(500, err, fmt.Sprintf("修改用户信息失败，\r\n失败信息 %s", errs[0].Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Import 导入
// @Summary 导入
// @Description 导入
// @Tags 用户中心
// @Accept application/json
// @Product application/json
// @Param data body dto.ImportReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/uc/admin/user-info/import [post]
// @Security Bearer
func (e UserInfo) Import(c *gin.Context) {
	s := service.UserInfo{}
	roleInfoService := service.RoleInfo{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		MakeService(&roleInfoService.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	importFile, err := c.FormFile("file")
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	excelApp := excel.NewExcel()

	templateFilePath := config.ExtConfig.ApiHost + "/static/exceltpl/customer_import.xlsx"
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
	roleList, _ := roleInfoService.GetStatusOkList()

	for i := 0; i < len(list); i++ {
		data := list[i]
		req := dto.UserInfoImportData{}
		_ = mapstructure.Decode(data, &req)
		errs := s.Import(&req, roleList)
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
		titleList2, data2 := excelApp.MergeErrMsgColumn(titleList, list, errMsgList)
		_ = excelApp.ExportExcelByMap(c, titleList2, data2, "用户导入部分失败", "sheet1")
	} else {
		//没有校验错误则进行数据操作并接口返回
		e.OK(nil, "导入成功")
	}

	return
}

// ProxyLogin 代登录
// @Summary 代登录
// @Description 代登录
// @Tags 用户中心
// @Accept application/json
// @Product application/json
// @Param UserID query string false "用户Id"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/uc/admin/user-info/proxy-login [get]
// @Security Bearer
func (e UserInfo) ProxyLogin(c *gin.Context) {
	req := dto.UserInfoProxyLoginReq{}
	s := service.UserInfo{}
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

	res, err := s.ProxyLogin(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("代登录失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(res, "准备登录...")
}

// UpdatePassword 修改密码
// @Summary 修改密码
// @Description 修改密码
// @Tags 用户中心
// @Accept application/json
// @Product application/json
// @Param data body dto.UserInfoUpdatePassword true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/uc/admin/user-info/update-password [put]
// @Security Bearer
func (e UserInfo) UpdatePassword(c *gin.Context) {
	req := dto.UserInfoUpdatePassword{}
	s := service.UserInfo{}
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

	err = s.UpdatePassword(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改用户中心失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.Id, "修改成功")
}

// UpdateUserPhone 修改手机号
// @Summary 修改手机号
// @Description 修改手机号
// @Tags 用户中心
// @Accept application/json
// @Product application/json
// @Param data body dto.UserInfoUpdateUserPhone true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/uc/admin/user-info/update-userphone [put]
// @Security Bearer
func (e UserInfo) UpdateUserPhone(c *gin.Context) {
	req := dto.UserInfoUpdateUserPhone{}
	s := service.UserInfo{}
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

	err = s.UpdateUserPhone(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改用户中心失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.Id, "修改成功")
}

//// Delete 删除用户中心
//// @Summary 删除用户中心
//// @Description 删除用户中心
//// @Tags 用户中心
//// @Param data body dto.UserInfoDeleteReq true "body"
//// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
//// @Router /api/v1/uc/admin/user-info [delete]
//// @Security Bearer
//func (e UserInfo) Delete(c *gin.Context) {
//    s := service.UserInfo{}
//    req := dto.UserInfoDeleteReq{}
//    err := e.MakeContext(c).
//        MakeOrm().
//        Bind(&req).
//        MakeService(&s.Service).
//        Errors
//    if err != nil {
//        e.Logger.Error(err)
//        e.Error(500, err, err.Error())
//        return
//    }
//
//	// req.SetUpdateBy(user.GetUserId(c))
//	p := actions.GetPermissionFromContext(c)
//
//	err = s.Remove(&req, p)
//	if err != nil {
//		e.Error(500, err, fmt.Sprintf("删除用户中心失败，\r\n失败信息 %s", err.Error()))
//        return
//	}
//	e.OK( req.GetId(), "删除成功")
//}
