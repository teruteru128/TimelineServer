package config

import (
	"log"
	"path"
	"runtime"

	"github.com/BurntSushi/toml"
)

// Config 設定構造体
type Config struct {
	API APIConfig
	DB  DBConfig
}

// APIConfig API設定構造体
type APIConfig struct {
	Port     uint   `toml:"port"`
	Version  string `toml:"version"`
	Endpoint string `toml:"endpoint"`
	Debug    bool   `toml:"debug"`
	Jwt      string `toml:"jwt"`
}

// DBConfig MongoDB設定構造体
type DBConfig struct {
	Server   string `toml:"server"`
	Database string `toml:"database"`
	User     string `toml:"user"`
	Password string `toml:"password"`
}

// GetConfig TOML設定ファイルから設定を取得
func GetConfig() Config {
	var config Config
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	dir := path.Dir(filename)
	if _, err := toml.DecodeFile(dir+"/../config.toml", &config); err != nil {
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
