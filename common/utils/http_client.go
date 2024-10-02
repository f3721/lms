package utils

import (
	"github.com/monaco-io/request"
	"github.com/monaco-io/request/response"
)

func PostJson(url string, post interface{}, header map[string]string) (resp *response.Sugar) {
	if header == nil {
		header = map[string]string{
			"Content-Type": "application/json;charset=utf-8",
		}
	}
	c := request.Client{
		URL:    url,
		Method: "POST",
		Header: header,
		JSON:   post,
	}
	resp = c.Send()
	return
}

func PostMultipartForm(url string, post, header map[string]string) (resp *response.Sugar) {
	c := request.Client{
		URL:    url,
		Method: "POST",
		Header: header,
		MultipartForm: request.MultipartForm{
			Fields: post,
		},
	}
	resp = c.Send()
	return
}
