package downloader_test

import (
	"testing"
	"video-downloader/downloader"

	"github.com/stretchr/testify/assert"
)

func TestGetVideo(t *testing.T) {
	// Wrong cases
	videoPath, err := downloader.GetVideo(nil, "")
	assert.Empty(t, videoPath)
	assert.EqualError(t, err, "Nil videoInfo passed in argument to 'getVideo'")

	vi := &downloader.VideoInfos{}
	videoPath, err = downloader.GetVideo(vi, "")
	assert.Empty(t, videoPath)
	assert.EqualError(t, err, "Empty video URL on 'getVideo'")

}
