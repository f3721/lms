package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"github.com/jinzhu/copier"

	"go-admin/app/uc/models"
	service "go-admin/app/uc/service/admin"
	"go-admin/app/uc/service/admin/dto"
	"go-admin/common/actions"
)

type CompanyInfo struct {
	api.Api
}

// GetPage 获取公司信息列表
// @Summary 获取公司信息列表
// @Description 获取公司信息列表
// @Tags 公司信息
// @Param companyStatus query int false "公司状态（1可用 0不可用）"
// @Param companyName query string false "公司名称"
// @Param companyNature query int false "公司性质 （2终端，3分销）"
// @Param companyType query int false "公司类型（1:KA 2:SME 3:DS 4:整包 5:央企 6:其他）"
// @Param parentId query int false "父级节点"
// @Param companyLevels query int false "新建公司等级1-5"
// @Param pid query int false "最高级公司"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=dto.CompanyInfoGetPageRes} "{"code": 200, "data": [...]}"
// @Router /api/v1/uc/admin/company-info [get]
// @Security Bearer
func (e CompanyInfo) GetPage(c *gin.Context) {
	req := dto.CompanyInfoGetPageReq{}
	res := dto.CompanyInfoGetPageRes{}
	s := service.CompanyInfo{}
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
	subsidiaryList := make([]models.CompanyInfo, 0)
	var count int64

	topInfo, err := s.GetTopInfo()
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	listData := &dto.CompanyInfoGetPageData{}
	_ = copier.Copy(listData, topInfo)
	req.ParentId = topInfo.Id
	req.PageSize = 1000
	err = s.GetPage(&req, p, &subsidiaryList, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取公司信息失败，\r\n失败信息 %s", err.Error()))
		return
	}
	for _, subCompanyInfo := range subsidiaryList {
		subCompanyInfoData := &dto.CompanyInfoGetPageData{}
		_ = copier.Copy(subCompanyInfoData, subCompanyInfo)
		listData.Children = append(listData.Children, subCompanyInfoData)
	}
	res.List = append(res.List, listData)
	e.OK(res, "查询成功")
}

// GetSelectList 查询公司列表（下拉）
// @Summary 查询公司列表（下拉）
// @Description 查询公司列表（下拉）
// @Tags 公司信息
// @Param companyStatus query int false "公司状态（1可用 0不可用）"
// @Param companyName query string false "公司名称"
// @Param companyNature query int false "公司性质 （2终端，3分销）"
// @Param companyType query int false "公司类型（1:KA 2:SME 3:DS 4:整包 5:央企 6:其他）"
// @Param parentId query int false "父级节点"
// @Param companyLevels query int false "新建公司等级1-5"
// @Param pid query int false "最高级公司"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]dto.CompanyInfoGetSelectPageData}} "{"code": 200, "data": [...]}"
// @Router /api/v1/uc/admin/company-info/select-list [get]
// @Security Bearer
func (e CompanyInfo) GetSelectList(c *gin.Context) {
	req := dto.CompanyInfoGetPageReq{}
	s := service.CompanyInfo{}
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
	list := make([]models.CompanyInfo, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取公司信息失败，\r\n失败信息 %s", err.Error()))
		return
	}

	var resList []dto.CompanyInfoGetSelectPageData

	_ = copier.Copy(&resList, list)

	e.PageOK(resList, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// GetInnerSelectList 查询公司列表（下拉）Inner
// @Summary 查询公司列表（下拉）Inner
// @Description 查询公司列表（下拉）Inner
// @Tags 公司信息
// @Param companyStatus query int false "公司状态（1可用 0不可用）"
// @Param companyName query string false "公司名称"
// @Param companyNature query int false "公司性质 （2终端，3分销）"
// @Param companyType query int false "公司类型（1:KA 2:SME 3:DS 4:整包 5:央企 6:其他）"
// @Param parentId query int false "父级节点"
// @Param companyLevels query int false "新建公司等级1-5"
// @Param pid query int false "最高级公司"
// @Param pageIndex query int false "页码"
// @Param pageSize query int false "页条数 支持-1不分页"
// @Param queryCompanyIds query string false "公司id多查 id逗号分隔"
// @Param queryCompanyNames query string false "公司name多查 name逗号分隔"
// @Success 200 {object} response.Response{data=response.Page{list=[]dto.CompanyInfoGetSelectPageData}} "{"code": 200, "data": [...]}"
// @Router /inner/uc/admin/company-info/select-list [get]
// @Security Bearer
func (e CompanyInfo) GetInnerSelectList(c *gin.Context) {
	req := dto.CompanyInfoGetPageReq{}
	s := service.CompanyInfo{}
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
	list := make([]models.CompanyInfo, 0)
	var count int64

	req.IgnoreUserPermission = true
	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取公司信息失败，\r\n失败信息 %s", err.Error()))
		return
	}

	var resList []dto.CompanyInfoGetSelectPageData

	_ = copier.Copy(&resList, list)

	e.PageOK(resList, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// IsAvailable 查询公司是否可用
// @Summary 查询公司是否可用
// @Description 查询公司是否可用
// @Tags 公司信息
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=dto.CompanyInfoIsAvailableRes} "{"code": 200, "data": [...]}"
// @Router /inner/uc/admin/company-info/is-available/{id} [get]
// @Security Bearer
func (e CompanyInfo) IsAvailable(c *gin.Context) {
	req := dto.CompanyInfoGetIdReq{}
	s := service.CompanyInfo{}
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

	IsAvailable, err := s.CompanyIsAvailable(req.Id)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取公司信息失败，\r\n失败信息 %s", err.Error()))
		return
	}

	res := dto.CompanyInfoIsAvailableRes{IsAvailable: IsAvailable}
	e.OK(res, "查询成功")
}

// Get 获取公司信息
// @Summary 获取公司信息
// @Description 获取公司信息
// @Tags 公司信息
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=dto.CompanyInfoGetRes} "{"code": 200, "data": [...]}"
// @Router /api/v1/uc/admin/company-info/{id} [get]
// @Security Bearer
func (e CompanyInfo) Get(c *gin.Context) {
	req := dto.CompanyInfoGetReq{}
	s := service.CompanyInfo{}
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
	object, err := s.GetInfo(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取公司信息失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// GetInner 获取公司信息-Inner
// @Summary 获取公司信息-Inner
// @Description 获取公司信息-Inner
// @Tags 公司信息
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=dto.CompanyInfoGetRes} "{"code": 200, "data": [...]}"
// @Router /inner/uc/company-info/{id} [get]
// @Security Bearer
func (e CompanyInfo) GetInner(c *gin.Context) {
	req := dto.CompanyInfoGetReq{}
	s := service.CompanyInfo{}
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

	object := models.CompanyInfo{}
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取公司信息失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// GetByName 获取公司信息(根据公司名)
// @Summary 获取公司信息(根据公司名)
// @Description 获取公司信息(根据公司名)
// @Tags 公司信息
// @Param companyName query string false "公司名称"
// @Success 200 {object} response.Response{data=models.CompanyInfo} "{"code": 200, "data": [...]}"
// @Router /api/v1/uc/admin/company-info/get-by-name [get]
// @Security Bearer
func (e CompanyInfo) GetByName(c *gin.Context) {
	req := dto.CompanyInfoGetByNameReq{}
	s := service.CompanyInfo{}
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

	companyInfo, err := s.CompanyByName(req.CompanyName)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取公司信息失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(companyInfo, "查询成功")
}

// Insert 创建公司信息
// @Summary 创建公司信息
// @Description 创建公司信息
// @Tags 公司信息
// @Accept application/json
// @Product application/json
// @Param data body dto.CompanyInfoInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/uc/admin/company-info [post]
// @Security Bearer
func (e CompanyInfo) Insert(c *gin.Context) {
	req := dto.CompanyInfoInsertReq{}
	s := service.CompanyInfo{}
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
	err = s.Insert(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建公司信息失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改公司信息
// @Summary 修改公司信息
// @Description 修改公司信息
// @Tags 公司信息
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.CompanyInfoUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/uc/admin/company-info/{id} [put]
// @Security Bearer
func (e CompanyInfo) Update(c *gin.Context) {
	req := dto.CompanyInfoUpdateReq{}
	s := service.CompanyInfo{}
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
		e.Error(500, err, fmt.Sprintf("修改公司信息失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除公司信息
// @Summary 删除公司信息
// @Description 删除公司信息
// @Tags 公司信息
// @Param data body dto.CompanyInfoDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/uc/admin/company-info [delete]
// @Security Bearer
func (e CompanyInfo) Delete(c *gin.Context) {
	s := service.CompanyInfo{}
	req := dto.CompanyInfoDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除公司信息失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}

// Parameters 获取公司参数
// @Summary 获取公司参数
// @Description 获取公司参数
// @Tags 公司信息
// @Success 200 {object} response.Response{data=interface{}} "{"code": 200, "data": [...]}"
// @Router /api/v1/uc/admin/company-info/parameters [get]
// @Security Bearer
func (e CompanyInfo) Parameters(c *gin.Context) {
	_ = e.MakeContext(c)
	e.OK(models.CompanyTypeData, "查询成功")
}
