package config

import (
	"encoding/json"
	"fmt"
	"github.com/GeneralKenobi/mailman/pkg/mdctx"
	"io/ioutil"
)

// Load prepares configuration. It starts with the default one as template (see default.go) and reads each file given in configFiles.
// When applying a configuration file the properties defined in it overwrite the current state. The undefined properties are ignored and
// don't overwrite the current state. The last config file has the highest priority (since it's the last one to be read and applied).
//
// Configuration files have to be in JSON format and have to match the structs defined in types.go.
// Typically, there are 2 configuration files - main configuration file with non-sensitive properties (e.g. HTTP port to listen on) and
// secret configuration file with sensitive data (e.g. DB credentials).
//
// All given configuration files have to be loaded successfully, otherwise an error.
func Load(configFiles []string) error {
	cfg := defaultConfig

	for _, configFile := range configFiles {
		configBytes, err := fileReadHook(configFile)
		if err != nil {
			return fmt.Errorf("error reading configuration file %q: %w", configFile, err)
		}

		// cfg already contains values, unmarshalling to it will partially overwrite them with the properties defined in the file
		err = json.Unmarshal(configBytes, &cfg)
		if err != nil {
			return fmt.Errorf("error unmarshaling json from configuration file %q: %w", configFile, err)
		}
	}

	mdctx.Debugf(nil, "Configuration loaded from files %v: %#v", configFiles, cfg)
	currentConfig = cfg
	return nil
}

// Get returns the active configuration. The first active configuration is the default one. Then it's replaced with the configuration loaded
// in Load.
func Get() Config {
	return currentConfig
}

var currentConfig = defaultConfig

// For mocking in unit tests.
var fileReadHook = ioutil.ReadFile
