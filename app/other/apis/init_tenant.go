package apis

import (
	"errors"
	"fmt"
	"go-admin/common/global"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/other/service"
	"go-admin/app/other/service/dto"
)

type InitTenant struct {
	api.Api
}

// GetTenantId 获取租户id
func (e InitTenant) GetTenantId(c *gin.Context) {
	s := service.InitTenant{}
	err := e.MakeContext(c).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	tenantName := c.Query("tenantName")
	if tenantName == "" {
		err = errors.New("租户名称必填")
		e.Error(500, err, err.Error())
		return
	}
	//tenants := database.GetInitOkTenantList(database.GetLubanToken())
	tenants := global.GetTenants()
	if len(tenants) <= 0 {
		err = errors.New("暂无租户")
		e.Error(500, err, err.Error())
		return
	}
	tenantId := ""
	for tenantKey, tenant := range tenants {
		if tenant.Name == tenantName {
			tenantId = tenantKey
			break
		}
	}
	if tenantId != "" {
		e.OK(tenantId, "查询成功")
	} else {
		err = errors.New("租户名称不存在")
		e.Error(500, err, err.Error())
	}
}

// InitTenant 初始化租户 表及数据
func (e InitTenant) InitTenant(c *gin.Context) {
	s := service.InitTenant{}
	req := dto.InitTenantReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	err = s.InitTenant(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("初始化租户失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(nil, "初始化租户成功")
}

// 去除租户名的登录
func (e InitTenant) LoginNew(c *gin.Context) {
	s := service.InitTenant{}
	req := dto.Login{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	// 逻辑处理
	data, err := s.LoginNew(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("前置登录，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(data, "去租户名登录查询成功")
}

// 去除租户名的登录-小程序
func (e InitTenant) LoginNewMini(c *gin.Context) {
	s := service.InitTenant{}
	req := dto.LoginMini{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	// 逻辑处理
	data, err := s.LoginNewMini(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("前置登录，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(data, "去租户名登录查询成功")
}
