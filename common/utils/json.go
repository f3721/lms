package utils

import (
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	"os"
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
					// 时间格式字符串 格式化
					if aValTime, err := time.Parse(time.RFC3339, fmt.Sprintf("%v", aVal)); err == nil {
						aVal = TimeFormat(aValTime)
					}
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
		case float64:
			float64Val := dVal.(float64)
			tmpVal := fmt.Sprintf("%v", dVal)
			// 浮点数为整型特殊处理
			if float64Val== float64(int64(float64Val)) {
				tmpVal = strconv.FormatFloat(float64Val, 'f', 0, 64)
			}
			logDiff := OperateLogDetail{
				Key:   getPrefix(prefix, dIndex),
				Desc:  "",
				Value: tmpVal,
			}
			if descStr, ok := desc.(string); ok {
				logDiff.Desc = descStr
			}
			*res = append(*res, logDiff)
		default:
			valstr := fmt.Sprintf("%v", dVal)
			if tmpVal, ok := dVal.(string); ok {
				// 时间格式字符串 格式化
				if tmpVal == "0001-01-01T00:00:00Z" ||  tmpVal == "<nil>" {
					valstr = ""
				} else if aValTime, err := time.Parse(time.RFC3339, fmt.Sprintf("%v", tmpVal)); err == nil {
					valstr = TimeFormat(aValTime)
				}
			}
			logDiff := OperateLogDetail{
				Key:   getPrefix(prefix, dIndex),
				Desc:  "",
				Value: valstr,
			}
			if descStr, ok := desc.(string); ok {
				logDiff.Desc = descStr
			}
			*res = append(*res, logDiff)
		}
	}
	*res = lo.Filter[OperateLogDetail](*res, func(diff OperateLogDetail, index int) bool {
		return !lo.Contains[string]([]string{
			"id",
			"Id",
			"ID",
			"created_at",
			"createdAt",
			"updated_at",
			"updatedAt",
			"deleted_at",
			"deletedAt",
			"createBy",
			"updateBy",
			"updateByName",
			"createByName",
			"create_time",
			"update_time",
		}, diff.Key)
	})
}

func getPrefix(prefix, dIndex string) string {
	if prefix == "" {
		return dIndex
	}
	return prefix + "." + dIndex
}

func marshal(j interface{}) string {
	value, _ := json.Marshal(j)
	return string(value)
}

func LoadJson(path string, dist interface{}) (err error) {
	var content []byte
	if content, err = os.ReadFile(path); err == nil {
		err = json.Unmarshal(content, dist)
	}
	return err
}
