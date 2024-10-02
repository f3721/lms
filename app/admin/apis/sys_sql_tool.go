package apis

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"go-admin/app/admin/service/dto"
	"go-admin/common/database"
	"strings"
)

type SysSqlTool struct {
	api.Api
}

// Transition sql生成 所有租户的sql
func (e SysSqlTool) Transition(c *gin.Context) {
	req := dto.SysSqlTooleTransitionReq{}
	err := e.MakeContext(c).
		Bind(&req).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	req.Sql = strings.TrimSpace(req.Sql)
	req.Sql = strings.Trim(req.Sql, "\n")
	req.Sql = strings.Trim(req.Sql, "\r")
	if req.Sql == "" {
		err = errors.New("sql不能为空")
		e.Error(500, err, err.Error())
		return
	}

	tenants := database.GetInitOkTenantList(database.GetLubanToken())
	if len(tenants) <= 0 {
		err = errors.New("暂无驿站租户")
		e.Error(500, err, err.Error())
		return
	}

	sqls := strings.Split(req.Sql, ";")
	modes := []string{"admin", "oc", "pc", "wc", "uc"}
	var resSql []string
	for _, sql := range sqls {
		sql = strings.TrimSpace(sql)
		sql = strings.Trim(sql, "\n")
		sql = strings.Trim(sql, "\r")
		if sql == "" {
			continue
		}
		tmpSql := sql + ";\r\n"
		for _, tenant := range tenants {
			for _, mode := range modes {
				if strings.Index(sql, "base_"+mode) != -1 {
					tmpSql = tmpSql + strings.Replace(sql, "base_"+mode, tenant.DatabaseName+"_"+mode, 1) + ";\r\n"
				}

			}
		}
		resSql = append(resSql, tmpSql)
	}
	res := strings.Join(resSql, "\r\n")

	e.OK(res, "转换成功")
}
