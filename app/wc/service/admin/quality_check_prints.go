package admin

import (
	"go-admin/app/wc/service/admin/dto"
	"go-admin/common/service"
)

type QualityCheckPrints struct {
	service.Service
	Prints  []interface{}
	dataMap map[string]interface{}
}

func (s *QualityCheckPrints) GetPrints(dataMap map[string]interface{}) {
	s.dataMap = dataMap
	// 标题
	s.Prints = append(s.Prints, s.getPrintsTitle())
	// 头部表格
	formUp := s.getPrintsFormUp()
	// 数据表
	table := s.getPrintsTable()
	// 尾部表单
	formDown := s.getPrintsFormDown()

	s.Prints = append(s.Prints, formUp, table, formDown)
}

func (s *QualityCheckPrints) getPrintsTitle() dto.PrintsTitle {

	return dto.PrintsTitle{
		Type:  "title",
		Title: "质检报告",
		Logo:  "",
	}
}

func (s *QualityCheckPrints) getPrintsFormUp() dto.PrintsForm {
	var listItems []dto.PrintsFormListItem

	listItems = append(listItems, dto.PrintsFormListItem{
		Label: "来源方：",
		Value: s.dataMap["sourceName"],
		Span:  8,
	}, dto.PrintsFormListItem{
		Label: "来源单号：",
		Value: s.dataMap["sourceCode"],
		Span:  8,
	}, dto.PrintsFormListItem{
		Label: "入库单号：",
		Value: s.dataMap["entryCode"],
		Span:  8,
	}, dto.PrintsFormListItem{
		Label: "质检数量：",
		Value: s.dataMap["quantityNum"],
		Span:  24,
	})

	return dto.PrintsForm{
		Type:       "form",
		LabelWidth: "100px",
		Top:        "5px",
		ItemTop:    "5px",
		List:       listItems,
	}
}

func (s *QualityCheckPrints) getPrintsTable() dto.PrintsTable {
	productsChunks, ok := s.dataMap["details"].([]map[string]interface{})
	if !ok {
		return dto.PrintsTable{}
	}
	var productsChunkAll []interface{}

	// 处理标题
	titleList := []dto.PrintsTitleValue{}
	titleList = append(titleList,
		dto.PrintsTitleValue{Value: "ID", W: "40px"},
		dto.PrintsTitleValue{Value: "质检项", W: "40px"},
		dto.PrintsTitleValue{Value: "检测结果", W: "40px"},
		dto.PrintsTitleValue{Value: "备注", W: ""},
		dto.PrintsTitleValue{Value: "检测人员", W: "40px"},
	)
	productsChunkAll = append(productsChunkAll, titleList)

	// 处理值
	for _, item := range productsChunks {
		var listItem []dto.PrintsValue
		listItem = append(listItem,
			dto.PrintsValue{Value: item["id"]},
			dto.PrintsValue{Value: item["qualityCheckOption"]},
			dto.PrintsValue{Value: item["qualityRes"]},
			dto.PrintsValue{Value: item["remark"]},
			dto.PrintsValue{Value: item["qualityByName"]},
		)
		productsChunkAll = append(productsChunkAll, listItem)
	}

	// 再包一层数组
	newList := []interface{}{productsChunkAll}

	// 返回Table
	return dto.PrintsTable{
		Type: "table",
		Top:  "20px",
		DefaultStyle: dto.PrintsTableDefaultStyle{
			W:           "80px",
			Align:       "center",
			BorderColor: "#000",
			HeadColor:   "#000",
			BodyColor:   "#000",
		},
		List: newList,
	}
}

func (s *QualityCheckPrints) getPrintsFormDown() dto.PrintsForm {
	var listItems []dto.PrintsFormListItem

	listItems = append(listItems, dto.PrintsFormListItem{
		Label: "质检结论:",
		Value: "",
		Span:  24,
	}, dto.PrintsFormListItem{
		Label: "质检结果：",
		Value: s.dataMap["qualityResName"],
		Span:  12,
	}, dto.PrintsFormListItem{
		Label: "质检时间：",
		Value: s.dataMap["qualityTime"],
		Span:  12,
	}, dto.PrintsFormListItem{
		Label: "质检人：",
		Value: s.dataMap["qualityByName"],
		Span:  24,
	})

	return dto.PrintsForm{
		Type:       "form",
		LabelWidth: "100px",
		Top:        "5px",
		ItemTop:    "5px",
		List:       listItems,
	}
}
