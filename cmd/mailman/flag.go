package main

import (
	"flag"
	"strings"
)

type argsConfig struct {
	configFiles []string
	logLevel    string
}

func commandLineArgsConfig() argsConfig {
	configFiles := flag.String("config-file", "",
		"Comma-separated paths to configuration file(s), the last one has the highest priority")
	logLevel := flag.String("log-level", "INFO", "Logging level: DEBUG, INFO, WARN, ERROR or FATAL")

	flag.Parse()

	return argsConfig{
		configFiles: strings.Split(*configFiles, ","),
		logLevel:    *logLevel,
	}
}
