package goco

import (
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-ini/ini"
)

// Initialize Config which decide to use docker or ini file
func InitializeConfig(config interface{}) error {
	if strings.ToLower(os.Getenv("DOCKER")) == "true" {
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
	v := reflect.ValueOf(config).Elem() // Get the struct value using reflection
	t := v.Type()                       // Get the struct type
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		dockerTag := os.Getenv(t.Field(i).Tag.Get("docker"))

		switch field.Kind() {
		case reflect.String:
			field.SetString(dockerTag)
		case reflect.Int:
			i, err := strconv.Atoi(dockerTag)
			if err != nil {
				return err
			}
			field.SetInt(int64(i))
		case reflect.Bool:
			b, err := strconv.ParseBool(dockerTag)
			if err != nil {
				return err
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
