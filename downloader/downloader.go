package downloader

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type VideoInfos struct {
	URL       string
	Extension string
	Title     string
}

// Based on information filled, attempt to dynamically create video filename
func generateVideoFilename(vi *VideoInfos) (string, error) {
	if nil == vi {
		return "", errors.New("Uninitalized parameters provided")
	}

	if vi.Title != "" && vi.Extension != "" {
		reg, err := regexp.Compile("[^a-zA-Z0-9_ ]+")
		if err != nil {
			log.Fatal(err)
		}
		replaceCharacterForbidden := reg.ReplaceAllString(vi.Title, "")
		slicedTitle := strings.Split(strings.ToLower(replaceCharacterForbidden), " ")
		if strings.HasPrefix(vi.Extension, ".") {
			return strings.Join(slicedTitle, "_") + vi.Extension, nil
		} else {
			return strings.Join(slicedTitle, "_") + "." + vi.Extension, nil
		}
	}

	return "", errors.New("No enough information to create video filename")
}

// GetVideo - effectively download the video once we got the right videoInfo
func GetVideo(vi *VideoInfos, destinationPath string) (string, error) {
	if nil == vi {
		return "", errors.New("Nil videoInfo passed in argument to 'getVideo'")
	}
	videoURL := vi.URL
	if "" == videoURL {
		return "", errors.New("Empty video URL on 'getVideo'")
	}

	resp, err := http.Get(videoURL)

	if nil != err {
		return "", errors.Wrapf(err, "Error fetching content [%s]", videoURL)
	}
	defer resp.Body.Close()

	var content []byte
	if content, err = ioutil.ReadAll(resp.Body); nil != err {
		return "", errors.Wrapf(err, "Error reading response of [%s]", videoURL)
	}

	var file *os.File
	var filename string
	if filename, err = generateVideoFilename(vi); nil != err {
		return "", errors.Wrap(err, "Error generating filename")
	}
	logrus.Debugf("generate file name [%s]", filename)
	file, err = os.OpenFile(filepath.Join(destinationPath, filename), os.O_CREATE|os.O_RDWR, 0777)
	if nil != err {
		return "", errors.Wrapf(err, "Error creating destination file for [%s]", videoURL)
	}
	defer file.Close()
	var totalWritten int
	if totalWritten, err = file.Write(content); nil != err || len(content) != totalWritten {
		return "", errors.Wrapf(err, "Error writing to destination file for [%s]", videoURL)
	}

	return file.Name(), nil
}
