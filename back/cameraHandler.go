package main

import (
	"log"
	"os"
	"os/exec"
)

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

func captureCamera(url string, dir string) {
	//FFMPEG command to capture video stream and save to HLS format
	_, err := exec.Command("ffmpeg", "-fflags", "nobuffer", "-rtsp_transport", "tcp", "-i", url, "-vsync", " 0", "-copyts", "-vcodec", "copy", "-movflags", "frag_keyframe+empty_moov", "-an", "-hls_flags", "delete_segments+append_list", "-f", "segment", "-segment_list_flags", "live", "-segment_time", "1", "-segment_list_size", "3", "-segment_format", "mpegts", "-segment_list", "stream/"+dir+"/index.m3u8", "-segment_list_type", "m3u8", "stream/"+dir+"/%d.ts").CombinedOutput()

	if err != nil {
		panic(err.Error())
	}

}

func camereHandle() {
	cameraList := db.get_camera_list()

	//Intitates a concurrent process for each camera
	for _, value := range cameraList {
		createDir("./stream/" + value.id)
		go captureCamera(value.url, value.id)
	}

}
