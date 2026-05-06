package main

import (
	// "github.com/Lord-lami/render-html"
	"encoding/json"
	"groupie-tracker/render"
	"log"
	"net/http"
	"net/url"
	"runtime/debug"
	"strconv"
)

type anArtistRelations struct {
	Id             render.Ignored `json:"id"`
	DatesLocations labeledMap     `json:"datesLocations"`
}

type anArtistDetail struct {
	Id           render.Ignored         `json:"id"`
	Image        render.ImageLinkString `json:"image"`
	Name         render.Ignored         `json:"name"`
	Members      labeledStringSlice     `json:"members"`
	CreationDate labeledInt             `json:"creationDate"`
	FirstAlbum   labeledDateString      `json:"firstAlbum"`
}

func handleAnArtistPage(w http.ResponseWriter, r *http.Request) {
	render.MapTypeToRenderFunc[labeledStringSlice](renderLabeledStringSlice)
	render.MapTypeToRenderFunc[labeledInt](renderLabeledInt)
	render.MapTypeToRenderFunc[labeledDateString](renderLabeledDateString)

	id := r.PathValue("id")
	if _, err := strconv.Atoi(id); err != nil {
		log.Println(err, string(debug.Stack()))
		return
	}

	artistRelationChan := make(chan anArtistRelations)
	go channelApiData(w, "/relation/"+id, artistRelationChan)

	artistDetailChan := make(chan anArtistDetail)
	go channelApiData(w, "/artists/"+id, artistDetailChan)

	artistRelation := <-artistRelationChan
	artistDetail := <-artistDetailChan
	if artistRelation.Id == nil || artistDetail.Id == nil {
		return
	}
	artistDetailHTML := render.RenderObj("ArtistDetails", artistDetail)
	artistRelationHTML := renderLabeledMap("ConcertLocationsAndDates",
		artistRelation.DatesLocations)

	page.Title = artistDetail.Name.(string)
	page.Content = artistDetailHTML + artistRelationHTML
	err := templates.ExecuteTemplate(w, "layout.html", page)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err, string(debug.Stack()))
		return
	}
}

func channelApiData[T any](w http.ResponseWriter, path string, dataChan chan T) {
	respBody, err := getApiResponseBody(path)

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
			if err == emptyBody {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
		var nothing T
		dataChan <- nothing
		return
	}

	var dataStruct T
	err = json.Unmarshal(respBody, &dataStruct)
	if err != nil {
		log.Println(err, string(debug.Stack()))
		w.WriteHeader(http.StatusInternalServerError)
		var nothing T
		dataChan <- nothing
		return
	}
	dataChan <- dataStruct
}
