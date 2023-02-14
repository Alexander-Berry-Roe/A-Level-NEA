package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Monitors struct {
	listofMonitors LinkedList
}

type Monitor struct {
	id               int
	url              string
	control          chan int
	running          bool
	playlist         string
	sequence         int
	mediaTag         int
	restartCounter   int
	automaticRestart bool
}

type cameraRequest struct {
	CameraIDs []int `json:"cameraIDs"`
}

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

func (monitors *Monitors) getMonitorById(id int) *Monitor {
	for i, e := range monitors.listofMonitors.toArray() {
		if e.id == id {
			return monitors.listofMonitors.get(i)
		}
	}
	var empty Monitor
	return &empty
}

func (monitors *Monitors) getExists(id int) bool {
	for _, e := range monitors.listofMonitors.toArray() {
		if e.id == id {
			return true
		}
	}
	return false
}

// Resoultion methods.
func (monitor *Monitor) getResolution() ([]int, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=width,height", "-of", "csv=s=x:p=0", monitor.url)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(stdout)
	scanner.Scan()
	resolution := scanner.Text()

	if resolution == "" {
		var empty []int
		return empty, errors.New("Unable to connect to camera")
	}

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

	resTemp := strings.Split(resolution, "x")
	if len(resTemp) != 2 {
		var empty []int
		return empty, errors.New("Failed to fetch resoltion")
	}
	var resoltionArray []int
	for _, e := range resTemp {
		tmp, err := strconv.Atoi(e)
		if err != nil {
			var empty []int
			return empty, errors.New("Failed to fetch resoltion")
		}
		resoltionArray = append(resoltionArray, tmp)
	}
	return resoltionArray, nil
}

func (monitor *Monitor) getResolutionFromDB() []int {
	return db.getResolution(monitor.id)
}

func (monitor *Monitor) captureCamera(dir string) {

	resolution, err := monitor.getResolution()
	if err != nil {
		log.Println(err)
		time.Sleep(10 * time.Second)
		monitor.restartCapture()
		return
	}
	db.setResolution(resolution, monitor.id)

	monitor.running = true
	//FFMPEG command to capture video stream and save to HLS format the playlist file is returned to go via the stdpipe
	cmd := exec.Command("ffmpeg", "-fflags", "nobuffer",
		"-rtsp_transport", "tcp", "-i", monitor.url, "-vsync", "0", "-copyts", "-vcodec",
		"copy", "-movflags", "frag_keyframe+empty_moov", "-an", "-hls_flags", "delete_segments+append_list", "-f",
		"segment", "-segment_list_flags", "live", "-segment_time", "4", "-segment_list_size", "1", "-segment_format", "mpegts",
		"-segment_list", "pipe:1", "-segment_list_type", "m3u8", "stream/"+dir+"/%d.ts", "-vf", "fps=1/8", "-update", "1", "stream/png/"+strconv.Itoa(monitor.id)+".png", "-y")

	//Open the stdout pipe
	stdout, _ := cmd.StdoutPipe()
	//Start ffmpeg
	cmd.Start()

	//Start scanner on stodut
	scanner := bufio.NewScanner(stdout)

	var tempPlaylistStore string
	var duration float64

	automaticRestartStop := make(chan int)
	//Starts the automatic restart if not already started.
	go monitor.automaticCameraReconecter(automaticRestartStop)

	//Scans every line of the playlist file and converts it to a series of recoridng records.
	for {
		select {
		//checks for stopping command
		case <-monitor.control:
			automaticRestartStop <- 0
			cmd.Process.Kill()
			monitor.running = false
			log.Println("STOPPING CAPTURE")
			return
		default:

			//Scans stdout, recieve playlist.
			if scanner.Scan() {
				//Explain := in golang secion
				m := scanner.Text()

				//Used to detect start of new playlist
				if m == "#EXTM3U" {
					monitor.playlist = tempPlaylistStore
					tempPlaylistStore = ""
				} else if len(m) == 0 {
					log.Println("Not really sure what this means")
				} else if len(m) > 16 && m[0:8] == "#EXTINF:" {

					duration, _ = strconv.ParseFloat(m[8:len(m)-1], 16)

				} else if m[0:1] != "#" {
					location := "/stream/" + dir + "/" + m
					tempPlaylistStore += location + "\n"
					db.createRecordingRecord(monitor.id, time.Now().Unix()-int64(duration), time.Now().Unix(), duration, location, false)
					monitor.mediaTag += 1
				} else {
					tempPlaylistStore += m + "\n"
				}
			}
			time.Sleep(50)
			break
		}
	}
}

func (monitors *Monitors) startCapture(id int) error {

	if monitors.getMonitorById(id).running {
		return errors.New("alreadyRunning")
	}
	go monitors.getMonitorById(id).startCapture()
	return nil

}

func (monitor *Monitor) startCapture() error {
	dir := strconv.Itoa(monitor.id)
	createDir("stream/" + strconv.Itoa(monitor.id))
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
	go monitor.captureCamera(dir + "/" + strconv.Itoa(highestDir))
	return nil
}

func (monitors *Monitors) stopCapture(id int) error {
	if !monitors.getMonitorById(id).running {
		return errors.New("notRunning")
	}

	monitors.getMonitorById(id).control <- 0
	return nil
}

func (monitor *Monitor) stopCapture() error {
	if !monitor.running {
		return errors.New("notRunning")
	}
	monitor.control <- 0
	return nil
}

func (monitor *Monitor) restartCapture() {
	log.Println("Stopping capture")
	monitor.stopCapture()
	time.Sleep(50)
	log.Println("Starting capture")
	monitor.startCapture()
}

func (monitors *Monitors) addMonitor(id int, url string) {
	var newMonitor Monitor
	newMonitor.id = id
	newMonitor.url = url
	newMonitor.control = make(chan int)
	monitors.listofMonitors.add(newMonitor)

}

func (monitors *Monitors) loadCameras() {
	cameras := db.get_camera_list()
	for _, e := range cameras {
		monitors.addMonitor(e.id, e.url)
		if e.enabled {
			monitors.startCapture(e.id)
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

func (monitor *Monitor) generateLivePlaylist() string {
	recordings := db.getLiveSegments(int64(monitor.id))

	var playlist string

	playlist = "#EXTM3U\n#EXT-X-VERSION:3\n#EXT-X-ALLOW-CACHE:NO\n#EXT-X-TARGETDURATION:4\n"

	playlist = playlist + "#EXT-X-MEDIA-SEQUENCE:" + strconv.Itoa(monitor.mediaTag) + "\n"

	for _, e := range recordings {
		playlist = playlist + "#EXTINF:" + strconv.FormatFloat(e.duration, 'f', -1, 64) + ",\n" + e.location + "\n"
	}

	return playlist

}

func mergeVideos(segments []recordingSegment, outputLocation string) error {
	// Build the input argument for FFmpeg
	input := ""
	for i, segment := range segments {
		input += fmt.Sprintf("-i %s", segment.location)
		if i < len(segments)-1 {
			input += " "
		}
	}

	// Build the filter complex argument for FFmpeg
	filterComplex := "concat=n=" + strconv.Itoa(len(segments)) + ":v=1:a=1"
	for i := range segments {
		filterComplex += fmt.Sprintf("[%d:0][%d:1]", i, i)
	}

	// Build the FFmpeg command
	cmd := exec.Command("ffmpeg", "-y", input, "-filter_complex", filterComplex, outputLocation)

	// Run the FFmpeg command
	return cmd.Run()
}

func mergetoMP4(start int64, end int64, id int64) {
	segmentList := db.getSegmentList(start, end, id)

	mergeVideos(segmentList, "stream/tmp/")
	return
}

func (Monitors *Monitors) reload(id int) error {
	if Monitors.getExists(id) {

	}
	return nil
}

func (monitor *Monitor) automaticCameraReconecter(control chan int) {
	unchangedCounter := 0
	mediaTagTemp := monitor.mediaTag
	for {
		select {
		case <-control:
			return
		default:
			if monitor.mediaTag == mediaTagTemp {
				unchangedCounter += 1

			} else {
				unchangedCounter = 0
				mediaTagTemp = monitor.mediaTag
			}

			if unchangedCounter >= 10 {
				log.Println("ALERT: Capture hasn't advanced for over 10 seconds will attempt restart")
				monitor.restartCounter += 1
				go monitor.restartCapture()
				return
			}
			time.Sleep(1 * time.Second)
		}
	}
}
