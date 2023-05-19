package goco

import (
	"errors"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-ini/ini"
)

// Initialize Config which decide to use docker or ini file
func InitializeConfig(config interface{}) error {
	if strings.ToLower(os.Getenv("DOCKER")) == "true" {
		log.Println("Docker config enabled..")
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
func getDockerTag(config interface{}) error {
	var err error
	v := reflect.ValueOf(config).Elem()
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		dockerTagValue := os.Getenv(t.Field(i).Tag.Get("docker"))
		switch field.Kind() {
		case reflect.String:
			field.SetString(dockerTagValue)
		case reflect.Int:
			i, err := strconv.Atoi(dockerTagValue)
			if err != nil {
				return errors.New("Error while converting " + dockerTagValue + " to " + t.Field(i).Name)
			}
			field.SetInt(int64(i))
		case reflect.Bool:
			b, err := strconv.ParseBool(dockerTagValue)
			if err != nil {
				log.Println(errors.New("Error while converting " + dockerTagValue + " to " + t.Field(i).Name))
				b = false
			}
			field.SetBool(b)
		case reflect.Struct:
			err = getDockerTag(field.Addr().Interface())
			if err != nil {
				return err
			}
		}
	}
	return err
}
