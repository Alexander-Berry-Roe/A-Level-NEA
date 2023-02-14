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
	id      int
	url     string
	name    string
	enabled bool
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

	res, err := mysql_db.db.Query("SELECT CameraID, url, name, enabled from cameras;")

	if err != nil {
		log.Println(err)
	}

	var camera camera

	//Adds the list of cameras from the database to the array
	for res.Next() {
		res.Scan(&camera.id, &camera.url, &camera.name, &camera.enabled)
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

func (mysql_db *Mysql_db) getLiveSegments(id int64) []recordingSegment {
	var recordList []recordingSegment
	res, err := mysql_db.db.Query("SELECT * FROM (SELECT start, end, duration, location FROM recordings WHERE CameraID = ? ORDER BY end DESC limit 3) AS `table`  ORDER BY end ASC;", id)
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

//Set and get Resolution

func (mysql_db *Mysql_db) setResolution(resoltion []int, cameraID int) {
	_, err := mysql_db.db.Query("UPDATE cameras SET width = ? , height = ? WHERE CameraID = ?;", resoltion[0], resoltion[1], cameraID)
	if err != nil {
		log.Println(err)
	}
}

func (mysql_db *Mysql_db) getResolution(cameraID int) []int {
	var resolution []int
	res, err := mysql_db.db.Query("SELECT (width, height) FROM cmaeras WHERE CameraID = ?", cameraID)
	if err != nil {
		log.Println(err)
	}
	res.Scan(resolution[0], resolution[1])

	return resolution

}
