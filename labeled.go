package main

import (
	"html/template"
	"unicode"

	"github.com/Lord-lami/render-html"
)

type labeledMap map[string][]render.DateString

type labeledStringSlice []string
type labeledInt int
type labeledDateString render.DateString

// The renderLabeled functions render a variable and add a 
// strong tag element containing the name before the variable's HTML

func renderLabeledMap(name string, data any) template.HTML {
	value := map[string][]render.DateString(data.(labeledMap))
	dataHTML := render.Render(name, value)
	labelHTML := template.HTML("<strong>" + spaceOutCamelCase(name) + ": </strong>")
	return labelHTML + dataHTML
}

func renderLabeledStringSlice(name string, data any) template.HTML {
	value := []string(data.(labeledStringSlice))
	dataHTML := render.Render(name, value)
	labelHTML := template.HTML("<strong>" + spaceOutCamelCase(name) + ":</strong>")
	return labelHTML + dataHTML
}

func renderLabeledInt(name string, data any) template.HTML {
	value := int(data.(labeledInt))
	dataHTML := render.Render(name, value)
	labelHTML := template.HTML("<strong>" + spaceOutCamelCase(name) + ": </strong>")
	return labelHTML + dataHTML
}

func renderLabeledDateString(name string, data any) template.HTML {
	value := render.DateString(data.(labeledDateString))
	dataHTML := render.Render(name, value)
	labelHTML := template.HTML("<strong>" + spaceOutCamelCase(name) + ": </strong>")
	return labelHTML + dataHTML
}

// spaceOutCamelCase adds spaces before the uppercase letters in a string.
func spaceOutCamelCase(str string) string {
	result := []rune{rune(str[0])}
	for _, char := range str[1:] {
		if unicode.IsUpper(char) {
			result = append(result, ' ', char)
		} else {
			result = append(result, char)
		}
	}
	return string(result)
}
