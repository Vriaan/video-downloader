package configloader_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"video-downloader/configloader"
	"video-downloader/parsingelement"
	"video-downloader/parsingelement/u"

	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	var pi *parsingelement.ParsingInformations
	var err error

	// Wrong case - Path is empty
	pi, err = configloader.ReadConfig("")
	assert.Nil(t, pi)
	assert.EqualError(t, err, "No configuraiton path provided")

	// Wrong case - path does not exist
	pathDoesNotExist := "/tmp/pathDoesNotExist/RE"
	pi, err = configloader.ReadConfig(pathDoesNotExist)
	assert.Nil(t, pi)
	assert.EqualError(t, err, fmt.Sprintf("Config file [%s] must exist on system", pathDoesNotExist))

	// Wrong case - file configuration is empty -> viper Unmarshal silently fails test
	pathToEmptyConf := "/tmp/pathToEmptyConf/"
	fileConfName := "conf.yml"

	if err = os.MkdirAll(pathToEmptyConf, 0777); nil != err {
		t.Errorf("Could not create temporary path [%s]", pathToEmptyConf)
	}
	var confFileHandler *os.File
	confFilePath := filepath.Join(pathToEmptyConf, fileConfName)
	if confFileHandler, err = os.OpenFile(confFilePath, os.O_CREATE|os.O_RDWR, 0777); nil != err {
		t.Errorf("Could not create temporary conf file [%s]", confFilePath)
	}
	defer confFileHandler.Close()
	defer os.RemoveAll(pathToEmptyConf)

	pi, err = configloader.ReadConfig(confFilePath)
	assert.Nil(t, pi)
	assert.EqualError(t, err, "Viper hasn't populated struct from configuration file")

	conf := `---
u:
    url_delimiter: "uNQS1!"
    id_number_character: 1
    get_video_info_urls:
        - "h"
        - "ht"
        - "htt"
        - "http"
    query_key_url: "?:lic!"`

	confFileHandler.WriteString(conf)
	pi, err = configloader.ReadConfig(confFilePath)

	assert.Equal(t, &parsingelement.ParsingInformations{
		&u.U{
			UrlsInfos:       []string{"h", "ht", "htt", "http"},
			Delimiter:       "uNQS1!",
			Uid_size:        1,
			QueryKeywordURL: "?:lic!",
		},
	}, pi)
	assert.Nil(t, err)
}
