package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"runtime/debug"
	"strconv"
	"sync"
)

type artistDetail struct {
	Id           ignored         `json:"id"`
	Image        imageLinkString `json:"image"`
	Name         string          `json:"name"`
	Members      []string        `json:"members"`
	CreationDate int             `json:"creationDate"`
	FirstAlbum   dateString      `json:"firstAlbum"`
	Locations    apiLinkString   `json:"locations"`
	ConcertDates apiLinkString   `json:"concertDates"`
	Relations    apiLinkString   `json:"relations"`
}

func handleArtistsPage(w http.ResponseWriter, r *http.Request) {
	// Receive and validate the page number
	pageNumStr := r.URL.Query().Get("page")
	if pageNumStr == "" {
		http.Redirect(w, r, "/artists?page=1", http.StatusSeeOther)
		return
	}
	pageNumInt, err := strconv.Atoi(pageNumStr)
	if err != nil {
		log.Println(err, string(debug.Stack()))
		return
	}

	// Get the response body from the API for 6 artists per page
	nbrOfItemsPerPage := 6
	var wg sync.WaitGroup
	artistsDetails := make([]artistDetail, nbrOfItemsPerPage)
	firstId := (pageNumInt-1)*nbrOfItemsPerPage + 1
	lastId := pageNumInt * nbrOfItemsPerPage
	for i := firstId; i <= lastId; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			id := strconv.Itoa(i)
			responseBody, err := getApiResponseBody("/artists/" + id)
			if err != nil {
				switch err := err.(type) {
				case *url.Error:
					if err.Timeout() {
						w.WriteHeader(http.StatusRequestTimeout)
					} else {
						w.WriteHeader(http.StatusServiceUnavailable)
					}
				default:
					w.WriteHeader(http.StatusInternalServerError)
				}
			}
			err = json.Unmarshal(responseBody, &artistsDetails[i-1])
			if err != nil {
				log.Println(err, string(debug.Stack()))
				return
			}
		}(i)
	}
	wg.Wait()
	artistsPageList := renderArr("artists-page", artistsDetails)

	page.Title = "Artist Page"
	page.Content = artistsPageList
	err = theTemplates.ExecuteTemplate(w, "layout.html", page)
	if err != nil {
		log.Println(err, string(debug.Stack()))
		return
	}
}
