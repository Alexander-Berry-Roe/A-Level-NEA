package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type user_info struct {
	ID   string `json:"id"`
	Pass string `json:"passwd"`
}

type user_info_name struct {
	ID string `json:"id"`
}

func authMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Get cookie, sent with request.
		token, err := r.Cookie("session_token")
		var user_id token_template
		//If no cookie sent with request
		if err != nil {
			if err == http.ErrNoCookie {
				var newToken string
				//Generate cookie value
				newToken = uuid.NewString()
				expires := time.Now().Add(120 * time.Minute)
				//Issue cookie
				http.SetCookie(w, &http.Cookie{
					Name:    "session_token",
					Value:   newToken,
					Expires: expires,
					Path:    "/",
				})
				//return unauthised response
				resp := simplejson.New()
				resp.Set("auth", false)
				payload, err := resp.MarshalJSON()
				if err != nil {
					log.Println(err)
				}

				w.Header().Add("Content-Type", "application/json")
				w.Write(payload)

				return

			} else {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		//Calls method to run a database querry to fetch the user data asociated with the users session cookie value (token)
		user_id = db.checkToken(token.Value)

		//Exemption, allows access the api functions required for login to unauthenticated clients, requests to /api/login are for logging in and should be allowed for unauthenticated users
		if len(r.URL.Path) > 10 {
			if r.URL.Path[0:10] == "/api/login" {
				next.ServeHTTP(w, r)
				return

			}
		}
		//Only serve response if client is logged in as a valid user.
		if user_id.user_ID != "" {
			next.ServeHTTP(w, r)
			return
		} else {
			w.WriteHeader(http.StatusForbidden)
		}
	})
}

// Checks if passowrd entered is correct and registers token to user.
func checkUser(w http.ResponseWriter, r *http.Request) {
	resp := simplejson.New()
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
			resp.Set("auth", true)

		}

	} else {
		resp.Set("auth", false)
	}

	payload, err := resp.MarshalJSON()
	if err != nil {
		log.Println(err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(payload)

}

func isLoggedIn(w http.ResponseWriter, r *http.Request) {
	resp := simplejson.New()
	token, err := r.Cookie("session_token")
	if err != nil {
		println(err)
	}

	user := db.getUserByToken(token.Value)

	if user.username != "" {
		resp.Set("auth", true)

	} else {
		resp.Set("auth", false)
	}

	payload, err := resp.MarshalJSON()
	if err != nil {
		log.Println(err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(payload)

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

// Function used to check if user has a particular permission.
func checkPermissions(r *http.Request, permissionName string) bool {
	token, err := r.Cookie("session_token")
	if err != nil {
		log.Println(err)
		return false
	}

	return db.checkUserPermission(db.getUserByToken(token.Value).ID, permissionName)
}
