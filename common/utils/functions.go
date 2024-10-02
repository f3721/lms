package utils

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"go-admin/config"
	"math"
	"math/big"
	mRand "math/rand"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/samber/lo"
)

func GetDateTimeString() string {
	now := time.Now()
	return now.Format("2006-01-02 15:04:05")
}

func GetDateTime() time.Time {
	now := time.Now().Unix()
	return time.Unix(now, 0)
}

func Rand(initNum int64) int64 {
	result, _ := rand.Int(rand.Reader, big.NewInt(initNum))
	return result.Int64()
}

// 生成N长度的随机数字
func GenRandNum(length int) int {
	mRand.Seed(time.Now().UnixNano())
	min := int(math.Pow10(length - 1))
	max := int(math.Pow10(length)) - 1
	return mRand.Intn(max-min+1) + min
}

// RandStr 生成长度为size的随机字符串
func RandStr(size int) string {
	words := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	var str string
	for i := 0; i < size; i++ {
		str += words[Rand(26)]
	}
	return str
}

// NumToString 生成某个位数的字符串
func NumToString(placeholder string, size int, num int) string {
	return fmt.Sprintf("%"+placeholder+strconv.Itoa(size)+"d", num)
}

// Decimal 保留两位小数
func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

func splitFunc() func(c rune) bool {
	return func(c rune) bool {
		if c == ',' || c == '，' || c == '|' || c == ' ' || c == ':' || c == '\n' {
			return true
		} else {
			return false
		}
	}
}

func Split(str string) []string {
	return strings.FieldsFunc(str, splitFunc())
}

func SplitToInt(str string) []int {
	slice := strings.FieldsFunc(str, splitFunc())
	return lo.Map(slice, func(item string, _ int) int {
		val, _ := strconv.Atoi(item)
		return val
	})
}

func SplitToInt32List(str string, sep string) (i32List []int64) {
	if str == "" {
		return
	}
	strList := strings.Split(str, sep)
	if len(strList) == 0 {
		return
	}
	for _, item := range strList {
		if item == "" {
			continue
		}
		val, err := strconv.ParseInt(item, 10, 64)
		if err != nil {
			continue
		}
		i32List = append(i32List, val)
	}
	return
}

func Splits(str string, sep string) (strList []string) {
	fields := make([]string, 0)
	if str == "" {
		return fields
	}
	strArr := strings.Split(str, sep)
	if len(strArr) == 0 {
		return fields
	}
	for _, item := range strArr {
		if item == "" {
			continue
		}
		strList = append(strList, item)
	}
	return
}

// StructColumn 参考：https:www.cnblogs.com/zhongweikang/p/12613389.html
func StructColumn(desk, input interface{}, columnKey, indexKey string) (err error) {
	structIndexColumn := func(desk, input interface{}, columnKey, indexKey string) (err error) {
		findStructValByIndexKey := func(curVal reflect.Value, elemType reflect.Type, indexKey, columnKey string) (indexVal, columnVal reflect.Value, err error) {
			indexExist := false
			columnExist := false
			for i := 0; i < elemType.NumField(); i++ {
				curField := curVal.Field(i)
				if elemType.Field(i).Name == indexKey {
					switch curField.Kind() {
					case reflect.String, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int, reflect.Float64, reflect.Float32:
						indexExist = true
						indexVal = curField
					default:
						return indexVal, columnVal, errors.New("indexKey must be int float or string")
					}
				}
				if elemType.Field(i).Name == columnKey {
					columnExist = true
					columnVal = curField
					continue
				}
			}
			if !indexExist {
				return indexVal, columnVal, errors.New(fmt.Sprintf("indexKey %s not found in %s's field", indexKey, elemType))
			}
			if len(columnKey) > 0 && !columnExist {
				return indexVal, columnVal, errors.New(fmt.Sprintf("columnKey %s not found in %s's field", columnKey, elemType))
			}
			return
		}

		deskValue := reflect.ValueOf(desk)
		if deskValue.Elem().Kind() != reflect.Map {
			return errors.New("desk must be map")
		}
		deskElem := deskValue.Type().Elem()
		if len(columnKey) == 0 && deskElem.Elem().Kind() != reflect.Struct {
			return errors.New(fmt.Sprintf("desk's elem expect struct, got %s", deskElem.Elem().Kind()))
		}

		rv := reflect.ValueOf(input)
		rt := reflect.TypeOf(input)
		elemType := rt.Elem()

		var indexVal, columnVal reflect.Value
		direct := reflect.Indirect(deskValue)
		mapReflect := reflect.MakeMap(deskElem)
		deskKey := deskValue.Type().Elem().Key()

		for i := 0; i < rv.Len(); i++ {
			curVal := rv.Index(i)
			indexVal, columnVal, err = findStructValByIndexKey(curVal, elemType, indexKey, columnKey)
			if err != nil {
				return
			}
			if deskKey.Kind() != indexVal.Kind() {
				return errors.New(fmt.Sprintf("cant't convert %s to %s, your map'key must be %s", indexVal.Kind(), deskKey.Kind(), indexVal.Kind()))
			}
			if len(columnKey) == 0 {
				mapReflect.SetMapIndex(indexVal, curVal)
				direct.Set(mapReflect)
			} else {
				if deskElem.Elem().Kind() != columnVal.Kind() {
					return errors.New(fmt.Sprintf("your map must be map[%s]%s", indexVal.Kind(), columnVal.Kind()))
				}
				mapReflect.SetMapIndex(indexVal, columnVal)
				direct.Set(mapReflect)
			}
		}
		return
	}

	structColumn := func(desk, input interface{}, columnKey string) (err error) {
		findStructValByColumnKey := func(curVal reflect.Value, elemType reflect.Type, columnKey string) (columnVal reflect.Value, err error) {
			columnExist := false
			for i := 0; i < elemType.NumField(); i++ {
				curField := curVal.Field(i)
				if elemType.Field(i).Name == columnKey {
					columnExist = true
					columnVal = curField
					continue
				}
			}
			if !columnExist {
				return columnVal, errors.New(fmt.Sprintf("columnKey %s not found in %s's field", columnKey, elemType))
			}
			return
		}

		if len(columnKey) == 0 {
			return errors.New("columnKey cannot not be empty")
		}

		deskElemType := reflect.TypeOf(desk).Elem()
		if deskElemType.Kind() != reflect.Slice {
			return errors.New("desk must be slice")
		}

		rv := reflect.ValueOf(input)
		rt := reflect.TypeOf(input)

		var columnVal reflect.Value
		deskValue := reflect.ValueOf(desk)
		direct := reflect.Indirect(deskValue)

		for i := 0; i < rv.Len(); i++ {
			columnVal, err = findStructValByColumnKey(rv.Index(i), rt.Elem(), columnKey)
			if err != nil {
				return
			}
			if deskElemType.Elem().Kind() != columnVal.Kind() {
				return errors.New(fmt.Sprintf("your slice must be []%s", columnVal.Kind()))
			}

			direct.Set(reflect.Append(direct, columnVal))
		}
		return
	}

	deskValue := reflect.ValueOf(desk)
	if deskValue.Kind() != reflect.Ptr {
		return errors.New("desk must be ptr")
	}

	rv := reflect.ValueOf(input)
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return errors.New("input must be map slice or array")
	}

	rt := reflect.TypeOf(input)
	if rt.Elem().Kind() != reflect.Struct {
		return errors.New("input's elem must be struct")
	}

	if len(indexKey) > 0 {
		return structIndexColumn(desk, input, columnKey, indexKey)
	}
	return structColumn(desk, input, columnKey)
}

func ContainChinese(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.Han, r) {
			return true
		}
	}
	return false
}

// EncryptMd5 Md5加密
func EncryptMd5(str, salt string) string {
	data := []byte(str + salt)
	md5New := md5.New()
	md5New.Write(data)
	return hex.EncodeToString(md5New.Sum(nil))
}

// Md5Uc 是用户密码算法
func Md5Uc(password string) string {
	salt := "rick123"
	data := []byte(strings.ToLower(salt) + password)
	hash := md5.Sum(data)
	return strings.ToUpper(hex.EncodeToString(hash[:]))
}

// 校验密码 | 8-20位，由大写字母、小写字母、数字和英文字符(除空格外)至少三种组成
func ValidatePassword(password string) bool {
	// 校验用户密码|
	if len(password) < 8 || len(password) > 20 {
		return false
	}
	// 检查是否包含大写字母
	hasUpper := regexp.MustCompile("[A-Z]").MatchString(password)
	// 检查是否包含小写字母
	hasLower := regexp.MustCompile("[a-z]").MatchString(password)
	// 检查是否包含数字
	hasDigit := regexp.MustCompile("[0-9]").MatchString(password)
	// 检查是否包含除空格以外的英文字符
	hasSpecial := regexp.MustCompile(`[^\s\w]`).MatchString(password)
	// 检查满足至少三种字符类型的要求
	typesCount := 0
	if hasLower {
		typesCount++
	}
	if hasUpper {
		typesCount++
	}
	if hasDigit {
		typesCount++
	}
	if hasSpecial {
		typesCount++
	}
	return typesCount >= 3
}

// 校验密码 | 8-20位，由大写字母、小写字母、数字和英文字符(除空格外)至少二种组成
func ValidateLoginName(password string) bool {
	// 校验用户密码|
	if len(password) < 8 || len(password) > 20 {
		return false
	}
	// 检查是否包含大写字母
	hasUpper := regexp.MustCompile("[A-Z]").MatchString(password)
	// 检查是否包含小写字母
	hasLower := regexp.MustCompile("[a-z]").MatchString(password)
	// 检查是否包含数字
	hasDigit := regexp.MustCompile("[0-9]").MatchString(password)
	// 检查是否包含除空格以外的英文字符
	hasSpecial := regexp.MustCompile(`[^\s\w]`).MatchString(password)
	// 检查满足至少三种字符类型的要求
	typesCount := 0
	if hasLower {
		typesCount++
	}
	if hasUpper {
		typesCount++
	}
	if hasDigit {
		typesCount++
	}
	if hasSpecial {
		typesCount++
	}
	return typesCount >= 2
}

// Sha1 加密值
func Sha1(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	hash := h.Sum(nil)
	return hex.EncodeToString(hash)
}

func MergeMap(mObj ...map[string]string) map[string]string {
	newObj := map[string]string{}
	for _, m := range mObj {
		for k, v := range m {
			newObj[k] = v
		}
	}
	return newObj
}

func IsAllDigits(s string) bool {
	re := regexp.MustCompile(`^\d+$`)
	return re.MatchString(s)
}

func IsPhoneNumber(s string) bool {
	re := regexp.MustCompile(`^1[3456789]\d{9}$`)
	return re.MatchString(s)
}

func ValidateEmailFormat(email string) bool {
	// 使用正则表达式匹配邮箱格式
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, err := regexp.MatchString(regex, email)
	if err != nil {
		// 正则表达式匹配出错，视为邮箱格式不正确
		return false
	}
	return match
}

func ReverseSlice(slice interface{}) {
	sliceValue := reflect.ValueOf(slice)
	if sliceValue.Kind() != reflect.Slice {
		panic("Slice type is required")
	}

	length := sliceValue.Len()
	swap := reflect.Swapper(slice)

	for i, j := 0, length-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}

func IsAlphanumeric(word string) bool {
	isAlphanumeric := regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(word)
	return isAlphanumeric
}

// 根据时间获取当前税率
func GetCurrentTaxRate() float64 {
	return 0.13
}

// 获取前后台Host地址 | hType: 0-前台，1-后台
func GetHostUrl(hType int) string {
	hostStr := ""
	if hType == 0 {
		hostStr = config.ExtConfig.MallHost
	}
	if hType == 1 {
		hostStr = config.ExtConfig.LmsHost
	}
	return hostStr
}

func SlicePage(page, pageSize, nums int64) (sliceStart, sliceEnd int64) {
	if page <= 0 {
		page = 1
	}
	if pageSize < 0 {
		pageSize = 20 //设置一页默认显示的记录数
	}
	if pageSize > nums {
		return 0, nums
	}
	// 总页数
	pageCount := int64(math.Ceil(float64(nums) / float64(pageSize)))
	if page > pageCount {
		return 0, 0
	}
	sliceStart = (page - 1) * pageSize
	sliceEnd = sliceStart + pageSize

	if sliceEnd > nums {
		sliceEnd = nums
	}
	return sliceStart, sliceEnd
}
