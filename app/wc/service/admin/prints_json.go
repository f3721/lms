package admin

import (
	"encoding/json"
	"fmt"
	"go-admin/app/wc/models"
	"go-admin/app/wc/service/admin/dto"
	"go-admin/common/actions"
	"go-admin/common/utils"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"github.com/samber/lo"
)

type PrintsJson struct {
	service.Service
}

func (e *PrintsJson) OutboundPrint(d *dto.CommonPrintsReq, p *actions.DataPermission, s *StockOutbound, outData *dto.CommonPrintsResp) error {
	req := &dto.StockOutboundGetReq{
		Id: d.Id,
	}
	resp := &dto.StockOutboundNoLocationResp{}
	if err := s.GetNoLocation(req, p, resp); err != nil {
		e.Log.Errorf("PrintsJson OutboundPrint error:%s \r\n", err)
		return err
	}
	outDataMap := utils.GetChunkProductsForPrint(resp, "stockOutboundProducts", map[string]int{
		//"productName": 10,
		//"mfgModel":    8,
		//"brandName":   6,
		"vendorName": 80,
		//"productNo":   10,
	}, 28, 50)
	outDataMap["addressFullName"] = resp.AddressFullName
	stockOutboundPrints := StockOutboundPrints{}
	stockOutboundPrints.GetPrints(outDataMap)
	outData.Prints = stockOutboundPrints.Prints
	return nil
}

func (e *PrintsJson) StockPrintPicking(d *dto.CommonPrintsReq, p *actions.DataPermission, s *StockOutbound, outData *dto.CommonPrintsResp) error {
	req := &dto.StockOutboundGetReq{
		Id: d.Id,
	}
	resp := &dto.StockOutboundGetResp{}
	if err := s.Get(req, p, resp); err != nil {
		e.Log.Errorf("PrintsJsonService PrintPicking error:%s \r\n", err)
		return err
	}
	// 部分出库-打印剩余未出库数量
	for _, item := range resp.StockOutboundProducts {
		item.OutboundProductSubCustom.LocationQuantity = item.OutboundProductSubCustom.LocationQuantity - item.OutboundProductSubCustom.LocationActQuantity
	}

	outDataMap := utils.GetChunkProductsForPrint(resp, "stockOutboundProducts", map[string]int{
		//"productName":  12,
		//"mfgModel":     6,
		//"brandName":    6,
		//"locationCode": 10,
		"vendorName": 80,
	}, 39, 46)

	stockPrintPicking := StockPrintsPicking{}
	stockPrintPicking.GetPrints(outDataMap, resp.Type)
	outData.Prints = stockPrintPicking.Prints

	return nil
}

func (e *PrintsJson) StockEntryPrints(d *dto.CommonPrintsReq, p *actions.DataPermission, s *StockEntry, outData *dto.CommonPrintsResp) error {
	req := &dto.StockEntryGetReq{
		Id:   d.Id,
		Type: "print",
	}
	resp := &dto.StockEntryGetResp{}
	if err := s.Get(req, p, resp); err != nil {
		e.Log.Errorf("PrintsJsonService StockEntryPrints error:%s \r\n", err)
		return err
	}
	outDataMap := utils.GetChunkProductsForPrint(resp, "stockEntryProducts", map[string]int{
		//"productName": 12,
		//"mfgModel":    6,
		//"brandName":   6,
		"vendorName": 80,
	}, 39, 46)
	stockEntryPrint := StockEntryPrints{}
	stockEntryPrint.GetPrints(outDataMap)
	outData.Prints = stockEntryPrint.Prints

	return nil
}

// 打印质检报告
func (e *PrintsJson) QualityReportPrint(d *dto.CommonPrintsReq, p *actions.DataPermission, s *QualityCheck, outData *dto.CommonPrintsResp) error {
	// 查询数据
	req := &dto.QualityCheckGetReq{
		Id: d.Id,
	}
	data := &dto.QualityCheckRes{}
	if err := s.Get(req, p, data); err != nil {
		e.Log.Errorf("PrintsJsonService StockEntryPrints error:%s \r\n", err)
		return err
	}
	bytes, _ := json.Marshal(data.QualityCheckDetail)
	fmt.Println(string(bytes))

	// 组合数据
	details := []map[string]any{}
	qualityByNames := map[string]bool{}
	for index, item := range data.QualityCheckDetail {
		// 统计质检详情
		detail := map[string]any{
			"id":                 index + 1,
			"qualityCheckOption": item.QualityCheckOption,
			"qualityRes":         models.QualityResName[item.QualityRes],
			"remark":             item.Remark,
			"qualityByName":      item.QualityByName,
		}
		details = append(details, detail)

		// 统计所有质检操作人员
		qualityByNames[item.QualityByName] = true
	}
	outDataMap := map[string]any{
		"sourceName":     data.SourceName,
		"sourceCode":     data.SourceCode,
		"entryCode":      data.EntryCode,
		"quantityNum":    data.QuantityNum,
		"qualityResName": data.QualityResName,
		"qualityTime":    data.QualityTime.Format("2006-01-02 15:04:05"),
		"qualityByName":  lo.Keys[string](qualityByNames),
		"details":        details,
	}

	// 拼接模板
	qualityCheckPrints := QualityCheckPrints{}
	qualityCheckPrints.GetPrints(outDataMap)
	outData.Prints = qualityCheckPrints.Prints
	return nil
}
