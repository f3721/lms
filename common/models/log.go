package models

type ModelLog struct {
	DataId     int    `json:"dataId" gorm:"type:int;comment:关联的主键ID"`
	Type       string `json:"type" gorm:"type:varchar(10);comment:变更类型"`
	Data       string `json:"data" gorm:"type:json;comment:源数据"`
	BeforeData string `json:"beforeData" gorm:"type:json;comment:变更前数据"`
	AfterData  string `json:"afterData" gorm:"type:json;comment:变更后数据"`
}

func (e *ModelLog) SetDataId(dataId int) {
	e.DataId = dataId
}

func (e *ModelLog) SetType(types string) {
	e.Type = types
}

func (e *ModelLog) SetData(data string) {
	e.Data = data
}

func (e *ModelLog) SetBeforeData(beforeData string) {
	e.BeforeData = beforeData
}

func (e *ModelLog) SetAfterData(afterData string) {
	e.AfterData = afterData
}
