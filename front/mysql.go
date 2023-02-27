package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type userTemplate struct {
	ID        int
	pass_hash string
	username  string
}

type token_template struct {
	token   string
	user_ID string
	exp     int
}

type cameraTemplate struct {
	id         int
	url        string
	name       string
	height     int
	width      int
	enabled    bool
	defaultExp int
}

type recording struct {
	id        string
	start     int
	end       int
	startFile int
	endFile   int
	exp       int
	protected bool
}

type layout struct {
	empty  bool
	width  int
	height int
	posX   int
	posY   int
}

type Mysql_db struct {
	db *sql.DB
}

// Open databse connection, must be called before any other methods is called
func (mysql_db *Mysql_db) openDb(username string, password string, address string, database string) {

	db, err := sql.Open("mysql", username+":"+password+"@tcp("+address+")/"+database)
	if err != nil {
		log.Fatalln("Unable to connect to databas")
		os.Exit(1)

	}

	//Sets database up databse connection pool settings
	db.SetConnMaxLifetime(time.Hour)
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(100)

	mysql_db.db = db

}

func ErrorCheck(err error) {
	if err != nil {
		log.Panic(err.Error())
	}
}

// Setup
func (database *Mysql_db) dbLoad() {

	create_user_table, err := database.db.Query("CREATE TABLE if not exists users (ID varchar(255) NOT NULL, pass_hash varchar(255), salt varchar(255) , PRIMARY KEY (ID));")
	if err != nil {
		log.Println(err.Error())
	}
	create_user_table.Close()

	create_token_table, err := database.db.Query("CREATE TABLE if not exists tokens (token varchar(255) NOT NULL, user_ID varchar(255) NOT NULL, PRIMARY KEY (token), FOREIGN KEY (user_ID) REFERENCES users(ID));")
	if err != nil {
		log.Println(err.Error())
	}
	create_token_table.Close()

}

// Add user to the user table
func (database *Mysql_db) addUser(user string, pass string) bool {

	pass_hash, err := hashPassword(pass)

	if err != nil {
		log.Println(err)
	}

	if database.getUserByUsername(user).username != "" {
		return false
	}

	res, err := database.db.Query("INSERT INTO users (userName, passHash) VALUES (?, ?)", user, pass_hash)
	if err != nil {
		println(err.Error)
		log.Println(err)
		res.Close()
		return false
	}

	res.Close()
	//database.db.Close()
	return true
}

// Check if a token is in the list of allowed tokens, return if true
func (database *Mysql_db) checkToken(token string) token_template {
	var token_info = token_template{}

	test_db := database.db

	//Delete tokens expired tokens
	res, err := test_db.Query("DELETE FROM tokens WHERE exp <= ? ;", time.Now().Unix())
	if err != nil {
		println(err)
	}

	res.Close()

	//Fetch the token if it exists
	res, err = test_db.Query("SELECT * FROM tokens WHERE token = ?", token)
	if err != nil {
		println(err)
	}

	if res == nil {
		return token_info
	}

	for res.Next() {
		err = res.Scan(&token_info.token, &token_info.user_ID, &token_info.exp)
		ErrorCheck(err)
	}

	res.Close()
	return token_info
}

// Add token to the table of allowed tokens
func (database *Mysql_db) allowToken(token string, userID int, exp int) {

	res, err := database.db.Query("insert into tokens (token, UserID, exp) VALUES (? , ?, ?);", token, userID, exp)
	if err != nil {
		log.Println(err)
	}
	res.Close()

}

// Return the user record from the database
func (database *Mysql_db) getUser(user_ID int) userTemplate {

	var user userTemplate

	res, err := database.db.Query("SELECT UserID, passHash, userName FROM users WHERE username = ?", user_ID)
	if err != nil {
		println(err)
	}

	for res.Next() {
		err = res.Scan(&user.ID, &user.pass_hash, &user.username)
		ErrorCheck(err)
	}

	res.Close()
	return user
}

// Add a record to the camera table

func (database *Mysql_db) getCamera(ID string) cameraTemplate {

	var camera cameraTemplate

	res, err := database.db.Query("SELECT (CameraID, url, name) FROM cameras WHERE ID = ?", ID)
	if err != nil {
		println(err)
	}

	for res.Next() {
		err = res.Scan(&camera.id, &camera.url, &camera.name)
		ErrorCheck(err)
	}

	res.Close()
	return camera
}

func (database *Mysql_db) getAllMointors() []cameraTemplate {
	var cameras []cameraTemplate

	res, err := database.db.Query("SELECT CameraID, url, name, height, width, enabled, defaultExp FROM cameras;")
	if err != nil {
		println(err)
	}

	var tempcamera cameraTemplate
	for res.Next() {
		res.Scan(&tempcamera.id, &tempcamera.url, &tempcamera.name, &tempcamera.height, &tempcamera.width, &tempcamera.enabled, &tempcamera.defaultExp)
		cameras = append(cameras, tempcamera)
	}

	return cameras

}

func (database *Mysql_db) getRecordings(CameraID string, Start int, End int) []recording {
	var recordings []recording

	res, err := database.db.Query("SELECT ID, url, name from cameras;")

	if err != nil {
		log.Println(err)
	}

	var recording recording

	for res.Next() {
		res.Scan(&recording.id)
	}

	res.Close()
	return recordings

}

func (db *Mysql_db) getUserByUsername(username string) userTemplate {
	var user userTemplate

	res, err := db.db.Query("SELECT UserID, userName, passHash from users where userName = ?;", username)
	if err != nil {
		log.Println(err)
	}

	for res.Next() {
		res.Scan(&user.ID, &user.username, &user.pass_hash)
	}

	return user
}

func (db *Mysql_db) getUserByToken(token string) userTemplate {
	var user userTemplate

	res, err := db.db.Query("SELECT users.userID, users.userName, users.passHash FROM users, tokens WHERE users.userID = tokens.userID AND tokens.token = ?;", token)
	if err != nil {
		log.Println(err)
	}

	for res.Next() {
		res.Scan(&user.ID, &user.username, &user.pass_hash)
	}

	return user

}

func (db *Mysql_db) getPermissionID(permissionName string) int {
	res, err := db.db.Query("SELECT permissionID FROM permissions WHERE permissionName = ?;", permissionName)
	if err != nil {
		log.Println(err)
		res.Close()
		return -1
	}

	var id int
	for res.Next() {
		res.Scan(&id)
	}
	return id
}

func (db *Mysql_db) giveUserPermission(userID int, permissionID int) bool {
	res, err := db.db.Query("INSERT INTO permissionLink (permissionID, userID) values(?, ?);", permissionID, userID)
	if err != nil {
		log.Println(err)
		return false
	}
	res.Close()
	return true
}

func (db *Mysql_db) checkUserPermission(userID int, permissionName string) bool {
	res, err := db.db.Query("SELECT users.userID FROM users, permissionLink, permissions WHERE users.userID = permissionLink.userID AND permissionLink.permissionID = permissions.permissionID AND permissions.permissionName = ? AND users.userID = ?;", permissionName, userID)
	if err != nil {
		log.Println(err)
		res.Close()
		return false
	}
	var id int
	for res.Next() {
		res.Scan(&id)
	}
	res.Close()
	if id == userID {
		return true
	}
	return false
}

func (db *Mysql_db) removeUserbyID(userID int) bool {
	if db.getUser(userID).username == "" {
		return false
	}

	res, err := db.db.Query("DELETE FROM users WHERE userID = ? ", userID)
	if err != nil {
		log.Println(err)
		return false
	}
	res.Close()
	return true
}

func (db *Mysql_db) removeUserbyName(username string) bool {
	if db.getUserByUsername(username).username == "" {
		return false
	}

	res, err := db.db.Query("DELETE FROM users WHERE userName = ? ", username)
	if err != nil {
		log.Println(err)
		return false
	}
	res.Close()
	return true

}

func (db *Mysql_db) removeToken(tokenValue string) bool {
	res, err := db.db.Query("DELETE FROM tokens WHERE token = ?", tokenValue)
	if err != nil {
		log.Println(err)
		return false
	}
	res.Close()
	return true
}

func (db *Mysql_db) changeUsernameByToken(token string, newUsername string) bool {
	user := db.getUserByToken(token)

	res, err := db.db.Query("UPDATE users SET userName = ? where UserID = ?;", newUsername, user.ID)
	if err != nil {
		log.Println(err)
		return false
	}
	res.Close()
	return true
}

func (db *Mysql_db) addCamera(url string, name string) {
	res, err := db.db.Query("INSERT INTO cameras (url, name) VALUES (?, ?)", url, name)
	if err != nil {
		log.Println(err)
	}
	res.Close()
}

// Methods for handling live player layout for each user.
func (db *Mysql_db) getCameraLayout(userID int, cameraID int) layout {
	var response layout
	response.empty = true
	res, err := db.db.Query("SELECT width, height, posX, posY FROM livePlayers WHERE userID = ? AND cameraID = ? ", userID, cameraID)
	if err != nil {
		log.Println(err)
	}
	if res.Next() {
		response.empty = false
		res.Scan(&response.width, &response.height, &response.posX, &response.posY)
	}
	res.Close()
	return response
}

// Unused TO BE DLETED LATTER AFTER DOCUMENTATION.
func (db *Mysql_db) setCameraLayout(userID int, cameraID int, width int, height int, posx int, posy int) bool {

	db.db.Query("DELETE FROM livePlayers WHERE cameraID = ?;", cameraID)
	res, err := db.db.Query("INSERT INTO livePlayers (userID, cameraID, width, height, posX, posY) VALUES (?, ?, ?, ?, ?, ?);", userID, cameraID, width, height, posx, posy)
	if err != nil {
		log.Println(err)
		return false
	}
	res.Close()
	return true

}
