package config

import (
	"errors"

	"github.com/spf13/viper"
)

var (
	defaultEnv   = "local"
	defaultPath  = "./config"
	LocalEnv     = "local"
	ProdEnv      = "prod"
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

// DatabaseObj is the "database"
type DatabaseObj struct {
	Type         string `json:"type,omitempty"`
	Host         string `json:"host,omitempty"`
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	Port         string `json:"port,omitempty"`
	DatabaseName string `json:"databaseName,omitempty"`
	LogSql       bool   `json:"logSql,omitempty"`
	FileName     string `json:"fileName,omitempty"`
}

type Authentication struct {
	EnableLogin          bool `json:"enableLogin"`
	EnableRegistration   bool `json:"enableRegistration"`
	EnableForgotPassword bool `json:"enableForgotPassword"`
}

type Mailer struct {
	SmtpHost   string `json:"smtpHost"`
	SmtpPort   int    `json:"smtpPort"`
	EmailId    string `json:"emailId"`
	BccEmailId string `json:"bccEmailId"`
	Username   string `json:"username"`
	Password   string `json:"password"`
}

// ConfigProvider ie converted value in go types
type ConfigProvider struct {
	EnableSSL         bool                   `json:"enableSSL"`
	ShareDataOverMail bool                   `json:"shareDataOverMail"`
	Database          map[string]DatabaseObj `json:"database"`
	AppConfig         AppConfig              `json:"appConfig"`
	Authentication    Authentication         `json:"authentication"`
	Mailer            Mailer                 `json:"mailer"`
}

var C ConfigProvider

func contains(source []string, target string) bool {
	for _, s := range source {
		if s == target {
			return true
		}
	}
	return false
}

func LoadConfig(environment string) error {
	if !contains(environments, environment) {
		return errors.New("environment is not provided. ie. local or prod")
	}

	var err error
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
	defaultEnv = environment
	return nil
}

func GetString(key string) string {
	return viper.GetString(key)
}

func GetStringMap(key string) map[string]any {
	return viper.GetStringMap(key)
}

func GetEnv() string {
	return defaultEnv
}

func SetKey(k string, v any) {
	viper.Set(k, v)
}

func GetDefaultPath() string {
	return defaultPath
}
