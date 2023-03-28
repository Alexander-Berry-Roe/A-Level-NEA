// Allow for reload camera data
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

var db MysqlDb

var recordings Monitors

func main() {

	r := mux.NewRouter()

	//Load the configuration file into memory
	config := loadConfig("config.yaml")

	//Creating instance of the Mysql class (this needs changing the cammel casing asap).
	db = MysqlDb{}

	//Open conneciton with database
	db.open_db(config.MysqlUsername, config.MysqlPassword, config.MysqlAddress, config.MysqlDatabase)

	recordings.loadCameras()
	go recordings.automaticDelete()
	//Set all requests to use authentication middleware

	//Sets web server configuation options
	srv := &http.Server{
		Handler: r,
		Addr:    config.Address,

		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	r.Path("/stream/record/{id:[0-9]+}/{start:[0-9]+}/{end:[0-9]+}.m3u8").HandlerFunc(getRecordedPlaylist)
	//Sets up routing for POST requests.
	r.HandleFunc("/api/startCapture", startCapture).Methods("POST")
	r.HandleFunc("/api/stopCapture", stopCapture).Methods("POST")
	r.HandleFunc("/stream/getLivePlaylist/{id:[0-9]+}.m3u8", getLivePlaylists).Methods("GET")
	r.HandleFunc("/stream/getRunning", getStreamUrls)
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

func getLivePlaylists(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)

	w.Header().Set("Content-Type", "application/x-mpegURL")
	fmt.Fprintf(w, recordings.getMonitorById(int(id)).generateLivePlaylist())

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

type StreamResponse struct {
	Id  int    `json:"id"`
	Url string `json:url`
}

func getStreamUrls(w http.ResponseWriter, r *http.Request) {
	var resp []StreamResponse

	for _, e := range recordings.listofMonitors.toArray() {
		if e.running {
			var tmpResp StreamResponse
			tmpResp.Id = e.id
			tmpResp.Url = "/stream/getLivePlaylist/" + strconv.Itoa(e.id) + "/"
			resp = append(resp, tmpResp)
		}
	}

	payload, err := json.Marshal(resp)
	if err != nil {
		log.Print(err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(payload)
}
