package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type artistDetail struct {
	Id           ignored   `json:"id"`
	Image        imageLinkString   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate string   `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    apiLinkString   `json:"locations"`
	ConcertDates apiLinkString   `json:"concertDates"`
	Relations    apiLinkString   `json:"relations"`
}

func getArtistsPage(w http.ResponseWriter, r *http.Request) {

	nbrOfItemsPerPage := 6
	artistsDetails := make([]artistDetail, nbrOfItemsPerPage)
	for i := range nbrOfItemsPerPage {
		go func(i int) {
			id := strconv.Itoa((pageNum-1)*nbrOfItemsPerPage + (i + 1))
			responseBody := groupieTrackerApiResponseBody("/artists/" + id)

			// fmt.Printf("As String: %s\n", body)
			if err := json.Unmarshal(responseBody, &artistsDetails[i]); err != nil {
				panic(err)
			}
		}(i)

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
	// Receive and validate the page number
	pageNum, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		panic(err)
	}

	// Get the response body from the API for 6 artists per page
	nbrOfItemsPerPage := 6
	artistsDetails := make([]artistDetail, nbrOfItemsPerPage)
	var wg sync.WaitGroup
	for i := range nbrOfItemsPerPage {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			id := strconv.Itoa(pageNum*nbrOfItemsPerPage + (i + 1))
			err := json.Unmarshal(groupieTrackerApiResponseBody("/artists/"+id),
				&artistsDetails[i])
			if err != nil {
				panic(err)
			}
		}(i)
	}
	wg.Wait()
}
