package configloader

import (
	"fmt"
	"os"
	"video-downloader/parsingelement"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// ReadConfig is used to parse the configuration file and returns the error encountered if there is
func ReadConfig(configPath string) (*parsingelement.ParsingInformations, error) {
	var err error
	if "" == configPath {
		return nil, errors.New("No configuraiton path provided")
	}
	// Remote source is not managed yet. So we manage the fact the file must exist on the system
	if _, err = os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("Config file [%s] must exist on system", configPath)
	} else if nil != err {
		return nil, errors.Wrapf(err, "Could not get stat on file [%s]", configPath)
	}

	// Read config from file
	viper.SetConfigFile(configPath)
	if err = viper.ReadInConfig(); err != nil {
		return nil, errors.Wrapf(err, "Error reading configuration from file %s", configPath)
	}

	pi := parsingelement.ParsingInformations{}
	if err = viper.Unmarshal(&pi); nil != err {
		return nil, errors.Wrapf(err, "Error Unmarshalling configuration file %s", configPath)
	}

	// Check viper has populates structure field
	if !checkConfigurationLoaded(&pi) {
		return nil, errors.New("Viper hasn't populated struct from configuration file")
	}
	return &pi, nil
}

// checkConfigurationLoaded - Unmarshal maybe not worked without error.
// true = is ok | false = not ok
func checkConfigurationLoaded(pi *parsingelement.ParsingInformations) bool {
	// default fields value on creation
	if "" == pi.U.Delimiter && 0 == pi.U.UidSize && nil == pi.U.UrlsInfos && "" == pi.U.QueryKeywordURL {
		return false
	}
	return true
}
