package render

import (
	"html/template"
	"io/fs"
	"log"
	"os"
	"reflect"
	"runtime/debug"
)

type Ignored any
type DateString string
type ImageLinkString string
type LinkString string

type RenderFunc func(name string, data any) template.HTML

var templateFS fs.FS = os.DirFS("templates")

var TheTemplates *template.Template = template.Must(template.ParseFS(templateFS,
	"*.html",
	"*/*.html",
	"*/*/*.html"))

var RenderTypeFunc map[string]RenderFunc = map[string]RenderFunc{
	TypeString[Ignored]():         func(name string, data any) template.HTML { return "" },
	"int":                         RenderType[int]("int.html"),
	"string":                      RenderType[string]("string.html"),
	"bool":                        RenderType[bool]("bool.html"),
	TypeString[DateString]():      renderDateString,
	TypeString[LinkString]():      RenderType[LinkString]("linkstring.html"),
	TypeString[ImageLinkString](): RenderType[ImageLinkString]("imagelinkstring.html"),
}

func TypeString[T any]() string {
	return reflect.TypeFor[T]().String()
}

func RenderType[T any](templateName string) RenderFunc {
	return func(name string, data any) (dataHTML template.HTML) {
		value := data.(T)
		var templateData struct {
			Name  string
			Value T
		}
		templateData.Name = name
		templateData.Value = value
		var err error
		dataHTML, err = renderData(templateData, templateName)
		if err != nil {
			log.Println(err, string(debug.Stack()))
			return
		}
		return
	}
}

func selectRenderFuncFor(value reflect.Value) (renderFunc RenderFunc) {
	switch value.Kind() {
	case reflect.Struct:
		renderFunc = RenderObj
	case reflect.Array, reflect.Slice:
		renderFunc = RenderArr
	case reflect.Map:
		renderFunc = RenderMap
	default:
		renderFunc = RenderTypeFunc[value.Type().String()]
	}
	return
}
