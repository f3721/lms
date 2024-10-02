package dto

import (
	"go-admin/app/wc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
	"gorm.io/gorm"
)

type QualityCheckGetPageReq struct {
	dto.Pagination `search:"-"`
	QualityCheck
	common.ControlBy
}

type QualityCheck struct {
	Id                 int    `form:"id"`
	QualityCheckCode   string `form:"qualityCheckCode" commit:"质检单号"`
	SourceCode         string `form:"sourceCode" commit:"来源单号"`
	EntryCode          string `form:"entryCode" commit:"入库单号"`
	SourceName         string `form:"sourceName" commit:"来源方"`
	Status             int    `form:"status" commit:"状态"`
	WarehouseCode      string `form:"warehouseCode" commit:"入库实体仓"`
	LogicWarehouseCode string `form:"logicWarehouseCode" commit:"入库逻辑仓"`
	QualityRes         int    `form:"qualityRes" commit:"质检结果"`
	Type               int    `form:"type" commit:"质检类型"`
	QualityStatus      int    `form:"qualityStatus" commit:"质检进度"`
	SkuCode            string `form:"skuCode" commit:"skuCode"`
	GoodsName          string `form:"goodsName" commit:"商品名称"`
	StartDate          string `form:"startDate" commit:"开始日期"`
	EndDate            string `form:"endDate" commit:"结束日期"`
}

type QualityCheckRes struct {
	models.QualityCheck
	StatusName        string `gorm:"-"` //状态:0待审批，1已审批
	TypeName          string `gorm:"-"` //质检类型,1全检，2抽检
	QualityStatusName string `gorm:"-"` //质检进度：0未质检，1部分质检，2全部质检
	QualityResName    string `gorm:"-"` //质检结果：1合格,2不合格
	UnqualifiedName   string `gorm:"-"` //不合格处理办法: 0-拒收 1-异常填报
}

type QualityExportRes struct {
	ID                 string `json:"id"`
	QualityCheckCode   string `json:"qualityCheckCode"`
	SourceCode         string `json:"sourceCode"`
	EntryCode          string `json:"entryCode"`
	WarehouseCode      string `json:"warehouseCode"`
	LogicWarehouseCode string `json:"logicWarehouseCode"`
	SourceName         string `json:"sourceName"`
	Status             string `json:"status"`
	Type               string `json:"type"`
	SkuCode            string `json:"skuCode"`
	QualityStatus      string `json:"qualityStatus"`
	StayQualityNum     string `json:"stayQualityNum"`
	QuantityNum        string `json:"quantityNum"`
	QualityRes         string `json:"qualityRes"`
	QualityTime        string `json:"qualityTime"`
}

func (m *QualityCheckGetPageReq) GetNeedSearch() interface{} {
	return *m
}

func (s *QualityCheckGetPageReq) Generate(model *models.QualityCheckConfig) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.Status = s.Status
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.CreateByName = s.CreateByName
	model.UpdateByName = s.UpdateByName
}

func (s *QualityCheckGetPageReq) GetId() interface{} {
	return s.Id
}

type QualityCheckGetReq struct {
	Id int `uri:"id"`
}

func (s *QualityCheckGetReq) GetId() interface{} {
	return s.Id
}

type QualityCheckDetailGetReq struct {
	Id int `uri:"id"`
}

func (s *QualityCheckDetailGetReq) GetId() interface{} {
	return s.Id
}

// 上传质检明细
type QualityCheckUpdateReq struct {
	Id            int             `uri:"id"`
	QuantityNum   int             `json:"quantityNum" vd:"$>0; msg:'质检数量不能为空'"`
	QualityRes    int             `json:"qualityRes" vd:"$>0; msg:'请输入质检结论'"`
	QualityBy     int             `json:"qualityBy" vd:"$>0; msg:'请输入质检人'"`
	Remark        string          `json:"remark"`
	QualityOption []QualityOption `json:"qualityOption"`
}

func (s *QualityCheckUpdateReq) GetId() interface{} {
	return s.Id
}

// 质检详细
type QualityOption struct {
	QualityCheckOption string `json:"qualityCheckOption"`
	QualityRes         int    `json:"qualityRes"`
	QualityBy          int    `json:"qualityBy"`
	Remark             string `json:"remark"`
}

func GenQualityCheckSearch(req *QualityCheckGetPageReq, tableAlias string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if req.QualityCheckCode != "" {
			db.Where(tableAlias+".quality_check_code IN ?", req.QualityCheckCode)
		}
		if req.SourceCode != "" {
			db.Where(".source_code IN ?", req.SourceCode)
		}
		if req.EntryCode != "" {
			db.Where(".entry_code IN ?", req.EntryCode)
		}
		if req.SourceName != "" {
			db.Where(".source_name IN ?", req.SourceName)
		}
		if req.WarehouseCode != "" {
			db.Where(".warehouse_code IN ?", req.WarehouseCode)
		}
		if req.LogicWarehouseCode != "" {
			db.Where(".logic_warehouse_code IN ?", req.LogicWarehouseCode)
		}
		if req.Status != 0 {
			db.Where(".status = ?", req.Status)
		}
		if req.Type != 0 {
			db.Where(".type = ?", req.Type)
		}
		if req.QualityStatus != 0 {
			db.Where(".quality_status = ?", req.QualityStatus)
		}
		if req.SkuCode != "" {
			db.Where(".sku_code = ?", req.SkuCode)
		}
		if req.StartDate != "" {
			db.Where(".quality_time > ?", req.StartDate)
		}
		if req.EndDate != "" {
			db.Where(".quality_time < ?", req.EndDate)
		}

		return db
	}
}
