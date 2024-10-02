package admin

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-admin/common/utils"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/pc/models"
	"go-admin/app/pc/service/admin/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type GoodsLog struct {
	service.Service
}

// GetPage 获取GoodsLog列表
func (e *GoodsLog) GetPage(c *dto.GoodsLogGetPageReq, p *actions.DataPermission, list *[]models.GoodsLog, count *int64) error {
	var err error
	var data models.GoodsLog

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("GoodsLogService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取GoodsLog对象
func (e *GoodsLog) Get(d *dto.GoodsLogGetReq, p *actions.DataPermission, detailResp *utils.OperateLogDetailResp) error {
	var model models.GoodsLog
	err := e.Orm.
		Scopes(
			actions.Permission(model.TableName(), p),
		).
		First(&model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetGoodsLog error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	if err := json.Unmarshal([]byte(fmt.Sprintf(`%s`, model.DiffData)), detailResp); err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}
