package config

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	defaultPath = "../../config/"
	tests := []struct {
		name string
		env  string
		want error
	}{
		{
			name: "should not panic on local",
			env:  "local",
			want: nil,
		},
		{
			name: "should not panic on prod",
			env:  "prod",
			want: nil,
		},
		{
			name: "should panic with wrong path",
			env:  "dummy",
			want: errors.New("environment is not provided. ie. local or prod"),
		},
		{
			name: "should panic with empty path",
			env:  "",
			want: errors.New("environment is not provided. ie. local or prod"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := LoadConfig(tt.env)
			if err != nil && err.Error() != tt.want.Error() {
				t.Errorf("LoadConfig() = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name string
		env  string
		want string
	}{
		{
			name: "should not panic on local",
			env:  "local",
			want: "local",
		},
		{
			name: "should not panic on prod",
			env:  "prod",
			want: "prod",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defaultEnv = tt.env
			got := GetEnv()
			if got != tt.want {
				t.Errorf("LoadConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetString(t *testing.T) {
	tests := []struct {
		name string
		key  string
		val  string
		want string
	}{
		{
			name: "should string be dummy",
			key:  "ok",
			val:  "dummy",
			want: "dummy",
		},
		{
			name: "should string be na",
			key:  "ok1",
			val:  "na",
			want: "na",
		},
		{
			name: "should string be empty",
			key:  "ok1",
			val:  "",
			want: "",
		},
	}
	defaultPath = "../../config/"
	defaultEnv = "local"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetKey(tt.key, tt.val)
			got := GetString(tt.key)
			if got != tt.want {
				t.Errorf("GetString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetStringMap(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		val   map[string]any
		want  map[string]any
		count int
	}{
		{
			name:  "should []string be dummy",
			key:   "ok",
			val:   map[string]any{"dummy": "dummy"},
			want:  map[string]any{"dummy": "dummy"},
			count: 1,
		},
		{
			name:  "should []string be na",
			key:   "ok1",
			val:   map[string]any{"na": "na1", "na2": "na2"},
			want:  map[string]any{"na": "na1", "na2": "na2"},
			count: 2,
		},
		{
			name:  "should []string be empty",
			key:   "ok1",
			val:   map[string]any{},
			want:  map[string]any{},
			count: 0,
		},
	}
	defaultPath = "../../config/"
	defaultEnv = "local"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetKey(tt.key, tt.val)
			got := GetStringMap(tt.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetString() = %v, want %v", got, tt.want)
			}
			if len(got) != tt.count {
				t.Errorf("GetString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigProviderStruct(t *testing.T) {
	c := ConfigProvider{
		EnableSSL:      false,
		Database:       nil,
		AppConfig:      AppConfig{},
		Authentication: Authentication{},
		Mailer:         Mailer{},
	}
	want := "config.ConfigProvider{EnableSSL:false, ShareDataOverMail:false, Database:map[string]config.DatabaseObj(nil), AppConfig:config.AppConfig{Debug:false, SecretKey:\"\", DefaultPort:\"\", CustomPort:\"\", UseCDN:false, Hostname:\"\"}, Authentication:config.Authentication{EnableLogin:false, OpenDisabledAccountByEmailLink:false, EnableRegistration:false, EnableForgotPassword:false}, Mailer:config.Mailer{SmtpHost:\"\", SmtpPort:0, EmailId:\"\", BccEmailId:\"\", Username:\"\", Password:\"\"}}"
	got := fmt.Sprintf("%#v", c)
	fmt.Println(got)
	if strings.Compare(want, got) != 0 {
		t.Errorf("GetString() = %s, want %s", got, want)
	}
}
