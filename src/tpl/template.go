package tpl

import (
	"github.com/P4elme6ka/go-http-media-server/src/tpl/util"
	"html/template"
)

func ParsePageTpl(tplText string) (tpl *template.Template, err error) {
	tpl = template.New("page")
	tpl = addFuncMap(tpl)
	tpl, err = tpl.Parse(tplText)

	return
}

func addFuncMap(tpl *template.Template) *template.Template {
	return tpl.Funcs(template.FuncMap{
		"fmtFilename": util.FormatFilename,
		"fmtSize":     util.FormatSize,
		"fmtTime":     util.FormatTime,
	})
}
