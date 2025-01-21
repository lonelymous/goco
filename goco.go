package goco

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-ini/ini"
)

// Debug Mode
var debugMode bool = false

// Initialize Config which decide to use docker or ini file
func InitializeConfig(config interface{}, debug ...bool) error {
	// Check if debug is passed
	if len(debug) != 0 {
		debugMode = debug[0]
	}

	// Check if DOCKER env is set to true
	if strings.ToLower(os.Getenv("DOCKER")) == "true" {
		debugLog("Docker config enabled..")
		return InitializeDockerConfig(config)
	}

	return InitializeIniConfig(config)
}

// Initialize Config from ini file
func InitializeIniConfig(config interface{}, filePaths ...string) error {
	filePath := "config.ini"
	if len(filePaths) != 0 {
		filePath = filePaths[0]
	}

	configFile, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	return ini.MapTo(config, configFile)
}

// Initialize Config from docker with docker tags
func InitializeDockerConfig(config interface{}) error {
	return getDockerTag(config)
}

// Get Docker Tag
func getDockerTag(config interface{}, structTag ...string) error {
	var err error
	v := reflect.ValueOf(config).Elem()
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		dockerTag := t.Field(i).Tag.Get("docker")

		// Check if structTag is passed
		if len(structTag) != 0 {
			dockerTag = fmt.Sprintf("%s_%s", structTag[0], t.Field(i).Name)
		}

		dockerTagValue := os.Getenv(dockerTag)
		debugLog("dockerTag:\t", dockerTag, "\t=>\tdockerTagValue:\t", dockerTagValue)

		switch field.Kind() {
		case reflect.String:
			field.SetString(dockerTagValue)
		case reflect.Int:
			i, err := strconv.Atoi(dockerTagValue)
			if err != nil {
				err = errors.New("Error while converting " + dockerTagValue + " to " + t.Field(i).Name + "\t" + err.Error())
				debugLog(err)
				return err
			}
			field.SetInt(int64(i))
		case reflect.Bool:
			b, err := strconv.ParseBool(dockerTagValue)
			if err != nil {
				err = errors.New("Error while converting " + dockerTagValue + " to " + t.Field(i).Name + "\t" + err.Error())
				debugLog(err)
				b = false
			}
			field.SetBool(b)
		case reflect.Struct:
			err = getDockerTag(field.Addr().Interface(), dockerTag)
			if err != nil {
				err = errors.New("Error while getting docker tag for " + t.Field(i).Name + "\t" + err.Error())
				debugLog(err)
				return err
			}
		}
	}
	return err
}

// Debug Log
func debugLog(message ...any) {
	if debugMode {
		fmt.Println(message...)
	}
}
