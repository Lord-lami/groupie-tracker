package main

import (
	"groupie-tracker/render"
	"html/template"
	"log"
	"runtime/debug"
	"strings"
)

type apiLinkString string

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
	linkHTML = render.RenderType[string]("apilinkstring.html")(name, url)
	return
}