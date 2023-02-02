package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type camera_proxy struct {
}

var db Mysql_db

var recordings monitors

func main() {

	r := mux.NewRouter()

	//
	db = Mysql_db{}

	//recordings.loadCameras()

	//Load the configuration file into memory
	config := loadConfig("config.yaml")

	//Open conneciton with database
	db.open_db(config.Mysql_username, config.Mysql_password, config.Mysql_address, config.Mysql_database)

	recordings.loadCameras()

	//Set all requests to use authentication middleware

	//Sets web server configuation options
	srv := &http.Server{
		Handler: r,
		Addr:    config.Address,

		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	r.Path("/stream/record/{id:[0-9]+}/{start:[0-9]+}/{end:[0-9]+}").HandlerFunc(getRecordedPlaylist)
	//Sets up routing for POST requests.
	r.HandleFunc("/api/startCapture", startCapture).Methods("POST")
	r.HandleFunc("/api/stopCapture", stopCapture).Methods("POST")
	r.HandleFunc("/api/getLivePlaylist", getLivePlaylits).Methods("GET")
	r.PathPrefix("/stream/").Handler(http.StripPrefix("/stream/", http.FileServer(http.Dir("./stream/"))))
	//Starts the webserver

	log.Println("Server listening on " + srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

func startCapture(w http.ResponseWriter, r *http.Request) {
	var cameras cameraRequest

	if err := json.NewDecoder(r.Body).Decode(&cameras); err != nil {
		log.Println(err)
		http.Error(w, "Error decoidng response object", http.StatusBadRequest)
		return
	}

	for _, e := range cameras.CameraIDs {
		recordings.startCapture(e)
	}
}

func stopCapture(w http.ResponseWriter, r *http.Request) {
	var cameras cameraRequest

	if err := json.NewDecoder(r.Body).Decode(&cameras); err != nil {
		log.Println(err)
		http.Error(w, "Error decoidng response object", http.StatusBadRequest)
		return
	}

	for _, e := range cameras.CameraIDs {
		recordings.stopCapture(e)
	}

}

func getLivePlaylits(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-mpegURL")
	fmt.Fprintf(w, recordings.listofMonitors[1].playlist)

}

func getRecordedPlaylist(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	start, _ := strconv.ParseInt(mux.Vars(r)["start"], 10, 64)
	end, _ := strconv.ParseInt(mux.Vars(r)["end"], 10, 64)

	//Generate a playlist file between requested times, not live
	playlist := gernatePlaylistForTime(start, end, id, false)

	w.Header().Set("Content-Type", "application/x-mpegURL")
	fmt.Fprint(w, playlist)
}
