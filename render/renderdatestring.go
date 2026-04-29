package render

import (
	"html/template"
	"log"
	"runtime/debug"
	"strings"
	"time"
)

func renderDateString(name string, data any) (dateStringHTML template.HTML) {
	raw := string(data.(DateString))
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
