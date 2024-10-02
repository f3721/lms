package models

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

type OperateLogDetailResp []OperateLogDetail

type OperateLogDetail struct {
	Key   string `json:"key"`
	Desc  string `json:"desc"`
	Value string `json:"value"`
}

type OperateLogs struct {
	Model

	DataId       string    `json:"dataId" gorm:"type:varchar(30);comment:数据id eg:WH0001"`
	ModelName    string    `json:"modelName" gorm:"type:varchar(50);comment:模型name"`
	Type         string    `json:"type" gorm:"type:varchar(20);comment:eg: create"`
	TypeName     string    `json:"typeName" gorm:"type:varchar(20);comment:eg: 创建实体仓"`
	DoStatus     string    `json:"doStatus" gorm:"type:varchar(50);comment:DoStatus"`
	Before       string    `json:"-" gorm:"type:text;comment:操作前数据"`
	Data         string    `json:"-" gorm:"type:text;comment:操作数据"`
	After        string    `json:"-" gorm:"type:text;comment:操作后数据"`
	Remark       string    `json:"remark" gorm:"type:varchar(255);comment:备注"`
	Diff         string    `json:"-" gorm:"type:text;comment:变更值"`
	OperatorId   int       `json:"operatorId" gorm:"type:int(11) unsigned;comment:操作人id"`
	OperatorName string    `json:"operatorName" gorm:"type:varchar(50);comment:操作人名称"`
	OperatorType int       `json:"operatorType" gorm:"type:tinyint(2);comment:操作人员类型(0:驿站admin)"`
	CreatedAt    time.Time `json:"createdAt" gorm:"comment:创建时间"`
}

func (OperateLogs) TableName() string {
	return "operate_logs"
}

func (e *OperateLogs) GetId() interface{} {
	return e.Id
}

func (e *OperateLogs) BeforeCreate(tx *gorm.DB) (err error) {
	if e.After == "" {
		e.After = "{}"
	}
	if e.Data == "" {
		e.Data = "{}"
	}
	if e.Before == "" {
		e.Before = "{}"
	}
	if e.Diff == "" {
		e.Diff = "[]"
	}
	return
}

func (e *OperateLogs) InsertItem(tx *gorm.DB) error {
	logType := OperateLogTypes{}
	beforeMap := map[string]interface{}{}
	afterMap := map[string]interface{}{}
	mappingMap := map[string]interface{}{}
	diffSlice := make(OperateLogDetailResp, 0)

	err := tx.Where("model_name = ?", e.ModelName).Where("type = ?", e.Type).Take(&logType).Error
	if err == nil {
		e.TypeName = logType.Name
		if e.After != "" {
			if errAf := json.Unmarshal([]byte(e.After), &afterMap); errAf == nil {
				if e.Before != "" {
					_ = json.Unmarshal([]byte(e.Before), &beforeMap)
				}
				if logType.Mapping != "" {
					_ = json.Unmarshal([]byte(logType.Mapping), &mappingMap)
				}
				diffData := CompareDiff(beforeMap, afterMap)
				FormatDiff(diffData, mappingMap, &diffSlice, "")
			}
		}
	}
	if diffBytes, err := json.Marshal(diffSlice); err == nil {
		e.Diff = string(diffBytes)
	}
	//数据库操作
	if err := tx.Create(e).Error; err != nil {
		return err
	}
	return nil
}

func CompareDiff(before, after map[string]interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	for aIndex, aVal := range after {
		if bVal, ok := before[aIndex]; ok {
			if assertAVal, ok := aVal.(map[string]interface{}); ok {
				if assertBVal, ok := bVal.(map[string]interface{}); ok {
					res[aIndex] = CompareDiff(assertBVal, assertAVal)
				} else {
					res[aIndex] = aVal
				}
			} else {
				if fmt.Sprintf("%v", aVal) != fmt.Sprintf("%v", bVal) {
					res[aIndex] = aVal
				}
			}
		} else {
			res[aIndex] = aVal
		}
	}
	return res
}

func SliceEmptyInterfacePrint(slice []interface{}) string {
	var res strings.Builder
	for _, val := range slice {
		res.WriteString(fmt.Sprintf("%v,", val))
	}
	return strings.TrimRight(res.String(), ",")
}

func FormatDiff(diffData, mapping map[string]interface{}, res *OperateLogDetailResp, prefix string) {
	for dIndex, dVal := range diffData {
		var desc interface{}
		assertMVal := map[string]interface{}{}
		if mVal, ok := mapping[dIndex]; ok {
			desc = mVal
		}
		switch dVal.(type) {
		case map[string]interface{}:
			assertDVal := dVal.(map[string]interface{})
			if desc != nil {
				if realAssertMVal, ok := desc.(map[string]interface{}); ok {
					assertMVal = realAssertMVal
				}
			}
			FormatDiff(assertDVal, assertMVal, res, getPrefix(prefix, dIndex))

		case []interface{}:
			assertDVal := dVal.([]interface{})
			logDiff := OperateLogDetail{
				Key:   getPrefix(prefix, dIndex),
				Desc:  "",
				Value: SliceEmptyInterfacePrint(assertDVal),
			}
			if descStr, ok := desc.(string); ok {
				logDiff.Desc = descStr
			}
			*res = append(*res, logDiff)

		default:
			logDiff := OperateLogDetail{
				Key:   getPrefix(prefix, dIndex),
				Desc:  "",
				Value: fmt.Sprintf("%v", dVal),
			}
			if descStr, ok := desc.(string); ok {
				logDiff.Desc = descStr
			}
			*res = append(*res, logDiff)
		}
	}
}

func getPrefix(prefix, dIndex string) string {
	if prefix == "" {
		return dIndex
	}
	return prefix + "." + dIndex
}

func AddOperateLog(tx *gorm.DB, dataId int, oldData, saveData, afterData interface{}, modelName, modelType string, operatorType, operatorId int, operatorName string) error {
	//记录操作日志
	oldDataStr := ""
	if saveData != nil {
		oldDataJson, _ := json.Marshal(oldData)
		oldDataStr = string(oldDataJson)
	}
	dataStr, _ := json.Marshal(saveData)
	afterDataStr, _ := json.Marshal(afterData)
	opLog := &OperateLogs{
		DataId:       strconv.Itoa(dataId),
		ModelName:    modelName,
		Type:         modelType,
		DoStatus:     "",
		Before:       oldDataStr,
		Data:         string(dataStr),
		After:        string(afterDataStr),
		OperatorId:   operatorId,
		OperatorName: operatorName,
		OperatorType: operatorType,
	}
	_ = opLog.InsertItem(tx)
	return nil
}
