// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// output 对解析后的数据进行渲染输出。
//
// 目前支持以下三种渲染方式：
//  - html: 以 html 格式输出文本，模板可自定义；
//  - html+: html 的调试模式，程序不会输出任何，而是在浏览器中展示相关页面；
//  - json: 以 JSON 格式输出内容。
package output

import (
	"time"

	"github.com/caixw/apidoc/app"
	"github.com/caixw/apidoc/doc"
	"github.com/issue9/utils"
)

// 支持的渲染方式
var renderTypes = []string{
	"html",
	"html+",
	"json",
}

// 渲染输出的相关设置项。
type Options struct {
	Dir      string        `json:"dir"`                // 文档的保存目录
	Type     string        `json:"type"`               // 渲染方式，默认为 html
	Template string        `json:"template,omitempty"` // 指定一个输出模板
	Port     string        `json:"port,omitempty"`     // 调试的端口
	Elapsed  time.Duration `json:"-"`                  // 编译用时
}

// 对 Options 作一些初始化操作。
func (o *Options) Init() *app.OptionsError {
	if len(o.Dir) == 0 {
		return &app.OptionsError{Field: "Dir", Message: "不能为空"}
	}

	if !isSuppertedType(o.Type) {
		return &app.OptionsError{Field: "Type", Message: "不支持的类型"}
	}

	// 只有 html 和 html+ 才需要判断模板文件是否存在
	if o.Type == "html" || o.Type == "html+" {
		if len(o.Template) > 0 && !utils.FileExists(o.Template) {
			return &app.OptionsError{Field: "Template", Message: "目录不存在"}
		}
	}

	// 调试模式，必须得有端口
	if o.Type == "html+" {
		if len(o.Port) == 0 {
			return &app.OptionsError{Field: "Port", Message: "不能为空"}
		}

		if o.Port[0] != ':' {
			o.Port = ":" + o.Port
		}
	}

	return nil
}

// Render 渲染 docs 的内容，具体的渲染参数由 o 指定。
func Render(docs *doc.Doc, o *Options) error {
	switch o.Type {
	case "html":
		return renderHTML(docs, o)
	case "html+":
		return renderHTMLPlus(docs, o)
	case "json":
		return renderJSON(docs, o)
	default:
		return &app.OptionsError{Field: "Type", Message: "不支持该类型"}
	}
}

func isSuppertedType(typ string) bool {
	for _, k := range renderTypes {
		if k == typ {
			return true
		}
	}

	return false
}