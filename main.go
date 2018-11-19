package main

import (
	"github.com/Hagelund96/paragliding/handler"
	"github.com/Hagelund96/paragliding/struct"
	"net/http"
	"os"
)

//main function for application. Initialises storage database
func main() {
	_struct.Db.Init()
	var p string
	if port := os.Getenv("PORT"); port != "" {
		p = ":" + port
	} else {
		p = ":8080"
	}
	//different handlers for urls
	http.HandleFunc("/paragliding", handler.Handler)
	http.HandleFunc("/paragliding/api/", handler.HandlerApi)
	//http.HandleFunc("/paragliding/api/igc/", handler.HandlerIgc)
	http.HandleFunc("/paragliding/api/track", handler.HandlerTrack)
	http.HandleFunc("/paragliding/api/track/{id}", handler.HandlerID)
	http.HandleFunc("/paragliding/api/track/{id}/{field}", handler.HandlerField)
	http.HandleFunc("/paragliding/admin/api/tracks_count", handler.AdminTracksCount)
	http.HandleFunc("/paragliding/admin/api/tracks", handler.AdminTracks)


	http.ListenAndServe(p, nil)
}