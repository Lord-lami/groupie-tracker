package render

import (
	"html/template"
	"strings"
)

func renderData(data any, templateName string) (template.HTML, error) {
	var dataHtml strings.Builder
	err := TheTemplates.ExecuteTemplate(&dataHtml,
		templateName,
		data)
	return template.HTML(dataHtml.String()), err
}
