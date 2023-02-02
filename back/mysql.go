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
	id        int
	url       string
	name      string
	streamURL string
}

type recordingSegment struct {
	start    int64
	end      int64
	duration float64
	location string
}

// Opens the connection pool
func (mysql_db *Mysql_db) open_db(username string, password string, address string, database string) {

	db, err := sql.Open("mysql", username+":"+password+"@tcp("+address+")/"+database)
	if err != nil {
		log.Println(err)
	}
	//Set connection pool parameters
	db.SetConnMaxLifetime(time.Hour)
	db.SetMaxOpenConns(1000)
	db.SetMaxIdleConns(10)

	mysql_db.db = db

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

// Returns an array containing a list of cameras
func (mysql_db *Mysql_db) get_camera_list() []camera {
	var cameras []camera

	res, err := mysql_db.db.Query("SELECT CameraID, url, name, streamURL from cameras;")

	if err != nil {
		log.Println(err)
	}

	var camera camera

	//Adds the list of cameras from the database to the array
	for res.Next() {
		res.Scan(&camera.id, &camera.url, &camera.name, &camera.streamURL)
		cameras = append(cameras, camera)
	}

	res.Close()
	return cameras

}

func (mysql_db *Mysql_db) createRecordingRecord(cameraID int, start int64, end int64, duration float64, location string, protected bool) {
	res, err := mysql_db.db.Query("INSERT INTO recordings (CameraID, start, end, duration, location, protected) VALUES (?, ?, ?, ?, ?, ?);", cameraID, start, end, duration, location, protected)
	if err != nil {
		log.Println(err)
		return
	}
	res.Close()
}

func (mysql_db *Mysql_db) getSegmentList(start int64, end int64, id int64) []recordingSegment {
	var recordList []recordingSegment
	res, err := mysql_db.db.Query("SELECT start, end, duration ,location FROM recordings WHERE start >= ? AND end <= ? AND CameraID = ?;", start, end, id)
	if err != nil {
		log.Println(err)
		return recordList
	}
	var tempRecord recordingSegment
	for res.Next() {
		res.Scan(&tempRecord.start, &tempRecord.end, &tempRecord.duration, &tempRecord.location)
		recordList = append(recordList, tempRecord)
	}
	res.Close()
	return recordList

}
