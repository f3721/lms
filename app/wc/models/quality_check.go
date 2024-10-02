package models

import (
	"go-admin/common/models"
	"time"
)

const (
	QualityStatusInit = 0
	QualityStatusOk   = 1

	QualityType1 = 1
	QualityType2 = 2

	QualityStatus     = 0
	QualityStatusPart = 1
	QualityStatusAll  = 2

	QualityResInit  = 0
	QualityResOk    = 1
	QualityResFalse = 2
)

var QsStatusName = map[int]string{
	QualityStatusInit: "未审批",
	QualityStatusOk:   "审批",
}

var QsTypeName = map[int]string{
	QualityType1: "全检",
	QualityType2: "抽检",
}

// 0未质检，1部分质检，2全部质检
var QualityStatusName = map[int]string{
	QualityStatus:     "未质检",
	QualityStatusPart: "部分质检",
	QualityStatusAll:  "全部质检",
}

var QualityResName = map[int]string{
	QualityResInit:  "未质检",
	QualityResOk:    "合格",
	QualityResFalse: "不合格",
}

type QualityCheck struct {
	models.Model

	QualityCheckCode   string               `json:"qualityCheckCode" gorm:"type:varchar(32);comment:质检单号code"`
	SourceCode         string               `json:"sourceCode" gorm:"type:varchar(32);comment:来源单据code"`
	EntryCode          string               `json:"entryCode" gorm:"type:varchar(20);comment:入库单编码"`
	WarehouseCode      string               `json:"warehouseCode" gorm:"type:varchar(20);comment:实体仓code"`
	LogicWarehouseCode string               `json:"logicWarehouseCode" gorm:"type:varchar(20);comment:逻辑仓code"`
	SourceName         string               `json:"sourceName" gorm:"type:varchar(100);comment:来源方"`
	Type               int                  `json:"type" gorm:"type:tinyint unsigned;comment:质检类型,1全检，2抽检"`
	SkuCode            string               `json:"skuCode" gorm:"type:varchar(20);comment:skuCode"`
	StayQualityNum     int                  `json:"stayQualityNum" gorm:"type:int;comment:待质检数量"`
	QuantityNum        int                  `json:"quantityNum" gorm:"type:int;comment:质检数量"`
	QualityStatus      int                  `json:"qualityStatus" gorm:"type:tinyint unsigned;comment:质检进度：0未质检，1部分质检，2全部质检"`
	QualityRes         int                  `json:"qualityRes" gorm:"type:tinyint unsigned;comment:质检结果：1合格,2不合格"`
	QualityTime        time.Time            `json:"qualityTime" gorm:"type:datetime;comment:质检时间"`
	QuantitySign       string               `json:"quantitySign" gorm:"type varchar(1000);comment:电子签照片"`
	Unqualified        int                  `json:"unqualified" gorm:"type:tinyint unsigned;comment:不合格处理办法: 0-拒收 1-异常填报"`
	Status             int                  `json:"status" gorm:"type:tinyint unsigned;comment:状态:0待审批，1已审批"`
	QualityCheckDetail []QualityCheckDetail `json:"-" gorm:"foreignKey:QualityCheckTaskId;references:Id"`
	models.ModelTime
	models.ControlBy
}

func (QualityCheck) TableName() string {
	return "quality_check_task"
}

func (e *QualityCheck) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *QualityCheck) GetId() interface{} {
	return e.Id
}
