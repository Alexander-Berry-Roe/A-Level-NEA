package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/bitly/go-simplejson"
)

type options struct {
	accountMenuOptions []AccountMenuOption `json:"accountMenuOptions"`
}

type AccountMenuOption struct {
	Title string `json:"title"`
}

type NewUserName struct {
	Username string `json:"username"`
}

// REST API add user funciton
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

// REST API fucntion to remove user
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

func apiGetOwnUserInfo(w http.ResponseWriter, r *http.Request) {
	resp := simplejson.New()
	token, err := r.Cookie("session_token")
	if err != nil {
		log.Println(err)
	}

	user := db.getUserByToken(token.Value)

	resp.Set(
		"id", user.username,
	)
	payload, err := resp.MarshalJSON()
	if err != nil {
		log.Println(err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(payload)

}

func apiLogout(w http.ResponseWriter, r *http.Request) {
	resp := simplejson.New()
	token, err := r.Cookie("session_token")
	if err != nil {
		log.Println(err)
	}

	if db.removeToken(token.Value) {
		resp.Set("success", true)
	} else {
		resp.Set("success", false)
	}

	payload, err := resp.MarshalJSON()
	if err != nil {
		log.Panicln(err)
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(payload)
}

func getAccountMenuOptions(w http.ResponseWriter, r *http.Request) {
	var response []AccountMenuOption
	_, err := r.Cookie("session_token")
	if err != nil {
		log.Println(err)
	}

	var option AccountMenuOption
	if true {
		option.Title = "Account settings"
		response = append(response, option)

	}
	option.Title = "Logout"
	response = append(response, option)

	payload, err := json.Marshal(response)
	if err != nil {
		log.Print(err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(payload)

}

func changeOwnUsername(w http.ResponseWriter, r *http.Request) {
	resp := simplejson.New()
	token, err := r.Cookie("session_token")
	if err != nil {
		log.Println(err)
		return
	}

	var user NewUserName
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Println(err)
		http.Error(w, "Error decoidng response object", http.StatusBadRequest)
		return
	}

	existingUser := db.getUserByUsername(user.Username)

	if existingUser.username != "" {
		resp.Set("success", false)
	} else if db.changeUsernameByToken(token.Value, user.Username) {
		resp.Set("success", true)
	} else {
		resp.Set("success", false)
	}

	payload, err := resp.MarshalJSON()
	if err != nil {
		log.Panicln(err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(payload)
}
