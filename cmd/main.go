package main

import (
	"os"
	"time"
	"video-downloader/configloader"
	"video-downloader/downloader"
	"video-downloader/parsingelement"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

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

var configurationPath = pflag.StringP("configuration", "c", "", "To run the program needs to get some configuration")
var loggerLvl = pflag.String("logLvl", "InfoLevel", "The level of log to show [default = InfoLevel]. Available are (PanicLevel, FatalLevel, ErrorLevel, WarnLevel, InfoLevel, DebugLevel). For more information, llo, at Logrus doc 'type level'")

//var fileURLDownload = pflag.StringP("file_infos", "f", "", "The program will try to download all video from this file")
var videoURlDownload = pflag.StringP("url", "u", "", "The program will try to download the video on this url")
var videoSiteOrigin = pflag.StringP("origin", "o", "", "The program will parse according to the origin web site you gave")
var destinationPath = pflag.StringP("destination", "d", "", "The program will write every final video in this directory")

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

	finishTime := time.Now()
	delta := finishTime.Sub(startTime)
	logrus.Infof("\nDownload accomplished in %v\n", delta)
}
