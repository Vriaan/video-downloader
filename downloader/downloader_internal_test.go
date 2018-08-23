package downloader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateVideoFilename(t *testing.T) {

	// Wrong cases
	filename, err := generateVideoFilename(nil)
	assert.Empty(t, filename)
	assert.EqualError(t, err, "Uninitalized parameters provided")

	vi := &VideoInfos{}
	filename, err = generateVideoFilename(vi)
	assert.Empty(t, filename)
	assert.EqualError(t, err, "No enough information to create video filename")

	// Right case
	vi.Extension = ".mp4"
	vi.Title = "I ve got friends yupi"
	filename, err = generateVideoFilename(vi)
	assert.Equal(t, "i_ve_got_friends_yupi.mp4", filename)
	assert.Nil(t, err)

	vi.Extension = "mp4"
	filename, err = generateVideoFilename(vi)
	assert.Equal(t, "i_ve_got_friends_yupi.mp4", filename)
	assert.Nil(t, err)
}
