package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Mysql_db struct {
	db *sql.DB
}

type camera struct {
	id   string
	url  string
	name string
}

//Opens the connection pool
func (mysql_db *Mysql_db) open_db(username string, password string, address string, database string) {

	db, err := sql.Open("mysql", username+":"+password+"@tcp("+address+")/"+database)
	if err != nil {
		panic(err)
	}
	//Set connection pool parameters
	db.SetConnMaxLifetime(time.Hour)
	db.SetMaxOpenConns(10000)
	db.SetMaxIdleConns(10)

	mysql_db.db = db

}

//Returns an array containing a list of cameras
func (mysql_db *Mysql_db) get_camera_list() []camera {
	var cameras []camera

	res, err := mysql_db.db.Query("SELECT ID, url, name from cameras;")

	if err != nil {
		log.Println(err)
	}

	var camera camera

	//Adds the list of cameras from the database to the array
	for res.Next() {
		res.Scan(&camera.id, &camera.url, &camera.name)
		cameras = append(cameras, camera)
	}

	res.Close()
	return cameras

}
