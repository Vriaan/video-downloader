package u

import (
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"strings"
	"video-downloader/downloader"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// This struct contains informations to parse a specific web page in order to get video (but which one ?)
// Those are fixed information. The structure can be access concurrently by many goroutine
type U struct {
	UrlsInfos       []string `mapstructure:"get_video_info_urls"`
	Delimiter       string   `mapstructure:"url_Delimiter"`
	UidSize         int      `mapstructure:"id_number_character"`
	QueryKeywordURL string   `mapstructure:"query_key_url"`
}

// Parse - Parse U page to get video. Can be used concurrently
func (u *U) Parse(url string) (*downloader.VideoInfos, error) {
	videoID, errConf := u.findVideoID(url)
	if nil != errConf {
		return nil, errors.Wrapf(errConf, "Run `findVideoID` error on url [%s]", url)
	}
	logrus.Debugf("Resolved findVideoID [%s] on url [%s]\n", videoID, url)

	videoInfo, err := u.parseVideoInfo(videoID)
	if nil != err {
		return nil, errors.Wrapf(err, "Run `parseVideoInfo` error on url [%s]", url)
	}
	logrus.Debugf("Parsed video information, obtained %v", videoInfo)
	return videoInfo, nil
}

// findVideoID - Used to get the UID of the video
func (u *U) findVideoID(url string) (string, error) {
	if "" == url {
		return "", errors.New("Empty URL given for `findVideoID`")
	}
	if "" == u.Delimiter {
		return "", errors.New("Empty delimiter for `findVideoID`")
	}
	if u.UidSize <= 0 {
		return "", errors.New("Empty ID size is <= 0 for `findVideoID`")
	}
	piecesOfURL := strings.Split(url, u.Delimiter)
	if len(piecesOfURL) < 2 {
		return "", fmt.Errorf("No delimiter [%s] found for within [%s] for `findVideoID`", u.Delimiter, url)
	}
	partContainingUID := piecesOfURL[1]
	if len(partContainingUID) < u.UidSize {
		return "", fmt.Errorf("Video ID size [%d] is bigger than the URL piece !", u.UidSize)
	}
	return partContainingUID[:u.UidSize], nil
}

func (u *U) fetchVideoInfoURL(urlInfo string) (url.Values, error) {
	resp, err := http.Get(urlInfo)
	if nil != err {
		return nil, errors.Wrapf(err, "Error fetching video infos [%s]", urlInfo)
	}
	defer resp.Body.Close()

	var content []byte
	if content, err = ioutil.ReadAll(resp.Body); nil != err {
		return nil, errors.Wrapf(err, "Error reading response of video infos [%s]", urlInfo)
	}
	videoFileInfoContent := string(content)

	// according to url doc, is a map[string][]string
	var query url.Values
	query, err = url.ParseQuery(videoFileInfoContent)
	if nil != err {
		return nil, errors.Wrapf(err, "Error parsing query from response [%s]", urlInfo)
	}

	return query, nil
}

// parseVideoInfo - Used to parse the information return by the forged URL that returns video information
func (u *U) parseVideoInfo(videoID string) (*downloader.VideoInfos, error) {
	var videoInfoURL string
	var query url.Values
	var err error

	// Try all video until having the information we want
	var foundURL bool
	var infosURLEncoded []string
	for _, videoInfoURL = range u.UrlsInfos {
		query, err = u.fetchVideoInfoURL(videoInfoURL + videoID)
		if infosURLEncoded, foundURL = query[u.QueryKeywordURL]; !foundURL {
			return nil, errors.Wrapf(err, "Error no '%s' key on query", u.QueryKeywordURL)
		}

		// The url didn't worked, try next one
		if len(infosURLEncoded) < 1 || infosURLEncoded[0] == "" {
			infosURLEncoded = nil
			logrus.Debugf("Info URL did not work with [%s], try next", videoInfoURL)
			continue
		}
		break
	}

	// None were found, need to trigger an error
	if nil == infosURLEncoded {
		return nil, errors.Wrapf(err, "Error no '%s' infos on key found from video information", u.QueryKeywordURL)
	}

	infosParsed, err := url.ParseQuery(infosURLEncoded[0])
	if nil != err {
		return nil, errors.Wrapf(err, "Error parsing '%s' from [%s]", u.QueryKeywordURL, videoInfoURL)
	}

	// For + switch because I might want more infos in the future
	extractedInfos := &downloader.VideoInfos{}
	for _, info := range []string{"url", "type"} {
		infoExtracted, infoExist := infosParsed[info]

		if !infoExist {
			return nil, fmt.Errorf("Error no url [%s] encountered for [%s]", info, videoInfoURL)
		}
		if len(infoExtracted) < 1 {
			return nil, fmt.Errorf("Error no info [%s] encountered for [%s]", info, videoInfoURL)
		}

		switch info {
		case "url":
			extractedInfos.URL = infoExtracted[0]
		case "type":
			extractedInfos.Extension, err = getExtensionFromType(infoExtracted[0])
			if nil != err {
				return nil, errors.Wrapf(err, "Error on info [%s] encountered for [%s]", info, videoInfoURL)
			}
		}
	}

	titleSlice, okTitle := query["title"]
	if !okTitle {
		return nil, errors.Wrap(err, "Error no 'title' key on query")
	}
	if len(titleSlice) < 1 {
		return nil, errors.Wrap(err, "Error no 'title'  value encountered")
	}
	extractedInfos.Title = titleSlice[0]

	return extractedInfos, err
}

// getExtensionFromType Parse the type value returned and try to gets the extension from it
func getExtensionFromType(typeValue string) (string, error) {
	parts := strings.Split(typeValue, "; ")
	extension, err := mime.ExtensionsByType(parts[0])
	if nil != extension && len(extension) > 0 {
		// Many extension may be possible, we just use the 1st one
		return extension[0], err
	} else {
		return "", fmt.Errorf("Type information contains no mime-type supported [%s]", typeValue)
	}

	return "", errors.Wrapf(err, "Error trying to get mime-type from [%s]", typeValue)
}
