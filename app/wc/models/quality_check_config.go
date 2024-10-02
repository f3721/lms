package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"go-admin/common/models"
)

// 约束字段配置项
var ConstraintMap = map[string]map[string]map[string]any{
	"skuConstraint": {
		"sku_code":          {"order": 1, "name": "产品SKU", "operator": "", "targetValue": []string{}},
		"supplier_sku_code": {"order": 2, "name": "供应商code", "operator": "", "targetValue": []string{}},
		"name_zh":           {"order": 3, "name": "中文名", "operator": "", "targetValue": []string{}},
		"sales_uom":         {"order": 4, "name": "售卖包装单位", "operator": "", "targetValue": []string{}},
		"physical_uom":      {"order": 5, "name": "物理单位", "operator": "", "targetValue": []string{}},
		"fragile_flag":      {"order": 6, "name": "易碎标志", "operator": "=", "targetValue": []string{"否", "是"}},
		"hazard_flag":       {"order": 7, "name": "危险品标志", "operator": "=", "targetValue": []string{"否", "是"}},
		"hazard_class":      {"order": 8, "name": "危险品等级", "operator": "=", "targetValue": []string{"1", "2", "3", "4", "5"}},
		"bulky_flag":        {"order": 9, "name": "抛货标志", "operator": "=", "targetValue": []string{"否", "是"}},
		"assemble_flag":     {"order": 10, "name": "拼装件标志", "operator": "=", "targetValue": []string{"否", "是"}},
		"is_valuables":      {"order": 11, "name": "是否贵重品", "operator": "=", "targetValue": []string{"否", "是"}},
		"is_fluid":          {"order": 12, "name": "是否液体", "operator": "=", "targetValue": []string{"否", "是"}},
		"consumptive_flag":  {"order": 13, "name": "耗材标志", "operator": "=", "targetValue": []string{"否", "是"}},
		"storage_flag":      {"order": 14, "name": "保存期标志", "operator": "=", "targetValue": []string{"否", "是"}},
		"storage_time":      {"order": 15, "name": "保存期限（月）", "operator": "", "targetValue": []string{}},
		"refund_flag":       {"order": 16, "name": "可退换货标志", "operator": "", "targetValue": []string{}},
		"custom_made_flag":  {"order": 17, "name": "是否定制品", "operator": "=", "targetValue": []string{"否", "是"}},
		"tax":               {"order": 18, "name": "税率", "operator": "", "targetValue": []string{"6%", "9%", "13%"}},
	},
	"categoryConstraint": {
		"cate_level1": {"order": 1, "name": "一级产线", "operator": "=", "targetValue": []string{"一级产线"}},
		"cate_level2": {"order": 2, "name": "二级产线", "operator": "=", "targetValue": []string{"二级产线"}},
		"cate_level3": {"order": 3, "name": "三级产线", "operator": "=", "targetValue": []string{"三级产线"}},
		"cate_level4": {"order": 4, "name": "四级产线", "operator": "=", "targetValue": []string{"四级产线"}},
	},
}

// 约束字段映射数据库值
var ConstraintValMap = map[string]string{
	"否":    "0",
	"是":    "1",
	"6%":   "0.6",
	"9%":   "0.9",
	"13%":  "0.13",
	"一级产线": "1",
	"二级产线": "2",
	"三级产线": "3",
	"四级产线": "4",
}

// 约束结构体
type Constraint struct {
	Relation    string `json:"relation" comment:"关联关系" vd:"@:in($,'and','or'); msg:'约束关系范围[and,or]'"`
	Field       string `json:"field" comment:"字段" vd:"@:len($)>0; msg:'约束字段必填'"`
	Operator    string `json:"operator" comment:"操作符" vd:"@:in($,'like','=','in'); msg:'约束操作符范围[like,=,in]'"`
	TargetValue string `json:"targetValue" comment:"对比值" vd:"@:len($)>0; msg:'约束值必填'"`
}

type Constraints []*Constraint

func (c *Constraints) Scan(value interface{}) error {
	if data, ok := value.([]byte); ok {
		return json.Unmarshal(data, c)
	}
	return errors.New("解析(Scan)约束错误")
}

func (c Constraints) Value() (driver.Value, error) {
	return json.Marshal(c)
}

type QualityCheckConfig struct {
	models.Model

	Status                   int                         `json:"status" gorm:"type:tinyint unsigned;comment:质检开关 0-关 1-开"`
	Type                     int                         `json:"type" gorm:"type:tinyint unsigned;comment:质检类型 0-全检  1-抽检"`
	SamplingNum              int                         `json:"samplingNum" gorm:"type:int unsigned;comment:抽检数量"`
	OrderType                int                         `json:"orderType" gorm:"type:tinyint unsigned;comment:质检订单类型:0-全部 1-采购入库 2-大货"`
	SkuConstraint            Constraints                 `json:"skuConstraint" gorm:"type:json;comment:sku约束"`
	CategoryConstraint       Constraints                 `json:"categoryConstraint" gorm:"type:json;comment:产线约束"`
	QualityCheckRoles        string                      `json:"qualityCheckRoles" gorm:"type:varchar(100);comment:质检角色IDs"`
	QualityCheckOptions      string                      `json:"qualityCheckOptions" gorm:"type:varchar(500);comment:质检内容"`
	Unqualified              int                         `json:"unqualified" gorm:"type:tinyint unsigned;comment:不合格处理办法: 0-拒收 1-异常填报"`
	QualityCheckConfigDetail []*QualityCheckConfigDetail `json:"-" gorm:"foreignKey:ConfigId;references:Id"`
	CreateByName             string                      `json:"createByName" gorm:"type:varchar(20);comment:创建人名称"`
	UpdateByName             string                      `json:"updateByName" gorm:"type:varchar(20);comment:修改人名称"`
	models.ModelTime
	models.ControlBy
}

func (QualityCheckConfig) TableName() string {
	return "quality_check_config"
}

func (e *QualityCheckConfig) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *QualityCheckConfig) GetId() interface{} {
	return e.Id
}
