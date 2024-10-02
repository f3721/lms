package admin

import (
	"encoding/json"
	"errors"
	"regexp"
	"strconv"

	"github.com/araddon/dateparse"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/uc/service/admin/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
	"go-admin/common/models"
)

type OperateLogs struct {
	service.Service
}

var RegTime = regexp.MustCompile(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z|\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}\+\d{2}:\d{2}`)

// GetPage 获取OperateLogs列表
func (e *OperateLogs) GetPage(c *dto.OperateLogsGetPageReq, p *actions.DataPermission, list *[]models.OperateLogs, count *int64) error {
	var err error
	var data models.OperateLogs

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Order("id desc").
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("OperateLogsService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取OperateLogs对象
func (e *OperateLogs) Get(d *dto.OperateLogsGetReq, p *actions.DataPermission, detailResp *models.OperateLogDetailResp) error {
	var data models.OperateLogs
	var model = &models.OperateLogs{}

	err := e.Orm.Model(&data).Debug().
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetOperateLogs error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	if err := json.Unmarshal([]byte(model.Diff), detailResp); err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}

	FormatTimeForLogDetail(detailResp)
	return nil
}

func FormatTimeForLogDetail(resp *models.OperateLogDetailResp) {
	for index, item := range *resp {
		t, err := dateparse.ParseAny(item.Value)
		if err == nil {
			(*resp)[index].Value = t.Format("2006-01-02 15:04:05")
		}
	}
}

func (e *OperateLogs) AddLog(dataId int, oldData, saveData interface{}, modelName, modelType string, operatorId int, operatorName string) error {
	//记录操作日志
	oldDataStr := ""
	if saveData != nil {
		oldDataJson, _ := json.Marshal(&oldData)
		oldDataStr = string(oldDataJson)
	}
	dataStr, _ := json.Marshal(&saveData)
	opLog := models.OperateLogs{
		DataId:       strconv.Itoa(dataId),
		ModelName:    modelName,
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
