package excel

import (
	"fmt"
	"math/rand"
	"net/url"
	"reflect"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

var (
	defaultSheetName = "Sheet1" //默认Sheet名称
	defaultHeight    = 25.0     //默认行高度
)

type lzExcel struct {
	file      *excelize.File
	sheetName string //可定义默认sheet名称
}

func NewExcel() *lzExcel {
	return &lzExcel{file: createFile(), sheetName: defaultSheetName}
}

// 导出基本的表格
func (l *lzExcel) ExportToPath(params []map[string]string, data []map[string]interface{}, path string) (string, error) {
	l.export(params, data)
	name := createFileName()
	filePath := path + "/" + name
	err := l.file.SaveAs(filePath)
	return filePath, err
}

// ExportToWeb 导出到浏览器。此处使用的gin框架 其他框架可自行修改ctx
func (l *lzExcel) ExportToWeb(params []map[string]string, data []map[string]interface{}, c *gin.Context) {
	l.export(params, data)
	buffer, _ := l.file.WriteToBuffer()
	//设置文件类型
	c.Header("Content-Type", "application/vnd.ms-excel;charset=utf8")
	//设置文件名称
	c.Header("Content-Disposition", "attachment; filename="+url.QueryEscape(createFileName()))
	_, _ = c.Writer.Write(buffer.Bytes())
}

// 设置首行
func (l *lzExcel) writeTop(params []map[string]string) {
	style := excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	}
	topStyle, _ := l.file.NewStyle(&style)
	var word = 'A'
	//首行写入
	for _, conf := range params {
		title := conf["title"]
		width, _ := strconv.ParseFloat(conf["width"], 30)
		line := fmt.Sprintf("%c1", word)
		//设置标题
		_ = l.file.SetCellValue(l.sheetName, line, title)
		//列宽
		_ = l.file.SetColWidth(l.sheetName, fmt.Sprintf("%c", word), fmt.Sprintf("%c", word), width)
		//设置样式
		_ = l.file.SetCellStyle(l.sheetName, line, line, topStyle)
		word++
	}
}

// 写入数据
func (l *lzExcel) writeData(params []map[string]string, data []map[string]interface{}) {
	style := excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	}
	lineStyle, _ := l.file.NewStyle(&style)
	//数据写入
	var j = 2 //数据开始行数
	for i, val := range data {
		//设置行高
		_ = l.file.SetRowHeight(l.sheetName, i+1, defaultHeight)
		//逐列写入
		var word = 'A'
		for _, conf := range params {
			valKey := conf["key"]
			line := fmt.Sprintf("%c%v", word, j)
			isNum := conf["is_num"]

			//设置值
			if isNum != "0" {
				valNum := fmt.Sprintf("'%v", val[valKey])
				_ = l.file.SetCellValue(l.sheetName, line, valNum)
			} else {
				_ = l.file.SetCellValue(l.sheetName, line, val[valKey])
			}

			//设置样式
			_ = l.file.SetCellStyle(l.sheetName, line, line, lineStyle)
			word++
		}
		j++
	}
	//设置行高 尾行
	_ = l.file.SetRowHeight(l.sheetName, len(data)+1, defaultHeight)
}

func (l *lzExcel) export(params []map[string]string, data []map[string]interface{}) {
	l.writeTop(params)
	l.writeData(params, data)
}

func createFile() *excelize.File {
	f := excelize.NewFile()
	// 创建一个默认工作表
	sheetName := defaultSheetName
	index, _ := f.NewSheet(sheetName)
	// 设置工作簿的默认工作表
	f.SetActiveSheet(index)
	return f
}

func createFileName() string {
	name := time.Now().Format("2006-01-02-15-04-05")
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("excle-%v-%v.xlsx", name, rand.Int63n(time.Now().Unix()))
}

// ExportExcelByStruct excel导出(数据源为Struct) []interface{}
func (l *lzExcel) ExportExcelByStruct(c *gin.Context, titleList []map[string]string, data []interface{}, fileName string, sheetName string) error {
	l.file.SetSheetName("Sheet1", sheetName)
	title := make([]string, 0)
	fields := make([]string, 0)
	for _, m := range titleList {
		for k, v := range m {
			title = append(title, v)
			fields = append(fields, k)
		}
	}
	style := excelize.Style{
		//Font: &excelize.Font{
		//	Family: "arial",
		//	Size:   13,
		//	Color:  "#666666",
		//},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	}
	rowStyleID, _ := l.file.NewStyle(&style)
	_ = l.file.SetSheetRow(sheetName, "A1", &title)
	_ = l.file.SetSheetRow(sheetName, "A2", &fields)
	_ = l.file.SetRowHeight("Sheet1", 1, 30)
	_ = l.file.SetRowHeight("Sheet1", 2, 30)
	length := len(titleList)
	headStyle := Letter(length)
	var lastRow string
	var widthRow string
	for k, v := range headStyle {
		if k == length-1 {
			lastRow = fmt.Sprintf("%s1", v)
			widthRow = v
		}
	}
	_ = l.file.SetColWidth(sheetName, "A", widthRow, 15)
	_ = l.file.SetCellStyle(sheetName, "A1", fmt.Sprintf("%s1", widthRow), rowStyleID)
	_ = l.file.SetCellStyle(sheetName, "A2", fmt.Sprintf("%s2", widthRow), rowStyleID)
	rowNum := 2
	for _, v := range data {
		t := reflect.TypeOf(v)
		value := reflect.ValueOf(v)
		row := make([]interface{}, 0)
		for l := 0; l < t.NumField(); l++ {

			val := value.Field(l).Interface()
			row = append(row, val)
		}
		rowNum++
		err := l.file.SetSheetRow(sheetName, "A"+strconv.Itoa(rowNum), &row)
		_ = l.file.SetCellStyle(sheetName, fmt.Sprintf("A%d", rowNum), fmt.Sprintf("%s", lastRow), rowStyleID)
		if err != nil {
			return err
		}
	}
	disposition := fmt.Sprintf("attachment; filename=%s-%s.xlsx", url.QueryEscape(fileName), time.Now().Format("20060102150405"))
	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	c.Writer.Header().Set("Content-Disposition", disposition)
	c.Writer.Header().Set("Content-Transfer-Encoding", "binary")
	c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
	return l.file.Write(c.Writer)
}

// ExportExcelByMap 导出excel 数据源为[]map
func (l *lzExcel) ExportExcelByMap(c *gin.Context, titleList []map[string]string, data []map[string]interface{}, fileName, sheetName string) error {
	l.file.SetSheetName("Sheet1", sheetName)
	title := make([]string, 0)
	fields := make([]string, 0)
	for _, m := range titleList {
		for k, v := range m {
			title = append(title, v)
			fields = append(fields, k)
		}
	}
	style := excelize.Style{
		//Font: &excelize.Font{
		//	Family: "arial",
		//	Size:   13,
		//	Color:  "#666666",
		//},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	}
	rowStyleID, _ := l.file.NewStyle(&style)
	_ = l.file.SetSheetRow(sheetName, "A1", &title)
	_ = l.file.SetSheetRow(sheetName, "A2", &fields)
	_ = l.file.SetRowHeight("Sheet1", 1, 30)
	_ = l.file.SetRowHeight("Sheet1", 2, 30)
	length := len(titleList)
	headStyle := Letter(length)
	var lastRow string
	var widthRow string
	for k, v := range headStyle {
		if k == length-1 {
			lastRow = fmt.Sprintf("%s1", v)
			widthRow = v
		}
	}
	_ = l.file.SetColWidth(sheetName, "A", widthRow, 15)
	_ = l.file.SetCellStyle(sheetName, "A1", fmt.Sprintf("%s1", widthRow), rowStyleID)
	_ = l.file.SetCellStyle(sheetName, "A2", fmt.Sprintf("%s2", widthRow), rowStyleID)
	rowNum := 2
	for _, value := range data {
		row := make([]interface{}, 0)
		var dataSlice []string
		for key := range value {
			dataSlice = append(dataSlice, key)
		}
		//sort.Strings(dataSlice)
		for _, v := range fields {
			if val, ok := value[v]; ok {
				row = append(row, val)
			}
		}
		rowNum++
		if err := l.file.SetSheetRow(sheetName, fmt.Sprintf("A%d", rowNum), &row); err != nil {
			return err
		}
		if err := l.file.SetCellStyle(sheetName, fmt.Sprintf("A%d", rowNum), fmt.Sprintf("%s", lastRow), rowStyleID); err != nil {
			return err
		}
	}
	disposition := fmt.Sprintf("attachment; filename=%s_%s.xlsx", url.QueryEscape(fileName), time.Now().Format("20060102150405"))
	//c.Writer.Header().Set("Content-Type", "application/octet-stream")
	c.Writer.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Writer.Header().Set("Content-Disposition", disposition)
	//c.Writer.Header().Set("Content-Transfer-Encoding", "binary")
	c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")

	return l.file.Write(c.Writer)
}

// Letter 遍历a-z
func Letter(length int) []string {
	var str []string
	for i := 0; i < length; i++ {
		if i >= 26 {
			str = append(str, "A"+string(rune('A'+(i-26))))
		} else {
			str = append(str, string(rune('A'+i)))
		}
	}
	return str
}

// MergeErrMsgColumn 将错误信息和data合并，返回新的titleList和data
func (l *lzExcel) MergeErrMsgColumn(titleList []map[string]string, data []map[string]interface{}, errMsg map[int]string) ([]map[string]string, []map[string]interface{}) {
	titleList = append(titleList, map[string]string{"errMsg": "错误提示"})
	for i, v := range errMsg {
		data[i]["errMsg"] = v
	}
	return titleList, data
}
