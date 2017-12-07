package config

import (
	"flag"
	"log"
	"os"
	"strconv"

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
	Env      string `toml:"env"`
}

// DBConfig MongoDB設定構造体
type DBConfig struct {
	Server   string `toml:"server"`
	Database string `toml:"database"`
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
			Version:  "1.0",
			Debug:    true,
			Endpoint: "tlstag.ddns.net",
			Jwt:      MockJwtToken,
			Secure:   false,
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

	if config.API.Env == "heroku" {
		port, _ := strconv.Atoi(os.Getenv("PORT"))
		redisPort, _ := strconv.Atoi(os.Getenv("REDIS_PORT"))
		herokuAPIConfig := APIConfig{
			Port:     port,
			Version:  "1.0",
			Debug:    false,
			Endpoint: "kittentlapi.herokuapp.com",
			Jwt:      os.Getenv("JWT_TOKEN"),
			Secure:   false,
		}
		herokuDBConfig := DBConfig{
			Server:   os.Getenv("MONGO_HOST"),
			Database: os.Getenv("MONGO_DBNAME"),
		}
		herokuCacheConfig := CacheConfig{
			Server:   os.Getenv("REDIS_SERVER"),
			Port:     redisPort,
			User:     os.Getenv("REDIS_USER"),
			Password: os.Getenv("REDIS_PASSWORD"),
		}
		herokuUploadImage := UploadImageConfig{
			Path: "uploads/img/",
		}
		return Config{
			API:         herokuAPIConfig,
			DB:          herokuDBConfig,
			Cache:       herokuCacheConfig,
			UploadImage: herokuUploadImage,
		}
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
