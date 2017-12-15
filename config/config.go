package config

import (
	"fmt"

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
	viper.AddConfigPath(".")
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
	conf.Server = host + ":" + port

	if user != "" {
		conf.Server = user + "@" + host + ":" + port
	}
	if user != "" && password != "" {
		conf.Server = user + ":" + password + "@" + host + ":" + port
	}
	conf.Database = viper.GetString("DB_NAME")

	return conf
}

func GetCacheConfig() CacheConfig {
	host := viper.GetString("REDIS_HOST")
	password := viper.GetString("REDIS_PASSWORD")
	port := viper.GetString("REDIS_PORT")

	conf := CacheConfig{}
	conf.Server = "redis://" + host + ":" + port
	if password != "" {
		conf.Server = "redis://" + host + ":" + port + "?password=" + password
	}
	return conf
}

func GetUploadImagePath() string {
	path := viper.GetString("IMAGE_UPLOAD_PATH")
	return path
}
