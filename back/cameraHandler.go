package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type monitors struct {
	listofMonitors []monitor
}

type monitor struct {
	id       int
	url      string
	control  chan int
	running  bool
	playlist string
}

type cameraRequest struct {
	CameraIDs []int `json:"cameraIDs"`
}

var monitorlist monitors

/*
Control codes for capture camera
0: stop
*/
func createDir(dir string) {
	_, err := os.Stat(dir)
	if err == nil {
		return
	}

	err = os.Mkdir(dir, 0755)
	if err != nil {
		log.Fatal(err)
	}
}

func (monitors *monitors) captureCamera(url string, dir string, control chan int, id int, index int) {

	//FFMPEG command to capture video stream and save to HLS format the playlist file is returned to go via the stdpipe
	cmd := exec.Command("ffmpeg", "-fflags", "nobuffer",
		"-rtsp_transport", "tcp", "-i", url, "-vsync", "0", "-copyts", "-vcodec",
		"copy", "-movflags", "frag_keyframe+empty_moov", "-an", "-hls_flags", "delete_segments+append_list", "-f",
		"segment", "-segment_list_flags", "live", "-segment_time", "4", "-segment_list_size", "1", "-segment_format", "mpegts",
		"-segment_list", "pipe:1", "-segment_list_type", "m3u8", "stream/"+dir+"/%d.ts")

	//Open the stdout pipe
	stdout, _ := cmd.StdoutPipe()

	//Start ffmpeg
	cmd.Start()

	//Start scanner on stodut
	scanner := bufio.NewScanner(stdout)

	var tempPlaylistStore string
	var duration float64
	//Scans every line of the playlist file and converts it to a series of recoridng records.
	for {
		select {
		case <-control:
			cmd.Process.Kill()
		default:
			scanner.Scan()
			m := scanner.Text()
			if m == "#EXTM3U" {
				monitors.listofMonitors[index].playlist = tempPlaylistStore
				tempPlaylistStore = ""
			} else if len(m) > 16 && m[0:8] == "#EXTINF:" {

				duration, _ = strconv.ParseFloat(m[8:len(m)-1], 16)

			} else if m[0:1] != "#" {
				location := "/stream/" + dir + "/" + m
				tempPlaylistStore += location + "\n"
				db.createRecordingRecord(id, time.Now().Unix()-int64(duration), time.Now().Unix(), duration, location, false)
			} else {
				tempPlaylistStore += m + "\n"
			}
			break
		}
	}
}

func (monitors *monitors) startCapture(id int) {
	for i, e := range monitors.listofMonitors {
		if e.id == id {
			dir := strconv.Itoa(id)
			files, err := ioutil.ReadDir("stream/" + dir)
			if err != nil {
				log.Fatal(err)
			}
			highestDir := 0
			for _, file := range files {
				if file.IsDir() {
					currentDir, _ := strconv.ParseInt(file.Name(), 10, 64)
					if highestDir <= int(currentDir) {
						highestDir = int(currentDir) + 1
					}
				}
			}

			createDir("stream/" + dir + "/" + strconv.Itoa(highestDir))
			go monitors.captureCamera(e.url, dir+"/"+strconv.Itoa(highestDir), e.control, e.id, i)
			return

		}
	}
}

func (monitors *monitors) addMonitor(id int, url string) {
	var newMonitor monitor
	newMonitor.id = id
	newMonitor.url = url
	newMonitor.control = make(chan int)
	monitors.listofMonitors = append(monitors.listofMonitors, newMonitor)

}

func (monitors *monitors) loadCameras() {
	cameras := db.get_camera_list()
	for _, e := range cameras {
		monitors.addMonitor(e.id, e.url)
	}

}

func (monitors *monitors) stopCapture(id int) {
	for _, e := range monitors.listofMonitors {
		if e.id == id {
			e.control <- 0
			return
		}
	}
}

func gernatePlaylistForTime(start int64, end int64, id int64, live bool) string {
	recordings := db.getSegmentList(start, end, id)

	var playlist string

	playlist = "#EXTM3U\n#EXT-X-VERSION:3\n#EXT-X-ALLOW-CACHE:NO\n#EXT-X-TARGETDURATION:4\n"

	for _, e := range recordings {
		playlist = playlist + "#EXTINF:" + strconv.FormatFloat(e.duration, 'f', -1, 64) + ",\n" + e.location + "\n"
	}

	if !live {
		playlist = playlist + "#EXT-X-ENDLIST"
	}

	return playlist
}
