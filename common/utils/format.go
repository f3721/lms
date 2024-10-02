package utils

import (
	"github.com/shopspring/decimal"
	"reflect"
	"strings"
	"time"
)

// StructToMap 结构体转换为map（会继承一级的结构体字段,并且字段首字母小写） 例如：
//type Person struct {
//	Name    string
//	Age     int
//	Address
//}
//
//type Address struct {
//	City    string
//	Country string
//}
// 转换为 map[Age:30 City:New York Country:USA Name:John]

func StructToMap(s interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	value := reflect.ValueOf(s)

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldName := FirstLower(value.Type().Field(i).Name)

		switch field.Kind() {
		case reflect.Struct:
			if _, ok := field.Interface().(time.Time); ok {
				result[fieldName] = field.Interface()
			} else {
				nestedValue := reflect.ValueOf(field.Interface())
				for j := 0; j < nestedValue.NumField(); j++ {
					nestedField := nestedValue.Field(j)
					nestedFieldName := FirstLower(nestedValue.Type().Field(j).Name)
					//nestedMap[nestedFieldName] = nestedField.Interface()
					result[nestedFieldName] = nestedField.Interface()
				}
			}
			//
			//result[fieldName] = nestedMap
		default:
			result[fieldName] = field.Interface()
		}
	}

	return result
}

// FirstUpper 字符串首字母大写
func FirstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// FirstLower 字符串首字母小写
func FirstLower(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

// FormatTime 将map中的time.Time 类型 格式化为 年-月-日 时:分:秒
func FormatMapTime(m map[string]interface{}) map[string]interface{} {
	for i, v := range m {
		if v2, ok := v.(time.Time); ok {
			m[i] = TimeFormat(v2)
		}
	}
	return m
}

// TrimMapSpace 将map数组中的每个值去除左右两边的空格
func TrimMapSpace(maps []map[string]interface{}) []map[string]interface{} {
	for _, m := range maps {
		for k, v := range m {
			m[k] = strings.TrimSpace(v.(string))
		}
	}
	return maps
}

func PriceToFixed(price string) string {
	newPrice, _ := decimal.NewFromString(price)
	return newPrice.RoundDown(2).StringFixed(2)
}
