package main

import (
	// "github.com/Lord-lami/render-html"
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
	Id      render.Ignored         `json:"id"`
	Image   render.ImageLinkString `json:"image"`
	Name    string                 `json:"name"`
	Details render.LinkString      `json:"-"`
}

var nbrOfItemsPerPage int = 6

func renderPageNav(pageNumInt int) template.HTML {
	type pageNavLinkString string
	var pageNavigator struct {
		LeftArrow, PageNumber, RightArrow pageNavLinkString
	}
	pageNumStr := strconv.Itoa(pageNumInt)
	lastId := pageNumInt * nbrOfItemsPerPage
	if pageNumInt > 1 {
		pageNavigator.LeftArrow = pageNavLinkString("/artists?page=" + strconv.Itoa(pageNumInt-1))
	}
	pageNavigator.PageNumber = pageNavLinkString("/artists?page=" + strconv.Itoa(pageNumInt))
	_, err := getApiResponseBody("/artists/" + strconv.Itoa(lastId+1))
	if err == nil {
		pageNavigator.RightArrow = pageNavLinkString("/artists?page=" + strconv.Itoa(pageNumInt+1))
	}
	render.MapTypeToRenderFunc[pageNavLinkString](func(name string, data any) (pageNavLinkHTML template.HTML) {
		linkText := ""
		if data.(pageNavLinkString) != "" {
			switch name {
			case "LeftArrow":
				linkText = "<"
			case "PageNumber":
				linkText = pageNumStr
			case "RightArrow":
				linkText = ">"
			}
		}
		return render.NewRenderFunc[pageNavLinkString]("linkstring.html")(linkText, data)
	})
	return render.RenderObj("page-navigator", pageNavigator)
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
	var wg sync.WaitGroup
	artistsDetails := make([]artistDetail, nbrOfItemsPerPage)
	firstId := (pageNumInt-1)*nbrOfItemsPerPage + 1
	lastId := pageNumInt * nbrOfItemsPerPage
	for i := 0; i <= lastId-firstId; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			id := strconv.Itoa(firstId + i)
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
					if err == emptyBody {
						return
					} else {
						w.WriteHeader(http.StatusInternalServerError)
					}
				}
				log.Println(err, string(debug.Stack()))
				return
			}
			err = json.Unmarshal(responseBody, &artistsDetails[i])
			if err != nil {
				log.Println(err, string(debug.Stack()))
				return
			}
			artistsDetails[i].Details = render.LinkString("/artists/" + id)
		}(i)
	}
	wg.Wait()

	filteredArtistsDetails := []artistDetail{}
	for _, ad := range artistsDetails {
		if ad.Id != render.Ignored(nil) {
			filteredArtistsDetails = append(filteredArtistsDetails, ad)
		}
	}

	var artistsPageList, pageNavigatorDiv template.HTML
	if len(filteredArtistsDetails) != 0 {
		artistsPageList = render.RenderArr("artist-list", filteredArtistsDetails)
		pageNavigatorDiv = renderPageNav(pageNumInt)
	} else {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	page.Title = "Artist Page"
	page.Content = artistsPageList + pageNavigatorDiv
	err = templates.ExecuteTemplate(w, "layout.html", page)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err, string(debug.Stack()))
		return
	}
}
