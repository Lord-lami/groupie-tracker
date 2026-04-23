package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
)

var API string = "https://groupietrackers.herokuapp.com/api"

var templateFS fs.FS = os.DirFS("templates")

var theTemplates *template.Template = template.Must(template.ParseFS(templateFS,
	"*.html",
	"*/*.html",
	"*/*/*.html"))

var pageNum int = 0

var page struct {
	Title   string
	Content template.HTML
}

type Link struct {
	Name string
	Url  string
}

type object struct {
	Name     string
	Elements []template.HTML
}

type apiLinkString string
type imageLinkString string

type mainPageDataHolder struct {
	Artists   apiLinkString `json:"artists"`
	Locations apiLinkString `json:"locations"`
	Dates     apiLinkString `json:"dates"`
	Relation  apiLinkString `json:"relation"`
}

func groupieTrackerApiResponseBody(path string) (body []byte) {
	apiLink := API + path
	resp, err := http.Get(apiLink)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return
}

func hanleMainPage(w http.ResponseWriter, r *http.Request) {
	// Get the response body from the API
	mainPageData := mainPageDataHolder{}
	err := json.Unmarshal(groupieTrackerApiResponseBody(""), &mainPageData)
	if err != nil {
		panic(err)
	}

	// Render the links from the API as html anchor tags
	// that link to their display page
	mainPageValue := reflect.ValueOf(&mainPageData).Elem()
	links := []template.HTML{}
	for i := range mainPageValue.NumField() {
		var linkHtml strings.Builder
		url := mainPageValue.Field(i).String()
		name := mainPageValue.Type().Field(i).Name
		var ok bool
		url, ok = strings.CutPrefix(url, "https://groupietrackers.herokuapp.com/api")
		if !ok {
			panic("Something is horribly wrong with this link")
		}
		err := theTemplates.ExecuteTemplate(&linkHtml,
			"pageurlstring.html",
			Link{name, url})
		if err != nil {
			panic(err)
		}
		links = append(links, template.HTML(linkHtml.String()))
	}

	// slices.Sort(links)

	// Render the html links as part of the index page object div
	var indexPageDiv strings.Builder
	err = theTemplates.ExecuteTemplate(&indexPageDiv,
		"object.html",
		object{"main-page", links})
	if err != nil {
		panic(err)
	}

	// Render the main page on the browser
	page.Title = "Index"
	page.Content = template.HTML(indexPageDiv.String())
	indexPageDiv.Reset()
	err = theTemplates.ExecuteTemplate(w, "layout.html", page)
	theTemplates.ExecuteTemplate(&indexPageDiv, "layout.html", page)
	if err != nil {
		panic(err)
	}
	fmt.Println(indexPageDiv.String())
}

func main() {
	// for _, tmp := range theTemplates.Templates() {
	// 	fmt.Println(tmp.Name())
	// }
	http.HandleFunc("GET /{$}", hanleMainPage)
	// http.HandleFunc("GET /artists{$}", hanleArtistsPage)

	log.Panicln(http.ListenAndServe(":8080", nil))
}
