package goco

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-ini/ini"
	"github.com/joho/godotenv"
)

// Debug Mode
var DebugMode bool = false
var ConfigMode string = "ini"

// Initialize Config which decide to use docker or ini file
// parameters: DebugMode bool = false, ConfigMode string(ini, env, docker) = ini
func InitializeConfig(config interface{}, parameters ...any) error {
	// Check if debug is passed
	if len(parameters) == 1 {
		DebugMode = parameters[0].(bool)
		debugLog("Debug Mode Enabled..")
	}

	if len(parameters) == 2 {
		ConfigMode = parameters[1].(string)
	} else {
		// Check if DOCKER env is set to true
		if strings.ToLower(os.Getenv("DOCKER")) == "true" {
			debugLog("Docker config enabled..")
			ConfigMode = "docker"
		}

		// Check if config.ini file is present
		if _, err := os.Stat("config.ini"); err == nil {
			debugLog("Config.ini file found..")
			ConfigMode = "ini"
		}

		// Check if .env file is present
		if _, err := os.Stat(".env"); err == nil {
			debugLog(".env file found..")
			ConfigMode = "env"
		}
	}

	switch ConfigMode {
	case "ini":
		return InitializeIniConfig(config)
	case "env":
		return InitializeEnvironmentConfig(config)
	case "docker":
		return InitializeDockerConfig(config)
	default:
		return errors.New("Invalid Config Mode")
	}
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

// Initialize Config from environment file
func InitializeEnvironmentConfig(config interface{}, filePaths ...string) error {
	// Load .env file if provided
	filePath := ".env"
	if len(filePaths) != 0 {
		filePath = filePaths[0]
	}
	if err := godotenv.Load(filePath); err != nil {
		return fmt.Errorf("failed to load env file: %w", err)
	}

	return getEnvironmentTag(config)
}

// Get Environment Tag
func getEnvironmentTag(config interface{}, structTag ...string) error {
	var err error

	// Populate struct fields with environment variables
	v := reflect.ValueOf(config)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return errors.New("config parameter must be a pointer to a struct")
	}
	v = v.Elem()

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		envTag := t.Field(i).Tag.Get("env")

		// Check if structTag is passed
		if len(structTag) != 0 {
			envTag = fmt.Sprintf("%s_%s", structTag[0], envTag)
		}

		envTagValue := os.Getenv(envTag)
		debugLog("envTag:\t", envTag, "\t=>\tenvTagValue:\t", envTagValue)

		switch field.Kind() {
		case reflect.String:
			field.SetString(envTagValue)
		case reflect.Int:
			intValue, err := strconv.Atoi(envTagValue)
			if err != nil {
				return fmt.Errorf("failed to convert %s to int for field %s: %w", envTagValue, fieldType.Name, err)
			}
			field.SetInt(int64(intValue))
		case reflect.Bool:
			boolValue, err := strconv.ParseBool(envTagValue)
			if err != nil {
				return fmt.Errorf("failed to convert %s to bool for field %s: %w", envTagValue, fieldType.Name, err)
			}
			field.SetBool(boolValue)
		case reflect.Struct:
			if err := getEnvironmentTag(field.Addr().Interface()); err != nil {
				return fmt.Errorf("failed to initialize nested struct %s: %w", fieldType.Name, err)
			}
		default:
			return fmt.Errorf("unsupported field type %s for field %s", field.Kind(), fieldType.Name)
		}
	}

	return err
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
		fieldType := t.Field(i)
		dockerTag := t.Field(i).Tag.Get("docker")

		// Check if structTag is passed
		if len(structTag) != 0 {
			dockerTag = fmt.Sprintf("%s_%s", structTag[0], dockerTag)
		}

		dockerTagValue := os.Getenv(dockerTag)
		debugLog("dockerTag:\t", dockerTag, "\t=>\tdockerTagValue:\t", dockerTagValue)

		switch field.Kind() {
		case reflect.String:
			field.SetString(dockerTagValue)
		case reflect.Int:
			intValue, err := strconv.Atoi(dockerTagValue)
			if err != nil {
				err = errors.New("Error while converting " + dockerTag + " to " + fieldType.Name + "\t" + err.Error())
				debugLog(err)
				return err
			}
			field.SetInt(int64(intValue))
		case reflect.Bool:
			boolValue, err := strconv.ParseBool(dockerTagValue)
			if err != nil {
				err = errors.New("Error while converting " + dockerTag + " to " + fieldType.Name + "\t" + err.Error())
				debugLog(err)
				boolValue = false
			}
			field.SetBool(boolValue)
		case reflect.Struct:
			err = getDockerTag(field.Addr().Interface(), dockerTag)
			if err != nil {
				err = errors.New("Error while getting docker tag for " + fieldType.Name + "\t" + err.Error())
				debugLog(err)
				return err
			}
		}
	}
	return err
}

// Debug Log
func debugLog(message ...any) {
	if DebugMode {
		fmt.Println(message...)
	}
}
