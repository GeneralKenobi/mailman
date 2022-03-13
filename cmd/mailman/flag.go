package main

import (
	"flag"
	"strings"
)

type argsConfig struct {
	configFiles []string
}

func commandLineArgsConfig() argsConfig {
	configFiles := flag.String("config-file", "",
		"Comma-separated paths to configuration file(s), the last one has the highest priority")

	flag.Parse()

	return argsConfig{
		configFiles: strings.Split(*configFiles, ","),
	}
}
