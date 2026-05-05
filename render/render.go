package render

import (
	"embed"
	"html/template"
	"log"
	"reflect"
	"runtime/debug"
)

type Ignored any
type DateString string
type ImageLinkString string
type LinkString string

type RenderFunc func(name string, data any) template.HTML

//go:embed templates/*
var templateFS embed.FS

var TheTemplates *template.Template = template.Must(template.ParseFS(templateFS,
	"templates/*.html",
	"templates/*/*.html",
	"templates/*/*/*.html"))

var RenderTypeFunc map[string]RenderFunc = map[string]RenderFunc{
	TypeString[Ignored]():         func(name string, data any) template.HTML { return "" },
	"int":                         NewRenderFunc[int]("int.html"),
	"string":                      NewRenderFunc[string]("string.html"),
	"bool":                        NewRenderFunc[bool]("bool.html"),
	TypeString[DateString]():      renderDateString,
	TypeString[LinkString]():      NewRenderFunc[LinkString]("linkstring.html"),
	TypeString[ImageLinkString](): NewRenderFunc[ImageLinkString]("imagelinkstring.html"),
}

func TypeString[T any]() string {
	return reflect.TypeFor[T]().String()
}

func MapTypeToRenderFunc[T any](f RenderFunc) {
	RenderTypeFunc[TypeString[T]()] = f
}

func NewRenderFunc[T any](templateName string) RenderFunc {
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
