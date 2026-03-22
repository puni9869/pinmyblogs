// Package config loads and provides access to application configuration.
package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

var (
	defaultEnv  = "local"
	defaultPath = "./config"
	// LocalEnv is the local environment identifier.
	LocalEnv = "local"
	// ProdEnv is the production environment identifier.
	ProdEnv      = "prod"
	environments = []string{"local", "prod"}
)

// AppConfig is the "appConfig"
type AppConfig struct {
	Debug        bool   `json:"debug"`
	SecretKey    string `json:"secretKey"`
	DefaultPort  string `json:"defaultPort"`
	CustomPort   string `json:"customPort"`
	UseCDN       bool   `json:"useCDN"`
	Hostname     string `json:"hostname"`
	DataBaseType string `json:"dataBaseType"`
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

// Authentication holds feature flags for authentication flows.
type Authentication struct {
	EnableLogin                    bool `json:"enableLogin"`
	OpenDisabledAccountByEmailLink bool `json:"openDisabledAccountByEmailLink"`
	EnableRegistration             bool `json:"enableRegistration"`
	EnableForgotPassword           bool `json:"enableForgotPassword"`
}

// Mailer holds SMTP configuration for outbound email.
type Mailer struct {
	SmtpHost   string `json:"smtpHost"`
	SmtpPort   int    `json:"smtpPort"`
	EmailId    string `json:"emailId"`
	BccEmailId string `json:"bccEmailId"`
	Username   string `json:"username"`
	Password   string `json:"password"`
}

// ConfigProvider is the top-level configuration structure.
type ConfigProvider struct {
	EnableSSL         bool                   `json:"enableSSL"`
	ShareDataOverMail bool                   `json:"shareDataOverMail"`
	Database          map[string]DatabaseObj `json:"database"`
	AppConfig         AppConfig              `json:"appConfig"`
	Authentication    Authentication         `json:"authentication"`
	Mailer            Mailer                 `json:"mailer"`
}

// C is the global configuration instance populated by LoadConfig.
var C ConfigProvider

func contains(source []string, target string) bool {
	for _, s := range source {
		if s == target {
			return true
		}
	}
	return false
}

// LoadConfig reads and unmarshals the configuration file for the given environment.
func LoadConfig(environment string) error {
	if !contains(environments, environment) {
		return errors.New("environment is not provided. ie. local or prod")
	}

	var err error
	viper.SetConfigName(environment)

	// load from the config directory
	viper.AddConfigPath(GetDefaultPath())
	if err = viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error in setting env variables: %w", err)
	}
	err = viper.Unmarshal(&C)
	if err != nil {
		return fmt.Errorf("error in unmarshal variables: %w", err)
	}
	defaultEnv = environment
	return nil
}

// GetString returns a configuration value as a string.
func GetString(key string) string {
	return viper.GetString(key)
}

// GetStringMap returns a configuration value as a map of string to any.
func GetStringMap(key string) map[string]any {
	return viper.GetStringMap(key)
}

// GetEnv returns the current environment name.
func GetEnv() string {
	return defaultEnv
}

// SetKey sets a configuration key to the given value.
func SetKey(k string, v any) {
	viper.Set(k, v)
}

// GetDefaultPath returns the default configuration directory path.
func GetDefaultPath() string {
	return defaultPath
}
