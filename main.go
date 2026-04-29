package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"runtime/debug"
)

var API string = "https://groupietrackers.herokuapp.com/api"

var page struct {
	Title   string
	Content template.HTML
}



func getApiResponseBody(path string) (body []byte, err error) {
	apiLink := API + path
	var resp *http.Response
	resp, err = http.Get(apiLink)
	if err != nil {
		log.Println(err, string(debug.Stack()))
		return nil, err
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err, string(debug.Stack()))
		return nil, err
	}
	return
}

func hanleMainPage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/artists?page=1", http.StatusSeeOther)
}

func main() {
	// for _, tmp := range theTemplates.Templates() {
	// 	fmt.Println(tmp.Name())
	// }
	http.HandleFunc("GET /{$}", hanleMainPage)
	http.HandleFunc("GET /artists", handleArtistsPage)

	log.Println("Server running on port 8080")
	log.Println(http.ListenAndServe(":8080", nil))
}
