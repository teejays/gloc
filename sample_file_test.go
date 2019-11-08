package main

var sampleFileA = `package config

import (
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/teejays/clog"
)

// Config defines the structure of the configuration file for the application.
type Config struct {
	InQueueBufferSize     int               
	AppLogSupressionLevel int               
	LogSources            []ConfigLogSource 
	Stats                 struct {
		Types []ConfigStatsType
	}
	Alert struct {
		Types []ConfigAlertType
	}
}

// ConfigLogSource is information from the config file regarding LogSources that the application need to use.
type ConfigLogSource struct {
	Name     string
	Type     string
	Path     string
	Disabled bool
	Settings ConfigLogSourceSettings
}

type ConfigLogSourceSettings struct {
	Format               string
	Headers              []string
	TimestampKey         string 
	TimestampFormat      string 
	UseFirstlineAsHeader bool   
}

type ConfigStatsType struct {
	Name            string
	DurationSeconds int64 
	Disabled        bool
	SourceSettings  []ConfigStatsTypeSourceSetting 
}

type ConfigStatsTypeSourceSetting struct {
	Name                string
	Key                 string
	ValueMutateFuncName string 
	OtherKeys           []string
}

type ConfigAlertType struct {
	Name            string
	DurationSeconds int64
	Threshold       int
	Disabled        bool
	SourceSettings  []ConfigAlertTypeSourceSetting
}

type ConfigAlertTypeSourceSetting struct {
	Name                string
	Key                 string
	ValueMutateFuncName string
	Values              []string
}

// ReadConfigTOML takes a path to a config file in TOML format, and parses it into a Config struct
func ReadConfigTOML(path string) (Config, error) {
	var cfg Config

	if strings.TrimSpace(path) == "" {
		return cfg, fmt.Errorf("empty config file path")
	}

	_, err := toml.DecodeFile(path, &cfg)
	if err != nil {
		return cfg, err
	}

	// Do some validation
	if cfg.InQueueBufferSize < 1 {
		return cfg, fmt.Errorf("InQueueBufferSize has an invalid value: %d", cfg.InQueueBufferSize)
	}

	clog.LogLevel = cfg.AppLogSupressionLevel

	clog.Debugf("Config Log Sources: %v", cfg.LogSources)
	if len(cfg.LogSources) < 1 {
		return cfg, fmt.Errorf("no log source detected from config file")
	}

	return cfg, nil
}`
