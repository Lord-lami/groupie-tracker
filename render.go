package main

import (
	"html/template"
	"reflect"
	"strings"
)

func renderApiLink(name, url string) template.HTML {
	var ok bool
	url, ok = strings.CutPrefix(url, API)
	if !ok {
		panic("This is not an api link string of API: " + API)
	}
	var apiLinkHtml strings.Builder
	err := theTemplates.ExecuteTemplate(&apiLinkHtml,
		"apilinkstring.html",
		Link{name, url})
	if err != nil {
		panic(err)
	}
	return template.HTML(apiLinkHtml.String())
}

func renderImageLink(name, url string) template.HTML {
	var imageLinkHtml strings.Builder
	err := theTemplates.ExecuteTemplate(&imageLinkHtml,
		"imagelinkstring.html",
		Link{name, url})
	if err != nil {
		panic(err)
	}
	return template.HTML(imageLinkHtml.String())
}

func renderString(name, value string) template.HTML {
	var stringHtml strings.Builder
	err := theTemplates.ExecuteTemplate(&stringHtml,
		"imagelinkstring.html",
		String{name, value})
	if err != nil {
		panic(err)
	}
	return template.HTML(stringHtml.String())
}


func renderApiObject(objName string, data any) string {
	apiObjVal := reflect.ValueOf(data).Elem()
	elements := []template.HTML{}
	for i := range apiObjVal.NumField() {
		switch apiObjVal.Field(i).Kind() {
		case reflect.TypeFor[apiLinkString]().Kind():
			name := apiObjVal.Type().Field(i).Name
			url := apiObjVal.Field(i).String()
			linkHTML := renderApiLink(name, url)
			elements = append(elements, linkHTML)
		case reflect.TypeFor[imageLinkString]().Kind():
			name := apiObjVal.Type().Field(i).Name
			url := apiObjVal.Field(i).String()
			linkHTML := renderImageLink(name, url)
			elements = append(elements, linkHTML)
		case reflect.String:
			name := apiObjVal.Type().Field(i).Name
			value := apiObjVal.Field(i).String()
			stringHTML := renderString(name, value)
			elements = append(elements, stringHTML)
		}
	}
	// Render the html links as part of the index page object div
	var apiObjHTML strings.Builder
	err := theTemplates.ExecuteTemplate(&apiObjHTML,
		"object.html",
		object{objName, elements})
	if err != nil {
		panic(err)
	}
	return apiObjHTML.String()
}
