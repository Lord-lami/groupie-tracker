package main

import (
	"encoding/json"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"runtime/debug"
)

var API string = "https://groupietrackers.herokuapp.com/api"

var templateFS fs.FS = os.DirFS("templates")

var theTemplates *template.Template = template.Must(template.ParseFS(templateFS,
	"*.html",
	"*/*.html",
	"*/*/*.html"))

var page struct {
	Title   string
	Content template.HTML
}

type String struct {
	Name  string
	Value string
}



type Link struct {
	Name string
	Url  string
}





type ignored any
type dateString string
type apiLinkString string
type imageLinkString string

type mainPageDataHolder struct {
	Artists   apiLinkString `json:"artists"`
	Locations apiLinkString `json:"locations"`
	Dates     apiLinkString `json:"dates"`
	Relation  apiLinkString `json:"relation"`
}

func getApiResponseBody(path string) (body []byte) {
	apiLink := API + path
	resp, err := http.Get(apiLink)
	if err != nil {
		log.Println(err, debug.Stack())
		return []byte{}
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err, debug.Stack())
		return []byte{}
	}
	return
}

func hanleMainPage(w http.ResponseWriter, r *http.Request) {
	// Get the response body from the API
	mainPageData := mainPageDataHolder{}
	responseBody := getApiResponseBody("")
	err := json.Unmarshal(responseBody, &mainPageData)
	if err != nil {
		log.Println(err, debug.Stack())
		return
	}

	// Render the links from the API as html anchor tags
	// that link to their display page
	indexPageDiv := renderObj("main-page", mainPageData)

	// Render the main page on the browser
	page.Title = "Stalk A Band"
	page.Content = indexPageDiv
	// indexPageDiv.Reset()
	err = theTemplates.ExecuteTemplate(w, "layout.html", page)
	// theTemplates.ExecuteTemplate(&indexPageDiv, "layout.html", page)
	if err != nil {
		log.Println(err, debug.Stack())
		return
	}
	// fmt.Println(indexPageDiv.String())
}

func main() {
	// for _, tmp := range theTemplates.Templates() {
	// 	fmt.Println(tmp.Name())
	// }
	http.HandleFunc("GET /{$}", hanleMainPage)
	http.HandleFunc("GET /artists", handleArtistsPage)

	log.Println("Server running on port 8080")
	log.Println(http.ListenAndServe(":8080", nil))
}
