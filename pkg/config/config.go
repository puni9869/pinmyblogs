package config

import (
	"errors"
	"github.com/spf13/viper"
)

var (
	environment  = "local"
	defaultPath  = "./config"
	environments = []string{"local", "prod"}
)

// AppConfig is the "appConfig"
type AppConfig struct {
	Debug       bool   `json:"debug"`
	SecretKey   string `json:"secretKey"`
	DefaultPort string `json:"defaultPort"`
	CustomPort  string `json:"customPort"`
	UseCDN      bool   `json:"useCDN"`
	Hostname    string `json:"hostname"`
}

// Database is the "database"
type Database struct {
	Type         string `json:"type"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Port         string `json:"port"`
	DatabaseName string `json:"databaseName"`
}

type Authentication struct {
	EnableLogin        bool `json:"enableLogin"`
	EnableRegistration bool `json:"enableRegistration"`
}

// Provider ie converted value in go types
type Provider struct {
	EnableSSL      bool           `json:"enableSSL"`
	Database       Database       `json:"database"`
	AppConfig      AppConfig      `json:"appConfig"`
	Authentication Authentication `json:"authentication"`
}

var C Provider

func LoadConfig(env string) error {
	if len(env) == 0 {
		return errors.New("environment is not provided. ie. local or prod")
	}
	var err error
	for _, e := range environments {
		if environment == e {
			break
		}
	}

	viper.SetConfigName(environment)
	// load from the config directory
	viper.AddConfigPath(GetDefaultPath())
	if err = viper.ReadInConfig(); err != nil {
		return err
	}
	err = viper.Unmarshal(&C)
	if err != nil {
		return err
	}
	return nil
}

func GetString(key string) string {
	return viper.GetString(key)
}

func GetStringMap(key string) map[string]any {
	return viper.GetStringMap(key)
}

func GetEnv() string {
	return environment
}

func GetDefaultPath() string {
	return defaultPath
}
