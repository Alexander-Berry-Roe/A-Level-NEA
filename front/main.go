package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type camera_proxy struct {
}

var db Mysql_db

func main() {

	r := mux.NewRouter()
	//
	db = Mysql_db{}

	//Load the configuration file into memory
	config := loadConfig("config.yaml")

	//Open conneciton with database
	db.open_db(config.Mysql_username, config.Mysql_password, config.Mysql_address, config.Mysql_database)

	//Set all requests to use authentication middleware
	r.Use(authMiddleware)

	//Sets web server configuation options
	srv := &http.Server{
		Handler: r,
		Addr:    config.Address,

		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	//Sets up routing for POST requests.
	r.HandleFunc("/api/manageUsers/newuser", apiAddUser).Methods("POST")
	r.HandleFunc("/api/manageUsers/removeuser", apiRemoveUser).Methods("POST")

	r.HandleFunc("/recording/{id:[a-zA-Z0-9].+}/{id:[a-zA-Z0-9].+}", getRecording)
	r.HandleFunc("/stream/{id:[a-zA-Z0-9].+}", getVideoFile)
	r.HandleFunc("/api/login/auth", checkUser).Methods("POST")
	r.HandleFunc("/api/login/status", isLoggedIn).Methods("GET")
	r.HandleFunc("/api/getSelfUser", apiGetOwnUserInfo).Methods("GET")
	r.HandleFunc("/api/logout", apiLogout).Methods("GET")
	r.HandleFunc("/api/getAccountMenu", getAccountMenuOptions).Methods("GET")
	r.HandleFunc("/api/setOwnUsername", changeOwnUsername).Methods("POST")
	r.HandleFunc("/api/getAllStreams", apiGetStreamUrls).Methods("GET")
	r.HandleFunc("/api/getCameraListSideMenu", apiGetCameraListSideMenu).Methods("GET")
	r.HandleFunc("/api/setCameraSettings", apiSetCameraSettings)
	r.HandleFunc("/api/getCameraSettings", apiGetCameraSettings).Methods("GET")
	//Starts the webserver

	log.Println("Server listening on " + srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

func getRecording(w http.ResponseWriter, r *http.Request) {
	//This is just for testing
	fmt.Fprintf(w, "OK")
}

func getVideoFile(w http.ResponseWriter, r *http.Request) {
	//This is still a very early prototype
	resp, err := http.Get("http://localhost:8081/" + r.URL.Path)
	if err != nil {
		log.Println("Falied to connect to backend")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(w, "Error")
		return
	}
	resp.Body.Close()
	w.Write(body)

}
