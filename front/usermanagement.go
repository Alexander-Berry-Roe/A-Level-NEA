package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func apiAddUser(w http.ResponseWriter, r *http.Request) {
	var new_user user_info

	if !checkPermissions(r, "manageUsers") {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&new_user); err != nil {
		log.Println(err)
		http.Error(w, "Error decoidng response object", http.StatusBadRequest)
		return
	}

	if db.addUser(new_user.ID, new_user.Pass) {
		fmt.Fprintf(w, "OK")
	} else {
		fmt.Fprintf(w, "Already exists")
	}

}

func apiRemoveUser(w http.ResponseWriter, r *http.Request) {
	if !checkPermissions(r, "manageUsers") {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	var user user_info_name
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Println(err)
		http.Error(w, "Error decoidng response object", http.StatusBadRequest)
		return
	}

	if db.removeUserbyName(user.ID) {
		fmt.Fprintf(w, "OK")
	} else {
		fmt.Fprintf(w, "User does not exist")
	}

}
