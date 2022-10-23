package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type user_info struct {
	ID   string `json:"id"`
	Pass string `json:"pass"`
}

type user_info_name struct {
	ID string `json:"id"`
}

func authMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token, err := r.Cookie("session_token")
		var newToken string
		var user_id token_template
		if err != nil {
			if err == http.ErrNoCookie {
				newToken = uuid.NewString()
				expires := time.Now().Add(120 * time.Minute)
				http.SetCookie(w, &http.Cookie{
					Name:    "session_token",
					Value:   newToken,
					Expires: expires,
					Path:    "/",
				})
				http.Redirect(w, r, "/login/login.html", http.StatusFound)
				return

			} else {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		user_id = db.checkToken(token.Value)

		if len(r.URL.Path) > 6 {
			if r.URL.Path[0:6] == "/login" {
				if user_id.token == "" {
					next.ServeHTTP(w, r)
					return
				} else {
					http.Redirect(w, r, "/", http.StatusFound)
					return
				}
			}
		}

		if user_id.user_ID != "" {
			next.ServeHTTP(w, r)
			return
		} else {
			http.Redirect(w, r, "/login/login.html", http.StatusFound)
			return
		}
	})
}

//Checks if passowrd entered is correct and registers token to user.
func checkUser(w http.ResponseWriter, r *http.Request) {

	token, err := r.Cookie("session_token")

	if err != nil {
		log.Println(err)
		return
	}

	login := user_info{}

	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		log.Println(err)
		http.Error(w, "Error decoidng response object", http.StatusBadRequest)
		return
	}

	user := db.getUserByUsername(login.ID)

	if user.username != "" {
		if doPasswordsMatch(user.pass_hash, login.Pass) {
			db.allowToken(token.Value, user.ID, int(time.Now().Unix()+3600))
			fmt.Fprintf(w, "OK")
			return
		}

	}

	w.Header().Add("Content-Type", "application/json")

}

// Calculate a new hash for a new password
func hashPassword(password string) (string, error) {
	var passwordBytes = []byte(password)

	hashedPasswordBytes, err := bcrypt.
		GenerateFromPassword(passwordBytes, bcrypt.MinCost)

	return string(hashedPasswordBytes), err
}

// check if the hash of the entered password matches the hash in the database
func doPasswordsMatch(hashedPassword, currPassword string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword), []byte(currPassword))

	if err == nil {
		return true
	} else {
		return false
	}
}

//Function used to check if user has a particular permission. 
func checkPermissions(r *http.Request, permissionName string) bool {
	token, err := r.Cookie("session_token")
	if err != nil {
		log.Println(err)
		return false
	}

	return db.checkUserPermission(db.getUserByToken(token.Value).ID, permissionName)
}
