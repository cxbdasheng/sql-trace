package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"sync"
)

const PathENV = "TRACE_CONFIG_FILE_PATH"

func GetConfigFilePathDefault() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "../.sql_trace_config.yaml"
	}
	return dir + string(os.PathSeparator) + ".sql_trace_config.yaml"
}
func GetDefaultPort() string {
	return ":8765"
}

// GetConfigFilePath 获得配置文件路径
func GetConfigFilePath() string {
	configFilePath := os.Getenv(PathENV)
	if configFilePath != "" {
		return configFilePath
	}
	return GetConfigFilePathDefault()
}

type Config struct {
	Settings
	// 语言
	Lang string
}

// ConfigCache ConfigCache
type cacheType struct {
	ConfigSingle *Config
	Err          error
	Lock         sync.Mutex
}

var cache = &cacheType{}

func GetConfigCached() (conf Config, err error) {
	cache.Lock.Lock()
	defer cache.Lock.Unlock()
	if cache.ConfigSingle != nil {
		return *cache.ConfigSingle, cache.Err
	}
	// init config
	cache.ConfigSingle = &Config{}
	configFilePath := GetConfigFilePath()
	_, err = os.Stat(configFilePath)
	if err != nil {
		cache.Err = err
		return *cache.ConfigSingle, err
	}

	byt, err := os.ReadFile(configFilePath)
	if err != nil {
		cache.Err = err
		return *cache.ConfigSingle, err
	}

	err = yaml.Unmarshal(byt, cache.ConfigSingle)
	if err != nil {
		cache.Err = err
		return *cache.ConfigSingle, err
	}

	// remove err
	cache.Err = nil
	return *cache.ConfigSingle, err
}

// SaveConfig 保存配置
func (conf *Config) SaveConfig() (err error) {
	cache.Lock.Lock()
	defer cache.Lock.Unlock()

	byt, err := yaml.Marshal(conf)
	if err != nil {
		log.Println(err)
		return err
	}

	configFilePath := GetConfigFilePath()
	err = os.WriteFile(configFilePath, byt, 0600)
	if err != nil {
		log.Println(err)
		return
	}

	// 清空配置缓存
	cache.ConfigSingle = nil

	return
}
func (conf *Config) GetPort() string {
	if conf.Settings.Port != "" {
		return ":" + conf.Settings.Port
	}
	return ":"
}
