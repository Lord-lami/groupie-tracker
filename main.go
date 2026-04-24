package main

import (
	"encoding/json"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
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
	Name string
	Value string
}

type Link struct {
	Name string
	Url  string
}

type object struct {
	Name     string
	Elements []template.HTML
}

type dateString string
type apiLinkString string
type imageLinkString string
type ignored any

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
	responseBody := groupieTrackerApiResponseBody("")
	err := json.Unmarshal(responseBody, &mainPageData)
	if err != nil {
		panic(err)
	}

	// Render the links from the API as html anchor tags
	// that link to their display page
	indexPageDiv := renderApiObject("main-page", &mainPageData)

	// Render the main page on the browser
	page.Title = "Stalk A Band"
	page.Content = template.HTML(indexPageDiv)
	// indexPageDiv.Reset()
	err = theTemplates.ExecuteTemplate(w, "layout.html", page)
	// theTemplates.ExecuteTemplate(&indexPageDiv, "layout.html", page)
	if err != nil {
		panic(err)
	}
	// fmt.Println(indexPageDiv.String())
}

func main() {
	// for _, tmp := range theTemplates.Templates() {
	// 	fmt.Println(tmp.Name())
	// }
	http.HandleFunc("GET /{$}", hanleMainPage)
	http.Handle("GET /artists{$}", http.RedirectHandler("/artists?page=1", http.StatusPermanentRedirect))
	http.HandleFunc("GET /artists?page=", hanleArtistsPage)
	

	log.Panicln(http.ListenAndServe(":8080", nil))
}
