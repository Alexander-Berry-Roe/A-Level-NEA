package main

import (
	"fmt"
	"log"
	"net/http"
)

var db Mysql_db

func main() {

	config := loadConfig("config.yaml")

	db.open_db(config.Mysql_username, config.Mysql_password, config.Mysql_address, config.Mysql_database)

	//Concurrently launches the camera handler
	go camereHandle()

	// configure the songs directory name and port
	const webDir = ""
	const port = 8081

	// add a handler for the song files
	http.Handle("/", addHeaders(http.FileServer(http.Dir(webDir))))
	fmt.Printf("Starting server on %v\n", port)
	log.Printf("Serving %s on HTTP port: %v\n", webDir, port)

	// serve and log errors
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

// addHeaders will act as middleware to give us CORS support
func addHeaders(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h.ServeHTTP(w, r)
	}
}
