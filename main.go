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
	"strconv"
	"strings"
)

var templateFS fs.FS = os.DirFS("templates")

var theTemplates *template.Template = template.Must(template.ParseFS(templateFS, "**/*.html"))

var pageNum int = 0

var page struct {
	Title   string
	Content template.HTML
}

type artistDetail struct {
	Id           uint16
	Image        string
	Name         string
	Members      []string
	CreationDate uint16
	FirstAlbum   string
	Locations    string
	ConcertDates string
	Relations    string
}

type pageLink struct {
	Name string
	Url  string
}

type object struct {
	Name     string
	Elements []string
}

func getArtistspage(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, r.URL.Path)
	nbrOfItemsInpage := 6
	artistsDetails := make([]artistDetail, nbrOfItemsInpage)
	for i := range 6 {
		id := strconv.Itoa(pageNum*nbrOfItemsInpage + (i + 1))
		artistLink := "https://groupietrackers.herokuapp.com/api/artists/" + id
		resp, err := http.Get(artistLink)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		// fmt.Println("Response Status:", resp.Status)
		body := []byte{}
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		// fmt.Printf("As String: %s\n", body)
		if err := json.Unmarshal(body, &artistsDetails[i]); err != nil {
			panic(err)
		}
		// fmt.Println()
		// fmt.Printf("%#v\n", artistDetails)

	}
	var content strings.Builder
	err := theTemplates.ExecuteTemplate(&content, "artists.html", artistsDetails)
	if err != nil {
		panic(err)
	}
	page.Title = "Artists"
	page.Content = template.HTML(content.String())
	err = theTemplates.ExecuteTemplate(w, "layout.html", page)
	if err != nil {
		panic(err)
	}
}

func groupieTrackerApiResponseBody(path string) (body []byte) {
	apiLink := "https://groupietrackers.herokuapp.com/api" + path
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

func hanleMainpage(w http.ResponseWriter, r *http.Request) {
	// Get the response body from the API
	var apiPathMainData map[string]string
	err := json.Unmarshal(groupieTrackerApiResponseBody(""), &apiPathMainData)
	if err != nil {
		panic(err)
	}

	// Render the links from the API as html anchor tags
	links := []string{}
	for name, url := range apiPathMainData {
		var linkHtml strings.Builder
		err := theTemplates.ExecuteTemplate(&linkHtml,
			"pageurlstring.html",
			pageLink{name, url})
		if err != nil {
			panic(err)
		}
		links = append(links, linkHtml.String())
	}

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
	err = theTemplates.ExecuteTemplate(w, "layout.html", page)
	if err != nil {
		panic(err)
	}
}

func main() {
	http.HandleFunc("GET /{$}", hanleMainpage)

	log.Panicln(http.ListenAndServe(":8080", nil))
}
