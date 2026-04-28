package main

import (
	"html/template"
	"log"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"
)

func renderData(data any, templateName string) (template.HTML, error) {
	var dataHtml strings.Builder
	err := theTemplates.ExecuteTemplate(&dataHtml,
		templateName,
		data)
	return template.HTML(dataHtml.String()), err
}

func renderType[T any](templateName string) func(name string, data any) (dataHTML template.HTML) {
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
			log.Println(err, debug.Stack())
			return
		}
		return
	}
}

func renderApiLink(name string, data any) (linkHTML template.HTML) {
	url := string(data.(apiLinkString))
	var ok bool
	url, ok = strings.CutPrefix(url, API)
	if !ok {
		log.Println("This is not an api link string of API: "+API, debug.Stack())
		return
	}
	data = url
	linkHTML = renderType[apiLinkString]("apilinkstring.html")(name, data)
	return
}

func renderDateString(name string, data any) (dateStringHTML template.HTML) {
	raw := string(data.(dateString))
	raw = strings.ReplaceAll(raw, "*", "")
	date, err := time.Parse("02-01-2006", raw)
	if err != nil {
		log.Println(err, debug.Stack())
		return
	}
	raw = date.Format(time.DateOnly)
	display := date.Format("Mon, 02 Jan 2006")
	type DateData struct {
		Raw     string
		Display string
	}
	dateStringHTML = renderType[DateData]("datestring.html")(name, DateData{raw, display})
	return
}

type RenderFunc func(name string, data any) template.HTML

var RenderStaticType map[string]RenderFunc = map[string]RenderFunc{
	"int":                                  renderType[int]("int.html"),
	"string":                               renderType[string]("string.html"),
	"bool":                                 renderType[bool]("bool.html"),
	reflect.TypeFor[dateString]().String(): renderDateString,
	reflect.TypeFor[apiLinkString]().String():   renderApiLink,
	reflect.TypeFor[imageLinkString]().String(): renderType[imageLinkString]("imagelinkstring.html"),
}

func selectRenderFuncFor(value reflect.Value) (renderFunc RenderFunc) {
	switch value.Kind() {
	case reflect.Struct:
		renderFunc = renderObj
	case reflect.Array, reflect.Slice:
		renderFunc = renderArr
	default:
		renderFunc = RenderStaticType[value.Type().String()]
	}
	return
}

func renderObj(objName string, data any) (objHTML template.HTML) {
	objVal := reflect.ValueOf(data)
	if objVal.Type().Kind() != reflect.Struct {
		log.Println("The object data was not passed as a struct")
		return ""
	}

	var wg sync.WaitGroup
	elements := make([]template.HTML, objVal.NumField())
	for i := range objVal.NumField() {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			renderFunc := selectRenderFuncFor(objVal.Field(i))
			name := objVal.Type().Field(i).Name
			// fmt.Println("The fieild is: " + name)
			data := objVal.Field(i).Interface()
			element := renderFunc(name, data)
			elements[i] = element
		}(i)
	}
	wg.Wait()

	// Render the html links as part of the object div
	objHTML = renderType[[]template.HTML]("object.html")(objName, elements)
	return
}

func renderArr(arrName string, data any) (arrHTML template.HTML) {
	if reflect.TypeOf(data).Kind() != reflect.Array &&
		reflect.TypeOf(data).Kind() != reflect.Slice {
		log.Panicln("array data was not passed as a slice or array")
		return ""
	}

	arrVal := reflect.ValueOf(data)
	renderFunc := selectRenderFuncFor(arrVal.Index(0))
	var wg sync.WaitGroup
	elements := make([]template.HTML, arrVal.Len())

	for i := range arrVal.Len() {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			index := strconv.Itoa(i)
			element := renderFunc(index, arrVal.Index(i).Interface())
			elements[i] = element
		}(i)
	}
	wg.Wait()

	// Render the html links as part of the array list
	arrHTML = renderType[[]template.HTML]("array.html")(arrName, elements)
	return

}
