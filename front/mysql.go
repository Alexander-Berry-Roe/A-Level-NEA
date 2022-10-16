package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type user_template struct {
	ID        string
	pass_hash string
	username  string
}

type token_template struct {
	token   string
	user_ID string
	exp     int
}

type camera_template struct {
	id   string
	url  string
	name string
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

type Mysql_db struct {
	db *sql.DB
}

// Open databse connection, must be called before any other methods is called
func (mysql_db *Mysql_db) open_db(username string, password string, address string, database string) {

	db, err := sql.Open("mysql", username+":"+password+"@tcp("+address+")/"+database)
	if err != nil {
		panic(err)
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
func (database *Mysql_db) db_load() {

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
func (database *Mysql_db) add_user(user string, pass string) bool {

	var user_info = user_template{}

	pass_hash, err := hashPassword(pass)

	if err != nil {
		log.Println(err)
	}

	if database.get_user(user).username != "" {
		return false
	}

	res, err := database.db.Query("INSERT INTO users (ID, username , pass_hash) VALUES (UUID(),?, ?)", user, pass_hash)
	if err != nil {
		log.Println(err)
	}

	for res.Next() {
		err = res.Scan(&user_info.ID)
		ErrorCheck(err)
	}

	res.Close()
	//database.db.Close()
	return true
}

// Check if a token is in the list of allowed tokens, return if true
func (database *Mysql_db) check_token(token string) token_template {
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
func (database *Mysql_db) allow_token(token string, userID string, exp int) {

	res, err := database.db.Query("insert into tokens (token, user_ID, exp) VALUES (? , ?, ?);", token, userID, exp)
	if err != nil {
		log.Println(err)
	}
	res.Close()
	//database.db.Close()
}

// Return the user record from the database
func (database *Mysql_db) get_user(user_ID string) user_template {

	var user user_template

	res, err := database.db.Query("SELECT ID, pass_hash, username FROM users WHERE username = ?", user_ID)
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
func (database *Mysql_db) add_camera(name string, url int) {

	res, err := database.db.Query("insert into cameras (ID, url) VALUES (UUID(), ?);", url)
	if err != nil {
		log.Println(err)
	}

	res.Close()

}

func (database *Mysql_db) get_camera(ID string) camera_template {

	var camera camera_template

	res, err := database.db.Query("SELECT * FROM cameras WHERE ID = ?", ID)
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
