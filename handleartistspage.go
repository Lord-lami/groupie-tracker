package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"runtime/debug"
	"strconv"
	"sync"

	"github.com/Lord-lami/render-html"
)

type artistDetail struct {
	Id      render.Ignored         `json:"id"`
	Image   render.ImageLinkString `json:"image"`
	Name    string                 `json:"name"`
	Details render.LinkString      `json:"-"`
}

var nbrOfArtistsPerPage int = 8

func handleArtistsPage(w http.ResponseWriter, r *http.Request) {
	// Receive and validate the page number
	pageNumStr := r.URL.Query().Get("page")
	if pageNumStr == "" {
		http.Redirect(w, r, "/artists?page=1", http.StatusSeeOther)
		return
	}
	pageNumInt, err := strconv.Atoi(pageNumStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err, string(debug.Stack()))
		return
	}

	// concurrently get the response body from the API for
	// nbrOfArtistsPerPage number of artists
	var wg sync.WaitGroup
	artistsDetails := make([]artistDetail, nbrOfArtistsPerPage)
	firstId := (pageNumInt-1)*nbrOfArtistsPerPage + 1
	lastId := pageNumInt * nbrOfArtistsPerPage
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

	// filter out empty artistsDetails elements
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

// renderPageNav renders the page navigation links
func renderPageNav(pageNumInt int) template.HTML {
	type pageNavLinkString string
	var pageNavigator struct {
		LeftArrow, PageNumber, RightArrow pageNavLinkString
	}

	pageNumStr := strconv.Itoa(pageNumInt)
	lastId := pageNumInt * nbrOfArtistsPerPage
	if pageNumInt > 1 {
		pageNavigator.LeftArrow = pageNavLinkString("/artists?page=" + strconv.Itoa(pageNumInt-1))
	}
	pageNavigator.PageNumber = pageNavLinkString("/artists?page=" + strconv.Itoa(pageNumInt))
	_, err := getApiResponseBody("/artists/" + strconv.Itoa(lastId+1))
	if err == nil {
		pageNavigator.RightArrow = pageNavLinkString("/artists?page=" + strconv.Itoa(pageNumInt+1))
	}

	// renderPageNavLinkString renders the pageNavLinkString with
	// the correct link text instead of using the field name.
	renderPageNavLinkString := func(name string, data any) (pageNavLinkHTML template.HTML) {
		linkText := ""
		if data.(pageNavLinkString) != "" {
			switch name {
			case "LeftArrow":
				linkText = "<"
			case "PageNumber":
				linkText = pageNumStr
			case "RightArrow":
				linkText = ">"
			default:
				panic("wrong field name for pageNavLinkString: " + name)
			}
			return render.RenderBasic("linkstring.html")(linkText, data)
		}
		return ""

	}

	render.MapTypeToRenderFunc[pageNavLinkString](renderPageNavLinkString)

	return render.RenderObj("page-navigator", pageNavigator)
}
