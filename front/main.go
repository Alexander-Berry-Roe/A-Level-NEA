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

	db = Mysql_db{}

	config := loadConfig("config.yaml")

	db.open_db(config.Mysql_username, config.Mysql_password, config.Mysql_address, config.Mysql_database)

	db.addUser("alex", "test")
	db.giveUserPermission(db.getUserByUsername("alex").ID, db.getPermissionID("manageUsers"))

	r.Use(authMiddleware)

	srv := &http.Server{
		Handler: r,
		Addr:    config.Address,

		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	r.HandleFunc("/api/manageUsers/newuser", apiAddUser).Methods("POST")
	r.HandleFunc("/api/manageUsers/removeuser", apiRemoveUser).Methods("POST")

	r.HandleFunc("/recording/{id:[a-zA-Z0-9].+}/{id:[a-zA-Z0-9].+}", getRecording)
	r.HandleFunc("/stream/{id:[a-zA-Z0-9].+}", test)
	r.HandleFunc("/login/auth", checkUser).Methods("POST")
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("static"))))

	log.Println("Server listening on " + srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

func getRecording(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func test(w http.ResponseWriter, r *http.Request) {

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
