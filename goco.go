package goco

import (
	"os"

	"github.com/go-ini/ini"
)

var Config interface{}

func InitializeConfig(filePaths ...string) error {
	filePath := "config.ini"
	if len(filePaths) != 0 {
		filePath = filePaths[0]
	}

	configFile, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	return ini.MapTo(&Config, configFile)
}
