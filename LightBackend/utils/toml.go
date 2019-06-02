package utils

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
)


type HotUpdate interface {
	UpdateConfig()
}

type ConfigManager struct {
	HotUpdateTarget []HotUpdate
}

var DefaultToml *DefaultTomlConfig
var DefaultConfigManager *ConfigManager

func init(){
	DefaultToml = &DefaultTomlConfig{}
	DefaultConfigManager = &ConfigManager{}
}

func LoadTomlConfig(localConfig interface{}, filePath string) {
	// load default config for global
	configPath, err := os.Getwd()
	var path string
	if err != nil {
		fmt.Printf("Read project dir error:%v \n", err)
		os.Exit(-1)
	}
	if filePath == "" {
		path = configPath + CONFIGPATH
	} else {
		path = filePath
	}
	_, err = toml.DecodeFile(path, DefaultToml)
	if err != nil {
		fmt.Printf("load config file error:%v \n", err)
		os.Exit(-1)
	}
	fmt.Println("Load config file success.")

	// load input config
	if localConfig != nil {
		_, err = toml.DecodeFile(CONFIGPATH, localConfig)
		if err != nil {
			fmt.Printf("load config file error:%v",err)
			os.Exit(-1)
		}
	}

}

func (cm *ConfigManager) HotUpdateConfig() {
	LoadTomlConfig(nil,"")
	for _, item := range cm.HotUpdateTarget {
		item.(HotUpdate).UpdateConfig()
	}
}
