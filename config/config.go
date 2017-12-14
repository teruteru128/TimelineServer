package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

// Config 設定構造体
type Config struct {
	API         APIConfig
	DB          DBConfig
	Cache       CacheConfig
	UploadImage UploadImageConfig
}

// APIConfig API設定構造体
type APIConfig struct {
	Port     int    `toml:"port"`
	Version  string `toml:"version"`
	Debug    bool   `toml:"debug"`
	Endpoint string `toml:"endpoint"`
	Secure   bool   `toml:"secure"`
	Jwt      string `toml:"jwt"`
}

// DBConfig MongoDB設定構造体
type DBConfig struct {
	Server   string `toml:"server"`
	Database string `toml:"database"`
}

type CacheConfig struct {
	Server string `toml:"server"`
}

type UploadImageConfig struct {
	Path string `toml:"path"`
}

const (
	MockJwtToken = "token"
)

// GetConfig TOML設定ファイルから設定を取得
func GetConfig() Config {
	var config Config
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		log.Fatal(err)
	}

	return config
}

// GetAPIConfig TOML設定ファイルからAPI設定を取得
func GetAPIConfig() APIConfig {
	baseConfig := GetConfig()
	return baseConfig.API
}

// GetDBConfig TOML設定ファイルからDB接続情報を取得
func GetDBConfig() DBConfig {
	baseConfig := GetConfig()
	return baseConfig.DB
}

func GetCacheConfig() CacheConfig {
	baseConfig := GetConfig()
	return baseConfig.Cache
}

func GetUploadImagePath() string {
	baseConfig := GetConfig().UploadImage
	return baseConfig.Path
}
