package dto

type PrintsTitle struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	Logo  string `json:"logo"`
}

type PrintsForm struct {
	Type       string `json:"type"`
	LabelWidth string `json:"labelWidth"`
	Top        string `json:"top"`
	ItemTop    string `json:"itemTop"`
	List       any    `json:"list"`
}

type PrintsFormListItem struct {
	Value any    `json:"value"`
	Label string `json:"label"`
	Span  int8   `json:"span"`
}

type PrintsValue struct {
	Value any `json:"value"`
}

type PrintsTitleValue struct {
	Value string `json:"value"`
	W     string `json:"w"`
}

type PrintsTableDefaultStyle struct {
	W           string `json:"w"`
	Align       string `json:"align"`
	BorderColor string `json:"borderColor"`
	HeadColor   string `json:"headColor"`
	BodyColor   string `json:"bodyColor"`
}
type PrintsTable struct {
	Type         string                  `json:"type"`
	Top          string                  `json:"top"`
	Title        string                  `json:"title"`
	Desc         string                  `json:"desc"`
	DefaultStyle PrintsTableDefaultStyle `json:"defaultStyle"`
	List         any                     `json:"list"`
}

type PrintsSign struct {
	Type    string `json:"type"`
	Bold    bool   `json:"bold"`
	Width   string `json:"width"`
	Top     string `json:"top"`
	ItemTop string `json:"itemTop"`
	Col     bool   `json:"col"`
	List    any    `json:"list"`
}

type PrintsSignLabel struct {
	Label string `json:"label"`
}

type FontWeigntStyle struct {
	FontWeight string `json:"fontWeight"`
}
type PrintsFormListItemStyle struct {
	Value      string          `json:"value"`
	Lable      string          `json:"lable"`
	Span       int8            `json:"span"`
	ValueStyle FontWeigntStyle `json:"valueStyle"`
}
