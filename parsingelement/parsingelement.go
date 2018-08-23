package parsingelement

import (
	"fmt"
	"strings"
	"video-downloader/downloader"
	"video-downloader/parsingelement/u"
)

// This struct contains all information parsing
type ParsingInformations struct {
	U *u.U `mapstructure:"u"`
}

// ParseOn Use the right parser for the page
func (pi *ParsingInformations) ParseOn(site string, url string) (vi *downloader.VideoInfos, err error) {

	switch strings.ToLower(site) {
	case "u":
		vi, err = pi.U.Parse(url)
	default:
		err = fmt.Errorf("No [%s] site available", site)
	}

	return
}
