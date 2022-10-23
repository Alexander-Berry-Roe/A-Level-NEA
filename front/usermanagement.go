package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

//REST API add user funciton
func apiAddUser(w http.ResponseWriter, r *http.Request) {
	var new_user user_info
	
	//Checks if user has permission to manage users.
	if !checkPermissions(r, "manageUsers") {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	
	//Decode json sent by post request.
	if err := json.NewDecoder(r.Body).Decode(&new_user); err != nil {
		log.Println(err)
		http.Error(w, "Error decoidng response object", http.StatusBadRequest)
		return
	}
	//Calls the database function to add user
	if db.addUser(new_user.ID, new_user.Pass) {
		fmt.Fprintf(w, "OK")
	} else {
		fmt.Fprintf(w, "Already exists")
	}

}

//REST API fucntion to remove user
func apiRemoveUser(w http.ResponseWriter, r *http.Request) {
	//Check if user has mangeUser permissions
	if !checkPermissions(r, "manageUsers") {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	var user user_info_name
	//Decode json sent by post request.
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Println(err)
		http.Error(w, "Error decoidng response object", http.StatusBadRequest)
		return
	}
	//Calls the database function to remove user
	if db.removeUserbyName(user.ID) {
		fmt.Fprintf(w, "OK")
	} else {
		fmt.Fprintf(w, "User does not exist")
	}

}
