package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MysqlDb struct {
	db *sql.DB
}

type camera struct {
	id      int
	url     string
	name    string
	enabled bool
	exp     int
}

type recordingSegment struct {
	cameraID  int
	start     int64
	end       int64
	duration  float64
	location  string
	protected bool
	exp       int64
}

// Opens the connection pool
func (mysqlDb *MysqlDb) openDb(username string, password string, address string, database string) {

	db, err := sql.Open("mysql", username+":"+password+"@tcp("+address+")/"+database)
	if err != nil {
		log.Println(err)
	}
	//Set connection pool parameters
	db.SetConnMaxLifetime(time.Hour)
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(100)

	mysqlDb.db = db

}

func (db *MysqlDb) getPermissionID(permissionName string) int {
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
	res.Close()
	return id
}

// Returns an array containing a list of cameras
func (mysqlDb *MysqlDb) getCameraList() []camera {
	var cameras []camera

	res, err := mysqlDb.db.Query("SELECT CameraID, url, name, enabled, defaultExp from cameras;")

	if err != nil {
		log.Println(err)
	}

	var camera camera

	//Adds the list of cameras from the database to the array
	for res.Next() {
		res.Scan(&camera.id, &camera.url, &camera.name, &camera.enabled, &camera.exp)
		cameras = append(cameras, camera)
	}

	res.Close()
	return cameras

}

func (mysqlDb *MysqlDb) createRecordingRecord(cameraID int, start int64, end int64, duration float64, location string, protected bool, exp int64) {
	res, err := mysqlDb.db.Query("INSERT INTO recordings (CameraID, start, end, duration, location, protected, exp) VALUES (?, ?, ?, ?, ?, ?, ?);", cameraID, start, end, duration, location, protected, exp)
	if err != nil {
		log.Println(err)
		return
	}
	res.Close()
}

func (mysqlDb *MysqlDb) getSegmentList(start int64, end int64, id int64) []recordingSegment {
	var recordList []recordingSegment
	res, err := mysqlDb.db.Query("SELECT start, end, duration ,location FROM recordings WHERE start >= ? AND end <= ? AND CameraID = ?;", start, end, id)
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

func (mysqlDb *MysqlDb) getLiveSegments(id int64) []recordingSegment {
	var recordList []recordingSegment
	res, err := mysqlDb.db.Query("SELECT * FROM (SELECT start, end, duration, location FROM recordings WHERE CameraID = ? ORDER BY end DESC limit 3) AS `table`  ORDER BY end ASC;", id)
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

//No longer in use rem

func (mysqlDb *MysqlDb) setResolution(resoltion []int, cameraID int) {
	res, err := mysqlDb.db.Query("UPDATE cameras SET width = ? , height = ? WHERE CameraID = ?;", resoltion[0], resoltion[1], cameraID)
	if err != nil {
		log.Println(err)
	}
	res.Close()
}

func (mysqlDb *MysqlDb) getResolution(cameraID int) []int {
	var resolution []int
	res, err := mysqlDb.db.Query("SELECT (width, height) FROM cmaeras WHERE CameraID = ?;", cameraID)
	if err != nil {
		log.Println(err)
	}
	res.Scan(resolution[0], resolution[1])
	res.Close()
	return resolution
}

func (mysqlDb *MysqlDb) getExpiredRecords(time int64) []recordingSegment {
	var recordingList []recordingSegment
	res, err := mysqlDb.db.Query("SELECT cameraID, start, end, duration, location, protected, exp FROM recordings WHERE exp <= ?;", time)
	if err != nil {
		log.Println(err)
		return recordingList
	}

	for res.Next() {
		var segment recordingSegment
		res.Scan(&segment.cameraID, &segment.start, &segment.end, &segment.duration, &segment.location, &segment.protected, &segment.exp)
		recordingList = append(recordingList, segment)
	}
	res.Close()
	return recordingList
}

func (mysqlDb *MysqlDb) deleteExpiredRecords(time int64) {
	res, err := mysqlDb.db.Query("DELETE FROM recordings WHERE exp <= ?;", time)
	if err != nil {
		log.Println(err)
	}
	res.Close()
}
