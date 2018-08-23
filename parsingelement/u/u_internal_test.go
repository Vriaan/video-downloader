package u

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetExtensionFromType(t *testing.T) {
	type expectedResult struct {
		err       error
		extension string
	}

	cases := map[string]expectedResult{
		"video/mp4; good": expectedResult{
			extension: ".mp4",
		},
		"video/x-msvideo; good": expectedResult{
			extension: ".avi",
		},
		"video/avi; wrong": expectedResult{
			err: errors.New("Type information contains no mime-type supported [video/avi; wrong]"),
		},
		"ouch; wrong": expectedResult{
			err: errors.New("Type information contains no mime-type supported [ouch; wrong]"),
		},
		"": expectedResult{
			err: errors.New("Type information contains no mime-type supported []"),
		},
	}

	for typeValue, expected := range cases {
		ext, err := getExtensionFromType(typeValue)

		assert.Equal(t, expected.extension, ext)
		if nil == expected.err {
			assert.Nil(t, err)
		} else {
			assert.EqualError(t, err, expected.err.Error())
		}
	}
}

func TestFindVideoID(t *testing.T) {
	UTestInstance := &U{Delimiter: "zzzzzzz=", UidSize: 3}
	type expectedResult struct {
		err     error
		videoID string
	}

	cases := map[string]expectedResult{
		"": expectedResult{
			videoID: "",
			err:     errors.New("Empty URL given for `findVideoID`"),
		},
		"http://testurl": expectedResult{
			videoID: "",
			err:     fmt.Errorf("No delimiter [%s] found for within [http://testurl] for `findVideoID`", UTestInstance.Delimiter),
		},
		"http://testurl" + UTestInstance.Delimiter + "qwerty": expectedResult{
			videoID: "qwe",
		},
		"http://testurl" + UTestInstance.Delimiter + "qw": expectedResult{
			err: errors.New("Video ID size [3] is bigger than the URL piece !"),
		},
	}

	for urlCase, expected := range cases {
		id, err := UTestInstance.findVideoID(urlCase)

		assert.Equal(t, expected.videoID, id)
		if nil == expected.err {
			assert.Nil(t, err)
		} else {
			assert.EqualError(t, err, expected.err.Error())
		}
	}

	// No ID Uid_size
	UTestInstance.UidSize = 0
	urlTest := "http://testurl" + UTestInstance.Delimiter + "qwerty"
	id, err := UTestInstance.findVideoID(urlTest)
	assert.Equal(t, "", id)
	assert.EqualError(t, err, "Empty ID size is <= 0 for `findVideoID`")

	// No delimiter set
	urlTest = "http://testurl"
	UTestInstance.UidSize = 3
	UTestInstance.Delimiter = ""
	id, err = UTestInstance.findVideoID(urlTest)

	assert.Equal(t, "", id)
	assert.EqualError(t, err, "Empty delimiter for `findVideoID`")
}
