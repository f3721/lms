package apis

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/utils"
	"strings"
	"text/template"

	"go-admin/app/other/apis/tools"
	adminTools "go-admin/app/other/models/tools"
)

type Gen struct {
	api.Api
}

func (e Gen) CustomGen(c *gin.Context) {
	e.Context = c
	var errors []error
	defer func() {
		if len(errors) > 0 {
			var errMsgs []string
			for _, err := range errors {
				errMsgs = append(errMsgs, err.Error())
			}
			e.Custom(gin.H{"status": 500, "msg": "gen code failed!!", "errors": errMsgs})
		}
	}()

	result, err := tools.SysTable{}.GetByTableNameCustom(c)
	if err != nil || result.TBName == "" {
		errors = append(errors, err)
		err = tools.SysTable{}.InsertCustom(c)
		if err != nil {
			errors = append(errors, err)
		}
		result, err = tools.SysTable{}.GetByTableNameCustom(c)
		if err != nil {
			errors = append(errors, err)
		}
	}
	log := e.GetLogger()
	table := adminTools.SysTables{}
	id := result.TableId

	db, err := pkg.GetOrm(c)

	if err != nil {
		log.Errorf("get db connection error, %s", err.Error())
		errors = append(errors, fmt.Errorf("数据库链接获取失败！错误详情：%s", err.Error()))
		return
	}

	if result.TBName != c.Query("tables") {
		errors = append(errors, fmt.Errorf("未找到记录，请检查配置文件中gen: dbname配置是否正确"))
		return
	}

	table.TableId = id
	tab, _ := table.Get(db, false)
	tab.MLTBName = strings.Replace(tab.TBName, "_", "-", -1)

	basePath := "template/v5/"
	t1, err := template.ParseFiles(basePath + "model.go.template")
	if err != nil {
		basePath = "../" + basePath
		t1, err = template.ParseFiles(basePath + "model.go.template")
	}
	if err != nil {
		log.Error(err)
		errors = append(errors, fmt.Errorf("model模版读取失败！错误详情：%s", err.Error()))
		return
	}
	t2, err := template.ParseFiles(basePath + "no_actions/apis.go.template")
	if err != nil {
		log.Error(err)
		errors = append(errors, fmt.Errorf("api模版读取失败！错误详情：%s", err.Error()))
		return
	}

	routerFile := basePath + "no_actions/router_check_role.go.template"
	t3, err := template.ParseFiles(routerFile)
	if err != nil {
		log.Error(err)
		errors = append(errors, fmt.Errorf("路由模版失败！错误详情：%s", err.Error()))
		return
	}

	t6, err := template.ParseFiles(basePath + "dto.go.template")
	if err != nil {
		log.Error(err)
		errors = append(errors, fmt.Errorf("dto模版解析失败失败！错误详情：%s", err.Error()))
		return
	}
	t7, err := template.ParseFiles(basePath + "no_actions/service.go.template")
	if err != nil {
		log.Error(err)
		errors = append(errors, fmt.Errorf("service模版失败！错误详情：%s", err.Error()))
		return
	}
	outBaseDir := c.DefaultQuery("out_dir", "./")
	if outBaseDir == "default" {
		outBaseDir = "./"
	}
	if !strings.HasSuffix(outBaseDir, "/") {
		outBaseDir += "/"
	}
	_ = pkg.PathCreate(outBaseDir + "app/" + tab.PackageName + "/apis/admin/")
	_ = pkg.PathCreate(outBaseDir + "app/" + tab.PackageName + "/models/")
	_ = pkg.PathCreate(outBaseDir + "app/" + tab.PackageName + "/router/")
	_ = pkg.PathCreate(outBaseDir + "app/" + tab.PackageName + "/service/admin/dto/")

	var b1 bytes.Buffer
	err = t1.Execute(&b1, tab)
	var b2 bytes.Buffer
	err = t2.Execute(&b2, tab)
	var b3 bytes.Buffer
	err = t3.Execute(&b3, tab)
	var b6 bytes.Buffer
	err = t6.Execute(&b6, tab)
	var b7 bytes.Buffer
	err = t7.Execute(&b7, tab)

	apiPath := outBaseDir + "app/" + tab.PackageName + "/apis/admin/" + tab.TBName + ".go"
	if fileExists := utils.CheckExist(apiPath); !fileExists {
		errors = append(errors, fmt.Errorf("api文件已存在：%s", apiPath))
	} else {
		pkg.FileCreate(b2, apiPath)
	}

	modelPath := outBaseDir + "app/" + tab.PackageName + "/models/" + tab.TBName + ".go"
	if fileExists := utils.CheckExist(modelPath); !fileExists {
		errors = append(errors, fmt.Errorf("model文件已存在：%s", modelPath))
	} else {
		pkg.FileCreate(b1, modelPath)
	}

	routerPath := outBaseDir + "app/" + tab.PackageName + "/router/" + tab.TBName + ".go"
	if fileExists := utils.CheckExist(routerPath); !fileExists {
		errors = append(errors, fmt.Errorf("router文件已存在：%s", routerPath))
	} else {
		pkg.FileCreate(b3, routerPath)
	}

	dtoPath := outBaseDir + "app/" + tab.PackageName + "/service/admin/dto/" + tab.TBName + ".go"
	if fileExists := utils.CheckExist(dtoPath); !fileExists {
		errors = append(errors, fmt.Errorf("dto文件已存在：%s", dtoPath))
	} else {
		pkg.FileCreate(b6, dtoPath)
	}

	servicePath := outBaseDir + "app/" + tab.PackageName + "/service/admin/" + tab.TBName + ".go"
	if fileExists := utils.CheckExist(servicePath); !fileExists {
		errors = append(errors, fmt.Errorf("service文件已存在：%s", servicePath))
	} else {
		pkg.FileCreate(b7, servicePath)
	}
	if len(errors) > 0 {
		return
	}
	e.OK("", "Code generated successfully！")
}
