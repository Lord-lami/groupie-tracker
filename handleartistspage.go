package main

import (
	"encoding/json"
	"groupie-tracker/render"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"runtime/debug"
	"strconv"
	"sync"
)

type artistDetail struct {
	Id           render.Ignored         `json:"id"`
	Image        render.ImageLinkString `json:"image"`
	Name         string                 `json:"name"`
	Members      []string               `json:"members"`
	CreationDate int                    `json:"creationDate"`
	FirstAlbum   render.DateString      `json:"firstAlbum"`
	Locations    apiLinkString          `json:"locations"`
	ConcertDates apiLinkString          `json:"concertDates"`
	Relations    apiLinkString          `json:"relations"`
}

// renderApiLink function is in renderapilink.go

func handleArtistsPage(w http.ResponseWriter, r *http.Request) {
	render.RenderTypeFunc[render.TypeString[apiLinkString]()] = renderApiLink
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
				log.Println(err, string(debug.Stack()))
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
				return
			}
			err = json.Unmarshal(responseBody, &artistsDetails[i-firstId])
			if err != nil {
				log.Println(err, string(debug.Stack()))
				return
			}
		}(i)
	}
	wg.Wait()

	filteredArtistsDetails := []artistDetail{}
	for _, ad := range artistsDetails {
		if ad.Id.(float64) != 0 {
			filteredArtistsDetails = append(filteredArtistsDetails, ad)
		}
	}

	var artistsPageList, pageNavigatorDiv template.HTML
	if len(filteredArtistsDetails) != 0 {
		artistsPageList = render.RenderArr("artist-list", filteredArtistsDetails)

		type pageNavLinkString string
		var pageNavigator struct {
			LeftArrow, PageNumber, RightArrow pageNavLinkString
		}
		if pageNumInt > 1 {
			pageNavigator.LeftArrow = pageNavLinkString("/artists?page=" + strconv.Itoa(pageNumInt-1))
		}
		pageNavigator.PageNumber = pageNavLinkString("/artists?page=" + strconv.Itoa(pageNumInt))
		if len(filteredArtistsDetails) == nbrOfItemsPerPage {
			pageNavigator.RightArrow = pageNavLinkString("/artists?page=" + strconv.Itoa(pageNumInt+1))
		}
		render.RenderTypeFunc["main.pageNavLinkString"] = func(name string, data any) (pageNavLinkHTML template.HTML) {
			linkText := ""
			switch name {
			case "LeftArrow":
				linkText = "<"
			case "PageNumber":
				linkText = pageNumStr
			case "RightArrow":
				linkText = ">"
			}
			return render.RenderType[pageNavLinkString]("linkstring.html")(linkText, data)
		}
		pageNavigatorDiv = render.RenderObj("page-navigator", pageNavigator)
	} else {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	page.Title = "Artist Page"
	page.Content = artistsPageList + pageNavigatorDiv
	err = render.TheTemplates.ExecuteTemplate(w, "layout.html", page)
	if err != nil {
		log.Println(err, string(debug.Stack()))
		return
	}
}
