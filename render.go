package main

import (
	"html/template"
	"log"
	"reflect"
	"runtime/debug"
	"slices"
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

func RenderType[T any](templateName string) func(name string, data any) (dataHTML template.HTML) {
	return func(name string, data any) (dataHTML template.HTML) {
		dataVal := reflect.ValueOf(data)
		if dataVal.Comparable() && data == reflect.Zero(dataVal.Type()).Interface() {
			return ""
		}
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

func renderApiLink(name string, data any) (linkHTML template.HTML) {
	url := string(data.(apiLinkString))
	if url == "" {
		return ""
	}
	var ok bool
	url, ok = strings.CutPrefix(url, API)
	if !ok {
		log.Println("This is not an api link string of API: "+API, string(debug.Stack()))
		return
	}
	linkHTML = RenderType[string]("apilinkstring.html")(name, url)
	return
}

func renderDateString(name string, data any) (dateStringHTML template.HTML) {
	raw := string(data.(dateString))
	if raw == "" {
		return ""
	}
	raw = strings.ReplaceAll(raw, "*", "")
	date, err := time.Parse("02-01-2006", raw)
	if err != nil {
		log.Println(err, string(debug.Stack()))
		return
	}
	raw = date.Format(time.DateOnly)
	display := date.Format("Mon, 02 Jan 2006")
	type DateData struct {
		Raw     string
		Display string
	}
	dateStringHTML = RenderType[DateData]("datestring.html")(name, DateData{raw, display})
	return
}

type RenderFunc func(name string, data any) template.HTML

var RenderTypeFunc map[string]RenderFunc = map[string]RenderFunc{
	reflect.TypeFor[ignored]().String():    func(name string, data any) template.HTML { return "" },
	"int":                                  RenderType[int]("int.html"),
	"string":                               RenderType[string]("string.html"),
	"bool":                                 RenderType[bool]("bool.html"),
	reflect.TypeFor[dateString]().String(): renderDateString,
	reflect.TypeFor[apiLinkString]().String():   renderApiLink,
	reflect.TypeFor[imageLinkString]().String(): RenderType[imageLinkString]("imagelinkstring.html"),
}

func selectRenderFuncFor(value reflect.Value) (renderFunc RenderFunc) {
	switch value.Kind() {
	case reflect.Struct:
		renderFunc = renderObj
	case reflect.Array, reflect.Slice:
		renderFunc = renderArr
	default:
		renderFunc = RenderTypeFunc[value.Type().String()]
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
		if objVal.Field(i).Type().String() == reflect.TypeFor[ignored]().String() {
			continue
		}
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			renderFunc := selectRenderFuncFor(objVal.Field(i))
			name := objVal.Type().Field(i).Name
			data := objVal.Field(i).Interface()
			if renderFunc == nil {
				log.Println(name)
			}
			element := renderFunc(name, data)
			elements[i] = element
		}(i)
	}
	wg.Wait()

	if slices.Equal(elements, make([]template.HTML, objVal.NumField())) {
		return ""
	}

	// Render the html links as part of the object div
	objHTML = RenderType[[]template.HTML]("object.html")(objName, elements)
	return
}

func renderArr(arrName string, data any) (arrHTML template.HTML) {
	if reflect.TypeOf(data).Kind() != reflect.Array &&
		reflect.TypeOf(data).Kind() != reflect.Slice {
		log.Println("array data was not passed as a slice or array", string(debug.Stack()))
		return ""
	}

	arrVal := reflect.ValueOf(data)
	if arrVal.Len() == 0 {
		return ""
	}
	renderFunc := selectRenderFuncFor(arrVal.Index(0))
	var wg sync.WaitGroup
	elements := make([]template.HTML, arrVal.Len())

	for i := range arrVal.Len() {
		if arrVal.Index(i).Type().String() == reflect.TypeFor[ignored]().String() {
			continue
		}
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			index := strconv.Itoa(i)
			element := renderFunc(index, arrVal.Index(i).Interface())
			elements[i] = element
		}(i)
	}
	wg.Wait()

	

	// Render the html as part of the array list
	arrHTML = RenderType[[]template.HTML]("array.html")(arrName, elements)
	return
}
