package configuration

import (
	"encoding/json"
	"io/ioutil"
)

// readFile reads the bytes from a given file
func readFile(cfgFile string) ([]byte, error) {
	source, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return nil, err
	}

	return source, nil
}

// readConfig tries to parse the byte data into a config file
func readConfig(fileContent []byte) (Config, error) {
	config := Config{}
	err := json.Unmarshal(fileContent, &config)
	return config, err
}

// GetConfig retrieves a new configuration from the given config file
// returns an error when config does not exist or cannot be read or validated
func GetConfig(cfgFile string) (Config, error) {
	content, err := readFile(cfgFile)
	if err != nil {
		return Config{}, err
	}

	conf, err := readConfig(content)
	if err != nil {
		return Config{}, err
	}

	err = conf.Validate()
	if err != nil {
		return Config{}, err
	}

	return conf, nil
}
