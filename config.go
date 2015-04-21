package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Ui   UiSection   `toml:"ui"`
	Core CoreSection `toml:"core"`
}

type UiSection struct {
	Editor string `toml:"editor"`
}

type CoreSection struct {
	Linefeed string `toml:"linefeed"`
}

func SetDefaultConfig(config *Config) {
	if len(config.Ui.Editor) <= 0 {
		config.Ui.Editor = "more"
	}
	if len(config.Core.Linefeed) <= 0 {
		config.Core.Linefeed = "unix"
	}
}

func LoadConfig() (Config, error) {
	var config Config
	var err error
	home, err := GetHomePath()
	if err == nil {
		configFile := filepath.Join(home, ".gookmarkrc")
		_, err = os.Stat(configFile)
		if err == nil {
			_, err = toml.DecodeFile(configFile, &config)
		}
		SetDefaultConfig(&config)
		err = nil
	}
	return config, err
}

type ConfigOption struct {
	Section  string
	Property string
	Value    string
}

func ParseConfigOption(option string) (ConfigOption, error) {
	var err error
	var buff []string

	var configOption ConfigOption

	buff = strings.SplitN(option, ".", 2)
	if len(buff) != 2 {
		err = errors.New("Invalid option")
	} else {
		configOption.Section = buff[0]
		buff = strings.SplitN(buff[1], "=", 2)
		if len(buff) != 2 {
			err = errors.New("Invalid option")
		} else {
			configOption.Property = buff[0]
			configOption.Value = buff[1]
		}
	}

	return configOption, err
}

func WriteNewConfig(config *Config, newConfig string) error {
	configOption, err := ParseConfigOption(newConfig)
	if err != nil {
		return err
	}

	if configOption.Section == "ui" {
		if configOption.Property == "editor" {
			config.Ui.Editor = configOption.Value
		} else {
			err = errors.New("Unknown property: " + configOption.Section + "." + configOption.Property)
		}
	} else if configOption.Section == "core" {
		if configOption.Property == "linefeed" {
			config.Core.Linefeed = configOption.Value
		} else {
			err = errors.New("Unknown property: " + configOption.Section + "." + configOption.Property)
		}
	} else {
		err = errors.New("Unknown section: " + configOption.Section)
	}
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	err = toml.NewEncoder(&buffer).Encode(config)
	if err != nil {
		return err
	}

	home, err := GetHomePath()
	if err != nil {
		return err
	}
	configFile := filepath.Join(home, ".gookmarkrc")
	fp, err := os.OpenFile(configFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer fp.Close()

	fmt.Fprintln(fp, buffer.String())

	return err
}
