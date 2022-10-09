package main

import (
	"encoding/json"
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

	//db.get_user("91ae8aae-398f-11ed-b9eb-8a280ec120f6")

	r.Use(authMiddleware)

	srv := &http.Server{
		Handler: r,
		Addr:    config.Address,

		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	r.HandleFunc("/api/users/newuser", api_add_user).Methods("POST")

	r.HandleFunc("/recording/{id:[a-zA-Z0-9].+}/{id:[a-zA-Z0-9].+}", getRecording)
	r.HandleFunc("/stream/{id:[a-zA-Z0-9].+}", test)
	r.HandleFunc("/login/auth", checkUser).Methods("POST")
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("static"))))

	log.Println("Server listening on " + srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

func api_add_user(w http.ResponseWriter, r *http.Request) {

	var new_user user_info

	if err := json.NewDecoder(r.Body).Decode(&new_user); err != nil {
		log.Println(err)
		http.Error(w, "Error decoidng response object", http.StatusBadRequest)
		return
	}

	if db.add_user(new_user.ID, new_user.Pass) {
		fmt.Fprintf(w, "OK")
	} else {
		fmt.Fprintf(w, "Already exists")
	}

}

func getRecording(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func test(w http.ResponseWriter, r *http.Request) {

	resp, err := http.Get("http://localhost:8081/" + r.URL.Path)
	if err != nil {
		print(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		print(err)
	}
	resp.Body.Close()
	w.Write(body)

}
