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

var page int = 0

var Page struct {
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

func getArtistsPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, r.URL.Path)
	nbrOfItemsInPage := 6
	artistsDetails := make([]artistDetail, nbrOfItemsInPage)
	for i := range 6 {
		id := strconv.Itoa(page*nbrOfItemsInPage + (i + 1))
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
	Page.Title = "Artists"
	Page.Content = template.HTML(content.String())
	err = theTemplates.ExecuteTemplate(w, "layout.html", Page)
	if err != nil {
		panic(err)
	}
}

func main() {
	http.HandleFunc("GET /", getArtistsPage)

	log.Panicln(http.ListenAndServe(":8080", nil))
}
