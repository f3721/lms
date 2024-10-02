package client

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"github.com/monaco-io/request"
	"github.com/monaco-io/request/response"
)

type HS map[string]string

type BaseClient struct {
	Context *gin.Context
}

func Get(ctx *gin.Context, url string, query HS) (resp *response.Sugar) {
	c := request.Client{
		URL:    url,
		Method: request.GET,
		Header: HS{
			//"Authorization": ctx.GetHeader("Authorization"),
			"tenant-id":    ctx.GetHeader("tenant-id"),
			"X-Request-Id": pkg.GenerateMsgIDFromContext(ctx),
		},
		Query: query,
	}
	resp = c.Send()
	return
}

func Post(ctx *gin.Context, url string, post HS, files []string) (resp *response.Sugar) {
	c := request.Client{
		URL:    url,
		Method: "POST",
		Header: HS{
			//"Authorization": ctx.GetHeader("Authorization"),
			"tenant-id":    ctx.GetHeader("tenant-id"),
			"X-Request-Id": pkg.GenerateMsgIDFromContext(ctx),
		},
		MultipartForm: request.MultipartForm{
			Fields: post,
			Files:  files,
		},
	}
	resp = c.Send()
	return
}

func PostJson(ctx *gin.Context, url string, post interface{}) (resp *response.Sugar) {
	c := request.Client{
		URL:    url,
		Method: "POST",
		Header: HS{
			//"Authorization": ctx.GetHeader("Authorization"),
			"tenant-id":    ctx.GetHeader("tenant-id"),
			"X-Request-Id": pkg.GenerateMsgIDFromContext(ctx),
		},
		JSON: post,
	}
	resp = c.Send()
	return
}

func PostCustomHeader(ctx *gin.Context, url string, post HS, header HS) (resp *response.Sugar) {
	c := request.Client{
		URL:    url,
		Method: "POST",
		Header: header,
		MultipartForm: request.MultipartForm{
			Fields: post,
		},
	}
	fmt.Println("GetShortLink[3]:", post)
	fmt.Println("GetShortLink[4]:", url)
	resp = c.Send()
	return
}

func Put(ctx *gin.Context, url string, put interface{}) (resp *response.Sugar) {
	c := request.Client{
		URL:    url,
		Method: request.PUT,
		Header: HS{
			//"Authorization": ctx.GetHeader("Authorization"),
			"tenant-id":    ctx.GetHeader("tenant-id"),
			"X-Request-Id": pkg.GenerateMsgIDFromContext(ctx),
		},
		JSON: put,
	}
	resp = c.Send()
	return
}
