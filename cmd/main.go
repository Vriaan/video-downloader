package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"video-downloader/configloader"
	"video-downloader/downloader"
	"video-downloader/parsingelement"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

var configurationPath = pflag.StringP("configuration", "c", "", "To run the program needs to get some configuration")
var loggerLvl = pflag.StringP("log-lvl", "l", "InfoLevel", "The level of log to show [default = InfoLevel]. Available are (PanicLevel, FatalLevel, ErrorLevel, WarnLevel, InfoLevel, DebugLevel). For more information, llo, at Logrus doc 'type level'")

//var fileURLDownload = pflag.StringP("file_infos", "f", "", "The program will try to download all video from this file")
var videoURlDownload = pflag.StringP("url", "u", "", "The program will try to download the video on this url")
var videoSiteOrigin = pflag.StringP("origin", "o", "", "The program will parse according to the origin web site you gave (this options is case insensitive). Available are [u]")
var destinationPath = pflag.StringP("destination", "d", "", "The program will write every final video in this directory")
var audioOnly = pflag.BoolP("audio-only", "a", false, "This flag indicate to create an audio file from the video extracted to .mp3 format. It uses the program `ffmpeg` to do so. Note that this flag will remove the video file")

// Instanciate logrus, parse the program flags and load Configuration
func setUp() (*parsingelement.ParsingInformations, time.Time) {
	startTime := time.Now()
	pflag.Parse()

	switch *loggerLvl {
	case "DebugLevel":
		logrus.SetLevel(logrus.DebugLevel)
	case "InfoLevel":
		logrus.SetLevel(logrus.InfoLevel)
	case "WarnLevel":
		logrus.SetLevel(logrus.WarnLevel)
	case "ErrorLevel":
		logrus.SetLevel(logrus.ErrorLevel)
	case "FatalLevel":
		logrus.SetLevel(logrus.FatalLevel)
	case "PanicLevel":
		logrus.SetLevel(logrus.PanicLevel)

	default:
		logrus.Fatalf("Logrus logger level doesn't exist")
	}

	parsingInformations, err := configloader.ReadConfig(*configurationPath)
	if nil != err {
		logrus.Fatalf("Error loading configuration [%s], reason: %v", *configurationPath, err)
	}
	logrus.Debugf("Loaded configuration: %v", parsingInformations)

	if err = os.MkdirAll(*destinationPath, 0777); nil != err {
		logrus.Fatalf("Could not create temporary path for Downloaded video, reason %v", err)
	}

	return parsingInformations, startTime
}

func main() {
	// Load conf, Load URL site where a video has to be downloaded and init what has to be
	parsingInformations, startTime := setUp()

	videoInfos, err := parsingInformations.ParseOn(*videoSiteOrigin, *videoURlDownload)
	if nil != err {
		logrus.Fatalf("Error working on site [%s]; url [%s], reason: %v", *videoSiteOrigin, *videoSiteOrigin, err)
	}

	var fileGenerated string
	if fileGenerated, err = downloader.GetVideo(videoInfos, *destinationPath); nil != err {
		logrus.Errorf("Run `getVideo` error on url [%s], reason: %v", *videoSiteOrigin, err)
	}
	logrus.Infof("Downloaded video [%s] in file [%s]\n", *videoSiteOrigin, fileGenerated)

	// Make an audio file from the video to .mp3 format.
	if *audioOnly {
		makeAudioFile(fileGenerated)
	}

	finishTime := time.Now()
	delta := finishTime.Sub(startTime)
	logrus.Infof("\nDownload accomplished in %v\n", delta)
}

// Generate audio file from video file
func makeAudioFile(videoFilename string) {
	// Remove the video file
	defer func(fileToRemove string) {
		if errorRemovingFile := os.Remove(fileToRemove); nil != errorRemovingFile {
			logrus.Errorf("Could not remove origin video file [%s], reason: %v", fileToRemove, errorRemovingFile)
		}
	}(videoFilename)

	videoExtension := filepath.Ext(videoFilename)
	audioFilename := videoFilename[:len(videoFilename)-len(videoExtension)] + ".mp3"

	// The exec is synchrone.
	cmd := exec.Command("ffmpeg", "-i", videoFilename, audioFilename)
	logrus.Infof("Launch command :`%s`", strings.Join([]string{"ffmpeg", "-i", videoFilename, audioFilename}, " "))

	var out []byte
	var err error
	if out, err = cmd.CombinedOutput(); nil != err {
		logrus.Errorf("Error launching command to generate audio file: %v", err)
	} else {
		logrus.Infof("Generated audio file [%s]", audioFilename)
	}
	logrus.Debugf("\n\nffmpeg output:\n\n%s", out)
}
