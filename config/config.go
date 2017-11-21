package config

import (
	"flag"
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
	Port     uint   `toml:"port"`
	Version  string `toml:"version"`
	Debug    bool   `toml:"debug"`
	Endpoint string `toml:"endpoint"`
	Jwt      string `toml:"jwt"`
}

// DBConfig MongoDB設定構造体
type DBConfig struct {
	Server   string `toml:"server"`
	Database string `toml:"database"`
	User     string `toml:"user"`
	Password string `toml:"password"`
}

type CacheConfig struct {
	Server   string `toml:"server"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
}

type UploadImageConfig struct {
	Path string `toml:"path"`
}

const (
	MockJwtToken = "token"
)

// GetConfig TOML設定ファイルから設定を取得
func GetConfig() Config {
	if flag.Lookup("test.v") != nil {
		mockAPIConfig := APIConfig{
			Port:     8080,
			Version:  "v1",
			Debug:    true,
			Endpoint: "",
			Jwt:      MockJwtToken,
		}
		mockDBConfig := DBConfig{
			Server:   "localhost:27017",
			Database: "timeline",
		}
		mockCacheConfig := CacheConfig{
			Server: "localhost",
			Port:   6379,
		}
		mockUploadImage := UploadImageConfig{
			Path: "",
		}
		return Config{
			API:         mockAPIConfig,
			DB:          mockDBConfig,
			Cache:       mockCacheConfig,
			UploadImage: mockUploadImage,
		}
	}

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
