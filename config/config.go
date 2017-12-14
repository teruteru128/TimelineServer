package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
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
	Version  string
	Debug    bool
	Endpoint string
	Secure   bool
	Jwt      string
}

// DBConfig MongoDB設定構造体
type DBConfig struct {
	Server   string
	Database string
}

type CacheConfig struct {
	Server string
}

type UploadImageConfig struct {
	Path string
}

func init() {
	viper.SetDefault("Version", "v1")
	viper.SetDefault("Debug", false)
	viper.SetDefault("Endpoint", "localhost")
	viper.SetDefault("JWT_TOKEN", "DEFAULT_JWT_TOKEN_CHANGE_ME")
	viper.SetDefault("Secure", false)

	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_USERNAME", "root")
	viper.SetDefault("DB_PASSWORD", "")
	viper.SetDefault("DB_PORT", "3306")
	viper.SetDefault("DB_NAME", "timeline_dev")

	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PASSWORD", "")
	viper.SetDefault("REDIS_PORT", "6379")

	viper.SetDefault("IMAGE_UPLOAD_PATH", "/uploads/img")

	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("config file not found. using default configurations.")
	}
}

// GetAPIConfig API設定を取得
func GetAPIConfig() APIConfig {
	conf := APIConfig{
		Version:  viper.GetString("Version"),
		Debug:    viper.GetBool("Debug"),
		Endpoint: viper.GetString("Endpoint"),
		Secure:   viper.GetBool("Secure"),
		Jwt:      viper.GetString("JWT_TOKEN"),
	}
	return conf
}

// GetDBConfig DB接続情報を取得
func GetDBConfig() DBConfig {
	conf := DBConfig{}

	host := viper.GetString("DB_HOST")
	user := viper.GetString("DB_USERNAME")
	password := viper.GetString("DB_PASSWORD")
	port := viper.GetString("DB_PORT")
	conf.Server = user + ":" + password + "@" + host + ":" + port
	conf.Database = viper.GetString("DB_NAME")

	return conf
}

func GetCacheConfig() CacheConfig {
	host := viper.GetString("REDIS_HOST")
	password := viper.GetString("REDIS_PASSWORD")
	port := os.Getenv("REDIS_PORT")

	conf := CacheConfig{}
	conf.Server = "redis://" + host + "/" + port
	if password != "" {
		conf.Server = "redis://" + host + "/" + port + "?password=" + password
	}
	return conf
}

func GetUploadImagePath() string {
	path := viper.GetString("IMAGE_UPLOAD_PATH")
	return path
}
