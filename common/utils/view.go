package utils

import (
	"bytes"
	"errors"
	"go-admin/config"
	"html/template"
	"io"
	"net/http"
	"strings"
	"time"
)

// 尽管*Template.Execute是协程安全的，但考虑到项目中View方法不是高频调用，故这里不对*Template做全局处理（减少静态内存占用）

func View(filePath string, model interface{}) (string, error) {
	// // 相对路径
	// var out bytes.Buffer
	// fh, err := os.Open(filePath)
	// if err != nil {
	// 	return "", err
	// }
	// defer fh.Close()
	// bytes, err := io.ReadAll(fh)
	// if err != nil {
	// 	return "", err
	// }

	// 网络路径
	var out bytes.Buffer
	url := config.ExtConfig.ApiHost + "/" + filePath
	response, err := http.Get(url)
	if err != nil {
		return "", errors.New("读取模板地址失败")
	}
	defer response.Body.Close()

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", errors.New("读取模板内容失败")
	}

	// 模板渲染
	tpl := template.New("").Funcs(TPLFuncMap())
	tpl, err = tpl.Parse(string(bytes))
	if err != nil {
		return "", err
	}
	if err = tpl.Execute(&out, model); err != nil {
		return "", err
	}
	return out.String(), nil
}

func TPLFuncMap() template.FuncMap {
	funcMap := template.FuncMap{}

	funcMap["upper"] = strings.ToUpper
	funcMap["lower"] = strings.ToLower
	funcMap["htmlUnescaped"] = func(x string) interface{} {
		return template.HTML(x)
	}
	funcMap["sub"] = func(x, y int) int {
		return x - y
	}
	funcMap["nowDateTime"] = func() string {
		return TimeFormat(time.Now())
	}
	funcMap["addOne"] = func(x int) int {
		return x + 1
	}
	return funcMap
}
