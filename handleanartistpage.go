package main

import (
	"encoding/json"
	"groupie-tracker/render"
	"log"
	"net/http"
	"net/url"
	"runtime/debug"
	"strconv"
)

type anArtistRelations struct {
	Id             render.Ignored                 `json:"id"`
	DatesLocations map[string][]render.DateString `json:"datesLocations"`
}

type anArtistDetail struct {
	Id           render.Ignored         `json:"id"`
	Image        render.ImageLinkString `json:"image"`
	Name         string                 `json:"name"`
	Members      []string               `json:"members"`
	CreationDate int                    `json:"creationDate"`
	FirstAlbum   render.DateString      `json:"firstAlbum"`
}

func channelApiData[T any](w http.ResponseWriter, path string, dataChan chan T) {
	respBody, err := getApiResponseBody(path)
	if err != nil {
		log.Println(err, string(debug.Stack()))
		var nothing T
		dataChan <- nothing
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
	if string(respBody) == `{"id":0,"image":"","name":"","members":null,"creationDate":0,"firstAlbum":"","locations":"","concertDates":"","relations":""}` {
		log.Println("Empty response body", string(debug.Stack()))
		w.WriteHeader(http.StatusNotFound)
		var nothing T
		dataChan <- nothing
		return
	}
	var dataStruct T
	err = json.Unmarshal(respBody, dataStruct)
	if err != nil {
		log.Println(err, string(debug.Stack()))
		w.WriteHeader(http.StatusInternalServerError)
		var nothing T
		dataChan <- nothing
		return
	}
	dataChan <- dataStruct
}

func handleAnArtistPage(w http.ResponseWriter, r *http.Request) {
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
	artistDetailHTML := render.RenderObj(artistDetail.Name+" Details", artistDetail)
	artistRelationHTML := render.RenderObj(artistDetail.Name+"Concert Locations and Dates",
		artistRelation)
	
	err := render.TheTemplates.ExecuteTemplate(w, "layout.html", artistDetailHTML+artistRelationHTML)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err, string(debug.Stack()))
		return
	}
}
