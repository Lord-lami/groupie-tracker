package main

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type artistDetail struct {
	Id           uint16   `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate string   `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

func getArtistspage(w http.ResponseWriter, r *http.Request) {
	nbrOfItemsPerpage := 6
	artistsDetails := make([]artistDetail, nbrOfItemsPerpage)
	for i := range 6 {
		id := strconv.Itoa(pageNum*nbrOfItemsPerpage + (i + 1))
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

func hanleArtistsPage(w http.ResponseWriter, r *http.Request) {
	// Get the response body from the API for 6 artists per page
	nbrOfItemsPerpage := 6
	artistsDetails := make([]artistDetail, nbrOfItemsPerpage)
	for i := range nbrOfItemsPerpage {
		id := strconv.Itoa(pageNum*nbrOfItemsPerpage + (i + 1))
		go func(id string) {
			err := json.Unmarshal(groupieTrackerApiResponseBody("/artists/"+id), &artistsDetails[i])
			if err != nil {
				panic(err)
			}
		}(id)
	}
}

func renderApiLink(name, url string) template.HTML {
	var ok bool
	url, ok = strings.CutPrefix(url, API)
	if !ok {
		panic("This is not an api link string of API: " + API)
	}
	var linkHtml strings.Builder
	err := theTemplates.ExecuteTemplate(&linkHtml,
		"linkstring.html",
		Link{name, url})
	if err != nil {
		panic(err)
	}
	return template.HTML(linkHtml.String())
}

func renderApiObject(data *any) {
	apiObjVal := reflect.ValueOf(data).Elem()
	elements := []template.HTML{}
	for i := range apiObjVal.NumField() {
		switch apiObjVal.Field(i).Kind() {
		case reflect.TypeFor[apiLinkString]().Kind():
			name := apiObjVal.Type().Field(i).Name
			url := apiObjVal.Field(i).String()
			linkHTML := renderApiLink(name, url)
			elements = append(elements, linkHTML)
		}
	}
	// Render the html links as part of the index page object div
	var apiObjDivBuilder strings.Builder
	err := theTemplates.ExecuteTemplate(&apiObjDivBuilder,
		"object.html",
		object{"main-page", elements})
	if err != nil {
		panic(err)
	}
	
}
