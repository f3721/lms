package excel

import (
	"github.com/xuri/excelize/v2"
	"mime/multipart"
	"net/http"
)

// GetExcelData 获取excel数据并返回 []map[string]string{}
func (l *lzExcel) GetExcelData(file *multipart.FileHeader) (err error, data []map[string]interface{}, titleList []map[string]string) {
	// 打开 Excel 文件
	fh, err := file.Open()
	if err != nil {
		return
	}
	defer fh.Close()

	f, err := excelize.OpenReader(fh)
	if err != nil {
		return
	}

	rows, err := f.Rows(l.sheetName) // 读取 sheetName 工作表的所有行
	if err != nil {
		return
	}
	defer rows.Close()

	// 读取第 1 行，title
	rows.Next()
	titles, err := rows.Columns()
	if err != nil {
		return
	}
	// 读取第 2 行，即字段信息，作为map中的key
	rows.Next()
	fields, err := rows.Columns()
	if err != nil {
		return
	}
	for i, field := range fields {
		titleList = append(titleList, map[string]string{field: titles[i]}) // 使用keyRow行作为key，对应单元格数据作为value
	}

	// 遍历keyRow下所有行，将每行中的单元格数据保存为map形式

	for rows.Next() {
		values, err := rows.Columns()
		if err != nil {
			return err, nil, nil
		}

		row := map[string]interface{}{} // 每行的数据保存在一个map中
		rowLen := len(values)
		for i, field := range fields {
			if rowLen-1 < i {
				row[field] = ""
			} else {
				row[field] = values[i]
			}
		}

		//for i, value := range values {
		//	row[fields[i]] = value  // 使用keyRow行作为key，对应单元格数据作为value
		//}
		data = append(data, row) // 将该行数据存储到数组中
	}
	return
}

func (l *lzExcel) ValidImportFieldsCorrect(importFile *multipart.FileHeader, tmpFilePath string) (res bool) {
	// 读取Excel文件内容
	importF, err := importFile.Open()
	if err != nil {
		return
	}
	defer importF.Close()

	f, err := excelize.OpenReader(importF)
	if err != nil {
		return
	}
	importRows, err := f.Rows(l.sheetName) // 读取 sheetName 工作表的所有行
	if err != nil {
		return
	}
	defer importRows.Close()
	// 读取第 2 行，即字段信息
	importRows.Next()
	importRows.Next()
	importFields, err := importRows.Columns()
	if err != nil {
		return
	}

	// 下载文件
	resp, err := http.Get(tmpFilePath)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	f2, err := excelize.OpenReader(resp.Body)
	if err != nil {
		return
	}

	// 读取根目录文件
	//f2, err := excelize.OpenFile(tmpFilePath)
	//if err != nil {
	//	return
	//}
	tmpRows, err := f2.Rows(l.sheetName) // 读取 sheetName 工作表的所有行
	if err != nil {
		return
	}
	defer tmpRows.Close()
	// 读取第 2 行，即字段信息
	tmpRows.Next()
	tmpRows.Next()
	tmpFields, err := tmpRows.Columns()
	if err != nil {
		return
	}
	if len(importFields) != len(tmpFields) {
		return
	}

	for i, _ := range tmpFields {
		if tmpFields[i] != importFields[i] {
			return
		}
	}

	return true
}
